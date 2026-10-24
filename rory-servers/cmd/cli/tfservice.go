package main

import (
	"fmt"
	"sort"
	"time"

	"example.invalid/spf13/cobra"
	"example.invalid/grpc"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
)

// TFServiceClient TF Service Client
var TFServiceClient pb.TFServiceClient

func tfConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "列出設定資訊",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			addr, _ := cmd.Flags().GetString("addr")
			port, _ := cmd.Flags().GetInt("port")
			address := fmt.Sprintf("%s:%d", addr, port)

			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
			if err != nil {
				return fmt.Errorf("GRPC Dial Error to %s, %v", address, err)
			}
			TFServiceClient = pb.NewTFServiceClient(conn)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			tfcfg, err := TFServiceClient.QueryConfig(cmd.Context(), new(pb.Empty))
			if err != nil {
				return err
			}

			sort.Slice(tfcfg.PortMaps, func(i, j int) bool {
				return tfcfg.PortMaps[i].SrcPort < tfcfg.PortMaps[j].SrcPort
			})

			cmd.Printf("遠端轉接設定: (%d)\n", len(tfcfg.PortMaps))
			for _, row := range tfcfg.PortMaps {
				cmd.Printf(" - %9d -> %s:%d\n", row.SrcPort, row.DistIP, row.DistPort)
			}

			return nil
		},
	}

	cmd.AddCommand(tfConfigRefreshCmd())

	return cmd
}

func tfConfigRefreshCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "更新資訊",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := TFServiceClient.RefreshConfig(cmd.Context(), new(pb.Empty))
			if err != nil {
				return err
			}
			cmd.Println("刷新遠程端口映射完成...")
			return nil
		},
	}
	return cmd
}
