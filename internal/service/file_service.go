package service

import (
	"io"
	"log"
	"my-grpc-file-service/internal/domain"
	"my-grpc-file-service/internal/infrastructure"
	"my-grpc-file-service/internal/repository"
)

type FileService struct {
	Repo    repository.FileRepositoryInterface
	Storage infrastructure.FileStorageInterface
}

func NewFileService(repo repository.FileRepositoryInterface, storage infrastructure.FileStorageInterface) *FileService {
	return &FileService{
		Repo:    repo,
		Storage: storage,
	}
}

func (s *FileService) SaveFileStream(filename string, stream func() ([]byte, error)) error {
	firstChunk := true
	for {
		chunk, err := stream()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("❌ error receiving chunk for file %s: %v", filename, err)
			return err
		}
		if err := s.Repo.SaveChunk(filename, chunk, firstChunk); err != nil {
			log.Printf("❌ error saving chunk to file %s: %v", filename, err)
			return err
		}
		firstChunk = false
	}
}

func (s *FileService) StreamFileToClient(filename string, send func([]byte) error) error {
	f, err := s.Repo.OpenFile(filename)
	if err != nil {
		log.Printf("❌ failed to open file %s: %v", filename, err)
		return err
	}
	defer f.Close()

	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("❌ error reading file %s: %v", filename, err)
			return err
		}

		if err := send(buf[:n]); err != nil {
			log.Printf("❌ error sending chunk for file %s: %v", filename, err)
			return err
		}
	}

	return nil
}

func (s *FileService) ListFiles() ([]*domain.FileInfo, error) {
	return s.Storage.ListFiles()
}
