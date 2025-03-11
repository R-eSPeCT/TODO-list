#!/bin/bash

# Создаем директории для сгенерированных файлов
mkdir -p api/proto/todo
mkdir -p api/proto/user

# Генерируем код из proto файлов
protoc --go_out=. \
       --go_opt=paths=source_relative \
       --go-grpc_out=. \
       --go-grpc_opt=paths=source_relative \
       api/proto/todo.proto \
       api/proto/user.proto 