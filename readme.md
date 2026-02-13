# Currency API

Микросервис для получения и управления курсами валют с поддержкой gRPC API.

## Возможности
- Получение актуальных курсов валют от ЦБ РФ
- Автоматическая синхронизация курсов (раз в 24 часа)
- gRPC API для интеграции с другими сервисами
- Кэширование данных для повышения производительности
- Уведомления о резких изменениях курсов
- Поддержка in-memory и PostgreSQL хранилищ

## Требования
- Go 1.25+
- Docker и Docker Compose
- PostgreSQL (опционально, для постоянного хранения)

## Установка и запуск
### Способ 1: Локальная разработка

```bash
# Клонирование репозитория
git clone https://github.com/yourusername/currency-api.git
cd currency-api

# Установка зависимостей
go mod download

# Запуск тестов
go test ./...

# Запуск приложения
go run cmd/currency/main.go
```

### Способ 2: Docker
```bash
# Сборка и запуск с Docker Compose
docker-compose up -d

# Проверка статуса
docker-compose ps

# Просмотр логов
docker-compose logs -f currency-api
```

### Способ 3: Kubernetes (опционально)
```bash
kubectl apply -f k8s/
```

## Конфигурация
Создайте файл .env в корне проекта:

env
# Основные настройки
GRPC_PORT=50051
LOG_MODE=dev
USE_POSTGRES=false

# Настройки ЦБ РФ
CBR_URL=https://www.cbr.ru/scripts/XML_daily.asp
CBR_REQUIRED_CODES=USD,EUR,AED,GBP,CNY,JPY
CBR_TIMEOUT=30s

# Настройки PostgreSQL
POSTGRES_DSN=postgres://currency:currency@postgres:5432/currency?sslmode=disable

# Настройки кэширования
CACHE_DURATION=24h
RATE_SPIKE_THRESHOLD=0.1  # 10%


