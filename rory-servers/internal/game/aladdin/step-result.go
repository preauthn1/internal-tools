package aladdin

import (
	"sort"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// ExpandOrder 拓展排序
type ExpandOrder string

// ExpandOrder
const (
	ExpandOrderASC  ExpandOrder = "ASC"
	ExpandOrderDESC ExpandOrder = "DESC"
)

// StepResult 每次轉完遊戲結果
type StepResult struct {
	expandSymbols map[string]*protodata.SymbolAttributes

	IsRespin bool
	IsExpand bool

	Order int

	ExpandReel   []int32
	ExpandSymbol string
	ExpandOrder  ExpandOrder

	AllExpandReels []int32

	UnexpandResult *slot.GameResultManager // 拓展前輪軸
	ExpandResult   *slot.GameResultManager // 拓展後輪軸

	GamePayResult []*slot.PayResult // 遊戲中獎內容

	GameBetMoney int64 // 下注倍率
	GamePayMoney int64 // 花費金額
	GameWinMoney int64 // 贏得金額

	PrevStep *StepResult
}

// NewStepResult create StepResult
func NewStepResult(order int, isRespin bool, betMoney, payMoney int64, attrs []*protodata.SymbolAttributes, result *slot.GameResultManager, prev *StepResult) *StepResult {
	step := new(StepResult)
	step.Order = order
	step.IsRespin = isRespin
	step.GameBetMoney = betMoney
	if step.IsRespin {
		step.GamePayMoney = 0
	} else {
		step.GamePayMoney = payMoney
	}
	step.PrevStep = prev

	// 檢查是否有拓展符號
	step.expandSymbols = make(map[string]*protodata.SymbolAttributes)
	for _, row := range attrs {
		if row.GetIsExpand() == 1 {
			step.expandSymbols[row.GetSymbol()] = row
		}
	}

	step.UnexpandResult = slot.CloneGameResultManager(result)

	step.Expand()

	return step
}

// Expand 產生遊戲結果包含拓展前後
func (step *StepResult) Expand() {

	symbolMap := make(map[int32]*protodata.SlotSymbolDistribution)

	// 檢查輪軸是否需要拓展
	step.UnexpandResult.ReelMode()
	for i, reel := range step.UnexpandResult.GameResult {

		var symbol *protodata.SlotSymbolDistribution

		for _, row := range reel {
			if _, ok := step.expandSymbols[row.GetSymbol()]; !ok {
				continue
			}
			symbol = row
			break
		}

		if symbol == nil {
			continue
		}

		// 檢查前次是否有拓展過
		if step.PrevStep == nil || step.PrevStep.ExpandResult.GameResult[i][0].GetSymbol() != symbol.GetSymbol() {
			symbolMap[int32(i+1)] = symbol
		}

	}

	step.ExpandResult = slot.CloneGameResultManager(step.UnexpandResult)

	step.UnexpandResult.ReelMode()

	for reelNum, symbol := range symbolMap {
		newReel := make([]*protodata.SlotSymbolDistribution, 0, step.ExpandResult.Row)

		// 只換符號～不換ID
		for _, org := range step.UnexpandResult.GameResult[reelNum-1] {
			smb := new(protodata.SlotSymbolDistribution)
			smb.Id = proto.Uint32(org.GetId())
			smb.GroupId = proto.Uint32(org.GetGroupId())
			smb.ReelNum = proto.Uint32(org.GetReelNum())
			smb.Order = proto.Uint32(org.GetOrder())
			smb.Proportion = proto.Uint32(org.GetProportion())
			smb.Symbol = proto.String(symbol.GetSymbol())
			newReel = append(newReel, smb)
		}

		var newResult slot.GameResult
		newResult = append(newResult, step.ExpandResult.GameResult[0:reelNum-1]...)
		newResult = append(newResult, newReel)
		newResult = append(newResult, step.ExpandResult.GameResult[reelNum:len(step.ExpandResult.GameResult)]...)

		step.IsExpand = true
		step.ExpandResult.GameResult = newResult
		step.ExpandReel = append(step.ExpandReel, int32(reelNum))
		step.ExpandSymbol = symbol.GetSymbol()
	}

	// 檢查所有拓展
	for i, reel := range step.ExpandResult.GameResult {
		// 檢查該軸是否全為拓展
		var hasDiff bool
		for _, row := range reel {
			if _, ok := step.expandSymbols[row.GetSymbol()]; !ok {
				hasDiff = true
				break
			}
		}
		if !hasDiff {
			step.AllExpandReels = append(step.AllExpandReels, int32(i+1))
			continue
		}
	}

	sort.Slice(step.AllExpandReels, func(i, j int) bool {
		return step.AllExpandReels[i] < step.AllExpandReels[j]
	})
}
