package config

import (
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Dir      string `env:"CFG_DIR" env-required:"true"`
	FileName string `env:"CFG_FILENAME" env-required:"true"`
	Bot      EnvBotConfig
	DB       EnvDBConfig
	RedisDB  EnvRedisConfig
}

type EnvDBConfig struct {
	Host string `env:"DB_HOST" env-required:"true"`
	Port string `env:"DB_PORT" env-required:"true"`
	User string `env:"DB_USER" env-required:"true"`
	Name string `env:"DB_NAME" env-required:"true"`
	Pass string `env:"DB_PASSWORD" env-required:"true" json:"-"`
}

type EnvRedisConfig struct {
	Host string `env:"REDIS_HOST" env-required:"true"`
	Port string `env:"REDIS_PORT" env-required:"true"`
	Name int    `env:"REDIS_DB" env-required:"true"`
}

type EnvBotConfig struct {
	Token string `env:"BOT_TOKEN" env-required:"true" json:"-"`
}

func mustLoadEnvConfig() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func addEnvInConfig(cfg *Config, botC *EnvBotConfig, db *EnvDBConfig, redisDB *EnvRedisConfig) {
	cfg.Bot.Token = botC.Token

	cfg.DB.Host = db.Host
	cfg.DB.Port = db.Port
	cfg.DB.User = db.User
	cfg.DB.DbName = db.Name
	cfg.DB.Password = db.Pass

	cfg.RedisDB.Host = redisDB.Host
	cfg.RedisDB.Port = redisDB.Port
	cfg.RedisDB.DB = redisDB.Name
}
