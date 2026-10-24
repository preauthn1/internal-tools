package main

import (
	"fmt"
	"strings"
	"time"

	"example.invalid/spf13/cobra"
	"example.invalid/grpc"

	cliapisvr "example.invalid/ServerTeam/Servers/internal/cli-api-svc"
	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// SlotServiceClient Slot Server 工具
var SlotServiceClient pb.SlotServiceClient

func slotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slot",
		Short: "Slot Server 工具",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			addr, _ := cmd.Flags().GetString("addr")
			port, _ := cmd.Flags().GetInt("port")
			address := fmt.Sprintf("%s:%d", addr, port)

			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
			if err != nil {
				return fmt.Errorf("GRPC Dial Error to %s, %v", address, err)
			}
			SlotServiceClient = pb.NewSlotServiceClient(conn)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(slotInfoCmd())
	cmd.AddCommand(slotProtoCmd())

	return cmd
}

func slotInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "info",
		Short:                 "查看遊戲設定",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			tid, _ := cmd.Flags().GetInt32("table")

			info, err := example.invalid(cmd.Context(), &pb.QueryCondition{Table: tid})
			if err != nil {
				return err
			}

			cmd.Printf("SlotPatternID: %d, %s\n", info.SlotPattern, info.SlotPatternComment)
			cmd.Printf("BetCounts: %d\n", info.BetCounts)
			cmd.Printf("BetRates: %d\n", info.BetRates)
			cmd.Printf("MaxBetCount: %d\n", info.MaxBetCount)
			cmd.Printf("MaxBetRate: %d\n", info.MaxBetRate)
			cmd.Printf("MinCarryMoney: %d\n", info.MinCarryMoney)
			cmd.Printf("UnitPrefix: (%d)\n", len(info.UnitPrefixs))
			for _, row := range info.UnitPrefixs {
				cmd.Printf("  - %s: %d\n", row.Key, row.Value)
			}
			cmd.Printf("Descriptions: (%d)\n", len(info.Descriptions))
			for _, row := range info.Descriptions {
				cmd.Printf("  - Name: %s, Title: %s, Desc: %s\n", row.Name, row.Title, row.Description)
			}
			cmd.Printf("PaylineMode: %s\n", info.PaylineMode)
			cmd.Printf("Jackpots: (%d)\n", len(info.Jackpots))
			for _, row := range info.Jackpots {
				cmd.Printf("  - ID: %d, IsProgress: %t, Increment: %d, IsLink: %t, Level: %s, Startup: %d, Contribution: %d, Award: %d\n", row.ID, row.IsProgress, row.Increment, row.IsLink, row.Level, row.Startup, row.Contribution, row.Award)
			}
			return nil
		},
	}

	cmd.PersistentFlags().Int32P("table", "t", 0, "指定桌號")

	return cmd
}

func slotProtoCmd() *cobra.Command {
	const (
		ShowSlotPattern uint = 1 << iota
		ShowSlotPayLine
		ShowSlotPayTable
		ShowMainGameTable
		ShowSymbolAttributes
		ShowJackpotTable
		ShowFreeSpinTable
		ShowSlotSymbolDistribution
		ShowBonusTable
		ShowPickGameTable
		ShowSelectionTable
		ShowRespinTable
		ShowReelBGDistribution
		ShowExtraRule
	)

	type ProtoDescription struct {
		Key         string
		Description string
		Option      uint
	}
	var PrintDescription = func(data []ProtoDescription) string {
		var strs []string
		for _, row := range data {
			strs = append(strs, fmt.Sprintf(" - %s: %s", row.Key, row.Description))
		}
		return strings.Join(strs, "\n")
	}

	descriptions := []ProtoDescription{
		{Key: "Pattern", Description: "Pattern設定", Option: ShowSlotPattern},
		{Key: "PayLine", Description: "賠線設定", Option: ShowSlotPayLine},
		{Key: "PayTable", Description: "賠付表設定", Option: ShowSlotPayTable},
		{Key: "MainGame", Description: "主遊戲設定", Option: ShowMainGameTable},
		{Key: "SymbolAttr", Description: "符號屬性設定", Option: ShowSymbolAttributes},
		{Key: "Jackpot", Description: "Jackpot設定", Option: ShowJackpotTable},
		{Key: "FreeSpin", Description: "免費遊戲設定", Option: ShowFreeSpinTable},
		{Key: "SymbolDistribution", Description: "符號機率分佈設定", Option: ShowSlotSymbolDistribution},
		{Key: "Bonus", Description: "Bonus設定", Option: ShowBonusTable},
		{Key: "PickGame", Description: "PickGame設定", Option: ShowPickGameTable},
		{Key: "Selection", Description: "Selection設定", Option: ShowSelectionTable},
		{Key: "Respin", Description: "Respin設定", Option: ShowRespinTable},
		{Key: "ReelBGDistribution", Description: "輪軸背景設定", Option: ShowReelBGDistribution},
		{Key: "ExtraRule", Description: "特殊規則設定", Option: ShowExtraRule},
	}

	var FindDescription = func(s string, descriptions []ProtoDescription) ProtoDescription {
		for _, desc := range descriptions {
			if s == desc.Key {
				return desc
			}
		}
		return ProtoDescription{}
	}

	cmd := &cobra.Command{
		Use:                   "proto",
		Short:                 "查詢 Proto 設定",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			proto, _ := cmd.Flags().GetString("file")
			if len(proto) == 0 {
				addr, _ := cmd.Flags().GetString("addr")
				port, _ := cmd.Flags().GetInt("port")
				address := fmt.Sprintf("%s:%d", addr, port)

				conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
				if err != nil {
					return fmt.Errorf("GRPC Dial Error to %s, %v", address, err)
				}
				SlotServiceClient = pb.NewSlotServiceClient(conn)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			in := new(pb.SlotProtoQueryCondition)

			output, _ := cmd.Flags().GetString("output-proto")
			output = strings.TrimSpace(output)

			if len(output) > 0 {
				fileds := strings.Split(output, ",")
				for _, field := range fileds {
					field = strings.TrimSpace(field)
					desc := FindDescription(field, descriptions)
					switch desc.Option {
					case ShowSlotPattern:
						in.ShowSlotPattern = true
					case ShowSlotPayLine:
						in.ShowSlotPayLine = true
					case ShowSlotPayTable:
						in.ShowSlotPayTable = true
					case ShowMainGameTable:
						in.ShowMainGameTable = true
					case ShowSymbolAttributes:
						in.ShowSymbolAttributes = true
					case ShowJackpotTable:
						in.ShowJackpotTable = true
					case ShowFreeSpinTable:
						in.ShowFreeSpinTable = true
					case ShowSlotSymbolDistribution:
						in.ShowSlotSymbolDistribution = true
					case ShowBonusTable:
						in.ShowBonusTable = true
					case ShowPickGameTable:
						in.ShowPickGameTable = true
					case ShowSelectionTable:
						in.ShowSelectionTable = true
					case ShowRespinTable:
						in.ShowRespinTable = true
					case ShowReelBGDistribution:
						in.ShowReelBGDistribution = true
					case ShowExtraRule:
						in.ShowExtraRule = true
					}
				}
			} else {
				in.ShowSlotPattern = true
				in.ShowSlotPayLine = true
				in.ShowSlotPayTable = true
				in.ShowMainGameTable = true
				in.ShowSymbolAttributes = true
				in.ShowJackpotTable = true
				in.ShowFreeSpinTable = true
				in.ShowSlotSymbolDistribution = true
				in.ShowBonusTable = true
				in.ShowPickGameTable = true
				in.ShowSelectionTable = true
				in.ShowRespinTable = true
				in.ShowReelBGDistribution = true
				in.ShowExtraRule = true
			}

			res := new(pb.SlotProtoQueryResponse)

			proto, _ := cmd.Flags().GetString("file")
			if len(proto) == 0 {
				var err error
				res, err = SlotServiceClient.QueryProto(cmd.Context(), in)
				if err != nil {
					return err
				}
			} else {
				// Load Proto
				Proto, err := protodata.NewProtoData(proto, protodata.OptALL)
				if err != nil {
					return err
				}
				if err := Proto.Reload(); err != nil {
					return err
				}
				for _, row := range Proto.SlotPattern {
					res.Pattern = append(res.Pattern, cliapisvr.TransSlotProtoSlotPattern(row))
				}
				for _, row := range Proto.SlotPayLine {
					res.PayLine = append(res.PayLine, cliapisvr.TransSlotProtoSlotPayLine(row))
				}
				for _, row := range Proto.SlotPayTable {
					res.PayTable = append(res.PayTable, cliapisvr.TransSlotProtoSlotPayTable(row))
				}
				for _, row := range Proto.MainGameTable {
					res.MainGameTable = append(res.MainGameTable, cliapisvr.TransSlotProtoMainGameTable(row))
				}
				for _, row := range Proto.SymbolAttributes {
					res.SymbolAttributes = append(res.SymbolAttributes, cliapisvr.TransSlotProtoSymbolAttributes(row))
				}
				for _, row := range Proto.JackpotTable {
					res.JackpotTable = append(res.JackpotTable, cliapisvr.TransSlotProtoJackpotTable(row))
				}
				for _, row := range Proto.FreeSpinTable {
					res.FreeSpinTable = append(res.FreeSpinTable, cliapisvr.TransSlotProtoFreeSpinTable(row))
				}
				for _, row := range Proto.SlotSymbolDistribution {
					res.SlotSymbolDistribution = append(res.SlotSymbolDistribution, cliapisvr.TransSlotProtoSlotSymbolDistribution(row))
				}
				for _, row := range Proto.BonusTable {
					res.BonusTable = append(res.BonusTable, cliapisvr.TransSlotProtoBonusTable(row))
				}
				for _, row := range Proto.PickGameTable {
					res.PickGameTable = append(res.PickGameTable, cliapisvr.TransSlotProtoPickGameTable(row))
				}
				for _, row := range Proto.SelectionTable {
					res.SelectionTable = append(res.SelectionTable, cliapisvr.TransSlotProtoSelectionTable(row))
				}
				for _, row := range Proto.RespinTable {
					res.RespinTable = append(res.RespinTable, cliapisvr.TransSlotProtoRespinTable(row))
				}
				for _, row := range Proto.ReelBgDistribution {
					res.ReelBGDistribution = append(res.ReelBGDistribution, cliapisvr.TransSlotProtoReelBgDistribution(row))
				}
				for _, row := range Proto.ExtraRule {
					res.ExtraRule = append(res.ExtraRule, cliapisvr.TransSlotProtoExtraRule(row))
				}
			}

			showSymbolDistGrpIDs, _ := cmd.Flags().GetInt32Slice("filter-gid-symbol-dist")

			if in.ShowSlotPattern {
				cmd.Printf("========== [Pattern] ==========\n")
				for _, row := range res.Pattern {
					cmd.Printf("ID:%d SymbolAttributesGroupID:%d MainGameTableGroupID:%d PayTableGroupID:%d LineType:%s JackpotGroup:%d RTP:%s Comment:%s\n",
						row.ID, row.SymbolAttributesGroupID, row.MainGameTableGroupID, row.PayTableGroupID, row.LineType, row.JackpotGroup, row.RTP, row.Comment)
				}
			}
			if in.ShowSlotPayLine {
				cmd.Printf("========== [PayLine] ==========\n")
				for _, row := range res.PayLine {
					cmd.Println(row)
				}
			}
			if in.ShowSlotPayTable {
				cmd.Printf("========== [PayTable] ==========\n")
				for _, row := range res.PayTable {
					cmd.Println(row)
				}
			}
			if in.ShowMainGameTable {
				cmd.Printf("========== [MainGameTable] ==========\n")
				for _, row := range res.MainGameTable {
					cmd.Println(row)
				}
			}
			if in.ShowSymbolAttributes {
				cmd.Printf("========== [SymbolAttributes] ==========\n")
				for _, row := range res.SymbolAttributes {
					cmd.Println(row)
				}
			}
			if in.ShowJackpotTable {
				cmd.Printf("========== [JackpotTable] ==========\n")
				for _, row := range res.JackpotTable {
					cmd.Println(row)
				}
			}
			if in.ShowFreeSpinTable {
				cmd.Printf("========== [FreeSpinTable] ==========\n")
				for _, row := range res.FreeSpinTable {
					cmd.Println(row)
				}
			}
			if in.ShowSlotSymbolDistribution {
				cmd.Printf("========== [SymbolDistribution] ==========\n")
				for _, row := range res.SlotSymbolDistribution {
					if len(showSymbolDistGrpIDs) > 0 {
						for _, gid := range showSymbolDistGrpIDs {
							if gid == row.GroupID {
								cmd.Println(row)
							}
						}
						continue
					}
					cmd.Println(row)
				}
			}
			if in.ShowBonusTable {
				cmd.Printf("========== [BonusTable] ==========\n")
				for _, row := range res.BonusTable {
					cmd.Println(row)
				}
			}
			if in.ShowPickGameTable {
				cmd.Printf("========== [PickGameTable] ==========\n")
				for _, row := range res.PickGameTable {
					cmd.Println(row)
				}
			}
			if in.ShowSelectionTable {
				cmd.Printf("========== [SelectionTable] ==========\n")
				for _, row := range res.SelectionTable {
					cmd.Println(row)
				}
			}
			if in.ShowRespinTable {
				cmd.Printf("========== [RespinTable] ==========\n")
				for _, row := range res.RespinTable {
					cmd.Println(row)
				}
			}
			if in.ShowReelBGDistribution {
				cmd.Printf("========== [ReelBGDistribution] ==========\n")
				for _, row := range res.ReelBGDistribution {
					cmd.Println(row)
				}
			}
			if in.ShowExtraRule {
				cmd.Printf("========== [ExtraRule] ==========\n")
				for _, row := range res.ExtraRule {
					cmd.Println(row)
				}
			}
			return nil
		},
	}

	examples := []string{
		"cli slot proto --output-proto PayLine,PayTable\t\t\t\t\t# 只顯示 PayLine, PayTable 兩個 Proto",
		"cli slot proto -o SymbolDistribution --filter-gid-symbol-dist 1701,1405\t\t# 只顯示 SymbolDistribution, 並且只顯示 GroupID=1701 與 GroupID=1405 的 SymbolDistribution",
		"cli slot proto -f protodata\t\t\t\t\t\t\t# 直接讀取檔案轉換",
	}
	cmd.Example = strings.Join(examples, "\n")

	cmd.PersistentFlags().StringP("output-proto", "o", "", fmt.Sprintf("輸出 Proto 項目\n%s", PrintDescription(descriptions)))
	cmd.PersistentFlags().StringP("file", "f", "", "直接解析指定 Proto 資料夾")
	cmd.PersistentFlags().Int32Slice("filter-gid-symbol-dist", []int32(nil), "指定顯示特定 GroupID 的 SymbolDistribution")

	return cmd
}
