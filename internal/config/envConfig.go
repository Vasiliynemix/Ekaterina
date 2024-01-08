package config

import (
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Dir      string `env:"CFG_DIR" env-required:"true"`
	FileName string `env:"CFG_FILENAME" env-required:"true"`
	Bot      EnvBotConfig
	DB       EnvDBConfig
	MongoDB  EnvMongoDBConfig
}

type EnvDBConfig struct {
	Host string `env:"DB_HOST" env-required:"true"`
	Port string `env:"DB_PORT" env-required:"true"`
	User string `env:"DB_USER" env-required:"true"`
	Name string `env:"DB_NAME" env-required:"true"`
	Pass string `env:"DB_PASSWORD" env-required:"true" json:"-"`
}

type EnvMongoDBConfig struct {
	Host string `env:"MONGO_HOST" env-required:"true"`
	Port int    `env:"MONGO_PORT" env-required:"true"`
	User string `env:"MONGO_USER" env-required:"true"`
	Name string `env:"MONGO_NAME" env-required:"true"`
	Pass string `env:"MONGO_PASSWORD" env-required:"true" json:"-"`
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

func addEnvInConfig(cfg *Config, botC *EnvBotConfig, db *EnvDBConfig, mongoDB *EnvMongoDBConfig) {
	cfg.Bot.Token = botC.Token

	cfg.DB.Host = db.Host
	cfg.DB.Port = db.Port
	cfg.DB.User = db.User
	cfg.DB.DbName = db.Name
	cfg.DB.Password = db.Pass

	cfg.MongoDB.Host = mongoDB.Host
	cfg.MongoDB.Port = mongoDB.Port
	cfg.MongoDB.User = mongoDB.User
	cfg.MongoDB.DbName = mongoDB.Name
	cfg.MongoDB.Password = mongoDB.Pass
}
