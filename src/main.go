package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	UploadFolder = "/app/upload"
	ImageFolder  = "/app/images"
	InputFolder  = "/app/input"
	SharedFolder = "/app/shared"
	OutputFolder = "/app/output"
)

func runContainerServer() error {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello, World!!!!")
	})

	router.GET("/images", getImages)
	router.POST("/images", addImage)

	router.GET("/shared", getShared)
	router.POST("/shared", addShared)

	//router.POST("/tasks", CreateTaskContainer)

	return router.Run(":8080")
}

func main() {
	if err := runContainerServer(); err != nil {
		log.Fatal(err)
	}
}
