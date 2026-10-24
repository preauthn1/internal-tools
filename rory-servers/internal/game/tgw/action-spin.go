package tgw

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
	bet := &netproto.TGW_MSG_Spin{}
	proto.Unmarshal(gmsg.MsgData.([]byte), bet)

	data := &netproto.TGW_MSG_Spin{}
	if err := proto.Unmarshal(gmsg.MsgData.([]byte), data); err != nil {
		mlog.Error("接收玩家选择封包解析错误. %v", err)
		return
	}

	ok, ret, money := logic.Bet(user, data.GetMoney(), int64(logic.MaxBetRate))

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

	logic.UpdateUserScore(user)

	// 計算流水
	logic.CalcFlowAmount(user)

	// 傳送開獎結果
	logic.SendGameResultMessage(user)

	// 調整狀態為等待
	logic.UpdateState(slot.GameStateWait)

	logic.GameEndLog()

	logic.GameEndFunc()
}
