package main

import (
	blog_db "sample_blog/blog_db"
	"sample_blog/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	blog_db.Init()
	router.Init()
}
