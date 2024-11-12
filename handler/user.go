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
// @Router /api/auth/login [post]
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
			"토큰 생성에 실패했습니다",
		))
		return
	}

	refreshToken, err := service.GenerateRefreshToken(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"토큰 생성에 실패했습니다",
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
// @Router /api/auth/register [post]
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

	refreshToken, err := service.GenerateRefreshToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error(
			http.StatusInternalServerError,
			"토큰 생성에 실패했습니다",
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
	userId, _ := c.Get("userId")

	err := h.userUsecase.DeleteUser(userId.(string))
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
			"회원탈퇴 중 오류가 발생했습니다",
		))
		return
	}

	// TODO: 토큰 블랙리스트에 추가하는 로직 필요

	c.JSON(http.StatusOK, common.Success(
		http.StatusOK,
		"회원탈퇴가 완료되었습니다",
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
	// token := c.GetHeader("Authorization")[7:] // "Bearer " 제거
	// err := service.BlacklistToken(token)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, common.Error(
	//         http.StatusInternalServerError,
	//         "로그아웃 처리 중 오류가 발생했습니다",
	//     ))
	//     return
	// }

	c.JSON(http.StatusOK, common.Success(
		http.StatusOK,
		"로그아웃 되었습니다",
		nil,
	))
}
