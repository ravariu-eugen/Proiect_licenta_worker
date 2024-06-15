package main

import (
	"errors"
	"os/exec"
	"path/filepath"
)

const (
	InputFolder  = "/input"
	SharedFolder = "/shared"
	OutputFolder = "/output"
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

func launchContainer(taskIndex int) error {

	input, err := getFolderByIndex(InputFolder, taskIndex)
	if err != nil {
		return err
	}

	cmd := exec.Command("docker", "run",
		"-v", InputFolder+"/"+input+":/input",
		"-v", SharedFolder+":/shared",
		"-v", OutputFolder+":/output",
		"your_docker_image_name")

	err = cmd.Start() // Start the command in the background
	err = cmd.Wait()  // Wait for the command to finish
	return err

}
