package inline

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	DataSchedule    = "schedule"
	MsgDataSchedule = "📅 Расписание"

	DataNotes    = "notes"
	MsgDataNotes = "📝 Заметки"
)

var StartKB = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MsgDataSchedule, DataSchedule),
		tgbotapi.NewInlineKeyboardButtonData(MsgDataNotes, DataNotes),
	),
)
