package gen

import (
	"sort"

	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

func (game *Game) mError(errType netproto.GEN_ErrorType, err error, src netproto.GEN_GameMessageClassID) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_ErrorID)

	data := new(netproto.GEN_Error)
	data.Type = errType
	data.Message = err.Error()
	data.Source = src

	m.MsgData = data

	return m
}

func (game *Game) mGameSetting() *msg.Message {

	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_GameSettingID)

	data := new(netproto.GEN_GameSetting)
	data.TableNo = game.TableManager.TableNo
	data.BetOdds = game.Bet
	data.ReelStrips = game.mDataReelStrips()

	m.MsgData = data

	return m
}

func (game *Game) mRoadMap() *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_RoadMapID)
	m.MsgData = game.mDataRoadMap(game.PachiSlotGame.RoadMap)
	return m
}

func (game *Game) mGameStatus() *msg.Message {

	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID)

	data := new(netproto.GEN_GameStatus)
	data.Scenario = game.transScenario(game.PachiSlotGame.Scenario)
	data.Mode = game.transMode(game.PachiSlotGame.Mode)
	data.Bonus = game.transScenario(game.PachiSlotGame.Bonus)
	data.RT = int32(game.PachiSlotGame.RT)
	data.Replay = game.PachiSlotGame.Replay
	data.BonusSure = game.PachiSlotGame.BonusSure
	data.RequireCoin = int32(game.PachiSlotGame.BetCoin())
	data.BonusLeft = int32(game.PachiSlotGame.BonusTimesLeft)
	data.JACLeft = int32(game.PachiSlotGame.JACTimesLeft)
	data.BonusStore = int32(game.PachiSlotGame.BonusStore)
	data.ChanceBuff = game.ChanceBuff
	data.ChanceBuffDeadLine = int32(game.ChanceDeadLine)

	m.MsgData = data

	return m
}

func (game *Game) mGameResult(s pb.Scenario, bet, win int64, result pb.ResultType, rsm pachislot.ReelStripMap, sg pachislot.SymbolGroup, payouts []pachislot.Payout, releaseRT bool, blinkPos []*pb.Position) *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_GameResultID)

	data := new(netproto.GEN_GameResult)
	data.Scenario = game.transScenario(s)
	data.BetMoney = bet
	data.WinMoney = win
	data.Result = game.transResultType(result)
	idx1, idx2, idx3 := rsm.StopIndex()
	data.StopIndex = []int32{int32(idx1), int32(idx2), int32(idx3)}
	for i := 0; i < 3; i++ {
		strip := new(netproto.GEN_ReelStrip)
		strip.Reel = int32(i) + 1
		for _, s := range sg[i] {
			strip.Symbols = append(strip.Symbols, game.transSymbol(s))
		}
		data.ReelStrip = append(data.ReelStrip, strip)
	}
	data.ReleaseRT = releaseRT
	for _, p := range payouts {
		data.Payouts = append(data.Payouts, game.mDataPayout(p))
	}
	for _, pos := range blinkPos {
		data.BlinkPositions = append(data.BlinkPositions, game.transPosition(pos))
	}

	m.MsgData = data

	return m
}

func (game *Game) mAppend1GBB0() *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_EventBB0JACID)
	data := new(netproto.GEN_Payout)
	data.Type = game.transPayType(pb.PayType_PT1G)
	data.Amount = 3
	m.MsgData = data
	return m
}

func (game *Game) mAppend1GBB1() *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_EventBB1JACID)
	data := new(netproto.GEN_Payout)
	data.Type = game.transPayType(pb.PayType_PT1G)
	data.Amount = 2
	m.MsgData = data
	return m
}

func (game *Game) mAppend1GBB2() *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_EventBB2JACID)
	data := new(netproto.GEN_Payout)
	data.Type = game.transPayType(pb.PayType_PT1G)
	data.Amount = 1
	m.MsgData = data
	return m
}

func (game *Game) mAppend1GJAC8Bell() *msg.Message {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_EventJAC8BellID)
	data := new(netproto.GEN_Payout)
	data.Type = game.transPayType(pb.PayType_PT1G)
	data.Amount = 1
	m.MsgData = data
	return m
}

func (game *Game) mDataPayLine() []*netproto.GEN_PayLine {
	var lines []*netproto.GEN_PayLine

	for _, _line := range game.PachiSlotGame.Proto.PayLine {
		line := new(netproto.GEN_PayLine)
		line.ID = _line.LineID

		for _, _pos := range _line.Positions {
			line.Position = append(line.Position, game.transPosition(_pos))
		}

		lines = append(lines, line)
	}

	return lines
}

func (game *Game) mDataReelStrips() []*netproto.GEN_ReelStrip {
	var strips []*netproto.GEN_ReelStrip
	for reel, symbols := range game.ReelStrips {
		strip := new(netproto.GEN_ReelStrip)
		strip.Reel = int32(reel)
		for _, s := range symbols {
			strip.Symbols = append(strip.Symbols, game.transSymbol(s))
		}
		strips = append(strips, strip)
	}
	return strips
}

func (game *Game) mDataRoadMap(rm *pachislot.RoadMap) *netproto.GEN_RoadMap {
	data := new(netproto.GEN_RoadMap)
	data.Total = int64(game.PachiSlotGame.RoadMap.TotalCnt)
	data.Spin = int64(game.PachiSlotGame.RoadMap.SpinCnt)
	data.Bonus = game.mDataRoadMapBonus(game.PachiSlotGame.RoadMap)
	data.Snapshots = game.mDataRoadMapSnapshots(game.PachiSlotGame.RoadMap)
	return data
}

func (game *Game) mDataRoadMapBonus(rm *pachislot.RoadMap) []*netproto.GEN_RoadMapBonus {
	var data []*netproto.GEN_RoadMapBonus
	for s, cnt := range rm.FeaturesCnt {
		row := new(netproto.GEN_RoadMapBonus)
		row.Type = game.transScenario(s)
		row.Count = int64(cnt)
		data = append(data, row)
	}
	return data
}

func (game *Game) mDataRoadMapSnapshots(rm *pachislot.RoadMap) []*netproto.GEN_RoadMapSnapshot {
	var data []*netproto.GEN_RoadMapSnapshot
	for _, s := range rm.Snapshots {
		if s == nil {
			continue
		}
		row := new(netproto.GEN_RoadMapSnapshot)
		row.Type = game.transScenario(s.Feature)
		row.Count = int64(s.Count)
		row.Time = s.Time.Unix()
		data = append(data, row)
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Time < data[j].Time
	})
	return data
}

func (game *Game) mDataPayout(p pachislot.Payout) *netproto.GEN_Payout {
	data := new(netproto.GEN_Payout)
	data.Type = game.transPayType(p.Type())
	switch p.Type() {
	case pb.PayType_PTChip:
		data.Amount = p.(*pachislot.PayoutChip).Odds()
	case pb.PayType_PTReplay:
		data.Amount = 1
	case pb.PayType_PTBonus:
		data.Amount = 1
		data.Bonus = game.transScenario(p.(*pachislot.PayoutBonus).Module)
	case pb.PayType_PT1G:
		data.Amount = int64(p.(*pachislot.Payout1G).Times)
	}
	return data
}

func (game *Game) transSymbol(s pb.Symbol) netproto.GEN_Symbol {
	return netproto.GEN_Symbol(s)
}

func (game *Game) transScenario(s pb.Scenario) netproto.GEN_Scenario {
	return netproto.GEN_Scenario(s)
}

func (game *Game) transMode(s pb.ModeType) netproto.GEN_Mode {
	return netproto.GEN_Mode(s)
}

func (game *Game) transResultType(rt pb.ResultType) netproto.GEN_ResultType {
	return netproto.GEN_ResultType(rt)
}

func (game *Game) transPayType(s pb.PayType) netproto.GEN_PayType {
	return netproto.GEN_PayType(s)
}

func (game *Game) transPosition(pos *pb.Position) *netproto.GEN_Position {
	return &netproto.GEN_Position{Reel: pos.Reel, Row: pos.Row}
}
