package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/R-eSPeCT/todo-list/internal/auth"
	"github.com/R-eSPeCT/todo-list/internal/config"
	"github.com/R-eSPeCT/todo-list/internal/repository"
	_ "github.com/lib/pq"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключаемся к базе данных
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Проверяем соединение с базой данных
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Создаем репозитории
	userRepo := repository.NewUserRepository(db)

	// Создаем конфигурацию gRPC сервера
	grpcConfig := auth.ServerConfig{
		MaxConnectionIdle:     15 * time.Minute,
		MaxConnectionAge:      30 * time.Minute,
		MaxConnectionAgeGrace: 5 * time.Second,
		Time:                  5 * time.Second,
		Timeout:               1 * time.Second,
		MaxRecvMsgSize:        4 * 1024 * 1024, // 4MB
	}

	// Создаем gRPC сервер
	grpcServer := auth.NewGRPCServer(userRepo, []byte(cfg.GRPC.JWTSecretKey), grpcConfig)

	// Создаем TCP listener для gRPC
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Создаем канал для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запускаем gRPC сервер в горутине
	go func() {
		log.Printf("Starting gRPC server on port %d", cfg.GRPC.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Failed to serve: %v", err)
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	<-stop
	log.Println("Shutting down server...")

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Останавливаем сервер
	grpcServer.Stop()

	// Ожидаем завершения всех горутин
	<-ctx.Done()
	log.Println("Server stopped")
}
