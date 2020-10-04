package image

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//读取图片
func ReadFile(filePath string) (image.Image, error) {
	var imgFile *os.File
	var img image.Image //background image

	//打开图片文件
	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开图片文件 %s", err)
	}
	defer func() {
		if err := imgFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	fSuffix := filepath.Ext(filePath)
	fSuffix = strings.ToLower(fSuffix)
	if fSuffix == ".png" {
		// 解析png图片
		img, err = png.Decode(imgFile)
		if err != nil {
			return nil, fmt.Errorf("解析png图片 %s", err)
		}
	} else if fSuffix == ".jpg" || fSuffix == ".jpeg" {
		// 解析jpeg图片
		img, err = jpeg.Decode(imgFile)
		if err != nil {
			return nil, fmt.Errorf("解析jpeg图片 %s", err)
		}
	} else {
		return nil, fmt.Errorf("未知的图片格式 %s", fSuffix)
	}

	return img, nil
}

//写入文件
func WriteFile(img image.Image, quality int, filename string) (err error) {
	if filename == "" {
		filename = "new_file.png"
	}

	var bs []byte
	fileSuffix := filepath.Ext(filename)
	switch fileSuffix {
	case ".png":
		bs, err = Encode(img, quality, "png")
	case ".jpg":
		bs, err = Encode(img, quality, "jpg")
	case ".jpeg":
		bs, err = Encode(img, quality, "jpeg")
	default:
		return fmt.Errorf("未知的图片扩展名 %s", fileSuffix)
	}
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bs, os.FileMode(0666))
}

func Clone(src image.Image) draw.Image {
	//RGBA
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()
	offset := image.Pt(0, 0)
	rectangle := image.Rect(0, 0, width, height)
	dst := image.NewRGBA(rectangle)
	draw.Draw(dst, src.Bounds().Add(offset), src, image.Point{}, draw.Over)

	/*
		//NRGBA
		img = image.NewNRGBA(bgImg.Bounds())
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 0; x < img.Bounds().Dx(); x++ {
				img.Set(x, y, bgImg.At(x, y))
			}
		}
	*/

	return dst
}

func Encode(img image.Image, quality int, typ string) ([]byte, error) {
	buff := new(bytes.Buffer)

	typ = strings.ToLower(typ)
	if typ == "png" {
		if err := png.Encode(buff, img); err != nil {
			return nil, err
		}
	} else if typ == "jpg" || typ == "jpeg" {
		if err := jpeg.Encode(buff, img, &jpeg.Options{Quality: quality}); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("未知的图片类型 %s", typ)
	}

	return buff.Bytes(), nil
}
