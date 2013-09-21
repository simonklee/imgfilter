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
	"github.com/simonz05/imgfilter/util"
	"github.com/gographics/imagick/imagick"
)

func Resize(data []byte, width, height int) ([]byte, error) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.ReadImageBlob(data)

	if err != nil {
		return nil, err
	}

	// Get original image size
	oWidth := mw.GetImageWidth()
	oHeight := mw.GetImageHeight()

	// Calculate new size
	hWidth := uint(util.IntMax(width, int(oWidth)))
	hHeight := uint(util.IntMax(height, int(oHeight)))

	// Resize the image using the Lanczos filter
	// The blur factor is a float, where > 1 is blurry, < 1 is sharp
	err = mw.ResizeImage(hWidth, hHeight, imagick.FILTER_LANCZOS, 1)

	if err != nil {
		return nil, err
	}

	// Set the compression quality to 95 (high quality = low compression)
	err = mw.SetImageCompressionQuality(95)

	if err != nil {
		return nil, err
	}

	return mw.GetImageBlob(), nil
}

func Crop(data []byte) ([]byte, error) {
	//err = mw.CropImage(200, 200, 0, 0)
	return nil, nil
}

func Thumbnail(data []byte) ([]byte, error) {
	return nil, nil
}
