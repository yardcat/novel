package util

import (
	"path/filepath"
	"strings"
)

func GetFileNameWithoutExt(path string) string {
	baseName := filepath.Base(path)
	ext := filepath.Ext(baseName)
	fileNameWithoutExt := strings.TrimSuffix(baseName, ext)
	return fileNameWithoutExt
}
