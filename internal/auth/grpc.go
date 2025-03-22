package auth

import (
	"context"
	"net"
	"time"

	"github.com/R-eSPeCT/todo-list/internal/models"
	"github.com/R-eSPeCT/todo-list/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCServer реализует gRPC сервер для аутентификации
type GRPCServer struct {
	userRepo   repository.UserRepository
	jwtManager *JWTManager
	grpcServer *grpc.Server
}

// ServerConfig содержит конфигурацию gRPC сервера
type ServerConfig struct {
	MaxConnectionIdle     time.Duration
	MaxConnectionAge      time.Duration
	MaxConnectionAgeGrace time.Duration
	Time                  time.Duration
	Timeout               time.Duration
	MaxRecvMsgSize        int
}

// NewGRPCServer создает новый экземпляр gRPC сервера
func NewGRPCServer(userRepo repository.UserRepository, jwtKey []byte, cfg ServerConfig) *GRPCServer {
	keepaliveParams := keepalive.ServerParameters{
		MaxConnectionIdle:     cfg.MaxConnectionIdle,
		MaxConnectionAge:      cfg.MaxConnectionAge,
		MaxConnectionAgeGrace: cfg.MaxConnectionAgeGrace,
		Time:                  cfg.Time,
		Timeout:               cfg.Timeout,
	}

	server := grpc.NewServer(
		grpc.KeepaliveParams(keepaliveParams),
		grpc.MaxRecvMsgSize(cfg.MaxRecvMsgSize),
		grpc.UnaryInterceptor(UnaryServerInterceptor()),
	)

	s := &GRPCServer{
		userRepo:   userRepo,
		jwtManager: NewJWTManager(jwtKey),
		grpcServer: server,
	}

	return s
}

// Serve запускает gRPC сервер
func (s *GRPCServer) Serve(lis net.Listener) error {
	return s.grpcServer.Serve(lis)
}

// Stop останавливает gRPC сервер gracefully
func (s *GRPCServer) Stop() {
	s.grpcServer.GracefulStop()
}

// Register регистрирует нового пользователя
func (s *GRPCServer) Register(ctx context.Context, email, password string) (*models.User, error) {
	if email == "" || password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	// Проверяем, существует ли пользователь
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	// Создаем нового пользователя
	userID := uuid.New()
	user := &models.User{
		ID:        userID,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return user, nil
}

// Login аутентифицирует пользователя
func (s *GRPCServer) Login(ctx context.Context, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", status.Error(codes.InvalidArgument, "email and password are required")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", status.Error(codes.NotFound, "user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", status.Error(codes.Unauthenticated, "invalid credentials")
	}

	token, err := s.jwtManager.Generate(user)
	if err != nil {
		return "", status.Error(codes.Internal, "failed to generate token")
	}

	return token, nil
}

// ValidateJWTToken проверяет JWT токен
func (s *GRPCServer) ValidateJWTToken(ctx context.Context, tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	claims, err := s.jwtManager.Validate(tokenString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return claims, nil
}

// UnaryServerInterceptor создает перехватчик для проверки JWT токена
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Пропускаем проверку для методов регистрации и логина
		if info.FullMethod == "/auth.AuthService/Register" || info.FullMethod == "/auth.AuthService/Login" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		// Добавляем информацию о пользователе в контекст
		userCtx := context.WithValue(ctx, "user_id", authHeader[0])
		return handler(userCtx, req)
	}
}
