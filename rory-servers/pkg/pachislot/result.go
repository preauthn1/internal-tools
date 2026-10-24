package pachislot

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
	"example.invalid/ServerTeam/Servers/pkg/prob"
)

// SymbolGroup 一整個盤面
//  盤面內容:
//  [     Row1      Row2      Row3
//    [ pb.Symbol, pb.Symbol, pb.Symbol ], // Reel 1
//    [ pb.Symbol, pb.Symbol, pb.Symbol ], // Reel 2
//    [ pb.Symbol, pb.Symbol, pb.Symbol ], // Reel 3
//  ]
type SymbolGroup [3][3]pb.Symbol

// Match 是否連線指定線與符號
func (sg SymbolGroup) Match(line *pb.PayLine, s pb.Symbol) bool {
	if line == nil {
		return false
	}
	for _, pos := range line.Positions {
		if sg[pos.Reel-1][pos.Row-1] != s {
			return false
		}
	}
	return true
}

// Symbol 取得指定位置的符號
func (sg SymbolGroup) Symbol(pos *pb.Position) pb.Symbol {
	return sg[pos.Reel-1][pos.Row-1]
}

// Positions 取得指定符號的位置
func (sg SymbolGroup) Positions(s pb.Symbol) []*pb.Position {
	var positions []*pb.Position
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if sg[i][j] == s {
				positions = append(positions, &pb.Position{Reel: int32(i + 1), Row: int32(j + 1)})
			}
		}
	}
	return positions
}

func (sg SymbolGroup) String() string {
	return fmt.Sprintf("\n%8s\t%8s\t%8s\n%8s\t%8s\t%8s\n%8s\t%8s\t%8s\n",
		sg[0][0], sg[1][0], sg[2][0],
		sg[0][1], sg[1][1], sg[2][1],
		sg[0][2], sg[1][2], sg[2][2])
}

// ReelStripGroup 輪軸帶群組
//  map:
//  {Reel}: {Symbol}, {Symbol}, {Symbol}, {Symbol}, ...
type ReelStripGroup map[int][]pb.Symbol

// NewReelStripGroup 建立輪軸帶群組
func NewReelStripGroup(strips []*pb.ReelStripSymbol) ReelStripGroup {

	sort.Slice(strips, func(i, j int) bool {
		if strips[i].Reel == strips[j].Reel {
			return strips[i].Row < strips[j].Row
		}
		return strips[i].Reel < strips[j].Reel
	})

	g := ReelStripGroup{}

	for _, opt := range strips {
		g[int(opt.Reel)] = append(g[int(opt.Reel)], opt.Symbol)
	}

	return g
}

// Valid 驗證是否正確
//  * Reel 要有 1, 2, 3
//  * 每個 Reel 要有3個以上的Symbol
func (g ReelStripGroup) Valid() bool {
	for i := 0; i < 3; i++ {
		if symbols, ok := g[i+1]; ok {
			if len(symbols) < 3 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

// SymbolIndex 定位符號在每個Reel的位置
func (g ReelStripGroup) SymbolIndex(s pb.Symbol) (idxMap map[int][]int) {
	idxMap = make(map[int][]int)
	for reel, symbols := range g {
		for i, symbol := range symbols {
			if symbol == s {
				idxMap[reel] = append(idxMap[reel], i)
			}
		}
	}
	return
}

// ReelStripMap 輪軸帶停輪中間一排的輪軸帶idx
//  盤面結果為:
//  BL(02)    BL(10)     BR(06)
//  B7(03)    BL(11)     BR(07)  --> (3, 11, 7)
//  B7(04)    BL(12)     BR(08)
type ReelStripMap [3]int

// StopIndex 取得各輪軸停輪位置
func (m ReelStripMap) StopIndex() (int, int, int) {
	return m[0], m[1], m[2]
}

// FixIndex 修正Index
func (m ReelStripMap) FixIndex(strips ReelStripGroup) ReelStripMap {

	var fix = func(idx int, symbols []pb.Symbol) int {
		if idx < 0 {
			return idx + len(symbols)
		}
		if idx >= len(symbols) {
			return idx - len(symbols)
		}
		return idx
	}

	m[0] = fix(m[0], strips[1])
	m[1] = fix(m[1], strips[2])
	m[2] = fix(m[2], strips[3])

	return m
}

// SymbolGroup 依照排列帶組合出對應停輪的結果
//  盤面內容:
//  [     Row1      Row2      Row3
//    [ pb.Symbol, pb.Symbol, pb.Symbol ], // Reel 1
//    [ pb.Symbol, pb.Symbol, pb.Symbol ], // Reel 2
//    [ pb.Symbol, pb.Symbol, pb.Symbol ], // Reel 3
//  ]
func (m ReelStripMap) SymbolGroup(strips ReelStripGroup) SymbolGroup {

	group := SymbolGroup{}

	for i := 0; i < 3; i++ {
		idx := m[i]
		// 中間軸
		group[i][1] = strips[i+1][idx]

		// 上一格
		prevIdx := idx - 1
		if prevIdx < 0 {
			prevIdx = len(strips[i+1]) - 1
		}
		group[i][0] = strips[i+1][prevIdx]

		// 下一格
		nextIdx := idx + 1
		if nextIdx >= len(strips[i+1]) {
			nextIdx = 0
		}
		group[i][2] = strips[i+1][nextIdx]
	}

	return group
}

// ResultSet 牌型結果
//  一般流程: 知道牌型結果 -> 在 PayLine 選一條線 -> 在輪軸帶裡的符號挑出符號，那格在線上
type ResultSet interface {
	Type() pb.ResultType                                                          // 牌型
	Symbol() pb.Symbol                                                            // 規則符號
	Match([]*pb.PayLine, SymbolGroup) bool                                        // 檢查盤面是否符合規則
	Generete([]*pb.PayLine, ReelStripGroup) ([]*pb.Position, ReelStripMap, error) // 建立牌型盤面，產生的盤面可能與其他牌型結果相符，請自行比對，並決定是否重新產生
}

// GenereteResultSet 依照牌型產生盤面結果
func GenereteResultSet(rt pb.ResultType) ResultSet {
	switch rt {
	case pb.ResultType_Empty:
		return nil
	case pb.ResultType_None: // 沒有中獎，雜牌
		return new(NoneResultSet)
	case pb.ResultType_FeatureNone: // 純槓，雜牌
		return new(NoneResultSet)
	case pb.ResultType_JACNone: // JAC槓，前面兩個Replay、第三個銅鐘 在線上
		return new(JACNoneResultSet)
	case pb.ResultType_BlueSevenLink: // 藍7連線
		return new(BlueSevenLinkResultSet)
	case pb.ResultType_RedSevenLink: // 紅7連線
		return new(RedSevenLinkResultSet)
	case pb.ResultType_BarLink: // Bar連線
		return new(BarLinkResultSet)
	case pb.ResultType_BellLink: // 銅鐘連線
		return new(BellLinkResultSet)
	case pb.ResultType_WaterMelonLink: // 西瓜連線
		return new(WaterMelonLinkResultSet)
	case pb.ResultType_BerryLink: // 櫻桃連線
		return new(BerryLinkResultSet)
	case pb.ResultType_ReplayLink: // Replay連線
		return new(ReplayLinkResultSet)
	case pb.ResultType_DoubleSevenLink: // 雙7(藍7紅7混合)連線
		return new(DoubleSevenLinkResultSet)
	case pb.ResultType_BlueSevenMiss: // 藍7 前兩格連線，第三格出現但沒連線
		return nil
	case pb.ResultType_RedSevenMiss: // 紅7 前兩格連線，第三格出現但沒連線
		return nil
	case pb.ResultType_BarMiss: // Bar 前兩格連線，第三格出現但沒連線
		return nil
	case pb.ResultType_BellMiss: // 銅鐘 前兩格連線，第三格出現但沒連線
		return nil
	case pb.ResultType_WaterMelonMiss: // 西瓜 前兩格連線，第三格出現但沒連線
		return nil
	case pb.ResultType_BerryMiss: // 櫻桃 前兩格連線，第三格出現但沒連線
		return nil
	case pb.ResultType_ReplayMiss: // Replay 前兩格連線，第三格出現但沒連線
		return nil
	case pb.ResultType_DoubleSevenMiss: // 雙7(藍7紅7混合) 前兩格連線，第三格出現但沒連線
		return nil
	case pb.ResultType_ReleaseRT: // 需要解除RT 骰出後要依照當前大獎是 RB / BB，對應到 ReleaseRTonBB / ReleaseRTonRB，雜牌
		return new(NoneResultSet)
	case pb.ResultType_ReleaseRTonBB: // 需要解除RT 當前大獎為BB，雜牌
		return new(NoneResultSet)
	case pb.ResultType_ReleaseRTonRB: // 需要解除RT 當前大獎為RB，雜牌
		return new(NoneResultSet)
	case pb.ResultType_BonusSure: // Bonus確定，雜牌
		return new(NoneResultSet)
	case pb.ResultType_BerryLeftAny: // 櫻桃出現在左邊任意位置
		return new(BerryLeftAnyResultSet)
	}

	return nil
}

// NoneResultSet 沒有中獎牌型，亂
type NoneResultSet struct{}

// Type 牌型
func (*NoneResultSet) Type() pb.ResultType {
	return pb.ResultType_None
}

// Symbol 規則符號
func (*NoneResultSet) Symbol() pb.Symbol {
	return pb.Symbol(-1)
}

// Match 檢查盤面是否符合規則
func (rs *NoneResultSet) Match(paylines []*pb.PayLine, sg SymbolGroup) bool {
	return false
}

// Generete 建立牌型盤面，產生的盤面可能與其他牌型結果相符，請自行比對，並決定是否重新產生
func (rs *NoneResultSet) Generete(paylines []*pb.PayLine, strips ReelStripGroup) ([]*pb.Position, ReelStripMap, error) {
	idx1 := rand.Intn(len(strips[1]))
	idx2 := rand.Intn(len(strips[2]))
	idx3 := rand.Intn(len(strips[3]))
	return []*pb.Position(nil), ReelStripMap{idx1, idx2, idx3}, nil
}

// JACNoneResultSet 在一條線中，Replay兩個，銅鐘一個
type JACNoneResultSet struct{}

// Type 牌型
func (*JACNoneResultSet) Type() pb.ResultType {
	return pb.ResultType_None
}

// Symbol 規則符號
func (*JACNoneResultSet) Symbol() pb.Symbol {
	return pb.Symbol(-1)
}

// Match 檢查盤面是否符合規則
func (rs *JACNoneResultSet) Match(paylines []*pb.PayLine, sg SymbolGroup) bool {

	var match = func(line *pb.PayLine, sg SymbolGroup) bool {
		positions := line.Positions
		sort.Slice(positions, func(i, j int) bool {
			return positions[i].Reel < positions[j].Reel
		})
		// Reel1
		if sg.Symbol(positions[0]) != pb.Symbol_Replay {
			return false
		}
		// Reel2
		if sg.Symbol(positions[1]) != pb.Symbol_Replay {
			return false
		}
		// Reel3
		if sg.Symbol(positions[2]) != pb.Symbol_Bell {
			return false
		}
		return true
	}

	for _, line := range paylines {
		if match(line, sg) {
			return true
		}
	}

	return false
}

// Generete 建立牌型盤面，產生的盤面可能與其他牌型結果相符，請自行比對，並決定是否重新產生
func (rs *JACNoneResultSet) Generete(paylines []*pb.PayLine, strips ReelStripGroup) ([]*pb.Position, ReelStripMap, error) {
	// 檢查 PayLine
	if len(paylines) == 0 {
		return []*pb.Position(nil), ReelStripMap{}, ErrPayLineNotFound
	}

	// 找出每個Reel的符號的idx
	replayMap := strips.SymbolIndex(pb.Symbol_Replay)
	bellMap := strips.SymbolIndex(pb.Symbol_Bell)

	// 檢查每個Reel都要有符號可以抓取
	if len(replayMap[1]) == 0 {
		return []*pb.Position(nil), ReelStripMap{}, ErrResultSetNotGeneratable
	}
	if len(replayMap[2]) == 0 {
		return []*pb.Position(nil), ReelStripMap{}, ErrResultSetNotGeneratable
	}
	if len(bellMap[3]) == 0 {
		return []*pb.Position(nil), ReelStripMap{}, ErrResultSetNotGeneratable
	}

	var pick = func(data []int) int {
		r := new(prob.Random)
		for _, idx := range data {
			r.Add(idx, 1)
		}
		return r.Choose().(int)
	}

	idx1 := pick(replayMap[1])
	idx2 := pick(replayMap[2])
	idx3 := pick(bellMap[3])

	m := ReelStripMap{idx1, idx2, idx3}
	m = m.FixIndex(strips)

	return []*pb.Position(nil), m, nil
}

// BellLinkResultSet 銅鐘連線牌型
type BellLinkResultSet struct{}

// Type 牌型
func (*BellLinkResultSet) Type() pb.ResultType {
	return pb.ResultType_BellLink
}

// Symbol 規則符號
func (*BellLinkResultSet) Symbol() pb.Symbol {
	return pb.Symbol_Bell
}

// Match 檢查盤面是否符合規則
func (rs *BellLinkResultSet) Match(paylines []*pb.PayLine, sg SymbolGroup) bool {
	return matchLinkResultSet(paylines, sg, rs.Symbol())
}

// Generete 建立牌型盤面，產生的盤面可能與其他牌型結果相符，請自行比對，並決定是否重新產生
func (rs *BellLinkResultSet) Generete(paylines []*pb.PayLine, strips ReelStripGroup) ([]*pb.Position, ReelStripMap, error) {
	return generateLinkResultSet(paylines, strips, rs.Symbol())
}

// BlueSevenLinkResultSet 藍7連線牌型
type BlueSevenLinkResultSet struct{}

// Type 牌型
func (*BlueSevenLinkResultSet) Type() pb.ResultType {
	return pb.ResultType_BlueSevenLink
}

// Symbol 規則符號
func (*BlueSevenLinkResultSet) Symbol() pb.Symbol {
	return pb.Symbol_BlueSeven
}

// Match 檢查盤面是否符合規則
func (rs *BlueSevenLinkResultSet) Match(paylines []*pb.PayLine, sg SymbolGroup) bool {
	return matchLinkResultSet(paylines, sg, rs.Symbol())
}

// Generete 建立牌型盤面，產生的盤面可能與其他牌型結果相符，請自行比對，並決定是否重新產生
func (rs *BlueSevenLinkResultSet) Generete(paylines []*pb.PayLine, strips ReelStripGroup) ([]*pb.Position, ReelStripMap, error) {
	return generateLinkResultSet(paylines, strips, rs.Symbol())
}

// RedSevenLinkResultSet 藍7連線牌型
type RedSevenLinkResultSet struct{}

// Type 牌型
func (*RedSevenLinkResultSet) Type() pb.ResultType {
	return pb.ResultType_RedSevenLink
}

// Symbol 規則符號
func (*RedSevenLinkResultSet) Symbol() pb.Symbol {
	return pb.Symbol_RedSeven
}

// Match 檢查盤面是否符合規則
func (rs *RedSevenLinkResultSet) Match(paylines []*pb.PayLine, sg SymbolGroup) bool {
	return matchLinkResultSet(paylines, sg, rs.Symbol())
}

// Generete 建立牌型盤面，產生的盤面可能與其他牌型結果相符，請自行比對，並決定是否重新產生
func (rs *RedSevenLinkResultSet) Generete(paylines []*pb.PayLine, strips ReelStripGroup) ([]*pb.Position, ReelStripMap, error) {
	return generateLinkResultSet(paylines, strips, rs.Symbol())
}

// BarLinkResultSet Bar連線牌型
type BarLinkResultSet struct{}

// Type 牌型
func (*BarLinkResultSet) Type() pb.ResultType {
	return pb.ResultType_BarLink
}

// Symbol 規則符號
func (*BarLinkResultSet) Symbol() pb.Symbol {
	return pb.Symbol_Bar
}

// Match 檢查盤面是否符合規則
func (rs *BarLinkResultSet) Match(paylines []*pb.PayLine, sg SymbolGroup) bool {
	return matchLinkResultSet(paylines, sg, rs.Symbol())
}

// Generete 建立牌型盤面，產生的盤面可能與其他牌型結果相符，請自行比對，並決定是否重新產生
func (rs *BarLinkResultSet) Generete(paylines []*pb.PayLine, strips ReelStripGroup) ([]*pb.Position, ReelStripMap, error) {
	return generateLinkResultSet(paylines, strips, rs.Symbol())
}

// WaterMelonLinkResultSet 西瓜連線牌型
type WaterMelonLinkResultSet struct{}

// Type 牌型
func (*WaterMelonLinkResultSet) Type() pb.ResultType {
	return pb.ResultType_WaterMelonLink
}

// Symbol 規則符號
func (*WaterMelonLinkResultSet) Symbol() pb.Symbol {
	return pb.Symbol_WaterMelon
}

// Match 檢查盤面是否符合規則
func (rs *WaterMelonLinkResultSet) Match(paylines []*pb.PayLine, sg SymbolGroup) bool {
	return matchLinkResultSet(paylines, sg, rs.Symbol())
}

// Generete 建立牌型盤面，產生的盤面可能與其他牌型結果相符，請自行比對，並決定是否重新產生
func (rs *WaterMelonLinkResultSet) Generete(paylines []*pb.PayLine, strips ReelStripGroup) ([]*pb.Position, ReelStripMap, error) {
	return generateLinkResultSet(paylines, strips, rs.Symbol())
}

// BerryLinkResultSet 櫻桃連線牌型
type BerryLinkResultSet struct{}

// Type 牌型
func (*BerryLinkResultSet) Type() pb.ResultType {
	return pb.ResultType_BerryLink
}

// Symbol 規則符號
func (*BerryLinkResultSet) Symbol() pb.Symbol {
	return pb.Symbol_Berry
}

// Match 檢查盤面是否符合規則
func (rs *BerryLinkResultSet) Match(paylines []*pb.PayLine, sg SymbolGroup) bool {
	return matchLinkResultSet(paylines, sg, rs.Symbol())
}

// Generete 建立牌型盤面，產生的盤面可能與其他牌型結果相符，請自行比對，並決定是否重新產生
func (rs *BerryLinkResultSet) Generete(paylines []*pb.PayLine, strips ReelStripGroup) ([]*pb.Position, ReelStripMap, error) {
	return generateLinkResultSet(paylines, strips, rs.Symbol())
}

// ReplayLinkResultSet Replay連線牌型
type ReplayLinkResultSet struct{}

// Type 牌型
func (*ReplayLinkResultSet) Type() pb.ResultType {
	return pb.ResultType_ReplayLink
}

// Symbol 規則符號
func (*ReplayLinkResultSet) Symbol() pb.Symbol {
	return pb.Symbol_Replay
}

// Match 檢查盤面是否符合規則
func (rs *ReplayLinkResultSet) Match(paylines []*pb.PayLine, sg SymbolGroup) bool {
	return matchLinkResultSet(paylines, sg, rs.Symbol())
}

// Generete 建立牌型盤面，產生的盤面可能與其他牌型結果相符，請自行比對，並決定是否重新產生
func (rs *ReplayLinkResultSet) Generete(paylines []*pb.PayLine, strips ReelStripGroup) ([]*pb.Position, ReelStripMap, error) {
	return generateLinkResultSet(paylines, strips, rs.Symbol())
}

// DoubleSevenLinkResultSet 7藍雙連線牌型
type DoubleSevenLinkResultSet struct{}

// Type 牌型
func (*DoubleSevenLinkResultSet) Type() pb.ResultType {
	return pb.ResultType_DoubleSevenLink
}

// Symbol 規則符號
func (*DoubleSevenLinkResultSet) Symbol() pb.Symbol {
	return pb.Symbol_BlueSeven
}

// Match 檢查盤面是否符合規則
func (rs *DoubleSevenLinkResultSet) Match(paylines []*pb.PayLine, sg SymbolGroup) bool {

	var cnt int
	for _, line := range paylines {
		if sg.Match(line, rs.Symbol()) {
			cnt++
		}
	}

	return cnt == 2
}

// Generete 建立牌型盤面，產生的盤面可能與其他牌型結果相符，請自行比對，並決定是否重新產生
func (rs *DoubleSevenLinkResultSet) Generete(paylines []*pb.PayLine, strips ReelStripGroup) ([]*pb.Position, ReelStripMap, error) {

	// 找出每個Reel的符號的idx
	idxMap := strips.SymbolIndex(rs.Symbol())

	// 檢查可能性必須是 Reel1 兩個連再一起，Reel2 一個，Reel3 兩個連再一起
	if len(idxMap[1]) < 2 || len(idxMap[2]) < 1 || len(idxMap[3]) < 2 {
		return []*pb.Position(nil), ReelStripMap{}, ErrResultSetNotGeneratable
	}

	// Reel1 不連續
	if math.Abs(float64(idxMap[1][0]-idxMap[1][1])) != 1 {
		return []*pb.Position(nil), ReelStripMap{}, ErrResultSetNotGeneratable
	}
	// Reel3 不連續
	if math.Abs(float64(idxMap[3][0]-idxMap[3][1])) != 1 {
		return []*pb.Position(nil), ReelStripMap{}, ErrResultSetNotGeneratable
	}

	var pick = func(data []int) int {
		r := new(prob.Random)
		for _, idx := range data {
			r.Add(idx, 1)
		}
		return r.Choose().(int)
	}

	drift := pick([]int{0, 1})
	idx1 := idxMap[1][0] + drift
	idx2 := pick(idxMap[2])
	idx3 := idxMap[3][1] - drift

	m := ReelStripMap{idx1, idx2, idx3}
	m = m.FixIndex(strips)

	sg := m.SymbolGroup(strips)

	return sg.Positions(rs.Symbol()), m, nil
}

// BerryLeftAnyResultSet 藍7連線牌型
type BerryLeftAnyResultSet struct{}

// Type 牌型
func (*BerryLeftAnyResultSet) Type() pb.ResultType {
	return pb.ResultType_BerryLeftAny
}

// Symbol 規則符號
func (*BerryLeftAnyResultSet) Symbol() pb.Symbol {
	return pb.Symbol_Berry
}

// Match 檢查盤面是否符合規則
func (rs *BerryLeftAnyResultSet) Match(paylines []*pb.PayLine, sg SymbolGroup) bool {
	return sg[0][0] == pb.Symbol_Berry || sg[0][1] == pb.Symbol_Berry || sg[0][2] == pb.Symbol_Berry
}

// Generete 建立牌型盤面，產生的盤面可能與其他牌型結果相符，請自行比對，並決定是否重新產生
func (rs *BerryLeftAnyResultSet) Generete(paylines []*pb.PayLine, strips ReelStripGroup) ([]*pb.Position, ReelStripMap, error) {
	// 找出每個Reel的符號的idx
	idxMap := strips.SymbolIndex(rs.Symbol())

	var pick = func(data []int) int {
		r := new(prob.Random)
		for _, idx := range data {
			r.Add(idx, 1)
		}
		return r.Choose().(int)
	}
	idx1 := pick(idxMap[1])
	idx2 := rand.Intn(len(strips[2]))
	idx3 := rand.Intn(len(strips[3]))

	// 隨機飄移上中下
	drift := pick([]int{-1, 0, 1})
	idx1 += drift

	m := ReelStripMap{idx1, idx2, idx3}
	m = m.FixIndex(strips)

	// 中心點 Row2 把飄移補回
	return []*pb.Position{{Reel: 1, Row: int32(-drift) + 2}}, m, nil
}

// generateLinkResultSet 建立連線牌型盤面
func generateLinkResultSet(paylines []*pb.PayLine, strips ReelStripGroup, symbol pb.Symbol) ([]*pb.Position, ReelStripMap, error) {
	// 檢查 PayLine
	if len(paylines) == 0 {
		return []*pb.Position(nil), ReelStripMap{}, ErrPayLineNotFound
	}

	// 找出每個Reel的符號的idx
	idxMap := strips.SymbolIndex(symbol)

	// 檢查每個Reel都要有符號可以抓取
	for i := 0; i < 3; i++ {
		if len(idxMap[i+1]) == 0 {
			return []*pb.Position(nil), ReelStripMap{}, ErrResultSetNotGeneratable
		}
	}

	var pick = func(data []int) int {
		r := new(prob.Random)
		for _, idx := range data {
			r.Add(idx, 1)
		}
		return r.Choose().(int)
	}
	idx1 := pick(idxMap[1])
	idx2 := pick(idxMap[2])
	idx3 := pick(idxMap[3])

	// 隨機選一條線
	r := new(prob.Random)
	for _, line := range paylines {
		r.Add(line, 1)
	}

	line := r.Choose().(*pb.PayLine)

	// 依照賠線修正idx
	var lineIdx [3]int
	for _, pos := range line.Positions {
		lineIdx[pos.Reel-1] = int(pos.Row - 1)
	}
	var fix = func(idx, std int) int { return idx + (1 - std) }
	idx1 = fix(idx1, lineIdx[0])
	idx2 = fix(idx2, lineIdx[1])
	idx3 = fix(idx3, lineIdx[2])

	m := ReelStripMap{idx1, idx2, idx3}
	m = m.FixIndex(strips)

	return line.Positions, m, nil
}

// matchLinkResultSet 檢查盤面是否符合連線規則
func matchLinkResultSet(paylines []*pb.PayLine, sg SymbolGroup, symbol pb.Symbol) bool {
	// 檢查每個連線是否達成，一條連線達成就是達成
	for _, line := range paylines {
		if sg.Match(line, symbol) {
			return true
		}
	}
	return false
}

// ResultSetValidator 盤面檢查器
type ResultSetValidator struct {
	ResultSets []ResultSet
}

// NewResultSetValidator 建立盤面檢查器
func NewResultSetValidator(types []pb.ResultType) *ResultSetValidator {
	v := new(ResultSetValidator)
	for _, rt := range types {
		v.ResultSets = append(v.ResultSets, GenereteResultSet(rt))
	}
	return v
}

// MatchOthers 檢查盤面是否也符合其他盤面
func (v *ResultSetValidator) MatchOthers(self pb.ResultType, sg SymbolGroup, paylines []*pb.PayLine) bool {
	for _, set := range v.ResultSets {
		if set == nil {
			continue
		}
		if set.Type() == self {
			continue
		}
		if set.Match(paylines, sg) {
			return true
		}
	}
	return false
}
