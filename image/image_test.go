package image

import (
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"testing"

	"github.com/LonelyPale/goutils/image/text"
)

func TestCut(t *testing.T) {
	filePath := "../testdata/text/background.jpg"

	//打开图片文件
	imgFile, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := imgFile.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	bgImg, err := jpeg.Decode(imgFile)
	if err != nil {
		t.Fatal(err)
	}

	width := bgImg.Bounds().Dx()
	height := bgImg.Bounds().Dy() - 50
	offset := image.Pt(0, -50)
	rectangle := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rectangle)
	draw.Draw(img, bgImg.Bounds().Add(offset), bgImg, image.Point{}, draw.Over)

	if err := text.WriteFile(img, 0, "../testdata/1.png"); err != nil {
		t.Fatal(err)
	}

	if err := text.WriteFile(img, 100, "../testdata/2.jpg"); err != nil {
		t.Fatal(err)
	}
}
