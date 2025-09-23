package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

// Config - стукрутра для конфигурации всего приложения
// 1) Engine - параметры движка
// 2) Network - параметры сети
// 3) logging - параметры логирования
// Так же есть вычисленные поля(maxMessageBytes, IdleTimeoutDur) что бы не парсить их каждый раз
type Config struct {
	Engine struct {
		Type string `yaml:"type"`
	} `yaml:"engine"`

	Network struct {
		Address        string `yaml:"address"`
		MaxConnections int    `yaml:"max_connections"`
		MaxMessageSize string `yaml:"max_message_size"`
		IdleTimeout    string `yaml:"idle_timeout"`
	} `yaml:"network"`

	Logging struct {
		Level  string `'yaml:"level"`
		Output string `'yaml:"output"`
	} `yaml:"logging"`

	MaxMessageBytes int
	IdleTimeoutDur  time.Duration
}

// Возвращает конфигурацию с дефолтными значениями
// Используется если файл конфигурации не задан или часть полей пуста
func DefaultConfig() Config {
	var defaultConfig Config

	defaultConfig.Engine.Type = "in_memory"

	defaultConfig.Network.Address = "localhost:8080"
	defaultConfig.Network.MaxConnections = 10
	defaultConfig.Network.MaxMessageSize = "4KB"
	defaultConfig.Network.IdleTimeout = "5m"

	defaultConfig.Logging.Level = "info"
	defaultConfig.Logging.Output = "stdout"

	defaultConfig.MaxMessageBytes = 4 * 1024
	defaultConfig.IdleTimeoutDur = time.Minute * 5

	return defaultConfig
}

func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	// Парсинг yaml файла поверх дефолтов
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return cfg, err
	}

	//Парсинг строки размера сообщения в байты
	if cfg.Network.MaxMessageSize != "" {
		b, err := parseSize(cfg.Network.MaxMessageSize)
		if err != nil {
			return cfg, fmt.Errorf("invalid max_message_size: %w", err)
		}
		cfg.MaxMessageBytes = b
	} else {
		cfg.MaxMessageBytes = 4 * 1024
	}

	//Парсинг таймаута в Duration
	if cfg.Network.IdleTimeout != "" {
		d, err := time.ParseDuration(cfg.Network.IdleTimeout)
		if err != nil {
			return cfg, fmt.Errorf("invalid idle_timeout: %w", err)
		}
		cfg.IdleTimeoutDur = d
	} else {
		cfg.IdleTimeoutDur = time.Minute * 5
	}

	if cfg.Network.MaxConnections <= 0 {
		cfg.Network.MaxConnections = 100
	}
	if cfg.Network.Address == "" {
		cfg.Network.Address = "127.0.0.1:3223"
	}
	if cfg.Engine.Type == "" {
		cfg.Engine.Type = "in_memory"
	}
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}

	return cfg, nil
}

// parseSize разбирает строку вида "4KB", "1MB" или просто "123" в количество байт.
func parseSize(s string) (int, error) {
	var n int
	var unit string
	// Попытка парсинга числа и единицы измерения
	_, err := fmt.Sscanf(s, "%d%s", &n, &unit)
	if err != nil {
		// если просто число без единицы
		var only int
		_, err2 := fmt.Sscanf(s, "%d", &only)
		if err2 == nil {
			return only, nil
		}
		return 0, err
	}
	switch unit {
	case "K", "k", "KB", "kb", "Kb":
		return n * 1024, nil
	case "M", "m", "MB", "mb", "Mb":
		return n * 1024 * 1024, nil
	case "B", "b":
		return n, nil
	default:
		return 0, errors.New("unknown size unit")
	}
}
