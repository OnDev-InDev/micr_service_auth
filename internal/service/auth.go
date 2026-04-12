package service

import (
	"micr_service_auth/internal/storage/repository"
)


type Identification struct {
	repo repository.AuthRepository
}


// идентификация
func (s *Identification)IndetificationUser(username, password string) bool {
	if s.repo.IdentificationRepo(username, password) {
		return true
	}

	return false
}