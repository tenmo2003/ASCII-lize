package main

import (
	"ascii-lize/internal/utils"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// with ascending order of brightness (black -> white)
// const DISPLAY_CHARACTERS = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
const DISPLAY_CHARACTERS = "8@$e*+!:.  "

const COLOR_RANGE_8_BIT = 256 // 2^8

func main() {
	targetedWidth := flag.Int("targeted-width", 100, "The width of the output that the program will try to fit the image in by characters (default to 100)")
	outputFileName := flag.String("output", "", "Write result to file along with stdout")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Please specify the path to image using argument (e.g ./script ~/test.jpg)")
		return
	}
	if *targetedWidth <= 0 {
		fmt.Println("Width must be positive")
		return
	}

	curWd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path := args[0]

	var writeToFile bool
	var outputPath string
	if *outputFileName == "" {
		writeToFile = false
	} else {
		writeToFile = true
		outputPath = filepath.Join(curWd, *outputFileName)
	}

	image, err := utils.GetImageFromFilePath(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	grayScale := [][]int{}
	imageBounds := image.Bounds()

	for y := range imageBounds.Max.Y {
		grayScale = append(grayScale, []int{})
		for x := range imageBounds.Max.X {
			color := image.At(x, y)
			r, g, b, _ := color.RGBA()
			// NOTE: Use 8-bit RGB as it works better with ASCII
			// Formula reference: Relative luminance - https://en.wikipedia.org/wiki/Relative_luminance
			grayScaleValue := (0.2126*float32(r>>8) + 0.7152*float32(g>>8) + 0.0722*float32(b>>8))
			grayScale[y] = append(grayScale[y], int(grayScaleValue))
		}
	}

	width := min(imageBounds.Max.X, *targetedWidth)
	scaleDownRate := imageBounds.Max.X / width

	aspectRatio := float64(imageBounds.Max.Y) / float64(imageBounds.Max.X)
	height := int(float64(width) * aspectRatio * 0.5)
	heightScaleRate := imageBounds.Max.Y / height

	fmt.Println("Scale", scaleDownRate)

	scaledGrayScale := [][]int{}

	for i := 0; i < len(grayScale)/heightScaleRate; i++ {
		row := []int{}
		scaledI := i * heightScaleRate
		for j := 0; j < len(grayScale[i])/scaleDownRate; j++ {
			scaledJ := j * scaleDownRate
			sum := 0
			pixelCount := 0

			for m := 0; m < heightScaleRate && scaledI+m < len(grayScale); m++ {
				for n := 0; n < scaleDownRate && scaledJ+n < len(grayScale[scaledI+m]); n++ {
					sum += grayScale[scaledI+m][scaledJ+n]
					pixelCount++
				}
			}

			if pixelCount > 0 {
				sum /= pixelCount
			}

			row = append(row, sum)
		}
		scaledGrayScale = append(scaledGrayScale, row)
	}

	var outputFile *os.File
	if writeToFile {
		outputFile, err = os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func() {
			if err := outputFile.Close(); err != nil {
				fmt.Printf("Error closing file: %v\n", err)
			}
		}()
	}

	for i := range scaledGrayScale {
		asciiLine := make([]byte, len(scaledGrayScale[i]))
		for j := range scaledGrayScale[i] {
			grayScaleValue := scaledGrayScale[i][j]

			charIndex := int(float32(grayScaleValue) / float32(COLOR_RANGE_8_BIT) * float32(len(DISPLAY_CHARACTERS)))
			if charIndex >= len(DISPLAY_CHARACTERS) {
				charIndex = len(DISPLAY_CHARACTERS) - 1
			}

			asciiLine[j] = DISPLAY_CHARACTERS[charIndex]
		}

		fmt.Println(string(asciiLine))

		if writeToFile {
			outputFile.Write(asciiLine)
			outputFile.WriteString("\n")
		}
	}
}
