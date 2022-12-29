package api

import (
	"net"
	"net/http"
	"os"
	"scraper-api-go/conf"
	"scraper-api-go/model"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"k8s.io/klog"
)

func (a *Api) initConf(cfg *conf.Cfg) {
	if a.err != nil {
		return
	}

	a.cfg = cfg
}

/*********** INIT VALIDATOR ***********/
func (a *Api) initValidator() {
	if a.err != nil {
		return
	}

	a.validator = validator.New()
}

/*********** INIT HTTP CLIENT ***********/
func (a *Api) initHttpClient() {
	if a.err != nil {
		return
	}

	transport := &http.Transport{
		DisableKeepAlives:   false,
		DisableCompression:  false,
		MaxIdleConnsPerHost: 1000,
		MaxConnsPerHost:     1000,
		Dial: (&net.Dialer{
			Timeout:   time.Duration(1) * time.Second,
			KeepAlive: 0 * time.Second,
		}).Dial,
		IdleConnTimeout: time.Duration(900) * time.Second,
	}

	a.httpClient = &http.Client{
		Transport: transport,
		Timeout:   time.Duration(5) * time.Second,
	}
}

/*********** INIT ROUTER ***********/
func (a *Api) initRouter() {
	if a.err != nil {
		return
	}

	file, err := os.OpenFile(a.cfg.Core.AccessLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		klog.Errorf("Unable to open log access log file: %s", err.Error())
	} else {
		gin.DefaultWriter = file
	}

	a.Engine = gin.New()
	a.Engine.Use(gin.Logger())
	a.Engine.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		a.abortWithError(c, model.ServerError(err).Prefix("Stop panicking: "))
	}))

	a.Engine.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "PUT", "PATCH", "DELETE", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	a.registerRoutes()
}
