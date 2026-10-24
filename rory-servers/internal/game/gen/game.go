package gen

import (
	"fmt"
	"math/rand"
	"runtime/debug"
	"strconv"
	"strings"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
	"example.invalid/ServerTeam/Servers/pkg/prob"
)

// Game 將軍
type Game struct {
	*gamelogic.LogicBase

	Bet               int64
	ChanceRangeOptMap map[pb.ModeType][][2]int // 機會模式區間設定 n組 [min, max]
	ChanceRate        int                      // 機會模式觸發機率
	ChanceCntOpt      [][2]int                 // 機會模式次數設定 n組 [count, weight]
	ChanceAlways      bool                     // 機會模式只要條件符合就進入，不骰1/256機率

	ReelStrips    pachislot.ReelStripGroup      // 輪軸帶資料
	Validator     *pachislot.ResultSetValidator // 牌型結果驗證器
	PachiSlotGame *pachislot.Game               // PachiSlot遊戲主體

	BellLinkInJAC  [3]int // 在JAC中獲得BellLink的次數
	JACInBB        int    // 在BB中獲得JAC的次數
	ChanceBuff     bool   // 機會模式
	ChanceDeadLine int    // 機會模式結束回合，RandMap.SpinCount相比

	isgaming bool
}

// NewGame 建立將軍遊戲
func NewGame(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group, protocfg *pachislot.Proto, pid int32, bet int64) *Game {
	game := new(Game)
	game.LogicBase = gamelogic.NewGameLogicBase(t, tfManager, cfgs)
	game.ChanceRangeOptMap = make(map[pb.ModeType][][2]int)

	// Load Config
	{ // Chance Range Options
		str := cfgs.GetConfigValue(gamelogic.ConfigFileNameGameServer, gamelogic.ConfigSectionNameGameServer, "chancerangeopt")
		modecfgs := strings.Split(str, ";")
		for _, modecfg := range modecfgs {
			modecfg = strings.TrimSpace(modecfg)
			strs := strings.Split(modecfg, ":")
			if len(strs) != 2 {
				mlog.Warn("gameserver.dat [gameserver] chancerangeopt get error group %s, want \"{mode}:{min}-{max},{min}-{max},...\"", modecfg)
				continue
			}
			strs[0] = strings.TrimSpace(strs[0])
			strs[1] = strings.TrimSpace(strs[1])
			mode := pb.ModeType(pb.ModeType_value[strs[0]])
			opts := strings.Split(strs[1], ",")
			for _, opt := range opts {
				opt = strings.TrimSpace(opt)
				_str := strings.Split(opt, "-")
				if len(_str) != 2 {
					mlog.Warn("gameserver.dat [gameserver] chancerangeopt get error option %s, want \"min:max\"", opt)
					continue
				}
				var min, max int
				if val, err := strconv.ParseInt(_str[0], 0, 32); err == nil {
					min = int(val)
				} else {
					mlog.Warn("gameserver.dat [gameserver] chancerangeopt get error option %s, %v", opt, err)
					continue
				}
				if val, err := strconv.ParseInt(_str[1], 0, 32); err == nil {
					max = int(val)
				} else {
					mlog.Warn("gameserver.dat [gameserver] chancerangeopt get error option %s, %v", opt, err)
					continue
				}
				game.ChanceRangeOptMap[mode] = append(game.ChanceRangeOptMap[mode], [2]int{min, max})
			}
		}
	}
	{ // Chance Rate
		str := cfgs.GetConfigValue(gamelogic.ConfigFileNameGameServer, gamelogic.ConfigSectionNameGameServer, "chancerate")
		if val, err := strconv.ParseInt(str, 0, 32); err == nil {
			game.ChanceRate = int(val)
		} else {
			mlog.Warn("gameserver.dat [gameserver] chancerate %s parse error, %v", str, err)
			game.ChanceRate = 256
		}
	}
	{ // Chance Count Options
		str := cfgs.GetConfigValue(gamelogic.ConfigFileNameGameServer, gamelogic.ConfigSectionNameGameServer, "chancecntopt")
		opts := strings.Split(str, ",")
		for _, opt := range opts {
			opt = strings.TrimSpace(opt)
			_str := strings.Split(opt, ":")
			if len(_str) != 2 {
				mlog.Warn("gameserver.dat [gameserver] chancecntopt get error option %s, want \"count:weight\"", opt)
				continue
			}
			var cnt, w int
			if val, err := strconv.ParseInt(_str[0], 0, 32); err == nil {
				cnt = int(val)
			} else {
				mlog.Warn("gameserver.dat [gameserver] chancecntopt get error option %s, %v", opt, err)
				continue
			}
			if val, err := strconv.ParseInt(_str[1], 0, 32); err == nil {
				w = int(val)
			} else {
				mlog.Warn("gameserver.dat [gameserver] chancecntopt get error option %s, %v", opt, err)
				continue
			}
			game.ChanceCntOpt = append(game.ChanceCntOpt, [2]int{cnt, w})
		}
	}

	game.Bet = bet

	pachigame, err := pachislot.NewGame(protocfg, pid)
	if err != nil {
		mlog.Fatal("Create PachiSlotGame(%d) Error, %v", pid, err)
	}
	pachigame.RoadMap.FeaturesCnt[pb.Scenario_BB] = 0
	pachigame.RoadMap.FeaturesCnt[pb.Scenario_RB] = 0
	pachigame.CostPerCoin = bet

	game.PachiSlotGame = pachigame

	game.ReelStrips = pachislot.NewReelStripGroup(pachigame.Proto.ReelStrips)
	game.Validator = pachislot.NewResultSetValidator(pachislot.UniqResultTypes(pachigame.Proto.PayTable))

	// 檢查是否達成在機會模式下
	game.PachiSlotGame.AddBeforeInitFunc(game.chkAssignableHeaven)

	// 檢查是否觸發機會模式
	game.PachiSlotGame.AddAfterSpinFunc(game.chkChanceBuff)

	// 檢查是否達成1G連
	game.PachiSlotGame.AddBeforeInitFunc(game.chkAvariable1G)

	// 初始時清除資料
	game.PachiSlotGame.AddAfterInitFunc(game.clean)

	// 檢查是否在JAC骰到Bell
	game.PachiSlotGame.AddAfterSpinFunc(game.chkBellInJAC)

	// 檢查在BB骰到JAC
	game.PachiSlotGame.AddAfterRewardFunc(game.chkJACInBB)

	// 遊戲開始後發送當前轉數表
	game.RegisterAfterGameStartFunc(func() {
		m := game.mRoadMap()
		game.BoardcastMessage(m)
	})

	// 遊戲開始後發送當前遊戲狀態
	game.RegisterAfterGameStartFunc(func() {
		m := game.mGameStatus()
		game.BoardcastMessage(m)
	})

	return game
}

// PiSaiMessage 遊戲封包
func (game *Game) PiSaiMessage(user *gamelogic.GameUser, m *msg.Message) {}

// GameMessage 遊戲封包
func (game *Game) GameMessage(user *gamelogic.GameUser, m *msg.Message) {
	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("GameMessage 發生錯誤, %v\n%s", rec, string(debug.Stack()))
		}
	}()

	if !game.CheckGameUserMessage(user) {
		game.KickUserFunc(user.UserID, "非法操作")
		return
	}

	// 空訊息，不處理
	if m == nil {
		return
	}

	// 非遊戲訊息，不處理
	if m.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	switch netproto.GEN_GameMessageClassID(m.SClassID) {
	case netproto.GEN_GameMessageClassID_GEN_Read_GameSettingID:
		m := game.mGameSetting()
		game.SendMessgeFunc(user.UserID, m)
	case netproto.GEN_GameMessageClassID_GEN_Read_GameStatusID:
		m := game.mGameStatus()
		game.SendMessgeFunc(user.UserID, m)
	case netproto.GEN_GameMessageClassID_GEN_Read_RoadMapID:
		m := game.mRoadMap()
		game.SendMessgeFunc(user.UserID, m)
	case netproto.GEN_GameMessageClassID_GEN_BetID:
		data := new(netproto.GEN_Bet)
		err := proto.Unmarshal(m.MsgData.([]byte), data)
		if err != nil {
			return
		}

		{ // 檢查是否倍率合法
			if data.Odds != game.Bet {
				err := fmt.Errorf("非法下注倍率 %d, 本機台倍率 %d", data.Odds, game.Bet)
				example.invalid("玩家[%d] 下注失敗, %v", user.UserID, err)
				m := game.mError(netproto.GEN_ErrorType_GEN_ErrBet_IllegalOdds, err, netproto.GEN_GameMessageClassID_GEN_BetID)
				game.SendMessgeFunc(user.UserID, m)
				return
			}
		}

		{ // 檢查所需代幣數是否足夠
			coin := game.PachiSlotGame.BetCoin()
			bet := coin * game.Bet
			if user.Score < bet {
				err := fmt.Errorf("下注 %d 倍，當前所需 %d 枚代幣，共需 %.2f，玩家金額 %.2f 不足所需金額", game.Bet, coin, float64(bet)/100, float64(user.Score)/100)
				example.invalid("玩家[%d] 下注失敗, %v", user.UserID, err)
				m := game.mError(netproto.GEN_ErrorType_GEN_ErrBet_MoneyNotEnough, err, netproto.GEN_GameMessageClassID_GEN_BetID)
				game.SendMessgeFunc(user.UserID, m)
				return
			}
		}

		game.Play(game.Bet, user)
	}
}

// Play 轉
func (game *Game) Play(odds int64, user *gamelogic.GameUser) {

	scenario := game.PachiSlotGame.Scenario
	coin := game.PachiSlotGame.BetCoin()
	bet := coin * odds

	game.GameLog.AddGameStep("遊戲: %v", game.PachiSlotGame)
	game.GameLog.AddGameStep("%v", game.PachiSlotGame.RoadMap)
	game.GameLog.AddGameStep("機台倍率: %d", game.Bet)

	_, result, releaseRT := game.PachiSlotGame.Spin()

	game.GameLog.AddGameStep("牌型結果[%v]", result)

	if releaseRT {
		game.GameLog.AddGameStep("玩家[%d] 牌型結果[%v] 解除RT", user.UserID, result)
		// 只要是透過 3松連線 這個 Result 解除的，下次Mode一定要是天國
		if result == pb.ResultType_WaterMelonLink {
			game.GameLog.AddGameStep("玩家[%d] 解除RT，牌型結果為3松連線，指定下次模式為天國", user.UserID)
			game.PachiSlotGame.NextMode = pb.ModeType_Heaven
		}
	}

	var winCoin int64
	var win int64

	payouts := game.PachiSlotGame.Reward(result)
	for _, p := range payouts {
		switch p.Type() {
		default:
			game.GameLog.AddGameStep("玩家[%d] 獲得獎勵[%v]", user.UserID, p.Type())
		case pb.PayType_PTChip:
			game.GameLog.AddGameStep("玩家[%d] 獲得獎勵[%v] %d 倍", user.UserID, p.Type(), p.Odds())
		case pb.PayType_PTBonus:
			game.GameLog.AddGameStep("玩家[%d] 獲得獎勵[%v] %v", user.UserID, p.Type(), game.PachiSlotGame.Scenario)
		case pb.PayType_PT1G:
			game.GameLog.AddGameStep("玩家[%d] 獲得獎勵[%v] %d 次，已累積 %d 次", user.UserID, p.Type(), p.(*pachislot.Payout1G).Times, game.PachiSlotGame.BonusStore)
		}
		winCoin += p.Odds()
	}
	win = winCoin * odds

	// 更新金幣
	user.BetAmount += bet
	user.FlowAmount += bet
	user.PayAmount += win
	game.UpdateUserScoreFunc(user.UserID, win-bet, 0)

	if game.GameLog != nil && game.GameLog.Context() != nil {
		game.GameLog.Context().UserScore(user.UserID, user.Score)
		game.GameLog.Context().UserUpdateScore(user.UserID, win-bet)
		game.GameLog.Context().UserBetScore(user.UserID, bet)
		game.GameLog.Context().UserPayoutScore(user.UserID, win)
		game.GameLog.Context().UserValidScore(user.UserID, bet)
	}
	game.GameLog.AddGameStep("玩家[%d] 當前(%v)消耗 %d 枚代幣，獲得 %d 枚代幣，更新 %.2f，剩餘 %.2f ", user.UserID, scenario, coin, winCoin, float64(win-bet)/100, float64(user.Score)/100)
	game.GameLog.GameEnd()

	// 產生盤面
	set := pachislot.GenereteResultSet(result)
	blinks, m, err := set.Generete(game.PachiSlotGame.Proto.PayLine, game.ReelStrips)
	if err != nil {

	}
	sg := m.SymbolGroup(game.ReelStrips)

	var cnt int
	for game.Validator.MatchOthers(set.Type(), sg, game.PachiSlotGame.Proto.PayLine) {
		if cnt > 1000 {
			mlog.Error("重新產生盤面超過 %d 次", 1000)
			break
		}
		if set.Type() == pb.ResultType_DoubleSevenLink {
			mlog.Debug("雙7連線，必定會與藍7連線重複，直接跳過")
			break
		}
		blinks, m, _ = set.Generete(game.PachiSlotGame.Proto.PayLine, game.ReelStrips)
		sg = m.SymbolGroup(game.ReelStrips)
		cnt++
	}

	{ // Send Result Message
		m := game.mGameResult(scenario, bet, win, result, m, sg, payouts, releaseRT, blinks)
		game.SendMessgeFunc(user.UserID, m)
	}

	game.GameEndFunc()
}

func (game *Game) clean(g *pachislot.Game, releaseBonus bool) {
	game.BellLinkInJAC = [3]int{}
	game.JACInBB = 0
	if releaseBonus {
		game.GameLog.AddGameStep("使用 1G連 庫存，剩餘 %d 次", game.PachiSlotGame.BonusStore)
	}
}

// 機會模式開啟的條件
func (game *Game) chkChanceBuff(g *pachislot.Game, result pb.ResultType, releaseRT bool) {

	if game.ChanceBuff {
		// 機會模式中，檢查是否關閉
		if game.PachiSlotGame.RoadMap.SpinCnt > game.ChanceDeadLine {
			game.ChanceBuff = false
			game.ChanceDeadLine = 0
			game.GameLog.AddGameStep("機會模式次數用盡，解除機會模式")
		}
		return
	}

	// 要在NG下
	if game.PachiSlotGame.Scenario != pb.Scenario_NG {
		return
	}

	// 要是Replay連線
	if result != pb.ResultType_ReplayLink {
		return
	}

	// 模式有設定
	rangeOpts, ok := game.ChanceRangeOptMap[game.PachiSlotGame.Mode]
	if !ok {
		return
	}

	// 要在區間中
	var inRange bool
	for _, r := range rangeOpts {
		if r[0] > game.PachiSlotGame.RoadMap.SpinCnt {
			continue
		}
		if r[1] < game.PachiSlotGame.RoadMap.SpinCnt {
			continue
		}
		inRange = true
		break
	}
	if !inRange {
		return
	}

	// 1/256 機率觸發
	n := rand.Intn(game.ChanceRate)
	if n != 0 && !game.ChanceAlways {
		return
	}

	// 骰出次數
	r := new(prob.Random)
	for _, opt := range game.ChanceCntOpt {
		r.Add(opt[0], opt[1])
	}
	cnt := r.Choose().(int)

	game.ChanceBuff = true
	game.ChanceDeadLine = game.PachiSlotGame.RoadMap.SpinCnt + cnt

	game.GameLog.AddGameStep("成功觸發機會模式，共 %d 次", cnt)
}

// 機會模式下，重新決定模式時，固定為天國
func (game *Game) chkAssignableHeaven(g *pachislot.Game) {

	if !game.ChanceBuff {
		return
	}

	g.NextMode = pb.ModeType_Heaven

	game.ChanceBuff = false
	game.ChanceDeadLine = 0
	game.GameLog.AddGameStep("機會模式下，指定下次模式為天國，解除機會模式")
}

// 在Bonus中，達成下述條件，即可在Bonus Game結束後，第一轉開出Bonus確定
//  * JAC中，連續8次銅鐘連線，可連莊1次
//  * BIG BONUS中 觸發JAC IN 0次，可連莊3次
//    BIG BONUS中 觸發JAC IN 1次，可連莊2次
//    BIG BONUS中 觸發JAC IN 2次，可連莊1次
func (game *Game) chkAvariable1G(g *pachislot.Game) {

	if g.Scenario != pb.Scenario_BB {
		return
	}

	switch game.JACInBB {
	case 0: // BIG BONUS中 觸發JAC IN 0次，可連莊3次
		g.BonusStore += 3
		m := game.mAppend1GBB0()
		game.BoardcastMessage(m)
		game.GameLog.AddGameStep("在 %v 中累積 [%v] 達 0 次，獲得 1G連 3 次，已累積 %d 次", pb.Scenario_BB, pb.Scenario_JAC, game.PachiSlotGame.BonusStore)
	case 1: // BIG BONUS中 觸發JAC IN 1次，可連莊2次
		g.BonusStore += 2
		m := game.mAppend1GBB1()
		game.BoardcastMessage(m)
		game.GameLog.AddGameStep("在 %v 中累積 [%v] 達 1 次，獲得 1G連 2 次，已累積 %d 次", pb.Scenario_BB, pb.Scenario_JAC, game.PachiSlotGame.BonusStore)
	case 2: // BIG BONUS中 觸發JAC IN 2次，可連莊1次
		g.BonusStore++
		m := game.mAppend1GBB2()
		game.BoardcastMessage(m)
		game.GameLog.AddGameStep("在 %v 中累積 [%v] 達 2 次，獲得 1G連 1 次，已累積 %d 次", pb.Scenario_BB, pb.Scenario_JAC, game.PachiSlotGame.BonusStore)
	}
}

// 在JAC下累積BellLink次數
func (game *Game) chkBellInJAC(g *pachislot.Game, result pb.ResultType, rt bool) {
	if g.Scenario == pb.Scenario_JAC && result == pb.ResultType_BellLink {
		game.BellLinkInJAC[g.JACTimes-1]++

		if game.BellLinkInJAC[g.JACTimes-1] == 8 {
			game.PachiSlotGame.BonusStore++
			m := game.mAppend1GJAC8Bell()
			game.BoardcastMessage(m)
			game.GameLog.AddGameStep("在 %v 中累積 [%v] 達 8 次，獲得 1G連 1 次，已累積 %d 次", pb.Scenario_JAC, pb.ResultType_BellLink, game.PachiSlotGame.BonusStore)
		}
	}
}

// 在BB下累積JAC的次數
func (game *Game) chkJACInBB(g *pachislot.Game, payouts []pachislot.Payout) {
	if g.Scenario == pb.Scenario_BB {
		return
	}
	for _, p := range payouts {
		if p.Type() != pb.PayType_PTBonus {
			continue
		}
		if p.(*pachislot.PayoutBonus).Module == pb.Scenario_JAC {
			game.JACInBB++
		}
	}
}
