package domain

import (
	"github.com/doyeon0307/tickit-backend/dto"
)

type ScheduleUsecase interface {
	GetSchedulePreviewsForTicket(userId, date string) ([]*dto.ScheduleTicketPreviewDTO, error)
	GetSchedulePreviewsForCalendar(userId, startDate, endDate string) ([]*dto.ScheduleCalendarPreviewDTO, error)
	GetScheduleById(userId, id string) (*dto.ScheduleResponseDTO, error)
	CreateSchedule(userId string, schedule *dto.ScheduleDTO) (*dto.ScheduleResponseDTO, error)
	UpdateSchedule(userId, id string, schedule *dto.ScheduleResponseDTO) (*dto.ScheduleResponseDTO, error)
	DeleteSchedule(userId, id string) error
}
