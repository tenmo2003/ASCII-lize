package image

type ScaleConfig struct {
	TargetWidth      int
	AspectCorrection float64
}

type Scaler struct{}

func NewScaler() *Scaler {
	return &Scaler{}
}

func (s *Scaler) CalculateScaling(width, height int, config ScaleConfig) (int, int) {
	targetWidth := min(width, config.TargetWidth)
	widthScaleRate := width / targetWidth

	var targetHeight, heightScaleRate int
	aspectRatio := float64(height) / float64(width)
	targetHeight = int(float64(targetWidth) * aspectRatio * config.AspectCorrection)
	heightScaleRate = height / targetHeight

	return widthScaleRate, heightScaleRate
}

func (s *Scaler) ScaleGrayscale(grayScale [][]int, widthScale, heightScale int) [][]int {
	if len(grayScale) == 0 || len(grayScale[0]) == 0 {
		return [][]int{}
	}

	scaledHeight := len(grayScale) / heightScale
	scaledWidth := len(grayScale[0]) / widthScale

	if scaledHeight <= 0 || scaledWidth <= 0 {
		return [][]int{}
	}

	scaledGrayScale := make([][]int, scaledHeight)

	for i := range scaledHeight {
		scaledGrayScale[i] = make([]int, scaledWidth)
		scaledI := i * heightScale

		for j := range scaledWidth {
			scaledJ := j * widthScale
			sum := 0
			pixelCount := 0

			for m := 0; m < heightScale && scaledI+m < len(grayScale); m++ {
				for n := 0; n < widthScale && scaledJ+n < len(grayScale[scaledI+m]); n++ {
					sum += grayScale[scaledI+m][scaledJ+n]
					pixelCount++
				}
			}

			if pixelCount > 0 {
				scaledGrayScale[i][j] = sum / pixelCount
			}
		}
	}

	return scaledGrayScale
}
