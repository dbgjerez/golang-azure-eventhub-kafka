package controllers

import (
	"golang-azure-eventhub-kafka/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
}

func HealthControllerHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		/* TODO: modify to return UP or DOWN depending of dependencies */
		h := models.HealthDTO{}
		h.Status = models.HealhStatusUp
		c.JSON(http.StatusOK, h)
	}
}
