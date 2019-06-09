package gifx_test

import (
	"github.com/disintegration/imaging"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/img"
	"github.com/hetianyi/gox/img/gifx"
	"github.com/hetianyi/gox/logger"
	log "github.com/sirupsen/logrus"
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
		log.Fatal(err)
	}
	watermark, err := img.OpenLocalFile("D:\\tmp\\4\\mark1.png")
	if err != nil {
		log.Fatal(err)
	}
	watermark = watermark.Resize(50, 50, imaging.Lanczos)
	g.AddWaterMark(watermark, imaging.BottomRight, 10, 10, 1)

	of, err := file.CreateFile("D:\\tmp\\4\\origin_watermark.gif")
	if err != nil {
		log.Panic(err)
	}
	defer of.Close()

	err = gif.EncodeAll(of, g.GetSource())
	if err != nil {
		log.Panic(err)
	}
}
