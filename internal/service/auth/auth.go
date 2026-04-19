package auth

import (
  "micr_service_auth/internal/storage/models"

	"golang.org/x/crypto/bcrypt"
)



type UserRepository interface {
	GetEmail(email string) (models.User, error)
}



type AuthService struct {
	userRepo UserRepository
}

func (s *AuthService) AuthenticateUser(input LoginInput) (string, error) {
  if err := input.Validate(); err != nil { 
		return "", err 
	}

	user, err := s.userRepo.GetEmail(input.Email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	)

	if err != nil {
		return "", ErrInvalidCredentials
	}

	return user.ID, nil
}
