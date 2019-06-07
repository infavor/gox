package img_test

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"log"
	"testing"
)

// 图片缩放示例
func TestResizeImage(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)

	// 图片按比例缩放，宽度固定，高度跟随
	dstImage128 := imaging.Resize(src, 128, 0, imaging.Lanczos)

	// Save the resulting image as JPEG.
	err = imaging.Save(dstImage128, "D:\\tmp\\123\\1_resize.jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}

// 图片裁剪示例
func TestCropImage(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	// Crop the original image to 300x300px size using the center anchor.
	src1 := imaging.CropAnchor(src, 100, 100, imaging.Center)
	src2 := imaging.CropAnchor(src, 100, 100, imaging.TopLeft)
	src3 := imaging.CropAnchor(src, 100, 100, imaging.TopRight)
	src4 := imaging.CropAnchor(src, 100, 100, imaging.Left)
	src5 := imaging.CropAnchor(src, 100, 100, imaging.Right)
	// Save the resulting image as JPEG.
	imaging.Save(src1, "D:\\tmp\\123\\1_crop_Center.jpg")
	imaging.Save(src2, "D:\\tmp\\123\\1_crop_TopLeft.jpg")
	imaging.Save(src3, "D:\\tmp\\123\\1_crop_TopRight.jpg")
	imaging.Save(src4, "D:\\tmp\\123\\1_crop_Left.jpg")
	imaging.Save(src5, "D:\\tmp\\123\\1_crop_Right.jpg")
}

// 使用高斯函数生成图像的模糊版本，Sigma参数必须为正，表示图像模糊的程度。
func TestBlurImage(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	img1 := imaging.Blur(src, 10)
	// Save the resulting image as JPEG.
	imaging.Save(img1, "D:\\tmp\\123\\1_blur.jpg")
}

// 缩小并高斯模糊
func TestResizeBlurImage(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	img1 := imaging.Resize(src, 128, 0, imaging.Lanczos)
	img1 = imaging.Blur(img1, 2)
	// Save the resulting image as JPEG.
	imaging.Save(img1, "D:\\tmp\\123\\1_resize_blur.jpg")
}

// 使图像变为黑白
func TestGray(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	img2 := imaging.Grayscale(src)
	// Save the resulting image as JPEG.
	imaging.Save(img2, "D:\\tmp\\123\\1_Grayscale.jpg")
}

// 调整图像对比度
func TestAdjustContrast(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	img2 := imaging.AdjustContrast(src, 50)
	// Save the resulting image as JPEG.
	imaging.Save(img2, "D:\\tmp\\123\\1_AdjustContrast.jpg")
}

// 锐化图像
func TestSharpen(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	img2 := imaging.Sharpen(src, 2)
	// Save the resulting image as JPEG.
	imaging.Save(img2, "D:\\tmp\\123\\1_Sharpen.jpg")
}

// 反转图像
func TestInvert(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	img3 := imaging.Invert(src) // Save the resulting image as JPEG.
	imaging.Save(img3, "D:\\tmp\\123\\1_Invert.jpg")
}

// 使用卷积滤镜创建图像的浮雕版本3x3。
func TestConvolve3x3(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	img4 := imaging.Convolve3x3(
		src,
		[9]float64{
			-4.3, -1, 5,
			-4.3, 1, 5,
			-4.3, -1, 5,
		},
		nil,
	)
	imaging.Save(img4, "D:\\tmp\\123\\1_Convolve3x3.jpg")
}

// 使用卷积滤镜创建图像的浮雕版本3x3。
func TestConvolve5x5(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	img4 := imaging.Convolve5x5(
		src,
		[25]float64{
			-2, -1, 0, 1, 1,
			-2, -1, 0, 1, 2,
			-2, -1, 1, 1, 2,
			-2, -1, 0, 1, 2,
			-1, -1, 0, 1, 2,
		},
		nil,
	)
	imaging.Save(img4, "D:\\tmp\\123\\1_Convolve5x5.jpg")
}

// 调整图像的亮度。
func TestAdjustBrightness(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	img4 := imaging.AdjustBrightness(src, 50)
	imaging.Save(img4, "D:\\tmp\\123\\1_Brightness.jpg")
}

// 调整图像的亮度。
func TestAdjustGamma(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	img4 := imaging.AdjustGamma(src, 0.3)
	imaging.Save(img4, "D:\\tmp\\123\\1_AdjustGamma.jpg")
}

// 调整图像的饱和度。
func TestAdjustSaturation(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	img4 := imaging.AdjustSaturation(src, 50)
	imaging.Save(img4, "D:\\tmp\\123\\1_AdjustSaturation.jpg")
}

// 调整图像的???。
func TestAdjustSigmoid(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	img4 := imaging.AdjustSigmoid(src, 0, -10)
	imaging.Save(img4, "D:\\tmp\\123\\1_AdjustSigmoid.jpg")
}

// 逆时针旋转图像。
func TestRotate(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	img4 := imaging.Rotate(src, 50, color.White)
	imaging.Save(img4, "D:\\tmp\\123\\1_Rotate.jpg")
}

// 逆时针旋转图像。
/*func Test2(t *testing.T) {
	// Open a test image.
	src1, _ := imaging.Open("D:\\tmp\\123\\1.jpg")
	src12, _ := imaging.Open("D:\\tmp\\123\\1.jpg")
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	img4 := imaging.Rotate(src, 50, color.White)
	imaging.Save(img4, "D:\\tmp\\123\\1_Rotate.jpg")
}*/

// 左右翻转图像
func TestTranspose(t *testing.T) {
	// Open a test image.
	src, err := imaging.Open("D:\\tmp\\123\\1.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println(src.Bounds().Size().X, ", ", src.Bounds().Size().Y)
	img4 := imaging.Transverse(src)
	img4 = imaging.Rotate90(img4)
	imaging.Save(img4, "D:\\tmp\\123\\1_Transpose.jpg")
}

// 覆盖子图像到背景图像(子图像不透明)
func TestPaste(t *testing.T) {
	// Open a test image.
	src1, _ := imaging.Open("D:\\tmp\\123\\1.jpg")
	src2, _ := imaging.Open("D:\\tmp\\123\\3.png")

	x := src1.Bounds().Size().X - src2.Bounds().Size().X
	y := src1.Bounds().Size().Y - src2.Bounds().Size().Y

	img3 := imaging.Paste(src1, src2, image.Pt(x, y))

	imaging.Save(img3, "D:\\tmp\\123\\1_Paste.jpg")
}

// 覆盖子图像到背景图像(子图像透明)(打水印)
func TestPaste1(t *testing.T) {
	// Open a test image.
	src1, _ := imaging.Open("D:\\tmp\\123\\1.jpg")
	src2, _ := imaging.Open("D:\\tmp\\123\\3.png")

	x := src1.Bounds().Size().X - src2.Bounds().Size().X - 20
	y := src1.Bounds().Size().Y - src2.Bounds().Size().Y - 20

	img3 := imaging.Overlay(src1, src2, image.Pt(x, y), 1)
	imaging.Save(img3, "D:\\tmp\\123\\1_Overlay.jpg")
}
