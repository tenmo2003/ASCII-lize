package image

import (
	"ascii-lize/internal/utils"
	"image"
)

const (
	ColorRange8Bit = 256 // 2^8
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) LoadImage(path string) (image.Image, error) {
	return utils.GetImageFromFilePath(path)
}

func (p *Processor) ConvertToGrayscale(img image.Image) [][]int {
	bounds := img.Bounds()
	grayScale := make([][]int, bounds.Max.Y)

	for y := 0; y < bounds.Max.Y; y++ {
		grayScale[y] = make([]int, bounds.Max.X)
		for x := 0; x < bounds.Max.X; x++ {
			color := img.At(x, y)
			r, g, b, _ := color.RGBA()

			// NOTE:
			// Use 8-bit RGB as it works better for ASCII art
			// Use luminance formula for better grayscale conversion - https://en.wikipedia.org/wiki/Relative_luminance
			grayScaleValue := uint32(0.2126*float64(r>>8) + 0.7152*float64(g>>8) + 0.0722*float64(b>>8))
			grayScale[y][x] = int(grayScaleValue)
		}
	}

	return grayScale
}

func (p *Processor) GetImageDimensions(img image.Image) (int, int) {
	bounds := img.Bounds()
	return bounds.Max.X, bounds.Max.Y
}
