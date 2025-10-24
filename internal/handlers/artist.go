package handlers

import (
	"net/http"
	"strconv"

	"github.com/bohdan-vykhovanets/artisan-hub/internal/database"
	"github.com/bohdan-vykhovanets/artisan-hub/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateArtist - POST /artists
func CreateArtist(c *gin.Context) {
	var artist models.Artist

	if err := c.ShouldBindJSON(&artist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Create(&artist)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, artist)
}

// GetArtists - GET /artists
func GetArtists(c *gin.Context) {
	var artists []models.Artist

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	result := database.DB.Limit(limit).Offset(offset).Find(&artists)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   artists,
		"limit":  limit,
		"offset": offset,
	})
}

// GetArtist - GET /artists/:id
func GetArtist(c *gin.Context) {
	var artist models.Artist
	id := c.Param("id")

	if err := database.DB.Preload("Artworks").First(&artist, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
		return
	}

	c.JSON(http.StatusOK, artist)
}

// UpdateArtist - PUT /artists/:id
func UpdateArtist(c *gin.Context) {
	var artist models.Artist
	id := c.Param("id")

	if err := database.DB.First(&artist, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
		return
	}

	var input models.Artist
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&artist).Updates(input)

	c.JSON(http.StatusOK, artist)
}

// DeleteArtist - DELETE /artists/:id
func DeleteArtist(c *gin.Context) {
	var artist models.Artist
	id := c.Param("id")

	if err := database.DB.First(&artist, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
		return
	}

	database.DB.Delete(&artist)

	c.JSON(http.StatusOK, gin.H{"message": "Artist deleted successfully"})
}
