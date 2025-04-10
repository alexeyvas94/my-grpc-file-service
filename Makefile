PROTOC = protoc

# Цель для генерации файлов
proto: check_dirs proto_generate

# Проверка и создание директорий, если их нет
check_dirs:
	@mkdir -p pkg/api

# Генерация файлов
proto_generate:
	$(PROTOC) --proto_path=proto \
		--go_out=paths=source_relative:pkg/api \
		--go-grpc_out=paths=source_relative:pkg/api \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		proto/file_service.proto
