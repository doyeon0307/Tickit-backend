package models

type Ticket struct {
	Id              string  `json:"id" bson:"_id,omitempty"`
	UserId          string  `json:"userId" bson:"userId"`
	Image           string  `json:"image" bson:"image"`
	Title           string  `json:"title" bson:"title"`
	Location        string  `json:"location" bson:"location"`
	Datetime        string  `json:"datetime" bson:"datetime"`
	BackgroundColor string  `json:"backgroundColor" bson:"backgroundColor"`
	ForegroundColor string  `json:"foregroundColor" bson:"foregroundColor"`
	Fields          []Field `json:"fields" bson:"fields"`
}

type Field struct {
	Subtitle string `json:"subtitle" bson:"subtitle"`
	Content  string `json:"content" bson:"content"`
}

type TicketPreview struct {
	Id    string `json:"id" bson:"_id"`
	Image string `json:"image" bson:"image"`
}
