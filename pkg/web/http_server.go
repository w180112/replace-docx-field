package web

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/w180112/docx-replacer/pkg/constants"
)

func HttpServer(htmlTemplatePath string) {
	httpAPIListener, err := net.Listen("tcp", fmt.Sprintf(":%d", constants.HTTPAPIListenPort))
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to listen on http api %d", constants.HTTPAPIListenPort)
	}
	defer func() {
		_ = httpAPIListener.Close()
	}()
	r := gin.Default()
	r.LoadHTMLGlob(htmlTemplatePath)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.POST("/docx/upload", UploadDocx)
	// listen and serve on 0.0.0.0:8080
	// on windows "localhost:8080"
	// can be overriden with the PORT env var

	err = r.RunListener(httpAPIListener)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to run http server")
	}
}

func UploadDocx(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		logrus.WithError(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	filename := header.Filename
	out, err := os.Create(filepath.Join(constants.DocxFilePath, filename))
	if err != nil {
		logrus.WithError(err).Error("create failed")
		c.Data(http.StatusInternalServerError, "text/plain", []byte("create failed"))
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		logrus.WithError(err).Error("copy file failed")
		c.Data(http.StatusInternalServerError, "text/plain", []byte("copy file failed"))
		return
	}
	filepath := []string{"file is at: ", constants.DocxFilePath + filename}
	c.HTML(http.StatusOK, "res.tmpl", gin.H{
		"Datas": filepath,
	})
}
