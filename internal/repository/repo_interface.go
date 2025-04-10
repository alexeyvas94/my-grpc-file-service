package repository

import (
	"os"
)

type FileRepositoryInterface interface {
	SaveChunk(filename string, chunk []byte, create bool) error
	OpenFile(filename string) (*os.File, error)
}
