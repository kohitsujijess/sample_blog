package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kohitsujijess/sample_blog/blog_db"
	"github.com/kohitsujijess/sample_blog/job"
	"github.com/kohitsujijess/sample_blog/models"
	"github.com/kohitsujijess/sample_blog/router"
)

func main() {
	db, err := blog_db.Connect()
	if err != nil {
		fmt.Println("Failed to connect to DB")
	} else {
		fmt.Println("Connected to DB")
	}
	db.AutoMigrate(&models.Entry{})

	e := router.Init()

	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	job.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
