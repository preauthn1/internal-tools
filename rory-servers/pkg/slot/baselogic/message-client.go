package baselogic

import (
	"encoding/json"
	"fmt"
	"math"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// GameMessage 接收遊戲訊息
func (logic *SlotLogicBase) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("GameMessageFunc 發生錯誤, %v\n%s", rec, string(debug.Stack()))
		}
	}()

	if !logic.CheckGameUserMessage(user) {
		logic.KickUserFunc(user.UserID, "非法操作")
		return
	}

	// 空訊息，不處理
	if gmsg == nil {
		return
	}

	// 非遊戲訊息，不處理
	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	// 玩家查詢記錄
	if gmsg.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LoadHistoryID) {

		data := &netproto.Slot_MSG_GameLogRead{}
		if err := proto.Unmarshal(gmsg.MsgData.([]byte), data); err != nil {
			mlog.Error("接收玩家选择封包解析错误. %v", err)
			return
		}

		logic.ReadSlotGameHistoryLog(user, data.GetPageIndex(), data.GetPageSize(), data.GetStartTime(), data.GetEndTime())
		return
	}

	// Client請求讀取Jackpot
	if gmsg.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_JackpotPoolID) {

		var ids []int32

		if gmsg.MsgData != nil && len(gmsg.MsgData.([]byte)) > 0 {
			data := &netproto.Slot_MSG_GetJackpot{}
			if err := proto.Unmarshal(gmsg.MsgData.([]byte), data); err != nil {
				mlog.Error("[玩家查詢Jackpot Pool] 接收玩家封包解析錯誤. %v", err)
				return
			}
			ids = append(ids, data.GetID()...)
		}

		logic.SendJackpotMessageToUser(user, ids...)
	}

	if logic.GameMessageFunc != nil {
		logic.GameMessageFunc(user, gmsg)
	}
}

// SendReelSettingMessage 發送遊戲排列帶內容
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 20010 (Slot_GameMessageClassID_Slot_GameReelSettingID)
func (logic *SlotLogicBase) SendReelSettingMessage(user *gamelogic.GameUser) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_GameReelSettingID)
	data := new(netproto.Slot_MSG_ReelSetting)

	for id, settings := range logic.NetprotoSlotReelsSetting {
		obj := new(netproto.Slot_GameReels)
		obj.ID = proto.Int32(id)
		obj.Name = proto.String(fmt.Sprintf("MainGame(%d)", id))
		obj.ReelSetting = settings
		for _, row := range logic.Proto.MainGameTable {
			if int32(row.GetId()) == obj.GetID() {
				obj.ConditionType = proto.String(row.GetGoldCondition())
				obj.ConditionValue = proto.String(row.GetConditionValue())
				break
			}
		}

		data.GameReels = append(data.GameReels, obj)
	}

	sort.Slice(data.GameReels, func(i, j int) bool {
		return data.GameReels[i].GetID() < data.GameReels[j].GetID()
	})

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendSettingMessage 發送遊戲設置內容
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 20001 (Slot_GameMessageClassID_Slot_SettingID)
//   * Message = { BetCount: [10, 20, 30, 60, 100, 200, 300, 600, 1000], RowCounts: 3, TotalPayLine: 1, BigWinSetp: [300, 400, 500, 500000], PayLineMode: 1, ReelSetting: [ { ReelID: 1, SymbolGroup: [ { SymbolID: 1, SymbolName: "p1" }, ... ] }, ... ], SymbolPayTable: [ { Rate: 5, SymbolName: "p1", Count: 3 }, ... ] }
func (logic *SlotLogicBase) SendSettingMessage(user *gamelogic.GameUser) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_SettingID)
	data := new(netproto.Slot_MSG_Setting)

	data.BetMoneys = logic.BetCounts
	data.BetRates = logic.BetRates
	data.RowCounts = proto.Int32(logic.DefaultRowNum)
	data.BigWinStep = logic.BigWinRule
	data.PayLineMode = proto.Int32(logic.PayLineMode.Int())
	data.MinCarryMoney = proto.Int64(logic.MinCarryMoney)

	switch logic.PayLineMode {
	default:
	case slot.PayLineModeAll:
		data.TotalPayLine = proto.Int32(int32(math.Pow(float64(logic.DefaultRowNum), float64(logic.DefaultReelNum))))
	case slot.PayLineModeFixed:
		data.TotalPayLine = proto.Int32(int32(len(logic.Proto.SlotPayLine)))
	}

	// UnitPrefix
	if logic.UnitPrefix != nil {
		for name, multiple := range logic.UnitPrefix {
			unit := new(netproto.Slot_UnitPrefix)
			unit.Name = proto.String(name)
			unit.Multiple = proto.Int64(multiple)
			data.UnitPrefix = append(data.UnitPrefix, unit)
		}
	}

	// Description
	if logic.Descriptions != nil {
		for key, val := range logic.Descriptions {
			desc := new(netproto.Slot_Description)
			desc.Name = proto.String(key)
			desc.Desc = proto.String(val.Description)
			desc.Title = proto.String(val.Title)
			data.Descriptions = append(data.Descriptions, desc)
		}
	}

	// PayLine
	data.PayLines = logic.NetprotoSlotPayLine

	// ReelSetting
	data.ReelSetting = logic.NetprotoSlotReelSetting

	// SymbolPayTable
	data.SymbolPayTable = logic.NetprotoSlotSymbolPayTable

	// FreeGameInfo
	data.FreeGames = logic.NetprotoFreeGameInfo

	// GoldBallSettings
	data.GoldBallSetting = logic.NetprotoGoldBallInfo

	// JackpotInfo
	for _, jackpot := range logic.Jackpots {
		row := new(netproto.Slot_JackpotInfo)
		row.Money = proto.Int64(jackpot.GetCurrentPool())
		row.Name = proto.String(jackpot.Level.String())
		row.ID = proto.Int32(jackpot.Origin.GetId())
		if jackpot.Origin.GetIsLinked() == 1 {
			row.IsLink = proto.Bool(true)
		} else {
			row.IsLink = proto.Bool(false)
		}
		data.Jackpots = append(data.Jackpots, row)
	}

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendStateMessageToUser 發送遊戲狀態資訊
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 20002 (Slot_GameMessageClassID_Slot_StatusID)
//   * Message = { State: 0 }
func (logic *SlotLogicBase) SendStateMessageToUser(user *gamelogic.GameUser) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_StatusID)

	data := new(netproto.Slot_MSG_Status)

	data.State = proto.Int32(int32(logic.GameState))

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendBetStatusMessageToUser 發送下注後檢查狀態
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 20003 (Slot_GameMessageClassID_Slot_BetStatusID)
//   * Message = { Ret: 0, Money: 1000 }
func (logic *SlotLogicBase) SendBetStatusMessageToUser(user *gamelogic.GameUser, ret int32, money int64) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_BetStatusID)

	data := new(netproto.Slot_MSG_BetStatus)
	data.Ret = proto.Int32(ret)
	data.Money = proto.Int64(money)

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendLeftFreeGameToUser 發送免費遊戲剩於次數
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 20005 (Slot_GameMessageClassID_Slot_LeftFreeGameID)
//   * Message = { ID: 12345, Count: 8, BetMoney: 20, BetRate: 1 }
func (logic *SlotLogicBase) SendLeftFreeGameToUser(user *gamelogic.GameUser) {

	if logic.FreeSpin == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_LeftFreeGameID)

	data := new(netproto.Slot_MSG_LeftFreeGameData)
	data.ID = proto.Int32(int32(logic.FreeSpin.Origin.GetId()))
	data.Count = proto.Int32(int32(logic.FreeSpin.Count))
	data.BetMoney = proto.Int64(logic.FreeSpin.BetMoney)
	data.BetRate = proto.Int64(logic.FreeSpin.BetRate)

	if logic.FreeSpin.FreeSpinExtAward != nil {
		data.FreeGameExtAward = new(netproto.Slot_FreeGameExtAward)
		data.FreeGameExtAward.TotalCount = proto.Int32(logic.FreeSpin.FreeSpinExtAward.TotalCount)
		data.FreeGameExtAward.TotalWinMoney = proto.Int64(logic.FreeSpin.FreeSpinExtAward.TotalWinMoney)
		data.FreeGameExtAward.Count = proto.Int32(logic.FreeSpin.FreeSpinExtAward.Count)
		data.FreeGameExtAward.WinMoney = proto.Int64(logic.FreeSpin.FreeSpinExtAward.WinMoney)
		data.FreeGameExtAward.EachCountMoney = proto.Int64(logic.FreeSpin.FreeSpinExtAward.EachCountMoney)
		data.FreeGameExtAward.IsFinished = proto.Bool(logic.FreeSpin.FreeSpinExtAward.IsFinished)

		for pos, symbol := range logic.FreeSpin.FreeSpinExtAward.Symbols {
			s := new(netproto.Slot_Symbol)
			s.SymbolID = proto.Int32(int32(symbol.GetId()))
			s.SymbolName = proto.String(symbol.GetSymbol())
			s.Position = &netproto.Slot_ReelPosition{
				Reel: proto.Int32(pos.Reel),
				Row:  proto.Int32(pos.Row),
			}
			data.FreeGameExtAward.Symbols = append(data.FreeGameExtAward.Symbols, s)
		}
	}

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// required int32 Count = 1;     // 剩餘次數
// required int64 BetMoney = 2;  // 下注金額
// optional int64 BetRate = 3;   // 下注倍率符號
// optional int32 ID = 4;        // 免費遊戲ID

// SendGameLogMessageToUser 傳送玩家遊戲紀錄
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 22002 (Slot_GameMessageClassID_Slot_LoadHistoryRetID)
//   * Message = { UserID: 12345, GameID: 26, PageCount: 2, SendHistoryLog: [ { RowNum: 1, CreateTime: 2019-10-10, OrderID: 1, BetMoney: 8, Profit: -8 }... ] }
func (logic *SlotLogicBase) SendGameLogMessageToUser(user *gamelogic.GameUser, userID, gameID, pageCnt int32, history []*netproto.Slot_GameLog) {
	// 收到資料後，直接轉傳給client
	// 給client的封包
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_LoadHistoryRetID)

	data := new(netproto.Slot_MSG_GameLogRecord)
	data.UserID = proto.Int32(userID)
	data.GameID = proto.Int32(gameID)
	data.Logs = history
	data.PageCount = proto.Int32(pageCnt)

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendJackpotMessageToUser 發送Jackpot獎池資訊
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 22004 (Slot_GameMessageClassID_Slot_JackpotPoolRetID)
//   * Message = [ { Money: 500000, Name: , ID: 1 }, ]
func (logic *SlotLogicBase) SendJackpotMessageToUser(user *gamelogic.GameUser, ids ...int32) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_JackpotPoolRetID)

	data := new(netproto.Slot_MSG_JackpotPool)

	for _, jackpot := range logic.Jackpots {

		if len(ids) > 0 {
			var exist bool
			for _, id := range ids {
				if jackpot.Origin.GetId() == id {
					exist = true
					break
				}
			}
			if !exist {
				continue
			}
		}

		row := new(netproto.Slot_JackpotInfo)
		row.Money = proto.Int64(jackpot.GetCurrentPool())
		row.Name = proto.String(jackpot.Level.String())
		row.ID = proto.Int32(jackpot.Origin.GetId())

		if jackpot.Origin.GetIsLinked() == 1 {
			row.IsLink = proto.Bool(true)
		} else {
			row.IsLink = proto.Bool(false)
		}

		data.Jackpots = append(data.Jackpots, row)
	}

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendJackpotGameDetailMessage 發送BonusGame內容
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 20009 (Slot_GameMessageClassID_Slot_BonusGameID)
//   * Message = { Money: 18181, BonusDetail: [ { Id: 1, Place: 0, Icon: 4 }, ... ] }
func (logic *SlotLogicBase) SendJackpotGameDetailMessage(user *gamelogic.GameUser) {

	if user == nil {
		return
	}

	// 沒有JackpotGame不做事
	if logic.JackpotGame == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_BonusGameID)

	data := new(netproto.Slot_MSG_BonusGame)
	data.Type = proto.String(string(slot.GameModeJackpotGame))
	data.Count = proto.Int32(int32(len(logic.JackpotGame.Process)))
	data.WinMoney = proto.Int64(logic.JackpotGame.WinMoney)
	data.Records = make([]*netproto.Slot_BonusGameRecord, 0, len(logic.JackpotGame.Process))
	data.BetMoney = proto.Int64(logic.JackpotGame.GameBetMoney)
	data.BetRate = proto.Int64(int64(logic.JackpotGame.GameBetRate))

	for _, detail := range logic.JackpotGame.Process {
		d := new(netproto.Slot_BonusGameRecord)
		d.ID = proto.Int32(detail.ID)
		d.Location = proto.Int32(detail.Place)
		d.RewareValue = proto.Int64(int64(detail.Icon))

		data.Records = append(data.Records, d)
	}

	sort.Slice(data.Records, func(i, j int) bool {
		return data.Records[i].GetID() < data.Records[j].GetID()
	})

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendPickGameDetailMessage 發送PickGame內容
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 20009 (Slot_GameMessageClassID_Slot_BonusGameID)
//   * Message = { Money: 18181, BonusDetail: [ { Id: 1, Place: 0, Icon: 4 }, ... ] }
func (logic *SlotLogicBase) SendPickGameDetailMessage(user *gamelogic.GameUser) {

	if user == nil {
		return
	}

	// 沒有 PickGame 不做事
	if logic.PickGame == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_BonusGameID)

	data := new(netproto.Slot_MSG_BonusGame)
	data.Type = proto.String(string(slot.GameModePickGame))
	data.Count = proto.Int32(int32(len(logic.PickGame.Process)))
	data.WinMoney = proto.Int64(logic.PickGame.WinMoney)
	data.Records = make([]*netproto.Slot_BonusGameRecord, 0, len(logic.PickGame.Process))
	data.BetMoney = proto.Int64(logic.PickGame.GameBetMoney)
	data.BetRate = proto.Int64(int64(logic.PickGame.GameBetRate))

	for _, detail := range logic.PickGame.Process {
		d := new(netproto.Slot_BonusGameRecord)
		d.ID = proto.Int32(detail.ID)
		d.Location = proto.Int32(detail.Location)
		d.RewareValue = proto.Int64(detail.Money)

		data.Records = append(data.Records, d)
	}

	for _, row := range logic.PickGame.All {
		data.All = append(data.All, int32(row))
	}

	for _, row := range logic.PickGame.Positions {
		data.Positions = append(data.Positions, &netproto.Slot_ReelPosition{Reel: proto.Int32(row.Reel), Row: proto.Int32(row.Row)})
	}

	sort.Slice(data.Records, func(i, j int) bool {
		return data.Records[i].GetID() < data.Records[j].GetID()
	})

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendRecvBonusGameRecord 返還玩家選擇BonusGame的內容
func (logic *SlotLogicBase) SendRecvBonusGameRecord(user *gamelogic.GameUser, data *netproto.Slot_MSG_PickBonus) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_PickBonusRetID)

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendSelectChanceMessage 發送遊戲結果
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 20008 (Slot_GameMessageClassID_Slot_SelectGameOptionsID)
func (logic *SlotLogicBase) SendSelectChanceMessage(user *gamelogic.GameUser) {

	if logic.SelectChances == nil || len(logic.SelectChances) == 0 {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_SelectGameOptionsID)

	data := new(netproto.Slot_MSG_ChooseGameOptions)
	data.Times = proto.Int32(int32(len(logic.SelectChances)))

	// data.Options
	chance := logic.SelectChances[0]

	data.BetMoney = proto.Int64(chance.BetMoney)
	data.BetRate = proto.Int64(chance.BetRate)

	for _, row := range chance.Options {
		opt := new(netproto.Slot_ChooseGameOption)
		opt.ID = proto.Int32(row.GameID)
		opt.Type = proto.String(string(row.GameMode))

		switch row.GameMode {
		case slot.GameModeFreeSpin:
			tb := row.GameTable.(*protodata.FreeSpinTable)
			dMap := map[string]interface{}{
				"count": tb.GetFreeSpinCount(),
				"info":  tb.GetInfo(),
			}
			if b, err := json.Marshal(dMap); err == nil {
				opt.Data = proto.String(string(b))
			}
		case slot.GameModeLightning:
			tb := row.GameTable.(*protodata.RespinTable)
			dMap := map[string]interface{}{
				"count": tb.GetRespinCount(),
			}
			if b, err := json.Marshal(dMap); err == nil {
				opt.Data = proto.String(string(b))
			}
		case slot.GameModeBonus:
			items := row.GameTable.([]*protodata.BonusTable)
			dMap := map[string]interface{}{
				"length": len(items),
			}
			if b, err := json.Marshal(dMap); err == nil {
				opt.Data = proto.String(string(b))
			}
		}

		data.Options = append(data.Options, opt)
	}
	if b, err := json.Marshal(chance.Data); err == nil {
		data.Data = proto.String(string(b))
	}

	data.ReelResult = logic.BuildReelResult(chance.GameResultManager, nil)

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}

// SendIncreaseGameTimesMessage 發送增加遊戲次數封包
func (logic *SlotLogicBase) SendIncreaseGameTimesMessage(user *gamelogic.GameUser, gtype string, originID, cnt int32, data string) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_IncreaseGameTimesID)

	m := new(netproto.Slot_MSG_IncreaseGameTimes)

	m.Type = proto.String(gtype)
	m.Count = proto.Int32(cnt)
	m.ID = proto.Int32(originID)
	m.Data = proto.String(data)

	pmsg.MsgData = m

	logic.SendMessageToUser(user, pmsg)
}

// SendLightningPrepareMessage 傳送Lightning準備畫面
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 20006 (Slot_GameMessageClassID_Slot_LightningPrepareID)
//   * Message = { LightningStep: 0, Count: 3, ReelResult: [{},{},...], LockReels: [ {Reel: 1, Row: 1}, {}, ... ] }
func (logic *SlotLogicBase) SendLightningPrepareMessage(user *gamelogic.GameUser) {

	if logic.LightningGame == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_LightningPrepareID)

	data := new(netproto.Slot_MSG_LightningPrepare)
	data.IsFreeGame = proto.Bool(logic.LightningGame.IsFreeGame)

	step := logic.LightningGame.Steps[logic.LightningGame.ShowStep]

	data.LightningStep = proto.Int32(int32(logic.LightningGame.ShowStep + 1))
	data.Count = proto.Int32(int32(step.LastChance))
	data.BetMoney = proto.Int64(logic.LightningGame.BetMoney)
	data.BetRate = proto.Int64(logic.LightningGame.BetRate)

	var symbolGroupID int32
	if logic.LightningGame.IsLimit && len(logic.LightningGame.Steps) > logic.LightningGame.LimitPoint {
		symbolGroupID = logic.LightningGame.LimitSymbolGroupID
	} else {
		symbolGroupID = logic.LightningGame.SymbolGroupID
	}
	data.ReelSetting = logic.BuildReelSetting(symbolGroupID, logic.DefaultReelNum, logic.DefaultRowNum, logic.Proto.SlotSymbolDistribution)

	step.GameResultManager.RowMode()
	for rowID, row := range step.GameResultManager.GameResult {
		for reelID, symbol := range row {
			if logic.LightningGame.IsGoldBall(symbol) {
				data.LockSymbols = append(data.LockSymbols, &netproto.Slot_LightningLock{
					SymbolID: proto.Int32(int32(symbol.GetId())),
					Position: &netproto.Slot_ReelPosition{
						Reel: proto.Int32(int32(reelID + 1)),
						Row:  proto.Int32(int32(rowID + 1)),
					},
				})
			}
		}
	}

	data.ReelResult = logic.BuildReelResult(step.GameResultManager, nil)

	// 把所有Step的Award都包進去
	for _, _step := range logic.LightningGame.Steps {
		for _, award := range _step.Awards {
			payout := new(netproto.Slot_LightningPayout)
			payout.SymbolID = proto.Int32(int32(award.Symbol.GetId()))

			payout.WinMoney = proto.Int64(award.WinMoney)
			payout.Type = proto.String(award.Type.String())
			if award.Type == slot.PayTypeJackpot {
				payout.JackpotType = proto.String(award.JackpotType.String())
			}
			if award.JackpotType != slot.JackpotMajor {
				payout.Position = &netproto.Slot_ReelPosition{Reel: proto.Int32(award.Position.Reel), Row: proto.Int32(award.Position.Row)}
			} else {
				continue
			}

			data.PastPayGroups = append(data.PastPayGroups, payout)
		}
	}

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)

	// 舊ID相容
	{
		pmsg := new(msg.Message)
		pmsg.BClassID = int32(netproto.MessageBClassID_Game)
		pmsg.SClassID = int32(10001)
		pmsg.MsgData = data
		logic.SendMessageToUser(user, pmsg)
	}
}

// SendLightningResultMessage 傳送Lightning輪軸結果
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 20007 (Slot_GameMessageClassID_Slot_LightningResultID)
//   * Message = { LightningStep: 0, Count: 3, ReelResult: [{},{},...], LockReels: [ {Reel: 1, Row: 1}, {}, ... ] }
func (logic *SlotLogicBase) SendLightningResultMessage(user *gamelogic.GameUser) {

	if logic.LightningGame == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_LightningResultID)
	data := new(netproto.Slot_MSG_LightningResult)

	step := logic.LightningGame.Steps[logic.LightningGame.ShowStep]

	data.BetMoney = proto.Int64(0)
	data.WinMoney = proto.Int64(step.WinMoney)
	data.FinalMoney = proto.Int64(logic.GameUser.Score)
	if logic.FreeSpin != nil {
		data.FreeGameWinMoney = proto.Int64(logic.FreeSpin.WinMoney)
	}

	data.LightningStep = proto.Int32(int32(logic.LightningGame.ShowStep))

	data.Count = proto.Int32(int32(step.LastChance))

	data.ReelResult = logic.BuildReelResult(step.GameResultManager, nil)

	for _, award := range step.Awards {
		payout := new(netproto.Slot_LightningPayout)

		payout.WinMoney = proto.Int64(award.WinMoney)
		payout.Type = proto.String(award.Type.String())
		if award.Type == slot.PayTypeJackpot {
			payout.JackpotType = proto.String(award.JackpotType.String())
		}
		if award.JackpotType != slot.JackpotGrand {
			payout.Position = &netproto.Slot_ReelPosition{Reel: proto.Int32(award.Position.Reel), Row: proto.Int32(award.Position.Row)}
		} else {
			data.JackpotMoney = proto.Int64(award.WinMoney)
			data.WinMoney = proto.Int64(step.WinMoney - award.WinMoney)
			continue
		}

		data.PayGroups = append(data.PayGroups, payout)
	}

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)

	// 舊ID相容
	{
		pmsg := new(msg.Message)
		pmsg.BClassID = int32(netproto.MessageBClassID_Game)
		pmsg.SClassID = int32(10002)

		_data := *data
		_data.PayGroups = make([]*netproto.Slot_LightningPayout, 0)

		// 把所有Step的Award都包進去
		for _, _step := range logic.LightningGame.Steps {
			for _, award := range _step.Awards {
				payout := new(netproto.Slot_LightningPayout)
				payout.SymbolID = proto.Int32(int32(award.Symbol.GetId()))

				payout.WinMoney = proto.Int64(award.WinMoney)
				payout.Type = proto.String(award.Type.String())
				if award.Type == slot.PayTypeJackpot {
					payout.JackpotType = proto.String(award.JackpotType.String())
				}
				if award.JackpotType != slot.JackpotMajor {
					payout.Position = &netproto.Slot_ReelPosition{Reel: proto.Int32(award.Position.Reel), Row: proto.Int32(award.Position.Row)}
				} else {
					continue
				}

				_data.PayGroups = append(_data.PayGroups, payout)
			}
		}

		pmsg.MsgData = &_data
		logic.SendMessageToUser(user, pmsg)
	}
}

// BuildReelSetting 建立傳給Client的遊戲設置
func (logic *SlotLogicBase) BuildReelSetting(groupID, reelNum, rowNum int32, symbols []*protodata.SlotSymbolDistribution) []*netproto.Slot_ReelSetting {

	var ans []*netproto.Slot_ReelSetting
	for reelNum := int32(1); reelNum <= logic.DefaultReelNum; reelNum++ {
		setting := new(netproto.Slot_ReelSetting)
		setting.ReelID = proto.Int32(int32(reelNum))

		// 只留一組，並且打亂
		sbMap := make(map[string]*netproto.Slot_Symbol)

		for rowNum := int32(1); rowNum <= logic.DefaultRowNum; rowNum++ {
			symbols := logic.GameResultManager.GetSymbolsByReelRow(int32(groupID), reelNum, rowNum, symbols)
			for _, symbol := range symbols {
				if symbol.GetProportion() == 0 {
					continue
				}
				sb := new(netproto.Slot_Symbol)
				sb.SymbolID = proto.Int32(int32(symbol.GetId()))
				sb.SymbolName = proto.String(symbol.GetSymbol())
				sbMap[sb.GetSymbolName()] = sb
			}
		}

		for _, row := range sbMap {
			setting.Symbols = append(setting.Symbols, row)
		}

		ans = append(ans, setting)
	}

	sort.Slice(ans, func(i, j int) bool {
		return ans[i].GetReelID() < ans[j].GetReelID()
	})

	return ans
}

// BuildSymbolPayTable 建立傳給Client的符號賠率設置
func (logic *SlotLogicBase) BuildSymbolPayTable(paytables []*protodata.SlotPayTable, descMap map[string]*slot.SymbolDesc) []*netproto.Slot_SymbolPayTable {

	var ans []*netproto.Slot_SymbolPayTable

	for _, row := range paytables {
		pay := new(netproto.Slot_SymbolPayTable)
		i, _ := strconv.ParseFloat(row.GetPayout(), 64)
		pay.Symbol = &netproto.Slot_Symbol{SymbolID: proto.Int32(0), SymbolName: proto.String(row.GetSymbol())}
		pay.Count = proto.Int32(row.GetCount())
		pay.Type = proto.String(row.GetPayType())

		switch pay.GetType() {
		default:
			pay.Payout = proto.Int32(0)
		case slot.PayTypeChip.String(): // 一般賠錢，回傳賠率
			pay.Payout = proto.Int32(int32(i))
		case slot.PayTypeChipXRate.String(): // 一般賠錢，回傳賠率
			pay.Payout = proto.Int32(int32(i))
		}

		// 若有設定符號文案則寫上去
		if descMap != nil {
			if desc, ok := descMap[row.GetSymbol()]; ok {
				pay.Title = proto.String(desc.Title)
				pay.Desc = proto.String(desc.Description)
			}
		}

		ans = append(ans, pay)
	}

	return ans
}

// BuildPayLines 建立傳給Client的賠線設置
func (logic *SlotLogicBase) BuildPayLines(lines []*protodata.SlotPayLine) []*netproto.Slot_PayLine {

	var ans []*netproto.Slot_PayLine
	// PayLines
	for _, row := range lines {
		payline := slot.NewPayLine(row)

		line := new(netproto.Slot_PayLine)
		line.PayLineID = proto.Int32(int32(row.GetLineId()))

		for _, p := range payline.Positions {
			pos := new(netproto.Slot_ReelPosition)
			pos.Reel = proto.Int32(p.Reel)
			pos.Row = proto.Int32(p.Row)
			line.Positions = append(line.Positions, pos)
		}
		ans = append(ans, line)
	}

	return ans
}

// BuildReelResult 依照GameResult建立傳給Client的開獎結果
func (logic *SlotLogicBase) BuildReelResult(result *slot.GameResultManager, payresults []*slot.PayResult) []*netproto.Slot_ReelResult {

	var results []*netproto.Slot_ReelResult

	reelMap := make(map[slot.Position]*netproto.Slot_ReelResult)
	orgReelMap := make(map[slot.Position]*protodata.SlotSymbolDistribution)

	// 輪軸結果
	result.RowMode()
	for row, reels := range result.GameResult {
		for reel, symbol := range reels {
			if symbol == nil {
				continue
			}
			symbolResult := &netproto.Slot_ReelResult{
				Position: &netproto.Slot_ReelPosition{Reel: proto.Int32(int32(reel + 1)), Row: proto.Int32(int32(row + 1))},
				Symbol:   &netproto.Slot_Symbol{SymbolID: proto.Int32(int32(symbol.GetId())), SymbolName: proto.String(symbol.GetSymbol())},
				Blink:    proto.Bool(false),
			}
			reelMap[slot.Position{Reel: int32(reel + 1), Row: int32(row + 1)}] = symbolResult
			orgReelMap[slot.Position{Reel: int32(reel + 1), Row: int32(row + 1)}] = symbol
		}
	}

	switch logic.PayLineMode {
	case slot.PayLineModeFixed:
		for _, row := range payresults {
			for _, line := range row.SymbolPayLines {

				// var blinkPositions []slot.Position
				// tmp := append([]slot.Position{}, line.PayLine.Positions...)

				// switch slot.LinkType(row.Rule.GetLinkType()) {
				// case slot.LinkTypeLeft:
				// 	sort.Slice(tmp, func(i, j int) bool {
				// 		return tmp[i].Reel < tmp[j].Reel
				// 	})
				// case slot.LinkTypeRight:
				// 	sort.Slice(tmp, func(i, j int) bool {
				// 		return tmp[i].Reel > tmp[j].Reel
				// 	})
				// }

				// for i := int32(0); i < row.PayCount; i++ {
				// 	blinkPositions = append(blinkPositions, tmp[i])
				// }

				// fmt.Println(row.SymbolName, row.Rule.GetLinkType(), line, blinkPositions, line.BlinkPositions)

				for _, pos := range line.BlinkPositions {
					if _, ok := reelMap[pos]; ok {
						reelMap[pos].Blink = proto.Bool(true)
					}
				}
			}
		}
	case slot.PayLineModeAll:
		blinkIDs := make(map[int32]bool)
		for _, row := range payresults {
			for _, id := range row.SymbolIDs {
				blinkIDs[id] = true
			}
		}

		for pos, org := range orgReelMap {
			if _, ok := blinkIDs[int32(org.GetId())]; ok {
				reelMap[pos].Blink = proto.Bool(true)
			}
			if val, ok := result.MultiSymbol[int32(org.GetId())]; ok {
				reelMap[pos].Multi = proto.Int64(val)
			} else {
				reelMap[pos].Multi = proto.Int64(1)
			}
		}

	}

	for _, symbol := range reelMap {
		results = append(results, symbol)
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].GetPosition().GetReel() == results[j].GetPosition().GetReel() {
			return results[i].GetPosition().GetRow() < results[j].GetPosition().GetRow()
		}

		return results[i].GetPosition().GetReel() < results[j].GetPosition().GetReel()
	})

	return results
}

// BuildBackgroundInfo 依照GameResult建立傳給Client的開獎結果
func (logic *SlotLogicBase) BuildBackgroundInfo(result *slot.GameResultManager, payresults []*slot.PayResult) *netproto.Slot_BackgroundInfo {

	if result.BackgroundResult == nil {
		return nil
	}

	info := new(netproto.Slot_BackgroundInfo)

	// 輪軸結果
	for pos, bg := range result.BackgroundResult.Result {
		unit := new(netproto.Slot_BackgroundUnit)
		unit.Position = &netproto.Slot_ReelPosition{Reel: proto.Int32(pos.Reel), Row: proto.Int32(pos.Row)}
		unit.Type = proto.String(string(bg))

		info.Result = append(info.Result, unit)
	}

	sort.Slice(info.Result, func(i, j int) bool {
		if info.Result[i].Position.GetReel() == info.Result[j].Position.GetReel() {
			return info.Result[i].Position.GetRow() < info.Result[j].Position.GetRow()
		}
		return info.Result[i].Position.GetReel() < info.Result[j].Position.GetReel()
	})

	blinkReelMap := make(map[int32]bool)
	// 獎勵結果
	for _, row := range payresults {

		if !row.FromBackground {
			continue
		}

		switch row.GetPayType() {
		case slot.PayTypeChip, slot.PayTypeChipXRate, slot.PayTypeBonus, slot.PayTypeJackpot:
			info.WinMoney = proto.Int64(row.WinMoney)
		}

		for _, pos := range row.BlinkPositions {
			blinkReelMap[pos.Reel] = true
		}

	}

	for reel := range blinkReelMap {
		info.BlinkPositions = append(info.BlinkPositions, &netproto.Slot_ReelPosition{Reel: proto.Int32(reel), Row: proto.Int32(0)})
	}

	sort.Slice(info.BlinkPositions, func(i, j int) bool {
		return info.BlinkPositions[i].GetReel() < info.BlinkPositions[j].GetReel()
	})

	return info
}

// BuildPayGroups 建立傳給Client的中獎結果
func (logic *SlotLogicBase) BuildPayGroups(results []*slot.PayResult) []*netproto.Slot_PayGroup {

	var grps []*netproto.Slot_PayGroup

	switch logic.PayLineMode {
	case slot.PayLineModeFixed:
		// for 指定線
		// 中獎內容，賠付線結果(派彩線)
		for _, row := range results {
			if row.GetPayType() == slot.PayTypeLock {
				continue
			}

			if row.FromBackground {
				continue
			}

			if row.Rule.GetLinkType() == slot.LinkTypeAny.String() {
				if row.GetPayType() != slot.PayTypeChip && row.GetPayType() != slot.PayTypeChipXRate {
					continue
				}
				paygroup := new(netproto.Slot_PayGroup)
				payline := new(netproto.Slot_PayLine)
				payline.PayLineID = proto.Int32(-1)
				_, win := row.CalcWinMoney(logic.GameBetMoney, logic.GameBetRate)
				paygroup.WinMoney = proto.Int64(win / int64(row.WayCount))
				paygroup.PayLines = payline
				paygroup.LinkType = proto.String(row.Rule.GetLinkType())
				paygroup.Count = proto.Int32(row.PayCount)
				paygroup.Way = proto.Int32(row.WayCount)
				paygroup.Multi = proto.Int32(row.Multi)
				paygroup.SymbolName = proto.String(row.SymbolName)
				for _, pos := range row.BlinkPositions {
					_pos := new(netproto.Slot_ReelPosition)
					_pos.Reel = proto.Int32(pos.Reel)
					_pos.Row = proto.Int32(pos.Row)

					paygroup.BlinkPositions = append(paygroup.BlinkPositions, _pos)
				}
				grps = append(grps, paygroup)

				continue
			}

			for _, line := range row.SymbolPayLines {
				paygroup := new(netproto.Slot_PayGroup)
				payline := new(netproto.Slot_PayLine)

				if row.IsWinJackpot() {
					paygroup.WinMoney = proto.Int64(logic.GameWinJackpot)
				} else {
					_, win := row.CalcWinMoney(logic.GameBetMoney, logic.GameBetRate)
					paygroup.WinMoney = proto.Int64(win / int64(row.WayCount))
				}
				payline.PayLineID = proto.Int32(int32(line.PayLine.Origin.GetLineId()))

				for _, orgpos := range line.PayLine.Positions {
					pos := new(netproto.Slot_ReelPosition)
					pos.Reel = proto.Int32(orgpos.Reel)
					pos.Row = proto.Int32(orgpos.Row)
					payline.Positions = append(payline.Positions, pos)
				}

				paygroup.PayLines = payline
				paygroup.LinkType = proto.String(row.Rule.GetLinkType())
				paygroup.Count = proto.Int32(row.PayCount)
				paygroup.Way = proto.Int32(row.WayCount)
				paygroup.SymbolName = proto.String(row.SymbolName)

				for _, pos := range line.BlinkPositions {
					paygroup.BlinkPositions = append(paygroup.BlinkPositions, &netproto.Slot_ReelPosition{Reel: proto.Int32(pos.Reel), Row: proto.Int32(pos.Row)})
				}
				// tmp := paygroup.PayLines.GetPositions()

				// switch slot.LinkType(row.Rule.GetLinkType()) {
				// case slot.LinkTypeLeft:
				// 	sort.Slice(tmp, func(i, j int) bool {
				// 		return tmp[i].GetReel() < tmp[j].GetReel()
				// 	})
				// case slot.LinkTypeRight:
				// 	sort.Slice(tmp, func(i, j int) bool {
				// 		return tmp[i].GetReel() > tmp[j].GetReel()
				// 	})
				// }

				// for i := int32(0); i < row.PayCount; i++ {
				// 	paygroup.BlinkPositions = append(paygroup.BlinkPositions, tmp[i])
				// }

				grps = append(grps, paygroup)
			}
		}

	case slot.PayLineModeAll:
		// for 全線
		for _, row := range results {
			if row.GetPayType() != slot.PayTypeChip && row.GetPayType() != slot.PayTypeChipXRate {
				continue
			}
			if row.FromBackground {
				continue
			}

			base, win := row.CalcWinMoney(logic.GameBetMoney, logic.GameBetRate)

			switch row.MultiType {
			default:

				paygroup := new(netproto.Slot_PayGroup)
				payline := new(netproto.Slot_PayLine)

				paygroup.PayLines = payline
				paygroup.LinkType = proto.String(row.Rule.GetLinkType())
				paygroup.Count = proto.Int32(row.PayCount)
				paygroup.Way = proto.Int32(row.WayCount)
				paygroup.WinMoney = proto.Int64(win)
				paygroup.Multi = proto.Int32(row.Multi)
				paygroup.SymbolName = proto.String(row.SymbolName)

				for _, pos := range row.BlinkPositions {
					paygroup.BlinkPositions = append(paygroup.BlinkPositions, &netproto.Slot_ReelPosition{Reel: proto.Int32(pos.Reel), Row: proto.Int32(pos.Row)})
				}

				grps = append(grps, paygroup)

			case slot.MultiTypeLinkWild:
				// 一般連線
				if row.NormalWayCount > 0 {
					paygroup := new(netproto.Slot_PayGroup)
					payline := new(netproto.Slot_PayLine)

					paygroup.PayLines = payline
					paygroup.LinkType = proto.String(row.Rule.GetLinkType())
					paygroup.Count = proto.Int32(row.PayCount)
					paygroup.Way = proto.Int32(row.NormalWayCount)
					paygroup.WinMoney = proto.Int64(base * int64(row.NormalWayCount))
					paygroup.Multi = proto.Int32(1)
					paygroup.SymbolName = proto.String(row.SymbolName)

					for _, pos := range row.BlinkPositions {
						paygroup.BlinkPositions = append(paygroup.BlinkPositions, &netproto.Slot_ReelPosition{Reel: proto.Int32(pos.Reel), Row: proto.Int32(pos.Row)})
					}

					grps = append(grps, paygroup)
				}

				//  Wild連線
				if row.WildWayCount > 0 {
					paygroup := new(netproto.Slot_PayGroup)
					payline := new(netproto.Slot_PayLine)

					paygroup.PayLines = payline
					paygroup.LinkType = proto.String(row.Rule.GetLinkType())
					paygroup.Count = proto.Int32(row.PayCount)
					paygroup.Way = proto.Int32(row.WildWayCount)
					paygroup.WinMoney = proto.Int64(base * int64(row.WildWayCount) * int64(row.Multi))
					paygroup.Multi = proto.Int32(row.Multi)
					paygroup.SymbolName = proto.String(row.SymbolName)

					for _, pos := range row.BlinkPositions {
						paygroup.BlinkPositions = append(paygroup.BlinkPositions, &netproto.Slot_ReelPosition{Reel: proto.Int32(pos.Reel), Row: proto.Int32(pos.Row)})
					}

					grps = append(grps, paygroup)
				}

			case slot.MultiTypeLinkWildTarget:
				for _, line := range row.SymbolPayLines {

					paygroup := new(netproto.Slot_PayGroup)

					if row.IsWinJackpot() {
						paygroup.WinMoney = proto.Int64(logic.GameWinJackpot)
					} else {
						paygroup.WinMoney = proto.Int64(line.Win)
					}

					paygroup.LinkType = proto.String(row.Rule.GetLinkType())
					paygroup.Count = proto.Int32(row.PayCount)
					paygroup.Way = proto.Int32(row.WayCount)
					paygroup.Multi = proto.Int32(line.GetMultipleValue(logic.MultiMode))
					paygroup.SymbolName = proto.String(row.SymbolName)

					for _, pos := range line.PayLine.Positions {
						paygroup.BlinkPositions = append(paygroup.BlinkPositions, &netproto.Slot_ReelPosition{Reel: proto.Int32(pos.Reel), Row: proto.Int32(pos.Row)})
					}

					grps = append(grps, paygroup)
				}
			}

		}
	}

	return grps
}

// BuildFreeGameInfos 以FreeSpinTable轉成FreeGameInfo
func (logic *SlotLogicBase) BuildFreeGameInfos(tbs []*protodata.FreeSpinTable) []*netproto.Slot_FreeGameInfo {

	infos := make([]*netproto.Slot_FreeGameInfo, 0, len(tbs))

	for _, row := range tbs {
		info := logic.BuildFreeGameInfo(row)
		infos = append(infos, info)
	}

	return infos
}

// BuildFreeGameInfo 以FreeSpinTable轉成FreeGameInfo
func (logic *SlotLogicBase) BuildFreeGameInfo(tb *protodata.FreeSpinTable) *netproto.Slot_FreeGameInfo {
	if tb == nil {
		return nil
	}

	info := new(netproto.Slot_FreeGameInfo)
	info.ID = proto.Int32(int32(tb.GetId()))
	info.Count = proto.Int32(int32(tb.GetFreeSpinCount()))
	info.ReelSetting = logic.BuildReelSetting(int32(tb.GetFreeSymbolGroupId()), int32(tb.GetFreeReelNum()), int32(tb.GetFreeRowNum()), logic.Proto.SlotSymbolDistribution)

	info.ReelNum = proto.Int32(int32(tb.GetFreeReelNum()))
	info.RowNum = proto.Int32(int32(tb.GetFreeRowNum()))
	info.TotalPayLine = proto.Int32(int32(math.Pow(float64(tb.GetFreeRowNum()), float64(tb.GetFreeReelNum()))))

	data := make(map[string]interface{})

	if tb.GetConditionValue() != "" {
		data["Condition"] = tb.GetGoldCondition()
		f, _ := strconv.ParseFloat(tb.GetConditionValue(), 64)
		data["ConditionValue"] = int64(f)
	}

	switch tb.GetMultiplierType() {
	case "wild":
		var key string
		var multi []int32

		strs := strings.Split(tb.GetMultiplierValue(), ",")
		for _, s := range strs {
			_datas := strings.Split(s, ":")

			f, _ := strconv.ParseFloat(_datas[0], 64)
			multi = append(multi, int32(f))
		}
		data["multi"] = multi

		for _, setting := range info.GetReelSetting() {
			for _, symbol := range setting.GetSymbols() {
				for _, attr := range logic.Proto.SymbolAttributes {
					if attr.GetSymbol() == symbol.GetSymbolName() {
						if attr.GetIsExpand() == 1 {
							continue
						}
						key = symbol.GetSymbolName()
						break
					}
				}
			}
		}
		data["key"] = key
	}

	if b, err := json.Marshal(data); err != nil {
		mlog.Warn("FreeGameInfo.Data marshal fail, %v", err)
	} else {
		info.Data = proto.String(string(b))
	}

	return info
}

// BuildGoldBallInfos 建立金球倍率設定
func (logic *SlotLogicBase) BuildGoldBallInfos(attrs []*protodata.SymbolAttributes) []*netproto.Slot_GoldBallInfo {

	var infos []*netproto.Slot_GoldBallInfo

	for _, attr := range attrs {
		if attr.GetIsGoldBall() == 1 {
			info := new(netproto.Slot_GoldBallInfo)
			info.Symbol = proto.String(attr.GetSymbol())

			f, _ := strconv.ParseFloat(attr.GetGoldballValue(), 64)

			switch slot.PayType(attr.GetGoldballType()) {
			case slot.PayTypeJackpot:
				info.Type = proto.String("jackpot")
				info.Value = proto.Int64(int64(f))
			case slot.PayTypeChip:
				info.Type = proto.String("multiple")
				info.Value = proto.Int64(int64(f))
			case slot.PayTypeChipXRate:
				info.Type = proto.String("multiple_rate")
				info.Value = proto.Int64(int64(f))
			case slot.PayTypeAccuMulate:
				info.Type = proto.String("accumulate")
				info.AccumulateSymbols = strings.Split(attr.GetGoldballValue(), ",")
			}

			infos = append(infos, info)
		}
	}

	return infos
}
