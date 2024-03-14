package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"subscriber-server/internal/app"
)

func GetNotificationsHandler(context *gin.Context) {
	app.SharedState.NotificationMutex.Lock()
	defer app.SharedState.NotificationMutex.Unlock()
	context.HTML(http.StatusOK, "notifications.html", gin.H{
		"FormattedNotifications": app.SharedState.FormattedNotifications,
	})
}

func RefreshNotificationsHandler(context *gin.Context) {
	app.SharedState.NotificationMutex.Lock()
	defer app.SharedState.NotificationMutex.Unlock()
	context.JSON(http.StatusOK, gin.H{
		"notifications": app.SharedState.FormattedNotifications,
	})
}
