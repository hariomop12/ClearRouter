package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hariomop12/clearrouter/apps/backend/internal/services"
)

func GetModelsHandler(c *gin.Context) {
	models := services.GetAllModels()
	c.JSON(http.StatusOK, gin.H{
		"data": models,
	})
}
