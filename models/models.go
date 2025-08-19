package models

import "time"

type Book struct {
	ID          int       `json:"id" bson:"id"`
	ISBN        string    `json:"isbn" bson:"isbn"`
	Title       string    `json:"title" bson:"title"`
	Author      string    `json:"author" bson:"author"`
	Publisher   string    `json:"publisher" bson:"publisher"`
	PublishedAt time.Time `json:"published_at" bson:"published_at"`
	Genre       string    `json:"genre" bson:"genre"`
	Language    string    `json:"language" bson:"language"`
	Pages       int       `json:"pages" bson:"pages"`
	Description string    `json:"description" bson:"description"`
	CoverURL    string    `json:"coverURL" bson:"coverURL"`

	Location  string    `json:"location" bson:"location"` // shelf location
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
