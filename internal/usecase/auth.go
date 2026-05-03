package usecase

import (
	"context"
	"time"

	"github.com/OnDev-InDev/micr_service_auth/internal/domain"
	"github.com/OnDev-InDev/micr_service_auth/internal/errors"
	"github.com/OnDev-InDev/micr_service_auth/internal/pkg"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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
	userRepo    UserRepository
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

func (uc *AuthUC) Login(ctx context.Context, input LoginInput) (string, error) {
	if err := input.Validate(); err != nil {
		return "", err
	}

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

	sessionID := pkg.GenerateID()

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
	if sessionID == "" {
		return errors.New(errors.CodeForbidden, "empty session id")
	}
	return uc.sessionRepo.Delete(ctx, sessionID)
}
