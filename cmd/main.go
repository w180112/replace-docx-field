package main

import (
	"flag"
	"path"
	"runtime"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/w180112/pod-exposer/pkg/text"
	"github.com/w180112/pod-exposer/pkg/web"
)

func main() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			//return frame.Function, fileName
			return "", fileName
		},
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go web.HttpServer()

	inputFile := flag.String("i", "", "input file name")
	outputFile := flag.String("o", "", "output file name")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		logrus.Errorf("argument error, input file path: %s, output file path: %s", *inputFile, *outputFile)
		return
	}

	text.FindAndReplace(*inputFile, *outputFile)
	wg.Wait()

}
