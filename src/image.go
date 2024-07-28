package main

import (
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func addImage(c *gin.Context) {

	imageDir, err := uploadAndExtractToDir(c, ImageFolder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() + "||| " + imageDir})
		return
	}
	imageName := filepath.Base(imageDir)
	// build the image
	err = buildImage(imageName, imageDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created the image successfully: " + imageName})
}

func getImages(c *gin.Context) {
	getFileList(c, ImageFolder)
}

func buildImage(imageName, imageDir string) error {
	cmd := exec.Command("docker", "build", "-q", "-t", imageName, imageDir)
	_, err := cmd.Output()
	return err
}
