package auth

import (
	"context"
	"fmt"
	"net"

	"github.com/yourusername/todo-list/internal/models"
	"github.com/yourusername/todo-list/internal/repository"
	pb "github.com/yourusername/todo-list/pkg/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedAuthServiceServer
	userRepo repository.UserRepository
	jwtKey   []byte
}

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

func (s *GRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	// Проверяем, существует ли пользователь
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}

	// Создаем нового пользователя
	user := &models.User{
		Email:    req.Email,
		Password: req.Password, // В реальном приложении пароль должен быть захеширован
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create user: %v", err))
	}

	// Генерируем токен
	token, err := GenerateToken(user.ID, s.jwtKey)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to generate token: %v", err))
	}

	return &pb.RegisterResponse{
		Token: token,
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

func (s *GRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	// Получаем пользователя
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Проверяем пароль (в реальном приложении нужно использовать безопасное сравнение)
	if user.Password != req.Password {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// Генерируем токен
	token, err := GenerateToken(user.ID, s.jwtKey)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to generate token: %v", err))
	}

	return &pb.LoginResponse{
		Token: token,
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

func (s *GRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	// Валидируем токен
	claims, err := ValidateToken(req.Token, s.jwtKey)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	// Получаем пользователя
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid: true,
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}
