package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"

	"github.com/nfnt/resize"
)

type resizerOwnImpl struct {
	imageBytes []byte
	width      int
	height     int
}

type resizerLibImpl struct {
	imageBytes []byte
	width      int
	height     int
}

// Resizer interace for resizing images
type resizer interface {
	resizeImage() []byte
}

// ResizeImage image
func ResizeImage(imageBytes []byte, width int, height int) []byte {
	// resizer := resizerOwnImpl{imageBytes: imageBytes, width: width, height: height}
	resizer := resizerLibImpl{imageBytes: imageBytes, width: width, height: height}
	return resizer.resizeImage()
}

func (r resizerLibImpl) resizeImage() []byte {
	img, _, err := image.Decode(bytes.NewReader(r.imageBytes))
	if err != nil {
		log.Fatal(err)
	}
	m := resize.Thumbnail(uint(r.width), uint(r.height), img, resize.Lanczos3)
	return imgToBytes(m)
}

func (r resizerOwnImpl) resizeImage() []byte {
	img, _, err := image.Decode(bytes.NewReader(r.imageBytes))
	if err != nil {
		log.Fatal(err)
	}
	resImg := resizeOwnImpl(img, r.width, r.height)
	return imgToBytes(resImg)
}

func resizeOwnImpl(img image.Image, width int, height int) image.Image {
	//truncate pixel size
	minX := img.Bounds().Min.X
	minY := img.Bounds().Min.Y
	maxX := img.Bounds().Max.X
	maxY := img.Bounds().Max.Y
	for (maxX-minX)%height != 0 {
		maxX--
	}
	for (maxY-minY)%width != 0 {
		maxY--
	}
	scaleX := (maxX - minX) / height
	scaleY := (maxY - minY) / width

	imgRect := image.Rect(0, 0, height, width)
	resImg := image.NewRGBA(imgRect)
	draw.Draw(resImg, resImg.Bounds(), &image.Uniform{C: color.White}, image.ZP, draw.Src)
	for y := 0; y < width; y++ {
		for x := 0; x < height; x++ {
			averageColor := getAverageColor(img, minX+x*scaleX, minX+(x+1)*scaleX, minY+y*scaleY, minY+(y+1)*scaleY)
			resImg.Set(x, y, averageColor)
		}
	}
	return resImg
}

func getAverageColor(img image.Image, minX int, maxX int, minY int, maxY int) color.Color {
	var averageRed float64
	var averageGreen float64
	var averageBlue float64
	var averageAlpha float64
	scale := 1.0 / float64((maxX-minX)*(maxY-minY))

	for i := minX; i < maxX; i++ {
		for k := minY; k < maxY; k++ {
			r, g, b, a := img.At(i, k).RGBA()
			averageRed += float64(r) * scale
			averageGreen += float64(g) * scale
			averageBlue += float64(b) * scale
			averageAlpha += float64(a) * scale
		}
	}

	averageRed = math.Sqrt(averageRed)
	averageGreen = math.Sqrt(averageGreen)
	averageBlue = math.Sqrt(averageBlue)
	averageAlpha = math.Sqrt(averageAlpha)

	averageColor := color.RGBA{
		R: uint8(averageRed),
		G: uint8(averageGreen),
		B: uint8(averageBlue),
		A: uint8(averageAlpha)}

	return averageColor
}

func imgToBytes(img image.Image) []byte {
	var opt jpeg.Options
	opt.Quality = 80

	buff := bytes.NewBuffer(nil)
	err := jpeg.Encode(buff, img, &opt)
	if err != nil {
		log.Fatal(err)
	}

	return buff.Bytes()
}
