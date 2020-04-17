package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("html/*")

	router.GET("/", getForm)
	router.POST("/", postForm)

	router.Run(":8080")
}
