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
	filename  string
	w, h      uint
	x, y      int
	direction string
}

func TestResize(t *testing.T) {
	imagick.Initialize()
	defer imagick.Terminate()

	tests := []*ResizeCase{
		{"fixture/circle.png", 30, 150},
	}

	if err := os.Mkdir("test-out", os.ModeDir|os.ModePerm); os.IsNotExist(err) {
		t.Fatal(err)
	}

	for i, x := range tests {
		before, err := ioutil.ReadFile(x.filename)

		if err != nil {
			t.Fatal(err)
		}

		after, err := Resize(before, x.w, x.h)

		if err != nil {
			t.Fatal(err)
		}

		filename := fmt.Sprintf("test-out/resize-%dx%d-%d%s", x.w, x.h, i, filepath.Ext(x.filename))
		err = ioutil.WriteFile(filename, after, os.ModePerm)

		if err != nil {
			t.Fatal(err)
		}
	}
}
func TestCrop(t *testing.T) {
	imagick.Initialize()
	defer imagick.Terminate()

	tests := []*CropCase{
		//{"fixture/gopher.png", 200, 200},
		{"fixture/circle.png", 300, 300, 40, 0, ""},
		{"fixture/circle.png", 30, 150, 0, 20, ""},
		{"fixture/circle.png", 30, 150, 0, 20, "northeast"},
		{"fixture/circle.png", 30, 150, 0, 0, "northeast"},
		{"fixture/circle.png", 30, 150, 20, 0, "northeast"},
		{"fixture/circle.png", 30, 150, 0, 0, "northeast"}, // should be equal as prev
	}

	if err := os.Mkdir("test-out", os.ModeDir|os.ModePerm); os.IsNotExist(err) {
		t.Fatal(err)
	}

	for i, x := range tests {
		before, err := ioutil.ReadFile(x.filename)

		if err != nil {
			t.Fatal(err)
		}

		after, err := Crop(before, x.w, x.h, x.x, x.y, x.direction)

		if err != nil {
			t.Fatal(err)
		}

		filename := fmt.Sprintf("test-out/crop-%dx%d+%d+%d-%s-%d%s", x.w, x.h, x.x, x.y, x.direction, i, filepath.Ext(x.filename))
		err = ioutil.WriteFile(filename, after, os.ModePerm)

		if err != nil {
			t.Fatal(err)
		}
	}
}
