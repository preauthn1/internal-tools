package main

import (
	"fmt"
	"time"

	"example.invalid/spf13/cobra"
	"example.invalid/grpc"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
)

// APIServiceClient API Server 工具
var APIServiceClient pb.APIServiceClient

func apiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "API Server 工具",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			addr, _ := cmd.Flags().GetString("addr")
			port, _ := cmd.Flags().GetInt("port-api")
			address := fmt.Sprintf("%s:%d", addr, port)

			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
			if err != nil {
				return fmt.Errorf("GRPC Dial Error to %s, %v", address, err)
			}
			APIServiceClient = pb.NewAPIServiceClient(conn)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(apiStatusCmd())
	cmd.AddCommand(apiOutputCmd())

	return cmd
}

func apiStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "查詢狀態",
		RunE: func(cmd *cobra.Command, args []string) error {

			status, err := APIServiceClient.Status(cmd.Context(), new(pb.Empty))
			if err != nil {
				return err
			}

			cmd.Printf("=== APIQueryHistroy ===\n")
			cmd.Printf("IOPS: %s\n", status.IPOS)
			cmd.Printf("Time Average: %s\n", status.Avg)
			cmd.Printf("Queue Size: %d\n", status.Size)
			cmd.Printf("Execute Min Time: %s\n", status.Min)
			cmd.Printf("Execute Max Time: %s\n", status.Max)

			return nil
		},
	}
	return cmd
}

func apiOutputCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "output",
		Short: "輸出最近查詢結果",
		RunE: func(cmd *cobra.Command, args []string) error {

			_, err := APIServiceClient.Output(cmd.Context(), new(pb.Empty))
			if err != nil {
				return err
			}

			cmd.Printf("Done\n")

			return nil
		},
	}
	return cmd
}
