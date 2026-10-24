package baselogic

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	apipb "example.invalid/ServerTeam/Servers/internal/api-service/apisvc"

	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/prob"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/defender"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// Bet 判斷下注是否
func (logic *SlotLogicBase) Bet(user *gamelogic.GameUser, money int64, rate int64) (ok bool, retcode int32, betmoney int64) {

	// 檢查狀態是否可下注
	if logic.GameState != slot.GameStateBet {
		return false, slot.BetStateStepError, logic.GamePayMoney
	}

	// 檢查下注金額是否合法
	var moneyOK, rateOK bool
	for _, m := range logic.BetCounts {
		if int64(m) == money {
			moneyOK = true
			break
		}
	}

	switch logic.PayLineMode {
	case slot.PayLineModeAll:
		for _, r := range logic.BetRates {
			if int64(r) == rate {
				rateOK = true
				break
			}
		}
	case slot.PayLineModeFixed:
		rate = int64(len(logic.Proto.SlotPayLine))
		rateOK = true
	}

	if !moneyOK || !rateOK {
		mlog.Warn("玩家 %v 非法下注，金額: %d，倍率: %d", user, money, rate)
		return false, slot.BetStateIllegal, logic.GamePayMoney
	}

	logic.GameBetMoney = money
	logic.GameBetRate = rate
	logic.GamePayMoney = money * rate

	if logic.FreeSpin != nil && logic.FreeSpin.Count > 0 {
		logic.GameBetMoney = logic.FreeSpin.BetMoney
		logic.GamePayMoney = 0
		logic.FreeSpin.Count--
		logic.FreeSpin.Used++
	}

	// 檢查最後金額是否為負數
	if logic.GamePayMoney < 0 {
		mlog.Warn("玩家 %v 非法下注，付款金額小於零 %d", user.UserID, logic.GamePayMoney)
		return false, slot.BetStateIllegal, logic.GamePayMoney
	}

	// 檢查玩家金額
	if user.Score < logic.GamePayMoney {
		return false, slot.BetStateMoneyNotEnough, logic.GamePayMoney
	}

	// 檢查最低攜帶金額
	if user.Score < logic.MinCarryMoney && logic.FreeSpin == nil {
		return false, slot.BetStateMoneyLessMin, logic.GamePayMoney
	}

	if user.IsRobot() {
		return true, slot.BetStateSuccess, logic.GamePayMoney
	}

	if err := app.Context.Invoke(func(c apipb.GameBetServiceClient, pm *player.Manager) error {
		p := pm.GetPlayerByUserID(user.UserID)
		req := new(apipb.PlaceBetReq)
		req.AccessToken = [REDACTED]
		req.ServerID = logic.TableManager.Room.ServerID
		req.TradeID = user.TradeID
		req.BetAmount = logic.GamePayMoney

		if logic.FreeSpin != nil {
			req.Memo = fmt.Sprintf("FreeSpin(%d)", logic.FreeSpin.Count)
		}

		res, err := c.PlaceBet(context.Background(), req)
		if err != nil {
			return err
		}

		if logic.GamePayMoney > 0 && logic.FreeSpin != nil {
			user.Score = res.AfterAmount + logic.GamePayMoney
			logic.BoardcastUserStateChange(user)
		}
		logic.TransactionNo = res.TransactionNo

		return nil
	}); err != nil {
		mlog.Error("Place Bet Error, %v", err)
		logic.KickUserFunc(user.UserID, "System Error ...")
		return false, slot.BetStateMoneyNotEnough, logic.GamePayMoney
	}

	return true, slot.BetStateSuccess, logic.GamePayMoney
}

// PrepareSpin 取得SymbolGroupID, BgGroupID, 累積JackpotPool
func (logic *SlotLogicBase) PrepareSpin() (symbolGroupID, bgGroupID, reelNum, rowNum int32) {

	// 檢查特殊模式是否過期
	logic.CheckModeExpire()

	if logic.FreeSpin != nil {
		symbolGroupID = logic.FreeSpin.SymbolGroupID
		bgGroupID = logic.FreeSpin.BgGroupID
		reelNum = int32(logic.FreeSpin.Origin.GetFreeReelNum())
		rowNum = int32(logic.FreeSpin.Origin.GetFreeRowNum())

		logic.GameLog.AddGameStep("玩家 %d; 底注: %d, 倍率: %d, 免費遊戲剩餘次數: %d", logic.GameUser.UserID, logic.FreeSpin.BetMoney, logic.FreeSpin.BetRate, logic.FreeSpin.Count)
	} else {
		symbolGroupID = logic.DefaultSymbolGroupID
		bgGroupID = logic.DefaultBgGroupID
		reelNum = logic.DefaultReelNum
		rowNum = logic.DefaultRowNum

		// 依照下注倍率選擇MainGameTable
		for _, table := range logic.Proto.MainGameTable {
			f, _ := strconv.ParseFloat(table.GetConditionValue(), 64)
			if int64(f) == logic.GameBetRate {
				symbolGroupID = int32(table.GetSymbolGroupId())
				bgGroupID = int32(table.GetGoldReelGroupId())
				reelNum = int32(table.GetReelNum())
				rowNum = int32(table.GetRowNum())
				break
			}
		}

		logic.GameLog.AddGameStep("玩家 %d; 底注: %d, 倍率: %d, 花費: %d", logic.GameUser.UserID, logic.GameBetMoney, logic.GameBetRate, logic.GamePayMoney)

		// 檢查是否有註冊JP防衛器
		if defender.JPDefender != nil {
			inc := defender.JPDefender.Increase(logic.GamePayMoney)
			mlog.Debug("Table %d, Increase defender.JPDefender %d", logic.TableNo(), inc)
		}
		logic.IncreaseJackpotPool()

		logic.SendJackpotMessageToUser(logic.GameUser)
	}
	return
}

// CalcGameResult 計算遊戲結果
func (logic *SlotLogicBase) CalcGameResult(grm *slot.GameResultManager) []*slot.PayResult {

	results := make([]*slot.PayResult, 0)

	// 依照SlotPayTable檢查中了哪些
	for _, rule := range logic.Proto.SlotPayTable {
		match, result := logic.PaylineMatcher.IsMatch(grm, rule)
		if !match {
			continue
		}
		results = append(results, result)
	}

	// 檢查 Background 獎勵
	for _, rule := range logic.Proto.SlotPayTable {
		match, result := logic.BackgroundMatcher.IsMatch(grm, rule)
		if !match {
			continue
		}
		results = append(results, result)
	}

	// 所有賠線中獎規則
	lineRewareMap := make(map[slot.LinkType]map[uint32]*slot.PayResult)
	freeRewareMap := make(map[slot.LinkType]map[*protodata.FreeSpinTable][]*slot.PayResult)

	// 依照連線個數做排序，多到少
	sort.Slice(results, func(i, j int) bool {
		return results[i].PayCount > results[j].PayCount
	})

	// 過濾已經出現過的連線，並取贏多的那筆規則
	for _, result := range results {
		linkType := slot.LinkType(result.Rule.GetLinkType())
		if !result.IsMatch || result.Rule == nil {
			continue
		}
		if result.Rule.GetPayType() == slot.PayTypeJackpot.String() {
			for _, line := range result.SymbolPayLines {
				if _, ok := lineRewareMap[linkType]; !ok {
					lineRewareMap[linkType] = make(map[uint32]*slot.PayResult)
				}
				lineRewareMap[linkType][line.PayLine.Origin.GetLineId()] = result
			}
		}
		if result.Rule.GetPayType() == slot.PayTypeFreeSpin.String() {
			if _, ok := freeRewareMap[linkType]; !ok {
				freeRewareMap[linkType] = make(map[*protodata.FreeSpinTable][]*slot.PayResult)
			}
			var fg *protodata.FreeSpinTable
			f1, _ := strconv.ParseFloat(result.Rule.GetPayout(), 64)
			reware1 := int32(f1)
			for _, fgtb := range logic.Proto.FreeSpinTable {
				if fgtb.GetGroupId() == reware1 {
					fg = fgtb
				}
			}
			if fg == nil {
				continue
			}
			freeRewareMap[linkType][fg] = append(freeRewareMap[linkType][fg], result)
		}
		if result.Rule.GetPayType() != slot.PayTypeChip.String() {
			continue
		}

		for _, line := range result.SymbolPayLines {
			if _, ok := lineRewareMap[linkType]; !ok {
				lineRewareMap[linkType] = make(map[uint32]*slot.PayResult)
			}

			if _, ok := lineRewareMap[linkType][line.PayLine.Origin.GetLineId()]; !ok {
				lineRewareMap[linkType][line.PayLine.Origin.GetLineId()] = result
				continue
			}

			orgResult := lineRewareMap[linkType][line.PayLine.Origin.GetLineId()]

			symbols1 := strings.Split(result.Rule.GetSymbol(), ",")
			symbols2 := strings.Split(orgResult.Rule.GetSymbol(), ",")

			if len(symbols1) != len(symbols2) {
				if len(symbols1) < len(symbols2) {
					orgResult.RemoveLine(line.PayLine.Origin.GetLineId())
				} else {
					result.RemoveLine(line.PayLine.Origin.GetLineId())
				}
				continue
			}

			f1, _ := strconv.ParseFloat(result.Rule.GetPayout(), 64)
			reware1 := int64(f1)

			f2, _ := strconv.ParseFloat(orgResult.Rule.GetPayout(), 64)
			reware2 := int64(f2)

			if reware1 > reware2 {
				orgResult.RemoveLine(line.PayLine.Origin.GetLineId())
			} else {
				result.RemoveLine(line.PayLine.Origin.GetLineId())
			}
		}
	}

	tmp := make([]*slot.PayResult, 0, len(results))
	// 整理沒有連線(被移除後)的規則
	for _, result := range results {
		if !result.IsMatch || result.Rule == nil {
			continue
		}
		if result.Rule.GetPayType() == slot.PayTypeFreeSpin.String() {
			continue
		}
		if result.Rule.GetPayType() != slot.PayTypeChip.String() {
			tmp = append(tmp, result)
			continue
		}
		if result.WayCount > 0 {
			tmp = append(tmp, result)
		}
	}

	if len(freeRewareMap) > 0 {
		var bestfree *protodata.FreeSpinTable
		for _, fgs := range freeRewareMap {
			for fg, rs := range fgs {
				if bestfree == nil {
					bestfree = fg
				} else if bestfree.GetFreeSpinCount() < fg.GetFreeSpinCount() {
					bestfree = fg
				}
				sort.Slice(rs, func(i, j int) bool {
					return rs[i].WayCount < rs[j].WayCount
				})
			}
		}

		for _, fgs := range freeRewareMap {
			if len(fgs[bestfree]) == 0 {
				continue
			}
			tmp = append(tmp, fgs[bestfree][0])
		}
	}

	return tmp
}

// CalcGameAward 計算遊戲獎勵
func (logic *SlotLogicBase) CalcGameAward(results []*slot.PayResult, betMoney, betRate int64) int64 {
	return logic.CalcGameAwardByMultiType(results, betMoney, betRate, 1, slot.MultiTypeNone, nil)
}

// CalcGameAwardByMultiType 計算遊戲獎勵
func (logic *SlotLogicBase) CalcGameAwardByMultiType(results []*slot.PayResult, betMoney, betRate, multi int64, multiType slot.MultiType, multiDetail map[slot.Position]int32) int64 {

	var winMoney int64

	for _, result := range results {

		result.MultiMode = logic.MultiMode

		if !result.IsMatch || result.Rule == nil {
			continue
		}

		switch result.GetPayType() {
		case slot.PayTypeChip, slot.PayTypeChipXRate:
			result.Multi = int32(multi)
			result.MultiType = multiType
			result.MultiDetail = multiDetail

			// 計算贏得的金額
			base, win := result.CalcWinMoney(betMoney, betRate)
			result.WinMoney = win
			winMoney += win

			logic.LogPayResultChip(result, base, win)

		case slot.PayTypeBonus:

			bonus := logic.WinBonus(result.Rule)

			if bonus == nil {
				mlog.Debug("No Bonus for [%s]", result.Rule.GetPayType())
				continue
			}

			switch slot.PayType(bonus.GetRewardType()) {
			default:
				continue
			case slot.PayTypeChip, slot.PayTypeChipXRate:
				win, rate := logic.CalcChipMoney(slot.PayType(bonus.GetRewardType()), bonus.GetRewardValue())
				result.WinMoney = win
				if win == 0 {
					continue
				}
				winMoney += win

				var base int64 = rate * logic.GameBetMoney
				if slot.PayType(bonus.GetRewardType()) == slot.PayTypeChipXRate {
					base *= logic.GameBetRate
				}

				logic.LogPayResultChip(result, base, win)
			}
		}
	}

	return winMoney
}

// CalcChipMoney 計算Chip or ChipXRate的錢
func (logic *SlotLogicBase) CalcChipMoney(paytype slot.PayType, value string) (win int64, rate int64) {
	switch paytype {
	default:
		return 0, 0
	case slot.PayTypeChip, slot.PayTypeChipXRate:
		f, _ := strconv.ParseFloat(value, 64)

		var betMoney int64 = logic.GameBetMoney
		var betRate int64 = 1
		if paytype == slot.PayTypeChipXRate {
			betRate = logic.GameBetRate
		}

		return betMoney * int64(f) * int64(betRate), int64(f)
	}
}

// CalcJackpot 傳入中獎的Jackpot處理
// 1. 計算Jackpot金額
// 2. 更新Jackpot金額
func (logic *SlotLogicBase) CalcJackpot(jackpot *slot.Jackpot) (winMoney int64) {
	return logic.CalcJackpotByMoney(jackpot, logic.GameBetMoney*logic.GameBetRate, int64(logic.MaxBetCount)*int64(logic.MaxBetRate))
}

// CalcJackpotByMoney 傳入中獎的Jackpot處理
// 1. 計算Jackpot金額
// 2. 更新Jackpot金額
func (logic *SlotLogicBase) CalcJackpotByMoney(jackpot *slot.Jackpot, currentBet, maxBet int64) (winMoney int64) {
	win := jackpot.GetCurrentPoolByBet(currentBet, maxBet)

	jackpot.AwardTotal += win

	if jackpot.Origin.GetIsProgress() == 1 {
		if jackpot.IsNeedReset() {
			// logic.GameLog.AddGameStep("重置獎池 %s", jackpot.Level.String())
			mlog.Debug("重置獎池[%s]", jackpot.Level.String())

			logic.UpdateJackpotPoolNoLink(0, -jackpot.TotalContribution, -jackpot.AwardTotal+win, jackpot.Level)
			jackpot.Reset()
		} else {
			// logic.GameLog.AddGameStep("更新獎池 %s 出獎金額: %d", jackpot.Level.String(), win)
			mlog.Debug("更新獎池[%s] 出獎金額: %d", jackpot.Level.String(), win)

			logic.UpdateJackpotPoolNoLink(0, 0, win, jackpot.Level)
		}
	}

	return win
}

// WinFreeSpin 處理贏得FreeSpin
func (logic *SlotLogicBase) WinFreeSpin(rule *protodata.SlotPayTable, log bool) *slot.FreeSpin {
	// 取得贏得FreeSpin的GroupID
	f, err := strconv.ParseFloat(rule.GetPayout(), 64)
	if err != nil {
		mlog.Debug("Parse [%s] error %v", rule.GetPayType(), err)
		return nil
	}

	freeSpinGroupID := int32(f)

	var freeSpinTables []*protodata.FreeSpinTable
	var freeSpinTable *protodata.FreeSpinTable

	// 取得FreeSpinTable
	for _, table := range logic.Proto.FreeSpinTable {
		if table.GetGroupId() != freeSpinGroupID {
			continue
		}
		freeSpinTables = append(freeSpinTables, table)
	}

	if len(freeSpinTables) == 0 {
		mlog.Warn("查無FreeSpinTable, FreeSpinGroupID: %d", freeSpinGroupID)
		return nil
	} else if len(freeSpinTables) == 1 {
		freeSpinTable = freeSpinTables[0]
	} else {
		// 依照倍率取得正確FreeSpinTable
		for _, table := range freeSpinTables {
			f, _ := strconv.ParseFloat(table.GetConditionValue(), 64)
			if int64(f) == logic.GameBetRate {
				freeSpinTable = table
				break
			}
		}
	}

	if freeSpinTable == nil {
		mlog.Warn("查無FreeSpinTable, FreeSpinGroupID: %d", freeSpinGroupID)
		return nil
	}

	// 檢查本身是否是FreeSpin
	if logic.FreeSpin != nil {
		if logic.FreeSpin.Origin != freeSpinTable {
			mlog.Warn("贏得的FreeSpin與原本的FreeSpin不同")
			return nil
		}

		logic.FreeSpin.Count += int(freeSpinTable.GetFreeSpinCount())

	} else {
		logic.FreeSpin = slot.NewFreeSpin(freeSpinTable, logic.Proto.SlotSymbolDistribution, logic.GameBetMoney, logic.GameBetRate)
	}

	if log {
		logic.GameLog.AddGameStep("開獎符號: %s:%s:%d, 贏得 FreeSpin(%d), Count: %d, 剩餘次數: %d", rule.GetSymbol(), rule.GetLinkType(), rule.GetCount(), logic.FreeSpin.Origin.GetFreeSymbolGroupId(), logic.FreeSpin.Origin.GetFreeSpinCount(), logic.FreeSpin.Count)
	}

	return logic.FreeSpin
}

// WinMysteryJP 處理贏得MysteryJP(機率中JP)該做的事
func (logic *SlotLogicBase) WinMysteryJP(rule *protodata.SlotPayTable) (jackpot *slot.Jackpot, win int64) {

	var jpid int32
	var winprop, nowinprop, denominator int64
	var err error
	switch slot.PayType(rule.GetPayType()) {
	case slot.PayTypeMysteryJP:
		jpid, winprop, denominator, err = slot.ParsePayoutMysteryJP(rule.GetPayout())
	case slot.PayTypeMysteryJPR:
		jpid, winprop, denominator, err = slot.ParsePayoutMysteryJPR(rule.GetPayout(), logic.GameBetRate)
	}

	if err != nil {
		mlog.Debug("Parse [%s] error %v", rule.GetPayType(), err)
		return nil, 0
	}

	nowinprop = denominator - winprop

	mlog.Debug("Prop => Win: %d,NoWin: %d", winprop, nowinprop)

	r := new(prob.Random)

	if winprop > 0 {
		r.Add(true, int(winprop))
	}

	if nowinprop > 0 {
		r.Add(false, int(nowinprop))
	}

	t := r.Choose().(bool)

	if !t {
		return nil, 0
	}

	// 找 Jackpot
	for _, jp := range logic.Jackpots {
		if jp.Origin.GetId() == int32(jpid) {
			jackpot = jp
		}
	}

	if jackpot == nil {
		return nil, 0
	}

	// 檢查是否有註冊JP防衛器
	if defender.JPDefender != nil {
		f := func() int64 {
			return jackpot.GetCurrentPoolByBet(logic.GameBetMoney*int64(logic.GameBetRate), int64(logic.MaxBetCount*logic.MaxBetRate))
		}
		if !defender.JPDefender.Eligible(jackpot, f) {
			example.invalid("[JPDefender] Talbe(%d) Illegible Jackpot[%v], JPDefender Pool: %d, Win: %d", logic.TableNo(), jackpot.Level, defender.JPDefender.Pool(), f())
			return nil, 0
		}
		defender.JPDefender.Award(f())
	}

	win = logic.CalcJackpot(jackpot)

	logic.GameLog.AddGameStep("開獎符號: %s:%s:%d, 贏得 Jackpot[%s], 獎金: %d", rule.GetSymbol(), rule.GetLinkType(), rule.GetCount(), jackpot.Level.String(), win)

	return jackpot, win
}

// WinMysteryJPG 處理贏得MysteryJPGR該做的事
func (logic *SlotLogicBase) WinMysteryJPG(rule *protodata.SlotPayTable) *slot.JackpotGame {

	var linkprop, nolinkprop, nowinprop, denominator int64
	var err error
	switch slot.PayType(rule.GetPayType()) {
	case slot.PayTypeMysteryJPG:
		linkprop, nolinkprop, denominator, err = slot.ParsePayoutMysteryJPG(rule.GetPayout())
	case slot.PayTypeMysteryJPGR:
		linkprop, nolinkprop, denominator, err = slot.ParsePayoutMysteryJPGR(rule.GetPayout(), logic.GameBetRate)
	}

	if err != nil {
		mlog.Debug("Parse [%s] error %v", rule.GetPayType(), err)
		return nil
	}

	// 連線水池機率要乘上總下注分數
	linkprop = linkprop * logic.GameBetMoney * logic.GameBetRate
	nowinprop = denominator - nolinkprop - linkprop

	type PropType int32
	const (
		LinkType PropType = iota
		NoLinkType
		NoWinType
	)

	// linkprop need x total bet
	r := new(prob.Random)

	mlog.Debug("LinkJP: %d, NoLinkJP: %d, NoWin: %d", linkprop, nolinkprop, nowinprop)

	if linkprop > 0 {
		r.Add(LinkType, int(linkprop))
	}

	if nolinkprop > 0 {
		r.Add(NoLinkType, int(nolinkprop))
	}

	if nowinprop > 0 {
		r.Add(NoWinType, int(nowinprop))
	}

	t := r.Choose().(PropType)

	if !conf.AnalysisMode {
		mlog.Debug("WinMysteryJPG, LinkProp(%v): %d, NoLinkProp(%v): %d, NoWinProp(%v): %d, 骰出: %v", LinkType, linkprop, NoLinkType, nolinkprop, NoWinType, nowinprop, t)
		mlog.Debug("%+v", r)
	}

	var jpg *slot.JackpotGame

	switch t {
	case NoWinType:
		mlog.Debug("JackpotGame 未中獎")
		return nil
	case LinkType:
		// 取出連線Jackpot獎池
		var jackpots, _jackpots []*slot.Jackpot
		for _, jp := range logic.Jackpots {
			if jp.Origin.GetIsMultiLevel() != 1 {
				continue
			}
			if jp.Origin.GetIsLinked() == 1 {
				// 檢查是否有註冊JP防衛器
				if defender.JPDefender != nil {
					f := func() int64 {
						return jp.GetCurrentPoolByBet(logic.GameBetMoney*int64(logic.GameBetRate), int64(logic.MaxBetCount*logic.MaxBetRate))
					}
					if !defender.JPDefender.Eligible(jp, f) {
						mlog.Debug("[JPDefender] Illegible Jackpot[%v] Skip Option ...", jp.Level)
						_jackpots = append(_jackpots, jp)
						continue
					}
				}
				jackpots = append(jackpots, jp)
			} else {
				_jackpots = append(_jackpots, jp)
			}
		}
		if len(jackpots) == 0 {
			mlog.Debug("無連線型Jackpot, 跳過 Jackpot Game")
			return nil
		}

		jpg = slot.NewJackpotGame(jackpots, logic.MaxBetRate, logic.MaxBetCount, int32(logic.GameBetRate), logic.GameBetMoney)

		jpg.WinJackpot = jpg.RandomJackpot()

		jpg.Jackpots = append(jpg.Jackpots, _jackpots...)

		jpg.Process = jpg.GenerateProcess()
	case NoLinkType:
		// 取出連線Jackpot獎池
		var jackpots, _jackpots []*slot.Jackpot

		for _, jp := range logic.Jackpots {
			if jp.Origin.GetIsMultiLevel() != 1 {
				continue
			}
			if jp.Origin.GetIsLinked() == 0 {
				// 檢查是否有註冊JP防衛器
				if defender.JPDefender != nil {
					f := func() int64 {
						return jp.GetCurrentPoolByBet(logic.GameBetMoney*int64(logic.GameBetRate), int64(logic.MaxBetCount*logic.MaxBetRate))
					}
					if !defender.JPDefender.Eligible(jp, f) {
						mlog.Debug("[JPDefender] Illegible Jackpot[%v] Skip Option ...", jp.Level)
						_jackpots = append(_jackpots, jp)
						continue
					}
				}
				jackpots = append(jackpots, jp)
			} else {
				_jackpots = append(_jackpots, jp)
			}
		}

		if len(jackpots) == 0 {
			mlog.Debug("無非連線型Jackpot, 跳過 Jackpot Game")
			return nil
		}

		jpg = slot.NewJackpotGame(jackpots, logic.MaxBetRate, logic.MaxBetCount, int32(logic.GameBetRate), logic.GameBetMoney)

		jpg.WinJackpot = jpg.RandomJackpot()

		jpg.Jackpots = append(jpg.Jackpots, _jackpots...)

		jpg.Process = jpg.GenerateProcess()
	}

	if jpg == nil {
		return nil
	}

	// 檢查是否有註冊JP防衛器
	if defender.JPDefender != nil {
		f := func() int64 {
			return jpg.WinJackpot.GetCurrentPoolByBet(logic.GameBetMoney*int64(logic.GameBetRate), int64(logic.MaxBetCount*logic.MaxBetRate))
		}
		if !defender.JPDefender.Eligible(jpg.WinJackpot, f) {
			example.invalid("[JPDefender] Talbe(%d) Illegible Jackpot[%v], JPDefender Pool: %d, Win: %d", logic.TableNo(), jpg.WinJackpot.Level, defender.JPDefender.Pool(), f())
			return nil
		}
		defender.JPDefender.Award(f())
	}

	jpg.WinMoney = logic.CalcJackpot(jpg.WinJackpot)
	jpg.Origin = rule

	logic.GameLog.AddGameStep("開獎符號: %s:%s:%d, 贏得 JackpotGame[%s], 獎金: %d, 將於翻牌結束後結算", rule.GetSymbol(), rule.GetLinkType(), rule.GetCount(), jpg.WinJackpot.Level.String(), jpg.WinMoney)

	return jpg
}

// WinSelection 處理贏得Selection該做的事
func (logic *SlotLogicBase) WinSelection(rule *protodata.SlotPayTable) *slot.SelectChance {
	selection, err := slot.NewSelectChance(logic.Proto, logic.GameResultManager, rule, logic.GameBetMoney, logic.GameBetRate)
	if err != nil {
		mlog.Error("建立選擇遊戲錯誤. %v", err)
		return nil
	}

	var opts []string
	// 依照本次下注倍率選擇該次選擇遊戲選項
	for _, row := range selection.Options {
		opts = append(opts, row.String())
	}

	logic.GameLog.AddGameStep("開獎符號: %s:%s:%d, 贏得選擇機會(%d), 選項: %s", rule.GetSymbol(), rule.GetLinkType(), rule.GetCount(), selection.GroupID, strings.Join(opts, "|"))

	return selection
}

// WinJackpot 處理贏得Jackpot該做的事
// 1. 找出對應Jackpot
// 2. 計算Jackpot金額
// 3. 更新Jackpot金額
// 4. 紀錄開獎符號
func (logic *SlotLogicBase) WinJackpot(rule *protodata.SlotPayTable) (jackpot *slot.Jackpot, winMoney int64) {
	f, _ := strconv.ParseFloat(rule.GetPayout(), 64)

	// 找 Jackpot
	for _, jp := range logic.Jackpots {
		if jp.Origin.GetId() == int32(f) {
			jackpot = jp
		}
	}

	if jackpot == nil {
		return nil, 0
	}

	win := logic.CalcJackpot(jackpot)

	if defender.JPDefender != nil {
		defender.JPDefender.Award(win)
	}

	logic.GameLog.AddGameStep("開獎符號: %s:%s:%d, 贏得 Jackpot[%s] 獎金: %d", rule.GetSymbol(), rule.GetLinkType(), rule.GetCount(), jackpot.Level.String(), win)

	return jackpot, win
}

// WinBonus 贏得Bonus選擇一筆
func (logic *SlotLogicBase) WinBonus(rule *protodata.SlotPayTable) (bonus *protodata.BonusTable) {

	f, err := strconv.ParseFloat(rule.GetPayout(), 64)
	if err != nil {
		mlog.Debug("Parse [%s] error %v", rule.GetPayType(), err)
		return nil
	}

	var bonuies []*protodata.BonusTable
	for _, row := range logic.Proto.BonusTable {
		// Filter by GroupID
		if row.GetGroupId() != int32(f) {
			continue
		}

		// Filter by Condition
		switch row.GetCondition() {
		default:
		case slot.ConditionRate:
			f, err := strconv.ParseFloat(row.GetConditionValue(), 32)
			if err != nil {
				mlog.Debug("Parse BonusTable ConditionValue error %v", err)
				continue
			}
			if logic.GameBetRate != int64(f) {
				continue
			}
		case slot.ConditionAmount:
		}

		bonuies = append(bonuies, row)
	}

	if len(bonuies) == 0 {
		mlog.Debug("No Bonus for [%s], GroupID: %d", rule.GetPayType(), int32(f))
		return nil
	}

	r := new(prob.Random)
	for _, row := range bonuies {
		r.Add(row, int(row.GetPerProb()))
	}

	bonus = r.Choose().(*protodata.BonusTable)
	return
}

// WinPickGame 贏得PickGame
func (logic *SlotLogicBase) WinPickGame(rule *protodata.SlotPayTable, size int32, log bool) (game *slot.PickGame) {

	game = BuildPickGame(rule, logic.Proto, size, logic.GameBetMoney, logic.GameBetRate)

	f, _ := strconv.ParseFloat(rule.GetPayout(), 64)

	if log {
		logic.GameLog.AddGameStep("開獎符號: %s:%s:%d, 贏得 PickGame[%d] 獎金: %d", rule.GetSymbol(), rule.GetLinkType(), rule.GetCount(), int(f), game.WinMoney)
	}

	return
}

// WinJackpotGameChance 贏得 JackpotGame 機會
func (logic *SlotLogicBase) WinJackpotGameChance(rule *protodata.SlotPayTable, log bool) *slot.JackpotGameChance {
	c, err := slot.NewJackpotGameChance(logic.GameResultManager, rule, logic.GameBetMoney, logic.GameBetRate)
	if err != nil {
		mlog.Warn("建立 Jackpot Game 機會錯誤 %v", err)
		return nil
	}

	if log {
		logic.GameLog.AddGameStep("開獎符號: %s:%s:%d, 贏得 JackpotGame, 結束後將決定獎項, 剩餘 %d 次", rule.GetSymbol(), rule.GetLinkType(), rule.GetCount(), len(logic.JackpotGameChances)+1)
	}

	return c
}

// WinLightningGame 贏得LightningGame
func (logic *SlotLogicBase) WinLightningGame(rule *protodata.SlotPayTable) (game *LightningGame) {
	// 取得贏得Lightning的GroupID
	f, err := strconv.ParseFloat(rule.GetPayout(), 64)
	if err != nil {
		mlog.Warn("SlotPayTable.Payout解析錯誤，%v", err)
		return nil
	}

	respinGroupID := int32(f)

	// 取得RespinTable
	var respinTable *protodata.RespinTable

	// 取得RespinTable
	for _, table := range logic.Proto.RespinTable {
		if table.GetGroupId() != respinGroupID {
			continue
		}
		respinTable = table
	}

	if respinTable == nil {
		mlog.Warn("查無RespinTable, RespinGroupID: %d", respinGroupID)
		return nil
	}

	game = NewLightningGame(logic, respinTable, logic.GameBetMoney, logic.GameBetRate, logic.IsFreeSpin())

	step := NewLightningStep(respinTable, rule, logic.GameResultManager, int(respinTable.GetRespinCount()))
	game.CalcStepAward(step)
	game.AddStep(step)
	game.WinMoney += step.WinMoney

	logic.GameLog.AddGameStep("開獎符號: %s:%s:%d, 贏得Lightning機會(%d)", rule.GetSymbol(), rule.GetLinkType(), rule.GetCount(), respinTable.GetRespinValue())
	logic.GameLog.AddGameStep("原始牌面獎勵:")
	step.PrintAward(logic)
	logic.GameLog.AddGameStep("Lightning 剩餘 %d 次, 累計金額: %d, 將於Lightning結束後計算", step.LastChance, game.WinMoney)

	return game
}

// WinFreeTimes 贏得免費遊戲次數
func (logic *SlotLogicBase) WinFreeTimes(rule *protodata.SlotPayTable, fg *slot.FreeSpin) int {
	if fg == nil {
		mlog.Error("增加免費遊戲次數錯誤, 沒有免費遊戲")
		return 0
	}

	f, _ := strconv.ParseFloat(rule.GetPayout(), 64)

	fg.Count += int(f)

	logic.GameLog.AddGameStep("開獎符號: %s:%s:%d, 贏得免費遊戲次數, 增加 %d 次, 剩餘 %d 次", rule.GetSymbol(), rule.GetLinkType(), rule.GetCount(), int(f), fg.Count)

	return int(f)
}

// IsGoldBall 判斷符號是否為金球
func (logic *SlotLogicBase) IsGoldBall(symbol string) bool {
	for _, attr := range logic.Proto.SymbolAttributes {
		if attr.GetIsGoldBall() == 1 && attr.GetSymbol() == symbol {
			return true
		}
	}
	return false
}
