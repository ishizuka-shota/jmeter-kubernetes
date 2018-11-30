package util

// GetSliceNotBlank 配列の中身から空白を除く
func GetSliceNotBlank(slice []string) []string {

	result := make([]string, 0)

	for _, e := range slice {
		if e != "" {
			result = append(result, e)
		}
	}
	return result
}
