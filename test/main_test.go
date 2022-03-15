package dockertest_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Entry struct {
	ID          string
	Title       string
	Description string
	Body        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

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
	db.AutoMigrate(&Entry{})
	return db, err
}

func TestCreate(t *testing.T) {
	db, _ := ConnectToTestDB()
	client, _ := db.DB()
	defer client.Close()
	t.Run(
		"Add new entry",
		func(t *testing.T) {
			data := &Entry{
				ID:          "qwertyuiop",
				Title:       "first test entry",
				Description: "first test entry",
				Body:        "first test entry by test blogger",
			}

			result := db.Create(data)
			if result.Error != nil {
				t.Error(result.Error)
			}

			addedData := &Entry{}
			resultData := db.Last(&addedData)
			if resultData.Error != nil {
				t.Error(resultData.Error)
			}
			if addedData.Title != data.Title {
				t.Errorf("expected: %s, addedData: %s", data.ID, addedData.Title)
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
			data := &Entry{
				ID:          "qwertyuiop",
				Title:       "first test entry",
				Description: "first test entry",
				Body:        "first test entry by test blogger",
			}

			result := db.Create(data)
			if result.Error != nil {
				t.Error(result.Error)
			}

			addedData := &Entry{}
			resultData := db.Last(&addedData)
			if resultData.Error != nil {
				t.Error(resultData.Error)
			}
			if addedData.Title != data.Title {
				t.Errorf("expected: %s, addedData: %s", data.ID, addedData.Title)
			}
		},
	)
}
