package pirate

// import (
// 	"example.invalid/golang/protobuf/proto"

// 	"example.invalid/ServerTeam/Servers/pkg/config"
// 	"example.invalid/ServerTeam/Servers/pkg/event"
// 	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
// 	"example.invalid/ServerTeam/Servers/pkg/mlog"
// 	"example.invalid/ServerTeam/Servers/pkg/msg"
// 	"example.invalid/ServerTeam/Servers/pkg/netproto"
// 	"example.invalid/ServerTeam/Servers/pkg/slot"
// 	"example.invalid/ServerTeam/Servers/pkg/slot/baselogic"
// )

// // Logic Aladdin
// type Logic struct {
// 	*baselogic.SlotLogicBase

// 	FreeSpin *slot.FreeSpin
// }

// // NewLogic create Aladdin game logic
// func NewLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.ConfigGroup) *Logic {
// 	logic := new(Logic)
// 	logic.SlotLogicBase = baselogic.NewSlotLogicBase(t, tfManager, cfgs, protodatas)
// 	logic.SlotLogicBase.GameMessageFunc = logic.GameMessageFunc
// 	logic.SlotLogicBase.OnRecvRPCMessageFunc = logic.OnRecvRPCMessageFunc

// 	logic.RegisterBeforeGameStartFunc(func() {})
// 	logic.RegisterAfterGameStartFunc(func() {})
// 	logic.RegisterBeforeGameEndFunc(func() {})
// 	logic.RegisterAfterGameEndFunc(func() {})

// 	return logic
// }

// // UserEnter 玩家進入
// func (logic *Logic) UserEnter(user *gamelogic.GameUser, tableno, seatno int32) {
// 	logic.SlotLogicBase.UserEnter(user, tableno, seatno)

// 	logic.SendSettingMessage(user)

// 	logic.UpdateState(slot.GameStateBet)
// }

// // UserRecover 玩家斷線重連
// func (logic *Logic) UserRecover(user *gamelogic.GameUser) {
// 	logic.SlotLogicBase.UserRecover(user)

// 	logic.SendSettingMessage(user)

// 	logic.UpdateState(slot.GameStateBet)
// }

// // UserLeave 玩家離開
// func (logic *Logic) UserLeave(user *gamelogic.GameUser) bool {

// 	logic.CleanFreeSpin()
// 	logic.Clean()

// 	return logic.SlotLogicBase.UserLeave(user)
// }

// // GameMessageFunc on recived client message
// func (logic *Logic) GameMessageFunc(user *gamelogic.GameUser, gmsg *msg.Message) {

// 	// 玩家Spin
// 	if gmsg.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID) {

// 		data := &netproto.Slot_MSG_Spin{}
// 		if err := proto.Unmarshal(gmsg.MsgData.([]byte), data); err != nil {
// 			mlog.Error("接收玩家选择封包解析错误. %v", err)
// 			return
// 		}

// 		ok, ret, money := logic.Bet(user, data.GetMoney())

// 		// 下注狀態
// 		logic.SendBetStatusMessageToUser(user, ret, money)

// 		if !ok {
// 			mlog.Debug("下注失敗. %v", ret)
// 			return
// 		}

// 		logic.UpdateState(slot.GameStateAward)

// 		symbolGroupID := logic.DefaultSymbolGroupID

// 		if logic.FreeSpin != nil {
// 			logic.GameLog.AddGameStep("玩家 %d; 底注: %d, 總線數: %d, 花費金額: %d, 此次為免費遊戲，剩餘次數: %d", user.UserID, logic.FreeSpin.BetMoney, len(logic.SlotPayLine), logic.GamePayMoney, logic.FreeSpin.Count)
// 			mlog.Debug("玩家 %d; 底注: %d, 總線數: %d, 花費金額: %d, 此次為免費遊戲，剩餘次數: %d", user.UserID, logic.FreeSpin.BetMoney, len(logic.SlotPayLine), logic.GamePayMoney, logic.FreeSpin.Count)
// 			symbolGroupID = logic.FreeSpin.SymbolGroupID

// 		} else {
// 			logic.GameLog.AddGameStep("玩家 %d; 底注: %d, 總線數: %d, 花費金額: %d", user.UserID, logic.GameBetMoney, len(logic.SlotPayLine), logic.GamePayMoney)
// 			mlog.Debug("玩家 %d; 底注: %d, 總線數: %d, 花費金額: %d", user.UserID, logic.GameBetMoney, len(logic.SlotPayLine), logic.GamePayMoney)

// 			// 累積Jackpot獎池
// 			for _, jackpot := range logic.Jackpots {
// 				// 手動更新獎池資訊
// 				jackpot.Increase(logic.GamePayMoney)
// 				logic.GameLog.AddGameStep("水池增加 (N/10000) N=%d , 累積獎金增加 %d, 當前最大出獎金額 %d", jackpot.TotalContribution, jackpot.AwardTotal, jackpot.GetCurrentPool())
// 				mlog.Debug("水池增加 (N/10000) N=%d , 累積獎金增加 %d, 當前最大出獎金額 %d", jackpot.TotalContribution, jackpot.AwardTotal, jackpot.GetCurrentPool())

// 				// 更新Jackpot
// 				logic.UpdateJackpotPoolNoLink(0, jackpot.GetIncrease(logic.GamePayMoney), 0, jackpot.Level)
// 			}

// 		}

// 		// 把手動更新得獎池資訊給Client
// 		// logic.SendJackpotPoolMessage(user)

// 		// 開始轉
// 		logic.Spin(symbolGroupID)

// 		// 計算流水
// 		logic.CalcFlowAmount(user)

// 		// 傳送開獎結果
// 		logic.SendGameResultMessage(user)

// 		// 調整狀態為等待
// 		logic.UpdateState(slot.GameStateWait)

// 		logic.GameEndLog()

// 		logic.GameEndFunc()

// 		logic.UpdateState(slot.GameStateBet)

// 		logic.Clean()

// 		return
// 	}
// }

// // PiSaiMessage on recived client message
// func (logic *Logic) PiSaiMessage(user *gamelogic.GameUser, gmsg *msg.Message) {}

// // Clean 清除遊戲資料
// func (logic *Logic) Clean() {
// 	logic.SlotLogicBase.Clean()
// }

// // CleanFreeSpin 清除FreeSpin資料
// func (logic *Logic) CleanFreeSpin() {
// 	logic.FreeSpin = nil
// }
