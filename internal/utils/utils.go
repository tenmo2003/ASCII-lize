package utils

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
)

func GetImageFromFilePath(filePath string) (image.Image, error) {
	filePath = ExpandPath(filePath)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func ExpandPath(filePath string) string {
	if strings.HasPrefix(filePath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		filePath = strings.Replace(filePath, "~", homeDir, 1)
	} else if strings.HasPrefix(filePath, ".") {
		currentDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		filePath = strings.Replace(filePath, ".", currentDir, 1)
	}
	return filePath
}
