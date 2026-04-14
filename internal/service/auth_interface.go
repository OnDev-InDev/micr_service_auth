package service

import "micr_service_auth/internal/storage/models"

type UserRepository interface {
	GetEmail(email string) (models.User, error)
}
