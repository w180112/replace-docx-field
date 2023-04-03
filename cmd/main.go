package main

import (
	"flag"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/w180112/docx-replacer/pkg/text"
	"github.com/w180112/docx-replacer/pkg/web"
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

	inputFile := flag.String("i", "", "input file name")
	outputFile := flag.String("o", "", "output file name")
	htmlTemplatePath := flag.String("p", "/etc/docx-replacer/template", "docx replacer web page template")
	flag.Parse()
	*htmlTemplatePath = filepath.Join(*htmlTemplatePath, "*")

	if *inputFile == "" || *outputFile == "" {
		logrus.Errorf("argument error, input file path: %s, output file path: %s", *inputFile, *outputFile)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		web.HttpServer(*htmlTemplatePath)
	}()

	text.FindAndReplace(*inputFile, *outputFile)
	wg.Wait()

}
