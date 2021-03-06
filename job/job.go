package job

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kohitsujijess/sample_blog/blog_db"
	"github.com/kohitsujijess/sample_blog/models"
)

type AutoGenerated struct {
	Sys struct {
		Type string `json:"type"`
	} `json:"sys"`
	Total int `json:"total"`
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
	Items []struct {
		Metadata struct {
			Tags []interface{} `json:"tags"`
		} `json:"metadata"`
		Sys struct {
			Space struct {
				Sys struct {
					Type     string `json:"type"`
					LinkType string `json:"linkType"`
					ID       string `json:"id"`
				} `json:"sys"`
			} `json:"space"`
			ID          string    `json:"id"`
			Type        string    `json:"type"`
			CreatedAt   time.Time `json:"createdAt"`
			UpdatedAt   time.Time `json:"updatedAt"`
			Environment struct {
				Sys struct {
					ID       string `json:"id"`
					Type     string `json:"type"`
					LinkType string `json:"linkType"`
				} `json:"sys"`
			} `json:"environment"`
			Revision    int `json:"revision"`
			ContentType struct {
				Sys struct {
					Type     string `json:"type"`
					LinkType string `json:"linkType"`
					ID       string `json:"id"`
				} `json:"sys"`
			} `json:"contentType"`
			Locale string `json:"locale"`
		} `json:"sys"`
		Fields struct {
			Title string `json:"title"`
			Body  struct {
				Data struct {
				} `json:"data"`
				Content []struct {
					Data struct {
					} `json:"data"`
					Content []struct {
						Data struct {
						} `json:"data"`
						Marks    []interface{} `json:"marks"`
						Value    string        `json:"value"`
						NodeType string        `json:"nodeType"`
					} `json:"content"`
					NodeType string `json:"nodeType"`
				} `json:"content"`
				NodeType string `json:"nodeType"`
			} `json:"body"`
		} `json:"fields"`
	} `json:"items"`
}

func GetEntriesFromAPI() {
	fmt.Println("GetEntriesFromAPI")
	spaceId := os.Getenv("SAMPLE_BLOG_SPACE_ID")
	accessToken := os.Getenv("SAMPLE_BLOG_ACCESS_TOKEN")
	url := "https://cdn.contentful.com/spaces/" + spaceId + "/entries?access_token=" + accessToken
	GetRequest(url)
}

func GetRequest(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll err=%s", err.Error())
	}
	jsonBytes := ([]byte)(string(body))

	data := new(AutoGenerated)
	if err := json.Unmarshal(jsonBytes, data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return err
	}

	db, e := blog_db.Connect()
	if e != nil {
		fmt.Println("failed to connect to DB:", err)
		return err
	}
	dbClient, _ := db.DB()
	defer dbClient.Close()

	for _, item := range data.Items {
		body, _ := json.Marshal(item.Fields.Body)
		uuid := item.Sys.ID
		title := item.Fields.Title
		description := ""
		bodyString := string(body)

		entry := models.Entry{ID: uuid, Title: title, Description: description, Body: bodyString}
		models.AddOrUpdateEntry(db, entry)
	}
	return nil
}

func Start() {
	for range time.Tick(5 * time.Minute) {
		GetEntriesFromAPI()
	}
}
