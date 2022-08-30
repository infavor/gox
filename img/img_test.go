package img_test

import (
	"github.com/disintegration/imaging"
	"github.com/infavor/gox/file"
	"github.com/infavor/gox/fontx"
	"github.com/infavor/gox/img"
	"github.com/infavor/gox/logger"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	logger.Init(nil)
}

func TestImage_Resize(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}
	im.Resize(500, 0, imaging.Lanczos)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Resize.jpg")
}

func TestImage_Crop(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}
	imNew := im.Clone()
	im.Crop(500, 200, imaging.Center)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Crop.jpg")
	imaging.Save(imNew.GetSource(), "E:\\test\\TestImage_Crop_clone.jpg")
}

func TestImage_Blur(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}
	im.Blur(16)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Blur.jpg")
}

func TestImage_Gray(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}
	im.Gray()
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Gray.jpg")
}

func TestImage_AdjustContrast(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}
	imNew := im.Clone()
	im.AdjustContrast(100)
	imNew.AdjustContrast(-50)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_AdjustContrast_100.jpg")
	imaging.Save(imNew.GetSource(), "E:\\test\\TestImage_AdjustContrast_-100.jpg")
}

func TestImage_Sharpen(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}
	im.Sharpen(100)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Sharpen.jpg")
}

func TestImage_Invert(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}
	im.Invert()
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Invert.jpg")
}

func TestImage_Convolve3x3(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}
	im.Convolve3x3(img.Default3x3Kernel)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Convolve3x3.jpg")
}

func TestImage_Convolve5x5(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}
	im.Convolve5x5(img.Default5x5Kernel)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Convolve5x5.jpg")
}

func TestImage_AdjustBrightness(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}
	imNew := im.Clone()
	im.AdjustBrightness(50)
	imNew.AdjustBrightness(-50)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_AdjustBrightness_50.jpg")
	imaging.Save(imNew.GetSource(), "E:\\test\\TestImage_AdjustBrightness-50.jpg")
}

func TestImage_AdjustGamma(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}
	imNew := im.Clone()
	im.AdjustGamma(50)
	imNew.AdjustBrightness(-50)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_AdjustGamma_50.jpg")
	imaging.Save(imNew.GetSource(), "E:\\test\\TestImage_AdjustGamma-50.jpg")
}

func TestImage_AdjustSaturation(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}
	imNew := im.Clone()
	im.AdjustSaturation(50)
	imNew.AdjustBrightness(-50)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_AdjustSaturation_50.jpg")
	imaging.Save(imNew.GetSource(), "E:\\test\\TestImage_AdjustSaturation-50.jpg")
}

func TestImage_Rotate(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}

	imNew1 := im.Clone()
	imNew2 := im.Clone()
	imNew3 := im.Clone()
	imNew4 := im.Clone()

	im.Rotate(45, color.White)
	imNew1.Rotate(90, color.White)
	imNew2.Rotate(180, color.White)
	imNew3.Rotate(270, color.White)
	imNew4.Rotate(360, color.White)

	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Rotate-45.jpg")
	imaging.Save(imNew1.GetSource(), "E:\\test\\TestImage_Rotate-90.jpg")
	imaging.Save(imNew2.GetSource(), "E:\\test\\TestImage_Rotate-180.jpg")
	imaging.Save(imNew3.GetSource(), "E:\\test\\TestImage_Rotate-270.jpg")
	imaging.Save(imNew4.GetSource(), "E:\\test\\TestImage_Rotate-360.jpg")
}

func TestImage_Transverse(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}
	im.Transverse()
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Transverse.jpg")
}

func TestPaste(t *testing.T) {
	im, _ := img.OpenLocalFile("E:\\test\\1.jpg")          // 1900x1283
	im1, _ := img.OpenLocalFile("E:\\test\\watermark.png") // 1900x1283
	im.Paste(im1, image.Pt(1000, 200))
	imaging.Save(im.GetSource(), "E:\\test\\TestPaste.jpg")
}

func TestOverlay(t *testing.T) {
	im, _ := img.OpenLocalFile("E:\\test\\1.jpg")          // 1900x1283
	im1, _ := img.OpenLocalFile("E:\\test\\watermark.png") // 1900x1283
	im.Overlay(im1, image.Pt(1000, 200), 1)
	imaging.Save(im.GetSource(), "E:\\test\\TestOverlay.jpg")
}

func TestImage_AddWaterMark(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}
	watermark, err := img.OpenLocalFile("E:\\test\\watermark.png")
	if err != nil {
		log.Fatal(err)
	}
	im.AddWaterMark(watermark, imaging.BottomRight, 20, 20, 0.5)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_AddWaterMark.jpg")
}

func TestImage_Compress(t *testing.T) {
	im, _ := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	im.Compress(10)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Compress.jpg")
}

func TestImage_Fit(t *testing.T) {
	im, _ := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	im.Fit(500, 500, imaging.Lanczos)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Fit.jpg")
}

func TestImage_Fill(t *testing.T) {
	im, _ := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	im.Fill(500, 500, imaging.Center, imaging.Lanczos)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Fill.jpg")
}

func TestSaveToFile(t *testing.T) {
	im, _ := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	im.Blur(8)
	img.SaveToFile(im, "E:\\test\\TestSaveToFile.jpg")
}

func TestSave(t *testing.T) {
	im, _ := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	im.Blur(8)

	out, _ := file.CreateFile("E:\\test\\TestSave.jpg")

	img.Save(im, out, imaging.JPEG)
}

/*
func TestText(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 300, 100))
	addLabel(img, 20, 30, "Hello Go")

	f, err := os.Create("E:\\test\\TestText.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{200, 100, 0, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}*/

func TestImage_DrawText(t *testing.T) {
	im, _ := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
	fo, _ := fontx.LoadFont("E:\\test\\Inkfree.ttf")
	fc := &fontx.FontConfig{
		Font:     fo.Font,
		FontSize: 200,
		Color:    color.Black,
	}
	im.Blur(8)
	metrics := fo.GetMetrics(fc)
	for {
		im.DrawText("Hello", fc, metrics, imaging.BottomRight, 500, 700)
		im.DrawText("From", fc, metrics, imaging.BottomRight, 500, 500)
		im.DrawText("Other Side", fc, metrics, imaging.BottomRight, 1100, 300)
	}
	out, _ := file.CreateFile("E:\\test\\TestImage_DrawText.jpg")

	img.Save(im, out, imaging.JPEG)
}

func TestCompressDir(t *testing.T) {

	filepath.Walk("C:\\spider\\test\\download", func(path string, info os.FileInfo, err error) error {

		if info.IsDir() && info.Name() == "download" {
			return nil
		}
		if info.IsDir() {
			return file.CreateDirs("D:\\out\\" + info.Name())
		}

		if !info.IsDir() {

			dirName := filepath.Base(filepath.Dir(path))

			im, err := img.OpenLocalFile(path) // 1900x1283
			if err != nil {
				logger.Error(err)
				return nil
			}
			x := im.GetSource().Bounds().Size().X
			y := im.GetSource().Bounds().Size().Y

			if x > y && x > 2000 {
				im.Resize(2000, 0, imaging.Lanczos)
				im.Compress(50)
				imaging.Save(im.GetSource(), "D:\\out\\"+dirName+"\\"+info.Name())
			} else {
				im.Resize(0, 2000, imaging.Lanczos)
				im.Compress(50)
				imaging.Save(im.GetSource(), "D:\\out\\"+dirName+"\\"+info.Name())
			}
		}
		return nil
	})
}
