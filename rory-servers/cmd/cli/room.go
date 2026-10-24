package main

import (
	"fmt"
	"io"
	"sort"
	"strconv"

	"example.invalid/spf13/cobra"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
)

func roomCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "room",
		Short:                 "房間資訊",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			room, err := CLIClient.RoomInfo(cmd.Context(), new(pb.Empty))
			if err != nil {
				return err
			}

			cmd.Printf("=== 房間資訊 ===\n")
			printRoom(cmd.OutOrStdout(), room, "")

			return nil
		},
	}

	return cmd
}

func tableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "table <tableid>",
		Short:                 "牌桌資訊",
		Args:                  cobra.ExactValidArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			tid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("桌號輸入錯誤, %v", err)
			}

			t, err := CLIClient.TableInfo(cmd.Context(), &pb.QueryCondition{Table: int32(tid)})
			if err != nil {
				return err
			}
			sort.Slice(t.Players, func(i, j int) bool {
				return t.Players[i].Seat < t.Players[j].Seat
			})

			cmd.Printf("=== 牌桌資訊 ===\n")
			printTable(cmd.OutOrStdout(), t, "")

			if len(t.GameSummary) > 0 {
				cmd.Printf("=== 遊戲狀態 ===\n")
				cmd.Println(t.GameSummary)
			}
			return nil
		},
	}

	cmd.AddCommand(tableListCmd())

	return cmd
}

func tableListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出牌桌",
		RunE: func(cmd *cobra.Command, args []string) error {

			room, err := CLIClient.RoomInfo(cmd.Context(), new(pb.Empty))
			if err != nil {
				return err
			}

			var gaming, free int

			for _, t := range room.Tables {
				if t.Gaming {
					gaming++
				} else {
					free++
				}
			}

			isgaming, _ := cmd.PersistentFlags().GetBool("gaming")
			isfree, _ := cmd.PersistentFlags().GetBool("free")

			detail, _ := cmd.PersistentFlags().GetBool("detail")

			cmd.Printf("桌數: %d = 遊戲中(%d) + 空閒中(%d)\n", len(room.Tables), gaming, free)
			for _, t := range room.Tables {
				if isgaming && !t.Gaming {
					continue
				}
				if isfree && t.Gaming {
					continue
				}
				if detail {
					printTable(cmd.OutOrStdout(), t, "  ")
				} else {
					printTableLite(cmd.OutOrStdout(), t, "  ")
				}
			}

			return nil
		},
	}

	cmd.PersistentFlags().BoolP("gaming", "g", false, "只列出遊戲中牌桌")
	cmd.PersistentFlags().BoolP("free", "f", false, "只列出空閒中牌桌")
	cmd.PersistentFlags().BoolP("all", "a", false, "列出所有牌桌")

	cmd.PersistentFlags().BoolP("detail", "d", false, "列出詳細資訊")

	return cmd
}

func printRoom(w io.Writer, room *pb.Room, indent string) {
	var gaming, free int

	for _, t := range room.Tables {
		if t.Gaming {
			gaming++
		} else {
			free++
		}
	}

	fmt.Fprintf(w, "%s遊戲ID: %d\n", indent, room.GameID)
	fmt.Fprintf(w, "%s遊戲名稱: %s\n", indent, room.GameName)
	fmt.Fprintf(w, "%s房間ID: %d\n", indent, room.ID)
	fmt.Fprintf(w, "%s房間名稱: %s\n", indent, room.Name)
	fmt.Fprintf(w, "%sFlag: %d\n", indent, room.Flag)
	fmt.Fprintf(w, "%s桌數: %d = 遊戲中(%d) + 空閒中(%d)\n", indent, len(room.Tables), gaming, free)
}

func printTable(w io.Writer, t *pb.Table, indent string) {

	printTableLite(w, t, indent)
	for _, p := range t.Players {
		printPlayerLite(w, p, fmt.Sprintf("%s  座位[%d]: ", indent, p.Seat))
	}
}

func printTableLite(w io.Writer, t *pb.Table, indent string) {

	var usercnt, rbtcnt int

	for _, p := range t.Players {
		switch p.Type {
		case gamelogic.UserTypeNormal:
			usercnt++
		case gamelogic.UserTypeRebot:
			rbtcnt++
		}
	}

	var gaming string = "空閒中"
	if t.Gaming {
		gaming = "遊戲中"
	}

	fmt.Fprintf(w, "%s牌桌[%d][%s]\t玩家: %3d 人\t機器人: %3d 人\n", indent, t.ID, gaming, usercnt, rbtcnt)
}
