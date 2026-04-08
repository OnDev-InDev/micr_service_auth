package service

import (
	"micr_service_auth/internal/storage/repository/postgres"
)

func AuthenticateUser(username, password string) bool {
	if postgres.IdentificationRepo(username, password) {
		return true
	}
	return false
}
