package middlewares

import (
	"bot/internal/bot/middlewares/mvLogger"
	"bot/pkg/logging"
)

type Middlewares struct {
	MvLog *mvLogger.LoggerMv
}

func InitMiddlewares(log *logging.Logger) *Middlewares {
	return &Middlewares{
		MvLog: mvLogger.New(log),
	}
}
