package vcode

import (
	"github.com/disintegration/imaging"
	"github.com/hetianyi/gox/font"
	"image"
	"image/color"
	"testing"
)

// https://stackoverflow.com/questions/27631736/meaning-of-top-ascent-baseline-descent-bottom-and-leading-in-androids-font
func TestVCode_Generate1(t *testing.T) {
	f, _ := font.LoadFont("E:\\test\\Inkfree.ttf")
	vc := &VCode{
		Font:            f,
		Size:            image.Rect(0, 0, 300, 150),
		Colors:          []color.Color{image.White},
		FontSize:        50,
		BackgroundColor: image.White,
	}
	im2 := vc.Generate("HelloLABCD", imaging.Top, 0, 0)
	im1 := vc.Generate("HelloLABCD", imaging.TopLeft, 0, 0)
	im3 := vc.Generate("HelloLABCD", imaging.Left, 0, 0)
	im4 := vc.Generate("HelloLABCD", imaging.Bottom, 0, 50)
	im5 := vc.Generate("HelloLABCD", imaging.BottomLeft, 0, 50)
	im6 := vc.Generate("HelloLABCD", imaging.Center, 0, 0)
	imaging.Save(im1, "E:\\test\\Inkfree-1.jpg")
	imaging.Save(im2, "E:\\test\\Inkfree-2.jpg")
	imaging.Save(im3, "E:\\test\\Inkfree-3.jpg")
	imaging.Save(im4, "E:\\test\\Inkfree-4.jpg")
	imaging.Save(im5, "E:\\test\\Inkfree-5.jpg")
	imaging.Save(im6, "E:\\test\\Inkfree-6.jpg")
}

func TestVCode_Generate2(t *testing.T) {
	f, _ := font.LoadFont("E:\\test\\simfang.ttf")
	vc := &VCode{
		Font:            f,
		Size:            image.Rect(0, 0, 300, 150),
		Colors:          []color.Color{image.White},
		FontSize:        50,
		BackgroundColor: image.White,
	}
	im2 := vc.Generate("HelloLABCD", imaging.Top, 0, 0)
	im1 := vc.Generate("HelloLABCD", imaging.TopLeft, 0, 0)
	im3 := vc.Generate("HelloLABCD", imaging.Left, 0, 0)
	im4 := vc.Generate("HelloLABCD", imaging.Bottom, 0, 36)
	im5 := vc.Generate("HelloLABCD", imaging.BottomLeft, 0, 36)
	im6 := vc.Generate("HelloLABCD", imaging.Center, 0, 0)
	imaging.Save(im1, "E:\\test\\simfang-1.jpg")
	imaging.Save(im2, "E:\\test\\simfang-2.jpg")
	imaging.Save(im3, "E:\\test\\simfang-3.jpg")
	imaging.Save(im4, "E:\\test\\simfang-4.jpg")
	imaging.Save(im5, "E:\\test\\simfang-5.jpg")
	imaging.Save(im6, "E:\\test\\simfang-6.jpg")
}

func TestVCode_Generate3(t *testing.T) {
	f, _ := font.LoadFont("E:\\test\\simkai.ttf")
	vc := &VCode{
		Font:            f,
		Size:            image.Rect(0, 0, 300, 150),
		Colors:          []color.Color{image.White},
		FontSize:        50,
		BackgroundColor: image.White,
	}
	im2 := vc.Generate("HelloLABCD", imaging.Top, 0, 0)
	im1 := vc.Generate("HelloLABCD", imaging.TopLeft, 0, 0)
	im3 := vc.Generate("HelloLABCD", imaging.Left, 0, 0)
	im4 := vc.Generate("HelloLABCD", imaging.Bottom, 0, 36)
	im5 := vc.Generate("HelloLABCD", imaging.BottomLeft, 0, 36)
	im6 := vc.Generate("HelloLABCD", imaging.Center, 0, 0)
	imaging.Save(im1, "E:\\test\\simkai-1.jpg")
	imaging.Save(im2, "E:\\test\\simkai-2.jpg")
	imaging.Save(im3, "E:\\test\\simkai-3.jpg")
	imaging.Save(im4, "E:\\test\\simkai-4.jpg")
	imaging.Save(im5, "E:\\test\\simkai-5.jpg")
	imaging.Save(im6, "E:\\test\\simkai-6.jpg")
}
