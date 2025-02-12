package config

// import (
// 	"log"
// 	"os"
// 	"time"

// 	"github.com/kelseyhightower/envconfig"
// )

// type Config struct {
// 	Env         string `yaml:"env" env-default:"local"`
// 	StoragePath string `yaml:"storage_path" env-required:"true"`
// 	HTTPServer  `yaml:"http_server"`
// }

// type HTTPServer struct {
// 	Address     string        `yaml:"address" env-default:"localhost:8080"`
// 	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
// 	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
// 	User        string        `yaml:"user" env-required:"true"`
// 	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
// }

// func NewConfig() *Config {
// 	configPath := os.Getenv("CONFIG")
// 	if configPath == "" {
// 		log.Fatal("Не указан путь к конфигурационному файлу")
// 	}

// 	if _, err := os.Stat(configPath); os.IsNotExist(err) {
// 		log.Fatalf("Файл конфигурации %s не найден", configPath)
// 	}

// 	var cfg Config

// 	if err := envconfig.Process("", &cfg); err != nil {
// 		log.Fatalf("Ошибка при загрузке конфигурации: %s", err.Error())
// 	}

// 	return &cfg
// }

