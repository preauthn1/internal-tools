package tgw

import (
	"runtime/debug"
	"strconv"
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

// Logic TGW
type Logic struct {
	*baselogic.SlotLogicBase

	PickGameSize int32
}

// NewLogic create tgw logic
func NewLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *Logic {
	logic := new(Logic)
	logic.SlotLogicBase = baselogic.NewSlotLogicBase(t, tfManager, cfgs)
	logic.SlotLogicBase.GameMessageFunc = logic.GameMessageFunc
	logic.SlotLogicBase.OnRecvRPCMessageFunc = logic.OnRecvRPCMessageFunc

	logic.loadConfig(cfgs)

	logic.RegisterBeforeGameStartFunc(func() {})
	logic.RegisterAfterGameStartFunc(func() {
		logic.Clean()
		logic.CheckGameState()
	})
	logic.RegisterBeforeGameEndFunc(func() {})
	logic.RegisterAfterGameEndFunc(func() {})

	return logic
}

func (logic *Logic) loadConfig(cfgs *config.Group) {

	if cfgs == nil {
		return
	}

	slotconfig := logic.Configs.GetConfig("slotconfig")

	// bonus
	{
		cfg := slotconfig.GetConfigSection("bonus")
		// use default 12
		logic.PickGameSize = 12

		if val, ok := cfg["size"]; ok {
			f, err := strconv.ParseFloat(val, 64)
			if err == nil {
				logic.PickGameSize = int32(f)
			}
		}
	}
}

// GameMessageFunc on recived client message
func (logic *Logic) GameMessageFunc(user *gamelogic.GameUser, gmsg *msg.Message) {

	// Client請求讀取Jackpot
	if gmsg.SClassID == int32(netproto.TGW_GameMessageClassID_TGW_ReadJackpotPoolID) {
		logic.SendJackpotPoolMessage(user)
	}

	// 玩家Spin
	if gmsg.SClassID == int32(netproto.TGW_GameMessageClassID_TGW_ActionSpinID) {
		logic.ActionSpin(user, gmsg)
		return
	}

	// 玩家選擇PickGame
	if gmsg.SClassID == int32(netproto.TGW_GameMessageClassID_TGW_ActionPickID) {

		if logic.PickGame == nil {
			return
		}

		data := &netproto.TGW_MSG_Pick{}
		if err := proto.Unmarshal(gmsg.MsgData.([]byte), data); err != nil {
			mlog.Error("PickGame Unmarshal TGW_MSG_Pick Error, %v", err)
			return
		}

		detail := &netproto.Slot_MSG_PickBonus{
			Record: &netproto.Slot_BonusGameRecord{
				ID:          proto.Int32(data.GetRecord().GetID()),
				Location:    proto.Int32(data.GetRecord().GetLocation()),
				RewareValue: proto.Int64(int64(data.GetRecord().GetRewareValue())),
			},
		}

		m := new(msg.Message)

		b, err := proto.Marshal(detail)
		if err != nil {
			mlog.Error("TGW PickGame Marshal Slot_MSG_PickBonus Error, %v", err)
			return
		}

		m.MsgData = b

		logic.ActionPickGame(user, m)
		logic.SendBonusGameRetMessage(user, data.GetRecord())
		return
	}

	// 玩家查詢記錄
	if gmsg.SClassID == int32(netproto.TGW_GameMessageClassID_TGW_LoadHistoryID) {

		data := &netproto.Slot_LoadHistoryData{}
		if err := proto.Unmarshal(gmsg.MsgData.([]byte), data); err != nil {
			mlog.Error("接收玩家选择封包解析错误. %v", err)
			return
		}

		logic.ReadSlotGameHistoryLog(user, data.GetPageIndex(), data.GetPageSize(), data.GetStartTime(), data.GetEndTime())
	}

}

// PiSaiMessage on recived client message
func (logic *Logic) PiSaiMessage(user *gamelogic.GameUser, gmsg *msg.Message) {}

// UserEnter 玩家進入
func (logic *Logic) UserEnter(user *gamelogic.GameUser, tableno, seatno int32) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("UserEnter 發生錯誤, %v\n%s", rec, string(debug.Stack()))
		}
	}()

	logic.SlotLogicBase.UserEnter(user, tableno, seatno)
	logic.GameUser = user

	logic.UpdateState(slot.GameStateWait)

	go logic.WaitAndSendData(user)
}

// UserRecover 玩家斷線重連
func (logic *Logic) UserRecover(user *gamelogic.GameUser) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("UserRecover 發生錯誤, %v\n%s", rec, string(debug.Stack()))
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
			mlog.Error("UserLeave 發生錯誤, %v\n%s", rec, string(debug.Stack()))
		}
	}()

	logic.RecordJackpotGame(user)

	logic.JackpotGame = nil

	return logic.SlotLogicBase.UserLeave(user)
}

// UpdateState 更新遊戲狀態
func (logic *Logic) UpdateState(state slot.GameState) {
	logic.GameState = state
	logic.SendStateMessage(logic.GameUser)
}

// CheckGameState 檢查遊戲狀態
func (logic *Logic) CheckGameState() {

	if logic.GameUser == nil {
		return
	}

	if logic.PickGame != nil {
		logic.SendBonusGameMessage(logic.GameUser)
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
		logic.ReadPickGame(user)

		logic.WG.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:

		logic.CheckFeatureGames()
		logic.CheckModeProto()
		logic.CheckModeExpire()

		logic.SendSettingMessage(user)
		logic.SendJackpotPoolMessage(logic.GameUser)
		logic.CheckGameState()
	case <-time.After(3 * time.Second):
		mlog.Error("讀取資料等待過久, 將玩家踢出遊戲")
		logic.KickUserFunc(user.UserID, "系統忙碌中")
	}
}
