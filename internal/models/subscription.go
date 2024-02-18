package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	UserID    uint
	ExpiresAt datatypes.Date
}
