package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegistryImage struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Registry  string `json:"registry"`
	ImageName string `json:"image"`
}

func GetRegistryImage(c *gin.Context) {
	var data RegistryImage
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

type ArchiveImage struct {
	Bucket    string `json:"bucket"`
	ImageFile string `json:"imageFile"`
}

func GetArchiveImage(c *gin.Context) {
	var data ArchiveImage
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
