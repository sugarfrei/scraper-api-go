package api

import (
	"net/http"
	"scraper-api-go/model"

	"github.com/gin-gonic/gin"
)

// route not found handler
func (a *Api) notFound(c *gin.Context) {
	a.abortWithError(c, model.NewHttpError(http.StatusNotFound, "The server has not found anything matching the Request-URI"))
}

// method not allowed handler
func (a *Api) methodNotAllowed(c *gin.Context) {
	a.abortWithError(c, model.NewHttpError(http.StatusMethodNotAllowed, "Method not allowed"))
}

func (a *Api) registerRoutes() {
	a.NoRoute(a.notFound)

	a.HandleMethodNotAllowed = true
	a.NoMethod(a.methodNotAllowed)

	v1 := a.Group("/v1")
	{
		/*********** PRIVATE ENDPOINTS ***********/
		private := v1.Group("/")
		private.Use(a.AuthorizeJWT())
		private.GET("instagram/:username", a.Instagram)
		private.GET("twitter/:username", a.Twitter)

		/*********** PUBLIC ENDPOINTS ***********/
		public := v1.Group("/")
		public.POST("login", a.Login)

	}
}
