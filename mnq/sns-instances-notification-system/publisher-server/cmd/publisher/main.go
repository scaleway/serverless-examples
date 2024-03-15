package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"publisher-server/internal/handlers"
	"publisher-server/internal/sns_client"
)

func main() {
	sns_client.InitAwsSnsClient()

	router := gin.Default()
	router.LoadHTMLGlob("./www/templates/*")
	router.Static("/", "./www")
	router.POST("/low-CPU-usage", handlers.LowCPUUsageHandler)
	router.POST("/high-CPU-usage", handlers.HighCPUUsageHandler)

	router.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusNotFound, "message.html", gin.H{
			"Title":   "404 Not Found",
			"Message": "Oh no 😢 The tutorial lost you",
		})
	})

	log.Println("Publisher-server listening on port 8081 ...")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
