package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Env      string       `yaml:"env" env-required:"true"`
	Debug    bool         `yaml:"debug"`
	Log      LoggerConfig `yaml:"logger"`
	Bot      BotConfig    `yaml:"bot"`
	DB       DBConfig     `yaml:"db"`
	RootPath string
	Paths    PathsConfig
}

type BotConfig struct {
	Token     string  `json:"-"`
	TimeOut   int     `yaml:"time_out" env-required:"true"`
	AdminsStr string  `yaml:"admins" env-required:"true" json:"-"`
	Admins    []int64 `yaml:"-" json:"-"`
}

type DBConfig struct {
	MigrationDirName string       `yaml:"migration_dir_name" env-required:"true"`
	Driver           string       `yaml:"driver" env-required:"true"`
	Password         string       `yaml:"password" json:"-"`
	SslMode          string       `yaml:"ssl_mode" env-required:"true"`
	Pool             DBPoolConfig `yaml:"pool"`
	Host             string
	Port             string
	User             string
	DbName           string
}

type DBPoolConfig struct {
	MaxIdleConns int           `yaml:"max_idle_conns"`
	MaxOpenConns int           `yaml:"max_open_conns"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type LoggerConfig struct {
	Dir           string `yaml:"dir" env-required:"true"`
	FIleInfoName  string `yaml:"file_info_name" env-required:"true"`
	FileDebugName string `yaml:"file_debug_name" env-required:"true"`
}

type PathsConfig struct {
	ConfigInfoPath  string
	ConfigDebugPath string
}

func (c *Config) ValidateEnv() error {
	switch c.Env {
	case "local", "prod":
		return nil
	default:
		return fmt.Errorf("invalid Env value: %s. Must be 'local' or 'prod'", c.Env)
	}
}

func (d *DBConfig) ConnString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		d.Host, d.Port, d.User, d.DbName, d.Password, d.SslMode,
	)
}

func MustLoad(levelsUp int) *Config {
	mustLoadEnvConfig()

	var cfgEnv EnvConfig

	err := cleanenv.ReadEnv(&cfgEnv)
	if err != nil {
		panic(err)
	}

	rootPath := getRootPath(levelsUp)

	pathToCfg := getPath(rootPath, cfgEnv.Dir, cfgEnv.FileName)

	return mustLoadCfg(rootPath, pathToCfg, &cfgEnv.Bot, &cfgEnv.DB)
}

func mustLoadCfg(rootPath string, pathToCfg string, botC *EnvBotConfig, db *EnvDBConfig) *Config {
	var cfg Config

	err := cleanenv.ReadConfig(pathToCfg, &cfg)
	if err != nil {
		panic(err)
	}

	err = cfg.ValidateEnv()
	if err != nil {
		panic(err)
	}

	cfg.RootPath = rootPath
	cfg.Paths.ConfigDebugPath = getPath(rootPath, cfg.Log.Dir, cfg.Log.FileDebugName)
	cfg.Paths.ConfigInfoPath = getPath(rootPath, cfg.Log.Dir, cfg.Log.FIleInfoName)

	cfg.Bot.Admins = parseAdmins(cfg.Bot.AdminsStr)

	addEnvInConfig(&cfg, botC, db)

	return &cfg
}

func parseAdmins(adminsStr string) []int64 {
	adminsSlice := strings.Split(adminsStr, ",")
	admins := make([]int64, 0, len(adminsSlice))

	for _, admin := range adminsSlice {
		adminID, err := strconv.ParseInt(strings.TrimSpace(admin), 10, 64)
		if err != nil {
			panic(err)
		}
		admins = append(admins, adminID)
	}

	return admins
}

func createPath(path string, fileName string) {
	_, err := os.Stat(path)
	dir := path
	if fileName != "" {
		dir = filepath.Dir(path)
	}
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
	}
}

func getPath(rootPath string, dir string, fileName string) string {
	path := filepath.Join(rootPath, dir, fileName)
	createPath(path, fileName)
	return path
}

func getRootPath(levelsUp int) string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("Failed to get root path")
	}

	parentPath := filename
	for i := 0; i < levelsUp; i++ {
		parentPath = filepath.Dir(parentPath)
	}
	return parentPath
}
