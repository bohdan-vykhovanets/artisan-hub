package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/bohdan-vykhovanets/artisan-hub/internal/cache"
	"github.com/bohdan-vykhovanets/artisan-hub/internal/database"
	"github.com/bohdan-vykhovanets/artisan-hub/internal/models"
	"github.com/bohdan-vykhovanets/artisan-hub/internal/websocket"
	"github.com/gin-gonic/gin"
)

type WebSocketMessage struct {
	Event   string         `json:"event"`
	Artwork models.Artwork `json:"artwork"`
}

// CreateArtwork - POST /artworks
func CreateArtwork(c *gin.Context) {
	var artwork models.Artwork

	if err := c.ShouldBindJSON(&artwork); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var artist models.Artist
	if err := database.DB.First(&artist, artwork.ArtistID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ArtistID does not exist"})
		return
	}

	result := database.DB.Create(&artwork)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	msg := WebSocketMessage{Event: "artwork_created", Artwork: artwork}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling websocket message:", err)
	} else {
		websocket.AppHub.Broadcast(msgBytes)
	}

	c.JSON(http.StatusCreated, artwork)
}

// GetArtworks - GET /artworks
func GetArtworks(c *gin.Context) {
	var artworks []models.Artwork

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	result := database.DB.Limit(limit).Offset(offset).Find(&artworks)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   artworks,
		"limit":  limit,
		"offset": offset,
	})
}

// GetArtwork - GET /artworks/:id
func GetArtwork(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "artwork_" + id

	if cachedArtwork, found := cache.AppCache.Get(cacheKey); found {
		fmt.Println("CACHE HIT for artwork:", id)
		c.JSON(http.StatusOK, cachedArtwork)
		return
	}

	fmt.Println("CACHE MISS for artwork:", id)
	var artwork models.Artwork

	if err := database.DB.First(&artwork, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artwork not found"})
		return
	}

	cache.AppCache.Set(cacheKey, artwork, cache.DefaultExpiration)

	c.JSON(http.StatusOK, artwork)
}

// UpdateArtwork - PUT /artworks/:id
func UpdateArtwork(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "artwork_" + id
	var artwork models.Artwork

	if err := database.DB.First(&artwork, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artwork not found"})
		return
	}

	var input models.Artwork
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&artwork).Updates(input)

	fmt.Println("CACHE DELETE for artwork:", id)
	cache.AppCache.Delete(cacheKey)

	var updatedArtwork models.Artwork
	database.DB.First(&updatedArtwork, artwork.ID)

	msg := WebSocketMessage{Event: "artwork_updated", Artwork: updatedArtwork}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling websocket message:", err)
	} else {
		websocket.AppHub.Broadcast(msgBytes)
	}

	c.JSON(http.StatusOK, updatedArtwork)
}

// DeleteArtwork - DELETE /artworks/:id
func DeleteArtwork(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "artwork_" + id
	var artwork models.Artwork

	if err := database.DB.First(&artwork, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artwork not found"})
		return
	}

	database.DB.Delete(&artwork)

	fmt.Println("CACHE DELETE for artwork:", id)
	cache.AppCache.Delete(cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Artwork deleted successfully"})
}

func ShowArtworkPage(c *gin.Context) {
	id := c.Param("id")
	var artwork models.Artwork

	if err := database.DB.First(&artwork, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{"error": "Artwork not found"})
		return
	}

	c.HTML(http.StatusOK, "artwork.html", artwork)
}
