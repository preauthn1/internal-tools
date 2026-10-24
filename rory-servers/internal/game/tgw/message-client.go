package tgw

import (
	"reflect"
	"sort"
	"strconv"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// SendSettingMessage 發送遊戲設置內容
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10001 (TGW_GameMessageClassID_TGW_SettingID)
//   * Message = { BetCount: [10, 20, 30, 60, 100, 200, 300, 600, 1000], RowCounts: 3, TotalPayLine: 1, BigWinSetp: [300, 400, 500, 500000], PayLineMode: 1, ReelSetting: [ { ReelID: 1, SymbolGroup: [ { SymbolID: 1, SymbolName: "p1" }, ... ] }, ... ], SymbolPayTable: [ { Rate: 5, SymbolName: "p1", Count: 3 }, ... ] }
func (logic *Logic) SendSettingMessage(user *gamelogic.GameUser) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TGW_GameMessageClassID_TGW_SettingID)
	data := new(netproto.TGW_MSG_Setting)

	data.BetCount = logic.BetCounts
	data.RowCounts = proto.Int32(logic.DefaultRowNum)
	data.TotalPayLine = proto.Int32(int32(len(logic.Proto.SlotPayLine)))
	data.BigWinStep = logic.BigWinRule
	data.PayLineMode = proto.Int32(logic.PayLineMode.Int())
	data.MinCarryMoney = proto.Int64(logic.MinCarryMoney)
	var payLineDesc string = ""
	if val, ok := logic.Descriptions["payline"]; ok {
		payLineDesc = val.Description
	}
	data.PayLineDesc = proto.String(payLineDesc)

	var pickGameDesc string = ""
	if val, ok := logic.Descriptions["pickgame"]; ok {
		pickGameDesc = val.Description
	}
	data.PickGameDesc = proto.String(pickGameDesc)

	// UnitPrefix
	if logic.UnitPrefix != nil {
		for name, multiple := range logic.UnitPrefix {
			unit := new(netproto.Slot_UnitPrefix)
			unit.Name = proto.String(name)
			unit.Multiple = proto.Int64(multiple)
			data.UnitPrefix = append(data.UnitPrefix, unit)
		}
	}

	// ReelSetting
	for reelNum := int32(1); reelNum <= logic.DefaultReelNum; reelNum++ {
		setting := new(netproto.TGW_ReelSetting)
		setting.RealID = proto.Int32(int32(reelNum))

		// 只留一組，並且打亂
		sbMap := make(map[string]*netproto.TGW_Symbol)

		for rowNum := int32(1); rowNum <= logic.DefaultRowNum; rowNum++ {
			symbols := logic.GameResultManager.GetSymbolsByReelRow(int32(logic.DefaultSymbolGroupID), reelNum, rowNum, logic.Proto.SlotSymbolDistribution)
			for _, symbol := range symbols {
				sb := new(netproto.TGW_Symbol)
				sb.SymbolID = proto.Int32(int32(symbol.GetId()))
				sb.SymbolName = proto.String(symbol.GetSymbol())
				sbMap[sb.GetSymbolName()] = sb
			}
		}

		for _, row := range sbMap {
			setting.Symbols = append(setting.Symbols, row)
		}

		data.ReelSetting = append(data.ReelSetting, setting)
	}

	sort.Slice(data.ReelSetting, func(i, j int) bool {
		return data.ReelSetting[i].GetRealID() < data.ReelSetting[j].GetRealID()
	})

	// SymbolPayTable
	for _, row := range logic.Proto.SlotPayTable {
		pay := new(netproto.TGW_SymbolPayTable)
		i, _ := strconv.ParseFloat(row.GetPayout(), 64)
		pay.Symbol = &netproto.TGW_Symbol{SymbolID: proto.Int32(0), SymbolName: proto.String(row.GetSymbol())}
		pay.Count = proto.Int32(row.GetCount())
		pay.Type = proto.String(row.GetPayType())

		switch pay.GetType() {
		case slot.PayTypeChip.String(): // 一般賠錢，回傳賠率
			pay.Payout = proto.Int32(int32(i))

		case slot.PayTypePickGame.String(): // BonusGame，找出對應的PickGameID，在找出可翻牌次數
			var pickgame *protodata.PickGameTable
			for _, row := range logic.Proto.PickGameTable {
				if row.GetId() == uint32(i) {
					pickgame = row
					break
				}
			}
			if pickgame != nil {
				pay.Payout = proto.Int32(int32(pickgame.GetPickCount()))
			}
		case slot.PayTypeJackpot.String(): // Jackpot，找出Jackpot最大出獎金額
			if jackpot, ok := logic.Jackpots[slot.JackpotSingle]; ok {
				pay.Payout = proto.Int32(int32(jackpot.GetCurrentPool()))
			}
		}

		// 若有設定符號文案則寫上去
		if desc, ok := logic.SymbolsDesc[row.GetSymbol()]; ok {
			pay.Title = proto.String(desc.Title)
			pay.Desc = proto.String(desc.Description)
		}

		data.SymbolPayTable = append(data.SymbolPayTable, pay)
	}

	// PayLines
	for _, row := range logic.Proto.SlotPayLine {
		payline := slot.NewPayLine(row)

		line := new(netproto.TGW_PayLine)
		line.PayLineID = proto.Int32(int32(row.GetLineId()))

		for _, p := range payline.Positions {
			pos := new(netproto.TGW_ReelPosition)
			pos.Reel = proto.Int32(p.Reel)
			pos.Row = proto.Int32(p.Row)
			line.Positions = append(line.Positions, pos)
		}
		data.PayLines = append(data.PayLines, line)
	}

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendStateMessage 發送遊戲狀態資訊
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10002 (TGW_GameMessageClassID_TGW_StatusID)
//   * Message = { State: 0 }
func (logic *Logic) SendStateMessage(user *gamelogic.GameUser) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TGW_GameMessageClassID_TGW_StatusID)

	data := new(netproto.DFDC_GameStatus)

	data.State = proto.Int32(int32(logic.GameState))

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendJackpotPoolMessage 發送Jackpot獎池資訊
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10003 (TGW_GameMessageClassID_TGW_JackpotPoolID)
//   * Message = { Money: 500000 }
func (logic *Logic) SendJackpotPoolMessage(user *gamelogic.GameUser) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TGW_GameMessageClassID_TGW_JackpotPoolID)

	data := new(netproto.TGW_MSG_JackpotPool)

	data.Money = make([]int64, 0, len(logic.Jackpots))

	for _, jackpot := range logic.Jackpots {
		data.Money = append(data.Money, jackpot.GetCurrentPool())
	}

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendBetRetMessage 發送下注後檢查狀態
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10004 (TGW_GameMessageClassID_TGW_BetStatusID)
//   * Message = { Ret: 0, Money: 1000 }
func (logic *Logic) SendBetRetMessage(user *gamelogic.GameUser, ret int32, money int64) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TGW_GameMessageClassID_TGW_BetStatusID)

	data := new(netproto.DFDC_BetRet)
	data.Ret = proto.Int32(ret)
	data.Money = proto.Int64(money)

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendGameResultMessage 發送遊戲結果
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10005 (TGW_GameMessageClassID_TGW_GameResultID)
//   * Message = { BetMoney: 100, WinMoney: 0, FinalMoney: 300, ReelResult: [{},{},...] }
func (logic *Logic) SendGameResultMessage(user *gamelogic.GameUser) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TGW_GameMessageClassID_TGW_GameResultID)

	data := new(netproto.TGW_MSG_GameResult)

	data.BetMoney = proto.Int64(logic.GameBetMoney)
	data.WinMoney = proto.Int64(logic.GameWinMoney - logic.GameWinJackpot)
	data.FinalMoney = proto.Int64(user.Score)
	data.JackpotMoney = proto.Int64(logic.GameWinJackpot)

	// 中獎內容，賠付線結果(派彩線)
	for _, row := range logic.GamePayResult {
		for _, line := range row.SymbolPayLines {
			paygroup := new(netproto.TGW_PayGroup)
			payline := new(netproto.TGW_PayLine)

			if row.IsWinJackpot() {
				paygroup.WinMoney = proto.Int64(logic.GameWinJackpot)
			} else {
				_, win := row.CalcWinMoney(logic.GameBetMoney, 1)
				paygroup.WinMoney = proto.Int64(win / int64(row.WayCount))
			}
			payline.PayLineID = proto.Int32(int32(line.PayLine.Origin.GetLineId()))

			for _, orgpos := range line.PayLine.Positions {
				pos := new(netproto.TGW_ReelPosition)
				pos.Reel = proto.Int32(orgpos.Reel)
				pos.Row = proto.Int32(orgpos.Row)
				payline.Positions = append(payline.Positions, pos)
			}

			paygroup.PayLines = payline
			data.PayGroups = append(data.PayGroups, paygroup)
		}
	}

	// 輪軸結果
	logic.GameResultManager.RowMode()
	for row, reels := range logic.GameResultManager.GameResult {
		for reel, symbol := range reels {
			symbolResult := &netproto.TGW_ReelResult{
				Position: &netproto.TGW_ReelPosition{Reel: proto.Int32(int32(reel + 1)), Row: proto.Int32(int32(row + 1))},
				Symbol:   &netproto.TGW_Symbol{SymbolID: proto.Int32(int32(symbol.GetId())), SymbolName: proto.String(symbol.GetSymbol())},
				Blink:    proto.Bool(false),
			}

			// 檢查這個位置是否有中獎
			for _, result := range logic.GamePayResult {
				for _, line := range result.SymbolPayLines {
					for _, pos := range line.PayLine.Positions {
						if pos.Reel > result.PayCount {
							continue
						}
						if reflect.DeepEqual(symbolResult.GetPosition(), &netproto.TGW_ReelPosition{Reel: proto.Int32(pos.Reel), Row: proto.Int32(pos.Row)}) {
							symbolResult.Blink = proto.Bool(true)
						}
					}
				}
			}

			data.ReelResult = append(data.ReelResult, symbolResult)
		}
	}

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendBonusGameMessage 發送BonusGame內容，傳給 Client 的 RewareValue 為選中的倍率
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10006 (TGW_GameMessageClassID_TGW_GameResultID)
//   * Message = { WinMoney: 120, Count: 3, Records: [ { Id: 1, Location: 0, RewareValue: 10 }, { Id: 2, Location: 0, RewareValue: 0 }, { Id: 3, Location: 0, RewareValue: 10 }, ... ] }
func (logic *Logic) SendBonusGameMessage(user *gamelogic.GameUser) {

	if logic.PickGame == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TGW_GameMessageClassID_TGW_BonusGameID)

	data := new(netproto.TGW_MSG_BonusGame)

	data.WinMoney = proto.Int64(logic.PickGame.WinMoney)
	data.Count = proto.Int32(logic.PickGame.GameCount)

	for _, detail := range logic.PickGame.Process {
		rate, _ := strconv.ParseFloat(detail.Origin.GetRewardValue(), 64)
		record := &netproto.TGW_BonusGameRecord{ID: proto.Int32(detail.ID), Location: proto.Int32(detail.Location), RewareValue: proto.Int64(int64(rate))}
		data.Records = append(data.Records, record)
	}

	sort.Slice(data.Records, func(i, j int) bool {
		return data.Records[i].GetID() < data.Records[j].GetID()
	})

	for _, row := range logic.PickGame.All {
		data.All = append(data.All, int32(row/logic.PickGame.GameBetMoney))
	}

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendBonusGameRetMessage 傳送確認玩家選擇封包
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 12004 (TGW_GameMessageClassID_TGW_ActionPickRetID)
//   * Message = { TGW_BonusGameRecord: { Id: 1, Location: 0, Money: 10 } }
func (logic *Logic) SendBonusGameRetMessage(user *gamelogic.GameUser, record *netproto.TGW_BonusGameRecord) {
	if user == nil || record == nil {
		return
	}
	// 收到資料後，直接轉傳給client
	// 給client的封包
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TGW_GameMessageClassID_TGW_ActionPickRetID)

	data := new(netproto.TGW_MSG_PickRet)
	data.Record = record

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendGameLogMessage 傳送玩家遊戲紀錄
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 12002 (TGW_GameMessageClassID_TGW_LoadHistoryRetID)
//   * Message = { WinMoney: 0, GoldenSymbol: 8, BetMoney: 10 }
func (logic *Logic) SendGameLogMessage(user *gamelogic.GameUser, userID, gameID, pageCnt int32, history []*netproto.Slot_SendHistoryLog) {
	// 收到資料後，直接轉傳給client
	// 給client的封包
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TGW_GameMessageClassID_TGW_LoadHistoryRetID)

	data := new(netproto.Slot_SendHistoryData)
	data.UserID = proto.Int32(userID)
	data.GameID = proto.Int32(gameID)
	data.SendHistoryLog = history
	data.PageCount = proto.Int32(pageCnt)

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}
