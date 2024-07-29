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
	jobName := c.PostForm("job")
	image := c.PostForm("image")

	uploadDir := filepath.Join(UploadFolder, jobName)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	taskDir, err := uploadAndExtractToDir(c, uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}
	taskName := filepath.Base(taskDir)

	containerID, err := launchContainer(image, jobName, taskName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error3": err.Error(), "dir": taskDir})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"containerID": containerID,
		"taskName":    taskName,
		"status":      "created",
	})
}

func getTask(c *gin.Context) {
	jobName := c.Param("job")
	taskName := c.Param("task")
	c.JSON(http.StatusOK, gin.H{"status": "ok", "job": jobName, "task": taskName})
}
