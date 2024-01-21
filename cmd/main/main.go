package main

import (
	"bot/internal/bot"
	"bot/internal/config"
	"bot/internal/storage/db"
	"bot/internal/storage/redisdb"
	"bot/pkg/logging"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
)

const (
	levelsUpToRootDir = 3
)

func main() {
	cfg := config.MustLoad(levelsUpToRootDir)

	log := setupLogger(cfg.Env, config.StructDateFormat, cfg.Paths.ConfigInfoPath, cfg.Paths.ConfigDebugPath)

	log.Info("Initializing logger and config...")
	log.Debug("Debug mode on...")

	b, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		panic(err)
	}

	b.Debug = cfg.Debug

	dbConn := connToDB(cfg, log)
	setupDBPool(cfg, dbConn)

	log.Info("RedisDB connected...")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisDB.Host, cfg.RedisDB.Port),
		Password: "",
		DB:       cfg.RedisDB.DB,
	})

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	redisClient := redisdb.New(rdb, log, cfg)

	dbClient := db.New(dbConn, log)

	go bot.Run(b, cfg, log, dbClient, redisClient)

	// Gracefully shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("Got signal", zap.String("signal", sign.String()))
	log.Info("Shutting down...")
}

func connToDB(cfg *config.Config, log *logging.Logger) *gorm.DB {
	dsn := cfg.DB.ConnString()

	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", zap.Error(err))

		os.Exit(1)
	}

	return dbConn
}

func setupDBPool(cfg *config.Config, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(cfg.DB.Pool.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DB.Pool.MaxOpenConns)
	sqlDB.SetConnMaxIdleTime(cfg.DB.Pool.IdleTimeout)
}

func setupLogger(env, structDateFormat, pathToInfoLogs, pathToDebugLogs string) *logging.Logger {
	log := logging.NewLogger(
		logging.InitLogger(
			env,
			structDateFormat,
			pathToInfoLogs,
			pathToDebugLogs,
		),
	)

	go logging.ClearLogFiles(
		pathToInfoLogs,
		pathToDebugLogs,
		structDateFormat,
		log,
	)

	return log
}
