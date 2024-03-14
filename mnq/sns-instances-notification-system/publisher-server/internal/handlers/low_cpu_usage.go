package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"publisher-server/internal/utils"
)

func LowCPUUsageHandler(context *gin.Context) {
	if err := utils.PublishMessage("Notice: 'publisher-server' CPU usage normalized to 30%"); err != nil {

		log.Printf("Error publishing message: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish message"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Low CPU Usage simulated"})
}
