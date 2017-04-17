package utils

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/oliamb/cutter"
	"github.com/x1ddos/imgdiff"
)

func decodeImage(data []byte) (image.Image, error) {
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func doFitImages(image1, image2 image.Image) (image.Image, image.Image, error) {
	bounds1 := image1.Bounds()
	bounds2 := image2.Bounds()
	width := bounds1.Dx()
	height := bounds1.Dy()

	if bounds1.Dx() != bounds2.Dx() {
		if bounds1.Dx() < bounds2.Dx() {
			width = bounds1.Dx()
		} else {
			width = bounds2.Dx()
		}
	}
	if bounds1.Dy() != bounds2.Dy() {
		if bounds1.Dy() < bounds2.Dy() {
			height = bounds1.Dy()
		} else {
			height = bounds2.Dy()
		}
	}
	if bounds1.Dx() != width || bounds1.Dy() != height {
		croppedImg, err := cutter.Crop(image1, cutter.Config{
			Width:  width,
			Height: height,
		})
		if err != nil {
			return nil, nil, err
		}
		image1 = croppedImg
	}
	if bounds2.Dx() != width || bounds2.Dy() != height {
		croppedImg, err := cutter.Crop(image2, cutter.Config{
			Width:  width,
			Height: height,
		})
		if err != nil {
			return nil, nil, err
		}
		image2 = croppedImg
	}

	return image1, image2, nil
}

// Compare - return percentage the difference of images
func Compare(bytes1, bytes2 []byte) (float64, error) {
	newDiffer := imgdiff.NewPerceptual(2.2, 100.0, 45.0, 1.0, true)
	image1, err := decodeImage(bytes1)
	if err != nil {
		return 0, err
	}
	image2, err := decodeImage(bytes2)
	if err != nil {
		return 0, err
	}

	image1, image2, err = doFitImages(image1, image2)
	if err != nil {
		return 0, err
	}

	var res image.Image
	var n int
	res, n, err = newDiffer.Compare(image1, image2)
	if err != nil {
		return 0, err
	}
	return float64(n) / float64(res.Bounds().Dx()*res.Bounds().Dy()), nil
}
