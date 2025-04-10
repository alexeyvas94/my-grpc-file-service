# 📦 my-grpc-file-service

gRPC-сервис на Go для загрузки, скачивания и просмотра изображений.

## 🚀 Возможности

- 📥 Загрузка файлов (stream upload)
- 📤 Скачивание файлов (stream download)
- 📃 Получение списка загруженных файлов
- 🔐 Ограничение конкурентных подключений:
  - Upload/Download — максимум 10
  - List — максимум 100

## 🛠️ Технологии

- Golang (gRPC, ProtoBuf)
- Docker / Docker Compose
- `.env` конфигурация

## 🗂 Структура

```
my-grpc-file-service/
├── cmd/server/main.go           # Точка входа
├── internal/                    # Логика сервиса
│   ├── delivery/                # gRPC обработчики
│   ├── domain/                  # Модели
│   ├── repository/              # Работа с файлами
│   ├── infrastructure/          # Метаданные
│   └── service/                 # Бизнес-логика
├── proto/file_service.proto     # gRPC контракты
├── api/                         # Сгенерированные gRPC-файлы
├── scripts/docker-compose.yml
├── .env                         # Настройки
├── Makefile                     # Утилиты
├── README.md
├── go.mod / go.sum
```

## ⚙️ Конфигурация

`.env` файл:

```env
GRPC_PORT=50051
UPLOAD_DIR=uploaded
```

## 🔧 Сборка и запуск

### Установка зависимостей:

```bash
go mod tidy
```

### Генерация gRPC кода:

```bash
make proto
```

Убедись, что установлен `protoc` и плагины:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Запуск локально:

```bash
go run cmd/server/main.go
```

### Запуск в Docker:

```bash
docker-compose -f scripts/docker-compose.yml up --build
```

## 📡 gRPC API

Файл контракта: `proto/file_service.proto`

```proto
rpc Upload(stream FileUploadRequest) returns (FileUploadResponse);
rpc Download(FileDownloadRequest) returns (stream FileDownloadResponse);
rpc ListFiles(ListFilesRequest) returns (ListFilesResponse);
```

## 🧪 Тестирование

Протестируй через:

- [grpcurl](https://github.com/fullstorydev/grpcurl)
- [BloomRPC](https://github.com/bloomrpc/bloomrpc)
- CLI клиент на Go (могу скинуть)

---