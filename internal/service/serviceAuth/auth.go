package serviceAuth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sensors-app/internal/entities"

	"time"

	"github.com/golang-jwt/jwt"
)

type TokenRepo interface {
	StoreToken(cxt context.Context, userId int64, token string, tokenTTL time.Duration) error
	DeleteToken(cxt context.Context, userId int64) error
	TokenExists(cxt context.Context, userId int64) (bool, error)
}

type AuthService struct {
	tokenRepo TokenRepo
}

func NewAuthService(tokenRepo TokenRepo) AuthService {
	return AuthService{
		tokenRepo: tokenRepo,
	}
}

type tokenClaims struct {
	jwt.StandardClaims
	userId int64
}

func ParseToken(token string, cnf entities.JWT) (int64, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(cnf.SignatureKey), nil
	})

	if err != nil {
		return 0, err
	}
	claims, ok := parsedToken.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.userId, nil
}

func (s *AuthService) CreateToken(cxt context.Context, userId int64, cnf entities.JWT) (string, error) {
	token, err := generateToken(userId, cnf)
	if err != nil {
		log.Printf("AuthService CreateToken GenerateToken err: %s", err)
		return "", fmt.Errorf("%w", err)
	}

	if err = s.tokenRepo.StoreToken(cxt, userId, token, cnf.TTL*time.Second); err != nil {
		log.Printf("AuthService CreateToken StoreToken err: %s", err)
		return "", fmt.Errorf("%w", err)
	}

	return token, nil
}

func generateToken(userId int64, cnf entities.JWT) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(cnf.TTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId})

	return token.SignedString([]byte(cnf.SignatureKey))
}

func (s *AuthService) DeleteToken(cxt context.Context, userId int64, token string) error {
	return s.tokenRepo.DeleteToken(cxt, userId)
}

func (s *AuthService) CheckToken(cxt context.Context, userId int64) (bool, error) {
	return s.tokenRepo.TokenExists(cxt, userId)
}
