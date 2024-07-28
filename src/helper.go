package main

import (
	"path/filepath"
	"strings"
)

func getFileNameWithoutExt(file string) string {
	file = filepath.Base(file)
	return strings.TrimSuffix(file, filepath.Ext(file))
}
