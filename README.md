# TODO List API

REST API и gRPC сервис для управления списком задач с поддержкой пользователей и аутентификации.

## Требования

- Go 1.21 или выше
- PostgreSQL 12 или выше
- Redis 6 или выше
- Protocol Buffers compiler (protoc)

## Установка protoc

### Windows
```bash
# Скачайте последнюю версию с https://github.com/protocolbuffers/protobuf/releases
# Распакуйте архив и добавьте bin директорию в PATH

# Установите Go плагины для protoc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Linux/macOS
```bash
# Linux
apt-get install -y protobuf-compiler

# macOS
brew install protobuf

# Установите Go плагины для protoc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/yourusername/TODO-list.git
cd TODO-list
```

2. Установите зависимости:
```bash
go mod download
```

3. Сгенерируйте код из proto файлов:
```bash
# Сделайте скрипт исполняемым (для Linux/macOS)
chmod +x scripts/protoc.sh

# Запустите генерацию
./scripts/protoc.sh
```

4. Создайте базу данных PostgreSQL:
```sql
CREATE DATABASE todo_list;
```

5. Настройте переменные окружения:
- Скопируйте файл `.env.example` в `.env`
- Отредактируйте параметры подключения к базе данных и Redis в `.env`

## Запуск

```bash
go run main.go
```

Сервер запустится на:
- REST API: порт 3000 (или указанный в PORT)
- gRPC: порт 50051 (или указанный в GRPC_PORT)

## API Endpoints

### Пользователи

- `POST /api/users/register` - Регистрация нового пользователя
- `POST /api/users/login` - Вход пользователя

### Задачи

- `GET /api/todos` - Получение списка всех задач пользователя
- `POST /api/todos` - Создание новой задачи
- `GET /api/todos/grouped` - Получение сгруппированных задач
- `GET /api/todos/:id` - Получение задачи по ID
- `PUT /api/todos/:id` - Обновление задачи
- `DELETE /api/todos/:id` - Удаление задачи

## Структура проекта

```
.
├── cmd/
├── internal/
│   ├── handlers/     # HTTP обработчики
│   ├── middleware/   # Промежуточное ПО
│   ├── models/       # Модели данных
│   ├── repository/   # Слой доступа к данным
│   └── services/     # Бизнес-логика
├── pkg/             # Общие пакеты
├── .env             # Конфигурация
├── .gitignore
├── go.mod
├── go.sum
├── main.go         # Точка входа
└── README.md
```

## Описание

TODO-list — это веб-приложение для управления задачами, разработанное с использованием Go и фреймворка Fiber. Оно позволяет пользователям создавать, обновлять, удалять и просматривать задачи.

## Технологии

- Go
- Fiber
- PostgreSQL 
 ## Переменные окружения

Перед запуском приложения убедитесь, что у вас настроены следующие переменные окружения:

- DATABASE_URL: URL для подключения к вашей базе данных. Пример: postgres://user:password@localhost:5432/todo_db

Вы можете создать файл .env в корне проекта и добавить туда переменные окружения:


DATABASE_URL=postgres://user:password@localhost:5432/todo_db

## Запуск

1. Запустите сервер:


Bash


go run cmd/server/main.go


2. Сервер будет доступен по адресу http://localhost:3000.

## Использование

- GET /tasks: Получить список всех задач.
- POST /tasks: Создать новую задачу.
- PUT /tasks/:id: Обновить задачу по ID.
- DELETE /tasks/:id: Удалить задачу по ID.

## Автор

- Нуриев Рустем Андреевич


## Лицензия

Этот проект лицензирован под MIT License. Подробности см. в файле LICENSE.


