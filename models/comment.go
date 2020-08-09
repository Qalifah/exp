package models

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	content string `gorm:"type:text"`
	PostID uint	
	UserID uint	
}