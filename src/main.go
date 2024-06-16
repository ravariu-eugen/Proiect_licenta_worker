package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	InputFolder  = "/input"
	SharedFolder = "/shared"
	OutputFolder = "/output"
)

func runContainerServer() error {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello, World!")
	})
	router.GET("/system_info", getSystemInfo)

	return router.Run(":8080")
}

func main() {
	if err := runContainerServer(); err != nil {
		log.Fatal(err)
	}
}
