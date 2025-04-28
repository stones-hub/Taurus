package util

import (
	"strings"
)

// InArray tests if a string is within an array. It is not case sensitive.
func InArray(k string, arr []string) bool {
	k = strings.ToLower(k)
	for _, v := range arr {
		if strings.ToLower(v) == k {
			return true
		}
	}

	return false
}
