package wlhc

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// ActionSpin 玩家下注動作
func (logic *Logic) ActionSpin(user *gamelogic.GameUser, gmsg *msg.Message) {

	data := &netproto.Slot_MSG_Spin{}
	if err := proto.Unmarshal(gmsg.MsgData.([]byte), data); err != nil {
		mlog.Error("[玩家Spin] 接收玩家封包解析錯誤. %v", err)
		return
	}

	// 檢查傳入的MainGameID
	if logic.FreeSpin == nil {
		var exist bool
		for _, row := range logic.Proto.MainGameTable {
			if int32(row.GetId()) == data.GetID() {
				exist = true
				break
			}
		}
		if !exist {
			logic.SendBetStatusMessageToUser(user, slot.BetStateBetTypeError, 0)
			return
		}
	}

	ok, ret, money := logic.Bet(user, data.GetMoney(), int64(logic.MaxBetRate))

	// 下注狀態
	logic.SendBetStatusMessageToUser(user, ret, money)

	if !ok {
		return
	}

	logic.GameLog.GameBegin()

	logic.UpdateState(slot.GameStateAward)

	symbolGroupID, _ := logic.PrepareSpin(data.GetID())

	logic.Spin(symbolGroupID)

	logic.GameResultManager.RowMode()

	// 紀錄
	logic.GameResultManager.ReelMode()
	if logic.IsFreeSpin() {
		logic.AddSlotGameHistoryLog(user, slot.HistoryTypeFREE, slot.LogTypeREEL, logic.GameResultManager.ToSymbolString())
	} else {
		logic.AddSlotGameHistoryLog(user, slot.HistoryTypeMAIN, slot.LogTypeREEL, logic.GameResultManager.ToSymbolString())
	}
	logic.GameResultManager.RowMode()

	// 計算流水
	logic.CalcFlowAmount(user)

	// 寫分
	logic.UpdateUserScore(user)

	// 傳送開獎結果
	logic.SendGameResultMessage(user)

	if logic.FreeSpin != nil {
		logic.UpdateFreeSpin(user)

		logic.GameLog.AddGameStep("玩家免費遊戲累積贏得 %d", logic.FreeSpin.WinMoney)
	}

	logic.GameEndLog()

	// 調整狀態為等待
	logic.UpdateState(slot.GameStateWait)

	logic.GameEndFunc()
}
