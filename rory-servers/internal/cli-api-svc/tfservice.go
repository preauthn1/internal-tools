package cliapisvr

import (
	"context"
	"fmt"

	"example.invalid/dig"
	"example.invalid/grpc/codes"
	"example.invalid/grpc/status"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
	"example.invalid/ServerTeam/Servers/internal/tfserver/router/portrouter"

	"example.invalid/ServerTeam/Servers/pkg/app"
)

// TFServiceServer TF Service Server
type TFServiceServer struct{}

// QueryConfig 查詢當前設定
func (s *TFServiceServer) QueryConfig(ctx context.Context, in *pb.Empty) (*pb.TFConfig, error) {
	if app.Context == nil {
		return nil, status.Error(codes.Aborted, "app.Context not set")
	}

	res := new(pb.TFConfig)

	err := app.Context.Invoke(func(loader portrouter.PortMapLoader) error {
		pm := loader.LoadPortMap()

		for srcPort, dist := range pm {
			row := new(pb.PortMap)
			row.SrcPort = int32(srcPort)
			row.DistIP = dist.TargetIP
			row.DistPort = int32(dist.TargetPort)
			res.PortMaps = append(res.PortMaps, row)
		}

		return nil
	})
	if err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}

	return res, nil
}

// RefreshConfig 刷新當前設定
func (s *TFServiceServer) RefreshConfig(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	if app.Context == nil {
		return nil, status.Error(codes.Aborted, "app.Context not set")
	}

	err := app.Context.Invoke(func(loader portrouter.PortMapLoader) error {
		loader.Refresh()
		return nil
	})
	if err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}

	return new(pb.Empty), nil
}
