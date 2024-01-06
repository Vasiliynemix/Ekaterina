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
	return nil
}
