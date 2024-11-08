package domain

import (
	"github.com/doyeon0307/tickit-backend/dto"
)

type ScheduleUsecase interface {
	GetSchedulePreviewsForTicket(date string) ([]*dto.ScheduleTicketPreviewDTO, error)
	GetSchedulePreviewsForCalendar(startDate, endDate string) ([]*dto.ScheduleCalendarPreviewDTO, error)
	GetScheduleById(id string) (*dto.ScheduleResponseDTO, error)
	CreateSchedule(schedule *dto.ScheduleDTO) (*dto.ScheduleResponseDTO, error)
	UpdateSchedule(id string, schedule *dto.ScheduleResponseDTO) (*dto.ScheduleResponseDTO, error)
	DeleteSchedule(id string) error
}
