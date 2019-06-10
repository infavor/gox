package img_test

import (
	"github.com/disintegration/imaging"
	"github.com/hetianyi/gox/img"
	"image"
	"image/color"
	"log"
	"testing"
)

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
	im, err := img.OpenLocalFile("E:\\test\\2.jpg") // 1900x1283
	if err != nil {
		log.Panic(err)
	}
	im.Convolve3x3(img.Default3x3Kernel)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Convolve3x3.jpg")
}

func TestImage_Convolve5x5(t *testing.T) {
	im, err := img.OpenLocalFile("E:\\test\\2.jpg") // 1900x1283
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
	im, err := img.OpenLocalFile("E:\\test\\2.jpg") // 1900x1283
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
	im, _ := img.OpenLocalFile("E:\\test\\2.jpg") // 1900x1283
	im.Compress(10)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Compress.jpg")
}

func TestImage_Fit(t *testing.T) {
	im, _ := img.OpenLocalFile("E:\\test\\2.jpg") // 1900x1283
	im.Fit(500, 500, imaging.Lanczos)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Fit.jpg")
}

func TestImage_Fill(t *testing.T) {
	im, _ := img.OpenLocalFile("E:\\test\\2.jpg") // 1900x1283
	im.Fill(500, 500, imaging.Center, imaging.Lanczos)
	imaging.Save(im.GetSource(), "E:\\test\\TestImage_Fill.jpg")
}
