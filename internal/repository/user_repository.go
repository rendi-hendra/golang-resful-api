package repository

import (
	"github.com/rendi-hendra/resful-api/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	Repository[entity.User]
	log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) UserRepository {
	return &UserRepositoryImpl{
		log: log,
	}
}

func (r *UserRepositoryImpl) FindByToken(db *gorm.DB, user *entity.User, token string) error {
	return db.Where("token = ?", token).First(user).Error
}
