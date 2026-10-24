package slot

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/prob"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// GameResultManager 專門控制GameResult
type GameResultManager struct {
	GameResult      GameResult
	MultiSymbol     map[int32]int64 // 特殊加倍符號 map[{SymbolID}] = {倍數}
	Reel            int32           // 共有多少直軸
	Row             int32           // 共有多少橫軸
	InvalidPosition []Position
	SymbolGroupID   int32 // 此圖由哪個SymbolGroupID轉出

	BackGroundGroupID int32 // 此圖由哪個GoldReelGroupID轉出
	BackgroundResult  *BackgroundResult

	AfterGenerateFunc []func(*GameResultManager) `json:"-"`
}

func (manager *GameResultManager) String() string {
	return fmt.Sprintf("Mode: %s, SymbolGroupoID: %d\n%v", manager.Mode(), manager.SymbolGroupID, manager.GameResult.String())
}

// NewGameResultManager create GameResultManager
func NewGameResultManager(reel, row int32) (*GameResultManager, error) {
	if reel == 0 {
		return nil, fmt.Errorf("Reel number could not be 0")
	}
	if row == 0 {
		return nil, fmt.Errorf("Row number could not be 0")
	}

	manager := new(GameResultManager)
	manager.Reel = reel
	manager.Row = row

	manager.MultiSymbol = make(map[int32]int64)

	return manager, nil
}

// CloneGameResultManager create GameResultManager
func CloneGameResultManager(origin *GameResultManager) *GameResultManager {

	manager := new(GameResultManager)
	manager.GameResult = origin.GameResult
	manager.Reel = origin.Reel
	manager.Row = origin.Row
	manager.SymbolGroupID = origin.SymbolGroupID

	return manager
}

// IsReelMode 目前第一維陣列為Reel組合
func (manager *GameResultManager) IsReelMode() bool {
	if len(manager.GameResult) == 0 {
		return false
	}

	if manager.Reel == manager.Row {

		// 第一筆
		s := manager.GameResult[0][0]

		var sameReel bool = true

		// 檢查之後每一筆是否皆為同一Reel
		for i := 1; i < len(manager.GameResult[0])-1; i++ {
			if manager.GameResult[0][i].GetReelNum() != s.GetReelNum() {
				sameReel = false
			}
		}

		return sameReel
	}

	return len(manager.GameResult) == int(manager.Reel) && len(manager.GameResult[0]) == int(manager.Row)
}

// IsRowMode 目前第一維陣列為Row組合
func (manager *GameResultManager) IsRowMode() bool {
	if len(manager.GameResult) == 0 {
		return false
	}

	if manager.Reel == manager.Row {

		// 第一筆
		s := manager.GameResult[0][0]

		var sameReel bool = true

		// 檢查之後每一筆是否皆為同一Reel, 是的話為Reel不是Row
		for i := 1; i < len(manager.GameResult[0])-1; i++ {
			if manager.GameResult[0][i].GetReelNum() != s.GetReelNum() {
				sameReel = false
			}
		}

		return !sameReel
	}

	return len(manager.GameResult) == int(manager.Row) && len(manager.GameResult[0]) == int(manager.Reel)
}

// Mode 取得目前模式
func (manager *GameResultManager) Mode() string {
	if manager.IsReelMode() {
		return "Reel"
	}

	if manager.IsRowMode() {
		return "Row"
	}

	return ""
}

// ReelMode 切換為ReelMode
func (manager *GameResultManager) ReelMode() {
	if len(manager.GameResult) == 0 {
		return
	}
	if manager.IsReelMode() {
		return
	}

	manager.GameResult = manager.GameResult.Transfor()
}

// RowMode 切換為RowMode
func (manager *GameResultManager) RowMode() {
	if len(manager.GameResult) == 0 {
		return
	}
	if manager.IsRowMode() {
		return
	}

	manager.GameResult = manager.GameResult.Transfor()
}

// Reverse 對調
func (manager *GameResultManager) Reverse() *GameResultManager {
	manager.GameResult = manager.GameResult.Reverse()
	return manager
}

// RegisterAfterGenerateFunc 註冊結束後執行邏輯
func (manager *GameResultManager) RegisterAfterGenerateFunc(f func(*GameResultManager)) {
	manager.AfterGenerateFunc = append(manager.AfterGenerateFunc, f)
}

// GenerateGameResult 產生指定GroupID的輪軸
// 一個Reel一個Reel骰
// 先決定停在最中間的Row是哪個Order
// 再向上與向下骰出對應的Symbol
func (manager *GameResultManager) GenerateGameResult(groupID int32, SymbolTable []*protodata.SlotSymbolDistribution) *GameResultManager {

	defer func() {
		if len(manager.AfterGenerateFunc) > 0 {
			for _, fn := range manager.AfterGenerateFunc {
				fn(manager)
			}
		}
	}()

	manager.SymbolGroupID = groupID

	result := make([][]*protodata.SlotSymbolDistribution, 0, manager.Reel)

	// 過濾同SymbolGroupID的符號
	var symbols []*protodata.SlotSymbolDistribution

	for _, symbol := range SymbolTable {
		if int32(symbol.GetGroupId()) != groupID {
			continue
		}

		symbols = append(symbols, symbol)
	}

	// 判斷最中間的Row
	var middleRow int32 = manager.Row / 2
	if manager.Row%2 == 1 {
		middleRow++
	}

	mlog.Trace("GroupID: %d, Reel: %d, Row: %d, Middle: %d", groupID, manager.Reel, manager.Row, middleRow)

	for reelNum := int32(1); reelNum <= manager.Reel; reelNum++ {
		mlog.Trace("Begin Reel: %d", reelNum)

		rows := make([]*protodata.SlotSymbolDistribution, manager.Row)

		// 整理此Reel所有Order內的符號
		ordersSymbols := make(map[int32][]*protodata.SlotSymbolDistribution)
		for _, symbol := range symbols {
			if int32(symbol.GetReelNum()) != reelNum {
				continue
			}
			dist := *symbol
			ordersSymbols[int32(symbol.GetOrder())] = append(ordersSymbols[int32(symbol.GetOrder())], &dist)
		}

		if len(ordersSymbols) == 0 {
			mlog.Warn("Reel: %d, has no symbols setting", reelNum)
			continue
		}

		var minOrder, maxOrder int32 = 1, 1
		r := new(prob.Random)
		for order, smbs := range ordersSymbols {
			mlog.Trace("Reel: %d, Order: %d, has %d item", reelNum, order, len(smbs))
			if len(smbs) == 0 {
				continue
			} else if len(smbs) == 1 {
				r.Add(order, int(smbs[0].GetProportion()))
			} else {
				r.Add(order, 1)
			}

			if minOrder > order {
				minOrder = order
			}
			if maxOrder < order {
				maxOrder = order
			}
		}

		// 由MiddleRow開始骰出停住的Order
		t := r.Choose().(int32)
		r.PrintOptions()
		if conf.TestMode {
			t = middleRow
		}
		mlog.Trace("Reel %d Stop at Order %d, MinOrder: %d, MaxOrder: %d", reelNum, t, minOrder, maxOrder)

		mlog.Trace("Dice MiddleRow")
		// 骰出MiddleRow符號
		if symbol := manager.DiceSymbol(ordersSymbols[t]); symbol == nil {
			mlog.Warn("Reel: %d, Row: %d, Order: %d, Dice nil symbol", reelNum, middleRow, t)
			continue
		} else {
			mlog.Trace("Reel: %d, Row: %d, Order: %d, Dice symbol %s, From: %v", reelNum, middleRow, t, symbol.GetSymbol(), manager.GetSymbolStrings(ordersSymbols[t]))
			rows[int(middleRow-1)] = symbol
		}

		// 向上查找Row
		mlog.Trace("Dice Upside")
		for i, j := t-1, middleRow-1; j > 0; i-- {
			if i < minOrder {
				i += int32(len(ordersSymbols))
			}
			// 骰出符號
			if symbol := manager.DiceSymbol(ordersSymbols[i]); symbol == nil {
				mlog.Warn("Reel: %d, Row: %d, Order: %d, Dice nil symbol", reelNum, j, i)
				continue
			} else {
				mlog.Trace("Reel: %d, Row: %d, Order: %d, Dice symbol %s, From: %v", reelNum, j, i, symbol.GetSymbol(), manager.GetSymbolStrings(ordersSymbols[i]))
				rows[int(j-1)] = symbol
			}
			j--
		}

		// 向下查找Row
		mlog.Trace("Dice Down")
		for i, j := t+1, middleRow+1; j <= manager.Row; i++ {
			if i > maxOrder {
				i -= int32(len(ordersSymbols))
			}
			// 骰出符號
			if symbol := manager.DiceSymbol(ordersSymbols[i]); symbol == nil {
				mlog.Warn("Reel: %d, Row: %d, Order: %d, Dice nil symbol", reelNum, j, i)
				continue
			} else {
				mlog.Trace("Reel: %d, Row: %d, Order: %d, Dice symbol %s, From: %v", reelNum, j, i, symbol.GetSymbol(), manager.GetSymbolStrings(ordersSymbols[i]))
				rows[int(j-1)] = symbol
			}
			j++
		}

		// 檢查是否有在不允許出現的位置內
		for j := range rows {
			for _, pos := range manager.InvalidPosition {
				if pos.Reel == reelNum && pos.Row == int32(j+1) {
					rows[j] = nil
					mlog.Trace("Remove Invalid Position %v", pos)
				}
			}
		}

		result = append(result, rows)
	}

	manager.GameResult = result
	manager.RowMode()

	return manager
}

// GenerateBackground 產生指定GroupID的背景
func (manager *GameResultManager) GenerateBackground(groupID int32, bgds []*protodata.ReelBgDistribution) *GameResultManager {

	manager.BackGroundGroupID = groupID

	bgr := new(BackgroundResult)
	bgr.Reel = manager.Reel
	bgr.Row = manager.Row

	for reelNum := int32(1); reelNum <= manager.Reel; reelNum++ {

		r := new(prob.Random)

		for _, bgd := range bgds {
			if manager.BackGroundGroupID != int32(bgd.GetGroupId()) {
				continue
			}
			if int32(bgd.GetReelNum()) == reelNum {
				r.Add(bgd, int(bgd.GetProportion()))
			}
		}

		t := r.Choose()

		if t != nil {
			dist := *t.(*protodata.ReelBgDistribution)
			bgr.AddDistribution(&dist)
		} else {
			bgr.AddDistribution(&protodata.ReelBgDistribution{
				ReelNum: proto.Uint32(uint32(reelNum)),
				BgEnum:  proto.String("bg0"),
			})
		}
	}

	manager.BackgroundResult = bgr

	return manager
}

// GetSymbolsByReelRow 取得指定GroupID, ReelNum, RowNum的Symbols
func (manager *GameResultManager) GetSymbolsByReelRow(groupID, reelNum, rowNum int32, SymbolTable []*protodata.SlotSymbolDistribution) []*protodata.SlotSymbolDistribution {
	var symbols []*protodata.SlotSymbolDistribution

	for _, row := range SymbolTable {
		if groupID != int32(row.GetGroupId()) {
			continue
		}
		if reelNum != int32(row.GetReelNum()) {
			continue
		}
		if rowNum != int32(row.GetOrder()) {
			continue
		}
		symbols = append(symbols, row)
	}

	return symbols
}

// GetSymbolString 取得符號二維陣列
func (manager *GameResultManager) GetSymbolString() [][]string {
	return manager.GameResult.GetSymbolString()
}

// ToSymbolString 取得 601;p1,p2,p3|p2,p3,p4|p3,p4,p5|p4,p5,p6|p5,p6,p7 字樣
func (manager *GameResultManager) ToSymbolString() string {
	return fmt.Sprintf("%d#%s", manager.SymbolGroupID, manager.GameResult.ToSymbolString())
}

// DiceSymbol 將傳入的Symbols依照Prop骰出一個
func (manager *GameResultManager) DiceSymbol(symbols []*protodata.SlotSymbolDistribution) *protodata.SlotSymbolDistribution {
	if len(symbols) == 0 {
		return nil
	}

	r := new(prob.Random)

	for _, row := range symbols {
		r.Add(row, int(row.GetProportion()))
	}

	t := r.Choose()

	r.PrintOptions()

	if t == nil {
		return nil
	}

	return t.(*protodata.SlotSymbolDistribution)
}

// GetSymbolStrings 取得傳入的Symbol
func (manager *GameResultManager) GetSymbolStrings(opts []*protodata.SlotSymbolDistribution) []string {
	var symbols []string
	for _, row := range opts {
		symbols = append(symbols, row.GetSymbol())
	}
	return symbols
}
