package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kohitsujijess/sample_blog/blog_db"
	"github.com/kohitsujijess/sample_blog/job"
	"github.com/kohitsujijess/sample_blog/models"
	"github.com/kohitsujijess/sample_blog/router"
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
	} else {
		fmt.Println("Connected to DB")
	}
	fmt.Println(db)
	db.AutoMigrate(&models.Entry{})

	/**
	jobrunner.Start()
	jobrunner.Schedule("@every 5m", GetEntries{})
	*/

	e := router.Init()
	e.Start(":1323")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
