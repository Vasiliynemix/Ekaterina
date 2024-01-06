package bot

import (
	"bot/internal/bot/middlewares"
	"bot/internal/bot/routers/start"
	"bot/internal/config"
	"bot/pkg/logging"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Routers struct {
	startRouter  StartRouter
	adminRouters AdminRouters
}

type StartRouter interface {
	CheckStartAdmin(msg *tgbotapi.Message) bool
	StartAdmin(msg *tgbotapi.Message)
	CheckStart(msg *tgbotapi.Message) bool
	Start(msg *tgbotapi.Message)
}

type AdminRouters interface {
}

func initRouters(b *tgbotapi.BotAPI, log *logging.Logger) Routers {
	log.Info("Initializing routers...")

	var r Routers

	startRouter := start.New(b, log)
	r.startRouter = startRouter

	r.adminRouters = startRouter

	return r
}

func Run(b *tgbotapi.BotAPI, cfg *config.Config, log *logging.Logger) {
	u := setupUpdateConfig(cfg.Bot)

	updates := b.GetUpdatesChan(u)

	log.Info(fmt.Sprintf("Authorized on account bot %s", b.Self.UserName))

	r := initRouters(b, log)

	mv := middlewares.InitMiddlewares(log)

	go checkUpdates(updates, r, mv)
}

func checkUpdates(
	updates tgbotapi.UpdatesChannel,
	r Routers,
	mv *middlewares.Middlewares,
) {
	for update := range updates {
		mv.MvLog.UpdateInfo(update)
		switch {
		case r.startRouter.CheckStartAdmin(update.Message):
			go r.startRouter.StartAdmin(update.Message)
		case r.startRouter.CheckStart(update.Message):
			go r.startRouter.Start(update.Message)
		default:
			continue
		}
	}
}

func setupUpdateConfig(BotCfg config.BotConfig) tgbotapi.UpdateConfig {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = BotCfg.TimeOut

	return u
}
