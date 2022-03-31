package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Code  string
	Name  string
	Email string
}
