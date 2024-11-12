package handler

import (
	"net/http"

	"github.com/doyeon0307/tickit-backend/common"
	"github.com/doyeon0307/tickit-backend/domain"
	"github.com/doyeon0307/tickit-backend/dto"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketUsecase domain.TicketUsecase
}

func NewTicketHandler(rg *gin.RouterGroup, usecase domain.TicketUsecase) {
	handler := &TicketHandler{
		ticketUsecase: usecase,
	}
	tickets := rg.Group("/tickets")
	{
		tickets.GET("", handler.GetTicketPreviews)
		tickets.GET("/:id", handler.GetTicketById)
		tickets.POST("", handler.MakeTicket)
		tickets.PUT("/:id", handler.UpdateTicket)
		tickets.DELETE("/:id", handler.DeleteTicket)
	}
}

// @Security ApiKeyAuth
// @Tags Tickets
// @Summary 티켓 목록 불러오기
// @Description 홈 화면에 작성한 티켓 목록을 불러옵니다
// @Accept json
// @Produce json
// @Success 200 {object} common.Response{data=dto.TicketPreview}
// @Router /api/tickets [get]
func (h *TicketHandler) GetTicketPreviews(c *gin.Context) {
	userId, _ := c.Get("userId")

	previews, err := h.ticketUsecase.GetTicketPreviews(userId.(string))
	if err != nil {
		if appErr, ok := err.(*common.AppError); ok {
			c.JSON(appErr.Code.StatusCode(), common.Error(
				appErr.Code.StatusCode(),
				appErr.Message,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"티켓 목록 불러오기에 실패했습니다",
		))
		return
	}
	c.JSON(http.StatusOK, common.Success(
		http.StatusOK,
		"티켓 목록 불러오기에 성공했습니다",
		previews,
	))
}

// @Security ApiKeyAuth
// @Tags Tickets
// @Summary 티켓 세부정보 불러오기
// @Description 티켓 아이디로 세부정보를 불러옵니다
// @Accept json
// @Produce json
// @Param id path string true "티켓 ID"
// @Success 200 {object} common.Response{data=dto.TicketDTO}
// @Router /api/tickets/{id} [get]
func (h *TicketHandler) GetTicketById(c *gin.Context) {
	userId, _ := c.Get("userId")

	id := c.Param("id")
	ticket, err := h.ticketUsecase.GetTicketByID(userId.(string), id)

	if err != nil {
		if appErr, ok := err.(*common.AppError); ok {
			c.JSON(appErr.Code.StatusCode(), common.Error(
				appErr.Code.StatusCode(),
				appErr.Message,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"티켓 불러오기에 실패했습니다",
		))
		return
	}
	c.JSON(http.StatusOK, common.Success(
		http.StatusOK,
		"티켓 불러오기에 성공했습니다",
		ticket,
	))
}

// @Security ApiKeyAuth
// @Tags Tickets
// @Summary 티켓 생성하기
// @Description 티켓을 생성합니다. presigned-url을 발급받아 이미지 업로드를 완료한 후에, s3 url을 image 값으로 저장합니다.
// @Accept json
// @Produce json
// @Param ticketDTO body dto.TicketDTO true "생성할 티켓 DTO"
// @Success 200 {object} common.Response{data=dto.TicketDTO}
// @Router /api/tickets [post]
func (h *TicketHandler) MakeTicket(c *gin.Context) {
	userId, _ := c.Get("userId")

	var req dto.TicketDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Error(
			http.StatusBadRequest,
			"Request Body가 잘못되었습니다",
		))
		return
	}

	if req.BackgroundColor == "" {
		req.BackgroundColor = "0xffFFFF"
	}
	if req.ForegroundColor == "" {
		req.ForegroundColor = "0xff000000"
	}

	ticket, err := h.ticketUsecase.CreateTicket(userId.(string), &req)
	if err != nil {
		if appErr, ok := err.(*common.AppError); ok {
			c.JSON(appErr.Code.StatusCode(), common.Error(
				appErr.Code.StatusCode(),
				appErr.Message,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"티켓 생성에 실패했습니다",
		))
		return
	}

	c.JSON(http.StatusCreated, common.Success(
		http.StatusCreated,
		"티켓이 생성되었습니다",
		ticket,
	))
}

// @Security ApiKeyAuth
// @Tags Tickets
// @Summary 티켓 수정하기
// @Description 티켓을 수정합니다. presigned-url을 발급받아 이미지 업로드를 완료한 후에, s3 url을 image 값으로 저장합니다.
// @Accept json
// @Produce json
// @Param id path string true "티켓 ID"
// @Param ticketDTO body dto.TicketUpdateDTO true "수정된 티켓 DTO"
// @Success 200 {object} common.Response{data=dto.TicketDTO}
// @Router /api/tickets/{id} [put]
func (h *TicketHandler) UpdateTicket(c *gin.Context) {
	userId, _ := c.Get("userId")
	id := c.Param("id")

	// 기존 티켓 존재 여부 확인
	if _, err := h.ticketUsecase.GetTicketByID(userId.(string), id); err != nil {
		if appErr, ok := err.(*common.AppError); ok {
			c.JSON(appErr.Code.StatusCode(), common.Error(
				appErr.Code.StatusCode(),
				appErr.Message,
			))
			return
		}
		c.JSON(http.StatusNotFound, common.Error(
			http.StatusNotFound,
			"티켓 조회에 실패했습니다. 아이디를 확인해주세요.",
		))
		return
	}

	var req dto.TicketUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Error(
			http.StatusBadRequest,
			"Request Body가 올바르지 않습니다",
		))
		return
	}

	req.Id = id
	if err := h.ticketUsecase.UpdateTicket(userId.(string), id, &req); err != nil {
		if appErr, ok := err.(*common.AppError); ok {
			c.JSON(appErr.Code.StatusCode(), common.Error(
				appErr.Code.StatusCode(),
				appErr.Message,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"티켓 수정에 실패했습니다",
		))
		return
	}

	// 수정된 티켓 정보 조회
	updatedTicket, err := h.ticketUsecase.GetTicketByID(userId.(string), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"수정된 티켓 정보 조회에 실패했습니다",
		))
		return
	}

	c.JSON(http.StatusAccepted, common.Success(
		http.StatusAccepted,
		"티켓이 수정되었습니다",
		updatedTicket,
	))
}

// @Security ApiKeyAuth
// @Tags Tickets
// @Summary 티켓 삭제하기
// @Description 티켓을 삭제합니다
// @Accept json
// @Produce json
// @Param id path string true "티켓 ID"
// @Success 200 {object} common.Response
// @Router /api/tickets/{id} [delete]
func (h *TicketHandler) DeleteTicket(c *gin.Context) {
	userId, _ := c.Get("userId")

	id := c.Param("id")

	if _, err := h.ticketUsecase.GetTicketByID(userId.(string), id); err != nil {
		if appErr, ok := err.(*common.AppError); ok {
			c.JSON(appErr.Code.StatusCode(), common.Error(
				appErr.Code.StatusCode(),
				appErr.Message,
			))
			return
		}
		c.JSON(http.StatusNotFound, common.Error(
			http.StatusNotFound,
			"티켓 조회에 실패했습니다. 아이디를 확인해주세요.",
		))
		return
	}

	err := h.ticketUsecase.DeleteTicket(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"티켓 삭제에 실패했습니다",
		))
		return
	}
	c.JSON(http.StatusOK, common.Success(
		http.StatusOK,
		"티켓이 삭제되었습니다",
		id,
	))
}
