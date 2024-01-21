package bot

import (
	"bot/internal/bot/middlewares"
	"bot/internal/bot/routers/admin/panel"
	"bot/internal/bot/routers/start"
	"bot/internal/bot/routers/user/schedule"
	"bot/internal/config"
	"bot/internal/storage/db"
	"bot/internal/storage/redisdb"
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
	CheckMainMenu(callback *tgbotapi.CallbackQuery) bool
	MenuMain(callback *tgbotapi.CallbackQuery, isAdmin bool, isModer bool)

	CheckStartAdmin(msg tgbotapi.Update) bool
	StartAdmin(msg *tgbotapi.Message)

	CheckStart(msg *tgbotapi.Message) bool
	Start(msg *tgbotapi.Message, isAdmin bool, isModer bool)

	CheckCancel(callback *tgbotapi.CallbackQuery) bool
	Cancel(callback *tgbotapi.CallbackQuery)
}

type AdminRouters interface {
	CheckAdminPanel(callback *tgbotapi.CallbackQuery) bool
	CheckBackToStartMenu(callback *tgbotapi.CallbackQuery) bool

	ShowAdminPanel(callback *tgbotapi.CallbackQuery)
}

type UserRouters interface {
	CheckSchedule(callback *tgbotapi.CallbackQuery) bool
	ShowSchedule(callback *tgbotapi.CallbackQuery, typeSchedule string)

	CheckTypeSchedule(callback *tgbotapi.CallbackQuery) bool
	TypeSchedule(callback *tgbotapi.CallbackQuery)

	CheckScheduleWeekEven(callback *tgbotapi.CallbackQuery) bool
	ScheduleWeekEven(callback *tgbotapi.CallbackQuery, typeSchedule string)

	CheckScheduleWeekOdd(callback *tgbotapi.CallbackQuery) bool
	ScheduleWeekOdd(callback *tgbotapi.CallbackQuery, typeSchedule string)

	CheckAddScheduleWeek(callback *tgbotapi.CallbackQuery) bool
	AddScheduleWeek(callback *tgbotapi.CallbackQuery)

	CheckStateSchedule(stateSchedule map[string]interface{}, msg *tgbotapi.Message) bool
	FileScheduleWeek(stateSchedule map[string]interface{}, typeSchedule string, msg *tgbotapi.Message)

	CheckBackToScheduleMenu(callback *tgbotapi.CallbackQuery) bool

	CheckDayMonday(callback *tgbotapi.CallbackQuery) bool
	CheckDayTuesday(callback *tgbotapi.CallbackQuery) bool
	CheckDayWednesday(callback *tgbotapi.CallbackQuery) bool
	CheckDayThursday(callback *tgbotapi.CallbackQuery) bool
	CheckDayFriday(callback *tgbotapi.CallbackQuery) bool
	CheckDaySaturday(callback *tgbotapi.CallbackQuery) bool
	CheckDaySunday(callback *tgbotapi.CallbackQuery) bool

	ShowDay(callback *tgbotapi.CallbackQuery)
}

func initRouters(
	b *tgbotapi.BotAPI,
	log *logging.Logger,
	cfg *config.Config,
	db *db.DB,
	redis *redisdb.RedisDB,
) Routers {
	log.Info("Initializing routers...")

	var r Routers

	startRouter := start.New(b, log, &cfg.Bot, db.User, redis)
	r.startRouter = startRouter

	scheduleRouter := schedule.New(b, log, cfg, db.Schedule, db.User, redis)
	r.userRouters = UserRouters(scheduleRouter)

	adminRouter := panel.New(b, log)
	r.adminRouters = AdminRouters(adminRouter)

	return r
}

func Run(b *tgbotapi.BotAPI, cfg *config.Config, log *logging.Logger, db *db.DB, redis *redisdb.RedisDB) {
	u := setupUpdateConfig(cfg.Bot)

	updates := b.GetUpdatesChan(u)

	log.Info(fmt.Sprintf("Authorized on account bot %s", b.Self.UserName))

	r := initRouters(b, log, cfg, db, redis)

	mv := middlewares.InitMiddlewares(log, db, cfg, redis)

	go checkUpdates(updates, r, mv)
}

func checkUpdates(
	updates tgbotapi.UpdatesChannel,
	r Routers,
	mv *middlewares.Middlewares,
) {
	for update := range updates {
		isAdmin, isModer, typeSchedule := mv.MvAddToDB.AddAndGetUserToDB(update)
		userState := mv.MvGetState.UserState(update)
		mv.MvLog.UpdateInfo(update)
		switch {
		case r.startRouter.CheckStartAdmin(update):
			go r.startRouter.StartAdmin(update.Message)

		case r.startRouter.CheckStart(update.Message):
			go r.startRouter.Start(update.Message, isAdmin, isModer)

		case r.startRouter.CheckMainMenu(update.CallbackQuery):
			go r.startRouter.MenuMain(update.CallbackQuery, isAdmin, isModer)

		case r.startRouter.CheckCancel(update.CallbackQuery):
			go r.startRouter.Cancel(update.CallbackQuery)

		case r.adminRouters.CheckAdminPanel(update.CallbackQuery):
			go r.adminRouters.ShowAdminPanel(update.CallbackQuery)

		case r.adminRouters.CheckBackToStartMenu(update.CallbackQuery):
			go r.startRouter.MenuMain(update.CallbackQuery, isAdmin, isModer)

		case r.userRouters.CheckBackToScheduleMenu(update.CallbackQuery):
			go r.userRouters.ShowSchedule(update.CallbackQuery, typeSchedule)

		case r.userRouters.CheckTypeSchedule(update.CallbackQuery):
			r.userRouters.TypeSchedule(update.CallbackQuery)

		case r.userRouters.CheckSchedule(update.CallbackQuery):
			go r.userRouters.ShowSchedule(update.CallbackQuery, typeSchedule)

		case r.userRouters.CheckScheduleWeekEven(update.CallbackQuery):
			go r.userRouters.ScheduleWeekEven(update.CallbackQuery, typeSchedule)

		case r.userRouters.CheckScheduleWeekOdd(update.CallbackQuery):
			go r.userRouters.ScheduleWeekOdd(update.CallbackQuery, typeSchedule)

		case r.userRouters.CheckAddScheduleWeek(update.CallbackQuery):
			go r.userRouters.AddScheduleWeek(update.CallbackQuery)

		case r.userRouters.CheckStateSchedule(userState, update.Message):
			go r.userRouters.FileScheduleWeek(userState, typeSchedule, update.Message)

		case r.userRouters.CheckDayMonday(update.CallbackQuery) ||
			r.userRouters.CheckDayTuesday(update.CallbackQuery) ||
			r.userRouters.CheckDayWednesday(update.CallbackQuery) ||
			r.userRouters.CheckDayThursday(update.CallbackQuery) ||
			r.userRouters.CheckDayFriday(update.CallbackQuery) ||
			r.userRouters.CheckDaySaturday(update.CallbackQuery) ||
			r.userRouters.CheckDaySunday(update.CallbackQuery):
			go r.userRouters.ShowDay(update.CallbackQuery)

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
