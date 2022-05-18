package db

import (
	"gorm.io/gorm"
)

type Reminder struct {
	gorm.Model
	ID        uint64 `gorm:"primaryKey,autoIncrement"`
	UserID    int    `gorm:"index:,unique,composite:usermonth"`
	Celebrant string
	Month     int `gorm:"index:,unique,composite:usermonth"`
	Day       int
	Created   int64 `gorm:"autoCreateTime"`
	Updated   int64 `gorm:"autoUpdateTime"`
}
