package sheduleRepo

import (
	"bot/internal/db/models"
	"bot/pkg/logging"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrDayExist = fmt.Errorf("day is exist")
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

	result = result.Preload("WeekEven.Days").Preload("WeekOdd.Days").First(&schedule)
	if result.Error != nil {
		r.log.Error("Failed to get schedule", zap.Error(result.Error))
		return models.Schedule{}, result.Error
	}

	return schedule, nil
}

func (r *ScheduleRepo) AddDay(TelegramID int64, dayName string, weekNum int) error {
	schedule, _ := r.GetScheduleByTelegramID(TelegramID)

	err := r.checkDayName(dayName, schedule)
	if err != nil {
		return err
	}

	newDay := models.Day{
		DayName:    dayName,
		WeekEvenID: schedule.WeekEven.ID,
		WeekOddID:  schedule.WeekOdd.ID,
	}

	if weekNum%2 == 1 {
		schedule.WeekOdd.Days = append(schedule.WeekOdd.Days, newDay)
	} else {
		schedule.WeekEven.Days = append(schedule.WeekEven.Days, newDay)
	}

	result := r.db.Save(&schedule)
	if result.Error != nil {
		r.log.Error("Failed to save schedule", zap.Error(result.Error))
		return result.Error
	}

	r.log.Info(
		"Day added",
		zap.String("dayName", dayName),
		zap.Int64("TelegramID", TelegramID),
	)

	return nil
}

func (r *ScheduleRepo) checkDayName(dayName string, schedule models.Schedule) error {
	for _, day := range schedule.WeekEven.Days {
		if day.DayName == dayName {
			return ErrDayExist
		}
	}

	for _, day := range schedule.WeekOdd.Days {
		if day.DayName == dayName {
			return ErrDayExist
		}
	}

	return nil
}
