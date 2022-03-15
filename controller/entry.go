package controller

import (
	"encoding/json"
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
	db, _ := blog_db.Connect()
	entries, err := models.SelectEntries(db, limitInt, offsetInt)
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
