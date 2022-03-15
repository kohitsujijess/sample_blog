package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kohitsujijess/sample_blog/blog_db"

	"github.com/kohitsujijess/sample_blog/models"

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

func GetEntryById(c echo.Context) error {
	uuid := c.Param("id")
	db, _ := blog_db.Connect()
	entry, err := models.SelectEntryWithId(uuid, db)
	if err != nil {
		return c.String(http.StatusBadRequest, "not found")
	}
	entryJson, _ := json.Marshal(entry)
	return c.String(http.StatusOK, string(entryJson))
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
