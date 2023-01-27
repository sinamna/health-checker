package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"health_checker/pkg/repository"
	"net/http"
	"strconv"
)

func AlertHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	alerts, dbErr := repository.Database.GetAlertByEndpointID(id)
	if dbErr != nil {
		fmt.Println(dbErr.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "could not get endpoint's status",
		})
		return
	}

	c.JSON(http.StatusOK, alerts)
}
