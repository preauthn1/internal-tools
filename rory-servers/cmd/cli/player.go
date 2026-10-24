package main

import (
	"fmt"
	"io"
	"strconv"

	"example.invalid/spf13/cobra"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
)

func playerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "player",
		Short:                 "玩家管理",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			summary, err := CLIClient.GetPlayerSummary(cmd.Context(), new(pb.Empty))
			if err != nil {
				return err
			}

			cmd.Printf("=== 玩家資訊 ===\n")
			cmd.Printf("總人數: %d\n", summary.PlayerCount)
			cmd.Printf("玩家: %d\n", summary.UserCount)
			cmd.Printf("遊客: %d\n", summary.GuestCount)

			return nil
		},
	}

	cmd.AddCommand(playerInfoCmd())
	cmd.AddCommand(playerKickCmd())
	cmd.AddCommand(playerListCmd())

	return cmd
}

func playerInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "find <pid>",
		Short:                 "查詢玩家",
		Args:                  cobra.ExactValidArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			uid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("玩家ID輸入錯誤, %v", err)
			}

			p, err := CLIClient.QueryPlayer(cmd.Context(), &pb.QueryCondition{PlayerID: int32(uid)})
			if err != nil {
				return err
			}

			printPlayer(cmd.OutOrStdout(), p)

			return nil
		},
	}

	return cmd
}

func playerKickCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "kick <pid>",
		Short:                 "剔除玩家",
		Args:                  cobra.ExactValidArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			uid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("玩家ID輸入錯誤, %v", err)
			}

			if _, err := CLIClient.KickPlayer(cmd.Context(), &pb.QueryCondition{PlayerID: int32(uid)}); err != nil {
				return err
			}

			cmd.Printf("玩家[%d]剔除成功\n", uid)

			return nil
		},
	}

	return cmd
}

func playerListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出玩家",
		RunE: func(cmd *cobra.Command, args []string) error {

			all, _ := cmd.PersistentFlags().GetBool("all")
			user, _ := cmd.PersistentFlags().GetBool("user")
			guest, _ := cmd.PersistentFlags().GetBool("guest")

			cond := new(pb.QueryCondition)
			cond.IsLogined = user
			cond.IsGuest = guest

			if all {
				cond.IsLogined = true
				cond.IsGuest = true
			}

			list, err := CLIClient.ListPlayer(cmd.Context(), cond)
			if err != nil {
				return err
			}

			cmd.Printf("總人數: %d\n", len(list.Players))

			for _, p := range list.Players {
				printPlayerLite(cmd.OutOrStdout(), p, "  ")
			}

			return nil
		},
	}

	cmd.PersistentFlags().BoolP("guest", "g", false, "只查詢訪客")
	cmd.PersistentFlags().BoolP("user", "u", false, "只查詢玩家")
	cmd.PersistentFlags().BoolP("all", "a", false, "查詢所有玩家")

	return cmd
}

func printPlayer(w io.Writer, p *pb.Player) {

	var utype string

	switch p.Type {
	case gamelogic.UserTypeGuest:
		utype = "訪客"
	case gamelogic.UserTypeNormal:
		utype = "玩家"
	case gamelogic.UserTypeRebot:
		utype = "機器人"
	}

	fmt.Fprintf(w, "============ %s[%d] ============\n", utype, p.PlayerID)
	if p.Conn != nil {
		fmt.Fprintf(w, "Connection Info:\n")
		fmt.Fprintf(w, "  IP Address: %s\n", p.Conn.IP)
		fmt.Fprintf(w, "  Port: %d\n", p.Conn.Port)
		fmt.Fprintf(w, "  Time: %s\n", p.Conn.CreateTime)
	}
	fmt.Fprintf(w, "User Data:\n")
	fmt.Fprintf(w, "  ID: %d\n", p.PlayerID)
	fmt.Fprintf(w, "  Name: %s\n", p.PlayerName)
	fmt.Fprintf(w, "  Table: %d\n", p.Table)
	fmt.Fprintf(w, "  Seat: %d\n", p.Seat)
	fmt.Fprintf(w, "  State: %s\n", p.State)
	fmt.Fprintf(w, "Game Data:\n")
	fmt.Fprintf(w, "  Score: %d\n", p.PlayerScore)
	fmt.Fprintf(w, "  BetAmount: %d\n", p.PlayerBetAmount)
	fmt.Fprintf(w, "  PayoutAmount: %d\n", p.PlayerPayoutAmount)
	fmt.Fprintf(w, "  FlowAmouont: %d\n", p.PlayerFlowAmount)
	fmt.Fprintf(w, "  TaxAmount: %d\n", p.PlayerTaxAmount)
}

func printPlayerLite(w io.Writer, p *pb.Player, indent string) {
	if p == nil {
		return
	}

	var utype string

	switch p.Type {
	case gamelogic.UserTypeGuest:
		utype = "訪客"
	case gamelogic.UserTypeNormal:
		utype = "玩家"
	case gamelogic.UserTypeRebot:
		utype = "機器人"
	}

	fmt.Fprintf(w, "%s[%s][%d] Name: %s, Table: %d, Seat: %d, State: %s, Score: %d\n", indent, utype, p.PlayerID, p.PlayerName, p.Table, p.Seat, p.State, p.PlayerScore)
}
