package start

import (
	"bot/internal/bot/keyboards/inline"
	"bot/internal/bot/lexicon/commands"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (s *RouterStart) CheckStartAdmin(msg tgbotapi.Update) bool {
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
	userShow, _ := s.userGetter.GetUserByTgID(tgID)

	return checkDataOrMsg && userShow.IsAdmin
}

func (s *RouterStart) StartAdmin(msg *tgbotapi.Message) {
	msgSend := tgbotapi.NewMessage(msg.Chat.ID, "Hello! you are admin!")

	msgSend.ReplyMarkup = inline.StartKB

	_, err := s.b.Send(msgSend)
	if err != nil {
		s.log.Error("Failed to send message", zap.Error(err))
	}
}

func (s *RouterStart) isAdmin(admins []int64, telegramID int64) bool {
	for _, adminID := range admins {
		if adminID == telegramID {
			return true
		}
	}
	return false
}
