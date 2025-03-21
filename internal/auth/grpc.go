package auth

import (
	"context"
	"net"
	"time"

	"github.com/R-eSPeCT/todo-list/internal/models"
	"github.com/R-eSPeCT/todo-list/internal/repository"
	pb "github.com/R-eSPeCT/todo-list/pkg/proto/auth"

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
	pb.UnimplementedAuthServiceServer
	userRepo   repository.UserRepository
	jwtKey     []byte
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
		jwtKey:     jwtKey,
		grpcServer: server,
	}

	pb.RegisterAuthServiceServer(server, s)
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
func (s *GRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	// Проверяем, существует ли пользователь
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &pb.RegisterResponse{
		Id:    user.ID.String(),
		Email: user.Email,
	}, nil
}

// Login аутентифицирует пользователя и возвращает JWT токен
func (s *GRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	token, err := GenerateToken(user.ID.String(), user.Email, s.jwtKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

// ValidateToken проверяет JWT токен
func (s *GRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	claims, err := ValidateToken(req.Token, s.jwtKey)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return &pb.ValidateTokenResponse{
		UserId: claims.UserID,
		Email:  claims.Email,
	}, nil
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
