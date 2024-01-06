package db

import (
	"bot/internal/db/repositories/userRepo"
	"bot/pkg/logging"
	"gorm.io/gorm"
)

type DB struct {
	log  *logging.Logger
	User *userRepo.UserRepo
}

func New(db *gorm.DB, log *logging.Logger) *DB {
	user := userRepo.New(db, log)

	return &DB{
		log:  log,
		User: user,
	}
}
