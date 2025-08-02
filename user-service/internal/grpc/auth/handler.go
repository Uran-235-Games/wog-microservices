package auth

import (
    "context"

    pb_auth "github.com/Uran-235-Games/wog-grpc-lib/gen/go/auth"
	auth_service "wog-server/internal/service/auth"

    "google.golang.org/grpc"
)

type serverAPI struct {
    pb_auth.UnimplementedAuthServiceServer
	authService *auth_service.AuthService
}

func Register(gRPCServer *grpc.Server, auth *auth_service.AuthService) {
	pb_auth.RegisterAuthServiceServer(gRPCServer, &serverAPI{authService: auth})
}

func (s *serverAPI) SignUp(ctx context.Context, in *pb_auth.SignUpRequest) (*pb_auth.SignUpResponse, error) {
    return &pb_auth.SignUpResponse{Id: "pon", Name: "ivan", Email: "cum@.exe"}, nil
}