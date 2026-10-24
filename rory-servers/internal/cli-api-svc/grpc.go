package cliapisvr

import (
	"example.invalid/grpc"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
)

// NewGRPC 建立grpc服務
func NewGRPC() *grpc.Server {
	srv := grpc.NewServer()
	pb.RegisterCLIServer(srv, new(CLIService))
	pb.RegisterLogServer(srv, new(LogService))
	pb.RegisterAPIServiceServer(srv, new(APIServiceServer))
	pb.RegisterSlotServiceServer(srv, new(SlotServiceServer))
	pb.RegisterMessageServiceServer(srv, new(MessageServiceServer))
	pb.RegisterPachiSlotServiceServer(srv, new(PachiSlotServiceServer))
	pb.RegisterTFServiceServer(srv, new(TFServiceServer))
	pb.RegisterTimerManagerServiceServer(srv, new(TimerManagerServiceServer))
	return srv
}
