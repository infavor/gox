package img_test

import (
	"bytes"
	"fmt"
	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
	"github.com/hetianyi/gox/file"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"runtime"
	"testing"
)

// GIF打水印
func TestGIFWaterMark4(t *testing.T) {
	// Open a test image.
	inputFile, _ := file.GetFile("D:\\tmp\\4\\origin.gif")
	im, err := gif.DecodeAll(inputFile)
	if err != nil {
		panic(err)
	}
	inputFile.Close()
	filter := gift.New(
		// high-quality resampling with pixel mixing
		gift.ResizeToFit(200, 200, gift.LanczosResampling),
		gift.Rotate(45, color.White, gift.LinearInterpolation),
	)

	src2, _ := imaging.Open("D:\\tmp\\4\\Office365LogoWLockup.scale-140.png")
	src2 = imaging.Resize(src2, 200, 50, imaging.Lanczos)

	firstFrame := im.Image[0]
	tmp := image.NewNRGBA(firstFrame.Bounds())

	for i := range im.Image {
		/*x := frame.Bounds().Size().X - src2.Bounds().Size().X - 20
		y := frame.Bounds().Size().Y - src2.Bounds().Size().Y - 20

		img3 := imaging.Overlay(frame, src2, image.Pt(x, y), 1)
		buf := &bytes.Buffer{}
		if err := gif.Encode(buf, img3, nil); err != nil {
			fmt.Println(err)
		}
		tmpimg, err := gif.Decode(buf)
		if err != nil {
			fmt.Println(err)
		}*/
		// draw current frame over previous:
		gift.New().DrawAt(tmp, im.Image[i], im.Image[i].Bounds().Min, gift.OverOperator)
		dst := image.NewPaletted(filter.Bounds(tmp.Bounds()), im.Image[i].Palette)
		filter.Draw(dst, tmp)
		im.Image[i] = dst
	}

	outputFile1, _ := os.Create("D:\\tmp\\4\\2_resize_new.gif")
	if err != nil {
		panic(err)
	}
	defer outputFile1.Close()

	fmt.Println(gif.EncodeAll(outputFile1, im))
}

// GIF打水印
func TestGIFWaterMark5(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	f, err := os.Open("D:\\tmp\\4\\origin.gif")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		panic(err)
	}

	imgWidth, imgHeight := getGifDimensions(g)
	overpaintImage := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(overpaintImage, overpaintImage.Bounds(), g.Image[0], image.ZP, draw.Src)

	src2, _ := imaging.Open("D:\\tmp\\4\\Office365LogoWLockup.scale-140.png")
	src2 = imaging.Resize(src2, 200, 50, imaging.Lanczos)

	for i, frame := range g.Image {

		x := frame.Bounds().Size().X - src2.Bounds().Size().X - 20
		y := frame.Bounds().Size().Y - src2.Bounds().Size().Y - 20

		img3 := imaging.Overlay(frame, src2, image.Pt(x, y), 1)
		draw.Draw(overpaintImage, overpaintImage.Bounds(), frame, image.ZP, draw.Over)

		//ut, _ := file.CreateFile("D:\\tmp\\4\\slice-" + convert.IntToStr(i) + ".bmp")
		//gif.Encode(out, img3, nil)

		buf := &bytes.Buffer{}
		if err := gif.Encode(buf, img3, nil); err != nil {
			fmt.Println(err)
		}
		tmpimg, err := gif.Decode(buf)
		if err != nil {
			fmt.Println(err)
		}
		g.Image[i] = tmpimg.(*image.Paletted)
	}

	of, err := os.Create("D:\\tmp\\4\\2_resize_new.gif")
	if err != nil {
		panic(err)
	}
	defer of.Close()

	err = gif.EncodeAll(of, g)
	if err != nil {
		panic(err)
	}
}

// ref: https://stackoverflow.com/questions/33295023/how-to-split-gif-into-images
func getGifDimensions(gif *gif.GIF) (x, y int) {
	var lowestX int
	var lowestY int
	var highestX int
	var highestY int
	for _, img := range gif.Image {
		if img.Rect.Min.X < lowestX {
			lowestX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowestY {
			lowestY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highestX {
			highestX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highestY {
			highestY = img.Rect.Max.Y
		}
	}
	return highestX - lowestX, highestY - lowestY
}
