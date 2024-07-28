package main

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func addShared(c *gin.Context) {

	newDir, err := uploadAndExtractToDir(c, SharedFolder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fileName := filepath.Base(newDir)

	c.JSON(http.StatusOK, gin.H{"message": "Uploaded the file successfully: " + fileName})
}

func getShared(c *gin.Context) {
	getFileList(c, SharedFolder)
}
