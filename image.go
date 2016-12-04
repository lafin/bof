package main

import (
	"bytes"
	"github.com/x1ddos/imgdiff"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func decodeImage(data []byte) (image.Image, error) {
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Compare - return percentage the difference of images
func Compare(bytes1, bytes2 []byte) (float64, error) {
	newDiffer := imgdiff.NewPerceptual(2.2, 100.0, 45.0, 1.0, false)
	image1, err := decodeImage(bytes1)
	if err != nil {
		return 0, err
	}
	image2, err := decodeImage(bytes2)
	if err != nil {
		return 0, err
	}
	var res image.Image
	var n int
	res, n, err = newDiffer.Compare(image1, image2)
	if err != nil {
		if err.Error() == "images have different sizes" {
			return 100, nil
		}
		return 0, err
	}
	return float64(n) / float64(res.Bounds().Dx()*res.Bounds().Dy()), nil
}
