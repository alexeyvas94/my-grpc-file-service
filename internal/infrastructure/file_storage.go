package infrastructure

import (
	"log"
	"my-grpc-file-service/internal/domain"
	"os"
	"path/filepath"
	"time"
)

type FileStorage struct {
	BasePath string
}

func NewFileStorage(path string) *FileStorage {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatalf("Ошибка при создании директории: %v", err)
	}

	return &FileStorage{BasePath: path}
}

func (fs *FileStorage) ListFiles() ([]*domain.FileInfo, error) {
	entries, err := os.ReadDir(fs.BasePath)
	if err != nil {
		return nil, err
	}

	var files []*domain.FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := os.Stat(filepath.Join(fs.BasePath, entry.Name()))
		if err != nil {
			continue
		}

		files = append(files, &domain.FileInfo{
			Name:    entry.Name(),
			Created: info.ModTime().Format(time.RFC3339),
			Updated: info.ModTime().Format(time.RFC3339),
		})
	}

	return files, nil
}
