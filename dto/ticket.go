package dto

import "github.com/doyeon0307/tickit-backend/models"

type Field struct {
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
}

type TicketDTO struct {
	Image           string  `json:"image" binding:"required"`
	Title           string  `json:"title" binding:"required"`
	Location        string  `json:"location" binding:"required"`
	Date            string  `json:"date" binding:"required"`
	Time            string  `json:"time" binding:"required"`
	BackgroundColor string  `json:"backgroundColor"`
	ForegroundColor string  `json:"foregroundColor"`
	Fields          []Field `json:"fields"`
}

type TicketResponseDTO struct {
	Id              string         `json:"id"`
	Image           string         `json:"image"`
	Title           string         `json:"title"`
	Location        string         `json:"location"`
	Date            string         `json:"date" binding:"required"`
	Time            string         `json:"time" binding:"required"`
	BackgroundColor string         `json:"backgroundColor"`
	ForegroundColor string         `json:"foregroundColor"`
	Fields          []models.Field `json:"fields"`
}

type TicketUpdateDTO struct {
	Id              string         `json:"id"`
	Image           string         `json:"image"`
	Title           string         `json:"title"`
	Location        string         `json:"location"`
	Date            string         `json:"date" binding:"required"`
	Time            string         `json:"time" binding:"required"`
	BackgroundColor string         `json:"backgroundColor"`
	ForegroundColor string         `json:"foregroundColor"`
	Fields          []models.Field `json:"fields"`
}

type TicketPreview struct {
	Id    string `json:"id"`
	Image string `json:"image"`
}
