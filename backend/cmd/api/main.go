package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"calculator-api/internal/handlers"
)

func main() {
	router := gin.Default()

	// ── CORS ──────────────────────────────────────────────────────────────────
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ── Health check ──────────────────────────────────────────────────────────
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ── API v1 routes ─────────────────────────────────────────────────────────
	v1 := router.Group("/api/v1")
	{
		v1.POST("/add",        handlers.Add)
		v1.POST("/subtract",   handlers.Subtract)
		v1.POST("/multiply",   handlers.Multiply)
		v1.POST("/divide",     handlers.Divide)
		v1.POST("/power",      handlers.Power)
		v1.POST("/percentage", handlers.Percentage)
		v1.POST("/sqrt",       handlers.Sqrt)
	}

	// ── Start server ──────────────────────────────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Calculator API starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
