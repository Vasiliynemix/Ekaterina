package start

import (
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RouterStart struct {
	b   *tgbotapi.BotAPI
	log *logging.Logger
}

func New(b *tgbotapi.BotAPI, log *logging.Logger) *RouterStart {
	return &RouterStart{
		b:   b,
		log: log,
	}
}
