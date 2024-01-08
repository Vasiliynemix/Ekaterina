package start

import (
	"bot/internal/bot/keyboards/inline"
	"bot/internal/bot/lexicon/commands"
	"bot/internal/bot/lexicon/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (r *RouterStart) CheckCancel(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}

	return callback.Data == inline.DataCancel
}

func (r *RouterStart) Cancel(callback *tgbotapi.CallbackQuery) {
	newCallback := tgbotapi.NewCallback(callback.ID, inline.MsgDataCancel)
	_, err := r.b.Request(newCallback)
	if err != nil {
		r.log.Error("Failed to send callback", zap.Error(err))
	}

	deleteMsg := tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)

	resp, err := r.b.Request(deleteMsg)
	if err != nil || !resp.Ok {
		r.log.Error(
			"Failed to delete message",
			zap.Error(err),
			zap.Bool("ok", resp.Ok),
			zap.Any("result", resp.Result),
			zap.Int("message_id", callback.Message.MessageID),
		)
	}
}

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

func (r *RouterStart) CheckMainMenu(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}

	return callback.Data == inline.DataMainMenu
}

func (r *RouterStart) MenuMain(callback *tgbotapi.CallbackQuery, isAdmin bool, isModer bool) {
	var msgText string

	if isAdmin || isModer {
		msgText = messages.MessageStartAdmin
	} else {
		msgText = messages.MessageStartUser
	}

	msgSend := tgbotapi.NewEditMessageTextAndMarkup(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		msgText,
		inline.StartKB(isAdmin, isModer),
	)

	_, err := r.b.Send(msgSend)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
	}
}
