package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kohitsujijess/sample_blog/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToTestDB() (*gorm.DB, error) {
	user := os.Getenv("SAMPLE_BLOG_DB_USER")
	pass := os.Getenv("SAMPLE_BLOG_DB_PASS")
	protocol := "tcp(db-test-container)"
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

			resultData := models.Entry{}
			result := db.First(&resultData, "id = ?", entryData.ID)
			if result.Error != nil {
				t.Error(result.Error)
			}
			t.Cleanup(func() {
				db.Delete(&resultData)
			})

			if entryData.Title != resultData.Title {
				t.Errorf("expected: %s, got: %s", entryData.Title, resultData.Title)
			}
		},
	)
}
