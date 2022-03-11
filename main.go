package main

import (
	"context"
	"fmt"
	"sample_blog/blog_db"
	"sample_blog/job"
	"sample_blog/router"
	"time"

	"github.com/bamzi/jobrunner"
	_ "github.com/go-sql-driver/mysql"
)

type GetEntries struct {
}

func (e GetEntries) Run() {
	job.GetEntriesFromAPI()
}

func main() {
	db, err := blog_db.Connect()
	if err != nil {
		fmt.Println("Failed to connect to DB")
	}
	db.Close()
	jobrunner.Start()
	jobrunner.Schedule("@every 5m", GetEntries{})

	e := router.Init()
	e.Start(":1323")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
