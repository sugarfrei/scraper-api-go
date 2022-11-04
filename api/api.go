package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Api struct {
	*gin.Engine
	cfg        *conf.Cfg
	storage    *storage.Provider
	httpClient *http.Client
	auth       *jwt.GinJWTMiddleware
	err        error
}
