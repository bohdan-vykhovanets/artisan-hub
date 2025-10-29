package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bohdan-vykhovanets/artisan-hub/internal/cache"
	"github.com/bohdan-vykhovanets/artisan-hub/internal/database"
	"github.com/bohdan-vykhovanets/artisan-hub/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, relying on environment variables.")
	}

	database.ConnectDatabase()
	cache.InitCache()

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/artists")
	})
	r.GET("/artists", handlers.ShowArtistsPage)
	r.GET("/artists/:id", handlers.ShowArtistPage)
	r.GET("/artworks/:id", handlers.ShowArtworkPage)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/artists", handlers.CreateArtist)
		v1.GET("/artists", handlers.GetArtists)
		v1.GET("/artists/:id", handlers.GetArtist)
		v1.PUT("/artists/:id", handlers.UpdateArtist)
		v1.DELETE("/artists/:id", handlers.DeleteArtist)

		v1.POST("/artworks", handlers.CreateArtwork)
		v1.GET("/artworks", handlers.GetArtworks)
		v1.GET("/artworks/:id", handlers.GetArtwork)
		v1.PUT("/artworks/:id", handlers.UpdateArtwork)
		v1.DELETE("/artworks/:id", handlers.DeleteArtwork)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Starting server on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
