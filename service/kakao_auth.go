package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/doyeon0307/tickit-backend/common"
	"github.com/doyeon0307/tickit-backend/dto"
)

// KakaoUserResponse는 실제 카카오 API 응답 구조체입니다
type KakaoUserResponse struct {
	NickName        string `json:"nickName"`
	ProfileImageURL string `json:"profileImageURL"`
	ThumbnailURL    string `json:"thumbnailURL"`
}

func GetUserInfoFromKakao(accessToken string) (*dto.KakaoProfile, error) {
	req, err := http.NewRequest("GET", "https://kapi.kakao.com/v1/api/talk/profile", nil)
	if err != nil {
		return nil, &common.AppError{
			Code:    common.ErrServer,
			Message: "Request 생성 실패",
			Err:     err,
		}
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &common.AppError{
			Code:    common.ErrServer,
			Message: "HTTP 요청 실패",
			Err:     err,
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &common.AppError{
			Code:    common.ErrServer,
			Message: "응답 바디 읽기 실패",
			Err:     err,
		}
	}

	// Debug: 응답 로깅
	fmt.Printf("Kakao API Response Status: %d\n", resp.StatusCode)
	fmt.Printf("Kakao API Response Body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, &common.AppError{
			Code:    common.ErrBadRequest,
			Message: fmt.Sprintf("카카오 API 오류 (상태 코드: %d): %s", resp.StatusCode, string(body)),
			Err:     fmt.Errorf("카카오 API 응답: %s", string(body)),
		}
	}

	var kakaoResp KakaoUserResponse
	if err := json.Unmarshal(body, &kakaoResp); err != nil {
		return nil, &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "카카오 응답 파싱 실패",
			Err:     err,
		}
	}

	if kakaoResp.NickName == "" {
		return nil, &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "사용자 닉네임을 찾을 수 없습니다",
			Err:     fmt.Errorf("empty nickname in response"),
		}
	}

	return &dto.KakaoProfile{
		NickName: kakaoResp.NickName,
	}, nil
}

func GetOAuthIdFromKakao(idToken string) (string, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return "", &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "ID Token의 형식이 잘못되었습니다",
		}
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "ID Token 디코딩에 실패했습니다",
			Err:     err,
		}
	}

	var claims struct {
		Sub string `json:"sub"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "ID Token 처리에 실패했습니다",
			Err:     err,
		}
	}

	if claims.Sub == "" {
		return "", &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "ID Token에서 OAuthId를 찾을 수 없습니다",
		}
	}

	return claims.Sub, nil
}
