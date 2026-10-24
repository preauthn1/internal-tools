package main

import (
	"example.invalid/spf13/cobra"
	"example.invalid/grpc/status"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
)

func summaryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "summary",
		Short: "列出總資訊",
		RunE: func(cmd *cobra.Command, args []string) error {

			// 版本資訊
			{
				v, err := CLIClient.ServerVer(cmd.Context(), new(pb.Empty))
				cmd.Printf("=== 版本資訊 ===\n")
				if err != nil {
					if e, ok := status.FromError(err); ok {
						cmd.Println(e.Message())
					} else {
						cmd.Println(err)
					}
				} else {
					cmd.Printf("當前版本：%s GitHash：%s\n", v.Version, v.Hash)
					cmd.Printf("建置時間：%s\n", v.BuildTime)
					cmd.Printf("建置編號：%s\n", v.BuildName)
				}
			}

			// 房間資訊
			{
				room, err := CLIClient.RoomInfo(cmd.Context(), new(pb.Empty))
				cmd.Printf("=== 房間資訊 ===\n")
				if err != nil {
					if e, ok := status.FromError(err); ok {
						cmd.Println(e.Message())
					} else {
						cmd.Println(err)
					}
				} else {
					printRoom(cmd.OutOrStdout(), room, "")
				}
			}

			// 玩家資訊
			{
				summary, err := CLIClient.GetPlayerSummary(cmd.Context(), new(pb.Empty))
				cmd.Printf("=== 玩家資訊 ===\n")
				if err != nil {
					if e, ok := status.FromError(err); ok {
						cmd.Println(e.Message())
					} else {
						cmd.Println(err)
					}
				} else {
					cmd.Printf("總人數: %d\n", summary.PlayerCount)
					cmd.Printf("玩家: %d\n", summary.UserCount)
					cmd.Printf("遊客: %d\n", summary.GuestCount)
				}
			}

			// 機器人資訊
			{
				summary, err := CLIClient.GetRobotSummary(cmd.Context(), new(pb.Empty))
				cmd.Printf("=== 機器人資訊 ===\n")
				if err != nil {
					if e, ok := status.FromError(err); ok {
						cmd.Println(e.Message())
					} else {
						cmd.Println(err)
					}
				} else {
					cmd.Printf("總數: %d\n", summary.RobotCount)
					cmd.Printf("遊戲中: %d\n", summary.GamingCount)
					cmd.Printf("空閒中: %d\n", summary.FreeCount)
				}
			}

			return nil
		},
	}
	return cmd
}
