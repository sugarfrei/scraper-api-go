package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"scraper-api-go/conf"
	"scraper-api-go/model"
	"syscall"
	"time"

	"k8s.io/klog"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	ErrMsgDefaultClientError = "Bad user input"
	ErrMsgValidation         = "Validation failed"
	ErrMsgBadRequest         = "Bad request"
	ErrMsgUnauthorized       = "Unathorized"
)

type Api struct {
	*gin.Engine
	httpClient *http.Client
	validator  *validator.Validate
	cfg        *conf.Cfg
	err        error
}

func New(cfg *conf.Cfg) (*Api, error) {
	a := &Api{}
	a.initConf(cfg)
	a.initHttpClient()
	a.initRouter()
	a.initValidator()

	return a, a.err
}

/*********** VALIDATOR HELPERS ***********/
func (a *Api) Validate(v interface{}) error {
	return a.validator.Struct(v)
}

/*********** ERROR HANDLING ***********/
func (a *Api) error(err *model.HttpError) {
	loggable := a.isErrorLoggable(err)

	if loggable {
		klog.Errorf("Stack: %v, %s", err.Caller(), err.WithTraceID())
	}
}

func (a *Api) abortWithError(c *gin.Context, err *model.HttpError) {
	a.error(err)
	c.Abort()
	c.JSON(err.Code, err)
}

func (a *Api) isErrorLoggable(err *model.HttpError) (ok bool) {
	switch a.cfg.Error.Log {
	case 2:
		ok = true
	case 1:
		if err.IsServer() {
			ok = true
		}
	}
	return
}

/*********** HTTP SERVER ***********/
func (a *Api) ListenAndServe() {
	executable, err := os.Executable()
	if err != nil {
		klog.Errorf("Unable to get execuatable, err=%v", err.Error())
		return
	}
	hostname, err := os.Hostname()
	if err != nil {
		klog.Errorf("Unable to get hostname, err=%v", err.Error())
		return
	}

	httpServer := &http.Server{
		Handler:           a.Engine,
		Addr:              a.cfg.Http.Listen,
		ReadTimeout:       time.Duration(a.cfg.Http.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(a.cfg.Http.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(a.cfg.Http.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(a.cfg.Http.IdleTimeout) * time.Second,
		MaxHeaderBytes:    a.cfg.Http.MaxHeaderBytes,
	}

	rand.Seed(time.Now().UnixNano())

	httpdone := make(chan bool)

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			klog.Errorf("Stopping listen and serve, err=%v.", err)
		}
		httpdone <- true
	}()

	startingMsg := fmt.Sprintf("Starting %v on %v!", executable, hostname)
	klog.V(2).Infoln(startingMsg)

	signals := make(chan os.Signal, 10)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
L:
	for {
		select {
		case sigval := <-signals:
			klog.Infof("Incoming signal is caught, num=%d.\n", sigval)

			switch sigval {
			case syscall.SIGINT, syscall.SIGTERM:
				ctx, cancelServer := context.WithTimeout(context.Background(), 5*time.Second)
				if err := httpServer.Shutdown(ctx); err != nil {
					klog.Errorf("Server Shutdown: %v", err.Error())
				}

				signal.Stop(signals)

				<-httpdone
				cancelServer()
				break L
			case syscall.SIGHUP:

				break
			}
		case <-httpdone:
			break L
		}
	}

	stoppingMsg := fmt.Sprintf("Stopping %v on %v!", executable, hostname)
	klog.V(2).Infoln(stoppingMsg)
	klog.Flush()
}
