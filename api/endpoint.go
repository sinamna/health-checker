package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"health_checker/pkg/model"
	"health_checker/pkg/repository"
	"net/http"
	"strconv"
)

func EndpointRetrieverHandler(c *gin.Context) {
	value, exists := c.Get("username")
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "could not authenticate user",
		})
		return
	}
	username := value.(string)
	endpoints, dbErr := repository.Database.GetEndpointsByUsername(username)
	if dbErr != nil {
		fmt.Println(dbErr.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "could not get user's endpoints",
		})
		return
	}

	c.JSON(http.StatusOK, endpoints)
}

func EndpointRegisterHandler(c *gin.Context) {
	value, exists := c.Get("username")
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "could not authenticate user",
		})
		return
	}
	username := value.(string)
	var endpoint model.Endpoint
	if err := c.ShouldBindJSON(&endpoint); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	id, dbErr := repository.Database.CreateNewEndpoint(endpoint.Url, username, endpoint.Threshold)
	if dbErr != nil {
		fmt.Println(dbErr.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "could not create endpoint",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("successfully registered endpoint %s for user %s", endpoint.Url, username),
		"id":      id,
	})
}

func EndpointStatusHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	success, failed, dbErr := repository.Database.GetEndpointStatusByID(id)
	if dbErr != nil {
		fmt.Println(dbErr.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "could not get endpoint's status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"success": success,
		"failed":  failed,
	})
}
