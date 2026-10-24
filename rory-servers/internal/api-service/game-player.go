package api

import (
	"context"
	"fmt"

	pb "example.invalid/ServerTeam/Servers/internal/api-service/apisvc"
	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	"example.invalid/grpc/codes"
	"example.invalid/grpc/status"
)

// GamePlayerServiceServer server of player service
type GamePlayerServiceServer struct{}

// Wallet 查詢玩家錢包
func (s *GamePlayerServiceServer) Wallet(ctx context.Context, in *pb.WalletReq) (*pb.WalletRes, error) {
	ok, amount, merchantWalletMode := service.GamePlayService.GetWallet(in.AccessToken)

	res := new(pb.WalletRes)

	if !ok {
		return nil, status.Error(codes.Internal, fmt.Sprintf("PlaceBet API Error"))
	}

	res.Amount = amount
	res.MerchantWalletMode = merchantWalletMode

	return res, nil
}
