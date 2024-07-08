package models

import "time"

type Recipe struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title   string `json:"title"`
	Content string `json:"content"`

	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:null"`
	DeletedAt time.Time `gorm:"default:null"`
}
