package service

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	apiutil "example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	"example.invalid/ServerTeam/Servers/pkg/db"
	dbutil "example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// SelectionService 玩家選擇遊戲相關
type SelectionService interface {
	Read(cer string, userID, gameID, flagID int32) (*netproto.Slot_SelectGameInfo, error)
	Record(cer string, userID, gameID, flagID, patternID, originID, times int32, data string) error
	Reset(cer string, userID, gameID, flagID int32) error
}

// DBSelectionService 玩家選擇遊戲相關
type DBSelectionService struct{}

// Read 讀取 SelectGame
func (s *DBSelectionService) Read(cer string, userID, gameID, flagID int32) (*netproto.Slot_SelectGameInfo, error) {
	const ProcedureName = "PrGs_Slot_SelectGameRead"

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

	data := new(netproto.Slot_SelectGameInfo)
	data.UserID = proto.Int32(userID)
	data.GameID = proto.Int32(gameID)
	data.FlagID = proto.Int32(flagID)
	data.Times = proto.Int32(0)
	data.Data = proto.String("")

	if len(dbret[0].Rows) == 0 {
		return data, nil
	}

	if val, ok := dbret[0].GetValueByColName(0, "UserID").(int64); ok {
		data.UserID = proto.Int32(int32(val))
	} else {
		mlog.Warn("Procedure %s Read UserID error", ProcedureName)
	}

	if val, ok := dbret[0].GetValueByColName(0, "GameID").(int64); ok {
		data.GameID = proto.Int32(int32(val))
	} else {
		mlog.Warn("Procedure %s Read GameID error", ProcedureName)
	}

	if val, ok := dbret[0].GetValueByColName(0, "FlagID").(int64); ok {
		data.FlagID = proto.Int32(int32(val))
	} else {
		mlog.Warn("Procedure %s Read FlagID error", ProcedureName)
	}

	if val, ok := dbret[0].GetValueByColName(0, "PatternID").(int64); ok {
		data.PatternID = proto.Int32(int32(val))
	} else {
		mlog.Warn("Procedure %s Read PatternID error", ProcedureName)
	}

	if val, ok := dbret[0].GetValueByColName(0, "OriginID").(int64); ok {
		data.OriginID = proto.Int32(int32(val))
	} else {
		mlog.Warn("Procedure %s Read OriginID error", ProcedureName)
	}

	if val, ok := dbret[0].GetValueByColName(0, "Times").(int64); ok {
		data.Times = proto.Int32(int32(val))
	} else {
		mlog.Warn("Procedure %s Read Times error", ProcedureName)
	}

	if val, ok := dbret[0].GetValueByColName(0, "Data").(string); ok {
		data.Data = proto.String(val)
	} else {
		mlog.Warn("Procedure %s Read Data error", ProcedureName)
	}

	return data, nil
}

// Record 紀錄 SelectGame
func (s *DBSelectionService) Record(cer string, userID, gameID, flagID, patternID, originID, times int32, data string) error {

	const ProcedureName = "PrGs_Slot_SelectGameRecord"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", userID)       // 使用者ID
	ps.AddIntInput("intGameID", gameID)       // GameID
	ps.AddIntInput("intFlagID", flagID)       // FlagID
	ps.AddIntInput("intPatternID", patternID) // PatternID
	ps.AddIntInput("intOriginID", originID)   // OriginID
	ps.AddIntInput("intTimes", times)         // Times
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

// Reset 清除 SelectGame
func (s *DBSelectionService) Reset(cer string, userID, gameID, flagID int32) error {

	const ProcedureName = "PrGs_Slot_SelectGameReset"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", userID) // 使用者ID
	ps.AddIntInput("intGameID", gameID) // GameID
	ps.AddIntInput("intFlagID", flagID) // FlagID
	ps.AddVarcharOutput("chvMsg", "")   // chvMsg的輸出

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

// APISelectionService 玩家選擇遊戲相關
type APISelectionService struct{}

// Read 讀取 SelectGame
func (s *APISelectionService) Read(cer string, userID, gameID, flagID int32) (*netproto.Slot_SelectGameInfo, error) {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	query := map[string]interface{}{
		"playerID": userID,
		"gameID":   gameID,
		"flagID":   flagID,
	}

	ok, res := apiutil.GameSlotService.GetProb(rpc, query)
	if !ok {
		return nil, fmt.Errorf("API GET SelectGame Error, %d %s", res.Code, res.Message)
	}

	data := new(netproto.Slot_SelectGameInfo)
	data.UserID = proto.Int32(userID)
	data.GameID = proto.Int32(gameID)
	data.FlagID = proto.Int32(flagID)
	data.Times = proto.Int32(0)
	data.Data = proto.String("")

	if len(res.Data.SlotProbRecord) == 0 {
		return data, nil
	}

	data.UserID = proto.Int32(res.Data.SlotProbRecord[0].PlayerID)
	data.GameID = proto.Int32(res.Data.SlotProbRecord[0].GameID)
	data.FlagID = proto.Int32(res.Data.SlotProbRecord[0].FlagID)
	data.PatternID = proto.Int32(res.Data.SlotProbRecord[0].PatternID)
	data.OriginID = proto.Int32(res.Data.SlotProbRecord[0].OriginID)
	data.Times = proto.Int32(res.Data.SlotProbRecord[0].Times)
	data.Data = proto.String(res.Data.SlotProbRecord[0].Data)

	return data, nil
}

// Record 紀錄 SelectGame
func (s *APISelectionService) Record(cer string, userID, gameID, flagID, patternID, originID, times int32, data string) error {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	params := map[string]interface{}{
		"playerID":  userID,
		"gameID":    gameID,
		"flagID":    flagID,
		"patternID": patternID,
		"originID":  originID,
		"times":     times,
		"data":      data,
	}

	ok, _ := apiutil.GameSlotService.PostProb(rpc, params)
	if !ok {
		return fmt.Errorf("API POST SelectGame Error")
	}

	return nil
}

// Reset 清除 SelectGame
func (s *APISelectionService) Reset(cer string, userID, gameID, flagID int32) error {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	query := map[string]interface{}{
		"playerID": userID,
		"gameID":   gameID,
		"flagID":   flagID,
	}

	ok, _ := apiutil.GameSlotService.DeleteProb(rpc, query, nil)
	if !ok {
		return fmt.Errorf("API DELETE SelectGame Error")
	}

	return nil
}
