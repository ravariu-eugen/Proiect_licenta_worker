package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func addImage(c *gin.Context) {

	// adds an image archive to the /app/images folder

	uploadFile(c, ImageFolder)
}

func addShared(c *gin.Context) {

	uploadFile(c, SharedFolder)

}

func getImages(c *gin.Context) {
	getFileList(c, ImageFolder)
}

func getShared(c *gin.Context) {
	getFileList(c, SharedFolder)
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

func uploadFile(c *gin.Context, dst string) {
	file, header, err := c.Request.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"stage": 1, "error": err.Error()})
		return
	}
	defer file.Close()
	newFile := dst + "/" + header.Filename
	out, err := os.Create(newFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"stage": 2, "file": newFile, "error": err.Error()})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"stage": 3, "error": err.Error()})
		return
	}

	extractInPlace(newFile)

	c.JSON(http.StatusOK, gin.H{
		"message": "Uploaded the image successfully: " + header.Filename,
	})

}
