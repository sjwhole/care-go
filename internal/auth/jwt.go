package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtManager struct {
	secretKey []byte
	validMin  int
}

func (j *JwtManager) IssueJwtToken(userId int) (string, error) {
	var (
		token  *jwt.Token
		signed string
	)

	token = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userId,
			"exp": time.Now().Add(time.Duration(j.validMin) * time.Minute).Unix(),
		},
	)
	signed, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (j *JwtManager) VerifyJwtToken(token string) (int, error) {
	var (
		claims jwt.MapClaims
		err    error
		userId int
	)

	claims = jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil {
		return 0, err
	}

	userId = int(claims["sub"].(float64))
	return userId, nil
}

func NewJwtManager(secretKey string, validMin int) *JwtManager {
	return &JwtManager{
		secretKey: []byte(secretKey),
		validMin:  validMin,
	}
}
