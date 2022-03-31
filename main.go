package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bamzi/jobrunner"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/kohitsujijess/sample_blog/blog_db"
	"github.com/kohitsujijess/sample_blog/job"
	"github.com/kohitsujijess/sample_blog/models"
	"github.com/kohitsujijess/sample_blog/router"
	"gorm.io/gorm"
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
	// db.AutoMigrate(&models.Entry{})

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "202203311700",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Author{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("authors")
			},
		},
		{
			ID: "202203311705",
			Migrate: func(tx *gorm.DB) error {
				db.Migrator().AddColumn(&models.Author{}, "Name")
				db.Migrator().AddColumn(&models.Author{}, "Code")
				db.Migrator().AddColumn(&models.Author{}, "Number")

				return tx.AutoMigrate(&models.Author{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropColumn(&models.Author{}, "Name")
				return tx.Migrator().DropColumn(&models.Author{}, "Code")
				return tx.Migrator().DropColumn(&models.Author{}, "Number")
			},
		},
	})
	if err = m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	jobrunner.Start()
	jobrunner.Schedule("@every 5m", GetEntries{})

	e := router.Init()

	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
