package start

import (
	"bot/internal/bot/keyboards/inline"
	"bot/internal/bot/lexicon/commands"
	"bot/internal/bot/lexicon/messages"
	"bot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (r *RouterStart) CheckStartAdmin(msg tgbotapi.Update) bool {
	var tgID int64
	var checkDataOrMsg bool

	switch {
	case msg.Message != nil:
		tgID = msg.Message.Chat.ID
		checkDataOrMsg = msg.Message.Command() == commands.MsgCommandStart
	case msg.CallbackQuery != nil:
		tgID = msg.CallbackQuery.Message.Chat.ID
		checkDataOrMsg = msg.CallbackQuery.Data == commands.MsgCommandStart
	}
	userShow, _ := r.userGetter.GetUserByTgID(tgID)

	return checkDataOrMsg && userShow.IsAdmin
}

func (r *RouterStart) StartAdmin(msg *tgbotapi.Message) {
	_ = r.stateCleaner.ClearState(msg.Chat.ID, config.ScheduleState)
	userShow, _ := r.userGetter.GetUserByTgID(msg.Chat.ID)

	msgSend := tgbotapi.NewMessage(msg.Chat.ID, messages.MessageStartAdmin)

	msgSend.ReplyMarkup = inline.StartKB(userShow.IsAdmin, userShow.IsModer)

	_, err := r.b.Send(msgSend)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
	}
}

func (r *RouterStart) isAdmin(admins []int64, telegramID int64) bool {
	for _, adminID := range admins {
		if adminID == telegramID {
			return true
		}
	}
	return false
}
