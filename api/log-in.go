package api

import (
	"net/http"
	"scraper-api-go/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (a *Api) Login(c *gin.Context) {

	var obj model.Authorization
	if err := c.ShouldBind(&obj); err != nil {
		a.abortWithError(c, model.ClientError(ErrMsgDefaultClientError).Prefix(err.Error()))
		return
	}

	if err := a.Validate(obj); err != nil {
		a.abortWithError(c, model.ClientError(err).Prefix(ErrMsgValidation))
		return
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(a.cfg.Core.PrvKey)
	if err != nil {
		a.abortWithError(c, model.ServerError(err).Prefix("Problem while parsing private key"))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"exp":        time.Now().Add(24 * time.Hour),
		"authorized": true,
		"user":       obj.Username,
		"password":   obj.Password,
	})

	tokenString, err := token.SignedString(key)
	if err != nil {
		a.abortWithError(c, model.ServerError(err).Prefix("Problem while generating the token"))
		return
	}

	c.JSON(http.StatusCreated, &model.Authorization{Status: "Success", Token: tokenString})
}
