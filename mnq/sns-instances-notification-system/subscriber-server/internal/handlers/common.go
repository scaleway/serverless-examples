package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NoRouteHandler(context *gin.Context) {
	context.HTML(http.StatusNotFound, "message.html", gin.H{
		"Title":   "404 Not Found",
		"Message": "Oh no 😢 The tutorial lost you",
	})
}

func GetIndexHandler(context *gin.Context) {
	context.File("../../www/index.html")
}
