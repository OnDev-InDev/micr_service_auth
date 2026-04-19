package usecase

import (
	"context"
	"micr_service_auth/internal/domain"
	"micr_service_auth/internal/service/auth"
	"micr_service_auth/internal/service/session"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)


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
	authService    *auth.AuthService
	sessionService *session.SessionService
}



func (uc *AuthUsecase) Login(ctx context.Context, input LoginInput) (string, error) {

	if err := input.Validate(); err != nil { 
		return "", err 
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



func (uc *AuthUsecase) ValidateSession(ctx context.Context, sessionID string) (domain.Session, error) {

	ses, err := uc.sessionService.CheckSession(ctx, sessionID)
	if err != nil {
		return domain.Session{}, err
	}

	return ses, nil

}



func (uc *AuthUsecase) Logout(ctx context.Context, sessionID string) error {
  uc.sessionService.DeleteSession(ctx, sessionID)
}