package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	db "sample_blog/blog_db"

	"github.com/labstack/echo"
)

type Entry struct {
	ID          int64
	Uuid        string
	Title       string
	Description string
	Body        string
}

func GetEntryList() echo.HandlerFunc {
	// TODO ent を使って entries から一覧を取得
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Entry list")
	}
}

func GetEntryByUuId() echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid := c.Param("uuid")
		// entry, err := entryByUuID(uuid)
		return c.String(http.StatusOK, "Entry by Uuid "+uuid)
	}
}

func entryByUuID(uuid string) (Entry, error) {
	var entry Entry

	// TODO ent を使って取得
	db := db.Connect()
	row := db.QueryRow("SELECT * FROM entries WHERE uuid = ?", uuid)
	if err := row.Scan(&entry.ID, &entry.Uuid, &entry.Title, &entry.Description, &entry.Body); err != nil {
		if err == sql.ErrNoRows {
			return entry, fmt.Errorf("uuid %v: not found", uuid)
		}
		return entry, fmt.Errorf("uuid %v: %v", uuid, err)
	}
	return entry, nil
}
