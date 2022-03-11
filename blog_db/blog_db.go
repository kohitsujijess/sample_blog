package blog_db

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

var db *sql.DB
var err error

func Connect() (*sql.DB, error) {
	driverName := "mysql"
	user := os.Getenv("SAMPLE_BLOG_DB_USER")
	pass := os.Getenv("SAMPLE_BLOG_DB_PASS")
	protocol := "tcp(db-container:3306)"
	dbName := "sample_blog"
	dataSourceName := user + ":" + pass + "@" + protocol + "/" + dbName + "?parseTime=True&loc=Asia%2FTokyo"
	db, err = sql.Open(driverName, dataSourceName)
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
			db, err = sql.Open(driverName, dataSourceName)
		}
	} else {
		fmt.Println("connect to DB using sql")
	}
	return db, err
}
