package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	
)

// AuthInterceptor представляет интерцептор для аутентификации
type AuthInterceptor struct {
	jwtManager *JWTManager
	// Методы, не требующие аутентификации
	publicMethods map[string]bool
}

// NewAuthInterceptor создает новый интерцептор аутентификации
func NewAuthInterceptor(jwtManager *JWTManager, publicMethods []string) *AuthInterceptor {
	methods := make(map[string]bool)
	for _, method := range publicMethods {
		methods[method] = true
	}

	return &AuthInterceptor{
		jwtManager:    jwtManager,
		publicMethods: methods,
	}
}

// Unary возвращает серверный унарный интерцептор для аутентификации
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if i.publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		userID, err := i.authorize(ctx)
		if err != nil {
			return nil, err
		}

		// Добавляем ID пользователя в контекст
		newCtx := context.WithValue(ctx, "user_id", userID)
		return handler(newCtx, req)
	}
}

// Stream возвращает серверный стрим интерцептор для аутентификации
func (i *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if i.publicMethods[info.FullMethod] {
			return handler(srv, stream)
		}

		userID, err := i.authorize(stream.Context())
		if err != nil {
			return err
		}

		// Оборачиваем стрим для добавления ID пользователя в контекст
		wrapped := newWrappedStream(stream, userID)
		return handler(srv, wrapped)
	}
}

// authorize проверяет токен и возвращает ID пользователя
func (i *AuthInterceptor) authorize(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	if !strings.HasPrefix(accessToken, "Bearer ") {
		return "", status.Errorf(codes.Unauthenticated, "invalid authorization format")
	}

	accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	claims, err := i.jwtManager.Verify(accessToken)
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	return claims.UserID, nil
}

// wrappedStream оборачивает grpc.ServerStream для добавления информации о пользователе
type wrappedStream struct {
	grpc.ServerStream
	userID string
}

func newWrappedStream(stream grpc.ServerStream, userID string) grpc.ServerStream {
	return &wrappedStream{
		ServerStream: stream,
		userID:       userID,
	}
}

func (w *wrappedStream) Context() context.Context {
	ctx := w.ServerStream.Context()
	return context.WithValue(ctx, "user_id", w.userID)
}
