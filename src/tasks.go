package main

import (
	"encoding/json"
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

type FileNameMapping struct {
	File string `json:"file"`
	Name string `json:"name"`
}

func CreateTaskContainer(c *gin.Context) {
	jobName := c.PostForm("job")
	image := c.PostForm("image")
	mappingJson := c.PostForm("mappings")

	var mappings []FileNameMapping
	err := json.Unmarshal([]byte(mappingJson), &mappings)

	jobDir := filepath.Join(InputFolder, jobName)
	if err := os.MkdirAll(jobDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	taskDir, err := uploadAndExtractToDir(c, jobDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error(), "mappings": mappingJson})
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

// getContainerID retrieves the container ID associated with a specific job and task.
//
// jobName and taskName are used to construct the key for the ID map.
// Returns the container ID as a string if found, otherwise an empty string.
func (idMap IDMap) getContainerID(jobName, taskName string) (containerID string) {

	key := jobName + "-" + taskName

	if val, ok := idMap.idmap[key]; ok {
		return val
	}

	return ""
}

// putContainerID maps a container ID to a job and task name.
//
// jobName and taskName are used to create a key for the mapping, and containerID is the value being mapped.
// No return value.
func (idMap IDMap) putContainerID(jobName, taskName, containerID string) {
	key := jobName + "-" + taskName

	idMap.idmap[key] = containerID

}

func launchContainer(imageName, job, task string, mappings ...FileNameMapping) (string, error) {

	// launches a container based on a task and an image

	taskInputDir := InputFolder + "/" + job + "/" + task
	taskOutputDir := OutputFolder + "/" + job + "/" + task

	if err := os.MkdirAll(taskOutputDir, 0755); err != nil {
		return "", err
	}

	argList := []string{"run", "-d", // Run in detached mode and remove container after it stops
		"-v", taskInputDir + ":" + containerInputFolder,
		"-v", taskOutputDir + ":" + containerOutputFolder,
	}

	for _, mapping := range mappings {
		localFileDir := SharedFolder + "/" + getFileNameWithoutExt(mapping.File)
		containerFileDir := containerSharedFolder + "/" + mapping.Name
		argList = append(argList, "-v", localFileDir+":"+containerFileDir)
	}

	argList = append(argList, imageName)

	cmd := exec.Command("docker", argList...)

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
	if containerID == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"job":   jobName,
			"task":  taskName,
			"error": "container not found",
		})
		return
	}
	status := getContainerStatus(containerID)
	switch status {
	case "running":
		c.JSON(http.StatusNoContent, gin.H{
			"job":         jobName,
			"task":        taskName,
			"containerID": containerID,
			"status":      status,
		})
	case "exited":

		returnResult(c)
	default:
		c.JSON(http.StatusNoContent, gin.H{
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
	jobOutputDir := filepath.Join(OutputFolder, jobName)
	err := os.Chdir(jobOutputDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error: chdir": err.Error()})
		return
	}
	taskOutputDir := taskName
	archiveName := filepath.Base(taskOutputDir) + ".zip"
	archivePath := filepath.Join(OutputFolder, jobName, archiveName)
	cmd := exec.Command("zip", "-q", "-r", archiveName, taskName)
	fmt.Print(archivePath)
	_, err = cmd.Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error: zip failed": err.Error()})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+archiveName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(archivePath)
	//getFileList(c, OutputFolder+"/"+jobName)
}

func clearTask(jobName, taskName string) {

}

func removeImage(imageName string) error {
	cmd := exec.Command("docker", "image", "rm", imageName)

	_, err := cmd.Output()
	return err
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
