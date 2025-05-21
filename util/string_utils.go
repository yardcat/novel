package util

import (
	"fmt"
	"strings"
)

func FormatString[T any](template string, items map[string]T) string {
	for key, value := range items {
		template = strings.ReplaceAll(template, fmt.Sprintf("{%v}", key), fmt.Sprintf("%v", value))
	}
	return template
}
