package models

import "time"

type User struct {
	ID         uint      `form:"id" json:"id" gorm:"primaryKey"`
	Username   string    `form:"username" json:"username" xml:"username" binding:"required" validate:"required,min=4,max=50"`
	Email      string    `form:"email" json:"email" xml:"email" binding:"required" validate:"required,email"`
	Password   string    `form:"password" json:"password" xml:"password" binding:"required" validate:"required,min=8"`
	Age        int       `form:"age" json:"age" xml:"age" binding:"required" validate:"required,numeric,min=8"`
	Created_at time.Time `form:"created_at" json:"created_at" xml:"created_at"`
	Updated_at time.Time `form:"updated_at" json:"updated_at" xml:"updated_at"`
}
