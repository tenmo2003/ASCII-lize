package ascii

import (
	"ascii-lize/internal/config"
	"ascii-lize/internal/image"
	"fmt"
	"os"
)

const (
	WIDTH_TO_HEIGHT_ASPECT_CORRECTION = 0.5 // Characters are typically taller than they are wide
)

type Converter struct {
	processor *image.Processor
	scaler    *image.Scaler
}

func NewConverter() *Converter {
	return &Converter{
		processor: image.NewProcessor(),
		scaler:    image.NewScaler(),
	}
}

func (c *Converter) ConvertToASCII(cfg *config.Config) error {
	img, err := c.processor.LoadImage(cfg.ImagePath)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	grayScale := c.processor.ConvertToGrayscale(img)
	width, height := c.processor.GetImageDimensions(img)

	scaleConfig := image.ScaleConfig{
		TargetWidth:      cfg.TargetedWidth,
		PreserveAspect:   true,
		AspectCorrection: WIDTH_TO_HEIGHT_ASPECT_CORRECTION,
	}

	widthScale, heightScale := c.scaler.CalculateScaling(width, height, scaleConfig)

	scaledGrayScale := c.scaler.ScaleGrayscale(grayScale, widthScale, heightScale)

	return c.outputASCII(scaledGrayScale, cfg)
}

func (c *Converter) outputASCII(grayScale [][]int, cfg *config.Config) error {
	characterSet := GetCharacterSet(cfg.CharacterSet)

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

	for _, row := range grayScale {
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

func (c *Converter) grayToCharIndex(grayValue, charSetLength int) int {
	charIndex := int(float32(grayValue) / float32(image.ColorRange8Bit) * float32(charSetLength))
	if charIndex >= charSetLength {
		charIndex = charSetLength - 1
	}
	return charIndex
}
