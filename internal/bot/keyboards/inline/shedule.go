package inline

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	MsgDataMainMenu = "🏠 Главное меню"
	DataMainMenu    = "main_menu"

	DataBackToWeek = "back_to_week"

	MsgDataScheduleWeekEven = "⚪️ Четная неделя"
	MsgDataScheduleWeekOdd  = "⚫️ Нечетная неделя"

	DataScheduleWeekEven = "week_even"
	DataScheduleWeekOdd  = "week_odd"

	MsgDataScheduleMonday    = "📅 Понедельник"
	MsgDataScheduleTuesday   = "🌟 Вторник"
	MsgDataScheduleWednesday = "🍃 Среда"
	MsgDataScheduleThursday  = "📚 Четверг"
	MsgDataScheduleFriday    = "🎉 Пятница"
	MsgDataScheduleSaturday  = "🎈 Суббота"
	MsgDataScheduleSunday    = "☀️ Воскресенье"

	DataScheduleMonday    = "monday"
	DataScheduleTuesday   = "tuesday"
	DataScheduleWednesday = "wednesday"
	DataScheduleThursday  = "thursday"
	DataScheduleFriday    = "friday"
	DataScheduleSaturday  = "saturday"
	DataScheduleSunday    = "sunday"
)

var ScheduleKB = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MsgDataScheduleWeekEven, DataScheduleWeekEven),
		tgbotapi.NewInlineKeyboardButtonData(MsgDataScheduleWeekOdd, DataScheduleWeekOdd),
	),

	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MsgDataBack, DataBackToStartMenu),
	),
)

func createButtonData(msg, data string, weekNum int) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(
		msg,
		fmt.Sprintf("%s|%d", data, weekNum),
	)
}

func ScheduleWeekKB(weekNum int) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()

	days := []struct {
		msg  string
		data string
	}{
		{MsgDataScheduleMonday, DataScheduleMonday},
		{MsgDataScheduleTuesday, DataScheduleTuesday},
		{MsgDataScheduleWednesday, DataScheduleWednesday},
		{MsgDataScheduleThursday, DataScheduleThursday},
		{MsgDataScheduleFriday, DataScheduleFriday},
		{MsgDataScheduleSaturday, DataScheduleSaturday},
		{MsgDataScheduleSunday, DataScheduleSunday},
	}

	var buttons []tgbotapi.InlineKeyboardButton
	for _, day := range days {
		button := createButtonData(day.msg, day.data, weekNum)
		buttons = append(buttons, button)

		if len(buttons) == 2 {
			row := tgbotapi.NewInlineKeyboardRow(buttons...)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			buttons = []tgbotapi.InlineKeyboardButton{}
		}
	}

	if len(buttons) > 0 {
		row := tgbotapi.NewInlineKeyboardRow(buttons...)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	kbBack := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MsgDataBack, DataBackToWeek),
	)

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, kbBack)

	kbMainMenu := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MsgDataMainMenu, DataMainMenu),
	)

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, kbMainMenu)

	return keyboard
}
