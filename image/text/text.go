package image

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

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

//文本：内容及坐标
type TextPt struct {
	Value string `json:"value"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
}

//字体：路径、像素、尺码、颜色
type TextFont struct {
	Path  string      `json:"path"`
	DPI   float64     `json:"dpi"`
	Size  float64     `json:"size"`
	Color color.Color `json:"color"`
}

//字体综合
type Text struct {
	TextFont
	Tps []TextPt `json:"tps"`
}

//读取图片信息
func readImage(filePath string) (*image.RGBA, error) {
	var imgFile *os.File
	var bgImg image.Image //background image

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
		bgImg, err = png.Decode(imgFile)
		if err != nil {
			return nil, fmt.Errorf("解析png图片 %s", err)
		}
	} else if fSuffix == ".jpg" || fSuffix == ".jpeg" {
		// 解析jpeg图片
		bgImg, err = jpeg.Decode(imgFile)
		if err != nil {
			return nil, fmt.Errorf("解析jpeg图片 %s", err)
		}
	}

	//RGBA
	width := bgImg.Bounds().Dx()
	height := bgImg.Bounds().Dy()
	offset := image.Pt(0, 0)
	rectangle := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rectangle)
	draw.Draw(img, bgImg.Bounds().Add(offset), bgImg, image.Point{}, draw.Over)

	/*
		//NRGBA
		img = image.NewNRGBA(bgImg.Bounds())
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 0; x < img.Bounds().Dx(); x++ {
				img.Set(x, y, bgImg.At(x, y))
			}
		}
	*/

	return img, nil
}

//读取字体
func readFont(filePath string) (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile(filePath) //读取字体
	if err != nil {
		return nil, fmt.Errorf("读取字体 ioutil.ReadFile error: %s", err)
	}

	font, err := freetype.ParseFont(fontBytes) //载入字体
	if err != nil {
		return nil, fmt.Errorf("载入字体 freetype.ParseFont error: %s", err)
	}

	return font, nil
}

//设置文字格式
func setTextType(font *truetype.Font, img *image.RGBA, dpi, fontSize float64, cr color.Color) *freetype.Context {
	f := freetype.NewContext()
	f.SetDPI(dpi)                  //设置分辨率
	f.SetFont(font)                //设置字体
	f.SetFontSize(fontSize)        //设置尺寸
	f.SetClip(img.Bounds())        //设置用于绘制的剪辑矩形。
	f.SetDst(img)                  //设置输出的图片
	f.SetSrc(image.NewUniform(cr)) //设置用于绘制操作的源图像(字体颜色)
	return f
}

//写入文字
func writeText(text Text, img *image.RGBA) error {
	//读取字体信息
	var font *truetype.Font
	font, err := readFont(text.Path)
	if err != nil {
		return fmt.Errorf("读取字体信息 %s", err)
	}

	//设置文字格式
	var ctx *freetype.Context
	ctx = setTextType(font, img, text.DPI, text.Size, text.Color)

	if len(text.Tps) > 0 {
		for _, tp := range text.Tps {
			pt := freetype.Pt(tp.X, tp.Y)
			_, err = ctx.DrawString(tp.Value, pt)
			if err != nil {
				return fmt.Errorf("设置文字格式 %s", err)
			}
		}
	}

	return nil
}

//写入文件
func writeFile(img *image.RGBA, quality int, filePath string) (err error) {
	if filePath == "" {
		filePath = "new_file.jpg"
	}

	var newFile *os.File
	newFile, err = os.Create(filePath)
	if err != nil {
		return err
	}

	err = jpeg.Encode(newFile, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return err
	}

	return nil
}

//生成文字图片
func GenerateTextImage(texts []Text, bgImagePath string, quality int) ([]byte, error) {
	bgimg, err := readImage(bgImagePath)
	if err != nil {
		return nil, err
	}

	if len(texts) > 0 {
		for _, text := range texts {
			if err := writeText(text, bgimg); err != nil {
				return nil, err
			}
		}
	}

	buff := new(bytes.Buffer)
	if err := jpeg.Encode(buff, bgimg, &jpeg.Options{Quality: quality}); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
