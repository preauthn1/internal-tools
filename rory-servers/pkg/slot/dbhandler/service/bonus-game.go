package service

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	apiutil "example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	"example.invalid/ServerTeam/Servers/pkg/db"
	dbutil "example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// BonusGameService 玩家紀錄相關
type BonusGameService interface {
	Read(cer string, userID, gameID, flagID int32) ([]*netproto.Slot_BonusGameDetail, error)
	Record(cer string, userID, gameID, flagID, patternID, originID, turnID, location, betMoney, betRate int32, data string) error
	Reset(cer string, serID, gameID, flagID, originID int32) error
}

// DBBonusGameService 玩家紀錄相關
type DBBonusGameService struct{}

// Read 讀取BonusGame
func (s *DBBonusGameService) Read(cer string, userID, gameID, flagID int32) ([]*netproto.Slot_BonusGameDetail, error) {

	const ProcedureName = "PrGs_Slot_BonusGameRead"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", userID) // UserID
	ps.AddIntInput("intGameID", gameID) // GameID
	ps.AddIntInput("intFlagID", flagID) // FlagID
	ps.AddVarcharOutput("chvMsg", "")   // chvMsg的輸出

	proc := db.NewProcedure(ProcedureName, ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		return nil, fmt.Errorf("Procedure %s execute error. %v; SQL: %s", ProcedureName, err, proc.GetSQLString())
	}

	msg := dbret.GetOutputParamValueString("chvMsg")

	if dbret.GetReturnValue() != 1 {
		return nil, fmt.Errorf("Procedure %s execute error. %v", ProcedureName, msg)
	}

	var result []*netproto.Slot_BonusGameDetail

	for i := 0; i < len(dbret[0].Rows); i++ {

		data := new(netproto.Slot_BonusGameDetail)

		if val, ok := dbret[0].GetValueByColName(i, "UserID").(int64); ok {
			data.UserID = proto.Int32(int32(val))
		} else {
			mlog.Warn("Procedure %s Read UserID error", ProcedureName)
		}

		if val, ok := dbret[0].GetValueByColName(i, "GameID").(int64); ok {
			data.GameID = proto.Int32(int32(val))
		} else {
			mlog.Warn("Procedure %s Read GameID error", ProcedureName)
		}

		if val, ok := dbret[0].GetValueByColName(i, "FlagID").(int64); ok {
			data.FlagID = proto.Int32(int32(val))
		} else {
			mlog.Warn("Procedure %s Read FlagID error", ProcedureName)
		}

		if val, ok := dbret[0].GetValueByColName(i, "PatternID").(int64); ok {
			data.PatternID = proto.Int32(int32(val))
		} else {
			mlog.Warn("Procedure %s Read PatternID error", ProcedureName)
		}

		if val, ok := dbret[0].GetValueByColName(i, "OriginID").(int64); ok {
			data.OriginID = proto.Int32(int32(val))
		} else {
			mlog.Warn("Procedure %s Read OriginID error", ProcedureName)
		}

		if val, ok := dbret[0].GetValueByColName(i, "TurnID").(int64); ok {
			data.TurnID = proto.Int32(int32(val))
		} else {
			mlog.Warn("Procedure %s Read TurnID error", ProcedureName)
		}

		if val, ok := dbret[0].GetValueByColName(i, "Location").(int64); ok {
			data.Location = proto.Int32(int32(val))
		} else {
			mlog.Warn("Procedure %s Read Location error", ProcedureName)
		}

		if val, ok := dbret[0].GetValueByColName(i, "BetMoney").(int64); ok {
			data.BetMoney = proto.Int32(int32(val))
		} else {
			mlog.Warn("Procedure %s Read BetMoney error", ProcedureName)
		}

		if val, ok := dbret[0].GetValueByColName(i, "BetRate").(int64); ok {
			data.BetRate = proto.Int32(int32(val))
		} else {
			mlog.Warn("Procedure %s Read BetRate error", ProcedureName)
		}

		if val, ok := dbret[0].GetValueByColName(i, "Data").(string); ok {
			data.Data = proto.String(val)
		} else {
			mlog.Warn("Procedure %s Read Data error", ProcedureName)
		}

		result = append(result, data)
	}

	return result, nil
}

// Record 紀錄BonusGame
func (s *DBBonusGameService) Record(cer string, userID, gameID, flagID, patternID, originID, turnID, location, betMoney, betRate int32, data string) error {
	const ProcedureName = "PrGs_Slot_BonusGameRecord"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", userID)       // UserID
	ps.AddIntInput("intGameID", gameID)       // GameID
	ps.AddIntInput("intFlagID", flagID)       // FlagID
	ps.AddIntInput("intPatternID", patternID) // PatternID
	ps.AddIntInput("intOriginID", originID)   // OriginID
	ps.AddIntInput("intTurnID", turnID)       // TurnID
	ps.AddIntInput("intLocation", location)   // Location
	ps.AddIntInput("intBetMoney", betMoney)   // BetMoney
	ps.AddIntInput("intBetRate", betRate)     // BetRate
	ps.AddVarcharInput("txtData", data, 255)  // Data
	ps.AddVarcharOutput("chvMsg", "")         // chvMsg的輸出

	proc := db.NewProcedure(ProcedureName, ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		return fmt.Errorf("Procedure %s execute error. %v; SQL: %s", ProcedureName, err, proc.GetSQLString())
	}

	msg := dbret.GetOutputParamValueString("chvMsg")

	if dbret.GetReturnValue() != 1 {
		return fmt.Errorf("Procedure %s execute error. %v", ProcedureName, msg)
	}

	return nil
}

// Reset 重置BonusGame
func (s *DBBonusGameService) Reset(cer string, userID, gameID, flagID, originID int32) error {
	const ProcedureName = "PrGs_Slot_BonusGameReset"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", userID)     // UserID
	ps.AddIntInput("intGameID", gameID)     // GameID
	ps.AddIntInput("intFlagID", flagID)     // FlagID
	ps.AddIntInput("intOriginID", originID) // OriginID
	ps.AddVarcharOutput("chvMsg", "")       // chvMsg的輸出

	proc := db.NewProcedure(ProcedureName, ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		return fmt.Errorf("Procedure %s execute error. %v; SQL: %s", ProcedureName, err, proc.GetSQLString())
	}

	msg := dbret.GetOutputParamValueString("chvMsg")

	if dbret.GetReturnValue() != 1 {
		return fmt.Errorf("Procedure %s execute error. %v", ProcedureName, msg)
	}

	return nil
}

// APIBonusGameService 玩家紀錄相關
type APIBonusGameService struct{}

// Read 讀取BonusGame
func (s *APIBonusGameService) Read(cer string, userID, gameID, flagID int32) ([]*netproto.Slot_BonusGameDetail, error) {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	query := map[string]interface{}{
		"playerID": userID,
		"gameID":   gameID,
		"flagID":   flagID,
	}

	ok, res := apiutil.GameSlotService.GetBonus(rpc, query)
	if !ok {
		return nil, fmt.Errorf("API Get BonusGame Error, %d %s", res.Code, res.Message)
	}

	var records []*netproto.Slot_BonusGameDetail

	for _, v := range res.Data.SlotBonusRecord {
		record := new(netproto.Slot_BonusGameDetail)
		record.UserID = proto.Int32(v.PlayerID)
		record.GameID = proto.Int32(v.GameID)
		record.FlagID = proto.Int32(v.FlagID)
		record.PatternID = proto.Int32(v.PatternID)
		record.OriginID = proto.Int32(v.OriginID)
		record.TurnID = proto.Int32(v.TurnID)
		record.Location = proto.Int32(v.Location)
		record.BetMoney = proto.Int32(int32(util.MoneyToCent(v.BetAmount)))
		record.BetRate = proto.Int32(v.BetRate)
		record.Data = proto.String(v.Data)

		records = append(records, record)
	}

	return records, nil
}

// Record 紀錄BonusGame
func (s *APIBonusGameService) Record(cer string, userID, gameID, flagID, patternID, originID, turnID, location, betMoney, betRate int32, data string) error {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	params := make([]map[string]interface{}, 0)

	params = append(params, map[string]interface{}{
		"playerID":  userID,
		"gameID":    gameID,
		"flagID":    flagID,
		"patternID": patternID,
		"originID":  originID,
		"turnID":    turnID,
		"location":  location,
		"betAmount": util.CentToMoney(int64(betMoney)),
		"betRate":   betRate,
		"data":      data,
	})

	if !apiutil.GameSlotService.PostBonus(rpc, params) {
		return fmt.Errorf("API POST BonusGame Error")
	}

	return nil
}

// Reset 重置BonusGame
func (s *APIBonusGameService) Reset(cer string, userID, gameID, flagID, originID int32) error {
	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	query := map[string]interface{}{
		"playerID": userID,
		"gameID":   gameID,
		"flagID":   flagID,
		"originID": originID,
	}

	if !apiutil.GameSlotService.DeleteBonus(rpc, query) {
		return fmt.Errorf("API POST BonusGame Error")
	}

	return nil
}
