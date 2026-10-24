package service

import (
	"fmt"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/model"
	apiutil "example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	"example.invalid/ServerTeam/Servers/pkg/db"
	dbutil "example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// UserAssignedModeService 玩家被指定遊戲模式
type UserAssignedModeService interface {
	Read(cer string, userID, gameID, flagID int32) (*netproto.Slot_UserAssignedMode, error)
	Record(cer string, userID, gameID, flagID, mode int32, expire string) error
	Delete(cer string, userID, gameID, flagID int32) error
}

// DBUserAssignedModeService 玩家被指定遊戲模式
type DBUserAssignedModeService struct{}

// Read 讀取 UserAssignedMode
func (s *DBUserAssignedModeService) Read(cer string, userID, gameID, flagID int32) (*netproto.Slot_UserAssignedMode, error) {

	const ProcedureName = "PrGs_Slot_UserAssignedMode_Read"

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

	if len(dbret[0].Rows) == 0 {
		return nil, nil
	}

	data := new(netproto.Slot_UserAssignedMode)

	if val, ok := dbret[0].GetValueByColName(0, "CreateTime").(time.Time); ok {
		data.CreateTime = proto.String(val.Format("2006-01-02 15:04:05.000"))
	} else {
		mlog.Warn("Procedure %s Read CreateTime error", ProcedureName)
	}

	if val, ok := dbret[0].GetValueByColName(0, "Mode").(int64); ok {
		data.Mode = proto.Int32(int32(val))
	} else {
		mlog.Warn("Procedure %s Read Mode error", ProcedureName)
	}

	if val, ok := dbret[0].GetValueByColName(0, "ExpireTime").(time.Time); ok {
		data.ExpireTime = proto.String(val.Format("2006-01-02 15:04:05.000"))
	} else {
		mlog.Warn("Procedure %s Read ExpireTime error", ProcedureName)
	}

	return data, nil
}

// Record 紀錄 UserAssignedMode
func (s *DBUserAssignedModeService) Record(cer string, userID, gameID, flagID, mode int32, expire string) error {

	const ProcedureName = "PrGs_Slot_UserAssignedMode_Assign"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", userID)              // UserID
	ps.AddIntInput("intGameID", gameID)              // GameID
	ps.AddIntInput("intFlagID", flagID)              // FlagID
	ps.AddIntInput("tinyMode", mode)                 // Mode
	ps.AddVarcharInput("chvExpireTime", expire, 255) // ExpireTime
	ps.AddVarcharOutput("chvMsg", "")                // chvMsg的輸出

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

// Delete 刪除 UserAssignedMode
func (s *DBUserAssignedModeService) Delete(cer string, userID, gameID, flagID int32) error {

	const ProcedureName = "PrGs_Slot_UserAssignedMode_Delete"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", userID) // UserID
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

// APIUserAssignedModeService 玩家被指定遊戲模式
type APIUserAssignedModeService struct{}

// Read 讀取 UserAssignedMode
func (s *APIUserAssignedModeService) Read(cer string, userID, gameID, flagID int32) (*netproto.Slot_UserAssignedMode, error) {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	query := map[string]interface{}{
		"playerID": userID,
		"gameID":   gameID,
		"flagID":   flagID,
	}

	ok, res := apiutil.GameSlotService.GetUserAssignedMode(rpc, query)
	if !ok {
		return nil, fmt.Errorf("API GET UserAssignedMode Error, %d %s", res.Code, res.Message)
	}

	if len(res.Data.SlotUserAssignedModeItem) == 0 {
		return nil, nil
	}

	data := new(netproto.Slot_UserAssignedMode)

	{
		// t, err := time.ParseInLocation(time.RFC3339Nano, res.Data.SlotUserAssignedModeItem[0].CreateTime, time.UTC)
		// if err != nil {
		// 	return nil, fmt.Errorf("API GET UserAssignedMode Parse CreateTime Error, %v", err)
		// }
		// data.CreateTime = proto.String(t.Local().Format("2006-01-02 15:04:05.000"))
		data.CreateTime = proto.String(res.Data.SlotUserAssignedModeItem[0].CreateTime.Format("2006-01-02 15:04:05.000"))
	}

	data.Mode = proto.Int32(res.Data.SlotUserAssignedModeItem[0].Mode)

	{
		// t, err := time.ParseInLocation(time.RFC3339Nano, res.Data.SlotUserAssignedModeItem[0].ExpireTime, time.UTC)
		// if err != nil {
		// 	return nil, fmt.Errorf("API GET UserAssignedMode Parse ExpireTime Error, %v", err)
		// }
		// data.ExpireTime = proto.String(t.Local().Format("2006-01-02 15:04:05.000"))
		data.ExpireTime = proto.String(res.Data.SlotUserAssignedModeItem[0].ExpireTime.Format("2006-01-02 15:04:05.000"))
	}

	return data, nil
}

// Record 紀錄 UserAssignedMode
func (s *APIUserAssignedModeService) Record(cer string, userID, gameID, flagID, mode int32, expire string) error {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	t, err := time.ParseInLocation("2006-01-02 15:04:05.000", expire, time.Local)
	if err != nil {
		return fmt.Errorf("API POST UserAssignedMode Parse expire Error, %v", err)
	}

	params := map[string]interface{}{
		"playerID": userID,
		"gameID":   gameID,
		"flagID":   flagID,
		"mode":     mode,
		// "expireTime": t.UTC().Format(time.RFC3339Nano),
		"expireTime": model.JSONTime{Time: t},
	}

	ok := apiutil.GameSlotService.PostUserAssignedMode(rpc, params)
	if !ok {
		return fmt.Errorf("API POST UserAssignedMode Error")
	}

	return nil
}

// Delete 刪除 UserAssignedMode
func (s *APIUserAssignedModeService) Delete(cer string, userID, gameID, flagID int32) error {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	query := map[string]interface{}{
		"playerID": userID,
		"gameID":   gameID,
		"flagID":   flagID,
	}

	ok := apiutil.GameSlotService.DeleteUserAssignedMode(rpc, query)
	if !ok {
		return fmt.Errorf("API DELETE UserAssignedMode Error")
	}

	return nil
}
