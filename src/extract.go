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
func extractFileFromPath(filePath, destinationFolder string) (string, error) {
	extension := filepath.Ext(filePath)
	_ = getFileNameWithoutExt(filePath)
	newDir := filepath.Join(destinationFolder, filepath.Base(filePath))
	var cmd *exec.Cmd

	switch extension {
	case ".tar":
		cmd = exec.Command("tar", "-xf", filePath, "-C", newDir)
	case ".zip":
		cmd = exec.Command("unzip", "-d", newDir, filePath)
	default:
		return "", fmt.Errorf("unsupported file type: %s", extension)
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to extract file: %v", err)
	}

	return newDir, nil
}

func saveLocal(file *multipart.FileHeader) (string, error) {
	tmpDir, err := os.MkdirTemp("", "temp")
	if err != nil {
		return "", err
	}

	newFilePath := filepath.Join(tmpDir, filepath.Base(file.Filename))
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(newFilePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return newFilePath, nil
}

func extractMultipartFile(file *multipart.FileHeader, destinationFolder string) (string, error) {

	// save the file
	newFile, err := saveLocal(file)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}
	defer os.Remove(newFile)

	return extractFileFromPath(newFile, destinationFolder)
}
