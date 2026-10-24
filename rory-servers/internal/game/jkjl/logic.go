package jkjl

import (
	"runtime/debug"
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/baselogic"
)

// Logic game logic
type Logic struct {
	*baselogic.SlotLogicBase
}

// NewLogic create game logic
func NewLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *Logic {
	logic := new(Logic)
	logic.SlotLogicBase = baselogic.NewSlotLogicBase(t, tfManager, cfgs)
	logic.SlotLogicBase.GameMessageFunc = logic.GameMessageFunc
	logic.SlotLogicBase.OnRecvRPCMessageFunc = logic.OnRecvRPCMessageFunc

	logic.RegisterBeforeGameStartFunc(func() {
		logic.Clean()
	})
	logic.RegisterAfterGameStartFunc(func() {
		logic.CheckFreeGame()
		logic.CheckGameState()
	})
	logic.RegisterBeforeGameEndFunc(func() {})
	logic.RegisterAfterGameEndFunc(func() {})

	return logic
}

// GameMessageFunc on recived client message
func (logic *Logic) GameMessageFunc(user *gamelogic.GameUser, gmsg *msg.Message) {
	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("GameMessageFunc 發生錯誤, %v, %s", rec, debug.Stack())
		}
	}()

	// 玩家Spin
	if gmsg.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID) {
		logic.ActionSpin(user, gmsg)
		return
	}

	// 玩家選擇SelectGame
	if gmsg.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_ActionChooseGameID) {
		logic.ActionChooseGame(user, gmsg)
		return
	}

	// 玩家選擇BonusGame ICON
	if gmsg.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_PickBonusID) {
		logic.ActionPickJackpot(user, gmsg)
		return
	}
}

// PiSaiMessage on recived client message
func (logic *Logic) PiSaiMessage(user *gamelogic.GameUser, gmsg *msg.Message) {}

// UserEnter 玩家進入
func (logic *Logic) UserEnter(user *gamelogic.GameUser, tableno, seatno int32) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("UserEnter 發生錯誤, %v", rec)
		}
	}()

	logic.SlotLogicBase.UserEnter(user, tableno, seatno)

	go logic.WaitAndSendData(user)
}

// UserRecover 玩家斷線重連
func (logic *Logic) UserRecover(user *gamelogic.GameUser) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("UserRecover 發生錯誤, %v", rec)
		}
	}()

	logic.SlotLogicBase.UserRecover(user)

	go logic.WaitAndSendData(user)
}

// UserLeave 玩家離開
func (logic *Logic) UserLeave(user *gamelogic.GameUser) bool {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("UserLeave 發生錯誤, %v", rec)
		}
	}()

	// 把殘留資料存DB
	logic.UpdateFreeSpin(user)
	logic.RecordSelectGameTimes(user)
	logic.RecordJackpotGame(user)

	logic.FreeSpin = nil
	logic.SelectChances = make([]*slot.SelectChance, 0)
	logic.JackpotGameChances = make([]*slot.JackpotGameChance, 0)
	logic.JackpotGame = nil

	return logic.SlotLogicBase.UserLeave(user)
}

// CheckGameState 檢查遊戲狀態
//
//   * 一般正常狀況: GameStateBet
//   * 有FreeSpin: GameStateBet
//   * 有SelectGame: GameStateUserAction
//   * 有JackpotGame: GameStateUserAction
func (logic *Logic) CheckGameState() {

	if logic.JackpotGame != nil {

		logic.SendJackpotGameDetailMessage(logic.GameUser)

		logic.UpdateState(slot.GameStateUserAction)
		return
	}

	if logic.FreeSpin != nil {
		logic.SendLeftFreeGameToUser(logic.GameUser)

		logic.UpdateState(slot.GameStateBet)
		return
	}

	if logic.SelectChances != nil && len(logic.SelectChances) > 0 {
		logic.GameLog.AddGameStep("還有未選擇遊戲 %d 次, 改變遊戲狀態 %d", len(logic.SelectChances), slot.GameStateUserAction)

		logic.SendSelectChanceMessage(logic.GameUser)

		logic.UpdateState(slot.GameStateUserAction)
		return
	}

	logic.UpdateState(slot.GameStateBet)
}

// WaitAndSendData 等待資料並傳送給玩家
func (logic *Logic) WaitAndSendData(user *gamelogic.GameUser) {

	select {
	case <-logic.ModeWait:
		mlog.Debug("玩家Mode查詢完畢, 開始查詢遊戲資料")
	case <-time.After(3 * time.Second):
		mlog.Error("玩家Mode查詢等待過久, 查詢遊戲資料")
	}

	// 設定等待
	logic.WG = new(sync.WaitGroup)
	defer func() { logic.WG = nil }()

	waitCh := make(chan struct{})

	go func() {
		close(waitCh)
	}()

	select {
	case <-waitCh:

		logic.CheckFeatureGames()
		logic.CheckModeProto()
		logic.CheckModeExpire()

		logic.SendReelSettingMessage(user)
		logic.SendSettingMessage(user)

		logic.CheckGameState()
	case <-time.After(3 * time.Second):
		mlog.Error("讀取資料等待過久, 將玩家踢出遊戲")
		logic.KickUserFunc(user.UserID, "系統忙碌中")
	}
}
