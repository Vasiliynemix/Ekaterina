package start

import (
	"bot/internal/bot/keyboards/inline"
	"bot/internal/bot/lexicon/commands"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (s *RouterStart) CheckStartAdmin(msg *tgbotapi.Message) bool {
	return msg != nil && msg.Command() == commands.MsgCommandStart && msg.Chat.ID == 5254091301
}

func (s *RouterStart) StartAdmin(msg *tgbotapi.Message) {
	msgSend := tgbotapi.NewMessage(msg.Chat.ID, "Hello! you are admin!")

	msgSend.ReplyMarkup = inline.StartKB

	_, err := s.b.Send(msgSend)
	if err != nil {
		s.log.Error("Failed to send message", zap.Error(err))
	}
}
