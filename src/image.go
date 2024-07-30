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
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	imageName := filepath.Base(imageDir)
	if err := buildImage(imageName, imageDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created the image successfully: " + imageName + " " + imageDir})
}
func getImages(c *gin.Context) {
	cmd := exec.Command("docker", "image", "ls", "--format", "{{.Repository}}:{{.Tag}}")
	out, err := cmd.Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"images": string(out)})
}

func buildImage(imageName, imageDir string) error {
	cmd := exec.Command("docker", "build", "-q", "-t", imageName, imageDir)
	_, err := cmd.Output()
	return err
}
