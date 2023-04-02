package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpServer() {
	r := gin.Default()
	r.LoadHTMLFiles("index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"books": "123",
		})
	})
	r.POST("/replace", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "User created"})
	})
	// listen and serve on 0.0.0.0:8080
	// on windows "localhost:8080"
	// can be overriden with the PORT env var
	r.Run()
}
