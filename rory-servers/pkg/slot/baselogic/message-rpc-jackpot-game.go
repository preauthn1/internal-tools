package baselogic

import (
	"encoding/json"
	"sort"
	"strconv"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// ReadJackpotGame 讀取JackpotGame, 使用BonusGame
func (logic *SlotLogicBase) ReadJackpotGame(user *gamelogic.GameUser) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadID)

	data := new(netproto.Slot_MSG_BonusGameRead)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.ParseInt(logic.TableManager.Room.Flag, 10, 32)
	data.FlagID = proto.Int32(int32(FlagID))

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// RecordJackpotGame 紀錄JackpotGame, 使用BonusGame
func (logic *SlotLogicBase) RecordJackpotGame(user *gamelogic.GameUser, turnIDs ...int32) {

	if logic.JackpotGame == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_BonusGameRecordID)

	record := new(netproto.Slot_MSG_BonusGameRecord)

	for _, row := range logic.JackpotGame.Process {
		if len(turnIDs) > 0 {
			var ok bool
			for _, turn := range turnIDs {
				if turn == row.ID {
					ok = true
				}
			}
			if !ok {
				continue
			}
		}
		detail := new(netproto.Slot_BonusGameDetail)
		detail.UserID = proto.Int32(user.UserID)
		detail.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
		FlagID, _ := strconv.ParseInt(logic.TableManager.Room.Flag, 10, 32)
		detail.FlagID = proto.Int32(int32(FlagID))

		detail.PatternID = proto.Int32(logic.Proto.PatternID)
		detail.OriginID = proto.Int32(logic.JackpotGame.Origin.GetId())
		detail.TurnID = proto.Int32(row.ID)
		detail.Location = proto.Int32(row.Place)
		detail.BetMoney = proto.Int32(int32(logic.GameBetMoney))
		detail.BetRate = proto.Int32(int32(logic.GameBetRate))

		data := make(map[string]interface{})
		data["Level"] = slot.NewJackpotLevel(row.Icon).String()
		data["Win"] = logic.JackpotGame.WinMoney

		b, _ := json.Marshal(data)

		// detail.Data = proto.String(slot.NewJackpotLevel(row.Icon).String())
		detail.Data = proto.String(string(b))

		record.Detail = append(record.Detail, detail)
	}

	sort.Slice(record.Detail, func(i, j int) bool {
		return record.Detail[i].GetTurnID() < record.Detail[j].GetTurnID()
	})

	pmsg.MsgData = record

	logic.SendRPCMessageFunc(true, pmsg)
}

// RestJackpotGame 重置, 使用BonusGame
func (logic *SlotLogicBase) RestJackpotGame(user *gamelogic.GameUser) {

	if logic.JackpotGame == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_BonusGameResetID)

	data := new(netproto.Slot_MSG_BonusGameRead)

	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.ParseInt(logic.TableManager.Room.Flag, 10, 32)
	data.FlagID = proto.Int32(int32(FlagID))
	data.OriginID = proto.Int32(logic.JackpotGame.Origin.GetId())

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// RecvJackpotGame 接收 Jackpot Game 復原
func (logic *SlotLogicBase) RecvJackpotGame(m *msg.Message) {

	defer logic.CheckWaitGroup()
	mlog.Debug("Table %d, RecvJackpotGame", logic.TableNo())

	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_BonusGameRecord{}
		return mm
	})

	rec := rmm.(*netproto.Slot_MSG_BonusGameRecord)

	if logic.GameUser == nil {
		return
	}

	if len(rec.Detail) == 0 {
		return
	}

	datas := make(map[int32][]*netproto.Slot_BonusGameDetail)

	for _, detail := range rec.Detail {
		if datas[detail.GetOriginID()] == nil {
			datas[detail.GetOriginID()] = make([]*netproto.Slot_BonusGameDetail, 0)
		}

		datas[detail.GetOriginID()] = append(datas[detail.GetOriginID()], detail)
	}

	for originID, details := range datas {
		sort.Slice(datas[originID], func(i, j int) bool {
			return datas[originID][i].GetTurnID() < datas[originID][j].GetTurnID()
		})

		if len(details) == 0 {
			continue
		}

		// 查 Default 表
		if Proto != nil && (Proto.PatternID == details[0].GetPatternID() || details[0].GetPatternID() == 0) {
			jackpotMap := make(map[slot.JackpotLevel]*slot.Jackpot)

			for _, row := range logic.Proto.JackpotTable {
				jackpot := new(slot.Jackpot)
				jackpot.Origin = row
				if len(row.GetLevel()) == 0 {
					jackpot.Level = slot.JackpotSingle
				} else {
					jackpot.Level = slot.JackpotLevel(row.GetLevel())
				}
				jackpotMap[jackpot.Level] = jackpot
			}

			for _, r := range Proto.SlotPayTable {
				switch slot.PayType(r.GetPayType()) {
				default:
					continue
				case slot.PayTypeMysteryJPG, slot.PayTypeMysteryJPGR, slot.PayTypeMysteryJackpot, slot.PayTypeMysteryJackpotR:
					if r.GetId() == originID {
						jpg := buildJackpotGame(jackpotMap, r, details, logic.MaxBetCount, logic.MaxBetRate)

						if jpg != nil && !jpg.IsFinished() {
							logic.bufFeatureGames[DefaultPattern] = append(logic.bufFeatureGames[DefaultPattern], jpg)
						}
					}
				}
			}
		}

		// 查 Release 表
		if ReleaseProto != nil && ReleaseProto.PatternID == details[0].GetPatternID() {
			jackpotMap := make(map[slot.JackpotLevel]*slot.Jackpot)

			for _, row := range logic.Proto.JackpotTable {
				jackpot := new(slot.Jackpot)
				jackpot.Origin = row
				if len(row.GetLevel()) == 0 {
					jackpot.Level = slot.JackpotSingle
				} else {
					jackpot.Level = slot.JackpotLevel(row.GetLevel())
				}
				jackpotMap[jackpot.Level] = jackpot
			}

			for _, r := range ReleaseProto.SlotPayTable {
				switch slot.PayType(r.GetPayType()) {
				default:
					continue
				case slot.PayTypeMysteryJPG, slot.PayTypeMysteryJPGR, slot.PayTypeMysteryJackpot, slot.PayTypeMysteryJackpotR:
					if r.GetId() == originID {
						jpg := buildJackpotGame(jackpotMap, r, details, logic.MaxBetCount, logic.MaxBetRate)

						if jpg != nil && !jpg.IsFinished() {
							logic.bufFeatureGames[ReleasePattern] = append(logic.bufFeatureGames[ReleasePattern], jpg)
						}
					}
				}
			}
		}

		// 查 Bite 表
		if BiteProto != nil && BiteProto.PatternID == details[0].GetPatternID() {
			jackpotMap := make(map[slot.JackpotLevel]*slot.Jackpot)

			for _, row := range logic.Proto.JackpotTable {
				jackpot := new(slot.Jackpot)
				jackpot.Origin = row
				if len(row.GetLevel()) == 0 {
					jackpot.Level = slot.JackpotSingle
				} else {
					jackpot.Level = slot.JackpotLevel(row.GetLevel())
				}
				jackpotMap[jackpot.Level] = jackpot
			}

			for _, r := range BiteProto.SlotPayTable {
				switch slot.PayType(r.GetPayType()) {
				default:
					continue
				case slot.PayTypeMysteryJPG, slot.PayTypeMysteryJPGR, slot.PayTypeMysteryJackpot, slot.PayTypeMysteryJackpotR:
					if r.GetId() == originID {
						jpg := buildJackpotGame(jackpotMap, r, details, logic.MaxBetCount, logic.MaxBetRate)

						if jpg != nil && !jpg.IsFinished() {
							logic.bufFeatureGames[BitePattern] = append(logic.bufFeatureGames[BitePattern], jpg)
						}
					}
				}
			}
		}

	}

}

func buildJackpotGame(jackpotMap map[slot.JackpotLevel]*slot.Jackpot, rule *protodata.SlotPayTable, datas []*netproto.Slot_BonusGameDetail, maxBetMoeny, maxBetRate int32) *slot.JackpotGame {

	tmp := datas[0]
	jackpots := make([]*slot.Jackpot, 0, len(jackpotMap))
	for _, jp := range jackpotMap {
		jackpots = append(jackpots, jp)
	}

	jpg := slot.NewJackpotGame(jackpots, maxBetRate, maxBetMoeny, tmp.GetBetRate(), int64(tmp.GetBetMoney()))
	jpg.Origin = rule

	jpg.Process = make(map[int32]*slot.JackpotGameDetail)

	jpMap := make(map[slot.JackpotLevel]int)

	for _, detail := range datas {

		data := make(map[string]interface{})

		if err := json.Unmarshal([]byte(detail.GetData()), &data); err != nil {
			continue
		}

		level := slot.JackpotLevel(data["Level"].(string))
		win := int64(data["Win"].(float64))

		jpMap[level]++
		jpg.WinMoney = win

		ps := new(slot.JackpotGameDetail)
		ps.ID = detail.GetTurnID()
		ps.Place = detail.GetLocation()
		ps.Icon = level.Int()

		jpg.Process[detail.GetTurnID()] = ps
	}

	for level, cnt := range jpMap {
		if cnt == 3 {
			jpg.WinJackpot = jackpotMap[level]
		}
	}

	return jpg
}
