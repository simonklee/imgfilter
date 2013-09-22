package image

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gographics/imagick/imagick"
)

type ResizeCase struct {
	filename string
	w, h     uint
}

type CropCase struct {
	filename string
	w, h     uint
	x, y     int
}

//func TestResize(t *testing.T) {
//	imagick.Initialize()
//	defer imagick.Terminate()
//
//	tests := []*ResizeCase{
//		//{"fixture/gopher.png", 200, 200},
//		{"fixture/gopher-1.jpg", 100, 100},
//	}
//
//	for _, x := range tests {
//		before, err := ioutil.ReadFile(x.filename)
//
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		after, err := Resize(before, x.w, x.h)
//
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		mw := imagick.NewMagickWand()
//		err = mw.ReadImageBlob(after)
//
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		mw.DisplayImage(os.Getenv("DYSPLAY"))
//		mw.Destroy()
//	}
//}

func TestCrop(t *testing.T) {
	imagick.Initialize()
	defer imagick.Terminate()

	tests := []*CropCase{
		//{"fixture/gopher.png", 200, 200},
		{"fixture/gopher-1.jpg", 30, 150, 0, 0},
	}

	if err := os.Mkdir("test-out", os.ModeDir|os.ModePerm); os.IsNotExist(err) {
		t.Fatal(err)
	}

	for i, x := range tests {
		before, err := ioutil.ReadFile(x.filename)

		if err != nil {
			t.Fatal(err)
		}

		after, err := Crop(before, x.w, x.h, x.x, x.y)

		if err != nil {
			t.Fatal(err)
		}

		filename := fmt.Sprintf("test-out/crop-%dx%d+%d+%d-%d%s", x.w, x.h, x.x, x.y, i, filepath.Ext(x.filename))
		err = ioutil.WriteFile(filename, after, os.ModePerm)

		if err != nil {
			t.Fatal(err)
		}
	}
}
