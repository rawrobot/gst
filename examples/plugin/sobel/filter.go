package main

import (
	"image"
	//"image/color"
	"math"
)

var (
	sobelXL = [9]int{
		-1, 0, 1,
		-2, 0, 2,
		-1, 0, 1,
	}

	sobelYL = [9]int{
		-1, -2, -1,
		0, 0, 0,
		1, 2, 1,
	}
)

const kernelSize = 3

func FilterGrayIP(grayImg *image.Gray) {
	max := grayImg.Bounds().Max
	min := grayImg.Bounds().Min

	//Copy source to work with, because we will modify incoming image
	imgCopy := CopySubImage(grayImg, image.Rect(max.X, max.Y, min.X, min.Y))
	/*
		 Filtered image must be two pixels shorter, because
		there must be a row of pixels on each side of a pixel for the sobel operator
		to work
	*/
	width := max.X - 1 //to provide a "border" of 1 pixel
	height := max.Y - 1

	var sv uint
	var v float64

	for x := 1; x < width; x++ {
		for y := 1; y < height; y++ {
			sv = applySobelFilter(imgCopy, x, y)
			//math.Sqrt works 30 times faster that
			//https://www.geeksforgeeks.org/square-root-of-an-integer/
			v = math.Sqrt(float64(sv))
			//clip value
			if v > 255.0 {
				v = 255.0
			} else if v < 0.0 {
				v = 0.0
			}
			//grayImg.SetGray(x, y, color.Gray{Y: uint8(v)})
			grayImg.Pix[grayImg.PixOffset(x-1, y-1)] = uint8(v)
		}
	}
}

func applySobelFilter(img *image.Gray, x int, y int) uint {
	var fX, fY, pixel, index int
	curX := x - 1
	curY := y - 1
	for i := 0; i < kernelSize; i++ {
		for j := 0; j < kernelSize; j++ {
			//it is unsafe but faster on 10% or so
			pixel = int(img.Pix[img.PixOffset(curX, curY)])
			fX += sobelXL[index] * pixel
			fY += sobelYL[index] * pixel
			curX++
			index++
		}
		curX = x - 1
		curY++
	}
	return uint(fX*fX + fY*fY)
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func CopySubImage(p *image.Gray, r image.Rectangle) *image.Gray {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &image.Gray{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	img := &image.Gray{
		Pix:    make([]byte, r.Max.X*r.Max.Y),
		Stride: p.Stride,
		Rect:   r,
	}
	copy(img.Pix, p.Pix[i:])
	return img
}
