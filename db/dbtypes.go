package db

import (
	"gorm.io/gorm"
)

type Date map[string]string

type Reminder struct {
	gorm.Model
	ID      uint64 `gorm:"primaryKey,autoIncrement"`
	UserID  string `gorm:"index:,unique,composite:usermonth"`
	Month   uint16 `gorm:"index:,unique,composite:usermonth"`
	Dates   []Date `gorm:"embedded"`
	Created int64  `gorm:"autoCreateTime"`
	Updated int64  `gorm:"autoUpdateTime"`
}
