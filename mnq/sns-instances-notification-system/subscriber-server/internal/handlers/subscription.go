package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"subscriber-server/internal/app"
)

func ConfirmSubscriptionHandler(context *gin.Context) {
	app.SharedState.ConfirmationMutex.Lock()
	defer app.SharedState.ConfirmationMutex.Unlock()

	if app.SharedState.ReceivedConfirmation.SubscribeURL == "" {
		context.HTML(http.StatusNotFound, "message.html", gin.H{
			"Title":   "Error",
			"Message": "Subscription URL not received",
		})
		return
	}

	if app.SharedState.ReceivedConfirmation.Confirmed {
		context.HTML(http.StatusGone, "message.html", gin.H{
			"Title":   "Error",
			"Message": "Subscription already confirmed",
		})
		return
	}

	context.HTML(http.StatusOK, "confirm-subscription.html", gin.H{"URL": app.SharedState.ReceivedConfirmation.SubscribeURL})
}

func ConfirmClickHandler(context *gin.Context) {
	app.SharedState.ConfirmationMutex.Lock()
	defer app.SharedState.ConfirmationMutex.Unlock()

	if app.SharedState.ReceivedConfirmation.Confirmed {
		context.JSON(http.StatusAlreadyReported, gin.H{"message": "Already confirmed"})
		return
	}

	app.SharedState.ReceivedConfirmation.Confirmed = true
	context.JSON(http.StatusOK, gin.H{"message": "Subscription confirmed successfully"})
}
