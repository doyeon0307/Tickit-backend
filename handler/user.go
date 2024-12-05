package handler

import (
	"net/http"

	"github.com/doyeon0307/tickit-backend/common"
	"github.com/doyeon0307/tickit-backend/domain"
	"github.com/doyeon0307/tickit-backend/dto"
	"github.com/doyeon0307/tickit-backend/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(rg *gin.RouterGroup, usecase domain.UserUsecase) {
	handler := &UserHandler{
		userUsecase: usecase,
	}

	users := rg.Group("/auth")
	{
		users.POST("/kakao/login", handler.Login)
		users.POST("/kakao/register", handler.Register)
		users.GET("/refresh", handler.RefreshToken)

		authorized := users.Group("")
		authorized.Use(service.AuthMiddleware())
		{
			authorized.DELETE("", handler.Withdraw)
			authorized.DELETE("/logout", handler.Logout)
			authorized.GET("", handler.GetProfile)
		}
	}
}

// @Tags Auth
// @Summary 로그인하기
// @Description 카카오 계정으로 로그인합니다
// @Accept json
// @Produce json
// @Param tokens body dto.KakaoTokens true "카카오 토큰"
// @Success 200 {object} common.Response{data=dto.TokenResponse} "성공"
// @Router /api/auth/kakao/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var tokens dto.KakaoTokens
	if err := c.ShouldBindJSON(&tokens); err != nil {
		c.JSON(http.StatusBadRequest, common.Error(
			http.StatusBadRequest,
			"Request Body가 올바르지 않습니다. accessToken, refreshToken, idToken이 포함되었는지 확인해주세요.",
		))
		return
	}

	oauthId, err := service.GetOAuthIdFromKakao(tokens.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.Error(
			http.StatusUnauthorized,
			"카카오 로그인 인증에 실패했습니다",
		))
		return
	}

	user, err := h.userUsecase.GetUserByOAuthId(oauthId)
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
			"가입되지 않은 사용자입니다. 회원가입이 필요합니다.",
		))
		return
	}

	id := user.Id
	accessToken, err := service.GenerateAccessToken(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"Access Token 생성에 실패했습니다",
		))
		return
	}

	refreshToken, expiryTime, err := service.GenerateRefreshToken(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"Refresh Token 생성에 실패했습니다",
		))
		return
	}

	if err := h.userUsecase.SaveRefreshToken(id, refreshToken, expiryTime); err != nil {
		if appErr, ok := err.(*common.AppError); ok {
			c.JSON(appErr.Code.StatusCode(), common.Error(
				appErr.Code.StatusCode(),
				appErr.Message,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"Refresh Token 저장에 오류가 발생했습니다",
		))
		return
	}

	resp := &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	c.JSON(http.StatusOK, common.Success(
		http.StatusOK,
		"로그인에 성공했습니다",
		resp,
	))
}

// @Tags Auth
// @Summary 회원가입하기
// @Description 카카오 계정으로 회원가입합니다
// @Accept json
// @Produce json
// @Param tokens body dto.KakaoTokens true "카카오 토큰"
// @Success 200 {object} common.Response{data=dto.TokenResponse}
// @Router /api/auth/kakao/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var tokens dto.KakaoTokens
	if err := c.ShouldBindJSON(&tokens); err != nil {
		c.JSON(http.StatusBadRequest, common.Error(
			http.StatusBadRequest,
			"Request Body가 올바르지 않습니다. accessToken과 idToken이 포함되었는지 확인해주세요.",
		))
		return
	}

	userId, err := h.userUsecase.CreateUser(tokens.IDToken, tokens.AccessToken)
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
			"회원가입 중 오류가 발생했습니다",
		))
		return
	}

	// 회원가입 성공 시 바로 JWT 토큰 발급
	accessToken, err := service.GenerateAccessToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"토큰 생성에 실패했습니다",
		))
		return
	}

	refreshToken, expiryTime, err := service.GenerateRefreshToken(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"Refresh Token 생성에 실패했습니다",
		))
		return
	}

	if err := h.userUsecase.SaveRefreshToken(userId, refreshToken, expiryTime); err != nil {
		if appErr, ok := err.(*common.AppError); ok {
			c.JSON(appErr.Code.StatusCode(), common.Error(
				appErr.Code.StatusCode(),
				appErr.Message,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"Refresh Token 저장에 오류가 발생했습니다",
		))
		return
	}

	resp := &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	c.JSON(http.StatusCreated, common.Success(
		http.StatusCreated,
		"회원가입에 성공했습니다",
		resp,
	))
}

// @Security ApiKeyAuth
// @Tags Auth
// @Summary 프로필 가져오기
// @Description 사용자 프로필 정보를 가져옵니다
// @Accept json
// @Produce json
// @Success 200 {object} common.Response{data=dto.KakaoProfile}
// @Router /api/auth [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userId, _ := c.Get("userId")

	profile, err := h.userUsecase.GetProfile(userId.(string))
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
			"프로필 조회 중 오류가 발생했습니다",
		))
		return
	}

	c.JSON(http.StatusOK, common.Success(
		http.StatusOK,
		"프로필 조회에 성공했습니다",
		profile,
	))
}

// @Security ApiKeyAuth
// @Tags Auth
// @Summary 탈퇴하기
// @Description 계정을 삭제합니다
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Router /api/auth [delete]
func (h *UserHandler) Withdraw(c *gin.Context) {
	userIdInterface, _ := c.Get("userId")
	userId, ok := userIdInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"사용자 ID 타입이 올바르지 않습니다",
		))
		return
	}

	if err := h.userUsecase.WithdrawUser(userId); err != nil {
		if appErr, ok := err.(*common.AppError); ok {
			c.JSON(appErr.Code.StatusCode(), common.Error(
				appErr.Code.StatusCode(),
				appErr.Message,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"회원 탈퇴 처리 중 오류가 발생했습니다",
		))
		return
	}

	c.JSON(http.StatusOK, common.Success(
		http.StatusOK,
		"회원 탈퇴가 완료되었습니다",
		nil,
	))
}

// @Security ApiKeyAuth
// @Tags Auth
// @Summary 로그아웃하기
// @Description 로그아웃합니다
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Router /api/auth/logout [delete]
func (h *UserHandler) Logout(c *gin.Context) {
	userIdInterface, _ := c.Get("userId")
	userId, ok := userIdInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"사용자 ID 타입이 올바르지 않습니다",
		))
		return
	}

	if err := h.userUsecase.Logout(userId); err != nil {
		if appErr, ok := err.(*common.AppError); ok {
			c.JSON(appErr.Code.StatusCode(), common.Error(
				appErr.Code.StatusCode(),
				appErr.Message,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"로그아웃 처리 중 오류가 발생했습니다",
		))
		return
	}

	c.JSON(http.StatusOK, common.Success(
		http.StatusOK,
		"로그아웃이 완료되었습니다",
		nil,
	))
}

// @Tags Auth
// @Summary Access Token 갱신하기
// @Description Refresh Token으로 새로운 Access Token을 합니다
// @Accept json
// @Produce json
// @Param tokens body dto.RefreshTokenRequest true "Refresh Token"
// @Success 200 {object} common.Response{data=dto.TokenResponse}
// @Router /api/auth/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest

	userIdInterface, _ := c.Get("userId")
	userId, ok := userIdInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"사용자 ID 타입이 올바르지 않습니다",
		))
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Error(
			http.StatusBadRequest,
			"Request Body가 올바르지 않습니다. refreshToken이 포함되었는지 확인해주세요.",
		))
		return
	}

	isValid, err := h.userUsecase.ValidateStoredRefreshToken(userId, req.RefreshToken)
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
			"Refresh Token 검증 도중 오류가 발생했습니다",
		))
		return
	}

	if !isValid {
		c.JSON(http.StatusUnauthorized, common.Error(
			http.StatusUnauthorized,
			"저장된 Refresh Token과 일치하지 않습니다",
		))
		return
	}

	accessToken, err := service.GenerateAccessToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"새로운 Access Token 생성에 실패했습니다",
		))
		return
	}

	resp := &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
	}
	c.JSON(http.StatusOK, common.Success(
		http.StatusOK,
		"토큰이 성공적으로 갱신되었습니다",
		resp,
	))
}
