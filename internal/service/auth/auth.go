package auth

import (
	"micr_service_auth/internal/domain"

	"golang.org/x/crypto/bcrypt"
)



type UserRepository interface {
	GetUser(email string) (domain.User, error)
}



type AuthService struct {
	userRepo UserRepository
}

func (s *AuthService) AuthenticateUser(email, password string) (string, error) {
  
	user, err := s.userRepo.GetUser(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", ErrInvalidCredentials
	}

	return user.ID, nil
}
