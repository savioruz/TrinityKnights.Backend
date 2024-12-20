package user

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	repository.RepositoryImpl[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(db *gorm.DB, log *logrus.Logger) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		RepositoryImpl: repository.RepositoryImpl[entity.User]{DB: db},
		Log:            log,
	}
}

func (r *UserRepositoryImpl) GetFirst(db *gorm.DB, user *entity.User) error {
	return db.First(&user).Error
}

func (r *UserRepositoryImpl) GetByID(db *gorm.DB, user *entity.User, id string) error {
	return db.Where("id = ?", id).Take(&user).Error
}

func (r *UserRepositoryImpl) GetByEmail(db *gorm.DB, user *entity.User, email string) error {
	return db.Where("email = ?", email).Take(&user).Error
}

func (r *UserRepositoryImpl) CountByRole(db *gorm.DB, role string) (int64, error) {
	var count int64
	err := db.Model(&entity.User{}).Where("role = ?", role).Count(&count).Error
	return count, err
}

func (r *UserRepositoryImpl) GetByResetPasswordToken(db *gorm.DB, user *entity.User, token string) error  {
	return db.Where("reset_password_token = ?", token).Take(&user).Error
}

func (r *UserRepositoryImpl) GetByVerifyEmailToken(db *gorm.DB, user *entity.User, token string) error  { 
	return db.Where("verify_email_token = ?", token).Take(&user).Error 
} 