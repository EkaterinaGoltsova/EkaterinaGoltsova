package main

import (
	"github.com/EkaterinaGoltsova/sprint-starter/web"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("web/template/*")

	router.GET("/", web.GetForm)
	router.POST("/", web.PostForm)

	router.Run(":8080")
}
