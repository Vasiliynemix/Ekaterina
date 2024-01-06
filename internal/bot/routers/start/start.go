package start

import (
	"bot/internal/config"
	"bot/internal/db/repositories/userRepo"
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RouterStart struct {
	b          *tgbotapi.BotAPI
	log        *logging.Logger
	cfg        *config.BotConfig
	userGetter UserGetter
}

type UserGetter interface {
	GetUserByTgID(telegramID int64) (*userRepo.UserShow, error)
}

func New(
	b *tgbotapi.BotAPI,
	log *logging.Logger,
	cfg *config.BotConfig,
	userGetter UserGetter,
) *RouterStart {
	return &RouterStart{
		b:          b,
		log:        log,
		cfg:        cfg,
		userGetter: userGetter,
	}
}
