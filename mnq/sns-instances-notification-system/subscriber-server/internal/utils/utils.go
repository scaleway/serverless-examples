package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"subscriber-server/internal/app"
	"subscriber-server/internal/types"
)

type ConfirmationRequest struct {
	SubscribeURL string `json:"SubscribeURL"`
}

func GetSubscribeURL(context *gin.Context) {
	var request ConfirmationRequest

	err := context.BindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app.SharedState.ConfirmationMutex.Lock()
	defer app.SharedState.ConfirmationMutex.Unlock()

	app.SharedState.ReceivedConfirmation.SubscribeURL = request.SubscribeURL
	app.SharedState.ReceivedConfirmation.Confirmed = false

	context.JSON(http.StatusOK, gin.H{"message": "URL received for confirmation"})
}

func FormatNotification(context *gin.Context) {
	var snsNotification types.Notification

	if err := context.BindJSON(&snsNotification); err != nil {
		context.String(http.StatusBadRequest, "Error parsing notification")
		return
	}

	app.SharedState.NotificationMutex.Lock()
	defer app.SharedState.NotificationMutex.Unlock()

	snsMessage := snsNotification.Message
	if snsNotification.Subject != "" {
		snsMessage = snsNotification.Subject + ": " + snsMessage
	}
	app.SharedState.FormattedNotifications = append(app.SharedState.FormattedNotifications, snsMessage)

	context.JSON(http.StatusOK, gin.H{"message": "Notification added"})
}
