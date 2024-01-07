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

	return nil
}
