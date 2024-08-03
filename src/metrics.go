package main

import (
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func getCPUUsageAndMemoryUtilization() (cpuUsage float64, memoryUtilization float64, err error) {
	// Get CPU usage
	cmd := exec.Command("mpstat", " 1 5 | awk 'END{print 100-$NF\" % \"}'")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	cpuUsageStr := strings.TrimSpace(string(output))
	cpuUsage, err = strconv.ParseFloat(cpuUsageStr, 64)
	if err != nil {
		return 0, 0, err
	}

	// Get total memory utilization
	cmd = exec.Command("sh", "-c", "top -bn1 | grep Mem | awk '{print $3}'")
	output, err = cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	memoryUtilizationStr := strings.TrimSpace(string(output))
	memoryUtilization, err = strconv.ParseFloat(memoryUtilizationStr, 64)
	if err != nil {
		return 0, 0, err
	}

	return cpuUsage, memoryUtilization, nil
}

func getMetrics(c *gin.Context) {
	cpuUsage, memoryUtilization, err := getCPUUsageAndMemoryUtilization()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"cpuUsage":          cpuUsage,
		"memoryUtilization": memoryUtilization,
	})

}
