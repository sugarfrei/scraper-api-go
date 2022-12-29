package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (a *Api) AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.Split(c.GetHeader("Authorization"), " ")
		if len(authHeader) != 2 || authHeader[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized,
				"message": "Bearer Authorization header is required",
			})
			return
		}

		token, err := a.ValidateToken(authHeader[1])
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized,
				"message": "Invalid token",
			})
			return
		}

	}
}

func (a *Api) ValidateToken(encodedToken string) (*jwt.Token, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(a.cfg.Core.PubKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
}
