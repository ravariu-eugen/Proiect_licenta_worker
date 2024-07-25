package main

import (
	"fmt"
	"os"
	"os/exec"
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

func unarchiveFile(filepath, destinationFolder string) error {

	filename := filepath[strings.LastIndex(filepath, "/")+1:]
	fileType := filename[:strings.Index(filepath, ".")]

	err := os.MkdirAll(destinationFolder, os.ModePerm)
	if err != nil {
		return err
	}
	var cmd *exec.Cmd
	switch fileType {
	case "tar":
		cmd = exec.Command("tar", "-xf", filepath, "-C", destinationFolder)
	case "zip":
		cmd = exec.Command("unzip", filepath, "-d", destinationFolder)
	default:
		return fmt.Errorf("unsupported file type: %s", fileType)
	}

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
