package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"subscriber-server/internal/handlers"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("../../www/templates/*")

	router.GET("/", handlers.GetIndexHandler)
	router.GET("/confirm-subscription", handlers.ConfirmSubscriptionHandler)
	router.POST("/notifications", handlers.SnsServiceHandler)
	router.POST("/confirm-subscription", handlers.ConfirmClickHandler)
	router.GET("/notifications", handlers.GetNotificationsHandler)
	router.GET("/notifications-refresh", handlers.RefreshNotificationsHandler)

	router.NoRoute(handlers.NoRouteHandler)

	fmt.Println("Subscriber-server listening on port 8081 ...")
	if err := router.Run(":8081"); err != nil {
		panic("Error when listening on port 8081: " + err.Error())
	}
}
