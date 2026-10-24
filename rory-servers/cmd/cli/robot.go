package main

import (
	"fmt"
	"strconv"

	"example.invalid/spf13/cobra"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
)

func robotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "robot",
		Short:                 "機器人管理",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			summary, err := CLIClient.GetRobotSummary(cmd.Context(), new(pb.Empty))
			if err != nil {
				return err
			}

			cmd.Printf("=== 機器人資訊 ===\n")
			cmd.Printf("總數: %d\n", summary.RobotCount)
			cmd.Printf("遊戲中: %d\n", summary.GamingCount)
			cmd.Printf("空閒中: %d\n", summary.FreeCount)

			return nil
		},
	}

	cmd.AddCommand(robotInfoCmd())
	cmd.AddCommand(robotListCmd())

	return cmd
}

func robotInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "find <pid>",
		Short:                 "查詢機器人",
		Args:                  cobra.ExactValidArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			uid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("機器人ID輸入錯誤, %v", err)
			}

			p, err := CLIClient.QueryRobot(cmd.Context(), &pb.QueryCondition{PlayerID: int32(uid)})
			if err != nil {
				return err
			}

			printPlayer(cmd.OutOrStdout(), p.Player)

			return nil
		},
	}

	return cmd
}

func robotListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出機器人",
		RunE: func(cmd *cobra.Command, args []string) error {

			all, _ := cmd.PersistentFlags().GetBool("all")
			free, _ := cmd.PersistentFlags().GetBool("free")
			gaming, _ := cmd.PersistentFlags().GetBool("gaming")

			cond := new(pb.QueryCondition)
			cond.IsFree = free
			cond.IsGaming = gaming

			if all {
				cond.IsFree = false
				cond.IsGaming = false
			}

			list, err := CLIClient.ListRobot(cmd.Context(), cond)
			if err != nil {
				return err
			}

			cmd.Printf("總人數: %d\n", len(list.Robots))

			for _, p := range list.Robots {
				printPlayerLite(cmd.OutOrStdout(), p.Player, "  ")
			}

			return nil
		},
	}

	cmd.PersistentFlags().BoolP("gaming", "g", false, "只查詢遊戲中機器人")
	cmd.PersistentFlags().BoolP("free", "f", false, "只查詢空閒中機器人")
	cmd.PersistentFlags().BoolP("all", "a", false, "查詢所有機器人")

	return cmd
}
