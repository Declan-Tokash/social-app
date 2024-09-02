package model

import (
	"gorm.io/gorm"
   )

type Post struct {
	gorm.Model
    Title    string    `json:"title"`
	Image    string    `json:image`
	Location Location  `gorm:"embedded"`
}

type NearbyPost struct {
    Post
}
