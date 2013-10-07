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

// Package server implements HTTP interface for imgfilter server.
package server

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/simonz05/imgfilter/backend"
	"github.com/simonz05/imgfilter/util"
)

var (
	router         *mux.Router
)

func sigTrapCloser(l net.Listener) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for _ = range c {
			// Once we close the listener the main loop will exit
			l.Close()
			util.Logf("Closed listener %s", l.Addr())
		}
	}()
}

func setupServer(imgBackend backend.ImageBackend) error {
	// HTTP endpoints
	router = mux.NewRouter()
	router.HandleFunc("/crop/", cropHandle).Methods("GET").Name("crop")
	router.HandleFunc("/resize/", resizeHandle).Methods("GET").Name("resize")
	router.HandleFunc("/thumbnail/", thumbnailHandle).Methods("GET").Name("thumbnail")
	router.StrictSlash(false)
	http.Handle("/", router)

	return nil
}

func ListenAndServe(laddr string, imgBackend backend.ImageBackend) error {
	if err := setupServer(imgBackend); err != nil {
		return err
	}

	l, err := net.Listen("tcp", laddr)

	if err != nil {
		return err
	}

	util.Logf("Listen on %s", l.Addr())

	sigTrapCloser(l)
	err = http.Serve(l, nil)
	util.Logf("Shutting down ..")
	return err
}
