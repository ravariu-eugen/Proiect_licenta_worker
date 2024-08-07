package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func getemoryUtilization() (float64, error) {

	// Get total memory utilization
	cmd := exec.Command("bash", "-c", "free | grep Mem | awk '{print $3/$2 * 100.0}'")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get total memory utilization: %v", err)
	}
	memoryUtilizationStr := strings.TrimSpace(string(output))
	memoryUtilization, err := strconv.ParseFloat(memoryUtilizationStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse total memory utilization: %v", err)
	}

	return memoryUtilization, nil
}

func getCPUUsage() (float64, error) {
	// Get CPU usage
	cmd := exec.Command("bash", "-c", "mpstat | tail -n 1 | awk '{print 100-$NF}'")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get CPU usage: %v", err)
	}
	cpuUsageStr := strings.TrimSpace(string(output))
	cpuUsage, err := strconv.ParseFloat(cpuUsageStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse CPU usage: %v", err)
	}
	return cpuUsage, nil
}

func getRemainingStorage() (int64, error) {
	cmd := exec.Command("bash", "-c", "df -m / | tail -n 1 | awk '{print $4}'")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get storage: %v", err)
	}
	remainingStorageStr := strings.TrimSpace(string(output))
	remainingStorage, err := strconv.ParseInt(remainingStorageStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse storage: %v", err)
	}
	return remainingStorage, nil
}

func getMetrics(c *gin.Context) {
	cpuUsage, err := getCPUUsage()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// memoryUtilization, err := getemoryUtilization()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	remainingStorage, err := getRemainingStorage()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cpuUsage":          cpuUsage,
		"memoryUtilization": 0,
		"remainingStorage":  remainingStorage,
	})

}
