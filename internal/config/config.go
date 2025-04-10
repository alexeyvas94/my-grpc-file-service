package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort  string // Порт для gRPC-сервера
	UploadDir string // Директория для загрузки файлов
}

// LoadConfig загружает конфигурацию из файла .env или использует значения по умолчанию
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Failed to load .env file, using default values")
	}

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051" // Порт по умолчанию
	}

	dir := os.Getenv("UPLOAD_DIR")
	if dir == "" {
		dir = "uploaded" // Директория по умолчанию
	}

	return &Config{
		GRPCPort:  port,
		UploadDir: dir,
	}
}
