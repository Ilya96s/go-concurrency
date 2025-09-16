package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"in-memory-kv/internal/compute"
	"in-memory-kv/internal/config"
	"in-memory-kv/internal/logger"
	"in-memory-kv/internal/network"
	"in-memory-kv/internal/storage"
	"os"
)

func main() {
	// путь к yaml конфигу через аргумент командной строки
	configPath := flag.String("config", "configs/config.yaml", "path to YAML config file")
	flag.Parse()

	// Загрузка конфига
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Println("error loading config:", err)
		os.Exit(1)
	}

	// Создание логгера
	log, err := logger.NewLogger(cfg)
	if err != nil {
		fmt.Println("error creating logger:", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Создание хранилища
	var engine storage.Engine
	switch cfg.Engine.Type {
	case "in_memory":
		engine = storage.NewMemoryEngine()
	default:
		log.Fatal("unknown engine type", zap.String("engine", cfg.Engine.Type))
	}

	// Создание compute (обработчик команд)
	comp := compute.NewCompute(engine, log)

	// Создание и запуск сервера
	server := network.NewServer(cfg, *comp, log)

	log.Info("starting server", zap.String("address", cfg.Network.Address))
	if err := server.StartServer(); err != nil {
		log.Fatal("server stopped", zap.Error(err))
	}
}
