package config

import "github.com/joho/godotenv"

type EnvConfig struct {
	Dir      string `env:"CFG_DIR" env-required:"true"`
	FileName string `env:"CFG_FILENAME" env-required:"true"`
}

func mustLoadEnvConfig() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func addEnvInConfig(cfg *Config) {

}
