package serviceAuth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sensors-app/internal/entities"
	"sensors-app/internal/repository/repoErrors"
	"sensors-app/internal/service/serviceErrors"

	"time"

	"github.com/golang-jwt/jwt"
)

type TokenRepo interface {
	StoreToken(cxt context.Context, userId int64, token string, tokenTTL time.Duration) error
	DeleteToken(cxt context.Context, userId int64) error
	GetTokenByUserID(cxt context.Context, userId int64) (string, error)
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
	UserId int64
}

func (s *AuthService) ParseToken(token string, cnf entities.JWT) (int64, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(cnf.SignatureKey), nil
	})

	if err != nil {
		var jwtErr *jwt.ValidationError
		if errors.As(err, &jwtErr) {
			return 0, fmt.Errorf("%w: %w", serviceErrors.ErrParseToken, jwtErr)
		}
		return 0, err
	}
	claims, ok := parsedToken.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
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

func (s *AuthService) DeleteToken(cxt context.Context, userId int64) error {
	if err := s.tokenRepo.DeleteToken(cxt, userId); err != nil {
		if errors.Is(err, repoErrors.ErrNoToken) {
			return serviceErrors.ErrTokenAlreadyRemoved
		}
		return err
	}

	return nil
}

func (s *AuthService) CheckToken(cxt context.Context, userId int64, token string) (bool, error) {
	tokenRepo, err := s.tokenRepo.GetTokenByUserID(cxt, userId)
	if err != nil {
		if errors.Is(err, repoErrors.ErrNoToken) {
			return false, serviceErrors.ErrNoTokenForCheck
		}
		return false, err
	}

	if token != tokenRepo {
		return false, nil
	}
	return true, nil
}

func (s *AuthService) GetUserIDFromCtx(ctx context.Context, key string) (int64, error) {
	id := ctx.Value(key)

	if id == nil {
		return 0, serviceErrors.ErrNoUserIDInCtx
	}

	IdInt64, ok := id.(int64)
	if !ok {
		return 0, serviceErrors.ErrUserIDNotInt64Type
	}
	return IdInt64, nil

}
