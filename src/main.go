package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
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
		"-v", inputFolder+":/input",
		"-v", SharedFolder+":/shared",
		"-v", OutputFolder+":/output",
		"your_docker_image_name")

	cmd := exec.Command("docker", "run",
		"-v", "/input:/input",
		"-v", "/shared:/shared",
		"-v", "/output:/output",
		"your_docker_image_name")
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read a number from the query string
		taskString := r.URL.Query().Get("task")
		task, err := strconv.Atoi(taskString)
		if err != nil {
			http.Error(w, "Invalid number", http.StatusBadRequest)
			return
		}
		// Launch a Docker container when an HTTP request is received
		cmd := exec.Command("docker", "run",
			"-v", "/input:/app/input",
			"-v", "/shared:/app/shared",
			"-v", "/output:/app/output",
			"your_docker_image_name")
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Failed to run Docker container: %v", err)
			http.Error(w, "Failed to run Docker container", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Docker container launched successfully")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
