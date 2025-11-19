package grpc_match_making_handler

import (
    "context"

    // "battle_service/domain"
    match_making_srv "battle_service/proto/gen/match_making"
    // "battle_service/internal/service/auth"

    "google.golang.org/grpc"
)

type serverAPI struct {
    match_making_srv.UnimplementedGoMatchMakingServer
}

func Register(gRPCServer *grpc.Server) {
	match_making_srv.RegisterGoMatchMakingServer(gRPCServer, &serverAPI{})
}

func (s *serverAPI) Request(context.Context, *match_making_srv.RequestRequest) (*match_making_srv.RequestResponse, error) {
	return &match_making_srv.RequestResponse{Id: 45452452}, nil
}