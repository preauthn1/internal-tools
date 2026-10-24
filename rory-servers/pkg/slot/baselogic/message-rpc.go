package baselogic

import (
	"reflect"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// OnRecvRPCMessage 接收DB Server回傳的Message訊息
func (logic *SlotLogicBase) OnRecvRPCMessage(m *msg.Message) {

	if m == nil {
		return
	}

	rm, ok := m.MsgData.(*rpcmsg.RPCMessage)
	if !ok {
		mlog.Error("接收RPC資訊錯誤, 訊息不為*rpcmsg.RPCMessage")
		return
	}

	mmbyte := rm.MsgData.([]byte)
	if len(mmbyte) == 0 {
		mlog.Warn("接收RPC資訊錯誤, MsgData為空")
		// return
	}

	var skip bool

	for _, id := range logic.SkipRPCMessageIDs {
		if id == m.SClassID {
			skip = true
			break
		}
	}

	if !skip {
		switch m.SClassID {
		case int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadRetID): // 從DB獲得的Jackpot資訊
			logic.RecvJackpotData(m)

		case int32(netproto.Slot_GameMessageClassID_Slot_GameLogReadRetID): // 從DB獲得的遊戲歷程紀錄
			logic.RecvSlotGameLog(m)

		case int32(netproto.Slot_GameMessageClassID_Slot_FreeGameReadRetID): // 從DB獲得的free spin資訊
			logic.RecvFreeSpinTimes(m)

		case int32(netproto.Slot_GameMessageClassID_Slot_SelectGameReadRetID): // 從DB獲得的SelectGame資訊
			logic.RecvSelectGameTimes(m)

		case int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadRetID):
			// 從DB獲得的BonusGame資訊
			// 從DB獲得的JackpotGame資訊
			// 由於不知道使用的遊戲會是BonusGame或是JackpotGame
			// 所以交給個遊戲自行處理
			// 	logic.recvJackpotGame(m)

		case int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadRetID): // 從DB獲得的DeferredAward資訊
			logic.RecvJackpotGameChance(m)

		case int32(netproto.Slot_GameMessageClassID_Slot_LightningReadRetID): // 從DB獲得的LightningGame資訊
			logic.recvLightningData(m)
		case int32(netproto.Slot_GameMessageClassID_Slot_UserAssignedModeReadRetID): // 從DB獲得的 UserAssignedMode 資訊
			logic.RecvUserAssignedMode(m)
		}
	}

	if logic.OnRecvRPCMessageFunc != nil {
		logic.OnRecvRPCMessageFunc(m)
	}
}

// CheckWaitGroup 若Logic有wait group存在，則跟wg說做好了
func (logic *SlotLogicBase) CheckWaitGroup() {
	if logic.WG == nil {
		return
	}
	logic.WG.Done()
}

// CheckFeatureGames 檢查殘留遊戲是否在buffer內，若是同個mode才能被解開
func (logic *SlotLogicBase) CheckFeatureGames() {

	mlog.Debug("Table %d, CheckFeatureGames", logic.TableNo())
	for mode, games := range logic.bufFeatureGames {
		mlog.Debug("  - %v", mode)
		for _, game := range games {
			mlog.Debug("    - %p, %v", game, reflect.TypeOf(game))
		}
	}

	tmplogic := new(SlotLogicBase)

	var injectGames = func(logic *SlotLogicBase, features []interface{}) *SlotLogicBase {
		for _, f := range features {
			if g, ok := f.(*slot.FreeSpin); ok {
				logic.FreeSpin = g
				continue
			}
			if g, ok := f.(*slot.JackpotGame); ok {
				logic.JackpotGame = g
				continue
			}
			if g, ok := f.(*slot.PickGame); ok {
				logic.PickGame = g
				continue
			}
			if g, ok := f.(*LightningGame); ok {
				logic.LightningGame = g
				continue
			}
			if g, ok := f.([]*slot.SelectChance); ok {
				logic.SelectChances = g
				continue
			}
			if g, ok := f.([]*slot.JackpotGameChance); ok {
				logic.JackpotGameChances = g
				continue
			}
		}
		return logic
	}

	// 檢查 DefaultPattern 是否有殘留資料
	tmplogic = injectGames(tmplogic, logic.bufFeatureGames[DefaultPattern])

	if tmplogic.ExistFeatureGame() {
		injectGames(logic, logic.bufFeatureGames[DefaultPattern])
		return
	}

	injectGames(logic, logic.bufFeatureGames[logic.Mode])
}
