package bot

import (
	"bot/internal/bot/middlewares"
	"bot/internal/bot/routers/admin/panel"
	"bot/internal/bot/routers/start"
	"bot/internal/bot/routers/user/schedule"
	"bot/internal/config"
	"bot/internal/db"
	"bot/pkg/logging"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Routers struct {
	startRouter  StartRouter
	adminRouters AdminRouters
	userRouters  UserRouters
}

type StartRouter interface {
	CheckStartAdmin(msg tgbotapi.Update) bool
	StartAdmin(msg *tgbotapi.Message)

	CheckStart(msg *tgbotapi.Message) bool
	Start(msg *tgbotapi.Message, isAdmin bool, isModer bool)

	MainMenu(msg *tgbotapi.CallbackQuery, isAdmin bool, isModer bool)
}

type AdminRouters interface {
	CheckAdminPanel(callback *tgbotapi.CallbackQuery) bool
	CheckBackToStartMenu(callback *tgbotapi.CallbackQuery) bool

	ShowAdminPanel(callback *tgbotapi.CallbackQuery)
}

type UserRouters interface {
	CheckSchedule(callback *tgbotapi.CallbackQuery) bool
	ShowSchedule(callback *tgbotapi.CallbackQuery)
}

func initRouters(
	b *tgbotapi.BotAPI,
	log *logging.Logger,
	cfg *config.Config,
	db *db.DB,
) Routers {
	log.Info("Initializing routers...")

	var r Routers

	startRouter := start.New(b, log, &cfg.Bot, db.User)
	r.startRouter = startRouter

	scheduleRouter := schedule.New(b, log)
	r.userRouters = UserRouters(scheduleRouter)

	adminRouter := panel.New(b, log)
	r.adminRouters = AdminRouters(adminRouter)

	return r
}

func Run(b *tgbotapi.BotAPI, cfg *config.Config, log *logging.Logger, db *db.DB) {
	u := setupUpdateConfig(cfg.Bot)

	updates := b.GetUpdatesChan(u)

	log.Info(fmt.Sprintf("Authorized on account bot %s", b.Self.UserName))

	r := initRouters(b, log, cfg, db)

	mv := middlewares.InitMiddlewares(log, db, cfg)

	go checkUpdates(updates, r, mv)
}

func checkUpdates(
	updates tgbotapi.UpdatesChannel,
	r Routers,
	mv *middlewares.Middlewares,
) {
	for update := range updates {
		isAdmin, isModer := mv.MvAddToDB.AddToDB(update)
		mv.MvLog.UpdateInfo(update)
		switch {
		case r.startRouter.CheckStartAdmin(update):
			go r.startRouter.StartAdmin(update.Message)

		case r.startRouter.CheckStart(update.Message):
			go r.startRouter.Start(update.Message, isAdmin, isModer)

		case r.adminRouters.CheckAdminPanel(update.CallbackQuery):
			go r.adminRouters.ShowAdminPanel(update.CallbackQuery)

		case r.adminRouters.CheckBackToStartMenu(update.CallbackQuery):
			go r.startRouter.MainMenu(update.CallbackQuery, isAdmin, isModer)

		case r.userRouters.CheckSchedule(update.CallbackQuery):
			go r.userRouters.ShowSchedule(update.CallbackQuery)

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
