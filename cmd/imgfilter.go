package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/simonz05/imgfilter/backend"
	"github.com/simonz05/imgfilter/server"
	"github.com/simonz05/imgfilter/util"
)

var (
	verbose            = flag.Bool("v", false, "verbose mode")
	help               = flag.Bool("h", false, "show help text")
	laddr              = flag.String("http", ":8080", "set bind address for the HTTP server")
	logLevel           = flag.Int("log", 0, "set log level")
	version            = flag.Bool("version", false, "show version number and exit")
	fsBaseDir          = flag.String("fs-base-dir", "", "file system base dir")
	awsAccessKeyId     = flag.String("aws-secret-access-key", "", "AWS access key id")
	awsSecretAccessKey = flag.String("aws-secret-access-key", "", "AWS secret access key")
	awsRegion          = flag.String("aws-region", "", "AWS region")
	awsBucket          = flag.String("aws-bucket", "", "AWS bucket")
	cpuprofile         = flag.String("debug.cpuprofile", "", "write cpu profile to file")
)

var Version = "0.1.0"

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Fprintln(os.Stdout, Version)
		return
	}

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	if *laddr == "" {
		fmt.Fprintln(os.Stderr, "listen address required")
		os.Exit(1)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if *verbose {
		util.LogLevel = 2
	} else {
		util.LogLevel = *logLevel
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var imgBackend backend.ImageBackend

	if *fsBaseDir != "" {
		imgBackend = backend.Dir(*fsBaseDir)
	} else if *awsAccessKeyId != "" && *awsSecretAccessKey != "" && *awsRegion != "" && *awsBucket != "" {
		imgBackend = backend.NewS3(*awsAccessKeyId, *awsSecretAccessKey, *awsRegion, *awsBucket)
	} else {
		util.Errln("Expected either aws-* or fs-* arguments")
		os.Exit(1)
	}

	err := server.ListenAndServe(*laddr, imgBackend)

	if err != nil {
		util.Logln(err)
	}
}
