package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	var router *gin.Engine = gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Calculator API is running",
			"status":  "succes",
		})
	})

	router.Run(":3000")
}
