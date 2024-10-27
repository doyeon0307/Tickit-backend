package handler

import (
	"net/http"
	"std/github.com/dodo/Tickit-backend/common"
	"std/github.com/dodo/Tickit-backend/domain"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketUseCase domain.TicketUsecase
}

func NewTicketHandler(r *gin.Engine, usecase domain.TicketUsecase) {
	handler := &TicketHandler{
		ticketUseCase: usecase,
	}

	r.GET("/ticket", handler.GetTicketPreviews)
	r.GET("/ticket/:id")
	r.POST("/ticket")
	r.PUT("/ticket/:id")
	r.DELETE("/ticket/:id")
}

func (h *TicketHandler) GetTicketPreviews(c *gin.Context) {
	previews, err := h.ticketUseCase.GetTicketPreviews()
	if err != nil {
		var appErr common.AppError
		switch appErr.Code {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, response.Error(
				http.StatusNotFound,
				appErr.Message,
			))
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, previews)
}
