package main

import (
	"errors"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func getFolderByIndex(folderPath string, taskIndex int) (string, error) {
	directories, err := filepath.Glob(filepath.Join(folderPath, "*"))
	if err != nil {
		return "", err
	}

	if taskIndex < 0 || taskIndex >= len(directories) {
		return "", errors.New("task index out of bounds")
	}

	return directories[taskIndex], nil
}

func launchContainer(taskIndex int, imageName string) (string, error) {

	input, err := getFolderByIndex(InputFolder, taskIndex)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("docker", "run", "-d", "--rm", // Run in detached mode and remove container after it stops
		"-v", InputFolder+"/"+input+":/input",
		"-v", SharedFolder+":/shared",
		"-v", OutputFolder+":/output",
		imageName)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	containerID := string(output)
	containerID = containerID[:len(containerID)-1] // Remove newline character

	return containerID, nil

}

var containerMap = make(map[int]string)

func CreateTaskContainer(c *gin.Context) {

	taskIndexSTR, ok := c.GetQuery("taskIndex")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing taskIndex query parameter",
		})
		return
	}

	taskIndex, err := strconv.Atoi(taskIndexSTR)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid taskIndex query parameter",
		})
		return
	}

	containerID, err := launchContainer(taskIndex, jobInfo.ImageName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to launch container: " + err.Error(),
		})
		return
	}

	containerMap[taskIndex] = containerID
	c.JSON(http.StatusOK, gin.H{
		"taskIndex": taskIndex,
	})
}

func CheckTaskContainer(c *gin.Context) {

	taskIndexSTR, ok := c.GetQuery("taskIndex")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing taskIndex query parameter",
		})
		return
	}

	taskIndex, err := strconv.Atoi(taskIndexSTR)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid taskIndex query parameter",
		})
		return
	}

	containerID, ok := containerMap[taskIndex]
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"taskIndex": taskIndex,
			"status":    "not created",
		})
		return
	}

	cmd := exec.Command("docker", "inspect", "-f", "{{.State.Running}} {{.State.ExitCode}}", containerID)
	output, err := cmd.Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check container status: " + err.Error(),
		})
		return
	}

	substrings := strings.Split(string(output), " ")
	isRunning := substrings[0] == "true"
	exitCode, _ := strconv.Atoi(substrings[1])

	if isRunning {
		c.JSON(http.StatusOK, gin.H{
			"taskIndex": taskIndex,
			"status":    "running",
			"exitCode":  exitCode,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"taskIndex": taskIndex,
			"status":    "finished",
			"exitCode":  exitCode,
		})
	}

}
