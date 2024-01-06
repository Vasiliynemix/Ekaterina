package inline

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var ScheduleKB = tgbotapi.NewInlineKeyboardMarkup(
	//tgbotapi.NewInlineKeyboardRow(
	//	tgbotapi.NewInlineKeyboardButtonData(MsgDataSchedule, DataSchedule),
	//	tgbotapi.NewInlineKeyboardButtonData(MsgDataNotes, DataNotes),
	//),

	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MsgDataBack, DataBackToStartMenu),
	),
)
