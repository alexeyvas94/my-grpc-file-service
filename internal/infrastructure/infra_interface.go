package infrastructure

import "my-grpc-file-service/internal/domain"

type FileStorageInterface interface {
	ListFiles() ([]*domain.FileInfo, error)
}
