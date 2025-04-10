package service_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"my-grpc-file-service/internal/infrastructure"
	"my-grpc-file-service/internal/repository"
	"my-grpc-file-service/internal/service"
)

func TestFileService(t *testing.T) {
	// Создаём временную директорию
	tempDir := t.TempDir()

	// Создаём репозиторий и хранилище на основе временной папки
	repo := repository.NewFileRepository(tempDir)
	storage := infrastructure.NewFileStorage(tempDir)

	// Создаём сервис
	svc := service.NewFileService(repo, storage)

	tests := []struct {
		name     string
		filename string
		content  []byte
	}{
		{
			name:     "Простой текстовый файл",
			filename: "test.txt",
			content:  []byte("hello test data"),
		},
		{
			name:     "Пустой файл",
			filename: "empty.txt",
			content:  []byte{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Загружаем файл (стримом)
			data := tc.content
			read := false
			err := svc.SaveFileStream(tc.filename, func() ([]byte, error) {
				if read {
					return nil, io.EOF
				}
				read = true
				return data, nil
			})
			if err != nil {
				t.Fatalf("Ошибка при сохранении файла: %v", err)
			}

			// Проверяем, что файл существует
			fullPath := filepath.Join(tempDir, tc.filename)
			if _, err := os.Stat(fullPath); err != nil {
				t.Fatalf("Файл %s не найден: %v", tc.filename, err)
			}

			// Читаем файл стримом
			var out bytes.Buffer
			err = svc.StreamFileToClient(tc.filename, func(chunk []byte) error {
				out.Write(chunk)
				return nil
			})
			if err != nil {
				t.Fatalf("Ошибка при чтении файла: %v", err)
			}

			// Проверяем содержимое
			if !bytes.Equal(out.Bytes(), tc.content) {
				t.Fatalf("Ожидалось: %q, Получено: %q", tc.content, out.Bytes())
			}
		})
	}
}
