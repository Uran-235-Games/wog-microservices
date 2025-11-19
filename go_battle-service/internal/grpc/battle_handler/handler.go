package grpc_battle_handler

import (
    "context"

    // "battle_service/domain"
    battle_srv "battle_service/proto/gen/battle_service"
    // "battle_service/internal/service/auth"

    "google.golang.org/grpc"
)

type serverAPI struct {
    battle_srv.UnimplementedBattleServiceServer
}

func Register(gRPCServer *grpc.Server) {
	battle_srv.RegisterBattleServiceServer(gRPCServer, &serverAPI{})
}

func (s *serverAPI) GetActiveBattles(ctx context.Context, in *battle_srv.GetActiveBattlesRequest) (*battle_srv.GetActiveBattlesResponse, error) {
	return &battle_srv.GetActiveBattlesResponse{}, nil
}