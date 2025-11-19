package grpc_auth

import (
    "context"

    "wog-server/domain"
    "wog-server/proto/gen"
    "wog-server/internal/service/auth"
    "wog-server/internal/lib/jwt"

    "google.golang.org/grpc"
)

type serverAPI struct {
    user_service.UnimplementedUserServiceServer
	auth *auth.AuthService
    jwt *jwt.JWTSrv
}

func Register(gRPCServer *grpc.Server, auth *auth.AuthService, jwt *jwt.JWTSrv) {
	user_service.RegisterUserServiceServer(gRPCServer, &serverAPI{auth: auth, jwt: jwt})
}

func (s *serverAPI) SignUp(ctx context.Context, in *user_service.SignUpRequest) (*user_service.SignUpResponse, error) {
    // TODO: validate input
    authData := domain.User{Name: in.Name, Email: in.Email, Password: in.Password}
    uid, err := s.auth.Register(authData)
    if err != nil {
        panic(err)
    }

    return &user_service.SignUpResponse{Id: uid, Name: in.GetName(), Email: in.GetEmail()}, nil
}