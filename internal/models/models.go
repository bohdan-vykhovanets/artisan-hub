package models

import "gorm.io/gorm"

type Artist struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Bio      string
	Artworks []Artwork `gorm:"foreignKey:ArtistID"`
}

type Artwork struct {
	gorm.Model
	Title       string `gorm:"not null"`
	ImageURL    string `gorm:"not null"`
	Description string
	ArtistID    uint
}
