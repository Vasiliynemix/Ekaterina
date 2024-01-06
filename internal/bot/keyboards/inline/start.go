package inline

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	DataSchedule    = "schedule"
	MsgDataSchedule = "📅 Расписание"

	DataNotes    = "notes"
	MsgDataNotes = "📝 Заметки"

	DataAdminPanel    = "admin_panel"
	MsgDataAdminPanel = "👥 Админка"

	DataModerPanel    = "moderator_panel"
	MsgDataModerPanel = "👥 Модераторка"
)

func StartKB(isAdmin bool, isModer bool) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()

	row1 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MsgDataSchedule, DataSchedule),
		tgbotapi.NewInlineKeyboardButtonData(MsgDataNotes, DataNotes),
	)

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row1)

	if isAdmin {
		adminButton := tgbotapi.NewInlineKeyboardButtonData(MsgDataAdminPanel, DataAdminPanel)
		row2 := tgbotapi.NewInlineKeyboardRow(adminButton)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row2)
	}

	if isModer {
		modButton := tgbotapi.NewInlineKeyboardButtonData(MsgDataModerPanel, DataModerPanel)
		row3 := tgbotapi.NewInlineKeyboardRow(modButton)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row3)
	}

	return keyboard
}
