package main

import (
	"fmt"
	"time"

	"example.invalid/spf13/cobra"
	"example.invalid/grpc"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
)

// TimerManagerServiceClient TimerManager Service Client
var TimerManagerServiceClient pb.TimerManagerServiceClient

func timerManagerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "timer",
		Short: "時間管理器",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			addr, _ := cmd.Flags().GetString("addr")
			port, _ := cmd.Flags().GetInt("port")
			address := fmt.Sprintf("%s:%d", addr, port)

			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
			if err != nil {
				return fmt.Errorf("GRPC Dial Error to %s, %v", address, err)
			}
			TimerManagerServiceClient = pb.NewTimerManagerServiceClient(conn)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			status, err := TimerManagerServiceClient.Status(cmd.Context(), new(pb.Empty))
			if err != nil {
				return err
			}

			cmd.Printf("=== TimerFuncManager ===\n")
			cmd.Printf("Length: %d\n", status.Length)
			if len(status.Keys) > 0 {
				cmd.Printf("Keys:\n")
				for _, key := range status.Keys {
					cmd.Printf(" - %s\n", key)
				}
			}

			return nil
		},
	}

	cmd.AddCommand(timerManagerRemoveCmd())

	return cmd
}

func timerManagerRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "rm <name>",
		Short:                 "移除指定 TimerFunc",
		Args:                  cobra.ExactValidArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			cond := new(pb.TimerQueryCondition)
			cond.Key = args[0]

			_, err := TimerManagerServiceClient.Remove(cmd.Context(), cond)
			if err != nil {
				return err
			}

			cmd.Printf("TimerFunc 『%s』 移除成功\n", cond.Key)

			return nil
		},
	}
	return cmd
}
