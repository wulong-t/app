package models

import "time"

type MessageModel struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SenderID   uint      `json:"sender_id" binding:"required"`
	ReceiverID uint      `json:"receiver_id" binding:"required"`
	Content    string    `json:"content" binding:"required"`
	CreatedAt  time.Time `json:"created_at"`
}

type UnReadNumModel struct {
	ID         uint `json:"id" gorm:"primaryKey"`
	SenderID   uint `json:"sender_id" binding:"required"`
	ReceiverID uint `json:"receiver_id" binding:"required"`
	Num        uint `json:"num" binding:"required"`
}
