package web

import (
	"fmt"
	"io"

	//"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/unrolled/secure"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/w180112/docx-replacer/pkg/constants"
	"github.com/w180112/docx-replacer/pkg/text"
)

var r *gin.Engine

func tlsHandler(port int) gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     fmt.Sprintf(":%d", constants.HTTPAPIListenPort),
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		if err != nil {
			return
		}

		c.Next()
	}
}

func HttpServer(htmlTemplatePath string) {
	/*httpAPIListener, err := net.Listen("tcp", fmt.Sprintf(":%d", constants.HTTPAPIListenPort))
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to listen on http api %d", constants.HTTPAPIListenPort)
	}
	defer func() {
		_ = httpAPIListener.Close()
	}()*/
	r = gin.Default()
	r.LoadHTMLGlob(htmlTemplatePath)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.POST("/docx/upload", UploadDocx)
	r.GET("/docx/download", DownloadDocx)
	// listen and serve on 0.0.0.0:8080
	// on windows "localhost:8080"
	// can be overriden with the PORT env var

	/*err = r.RunListener(httpAPIListener)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to run http server")
	}*/
	path, err := os.Getwd()
	if err != nil {
		logrus.WithError(err).Error("get pwd failed")
		return
	}
	fmt.Printf("pwd = %s\n", path)
	r.Use(tlsHandler(constants.HTTPAPIListenPort))
	r.RunTLS(fmt.Sprintf(":%d", constants.HTTPAPIListenPort), filepath.Join(path, "certs/server.crt"), filepath.Join(path, "certs/server.key"))
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
	text.FindAndReplace(constants.DocxFilePath+filename, constants.DocxFilePath+"cht_"+filename)
	/*filepath := []string{"file is at: ", constants.DocxFilePath + filename + ", redirecting..."}
	c.Writer.Header().Set("filename", filename)
	c.HTML(http.StatusFound, "res.tmpl", gin.H{
		"Datas": filepath,
	})*/
	filenameArr := strings.Split(filename, ".")
	q := url.Values{}
	q.Set("filename", filenameArr[0])
	location := url.URL{Path: "/docx/download", RawQuery: q.Encode()}
	//c.Request.URL.Path = "/docx/download"
	//c.Request.URL.RawQuery = q.Encode()
	//r.HandleContext(c)
	//c.Writer.Header().Set("filename", filename)
	c.Redirect(http.StatusFound, location.RequestURI())
}

func DownloadDocx(c *gin.Context) {
	filename := c.Query("filename")
	//filename := c.Param("filename")
	if filename == "" {
		c.Data(http.StatusBadRequest, "text/plain", []byte("filename name cannot be empty"))
		return
	}
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+"cht_"+filename+".docx")
	c.Header("Content-Type", "application/octet-stream")
	c.File(constants.DocxFilePath + "cht_" + filename + ".docx")
}
