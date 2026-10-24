package infc

import (
	"context"
	"runtime/debug"
	"sort"
	"strconv"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// DeathGame 玩到死
//  進行切水果遊戲取得分數
//  * 切到1個水果5分
//  * 連切到2個水果15分
//  * 連切到3個水果35分
//  * 依此類推
//  禁止切炸彈，切到立即結束遊戲
//  當有3個水果沒被切到掉出畫面，立即結束遊戲
//  等所有參與的玩家都結束遊戲，比較分數判斷輸贏
//  如果對手提早結束遊戲，當你進行到分數高於對手，則直接結束遊戲，宣告玩家獲勝
type DeathGame struct {
	*gamelogic.LogicBase

	Proto *Proto

	TaxRate       float64
	MissLimit     int64
	BetAmount     int64
	StartDuration time.Duration
	StartAt       time.Time

	isgaming bool
	ctx      *RoundContext
}

// NewDeathGame create game logic
func NewDeathGame(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group, proto *Proto) *DeathGame {
	game := new(DeathGame)
	game.LogicBase = gamelogic.NewGameLogicBase(t, tfManager, cfgs)
	game.Proto = proto
	game.initConfig(cfgs)

	game.RegisterAfterGameStartFunc(func() {
		game.ctx = NewRoundContext(game.AllGameUser())
		game.Start(game.ctx)
	})

	return game
}

func (game *DeathGame) initConfig(cfgs *config.Group) {
	game.TaxRate = float64(game.GetServerConfig().TaxScale) / 100
	game.BetAmount = int64(game.GetServerConfig().BaseScore)

	sec := cfgs.GetConfigSection("gameserver", "gameserver")

	if val, err := strconv.ParseInt(sec["misslimit"], 0, 64); err != nil {
		game.MissLimit = 3
	} else {
		game.MissLimit = val
	}

	if val, err := time.ParseDuration(sec["startduration"]); err != nil {
		game.StartDuration = 5 * time.Second
	} else {
		game.StartDuration = val
	}
}

// GameMessage on recived client message
func (game *DeathGame) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("GameMessage 發生錯誤, %v, %s", rec, debug.Stack())
		}
	}()

	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	if game.ctx == nil {
		mlog.Warn("收到封包 [%d:%d], 目前沒有比賽進行中", gmsg.BClassID, gmsg.SClassID)
		return
	}
	if game.ctx.GameState != GameStateGaming {
		mlog.Warn("收到封包 [%d:%d], 目前不在比賽中", gmsg.BClassID, gmsg.SClassID)
		return
	}

	switch netproto.INFC_GameMessageClassID(gmsg.SClassID) {
	case netproto.INFC_GameMessageClassID_INFC_QueryGameSettingID:
		m := game.mGameSetting(game.StartAt, game.ctx)
		game.SendMessgeFunc(user.UserID, m)
	case netproto.INFC_GameMessageClassID_INFC_QueryProtoSettingID:
		m := game.mProtoSetting()
		game.SendMessgeFunc(user.UserID, m)
	case netproto.INFC_GameMessageClassID_INFC_IncScoreID:
		p := game.ctx.Player(user.UserID)
		if p.State != StateGaming {
			return
		}
		data := new(netproto.INFC_IncScore)
		err := proto.Unmarshal(gmsg.MsgData.([]byte), data)
		if err != nil {
			mlog.Warn("封包解析錯誤錯誤 [%d:%d], %v", gmsg.BClassID, gmsg.SClassID, err)
			return
		}
		mlog.Debug("玩家[%d] 回報遊戲結果, 加 %d 分", user.UserID, data.Score)
		game.ctx.AddScore(user.UserID, data.Score)

		game.checkFinish(game.ctx)
		if game.ctx.GameState == GameStateGaming {
			m := game.mRank(game.ctx)
			game.BoardcastMessage(m)
		}
	case netproto.INFC_GameMessageClassID_INFC_MissID:
		p := game.ctx.Player(user.UserID)
		if p.State != StateGaming {
			return
		}
		mlog.Debug("玩家[%d] 回報Miss結果, 加 %d 次, 累計: %d 次", user.UserID, 1, game.ctx.Player(user.UserID).Miss+1)
		game.GameLog.AddGameStep("玩家[%d] (M) 回報Miss結果, 加 %d 次, 累計: %d 次", user.UserID, 1, game.ctx.Player(user.UserID).Miss+1)
		game.ctx.AddMiss(user.UserID, 1)

		game.checkFinish(game.ctx)
		if game.ctx.GameState == GameStateGaming {
			m := game.mRank(game.ctx)
			game.BoardcastMessage(m)
		}
	case netproto.INFC_GameMessageClassID_INFC_BoomID:
		p := game.ctx.Player(user.UserID)
		if p.State != StateGaming {
			return
		}
		mlog.Debug("玩家[%d] 刷到炸彈，掛了, 死於 %d 分", p.UID, p.Score)
		game.GameLog.AddGameStep("玩家[%d] 刷到炸彈，掛了, 死於 %d 分", p.UID, p.Score)
		p.State = StateDeath

		game.checkFinish(game.ctx)
		if game.ctx.GameState == GameStateGaming {
			m := game.mRank(game.ctx)
			game.BoardcastMessage(m)
		}
	case netproto.INFC_GameMessageClassID_INFC_GiveupID:
		p := game.ctx.Player(user.UserID)
		mlog.Debug("玩家[%d] 棄權於 %d 分", p.UID, p.Score)
		game.GameLog.AddGameStep("玩家[%d] 棄權於 %d 分", p.UID, p.Score)
		p.State = StateGiveup

		game.checkFinish(game.ctx)
		if game.ctx.GameState == GameStateGaming {
			m := game.mRank(game.ctx)
			game.BoardcastMessage(m)
		}
	}
}

// PiSaiMessage on recived client message
func (game *DeathGame) PiSaiMessage(user *gamelogic.GameUser, gmsg *msg.Message) {}

// Start 開始遊戲
func (game *DeathGame) Start(ctx *RoundContext) {

	game.GameLog.GameBegin()

	now := time.Now()
	start := now.Add(game.StartDuration / time.Duration(conf.TimerFastMode))
	game.StartAt = start

	{ // Game Setting
		m := game.mGameSetting(start, ctx)
		game.BoardcastMessage(m)
	}
	{ // Proto Setting
		m := game.mProtoSetting()
		game.BoardcastMessage(m)
	}

	for _, u := range game.AllGameUser() {
		p := ctx.Player(u.UserID)
		p.State = StateGaming
		game.Bet(ctx, u)
	}
	game.mBet(game.BetAmount)

	game.GameLog.AddGameStep("遊戲水池 %d", ctx.Pool)

	startCtx, cancel := context.WithDeadline(ctx.ctx, start)
	defer cancel()
	<-startCtx.Done()

	ctx.GameState = GameStateGaming
	game.mStart()

	go func() { // 等待遊戲結束
		select {
		case <-ctx.ctx.Done():
			game.Calc(ctx)
			game.GameLog.GameEnd()
			game.GameEndFunc()
		}
	}()
}

// 檢查玩家狀態判斷是否須終止遊戲
func (game *DeathGame) checkFinish(ctx *RoundContext) {
	if ctx.GameState != GameStateGaming {
		return
	}

	for _, p := range ctx.Rank {
		if p.State == StateGiveup {
			continue
		}
		if p.State != StateDeath && p.Miss >= game.MissLimit {
			game.GameLog.AddGameStep("玩家[%d] (死) 失誤次數達到上限, 死於 %d 分", p.UID, p.Score)
			p.State = StateDeath
			// game.mDeath(p.UID, "失誤次數達到上限")
		}
	}

	// 取得活著的人
	survivors := ctx.Rank.Survivors()
	if len(survivors) == 1 {
		// 其餘玩家皆放棄或離線
		if ctx.Rank.IsOthersIncalculable(survivors[0]) {
			game.GameLog.AddGameStep("遊戲僅存一位玩家[%d], 其餘玩家皆不可被算分, 結束遊戲", survivors[0].UID)
			ctx.cancel()
			return
		}
		// 存活者分數達到最高
		if survivors[0].Score == ctx.Rank.TopScore() {
			game.GameLog.AddGameStep("遊戲僅存一位玩家[%d], 分數達到第一名 %d, 結束遊戲", survivors[0].UID, ctx.Rank.TopScore())
			ctx.cancel()
			return
		}
	}

	// 所有玩家都斷線
	if len(survivors) == 0 {
		game.GameLog.AddGameStep("已經沒有存活玩家, 結束遊戲")
		example.invalid("桌子[%d] 已經沒有存活玩家, 結束遊戲", game.TableManager.TableNo)
		ctx.cancel()
		return
	}

	for _, u := range game.AllTableUser() {
		p := ctx.Player(u.UserID)
		if p.State == StateGaming {
			p.LastTime = time.Now()
		}
	}
}

// Bet 玩家下注
func (game *DeathGame) Bet(ctx *RoundContext, user *gamelogic.GameUser) {

	game.GameLog.AddGameStep("玩家[%d] 支付遊戲底分 %d", user.UserID, game.BetAmount)

	user.BetAmount = game.BetAmount
	game.GameLog.Context().UserUpdateScore(user.UserID, -game.BetAmount)
	game.GameLog.Context().UserBetScore(user.UserID, game.BetAmount)
	ctx.Pool += game.BetAmount

	game.UpdateUserScoreFunc(user.UserID, -game.BetAmount, 0)
}

// Calc 結算
func (game *DeathGame) Calc(ctx *RoundContext) {

	ctx.GameState = GameStateCalculating

	var winner *Player
	topScorePlayers := ctx.Rank.Winner()
	sort.Slice(topScorePlayers, func(i, j int) bool {
		return topScorePlayers[i].LastTime.Before(topScorePlayers[j].LastTime)
	})
	winner = topScorePlayers[0]

	game.GameLog.AddGameStep("==================== 結算 ====================")
	for _, u := range game.AllTableUser() {
		var winloss string

		bet := game.GameLog.Context().Users[u.UserID].BetScore

		if winner == nil {
			winloss = "平手"
		} else if u.UserID == winner.UID {
			win := ctx.Pool
			tax := int64(float64(win) * game.TaxRate)
			pay := win - tax

			u.PayAmount += pay
			game.GameLog.Context().UserUpdateScore(u.UserID, pay)
			game.GameLog.Context().UserPayoutScore(u.UserID, win)
			game.GameLog.Context().UserTaxAmount(u.UserID, tax)
			game.UpdateUserScoreFunc(u.UserID, pay, tax)

			winloss = "勝利"
			game.GameLog.AddGameStep("玩家[%d] 獲得 %8d 分 %s, 注額: %5d, 獎金: %5d, 稅金: %5d, 實收: %5d, 流水: %5d", u.UserID, ctx.Player(u.UserID).Score, winloss, bet, win, tax, game.GameLog.Context().Users[u.UserID].UpdateScore, 0)

		} else {
			winloss = "戰敗"

			u.FlowAmount += bet
			game.GameLog.Context().UserValidScore(u.UserID, bet)

			game.GameLog.AddGameStep("玩家[%d] 獲得 %8d 分 %s, 注額: %5d, 獎金: %5d, 稅金: %5d, 實收: %5d, 流水: %5d", u.UserID, ctx.Player(u.UserID).Score, winloss, bet, 0, 0, game.GameLog.Context().Users[u.UserID].UpdateScore, bet)
		}
	}
	game.GameLog.AddGameStep("==============================================")

	game.mResult(ctx)
}

// UserEnter 玩家進入
func (game *DeathGame) UserEnter(user *gamelogic.GameUser, tableno, seatno int32) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("UserEnter 發生錯誤, %v, %s", rec, debug.Stack())
		}
	}()

	// 用戶坐下
	user.UserState.OnUserSit(user, tableno, seatno, game.TableManager.Room)
	// 用戶準備
	user.UserState.OnUserReady(user, game.TableManager.Room)

	// 通知用戶進入
	game.BoardcastUserEnter(user)

	game.UserReady(user)
}

// UserReady 用戶準備
func (game *DeathGame) UserReady(user *gamelogic.GameUser) {

	user.UserState.OnUserReady(user, game.TableManager.Room)

	game.BoardcastUserStateChange(user)

	// 檢查所有用戶都準備及開始
	if len(game.TableManager.ReadyUser()) == int(game.TableManager.Room.SeatCount) {
		game.GameStartFunc()
	}
}

// UserRecover 斷線返回
func (game *DeathGame) UserRecover(user *gamelogic.GameUser) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("UserRecover 發生錯誤, %v, %s", rec, debug.Stack())
		}
	}()

	game.LogicBase.UserRecover(user)

	if !game.IsGaming() {
		game.UserLeave(user)
		return
	}

	if game.ctx != nil {
		p := game.ctx.Player(user.UserID)
		p.State = StateLocked
		m := game.mRank(game.ctx)
		game.SendMessgeFunc(user.UserID, m)
	}
}

// UserLeave 玩家離開
func (game *DeathGame) UserLeave(user *gamelogic.GameUser) bool {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("UserLeave 發生錯誤, %v, %s", rec, debug.Stack())
		}
	}()

	if game.IsGaming() {
		user.UserState.OnUserOffLine(user, game.TableManager.Room)
		p := game.ctx.Player(user.UserID)
		if p != nil {
			if p.State == StateGaming {
				p.State = StateOffline
			}
			p.LastTime = time.Now()
		}
		game.GameLog.AddGameStep("玩家[%d] (L) 離開遊戲", user.UserID)
		game.checkFinish(game.ctx)
		return false
	}

	user.UserState.OnUserLeave(user, game.TableManager.Room)

	return true
}

// IsGaming 是否遊戲中
func (game *DeathGame) IsGaming() bool {
	return game.isgaming
}

// OnGameStart 遊戲開始
func (game *DeathGame) OnGameStart(f func()) {
	game.GameStartFunc = func() {
		if game.IsGaming() {
			return
		}

		// 遊戲開始前調用函數
		for _, f := range game.BeforeGameStartFunc {
			f()
		}

		game.isgaming = true
		game.GameStartTime = time.Now()
		game.PlayerActionManager.Clear()

		// 通知用戶狀態改變
		for _, usr := range game.AllTableUser() {
			usr.UserState.OnUserGame(usr, game.TableManager.Room)
		}

		// 廣播遊戲開始
		game.BoardcastGameStart()
		mlog.Debug("遊戲開始...")

		// 牌局記錄
		game.GameLog.GameBegin()

		f()

		// 遊戲開始後調用函數
		for _, f := range game.AfterGameStartFunc {
			f()
		}
	}
}

// OnGameEnd 遊戲結束
func (game *DeathGame) OnGameEnd(f func(int32)) {
	game.GameEndFunc = func() {
		// 遊戲結束前調用函數
		for _, f := range game.BeforeGameEndFunc {
			f()
		}

		game.isgaming = false
		game.LastGameEndTime = time.Now()
		game.PlayerActionManager.Clear()

		// 廣播遊戲結束
		game.BoardcastGameEnd()

		for _, usr := range game.AllGameUser() {
			usr.UserState.OnGameEnd(usr, game.TableManager.Room)
		}

		gametime := int32(time.Now().Sub(game.GameStartTime).Seconds())

		f(gametime)

		// 遊戲結束後調用函數
		for _, f := range game.AfterGameEndFunc {
			f()
		}

		mlog.Debug("遊戲結束...")

		// 遊戲結束剔除所有玩家，重新排隊
		for _, usr := range game.AllTableUser() {
			game.UserLeave(usr)
		}
	}
}

// 遊戲設定 封包
func (game *DeathGame) mGameSetting(start time.Time, ctx *RoundContext) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_GameSettingID)
	data := new(netproto.INFC_GameSetting)
	data.TableNo = game.TableManager.TableNo
	data.Bet = int64(game.GetServerConfig().BaseScore)
	data.StartTime = start.Unix()
	for _, u := range game.AllGameUser() {
		p := ctx.Player(u.UserID)
		data.Players = append(data.Players, game.mDataPlayer(p, ctx.Rank))
	}
	m.MsgData = data
	return m
}

// Proto設定 封包
func (game *DeathGame) mProtoSetting() *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_ProtoSettingID)
	data := new(netproto.INFC_ProtoSetting)
	for _, row := range game.Proto.Patterns {
		data.Patterns = append(data.Patterns, TransforPattern(row))
	}
	for _, row := range game.Proto.FruitAttributes {
		data.FruitAttributes = append(data.FruitAttributes, TransforFruitAttribute(row))
	}
	for _, row := range game.Proto.FruitEmitters {
		data.FruitEmitters = append(data.FruitEmitters, TransforFruitEmitter(row))
	}
	for _, row := range game.Proto.StageGroups {
		data.StageGroups = append(data.StageGroups, TransforStageGroup(row))
	}
	m.MsgData = data
	return m
}

// 扣你錢 封包
func (game *DeathGame) mBet(bet int64) {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_BetStepID)
	// m.MsgData = fmt.Sprintf("玩家[%d] 支付 %d 籌碼", user.UserID, bet)
	game.BoardcastMessage(m)
}

// Game Start 封包
func (game *DeathGame) mStart() {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_GameStartID)
	// m.MsgData = fmt.Sprintf("開幹")
	game.BoardcastMessage(m)
}

// 玩家分數排行 封包
func (game *DeathGame) mRank(ctx *RoundContext) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_GameRankID)
	m.MsgData = game.mDataRank(ctx)
	return m
}

// 玩家分數排行 封包
func (game *DeathGame) mDataRank(ctx *RoundContext) *netproto.INFC_GameRank {
	data := new(netproto.INFC_GameRank)
	for _, u := range game.AllGameUser() {
		p := ctx.Player(u.UserID)
		data.Players = append(data.Players, game.mDataPlayer(p, ctx.Rank))
	}
	return data
}

// 比賽結果 封包
func (game *DeathGame) mResult(ctx *RoundContext) {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_GameResultID)

	data := new(netproto.INFC_GameResult)

	win := ctx.Pool
	tax := int64(float64(win) * game.TaxRate)
	pay := win - tax

	data.WinMoney = pay
	data.Winner = game.mDataPlayer(ctx.Rank.Winner()[0], ctx.Rank)
	data.Rank = game.mDataRank(ctx)

	m.MsgData = data
	game.BoardcastMessage(m)
}

func (game *DeathGame) mDataPlayer(p *Player, r Rank) *netproto.INFC_Player {
	player := new(netproto.INFC_Player)
	switch p.State {
	case StateGaming:
		player.State = netproto.INFC_PlayerState_INFC_State_Gaming
	case StateDeath, StateGiveup, StateLocked:
		player.State = netproto.INFC_PlayerState_INFC_State_Dead
	case StateOffline:
		player.State = netproto.INFC_PlayerState_INFC_State_Offline
	}
	player.UserID = p.user.UserID
	player.UserName = p.user.NickName
	player.UserScore = p.user.Score
	player.Score = p.Score
	player.Miss = int32(p.Miss)
	player.Rank = int32(r.Rank(p.UID))

	return player
}

// // 通知死亡 封包
// func (game *DeathGame) mDeath(uid int32, way string) {
// 	m := new(msg.Message)
// 	m.BClassID = int32(netproto.MessageBClassID_Game)
// 	m.SClassID = int32(GameMessageID_Death)

// 	m.MsgData = fmt.Sprintf("玩家[%d] 掛了, 死法: %s", uid, way)

// 	game.BoardcastMessage(m)
// }

// // 通知封鎖 封包
// func (game *DeathGame) mLock(user *gamelogic.GameUser) {
// 	m := new(msg.Message)
// 	m.BClassID = int32(netproto.MessageBClassID_Game)
// 	m.SClassID = int32(GameMessageID_Lock)

// 	m.MsgData = "幹你就被鎖了來幹嘛"

// 	game.SendMessgeFunc(user.UserID, m)
// }
