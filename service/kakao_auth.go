package service

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/doyeon0307/tickit-backend/common"
	"github.com/doyeon0307/tickit-backend/dto"
)

type Payload struct {
	Aud       string `json:"aud"`
	Sub       string `json:"sub"`
	Auth_time string `json:"auth_time"`
	Iss       string `json:"iss"`
	Exp       string `json:"exp"`
	Iat       string `json:"iat"`
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

	var claims Payload
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

func GetUserInfoFromKakao(accessToken string) (*dto.KakaoProfile, error) {
	req, err := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "카카오 서버에서 사용자 정보를 불러오는데 실패했습니다",
			Err:     err,
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var profile dto.KakaoProfile
	convertErr := json.Unmarshal(body, &profile)
	if convertErr != nil {
		return nil, &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "사용자 정보 처리에 실패했습니다",
			Err:     err,
		}
	}

	return &profile, nil
}
