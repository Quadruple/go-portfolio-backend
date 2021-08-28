package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"atakan-portfolio.com/image_retriever"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	}
}

func main() {
	router := gin.Default()

	router.Use(CORSMiddleware())

	router.GET("/", HelloThere)
	router.GET("/media", image_retriever.GetImages)

	router.Run("localhost:8080")
}

func HelloThere(c *gin.Context) {
	c.String(http.StatusOK, `
		Hello there, welcome to my portfolio's backend service written with Go.
		I created this API in order to practise with Go, Gin and Firebase.
		Enjoy. :)`)
}
