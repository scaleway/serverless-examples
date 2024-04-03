package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"publisher-server/internal/utils"
)

func HighCPUUsageHandler(context *gin.Context) {
	if err := utils.PublishMessage("Alert: 'publisher-server' CPU usage is above 70%"); err != nil {
		log.Printf("Error publishing message: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish message"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "High CPU Usage simulated"})
}
