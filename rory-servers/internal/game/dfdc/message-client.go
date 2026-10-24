package dfdc

import (
	"sort"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// SendSettingMessage 發送遊戲設置內容
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10001 (DFDC_GameMessageClassID_DFDCSettingID)
//   * Message = { BetCount: [10, 20, 30, 60, 100, 200, 300, 600, 1000], RowCounts: 3, TotalPayLine: 1, BigWinSetp: [300, 400, 500, 500000], PayLineMode: 1, ReelSetting: [ { ReelID: 1, SymbolGroup: [ { SymbolID: 1, SymbolName: "p1" }, ... ] }, ... ], SymbolPayTable: [ { Rate: 5, SymbolName: "p1", Count: 3 }, ... ] }
func (logic *Logic) SendSettingMessage(user *gamelogic.GameUser) {

	if user == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DFDC_GameMessageClassID_DFDCSettingID)
	data := new(netproto.DFDC_Setting)

	data.BetCount = logic.BetCounts
	data.RowCounts = proto.Int32(logic.DefaultRowNum)
	data.TotalPayLine = proto.Int32(1)
	data.BigWinStep = logic.BigWinRule
	data.PayLineMode = proto.Int32(logic.PayLineMode.Int())
	data.MinCarryMoney = proto.Int64(logic.MinCarryMoney)

	// ReelSetting
	for _, row := range logic.NetprotoSlotReelSetting {
		setting := new(netproto.DFDC_ReelSetting)
		setting.RealID = proto.Int32(row.GetReelID())
		for _, symbol := range row.GetSymbols() {
			sg := new(netproto.DFDC_SymbolGroup)
			sg.SymbolID = proto.Int32(symbol.GetSymbolID())
			sg.SymbolName = proto.String(symbol.GetSymbolName())
			setting.SymbolGroup = append(setting.SymbolGroup, sg)
		}
		data.ReelSetting = append(data.ReelSetting, setting)
	}

	// SymbolPayTable
	for _, row := range logic.NetprotoSlotSymbolPayTable {

		pay := new(netproto.DFDC_SymbolPayTable)

		pay.Rate = proto.Int32(row.GetPayout())
		pay.SymbolName = proto.String(row.GetSymbol().GetSymbolName())
		pay.Count = proto.Int32(row.GetCount())

		data.SymbolPayTable = append(data.SymbolPayTable, pay)
	}

	for key, val := range logic.UnitPrefix {
		data.UnitPrefix = append(data.UnitPrefix, &netproto.Slot_UnitPrefix{Name: proto.String(key), Multiple: proto.Int64(val)})
	}

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendGoldenSymbolMessage 發送黃金符號，設定檔中的倍率BetRate列表
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10011 (DFDC_GameMessageClassID_DFDCGoldenSymbolID)
//   * Message = { GoldenDetail: [ { Symbols: 1, Credits: 8 }, { Symbols: 2, Credits: 18 }, { Symbols: 3, Credits: 38 }, { Symbols: 4, Credits: 68 }, { Symbols: 5, Credits: 88 } ] }
func (logic *Logic) SendGoldenSymbolMessage(user *gamelogic.GameUser) {

	if user == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DFDC_GameMessageClassID_DFDCGoldenSymbolID)
	data := new(netproto.DFDC_GoldenSymbolList)

	data.GoldenDetail = make([]*netproto.DFDC_GoldenSymbolDetail, 0, len(logic.BetRates))

	for i, rate := range logic.BetRates {
		detail := new(netproto.DFDC_GoldenSymbolDetail)
		detail.Symbols = proto.Int32(int32(i + 1))
		detail.Credits = proto.Int32(rate)
		data.GoldenDetail = append(data.GoldenDetail, detail)
	}

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendJackpotPoolMessage 發送Jackpot獎池資訊
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10009 (DFDC_GameMessageClassID_DFDCJackpotPoolS2CID)
//   * Message = { JackpotPoolID: [ { JackpotID: 1, Level: "Grand", Syspool: 20000000 }, { JackpotID: 2, Level: "Major", Syspool: 7500000 }, { JackpotID: 3, Level: "Minor", Syspool: 375000 }, { JackpotID: 4, Level: "Mini", Syspool: 200000 } ] }
func (logic *Logic) SendJackpotPoolMessage(user *gamelogic.GameUser) {

	if user == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DFDC_GameMessageClassID_DFDCJackpotPoolS2CID)

	data := new(netproto.DFDC_JackpotPoolS2CMsg)

	jackpots := make([]*netproto.DFDC_JackpotPoolDetail, 0, len(logic.Jackpots))

	for _, jackpot := range logic.Jackpots {
		detail := new(netproto.DFDC_JackpotPoolDetail)
		detail.JackpotID = proto.Int32(jackpot.Level.Int())
		detail.Syspool = proto.Int64(jackpot.GetCurrentPool())
		detail.Level = proto.String(jackpot.Level.String())

		jackpots = append(jackpots, detail)
	}

	data.JackpotPoolID = jackpots

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendStateMessage 發送遊戲狀態資訊
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10004 (DFDC_GameMessageClassID_DFDCStateID)
//   * Message = { State: 0 }
func (logic *Logic) SendStateMessage(user *gamelogic.GameUser) {

	if user == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DFDC_GameMessageClassID_DFDCStateID)

	data := new(netproto.DFDC_GameStatus)

	data.State = proto.Int32(int32(logic.GameState))

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendBetRetMessage 發送下注後檢查狀態
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10010 (DFDC_GameMessageClassID_DFDCBetRetID)
//   * Message = { Ret: 0, Money: 1000 }
func (logic *Logic) SendBetRetMessage(user *gamelogic.GameUser, ret int32, money int64) {

	if user == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DFDC_GameMessageClassID_DFDCBetRetID)

	data := new(netproto.DFDC_BetRet)
	data.Ret = proto.Int32(ret)
	data.Money = proto.Int64(money)

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendGameResultMessage 發送遊戲結果
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10003 (DFDC_GameMessageClassID_DFDCGameResult)
//   * Message = { Ret: 0, Money: 1000 }
func (logic *Logic) SendGameResultMessage(user *gamelogic.GameUser) {

	if user == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DFDC_GameMessageClassID_DFDCGameResult)

	data := new(netproto.DFDC_GameResultReel)

	// Symbol ID不等於開獎結果，而是透過Symbol ID骰出該order id 之結果
	data.SymbolID = logic.GameResultManager.GameResult.GetSymbolIDByRow(1)

	// ReelSetting
	reelResult := logic.BuildReelResult(logic.GameResultManager, logic.GamePayResult)
	for _, row := range reelResult {
		setting := new(netproto.DFDC_ReelSetting)
		setting.RealID = proto.Int32(row.GetPosition().GetReel())
		setting.RowIndex = proto.Int32(row.GetPosition().GetRow())
		setting.SymbolGroup = []*netproto.DFDC_SymbolGroup{
			{SymbolID: proto.Int32(int32(row.GetSymbol().GetSymbolID())), SymbolName: proto.String(row.GetSymbol().GetSymbolName())},
		}
		data.ReelSetting = append(data.ReelSetting, setting)
	}

	// PayAllGroup
	payAllGroup := new(netproto.DFDC_PayAllGroup)

	var SymbolIDList []int32

	for _, result := range logic.GamePayResult {
		SymbolIDList = append(SymbolIDList, result.SymbolIDs...)

		switch result.GetPayType() {
		default:
			continue
		case slot.PayTypeChip, slot.PayTypeChipXRate:
		}

		detail := new(netproto.DFDC_PayAllDetail)
		detail.SymbolIDList = result.SymbolIDs
		detail.SymbolName = proto.String(result.SymbolName)
		detail.PayCount = proto.Int32(result.PayCount)
		detail.WayCount = proto.Int32(result.WayCount)
		_, money := result.CalcWinMoney(logic.GameBetMoney, logic.GameBetRate)

		detail.PayRate = proto.Int32(int32(money))
		detail.FeaturePays = proto.Bool(false)
		payAllGroup.PayAllDetail = append(payAllGroup.PayAllDetail, detail)
	}

	symbolData := logic.GameResultManager.GameResult.GetReelSymbol()
	for reelID, reel := range symbolData {
		for _, row := range reel {
			symbol := row[0]
			detail := new(netproto.DFDC_PayLineDetail)
			detail.SymbolID = proto.Int32(int32(symbol.GetId()))
			detail.ReelID = proto.Int32(int32(reelID))
			detail.Line = proto.Bool(false)
			detail.Blink = proto.Bool(false)
			// 檢查是否為連線符號
			for _, line := range SymbolIDList {
				if int32(symbol.GetId()) == line {
					detail.Blink = proto.Bool(true)
					break
				}
			}
			payAllGroup.PayLineDetail = append(payAllGroup.PayLineDetail, detail)
		}
	}

	data.PayAllGroup = []*netproto.DFDC_PayAllGroup{payAllGroup}
	data.PayLineGroup = []*netproto.DFDC_PayLineGroup{}

	if logic.FreeSpin != nil {
		data.FreeGameWinMoney = proto.Int64(logic.FreeSpin.WinMoney)
	} else {
		data.FreeGameWinMoney = proto.Int64(0)
	}

	data.WinMoney = proto.Int64(logic.GameWinMoney)

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendFreeSpinCountMessage 傳送FreeSpin數量給Client
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10005 (DFDC_GameMessageClassID_DFDCFreeOrPayID)
//   * Message = { FreeCount: 0 }
func (logic *Logic) SendFreeSpinCountMessage(user *gamelogic.GameUser) {

	if user == nil {
		return
	}

	if logic.FreeSpin == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DFDC_GameMessageClassID_DFDCFreeOrPayID)

	data := new(netproto.DFDC_FreeOrPayMsg)

	data.FreeCount = proto.Int32(int32(logic.FreeSpin.Count))

	// ReelSetting
	var info *netproto.Slot_FreeGameInfo
	for _, row := range logic.NetprotoFreeGameInfo {
		if row.GetID() == int32(logic.FreeSpin.Origin.GetId()) {
			info = row
		}
	}

	if info != nil {
		for _, row := range info.ReelSetting {
			setting := new(netproto.DFDC_ReelSetting)
			setting.RealID = proto.Int32(row.GetReelID())
			for _, symbol := range row.GetSymbols() {
				sg := new(netproto.DFDC_SymbolGroup)
				sg.SymbolID = proto.Int32(symbol.GetSymbolID())
				sg.SymbolName = proto.String(symbol.GetSymbolName())
				setting.SymbolGroup = append(setting.SymbolGroup, sg)
			}
			data.ReelSetting = append(data.ReelSetting, setting)
		}
	}

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendJackpotGameDetailMessage 發送BonusGame內容
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10012 (DFDC_GameMessageClassID_DFDCBonusGameID)
//   * Message = { Money: 18181, BonusDetail: [ { Id: 1, Place: 0, Icon: 4 }, ... ] }
func (logic *Logic) SendJackpotGameDetailMessage(user *gamelogic.GameUser) {

	if user == nil {
		return
	}

	// 沒有JackpotGame不做事
	if logic.JackpotGame == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DFDC_GameMessageClassID_DFDCBonusGameID)

	data := new(netproto.DFDC_BonusGame)
	data.Money = proto.Int64(logic.JackpotGame.WinMoney)
	data.BonusDetail = make([]*netproto.DFDC_BonusGameDetail, 0, len(logic.JackpotGame.Process))

	for _, detail := range logic.JackpotGame.Process {
		d := new(netproto.DFDC_BonusGameDetail)
		d.Id = proto.Int32(detail.ID)
		d.Place = proto.Int32(detail.Place)
		d.Icon = proto.Int32(detail.Icon)

		data.BonusDetail = append(data.BonusDetail, d)
	}

	sort.Slice(data.BonusDetail, func(i, j int) bool {
		return data.BonusDetail[i].GetId() < data.BonusDetail[j].GetId()
	})

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendFreeSpinBetMessage 傳送玩家上局下注資訊
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10019 (DFDC_GameMessageClassID_DFDCUserRecoverID)
//   * Message = { WinMoney: 0, GoldenSymbol: 8, BetMoney: 10 }
func (logic *Logic) SendFreeSpinBetMessage(user *gamelogic.GameUser) {

	if user == nil {
		return
	}

	if logic.FreeSpin == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DFDC_GameMessageClassID_DFDCUserRecoverID)

	data := new(netproto.DFDC_UserRecover)
	data.WinMoey = proto.Int64(logic.FreeSpin.WinMoney)
	data.BetMoney = proto.Int64(logic.FreeSpin.BetMoney)
	data.GoldenSymbol = proto.Int32(int32(logic.FreeSpin.BetRate))

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}
