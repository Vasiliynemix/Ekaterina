package panel

import (
	"bot/internal/bot/keyboards/inline"
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type RouterAdminPanel struct {
	b   *tgbotapi.BotAPI
	log *logging.Logger
}

func New(b *tgbotapi.BotAPI, log *logging.Logger) *RouterAdminPanel {
	return &RouterAdminPanel{
		b:   b,
		log: log,
	}
}

func (r *RouterAdminPanel) CheckAdminPanel(
	callback *tgbotapi.CallbackQuery,
) bool {
	if callback == nil {
		return false
	}

	return callback.Data == inline.DataAdminPanel
}

func (r *RouterAdminPanel) ShowAdminPanel(callback *tgbotapi.CallbackQuery) {
	newCallback := tgbotapi.NewCallback(callback.ID, "")
	_, err := r.b.Request(newCallback)
	if err != nil {
		r.log.Error("Failed to send callback", zap.Error(err))
	}

	msgSend := tgbotapi.NewEditMessageTextAndMarkup(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		inline.MsgDataAdminPanel, inline.AdminPanelKB,
	)

	_, err = r.b.Send(msgSend)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
	}
}

func (r *RouterAdminPanel) CheckBackToStartMenu(
	callback *tgbotapi.CallbackQuery,
) bool {
	if callback == nil {
		return false
	}

	return callback.Data == inline.DataBackToStartMenu
}
