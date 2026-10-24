package cliapisvr

import (
	"context"
	"fmt"
	"time"

	"example.invalid/dig"
	"example.invalid/grpc/codes"
	"example.invalid/grpc/status"

	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/pachislot"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
	pspb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// PachiSlotServiceServer PachiSlotService Server
type PachiSlotServiceServer struct{}

// MachineList 列出所有機台
func (s *PachiSlotServiceServer) MachineList(ctx context.Context, arg *pb.Empty) (machines *pb.PachiSlotMachineList, err error) {
	if app.Context == nil {
		return nil, status.Error(codes.Aborted, "app.Context not set")
	}
	if err := app.Context.Invoke(func(f pachislot.PachiGameFunc, r *gamelogic.Room) error {
		machines = new(pb.PachiSlotMachineList)
		for _, t := range r.Tables {
			game := f(t.GameLogic)
			m := TransPachiSlotMachine(t, game)
			machines.Machines = append(machines.Machines, m)
		}
		return nil
	}); err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}
	return
}

// MachineInfo 查詢機台資訊
func (s *PachiSlotServiceServer) MachineInfo(ctx context.Context, cond *pb.QueryCondition) (machine *pb.PachiSlotMachine, err error) {
	if app.Context == nil {
		return nil, status.Error(codes.Aborted, "app.Context not set")
	}
	if err := app.Context.Invoke(func(f pachislot.PachiGameFunc, r *gamelogic.Room) error {
		table := r.GetTableByNo(cond.Table)
		if table == nil {
			return status.Error(codes.NotFound, fmt.Sprintf("Machine %d not found", cond.Table))
		}
		game := f(table.GameLogic)
		machine = TransPachiSlotMachine(table, game)
		return nil
	}); err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}
	return
}

// AssignBonus 指定當前大獎
func (s *PachiSlotServiceServer) AssignBonus(ctx context.Context, req *pb.PachiSlotMachine) (*pb.Empty, error) {
	if app.Context == nil {
		return nil, status.Error(codes.InvalidArgument, "app.Context not set")
	}
	err := app.Context.Invoke(func(f pachislot.PachiGameFunc, r *gamelogic.Room) error {
		table := r.GetTableByNo(req.No)
		if table == nil {
			return status.Error(codes.NotFound, fmt.Sprintf("Machine %d not found", req.No))
		}

		game := f(table.GameLogic)

		if game.Scenario != pspb.Scenario_NG {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("Machine[%d] is not in NG, can't change Bonus", req.No))
		}

		// 檢查大獎是否在遊戲轉數表內
		b := ReversePachiSlotScenario(req.Bonus)
		if _, ok := game.RoadMap.FeaturesCnt[b]; !ok {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("Machine[%d] has no Bonus with [%v]", req.No, b))
		}

		game.Bonus = b
		example.invalid("機台[%d] 指定當前大獎[%v]", req.No, game.Bonus)

		return nil
	})
	if err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}

	return new(pb.Empty), nil
}

// AssignMode 指定下個模式
func (s *PachiSlotServiceServer) AssignMode(ctx context.Context, req *pb.PachiSlotMachine) (*pb.Empty, error) {
	if app.Context == nil {
		return nil, status.Error(codes.Aborted, "app.Context not set")
	}

	err := app.Context.Invoke(func(f pachislot.PachiGameFunc, r *gamelogic.Room) error {
		table := r.GetTableByNo(req.No)
		if table == nil {
			return status.Error(codes.NotFound, fmt.Sprintf("Machine %d not found", req.No))
		}

		game := f(table.GameLogic)

		game.NextMode = ReversePachiSlotMode(req.NextMode)
		example.invalid("機台[%d] 指定模式[%v]", req.No, game.NextMode)

		return nil
	})
	if err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}

	return new(pb.Empty), nil
}

// AssignResult 指定下個牌型
func (s *PachiSlotServiceServer) AssignResult(ctx context.Context, req *pb.PachiSlotMachine) (*pb.Empty, error) {
	if app.Context == nil {
		return nil, status.Error(codes.Aborted, "app.Context not set")
	}

	err := app.Context.Invoke(func(f pachislot.PachiGameFunc, r *gamelogic.Room) error {
		table := r.GetTableByNo(req.No)
		if table == nil {
			return status.Error(codes.NotFound, fmt.Sprintf("Machine %d not found", req.No))
		}

		game := f(table.GameLogic)

		game.NextResult = ReversePachiSlotResultType(req.NextResult)
		example.invalid("機台[%d] 指定牌型[%v]", req.No, game.NextResult)

		return nil
	})
	if err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}

	return new(pb.Empty), nil
}

// AssignReplay 指定為Replay
func (s *PachiSlotServiceServer) AssignReplay(ctx context.Context, req *pb.PachiSlotMachine) (*pb.Empty, error) {
	if app.Context == nil {
		return nil, status.Error(codes.InvalidArgument, "app.Context not set")
	}
	err := app.Context.Invoke(func(f pachislot.PachiGameFunc, r *gamelogic.Room) error {
		table := r.GetTableByNo(req.No)
		if table == nil {
			return status.Error(codes.NotFound, fmt.Sprintf("Machine %d not found", req.No))
		}

		game := f(table.GameLogic)

		if game.Scenario != pspb.Scenario_NG {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("Machine[%d] is not in NG, can't assign replay", req.No))
		}

		game.Replay = true
		example.invalid("機台[%d] 指定為 Replay", req.No)

		return nil
	})
	if err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}

	return new(pb.Empty), nil
}

// AssignBonusSure 指定為Bonus確定
func (s *PachiSlotServiceServer) AssignBonusSure(ctx context.Context, req *pb.PachiSlotMachine) (*pb.Empty, error) {
	if app.Context == nil {
		return nil, status.Error(codes.InvalidArgument, "app.Context not set")
	}
	err := app.Context.Invoke(func(f pachislot.PachiGameFunc, r *gamelogic.Room) error {
		table := r.GetTableByNo(req.No)
		if table == nil {
			return status.Error(codes.NotFound, fmt.Sprintf("Machine %d not found", req.No))
		}

		game := f(table.GameLogic)

		if game.Scenario != pspb.Scenario_NG {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("Machine[%d] is not in NG, can't assign Bonus Sure", req.No))
		}

		game.BonusSure = true
		example.invalid("機台[%d] 指定為 Bonus確定", req.No)

		return nil
	})
	if err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}

	return new(pb.Empty), nil
}

// AssignRT 指定RT數
func (s *PachiSlotServiceServer) AssignRT(ctx context.Context, req *pb.PachiSlotMachine) (*pb.Empty, error) {
	if app.Context == nil {
		return nil, status.Error(codes.InvalidArgument, "app.Context not set")
	}
	err := app.Context.Invoke(func(f pachislot.PachiGameFunc, r *gamelogic.Room) error {
		table := r.GetTableByNo(req.No)
		if table == nil {
			return status.Error(codes.NotFound, fmt.Sprintf("Machine %d not found", req.No))
		}

		game := f(table.GameLogic)

		if game.Scenario != pspb.Scenario_NG {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("Machine[%d] is not in NG, can't change RT", req.No))
		}

		game.RT = int(req.RT)
		example.invalid("機台[%d] 指定剩餘RT為 %d", req.No, game.RT)

		return nil
	})
	if err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}

	return new(pb.Empty), nil
}

// AssignBonusTimesLeft 指定Bonus剩餘次數
func (s *PachiSlotServiceServer) AssignBonusTimesLeft(ctx context.Context, req *pb.PachiSlotMachine) (*pb.Empty, error) {
	if app.Context == nil {
		return nil, status.Error(codes.InvalidArgument, "app.Context not set")
	}
	err := app.Context.Invoke(func(f pachislot.PachiGameFunc, r *gamelogic.Room) error {
		table := r.GetTableByNo(req.No)
		if table == nil {
			return status.Error(codes.NotFound, fmt.Sprintf("Machine %d not found", req.No))
		}

		game := f(table.GameLogic)

		if game.Scenario == pspb.Scenario_NG {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("Machine[%d] is in NG, can't change BonusTimesLeft", req.No))
		}

		game.BonusTimesLeft = int(req.BonusTimesLeft)
		example.invalid("機台[%d] 指定剩餘Bonus次數為 %d", req.No, game.BonusTimesLeft)

		return nil
	})
	if err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}

	return new(pb.Empty), nil
}

// AssignJACTimesLeft 指定JAC剩餘次數
func (s *PachiSlotServiceServer) AssignJACTimesLeft(ctx context.Context, req *pb.PachiSlotMachine) (*pb.Empty, error) {
	if app.Context == nil {
		return nil, status.Error(codes.InvalidArgument, "app.Context not set")
	}
	err := app.Context.Invoke(func(f pachislot.PachiGameFunc, r *gamelogic.Room) error {
		table := r.GetTableByNo(req.No)
		if table == nil {
			return status.Error(codes.NotFound, fmt.Sprintf("Machine %d not found", req.No))
		}

		game := f(table.GameLogic)

		if game.Scenario == pspb.Scenario_NG {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("Machine[%d] is in NG, can't change JACTimesLeft", req.No))
		}

		if game.Bonus == pspb.Scenario_BB {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("Machine[%d] Bonus is not BB, can't change JACTimesLeft", req.No))
		}

		game.JACTimesLeft = int(req.JACTimesLeft)
		example.invalid("機台[%d] 指定剩餘JAC次數為 %d", req.No, game.JACTimesLeft)

		return nil
	})
	if err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}

	return new(pb.Empty), nil
}

// AssignBonusStore 指定BonusStore剩餘次數
func (s *PachiSlotServiceServer) AssignBonusStore(ctx context.Context, req *pb.PachiSlotMachine) (*pb.Empty, error) {
	if app.Context == nil {
		return nil, status.Error(codes.InvalidArgument, "app.Context not set")
	}
	err := app.Context.Invoke(func(f pachislot.PachiGameFunc, r *gamelogic.Room) error {
		table := r.GetTableByNo(req.No)
		if table == nil {
			return status.Error(codes.NotFound, fmt.Sprintf("Machine %d not found", req.No))
		}

		game := f(table.GameLogic)

		game.BonusStore = int(req.BonusStore)
		example.invalid("機台[%d] 指定剩餘Bonus次數為 %d", req.No, game.BonusStore)

		return nil
	})
	if err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}

	return new(pb.Empty), nil
}

// AssignRoadMap 指定轉數表
func (s *PachiSlotServiceServer) AssignRoadMap(ctx context.Context, req *pb.PachiSlotMachine) (*pb.Empty, error) {
	if app.Context == nil {
		return nil, status.Error(codes.InvalidArgument, "app.Context not set")
	}
	err := app.Context.Invoke(func(f pachislot.PachiGameFunc, r *gamelogic.Room) error {
		table := r.GetTableByNo(req.No)
		if table == nil {
			return status.Error(codes.NotFound, fmt.Sprintf("Machine %d not found", req.No))
		}

		game := f(table.GameLogic)

		game.RoadMap = ReversePachiSlotRoadMap(req.RoadMap)
		example.invalid("機台[%d] 指定轉數表 %v", req.No, game.RoadMap)

		return nil
	})
	if err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}

	return new(pb.Empty), nil
}

// QueryProto 查詢Proto設定
func (s *PachiSlotServiceServer) QueryProto(ctx context.Context, req *pb.PachiSlotQueryProtoCondition) (*pb.PachiSlotQueryProtoResponse, error) {
	if app.Context == nil {
		return nil, status.Error(codes.InvalidArgument, "app.Context not set")
	}

	res := new(pb.PachiSlotQueryProtoResponse)

	err := app.Context.Invoke(func(proto *pachislot.Proto) error {

		if req.ShowGameConfig {
			res.GameConfig = TransPachiSlotProtoGameConfig(proto.GameConfig)
		}
		if req.ShowPayLine {
			for _, row := range proto.PayLine {
				res.PayLine = append(res.PayLine, TransPachiSlotProtoPayLine(row))
			}
		}
		if req.ShowPayTable {
			for _, row := range proto.PayTable {
				res.PayTable = append(res.PayTable, TransPachiSlotProtoPayoutOption(row))
			}
		}
		if req.ShowReelStrips {
			for _, row := range proto.ReelStrips {
				res.ReelStrips = append(res.ReelStrips, TransPachiSlotProtoReelStripSymbol(row))
			}
		}
		if req.ShowPattern {
			for _, row := range proto.Pattern {
				res.Pattern = append(res.Pattern, TransPachiSlotProtoPattern(row))
			}
		}
		if req.ShowModule {
			for _, row := range proto.ModuleConfigs {
				res.Module = append(res.Module, TransPachiSlotProtoModule(row))
			}
		}
		if req.ShowMode {
			for _, row := range proto.ModeConfigs {
				res.Mode = append(res.Mode, TransPachiSlotProtoMode(row))
			}
		}
		if req.ShowRTRelease {
			for _, row := range proto.RTReleaseDistribution {
				res.RTRelease = append(res.RTRelease, TransPachiSlotProtoRTReleaseOption(row))
			}
		}
		if req.ShowResultDistribution {
			for _, row := range proto.ResultDistribution {
				res.ResultDistribution = append(res.ResultDistribution, TransPachiSlotProtoResultDistributionOption(row))
			}
		}
		if req.ShowNextModeDistribution {
			for _, row := range proto.NextModeDistribution {
				res.NextModeDistribution = append(res.NextModeDistribution, TransPachiSlotProtoNextModeDistributionOption(row))
			}
		}
		if req.ShowRTDistribution {
			for _, row := range proto.RTDistribution {
				res.RTDistribution = append(res.RTDistribution, TransPachiSlotProtoRTDistributionOption(row))
			}
		}
		if req.ShowBRDistribution {
			for _, row := range proto.BRDistribution {
				res.BRDistribution = append(res.BRDistribution, TransPachiSlotProtoBRDistributionOption(row))
			}
		}

		return nil
	})
	if err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}

	return res, nil
}

// TransPachiSlotProtoGameConfig 轉換 Proto GameConfig
func TransPachiSlotProtoGameConfig(src *pspb.GameConfig) *pb.PachiSlot_Proto_GameConfig {
	if src == nil {
		return nil
	}
	dist := new(pb.PachiSlot_Proto_GameConfig)
	dist.ID = src.ID
	dist.Name = src.Name
	dist.PayLineGroupID = src.PayLineGroupID
	dist.PayTableGroupID = src.PayTableGroupID
	dist.PatternGroupID = src.PatternGroupID
	dist.ReelStripsGroupID = src.ReelStripsGroupID
	return dist
}

// ReversePachiSlotProtoGameConfig 轉換 Proto GameConfig
func ReversePachiSlotProtoGameConfig(src *pb.PachiSlot_Proto_GameConfig) *pspb.GameConfig {
	if src == nil {
		return nil
	}
	dist := new(pspb.GameConfig)
	dist.ID = src.ID
	dist.Name = src.Name
	dist.PayLineGroupID = src.PayLineGroupID
	dist.PayTableGroupID = src.PayTableGroupID
	dist.PatternGroupID = src.PatternGroupID
	dist.ReelStripsGroupID = src.ReelStripsGroupID
	return dist
}

// TransPachiSlotProtoPayLine 轉換 Proto PayLine
func TransPachiSlotProtoPayLine(src *pspb.PayLine) *pb.PachiSlot_Proto_PayLine {
	if src == nil {
		return nil
	}
	dist := new(pb.PachiSlot_Proto_PayLine)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.LineID = src.LineID
	for _, pos := range src.Positions {
		dist.Positions = append(dist.Positions, TransPachiSlotProtoPosition(pos))
	}
	return dist
}

// ReversePachiSlotProtoPayLine 轉換 Proto PayLine
func ReversePachiSlotProtoPayLine(src *pb.PachiSlot_Proto_PayLine) *pspb.PayLine {
	if src == nil {
		return nil
	}
	dist := new(pspb.PayLine)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.LineID = src.LineID
	for _, pos := range src.Positions {
		dist.Positions = append(dist.Positions, ReversPachiSlotProtoPosition(pos))
	}
	return dist
}

// TransPachiSlotProtoPosition 轉換 Proto Position
func TransPachiSlotProtoPosition(src *pspb.Position) *pb.PachiSlot_Proto_Position {
	if src == nil {
		return nil
	}
	dist := new(pb.PachiSlot_Proto_Position)
	dist.Reel = src.Reel
	dist.Row = src.Row
	return dist
}

// ReversPachiSlotProtoPosition 轉換 Proto Position
func ReversPachiSlotProtoPosition(src *pb.PachiSlot_Proto_Position) *pspb.Position {
	if src == nil {
		return nil
	}
	dist := new(pspb.Position)
	dist.Reel = src.Reel
	dist.Row = src.Row
	return dist
}

// TransPachiSlotProtoPayoutOption 轉換 Proto PayoutOption
func TransPachiSlotProtoPayoutOption(src *pspb.PayoutOption) *pb.PachiSlot_Proto_PayoutOption {
	if src == nil {
		return nil
	}
	dist := new(pb.PachiSlot_Proto_PayoutOption)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.ResultType = TransPachiSlotResultType(src.ResultType)
	dist.PayType = TransPachiSlotPayType(src.PayType)
	dist.Payout = src.Payout
	dist.UseType = TransPachiSlotScenario(src.UseType)
	return dist
}

// ReversePachiSlotProtoPayoutOption 轉換 Proto PayoutOption
func ReversePachiSlotProtoPayoutOption(src *pb.PachiSlot_Proto_PayoutOption) *pspb.PayoutOption {
	if src == nil {
		return nil
	}
	dist := new(pspb.PayoutOption)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.ResultType = ReversePachiSlotResultType(src.ResultType)
	dist.PayType = ReversePachiSlotPayType(src.PayType)
	dist.Payout = src.Payout
	dist.UseType = ReversePachiSlotScenario(src.UseType)
	return dist
}

// TransPachiSlotProtoReelStripSymbol 轉換 Proto ReelStripSymbol
func TransPachiSlotProtoReelStripSymbol(src *pspb.ReelStripSymbol) *pb.PachiSlot_Proto_ReelStripSymbol {
	if src == nil {
		return nil
	}
	dist := new(pb.PachiSlot_Proto_ReelStripSymbol)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.Reel = src.Reel
	dist.Row = src.Row
	dist.Symbol = TransPachiSlotSymbol(src.Symbol)
	return dist
}

// ReversePachiSlotProtoReelStripSymbol 轉換 Proto ReelStripSymbol
func ReversePachiSlotProtoReelStripSymbol(src *pb.PachiSlot_Proto_ReelStripSymbol) *pspb.ReelStripSymbol {
	if src == nil {
		return nil
	}
	dist := new(pspb.ReelStripSymbol)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.Reel = src.Reel
	dist.Row = src.Row
	dist.Symbol = ReversePachiSlotSymbol(src.Symbol)
	return dist
}

// TransPachiSlotProtoPattern 轉換 Proto PatternOption
func TransPachiSlotProtoPattern(src *pspb.PatternOption) *pb.PachiSlot_Proto_PatternOption {
	if src == nil {
		return nil
	}
	dist := new(pb.PachiSlot_Proto_PatternOption)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.ModeGroupID = src.ModeGroupID
	dist.NGResultGroupID = src.NGResultGroupID
	dist.ModuleGroupID = src.ModuleGroupID
	dist.RTReleaseGroupID = src.RTReleaseGroupID
	return dist
}

// ReversePachiSlotProtoPattern 轉換 Proto PatternOption
func ReversePachiSlotProtoPattern(src *pb.PachiSlot_Proto_PatternOption) *pspb.PatternOption {
	if src == nil {
		return nil
	}
	dist := new(pspb.PatternOption)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.ModeGroupID = src.ModeGroupID
	dist.NGResultGroupID = src.NGResultGroupID
	dist.ModuleGroupID = src.ModuleGroupID
	dist.RTReleaseGroupID = src.RTReleaseGroupID
	return dist
}

// TransPachiSlotProtoModule 轉換 Proto Module
func TransPachiSlotProtoModule(src *pspb.Module) *pb.PachiSlot_Proto_Module {
	if src == nil {
		return nil
	}
	dist := new(pb.PachiSlot_Proto_Module)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.ResultProbGroupID = src.ResultProbGroupID
	dist.Type = TransPachiSlotScenario(src.Type)
	dist.Count = src.Count
	dist.Cost = src.Cost
	return dist
}

// ReversePachiSlotProtoModule 轉換 Proto Module
func ReversePachiSlotProtoModule(src *pb.PachiSlot_Proto_Module) *pspb.Module {
	if src == nil {
		return nil
	}
	dist := new(pspb.Module)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.ResultProbGroupID = src.ResultProbGroupID
	dist.Type = ReversePachiSlotScenario(src.Type)
	dist.Count = src.Count
	dist.Cost = src.Cost
	return dist
}

// TransPachiSlotProtoMode 轉換 Proto Mode
func TransPachiSlotProtoMode(src *pspb.Mode) *pb.PachiSlot_Proto_Mode {
	if src == nil {
		return nil
	}
	dist := new(pb.PachiSlot_Proto_Mode)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.ModeType = TransPachiSlotMode(src.ModeType)
	dist.NextModeGroupID = src.NextModeGroupID
	dist.RTGroupID = src.RTGroupID
	dist.BRGroupID = src.BRGroupID
	return dist
}

// ReversePachiSlotProtoMode 轉換 Proto Mode
func ReversePachiSlotProtoMode(src *pb.PachiSlot_Proto_Mode) *pspb.Mode {
	if src == nil {
		return nil
	}
	dist := new(pspb.Mode)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.ModeType = ReversePachiSlotMode(src.ModeType)
	dist.NextModeGroupID = src.NextModeGroupID
	dist.RTGroupID = src.RTGroupID
	dist.BRGroupID = src.BRGroupID
	return dist
}

// TransPachiSlotProtoRTReleaseOption 轉換 Proto RTReleaseOption
func TransPachiSlotProtoRTReleaseOption(src *pspb.RTReleaseOption) *pb.PachiSlot_Proto_RTReleaseOption {
	if src == nil {
		return nil
	}
	dist := new(pb.PachiSlot_Proto_RTReleaseOption)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.ResultType = TransPachiSlotResultType(src.ResultType)
	dist.Rate = src.Rate
	dist.Denominator = src.Denominator
	dist.Times = src.Times
	dist.UseType = TransPachiSlotScenario(src.UseType)
	return dist
}

// ReversePachiSlotProtoRTReleaseOption 轉換 Proto RTReleaseOption
func ReversePachiSlotProtoRTReleaseOption(src *pb.PachiSlot_Proto_RTReleaseOption) *pspb.RTReleaseOption {
	if src == nil {
		return nil
	}
	dist := new(pspb.RTReleaseOption)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.ResultType = ReversePachiSlotResultType(src.ResultType)
	dist.Rate = src.Rate
	dist.Denominator = src.Denominator
	dist.Times = src.Times
	dist.UseType = ReversePachiSlotScenario(src.UseType)
	return dist
}

// TransPachiSlotProtoResultDistributionOption 轉換 Proto ResultDistributionOption
func TransPachiSlotProtoResultDistributionOption(src *pspb.ResultDistributionOption) *pb.PachiSlot_Proto_ResultDistributionOption {
	if src == nil {
		return nil
	}
	dist := new(pb.PachiSlot_Proto_ResultDistributionOption)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.Result = TransPachiSlotResultType(src.Result)
	dist.Weight = src.Weight
	return dist
}

// ReversePachiSlotProtoResultDistributionOption 轉換 Proto ResultDistributionOption
func ReversePachiSlotProtoResultDistributionOption(src *pb.PachiSlot_Proto_ResultDistributionOption) *pspb.ResultDistributionOption {
	if src == nil {
		return nil
	}
	dist := new(pspb.ResultDistributionOption)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.Result = ReversePachiSlotResultType(src.Result)
	dist.Weight = src.Weight
	return dist
}

// TransPachiSlotProtoNextModeDistributionOption 轉換 Proto NextModeDistributionOption
func TransPachiSlotProtoNextModeDistributionOption(src *pspb.NextModeDistributionOption) *pb.PachiSlot_Proto_NextModeDistributionOption {
	if src == nil {
		return nil
	}
	dist := new(pb.PachiSlot_Proto_NextModeDistributionOption)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.Mode = TransPachiSlotMode(src.Mode)
	dist.Weight = src.Weight
	return dist
}

// ReversePachiSlotProtoNextModeDistributionOption 轉換 Proto NextModeDistributionOption
func ReversePachiSlotProtoNextModeDistributionOption(src *pb.PachiSlot_Proto_NextModeDistributionOption) *pspb.NextModeDistributionOption {
	if src == nil {
		return nil
	}
	dist := new(pspb.NextModeDistributionOption)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.Mode = ReversePachiSlotMode(src.Mode)
	dist.Weight = src.Weight
	return dist
}

// TransPachiSlotProtoRTDistributionOption 轉換 Proto RTDistributionOption
func TransPachiSlotProtoRTDistributionOption(src *pspb.RTDistributionOption) *pb.PachiSlot_Proto_RTDistributionOption {
	if src == nil {
		return nil
	}
	dist := new(pb.PachiSlot_Proto_RTDistributionOption)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.Min = src.Min
	dist.Max = src.Max
	dist.Weight = src.Weight
	return dist
}

// ReversePachiSlotProtoRTDistributionOption 轉換 Proto RTDistributionOption
func ReversePachiSlotProtoRTDistributionOption(src *pb.PachiSlot_Proto_RTDistributionOption) *pspb.RTDistributionOption {
	if src == nil {
		return nil
	}
	dist := new(pspb.RTDistributionOption)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.Min = src.Min
	dist.Max = src.Max
	dist.Weight = src.Weight
	return dist
}

// TransPachiSlotProtoBRDistributionOption 轉換 Proto RTDistributionOption
func TransPachiSlotProtoBRDistributionOption(src *pspb.BRDistributionOption) *pb.PachiSlot_Proto_BRDistributionOption {
	if src == nil {
		return nil
	}
	dist := new(pb.PachiSlot_Proto_BRDistributionOption)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.Type = TransPachiSlotScenario(src.Type)
	dist.Weight = src.Weight
	return dist
}

// ReversePachiSlotProtoBRDistributionOption 轉換 Proto BRDistributionOption
func ReversePachiSlotProtoBRDistributionOption(src *pb.PachiSlot_Proto_BRDistributionOption) *pspb.BRDistributionOption {
	if src == nil {
		return nil
	}
	dist := new(pspb.BRDistributionOption)
	dist.ID = src.ID
	dist.GroupID = src.GroupID
	dist.Type = ReversePachiSlotScenario(src.Type)
	dist.Weight = src.Weight
	return dist
}

// TransPachiSlotMachine 轉換PachiSlot機台資料
func TransPachiSlotMachine(tb *gamelogic.Table, g *pachislot.Game) *pb.PachiSlotMachine {
	m := new(pb.PachiSlotMachine)

	m.No = tb.TableNo
	m.BetOdds = int32(g.CostPerCoin)
	m.Pattern = g.Pattern.ID
	m.Scenario = TransPachiSlotScenario(g.Scenario)
	m.Bonus = TransPachiSlotScenario(g.Bonus)
	m.Mode = TransPachiSlotMode(g.Mode)
	m.RT = int32(g.RT)
	m.Replay = g.Replay
	m.BonusSure = g.BonusSure
	m.JACTimes = int32(g.JACTimes)
	m.BonusTimesLeft = int32(g.BonusTimesLeft)
	m.JACTimesLeft = int32(g.JACTimesLeft)
	m.NextMode = TransPachiSlotMode(g.NextMode)
	m.BonusStore = int32(g.BonusStore)
	m.NextResult = TransPachiSlotResultType(g.NextResult)
	m.RoadMap = TransPachiSlotRoadMap(g.RoadMap)
	for _, u := range tb.GetTableUser() {
		m.Player = TransGameUser(u, nil)
	}
	return m
}

// TransPachiSlotRoadMap 轉換PachiSlot 轉數表資料
func TransPachiSlotRoadMap(src *pachislot.RoadMap) *pb.PachiSlotRoadMap {
	if src == nil {
		return nil
	}
	dist := new(pb.PachiSlotRoadMap)

	dist.Total = int64(src.TotalCnt)
	dist.Spin = int64(src.SpinCnt)
	for s, cnt := range src.FeaturesCnt {
		b := new(pb.PachiSlotRoadMapBonus)
		b.Count = int64(cnt)
		b.Type = TransPachiSlotScenario(s)
		dist.Bonus = append(dist.Bonus, b)
	}
	for _, row := range src.Snapshots {
		if row == nil {
			continue
		}
		shot := new(pb.PachiSlotRoadMapSnapshot)
		shot.Count = int64(row.Count)
		shot.Type = TransPachiSlotScenario(row.Feature)
		shot.Time = row.Time.Unix()
	}

	return dist
}

// ReversePachiSlotRoadMap 轉換PachiSlot 轉數表資料
func ReversePachiSlotRoadMap(src *pb.PachiSlotRoadMap) *pachislot.RoadMap {
	dist := pachislot.NewRoadMap(pachislot.DefaultRoadMapBuf)

	dist.TotalCnt = int(src.Total)
	dist.SpinCnt = int(src.Spin)
	for _, row := range src.Bonus {
		dist.FeaturesCnt[ReversePachiSlotScenario(row.Type)] = int(row.Count)
	}
	for i, row := range src.Snapshots {
		snapshot := &pachislot.RoadMapSnapshot{
			Feature: ReversePachiSlotScenario(row.Type),
			Count:   int(row.Count),
			Time:    time.Unix(row.Time, 0),
		}
		dist.Snapshots[i] = snapshot
	}

	return dist
}

// TransPachiSlotScenario 轉換 PachiSlot Scenario
func TransPachiSlotScenario(src pspb.Scenario) pb.PachiSlot_Scenario {
	return pb.PachiSlot_Scenario(src)
}

// ReversePachiSlotScenario 轉換 PachiSlot Scenario
func ReversePachiSlotScenario(src pb.PachiSlot_Scenario) pspb.Scenario {
	return pspb.Scenario(src)
}

// TransPachiSlotMode 轉換 PachiSlot Mode
func TransPachiSlotMode(src pspb.ModeType) pb.PachiSlot_ModeType {
	return pb.PachiSlot_ModeType(src)
}

// ReversePachiSlotMode 轉換 PachiSlot Mode
func ReversePachiSlotMode(src pb.PachiSlot_ModeType) pspb.ModeType {
	return pspb.ModeType(src)
}

// TransPachiSlotResultType 轉換 PachiSlot ResultType
func TransPachiSlotResultType(src pspb.ResultType) pb.PachiSlot_ResultType {
	return pb.PachiSlot_ResultType(src)
}

// ReversePachiSlotResultType 轉換 PachiSlot ResultType
func ReversePachiSlotResultType(src pb.PachiSlot_ResultType) pspb.ResultType {
	return pspb.ResultType(src)
}

// TransPachiSlotPayType 轉換 PachiSlot PayType
func TransPachiSlotPayType(src pspb.PayType) pb.PachiSlot_PayType {
	return pb.PachiSlot_PayType(src)
}

// ReversePachiSlotPayType 轉換 PachiSlot PayType
func ReversePachiSlotPayType(src pb.PachiSlot_PayType) pspb.PayType {
	return pspb.PayType(src)
}

// TransPachiSlotSymbol 轉換 PachiSlot Symbol
func TransPachiSlotSymbol(src pspb.Symbol) pb.PachiSlot_Symbol {
	return pb.PachiSlot_Symbol(src)
}

// ReversePachiSlotSymbol 轉換 PachiSlot Symbol
func ReversePachiSlotSymbol(src pb.PachiSlot_Symbol) pspb.Symbol {
	return pspb.Symbol(src)
}
