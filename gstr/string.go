package gstr

import (
	"strconv"
	"strings"
)

func Str2Int(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		// handle error
		return 0
	}
	return i
}

func Str2Float(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0.0
	}
	return f
}

func Str2Bool(str string) bool {
	v := strings.ToLower(str)
	if v == "false" || v == "0" || v == "" || v == "no" || v == "n" || v == "f" || v == "off" || v == "[]" || v == "{}" || v == "null" {
		return false
	} else {
		return true
	}
}

func Join(str ...string) string {
	return strings.Join(str, "")
}
