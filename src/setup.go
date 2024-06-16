package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Credentials struct {
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

var cred Credentials
var jobInfo JobInfo

func setCredentials(awsAccessKey, awsSecretKey string) {

	cred.AccessKeyID = awsAccessKey
	cred.SecretAccessKey = awsSecretKey
}

func setInstanceInfo(c *gin.Context) {

	// sets instance info
	c.JSON(http.StatusOK, gin.H{
		"instanceId":      "i-1234567890abcdef0",
		"instanceType":    "t2.micro",
		"instanceState":   "running",
		"publicIpAddress": "1.2.3.4",
	})
}
