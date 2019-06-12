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
}

func (v *VCode) Generate(content string, anchor imaging.Anchor, paddingX int, paddingY int) *image.RGBA {
	target := img.NewRGBA(v.Size)
	ctx := freetype.NewContext()
	ctx = freetype.NewContext()
	ctx.SetDPI(DefaultDPI)
	ctx.SetFont(v.Font)
	ctx.SetFontSize(v.FontSize)
	ctx.SetClip(v.Size)
	ctx.SetDst(target) // 将freetype绑定到该画布
	ctx.SetSrc(image.NewUniform(v.Colors[0]))
	ctx.SetHinting(font.HintingNone)

	opt := truetype.Options{
		Size:    v.FontSize,
		DPI:     DefaultDPI,
		Hinting: font.HintingNone,
	}
	face := truetype.NewFace(v.Font, &opt)
	m := face.Metrics()
	fmt.Println("Ascent: ", m.Ascent.Ceil(), "Descent: ",
		m.Descent.Ceil(),
		"XHeight: ", m.XHeight.Ceil(),
		"CapHeight: ", m.CapHeight.Ceil(),
		"Height: ", m.Height.Ceil(),
	)

	offset := 0
	if anchor == imaging.Top || anchor == imaging.TopLeft {
		offset = m.Ascent.Ceil() - m.Descent.Ceil()
	} else if anchor == imaging.Left || anchor == imaging.Right || anchor == imaging.Center {
		offset = (m.Ascent.Ceil() - m.Descent.Ceil()) / 2
	} else if anchor == imaging.BottomLeft || anchor == imaging.Bottom || anchor == imaging.BottomRight {
		offset = -m.Descent.Ceil() + m.Descent.Ceil()
	}

	pot := img.CalculatePt(v.Size.Max, image.Point{0, 0}, anchor, paddingX, paddingY)
	_, err := ctx.DrawString(content, freetype.Pt(pot.X, pot.Y+offset))
	if err != nil {
		return nil
	}
	ctx.DrawString("☻", freetype.Pt(pot.X, pot.Y+offset))
	return target
}
