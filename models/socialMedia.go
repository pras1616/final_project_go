package models

import "time"

type SocialMedia struct {
	ID               uint      `form:"id" json:"id" gorm:"primaryKey"`
	Name             string    `form:"name" json:"name" xml:"name" binding:"required"`
	Social_media_url string    `form:"social_media_url" json:"social_media_url" xml:"social_media_url" binding:"required"`
	User_id          uint      `form:"user_id" json:"user_id" xml:"user_id" binding:"required"`
	Created_at       time.Time `form:"created_at" json:"created_at" xml:"created_at"`
	Updated_at       time.Time `form:"updated_at" json:"updated_at" xml:"updated_at"`
}
