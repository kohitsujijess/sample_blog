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

func SelectEntries(db *gorm.DB, limit, offset int) ([]Entry, error) {
	var entries []Entry
	if limit == 0 {
		limit = 40
	}
	result := db.Limit(limit).Offset(offset).Find(&entries)
	if result.Error == gorm.ErrRecordNotFound {
		return entries, fmt.Errorf("SelectEntries ErrRecordNotFound")
	}
	return entries, nil
}

func AddOrUpdateEntry(db *gorm.DB, entry Entry) {
	resultData := Entry{}
	result := db.First(&resultData, "id = ?", entry.ID)
	if result.Error == gorm.ErrRecordNotFound {
		db.Create(&entry)
	} else {
		db.Save(&entry)
	}
}
