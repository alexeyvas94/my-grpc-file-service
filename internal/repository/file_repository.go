package repository

import (
	"log"
	"os"
	"path/filepath"
)

type FileRepository struct {
	BasePath string
}

// NewFileRepository создаёт новый объект FileRepository и создаёт директорию, если её ещё нет
func NewFileRepository(path string) *FileRepository {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatalf("Ошибка при создании директории: %v", err)
	}
	return &FileRepository{BasePath: path}
}

func (r *FileRepository) SaveChunk(filename string, chunk []byte, create bool) error {
	path := filepath.Join(r.BasePath, filename)
	var f *os.File
	var err error

	if create {
		f, err = os.Create(path)
	} else {
		f, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	}
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(chunk)
	return err
}

func (r *FileRepository) OpenFile(filename string) (*os.File, error) {
	return os.Open(filepath.Join(r.BasePath, filename))
}
