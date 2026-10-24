package service

import (
	"fmt"
	"strings"
	"time"

	"example.invalid/golang/protobuf/proto"

	apiutil "example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	"example.invalid/ServerTeam/Servers/pkg/db"
	dbutil "example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// GameLogService 玩家紀錄相關
type GameLogService interface {
	Add(cer string, log *netproto.Slot_MSG_GameLogAdd) error
	BatchAdd(cer string, logs []*netproto.Slot_MSG_GameLogAdd) error
	Search(cer string, userID, gameID, index, size int32, start, end string) (int32, []*netproto.Slot_GameLog, error)
	MarshalGameLog(log netproto.Slot_MSG_GameLogAdd) string
}

// DBGameLogService 玩家紀錄相關
type DBGameLogService struct{}

// Add 新增遊戲紀錄
func (s *DBGameLogService) Add(cer string, log *netproto.Slot_MSG_GameLogAdd) error {

	if log == nil {
		return nil
	}

	const ProcedureName = "PrGs_Slot_AddGameLog"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", log.GetUserID())                  // 玩家ID
	ps.AddIntInput("intGameID", log.GetGameID())                  // 遊戲ID
	ps.AddIntInput("intServerID", log.GetServerID())              // ServerID
	ps.AddBigIntInput("lngCurrentAmount", log.GetCurrentAmount()) // 當前金額
	ps.AddBigIntInput("lngBetMoney", log.GetBetMoney())           // 花費金額
	ps.AddBigIntInput("lngWinMoney", log.GetWinMoney())           // 贏得金額
	ps.AddBigIntInput("lngFinalScore", log.GetFinalScore())       // 最終金額
	ps.AddIntInput("intGameType", log.GetGameType())              // 遊戲類型
	ps.AddIntInput("intLogType", log.GetLogType())                // 紀錄類型
	ps.AddVarcharInput("chvGameData", log.GetGameData(), 100)     // 遊戲資料

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

// BatchAdd 新增遊戲紀錄(多筆)
func (s *DBGameLogService) BatchAdd(cer string, logs []*netproto.Slot_MSG_GameLogAdd) error {

	if len(logs) == 0 {
		return nil
	}

	const ProcedureName = "PrGs_Slot_AddGameLog_Batch"

	ps := db.NewSQLParameters()

	// CreateTime,GameID,ServerID,UserID,CurrentAmount,BetMoney,WinMoney,FinalScore,GameType,LogType,GameData
	// 2019-12-31 10:00:55.960;1005;10104;15267008;1057380;280;0;1057100;1;1;p7,p9,p4,p9|p5,p7,p11,p5|p11,p9,p8,p11|p10,p3,p7,p11|p8,p8,p9,p7
	// 按照這個順序串起來(共11個欄位)，欄位間用 ; 區隔，資料間用 ! 區隔

	var strs []string
	for _, row := range logs {
		strs = append(strs, s.MarshalGameLog(*row))
	}

	ps.AddVarcharInput("chvGameDataList", strings.Join(strs, "!"), -1)
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

// Search 查詢遊戲紀錄
func (s *DBGameLogService) Search(cer string, userID, gameID, index, size int32, start, end string) (int32, []*netproto.Slot_GameLog, error) {

	const ProcedureName = "PrGs_Slot_GetGameLog"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", userID)           // 玩家ID
	ps.AddIntInput("intGameID", gameID)           // 遊戲ID
	ps.AddVarcharInput("dtStartTime", start, 100) // 開始時間
	ps.AddVarcharInput("dtEndTime", end, 100)     // 結束時間
	ps.AddIntInput("intPageIndex", index)         // 起始筆數
	ps.AddIntInput("intPageSize", size)           // 每頁幾筆
	ps.AddIntOutput("intPageCount", 0)
	ps.AddVarcharOutput("chvMsg", "")

	proc := db.NewProcedure(ProcedureName, ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		return 0, nil, fmt.Errorf("Procedure %s execute error. %v; SQL: %s", ProcedureName, err, proc.GetSQLString())
	}

	msg := dbret.GetOutputParamValueString("chvMsg")

	if dbret.GetReturnValue() != 1 {
		return 0, nil, fmt.Errorf("Procedure %s execute error. %v", ProcedureName, msg)
	}

	pageCount := dbret.GetOutputParamValueInt32("intPageCount")

	var result []*netproto.Slot_GameLog

	for i := 0; i < len(dbret[0].Rows); i++ {
		data := new(netproto.Slot_GameLog)

		if val, ok := dbret[0].GetValueByColName(i, "RowNum").(int64); ok {
			data.RowNum = proto.Int32(int32(val))
		}

		if val, ok := dbret[0].GetValueByColName(i, "OrderID").(int64); ok {
			data.OrderID = proto.Int64(val)
		}

		if val, ok := dbret[0].GetValueByColName(i, "CreateTime").(time.Time); ok {
			data.CreateTime = proto.String(val.Format("2006-01-02 15:04:05"))
		}

		if val, ok := dbret[0].GetValueByColName(i, "BetMoney").(int64); ok {
			data.BetMoney = proto.Int64(val)
		}

		if val, ok := dbret[0].GetValueByColName(i, "FinalScore").(int64); ok {
			data.Profit = proto.Int64(val)
		}

		result = append(result, data)
	}

	return pageCount, result, nil
}

// MarshalGameLog 將 netproto.Slot_MSG_GameLogAdd 包成字串
func (s *DBGameLogService) MarshalGameLog(log netproto.Slot_MSG_GameLogAdd) string {
	return fmt.Sprintf("%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v",
		log.GetCreateTime(),
		log.GetGameID(),
		log.GetServerID(),
		log.GetUserID(),
		log.GetCurrentAmount(),
		log.GetBetMoney(),
		log.GetWinMoney(),
		log.GetFinalScore(),
		log.GetGameType(),
		log.GetLogType(),
		log.GetGameData())
}

// APIGameLogService 玩家紀錄相關
type APIGameLogService struct{}

// Add 新增遊戲紀錄
func (s *APIGameLogService) Add(cer string, log *netproto.Slot_MSG_GameLogAdd) error {
	return s.BatchAdd(cer, []*netproto.Slot_MSG_GameLogAdd{log})
}

// BatchAdd 新增遊戲紀錄(多筆)
func (s *APIGameLogService) BatchAdd(cer string, logs []*netproto.Slot_MSG_GameLogAdd) error {

	if len(logs) == 0 {
		return nil
	}

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	var data []string

	for _, log := range logs {
		data = append(data, s.MarshalGameLog(*log))
	}

	params := map[string]interface{}{
		"data": data,
	}

	ok := apiutil.GameSlotService.PostLog(rpc, params)
	if !ok {
		return fmt.Errorf("API POST GameLog Error")
	}

	return nil
}

// Search 查詢遊戲紀錄 Client 傳 JS Date 物件 的 toISOString() 進來，2020-06-16T16:00:00.000Z
func (s *APIGameLogService) Search(cer string, userID, gameID, index, size int32, start, end string) (int32, []*netproto.Slot_GameLog, error) {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	query := map[string]interface{}{
		"playerID":      userID,
		"gameID":        gameID,
		"gameTimeBegin": start,
		"gameTimeEnd":   end,
		"pageIndex":     index,
		"pageSize":      size,
	}

	ok, res := apiutil.GameSlotService.GetLog(rpc, query)
	if !ok {
		return 0, nil, fmt.Errorf("API GET GameLog Error, %d %s", res.Code, res.Message)
	}

	if len(res.Data.Items) == 0 {
		return 0, nil, nil
	}

	var logs []*netproto.Slot_GameLog

	for _, row := range res.Data.Items {
		log := new(netproto.Slot_GameLog)
		log.CreateTime = proto.String(row.GameTime.Format(time.RFC3339Nano))
		log.OrderID = proto.Int64(row.LogID)
		log.RowNum = proto.Int32(row.RowNum)
		log.Profit = proto.Int64(util.MoneyToCent(row.FinalAmount))
		log.BetMoney = proto.Int64(util.MoneyToCent(row.BetAmount))

		logs = append(logs, log)
	}

	return res.Data.TotalCount, logs, nil
}

// MarshalGameLog 將 netproto.Slot_MSG_GameLogAdd 包成字串
func (s *APIGameLogService) MarshalGameLog(log netproto.Slot_MSG_GameLogAdd) string {

	// 日期轉換, API 那邊吃 UTC 日期格式
	if t, err := time.ParseInLocation("2006-01-02 15:04:05.000", log.GetCreateTime(), time.Local); err == nil {
		log.CreateTime = proto.String(t.Format(time.RFC3339Nano))
	}

	return fmt.Sprintf("%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v",
		log.GetCreateTime(),
		log.GetGameID(),
		log.GetServerID(),
		log.GetUserID(),
		util.CentToMoney(log.GetCurrentAmount()),
		util.CentToMoney(log.GetBetMoney()),
		util.CentToMoney(log.GetWinMoney()),
		util.CentToMoney(log.GetFinalScore()),
		log.GetGameType(),
		log.GetLogType(),
		log.GetGameData())
}
