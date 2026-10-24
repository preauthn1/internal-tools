package xmbx

import (
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

	logic.RegisterBeforeGameStartFunc(func() {})
	logic.RegisterAfterGameStartFunc(func() {
		logic.Clean()
		logic.CheckGameState()
	})
	logic.RegisterBeforeGameEndFunc(func() {})
	logic.RegisterAfterGameEndFunc(func() {
		logic.CheckFreeGame()
	})

	return logic
}

// GameMessageFunc on recived client message
func (logic *Logic) GameMessageFunc(user *gamelogic.GameUser, gmsg *msg.Message) {

	// 玩家Spin
	if gmsg.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID) {
		logic.ActionSpin(user, gmsg)
		return
	}

	// 玩家演完Lightning
	if gmsg.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_ActionLightningSpinID) || gmsg.SClassID == int32(11001) {
		logic.ActionLightningSpin(user, gmsg)
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

	logic.FreeSpin = nil
	logic.LightningGame = nil

	return logic.SlotLogicBase.UserLeave(user)
}

// CheckGameState 檢查遊戲狀態
//
//   * 一般正常狀況: GameStateBet
//   * 有FreeSpin: GameStateBet
//   * 有Lightning: GameStateUserAction
func (logic *Logic) CheckGameState() {

	if logic.LightningGame != nil {
		logic.SendLightningPrepareMessage(logic.GameUser)
		logic.UpdateState(slot.GameStateUserAction)
		return
	}

	if logic.FreeSpin != nil {
		logic.SendLeftFreeGameToUser(logic.GameUser)
		logic.UpdateState(slot.GameStateBet)
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
		logic.WG.Add(1)
		logic.ReadFreeSpinTimes(user)

		logic.WG.Add(1)
		logic.ReadJackpotPoolNoLink()

		logic.WG.Add(1)
		logic.ReadLightningData(user)

		logic.WG.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:

		logic.CheckFeatureGames()
		logic.CheckModeProto()
		logic.CheckModeExpire()

		logic.SendSettingMessage(user)
		logic.SendJackpotMessageToUser(user)

		logic.CheckFreeGame()

		logic.CheckGameState()
	case <-time.After(3 * time.Second):
		mlog.Error("讀取資料等待過久, 將玩家踢出遊戲")
		logic.KickUserFunc(user.UserID, "系統忙碌中")
	}
}
