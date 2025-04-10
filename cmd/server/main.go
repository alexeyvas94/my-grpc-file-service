package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"my-grpc-file-service/internal/api"
	"my-grpc-file-service/internal/config"
	"my-grpc-file-service/internal/infrastructure"
	"my-grpc-file-service/internal/repository"
	"my-grpc-file-service/internal/service"
	pb "my-grpc-file-service/pkg/api"
)

func main() {
	cfg := config.LoadConfig()

	// Создание репозитория и хранилища через интерфейсы
	repo := repository.NewFileRepository(cfg.UploadDir)
	storage := infrastructure.NewFileStorage(cfg.UploadDir)

	// Инициализация сервисного слоя с интерфейсами
	svc := service.NewFileService(repo, storage)

	// Запуск gRPC сервера
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("❌ Не удалось запустить прослушивание порта %s: %v", cfg.GRPCPort, err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Инициализация обработчика
	handler := api.NewGRPCHandler(svc)

	pb.RegisterFileServiceServer(grpcServer, handler)

	log.Printf("✅ gRPC сервер запущен на порту %s...\n", cfg.GRPCPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Ошибка при запуске сервера: %v", err)
	}
}
