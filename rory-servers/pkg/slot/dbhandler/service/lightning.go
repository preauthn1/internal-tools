package service

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"
	apiutil "example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	"example.invalid/ServerTeam/Servers/pkg/db"
	dbutil "example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// LightningService Lightning相關功能
type LightningService interface {
	Read(cer string, userID, gameID, flagID int32) (int32, string, error)
	Record(cer string, userID, gameID, flagID, patternID int32, data string) error
}

// DBLightningService Lightning相關功能
type DBLightningService struct{}

// Read 查詢Lightning資料
func (s *DBLightningService) Read(cer string, userID, gameID, flagID int32) (int32, string, error) {

	const ProcedureName = "PrGs_Slot_LightningGameRead"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", userID) // 玩家ID
	ps.AddIntInput("intGameID", gameID) // 遊戲ID
	ps.AddIntInput("intFlagID", flagID) // FlagID

	ps.AddVarcharOutput("chvMsg", "")

	proc := db.NewProcedure(ProcedureName, ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		return 0, "", fmt.Errorf("Procedure %s execute error. %v; SQL: %s", ProcedureName, err, proc.GetSQLString())
	}

	msg := dbret.GetOutputParamValueString("chvMsg")

	if dbret.GetReturnValue() != 1 {
		return 0, "", fmt.Errorf("Procedure %s execute error. %v", ProcedureName, msg)
	}

	var data string
	var patternID int32

	for i := 0; i < len(dbret[0].Rows); i++ {
		if val, ok := dbret[0].GetValueByColName(i, "Data").(string); ok {
			data = val
		}
		if val, ok := dbret[0].GetValueByColName(i, "PatternID").(int64); ok {
			patternID = int32(val)
		}
	}

	return patternID, data, nil
}

// Record 更新Lightning資料
func (s *DBLightningService) Record(cer string, userID, gameID, flagID, patternID int32, data string) error {

	const ProcedureName = "PrGs_Slot_LightningGameRecord"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", userID)          // 玩家ID
	ps.AddIntInput("intGameID", gameID)          // 遊戲ID
	ps.AddIntInput("intFlagID", flagID)          // FlagID
	ps.AddIntInput("intPatternID", patternID)    // PatternID
	ps.AddVarcharInput("txtData", data, 1<<16-1) // Data

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

// APILightningService Lightning相關功能
type APILightningService struct{}

// Read 查詢Lightning資料
func (s *APILightningService) Read(cer string, userID, gameID, flagID int32) (int32, string, error) {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	query := map[string]interface{}{
		"playerID": userID,
		"gameID":   gameID,
		"flagID":   flagID,
	}

	ok, res := apiutil.GameSlotService.GetLightning(rpc, query)
	if !ok {
		return 0, "", fmt.Errorf("API GET Lightning Error, %d %s", res.Code, res.Message)
	}

	if len(res.Data.SlotLightningRecord) == 0 {
		return 0, "", nil
	}

	return res.Data.SlotLightningRecord[0].PatternID, res.Data.SlotLightningRecord[0].Data, nil
}

// Record 更新Lightning資料
func (s *APILightningService) Record(cer string, userID, gameID, flagID, patternID int32, data string) error {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	params := map[string]interface{}{
		"playerID":  userID,
		"gameID":    gameID,
		"flagID":    flagID,
		"patternID": patternID,
		"data":      data,
	}

	ok := apiutil.GameSlotService.PostLightning(rpc, params)
	if !ok {
		return fmt.Errorf("API POST Lightning Error")
	}

	return nil
}
