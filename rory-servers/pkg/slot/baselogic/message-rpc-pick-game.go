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

// ReadPickGame 讀取PickGame, 使用BonusGame
func (logic *SlotLogicBase) ReadPickGame(user *gamelogic.GameUser) {
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

// RecordPickGame 紀錄PickGame, 使用BonusGame
func (logic *SlotLogicBase) RecordPickGame(user *gamelogic.GameUser, turnIDs ...int32) {

	if logic.PickGame == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_BonusGameRecordID)

	record := new(netproto.Slot_MSG_BonusGameRecord)

	for _, row := range logic.PickGame.Process {
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
		detail.OriginID = proto.Int32(logic.PickGame.Origin.GetId())
		detail.TurnID = proto.Int32(row.ID)
		detail.Location = proto.Int32(row.Location)
		detail.BetMoney = proto.Int32(int32(logic.GameBetMoney))
		detail.BetRate = proto.Int32(int32(logic.GameBetRate))

		data := make(map[string]interface{})
		data["Bonus"] = row.BonusID
		data["Position"] = row.Positon

		data["TotalWin"] = logic.PickGame.WinMoney
		data["AllPosition"] = logic.PickGame.Positions

		b, _ := json.Marshal(data)

		detail.Data = proto.String(string(b))

		record.Detail = append(record.Detail, detail)
	}

	sort.Slice(record.Detail, func(i, j int) bool {
		return record.Detail[i].GetTurnID() < record.Detail[j].GetTurnID()
	})

	pmsg.MsgData = record

	logic.SendRPCMessageFunc(true, pmsg)
}

// ResetPickGame 重置, 使用BonusGame
func (logic *SlotLogicBase) ResetPickGame(user *gamelogic.GameUser) {

	if logic.PickGame == nil {
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
	data.OriginID = proto.Int32(logic.PickGame.Origin.GetId())

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// RecvPickGame 接收 Pick Game 復原
func (logic *SlotLogicBase) RecvPickGame(m *msg.Message) {

	defer logic.CheckWaitGroup()
	mlog.Debug("Table %d, RecvPickGame", logic.TableNo())

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
			for _, r := range Proto.SlotPayTable {
				switch slot.PayType(r.GetPayType()) {
				default:
					continue
				case slot.PayTypePickGame:
					if r.GetId() == originID {
						pg := buildPickGame(r, Proto, datas[originID])

						if pg != nil && !pg.IsFinished() {
							logic.bufFeatureGames[DefaultPattern] = append(logic.bufFeatureGames[DefaultPattern], pg)
						}
					}
				}
			}
		}

		// 查 Release 表
		if ReleaseProto != nil && ReleaseProto.PatternID == details[0].GetPatternID() {
			for _, r := range ReleaseProto.SlotPayTable {
				switch slot.PayType(r.GetPayType()) {
				default:
					continue
				case slot.PayTypePickGame:
					if r.GetId() == originID {
						pg := buildPickGame(r, ReleaseProto, datas[originID])

						if pg != nil && !pg.IsFinished() {
							logic.bufFeatureGames[ReleasePattern] = append(logic.bufFeatureGames[ReleasePattern], pg)
						}
					}
				}
			}
		}

		// 查 Bite 表
		if BiteProto != nil && BiteProto.PatternID == details[0].GetPatternID() {
			for _, r := range BiteProto.SlotPayTable {
				switch slot.PayType(r.GetPayType()) {
				default:
					continue
				case slot.PayTypePickGame:
					if r.GetId() == originID {
						pg := buildPickGame(r, BiteProto, datas[originID])

						if pg != nil && !pg.IsFinished() {
							logic.bufFeatureGames[BitePattern] = append(logic.bufFeatureGames[BitePattern], pg)
						}
					}
				}
			}
		}
	}
}

func buildPickGame(rule *protodata.SlotPayTable, proto *protodata.ProtoData, datas []*netproto.Slot_BonusGameDetail) *slot.PickGame {

	tmp := datas[0]
	var positions []slot.Position
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(tmp.GetData()), &data); err != nil {
		return nil
	}
	if val, ok := data["AllPosition"]; ok {
		if vals, ok := val.([]interface{}); ok {
			for _, _pos := range vals {
				_pos := _pos.(map[string]interface{})

				var p slot.Position

				p.Reel = int32(_pos["Reel"].(float64))
				p.Row = int32(_pos["Row"].(float64))

				positions = append(positions, p)
			}
		}
	}

	game := BuildPickGame(rule, proto, int32(len(positions)), int64(tmp.GetBetMoney()), int64(tmp.GetBetRate()))
	game.Positions = positions

	// 重置Process
	game.Process = make(map[int32]*slot.PickGameDetail)
	for _, detail := range datas {

		data := make(map[string]interface{})

		if err := json.Unmarshal([]byte(detail.GetData()), &data); err != nil {
			continue
		}

		bonusID := int32((data["Bonus"].(float64)))

		win := int64(data["TotalWin"].(float64))

		var pos slot.Position

		if val, ok := data["Position"]; ok {
			_pos := val.(map[string]interface{})
			pos.Reel = int32(_pos["Reel"].(float64))
			pos.Row = int32(_pos["Row"].(float64))
		}

		game.WinMoney = win

		var bonus *protodata.BonusTable
		for _, row := range proto.BonusTable {
			if row.GetId() == bonusID {
				bonus = row
				break
			}
		}

		if bonus == nil {
			mlog.Warn("[PickGame] Recover錯誤, 查無對應Bonus(%d)", bonusID)
			continue
		}

		ps := new(slot.PickGameDetail)
		ps.Origin = bonus
		ps.BonusID = bonus.GetId()
		ps.ID = detail.GetTurnID()
		ps.Location = detail.GetLocation()
		ps.Positon = pos

		// 回推金額
		switch slot.PayType(bonus.GetRewardType()) {
		case slot.PayTypeChip:
			f, _ := strconv.ParseFloat(bonus.GetRewardValue(), 64)
			ps.Money = int64(f) * game.GameBetMoney
		case slot.PayTypeChipXRate:
			f, _ := strconv.ParseFloat(bonus.GetRewardValue(), 64)
			ps.Money = int64(f) * game.GameBetMoney * game.GameBetRate
		}

		game.Process[detail.GetTurnID()] = ps

		game.All[int(ps.ID-1)] = ps.Money
	}

	if game.IsFinished() {
		return nil
	}

	return game
}

// BuildPickGame 建立PickGame
func BuildPickGame(rule *protodata.SlotPayTable, proto *protodata.ProtoData, size int32, betMoney, betRate int64) *slot.PickGame {
	f, err := strconv.ParseFloat(rule.GetPayout(), 64)
	if err != nil {
		mlog.Debug("Parse [%s] error %v", rule.GetPayType(), err)
		return nil
	}

	var picktb *protodata.PickGameTable
	for _, row := range proto.PickGameTable {
		if int32(row.GetId()) == int32(f) {
			picktb = row
			break
		}
	}

	if picktb == nil {
		mlog.Debug("No PickGame for [%s], ID: %d", rule.GetPayType(), int32(f))
		return nil
	}

	var bonuies []*protodata.BonusTable
	for _, row := range proto.BonusTable {
		if row.GetGroupId() != int32(picktb.GetRewardValue()) {
			continue
		}
		bonuies = append(bonuies, row)
	}

	if len(bonuies) == 0 {
		mlog.Debug("No Bonus for PickGame(%d), GroupID: %d", picktb.GetId(), picktb.GetRewardValue())
		return nil
	}

	game := slot.NewPickGame(bonuies, size, int32(picktb.GetPickCount()), betMoney, betRate)
	game.Origin = rule

	return game
}
