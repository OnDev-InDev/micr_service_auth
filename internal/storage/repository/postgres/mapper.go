package postgres

import "micr_service_auth/internal/domain"

func toDomain(u UserModel) domain.User {
	return domain.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
	}
}
