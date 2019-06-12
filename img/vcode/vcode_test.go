package vcode

import (
	"github.com/disintegration/imaging"
	"github.com/hetianyi/gox/font"
	"image"
	"image/color"
	"testing"
)

var words = "HelloL你好"

// https://stackoverflow.com/questions/27631736/meaning-of-top-ascent-baseline-descent-bottom-and-leading-in-androids-font
func TestVCode_Generate1(t *testing.T) {
	f, _ := font.LoadFont("E:\\test1\\Inkfree.ttf")
	vc := &VCode{
		Font:            f,
		Size:            image.Rect(0, 0, 300, 150),
		Colors:          []color.Color{image.White},
		FontSize:        50,
		BackgroundColor: image.White,
	}
	im1 := vc.Generate(words, imaging.TopLeft, 0, 0)
	im2 := vc.Generate(words, imaging.Top, 0, 0)
	im3 := vc.Generate(words, imaging.Left, 0, 0)
	im4 := vc.Generate(words, imaging.Bottom, 0, 0)
	im5 := vc.Generate(words, imaging.BottomLeft, 0, 0)
	im6 := vc.Generate(words, imaging.Center, 0, 0)
	imaging.Save(im1, "E:\\test1\\Inkfree-TopLeft.jpg")
	imaging.Save(im2, "E:\\test1\\Inkfree-Top.jpg")
	imaging.Save(im3, "E:\\test1\\Inkfree-Left.jpg")
	imaging.Save(im4, "E:\\test1\\Inkfree-Bottom.jpg")
	imaging.Save(im5, "E:\\test1\\Inkfree-BottomLeft.jpg")
	imaging.Save(im6, "E:\\test1\\Inkfree-Center.jpg")
}

func TestVCode_Generate2(t *testing.T) {
	f, _ := font.LoadFont("E:\\test1\\SIMLI.TTF")
	vc := &VCode{
		Font:            f,
		Size:            image.Rect(0, 0, 300, 150),
		Colors:          []color.Color{image.White},
		FontSize:        50,
		BackgroundColor: image.White,
	}
	im1 := vc.Generate(words, imaging.TopLeft, 0, 0)
	im2 := vc.Generate(words, imaging.Top, 0, 0)
	im3 := vc.Generate(words, imaging.Left, 0, 0)
	im4 := vc.Generate(words, imaging.Bottom, 0, 0)
	im5 := vc.Generate(words, imaging.BottomLeft, 0, 0)
	im6 := vc.Generate(words, imaging.Center, 0, 0)
	imaging.Save(im1, "E:\\test1\\SIMLI-TopLeft.jpg")
	imaging.Save(im2, "E:\\test1\\SIMLI-Top.jpg")
	imaging.Save(im3, "E:\\test1\\SIMLI-Left.jpg")
	imaging.Save(im4, "E:\\test1\\SIMLI-Bottom.jpg")
	imaging.Save(im5, "E:\\test1\\SIMLI-BottomLeft.jpg")
	imaging.Save(im6, "E:\\test1\\SIMLI-Center.jpg")
}

func TestVCode_Generate3(t *testing.T) {
	f, _ := font.LoadFont("E:\\test1\\STXINGKA.TTF")
	vc := &VCode{
		Font:            f,
		Size:            image.Rect(0, 0, 300, 150),
		Colors:          []color.Color{image.White},
		FontSize:        50,
		BackgroundColor: image.White,
	}
	im1 := vc.Generate(words, imaging.TopLeft, 0, 0)
	im2 := vc.Generate(words, imaging.Top, 0, 0)
	im3 := vc.Generate(words, imaging.Left, 0, 0)
	im4 := vc.Generate(words, imaging.Bottom, 0, 0)
	im5 := vc.Generate(words, imaging.BottomLeft, 0, 0)
	im6 := vc.Generate(words, imaging.Center, 0, 0)
	imaging.Save(im1, "E:\\test1\\STXINGKA-TopLeft.jpg")
	imaging.Save(im2, "E:\\test1\\STXINGKA-Top.jpg")
	imaging.Save(im3, "E:\\test1\\STXINGKA-Left.jpg")
	imaging.Save(im4, "E:\\test1\\STXINGKA-Bottom.jpg")
	imaging.Save(im5, "E:\\test1\\STXINGKA-BottomLeft.jpg")
	imaging.Save(im6, "E:\\test1\\STXINGKA-Center.jpg")
}

func TestVCode_Generate4(t *testing.T) {
	f, _ := font.LoadFont("E:\\test1\\BAUHS93.TTF")
	vc := &VCode{
		Font:            f,
		Size:            image.Rect(0, 0, 300, 150),
		Colors:          []color.Color{image.White},
		FontSize:        50,
		BackgroundColor: image.White,
	}
	im1 := vc.Generate(words, imaging.TopLeft, 0, 0)
	im2 := vc.Generate(words, imaging.Top, 0, 0)
	im3 := vc.Generate(words, imaging.Left, 0, 0)
	im4 := vc.Generate(words, imaging.Bottom, 0, 0)
	im5 := vc.Generate(words, imaging.BottomLeft, 0, 0)
	im6 := vc.Generate(words, imaging.Center, 0, 0)
	imaging.Save(im1, "E:\\test1\\BAUHS93-TopLeft.jpg")
	imaging.Save(im2, "E:\\test1\\BAUHS93-Top.jpg")
	imaging.Save(im3, "E:\\test1\\BAUHS93-Left.jpg")
	imaging.Save(im4, "E:\\test1\\BAUHS93-Bottom.jpg")
	imaging.Save(im5, "E:\\test1\\BAUHS93-BottomLeft.jpg")
	imaging.Save(im6, "E:\\test1\\BAUHS93-Center.jpg")
}

func TestVCode_Generate5(t *testing.T) {
	f, _ := font.LoadFont("E:\\test1\\Inkfree.ttf")
	vc := &VCode{
		Font:            f,
		Size:            image.Rect(0, 0, 300, 150),
		Colors:          []color.Color{image.White},
		FontSize:        50,
		BackgroundColor: image.White,
	}
	words := ""
	im1 := vc.Generate(words, imaging.TopLeft, 0, 0)
	im2 := vc.Generate(words, imaging.Top, 0, 0)
	im3 := vc.Generate(words, imaging.Left, 0, 0)
	im4 := vc.Generate(words, imaging.Bottom, 0, 0)
	im5 := vc.Generate(words, imaging.BottomLeft, 0, 0)
	im6 := vc.Generate(words, imaging.Center, 0, 0)
	imaging.Save(im1, "E:\\test1\\BAUHS93-TopLeft.jpg")
	imaging.Save(im2, "E:\\test1\\BAUHS93-Top.jpg")
	imaging.Save(im3, "E:\\test1\\BAUHS93-Left.jpg")
	imaging.Save(im4, "E:\\test1\\BAUHS93-Bottom.jpg")
	imaging.Save(im5, "E:\\test1\\BAUHS93-BottomLeft.jpg")
	imaging.Save(im6, "E:\\test1\\BAUHS93-Center.jpg")
}
