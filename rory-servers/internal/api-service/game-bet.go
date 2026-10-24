package api

import (
	"context"
	"fmt"

	pb "example.invalid/ServerTeam/Servers/internal/api-service/apisvc"
	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	"example.invalid/grpc"
	"example.invalid/grpc/codes"
	"example.invalid/grpc/status"
)

// GameBetServiceServer server of game bet service
type GameBetServiceServer struct{}

// PlaceBet 下注扣錢
func (s *GameBetServiceServer) PlaceBet(ctx context.Context, in *pb.PlaceBetReq) (*pb.PlaceBetRes, error) {

	ok, afterAmount, transactionNo := service.GamePlayService.PostPlaceBet(in.AccessToken, in.ServerID, in.TradeID, in.BetAmount, in.IPAddress, in.Memo)

	res := new(pb.PlaceBetRes)

	if !ok {
		return nil, status.Error(codes.Internal, fmt.Sprintf("PlaceBet API Error"))
	}

	res.AfterAmount = afterAmount
	res.TransactionNo = transactionNo

	return res, nil
}

// MockGameBetServiceServer server of game bet service
type MockGameBetServiceServer struct{}

// PlaceBet 下注扣錢
func (s *MockGameBetServiceServer) PlaceBet(ctx context.Context, in *pb.PlaceBetReq, options ...grpc.CallOption) (*pb.PlaceBetRes, error) {
	res := new(pb.PlaceBetRes)
	res.AfterAmount = 0
	res.TransactionNo = "mock"
	return res, nil
}
