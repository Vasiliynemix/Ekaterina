package userRepo

import (
	"bot/internal/db/models"
	"bot/pkg/logging"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

const (
	ErrDuplicateKeyStr = "duplicate key value violates unique constraint"
)

var (
	ErrUserExists = fmt.Errorf("user already exists")
)

type UserRepo struct {
	db  *gorm.DB
	log *logging.Logger
}

var (
	ErrNotFound = fmt.Errorf("not found")
)

func New(db *gorm.DB, log *logging.Logger) *UserRepo {
	return &UserRepo{
		db:  db,
		log: log,
	}
}

func (r *UserRepo) AddUser(user UserAddParams) error {
	userAdd := models.User{
		TelegramID: user.TelegramID,
		UserName:   user.UserName,
		CreatedAt:  user.CreatedAt,
		IsAdmin:    user.IsAdmin,
	}

	result := r.db.Save(&userAdd)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), ErrDuplicateKeyStr) {
			r.log.Info("User already exists", zap.Error(result.Error))

			return ErrUserExists
		}
		r.log.Error("Failed to add user", zap.Error(result.Error))

		return result.Error
	}

	return nil
}

func (r *UserRepo) GetUserByTgID(telegramID int64) (*UserShow, error) {
	var userFind models.User

	mapFind := map[string]interface{}{
		"telegram_id": telegramID,
	}

	result := r.db.Find(&userFind, mapFind).First(&userFind)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.log.Info("User not found", zap.Error(result.Error))

			return &UserShow{}, ErrNotFound
		}

		r.log.Error("Failed to get user", zap.Error(result.Error))

		return &UserShow{}, result.Error
	}

	userShow := UserShow{
		TelegramID: userFind.TelegramID,
		UserName:   userFind.UserName,
		IsAdmin:    userFind.IsAdmin,
		IsModer:    userFind.IsModer,
	}

	return &userShow, nil
}
