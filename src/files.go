package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// downloadFile downloads a file from a specified S3 bucket and saves it to a given destination folder.
//
// Parameters:
//   - bucketName: the name of the S3 bucket from which to download the file (string)
//   - filename: the name of the file to download (string)
//   - destinationFolder: the folder where the downloaded file will be saved (string)
//
// Returns:
//   - error: an error if the file download fails, otherwise nil (error)

func extractInPlace(file string) error {
	// Extract the name of the archive
	archiveName := filepath.Base(file)
	// Remove the .zip extension
	archiveName = archiveName[:len(archiveName)-len(".zip")]
	// Create the destination folder
	destinationFolder := filepath.Dir(file) + "/" + archiveName
	fmt.Println(destinationFolder)
	err := os.Mkdir(destinationFolder, 0755)
	if err != nil {
		return err
	}

	return extractFile(file, destinationFolder)
}

func extractFile(filepath, destinationFolder string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	var cmd *exec.Cmd

	switch {
	case strings.HasSuffix(filepath, ".tar"):
		cmd = exec.Command("tar", "-xf", "-", "-C", destinationFolder)
		cmd.Stdin = file
	case strings.HasSuffix(filepath, ".zip"):
		cmd = exec.Command("unzip", "-d", destinationFolder)
		cmd.Stdin = file
	default:
		return fmt.Errorf("unsupported file type: %s", filepath)
	}

	cmd.Stderr = os.Stderr
	return cmd.Run()
}
