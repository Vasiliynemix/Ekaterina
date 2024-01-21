package mvAddToDB

import (
	"bot/internal/config"
	"bot/internal/storage/db/repositories/userRepo"
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"time"
)

type AddToDBMv struct {
	log           *logging.Logger
	cfg           *config.Config
	userSaver     UserSaver
	scheduleSaver ScheduleSaver
}

type UserSaver interface {
	AddUser(user userRepo.UserAddParams) error
	GetUserByTgID(telegramID int64) (*userRepo.UserShow, error)
}

type ScheduleSaver interface {
	AddSchedule(TelegramID int64) error
}

func New(
	log *logging.Logger,
	userSaver UserSaver,
	scheduleSaver ScheduleSaver,
	cfg *config.Config,
) *AddToDBMv {
	return &AddToDBMv{
		log:           log,
		userSaver:     userSaver,
		scheduleSaver: scheduleSaver,
		cfg:           cfg,
	}
}

func (l *AddToDBMv) AddAndGetUserToDB(msg tgbotapi.Update) (bool, bool, string) {
	var telegramID int64
	var userName string

	switch {
	case msg.Message != nil:
		telegramID = msg.Message.Chat.ID
		userName = msg.Message.From.UserName

	case msg.CallbackQuery != nil:
		telegramID = msg.CallbackQuery.Message.Chat.ID
		userName = msg.CallbackQuery.From.UserName
	}

	userShow, err := l.userSaver.GetUserByTgID(telegramID)

	if err == nil {
		return userShow.IsAdmin, userShow.IsModer, userShow.TypeSchedule
	}

	isAdmin := false

	for _, adminID := range l.cfg.Bot.Admins {
		if adminID == telegramID {
			l.log.Info("User is admin", zap.Uint("telegramID", uint(telegramID)))
			isAdmin = true
			break
		}
	}

	userAdd := userRepo.UserAddParams{
		TelegramID: telegramID,
		UserName:   userName,
		CreatedAt:  time.Now().UTC().Unix(),
		IsAdmin:    isAdmin,
	}

	_ = l.userSaver.AddUser(userAdd)

	_ = l.scheduleSaver.AddSchedule(telegramID)

	l.log.Info("User added to DB and", zap.Uint("telegramID", uint(telegramID)))

	return userShow.IsAdmin, userShow.IsModer, userShow.TypeSchedule
}
