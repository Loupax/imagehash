package imagehash

import (
	"github.com/disintegration/imaging"
	"image"
)

func grayscale(img image.Image) ([][]float64, error) {
	bounds := img.Bounds()
	imgScale := int(floorp2(min(bounds.Dx(), bounds.Dy())))
	resized := imaging.Resize(img, imgScale, imgScale, imaging.Cosine)

	data := make([][]float64, imgScale)
	for y := 0; y < imgScale; y++ {
		data[y] = make([]float64, imgScale)
		for x := 0; x < imgScale; x++ {
			pixel := resized.At(x, y)
			r, g, b, _ := pixel.RGBA()
			data[y][x] = (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 65535.0
		}
	}
	return data, nil
}
