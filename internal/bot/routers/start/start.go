package start

import (
	"bot/internal/config"
	"bot/internal/storage/db/repositories/userRepo"
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RouterStart struct {
	b            *tgbotapi.BotAPI
	log          *logging.Logger
	cfg          *config.BotConfig
	userGetter   UserGetter
	stateCleaner StateCleaner
}

type UserGetter interface {
	GetUserByTgID(telegramID int64) (*userRepo.UserShow, error)
}

type StateCleaner interface {
	ClearState(telegramID int64, stateName string) error
}

func New(
	b *tgbotapi.BotAPI,
	log *logging.Logger,
	cfg *config.BotConfig,
	userGetter UserGetter,
	stateCleaner StateCleaner,
) *RouterStart {
	return &RouterStart{
		b:            b,
		log:          log,
		cfg:          cfg,
		userGetter:   userGetter,
		stateCleaner: stateCleaner,
	}
}
