package qrcode

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
)

type Options struct {
	Size            int
	BackgroundColor color.Color
	ForegroundColor color.Color
	Level           qrcode.RecoveryLevel
	LogoImage       image.Image
	LogoPoint       image.Point
	LogoResize      struct {
		Width  int
		Height int
	}
}

func DefaultOptions() *Options {
	return &Options{
		Size:            100,
		BackgroundColor: color.White,
		ForegroundColor: color.Black,
		Level:           qrcode.Highest,
	}
}

func Generate(text string, opts ...*Options) (image.Image, error) {
	var opt *Options
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = DefaultOptions()
	}

	qrc, err := qrcode.New(text, opt.Level)
	if err != nil {
		return nil, err
	}

	qrc.BackgroundColor = opt.BackgroundColor
	qrc.ForegroundColor = opt.ForegroundColor
	qrImg := qrc.Image(opt.Size)

	if opt.LogoImage == nil {
		return qrImg, nil
	} else {
		//todo: logo img
		newImg := image.NewRGBA(image.Rect(8, 8, 93, 93))
		backgroundImg := qrImg.Bounds()
		draw.Draw(newImg, backgroundImg, qrImg, image.Point{}, draw.Src)

		var logoImg image.Image
		if opt.LogoResize.Width > 0 && opt.LogoResize.Height > 0 {
			logoImg = resize.Resize(20, 20, logoImg, resize.Bilinear) //缩放logo尺寸
		}

		if logoImg != nil {
			offset := image.Pt(opt.LogoPoint.X, opt.LogoPoint.Y)
			draw.Draw(newImg, logoImg.Bounds().Add(offset), logoImg, image.Point{}, draw.Over)
		}

		return newImg, nil
	}
}

func WriteFile(filename string, text string, opts ...*Options) error {
	img, err := Generate(text, opts...)
	if err != nil {
		return err
	}

	buff := new(bytes.Buffer)
	suffix := filepath.Ext(filename)
	suffix = strings.ToLower(suffix)
	if suffix == ".png" {
		if err := png.Encode(buff, img); err != nil {
			return err
		}
	} else if suffix == ".jpg" || suffix == ".jpeg" {
		//todo: 有噪点
		if err := jpeg.Encode(buff, img, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("未知的图片类型 %s", suffix)
	}

	if err := ioutil.WriteFile(filename, buff.Bytes(), os.ModePerm); err != nil {
		return err
	}

	return nil
}
