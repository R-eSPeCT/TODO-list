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
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCServer реализует gRPC сервер для аутентификации
type GRPCServer struct {
	pb.UnimplementedAuthServiceServer
	userRepo repository.UserRepository
	jwtKey   []byte
}

// NewGRPCServer создает новый экземпляр gRPC сервера
func NewGRPCServer(userRepo repository.UserRepository, jwtKey []byte) *GRPCServer {
	return &GRPCServer{
		userRepo: userRepo,
		jwtKey:   jwtKey,
	}
}

func (s *GRPCServer) Serve(lis net.Listener) error {
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, s)
	return server.Serve(lis)
}

func (s *GRPCServer) Stop() {
	// Реализация остановки сервера
}

// Register регистрирует нового пользователя
func (s *GRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
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

	// Создаем нового пользователя
	user := &models.User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &pb.RegisterResponse{
		UserId: user.ID,
	}, nil
}

// Login аутентифицирует пользователя и возвращает JWT токен
func (s *GRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Получаем пользователя по email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	// Генерируем JWT токен
	token, err := GenerateToken(user.ID, s.jwtKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &pb.LoginResponse{
		Token: token,
		User: &pb.User{
			Id:    user.ID,
			Email: user.Email,
		},
	}, nil
}

// ValidateToken проверяет JWT токен
func (s *GRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	// Получаем токен из метаданных
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	token := md.Get("authorization")
	if len(token) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	// Проверяем токен
	claims, err := ValidateToken(token[0], s.jwtKey)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// Получаем пользователя
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.ValidateTokenResponse{
		Valid: true,
		User: &pb.User{
			Id:    user.ID,
			Email: user.Email,
		},
	}, nil
}

// UnaryInterceptor для проверки JWT токена
func (s *GRPCServer) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Пропускаем проверку для методов Register и Login
	if info.FullMethod == "/auth.AuthService/Register" || info.FullMethod == "/auth.AuthService/Login" {
		return handler(ctx, req)
	}

	// Проверяем токен для остальных методов
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	token := md.Get("authorization")
	if len(token) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	claims, err := ValidateToken(token[0], s.jwtKey)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// Добавляем ID пользователя в контекст
	ctx = context.WithValue(ctx, "user_id", claims.UserID)
	return handler(ctx, req)
}
