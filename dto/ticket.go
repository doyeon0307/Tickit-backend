package dto

type Field struct {
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
}

type TicketDTO struct {
	UserId          string  `json:"userId"`
	Image           string  `json:"image" binding:"required"`
	Title           string  `json:"title" binding:"required"`
	Location        string  `json:"location" binding:"required"`
	Datetime        string  `json:"datetime" binding:"required"`
	BackgroundColor string  `json:"backgroundColor"`
	ForegroundColor string  `json:"foregroundColor"`
	Fields          []Field `json:"fields"`
}

type TicketResponseDTO struct {
	Id              string  `json:"id"`
	UserId          string  `json:"userId"`
	Image           string  `json:"image"`
	Title           string  `json:"title"`
	Location        string  `json:"location"`
	Datetime        string  `json:"datetime"`
	BackgroundColor string  `json:"backgroundColor"`
	ForegroundColor string  `json:"foregroundColor"`
	Fields          []Field `json:"fields"`
}

type TicketUpdateDTO struct {
	Id              string  `json:"id"`
	UserId          string  `json:"userId"`
	Image           string  `json:"image"`
	Title           string  `json:"title"`
	Location        string  `json:"location"`
	Datetime        string  `json:"datetime"`
	BackgroundColor string  `json:"backgroundColor"`
	ForegroundColor string  `json:"foregroundColor"`
	Fields          []Field `json:"fields"`
}
