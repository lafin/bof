package main

import (
	"bytes"
	"github.com/x1ddos/imgdiff"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
)

func decodeImage(data []byte) image.Image {
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return img
}

// Compare - return percentage the difference of images
func Compare(img1, img2 []byte) float64 {
	newDiffer := imgdiff.NewPerceptual(2.2, 100.0, 45.0, 1.0, false)
	res, n, err := newDiffer.Compare(decodeImage(img1), decodeImage(img2))
	if err != nil {
		log.Fatal(err)
	}
	return float64(n) / float64(res.Bounds().Dx()*res.Bounds().Dy())
}
