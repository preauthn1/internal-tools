package wlzb

import (
	"sync"
	"time"

	"example.invalid/golang/protobuf/proto"
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

	if len(logic.BetRates) > 0 {
		logic.GameBetRate = int64(logic.BetRates[0])
	}

	logic.RegisterBeforeGameStartFunc(func() {
		logic.Clean()
	})
	logic.RegisterAfterGameStartFunc(func() {
		// logic.UpdateState(slot.GameStateBet)
		logic.CheckGameState()
	})
	logic.RegisterBeforeGameEndFunc(func() {})
	logic.RegisterAfterGameEndFunc(func() {
		logic.CheckFreeGame()
	})

	logic.RegisterAfterSwitchProtoFunc(func() {
		logic.SendSettingMessage(logic.GameUser)
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

	// 玩家選擇遊戲
	if gmsg.SClassID == int32(netproto.WLZB_GameMessageClassID_WLZB_ChooseGameID) {
		data := &netproto.WLZB_MSG_ChooseGame{}
		if err := proto.Unmarshal(gmsg.MsgData.([]byte), data); err != nil {
			mlog.Error("[玩家選擇遊戲] 接收玩家封包解析錯誤. %v", err)
			return
		}

		// 製作Option給底層處理
		_data := &netproto.Slot_MSG_ChooseGame{
			Option: &netproto.Slot_ChooseGameOption{
				ID:   proto.Int32(data.GetID()),
				Type: proto.String(string(slot.GameModeFreeSpin)),
			},
		}

		m := new(msg.Message)

		b, err := proto.Marshal(_data)
		if err != nil {
			mlog.Error("WLZB_ChooseGame Marshal Slot_MSG_ChooseGame Error, %v", err)
			return
		}

		m.MsgData = b

		logic.ActionChooseGame(user, m)
		logic.SendSelectChanceMessage(user)

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

	logic.UpdateFreeSpin(user)
	logic.RecordSelectGameTimes(user)

	return logic.SlotLogicBase.UserLeave(user)
}

// CheckGameState 檢查遊戲狀態
//
//   * 一般正常狀況: GameStateBet
//   * 有FreeSpin: GameStateBet
//   * 沒有FreeSpin, 有SelectChance: GameStateUserAction
func (logic *Logic) CheckGameState() {

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
		logic.WG.Add(1)
		logic.ReadFreeSpinTimes(user)

		logic.WG.Add(1)
		logic.ReadSelectGameTimes(user)

		logic.WG.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
		logic.CheckFeatureGames()
		logic.CheckModeProto()
		logic.CheckModeExpire()

		logic.SendSettingMessage(user)

		logic.CheckFreeGame()
		logic.CheckGameState()
	case <-time.After(3 * time.Second):
		mlog.Error("讀取資料等待過久, 將玩家踢出遊戲")
		logic.KickUserFunc(user.UserID, "系統忙碌中")
	}
}
