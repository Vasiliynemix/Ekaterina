package inline

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	DataModerators    = "moderators"
	MsgDataModerators = "Модераторы"

	DataBackToStartMenu = "back_start_menu"
	MsgDataBack         = "⏪ Назад"
)

var AdminPanelKB = tgbotapi.NewInlineKeyboardMarkup(
	//tgbotapi.NewInlineKeyboardRow(
	//	tgbotapi.NewInlineKeyboardButtonData("Admin panel", DataAdminPanel),
	//	tgbotapi.NewInlineKeyboardButtonData(MsgDataNotes, DataNotes),
	//),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MsgDataModerators, DataModerators),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MsgDataBack, DataBackToStartMenu),
	),
)
