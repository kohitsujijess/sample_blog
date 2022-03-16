package blog_db

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	user := os.Getenv("SAMPLE_BLOG_DB_USER")
	pass := os.Getenv("SAMPLE_BLOG_DB_PASS")
	protocol := "tcp(db-container:3306)"
	dbName := "sample_blog"
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
				fmt.Println("failed to connect to DB:", err)
			}
			db, err = gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
		}
	}
	return db, err
}
