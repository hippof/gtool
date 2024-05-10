package gfont

import (
	"errors"
	"strings"

	"github.com/flopp/go-findfont"
)

var (
	fontPaths []string
)

func init() {
	fontPaths = findfont.List()
}

func GetCnFont() (string, error) {
	for _, path := range fontPaths {
		if strings.Contains(path, "msyh.ttf") || strings.Contains(path, "simhei.ttf") || strings.Contains(path, "simsun.ttc") || strings.Contains(path, "simkai.ttf") {
			return path, nil
		}
	}
	return "", errors.New("未找到中文字体")
}
