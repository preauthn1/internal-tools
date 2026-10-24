package s888

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

	ok, ret, money := logic.Bet(user, data.GetMoney(), int64(logic.MaxBetRate))

	// 下注狀態
	logic.SendBetStatusMessageToUser(user, ret, money)

	if !ok {
		return
	}

	logic.GameLog.GameBegin()

	logic.UpdateState(slot.GameStateAward)

	symbolGroupID, _, _, _ := logic.PrepareSpin()

	logic.Spin(symbolGroupID)

	logic.GameResultManager.RowMode()

	// 計算流水
	logic.CalcFlowAmount(user)

	// 寫分
	logic.UpdateUserScore(user)

	// 傳送開獎結果
	logic.SendGameResultMessage(user)

	if logic.FreeSpin != nil {
		logic.UpdateFreeSpin(user)
		logic.SendLeftFreeGameToUser(user)

		logic.GameLog.AddGameStep("玩家免費遊戲累積贏得 %d", logic.FreeSpin.WinMoney)
	}

	logic.GameEndLog()

	// 調整狀態為等待
	logic.UpdateState(slot.GameStateWait)

	logic.GameEndFunc()
}
