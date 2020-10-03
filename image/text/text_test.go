package image

import (
	"image/color"
	"io/ioutil"
	"os"
	"testing"
)

func TestGenerateTextImage(t *testing.T) {
	text := Text{
		TextFont: TextFont{
			Path:  "/Users/wyb/project/github/byinfo/testdata/font/FZSJ-BAOZNH.TTF",
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

	bgimg := "/Users/wyb/project/github/byinfo/testdata/image/background.jpg"
	bs, err := GenerateTextImage([]Text{text}, bgimg, 80)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("/Users/wyb/project/github/byinfo/testdata/text.jpg", bs, os.FileMode(0666)); err != nil {
		t.Fatal(err)
	}
}
