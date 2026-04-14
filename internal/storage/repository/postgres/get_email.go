package postgres

import (
  "github.com/jinzhu/gorm"
  "micr_service_auth/internal/storage/models"
  "fmt"
)



type PostgresAuthRepo struct {
	db *gorm.DB
}

func (p *PostgresAuthRepo) GetEmail(email string) (*models.User, error) {
  var user models.User
  err := p.db.Where("email = ?", email).First(&user).Error 
  if err != nil {
    return nil, fmt.Errorf("repo GetByEmail: %w", err)
  }
  return &user, nil
}
