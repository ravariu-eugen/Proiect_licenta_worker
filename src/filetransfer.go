package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// uploadFile uploads a file to the server.
//
// It takes a gin.Context as a parameter and returns a string and an error.
// The string is the path of the uploaded file, and the error is any error that occurred during the upload.
func uploadFile(c *gin.Context, destinationDir string) (string, error) {
	_, header, err := c.Request.FormFile("file")

	if err != nil {
		return "", err
	}

	newFile := destinationDir + "/" + header.Filename
	err = c.SaveUploadedFile(header, newFile)
	if err != nil {
		return "", err
	}
	return newFile, nil

}

func uploadAndExtractToDir(c *gin.Context, destinationDir string) (string, error) {
	file, _, err := c.Request.FormFile("file")

	if err != nil {
		return "", err
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "err.Error()1"})
	return extractFile(file.(*os.File), destinationDir)

}

func getFileList(c *gin.Context, dir string) {

	files, err := listFiles(dir)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})

}

type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func listFiles(dir string) ([]FileInfo, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var result []FileInfo
	for _, file := range files {
		if !file.IsDir() {

			info, err := file.Info()
			if err != nil {
				return nil, err
			}

			result = append(result, FileInfo{
				Name: file.Name(),
				Size: info.Size(),
			})
		} else {

			result = append(result, FileInfo{
				Name: file.Name() + "/",
				Size: 0,
			})
		}
	}

	return result, nil
}
