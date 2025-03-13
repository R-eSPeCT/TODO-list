package auth

import (
	"testing"

	"github.com/your-project/pb"
)

func TestGRPCServer_Register(t *testing.T) {
	// Создаем тестовый сервер
	grpcTestServer := NewTestServer(t)
	defer grpcTestServer.Cleanup()

	// Создаем тестовый клиент
	conn := grpcTestServer.ClientConn()
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	// ... existing code ...
}
