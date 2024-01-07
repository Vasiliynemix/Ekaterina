package schedule

import (
	"bot/internal/bot/keyboards/inline"
	"bot/internal/db/models"
	"bot/pkg/logging"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strings"
)

const (
	msgNotExistsDaysForEvenWeek = "😥 Расписание для четной недели пока не добавлено"
	msgNotExistsDaysForOddWeek  = "😥 Расписание для нечетной недели пока не добавлено"
)

type RouterSchedule struct {
	b              *tgbotapi.BotAPI
	log            *logging.Logger
	scheduleGetter GetterSchedule
}

type GetterSchedule interface {
	GetScheduleByTelegramID(TelegramID int64) (models.Schedule, error)
}

func New(
	b *tgbotapi.BotAPI,
	log *logging.Logger,
	scheduleGetter GetterSchedule,
) *RouterSchedule {
	return &RouterSchedule{
		b:              b,
		log:            log,
		scheduleGetter: scheduleGetter,
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

func (r *RouterSchedule) CheckScheduleWeekEven(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}
	return callback.Data == inline.DataScheduleWeekEven
}

func (r *RouterSchedule) ScheduleWeekEven(callback *tgbotapi.CallbackQuery) {

	newCallback := tgbotapi.NewCallback(callback.ID, "")
	_, err := r.b.Request(newCallback)
	if err != nil {
		r.log.Error("Failed to send callback", zap.Error(err))
	}

	r.sendMsgWeek(callback, inline.MsgDataScheduleWeekEven, inline.ScheduleWeekKB(2))
}

func (r *RouterSchedule) CheckScheduleWeekOdd(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}
	return callback.Data == inline.DataScheduleWeekOdd
}

func (r *RouterSchedule) ScheduleWeekOdd(callback *tgbotapi.CallbackQuery) {
	newCallback := tgbotapi.NewCallback(callback.ID, "")
	_, err := r.b.Request(newCallback)
	if err != nil {
		r.log.Error("Failed to send callback", zap.Error(err))
	}

	r.sendMsgWeek(callback, inline.MsgDataScheduleWeekOdd, inline.ScheduleWeekKB(1))
}

func (r *RouterSchedule) CheckBackToScheduleMenu(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}
	return callback.Data == inline.DataBackToWeek
}

func (r *RouterSchedule) sendMsgWeek(callback *tgbotapi.CallbackQuery, msg string, kb tgbotapi.InlineKeyboardMarkup) {
	schedule, _ := r.scheduleGetter.GetScheduleByTelegramID(callback.Message.Chat.ID)

	var msgGetDays string

	switch {
	case callback.Data == inline.DataScheduleWeekEven:
		if schedule.WeekEven.Days == nil {
			msgGetDays = msgNotExistsDaysForEvenWeek
		}
	case callback.Data == inline.DataScheduleWeekOdd:
		if schedule.WeekOdd.Days == nil {
			msgGetDays = msgNotExistsDaysForOddWeek
		}
	}

	var msgText string
	if msgGetDays == "" {
		msgText = msg
	} else {
		msgText = fmt.Sprintf("%s\n\n%s", msg, msgGetDays)
	}

	msgSend := tgbotapi.NewEditMessageTextAndMarkup(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		msgText,
		kb,
	)

	_, err := r.b.Send(msgSend)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
	}
}

func (r *RouterSchedule) CheckDayMonday(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}
	return strings.HasPrefix(callback.Data, inline.DataScheduleMonday)
}

func (r *RouterSchedule) CheckDayTuesday(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}
	return strings.HasPrefix(callback.Data, inline.DataScheduleTuesday)
}

func (r *RouterSchedule) CheckDayWednesday(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}
	return strings.HasPrefix(callback.Data, inline.DataScheduleWednesday)
}

func (r *RouterSchedule) CheckDayThursday(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}
	return strings.HasPrefix(callback.Data, inline.DataScheduleThursday)
}

func (r *RouterSchedule) CheckDayFriday(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}
	return strings.HasPrefix(callback.Data, inline.DataScheduleFriday)
}

func (r *RouterSchedule) CheckDaySaturday(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}
	return strings.HasPrefix(callback.Data, inline.DataScheduleSaturday)
}

func (r *RouterSchedule) CheckDaySunday(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}
	return strings.HasPrefix(callback.Data, inline.DataScheduleSunday)
}

func (r *RouterSchedule) ShowDay(callback *tgbotapi.CallbackQuery) {
	newCallback := tgbotapi.NewCallback(callback.ID, callback.Data)
	_, err := r.b.Request(newCallback)
	if err != nil {
		r.log.Error("Failed to send callback", zap.Error(err))
	}
}
