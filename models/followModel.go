package models

import "time"

type FollowModel struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	FollowerID uint      `json:"follower_id" gorm:"not null; uniqueIndex:follow_unique"`
	FollowedID uint      `json:"followed_id" gorm:"not null; uniqueIndex:follow_unique"`
	CreatedAt  time.Time `json:"created_at"`
}
