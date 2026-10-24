package baselogic

import (
	"fmt"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// LogBackgroundResult 印出背景結果
func (logic *SlotLogicBase) LogBackgroundResult(manager *slot.GameResultManager) {
	logic.GameLog.AddGameStep("背景(%d) %v", manager.BackGroundGroupID, manager.BackgroundResult)
}

// LogGameResult 印出盤面結果
func (logic *SlotLogicBase) LogGameResult(manager *slot.GameResultManager) {
	manager.RowMode()
	logic.GameLog.AddGameStep("盤面(%d) %s", manager.SymbolGroupID, manager.ToSymbolString())
}

// LogPayResultChip 印出賠付結果
func (logic *SlotLogicBase) LogPayResultChip(result *slot.PayResult, base, win int64) {

	var desc string

	switch result.MultiType {
	default:
		desc = fmt.Sprintf("開獎符號: %s:%s:%d, %v x %d = %d", result.SymbolName, result.Rule.GetLinkType(), result.PayCount, base, result.WayCount, win)
	case slot.MultiTypeApearWild:
		desc = fmt.Sprintf("開獎符號: %s:%s:%d, %v x %d(x%d) = %d", result.SymbolName, result.Rule.GetLinkType(), result.PayCount, base, result.WayCount, result.Multi, win)
	case slot.MultiTypeLinkWild:
		desc = fmt.Sprintf("開獎符號: %s:%s:%d, %v x [%d + %d(x%d)] = %d", result.SymbolName, result.Rule.GetLinkType(), result.PayCount, base, result.NormalWayCount, result.WildWayCount, result.Multi, win)
	case slot.MultiTypeLinkWildTarget:
		mMap := make(map[int32]int)
		for _, line := range result.SymbolPayLines {
			m := line.GetMultipleValue(result.MultiMode)
			mMap[m]++
		}

		var opts []string
		for multi, times := range mMap {
			opts = append(opts, fmt.Sprintf("%d(x%d)", times, multi))
		}
		desc = fmt.Sprintf("開獎符號: %s:%s:%d, %v x [%s] = %d", result.SymbolName, result.Rule.GetLinkType(), result.PayCount, base, strings.Join(opts, " + "), win)
	}

	var lines []string
	for _, line := range result.SymbolPayLines {

		var _m string
		if line.PayLine.Origin != nil {
			_m = fmt.Sprintf("%d", line.PayLine.Origin.GetLineId())
		} else {
			_m = line.PayLine.String()
		}

		lines = append(lines, fmt.Sprintf("%sx%d", _m, line.GetMultipleValue(result.MultiMode)))
	}

	if len(lines) > 0 {
		desc = fmt.Sprintf("%s, 連線內容: %s", desc, strings.Join(lines, ","))
	}

	logic.GameLog.AddGameStep(desc)
}
