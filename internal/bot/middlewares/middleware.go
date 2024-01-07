package middlewares

import (
	"bot/internal/bot/middlewares/mvAddToDB"
	"bot/internal/bot/middlewares/mvLogger"
	"bot/internal/config"
	"bot/internal/db"
	"bot/pkg/logging"
)

type Middlewares struct {
	MvLog     *mvLogger.LoggerMv
	MvAddToDB *mvAddToDB.AddToDBMv
}

func InitMiddlewares(log *logging.Logger, db *db.DB, cfg *config.Config) *Middlewares {
	return &Middlewares{
		MvLog:     mvLogger.New(log),
		MvAddToDB: mvAddToDB.New(log, db.User, db.Schedule, cfg),
	}
}
