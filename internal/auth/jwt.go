package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
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

func (j *JwtManager) VerifyJwtToken(token string) (uint, error) {
	var (
		claims jwt.MapClaims
		err    error
		userId uint
	)

	claims = jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil {
		return 0, err
	}

	_, ok := claims["sub"].(float64)
	if !ok {
		return 0, fmt.Errorf("got data of type %T but expectes float64 type for userId", userId)
	}
	userId = uint(claims["sub"].(float64))

	return userId, nil
}

func NewJwtManager() *JwtManager {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("SECRET_KEY")
	validMin, err := strconv.Atoi(os.Getenv("VALID_MIN"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &JwtManager{
		secretKey: []byte(secretKey),
		validMin:  validMin,
	}
}
