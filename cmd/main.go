package main

import (
	"flag"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
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

	//inputFile := flag.String("i", "", "input file name")
	//outputFile := flag.String("o", "", "output file name")
	htmlTemplatePath := flag.String("p", "/etc/docx-replacer/template", "docx replacer web page template")
	flag.Parse()
	*htmlTemplatePath = filepath.Join(*htmlTemplatePath, "*")

	/*if *inputFile == "" || *outputFile == "" {
		logrus.Errorf("argument error, input file path: %s, output file path: %s", *inputFile, *outputFile)
		return
	}*/

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		web.HttpServer(*htmlTemplatePath)
	}()

	//text.FindAndReplace("/home/the/Desktop/docx-test/LittleBabyBumEp5.docx", "/home/the/Desktop/docx-test/edit.docx")
	//text.FindAndReplace("/home/the/Downloads/en_DIALOG_LIST_LittleBabyBumMusicTime_LittleBabyBumMusicTimeSeason1HaveYouEverSeenALassieThisOldManTeddyBear_TeddyBear_en_V1_Final_Cut_v6_0_xlsx_36_44_88.doc", "/home/the/Desktop/docx-test/edit.docx")
	wg.Wait()

}
