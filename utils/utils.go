package utils

import (
	"strconv"
	"strings"
)

// IntSliceToString converts a slice of integers to a comma-separated string.
func IntSliceToString(slice []int) string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = strconv.Itoa(v)
	}
	return strings.Join(strSlice, ",")
}

// StringToIntSlice converts a comma-separated string to a slice of integers.
func StringToIntSlice(str string) ([]int, error) {
	if str == "" {
		return []int{}, nil
	}
	strSlice := strings.Split(str, ",")
	intSlice := make([]int, len(strSlice))
	for i, s := range strSlice {
		val, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		intSlice[i] = val
	}
	return intSlice, nil
}
