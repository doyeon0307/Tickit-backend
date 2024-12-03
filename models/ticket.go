package models

import "time"

type Ticket struct {
	Id              string    `json:"id" bson:"_id,omitempty"`
	UserId          string    `json:"userId" bson:"userId"`
	Image           string    `json:"image" bson:"image"`
	Title           string    `json:"title" bson:"title"`
	Location        string    `json:"location" bson:"location"`
	DateTime        time.Time `json:"dateTime" bson:"dateTime"`
	BackgroundColor string    `json:"backgroundColor" bson:"backgroundColor"`
	ForegroundColor string    `json:"foregroundColor" bson:"foregroundColor"`
	Fields          []Field   `json:"fields" bson:"fields"`
	CreatedAt       time.Time `json:"createdAt" bson:"createdAt"`
}

type Field struct {
	Subtitle string `json:"subtitle" bson:"subtitle"`
	Content  string `json:"content" bson:"content"`
}
