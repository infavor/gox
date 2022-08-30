package gifx_test

import (
	"github.com/disintegration/imaging"
	"github.com/infavor/gox/file"
	"github.com/infavor/gox/img"
	"github.com/infavor/gox/img/gifx"
	"github.com/infavor/gox/logger"
	"image/gif"
	"testing"
)

func init() {
	logger.Init(nil)
}

// GIF打水印
func TestGIFWaterMark1(t *testing.T) {
	g, err := gifx.LoadFromLocalFile("D:\\tmp\\4\\origin.gif")
	if err != nil {
		logger.Fatal(err)
	}
	watermark, err := img.OpenLocalFile("D:\\tmp\\4\\mark1.png")
	if err != nil {
		logger.Fatal(err)
	}
	watermark = watermark.Resize(50, 50, imaging.Lanczos)
	g.AddWaterMark(watermark, imaging.BottomRight, 10, 10, 1)

	of, err := file.CreateFile("D:\\tmp\\4\\origin_watermark.gif")
	if err != nil {
		logger.Panic(err)
	}
	defer of.Close()

	err = gif.EncodeAll(of, g.GetSource())
	if err != nil {
		logger.Panic(err)
	}
}

func TestGenerate(t *testing.T) {
	images := make([]*img.Image, 4)
	images[0], _ = img.OpenLocalFile("D:\\tmp\\5\\1.jpg")
	images[1], _ = img.OpenLocalFile("D:\\tmp\\5\\2.jpg")
	images[2], _ = img.OpenLocalFile("D:\\tmp\\5\\3.jpg")
	images[3], _ = img.OpenLocalFile("D:\\tmp\\5\\4.jpg")
	g, err := gifx.Generate(images, []int{100, 100, 100, 100}, 0)
	if err != nil {
		logger.Panic(err)
	}
	watermark, err := img.OpenLocalFile("D:\\tmp\\4\\mark1.png")
	if err != nil {
		logger.Fatal(err)
	}
	watermark = watermark.Resize(50, 50, imaging.Lanczos)
	g.AddWaterMark(watermark, imaging.BottomRight, 10, 10, 1)

	of, err := file.CreateFile("D:\\tmp\\5\\merge.gif")
	if err != nil {
		logger.Panic(err)
	}
	defer of.Close()

	err = gif.EncodeAll(of, g.GetSource())
	if err != nil {
		logger.Panic(err)
	}
	logger.Info("merge success")
}
