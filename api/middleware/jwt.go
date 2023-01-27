package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	jwt2 "health_checker/pkg/jwt"
	"health_checker/pkg/repository"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimSpace(strings.SplitN(authHeader, "Bearer", 2)[1])
		claims, err := jwt2.ParseToken(tokenString)
		if err == nil {
			token, err2 := repository.Database.GetTokenByUsername(claims.Username)
			if err2 != nil {
				fmt.Println(err2.Error())
				return
			}
			if token != tokenString {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
			}
			c.Set("username", claims.Username)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		}
	}
}
