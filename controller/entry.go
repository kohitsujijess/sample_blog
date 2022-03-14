package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kohitsujijess/sample_blog/blog_db"

	"github.com/kohitsujijess/sample_blog/models"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func GetEntryList(c echo.Context) error {
	limit := c.Param("limit")
	offset := c.Param("offset")
	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)
	entries, err := SelectEntries(limitInt, offsetInt)
	if err != nil {
		fmt.Println("failed to get entries", err)
	}
	entriesJson, _ := json.Marshal(entries)
	return c.String(http.StatusOK, string(entriesJson))
}

func GetEntryByUuId(c echo.Context) error {
	uuid := c.Param("uuid")
	entry, err := SelectEntryWithUuid(uuid)
	if err != nil {
		return c.String(http.StatusBadRequest, "not found")
	}
	entryJson, _ := json.Marshal(entry)
	return c.String(http.StatusOK, string(entryJson))
}

func SelectEntryWithUuid(uuid string) (models.Entry, error) {
	var entry models.Entry
	db, err := blog_db.Connect()
	if err != nil {
		return models.Entry{}, errors.New("failed to connect to DB")
	}
	result := db.First(&entry, "id = ?", uuid)
	if result.Error == gorm.ErrRecordNotFound {
		return entry, fmt.Errorf("SelectEntryWithUuid: %v", uuid)
	}
	return entry, nil
}

func SelectEntries(limit, offset int) ([]models.Entry, error) {
	var entries []models.Entry
	if limit == 0 {
		limit = 40
	}
	db, err := blog_db.Connect()
	if err != nil {
		return entries, errors.New("failed to connect to DB")
	}
	db.Limit(limit).Offset(offset).Find(&entries)
	return entries, nil
}
