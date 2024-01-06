package mvLogger

import (
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type LoggerMv struct {
	log *logging.Logger
}

func New(log *logging.Logger) *LoggerMv {
	return &LoggerMv{
		log: log,
	}
}

func (l *LoggerMv) UpdateInfo(msg tgbotapi.Update) {
	switch {
	case msg.Message != nil:
		l.log.Info(
			"Update info",
			zap.Int("chat_id", int(msg.Message.Chat.ID)),
			zap.String("username", msg.Message.From.UserName),
			zap.String("text", msg.Message.Text),
		)
	case msg.CallbackQuery != nil:
		l.log.Info(
			"Update info",
			zap.Int("chat_id", int(msg.CallbackQuery.Message.Chat.ID)),
			zap.String("username", msg.CallbackQuery.From.UserName),
			zap.String("data", msg.CallbackQuery.Data),
		)
	}
}
