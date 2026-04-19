package postgres

import (
	"fmt"
	"micr_service_auth/internal/domain"

	"github.com/jinzhu/gorm"
)

type PostgresAuthRepo struct {
	db *gorm.DB
}

func (p *PostgresAuthRepo) GetUser(email string) (domain.User, error) {
	var userModel UserModel

	err := p.db.Where("email = ?", email).First(&userModel).Error
	if err != nil {
		return domain.User{}, fmt.Errorf("repo GetUser: %w", err)
	}

	return toDomain(userModel), nil
}
