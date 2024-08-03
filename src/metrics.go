package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func getCPUUsageAndMemoryUtilization() (cpuUsage float64, memoryUtilization float64, err error) {
	// Get CPU usage
	cmd := exec.Command("bash", "-c", "mpstat 1 5 | tail -n 1 | awk '{print 100-$NF}'")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get CPU usage: %v", err)
	}
	cpuUsageStr := strings.TrimSpace(string(output))
	cpuUsage, err = strconv.ParseFloat(cpuUsageStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse CPU usage: %v", err)
	}

	// Get total memory utilization
	cmd = exec.Command("bash", "-c", "free | grep Mem | awk '{print $3/$2 * 100.0}'")
	output, err = cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get total memory utilization: %v", err)
	}
	memoryUtilizationStr := strings.TrimSpace(string(output))
	memoryUtilization, err = strconv.ParseFloat(memoryUtilizationStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse total memory utilization: %v", err)
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
