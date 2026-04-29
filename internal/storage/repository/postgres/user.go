package postgres

import (
	"micr_service_auth/internal/domain"
	"micr_service_auth/internal/errors"

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

func (p *PostgresRepo) Get(email string) (domain.User, error) {
	var userModel UserModel

	err := p.db.Where("email = ?", email).First(&userModel).Error
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			return domain.User{}, errors.Wrap(errors.ErrDB, "user lookup failed", err)
		}

		return domain.User{}, errors.Wrap(errors.ErrDB, "db error", err)
	}

	return toDomain(userModel), nil
}
