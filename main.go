package main

import (
	"fmt"
	"sample_blog/blog_db"
	"sample_blog/cron"
	"sample_blog/model"
	"sample_blog/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	blog_db.Init()
	db, err := blog_db.Connect()
	if err != nil {
		fmt.Println("Failed to connect to DB")
	}
	db.AutoMigrate(model.Entry{})
	fmt.Println("Migrated")
	router.Init()
	cron.Start()
}
