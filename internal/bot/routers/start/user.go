package start

import (
	"bot/internal/bot/keyboards/inline"
	"bot/internal/bot/lexicon/commands"
	"bot/internal/bot/lexicon/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (r *RouterStart) CheckStart(msg *tgbotapi.Message) bool {
	return msg != nil && msg.Command() == commands.MsgCommandStart
}

func (r *RouterStart) Start(msg *tgbotapi.Message, isAdmin bool, isModer bool) {
	msgSend := tgbotapi.NewMessage(msg.Chat.ID, messages.MessageStartUser)
	msgSend.ReplyMarkup = inline.StartKB(isAdmin, isModer)

	_, err := r.b.Send(msgSend)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
	}
}

func (r *RouterStart) MainMenu(msg *tgbotapi.CallbackQuery, isAdmin bool, isModer bool) {
	var msgText string

	if isAdmin || isModer {
		msgText = messages.MessageStartAdmin
	} else {
		msgText = messages.MessageStartUser
	}

	msgSend := tgbotapi.NewEditMessageTextAndMarkup(
		msg.Message.Chat.ID,
		msg.Message.MessageID,
		msgText,
		inline.StartKB(isAdmin, isModer),
	)

	_, err := r.b.Send(msgSend)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
	}
}
