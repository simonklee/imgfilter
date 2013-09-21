package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/simonz05/imgfilter/util"
	"github.com/simonz05/imgfilter/backend"
)

var (
	once       sync.Once
	serverAddr string
	server     *httptest.Server
)

func startServer() {
	util.LogLevel = 0

	fs := NewFileSystem("")
	err := setupServer("travis@tcp(localhost:3306)/myapp_test?charset=utf8", fs)

	if err != nil {
		panic(err)
	}

	server = httptest.NewServer(nil)
	serverAddr = server.Listener.Addr().String()
}
