package main

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

func getSystemInfo(c *gin.Context) {
	// return system cpu and memory usage
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	c.JSON(http.StatusOK, gin.H{
		"cpuUsage": stats.Sys,
		"memUsage": stats.HeapAlloc,
		"totalMem": stats.HeapSys,
	})

}
