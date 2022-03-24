package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kohitsujijess/sample_blog/controller"
	"github.com/kohitsujijess/sample_blog/job"
	"github.com/kohitsujijess/sample_blog/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToTestDB() (*gorm.DB, error) {
	user := "test_blogger"
	pass := "tset_reggolb"
	protocol := "tcp(db-test-container)"
	protocol = "tcp(localhost)"
	protocol = "tcp(" + os.Getenv("TEST_DB") + ")"
	dbName := "sample_blog_test"
	dataSourceName := user + ":" + pass + "@" + protocol + "/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	count := 0
	if err != nil {
		for {
			if err == nil {
				break
			}
			time.Sleep(time.Second)
			count++
			if count > 120 {
				fmt.Println("failed to connect to test DB:", err)
			}
			db, err = gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
		}
	} else {
		fmt.Println("connect to test DB using gorm")
	}
	db.AutoMigrate(&models.Entry{})
	return db, err
}

func TestCreate(t *testing.T) {
	db, _ := ConnectToTestDB()
	client, _ := db.DB()
	defer client.Close()
	t.Run(
		"Add new entry",
		func(t *testing.T) {
			data := models.Entry{
				ID:          "qwertyuiop",
				Title:       "first test entry",
				Description: "first test entry",
				Body:        "first test entry by test blogger",
			}

			result := db.Create(&data)
			if result.Error != nil {
				t.Error(result.Error)
			}
			t.Cleanup(func() {
				db.Delete(&data)
			})

			addedData := &models.Entry{}
			resultData := db.Last(&addedData)
			t.Cleanup(func() {
				db.Delete(&addedData)
			})

			if resultData.Error != nil {
				t.Error(resultData.Error)
			}
			if addedData.Title != data.Title {
				t.Errorf("expected: %s, got: %s", data.ID, addedData.Title)
			}
		},
	)
}

func TestSave(t *testing.T) {
	db, _ := ConnectToTestDB()
	client, _ := db.DB()
	defer client.Close()
	t.Run(
		"Save entry",
		func(t *testing.T) {
			id := "qwertyuiop"
			data := models.Entry{
				ID:          id,
				Title:       "first test entry",
				Description: "first test entry",
				Body:        "first test entry by test blogger",
			}
			result := db.Create(&data)
			if result.Error != nil {
				t.Error(result.Error)
			}
			t.Cleanup(func() {
				db.Delete(&data)
			})

			entry := &models.Entry{}
			db.First(&entry, "id = ?", id)
			entry.Title = "updated entry"
			entry.Description = "updated entry"
			entry.Body = "updated entry"
			db.Save(&entry)
			t.Cleanup(func() {
				db.Delete(&entry)
			})

			if entry.Title == data.Title {
				t.Errorf("expected: %s, got: %s", entry.Title, data.Title)
			}
		},
	)
}

func TestSelectEntries(t *testing.T) {
	db, _ := ConnectToTestDB()
	client, _ := db.DB()
	defer client.Close()
	t.Run(
		"Select entry",
		func(t *testing.T) {
			data1 := &models.Entry{
				ID:          "1qaz2wsx",
				Title:       "test entry",
				Description: "test entry",
				Body:        "test entry",
			}
			db.Create(&data1)
			t.Cleanup(func() {
				db.Delete(&data1)
			})
			data2 := &models.Entry{
				ID:          "3edc4rfv",
				Title:       "test entry",
				Description: "test entry",
				Body:        "test entry",
			}
			db.Create(&data2)
			t.Cleanup(func() {
				db.Delete(&data2)
			})

			entries, err := models.SelectEntries(db, 2, 0)
			if err != nil {
				t.Error(err.Error())
			}
			if len(entries) != 2 {
				t.Errorf("expected: two entries, got: %d entries", len(entries))
			}
			if entries[0].ID != data2.ID {
				t.Errorf("expected: %s, got: %s", data2.ID, entries[0].ID)
			}
			if entries[1].ID != data1.ID {
				t.Errorf("expected: %s, got: %s", data1.ID, entries[1].ID)
			}
		},
	)
	t.Run(
		"No entries",
		func(t *testing.T) {
			entries, err := models.SelectEntries(db, 0, 0)
			if len(entries) != 0 {
				t.Errorf("expected: zero entries, got: %d entries", len(entries))
			}
			if err != nil {
				t.Errorf("expected: nil, got: %s", err)
			}
		},
	)
}

func TestSelectEntryWithId(t *testing.T) {
	db, _ := ConnectToTestDB()
	client, _ := db.DB()
	defer client.Close()
	t.Run(
		"Select entry",
		func(t *testing.T) {
			data := &models.Entry{
				ID:          "asdfghjkl",
				Title:       "test entry title",
				Description: "test entry description",
				Body:        "test entry body",
			}
			db.Create(&data)
			t.Cleanup(func() {
				db.Delete(&data)
			})

			resultData, err := models.SelectEntryWithId(data.ID, db)
			if err != nil {
				t.Error(err.Error())
			}
			t.Cleanup(func() {
				db.Delete(&resultData)
			})

			if data.Title != resultData.Title {
				t.Errorf("expected: %s, got: %s", data.Title, resultData.Title)
			}
		},
	)
	t.Run(
		"Record is not found",
		func(t *testing.T) {
			data := &models.Entry{
				ID:          "aadfghjkl",
				Title:       "test entry title",
				Description: "test entry description",
				Body:        "test entry body",
			}
			db.Create(&data)
			t.Cleanup(func() {
				db.Delete(&data)
			})

			_, err := models.SelectEntryWithId("bbdfghjkl", db)
			if err == nil {
				t.Errorf("expected: error, got: nil")
			}
		},
	)
}

func TestAddOrUpdateEntry(t *testing.T) {
	db, _ := ConnectToTestDB()
	client, _ := db.DB()
	defer client.Close()
	t.Run(
		"Insert or update entry",
		func(t *testing.T) {
			data := models.Entry{
				ID:          "zxcvbnm",
				Title:       "test entry title",
				Description: "test entry description",
				Body:        "test entry body",
			}
			db.Create(&data)
			t.Cleanup(func() {
				db.Delete(&data)
			})

			entryData := models.Entry{
				ID:          data.ID,
				Title:       "updated entry title",
				Description: "updated entry description",
				Body:        "updated entry body",
			}
			models.AddOrUpdateEntry(db, entryData)
			t.Cleanup(func() {
				db.Delete(&entryData)
			})

			resultData := models.Entry{}
			result := db.First(&resultData, "id = ?", entryData.ID)
			if result.Error != nil {
				t.Error(result.Error)
			}

			if entryData.Title != resultData.Title {
				t.Errorf("expected: %s, got: %s", entryData.Title, resultData.Title)
			}
		},
	)
	t.Run(
		"No data is created nor updated",
		func(t *testing.T) {
			data := models.Entry{
				ID:          "mlpnkobji",
				Title:       "test entry title",
				Description: "test entry description",
				Body:        "test entry body",
			}
			db.Create(&data)
			t.Cleanup(func() {
				db.Delete(&data)
			})
			var entries []models.Entry
			result := db.Order("id desc").Find(&entries)
			if result.Error != nil {
				t.Error(result.Error)
			}
			entryData := models.Entry{
				ID:   "",
				Body: "updated entry body",
			}
			models.AddOrUpdateEntry(db, entryData)

			var entries2 []models.Entry
			result2 := db.Order("id desc").Find(&entries2)
			if result2.Error != nil {
				t.Error(result.Error)
			}
			if len(entries) != len(entries2) {
				t.Errorf("expected: %d, got: %d", len(entries), len(entries2))
			}
		},
	)
}

func TestAuthenticate(t *testing.T) {
	t.Run(
		"With invalid parameters",
		func(t *testing.T) {
			e := echo.New()
			userJSON := `{"username":"wrong_value","password":"wrong_value"}`
			req := httptest.NewRequest(http.MethodPost, "/authenticate", strings.NewReader(userJSON))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			assert.Error(t, controller.Authenticate(c))
		},
	)
}
func TestGetRequest(t *testing.T) {
	t.Run(
		"Url is not correct",
		func(t *testing.T) {
			url := "https://contentful.com/spaces/"
			result := job.GetRequest(url)
			if result == nil {
				t.Errorf("expected: error, got: nil")
			}
		},
	)
}
