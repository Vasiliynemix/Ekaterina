package main

import (
	"bot/internal/bot"
	"bot/internal/config"
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
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

	log.Debug("config: ", zap.Any("config", cfg))
	log.Info("Initializing logger and config...")
	log.Debug("Debug mode on...")

	b, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		panic(err)
	}

	b.Debug = cfg.Debug

	go bot.Run(b, cfg, log)

	// Gracefully shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("Got signal", zap.String("signal", sign.String()))
	log.Info("Shutting down...")
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
