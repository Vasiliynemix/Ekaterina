package sheduleRepo

import (
	"bot/internal/db/models"
	"bot/pkg/logging"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ScheduleRepo struct {
	db  *gorm.DB
	log *logging.Logger
}

func New(db *gorm.DB, log *logging.Logger) *ScheduleRepo {
	return &ScheduleRepo{
		db:  db,
		log: log,
	}
}

func (r *ScheduleRepo) AddSchedule(TelegramID int64) error {
	newSchedule := models.Schedule{
		TelegramID: TelegramID,
	}

	result := r.db.Create(&newSchedule)
	if result.Error != nil {
		r.log.Error("Failed to add schedule", zap.Error(result.Error))

		return result.Error
	}

	newWeekEven := models.WeekEven{
		ScheduleID: newSchedule.ID,
	}

	result = r.db.Create(&newWeekEven)
	if result.Error != nil {
		r.log.Error("Failed to add weekEven", zap.Error(result.Error))

		return result.Error
	}

	newWeekOdd := models.WeekOdd{
		ScheduleID: newSchedule.ID,
	}

	result = r.db.Create(&newWeekOdd)
	if result.Error != nil {
		r.log.Error("Failed to add weekOdd", zap.Error(result.Error))

		return result.Error
	}

	r.log.Info("Schedule, weekEven and weekOdd added", zap.Int64("TelegramID", TelegramID))

	return nil
}

func (r *ScheduleRepo) GetScheduleByTelegramID(TelegramID int64) (models.Schedule, error) {
	var schedule models.Schedule

	result := r.db.Where("telegram_id = ?", TelegramID)
	if result.Error != nil {
		r.log.Error("Failed to get schedule", zap.Error(result.Error))
		return models.Schedule{}, result.Error
	}

	result = result.Preload("WeekEven").Preload("WeekOdd").First(&schedule)
	if result.Error != nil {
		r.log.Error("Failed to get schedule", zap.Error(result.Error))
		return models.Schedule{}, result.Error
	}

	return schedule, nil
}
