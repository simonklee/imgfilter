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
	"errors"
	"fmt"

	"github.com/gographics/imagick/imagick"
)

var maxOutputSize uint = 2000

type image struct {
	mw *imagick.MagickWand
}

func newImageFromBlob(blob []byte) (*image, error) {
	im := new(image)
	imagick.Initialize()
	im.mw = imagick.NewMagickWand()
	return im, im.mw.ReadImageBlob(blob)
}

func (im *image) normalizeSize(width, height uint) (w, h uint) {
	if width > height {
		h = maxOutputSize * height / width
		w = maxOutputSize
	} else {
		w = maxOutputSize * width / height
		h = maxOutputSize
	}

	if w >= width && h >= height {
		w = width
		h = height
	}

	//// Get original image size
	//oWidth := mw.GetImageWidth()
	//oHeight := mw.GetImageHeight()

	//// Calculate new size
	//hWidth := util.UintMin(width, oWidth)
	//hHeight := util.UintMin(height, oHeight)

	return w, h
}

func (im *image) destroy() {
	im.mw.Destroy()
	imagick.Terminate()
}

func (im *image) resize(width, height uint) error {
	if err := im.mw.SetGravity(imagick.GRAVITY_CENTER); err != nil {
		return err
	}

	if err := im.mw.ResizeImage(width, height, imagick.FILTER_LANCZOS, 1); err != nil {
		return err
	}

	return nil
}

func (im *image) crop(width, height uint, x, y int) (err error) {
	if err = im.mw.CropImage(width, height, x, y); err != nil {
		return
	}

	//if err = im.mw.ResetImagePage("0x0"); err != nil {
	//	return
	//}

	return
}

func (im *image) transform(crop, resize string) error {
	if mw1 := im.mw.TransformImage(crop, resize); mw1 == nil {
		return errors.New("Image transformation failed")
	}

	return nil
}

// Set the compression quality (high quality = low compression)
func (im *image) compress(level uint) error {
	return im.mw.SetImageCompressionQuality(level)
}

func Resize(data []byte, width, height uint) ([]byte, error) {
	im, err := newImageFromBlob(data)
	defer im.destroy()

	if err != nil {
		return nil, err
	}

	w, h := im.normalizeSize(width, height)

	if err = im.resize(w, h); err != nil {
		return nil, err
	}

	return im.mw.GetImageBlob(), nil
}

func Crop(data []byte, width, height uint, x, y int) ([]byte, error) {
	im, err := newImageFromBlob(data)
	defer im.destroy()

	if err != nil {
		return nil, err
	}

	w, h := im.normalizeSize(width, height)

	if err = im.crop(w, h, x, y); err != nil {
		return nil, err
	}

	return im.mw.GetImageBlob(), nil
}

func Transform(data []byte, resizeWidth, resizeHeight, cropWidth, cropHeight uint, x, y int) ([]byte, error) {
	im, err := newImageFromBlob(data)
	defer im.destroy()

	if err != nil {
		return nil, err
	}

	w, h := im.normalizeSize(resizeWidth, resizeHeight)

	crop := fmt.Sprintf("%dx%d+%d+%d", cropWidth, cropHeight, x, y)
	resize := fmt.Sprintf("%dx%d>", w, h)

	if err = im.transform(crop, resize); err != nil {
		return nil, err
	}

	return im.mw.GetImageBlob(), nil
}
