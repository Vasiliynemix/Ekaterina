package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	RootPath string
	Env      string       `yaml:"env" env-required:"true"`
	Log      LoggerConfig `yaml:"logger"`
	Paths    PathsConfig
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

func MustLoad(levelsUp int) *Config {
	mustLoadEnvConfig()

	var cfgEnv EnvConfig

	err := cleanenv.ReadEnv(&cfgEnv)
	if err != nil {
		panic(err)
	}

	rootPath := getRootPath(levelsUp)

	pathToCfg := getPath(rootPath, cfgEnv.Dir, cfgEnv.FileName)

	return mustLoadCfg(rootPath, pathToCfg)
}

func mustLoadCfg(rootPath string, pathToCfg string) *Config {
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

	addEnvInConfig(&cfg)

	return &cfg
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
