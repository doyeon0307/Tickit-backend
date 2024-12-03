package handler

import (
	"net/http"
	"regexp"
	"time"

	"github.com/doyeon0307/tickit-backend/common"
	"github.com/doyeon0307/tickit-backend/domain"
	"github.com/doyeon0307/tickit-backend/dto"
	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	scheduleUsecase domain.ScheduleUsecase
}

func NewScheduleHandler(rg *gin.RouterGroup, usecase domain.ScheduleUsecase) {
	handler := &ScheduleHandler{
		scheduleUsecase: usecase,
	}
	schedules := rg.Group("/schedules")
	{
		schedules.GET("/for-ticket", handler.GetSchedulePreviewsForTicket)
		schedules.GET("", handler.GetSchedulePreviewsForCalendar)
		schedules.GET("/:id", handler.GetScheduleById)
		schedules.POST("", handler.CreateSchedule)
		schedules.PUT("/:id", handler.UpdateSchedule)
		schedules.DELETE("/:id", handler.DeleteSchedule)
	}
}

// @Security ApiKeyAuth
// @Tags Schedules
// @Summary 티켓 생성 가능한 일정 목록 불러오기
// @Description 현 날짜 이전의 일정 목록을 불러옵니다. 티켓 생성 화면의'일정 불러오기' 버튼에서 사용됩니다.
// @Accept json
// @Produce json
// @Param date query string false "오늘 날짜"
// @Success 200 {object} common.Response{data=dto.ScheduleTicketPreviewDTO}
// @Router /api/schedules/for-ticket [get]
func (h *ScheduleHandler) GetSchedulePreviewsForTicket(c *gin.Context) {
	userId, _ := c.Get("userId")
	date := c.Query("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	} else {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Error(
				http.StatusBadRequest,
				"날짜 형식이 잘못되었습니다. YYYY-MM-DD 형식으로 입력해주세요.",
			))
			return
		}
	}
	previews, err := h.scheduleUsecase.GetSchedulePreviewsForTicket(userId.(string), date)
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
			"일정 불러오기에 실패했습니다",
		))
		return
	}
	c.JSON(http.StatusOK, common.Success(
		http.StatusOK,
		"티켓으로 만들 수 있는 일정 목록 불러오기에 성공했습니다",
		previews,
	))
}

// @Security ApiKeyAuth
// @Tags Schedules
// @Summary 달력에 일정 목록 불러오기
// @Description 시작 날짜와 종료 날짜 사이의 일정 목록을 불러옵니다
// @Accept json
// @Produce json
// @Param startDate query string true "시작 날짜"
// @Param endDate query string true "종료 날짜"
// @Success 200 {object} common.Response{data=dto.ScheduleCalendarPreviewDTO}
// @Router /api/schedules [get]
func (h *ScheduleHandler) GetSchedulePreviewsForCalendar(c *gin.Context) {
	userId, _ := c.Get("userId")

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		c.JSON(http.StatusBadRequest, common.Error(
			http.StatusBadRequest,
			"시작 날짜 형식이 잘못되었습니다. YYYY-MM-DD 형식으로 입력해주세요.",
		))
		return
	}

	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		c.JSON(http.StatusBadRequest, common.Error(
			http.StatusBadRequest,
			"종료 날짜 형식이 잘못되었습니다. YYYY-MM-DD 형식으로 입력해주세요.",
		))
		return
	}

	previews, err := h.scheduleUsecase.GetSchedulePreviewsForCalendar(userId.(string), startDate, endDate)
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
			"일정 불러오기에 실패했습니다",
		))
		return
	}
	c.JSON(http.StatusOK, common.Success(
		http.StatusOK,
		"일정 목록 불러오기에 성공했습니다",
		previews,
	))
}

// @Security ApiKeyAuth
// @Tags Schedules
// @Summary 세부 일정 불러오기
// @Description 세부 일정을 불러옵니다
// @Accept json
// @Produce json
// @Param id path string true "일정 ID"
// @Success 200 {object} common.Response{data=dto.ScheduleResponseDTO}
// @Router /api/schedules/{id} [get]
func (h *ScheduleHandler) GetScheduleById(c *gin.Context) {
	userId, _ := c.Get("userId")

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, common.Error(
			http.StatusBadRequest,
			"아이디를 입력해주세요",
		))
	}

	schedule, err := h.scheduleUsecase.GetScheduleById(userId.(string), id)
	if err != nil {
		if appErr, ok := err.(*common.AppError); ok {
			c.JSON(appErr.Code.StatusCode(), common.Error(
				appErr.Code.StatusCode(),
				appErr.Message,
			))
			return
		}
		c.JSON(http.StatusNotFound, common.Error(
			http.StatusNotFound,
			"일정 조회에 실패했습니다. 아이디를 확인해주세요.",
		))
		return
	}
	c.JSON(http.StatusAccepted, common.Success(
		http.StatusAccepted,
		"세부 일정 조회에 성공했습니다",
		schedule,
	))
}

// @Security ApiKeyAuth
// @Tags Schedules
// @Summary 일정 생성하기
// @Description 일정을 생성합니다. presigned-url을 발급받아 이미지 업로드를 완료한 후에, s3 url을 image 값으로 저장합니다. 날짜 형식은 YYYY-MM-DD, 시간 형식은 AM/PM-HH-MM입니다.
// @Accept json
// @Produce json
// @Param scheduleDTO body dto.ScheduleDTO true "일정 DTO"
// @Success 200 {object} common.Response{data=dto.ScheduleResponseDTO}
// @Router /api/schedules [post]
func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	userId, _ := c.Get("userId")

	var schedule dto.ScheduleDTO
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, common.Error(
			http.StatusBadRequest,
			"Request Body가 올바르지 않습니다",
		))
		return
	}

	if _, err := time.Parse("2006-01-02", schedule.Date); err != nil {
		c.JSON(http.StatusBadRequest, common.Error(
			http.StatusBadRequest,
			"날짜 형식이 잘못되었습니다. YYYY-MM-DD 형식으로 입력해주세요.",
		))
		return
	}

	timePattern := `^(AM|PM)-(?:0[1-9]|1[0-2])-(?:[0-5][0-9])$`
	if !regexp.MustCompile(timePattern).MatchString(schedule.Time) {
		c.JSON(http.StatusBadRequest, common.Error(
			http.StatusBadRequest,
			"시간 형식이 잘못되었습니다. AM/PM-HH-MM 형식으로 입력해주세요.",
		))
		return
	}

	resp, err := h.scheduleUsecase.CreateSchedule(userId.(string), &schedule)
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
			"일정 생성에 실패했습니다",
		))
		return
	}
	c.JSON(http.StatusAccepted, common.Success(
		http.StatusAccepted,
		"일정 생성에 성공했습니다",
		resp,
	))
}

// @Security ApiKeyAuth
// @Tags Schedules
// @Summary 일정 수정하기
// @Description 일정을 수정합니다. presigned-url을 발급받아 이미지 업로드를 완료한 후에, s3 url을 image 값으로 저장합니다.
// @Accept json
// @Produce json
// @Param id path string true "일정 ID"
// @Param scheduleDTO body dto.ScheduleDTO true "일정 DTO"
// @Success 200 {object} common.Response{data=dto.ScheduleResponseDTO}
// @Router /api/schedules/{id} [put]
func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	userId, _ := c.Get("userId")

	id := c.Param("id")
	var schedule dto.ScheduleResponseDTO

	if id == "" {
		c.JSON(http.StatusBadRequest, common.Error(
			http.StatusBadRequest,
			"아이디를 입력해주세요",
		))
	}

	if _, err := h.scheduleUsecase.GetScheduleById(userId.(string), id); err != nil {
		if appErr, ok := err.(*common.AppError); ok {
			c.JSON(appErr.Code.StatusCode(), common.Error(
				appErr.Code.StatusCode(),
				appErr.Message,
			))
			return
		}
		c.JSON(http.StatusNotFound, common.Error(
			http.StatusNotFound,
			"일정 조회에 실패했습니다. 아이디를 확인해주세요.",
		))
		return
	}

	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, common.Error(
			http.StatusBadRequest,
			"Request Body가 올바르지 않습니다",
		))
		return
	}

	resp, err := h.scheduleUsecase.UpdateSchedule(userId.(string), id, &schedule)
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
			"일정 수정에 실패했습니다",
		))
		return
	}
	c.JSON(http.StatusAccepted, common.Success(
		http.StatusAccepted,
		"일정이 수정되었습니다",
		resp,
	))
}

// @Security ApiKeyAuth
// @Tags Schedules
// @Summary 일정 삭제하기
// @Description 일정을 삭제합니다
// @Accept json
// @Produce json
// @Param id path string true "일정 ID"
// @Success 200 {object} common.Response
// @Router /api/schedules/{id} [delete]
func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	userId, _ := c.Get("userId")

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, common.Error(
			http.StatusBadRequest,
			"아이디를 입력해주세요",
		))
	}

	if _, err := h.scheduleUsecase.GetScheduleById(userId.(string), id); err != nil {
		if appErr, ok := err.(*common.AppError); ok {
			c.JSON(appErr.Code.StatusCode(), common.Error(
				appErr.Code.StatusCode(),
				appErr.Message,
			))
			return
		}
		c.JSON(http.StatusNotFound, common.Error(
			http.StatusNotFound,
			"일정 조회에 실패했습니다. 아이디를 확인해주세요.",
		))
		return
	}

	err := h.scheduleUsecase.DeleteSchedule(userId.(string), id)
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
			"일정 삭제에 실패했습니다",
		))
		return
	}
	c.JSON(http.StatusAccepted, common.Success(
		http.StatusAccepted,
		"일정이 삭제되었습니다",
		id,
	))
}
