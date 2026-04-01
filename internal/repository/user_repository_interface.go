package repository

import (
	"github.com/rendi-hendra/resful-api/internal/entity"
	"gorm.io/gorm"
)

// IUserRepository defines the contract for all User data access operations.
type IUserRepository interface {
	Create(db *gorm.DB, user *entity.User) error
	Update(db *gorm.DB, user *entity.User) error
	Delete(db *gorm.DB, user *entity.User) error
	CountById(db *gorm.DB, id any) (int64, error)
	FindById(db *gorm.DB, user *entity.User, id any) error
	FindByToken(db *gorm.DB, user *entity.User, token string) error
}
