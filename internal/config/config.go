package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
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
func MustLoad() *Config {

	path := fetchConfigPath()

	if path == "" {
		panic("path to config file is empty")
	}

	_, err := os.Stat(path)
	if err != nil {
		panic("path to config file not exist")
	}

	var cfg Config

	err = cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic("failed to read config" + err.Error())
	}

	return &cfg
}

// Функция для определения какой файл конфига читать (local, dev, prod)
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
