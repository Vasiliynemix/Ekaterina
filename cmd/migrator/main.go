package main

import (
	"bot/internal/config"
	models2 "bot/internal/storage/db/models"
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
	ok := db.Migrator().HasTable(&models2.User{})
	if !ok {
		err := db.Migrator().CreateTable(&models2.User{})
		if err != nil {
			return err
		}
		fmt.Println("user table created")
	}

	ok = db.Migrator().HasTable(&models2.Schedule{})
	if !ok {
		err := db.Migrator().CreateTable(&models2.Schedule{})
		if err != nil {
			return err
		}
		fmt.Println("schedule table created")
	}

	ok = db.Migrator().HasTable(&models2.WeekEven{})
	if !ok {
		err := db.Migrator().CreateTable(&models2.WeekEven{})
		if err != nil {
			return err
		}
		fmt.Println("weekEven table created")
	}

	ok = db.Migrator().HasTable(&models2.WeekOdd{})
	if !ok {
		err := db.Migrator().CreateTable(&models2.WeekOdd{})
		if err != nil {
			return err
		}
		fmt.Println("weekOdd table created")
	}

	ok = db.Migrator().HasTable(&models2.Day{})
	if !ok {
		err := db.Migrator().CreateTable(&models2.Day{})
		if err != nil {
			return err
		}
		fmt.Println("day table created")
	}

	_ = dropColumns(db)

	err := db.AutoMigrate(
		&models2.User{},
		&models2.Schedule{},
		&models2.WeekEven{},
		&models2.WeekOdd{},
		&models2.Day{},
	)
	if err != nil {
		return err
	}

	//users := getUsers(db)
	//for _, user := range users {
	//	_ = addSchedule(db, user.TelegramID)
	//}

	return nil
}

func dropColumns(db *gorm.DB) error {
	query := "ALTER TABLE days DROP COLUMN IF EXISTS prepod"
	if err := db.Exec(query).Error; err != nil {
		fmt.Println("Failed to drop column", err)
		return err
	}

	return nil
}

func getUsers(db *gorm.DB) []models2.User {
	var users []models2.User
	result := db.Find(&users)
	if result.Error != nil {
		fmt.Println("Failed to get users", result.Error)
		return []models2.User{}
	}

	return users
}

func addSchedule(db *gorm.DB, TelegramID int64) error {
	newSchedule := models2.Schedule{
		TelegramID: TelegramID,
	}

	result := db.Create(&newSchedule)
	if result.Error != nil {
		fmt.Println("Failed to add schedule", result.Error)
		return result.Error
	}

	newWeekEven := models2.WeekEven{
		ScheduleID: newSchedule.ID,
	}

	result = db.Create(&newWeekEven)
	if result.Error != nil {
		fmt.Println("Failed to add weekEven", result.Error)
		return result.Error
	}

	newWeekOdd := models2.WeekOdd{
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
