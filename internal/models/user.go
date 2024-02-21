package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	KakaoId       uint64 `gorm:"uniqueIndex"`
	Name          string `gorm:"type:varchar(100)"`
	PhoneNo       string `gorm:"type:varchar(100)"`
	Subscriptions []Subscription
	Parents       []Parent
}
