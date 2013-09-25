package server

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/simonz05/imgfilter/image"
	"github.com/simonz05/imgfilter/util"
)

var (
	resizeRe    = regexp.MustCompile("([0-9]+)x([0-9]+)(.+)")
	cropRe      = regexp.MustCompile("([0-9]+)x([0-9]+)(\\+([-0-9]+)\\+([-0-9])+)?(/(northwest|northeast|southwest|southeast|north|west|south|east|center))?(.+)")
	thumbnailRe = regexp.MustCompile("([0-9]+)x([0-9]+)(/(northwest|northeast|southwest|southeast|north|west|south|east|center))?(.+)")
)

type ImageFilter interface {
	SizeParser(string) (*FileInfo, error)
	Filter([]byte, *FileInfo) ([]byte, error)
}

func makeHandler(im ImageFilter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		imageHandle(w, r, im)
	}
}

type ThumbnailFilter struct {
	re *regexp.Regexp
}

func NewThumbnailFilter() *ThumbnailFilter {
	return &ThumbnailFilter{re: thumbnailRe}
}

// parseThumbnail validates the file info for a thumbnail.
func (t *ThumbnailFilter) SizeParser(v string) (f *FileInfo, err error) {
	result := t.re.FindStringSubmatch(v)

	if len(result) != 6 {
		err = errors.New("string mismatch")
		return
	}

	width, err := strconv.ParseUint(result[1], 10, 16)

	if err != nil {
		return
	}

	height, err := strconv.ParseUint(result[2], 10, 16)

	if err != nil {
		return
	}

	filepath := path.Clean(result[5][1:]) // remove first slash

	f = &FileInfo{
		width:     uint(width),
		height:    uint(height),
		direction: result[4],
		filepath:  filepath,
	}
	return
}

func (t *ThumbnailFilter) Filter(data []byte, f *FileInfo) ([]byte, error) {
	return image.Thumbnail(data, f.width, f.height, f.direction)
}

type CropFilter struct {
	re *regexp.Regexp
}

func NewCropFilter() *CropFilter {
	return &CropFilter{re: cropRe}
}

// parseThumbnail validates the file info for a thumbnail.
func (t *CropFilter) SizeParser(v string) (f *FileInfo, err error) {
	result := t.re.FindStringSubmatch(v)
	fmt.Println("res", result, len(result))

	if len(result) != 9 {
		err = errors.New("string mismatch")
		return
	}

	width, err := strconv.ParseUint(result[1], 10, 16)

	if err != nil {
		return
	}

	height, err := strconv.ParseUint(result[2], 10, 16)

	if err != nil {
		return
	}

	x, _ := strconv.ParseInt(result[4], 10, 16)
	y, _ := strconv.ParseInt(result[5], 10, 16)
	filepath := path.Clean(result[8][1:]) // remove first slash

	f = &FileInfo{
		width:     uint(width),
		height:    uint(height),
		direction: result[7],
		x:         int(x),
		y:         int(y),
		filepath:  filepath,
	}
	return
}

func (t *CropFilter) Filter(data []byte, f *FileInfo) ([]byte, error) {
	return image.Crop(data, f.width, f.height, f.x, f.y, f.direction)
}

type ResizeFilter struct {
	re *regexp.Regexp
}

func NewResizeFilter() *ResizeFilter {
	return &ResizeFilter{re: resizeRe}
}

// parseThumbnail validates the file info for a thumbnail.
func (t *ResizeFilter) SizeParser(v string) (f *FileInfo, err error) {
	result := t.re.FindStringSubmatch(v)

	if len(result) != 4 {
		err = errors.New("string mismatch")
		return
	}

	width, err := strconv.ParseUint(result[1], 10, 16)

	if err != nil {
		return
	}

	height, err := strconv.ParseUint(result[2], 10, 16)

	if err != nil {
		return
	}

	filepath := path.Clean(result[3][1:]) // remove first slash

	f = &FileInfo{
		width:    uint(width),
		height:   uint(height),
		filepath: filepath,
	}
	return
}

func (t *ResizeFilter) Filter(data []byte, f *FileInfo) ([]byte, error) {
	return image.Resize(data, f.width, f.height)
}

type FileInfo struct {
	width, height uint
	x, y          int
	direction     string
	filepath      string
}

func (f *FileInfo) String() string {
	return fmt.Sprintf("%dx%d+%d+%d:%s\n%s", f.width, f.height, f.x, f.y, f.direction, f.filepath)
}

func validContentType(mime string) error {
	if mime == "image/png" || mime == "image/jpeg" {
		return nil
	}
	return errors.New("Invalid MIME type")
}

func writeError(w http.ResponseWriter, err string, statusCode int) {
	util.Logf("err: %v", err)
	w.WriteHeader(statusCode)
	w.Write([]byte(err))
}

func imageHandle(w http.ResponseWriter, r *http.Request, f ImageFilter) {
	start := time.Now()
	m := mux.Vars(r)
	util.Logln(m["fileinfo"])

	fi, err := f.SizeParser(m["fileinfo"])

	if err != nil {
		writeError(w, err.Error(), 400)
		return
	}

	util.Logln(fi)

	data, err := imageBackend.ReadFile(fi.filepath)

	if err != nil {
		writeError(w, err.Error(), 400)
		return
	}

	mimeType := http.DetectContentType(data)

	if err := validContentType(mimeType); err != nil {
		writeError(w, err.Error(), 400)
		return
	}

	thumb, err := f.Filter(data, fi)

	if err != nil {
		writeError(w, err.Error(), 400)
		return
	}

	w.Header().Set("Content-Type", mimeType)
	w.Write(thumb)
	util.Logf("Image Handle OK %v", time.Since(start))
}
