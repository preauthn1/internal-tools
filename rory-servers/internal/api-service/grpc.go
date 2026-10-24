package api

import (
	"example.invalid/grpc"

	pb "example.invalid/ServerTeam/Servers/internal/api-service/apisvc"
)

// NewGRPC 建立grpc服務
func NewGRPC() *grpc.Server {
	srv := grpc.NewServer()
	pb.RegisterGameBetServiceServer(srv, new(GameBetServiceServer))
	pb.RegisterGameResultServiceServer(srv, new(GameResultServiceServer))
	pb.RegisterGamePlayerServiceServer(srv, new(GamePlayerServiceServer))
	return srv
}
