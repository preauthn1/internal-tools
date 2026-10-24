package main

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"example.invalid/spf13/cobra"
	"example.invalid/grpc"

	cliapisvr "example.invalid/ServerTeam/Servers/internal/cli-api-svc"
	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
	"example.invalid/ServerTeam/Servers/pkg/pachislot/proto-provider/protofile"
)

// PachiSlotServiceClient PachiSlot Server 工具
var PachiSlotServiceClient pb.PachiSlotServiceClient

func pachiSlotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pachislot",
		Short: "PachiSlot Server 工具",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			addr, _ := cmd.Flags().GetString("addr")
			port, _ := cmd.Flags().GetInt("port")
			address := fmt.Sprintf("%s:%d", addr, port)

			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
			if err != nil {
				return fmt.Errorf("GRPC Dial Error to %s, %v", address, err)
			}
			PachiSlotServiceClient = pb.NewPachiSlotServiceClient(conn)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(pachiSlotProtoCmd())
	cmd.AddCommand(pachiSlotAssignCmd())
	cmd.AddCommand(pachiSlotMachineCmd())

	return cmd
}

func pachiSlotProtoCmd() *cobra.Command {
	const (
		ShowGameConfig uint = 1 << iota
		ShowPayLine
		ShowPayTable
		ShowReelStrips
		ShowPattern
		ShowModule
		ShowMode
		ShowRTRelease
		ShowResultDistribution
		ShowNextModeDistribution
		ShowRTDistribution
		ShowBRDistribution
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
		{Key: "GameConfig", Description: "遊戲設定", Option: ShowGameConfig},
		{Key: "PayLine", Description: "賠線設定", Option: ShowPayLine},
		{Key: "PayTable", Description: "賠付表設定", Option: ShowPayTable},
		{Key: "ReelStrips", Description: "輪軸帶設定", Option: ShowReelStrips},
		{Key: "Pattern", Description: "Pattern設定", Option: ShowPattern},
		{Key: "Module", Description: "模組設定", Option: ShowModule},
		{Key: "Mode", Description: "模式設定", Option: ShowMode},
		{Key: "RTRelease", Description: "RT解除分佈設定", Option: ShowRTRelease},
		{Key: "ResultDistribution", Description: "牌型分佈設定", Option: ShowResultDistribution},
		{Key: "NextModeDistribution", Description: "下次模式分佈設定", Option: ShowNextModeDistribution},
		{Key: "RTDistribution", Description: "RT次數分佈設定", Option: ShowRTDistribution},
		{Key: "BRDistribution", Description: "大獎分佈設定", Option: ShowBRDistribution},
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
				PachiSlotServiceClient = pb.NewPachiSlotServiceClient(conn)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			in := new(pb.PachiSlotQueryProtoCondition)

			output, _ := cmd.Flags().GetString("output-proto")
			output = strings.TrimSpace(output)

			if len(output) > 0 {
				fileds := strings.Split(output, ",")
				for _, field := range fileds {
					field = strings.TrimSpace(field)
					desc := FindDescription(field, descriptions)
					switch desc.Option {
					case ShowGameConfig:
						in.ShowGameConfig = true
					case ShowPayLine:
						in.ShowPayLine = true
					case ShowPayTable:
						in.ShowPayTable = true
					case ShowReelStrips:
						in.ShowReelStrips = true
					case ShowPattern:
						in.ShowPattern = true
					case ShowModule:
						in.ShowModule = true
					case ShowMode:
						in.ShowMode = true
					case ShowRTRelease:
						in.ShowRTRelease = true
					case ShowResultDistribution:
						in.ShowResultDistribution = true
					case ShowNextModeDistribution:
						in.ShowNextModeDistribution = true
					case ShowRTDistribution:
						in.ShowRTDistribution = true
					case ShowBRDistribution:
						in.ShowBRDistribution = true
					}
				}
			} else {
				in.ShowGameConfig = true
				in.ShowPayLine = true
				in.ShowPayTable = true
				in.ShowReelStrips = true
				in.ShowPattern = true
				in.ShowModule = true
				in.ShowMode = true
				in.ShowRTRelease = true
				in.ShowResultDistribution = true
				in.ShowNextModeDistribution = true
				in.ShowRTDistribution = true
				in.ShowBRDistribution = true
			}

			var res *pb.PachiSlotQueryProtoResponse

			proto, _ := cmd.Flags().GetString("file")
			if len(proto) == 0 {
				var err error
				res, err = PachiSlotServiceClient.QueryProto(cmd.Context(), in)
				if err != nil {
					return err
				}
			} else {
				p, err := protofile.NewProvider(proto)
				if err != nil {
					return err
				}
				for _, t := range protofile.DefaultTransformers {
					if err := p.RegisterTransformer(t); err != nil {
						return nil
					}
				}
				proto, err := p.Proto()
				if err != nil {
					return err
				}
				res = new(pb.PachiSlotQueryProtoResponse)
				res.GameConfig = cliapisvr.TransPachiSlotProtoGameConfig(proto.GameConfig)
				for _, row := range proto.PayLine {
					res.PayLine = append(res.PayLine, cliapisvr.TransPachiSlotProtoPayLine(row))
				}
				for _, row := range proto.PayTable {
					res.PayTable = append(res.PayTable, cliapisvr.TransPachiSlotProtoPayoutOption(row))
				}
				for _, row := range proto.ReelStrips {
					res.ReelStrips = append(res.ReelStrips, cliapisvr.TransPachiSlotProtoReelStripSymbol(row))
				}
				for _, row := range proto.Pattern {
					res.Pattern = append(res.Pattern, cliapisvr.TransPachiSlotProtoPattern(row))
				}
				for _, row := range proto.ModuleConfigs {
					res.Module = append(res.Module, cliapisvr.TransPachiSlotProtoModule(row))
				}
				for _, row := range proto.ModeConfigs {
					res.Mode = append(res.Mode, cliapisvr.TransPachiSlotProtoMode(row))
				}
				for _, row := range proto.RTReleaseDistribution {
					res.RTRelease = append(res.RTRelease, cliapisvr.TransPachiSlotProtoRTReleaseOption(row))
				}
				for _, row := range proto.ResultDistribution {
					res.ResultDistribution = append(res.ResultDistribution, cliapisvr.TransPachiSlotProtoResultDistributionOption(row))
				}
				for _, row := range proto.NextModeDistribution {
					res.NextModeDistribution = append(res.NextModeDistribution, cliapisvr.TransPachiSlotProtoNextModeDistributionOption(row))
				}
				for _, row := range proto.RTDistribution {
					res.RTDistribution = append(res.RTDistribution, cliapisvr.TransPachiSlotProtoRTDistributionOption(row))
				}
				for _, row := range proto.BRDistribution {
					res.BRDistribution = append(res.BRDistribution, cliapisvr.TransPachiSlotProtoBRDistributionOption(row))
				}
			}

			showPatternGrpIDs, _ := cmd.Flags().GetInt32Slice("filter-gid-pattern")
			showModuleGrpIDs, _ := cmd.Flags().GetInt32Slice("filter-gid-module")
			showModeGrpIDs, _ := cmd.Flags().GetInt32Slice("filter-gid-mode")
			showRTReleaseGrpIDs, _ := cmd.Flags().GetInt32Slice("filter-gid-rt-release")
			showResultDistGrpIDs, _ := cmd.Flags().GetInt32Slice("filter-gid-result-dist")
			showNextModeDistGrpIDs, _ := cmd.Flags().GetInt32Slice("filter-gid-next-mode-dist")
			showRTDistGrpIDs, _ := cmd.Flags().GetInt32Slice("filter-gid-rt-dist")
			showBRDistGrpIDs, _ := cmd.Flags().GetInt32Slice("filter-gid-br-dist")

			if in.ShowGameConfig {
				cmd.Printf("========== [GameConfig] ==========\n")
				cmd.Println(res.GameConfig)
			}
			if in.ShowPayLine {
				cmd.Printf("========== [PayLine] ==========\n")
				for _, row := range res.PayLine {
					cmd.Println(row)
				}
			}
			if in.ShowPayTable {
				cmd.Printf("========== [PayTable] ==========\n")
				for _, row := range res.PayTable {
					cmd.Println(row)
				}
			}
			if in.ShowReelStrips {
				cmd.Printf("========== [ReelStrips] ==========\n")
				for _, row := range res.ReelStrips {
					cmd.Printf("ID:%d GroupID:%d Reel:%d Row:%d Symbol:%v\n", row.ID, row.GroupID, row.Reel, row.Row, row.Symbol)
				}
			}
			if in.ShowPattern {
				cmd.Printf("========== [Pattern] ==========\n")
				for _, row := range res.Pattern {
					if len(showPatternGrpIDs) > 0 {
						for _, gid := range showPatternGrpIDs {
							if gid == row.GroupID {
								cmd.Println(row)
							}
						}
						continue
					}
					cmd.Println(row)
				}
			}
			if in.ShowModule {
				cmd.Printf("========== [Module] ==========\n")
				for _, row := range res.Module {
					if len(showModuleGrpIDs) > 0 {
						for _, gid := range showModuleGrpIDs {
							if gid == row.GroupID {
								cmd.Println(row)
							}
						}
						continue
					}
					cmd.Println(row)
				}
			}
			if in.ShowMode {
				cmd.Printf("========== [Mode] ==========\n")
				for _, row := range res.Mode {
					if len(showModeGrpIDs) > 0 {
						for _, gid := range showModeGrpIDs {
							if gid == row.GroupID {
								cmd.Println(row)
							}
						}
						continue
					}
					cmd.Println(row)
				}
			}
			if in.ShowRTRelease {
				cmd.Printf("========== [RTRelease] ==========\n")
				for _, row := range res.RTRelease {
					if len(showRTReleaseGrpIDs) > 0 {
						for _, gid := range showRTReleaseGrpIDs {
							if gid == row.GroupID {
								cmd.Println(row)
							}
						}
						continue
					}
					cmd.Println(row)
				}
			}
			if in.ShowResultDistribution {
				cmd.Printf("========== [ResultDistribution] ==========\n")
				for _, row := range res.ResultDistribution {
					if len(showResultDistGrpIDs) > 0 {
						for _, gid := range showResultDistGrpIDs {
							if gid == row.GroupID {
								cmd.Println(row)
							}
						}
						continue
					}
					cmd.Println(row)
				}
			}
			if in.ShowNextModeDistribution {
				cmd.Printf("========== [NextModeDistribution] ==========\n")
				for _, row := range res.NextModeDistribution {
					if len(showNextModeDistGrpIDs) > 0 {
						for _, gid := range showNextModeDistGrpIDs {
							if gid == row.GroupID {
								cmd.Println(row)
							}
						}
						continue
					}
					cmd.Println(row)
				}
			}
			if in.ShowRTDistribution {
				cmd.Printf("========== [RTDistribution] ==========\n")
				for _, row := range res.RTDistribution {
					if len(showRTDistGrpIDs) > 0 {
						for _, gid := range showRTDistGrpIDs {
							if gid == row.GroupID {
								cmd.Println(row)
							}
						}
						continue
					}
					cmd.Println(row)
				}
			}
			if in.ShowBRDistribution {
				cmd.Printf("========== [BRDistribution] ==========\n")
				for _, row := range res.BRDistribution {
					if len(showBRDistGrpIDs) > 0 {
						for _, gid := range showBRDistGrpIDs {
							if gid == row.GroupID {
								cmd.Println(row)
							}
						}
						continue
					}
					cmd.Println(row)
				}
			}

			return nil
		},
	}

	examples := []string{
		"cli pachislot proto --output-proto PayLine,PayTable\t\t\t# 只顯示 PayLine, PayTable 兩個 Proto",
		"cli pachislot proto -o Mode --filter-gid-mode 600103\t# 只顯示 Mode, 並且只顯示 GroupID=600103 的 Mode",
		"cli pachislot proto -f protodata\t# 直接讀取檔案轉換",
	}
	cmd.Example = strings.Join(examples, "\n")

	cmd.PersistentFlags().StringP("output-proto", "o", "", fmt.Sprintf("輸出 Proto 項目\n%s", PrintDescription(descriptions)))
	cmd.PersistentFlags().StringP("file", "f", "", "直接解析指定 Proto 資料夾")
	cmd.PersistentFlags().Int32Slice("filter-gid-pattern", []int32(nil), "指定顯示特定 GroupID 的 Pattern")
	cmd.PersistentFlags().Int32Slice("filter-gid-module", []int32(nil), "指定顯示特定 GroupID 的 Module")
	cmd.PersistentFlags().Int32Slice("filter-gid-mode", []int32(nil), "指定顯示特定 GroupID 的 Mode")
	cmd.PersistentFlags().Int32Slice("filter-gid-rt-release", []int32(nil), "指定顯示特定 GroupID 的 RTRelease")
	cmd.PersistentFlags().Int32Slice("filter-gid-result-dist", []int32(nil), "指定顯示特定 GroupID 的 ResultDistribution")
	cmd.PersistentFlags().Int32Slice("filter-gid-next-mode-dist", []int32(nil), "指定顯示特定 GroupID 的 NextModeDistribution")
	cmd.PersistentFlags().Int32Slice("filter-gid-rt-dist", []int32(nil), "指定顯示特定 GroupID 的 RTDistribution")
	cmd.PersistentFlags().Int32Slice("filter-gid-br-dist", []int32(nil), "指定顯示特定 GroupID 的 BRDistribution")

	return cmd
}

func pachiSlotAssignCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "assign",
		Short:                 "修改機台狀態",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(pachiSlotAssignBonusCmd())
	cmd.AddCommand(pachiSlotAssignResultCmd())
	cmd.AddCommand(pachiSlotAssignModeCmd())
	cmd.AddCommand(pachiSlotAssignReplayCmd())
	cmd.AddCommand(pachiSlotAssignBonusSureCmd())
	cmd.AddCommand(pachiSlotAssignRTCmd())
	cmd.AddCommand(pachiSlotAssignBonusTimesLeftCmd())
	cmd.AddCommand(pachiSlotAssignJACTimesLeftCmd())
	cmd.AddCommand(pachiSlotAssignBonusStoreCmd())
	cmd.AddCommand(pachiSlotAssignRoadMapCmd())

	return cmd
}

func pachiSlotAssignBonusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bonus <machineID> <Bonus>",
		Short: "指定機台當前大獎",
		Long: `指定機台當前大獎

Bonus:
  BB
  RB
`,
		Args:                  cobra.ExactValidArgs(2),
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			tid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("桌號輸入錯誤, %v", err)
			}

			key, ok := pb.PachiSlot_Scenario_value["PachiSlot_Scenario_"+args[1]]
			if !ok {
				return fmt.Errorf("大獎輸入錯誤")
			}

			s := pb.PachiSlot_Scenario(key)

			if s == pb.PachiSlot_Scenario_PachiSlot_Scenario_NG || s == pb.PachiSlot_Scenario_PachiSlot_Scenario_ALL {
				return fmt.Errorf("大獎輸入錯誤")
			}

			in := new(pb.PachiSlotMachine)
			in.No = int32(tid)
			in.Bonus = s

			if _, err := PachiSlotServiceClient.AssignBonus(cmd.Context(), in); err != nil {
				return err
			}

			cmd.Printf("機台[%d] 指定當前大獎[%v]\n", in.No, cliapisvr.ReversePachiSlotScenario(in.Bonus))

			return nil
		},
	}

	return cmd
}

func pachiSlotAssignResultCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "result <machineID> <ResultType>",
		Short: "指定機台下一把牌型",
		Long: `指定機台下一把牌型

ResultType:
  None
  FeatureNone
  JACNone
  BlueSevenLink
  RedSevenLink
  BarLink
  BellLink
  WaterMelonLink
  BerryLink
  ReplayLink
  DoubleSevenLink
  ReleaseRT
  BonusSure
  BerryLeftAny
`,
		Args:                  cobra.ExactValidArgs(2),
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			tid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("桌號輸入錯誤, %v", err)
			}

			key, ok := pb.PachiSlot_ResultType_value["PachiSlot_ResultType_"+args[1]]
			if !ok {
				return fmt.Errorf("牌型輸入錯誤")
			}

			rt := pb.PachiSlot_ResultType(key)

			if rt == pb.PachiSlot_ResultType_PachiSlot_ResultType_Empty {
				return fmt.Errorf("牌型輸入錯誤, 不可指定為 Empty")
			}

			if rt == pb.PachiSlot_ResultType_PachiSlot_ResultType_ReleaseRTonBB {
				return fmt.Errorf("牌型輸入錯誤, 不可指定為 ReleaseRTonBB")
			}

			if rt == pb.PachiSlot_ResultType_PachiSlot_ResultType_ReleaseRTonRB {
				return fmt.Errorf("牌型輸入錯誤, 不可指定為 ReleaseRTonRB")
			}

			in := new(pb.PachiSlotMachine)
			in.No = int32(tid)
			in.NextResult = rt

			if _, err := PachiSlotServiceClient.AssignResult(cmd.Context(), in); err != nil {
				return err
			}

			cmd.Printf("機台[%d] 指定牌型[%v]\n", in.No, cliapisvr.ReversePachiSlotResultType(in.NextResult))

			return nil
		},
	}

	return cmd
}

func pachiSlotAssignModeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mode <machineID> <Mode>",
		Short: "指定機台下一次模式",
		Long: `指定機台下一次模式

Mode:
  A         
  B
  Heaven
`,
		Args:                  cobra.ExactValidArgs(2),
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			tid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("桌號輸入錯誤, %v", err)
			}

			key, ok := pb.PachiSlot_ModeType_value["PachiSlot_ModeType_"+args[1]]
			if !ok {
				return fmt.Errorf("模式輸入錯誤")
			}

			m := pb.PachiSlot_ModeType(key)

			if m == pb.PachiSlot_ModeType_PachiSlot_ModeType_Unassigned {
				return fmt.Errorf("模式輸入錯誤")
			}

			in := new(pb.PachiSlotMachine)
			in.No = int32(tid)
			in.NextMode = m

			if _, err := PachiSlotServiceClient.AssignMode(cmd.Context(), in); err != nil {
				return err
			}

			cmd.Printf("機台[%d] 指定模式[%v]\n", in.No, cliapisvr.ReversePachiSlotMode(in.NextMode))

			return nil
		},
	}

	return cmd
}

func pachiSlotAssignReplayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "replay <machineID>",
		Short:                 "指定機台為Replay",
		Args:                  cobra.ExactValidArgs(1),
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			tid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("桌號輸入錯誤, %v", err)
			}

			in := new(pb.PachiSlotMachine)
			in.No = int32(tid)
			in.Replay = true

			if _, err := PachiSlotServiceClient.AssignReplay(cmd.Context(), in); err != nil {
				return err
			}

			cmd.Printf("機台[%d] 指定為 Replay\n", in.No)

			return nil
		},
	}

	return cmd
}

func pachiSlotAssignBonusSureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "bonus-sure <machineID>",
		Short:                 "指定機台為Bonus確定",
		Args:                  cobra.ExactValidArgs(1),
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			tid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("桌號輸入錯誤, %v", err)
			}

			in := new(pb.PachiSlotMachine)
			in.No = int32(tid)
			in.BonusSure = true

			if _, err := PachiSlotServiceClient.AssignBonusSure(cmd.Context(), in); err != nil {
				return err
			}

			cmd.Printf("機台[%d] 指定為 Bonus確定\n", in.No)

			return nil
		},
	}

	return cmd
}

func pachiSlotAssignRTCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "RT <machineID> <times>",
		Short:                 "指定機台剩餘RT數",
		Args:                  cobra.ExactValidArgs(2),
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			tid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("桌號輸入錯誤, %v", err)
			}

			rt, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("剩餘RT數輸入錯誤, %v", err)
			}

			in := new(pb.PachiSlotMachine)
			in.No = int32(tid)
			in.RT = int32(rt)

			if _, err := PachiSlotServiceClient.AssignRT(cmd.Context(), in); err != nil {
				return err
			}

			cmd.Printf("機台[%d] 指定剩餘RT數為 %d\n", in.No, in.RT)

			return nil
		},
	}

	return cmd
}

func pachiSlotAssignBonusTimesLeftCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "bonus-times <machineID> <times>",
		Short:                 "指定機台剩餘Bonus次數",
		Args:                  cobra.ExactValidArgs(2),
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			tid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("桌號輸入錯誤, %v", err)
			}

			times, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("Bonus次數輸入錯誤, %v", err)
			}

			in := new(pb.PachiSlotMachine)
			in.No = int32(tid)
			in.BonusTimesLeft = int32(times)

			if _, err := PachiSlotServiceClient.AssignBonusTimesLeft(cmd.Context(), in); err != nil {
				return err
			}

			cmd.Printf("機台[%d] 指定剩餘Bonus次數為 %d\n", in.No, in.BonusTimesLeft)

			return nil
		},
	}

	return cmd
}

func pachiSlotAssignJACTimesLeftCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "jac-times <machineID> <times>",
		Short:                 "指定機台剩餘JAC次數",
		Args:                  cobra.ExactValidArgs(2),
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			tid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("桌號輸入錯誤, %v", err)
			}

			times, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("JAC次數輸入錯誤, %v", err)
			}

			in := new(pb.PachiSlotMachine)
			in.No = int32(tid)
			in.JACTimesLeft = int32(times)

			if _, err := PachiSlotServiceClient.AssignJACTimesLeft(cmd.Context(), in); err != nil {
				return err
			}

			cmd.Printf("機台[%d] 指定剩餘JAC次數為 %d\n", in.No, in.JACTimesLeft)

			return nil
		},
	}

	return cmd
}

func pachiSlotAssignBonusStoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "bonus-store <machineID> <times>",
		Short:                 "指定機台剩餘Bonus次數",
		Args:                  cobra.ExactValidArgs(2),
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			tid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("桌號輸入錯誤, %v", err)
			}

			times, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("Bonus次數輸入錯誤, %v", err)
			}

			in := new(pb.PachiSlotMachine)
			in.No = int32(tid)
			in.BonusStore = int32(times)

			if _, err := PachiSlotServiceClient.AssignBonusStore(cmd.Context(), in); err != nil {
				return err
			}

			cmd.Printf("機台[%d] 指定剩餘Bonus次數為 %d\n", in.No, in.BonusStore)

			return nil
		},
	}

	return cmd
}

func pachiSlotAssignRoadMapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "road-map <machineID>",
		Short:                 "指定機台轉數表",
		Args:                  cobra.ExactValidArgs(1),
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			tid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("桌號輸入錯誤, %v", err)
			}

			total, _ := cmd.Flags().GetInt64("total")
			spin, _ := cmd.Flags().GetInt64("spin")
			bonus, _ := cmd.Flags().GetStringSlice("bonus")
			snapshots, _ := cmd.Flags().GetStringSlice("snapshots")

			in := new(pb.PachiSlotMachine)
			in.No = int32(tid)

			m := new(pb.PachiSlotRoadMap)
			m.Total = total
			m.Spin = spin

			for _, row := range bonus {
				row = strings.TrimSpace(row)
				strs := strings.Split(row, ":")
				if len(strs) != 2 {
					continue
				}
				strs[0] = strings.TrimSpace(strs[0])
				strs[1] = strings.TrimSpace(strs[1])

				cnt, err := strconv.Atoi(strs[1])
				if err != nil {
					continue
				}

				val, ok := pb.PachiSlot_Scenario_value["PachiSlot_Scenario_"+strs[0]]
				if !ok {
					continue
				}

				b := new(pb.PachiSlotRoadMapBonus)
				b.Type = pb.PachiSlot_Scenario(val)
				b.Count = int64(cnt)

				var exist bool
				for _, row := range m.Bonus {
					if row.Type == b.Type {
						exist = true
						break
					}
				}

				if exist {
					continue
				}

				m.Bonus = append(m.Bonus, b)
			}

			now := time.Now()

			for i, row := range snapshots {
				row = strings.TrimSpace(row)
				strs := strings.Split(row, ":")
				if len(strs) != 2 {
					continue
				}
				strs[0] = strings.TrimSpace(strs[0])
				strs[1] = strings.TrimSpace(strs[1])

				cnt, err := strconv.Atoi(strs[1])
				if err != nil {
					continue
				}

				val, ok := pb.PachiSlot_Scenario_value["PachiSlot_Scenario_"+strs[0]]
				if !ok {
					continue
				}

				snapshot := new(pb.PachiSlotRoadMapSnapshot)
				snapshot.Type = pb.PachiSlot_Scenario(val)
				snapshot.Count = int64(cnt)
				snapshot.Time = now.Add(time.Duration(int64(-i)*5) * time.Minute).Unix()

				m.Snapshots = append(m.Snapshots, snapshot)
			}

			in.RoadMap = m

			if _, err := PachiSlotServiceClient.AssignRoadMap(cmd.Context(), in); err != nil {
				return err
			}

			cmd.Printf("機台[%d] 指定轉數表\n", in.No)
			cmd.Printf("總轉數: %d 次\n", in.RoadMap.Total)
			cmd.Printf("當前轉數: %d 次\n", in.RoadMap.Spin)
			for _, row := range in.RoadMap.Bonus {
				cmd.Printf("大獎[%v] %d 次\n", cliapisvr.ReversePachiSlotScenario(row.Type), row.Count)
			}
			for _, row := range in.RoadMap.Snapshots {
				cmd.Printf("快照[%v] %d 次 %s\n", cliapisvr.ReversePachiSlotScenario(row.Type), row.Count, time.Unix(row.Time, 0).Format(time.RFC3339))
			}

			return nil
		},
	}

	cmd.Flags().Int64("total", 0, "總轉數")
	cmd.Flags().Int64("spin", 0, "當前轉數")
	cmd.Flags().StringSlice("bonus", []string(nil), "Bonus次數. ex: BB:123,RB:323")
	cmd.Flags().StringSlice("snapshots", []string(nil), "歷史快照. ex: BB:134,BB:235,RB:12")

	return cmd
}

func pachiSlotMachineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "machine <machineID>",
		Short:                 "查看遊戲機台",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			addr, _ := cmd.Flags().GetString("addr")
			port, _ := cmd.Flags().GetInt("port")
			address := fmt.Sprintf("%s:%d", addr, port)

			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
			if err != nil {
				return fmt.Errorf("GRPC Dial Error to %s, %v", address, err)
			}
			PachiSlotServiceClient = pb.NewPachiSlotServiceClient(conn)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) == 0 {
				machines, err := PachiSlotServiceClient.MachineList(cmd.Context(), new(pb.Empty))
				if err != nil {
					return err
				}

				p, _ := cmd.Flags().GetInt32("filter-pattern")
				isgaming, _ := cmd.PersistentFlags().GetBool("filter-gaming")
				isfree, _ := cmd.PersistentFlags().GetBool("filter-free")

				var freecnt, gamecnt int
				for _, m := range machines.Machines {
					if m.Player == nil {
						freecnt++
					} else {
						gamecnt++
					}
				}

				cmd.Printf("=== 機台列表(%d) ===\n", len(machines.Machines))
				cmd.Printf("遊戲中[%d]\t空閒中[%d]\n", gamecnt, freecnt)
				for _, m := range machines.Machines {
					if p > 0 && m.Pattern != p {
						continue
					}
					if isgaming && m.Player == nil {
						continue
					}
					if isfree && m.Player != nil {
						continue
					}
					printMachineLite(cmd.OutOrStdout(), m, "")
				}
				return nil
			}

			tid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("機台輸入錯誤, %v", err)
			}

			cond := new(pb.QueryCondition)
			cond.Table = int32(tid)

			m, err := PachiSlotServiceClient.MachineInfo(cmd.Context(), cond)
			if err != nil {
				return err
			}

			printMachine(cmd.OutOrStdout(), m, "")
			return nil
		},
	}

	cmd.PersistentFlags().BoolP("detail", "d", false, "列出詳細資訊")
	cmd.PersistentFlags().Int32P("filter-pattern", "p", 0, "指定輸出Pattern")
	cmd.PersistentFlags().BoolP("filter-gaming", "g", false, "只列出遊戲中牌桌")
	cmd.PersistentFlags().BoolP("filter-free", "f", false, "只列出空閒中牌桌")

	return cmd
}

func printMachine(w io.Writer, m *pb.PachiSlotMachine, indent string) {

	fmt.Fprintf(w, "%s機台編號: %04d\n", indent, m.No)
	fmt.Fprintf(w, "%s硬幣花費: %4d\n", indent, m.BetOdds)
	fmt.Fprintf(w, "%sPattern: %5d\n", indent, m.Pattern)
	fmt.Fprintf(w, "%s====================================\n", indent)
	fmt.Fprintf(w, "%s情境[%v] 大獎[%v] 模式[%v] \n", indent, cliapisvr.ReversePachiSlotScenario(m.Scenario), cliapisvr.ReversePachiSlotScenario(m.Bonus), cliapisvr.ReversePachiSlotMode(m.Mode))
	fmt.Fprintf(w, "%s剩餘RT[%d] Replay[%t] Bonus確定[%t] 剩餘連莊[%d]\n", indent, m.RT, m.Replay, m.BonusSure, m.BonusStore)
	if m.Scenario == pb.PachiSlot_Scenario_PachiSlot_Scenario_BB || m.Scenario == pb.PachiSlot_Scenario_PachiSlot_Scenario_RB {
		fmt.Fprintf(w, "%s[%v] 剩餘次數 %d\n", indent, cliapisvr.ReversePachiSlotScenario(m.Scenario), m.BonusTimesLeft)
		if m.Scenario == pb.PachiSlot_Scenario_PachiSlot_Scenario_BB {
			fmt.Fprintf(w, "%s[%v] 贏得JAC次數 %d\n", indent, cliapisvr.ReversePachiSlotScenario(m.Scenario), m.JACTimes)
		}
	}
	if m.Scenario == pb.PachiSlot_Scenario_PachiSlot_Scenario_JAC {
		fmt.Fprintf(w, "%s[%v] 剩餘次數 %d\n", indent, cliapisvr.ReversePachiSlotScenario(m.Bonus), m.BonusTimesLeft)
		fmt.Fprintf(w, "%s[%v] 剩餘次數 %d\n", indent, cliapisvr.ReversePachiSlotScenario(m.Scenario), m.JACTimesLeft)
	}
	switch m.NextMode {
	case pb.PachiSlot_ModeType_PachiSlot_ModeType_A, pb.PachiSlot_ModeType_PachiSlot_ModeType_B, pb.PachiSlot_ModeType_PachiSlot_ModeType_Heaven:
		fmt.Fprintf(w, "%s指定下次模式[%v]\n", indent, cliapisvr.ReversePachiSlotMode(m.NextMode))
	}
	if m.NextResult != pb.PachiSlot_ResultType_PachiSlot_ResultType_Empty && m.NextResult != pb.PachiSlot_ResultType(-1) {
		fmt.Fprintf(w, "%s指定下次牌型[%v]\n", indent, cliapisvr.ReversePachiSlotResultType(m.NextResult))
	}
	fmt.Fprintf(w, "%s====================================\n", indent)
	fmt.Fprintf(w, "%s轉數表:\n", indent)
	fmt.Fprintf(w, "%s  總轉數: %d 轉\n", indent, m.RoadMap.Total)
	fmt.Fprintf(w, "%s  當前轉數: %d 轉\n", indent, m.RoadMap.Spin)
	for _, row := range m.RoadMap.Bonus {
		var t string
		switch row.Type {
		case pb.PachiSlot_Scenario_PachiSlot_Scenario_BB:
			t = "BB"
		case pb.PachiSlot_Scenario_PachiSlot_Scenario_RB:
			t = "RB"
		}
		fmt.Fprintf(w, "%s  [%s]: %d 次\n", indent, t, row.Count)
	}
	fmt.Fprintf(w, "%s====================================\n", indent)
	fmt.Fprintf(w, "%s玩家:\n", indent)
	if m.Player == nil {
		fmt.Fprintf(w, "%s  空閒中\n", indent)
	} else {
		printPlayerLite(w, m.Player, indent+"  ")
	}
}

func printMachineLite(w io.Writer, m *pb.PachiSlotMachine, indent string) {

	var gaming string = "空閒中"
	if m.Player != nil {
		gaming = "遊戲中"
	}

	// 總轉數[%d] 當前轉數[%d] BB[%d] RB[%d] AT[%d]

	var strs []string
	for _, row := range m.RoadMap.Bonus {
		strs = append(strs, fmt.Sprintf("%v[%d]", cliapisvr.ReversePachiSlotScenario(row.Type), row.Count))
	}
	sort.Strings(strs)

	fmt.Fprintf(w, "%s機台[%04d][%s]  硬幣花費: %4d  Pattern: %5d 總轉數[%d] 當前轉數[%d] %s\n", indent, m.No, gaming, m.BetOdds, m.Pattern, m.RoadMap.Total, m.RoadMap.Spin, strings.Join(strs, " "))
}
