package image_retriever

import (
	"net/http"

	"atakan-portfolio.com/firebase_connector"
	"github.com/gin-gonic/gin"
)

func GetImages(c *gin.Context) {
	c.JSON(http.StatusOK, firebase_connector.GetImages())
}
