package handlers

import (
	"net/http"
	"strconv"

	"github.com/bohdan-vykhovanets/artisan-hub/internal/database"
	"github.com/bohdan-vykhovanets/artisan-hub/internal/models"
	"github.com/gin-gonic/gin"
)

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
	var artwork models.Artwork
	id := c.Param("id")

	if err := database.DB.First(&artwork, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artwork not found"})
		return
	}

	c.JSON(http.StatusOK, artwork)
}

// UpdateArtwork - PUT /artworks/:id
func UpdateArtwork(c *gin.Context) {
	var artwork models.Artwork
	id := c.Param("id")

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

	c.JSON(http.StatusOK, artwork)
}

// DeleteArtwork - DELETE /artworks/:id
func DeleteArtwork(c *gin.Context) {
	var artwork models.Artwork
	id := c.Param("id")

	if err := database.DB.First(&artwork, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artwork not found"})
		return
	}

	database.DB.Delete(&artwork)

	c.JSON(http.StatusOK, gin.H{"message": "Artwork deleted successfully"})
}
