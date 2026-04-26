package auth

import (
	"micr_service_auth/internal/domain"
	"micr_service_auth/internal/errors"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetUser(email string) (domain.User, error)
}

type AuthService struct {
	userRepo UserRepository
}

func NewAuthService(userRepo UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) AuthenticateUser(email, password string) (string, error) {

	user, err := s.userRepo.GetUser(email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", errors.New(errors.CodeInvalidCredentials, "invalid credentials")
	}

	return user.ID, nil
}
