package util

import (
	"log"
	"os"
)

var (
	LogLevel int = 0
	Logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	Error  = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)
)

func Errf(fmt string, args ...interface{}) {
	Error.Printf(fmt, args...)
}

func Errln(args ...interface{}) {
	Logger.Println(args...)
}

func Logf(fmt string, args ...interface{}) {
	if LogLevel == 0 {
		return
	}

	Logger.Printf(fmt, args...)
}

func Logln(args ...interface{}) {
	if LogLevel == 0 {
		return
	}

	Logger.Println(args...)
}
