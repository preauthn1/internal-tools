package inrm

import (
	"bytes"
	"fmt"
	"math"
	"runtime/debug"
	"strconv"
	"time"

	"example.invalid/dustin/go-humanize"
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
)

// Timers 每桌要區隔 Timer Key 不然會跟其他桌打架
const (
	TimerWaitPlayer = "WaitPlayer"
	TimerStartCD    = "StartCD"
	TimerWaitAction = "WaitAction"
	TimerActionCD   = "ActionCD"
	TimerWaitCalc   = "WaitCalc"
	TimerCalcCD     = "CalcCD"
)

// TableState 桌子狀態
type TableState uint

// TableStates
const (
	// 等玩家進入
	WaitPlayer TableState = iota
	// 倒數開始
	CountDown
	// 遊戲中
	Gaming
)

// GameServer 印度拉米遊戲(Server用)
type GameServer struct {
	*gamelogic.LogicBase
	isgaming bool
	State    TableState

	Game        *Game
	GameOptions []GameOpt
	LastRound   int

	ActionTimeoutAt time.Time // 玩家超時時間
	StartAt         time.Time // 遊戲實際開始時間
	FinalAt         time.Time // 結算時間

	StartWaitTime  time.Duration // 等待開始時間
	ActionWaitTime time.Duration // 每次動作等待時間
	ScorePerPoint  int64         // 每一點使用的分數
}

// NewGameServer 印度拉米遊戲(Server用)
func NewGameServer(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *GameServer {
	gs := new(GameServer)
	gs.LogicBase = gamelogic.NewGameLogicBase(t, tfManager, cfgs)

	gs.initGameOptions(cfgs) // 初始化設定

	gs.RegisterAfterGameStartFunc(gs.Start)
	gs.RegisterAfterGameEndFunc(gs.Clean)

	return gs
}

// Summary 此遊戲邏輯目前的狀態描述
func (gs *GameServer) Summary() string {
	switch gs.State {
	default:
		return ""
	case WaitPlayer:
		return "等待玩家中"
	case CountDown:
		return "遊戲開始 倒數中"
	case Gaming:
		buf := bytes.NewBuffer(nil)
		fmt.Fprintf(buf, "遊戲中 ...\n")
		fmt.Fprintf(buf, "%v\n", gs.Game)
		for _, p := range gs.Game.Players {
			fmt.Fprintf(buf, " - %v\n", p)
		}
		return buf.String()
	}
}

func (gs *GameServer) initGameOptions(cfgs *config.Group) {
	{ // startwaittime=10s		// 倒數等待開始時間 Default: 10s
		s := cfgs.GetConfigValue(gamelogic.ConfigFileNameGameServer, gamelogic.ConfigSectionNameGameServer, "startwaittime")
		if d, err := time.ParseDuration(s); err != nil {
			gs.StartWaitTime = 10 * time.Second
		} else {
			gs.StartWaitTime = d
		}
	}
	{ // actionwaittime=15s		// 每次動作等待時間 Default: 15s
		s := cfgs.GetConfigValue(gamelogic.ConfigFileNameGameServer, gamelogic.ConfigSectionNameGameServer, "actionwaittime")
		if d, err := time.ParseDuration(s); err != nil {
			gs.ActionWaitTime = 15 * time.Second
		} else {
			gs.ActionWaitTime = d
		}
	}
	{ // scoerperpoint=50		// 每一點使用的分數 Default: 50
		s := cfgs.GetConfigValue(gamelogic.ConfigFileNameGameServer, gamelogic.ConfigSectionNameGameServer, "scoerperpoint")
		if val, err := strconv.ParseInt(s, 0, 64); err != nil {
			mlog.Warn("config scoerperpoint parse error, %v", err)
			gs.ScorePerPoint = 50
		} else {
			gs.ScorePerPoint = val
		}
	}
	{ // handcardlehgth=13		// 手牌張數 Default: 13
		s := cfgs.GetConfigValue(gamelogic.ConfigFileNameGameServer, gamelogic.ConfigSectionNameGameServer, "handcardlehgth")
		if val, err := strconv.ParseInt(s, 0, 8); err != nil {
			mlog.Warn("config handcardlehgth parse error, %v", err)
		} else {
			gs.GameOptions = append(gs.GameOptions, HandCardLength(int(val)))
		}
	}
	{ // maxpoint=80			// 最大失敗點數 Default: 80
		s := cfgs.GetConfigValue(gamelogic.ConfigFileNameGameServer, gamelogic.ConfigSectionNameGameServer, "maxpoint")
		if val, err := strconv.ParseInt(s, 0, 8); err != nil {
			mlog.Warn("config maxpoint parse error, %v", err)
		} else {
			gs.GameOptions = append(gs.GameOptions, MaxPoint(int(val)))
		}
	}
	{ // giveuppoint=40			// 棄牌點數 Default: 40
		s := cfgs.GetConfigValue(gamelogic.ConfigFileNameGameServer, gamelogic.ConfigSectionNameGameServer, "giveuppoint")
		if val, err := strconv.ParseInt(s, 0, 8); err != nil {
			mlog.Warn("config giveuppoint parse error, %v", err)
		} else {
			gs.GameOptions = append(gs.GameOptions, GiveupPoint(int(val)))
		}
	}
	{ // misslimit=3			// 失誤次數 Default: 3
		s := cfgs.GetConfigValue(gamelogic.ConfigFileNameGameServer, gamelogic.ConfigSectionNameGameServer, "misslimit")
		if val, err := strconv.ParseInt(s, 0, 8); err != nil {
			mlog.Warn("config misslimit parse error, %v", err)
		} else {
			gs.GameOptions = append(gs.GameOptions, MissLimit(int(val)))
		}
	}
	{ // shuffletimes=3			// 每局洗牌次數 Default: 3
		s := cfgs.GetConfigValue(gamelogic.ConfigFileNameGameServer, gamelogic.ConfigSectionNameGameServer, "shuffletimes")
		if val, err := strconv.ParseInt(s, 0, 8); err != nil {
			mlog.Warn("config shuffletimes parse error, %v", err)
		} else {
			gs.GameOptions = append(gs.GameOptions, ShuffleTimes(int(val)))
		}
	}
	{ // usepokernum=2			// 使用幾副牌 Default: 2
		s := cfgs.GetConfigValue(gamelogic.ConfigFileNameGameServer, gamelogic.ConfigSectionNameGameServer, "usepokernum")
		if val, err := strconv.ParseInt(s, 0, 8); err != nil {
			mlog.Warn("config usepokernum parse error, %v", err)
		} else {
			gs.GameOptions = append(gs.GameOptions, UsePokerNum(int(val)))
		}
	}
}

// PiSaiMessage 遊戲封包
func (gs *GameServer) PiSaiMessage(user *gamelogic.GameUser, m *msg.Message) {}

// GameMessage 遊戲封包
func (gs *GameServer) GameMessage(user *gamelogic.GameUser, m *msg.Message) {
	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("GameMessage 發生錯誤, %v\n%s", rec, string(debug.Stack()))
		}
	}()

	if !gs.CheckGameUserMessage(user) {
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

	switch netproto.INRM_GameMessageClassID(m.SClassID) {
	case netproto.INRM_GameMessageClassID_INRM_RefreshHandID: // 玩家更新手牌組合
		sclass := netproto.INRM_GameMessageClassID_INRM_RefreshHandRetID

		data := new(netproto.Poker_Hand)
		if err := proto.Unmarshal(m.MsgData.([]byte), data); err != nil {
			m := gs.mActionRet(false, err.Error(), sclass)
			gs.SendMessgeFunc(user.UserID, m)
			return
		}

		var hands []*poker.Deck
		for _, d := range data.Decks {
			hands = append(hands, poker.ReversePokerDeck(d))
		}
		if err := gs.Game.Action(user.UserID, RefreshHandsActionOpt(hands...)); err != nil {
			m := gs.mActionRet(false, err.Error(), sclass)
			gs.SendMessgeFunc(user.UserID, m)
			return
		}
	case netproto.INRM_GameMessageClassID_INRM_PopCloseDeckID: // 玩家抽取蓋牌區
		sclass := netproto.INRM_GameMessageClassID_INRM_PopCloseDeckRetID
		if err := gs.Game.Action(user.UserID, PopCloseDeckActionOpt()); err != nil {
			m := gs.mPopCardRet(false, err.Error(), 0, sclass)
			gs.SendMessgeFunc(user.UserID, m)
			return
		}
	case netproto.INRM_GameMessageClassID_INRM_PopOpenDeckID: // 玩家抽取開牌區
		sclass := netproto.INRM_GameMessageClassID_INRM_PopOpenDeckRetID
		if err := gs.Game.Action(user.UserID, PopOpenDeckActionOpt()); err != nil {
			m := gs.mPopCardRet(false, err.Error(), 0, sclass)
			gs.SendMessgeFunc(user.UserID, m)
			return
		}
	case netproto.INRM_GameMessageClassID_INRM_GiveupID: // 玩家棄牌
		sclass := netproto.INRM_GameMessageClassID_INRM_GiveupRetID
		if err := gs.Game.Action(user.UserID, GiveupActionOpt()); err != nil {
			m := gs.mActionRet(false, err.Error(), sclass)
			gs.SendMessgeFunc(user.UserID, m)
			return
		}
	case netproto.INRM_GameMessageClassID_INRM_ThrowCardID: // 玩家丟牌
		sclass := netproto.INRM_GameMessageClassID_INRM_ThrowCardRetID
		data := new(netproto.INRM_ThrowCard)
		if err := proto.Unmarshal(m.MsgData.([]byte), data); err != nil {
			m := gs.mActionRet(false, err.Error(), sclass)
			gs.SendMessgeFunc(user.UserID, m)
			return
		}
		if err := gs.Game.Action(user.UserID, ThrowCardActionOpt(poker.ReversePokerCard(data.Throw))); err != nil {
			m := gs.mActionRet(false, err.Error(), sclass)
			gs.SendMessgeFunc(user.UserID, m)
			return
		}
	case netproto.INRM_GameMessageClassID_INRM_DeclareID: // 玩家宣告勝利
		sclass := netproto.INRM_GameMessageClassID_INRM_DeclareRetID
		data := new(netproto.INRM_Declare)
		if err := proto.Unmarshal(m.MsgData.([]byte), data); err != nil {
			m := gs.mActionRet(false, err.Error(), sclass)
			gs.SendMessgeFunc(user.UserID, m)
			return
		}

		c := poker.ReversePokerCard(data.Throw)
		var hands []*poker.Deck
		for _, d := range data.Hands.Decks {
			hands = append(hands, poker.ReversePokerDeck(d))
		}

		if err := gs.Game.Action(user.UserID, VictoryActionOpt(c, hands...)); err != nil {
			m := gs.mActionRet(false, err.Error(), sclass)
			gs.SendMessgeFunc(user.UserID, m)
			return
		}
	case netproto.INRM_GameMessageClassID_INRM_MissID: // 玩家Miss
		sclass := netproto.INRM_GameMessageClassID_INRM_MissRetID
		if err := gs.Game.Action(user.UserID, MissActionOpt()); err != nil {
			m := gs.mActionRet(false, err.Error(), sclass)
			gs.SendMessgeFunc(user.UserID, m)
			return
		}
	}
}

// Start 開始
func (gs *GameServer) Start() {

	gs.StartAt = time.Now()
	gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerWaitPlayer))
	gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerStartCD))

	users := gs.AllGameUser()
	var pids []int32
	for _, u := range users {
		pids = append(pids, u.UserID)
	}

	game, err := NewGame(pids, gs.GameOptions...)
	if err != nil {
		mlog.Error("牌桌[%d] 建立遊戲失敗, %v", gs.TableManager.TableNo, err)
		gs.GameEndFunc()
		return
	}
	game.OnGameStart = gs.OnStart
	game.OnRefreshHands = gs.OnRefreshHands
	game.OnPlayerPopClose = gs.OnPlayerPopClose
	game.OnPlayerPopOpen = gs.OnPlayerPopOpen
	game.OnPlayerThrow = gs.OnPlayerThrow
	game.OnPlayerGiveup = gs.OnPlayerGiveup
	game.OnPlayerMiss = gs.OnPlayerMiss
	game.OnPlayerBluff = gs.OnPlayerBluff
	game.OnPlayerVictory = gs.OnPlayerVictory
	game.OnPlayerPass = gs.OnPlayerPass
	game.OnDeckReset = gs.OnDeckReset
	game.OnGameOver = gs.OnGameOver

	gs.Game = game
	if err := gs.Game.Start(); err != nil {
		mlog.Error("牌桌[%d] 遊戲開始失敗, %v", gs.TableManager.TableNo, err)
		gs.GameEndFunc()
		return
	}
}

// Clean 清理遊戲資訊
func (gs *GameServer) Clean() {
	gs.Game = nil
}

// OnStart 當遊戲開始
func (gs *GameServer) OnStart() {
	{ // 開始第一人的等待時間
		p, _ := gs.Game.CurrentPlayer()
		gs.timerTimerWaitActionStart(p)
		gs.tickerActionCDStart(p)
	}
	gs.GameLog.GameBegin()

	gs.GameLog.AddGameStep("每點: %s 分", humanize.Commaf(float64(gs.ScorePerPoint)/100))
	gs.GameLog.AddGameStep("指定鬼牌點數: %v", gs.Game.RandJoker)

	gs.GameLog.AddGameStep("參與玩家:")
	for _, p := range gs.Game.Players {
		u := gs.GetUserByUserID(p.ID)
		var uType string
		if u.IsRobot() {
			uType = "[機]"
		}
		gs.GameLog.AddGameStep("  玩家[%d]%s(%s): %v", p.ID, uType, humanize.Commaf(float64(u.Score)/100), p.Deck)
		// 發送告知玩家手牌
		m := gs.mPlayerCards(p)
		gs.SendMessgeFunc(p.ID, m)
	}

	gs.LastRound = gs.Game.Round
	gs.GameLog.AddGameStep("第 %02d 局:", gs.Game.Round)

	gs.BoardcastMessage(gs.mGameState())
}

// OnRefreshHands 當玩家更新手牌
func (gs *GameServer) OnRefreshHands(p *Player) {
	m := gs.mRefreshHandRet(true, "", p.Hands)
	gs.SendMessgeFunc(p.ID, m)
}

// OnPlayerPopClose 玩家抽蓋牌區
func (gs *GameServer) OnPlayerPopClose(c poker.Card, p *Player) {
	u := gs.GetUserByUserID(p.ID)
	var uType string
	if u.IsRobot() {
		uType = "[機]"
	}
	gs.GameLog.AddGameStep("  玩家[%d]%s 抽取蓋牌區, 抽到 %v", p.ID, uType, c)
	mlog.Debug("  玩家[%d]%s 抽取蓋牌區, 抽到 %v", p.ID, uType, c)

	{ // 發送玩家抽牌封包
		m := gs.mPopCardRet(true, "", p.LastPop, netproto.INRM_GameMessageClassID_INRM_PopCloseDeckRetID)
		gs.SendMessgeFunc(p.ID, m)
	}

	{ // 通知其他玩家
		gs.BoardcastMessage(gs.mPopCloseDeck(p.ID))
	}
}

// OnPlayerPopOpen 玩家抽亮牌區
func (gs *GameServer) OnPlayerPopOpen(c poker.Card, p *Player) {
	u := gs.GetUserByUserID(p.ID)
	var uType string
	if u.IsRobot() {
		uType = "[機]"
	}
	gs.GameLog.AddGameStep("  玩家[%d]%s 抽取亮牌區, 抽到 %v", p.ID, uType, c)
	mlog.Debug("  玩家[%d]%s 抽取亮牌區, 抽到 %v", p.ID, uType, c)

	{ // 發送玩家抽牌封包
		m := gs.mPopCardRet(true, "", p.LastPop, netproto.INRM_GameMessageClassID_INRM_PopOpenDeckRetID)
		gs.SendMessgeFunc(p.ID, m)
	}

	{ // 通知其他玩家
		gs.BoardcastMessage(gs.mPopOpenDeck(p.ID))
	}
}

// OnDeckReset 當蓋牌區用完，重新洗牌時
func (gs *GameServer) OnDeckReset() {
	m := gs.mDeckReset()
	gs.BoardcastMessage(m)
}

// OnPlayerThrow 玩家丟牌
func (gs *GameServer) OnPlayerThrow(c poker.Card, p *Player) {
	u := gs.GetUserByUserID(p.ID)
	var uType string
	if u.IsRobot() {
		uType = "[機]"
	}
	gs.GameLog.AddGameStep("  玩家[%d]%s 丟一張 %v 到亮牌區", p.ID, uType, c)
	mlog.Debug("  玩家[%d]%s 丟一張 %v 到亮牌區", p.ID, uType, c)

	{ // 發送玩家丟牌封包
		m := gs.mActionRet(true, "", netproto.INRM_GameMessageClassID_INRM_ThrowCardRetID)
		gs.SendMessgeFunc(p.ID, m)
	}

	{ // 通知其他玩家
		gs.BoardcastMessage(gs.mThrowCard(p.ID, c))
	}
}

// OnPlayerGiveup 玩家棄牌
func (gs *GameServer) OnPlayerGiveup(p *Player) {
	u := gs.GetUserByUserID(p.ID)
	var uType string
	if u.IsRobot() {
		uType = "[機]"
	}

	gs.GameLog.AddGameStep("  玩家[%d]%s 放棄牌局", p.ID, uType)
	mlog.Debug("牌局[%02d] 玩家[%d]%s 放棄牌局", gs.Game.Round, p.ID, uType)

	{ // 發送玩家棄牌封包
		m := gs.mActionRet(true, "", netproto.INRM_GameMessageClassID_INRM_GiveupRetID)
		gs.SendMessgeFunc(p.ID, m)
	}

	{ // 通知其他玩家
		gs.BoardcastMessage(gs.mGiveup(p.ID))
	}
}

// OnPlayerMiss 玩家Miss
func (gs *GameServer) OnPlayerMiss(p *Player) {
	u := gs.GetUserByUserID(p.ID)
	var uType string
	if u.IsRobot() {
		uType = "[機]"
	}

	gs.GameLog.AddGameStep("  玩家[%d]%s 動作超時, 已累計 %d 次", p.ID, uType, p.Miss)

	mlog.Debug("牌局[%02d] 玩家[%d]%s 動作超時, 已累計 %d 次", gs.Game.Round, p.ID, uType, p.Miss)
	{ // 發送玩家Miss封包
		m := gs.mActionRet(true, "", netproto.INRM_GameMessageClassID_INRM_MissRetID)
		gs.SendMessgeFunc(p.ID, m)
	}

	{ // 通知其他玩家
		gs.BoardcastMessage(gs.mMiss(p.ID))
	}
}

// OnPlayerBluff 玩家詐胡
func (gs *GameServer) OnPlayerBluff(p *Player, c poker.Card, hands []*poker.Deck) {
	u := gs.GetUserByUserID(p.ID)
	var uType string
	if u.IsRobot() {
		uType = "[機]"
	}

	gs.GameLog.AddGameStep("  玩家[%d]%s 宣告勝利 詐胡, 手牌 %v", p.ID, uType, hands)

	{ // 發送玩家詐胡封包
		m := gs.mActionRet(false, "詐胡", netproto.INRM_GameMessageClassID_INRM_DeclareRetID)
		gs.SendMessgeFunc(p.ID, m)
	}

	{ // 通知其他玩家
		gs.BoardcastMessage(gs.mBluff(p.ID))
	}
}

// OnPlayerVictory 玩家勝利
func (gs *GameServer) OnPlayerVictory(p *Player, c poker.Card, hands []*poker.Deck) {
	u := gs.GetUserByUserID(p.ID)
	var uType string
	if u.IsRobot() {
		uType = "[機]"
	}

	gs.GameLog.AddGameStep("  玩家[%d]%s 宣告勝利 勝利, 手牌 %v", p.ID, uType, hands)

	{ // 發送玩家勝利封包
		m := gs.mActionRet(true, "", netproto.INRM_GameMessageClassID_INRM_DeclareRetID)
		gs.SendMessgeFunc(p.ID, m)
	}

	{ // 通知其他玩家
		gs.BoardcastMessage(gs.mDeclare(p.ID, c, hands))
	}

	gs.TimerFuncManager.RemoveTimer(TimerWaitAction)
	gs.TimerFuncManager.RemoveTimer(TimerActionCD)

	// 等待玩家最後更新手牌
	mlog.Debug("牌桌[%d] 計時等待玩家最後更新手牌", gs.TableManager.TableNo)
	gs.FinalAt = time.Now().Add(gs.ActionWaitTime)
	gs.timerTimerWaitCalcStart(p)
	gs.tickerCalcCDStart()
}

// OnPlayerPass 當玩家換人
func (gs *GameServer) OnPlayerPass(prev, next *Player) {
	gs.TimerFuncManager.RemoveTimer(TimerWaitAction)
	gs.TimerFuncManager.RemoveTimer(TimerActionCD)

	if next != nil { // 開始下一人的等待時間
		gs.timerTimerWaitActionStart(next)
		gs.tickerActionCDStart(next)
	}
	if gs.LastRound != gs.Game.Round {
		gs.LastRound = gs.Game.Round
		gs.GameLog.AddGameStep("第 %02d 局:", gs.Game.Round)
	}
	gs.BoardcastMessage(gs.mGameState())
}

// OnGameOver 遊戲結束
func (gs *GameServer) OnGameOver() {
	gs.GameLog.AddGameStep("牌局結束")

	// 移除所有計時
	gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerWaitPlayer))
	gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerStartCD))
	gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerWaitCalc))
	gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerCalcCD))
	gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerWaitAction))
	gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerActionCD))

	var point int
	var winner *Player
	ctx := gs.GameLog.Context()

	for _, p := range gs.Game.Players {
		u := gs.GetUserByUserID(p.ID)
		var uType string
		if u.IsRobot() {
			uType = "[機]"
		}

		{ // 移除所有人等待時間
			name := fmt.Sprintf("%s-%d", TimerWaitAction, p.ID)
			mlog.Debug("牌桌[%d] 解除玩家[%d]%s動作超時", gs.TableManager.TableNo, p.ID, uType)
			gs.TimerFuncManager.RemoveTimer(name)

			cdname := fmt.Sprintf("%s-%d", TimerActionCD, p.ID)
			mlog.Debug("牌桌[%d] 解除玩家[%d]動作超時 %s", gs.TableManager.TableNo, p.ID, cdname)
			gs.TimerFuncManager.RemoveTimer(cdname)
		}

		// 計算輸家
		if p.State != PStateVictory {
			point += p.LossPoint

			loss := int64(p.LossPoint) * gs.ScorePerPoint

			u := gs.GetUserByUserID(p.ID)
			u.BetAmount += loss
			u.FlowAmount += loss

			ctx.UserBetScore(p.ID, loss)
			ctx.UserValidScore(p.ID, loss)
			ctx.UserUpdateScore(p.ID, -loss)

			gs.UpdateUserScoreFunc(p.ID, -loss, 0)
			gs.GameLog.AddGameStep("玩家[%d]%s[%v] 輸 %s 分 最終 %s 分，點數 %d 牌型 %v", p.ID, uType, p.State, humanize.Commaf(float64(loss)/100), humanize.Commaf(float64(u.Score)/100), p.LossPoint, p.Hands)
		} else {
			winner = p
			continue
		}
	}

	{ // 計算贏家
		u := gs.GetUserByUserID(winner.ID)
		var uType string
		if u.IsRobot() {
			uType = "[機]"
		}

		win := int64(point) * gs.ScorePerPoint
		tax := win * int64(gs.GetServerConfig().TaxScale) / 100

		u.PayAmount += win - tax

		ctx.UserUpdateScore(winner.ID, win-tax)
		ctx.UserPayoutScore(winner.ID, win)
		ctx.UserTaxAmount(winner.ID, tax)
		gs.UpdateUserScoreFunc(winner.ID, win-tax, tax)

		gs.GameLog.AddGameStep("玩家[%d]%s[%v] 贏 %s 分 稅收 %s 分 最終 %s 分，點數 %d 牌型 %v", winner.ID, uType, winner.State, humanize.Commaf(float64(win)/100), humanize.Commaf(float64(tax)/100), humanize.Commaf(float64(u.Score)/100), winner.LossPoint, winner.Hands)
	}

	// 傳送遊戲結果
	gs.BoardcastMessage(gs.mGameResult())

	gs.GameLog.GameEnd()
	gs.GameEndFunc()
}

// IsGaming 是否遊戲中
func (gs *GameServer) IsGaming() bool {
	return gs.isgaming
}

// OnGameStart 遊戲開始
func (gs *GameServer) OnGameStart(f func()) {
	gs.GameStartFunc = func() {
		if gs.IsGaming() {
			return
		}

		// 遊戲開始前調用函數
		for _, f := range gs.BeforeGameStartFunc {
			f()
		}

		// 通知用戶狀態改變
		for _, usr := range gs.AllTableUser() {
			if usr.IsReady() {
				usr.UserState.OnUserGame(usr, gs.TableManager.Room)
				continue
			}
			gs.UserLeave(usr)
		}

		gs.isgaming = true
		gs.State = Gaming
		gs.GameStartTime = time.Now()
		gs.PlayerActionManager.Clear()

		// 廣播遊戲開始
		gs.BoardcastGameStart()
		mlog.Debug("牌桌[%d] 遊戲開始...", gs.TableManager.TableNo)

		// 牌局記錄
		gs.GameLog.GameBegin()

		f()

		// 遊戲開始後調用函數
		for _, f := range gs.AfterGameStartFunc {
			f()
		}
	}
}

// OnGameEnd 遊戲結束
func (gs *GameServer) OnGameEnd(f func(int32)) {
	gs.GameEndFunc = func() {
		// 遊戲結束前調用函數
		for _, f := range gs.BeforeGameEndFunc {
			f()
		}

		gs.isgaming = false
		gs.State = WaitPlayer
		gs.LastGameEndTime = time.Now()
		gs.PlayerActionManager.Clear()

		// 廣播遊戲結束
		gs.BoardcastGameEnd()

		for _, usr := range gs.AllGameUser() {
			usr.UserState.OnGameEnd(usr, gs.TableManager.Room)

			if usr.IsOffLine() {
				continue
			}
			// 加入準備超時檢查
			gs.CheckUserReadyTimer(usr)
		}

		gametime := int32(time.Now().Sub(gs.GameStartTime).Seconds())

		f(gametime)

		// 遊戲結束後調用函數
		for _, f := range gs.AfterGameEndFunc {
			f()
		}

		mlog.Debug("牌桌[%d] 遊戲結束...", gs.TableManager.TableNo)
	}
}

// UserEnter 玩家進入
func (gs *GameServer) UserEnter(user *gamelogic.GameUser, tableno, seatno int32) {
	if gs.IsGaming() {
		return
	}

	// 用戶坐下
	user.UserState.OnUserSit(user, tableno, seatno, gs.TableManager.Room)

	// 通知用戶進入
	gs.BoardcastUserEnter(user)

	gs.UserReady(user)

	// 傳送遊戲設定
	gs.SendMessgeFunc(user.UserID, gs.mGameSetting())
	// 傳送遊戲狀態
	gs.SendMessgeFunc(user.UserID, gs.mGameStateLite())
}

// UserRecover 斷線返回
func (gs *GameServer) UserRecover(user *gamelogic.GameUser) {
	gs.LogicBase.UserRecover(user)

	p, err := gs.Game.Player(user.UserID)
	if err != nil {
		mlog.Warn("玩家[%d] 斷線返回, 但是取得玩家錯誤, %v", user.UserID, err)
		return
	}

	// 傳送遊戲設定
	gs.SendMessgeFunc(user.UserID, gs.mGameSetting())
	// 傳送遊戲狀態
	gs.SendMessgeFunc(user.UserID, gs.mGameState())
	// 傳送玩家手牌
	gs.SendMessgeFunc(user.UserID, gs.mPlayerCards(p))
}

// UserLeave 用戶離開
func (gs *GameServer) UserLeave(user *gamelogic.GameUser) bool {

	if gs.IsGaming() {
		return false
	}

	// 通知用戶離開
	user.UserState.OnUserLeave(user, gs.TableManager.Room)

	// 廣播用戶離開
	gs.BoardcastUserLeave(user)

	// 檢查人數超過2人就開始倒數
	if gs.State == CountDown && gs.GetGameReadyUserCount() < gs.GetServerConfig().MinPlayer {
		mlog.Debug("牌桌[%d] 停止倒數計時遊戲開始", gs.TableManager.TableNo)
		gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerWaitPlayer))
		gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerStartCD))
		gs.State = WaitPlayer
		gs.BoardcastMessage(gs.mGameStartCountDownCancel())
	}

	// 檢查是否有活人，沒有就把機器人踢光
	if len(gs.TableManager.GetTableRealUser()) == 0 && !user.IsRobot() {
		for _, p := range gs.TableManager.GetTableRealRobot() {
			example.invalid("牌桌[%d] 已無真實玩家，機器人[%d]離開", gs.TableManager.TableNo, p.UserID)
			gs.UserLeave(p)
		}
	}

	return true
}

// UserReady 用戶準備
func (gs *GameServer) UserReady(user *gamelogic.GameUser) {

	user.UserState.OnUserReady(user, gs.TableManager.Room)

	gs.BoardcastUserStateChange(user)

	// 檢查人數超過2人就開始倒數
	gs.checkGameStartable()
}

func (gs *GameServer) checkGameStartable() {
	if gs.State == WaitPlayer && gs.GetGameReadyUserCount() >= gs.GetServerConfig().MinPlayer {
		gs.timerWaitPlayerStart()
		gs.State = CountDown
		gs.StartAt = time.Now().Add(gs.StartWaitTime / time.Duration(conf.TimerFastMode))
		gs.BoardcastMessage(gs.mGameStartCountDown())
		mlog.Debug("牌桌[%d] 開始倒數計時遊戲開始 %s", gs.TableManager.TableNo, gs.StartAt.Format(time.RFC3339))
		gs.tickerStartCDStart()
	}
}

func (gs *GameServer) timerWaitPlayerStart() {
	tf := gs.TimerFuncManager.GetTimerFunc(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerWaitPlayer))
	if tf != nil {
		gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerWaitPlayer))
	}
	gs.TimerFuncManager.NewNamedTimer(int32(gs.StartWaitTime.Seconds()), gs.GameStartFunc, fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerWaitPlayer))
}

func (gs *GameServer) tickerStartCDStart() {
	tf := gs.TimerFuncManager.GetTimerFunc(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerStartCD))
	if tf != nil {
		gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerStartCD))
	}
	gs.TimerFuncManager.NewNamedTicker(1, func() { gs.BoardcastMessage(gs.mGameStateLite()) }, fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerStartCD))
}

func (gs *GameServer) timerTimerWaitActionStart(p *Player) {
	if p == nil {
		return
	}
	tf := gs.TimerFuncManager.GetTimerFunc(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerWaitAction))
	if tf != nil {
		gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerWaitAction))
	}
	gs.TimerFuncManager.NewNamedTicker(int32(gs.ActionWaitTime.Seconds()), func() { gs.Game.Action(p.ID, MissActionOpt()) }, fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerWaitAction))
}

func (gs *GameServer) tickerActionCDStart(p *Player) {
	if p == nil {
		return
	}
	gs.ActionTimeoutAt = time.Now().Add(gs.ActionWaitTime)
	tf := gs.TimerFuncManager.GetTimerFunc(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerActionCD))
	if tf != nil {
		gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerActionCD))
	}
	gs.TimerFuncManager.NewNamedTicker(1, func() { gs.BoardcastMessage(gs.mPlayerActionCD(p)) }, fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerActionCD))
}

func (gs *GameServer) timerTimerWaitCalcStart(p *Player) {
	if p == nil {
		return
	}

	tf := gs.TimerFuncManager.GetTimerFunc(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerWaitCalc))
	if tf != nil {
		gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerWaitCalc))
	}
	gs.TimerFuncManager.NewNamedTicker(int32(gs.ActionWaitTime.Seconds()), func() { gs.Game.Action(p.ID, GameOverActionOpt()) }, fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerWaitCalc))
}

func (gs *GameServer) tickerCalcCDStart() {
	tf := gs.TimerFuncManager.GetTimerFunc(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerCalcCD))
	if tf != nil {
		gs.TimerFuncManager.RemoveTimer(fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerCalcCD))
	}
	gs.TimerFuncManager.NewNamedTicker(1, func() { gs.BoardcastMessage(gs.mWaitFinalRefreshCD()) }, fmt.Sprintf("%d-%s", gs.TableManager.TableNo, TimerCalcCD))
}

func (gs *GameServer) mGameStartCountDown() *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_GameStartCountDownID)
	return m
}

func (gs *GameServer) mGameStartCountDownCancel() *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_GameStartCountDownCancelID)
	return m
}

func (gs *GameServer) mPlayerActionCD(p *Player) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PlayerActionCountDownID)

	data := new(netproto.INRM_Player)
	data.ID = p.ID
	data.ActionCD = int64(math.Round(float64(gs.ActionTimeoutAt.Sub(time.Now()).Milliseconds())/1000)) * 1000

	m.MsgData = data
	return m
}

func (gs *GameServer) mPlayerActionCDCancel(p *Player) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PlayerActionCountDownCancelID)
	return m
}

func (gs *GameServer) mWaitFinalRefreshCD() *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_WaitFinalRefreshCountDownID)

	data := new(netproto.INRM_TimeCD)
	data.CD = int64(math.Round(float64(gs.FinalAt.Sub(time.Now()).Milliseconds())/1000)) * 1000

	m.MsgData = data
	return m
}

func (gs *GameServer) mActionRet(success bool, message string, sclass netproto.INRM_GameMessageClassID) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(sclass)

	data := new(netproto.INRM_ActionRet)
	data.Success = success
	data.Message = message

	m.MsgData = data

	return m
}

func (gs *GameServer) mPopCardRet(success bool, message string, c poker.Card, sclass netproto.INRM_GameMessageClassID) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(sclass)

	data := new(netproto.INRM_PopCardRet)
	data.Success = success
	data.Message = message
	data.Card = poker.TransePokerCard(c)

	m.MsgData = data

	return m
}

func (gs *GameServer) mRefreshHandRet(success bool, message string, hands []*poker.Deck) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandRetID)

	data := new(netproto.INRM_RefreshHandRet)
	data.Success = success
	data.Message = message
	data.Result = gs.mDataHandResult(hands)

	m.MsgData = data

	return m
}

func (gs *GameServer) mPopCloseDeck(pid int32) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopCloseDeckHappenedID)
	data := new(netproto.INRM_Player)
	data.ID = pid
	m.MsgData = data
	return m
}

func (gs *GameServer) mPopOpenDeck(pid int32) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopOpenDeckHappenedID)
	data := new(netproto.INRM_Player)
	data.ID = pid
	m.MsgData = data
	return m
}

func (gs *GameServer) mGiveup(pid int32) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_GiveupHappenedID)
	data := new(netproto.INRM_Player)
	data.ID = pid
	m.MsgData = data
	return m
}

func (gs *GameServer) mMiss(pid int32) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_MissHappenedID)
	data := new(netproto.INRM_Player)
	data.ID = pid
	m.MsgData = data
	return m
}

func (gs *GameServer) mBluff(pid int32) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_BluffHappenedID)
	data := new(netproto.INRM_Player)
	data.ID = pid
	m.MsgData = data
	return m
}

func (gs *GameServer) mThrowCard(pid int32, c poker.Card) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_ThrowCardHappenedID)

	data := new(netproto.INRM_ThrowCard)
	data.Throw = poker.TransePokerCard(c)
	data.PlayerID = pid

	m.MsgData = data

	return m
}

func (gs *GameServer) mDeclare(pid int32, c poker.Card, hands []*poker.Deck) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_DeclareHappenedID)

	data := new(netproto.INRM_Declare)
	data.Throw = poker.TransePokerCard(c)
	h := new(netproto.Poker_Hand)
	for _, d := range hands {
		h.Decks = append(h.Decks, poker.TransePokerDeck(d))
	}
	data.Hands = h
	data.PlayerID = pid

	m.MsgData = data

	return m
}

func (gs *GameServer) mGameSetting() *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_GameSettingID)

	data := new(netproto.INRM_GameSetting)
	data.TableNo = gs.TableManager.TableNo
	data.StartWaitTime = int32(gs.StartWaitTime.Milliseconds())
	data.ActionWaitTime = int32(gs.ActionWaitTime.Milliseconds())
	data.ScorePerPoint = gs.ScorePerPoint

	data.HandCardLength = 13
	data.MaxPoint = 80
	data.GiveupPoint = 40
	data.MissLimit = 3

	for _, opt := range gs.GameOptions {
		if o, ok := opt.(*handCardLengthGameOpt); ok {
			data.HandCardLength = int64(o.length)
			continue
		}
		if o, ok := opt.(*maxPointGameOpt); ok {
			data.MaxPoint = int64(o.point)
			continue
		}
		if o, ok := opt.(*giveupPointGameOpt); ok {
			data.GiveupPoint = int64(o.point)
			continue
		}
		if o, ok := opt.(*missGameOpt); ok {
			data.MissLimit = int32(o.n)
			continue
		}
	}

	m.MsgData = data

	return m
}

func (gs *GameServer) mGameState() *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_GameStateID)

	data := new(netproto.INRM_GameState)
	data.State = gs.mTransTableState(gs.State)
	switch gs.State {
	case CountDown:
		data.StartAt = gs.StartAt.Unix()
		cd := gs.StartAt.Sub(time.Now())
		data.StartCD = cd.Milliseconds()
	case Gaming:
		data.StartAt = gs.StartAt.Unix()
	}
	data.Round = int32(gs.Game.Round)
	p, err := gs.Game.CurrentPlayer()
	if err == nil {
		data.PlayerID = p.ID
	}
	data.RandJoker = netproto.Poker_Rank(gs.Game.RandJoker)
	data.CloseDeckCount = int32(gs.Game.CloseDeck.Len())
	data.OpenDeckCount = int32(gs.Game.OpenDeck.Len())
	data.FirstOpenDeck = poker.TransePokerCard(gs.Game.OpenDeck.Cards[0])

	for _, p := range gs.Game.Players {
		data.Players = append(data.Players, gs.mDataPlayerLite(p))
	}

	m.MsgData = data

	return m
}

func (gs *GameServer) mGameStateLite() *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_GameStateID)

	data := new(netproto.INRM_GameState)
	data.State = gs.mTransTableState(gs.State)
	switch gs.State {
	case CountDown:
		data.StartAt = gs.StartAt.Unix()
		cd := gs.StartAt.Sub(time.Now())

		// 取整數
		data.StartCD = int64(math.Round(float64(cd.Milliseconds())/1000)) * 1000
	case Gaming:
		data.StartAt = gs.StartAt.Unix()
	}

	m.MsgData = data

	return m
}

func (gs *GameServer) mGameResult() *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_GameResultID)

	ctx := gs.GameLog.Context()

	data := new(netproto.INRM_GameResult)
	for _, p := range gs.Game.Players {
		ud := ctx.Users[p.ID]
		dataP := gs.mDataPlayer(p, ud.UpdateScore, ud.PayoutScore, ud.TaxAmount)
		if p.State == PStateVictory {
			data.Winner = dataP
			continue
		}
		data.Losser = append(data.Losser, dataP)
	}

	m.MsgData = data

	return m
}

func (gs *GameServer) mPlayerCards(p *Player) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PlayerCardID)

	data := new(netproto.INRM_PlayerCards)

	p.Deck.Range(func(i int, c poker.Card) bool {
		data.Cards = append(data.Cards, poker.TransePokerCard(c))
		return true
	})
	data.Decks = append(data.Decks, gs.mDataHandResult(p.Hands).Decks...)

	m.MsgData = data

	return m
}

func (gs *GameServer) mDeckReset() *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_DeckResetID)
	return m
}

func (gs *GameServer) mTransTableState(s TableState) netproto.INRM_GState {
	switch s {
	default:
		return netproto.INRM_GState_INRM_GState_UNSPECIFIED
	case WaitPlayer:
		return netproto.INRM_GState_INRM_GState_WaitPlayer
	case CountDown:
		return netproto.INRM_GState_INRM_GState_CountDown
	case Gaming:
		return netproto.INRM_GState_INRM_GState_Gaming
	}
}

func (gs *GameServer) mTransPlayerState(s PState) netproto.INRM_PState {
	switch s {
	default:
		return netproto.INRM_PState_INRM_PState_UNSPECIFIED
	case PStateWait:
		return netproto.INRM_PState_INRM_PState_Wait
	case PStateGameV1:
		return netproto.INRM_PState_INRM_PState_GameV1
	case PStateGameV2:
		return netproto.INRM_PState_INRM_PState_GameV2
	case PStateFirstGiveup:
		return netproto.INRM_PState_INRM_PState_FirstGiveup
	case PStateGiveup:
		return netproto.INRM_PState_INRM_PState_Giveup
	case PStateBluff:
		return netproto.INRM_PState_INRM_PState_Bluff
	case PStateVictory:
		return netproto.INRM_PState_INRM_PState_Victory
	case PStateLoss:
		return netproto.INRM_PState_INRM_PState_Loss
	}
}

func (gs *GameServer) mDataPlayerLite(p *Player) *netproto.INRM_Player {
	data := new(netproto.INRM_Player)
	data.ID = p.ID
	data.State = gs.mTransPlayerState(p.State)
	data.Miss = int32(p.Miss)
	return data
}

func (gs *GameServer) mDataPlayer(p *Player, update, win, tax int64) *netproto.INRM_Player {
	data := gs.mDataPlayerLite(p)
	data.Hands = gs.mDataHandResult(p.Hands)
	data.LossPoint = int32(p.LossPoint)
	data.UpdateScore = update
	data.Win = win
	data.Tax = tax
	return data
}

func (gs *GameServer) mDataHandResult(hands []*poker.Deck) *netproto.INRM_HandResult {
	data := new(netproto.INRM_HandResult)

	pr := NewPointResult(gs.Game.RandJoker, hands...)
	data.Point = int32(pr.Points())
	if data.Point > int32(gs.Game.MaxPoint) {
		data.Point = int32(gs.Game.MaxPoint)
	}
	if pr.FirstLife != nil {
		iDeck := new(netproto.INRM_Deck)
		iDeck.Deck = poker.TransePokerDeck(pr.FirstLife)
		iDeck.Category = netproto.INRM_Category_INRM_Category_1stLife
		data.Decks = append(data.Decks, iDeck)
	}
	if pr.SecondLife != nil {
		iDeck := new(netproto.INRM_Deck)
		iDeck.Deck = poker.TransePokerDeck(pr.SecondLife)
		iDeck.Category = netproto.INRM_Category_INRM_Category_2stLife
		data.Decks = append(data.Decks, iDeck)
	}
	for _, d := range pr.PureSequence {
		if d == nil {
			continue
		}
		iDeck := new(netproto.INRM_Deck)
		iDeck.Deck = poker.TransePokerDeck(d)
		if pr.FirstLife == nil {
			iDeck.Category = netproto.INRM_Category_INRM_Category_1stLifeNeeded
		} else if pr.FirstLife != nil && pr.SecondLife == nil {
			iDeck.Category = netproto.INRM_Category_INRM_Category_2stLifeNeeded
		} else {
			iDeck.Category = netproto.INRM_Category_INRM_Category_PureSequence
		}

		data.Decks = append(data.Decks, iDeck)
	}
	for _, d := range pr.ImpureSequence {
		if d == nil {
			continue
		}
		iDeck := new(netproto.INRM_Deck)
		iDeck.Deck = poker.TransePokerDeck(d)
		if pr.FirstLife == nil {
			iDeck.Category = netproto.INRM_Category_INRM_Category_1stLifeNeeded
		} else if pr.FirstLife != nil && pr.SecondLife == nil {
			iDeck.Category = netproto.INRM_Category_INRM_Category_2stLifeNeeded
		} else {
			iDeck.Category = netproto.INRM_Category_INRM_Category_ImpureSequence
		}

		data.Decks = append(data.Decks, iDeck)
	}
	for _, d := range pr.Set {
		if d == nil {
			continue
		}
		iDeck := new(netproto.INRM_Deck)
		iDeck.Deck = poker.TransePokerDeck(d)
		if pr.FirstLife == nil {
			iDeck.Category = netproto.INRM_Category_INRM_Category_1stLifeNeeded
		} else if pr.FirstLife != nil && pr.SecondLife == nil {
			iDeck.Category = netproto.INRM_Category_INRM_Category_2stLifeNeeded
		} else {
			iDeck.Category = netproto.INRM_Category_INRM_Category_Set
		}

		data.Decks = append(data.Decks, iDeck)
	}
	for _, d := range pr.None {
		if d == nil {
			continue
		}
		iDeck := new(netproto.INRM_Deck)
		iDeck.Deck = poker.TransePokerDeck(d)
		if d.Len() > 2 {
			iDeck.Category = netproto.INRM_Category_INRM_Category_Invalid
		} else {
			iDeck.Category = netproto.INRM_Category_INRM_Category_None
		}

		data.Decks = append(data.Decks, iDeck)
	}

	return data
}
