package domain

import (
	"github.com/doyeon0307/tickit-backend/dto"
)

type TicketUsecase interface {
	GetTicketPreviews(userId string) ([]*dto.TicketPreview, error)
	GetTicketByID(userId, id string) (*dto.TicketResponseDTO, error)
	CreateTicket(userId string, ticket *dto.TicketDTO) (string, error)
	UpdateTicket(userId, id string, ticket *dto.TicketUpdateDTO) error
	DeleteTicket(id string) error
}
