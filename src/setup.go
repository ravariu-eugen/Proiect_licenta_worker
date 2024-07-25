package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

type JobInfo struct {
	ImageName   string   `json:"imageName"`
	OutputFiles []string `json:"outputFiles"`
	SharedFiles []string `json:"sharedFiles"`
}

type JobInfoMap map[string]JobInfo

var jobInfo JobInfoMap = make(JobInfoMap)

func SetJobInfo(c *gin.Context) {
	if err := c.ShouldBindJSON(&jobInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": jobInfo})
}
