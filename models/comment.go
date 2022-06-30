package models

import "time"

type Comment struct {
	ID         uint      `form:"id" json:"id" gorm:"primaryKey"`
	User_id    uint      `form:"user_id" json:"user_id" xml:"user_id" binding:"required"`
	Photo_id   uint      `form:"photo_id" json:"photo_id" xml:"photo_id" binding:"required"`
	Message    uint      `form:"message" json:"message" gorm:"primaryKey"`
	Created_at time.Time `form:"created_at" json:"created_at" xml:"created_at"`
	Updated_at time.Time `form:"updated_at" json:"updated_at" xml:"updated_at"`
}
