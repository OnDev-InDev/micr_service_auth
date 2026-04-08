package service

import (
	"micr_service_auth/internal/storage/repository/postgres"
)

// идентификация
func IndetificationUser(username, password string) bool {
	if !postgres.IdentificationRepo(username, password) {
		return false
	}

	return true
}
