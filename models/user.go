package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Email string `gorm:"type:varchar(100);unique_index"`
	Gender string `json:"Gender"`
	Password string `json:"Password"`
}