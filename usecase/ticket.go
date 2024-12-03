package usecase

import (
	"context"
	"time"

	"github.com/doyeon0307/tickit-backend/domain"
	"github.com/doyeon0307/tickit-backend/dto"
	"github.com/doyeon0307/tickit-backend/models"
	"github.com/doyeon0307/tickit-backend/utils"
)

type ticketUsecase struct {
	ticketRepo domain.TicketRepository
}

func NewTicketUseCase(repo domain.TicketRepository) domain.TicketUsecase {
	return &ticketUsecase{
		ticketRepo: repo,
	}
}

func (u ticketUsecase) GetTicketPreviews(userId string) ([]*dto.TicketPreview, error) {
	models, err := u.ticketRepo.GetPreviews(context.Background(), userId)
	if err != nil {
		return nil, err
	}

	previews := make([]*dto.TicketPreview, len(models))
	for i, model := range models {
		previews[i] = &dto.TicketPreview{
			Id:    model.Id,
			Image: model.Image,
		}
	}
	return previews, nil
}

func (u ticketUsecase) GetTicketByID(userId, id string) (*dto.TicketResponseDTO, error) {
	model, err := u.ticketRepo.GetById(context.Background(), userId, id)
	if err != nil {
		return nil, err
	}

	date, time := utils.SplitDateTime(model.DateTime)

	ticket := &dto.TicketResponseDTO{
		Id:              model.Id,
		Image:           model.Image,
		Title:           model.Title,
		Location:        model.Location,
		Date:            date,
		Time:            time,
		BackgroundColor: model.BackgroundColor,
		ForegroundColor: model.ForegroundColor,
		Fields:          model.Fields,
	}
	return ticket, nil
}

func (u ticketUsecase) CreateTicket(userId string, ticket *dto.TicketDTO) (string, error) {
	fields := make([]models.Field, len(ticket.Fields))
	for i, f := range ticket.Fields {
		fields[i] = models.Field{
			Subtitle: f.Subtitle,
			Content:  f.Content,
		}
	}

	dateTime, err := utils.CombineDateTime(ticket.Date, ticket.Time)

	if err != nil {
		return "", err
	}

	model := &models.Ticket{
		UserId:          userId,
		Image:           ticket.Image,
		Title:           ticket.Title,
		Location:        ticket.Location,
		DateTime:        dateTime,
		BackgroundColor: ticket.BackgroundColor,
		ForegroundColor: ticket.ForegroundColor,
		Fields:          fields,
		CreatedAt:       time.Now(),
	}

	id, err := u.ticketRepo.Create(context.Background(), userId, model)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (u ticketUsecase) UpdateTicket(userId, id string, ticket *dto.TicketUpdateDTO) error {
	dateTime, err := utils.CombineDateTime(ticket.Date, ticket.Time)

	if err != nil {
		return err
	}

	model := &models.Ticket{
		UserId:          userId,
		Image:           ticket.Image,
		Title:           ticket.Title,
		Location:        ticket.Location,
		DateTime:        dateTime,
		BackgroundColor: ticket.BackgroundColor,
		ForegroundColor: ticket.ForegroundColor,
		Fields:          ticket.Fields,
	}
	return u.ticketRepo.Update(context.Background(), userId, id, model)
}

func (u ticketUsecase) DeleteTicket(id string) error {
	return u.ticketRepo.Delete(context.Background(), id)
}
