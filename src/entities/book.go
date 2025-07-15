package models

import (
	"time"
)

type Book struct {
	ID             int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	ISBN          string    `json:"isbn"`
	PublicationDate time.Time `json:"publication_date"`
}