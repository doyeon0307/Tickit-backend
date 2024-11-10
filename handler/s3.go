package handler

import (
	"net/http"

	"github.com/doyeon0307/tickit-backend/common"
	"github.com/doyeon0307/tickit-backend/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type S3Handler struct {
	s3Config *config.S3Config
}

func NewS3Handler(rg *gin.RouterGroup, s3Config *config.S3Config) {
	handler := &S3Handler{
		s3Config: s3Config,
	}
	s3 := rg.Group("/s3")
	{
		s3.GET("/presigned-url", handler.GetPresignedUrl)
	}
}

// @Security ApiKeyAuth
// @Tags S3
// @Summary Presigend URL 불러오기
// @Description Presigend URL를 얻고, 해당 URL을 통해 S3 이미지 업로드를 수행합니다
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Router /api/s3/presigned-url [get]
func (h *S3Handler) GetPresignedUrl(c *gin.Context) {
	key := uuid.New().String()

	url, err := h.s3Config.MakePresignURL(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"URL 생성에 실패했습니다",
		))
		return
	}
	c.JSON(http.StatusOK, common.Success(
		http.StatusOK,
		"URL 생성에 성공했습니다",
		url,
	))
}
