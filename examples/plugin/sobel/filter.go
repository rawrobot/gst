package main

import (
	"image"
	"image/color"
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

	/* filtered image must be two pixels shorter, because
	there must be a row of pixels on each side of a pixel for the sobel operator
	to work*/
	//Copy source to work with, because we will modify incomeng image
	imgCopy := CopySubImage(grayImg, image.Rect(max.X, max.Y, min.X, min.Y))
	width := max.X - 1 //to provide a "border" of 1 pixel
	height := max.Y - 1

	var v uint32

	for x := 1; x < width; x++ {
		for y := 1; y < height; y++ {
			fX, fY := applySobelFilter(imgCopy, x, y)
			v = FloorSqrt((fX*fX)+(fY*fY)) + 1 // +1 to make it ceil
			grayImg.SetGray(x, y, color.Gray{Y: uint8(v)})
		}
	}
}

func FilterGrayFast(grayImg *image.Gray) (filtered *image.Gray) {
	max := grayImg.Bounds().Max
	min := grayImg.Bounds().Min

	/* filtered image must be two pixels shorter, because
	there must be a row of pixels on each side of a pixel for the sobel operator
	to work*/
	filtered = image.NewGray(image.Rect(max.X-2, max.Y-2, min.X, min.Y))
	width := max.X - 1 //to provide a "border" of 1 pixel
	height := max.Y - 1

	var v uint32

	for x := 1; x < width; x++ {
		for y := 1; y < height; y++ {
			fX, fY := applySobelFilter(grayImg, x, y)
			v = FloorSqrt((fX*fX)+(fY*fY)) + 1 // +1 to make it ceil
			filtered.SetGray(x, y, color.Gray{Y: uint8(v)})
		}
	}

	return filtered
}

func applySobelFilter(img *image.Gray, x int, y int) (uint32, uint32) {
	var fX, fY, pixel, index int
	curX := x - 1
	curY := y - 1
	for i := 0; i < kernelSize; i++ {
		//index = i * kernelSize
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
	return Abs(fX), Abs(fY)
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

//unint32 math

// Abs returns the absolute value of the given int.
func Abs(x int) uint32 {
	if x < 0 {
		return uint32(-x)
	} else {
		return uint32(x)
	}
}

//https://www.geeksforgeeks.org/square-root-of-an-integer/
//Time Complexity: O(Log x)
//Note: The Binary Search can be further optimized to start with ‘start’ = 0 and ‘end’ = x/2.
//Floor of square root of x cannot be more than x/2 when x > 1.
func FloorSqrt(x uint32) (ans uint32) {
	// Base Cases
	if x == 0 || x == 1 {
		return x
	}

	// Do Binary Search for floor(sqrt(x))
	var (
		start uint32 = 0
		mid   uint32
	)
	end := x / 2

	for start <= end {
		mid = (start + end) / 2

		// If x is a perfect square
		if mid*mid == x {
			return mid
		}

		// Since we need floor, we update answer when mid*mid is
		// smaller than x, and move closer to sqrt(x)
		if mid*mid < x {

			start = mid + 1
			ans = mid
		} else { // If mid*mid is greater than x
			end = mid - 1
		}
	}
	return ans
}
