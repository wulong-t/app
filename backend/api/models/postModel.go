package models

import "time"

type PostModel struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	CreatorID    uint      `json:"creator_id" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	Title        string    `json:"title" binding:"required"`
	Content      string    `json:"content" binding:"required"`
	SelectedFile string    `json:"selected_file"`
	// Likes     []string
	// Comments  []string
}
