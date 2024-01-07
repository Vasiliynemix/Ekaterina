package main

import (
	"bot/internal/config"
	"bot/internal/db/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.MustLoad(3)

	dsn := cfg.DB.ConnString()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = runMigrations(db)
	if err != nil {
		panic(err)
	}
}

func runMigrations(db *gorm.DB) error {
	ok := db.Migrator().HasTable(&models.User{})
	if !ok {
		err := db.Migrator().CreateTable(&models.User{})
		if err != nil {
			return err
		}
		fmt.Println("user table created")
	}

	ok = db.Migrator().HasTable(&models.Schedule{})
	if !ok {
		err := db.Migrator().CreateTable(&models.Schedule{})
		if err != nil {
			return err
		}
		fmt.Println("schedule table created")
	}

	ok = db.Migrator().HasTable(&models.WeekEven{})
	if !ok {
		err := db.Migrator().CreateTable(&models.WeekEven{})
		if err != nil {
			return err
		}
		fmt.Println("weekEven table created")
	}

	ok = db.Migrator().HasTable(&models.WeekOdd{})
	if !ok {
		err := db.Migrator().CreateTable(&models.WeekOdd{})
		if err != nil {
			return err
		}
		fmt.Println("weekOdd table created")
	}

	ok = db.Migrator().HasTable(&models.Day{})
	if !ok {
		err := db.Migrator().CreateTable(&models.Day{})
		if err != nil {
			return err
		}
		fmt.Println("day table created")
	}

	err := db.AutoMigrate(
		&models.User{},
		&models.Schedule{},
		&models.WeekEven{},
		&models.WeekOdd{},
		&models.Day{},
	)
	if err != nil {
		return err
	}

	users := getUsers(db)
	for _, user := range users {
		_ = addSchedule(db, user.TelegramID)
	}

	return nil
}

func getUsers(db *gorm.DB) []models.User {
	var users []models.User
	result := db.Find(&users)
	if result.Error != nil {
		fmt.Println("Failed to get users", result.Error)
		return []models.User{}
	}

	return users
}

func addSchedule(db *gorm.DB, TelegramID int64) error {
	newSchedule := models.Schedule{
		TelegramID: TelegramID,
	}

	result := db.Create(&newSchedule)
	if result.Error != nil {
		fmt.Println("Failed to add schedule", result.Error)
		return result.Error
	}

	newWeekEven := models.WeekEven{
		ScheduleID: newSchedule.ID,
	}

	result = db.Create(&newWeekEven)
	if result.Error != nil {
		fmt.Println("Failed to add weekEven", result.Error)
		return result.Error
	}

	newWeekOdd := models.WeekOdd{
		ScheduleID: newSchedule.ID,
	}

	result = db.Create(&newWeekOdd)
	if result.Error != nil {
		fmt.Println("Failed to add weekOdd", result.Error)
		return result.Error
	}

	fmt.Println("Schedule, weekEven and weekOdd added", TelegramID)

	return nil
}
