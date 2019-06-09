package img_test

import (
	"github.com/disintegration/imaging"
	"github.com/hetianyi/gox/img"
	"log"
	"testing"
)

func TestImage_AddWaterMark(t *testing.T) {
	im, err := img.OpenLocalFile("D:\\tmp\\5\\1.jpg")
	if err != nil {
		log.Panic(err)
	}
	watermark, err := img.OpenLocalFile("D:\\tmp\\4\\mark1.png")
	if err != nil {
		log.Fatal(err)
	}
	im.AddWaterMark(watermark, imaging.TopLeft, 20, 20, 0.5)
	imaging.Save(im.GetSource(), "D:\\tmp\\5\\test_watermark.jpg")
}
