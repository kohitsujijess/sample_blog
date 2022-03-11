package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sample_blog/blog_db"
	"sample_blog/model"
	"strconv"

	"github.com/labstack/echo"
)

func GetEntryList() echo.HandlerFunc {
	return func(c echo.Context) error {
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
}

func GetEntryByUuId() echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid := c.Param("uuid")
		entry, error := SelectEntryWithUuid(uuid)
		if error != nil {
			return c.String(http.StatusBadRequest, "Failed to connect to DB")
		}
		entryJson, _ := json.Marshal(entry)
		return c.String(http.StatusOK, string(entryJson))
	}
}

func SelectEntryWithUuid(uuid string) (model.Entry, error) {
	var entry model.Entry
	db, err := blog_db.Connect()
	if err != nil {
		return model.Entry{}, errors.New("failed to connect to DB")
	}
	// result := db.First(&entry, "uuid = ?", uuid)
	row := db.QueryRow("SELECT * FROM entries WHERE uuid = ?", uuid)
	if err := row.Scan(&entry.Id, &entry.Uuid, &entry.Title,
		model.SkippedScanner{}, model.SkippedScanner{}, model.SkippedScanner{}, model.SkippedScanner{}); err != nil {
		if err == sql.ErrNoRows {
			return entry, fmt.Errorf("SelectEntryWithUuid %v: not found", uuid)
		}
		return entry, fmt.Errorf("SelectEntryWithUuid %v: %v", uuid, err)
	}
	return entry, nil
}

func SelectEntries(limit, offset int) ([]model.Entry, error) {
	var entries []model.Entry
	if limit == 0 {
		limit = 40
	}
	db, err := blog_db.Connect()
	if err != nil {
		return entries, errors.New("failed to connect to DB")
	}
	// db.Limit(limit).Offset(offset).Find(&entries)

	rows, err := db.Query("SELECT * FROM entries LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("SelectEntries %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry model.Entry
		if err := rows.Scan(&entry.Id, &entry.Uuid, &entry.Title,
			model.SkippedScanner{}, model.SkippedScanner{}, model.SkippedScanner{}, model.SkippedScanner{}); err != nil {
			return nil, fmt.Errorf("SelectEntries %v", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("SelectEntries %v", err)
	}
	return entries, nil
}
