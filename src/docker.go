package main

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

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
