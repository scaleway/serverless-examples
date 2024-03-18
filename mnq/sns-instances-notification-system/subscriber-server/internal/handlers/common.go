package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NoRouteHandler(context *gin.Context) {
	context.HTML(http.StatusNotFound, "message.html", gin.H{
		"Title":   "404 Not Found",
		"Message": "Oh no 😢 The tutorial lost you",
	})
}

func GetIndexHandler(context *gin.Context) {
	context.File("./www/index.html")
}