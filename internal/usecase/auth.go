package usecase

import (
	"context"
	"micr_service_auth/internal/domain"
	"micr_service_auth/internal/errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)


// нужно интепретировать ошибки из репо, так как они не должны подниматься наверх

type UserRepository interface {
	Get(email string) (domain.User, error)
}


type SessionRepository interface {
	Create(ctx context.Context, session domain.Session) error
	Get(ctx context.Context, value string) (domain.Session, error)
	Delete(ctx context.Context, value string) error
}

type AuthUC struct {
	userRepo UserRepository
	sessionRepo SessionRepository
}  



type LoginInput struct {
	Email    string
	Password string
}




func (i LoginInput) Validate() error {
	err := validation.ValidateStruct(&i,
		validation.Field(&i.Email,
			validation.Required,
			validation.Length(3, 100),
		),
		validation.Field(&i.Password,
			validation.Required,
			validation.Length(6, 100),
		),
	)

	if err != nil {
		return errors.Wrap(errors.CodeValidationError, "validation failed", err)
	}

	return nil
}



func NewUsecase(userRepo UserRepository, sessionRepo SessionRepository) *AuthUC {
	return &AuthUC{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}



// генерация ID
func generateSessionID() string {
	newID := uuid.New().String()
	return newID
}




func (uc *AuthUC) Login(ctx context.Context, input LoginInput) (string, error) {
	if err := input.Validate(); err != nil {
		return "", err
	}

	// 1. идентификация / авторизация
	user, err := uc.userRepo.Get(input.Email)
	if err != nil {
		return "", err
	}
  
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	)
  if err != nil {
		return "", errors.New(errors.CodeInvalidCredentials, "invalid credentials")
	}

  sessionID := generateSessionID()

	//описываем новую сессию
	session := domain.Session{
		ID:        sessionID,
		UserID:    user.ID,
		Role:      "user",
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

  if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return "", errors.Wrap(errors.CodeInternal, "session repository failure", err)
	}
	return sessionID, nil
	
	
}




func (uc *AuthUC) GetSession(ctx context.Context, sessionID string) (domain.Session, error) {
	// Получаем JSON из Redis
	session, err := uc.sessionRepo.Get(ctx, sessionID)
	if err != nil {
		return domain.Session{}, err
	}

	if time.Now().After(session.ExpiresAt) {
		return domain.Session{}, errors.New(errors.CodeSessionExpired, "session expired")
	}

	return session, nil
}



func (uc *AuthUC) Logout(ctx context.Context, sessionID string) error {
	return uc.sessionRepo.Delete(ctx, sessionID)
}
