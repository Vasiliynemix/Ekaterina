package middlewares

import (
	"bot/internal/bot/middlewares/mvAddToDB"
	"bot/internal/bot/middlewares/mvGetState"
	"bot/internal/bot/middlewares/mvLogger"
	"bot/internal/config"
	"bot/internal/storage/db"
	"bot/internal/storage/redisdb"
	"bot/pkg/logging"
)

type Middlewares struct {
	MvLog      *mvLogger.LoggerMv
	MvAddToDB  *mvAddToDB.AddToDBMv
	MvGetState *mvGetState.GetStateMv
}

func InitMiddlewares(log *logging.Logger, db *db.DB, cfg *config.Config, redis *redisdb.RedisDB) *Middlewares {
	return &Middlewares{
		MvLog:      mvLogger.New(log),
		MvAddToDB:  mvAddToDB.New(log, db.User, db.Schedule, cfg),
		MvGetState: mvGetState.New(log, cfg, redis),
	}
}
