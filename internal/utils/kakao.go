package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type KakaoAccount struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type KakaoResponse struct {
	Id           uint64       `json:"id"`
	KakaoAccount KakaoAccount `json:"kakao_account"`
}

func GetUserInfoFromKakao(accessToken string) (KakaoResponse, error) {
	url := "https://kapi.kakao.com/v2/user/me"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, _ := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)

	var kakaoResponse KakaoResponse

	err := json.NewDecoder(res.Body).Decode(&kakaoResponse)
	if err != nil {
		return KakaoResponse{}, err
	}
	return kakaoResponse, nil
}
