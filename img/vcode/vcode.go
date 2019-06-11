package vcode

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/hetianyi/gox/img"
	"golang.org/x/image/font"
	"image"
	"image/color"
)

const (
	DefaultDPI = 72
)

type VCode struct {
	Font            *truetype.Font
	FontSize        float64
	Colors          []color.Color
	BackgroundColor color.Color
	Size            image.Rectangle
	ctx             *freetype.Context
}

func (v *VCode) Generate(content string, anchor imaging.Anchor, paddingX int, paddingY int) *image.RGBA {
	target := img.NewRGBA(v.Size)
	if v.ctx == nil {
		v.ctx = freetype.NewContext()
		v.ctx.SetDPI(DefaultDPI)
		v.ctx.SetFont(v.Font)
		v.ctx.SetFontSize(v.FontSize)
		v.ctx.SetClip(v.Size)
		v.ctx.SetDst(target) // 将freetype绑定到该画布
		v.ctx.SetSrc(image.NewUniform(v.Colors[0]))
		v.ctx.SetHinting(font.HintingNone)
	}

	opt := truetype.Options{
		Size: v.FontSize,
		DPI:  DefaultDPI,
	}
	face := truetype.NewFace(v.Font, &opt)
	m := face.Metrics()
	fmt.Println("Ascent: ", m.Ascent.Floor(), "Descent: ", m.Descent.Ceil(), "Height: ", m.Height)

	pot := img.CalculatePt(v.Size.Max, image.Point{0, 0}, anchor, paddingX, paddingY)
	fmt.Println(freetype.Pt(pot.X+int(v.FontSize), pot.Y))
	_, err := v.ctx.DrawString(content, freetype.Pt(pot.X, pot.Y+m.Ascent.Ceil()-m.Descent.Ceil()/2))
	if err != nil {
		return nil
	}
	return target
}
