package postgres

import (
	"fmt"
	"micr_service_auth/internal/domain"

	"github.com/jinzhu/gorm"
)

type PostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (p *PostgresRepo) GetUser(email string) (domain.User, error) {
	var userModel UserModel

	err := p.db.Where("email = ?", email).First(&userModel).Error
	if err != nil {
		return domain.User{}, fmt.Errorf("repo GetUser: %w", err)
	}

	return toDomain(userModel), nil
}
