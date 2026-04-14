package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo UserRepository
}


func (s *AuthService) AuthenticateUser(email, password string) error {
	user, err := s.userRepo.GetEmail(email) 
	if err != nil {
		return fmt.Errorf("service AuthenticateUser: %w", err)
	}
	
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return fmt.Errorf("Invalid credentials: %w", err)
	}

	return nil
}
