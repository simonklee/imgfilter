package server

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"time"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/simonz05/imgfilter/util"
	"github.com/simonz05/imgfilter/image"
)

var (
	resizeRe    = regexp.MustCompile("([0-9]+)x([0-9]+)")
	thumbnailRe = regexp.MustCompile("([0-9]+)x([0-9]+)(northwest|northeast|southwest|southeast|north|west|south|east)?(.*)")
)

type FileInfo struct {
	width, height uint
	x, y          int
	direction     string
	filepath      string
}

func (f *FileInfo) String() string {
	return fmt.Sprintf("%dx%d+%d+%d:%s\n%s", f.width, f.height, f.x, f.y, f.direction, f.filepath)
}

// parsePath cleans and validates a filepath
func parsePath(v string) (string, error) {
	p := path.Clean(v)

	switch path.Ext(p) {
	case ".jpeg", ".jpg", ".png":
	default:
		return "", errors.New("invalid ext")
	}
	p = p[1:] // remove first slash
	return p, nil
}

// parseThumbnail validates the file info for a thumbnail.
func parseThumbnail(v string) (f *FileInfo, err error) {
	result := thumbnailRe.FindStringSubmatch(v)

	if len(result) != 5 {
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

	p, err := parsePath(result[4])

	if err != nil {
		return
	}

	f = &FileInfo{
		width:     uint(width),
		height:    uint(height),
		direction: result[3],
		filepath:  p,
	}
	return
}

func writeError(w http.ResponseWriter, err string, statusCode int) {
	util.Logf("err: %v", err)
	w.WriteHeader(statusCode)
	w.Write([]byte(err))
}

func cropHandle(w http.ResponseWriter, r *http.Request) {
	util.Logf("Crop Handle")
}

func thumbnailHandle(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	util.Logf("Thumbnail Handle %v", start)
	m := mux.Vars(r)
	util.Logln(m["fileinfo"])

	fi, err := parseThumbnail(m["fileinfo"])

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

	thumb, err := image.Thumbnail(data, fi.width, fi.height, fi.direction)
	if err != nil {
		writeError(w, err.Error(), 400)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(thumb)
	util.Logf("Thumbnail Handle OK %v", time.Since(start))
}

func resizeHandle(w http.ResponseWriter, r *http.Request) {
	util.Logf("Resize Handle")
}
