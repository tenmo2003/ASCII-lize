package ascii

import (
	"ascii-lize/internal/config"
	internalImage "ascii-lize/internal/image"
	"ascii-lize/internal/utils"
	"fmt"
	"image"
	"image/gif"
	"os"
	"strings"
	"time"
)

const CLEAR_SCREEN = "\033[2J\033[H"

const (
	WIDTH_TO_HEIGHT_ASPECT_CORRECTION = 0.5 // Characters are typically taller than they are wide
)

type Converter struct {
	processor *internalImage.Processor
	scaler    *internalImage.Scaler
}

func NewConverter() *Converter {
	return &Converter{
		processor: internalImage.NewProcessor(),
		scaler:    internalImage.NewScaler(),
	}
}

func (c *Converter) ConvertToASCII(cfg *config.Config) error {
	img, format, err := utils.LoadMediaFromFilePath(cfg.MediaPath)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	fmt.Println(format)

	switch format {
	case "jpeg", "png":
		return c.outputASCIIForImage(img, cfg)
	case "gif":
		return c.outputASCIIForGif(cfg)
	default:
		return fmt.Errorf("Invalid image format")
	}
}

func (c *Converter) extractGrayScaleAndApplyScaling(img image.Image, cfg *config.Config) [][]int {
	grayScale := c.processor.ConvertToGrayscale(img)
	width, height := c.processor.GetImageDimensions(img)

	scaleConfig := internalImage.ScaleConfig{
		TargetWidth:      cfg.TargetedWidth,
		AspectCorrection: WIDTH_TO_HEIGHT_ASPECT_CORRECTION,
	}

	widthScale, heightScale := c.scaler.CalculateScaling(width, height, scaleConfig)

	scaledGrayScale := c.scaler.ScaleGrayscale(grayScale, widthScale, heightScale)
	return scaledGrayScale
}

func (c *Converter) outputASCIIForImage(img image.Image, cfg *config.Config) error {
	scaledGrayScale := c.extractGrayScaleAndApplyScaling(img, cfg)

	characterSet := GetCharacterSet(cfg.CharacterSet)

	characterSet = characterSet + strings.Repeat(" ", cfg.SpaceDensity)

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	outputCfg := cfg.ResolveOutputConfig(currentDir)

	var fileHandle *os.File
	if outputCfg.WriteToFile {
		fileHandle, err = os.OpenFile(outputCfg.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer func() {
			if closeErr := fileHandle.Close(); closeErr != nil {
				fmt.Printf("Error closing file: %v\n", closeErr)
			}
		}()
	}

	fmt.Println(CLEAR_SCREEN)

	for _, row := range scaledGrayScale {
		asciiLine := make([]byte, len(row))

		for j, grayValue := range row {
			charIndex := c.grayToCharIndex(grayValue, len(characterSet))
			asciiLine[j] = characterSet[charIndex]
		}

		if outputCfg.WriteToFile {
			fileHandle.Write(asciiLine)
			fileHandle.WriteString("\n")
		}

		fmt.Println(string(asciiLine))
	}

	if outputCfg.WriteToFile {
		fmt.Printf("ASCII art saved to: %s\n", outputCfg.Path)
	}

	return nil
}

func (c *Converter) outputASCIIForGif(cfg *config.Config) error {
	file, err := os.Open(cfg.MediaPath)
	if err != nil {
		return err
	}

	gif, err := gif.DecodeAll(file)
	for index := range len(gif.Image) {
		curImg := gif.Image[index]
		curDelay := gif.Delay[index]

		c.outputASCIIForImage(curImg, cfg)
		time.Sleep(time.Duration(curDelay) * 10 * time.Millisecond)
	}

	return nil
}

func (c *Converter) grayToCharIndex(grayValue, charSetLength int) int {
	charIndex := int(float32(grayValue) / float32(internalImage.ColorRange8Bit) * float32(charSetLength))
	if charIndex >= charSetLength {
		charIndex = charSetLength - 1
	}
	return charIndex
}
