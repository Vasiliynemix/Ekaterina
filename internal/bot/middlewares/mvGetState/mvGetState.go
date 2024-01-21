package mvGetState

import (
	"bot/internal/config"
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GetStateMv struct {
	log         *logging.Logger
	cfg         *config.Config
	stateGetter StateGetter
}

type StateGetter interface {
	GetState(telegramID int64, stateName string) (map[string]interface{}, error)
}

func New(
	log *logging.Logger,
	cfg *config.Config,
	stateGetter StateGetter,
) *GetStateMv {
	return &GetStateMv{
		log:         log,
		cfg:         cfg,
		stateGetter: stateGetter,
	}
}

func (l *GetStateMv) UserState(msg tgbotapi.Update) map[string]interface{} {
	var telegramID int64

	switch {
	case msg.Message != nil:
		telegramID = msg.Message.Chat.ID

	case msg.CallbackQuery != nil:
		telegramID = msg.CallbackQuery.Message.Chat.ID
	}

	userState, _ := l.stateGetter.GetState(telegramID, config.ScheduleState)

	if userState == nil {
		return nil
	}

	return userState
}
