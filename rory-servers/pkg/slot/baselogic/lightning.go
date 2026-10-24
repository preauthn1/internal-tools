package baselogic

import (
	"math/rand"
	"strconv"
	"strings"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/defender"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// LightningAward 獎勵
type LightningAward struct {
	WinMoney    int64                             `json:"win_money"`
	Position    slot.Position                     `json:"position"`
	Type        slot.PayType                      `json:"type"`
	JackpotType slot.JackpotLevel                 `json:"jp_type"`
	Symbol      *protodata.SlotSymbolDistribution `json:"symbol"`
}

// LightningStep 每次Step相關資訊
type LightningStep struct {
	Origin            *protodata.RespinTable  `json:"-"`
	PayResults        []*slot.PayResult       `json:"-"`
	RespinID          int32                   `json:"respin_id"`
	Rule              *protodata.SlotPayTable `json:"-"`
	RuleID            int32                   `json:"rule_id"`
	GameResultManager *slot.GameResultManager `json:"game_result"`
	LastChance        int                     `json:"last_chance"`
	Awards            []*LightningAward       `json:"awards"`
	WinMoney          int64                   `json:"win_money"` // 總共贏得金額
	Log               []string                `json:"log"`
	NextAdd           []string                `json:"next_add"` // 下一次Spin要塞入的符號
}

// NewLightningStep create LightningStep
func NewLightningStep(table *protodata.RespinTable, rule *protodata.SlotPayTable, result *slot.GameResultManager, lastChance int) *LightningStep {
	step := new(LightningStep)

	if table != nil {
		step.Origin = table
		step.RespinID = int32(table.GetId())
	}

	if rule != nil {
		step.Rule = rule
		step.RuleID = int32(rule.GetId())
	}

	step.GameResultManager = result
	step.LastChance = lastChance

	return step
}

func (s *LightningStep) String() string {
	return s.GameResultManager.String()
}

// PrintAward 紀錄獎勵
func (s *LightningStep) PrintAward(logic *SlotLogicBase) {
	for _, award := range s.Awards {
		switch award.Type {
		case slot.PayTypeChip, slot.PayTypeChipXRate, slot.PayTypeAccuMulate:
			logic.GameLog.AddGameStep("%d_%d:%s(%s)=%d", award.Position.Reel, award.Position.Row, award.Symbol.GetSymbol(), award.Type, award.WinMoney)
		case slot.PayTypeJackpot:
			// 15 科
			if award.Position.Reel == 0 && award.Position.Row == 0 {
				// logic.GameLog.AddGameStep("贏得 Jackpot[%s], 金額: %d, 提前結束Lightning", award.JackpotType.String(), award.WinMoney)
				continue
			}
			logic.GameLog.AddGameStep("%d_%d:%s(Jackpot[%s])=%d", award.Position.Reel, award.Position.Row, award.Symbol.GetSymbol(), award.JackpotType.String(), award.WinMoney)
		}
	}
}

// PrintPayResult 紀錄獎勵
func (s *LightningStep) PrintPayResult(logic *SlotLogicBase) {
	for _, result := range s.PayResults {
		switch result.GetPayType() {
		case slot.PayTypeLightning:
			logic.GameLog.AddGameStep("重置 Lightning 次數(%d)", int(s.LastChance))
		case slot.PayTypeAddTimes:
			f, _ := strconv.ParseFloat(result.Rule.GetPayout(), 64)
			logic.GameLog.AddGameStep("Lightning 出現獎勵符號, 次數增加 %d 次", int(f))
		case slot.PayTypeJackpot:
			award := s.Awards[len(s.Awards)-1]
			logic.GameLog.AddGameStep("贏得 Jackpot[%s], 金額: %d, 提前結束Lightning", award.JackpotType.String(), award.WinMoney)
		}
	}
}

// LightningGame Lightning 模式遊戲
type LightningGame struct {
	Logic              *SlotLogicBase         `json:"-"`
	Steps              []*LightningStep       `json:"steps"` // 每次轉出的結果
	Origin             *protodata.RespinTable `json:"-"`
	OriginID           int32                  `json:"origin_id"`
	BetMoney           int64                  `json:"bet_money"`             // 下注金額
	BetRate            int64                  `json:"bet_rate"`              // 下注倍率
	ShowStep           int                    `json:"show_step"`             // 以表演完階段
	IsFreeGame         bool                   `json:"is_free_game"`          // 是否由FreeGame贏得
	SymbolGroupID      int32                  `json:"symbol_group_id"`       // 預設SymbolGroupID
	IsLimit            bool                   `json:"is_limit"`              // 是否有臨界值
	LimitSymbolGroupID int32                  `json:"limit_symbol_group_id"` // 當超過臨界時要轉的SymbolGroupID
	LimitPoint         int                    `json:"limit_point"`           // 重轉臨界點，超過換SymbolGroupID
	WinMoney           int64                  `json:"win_money"`             // 累計共贏金額
}

// NewLightningGame create LightningGame
func NewLightningGame(logic *SlotLogicBase, table *protodata.RespinTable, betMoney, betRate int64, isFreeGame bool) *LightningGame {
	g := new(LightningGame)
	g.Logic = logic
	g.IsFreeGame = isFreeGame
	g.Origin = table
	g.OriginID = int32(table.GetId())
	g.SymbolGroupID = int32(table.GetRespinValue())
	g.BetMoney = betMoney
	g.BetRate = betRate

	if table.GetIsLimit() == 1 {
		g.IsLimit = true
		if table.GetLimitType() == "changeDistribution" {
			strs := strings.Split(table.GetLimitValue(), ":")
			if f, err := strconv.ParseFloat(strs[0], 64); err == nil {
				g.LimitPoint = int(f)
			}
			if f, err := strconv.ParseFloat(strs[1], 64); err == nil {
				g.LimitSymbolGroupID = int32(f)
			}
		}
	}

	return g
}

// GetLastStep 取得最後一筆Step
func (g *LightningGame) GetLastStep() *LightningStep {
	if len(g.Steps) > 0 {
		return g.Steps[len(g.Steps)-1]
	}
	return nil
}

// AddStep 加入一個Step
func (g *LightningGame) AddStep(result *LightningStep) {
	if result == nil {
		return
	}
	g.Steps = append(g.Steps, result)
}

// Spin 轉出一個新輪軸結果加入
func (g *LightningGame) Spin() {

	last := g.GetLastStep()
	if last == nil {
		mlog.Error("查無前一次輪軸結果")
		return
	}

	var symbolGroupID int32
	if g.IsLimit && len(g.Steps) > g.LimitPoint {
		symbolGroupID = g.LimitSymbolGroupID
	} else {
		symbolGroupID = g.SymbolGroupID
	}

	last.LastChance--
	lastChance := last.LastChance

	// 建立全新的輪軸結果
	manager, _ := slot.NewGameResultManager(g.Logic.DefaultReelNum, g.Logic.DefaultRowNum)
	manager.RegisterAfterGenerateFunc(func(grm *slot.GameResultManager) {
		mlog.Debug("Extend Last GameResult")
		if last := g.GetLastStep(); last != nil {
			last.GameResultManager.RowMode()
			grm.RowMode()

			for rowID, row := range last.GameResultManager.GameResult {
				for reelID, symbol := range row {
					if g.IsGoldBall(symbol) {
						grm.GameResult[rowID][reelID] = symbol
					}
				}
			}
		}
	})
	manager.RegisterAfterGenerateFunc(func(grm *slot.GameResultManager) {
		mlog.Debug("Add Symbol")
		if len(last.NextAdd) == 0 {
			return
		}
		grm.RowMode()
		for _, symbol := range last.NextAdd {
			var inserted bool
			var retryLimit int = 100
			var retry int

			for ; !inserted; retry++ {
				//  骰一個ReelID, RowID, 檢查該位置是不是金球, 不是金球就把符號塞入

				randomReelID := rand.Intn(int(grm.Reel))
				randomRowID := rand.Intn(int(grm.Row))

				mlog.Debug("randomReelID: %d, randomRowID: %d, Retry: %d", randomReelID, randomRowID, retry)

				if !g.IsGoldBall(grm.GameResult[randomRowID][randomReelID]) {

					smb := *grm.GameResult[randomRowID][randomReelID]
					smb.Symbol = proto.String(symbol)
					// smb := &protodata.SlotSymbolDistribution{Id: proto.Uint32(uint32(100000 + i)), Symbol: proto.String(symbol)}

					grm.GameResult[randomRowID][randomReelID] = &smb

					mlog.Debug("加入符號 %s, ReelID: %d, RowID: %d, Retry: %d", symbol, randomReelID+1, randomRowID+1, retry)

					inserted = true
					break
				}

				if retry > retryLimit {
					break
				}
			}
		}
	})
	manager.GenerateGameResult(symbolGroupID, g.Logic.Proto.SlotSymbolDistribution)

	// 盤面檢查
	g.Logic.GameBetMoney = g.BetMoney
	g.Logic.GameBetRate = g.BetRate
	g.Logic.HandleAfterSpin(manager)

	var origin *protodata.RespinTable = last.Origin

	// 檢查是否有中獎
	_results := g.Logic.CalcGameResult(manager)
	results := make([]*slot.PayResult, 0)

	for _, result := range _results {
		switch slot.UseType(result.Rule.GetUseType()) {
		case slot.UseTypeMainGame:
			if g.IsFreeGame {
				continue
			}
		case slot.UseTypeFreeGame:
			if !g.IsFreeGame {
				continue
			}
		}
		results = append(results, result)
	}

	manager.RowMode()

	step := NewLightningStep(origin, last.Rule, manager, lastChance)
	step.PayResults = results

	g.CalcStepAward(step)
	g.CalcStepPayResult(step)

	g.WinMoney += step.WinMoney
	g.AddStep(step)
}

// IsGoldBall 檢查是否為金球
func (g *LightningGame) IsGoldBall(symbol *protodata.SlotSymbolDistribution) bool {

	for _, attr := range g.Logic.Proto.SymbolAttributes {
		if attr.GetIsGoldBall() == 1 && attr.GetSymbol() == symbol.GetSymbol() {
			return true
		}
	}

	return false
}

// CalcStepAward 計算此Step的獎勵金額
func (g *LightningGame) CalcStepAward(step *LightningStep) {

	step.GameResultManager.RowMode()

	// 把AccuMulate挑出來特別處理
	type AccuMulateObj struct {
		Position  slot.Position
		Symbol    *protodata.SlotSymbolDistribution
		Attribute *protodata.SymbolAttributes
	}
	var accuMulates []*AccuMulateObj

	for rowID, row := range step.GameResultManager.GameResult {
		for reelID, symbol := range row {
			// 檢查此位置是否計算過獎勵
			var isCalc bool
			for _, _step := range g.Steps {
				for _, award := range _step.Awards {
					if award.Position.Reel == int32(reelID+1) && award.Position.Row == int32(rowID+1) {
						isCalc = true
						break
					}
				}
			}
			if isCalc {
				continue
			}
			for _, attr := range g.Logic.Proto.SymbolAttributes {
				if attr.GetIsGoldBall() == 1 && attr.GetSymbol() == symbol.GetSymbol() {
					pos := slot.Position{Reel: int32(reelID + 1), Row: int32(rowID + 1)}

					switch slot.PayType(attr.GetGoldballType()) {
					case slot.PayTypeChip:
						f, _ := strconv.ParseFloat(attr.GetGoldballValue(), 64)
						win := g.BetMoney * int64(f)

						award := &LightningAward{
							WinMoney: win,
							Type:     slot.PayTypeChip,
							Symbol:   symbol,
							Position: pos,
						}

						step.Awards = append(step.Awards, award)
						step.WinMoney += win
					case slot.PayTypeChipXRate:
						f, _ := strconv.ParseFloat(attr.GetGoldballValue(), 64)

						win := g.BetMoney * g.BetRate * int64(f)

						award := &LightningAward{
							WinMoney: win,
							Type:     slot.PayTypeChipXRate,
							Symbol:   symbol,
							Position: pos,
						}

						step.Awards = append(step.Awards, award)
						step.WinMoney += win
					case slot.PayTypeJackpot:
						f, _ := strconv.ParseFloat(attr.GetGoldballValue(), 64)

						var jackpot *slot.Jackpot
						// 找 Jackpot
						for _, jp := range g.Logic.Jackpots {
							if jp.Origin.GetId() == int32(f) {
								jackpot = jp
							}
						}

						if jackpot != nil {

							win := jackpot.GetCurrentPoolByBet(g.BetMoney*g.BetRate, int64(g.Logic.MaxBetCount)*int64(g.Logic.MaxBetRate))

							award := &LightningAward{
								WinMoney:    win,
								Type:        slot.PayTypeJackpot,
								JackpotType: jackpot.Level,
								Symbol:      symbol,
								Position:    pos,
							}

							if defender.JPDefender != nil {
								defender.JPDefender.Award(win)
							}

							step.Awards = append(step.Awards, award)
							step.WinMoney += win
						}
					case slot.PayTypeAccuMulate:
						accuMulates = append(accuMulates, &AccuMulateObj{
							Position:  pos,
							Symbol:    symbol,
							Attribute: attr,
						})
					}
				}
			}
		}
	}

	for _, row := range accuMulates {
		if row.Symbol.GetSymbol() == "g7" {
			symbols := strings.Split(row.Attribute.GetGoldballValue(), ",")

			var win int64

			// 抓出已出現獎勵統計金額
			var _awards []*LightningAward
			for _, _step := range g.Steps {
				_awards = append(_awards, _step.Awards...)
			}
			// 本次Step的Award
			_awards = append(_awards, step.Awards...)

			for _, award := range _awards {
				for _, s := range symbols {
					if award.Symbol.GetSymbol() == s {
						win += award.WinMoney
						break
					}
				}
			}

			award := &LightningAward{
				WinMoney: win,
				Type:     slot.PayTypeAccuMulate,
				Symbol:   row.Symbol,
				Position: row.Position,
			}

			step.Awards = append(step.Awards, award)
			step.WinMoney += win
		}
	}

	for _, row := range accuMulates {
		if row.Symbol.GetSymbol() == "g8" {
			symbols := strings.Split(row.Attribute.GetGoldballValue(), ",")

			var win int64

			// 抓出已出現獎勵統計金額
			var _awards []*LightningAward
			for _, _step := range g.Steps {
				_awards = append(_awards, _step.Awards...)
			}
			// 本次Step的Award
			_awards = append(_awards, step.Awards...)

			for _, award := range _awards {
				for _, s := range symbols {
					if award.Symbol.GetSymbol() == s {
						win += award.WinMoney
						break
					}
				}
			}

			award := &LightningAward{
				WinMoney: win,
				Type:     slot.PayTypeAccuMulate,
				Symbol:   row.Symbol,
				Position: row.Position,
			}

			step.Awards = append(step.Awards, award)
			step.WinMoney += win
		}
	}
}

// CalcStepPayResult 計算此Step的非金球獎利
func (g *LightningGame) CalcStepPayResult(step *LightningStep) {

	var rules []*slot.PayResult

	for _, result := range step.PayResults {
		switch result.GetPayType() {
		case slot.PayTypeAddTimes:
			f, _ := strconv.ParseFloat(result.Rule.GetPayout(), 64)
			step.LastChance += int(f)
			rules = append(rules, result)
		case slot.PayTypeJackpot:

			jackpot, win := g.Logic.WinJackpot(result.Rule)

			if jackpot != nil {
				award := &LightningAward{
					WinMoney:    win,
					Type:        slot.PayTypeJackpot,
					JackpotType: jackpot.Level,
				}
				step.Awards = append(step.Awards, award)
				step.WinMoney += win
				step.LastChance = 0
			}

			rules = append(rules, result)
		case slot.PayTypeLightning:

			if step.Rule != nil && step.Rule == result.Rule {
				continue
			}

			// 取得贏得Lightning的GroupID
			f, err := strconv.ParseFloat(result.Rule.GetPayout(), 64)
			if err != nil {
				mlog.Warn("SlotPayTable.Payout解析錯誤，%v", err)
				continue
			}

			respinGroupID := int32(f)

			if respinGroupID != step.Origin.GetGroupId() {
				continue
			}

			step.LastChance = int(step.Origin.GetRespinCount())
			step.Rule = result.Rule
			step.RuleID = result.Rule.GetId()
			rules = append(rules, result)
		}
	}

	step.PayResults = rules
}
