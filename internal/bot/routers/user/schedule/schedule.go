package schedule

import (
	"bot/internal/bot/keyboards/inline"
	"bot/internal/config"
	"bot/internal/storage/db/models"
	"bot/pkg/logging"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	msgNotExistsDaysForEvenWeek = "😥 Расписание для четной недели пока не добавлено"
	msgNotExistsDaysForOddWeek  = "😥 Расписание для нечетной недели пока не добавлено"
	msgTypeScheduleNotSelected  = "😥 Я вижу тип раписания не выбран.\n" +
		"Давай сперва выберем тип для удобства? 😋"
)

type RouterSchedule struct {
	b              *tgbotapi.BotAPI
	log            *logging.Logger
	cfg            *config.Config
	scheduleGetter GetterSchedule
	userProvider   UserProvider
	stateProvider  StateProvider
}

type GetterSchedule interface {
	GetScheduleByTelegramID(TelegramID int64) (models.Schedule, error)
	AddDay(TelegramID int64, schedule models.Schedule, dayName string, weekNum int) error
	CheckDayNameExistByTelegramID(TelegramID int64, dayName string, weekNum int) (models.Schedule, bool)
}

type UserProvider interface {
	UpdateTypeSchedule(TelegramID int64, typeSchedule string) error
}

type StateProvider interface {
	ClearState(telegramID int64, stateName string) error
	SetState(telegramID int64, stateName string, stateData *map[string]interface{}) error
	GetState(telegramID int64, stateName string) (map[string]interface{}, error)
	UpdateState(telegramID int64, stateName string, fieldName string, fieldValue interface{}) error
}

func New(
	b *tgbotapi.BotAPI,
	log *logging.Logger,
	cfg *config.Config,
	scheduleGetter GetterSchedule,
	userProvider UserProvider,
	stateProvider StateProvider,
) *RouterSchedule {
	return &RouterSchedule{
		b:              b,
		log:            log,
		cfg:            cfg,
		scheduleGetter: scheduleGetter,
		userProvider:   userProvider,
		stateProvider:  stateProvider,
	}
}

func (r *RouterSchedule) CheckSchedule(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}
	return callback.Data == inline.DataSchedule
}

func (r *RouterSchedule) CheckTypeSchedule(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}
	return strings.HasPrefix(callback.Data, inline.DataScheduleTypeFind)
}

func (r *RouterSchedule) TypeSchedule(callback *tgbotapi.CallbackQuery) {
	newCallback := tgbotapi.NewCallback(callback.ID, "")
	_, err := r.b.Request(newCallback)
	if err != nil {
		r.log.Error("Failed to send callback", zap.Error(err))
	}

	typeInfo := strings.Split(callback.Data, "|")
	typeSchedule := typeInfo[1]

	_ = r.userProvider.UpdateTypeSchedule(callback.From.ID, typeSchedule)

	msgSend := tgbotapi.NewEditMessageTextAndMarkup(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		inline.MsgDataSchedule,
		inline.ScheduleKB(typeSchedule),
	)

	_, err = r.b.Send(msgSend)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
	}
}

func (r *RouterSchedule) ShowSchedule(callback *tgbotapi.CallbackQuery, typeSchedule string) {
	newCallback := tgbotapi.NewCallback(callback.ID, "")
	_, err := r.b.Request(newCallback)
	if err != nil {
		r.log.Error("Failed to send callback", zap.Error(err))
	}

	if typeSchedule == "" {
		msgSend := tgbotapi.NewEditMessageTextAndMarkup(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			msgTypeScheduleNotSelected,
			inline.TypeScheduleKB,
		)
		_, err = r.b.Send(msgSend)
		if err != nil {
			r.log.Error("Failed to send message", zap.Error(err))
		}
		return
	}

	msgSend := tgbotapi.NewEditMessageTextAndMarkup(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		inline.MsgDataSchedule,
		inline.ScheduleKB(typeSchedule),
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

func (r *RouterSchedule) ScheduleWeekEven(callback *tgbotapi.CallbackQuery,
	typeSchedule string,
) {
	newCallback := tgbotapi.NewCallback(callback.ID, "")
	_, err := r.b.Request(newCallback)
	if err != nil {
		r.log.Error("Failed to send callback", zap.Error(err))
	}

	r.sendMsgWeek(
		callback,
		inline.MsgDataScheduleWeekEven,
		inline.ScheduleWeekKB(config.WeekEvenAndDefault, typeSchedule),
	)
}

func (r *RouterSchedule) CheckScheduleWeekOdd(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}
	return callback.Data == inline.DataScheduleWeekOdd
}

func (r *RouterSchedule) ScheduleWeekOdd(callback *tgbotapi.CallbackQuery,
	typeSchedule string,
) {
	newCallback := tgbotapi.NewCallback(callback.ID, "")
	_, err := r.b.Request(newCallback)
	if err != nil {
		r.log.Error("Failed to send callback", zap.Error(err))
	}

	r.sendMsgWeek(
		callback,
		inline.MsgDataScheduleWeekOdd,
		inline.ScheduleWeekKB(config.WeekOdd, typeSchedule),
	)
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

	dayInfo := strings.Split(callback.Data, "|")

	dayName := dayInfo[0]
	weekNum, err := strconv.Atoi(dayInfo[1])
	if err != nil {
		r.log.Error("Failed to convert week number", zap.Error(err))
		return
	}

	schedule, ok := r.scheduleGetter.CheckDayNameExistByTelegramID(callback.Message.Chat.ID, dayName, weekNum)
	if !ok {
		_ = r.scheduleGetter.AddDay(callback.Message.Chat.ID, schedule, dayName, weekNum)
	}
}

func (r *RouterSchedule) CheckAddScheduleWeek(callback *tgbotapi.CallbackQuery) bool {
	if callback == nil {
		return false
	}
	return strings.HasPrefix(callback.Data, inline.DataAddScheduleWeek)
}

func (r *RouterSchedule) AddScheduleWeek(callback *tgbotapi.CallbackQuery) {
	newCallback := tgbotapi.NewCallback(callback.ID, "")
	_, err := r.b.Request(newCallback)
	if err != nil {
		r.log.Error("Failed to send callback", zap.Error(err))
	}

	msgSend := tgbotapi.NewMessage(
		callback.Message.Chat.ID,
		"Добавь файл с расписанием на неделю.",
	)
	msgSend.ReplyMarkup = inline.CancelKB

	msg, err := r.b.Send(msgSend)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
	}

	startState := map[string]interface{}{
		"msg_old_id": callback.Message.MessageID,
		"msg_new_id": msg.MessageID,
	}

	_ = r.stateProvider.SetState(callback.Message.Chat.ID, config.ScheduleState, &startState)
}

func (r *RouterSchedule) CheckStateSchedule(stateSchedule map[string]interface{}, msg *tgbotapi.Message) bool {
	if stateSchedule == nil {
		return false
	}
	if msg.Document == nil {
		return false
	}
	return true
}

func (r *RouterSchedule) FileScheduleWeek(stateSchedule map[string]interface{}, typeSchedule string, msg *tgbotapi.Message) {
	file, err := r.b.GetFile(tgbotapi.FileConfig{FileID: msg.Document.FileID})
	if err != nil {
		r.log.Error("Failed to get file", zap.Error(err))
		return
	}

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", r.cfg.Bot.Token, file.FilePath)

	fileName := "schedule.pdf"

	err = downloadFile(r.log, fileURL, fileName)
	if err != nil {
		return
	}

	msgOldID := int(stateSchedule["msg_old_id"].(float64))
	msgNewID := int(stateSchedule["msg_new_id"].(float64))
	r.deleteMsg(msg.Chat.ID, msgNewID)
	r.deleteMsg(msg.Chat.ID, msgOldID)

	msgSend := tgbotapi.NewMessage(
		msg.Chat.ID,
		inline.MsgDataSchedule,
	)
	msgSend.ReplyMarkup = inline.ScheduleKB(typeSchedule)

	_, err = r.b.Send(msgSend)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
	}

	_ = r.stateProvider.ClearState(msg.Chat.ID, config.ScheduleState)
}

func downloadFile(log *logging.Logger, url, fileName string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(fileName)
	if err != nil {
		log.Error("Failed to create file", zap.Error(err))
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Error("Failed to copy file", zap.Error(err))
		return err
	}

	return nil
}

func (r *RouterSchedule) removeReplyMarkup(chatID int64, msgID int) {
	emptyKeyboard := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
	}

	editMsg := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, emptyKeyboard)
	_, err := r.b.Send(editMsg)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
	}
}

func (r *RouterSchedule) deleteMsg(chatID int64, msgID int) {
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, msgID)
	_, _ = r.b.Send(deleteMsg)
}
