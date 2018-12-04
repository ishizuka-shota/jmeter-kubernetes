package util

import (
	"strings"
)

// GetSliceNotBlank 配列の中身から空白を除く
func GetSliceNotBlank(slice []string) []string {

	result := make([]string, 0)

	for _, e := range slice {
		switch e {
		case "":
			break
		case "\"":
			break
		default:
			if strings.Contains(e, "\"") {
				result = append(result, strings.Replace(e, "\"", "", 2))
			} else {
				result = append(result, e)
			}
		}
	}
	return result
}
