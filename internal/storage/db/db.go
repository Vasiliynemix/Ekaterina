package db

import (
	"bot/internal/storage/db/repositories/sheduleRepo"
	"bot/internal/storage/db/repositories/userRepo"
	"bot/pkg/logging"
	"gorm.io/gorm"
)

type DB struct {
	log      *logging.Logger
	User     *userRepo.UserRepo
	Schedule *sheduleRepo.ScheduleRepo
}

func New(db *gorm.DB, log *logging.Logger) *DB {
	user := userRepo.New(db, log)
	schedule := sheduleRepo.New(db, log)

	return &DB{
		log:      log,
		User:     user,
		Schedule: schedule,
	}
}
