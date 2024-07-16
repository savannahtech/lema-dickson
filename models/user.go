package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName string `gorm:"full_name"`
	Username string `gorm:"username"`
}
