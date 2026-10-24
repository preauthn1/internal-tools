package pirate

// import (
// 	"example.invalid/golang/protobuf/proto"

// 	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
// 	"example.invalid/ServerTeam/Servers/pkg/msg"
// 	"example.invalid/ServerTeam/Servers/pkg/netproto"
// )

// // SendSettingMessage 發送遊戲設置內容
// //
// //   * BClassID = 10001 (MessageBClassID_Game)
// //   * SClassID = 20001 (Slot_GameMessageClassID_Slot_SettingID)
// //   * Message = { BetCount: [10, 20, 30, 60, 100, 200, 300, 600, 1000], RowCounts: 3, TotalPayLine: 1, BigWinSetp: [300, 400, 500, 500000], PayLineMode: 1, ReelSetting: [ { ReelID: 1, SymbolGroup: [ { SymbolID: 1, SymbolName: "p1" }, ... ] }, ... ], SymbolPayTable: [ { Rate: 5, SymbolName: "p1", Count: 3 }, ... ] }
// func (logic *Logic) SendSettingMessage(user *gamelogic.GameUser) {
// 	pmsg := new(msg.Message)
// 	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
// 	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_SettingID)
// 	data := new(netproto.ALD_MSG_Setting)

// 	data.BetCount = logic.BetCounts
// 	data.RowCounts = proto.Int32(logic.DefaultRowNum)
// 	data.TotalPayLine = proto.Int32(int32(len(logic.SlotPayLine)))
// 	data.BigWinStep = logic.BigWinRule
// 	data.PayLineMode = proto.Int32(logic.PayLineMode.Int())
// 	data.MinCarryMoney = proto.Int64(logic.MinCarryMoney)
// 	data.PayLineDesc = proto.String(logic.PayLineDesc)

// 	// ReelSetting
// 	data.ReelSetting = logic.NetprotoSlotReelSetting

// 	// SymbolPayTable
// 	data.SymbolPayTable = logic.NetprotoSlotSymbolPayTable

// 	// PayLine
// 	data.PayLines = logic.NetprotoSlotPayLine

// 	pmsg.MsgData = data

// 	logic.SendMessageToUser(user, pmsg)
// }

// // SendGameResultMessage 發送遊戲結果
// //
// //   * BClassID = 10001 (MessageBClassID_Game)
// //   * SClassID = 20004 (Slot_GameMessageClassID_Slot_GameResultID)
// //   * Message = { BetMoney: 100, WinMoney: 0, FinalMoney: 300, ReelResult: [{},{},...] }
// func (logic *Logic) SendGameResultMessage(user *gamelogic.GameUser) {

// 	pmsg := new(msg.Message)
// 	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
// 	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_GameResultID)
// 	data := new(netproto.ALD_MSG_GameResult)

// 	data.BetMoney = proto.Int64(logic.GamePayMoney)
// 	data.WinMoney = proto.Int64(logic.GameWinMoney)
// 	data.FinalMoney = proto.Int64(logic.GameUser.Score)

// 	pmsg.MsgData = data

// 	logic.SendMessageToUser(user, pmsg)
// }
