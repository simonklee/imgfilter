// Copyright (c) 2013 Simon Zimmermann
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package image provides an interface to imagemagick c library
package image

import (
	"github.com/gographics/imagick/imagick"
	"github.com/simonz05/util/math"
)

type Image struct {
	mw        *imagick.MagickWand
	w, h      uint
	nW, nH    uint
	direction string
}

// Create a new image from raw image source.
//
// Example
//
//		im, err := NewImageFromBlob(data)
// 		defer im.Destroy()
//
// 		if err != nil {
// 			return nil, err
// 		}
func NewImageFromBlob(blob []byte) (*Image, error) {
	im := new(Image)
	imagick.Initialize()

	im.mw = imagick.NewMagickWand()
	err := im.mw.ReadImageBlob(blob)

	if err != nil {
		return im, err
	}

	im.w = im.mw.GetImageWidth()
	im.h = im.mw.GetImageHeight()
	return im, nil
}

// Ensure that width and height does contain image.
func (im *Image) normalizeSize(width, height uint) (w, h uint) {
	if width > height {
		h = im.h * height / width
		w = im.w
	} else {
		w = im.w * width / height
		h = im.h
	}

	if w >= width && h >= height {
		w = width
		h = height
	}

	return
}

// Ensure that x and y offsets does contain image.
func (im *Image) normalizeOffset(w, h uint, xOffset, yOffset int) (x int, y int) {
	x, y = im.gravity(w, h)
	x += xOffset
	y += yOffset
	x -= math.IntMax(0, x+int(w)-int(im.w))
	y -= math.IntMax(0, y+int(h)-int(im.h))
	return math.IntMax(0, x), math.IntMax(0, y)
}

// Calculate x and y offset based on gravity. ImageMagick's SetImageGravity
// function doesn't seem to work.
func (im *Image) gravity(w, h uint) (x, y int) {
	switch im.direction {
	case "northwest":
		break
	case "north":
		x = int((im.w / 2) - (w / 2))
	case "northeast":
		x = int(im.w - w)
	case "west":
		y = int((im.h / 2) - (h / 2))
	case "east":
		x = int(im.w - w)
		y = int((im.h / 2) - (h / 2))
	case "southwest":
		y = int(im.h - h)
	case "south":
		x = int((im.w / 2) - (w / 2))
		y = int(im.h - h)
	case "southeast":
		x = int(im.w - w)
		y = int(im.h - h)
	case "center":
		x = int((im.w / 2) - (w / 2))
		y = int((im.h / 2) - (h / 2))
	default:
		x = int((im.w / 2) - (w / 2))
		y = int((im.h / 2) - (h / 2))
	}
	return
}

// cropSize computes a crop area for making thumbnail images.
func (im *Image) cropSize(width, height uint) (uint, uint) {
	fWidth := float64(width)
	fHeight := float64(height)
	fx := float64(im.w) / fWidth
	fy := float64(im.h) / fHeight

	var ratio float64

	if fx < fy {
		ratio = fx
	} else {
		ratio = fy
	}

	return uint(fWidth * ratio), uint(fHeight * ratio)
}

// Free image resource. Always call this.
func (im *Image) Destroy() {
	im.mw.Destroy()
	imagick.Terminate()
}

// Return width
func (im *Image) Width() uint {
	return im.w
}

// Return height
func (im *Image) Height() uint {
	return im.h
}

// SetDirection sets the crop gravity direction
func (im *Image) SetDirection(direction string) {
	im.direction = direction
}

// Resize formats an image according to width and height.
func (im *Image) Resize(width, height uint) error {
	w, h := im.normalizeSize(width, height)

	if err := im.mw.ResizeImage(w, h, imagick.FILTER_LANCZOS, 1); err != nil {
		return err
	}

	return nil
}

// Crop formats an image according to width and height.
func (im *Image) Crop(width, height uint, x, y int) (err error) {
	w, h := im.normalizeSize(width, height)
	x, y = im.normalizeOffset(w, h, x, y)

	if err = im.mw.CropImage(w, h, x, y); err != nil {
		return
	}

	//if err = im.mw.ResetImagePage("0x0"); err != nil {
	//	return
	//}

	return
}

// Thumbnail fits an image to a given size. It first calls crop, then resize.
func (im *Image) Thumbnail(width, height uint, x, y int) (err error) {
	w, h := im.normalizeSize(width, height)
	cw, ch := im.cropSize(w, h)
	x, y = im.normalizeOffset(cw, ch, x, y)

	if err = im.mw.CropImage(cw, ch, x, y); err != nil {
		return
	}

	if err = im.mw.ResizeImage(w, h, imagick.FILTER_LANCZOS, 1); err != nil {
		return
	}

	return
}

// Set the compression quality (high quality = low compression)
func (im *Image) Compress(level uint) error {
	return im.mw.SetImageCompressionQuality(level)
}

// Resize formats an image according to width and height. It does not preserve
// the aspect ratio of the image.
func Resize(data []byte, width, height uint) ([]byte, error) {
	im, err := NewImageFromBlob(data)
	defer im.Destroy()

	if err != nil {
		return nil, err
	}

	if err = im.Resize(width, height); err != nil {
		return nil, err
	}

	return im.mw.GetImageBlob(), nil
}

// Crop formats an image according to width and height.
func Crop(data []byte, width, height uint, x, y int, direction string) ([]byte, error) {
	im, err := NewImageFromBlob(data)
	defer im.Destroy()

	if err != nil {
		return nil, err
	}

	im.direction = direction

	if err = im.Crop(width, height, x, y); err != nil {
		return nil, err
	}

	return im.mw.GetImageBlob(), nil
}

// Thumbnail fits an image to a given size. It first calls Crop, then Resize.
func Thumbnail(data []byte, width, height uint, direction string) ([]byte, error) {
	im, err := NewImageFromBlob(data)
	defer im.Destroy()

	if err != nil {
		return nil, err
	}

	im.SetDirection(direction)

	if err = im.Thumbnail(width, height, 0, 0); err != nil {
		return nil, err
	}

	return im.mw.GetImageBlob(), nil
}
