package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AWSCredentials struct {
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type JobInfo struct {
	ImageName   string   `json:"imageName"`
	ImageInFile bool     `json:"imageInFile"`
	InputFile   string   `json:"inputFile"`
	OutputFile  string   `json:"outputFile"`
	SharedFiles []string `json:"sharedFiles"`
}

// needed credentials

// aws credentials:
//   - access key id
//   - secret access key

// container registry credentials:
//   - username
//   - password
//   - registry

// needed info:
//   - executed image - either in registry or in file
//   - input file
//   - shared files
//   - output file name

var awsCred AWSCredentials
var jobInfo JobInfo

func SetAWSCredentials(c *gin.Context) {
	if err := c.ShouldBindJSON(&awsCred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": awsCred})
}

func SetJobInfo(c *gin.Context) {
	if err := c.ShouldBindJSON(&jobInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": jobInfo})
}
