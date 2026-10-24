package dfdc

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// ActionSpin spin
func (logic *Logic) ActionSpin(user *gamelogic.GameUser, gmsg *msg.Message) {

	bet := &netproto.DFDC_StartBet{}
	proto.Unmarshal(gmsg.MsgData.([]byte), bet)

	data := &netproto.DFDC_StartBet{}
	if err := proto.Unmarshal(gmsg.MsgData.([]byte), data); err != nil {
		mlog.Error("接收玩家选择封包解析错误. %v", err)
		return
	}

	ok, ret, money := logic.Bet(user, data.GetMoney(), int64(data.GetGoldenSymbol()))

	// 下注狀態
	logic.SendBetRetMessage(user, ret, money)

	if !ok {
		mlog.Debug("下注失敗. %v", ret)
		return
	}

	logic.UpdateState(slot.GameStateAward)

	// 依照選擇下注倍率選擇 SymbolGroupID
	SymbolGroupID, _, _, _ := logic.PrepareSpin()

	logic.SendJackpotPoolMessage(user)

	logic.Spin(SymbolGroupID)

	// 紀錄
	if !logic.IsFreeSpin() {
		logic.GameResultManager.ReelMode()
		logic.AddSlotGameHistoryLog(user, slot.HistoryTypeMAIN, slot.LogTypeREEL, logic.GameResultManager.ToSymbolString())
	} else {
		logic.GameResultManager.ReelMode()
		logic.AddSlotGameHistoryLog(user, slot.HistoryTypeFREE, slot.LogTypeREEL, logic.GameResultManager.ToSymbolString())
	}

	logic.UpdateUserScore(user)

	// 計算流水
	logic.CalcFlowAmount(user)

	// 傳送開獎結果
	logic.SendGameResultMessage(user)

	// 調整狀態為等待
	logic.UpdateState(slot.GameStateWait)

	if logic.IsFreeSpin() {
		logic.UpdateFreeSpin(user)

		logic.GameLog.AddGameStep("玩家免費遊戲累積贏得 %d", logic.FreeSpin.WinMoney)
	}

	logic.GameEndLog()

	logic.GameEndFunc()
}
