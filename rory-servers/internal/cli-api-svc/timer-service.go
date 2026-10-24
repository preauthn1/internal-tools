package cliapisvr

import (
	"context"
	"fmt"

	"example.invalid/dig"
	"example.invalid/grpc/codes"
	"example.invalid/grpc/status"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"

	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/event"
)

// TimerManagerServiceServer 時間管理器
type TimerManagerServiceServer struct{}

// Status 取得時間管理器狀態
func (s *TimerManagerServiceServer) Status(ctx context.Context, in *pb.Empty) (*pb.TimerManagerStatus, error) {
	if app.Context == nil {
		return nil, status.Error(codes.Aborted, "app.Context not set")
	}

	res := new(pb.TimerManagerStatus)

	err := app.Context.Invoke(func(tfm *event.TimerFuncManager) error {
		res.Length = int32(tfm.Len())
		res.Keys = tfm.Keys()
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

// Remove 移除指定計時器
func (s *TimerManagerServiceServer) Remove(ctx context.Context, in *pb.TimerQueryCondition) (*pb.Empty, error) {
	if app.Context == nil {
		return nil, status.Error(codes.Aborted, "app.Context not set")
	}

	res := new(pb.Empty)

	err := app.Context.Invoke(func(tfm *event.TimerFuncManager) error {
		f := tfm.GetTimerFunc(in.Key)
		if f == nil {
			return status.Error(codes.NotFound, fmt.Sprintf("TimerFunc %s not found", in.Key))
		}
		tfm.RemoveTimer(in.Key)
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
