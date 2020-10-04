package text

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"

	"github.com/LonelyPale/goutils/errors"
	goimage "github.com/LonelyPale/goutils/image"
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
func setTextType(font *truetype.Font, img image.Image, dpi, fontSize float64, cr color.Color) *freetype.Context {
	f := freetype.NewContext()
	f.SetDPI(dpi)                  //设置分辨率
	f.SetFont(font)                //设置字体
	f.SetFontSize(fontSize)        //设置尺寸
	f.SetClip(img.Bounds())        //设置用于绘制的剪辑矩形。
	f.SetDst(img.(draw.Image))     //设置输出的图片
	f.SetSrc(image.NewUniform(cr)) //设置用于绘制操作的源图像(字体颜色)
	return f
}

//写入文字
func writeText(text Text, img image.Image) error {
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

//生成文字图片
func GenerateTextImage(texts []Text, bgImage image.Image) (image.Image, error) {
	if bgImage == nil {
		return nil, errors.New("background is not nil")
	}

	img, ok := bgImage.(draw.Image)
	if !ok {
		img = goimage.Clone(bgImage)
	}

	if len(texts) > 0 {
		for _, text := range texts {
			if err := writeText(text, img); err != nil {
				return nil, err
			}
		}
	}

	return img, nil
}

//生成文字图片文件
func GenerateTextImageFile(texts []Text, bgImagePath string) (image.Image, error) {
	img, err := goimage.ReadFile(bgImagePath)
	if err != nil {
		return nil, err
	}

	return GenerateTextImage(texts, img)
}
