package schedule

import (
	"bot/internal/bot/keyboards/inline"
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type RouterSchedule struct {
	b   *tgbotapi.BotAPI
	log *logging.Logger
}

func New(b *tgbotapi.BotAPI, log *logging.Logger) *RouterSchedule {
	return &RouterSchedule{
		b:   b,
		log: log,
	}
}

func (r *RouterSchedule) CheckSchedule(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}
	return callback.Data == inline.DataSchedule
}

func (r *RouterSchedule) ShowSchedule(callback *tgbotapi.CallbackQuery) {
	newCallback := tgbotapi.NewCallback(callback.ID, "")
	_, err := r.b.Request(newCallback)
	if err != nil {
		r.log.Error("Failed to send callback", zap.Error(err))
	}

	msgSend := tgbotapi.NewEditMessageTextAndMarkup(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		inline.MsgDataSchedule, inline.ScheduleKB,
	)

	_, err = r.b.Send(msgSend)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
	}
}
