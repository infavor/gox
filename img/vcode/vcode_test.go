package vcode

import (
	"github.com/disintegration/imaging"
	"github.com/hetianyi/gox/font"
	"image"
	"image/color"
	"testing"
)

func TestVCode_Generate(t *testing.T) {
	f, _ := font.LoadFont("E:\\test\\STXINGKA.TTF")
	vc := &VCode{
		Font:            f,
		Size:            image.Rect(0, 0, 300, 150),
		Colors:          []color.Color{image.White},
		FontSize:        50,
		BackgroundColor: image.White,
	}
	im := vc.Generate("HelloLABCDEFGHIJKLMN", imaging.BottomLeft, 0, 50)
	imaging.Save(im, "E:\\test\\vcode.jpg")
}
