package model

import (
	"fmt"
	"net/http"
	"runtime"

	uuid "github.com/satori/go.uuid"
)

// curentlly only 1 frame is supported
var ErrorMaxStackFrames = 50

var errorsMap = map[int]string{
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	409: "Conflict",
	422: "Unprocessable Entity",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeuot",
}

type HttpError struct {
	err     error
	Code    int    `json:"code" xml:"code" yaml:"code"`
	Title   string `json:"title" xml:"title" yaml:"title"`
	Message string `json:"message,omitempty" xml:"message,omitempty" yaml:"message,omitempty"`
	TraceID string `json:"trace_id,omitempty" xml:"trace_id" yaml:"trace_id"`
	Details struct {
		Field string      `json:"field,omitempty" xml:"field,omitempty" yaml:"field,omitempty"`
		Value interface{} `json:"value,omitempty" xml:"value,omitempty" yaml:"value,omitempty"`
		Rule  string      `json:"rule,omitempty" xml:"rule,omitempty" yaml:"rule,omitempty"`
		Issue string      `json:"issue,omitempty" xml:"issue,omitempty" yaml:"issue,omitempty"`
	} `json:"details,omitempty" xml:"details,omitempty" yaml:"details,omitempty"`
	developerMsg string
	caller       *StackFrame
	prefix       string
}

func err(skip, code int, e interface{}) *HttpError {
	//todo return nil if error code is not 4xx or 5xx
	var err error

	switch e := e.(type) {
	case error:
		err = e
	default:
		err = fmt.Errorf("%v", e)
	}

	httpError := &HttpError{
		err:          err,
		Code:         code,
		Title:        errorsMap[code],
		TraceID:      uuid.NewV4().String(),
		developerMsg: err.Error(),
	}

	if httpError.IsClient() {
		if httpError.Title == "" {
			httpError.Title = errorsMap[400]
		}
		httpError.Message = err.Error()
	}

	if httpError.Title == "" {
		httpError.Title = errorsMap[500]
	}

	//todo go MaxStackFrames deep in the stack
	//currently only caller -> first frame in the stack is supported
	if ErrorMaxStackFrames > 0 {
		pc, file, line, ok := runtime.Caller(skip)
		if ok {
			httpError.caller = &StackFrame{
				File: file,
				Line: line,
				Func: runtime.FuncForPC(pc).Name(),
				Pc:   pc,
			}
		}
	}

	return httpError
}

func NewHttpError(code int, e interface{}) *HttpError {
	return err(2, code, e)
}

func ServerError(e interface{}) *HttpError {
	return err(2, http.StatusInternalServerError, e)
}

func ClientError(e interface{}) *HttpError {
	return err(2, http.StatusBadRequest, e)
}

func (e *HttpError) Error() string {
	errString := e.developerMsg

	if e.prefix != "" {
		errString = e.prefix + " " + errString
	}

	return errString
}

func (e *HttpError) WithTraceID() string {
	return fmt.Sprintf("Trace: %s, error: %s", e.TraceID, e.Error())
}

func (e *HttpError) Caller() *StackFrame {
	return e.caller
}

func (e *HttpError) Prefix(prefix string) *HttpError {
	e.prefix = prefix
	return e
}

func (e *HttpError) Prefixf(format string, a ...interface{}) *HttpError {
	e.prefix = fmt.Sprintf(format, a...)
	return e
}

func (e *HttpError) IsClient() bool {
	return e.Code >= 400 && e.Code < 500
}

func (e *HttpError) IsServer() bool {
	return !e.IsClient()
}

type StackFrame struct {
	File string
	Line int
	Func string
	Pc   uintptr
}
