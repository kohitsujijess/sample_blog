package blog_db

import (
	"database/sql"
	"fmt"
	"os"
)

var client *sql.DB

func Init() {
	driverName := "mysql"
	user := os.Getenv("SAMPLE_BLOG_DB_USER")
	pass := os.Getenv("SAMPLE_BLOG_DB_PASS")
	dbName := "sample_blog"
	dataSourceName := user + ":" + pass + "@tcp(db-container:3306)/" + dbName
	var err error
	client, err = sql.Open(driverName, dataSourceName)
	// client, err = ent.Open(driverName, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
	//	"db-container", "3306", user, dbName, pass))

	if err != nil {
		fmt.Println("error connecting to database:", err)
	}
	fmt.Println("connect to DB by ent")
}

func Connect() *sql.DB {
	return client
}
