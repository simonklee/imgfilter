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
	"github.com/simonz05/imgfilter/util"
)

type image struct {
	mw        *imagick.MagickWand
	w, h      uint
	direction string
}

func newImageFromBlob(blob []byte) (*image, error) {
	im := new(image)
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
func (im *image) normalizeSize(width, height uint) (w, h uint) {
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
func (im *image) normalizeOffset(w, h uint, xOffset, yOffset int) (x int, y int) {
	x, y = im.gravity(w, h)
	x += xOffset
	y += yOffset
	x -= util.IntMax(0, x+int(w)-int(im.w))
	y -= util.IntMax(0, y+int(h)-int(im.h))
	return util.IntMax(0, x), util.IntMax(0, y)
}

// Calculate x and y offset based on gravity. ImageMagick's SetImageGravity
// function doesn't seem to work.
func (im *image) gravity(w, h uint) (x, y int) {
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

func (im *image) destroy() {
	im.mw.Destroy()
	imagick.Terminate()
}

func (im *image) resize(width, height uint) error {
	w, h := im.normalizeSize(width, height)

	if err := im.mw.ResizeImage(w, h, imagick.FILTER_LANCZOS, 1); err != nil {
		return err
	}

	return nil
}

func (im *image) crop(width, height uint, x, y int) (err error) {
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

// Set the compression quality (high quality = low compression)
func (im *image) compress(level uint) error {
	return im.mw.SetImageCompressionQuality(level)
}

// Resize formats an image according to width and height. It does not preserve
// the aspect ratio of the image.
func Resize(data []byte, width, height uint) ([]byte, error) {
	im, err := newImageFromBlob(data)
	defer im.destroy()

	if err != nil {
		return nil, err
	}

	if err = im.resize(width, height); err != nil {
		return nil, err
	}

	return im.mw.GetImageBlob(), nil
}

// Crop formats an image according to width and height.
func Crop(data []byte, width, height uint, x, y int, direction string) ([]byte, error) {
	im, err := newImageFromBlob(data)
	defer im.destroy()

	if err != nil {
		return nil, err
	}

	im.direction = direction

	if err = im.crop(width, height, x, y); err != nil {
		return nil, err
	}

	return im.mw.GetImageBlob(), nil
}
