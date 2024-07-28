package main

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const (
	containerInputFolder  = "/run/input"
	containerOutputFolder = "/run/output"
	containerSharedFolder = "/run/shared"
)

func launchContainer(imageName, job, task string) (string, error) {

	// launches a container based on a task and an image

	taskInputDir := InputFolder + "/" + job + "/" + task
	taskOutputDir := OutputFolder + "/" + job + "/" + task

	cmd := exec.Command("docker", "run", "-d", // Run in detached mode and remove container after it stops
		"-v", taskInputDir+":"+containerInputFolder,
		"-v", SharedFolder+":"+containerSharedFolder,
		"-v", taskOutputDir+":"+containerOutputFolder,
		imageName)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	containerID := string(output)
	containerID = containerID[:len(containerID)-1] // Remove newline character

	return containerID, nil

}

func CreateTaskContainer(c *gin.Context) {

	jobName := c.Request.FormValue("job")
	image := c.Request.FormValue("image")

	uploadDir := UploadFolder + "/" + jobName
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	taskDir, err := uploadAndExtractToDir(c, uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	taskName := filepath.Base(taskDir)

	containerID, err := launchContainer(image, jobName, taskName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"containerID": containerID,
		"taskName":    taskName,
		"status":      "created",
	})

}
