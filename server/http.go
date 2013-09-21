package server

import (
	"net/http"

	"github.com/simonz05/imgfilter/util"
)

func writeError(w http.ResponseWriter, err string, statusCode int) {
	util.Logf("err: %v", err)
	w.WriteHeader(statusCode)
	w.Write([]byte(err))
}

func cropHandle(w http.ResponseWriter, r *http.Request) {
	util.Logf("Crop Handle")
}

func thumbnailHandle(w http.ResponseWriter, r *http.Request) {
	util.Logf("Thumbnail Handle")
}

func resizeHandle(w http.ResponseWriter, r *http.Request) {
	util.Logf("Resize Handle")
}
