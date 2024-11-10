package models

type Schedule struct {
	Id        string `json:"id" bson:"_id,omitempty"`
	UserId    string `json:"userId" bson:"userId"`
	Date      string `json:"date" bson:"date"`
	Title     string `json:"title" bson:"title"`
	Number    int    `json:"number" bson:"number"`
	Image     string `json:"image" bson:"image"`
	Thumbnail bool   `json:"thumbnail" bson:"thumbnail"`
	Location  string `json:"location" bson:"location"`
	Time      string `json:"time" bson:"time"`
	Seat      string `json:"seat" bson:"seat"`
	Casting   string `json:"casting" bson:"casting"`
	Company   string `json:"company" bson:"company"`
	Link      string `json:"link" bson:"link"`
	Memo      string `json:"memo" bson:"memo"`
}
