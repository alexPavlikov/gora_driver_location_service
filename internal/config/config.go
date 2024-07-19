package config

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Config struct {
	Env        string        `mapstructure:"env"`
	Timeout    time.Duration `mapstructure:"timeout"`
	Server     Server        `mapstructure:"server"`
	LogLevel   int           `mapstructure:"loglevel"`
	Kafka      Server        `mapstructure:"kafka"`
	KafkaTopic string        `mapstructure:"kafkatopic"`
}

type Server struct {
	Path string `mapstructure:"path"`
	Port int    `mapstructure:"port"`
}

func (s *Server) ToString() string {
	return fmt.Sprintf("%s:%d", s.Path, s.Port)
}

// Функция чтения конфиг файла
func Load() (*Config, error) {

	path, fileName := fetchConfigPath()

	if path == "" || fileName == "" {
		return nil, errors.New("path or filename is empty")
	}

	var cfg Config

	cfg, err := initViper(path, fileName, cfg)
	if err != nil {
		return nil, err
	}

	var level slog.Level = slog.Level(cfg.LogLevel)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	slog.SetDefault(log)

	return &cfg, nil
}

// Функция инициализации viper
func initViper(path string, filename string, cfg Config) (Config, error) {

	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	err := viper.ReadInConfig()
	if err != nil {
		return cfg, err
	}

	viper.Unmarshal(&cfg)

	return cfg, nil
}

// Функция для определения какой файл конфига читать (local, dev, prod)
func fetchConfigPath() (path string, file string) {

	flag.StringVar(&path, "config_path", "", "path to config file")
	flag.StringVar(&file, "config_file", "", "config file name")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	if file == "" {
		file = os.Getenv("CONFIG_FILE")
	}

	return path, file
}
