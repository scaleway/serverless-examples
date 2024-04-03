package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"subscriber-server/internal/utils"
)

func SnsServiceHandler(context *gin.Context) {
	messageType := context.GetHeader("x-amz-sns-message-type")

	switch messageType {
	case "SubscriptionConfirmation":
		utils.GetSubscribeURL(context)
	case "Notification":
		utils.FormatNotification(context)
	default:
		context.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported message type"})
	}
}
