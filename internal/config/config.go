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
	Env     string        `yaml:"env"`
	Timeout time.Duration `yaml:"timeout"`

	ServerPath string `yaml:"server_path"`
	ServerPort int    `yaml:"server_port"`

	Kafka Kafka `yaml:"kafka"`
}

type Kafka struct {
	Path string `yaml:"kafka_path"`
	Port int    `yaml:"kafka_port"`
}

// Функция чтения конфиг файла
func Load() (*Config, error) {

	path, fileName := fetchConfigPath()

	if path == "" || fileName == "" {
		return nil, errors.New("path or filename is empty")
	}

	var cfg Config

	err := initViper(path, fileName)
	if err != nil {
		return nil, err
	}

	cfg.Env = viper.GetString("env")
	cfg.Timeout = viper.GetDuration("timeout")
	cfg.ServerPort = viper.GetInt("server_port")
	cfg.ServerPath = viper.GetString("server_path")
	cfg.Kafka.Port = viper.GetInt("kafka.kafka_port")
	cfg.Kafka.Path = viper.GetString("kafka.kafka_path")

	return &cfg, nil
}

// Функция инициализации viper
func initViper(path string, filename string) error {

	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
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
func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		slog.SetDefault(log)
	}

	return log
}
