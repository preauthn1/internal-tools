package cliapisvr

import (
	"context"
	"fmt"

	"example.invalid/dig"
	"example.invalid/grpc/codes"
	"example.invalid/grpc/status"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"

	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/slot/helper"
	ppb "example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// SlotServiceServer Slot Service Server
type SlotServiceServer struct{}

// Info 取得 房間邏輯資訊
func (s *SlotServiceServer) Info(ctx context.Context, cond *pb.QueryCondition) (*pb.SlotInfo, error) {

	if helper.SlotHelperCmdP == nil {
		return nil, status.Error(codes.Aborted, "此遊戲不是老虎機")
	}

	g := helper.SlotHelperCmdP.GameLogic(cond.Table)

	if g == nil {
		return nil, status.Error(codes.NotFound, "查無機台")
	}

	game := helper.SlotHelperCmdP.BaseLogicFunc(g)

	info := new(pb.SlotInfo)
	info.SlotPattern = int32(game.SlotPatternID)
	info.SlotPatternComment = game.Proto.SlotPattern[0].GetMyComment()
	for _, r := range game.BetCounts {
		info.BetCounts = append(info.BetCounts, int64(r))
	}
	for _, r := range game.BetRates {
		info.BetRates = append(info.BetRates, int64(r))
	}
	info.MaxBetCount = int64(game.MaxBetCount)
	info.MaxBetRate = int64(game.MaxBetRate)
	info.MinCarryMoney = game.MinCarryMoney
	for key, val := range game.UnitPrefix {
		info.UnitPrefixs = append(info.UnitPrefixs, &pb.SlotUnitPrefix{Key: key, Value: val})
	}
	for _, desc := range game.Descriptions {
		info.Descriptions = append(info.Descriptions, &pb.SlotDescription{Name: desc.Name, Title: desc.Title, Description: desc.Description})
	}
	info.PaylineMode = game.PayLineMode.String()
	for _, jackpot := range game.Jackpots {
		j := new(pb.Jackpot)
		j.ID = jackpot.Origin.GetId()
		if jackpot.Origin.GetIsProgress() == 1 {
			j.IsProgress = true
			j.Increment = int64(jackpot.Origin.GetIncrement())
		}
		if jackpot.Origin.GetIsLinked() == 1 {
			j.IsLink = true
		}
		j.Level = jackpot.Level.String()
		j.Startup = jackpot.Origin.GetStartUp()
		j.Contribution = jackpot.TotalContribution
		j.Award = jackpot.AwardTotal

		info.Jackpots = append(info.Jackpots, j)
	}

	return info, nil
}

// QueryProto 查詢 Proto
func (s *SlotServiceServer) QueryProto(ctx context.Context, req *pb.SlotProtoQueryCondition) (*pb.SlotProtoQueryResponse, error) {
	if app.Context == nil {
		return nil, status.Error(codes.Aborted, "app.Context not set")
	}

	res := new(pb.SlotProtoQueryResponse)

	err := app.Context.Invoke(func(proto *ppb.ProtoData) error {
		if req.ShowSlotPattern {
			for _, row := range proto.SlotPattern {
				res.Pattern = append(res.Pattern, TransSlotProtoSlotPattern(row))
			}
		}
		if req.ShowSlotPayLine {
			for _, row := range proto.SlotPayLine {
				res.PayLine = append(res.PayLine, TransSlotProtoSlotPayLine(row))
			}
		}
		if req.ShowSlotPayTable {
			for _, row := range proto.SlotPayTable {
				res.PayTable = append(res.PayTable, TransSlotProtoSlotPayTable(row))
			}
		}
		if req.ShowMainGameTable {
			for _, row := range proto.MainGameTable {
				res.MainGameTable = append(res.MainGameTable, TransSlotProtoMainGameTable(row))
			}
		}
		if req.ShowSymbolAttributes {
			for _, row := range proto.SymbolAttributes {
				res.SymbolAttributes = append(res.SymbolAttributes, TransSlotProtoSymbolAttributes(row))
			}
		}
		if req.ShowJackpotTable {
			for _, row := range proto.JackpotTable {
				res.JackpotTable = append(res.JackpotTable, TransSlotProtoJackpotTable(row))
			}
		}
		if req.ShowFreeSpinTable {
			for _, row := range proto.FreeSpinTable {
				res.FreeSpinTable = append(res.FreeSpinTable, TransSlotProtoFreeSpinTable(row))
			}
		}
		if req.ShowSlotSymbolDistribution {
			for _, row := range proto.SlotSymbolDistribution {
				res.SlotSymbolDistribution = append(res.SlotSymbolDistribution, TransSlotProtoSlotSymbolDistribution(row))
			}
		}
		if req.ShowBonusTable {
			for _, row := range proto.BonusTable {
				res.BonusTable = append(res.BonusTable, TransSlotProtoBonusTable(row))
			}
		}
		if req.ShowPickGameTable {
			for _, row := range proto.PickGameTable {
				res.PickGameTable = append(res.PickGameTable, TransSlotProtoPickGameTable(row))
			}
		}
		if req.ShowSelectionTable {
			for _, row := range proto.SelectionTable {
				res.SelectionTable = append(res.SelectionTable, TransSlotProtoSelectionTable(row))
			}
		}
		if req.ShowRespinTable {
			for _, row := range proto.RespinTable {
				res.RespinTable = append(res.RespinTable, TransSlotProtoRespinTable(row))
			}
		}
		if req.ShowReelBGDistribution {
			for _, row := range proto.ReelBgDistribution {
				res.ReelBGDistribution = append(res.ReelBGDistribution, TransSlotProtoReelBgDistribution(row))
			}
		}
		if req.ShowExtraRule {
			for _, row := range proto.ExtraRule {
				res.ExtraRule = append(res.ExtraRule, TransSlotProtoExtraRule(row))
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

// TransSlotProtoSlotPattern 轉換 SlotPattern
func TransSlotProtoSlotPattern(src *ppb.SlotPattern) *pb.Slot_Proto_SlotPattern {
	dist := new(pb.Slot_Proto_SlotPattern)
	dist.ID = int32(src.GetId())
	dist.SymbolAttributesGroupID = int32(src.GetSymbolAttributesGroupId())
	dist.MainGameTableGroupID = int32(src.GetMainGameTableGroupId())
	dist.PayTableGroupID = int32(src.GetPayTableGroupId())
	dist.LineType = src.GetLineType()
	dist.PayLineGroupID = int32(src.GetPayLineGroupId())
	dist.ExtraRuleGroup = int32(src.GetExtraRuleGroup())
	dist.JackpotGroup = int32(src.GetJackpotGroup())
	dist.RTP = src.GetRtp()
	dist.Comment = src.GetMyComment()
	return dist
}

// TransSlotProtoSlotPayLine 轉換 SlotPayLine
func TransSlotProtoSlotPayLine(src *ppb.SlotPayLine) *pb.Slot_Proto_SlotPayLine {
	dist := new(pb.Slot_Proto_SlotPayLine)
	dist.ID = int32(src.GetId())
	dist.GroupID = int32(src.GetGroupId())
	dist.LineID = int32(src.GetLineId())
	dist.Content = src.GetContent()
	return dist
}

// TransSlotProtoSlotPayTable 轉換 SlotPayTable
func TransSlotProtoSlotPayTable(src *ppb.SlotPayTable) *pb.Slot_Proto_SlotPayTable {
	dist := new(pb.Slot_Proto_SlotPayTable)
	dist.ID = int32(src.GetId())
	dist.GroupID = int32(src.GetGroupId())
	dist.LinkType = src.GetLinkType()
	dist.Symbol = src.GetSymbol()
	dist.PayType = src.GetPayType()
	dist.Count = int32(src.GetCount())
	dist.Payout = src.GetPayout()
	dist.UseType = src.GetUseType()
	return dist
}

// TransSlotProtoMainGameTable 轉換 MainGameTable
func TransSlotProtoMainGameTable(src *ppb.MainGameTable) *pb.Slot_Proto_MainGameTable {
	dist := new(pb.Slot_Proto_MainGameTable)
	dist.ID = int32(src.GetId())
	dist.GroupID = int32(src.GetGroupId())
	dist.SymbolGroupID = int32(src.GetSymbolGroupId())
	dist.ReelNum = int32(src.GetReelNum())
	dist.RowNum = int32(src.GetRowNum())
	dist.MultiplierType = src.GetMultiplierType()
	dist.MultiplierValue = src.GetMultiplierValue()
	dist.IsGoldsymbol = int32(src.GetIsGoldsymbol())
	dist.GoldCondition = src.GetGoldCondition()
	dist.ConditionValue = src.GetConditionValue()
	dist.IsGoldReel = int32(src.GetIsGoldReel())
	dist.GoldReelGroupId = int32(src.GetGoldReelGroupId())
	dist.InvalidPosition = src.GetInvalidPosition()
	return dist
}

// TransSlotProtoSymbolAttributes 轉換 SymbolAttributes
func TransSlotProtoSymbolAttributes(src *ppb.SymbolAttributes) *pb.Slot_Proto_SymbolAttributes {
	dist := new(pb.Slot_Proto_SymbolAttributes)
	dist.ID = int32(src.GetId())
	dist.GroupID = int32(src.GetGroupId())
	dist.Symbol = src.GetSymbol()
	dist.IsGoldsymbol = int32(src.GetIsGoldsymbol())
	dist.GoldCondition = src.GetGoldCondition()
	dist.ConditionValue = src.GetConditionValue()
	dist.ReplaceGoldsymbol = src.GetReplaceGoldsymbol()
	dist.IsWild = int32(src.GetIsWild())
	dist.WildEx = src.GetWildEx()
	dist.Multiplier = int32(src.GetMultiplier())
	dist.IsExpand = int32(src.GetIsExpand())
	dist.IsSticky = int32(src.GetIsSticky())
	dist.StickyType = src.GetStickyType()
	dist.IsGoldBall = int32(src.GetIsGoldBall())
	dist.GoldballType = src.GetGoldballType()
	dist.GoldballValue = src.GetGoldballValue()
	dist.IsPluri = int32(src.GetIsPluri())
	dist.PluriCount = int32(src.GetPluriCount())
	dist.SubstituteSymbol = src.GetSubstituteSymbol()
	return dist
}

// TransSlotProtoJackpotTable 轉換 JackpotTable
func TransSlotProtoJackpotTable(src *ppb.JackpotTable) *pb.Slot_Proto_JackpotTable {
	dist := new(pb.Slot_Proto_JackpotTable)
	dist.ID = int32(src.GetId())
	dist.GroupID = int32(src.GetGroupId())
	dist.IsProgress = int32(src.GetIsProgress())
	dist.IsLinked = int32(src.GetIsLinked())
	dist.Increment = int32(src.GetIncrement())
	dist.PerBetProb = int32(src.GetPerBetProb())
	dist.StartUp = int64(src.GetStartUp())
	dist.IsMultiLevel = int32(src.GetIsMultiLevel())
	dist.Level = src.GetLevel()
	dist.IsCeling = int32(src.GetIsCeling())
	dist.CelingValue = int64(src.GetCelingValue())
	dist.PerProb = int32(src.GetPerProb())
	return dist
}

// TransSlotProtoFreeSpinTable 轉換 FreeSpinTable
func TransSlotProtoFreeSpinTable(src *ppb.FreeSpinTable) *pb.Slot_Proto_FreeSpinTable {
	dist := new(pb.Slot_Proto_FreeSpinTable)
	dist.ID = int32(src.GetId())
	dist.GroupID = int32(src.GetGroupId())
	dist.FreeSpinCount = int32(src.GetFreeSpinCount())
	dist.FreeSymbolGroupID = int32(src.GetFreeSymbolGroupId())
	dist.FreeReelNum = int32(src.GetFreeReelNum())
	dist.FreeRowNum = int32(src.GetFreeRowNum())
	dist.MultiplierType = src.GetMultiplierType()
	dist.MultiplierValue = src.GetMultiplierValue()
	dist.IsGoldsymbol = int32(src.GetIsGoldsymbol())
	dist.GoldCondition = src.GetGoldCondition()
	dist.ConditionValue = src.GetConditionValue()
	dist.IsInfo = int32(src.GetIsInfo())
	example.invalid = src.GetInfo()
	dist.IsGoldReel = int32(src.GetIsGoldReel())
	dist.GoldReelGroupID = int32(src.GetGoldReelGroupId())
	dist.InvalidPosition = src.GetInvalidPosition()
	return dist
}

// TransSlotProtoSlotSymbolDistribution 轉換 SlotSymbolDistribution
func TransSlotProtoSlotSymbolDistribution(src *ppb.SlotSymbolDistribution) *pb.Slot_Proto_SlotSymbolDistribution {
	dist := new(pb.Slot_Proto_SlotSymbolDistribution)
	dist.ID = int32(src.GetId())
	dist.GroupID = int32(src.GetGroupId())
	dist.ReelNum = int32(src.GetReelNum())
	dist.Symbol = src.GetSymbol()
	dist.Proportion = int32(src.GetProportion())
	dist.Order = int32(src.GetOrder())
	return dist
}

// TransSlotProtoBonusTable 轉換 BonusTable
func TransSlotProtoBonusTable(src *ppb.BonusTable) *pb.Slot_Proto_BonusTable {
	dist := new(pb.Slot_Proto_BonusTable)
	dist.ID = int32(src.GetId())
	dist.GroupID = int32(src.GetGroupId())
	dist.RewardType = src.GetRewardType()
	dist.RewardValue = src.GetRewardValue()
	dist.PerProb = int32(src.GetPerProb())
	dist.Condition = src.GetCondition()
	dist.ConditionValue = src.GetConditionValue()
	return dist
}

// TransSlotProtoPickGameTable 轉換 PickGameTable
func TransSlotProtoPickGameTable(src *ppb.PickGameTable) *pb.Slot_Proto_PickGameTable {
	dist := new(pb.Slot_Proto_PickGameTable)
	dist.ID = int32(src.GetId())
	dist.GroupID = int32(src.GetGroupId())
	dist.PickCount = int32(src.GetPickCount())
	dist.RewardValue = int32(src.GetRewardValue())
	dist.IsGhost = int32(src.GetIsGhost())
	return dist
}

// TransSlotProtoSelectionTable 轉換 SelectionTable
func TransSlotProtoSelectionTable(src *ppb.SelectionTable) *pb.Slot_Proto_SelectionTable {
	dist := new(pb.Slot_Proto_SelectionTable)
	dist.ID = int32(src.GetId())
	dist.GroupID = int32(src.GetGroupId())
	dist.SelectionType = src.GetSelectionType()
	dist.SelectionValue = int32(src.GetSelectionValue())
	dist.IsGoldsymbol = int32(src.GetIsGoldsymbol())
	dist.GoldCondition = src.GetGoldCondition()
	dist.ConditionValue = src.GetConditionValue()
	return dist
}

// TransSlotProtoRespinTable 轉換 RespinTable
func TransSlotProtoRespinTable(src *ppb.RespinTable) *pb.Slot_Proto_RespinTable {
	dist := new(pb.Slot_Proto_RespinTable)
	dist.ID = int32(src.GetId())
	dist.GroupID = int32(src.GetGroupId())
	dist.RespinCount = int32(src.GetRespinCount())
	dist.AppearReel = src.GetAppearReel()
	dist.RespinType = src.GetRespinType()
	dist.RespinValue = int32(src.GetRespinValue())
	dist.IsLimit = int32(src.GetIsLimit())
	dist.LimitType = src.GetLimitType()
	dist.LimitValue = src.GetLimitValue()
	return dist
}

// TransSlotProtoReelBgDistribution 轉換 ReelBgDistribution
func TransSlotProtoReelBgDistribution(src *ppb.ReelBgDistribution) *pb.Slot_Proto_ReelBGDistribution {
	dist := new(pb.Slot_Proto_ReelBGDistribution)
	dist.ID = int32(src.GetId())
	dist.GroupID = int32(src.GetGroupId())
	dist.ReelNum = int32(src.GetReelNum())
	dist.BgEnum = src.GetBgEnum()
	dist.Proportion = int32(src.GetProportion())
	return dist
}

// TransSlotProtoExtraRule 轉換 ExtraRule
func TransSlotProtoExtraRule(src *ppb.ExtraRule) *pb.Slot_Proto_ExtraRule {
	dist := new(pb.Slot_Proto_ExtraRule)
	dist.ID = int32(src.GetId())
	dist.GroupID = int32(src.GetGroupId())
	dist.UseType = src.GetUseType()
	dist.Condition = src.GetCondition()
	dist.ConditionValue = src.GetConditionValue()
	dist.Action = src.GetAction()
	dist.ActionValue = src.GetActionValue()
	return dist
}
