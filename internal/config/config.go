package config

import (
	"errors"
	"flag"
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
	Env        string        `yaml:"env"`
	Timeout    time.Duration `yaml:"timeout"`
	ServerPath string        `yaml:"serverpath"`
	ServerPort int           `yaml:"serverport"`
	LogLevel   string        `yaml:"loglevel"`
	Kafka      Kafka         `yaml:"kafka"`
}

type Kafka struct {
	Path string `yaml:"kafka.path"`
	Port int    `yaml:"kafka.port"`
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

// Функция для инициализации слоя логирования
func SetupLogger(logLevel string) *slog.Logger {
	var log *slog.Logger

	switch logLevel {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	slog.SetDefault(log)

	return log
}
