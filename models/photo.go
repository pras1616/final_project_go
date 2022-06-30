package models

import "time"

type Photo struct {
	ID         uint      `form:"id" json:"id" gorm:"primaryKey"`
	Title      string    `form:"title" json:"title" xml:"title" binding:"required"`
	Caption    string    `form:"caption" json:"caption" xml:"caption" binding:"required"`
	Photo_url  string    `form:"photo_url" json:"photo_url" xml:"photo_url" binding:"required"`
	User_id    uint      `form:"user_id" json:"user_id" xml:"user_id" binding:"required"`
	Created_at time.Time `form:"created_at" json:"created_at" xml:"created_at"`
	Updated_at time.Time `form:"updated_at" json:"updated_at" xml:"updated_at"`
}
