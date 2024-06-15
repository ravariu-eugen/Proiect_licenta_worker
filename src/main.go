package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

const (
	InputFolder  = "/input"
	SharedFolder = "/shared"
	OutputFolder = "/output"
)

func runContainerServer() error {
	router := gin.Default()
	router.GET("/system_info", getSystemInfo)

	return router.Run(":8080")
}

func main() {
	if err := runContainerServer(); err != nil {
		log.Fatal(err)
	}
}
