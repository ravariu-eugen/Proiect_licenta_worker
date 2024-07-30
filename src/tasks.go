package main

import (
	"fmt"
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

func CreateTaskContainer(c *gin.Context) {
	jobName := c.PostForm("job")
	image := c.PostForm("image")

	jobDir := filepath.Join(InputFolder, jobName)
	if err := os.MkdirAll(jobDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	taskDir, err := uploadAndExtractToDir(c, jobDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}
	taskName := filepath.Base(taskDir)

	containerID, err := launchContainer(image, jobName, taskName)
	if err != nil {
		//getFileList(c, taskDir)
		c.JSON(http.StatusInternalServerError, gin.H{"error3": err.Error(), "dir": taskDir})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"containerID": containerID,
		"taskName":    taskName,
		"status":      "created",
	})
}

type IDMap struct {
	idmap map[string]string
}

var globalIDMap IDMap = IDMap{make(map[string]string)}

func (idMap IDMap) getContainerID(jobName, taskName string) (containerID string) {

	key := jobName + "-" + taskName

	if val, ok := idMap.idmap[key]; ok {
		return val
	}

	return ""
}

func (idMap IDMap) putContainerID(jobName, taskName, containerID string) {
	key := jobName + "-" + taskName

	idMap.idmap[key] = containerID

}

func launchContainer(imageName, job, task string) (string, error) {

	// launches a container based on a task and an image

	taskInputDir := InputFolder + "/" + job + "/" + task
	taskOutputDir := OutputFolder + "/" + job + "/" + task

	if err := os.MkdirAll(taskOutputDir, 0755); err != nil {
		return "", err
	}

	cmd := exec.Command("docker", "run", "-d", // Run in detached mode and remove container after it stops
		"-v", taskInputDir+":"+containerInputFolder,
		"-v", SharedFolder+":"+containerSharedFolder,
		"-v", taskOutputDir+":"+containerOutputFolder,
		imageName)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to launch container: %v", err)
	}

	containerID := string(output)
	containerID = containerID[:len(containerID)-1] // Remove newline character

	globalIDMap.putContainerID(job, task, containerID)

	return containerID, nil

}

func getTask(c *gin.Context) {
	jobName := c.Param("job")
	taskName := c.Param("task")
	containerID := globalIDMap.getContainerID(jobName, taskName)
	status := getContainerStatus(containerID)
	switch status {
	case "running":
		c.JSON(http.StatusOK, gin.H{
			"job":         jobName,
			"task":        taskName,
			"containerID": containerID,
			"status":      status,
		})
	case "exited":

		returnResult(c)
	default:
		c.JSON(http.StatusOK, gin.H{
			"job":         jobName,
			"task":        taskName,
			"containerID": containerID,
			"status":      status,
		})
	}

}

func returnResult(c *gin.Context) {
	jobName := c.Param("job")
	taskName := c.Param("task")

	// create output archive
	taskOutputDir := OutputFolder + "/" + jobName + "/" + taskName
	archiveName := filepath.Base(taskOutputDir) + ".zip"
	archivePath := filepath.Join(OutputFolder, jobName, archiveName)
	cmd := exec.Command("zip", "-q", "-r", archiveName, taskOutputDir)

	_, err := cmd.Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error: zip failed": err.Error()})
		return
	}

	// c.Header("Content-Description", "File Transfer")
	// c.Header("Content-Transfer-Encoding", "binary")
	// c.Header("Content-Disposition", "attachment; filename="+archiveName)
	// c.Header("Content-Type", "application/octet-stream")
	// c.File(archivePath)
	c.JSON(http.StatusOK, gin.H{"message": "Downloaded the file successfully: " + archiveName, "archivePath": archivePath})
}

func getContainerStatus(containerID string) string {

	cmd := exec.Command("docker", "container", "inspect", "-f", "{{.State.Status}}", containerID)

	out, err := cmd.Output()

	if err != nil {
		return ""
	}
	status := string(out)

	return status[:len(status)-1]
}
