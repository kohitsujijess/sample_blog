package controller

import (
	"errors"
	"net/http"
	"sample_blog/blog_db"
	"sample_blog/model"

	"github.com/labstack/echo"
)

func GetEntryList() echo.HandlerFunc {
	// TODO ent を使って entries から一覧を取得
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Entry list")
	}
}

func GetEntryByUuId() echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid := c.Param("uuid")
		entry, error := SelectEntryWithUuid(uuid)
		if error != nil {
			return c.String(http.StatusBadRequest, "Failed to connect to DB")
		}
		return c.String(http.StatusOK, "Entry by Uuid "+entry.Uuid)
	}
}

func SelectEntryWithUuid(uuid string) (model.Entry, error) {
	var entry model.Entry
	db, err := blog_db.Connect()
	if err != nil {
		return model.Entry{}, errors.New("failed to connect to DB")
	}
	result := db.First(&entry, "uuid = ?", uuid)
	if result.RowsAffected == 0 {
		return model.Entry{}, errors.New("entry data not found")
	}
	return entry, nil
}

func SelectEntries() (model.Entry, error) {
	var entries model.Entry
	db, err := blog_db.Connect()
	if err != nil {
		return model.Entry{}, errors.New("failed to connect to DB")
	}
	result := db.Find(&entries)
	if result.RowsAffected == 0 {
		return model.Entry{}, errors.New("entry data not found")
	}
	return entries, nil
}
