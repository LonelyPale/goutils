package image

import (
	"image"
	"image/draw"
	"testing"
)

func TestCut(t *testing.T) {
	filePath := "../testdata/text/background.jpg"

	bgImg, err := ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}

	width := bgImg.Bounds().Dx()
	height := bgImg.Bounds().Dy() - 50
	offset := image.Pt(0, -50)
	rectangle := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rectangle)
	draw.Draw(img, bgImg.Bounds().Add(offset), bgImg, image.Point{}, draw.Over)

	if err := WriteFile(img, 0, "../testdata/1.png"); err != nil {
		t.Fatal(err)
	}

	if err := WriteFile(img, 100, "../testdata/2.jpg"); err != nil {
		t.Fatal(err)
	}
}
