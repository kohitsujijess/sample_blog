package model

import "time"

type Entry struct {
	ID          string `gorm:"primary_key"`
	Uuid        string
	Title       string
	Description string
	Body        string
	Created     time.Time `gorm:"autoCreateTime"`
	Updated     time.Time `gorm:"autoUpdateTime"`
}
