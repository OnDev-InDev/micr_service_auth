package service

import (
	"micr_service_auth/internal/storage/repository/postgres"
)


type Identification struct {
	repo postgres.AuthRepository
}


// идентификация
func (s *Identification)IndetificationUser(username, password string) bool {
	if s.repo.IdentificationRepo(username, password) {
		return true
	}

	return false
}
