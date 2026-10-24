package slot

import (
	"fmt"
	"strconv"
	"strings"
)

// ParsePayoutMysteryJP 解析MysteryJP(機率中Jackpot)的Payout內容
// @Argument
//   * pattern 		Payout文字內容	{JackpotID}|{機率}:{分母} ex: 7|100:10000
// @Return
//   * jackpotID 	應中獎的JackpotID
//   * prop			中獎機率
//   * denominator	機率分母
func ParsePayoutMysteryJP(pattern string) (jackpotID int32, prop, denominator int64, err error) {
	pattern = strings.TrimSpace(pattern)

	strs := strings.Split(pattern, "|")

	if len(strs) != 2 {
		return 0, 0, 0, fmt.Errorf("Payout format error")
	}

	if val, err := strconv.ParseFloat(strs[0], 32); err == nil {
		jackpotID = int32(val)
	} else {
		return 0, 0, 0, fmt.Errorf("Payout format error. JackpoID is invalid. %v", err)
	}

	data := strings.Split(strs[1], ":")

	if len(data) != 2 {
		return 0, 0, 0, fmt.Errorf("Payout format error. prop & denominator error")
	}

	if val, err := strconv.ParseFloat(data[0], 64); err == nil {
		prop = int64(val)
	} else {
		return 0, 0, 0, fmt.Errorf("Payout format error. prop is invalid. %v", err)
	}

	if val, err := strconv.ParseFloat(data[1], 64); err == nil {
		denominator = int64(val)
	} else {
		return 0, 0, 0, fmt.Errorf("Payout format error. denominator is invalid. %v", err)
	}

	return
}

// ParsePayoutMysteryJPR 解析MysteryJPR(依照不同倍率機率中Jackpot)的Payout內容
// @Argument
//   * pattern 		Payout文字內容	{JackpotID}|{下注倍率}:{機率}:{分母},{下注倍率}:{機率}:{分母}... ex: 7|8:100:10000,18:200:10000,38:300:10000,68:400:10000,88:500:10000
//   * betRate 		當前下注倍率
// @Return
//   * jackpotID 	應中獎的JackpotID
//   * prop			中獎機率
//   * denominator	機率分母
func ParsePayoutMysteryJPR(pattern string, betRate int64) (jackpotID int32, prop, denominator int64, err error) {
	pattern = strings.TrimSpace(pattern)

	strs := strings.Split(pattern, "|")

	if len(strs) != 2 {
		return 0, 0, 0, fmt.Errorf("Payout format error")
	}

	if val, err := strconv.ParseFloat(strs[0], 32); err == nil {
		jackpotID = int32(val)
	} else {
		return 0, 0, 0, fmt.Errorf("Payout format error. JackpoID is invalid. %v", err)
	}

	type Group struct {
		BetRate    int64
		prop       int64
		Denminator int64
	}

	rateMap := make(map[int64]*Group)
	groups := strings.Split(strs[1], ",")

	for _, group := range groups {
		data := strings.Split(group, ":")

		g := new(Group)

		if val, err := strconv.ParseFloat(data[0], 64); err == nil {
			g.BetRate = int64(val)
		} else {
			continue
		}

		if val, err := strconv.ParseFloat(data[1], 64); err == nil {
			g.prop = int64(val)
		} else {
			continue
		}

		if val, err := strconv.ParseFloat(data[2], 64); err == nil {
			g.Denminator = int64(val)
		} else {
			continue
		}

		rateMap[g.BetRate] = g
	}

	if g, ok := rateMap[betRate]; ok {
		prop = g.prop
		denominator = g.Denminator
	} else {
		return 0, 0, 0, fmt.Errorf("Payout format error")
	}

	return
}

// ParsePayoutMysteryJPG 解析MysteryJPG(機率中JackpotGame)的Payout內容
// @Argument
//   * pattern 		Payout文字內容  {連線機率}:{非連線機率}:{分母}	ex: 7:100:10000
// @Return
//   * linkprop 	連線JP機率 (未來要乘 TotalBet)
//   * nolinkprop	非連線JP機率
//   * denominator	機率分母
func ParsePayoutMysteryJPG(pattern string) (linkprop, nolinkprop, denominator int64, err error) {
	pattern = strings.TrimSpace(pattern)

	strs := strings.Split(pattern, ":")

	if len(strs) != 3 {
		return 0, 0, 0, fmt.Errorf("Payout format error")
	}

	if val, err := strconv.ParseFloat(strs[0], 64); err == nil {
		linkprop = int64(val)
	} else {
		return 0, 0, 0, fmt.Errorf("Payout format error. linkprop is invalid. %v", err)
	}

	if val, err := strconv.ParseFloat(strs[1], 64); err == nil {
		nolinkprop = int64(val)
	} else {
		return 0, 0, 0, fmt.Errorf("Payout format error. nolinkprop is invalid. %v", err)
	}

	if val, err := strconv.ParseFloat(strs[2], 64); err == nil {
		denominator = int64(val)
	} else {
		return 0, 0, 0, fmt.Errorf("Payout format error. denominator is invalid. %v", err)
	}

	return
}

// ParsePayoutMysteryJPGR 解析MysteryJPGR(依照不同倍率機率中JackpotGame)的Payout內容
// @Argument
//   * pattern 		Payout文字內容  {下注倍率}:{連線機率}:{非連線機率}:{分母},...	ex: 8:100:1e5:1e10,18:200:1e5:1e10,38:300:1e5:1e10,88:400:1e5:1e10,88:500:1e5:1e10
//   * betRate 		當前下注倍率
// @Return
//   * linkprop 	連線JP機率 (未來要乘 TotalBet)
//   * nolinkprop	非連線JP機率
//   * denominator	機率分母
func ParsePayoutMysteryJPGR(pattern string, betRate int64) (linkprop, nolinkprop, denominator int64, err error) {
	pattern = strings.TrimSpace(pattern)

	strs := strings.Split(pattern, ",")

	data := make(map[int64]string)

	for _, str := range strs {
		_strs := strings.Split(str, ":")

		if len(_strs) != 4 {
			return 0, 0, 0, fmt.Errorf("Payout format error")
		}

		var rate int64
		var p string

		if val, err := strconv.ParseFloat(_strs[0], 64); err == nil {
			rate = int64(val)
		} else {
			return 0, 0, 0, fmt.Errorf("Payout format error. BetRate is invalid. %v", err)
		}

		p = strings.Join(_strs[1:], ":")

		data[rate] = p
	}

	str, ok := data[betRate]
	if !ok {
		return 0, 0, 0, fmt.Errorf("BetRate not exist")
	}

	return ParsePayoutMysteryJPG(str)
}

// ParseExtraRuleConditionAfterSpin 解析AfterSpin內容
// @Argument
//   * value	ConditionValue文字內容	{連線類型}:{符號名稱}:{連線個數}	ex: any:g1:15
// @Return
//   * linkType	連線類型
//   * symbol	符號名稱
//   * count	連線個數
func ParseExtraRuleConditionAfterSpin(pattern string) (linkType, symbol string, count int, err error) {
	pattern = strings.TrimSpace(pattern)

	strs := strings.Split(pattern, ":")

	if len(strs) != 3 {
		return "", "", 0, fmt.Errorf("Payout format error")
	}

	linkType = strs[0]
	symbol = strs[1]

	if val, err := strconv.ParseFloat(strs[2], 8); err == nil {
		count = int(val)
	} else {
		return "", "", 0, fmt.Errorf("Payout format error. count is invalid. %v", err)
	}

	return
}

// ParseExtraRuleConditionAppearSymbol 解析AppearSymbol內容
// @Argument
//   * pattern	ConditionValue文字內容	{Symbol}|{Reel/Row}|{ReelNum/RowNum}...	ex: w1|Reel|2,3,4
// @Return
//   * symbol	符號名稱
//   * ror		Reel or Row or Any
//   * vals		對應的 Reel or Row
func ParseExtraRuleConditionAppearSymbol(pattern string) (symbol, ror string, vals []int, err error) {
	pattern = strings.TrimSpace(pattern)

	strs := strings.Split(pattern, "|")

	if len(strs) != 2 && len(strs) != 3 {
		return "", "", []int(nil), fmt.Errorf("Payout format error")
	}

	symbol = strs[0]
	ror = strs[1]

	if ror == "Any" {
		return
	}

	if len(strs) != 3 {
		return "", "", []int(nil), fmt.Errorf("Payout format error")
	}

	valstrs := strings.Split(strs[2], ",")

	for _, s := range valstrs {
		if len(s) == 0 {
			continue
		}
		if val, err := strconv.ParseFloat(s, 8); err == nil {
			vals = append(vals, int(val))
		} else {
			return "", "", []int(nil), fmt.Errorf("Payout format error. vals is invalid. %v", err)
		}
	}

	return
}

// ParseExtraRuleActionRespin 解析Respin內容
// @Argument
//   * pattern	ActionValue文字內容	{機率}:{分母}	ex: 100:100000
// @Return
//   * prop			中獎機率
//   * denominator	機率分母
func ParseExtraRuleActionRespin(pattern string) (prop, denominator int64, err error) {
	pattern = strings.TrimSpace(pattern)

	strs := strings.Split(pattern, ":")

	if len(strs) != 2 {
		return 0, 0, fmt.Errorf("Payout format error")
	}

	if val, err := strconv.ParseFloat(strs[0], 64); err == nil {
		prop = int64(val)
	} else {
		return 0, 0, fmt.Errorf("Payout format error. prop is invalid. %v", err)
	}

	if val, err := strconv.ParseFloat(strs[1], 64); err == nil {
		denominator = int64(val)
	} else {
		return 0, 0, fmt.Errorf("Payout format error. denominator is invalid. %v", err)
	}

	return
}

// ParseExtraRuleActionRespinR 解析Respin內容
// @Argument
//   * pattern	ActionValue文字內容	{下注倍率}:{機率}:{分母},...	ex: 8:100:1e5,18:200:1e5,...
// @Return
//   * prop			中獎機率
//   * denominator	機率分母
func ParseExtraRuleActionRespinR(pattern string, rate int64) (prop, denominator int64, err error) {
	pattern = strings.TrimSpace(pattern)

	patterns := strings.Split(pattern, ",")

	patternMap := make(map[int64]string)

	for _, _pattern := range patterns {

		strs := strings.Split(_pattern, ":")

		if val, err := strconv.ParseFloat(strs[0], 64); err == nil {
			patternMap[int64(val)] = strings.Join(strs[1:], ":")
		} else {
			return 0, 0, fmt.Errorf("Payout format error. BetRate is invalid. %v", err)
		}
	}

	if str, ok := patternMap[rate]; ok {
		return ParseExtraRuleActionRespin(str)
	}

	return 0, 0, fmt.Errorf("Payout format error. Could not find Rate %d", rate)
}

// ParseExtraRuleActionUniGoldBall 解析Respin內容
// @Argument
//   * pattern	ActionValue文字內容	{機率}:{分母}	ex: 100:100000
// @Return
//   * symbol		要增加的符號
//   * prop			中獎機率
//   * denominator	機率分母
func ParseExtraRuleActionUniGoldBall(pattern string) (symbol string, prop, denominator int64, err error) {
	pattern = strings.TrimSpace(pattern)

	strs := strings.Split(pattern, ":")

	if len(strs) != 3 {
		return "", 0, 0, fmt.Errorf("Payout format error")
	}

	symbol = strs[0]

	if val, err := strconv.ParseFloat(strs[1], 64); err == nil {
		prop = int64(val)
	} else {
		return "", 0, 0, fmt.Errorf("Payout format error. prop is invalid. %v", err)
	}

	if val, err := strconv.ParseFloat(strs[2], 64); err == nil {
		denominator = int64(val)
	} else {
		return "", 0, 0, fmt.Errorf("Payout format error. denominator is invalid. %v", err)
	}

	return
}

// ParseExtraRuleActionSwitchReelRandom 解析SwitchReelRandom內容
// @Argument
//   * pattern	ActionValue文字內容	{機率}:{分母}	ex: 1:4
// @Return
//   * prop			中獎機率
//   * denominator	機率分母
func ParseExtraRuleActionSwitchReelRandom(pattern string) (prop, denominator int64, err error) {
	pattern = strings.TrimSpace(pattern)

	strs := strings.Split(pattern, ":")

	if len(strs) != 2 {
		return 0, 0, fmt.Errorf("Payout format error")
	}

	if val, err := strconv.ParseFloat(strs[0], 64); err == nil {
		prop = int64(val)
	} else {
		return 0, 0, fmt.Errorf("Payout format error. prop is invalid. %v", err)
	}

	if val, err := strconv.ParseFloat(strs[1], 64); err == nil {
		denominator = int64(val)
	} else {
		return 0, 0, fmt.Errorf("Payout format error. denominator is invalid. %v", err)
	}

	return
}

// ParseExtraRuleActionReplaceBG 解析ReplaceBG內容
// @Argument
//   * pattern	ActionValue文字內容	{Reel/Row}:{val}:{BGName}:{機率}:{分母}	ex: Reel:5:bg4:100:1000000
// @Return
//   * ror			Reel or Row
//   * target		第幾個 Reel or Row
//   * symbol		替換符號
//   * prop			中獎機率
//   * denominator	機率分母
func ParseExtraRuleActionReplaceBG(pattern string) (ror string, target int, symbol string, prop, denominator int64, err error) {
	pattern = strings.TrimSpace(pattern)

	strs := strings.Split(pattern, ":")

	if len(strs) != 5 {
		return "", 0, "", 0, 0, fmt.Errorf("Payout format error")
	}

	ror = strs[0]

	if val, err := strconv.ParseFloat(strs[1], 8); err == nil {
		target = int(val)
	} else {
		return "", 0, "", 0, 0, fmt.Errorf("Payout format error. taget is invalid. %v", err)
	}

	symbol = strs[2]

	if val, err := strconv.ParseFloat(strs[3], 64); err == nil {
		prop = int64(val)
	} else {
		return "", 0, "", 0, 0, fmt.Errorf("Payout format error. prop is invalid. %v", err)
	}

	if val, err := strconv.ParseFloat(strs[4], 64); err == nil {
		denominator = int64(val)
	} else {
		return "", 0, "", 0, 0, fmt.Errorf("Payout format error. denominator is invalid. %v", err)
	}

	return
}

// ParseExtraRuleActionReplaceBGR 解析ReplaceBGR內容
// @Argument
//   * pattern	ActionValue文字內容	{下注倍率}:{Reel/Row}:{val}:{BGName}:{機率}:{分母},...	ex: 28:Reel:5:bg4:100:1000000,...
//   * betRate 		當前下注倍率
// @Return
//   * ror			Reel or Row
//   * target		第幾個 Reel or Row
//   * symbol		替換符號
//   * prop			中獎機率
//   * denominator	機率分母
func ParseExtraRuleActionReplaceBGR(pattern string, betRate int64) (ror string, target int, symbol string, prop, denominator int64, err error) {
	pattern = strings.TrimSpace(pattern)

	strs := strings.Split(pattern, ",")

	data := make(map[int64]string)

	for _, str := range strs {
		_strs := strings.Split(str, ":")

		if len(_strs) != 6 {
			return "", 0, "", 0, 0, fmt.Errorf("Payout format error")
		}

		var rate int64
		var p string

		if val, err := strconv.ParseFloat(_strs[0], 64); err == nil {
			rate = int64(val)
		} else {
			return "", 0, "", 0, 0, fmt.Errorf("Payout format error. BetRate is invalid. %v", err)
		}

		p = strings.Join(_strs[1:], ":")

		data[rate] = p
	}

	str, ok := data[betRate]
	if !ok {
		return "", 0, "", 0, 0, fmt.Errorf("BetRate not exist")
	}

	return ParseExtraRuleActionReplaceBG(str)
}

// ParseExtraRuleActionConditionalAction 解析ConditionalAction內容
// @Argument
//   * pattern	ActionValue文字內容	{Condition}:{Action}	ex: JPDefense:Respin
// @Return
//   * condition	條件
//   * action		動作
func ParseExtraRuleActionConditionalAction(pattern string) (condition ExtraRuleCondition, action ExtraRuleAction, err error) {
	pattern = strings.TrimSpace(pattern)

	strs := strings.Split(pattern, ":")

	if len(strs) != 2 {
		return ExtraRuleCondition(""), ExtraRuleAction(""), fmt.Errorf("Payout format error")
	}

	condition = ExtraRuleCondition(strs[0])
	action = ExtraRuleAction(strs[1])

	return
}
