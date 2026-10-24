package infc

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

var mu sync.Mutex

// FixedTimeGame game logic
type FixedTimeGame struct {
	*gamelogic.FixedGameLogicBase

	roundTime time.Duration
	taxrate   float64
	autoReady bool
	missLimit int64

	ctx *RoundContext
}

// NewFixedTimeGame create game logic
func NewFixedTimeGame(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *FixedTimeGame {
	game := new(FixedTimeGame)
	game.FixedGameLogicBase = gamelogic.NewFixedGameLogicBase(t, tfManager, cfgs)
	game.initConfig(cfgs)

	game.RegisterBeforeGameStartFunc(func() {})
	game.RegisterAfterGameStartFunc(func() {
		game.ctx = NewRoundContext(game.AllGameUser())
		game.Start(game.ctx)
	})
	game.RegisterBeforeGameEndFunc(func() {})
	game.RegisterAfterGameEndFunc(func() {
		// 檢查是否自動準備
		if game.autoReady {
			for _, u := range game.GetTable().GetTableUser() {
				game.FixedGameLogicBase.UserReady(u)
			}
		} else {
			// game.mReadyCheck()
		}
	})

	return game
}

func (game *FixedTimeGame) initConfig(cfgs *config.Group) {
	game.taxrate = float64(game.GetServerConfig().TaxScale) / 100

	sec := cfgs.GetConfigSection("gameserver", "gameserver")

	if val, err := strconv.ParseInt(sec["misslimit"], 0, 64); err != nil {
		game.missLimit = 3
	} else {
		game.missLimit = val
	}

	if val, err := strconv.ParseBool(sec["autoready"]); err != nil {
		game.autoReady = false
	} else {
		game.autoReady = val
	}

	if val, err := time.ParseDuration(sec["roundtime"]); err != nil {
		game.roundTime = 30 * time.Second
	} else {
		game.roundTime = val
	}

	game.taxrate = float64(game.GetServerConfig().TaxScale) / 100

}

// GameMessage on recived client message
func (game *FixedTimeGame) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("GameMessage 發生錯誤, %v, %s", rec, debug.Stack())
		}
	}()

	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	// switch GameMessageID(gmsg.SClassID) {
	// case GameMessageID_Ready:
	// 	game.FixedGameLogicBase.UserReady(user)
	// case GameMessageID_Score:
	// 	data := gmsg.MsgData.(map[string]interface{})
	// 	switch data["type"].(string) {
	// 	case "add":
	// 		example.invalid("玩家[%d] 回報遊戲結果, 加 %d 分", user.UserID, data["score"])
	// 		game.ctx.AddScore(user.UserID, data["score"].(int64))
	// 	case "set":
	// 		example.invalid("玩家[%d] 回報遊戲結果, %d 分", user.UserID, data["score"])
	// 		game.ctx.SetScore(user.UserID, data["score"].(int64))
	// 	}
	// case GameMessageID_Miss:
	// 	data := gmsg.MsgData.(map[string]interface{})
	// 	switch data["type"].(string) {
	// 	case "add":
	// 		example.invalid("玩家[%d] 回報Miss結果, 加 %d 次", user.UserID, data["miss"])
	// 		game.ctx.AddMiss(user.UserID, data["miss"].(int64))
	// 	case "set":
	// 		example.invalid("玩家[%d] 回報Miss結果, %d 次", user.UserID, data["miss"])
	// 		game.ctx.SetMiss(user.UserID, data["miss"].(int64))
	// 	}
	// }
}

// PiSaiMessage on recived client message
func (game *FixedTimeGame) PiSaiMessage(user *gamelogic.GameUser, gmsg *msg.Message) {}

// Start 開始遊戲
func (game *FixedTimeGame) Start(ctx *RoundContext) {

	game.GameLog.GameBegin()

	for _, u := range game.AllGameUser() {
		p := ctx.Player(u.UserID)
		p.State = StateGaming
		game.Bet(ctx, u)
	}

	game.GameLog.AddGameStep("遊戲水池 %d", ctx.Pool)
	game.mStart()

	go func() { // 定時通知所有玩家遊戲訊息
		for {
			select {
			case <-time.After(500 * time.Millisecond):
				game.mRank(game.ctx)
			case <-ctx.ctx.Done():
				mlog.Debug("停止定期通知資訊")
				return
			}
		}
	}()

	go func() { // 計時遊戲
		select {
		case <-time.After(game.roundTime):
			game.Calc(ctx)
		}
	}()

	go func() { // 等待遊戲結束
		select {
		case <-ctx.ctx.Done():
			game.GameLog.GameEnd()
			game.GameEndFunc()
		}
	}()
}

// Bet 玩家下注
func (game *FixedTimeGame) Bet(ctx *RoundContext, user *gamelogic.GameUser) {

	bet := int64(game.GetServerConfig().BaseScore)

	game.GameLog.AddGameStep("玩家[%d] 支付遊戲底分 %d", user.UserID, bet)

	user.BetAmount = bet
	ctx.Pool += bet

	game.UpdateUserScoreFunc(user.UserID, -bet, 0)
	game.mBet(user, bet)
}

// Calc 結算
func (game *FixedTimeGame) Calc(ctx *RoundContext) {

	// winner := ctx.Winner()

	// win := ctx.Pool
	// tax := int64(float64(win) * game.taxrate)
	// pay := win - tax

	// winner.user.PayAmount = pay

	// for _, u := range game.AllGameUser() {
	// 	var winloss string

	// 	if u.UserID == winner.UID {
	// 		ctx.UpdateAmount[u.UserID] += pay
	// 		ctx.PayAmount[u.UserID] = pay
	// 		winloss = "勝利"
	// 		game.UpdateUserScoreFunc(u.UserID, pay, tax)
	// 	} else {
	// 		winloss = "戰敗"
	// 	}

	// 	if ctx.UpdateAmount[u.UserID] < 0 {
	// 		ctx.ValidAmount[u.UserID] = ctx.BetAmount[u.UserID]
	// 	}

	// 	game.GameLog.AddGameStep("玩家[%d] 獲得 %8d 分 %s, 注額: %5d, 獎金: %5d, 稅金: %5d, 實收: %5d, 流水: %5d", u.UserID, ctx.Player(u.UserID).Score, winloss, ctx.BetAmount[u.UserID], win, tax, ctx.PayAmount[u.UserID], ctx.ValidAmount[u.UserID])
	// }

	ctx.cancel()
}

// UserEnter 玩家進入
func (game *FixedTimeGame) UserEnter(user *gamelogic.GameUser, tableno, seatno int32) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("UserEnter 發生錯誤, %v, %s", rec, debug.Stack())
		}
	}()

	game.FixedGameLogicBase.UserEnter(user, tableno, seatno)
	game.FixedGameLogicBase.UserReady(user)
}

// UserRecover 玩家斷線重連
func (game *FixedTimeGame) UserRecover(user *gamelogic.GameUser) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("UserRecover 發生錯誤, %v, %s", rec, debug.Stack())
		}
	}()

	game.FixedGameLogicBase.UserRecover(user)

	if game.ctx == nil {
		example.invalid("UserRecover but round is finish")
		return
	}

	p := game.ctx.Player(user.UserID)
	if p == nil {
		example.invalid("UserRecover but round can not find user %d", user.UserID)
		return
	}

	example.invalid("UserRecover Lock user %d", user.UserID)
	// p.State = StateLocked
	// game.mLock(user)
}

// UserLeave 玩家離開
func (game *FixedTimeGame) UserLeave(user *gamelogic.GameUser) bool {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("UserLeave 發生錯誤, %v, %s", rec, debug.Stack())
		}
	}()

	return game.FixedGameLogicBase.UserLeave(user)
}

// 扣你錢 封包
func (game *FixedTimeGame) mBet(user *gamelogic.GameUser, bet int64) {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_BetStepID)
	// m.MsgData = fmt.Sprintf("玩家[%d] 支付 %d 籌碼", user.UserID, bet)
	game.BoardcastMessage(m)
}

// Game Start 封包
func (game *FixedTimeGame) mStart() {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_GameStartID)
	m.MsgData = fmt.Sprintf("開幹")
	game.BoardcastMessage(m)
}

// 發送請求準備 封包
// func (game *FixedTimeGame) mReadyCheck() {
// 	m := new(msg.Message)
// 	m.BClassID = int32(netproto.MessageBClassID_Game)
// 	m.SClassID = int32(GameMessageID_ReadyCheck)
// 	m.MsgData = fmt.Sprintf("準備嗎？")
// 	game.BoardcastMessage(m)
// }

// 玩家分數排行 封包
func (game *FixedTimeGame) mRank(ctx *RoundContext) {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_GameRankID)

	data := make(map[string]interface{})

	for _, u := range game.AllGameUser() {
		p := ctx.Player(u.UserID)
		if u.IsOffLine() {
			p.State = StateOffline
		}
		data[fmt.Sprintf("%d", u.UserID)] = map[string]interface{}{
			"score": p.Score,
			"miss":  p.Miss,
			"state": p.State,
		}
	}

	m.MsgData = data
	game.BoardcastMessage(m)
}

// 通知封鎖 封包
// func (game *FixedTimeGame) mLock(user *gamelogic.GameUser) {
// 	m := new(msg.Message)
// 	m.BClassID = int32(netproto.MessageBClassID_Game)
// 	m.SClassID = int32(GameMessageID_Lock)

// 	m.MsgData = "幹你就被鎖了來幹嘛"

// 	game.SendMessgeFunc(user.UserID, m)
// }
