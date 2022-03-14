package test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Entry struct {
	ID          string
	Title       string
	Description string
	Body        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func createContainer() (*dockertest.Resource, *dockertest.Pool) {
	pwd, _ := os.Getwd()

	pool, err := dockertest.NewPool("")
	pool.MaxWait = time.Minute * 2
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	runOptions := &dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.0",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=test_blog_echo",
		},
		Mounts: []string{
			pwd + "/my.cnf:/etc/mysql/my.cnf",                            // MySQLの設定ファイル
			pwd + "/entries.sql:/docker-entrypoint-initdb.d/entries.sql", // コンテナ起動時に実行したいSQL
		},
		Cmd: []string{
			"mysqld",
			"--character-set-server=utf8mb4",
			"--collation-server=utf8mb4_unicode_ci",
		},
	}

	resource, err := pool.RunWithOptions(runOptions)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	return resource, pool
}

func closeContainer(resource *dockertest.Resource, pool *dockertest.Pool) {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func connectToTestDB(resource *dockertest.Resource, pool *dockertest.Pool) *gorm.DB {
	var db *gorm.DB
	if err := pool.Retry(func() error {
		var e error
		time.Sleep(time.Second * 100)

		user := "test_blogger"
		pass := "test_blog_echo"
		protocol := "(localhost:" + resource.GetPort("3306/tcp") + ")"
		dbName := "sample_blog_test"
		dataSourceName := user + ":" + pass + "@" + protocol + "/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, e = gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})

		count := 0
		if e != nil {
			for {
				if e == nil {
					break
				}
				time.Sleep(time.Second)
				count++
				if count > 500 {
					fmt.Println("ERROR")
					return e
				}
				db, e = gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
			}
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return db
}

func TestCreateWithDB(t *testing.T) {
	t.Run(
		"Add new entry",
		func(t *testing.T) {
			resource, pool := createContainer()
			defer closeContainer(resource, pool)
			db := connectToTestDB(resource, pool)
			data := &Entry{
				ID:          "qwertyuiop",
				Title:       "first test entry",
				Description: "first test entry",
				Body:        "first test entry by test blogger",
			}

			result := db.Create(data)
			if result.Error != nil {
				t.Error(result.Error)
			}

			addedData := &Entry{}
			resultData := db.Last(&addedData)
			if resultData.Error != nil {
				t.Error(resultData.Error)
			}
			if addedData.Title != data.Title {
				t.Errorf("expected: %s, addedData: %s", data.ID, addedData.Title)
			}
		},
	)
}
