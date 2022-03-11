package model

import "time"

type Entry struct {
	Id          string    `json:"id"`
	Uuid        string    `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
