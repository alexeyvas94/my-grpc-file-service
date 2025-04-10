package service

import "my-grpc-file-service/internal/domain"

type FileServiceInterface interface {
	SaveFileStream(filename string, stream func() ([]byte, error)) error
	StreamFileToClient(filename string, send func([]byte) error) error
	ListFiles() ([]*domain.FileInfo, error)
}
