package service

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	apiutil "example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	"example.invalid/ServerTeam/Servers/pkg/db"
	dbutil "example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// DeferredAwardService 延遲獎勵
type DeferredAwardService interface {
	Read(cer string, userID, gameID, flagID int32, awardType string) ([]*netproto.Slot_DeferredAward, error)
	Record(cer string, userID, gameID, flagID int32, awardType string, patternID, originID int32, data string) error
	Pop(cer string, userID, gameID, flagID int32, awardType string) (*netproto.Slot_DeferredAward, error)
}

// DBDeferredAwardService 延遲獎勵
type DBDeferredAwardService struct{}

// Read 讀取 DeferredAward
func (s *DBDeferredAwardService) Read(cer string, userID, gameID, flagID int32, awardType string) ([]*netproto.Slot_DeferredAward, error) {

	const ProcedureName = "PrGs_Slot_DeferredAwardRead"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", userID)                // UserID
	ps.AddIntInput("intGameID", gameID)                // GameID
	ps.AddIntInput("intFlagID", flagID)                // FlagID
	ps.AddVarcharInput("chvAwardType", awardType, 100) // Award Type
	ps.AddVarcharOutput("chvMsg", "")                  // chvMsg的輸出

	proc := db.NewProcedure(ProcedureName, ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		return nil, fmt.Errorf("Procedure %s execute error. %v; SQL: %s", ProcedureName, err, proc.GetSQLString())
	}

	msg := dbret.GetOutputParamValueString("chvMsg")

	if dbret.GetReturnValue() != 1 {
		return nil, fmt.Errorf("Procedure %s execute error. %v", ProcedureName, msg)
	}

	var result []*netproto.Slot_DeferredAward

	for i := 0; i < len(dbret[0].Rows); i++ {

		data := new(netproto.Slot_DeferredAward)

		if val, ok := dbret[0].GetValueByColName(i, "AwardType").(string); ok {
			data.AwardType = proto.String(val)
		}

		if val, ok := dbret[0].GetValueByColName(i, "PatternID").(int64); ok {
			data.PatternID = proto.Int32(int32(val))
		}

		if val, ok := dbret[0].GetValueByColName(i, "OriginID").(int64); ok {
			data.OriginID = proto.Int32(int32(val))
		}

		if val, ok := dbret[0].GetValueByColName(i, "Data").(string); ok {
			data.Data = proto.String(val)
		}

		result = append(result, data)
	}

	return result, nil
}

// Record 新增獎勵紀錄
func (s *DBDeferredAwardService) Record(cer string, userID, gameID, flagID int32, awardType string, patternID, originID int32, data string) error {

	const ProcedureName = "PrGs_Slot_DeferredAwardRecord"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", userID)                // 玩家ID
	ps.AddIntInput("intGameID", gameID)                // 遊戲ID
	ps.AddIntInput("intFlagID", flagID)                // FlagID
	ps.AddVarcharInput("chvAwardType", awardType, 100) // Award Type
	ps.AddIntInput("intPatternID", patternID)          // PatternID
	ps.AddIntInput("intOriginID", originID)            // OriginID
	ps.AddVarcharInput("txtData", data, 255)           // Data

	ps.AddVarcharOutput("chvMsg", "")

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

// Pop 取用一筆 DeferredAward
func (s *DBDeferredAwardService) Pop(cer string, userID, gameID, flagID int32, awardType string) (*netproto.Slot_DeferredAward, error) {

	const ProcedureName = "PrGs_Slot_DeferredAwardPop"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", userID)                // UserID
	ps.AddIntInput("intGameID", gameID)                // GameID
	ps.AddIntInput("intFlagID", flagID)                // FlagID
	ps.AddVarcharInput("chvAwardType", awardType, 100) // Award Type
	ps.AddVarcharOutput("chvMsg", "")                  // chvMsg的輸出

	proc := db.NewProcedure(ProcedureName, ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		return nil, fmt.Errorf("Procedure %s execute error. %v; SQL: %s", ProcedureName, err, proc.GetSQLString())
	}

	msg := dbret.GetOutputParamValueString("chvMsg")

	if dbret.GetReturnValue() != 1 {
		return nil, fmt.Errorf("Procedure %s execute error. %v", ProcedureName, msg)
	}

	var result []*netproto.Slot_DeferredAward

	for i := 0; i < len(dbret[0].Rows); i++ {

		data := new(netproto.Slot_DeferredAward)

		if val, ok := dbret[0].GetValueByColName(i, "AwardType").(string); ok {
			data.AwardType = proto.String(val)
		}

		if val, ok := dbret[0].GetValueByColName(i, "OriginID").(int64); ok {
			data.OriginID = proto.Int32(int32(val))
		}

		if val, ok := dbret[0].GetValueByColName(i, "Data").(string); ok {
			data.Data = proto.String(val)
		}

		result = append(result, data)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result[0], nil
}

// APIDeferredAwardService 延遲獎勵
type APIDeferredAwardService struct{}

// Read 讀取 DeferredAward
func (s *APIDeferredAwardService) Read(cer string, userID, gameID, flagID int32, awardType string) ([]*netproto.Slot_DeferredAward, error) {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	query := map[string]interface{}{
		"playerID":  userID,
		"gameID":    gameID,
		"flagID":    flagID,
		"awardType": awardType,
	}

	ok, res := apiutil.GameSlotService.GetDeferredAward(rpc, query)
	if !ok {
		return nil, fmt.Errorf("API GET DeferredAward Error, %d %s", res.Code, res.Message)
	}

	var result []*netproto.Slot_DeferredAward

	for _, row := range res.Data.SlotDeferredAward {
		data := new(netproto.Slot_DeferredAward)
		data.AwardType = proto.String(row.AwardType)
		data.PatternID = proto.Int32(row.PatternID)
		data.OriginID = proto.Int32(row.OriginID)
		data.Data = proto.String(row.Data)
		result = append(result, data)
	}

	return result, nil
}

// Record 新增獎勵紀錄
func (s *APIDeferredAwardService) Record(cer string, userID, gameID, flagID int32, awardType string, patternID, originID int32, data string) error {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	params := map[string]interface{}{
		"playerID":  userID,
		"gameID":    gameID,
		"flagID":    flagID,
		"awardType": awardType,
		"patternID": patternID,
		"originID":  originID,
		"data":      data,
	}

	ok := apiutil.GameSlotService.PostDeferredAward(rpc, params)
	if !ok {
		return fmt.Errorf("API POST DeferredAward Error")
	}

	return nil
}

// Pop 取用一筆 DeferredAward
func (s *APIDeferredAwardService) Pop(cer string, userID, gameID, flagID int32, awardType string) (*netproto.Slot_DeferredAward, error) {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	query := map[string]interface{}{
		"playerID":  userID,
		"gameID":    gameID,
		"flagID":    flagID,
		"awardType": awardType,
	}

	ok, res := apiutil.GameSlotService.DeleteDeferredAward(rpc, query)
	if !ok {
		return nil, fmt.Errorf("API GET DeferredAward Error, %d %s", res.Code, res.Message)
	}

	if res.Data.PlayerID == 0 {
		return nil, nil
	}

	data := new(netproto.Slot_DeferredAward)
	data.AwardType = proto.String(res.Data.AwardType)
	data.PatternID = proto.Int32(res.Data.PatternID)
	data.OriginID = proto.Int32(res.Data.OriginID)
	data.Data = proto.String(res.Data.Data)

	return data, nil
}
