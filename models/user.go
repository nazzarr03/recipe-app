package models

import "time"

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`

	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:null"`
	DeletedAt time.Time `gorm:"default:null"`
}
