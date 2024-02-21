package models

import "gorm.io/gorm"

type Parent struct {
	gorm.Model
	UserID  uint
	Name    string `gorm:"type:varchar(100)"`
	PhoneNo string `gorm:"type:varchar(100)"`
}
