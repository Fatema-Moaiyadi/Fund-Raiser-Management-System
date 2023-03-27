package service

import (
	"errors"
	"github.com/fatema-moaiyadi/fund-raiser-system/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type TokenService interface {
	GenerateToken(userId int64, isAdmin bool) (string, error)
	VerifyToken(jwtToken string) (*TokenPayload, error)
}

type TokenPayload struct {
	UserID   int64
	IsAdmin  bool
	IssuedAt time.Time
	*jwt.RegisteredClaims
}

type jwtTokenService struct {
}

func NewJWTTokenService() TokenService {
	return &jwtTokenService{}
}

func (jt *jwtTokenService) GenerateToken(userId int64, isAdmin bool) (string, error) {
	payload := &TokenPayload{
		UserID:   userId,
		IsAdmin:  isAdmin,
		IssuedAt: time.Now(),
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(config.GetJWTConfig().Duration))),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(config.GetJWTConfig().Secret))
}

func (jt *jwtTokenService) VerifyToken(jwtToken string) (*TokenPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Invalid token")
		}

		return []byte(config.GetJWTConfig().Secret), nil
	}
	token, err := jwt.ParseWithClaims(jwtToken, &TokenPayload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	return token.Claims.(*TokenPayload), nil
}
