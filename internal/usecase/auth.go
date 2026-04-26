package usecase

import (
	"context"
	"micr_service_auth/internal/domain"
	"micr_service_auth/internal/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AuthService interface {
	AuthenticateUser(email, password string) (string, error)
}

type SessionService interface {
	CreateSession(ctx context.Context, userID string) (string, error)
	CheckSession(ctx context.Context, sessionID string) (domain.Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
}

type LoginInput struct {
	Email    string
	Password string
}

func (i LoginInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Email,
			validation.Required,
			validation.Length(3, 100),
		),
		validation.Field(&i.Password,
			validation.Required,
			validation.Length(6, 100),
		),
	)
}

type AuthUsecase struct {
	authService    AuthService
	sessionService SessionService
}

func NewUsecase(authService AuthService, sessionService SessionService) *AuthUsecase {
	return &AuthUsecase{
		authService:    authService,
		sessionService: sessionService,
	}
}

func (uc *AuthUsecase) Login(ctx context.Context, input LoginInput) (string, error) {

	if err := input.Validate(); err != nil {
		return "", errors.New(errors.CodeValidationError, "invalid input")
	}

	// 1. идентификация / авторизация
	userID, err := uc.authService.AuthenticateUser(input.Email, input.Password)
	if err != nil {
		return "", err
	}

	// 2. создание сессии
	sessionID, err := uc.sessionService.CreateSession(ctx, userID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (uc *AuthUsecase) GetSession(ctx context.Context, sessionID string) (domain.Session, error) {

	ses, err := uc.sessionService.CheckSession(ctx, sessionID)
	if err != nil {
		return domain.Session{}, err
	}

	return ses, nil

}

func (uc *AuthUsecase) Logout(ctx context.Context, sessionID string) error {
	return uc.sessionService.DeleteSession(ctx, sessionID)
}
