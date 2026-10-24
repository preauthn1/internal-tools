package main

import (
	"io"

	"example.invalid/spf13/cobra"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
)

// MessageClient 封包工具
var MessageClient pb.MessageServiceClient

func messageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "message",
		Short:                 "封包工具",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(msgtraceCmd())

	return cmd
}

func msgtraceCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "trace",
		Short: "追蹤封包",
		RunE: func(cmd *cobra.Command, args []string) error {

			table, _ := cmd.PersistentFlags().GetInt32("table")
			player, _ := cmd.PersistentFlags().GetInt32("player")

			q := new(pb.QueryCondition)
			q.Table = table
			q.PlayerID = player

			stream, err := MessageClient.Trace(cmd.Context(), q)
			if err != nil {
				return err
			}

			if player > 0 {
				cmd.Printf("開始追蹤遊戲封包, 桌號: %d, 玩家: %d\n", table, player)
			} else {
				cmd.Printf("開始追蹤遊戲封包, 桌號: %d\n", table)
			}

			for {
				select {
				default:
					resp, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							return nil
						}
						return err
					}

					var vec string
					if resp.Src == pb.Endpoint_Server && resp.Dest == pb.Endpoint_Client {
						vec = "<--"
					} else if resp.Src == pb.Endpoint_Client && resp.Dest == pb.Endpoint_Server {
						vec = "-->"
					}

					cmd.Printf("玩家[%d] %s [%d:%d] %s\n", resp.PlayerID, vec, resp.BClassID, resp.SClassID, resp.Data)
				case <-cmd.Context().Done():
					cmd.Println("Command Close")
				case <-stream.Context().Done():
					cmd.Println("Server Close Connect")
					return nil
				}
			}
		},
	}

	cmd.PersistentFlags().Int32P("table", "t", 0, "指定桌號")
	cmd.PersistentFlags().Int32P("player", "p", 0, "指定玩家")

	return cmd
}
