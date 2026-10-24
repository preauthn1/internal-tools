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

// FreeSpinService 玩家免費遊戲相關
type FreeSpinService interface {
	Read(cer string, userID, gameID, flagID int32) (*netproto.Slot_FreeSpinDetail, error)
	Record(cer string, userID, gameID, flagID, patternID int32, originID, betMoney int64, betRate, symbolGroupID, times int32, winMoney int64, winTimes int32, data string) error
}

// DBFreeSpinService 玩家免費遊戲相關
type DBFreeSpinService struct{}

// Read 讀取 FreeSpin
func (s *DBFreeSpinService) Read(cer string, userID, gameID, flagID int32) (*netproto.Slot_FreeSpinDetail, error) {
	const ProcedureName = "PrGs_Slot_FreeSpinRead"

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

	data := new(netproto.Slot_FreeSpinDetail)

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
		data.OriginID = proto.Int64(val)
	} else {
		mlog.Warn("Procedure %s Read OriginID error", ProcedureName)
	}

	if val, ok := dbret[0].GetValueByColName(0, "FreeSpinGroup").(int64); ok {
		data.FreeSpinGroup = proto.Int32(int32(val))
	} else {
		mlog.Warn("Procedure %s Read FreeSpinGroup error", ProcedureName)
	}

	if val, ok := dbret[0].GetValueByColName(0, "Times").(int64); ok {
		data.Times = proto.Int32(int32(val))
	} else {
		mlog.Warn("Procedure %s Read Times error", ProcedureName)
	}

	if val, ok := dbret[0].GetValueByColName(0, "UserBet").(int64); ok {
		data.UserBet = proto.Int64(val)
	} else {
		mlog.Warn("Procedure %s Read UserBet error", ProcedureName)
	}

	if val, ok := dbret[0].GetValueByColName(0, "ChooseRate").(int64); ok {
		data.ChooseRate = proto.Int32(int32(val))
	} else {
		mlog.Warn("Procedure %s Read ChooseRate error", ProcedureName)
	}

	if val, ok := dbret[0].GetValueByColName(0, "WinMoney").(int64); ok {
		data.WinMoney = proto.Int64(val)
	} else {
		mlog.Warn("Procedure %s Read WinMoney error", ProcedureName)
	}

	if val, ok := dbret[0].GetValueByColName(0, "WinTimes").(int64); ok {
		data.WinTimes = proto.Int32(int32(val))
	} else {
		mlog.Warn("Procedure %s Read WinTimes error", ProcedureName)
	}

	if val, ok := dbret[0].GetValueByColName(0, "Data").(string); ok {
		data.Data = proto.String(val)
	} else {
		mlog.Warn("Procedure %s Read Data error", ProcedureName)
	}

	return data, nil
}

// Record 紀錄 FreeSpin
func (s *DBFreeSpinService) Record(cer string, userID, gameID, flagID, patternID int32, originID, betMoney int64, betRate, symbolGroupID, times int32, winMoney int64, winTimes int32, data string) error {

	const ProcedureName = "PrGs_Slot_FreeSpinRecord"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", userID)               // UserID
	ps.AddIntInput("intGameID", gameID)               // GameID
	ps.AddIntInput("intFlagID", flagID)               // FlagID
	ps.AddIntInput("intPatternID", patternID)         // PatternID
	ps.AddBigIntInput("intOriginID", originID)        // OriginID
	ps.AddBigIntInput("intUserBet", betMoney)         // UserBet
	ps.AddIntInput("intChooseRate", betRate)          // ChooseRate
	ps.AddIntInput("intFreeSpinGroup", symbolGroupID) // FreeSpinGroup
	ps.AddIntInput("intTimes", times)                 // Times
	ps.AddBigIntInput("lngWinMoney", winMoney)        // WinMoney
	ps.AddIntInput("intWinTimes", winTimes)           // WinTimes
	ps.AddVarcharInput("txtData", data, 255)          // 額外資訊
	ps.AddVarcharOutput("chvMsg", "")                 // chvMsg的輸出

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

// APIFreeSpinService 玩家免費遊戲相關
type APIFreeSpinService struct{}

// Read 讀取 FreeSpin
func (s *APIFreeSpinService) Read(cer string, userID, gameID, flagID int32) (*netproto.Slot_FreeSpinDetail, error) {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	query := map[string]interface{}{
		"playerID": userID,
		"gameID":   gameID,
		"flagID":   flagID,
	}

	ok, res := apiutil.GameSlotService.GetFreeSpin(rpc, query)
	if !ok {
		return nil, fmt.Errorf("API Get FreeSpin Error, %d %s", res.Code, res.Message)
	}

	if len(res.Data.SlotFreeSpinRecord) == 0 {
		return nil, nil
	}

	data := new(netproto.Slot_FreeSpinDetail)

	data.UserID = proto.Int32(res.Data.SlotFreeSpinRecord[0].PlayerID)
	data.GameID = proto.Int32(res.Data.SlotFreeSpinRecord[0].GameID)
	data.FlagID = proto.Int32(res.Data.SlotFreeSpinRecord[0].FlagID)
	data.PatternID = proto.Int32(res.Data.SlotFreeSpinRecord[0].PatternID)
	data.OriginID = proto.Int64(res.Data.SlotFreeSpinRecord[0].OriginID)
	data.FreeSpinGroup = proto.Int32(res.Data.SlotFreeSpinRecord[0].FreeSpinGroup)
	data.Times = proto.Int32(res.Data.SlotFreeSpinRecord[0].Times)
	data.UserBet = proto.Int64(util.MoneyToCent(res.Data.SlotFreeSpinRecord[0].BetAmount))
	data.ChooseRate = proto.Int32(res.Data.SlotFreeSpinRecord[0].ChooseRate)
	data.WinMoney = proto.Int64(util.MoneyToCent(res.Data.SlotFreeSpinRecord[0].WinAmount))
	data.WinTimes = proto.Int32(res.Data.SlotFreeSpinRecord[0].WinTimes)
	data.Data = proto.String(res.Data.SlotFreeSpinRecord[0].Data)

	return data, nil
}

// Record 紀錄 FreeSpin
func (s *APIFreeSpinService) Record(cer string, userID, gameID, flagID, patternID int32, originID, betMoney int64, betRate, symbolGroupID, times int32, winMoney int64, winTimes int32, data string) error {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	params := map[string]interface{}{
		"playerID":      userID,
		"gameID":        gameID,
		"flagID":        flagID,
		"patternID":     patternID,
		"originID":      originID,
		"betAmount":     util.CentToMoney(betMoney),
		"chooseRate":    betRate,
		"freeSpinGroup": symbolGroupID,
		"times":         times,
		"winAmount":     util.CentToMoney(winMoney),
		"winTimes":      winTimes,
		"data":          data,
	}

	ok := apiutil.GameSlotService.PostFreeSpin(rpc, params)
	if !ok {
		return fmt.Errorf("API POST FreeSpin Error")
	}

	return nil
}
