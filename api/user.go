package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"health_checker/pkg/jwt"
	"health_checker/pkg/model"
	"health_checker/pkg/repository"
	"health_checker/pkg/utils"
	"net/http"
)

func SignupHandler(c *gin.Context) {
	var userReq model.User
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Printf("user %s wants to signup \n", userReq.Username)

	dbErr := repository.Database.CreateNewUser(userReq.Username, userReq.Password)
	if dbErr != nil {
		fmt.Println(dbErr.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s successfully signed up. please login to retrieve your token", userReq.Username),
	})
}
func LoginHandler(c *gin.Context) {
	var loginReq model.User
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Printf("user %s wants to login \n", loginReq.Username)

	user, dbErr := repository.Database.GetUserByID(loginReq.Username)
	if dbErr != nil {
		fmt.Println(dbErr.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to login",
		})
		return
	}

	// comparing passwords
	hashedPass := utils.HashString(loginReq.Password)
	if user.Password != hashedPass {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "password mismatched",
		})
		return
	}

	jwtToken, err := jwt.GenerateToken(user.Username)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "could not create jwt token",
		})
		return
	}

	err = repository.Database.FlushUserTokens(user.Username)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "could not flush previous token",
		})
		return
	}

	err = repository.Database.CreateNewToken(user.Username, jwtToken)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "could not insert token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s successfully logged in", user.Username),
		"token":   jwtToken,
	})
}
