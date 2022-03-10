package blog_db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

func Init() {
	driverName := "mysql"
	user := os.Getenv("SAMPLE_BLOG_DB_USER")
	pass := os.Getenv("SAMPLE_BLOG_DB_PASS")
	dbName := "sample_blog"
	dataSourceName := user + ":" + pass + "@tcp(db-container:3306)/" + dbName + "?parseTime=True&loc=Asia%2FTokyo"
	db, err = gorm.Open(driverName, dataSourceName)

	if err != nil {
		fmt.Println("error connecting to database:", err)
	}
	fmt.Println("connect to DB using gorm")
}

func Connect() (*gorm.DB, error) {
	return db, err
}
