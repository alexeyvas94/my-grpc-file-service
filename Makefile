proto:
	protoc -I proto \
      --go_out=api --go_opt=paths=source_relative \
      --go-grpc_out=api --go-grpc_opt=paths=source_relative \
      proto/file_service.proto
