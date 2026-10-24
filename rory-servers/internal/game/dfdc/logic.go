package dfdc

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

// Logic DFDC
type Logic struct {
	*baselogic.SlotLogicBase
}

// NewLogic create dfdc logic
func NewLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *Logic {
	logic := new(Logic)
	logic.SlotLogicBase = baselogic.NewSlotLogicBase(t, tfManager, cfgs)
	logic.SlotLogicBase.GameMessageFunc = logic.GameMessageFunc
	logic.SlotLogicBase.OnRecvRPCMessageFunc = logic.OnRecvRPCMessageFunc
	logic.CheckFreeGameFunc = logic.CheckFreeGame

	logic.RegisterBeforeGameStartFunc(func() {})
	logic.RegisterAfterGameStartFunc(func() {
		logic.Clean()
		logic.CheckGameState()
	})
	logic.RegisterBeforeGameEndFunc(func() {})
	logic.RegisterAfterGameEndFunc(func() {})

	return logic
}

// GameMessageFunc on recived client message
func (logic *Logic) GameMessageFunc(user *gamelogic.GameUser, gmsg *msg.Message) {
	// 遊戲開始
	if gmsg.SClassID == int32(netproto.DFDC_GameMessageClassID_DFDCGameStartID) {
		logic.ActionSpin(user, gmsg)
		return
	}

	// 玩家Bonus Game選擇
	if gmsg.SClassID == int32(netproto.DFDC_GameMessageClassID_DFDCBonusGameIconID) {

		detail := &netproto.DFDC_BonusGameDetail{}
		if err := proto.Unmarshal(gmsg.MsgData.([]byte), detail); err != nil {
			mlog.Error("DFDCBonusGame Unmarshal DFDC_BonusGameDetail Error, %v", err)
			return
		}

		bonus := &netproto.Slot_MSG_PickBonus{
			Record: &netproto.Slot_BonusGameRecord{
				ID:          proto.Int32(detail.GetId()),
				Location:    proto.Int32(detail.GetPlace()),
				RewareValue: proto.Int64(int64(detail.GetIcon())),
			},
		}

		m := new(msg.Message)

		b, err := proto.Marshal(bonus)
		if err != nil {
			mlog.Error("DFDCBonusGame Marshal Slot_MSG_PickBonus Error, %v", err)
			return
		}

		m.MsgData = b

		if logic.ActionPickJackpot(user, m) {
			logic.SendJackpotPoolMessage(user)
		}
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

	logic.UpdateState(slot.GameStateWait)

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

	logic.UpdateState(slot.GameStateWait)

	go logic.WaitAndSendData(user)
}

// UserLeave 玩家離開
func (logic *Logic) UserLeave(user *gamelogic.GameUser) bool {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("UserLeave 發生錯誤, %v", rec)
		}
	}()

	logic.RecordJackpotGame(user)
	logic.UpdateFreeSpin(user)

	logic.JackpotGame = nil
	logic.FreeSpin = nil

	return logic.SlotLogicBase.UserLeave(user)
}

// UpdateState 更新遊戲狀態
func (logic *Logic) UpdateState(state slot.GameState) {
	logic.GameState = state
	logic.SendStateMessage(logic.GameUser)
}

// CheckFreeGame 檢查 FreeSpin, 檢查免費遊戲是否結束沒次數
func (logic *Logic) CheckFreeGame() {

	if logic.FreeSpin == nil {
		return
	}

	if logic.FreeSpin.Count > 0 {
		return
	}

	logic.SendFreeSpinCountMessage(logic.GameUser)
	logic.FreeSpin = nil
}

// CheckGameState 檢查遊戲狀態
func (logic *Logic) CheckGameState() {

	if logic.JackpotGame != nil {
		logic.SendJackpotGameDetailMessage(logic.GameUser)
		logic.UpdateState(slot.GameStateBet)
		return
	}

	if logic.FreeSpin != nil {
		logic.SendFreeSpinCountMessage(logic.GameUser)
		logic.SendFreeSpinBetMessage(logic.GameUser)

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
		logic.ReadJackpotPoolNoLink()

		logic.WG.Add(1)
		logic.ReadFreeSpinTimes(user)

		logic.WG.Add(1)
		logic.ReadJackpotGame(user)

		logic.WG.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:

		logic.CheckFeatureGames()
		logic.CheckModeProto()
		logic.CheckModeExpire()

		logic.SendGoldenSymbolMessage(user)
		logic.SendSettingMessage(user)

		logic.SendJackpotPoolMessage(user)

		logic.CheckGameState()
	case <-time.After(3 * time.Second):
		mlog.Error("讀取資料等待過久, 將玩家踢出遊戲")
		logic.KickUserFunc(user.UserID, "系統忙碌中")
	}
}
