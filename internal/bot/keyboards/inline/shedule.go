package inline

import (
	"bot/internal/config"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	MsgDataMainMenu = "🏠 Главное меню"
	DataMainMenu    = "main_menu"

	MsgDataCancel = "❌ Отмена"
	DataCancel    = "cancel"

	DataBackToWeek = "back_to_week"

	MsgDataScheduleWeekEven = "⚪️ Четная неделя"
	DataScheduleWeekEven    = "week_even"

	MsgDataScheduleWeekOdd = "⚫️ Нечетная неделя"
	DataScheduleWeekOdd    = "week_odd"

	MsgDataScheduleTypeEvenOdd = "🤔 Четное, нечетное расписание"
	MsgDataScheduleTypeDefault = "🤔 Еженедельное расписание"
	DataScheduleTypeFind       = "type_schedule"

	MsgDataScheduleMonday    = "📚 Понедельник"
	MsgDataScheduleTuesday   = "🌟 Вторник"
	MsgDataScheduleWednesday = "🍃 Среда"
	MsgDataScheduleThursday  = "📚 Четверг"
	MsgDataScheduleFriday    = "🎉 Пятница"
	MsgDataScheduleSaturday  = "🎈 Суббота"
	MsgDataScheduleSunday    = "☀️ Воскресенье"

	MsgDataAddScheduleWeek = "🌎 Добавить расписание"
	DataAddScheduleWeek    = "add_schedule_week"

	DataScheduleMonday    = "monday"
	DataScheduleTuesday   = "tuesday"
	DataScheduleWednesday = "wednesday"
	DataScheduleThursday  = "thursday"
	DataScheduleFriday    = "friday"
	DataScheduleSaturday  = "saturday"
	DataScheduleSunday    = "sunday"
)

var (
	DataScheduleTypeEvenOdd = fmt.Sprintf("%s|%s", DataScheduleTypeFind, config.TypeScheduleEvenOdd)
	DataScheduleTypeDefault = fmt.Sprintf("%s|%s", DataScheduleTypeFind, config.TypeScheduleDefault)
)

var CancelKB = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MsgDataCancel, DataCancel),
	),
)

var TypeScheduleKB = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MsgDataScheduleTypeDefault, DataScheduleTypeDefault),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MsgDataScheduleTypeEvenOdd, DataScheduleTypeEvenOdd),
	),
)

var scheduleKBEvenOdd = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MsgDataScheduleWeekEven, DataScheduleWeekEven),
		tgbotapi.NewInlineKeyboardButtonData(MsgDataScheduleWeekOdd, DataScheduleWeekOdd),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MsgDataBack, DataBackToStartMenu),
	),
)

func ScheduleKB(typeSchedule string) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	switch {
	case typeSchedule == config.TypeScheduleEvenOdd:
		return scheduleKBEvenOdd
	case typeSchedule == config.TypeScheduleDefault:
		return ScheduleWeekKB(config.WeekEvenAndDefault, typeSchedule)
	}

	return keyboard
}

func ScheduleWeekKB(weekNum int, typeSchedule string) tgbotapi.InlineKeyboardMarkup {
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

	var kbBack []tgbotapi.InlineKeyboardButton
	var kbAddScheduleWeek []tgbotapi.InlineKeyboardButton

	switch {
	case typeSchedule == config.TypeScheduleEvenOdd:
		kbAddScheduleWeek = tgbotapi.NewInlineKeyboardRow(
			createButtonData(MsgDataAddScheduleWeek, DataAddScheduleWeek, weekNum),
		)

		kbBack = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(MsgDataBack, DataBackToWeek),
		)
	case typeSchedule == config.TypeScheduleDefault:
		kbAddScheduleWeek = tgbotapi.NewInlineKeyboardRow(
			createButtonData(MsgDataAddScheduleWeek, DataAddScheduleWeek, weekNum),
		)

		kbBack = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(MsgDataBack, DataBackToStartMenu),
		)
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, kbAddScheduleWeek)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, kbBack)

	kbMainMenu := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MsgDataMainMenu, DataMainMenu),
	)

	if typeSchedule != config.TypeScheduleDefault {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, kbMainMenu)
	}

	return keyboard
}

func createButtonData(msg, data string, weekNum int) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(
		msg,
		fmt.Sprintf("%s|%d", data, weekNum),
	)
}
