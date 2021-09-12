package text

import (
	"image/color"
	"testing"

	goimage "github.com/lonelypale/goutils/image"
)

func TestGenerateTextImage(t *testing.T) {
	text := Text{
		TextFont: TextFont{
			Path:  "../../testdata/text/FZSJ-BAOZNH.TTF",
			DPI:   72,
			Size:  30,
			Color: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		},
		Tps: []TextPt{
			{
				Value: "姓名：文字水印",
				X:     10,
				Y:     500,
			},
			{
				Value: "单位：ABCD",
				X:     10,
				Y:     600,
			},
			{
				Value: "交易ID：baf3605108ebe469b9d2381299b02619d62e2dd0962e4050b5c702b287a30408",
				X:     10,
				Y:     700,
			},
			{
				Value: "区块高度：确认中",
				X:     10,
				Y:     800,
			},
		},
	}

	bgimg := "../../testdata/text/background.jpg"
	img, err := GenerateTextImageFile([]Text{text}, bgimg)
	if err != nil {
		t.Fatal(err)
	}

	if err := goimage.WriteFile(img, 0, "../../testdata/text.png"); err != nil {
		t.Fatal(err)
	}

	//jpeg.DefaultQuality
	if err := goimage.WriteFile(img, 100, "../../testdata/text.jpg"); err != nil {
		t.Fatal(err)
	}
}
