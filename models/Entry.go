package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Entry struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `gorm:"type:text" json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func SelectEntryWithId(id string, db *gorm.DB) (Entry, error) {
	var entry Entry
	result := db.First(&entry, "id = ?", id)
	if result.Error == gorm.ErrRecordNotFound {
		return entry, fmt.Errorf("SelectEntryWithUuid: %v", id)
	}
	return entry, nil
}
