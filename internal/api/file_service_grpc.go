package api

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"my-grpc-file-service/internal/service"
	pb "my-grpc-file-service/pkg/api"
)

var (
	uploadLimiter = make(chan struct{}, 10)
	listLimiter   = make(chan struct{}, 100)
)

type grpcHandler struct {
	pb.UnimplementedFileServiceServer
	service service.FileServiceInterface
}

// NewGRPCHandler - создает новый обработчик для gRPC
func NewGRPCHandler(s service.FileServiceInterface) pb.FileServiceServer {
	return &grpcHandler{service: s}
}

// Upload - обрабатывает запросы на загрузку файлов
func (h *grpcHandler) Upload(stream pb.FileService_UploadServer) error {
	if !acquire(uploadLimiter) {
		return status.Error(codes.ResourceExhausted, "Достигнут лимит загрузки")
	}
	defer release(uploadLimiter)

	req, err := stream.Recv()
	if err != nil {
		return err
	}
	filename := req.GetFilename()
	firstChunk := req.GetChunk()

	firstChunkSent := false

	err = h.service.SaveFileStream(filename, func() ([]byte, error) {
		if !firstChunkSent {
			firstChunkSent = true
			return firstChunk, nil
		}
		req, err := stream.Recv()
		if err != nil {
			return nil, err
		}
		return req.GetChunk(), nil
	})
	if err != nil {
		return err
	}

	return stream.SendAndClose(&pb.FileUploadResponse{Message: "Загрузка успешна"})
}

// Download - обрабатывает запросы на скачивание файлов
func (h *grpcHandler) Download(req *pb.FileDownloadRequest, stream pb.FileService_DownloadServer) error {
	if !acquire(uploadLimiter) {
		return status.Error(codes.ResourceExhausted, "Достигнут лимит скачивания")
	}
	defer release(uploadLimiter)

	return h.service.StreamFileToClient(req.GetFilename(), func(chunk []byte) error {
		return stream.Send(&pb.FileDownloadResponse{Chunk: chunk})
	})
}

// ListFiles - возвращает список файлов
func (h *grpcHandler) ListFiles(ctx context.Context, _ *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	if !acquire(listLimiter) {
		return nil, status.Error(codes.ResourceExhausted, "Достигнут лимит на запросы списка файлов")
	}
	defer release(listLimiter)

	data, err := h.service.ListFiles()
	if err != nil {
		return nil, err
	}

	var result []*pb.FileInfo
	for _, f := range data {
		result = append(result, &pb.FileInfo{
			Name:    f.Name,
			Created: f.Created,
			Updated: f.Updated,
		})
	}

	return &pb.ListFilesResponse{Files: result}, nil
}

// acquire - захватывает семафор для ограничения параллельных запросов
func acquire(sem chan struct{}) bool {
	select {
	case sem <- struct{}{}:
		return true
	default:
		return false
	}
}

// release - освобождает семафор
func release(sem chan struct{}) {
	<-sem
}
