package usecase

import (
	"context"

	"github.com/doyeon0307/tickit-backend/domain"
	"github.com/doyeon0307/tickit-backend/dto"
	"github.com/doyeon0307/tickit-backend/models"
)

type scheduleUsecase struct {
	scheduleRepo domain.ScheduleRepository
}

func NewScheduleUsecase(repo domain.ScheduleRepository) domain.ScheduleUsecase {
	return &scheduleUsecase{
		scheduleRepo: repo,
	}
}

func (u scheduleUsecase) GetSchedulePreviewsForTicket(userId, date string) ([]*dto.ScheduleTicketPreviewDTO, error) {
	schedules, err := u.scheduleRepo.GetPreviewsForTicket(context.Background(), userId, date)
	if err != nil {
		return nil, err
	}

	previews := make([]*dto.ScheduleTicketPreviewDTO, len(schedules))
	for i, schedule := range schedules {
		previews[i] = &dto.ScheduleTicketPreviewDTO{
			Id:    schedule.Id,
			Title: schedule.Title,
			Date:  schedule.Date,
		}
	}

	return previews, nil
}

func (u scheduleUsecase) GetSchedulePreviewsForCalendar(userId, startDate, endDate string) ([]*dto.ScheduleCalendarPreviewDTO, error) {
	schedules, err := u.scheduleRepo.GetPreviewsForCalendar(context.Background(), userId, startDate, endDate)
	if err != nil {
		return nil, err
	}

	previews := make([]*dto.ScheduleCalendarPreviewDTO, len(schedules))
	for i, schedule := range schedules {
		previews[i] = &dto.ScheduleCalendarPreviewDTO{
			Id:    schedule.Id,
			Title: schedule.Title,
			Image: schedule.Image,
			Date:  schedule.Date,
		}
	}

	return previews, nil
}

func (u scheduleUsecase) GetScheduleById(userId, id string) (*dto.ScheduleResponseDTO, error) {
	model, err := u.scheduleRepo.GetById(context.Background(), userId, id)
	if err != nil {
		return nil, err
	}

	schedule := &dto.ScheduleResponseDTO{
		Id:        model.Id,
		Date:      model.Date,
		Title:     model.Title,
		Number:    model.Number,
		Image:     model.Image,
		Thumbnail: model.Thumbnail,
		Location:  model.Location,
		Time:      model.Time,
		Seat:      model.Seat,
		Casting:   model.Casting,
		Company:   model.Company,
		Link:      model.Link,
		Memo:      model.Memo,
	}

	return schedule, nil
}

func (u scheduleUsecase) CreateSchedule(userId string, schedule *dto.ScheduleDTO) (*dto.ScheduleResponseDTO, error) {
	model := &models.Schedule{
		UserId:    userId,
		Date:      schedule.Date,
		Title:     schedule.Title,
		Number:    schedule.Number,
		Image:     schedule.Image,
		Thumbnail: schedule.Thumbnail,
		Location:  schedule.Location,
		Time:      schedule.Time,
		Seat:      schedule.Seat,
		Casting:   schedule.Casting,
		Company:   schedule.Company,
		Link:      schedule.Link,
		Memo:      schedule.Memo,
	}

	id, err := u.scheduleRepo.Create(context.Background(), model)
	if err != nil {
		return &dto.ScheduleResponseDTO{}, err
	}

	result := &dto.ScheduleResponseDTO{
		Id:        id,
		Date:      schedule.Date,
		Title:     schedule.Title,
		Number:    schedule.Number,
		Image:     schedule.Image,
		Thumbnail: schedule.Thumbnail,
		Location:  schedule.Location,
		Time:      schedule.Time,
		Seat:      schedule.Seat,
		Casting:   schedule.Casting,
		Company:   schedule.Company,
		Link:      schedule.Link,
		Memo:      schedule.Memo,
	}
	return result, nil
}

func (u scheduleUsecase) UpdateSchedule(userId, id string, schedule *dto.ScheduleResponseDTO) (*dto.ScheduleResponseDTO, error) {
	model := &models.Schedule{
		UserId:    userId,
		Date:      schedule.Date,
		Title:     schedule.Title,
		Number:    schedule.Number,
		Image:     schedule.Image,
		Thumbnail: schedule.Thumbnail,
		Location:  schedule.Location,
		Time:      schedule.Time,
		Seat:      schedule.Seat,
		Casting:   schedule.Casting,
		Company:   schedule.Company,
		Link:      schedule.Link,
		Memo:      schedule.Memo,
	}

	err := u.scheduleRepo.Update(context.Background(), userId, id, model)
	if err != nil {
		return nil, err
	}

	result := &dto.ScheduleResponseDTO{
		Id:        id,
		Date:      schedule.Date,
		Title:     schedule.Title,
		Number:    schedule.Number,
		Image:     schedule.Image,
		Thumbnail: schedule.Thumbnail,
		Location:  schedule.Location,
		Time:      schedule.Time,
		Seat:      schedule.Seat,
		Casting:   schedule.Casting,
		Company:   schedule.Company,
		Link:      schedule.Link,
		Memo:      schedule.Memo,
	}
	return result, nil
}

func (u scheduleUsecase) DeleteSchedule(userId, id string) error {
	return u.scheduleRepo.Delete(context.Background(), userId, id)
}
