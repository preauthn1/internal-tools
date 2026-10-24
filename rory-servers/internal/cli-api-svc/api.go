package cliapisvr

import (
	"context"
	"fmt"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"

	"example.invalid/ServerTeam/Servers/pkg/httputil"
)

// APIServiceServer API Service Server
type APIServiceServer struct{}

// Status get api service info
func (s *APIServiceServer) Status(ctx context.Context, arg *pb.Empty) (*pb.APIStatus, error) {
	if httputil.APIUtilCmdP == nil {
		return nil, fmt.Errorf("APIUtilCmdP not init")
	}

	status := new(pb.APIStatus)
	status.IPOS = httputil.APIUtilCmdP.IOPS()
	status.Avg = httputil.APIUtilCmdP.Avg().String()
	status.Size = int32(httputil.APIUtilCmdP.Size())
	min, max := httputil.APIUtilCmdP.MinMax()
	status.Min = min.String()
	status.Max = max.String()

	return status, nil
}

// Output history log
func (s *APIServiceServer) Output(ctx context.Context, arg *pb.Empty) (*pb.Empty, error) {
	if httputil.APIUtilCmdP == nil {
		return nil, fmt.Errorf("APIUtilCmdP not init")
	}
	httputil.APIUtilCmdP.Output()

	return new(pb.Empty), nil
}
