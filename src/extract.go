package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
)

func extractInPlace(file string) (string, error) {
	// get the destination folder
	destinationFolder := filepath.Dir(file)

	return extractFileFromPath(file, destinationFolder)
}

// extractFile extracts the contents of an archive to a destination folder.
// the archive type is determined by the file extension.
// the archive <name>.<ext> creates a new directory named <name> in the destination folder.
//
// Parameters:
// - filepath: the path of the file to extract.
// - destinationFolder: the folder where the extracted contents will be placed.
//
// Returns:
// - string: the path of the newly created directory where the extracted contents were placed.
// - error: an error if the file could not be opened or if the file type is not supported.
func extractFileFromPath(filepath, destinationFolder string) (string, error) {

	// open the file
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return extractFile(file, destinationFolder)
}

func extractFile(file *os.File, destinationFolder string) (string, error) {
	var cmd *exec.Cmd

	// get the extension
	extension := filepath.Ext(file.Name())
	switch extension { // choose the correct command to execute
	case ".tar":
		cmd = exec.Command("tar", "-xf", "-", "-C", destinationFolder)
	case ".zip":
		cmd = exec.Command("unzip", "-d", destinationFolder)
	default:
		return "", fmt.Errorf("unsupported file type: %s", extension)
	}
	cmd.Stdin = file
	cmd.Stderr = os.Stderr

	// get the path of the new directory
	newDir := destinationFolder + "/" + filepath.Base(file.Name())
	return newDir, cmd.Run()
}

func saveLocal(file *multipart.FileHeader) (string, error) {
	// create temporary folder

	tmpDir, err := os.MkdirTemp("", "temp")
	if err != nil {
		return "", err
	}

	newFile := tmpDir + "/" + filepath.Base(file.Filename)
	con, err := file.Open()
	if err != nil {
		return "", err
	}
	defer con.Close()

	out, err := os.Create(newFile)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, con)
	if err != nil {
		return "", err
	}

	return newFile, nil

}

func extractMultipartFile(file *multipart.FileHeader, destinationFolder string) (string, error) {

	// save the file
	newFile, err := saveLocal(file)
	if err != nil {
		return "", err
	}

	extractFileFromPath(newFile, destinationFolder)
	return newFile, nil
}
