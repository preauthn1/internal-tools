package service

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	apiutil "example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	"example.invalid/ServerTeam/Servers/pkg/db"
	dbutil "example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// JackpotNoLinkService Jackpot相關存取DB服務
type JackpotNoLinkService interface {
	GetJackpotNoLink(cer string, gameID, flagID, tableNo, levelID int32) ([]*netproto.Slot_JackpotData, error)
	UpdateJackpotNoLink(cer string, gameID, flagID, tableNo, levelID int32, contribution, totalContribution, awardTotal int64) ([]*netproto.Slot_JackpotData, error)
}

// DBJackpotNoLinkService Jackpot相關存取DB服務
type DBJackpotNoLinkService struct{}

// GetJackpotNoLink 取得非連線型Jackpot資訊
func (s *DBJackpotNoLinkService) GetJackpotNoLink(cer string, gameID, flagID, tableNo, levelID int32) ([]*netproto.Slot_JackpotData, error) {

	const ProcedureName = "PrGs_Slot_JackpotPool_NoLink"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intGameID", gameID)   // GameID
	ps.AddIntInput("intFlagID", flagID)   // FlagID
	ps.AddIntInput("intTableNo", tableNo) // 桌子編號
	ps.AddIntInput("intLevelID", levelID) // Jackpot等級
	ps.AddVarcharOutput("chvMsg", "")     // chvMsg的輸出

	proc := db.NewProcedure(ProcedureName, ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		return nil, fmt.Errorf("Procedure %s execute error. %v; SQL: %s", ProcedureName, err, proc.GetSQLString())
	}

	msg := dbret.GetOutputParamValueString("chvMsg")

	if dbret.GetReturnValue() != 1 {
		return nil, fmt.Errorf("Procedure %s execute error. %v", ProcedureName, msg)
	}

	var result []*netproto.Slot_JackpotData

	for i := 0; i < len(dbret[0].Rows); i++ {

		data := new(netproto.Slot_JackpotData)

		if val, ok := dbret[0].GetValueByColName(i, "GameID").(int64); ok {
			data.GameID = proto.Int32(int32(val))
		}

		if val, ok := dbret[0].GetValueByColName(i, "FlagID").(int64); ok {
			data.FlagID = proto.Int32(int32(val))
		}

		if val, ok := dbret[0].GetValueByColName(i, "TableNo").(int64); ok {
			data.TableNo = proto.Int32(int32(val))
		}

		if val, ok := dbret[0].GetValueByColName(i, "LevelID").(int64); ok {
			data.Level = proto.Int32(int32(val))
		}

		if val, ok := dbret[0].GetValueByColName(i, "Contribution").(int64); ok {
			data.Contribution = proto.Int64(val)
		}

		if val, ok := dbret[0].GetValueByColName(i, "TotalContribution").(int64); ok {
			data.TotalContribution = proto.Int64(val)
		}

		if val, ok := dbret[0].GetValueByColName(i, "AwardTotal").(int64); ok {
			data.AwardTotal = proto.Int64(val)
		}

		result = append(result, data)
	}

	return result, nil
}

// UpdateJackpotNoLink 取得非連線型Jackpot資訊
func (s *DBJackpotNoLinkService) UpdateJackpotNoLink(cer string, gameID, flagID, tableNo, levelID int32, contribution, totalContribution, awardTotal int64) ([]*netproto.Slot_JackpotData, error) {

	const ProcedureName = "PrGs_Slot_JackpotPool_NoLink"

	ps := db.NewSQLParameters()

	ps.AddIntInput("intGameID", gameID)                          // GameID
	ps.AddIntInput("intFlagID", flagID)                          // FlagID
	ps.AddIntInput("intTableNo", tableNo)                        // 桌子編號
	ps.AddIntInput("intLevelID", levelID)                        // Jackpot等級
	ps.AddBigIntInput("lngContribution", contribution)           // 玩家貢獻值
	ps.AddBigIntInput("lngTotalContribution", totalContribution) // 累積貢獻值
	ps.AddBigIntInput("lngAwardTotal", awardTotal)               // 已發送獎金總和
	ps.AddVarcharOutput("chvMsg", "")                            // chvMsg的輸出

	proc := db.NewProcedure(ProcedureName, ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		return nil, fmt.Errorf("Procedure %s execute error. %v; SQL: %s", ProcedureName, err, proc.GetSQLString())
	}

	msg := dbret.GetOutputParamValueString("chvMsg")

	if dbret.GetReturnValue() != 1 {
		return nil, fmt.Errorf("Procedure %s execute error. %v", ProcedureName, msg)
	}

	var result []*netproto.Slot_JackpotData

	for i := 0; i < len(dbret[0].Rows); i++ {

		data := new(netproto.Slot_JackpotData)

		if val, ok := dbret[0].GetValueByColName(i, "GameID").(int64); ok {
			data.GameID = proto.Int32(int32(val))
		}

		if val, ok := dbret[0].GetValueByColName(i, "FlagID").(int64); ok {
			data.FlagID = proto.Int32(int32(val))
		}

		if val, ok := dbret[0].GetValueByColName(i, "TableNo").(int64); ok {
			data.TableNo = proto.Int32(int32(val))
		}

		if val, ok := dbret[0].GetValueByColName(i, "LevelID").(int64); ok {
			data.Level = proto.Int32(int32(val))
		}

		if val, ok := dbret[0].GetValueByColName(i, "Contribution").(int64); ok {
			data.Contribution = proto.Int64(val)
		}

		if val, ok := dbret[0].GetValueByColName(i, "TotalContribution").(int64); ok {
			data.TotalContribution = proto.Int64(val)
		}

		if val, ok := dbret[0].GetValueByColName(i, "AwardTotal").(int64); ok {
			data.AwardTotal = proto.Int64(val)
		}

		result = append(result, data)
	}

	return result, nil
}

// APIJackpotNoLinkService Jackpot相關存取DB服務
type APIJackpotNoLinkService struct{}

// GetJackpotNoLink 取得非連線型Jackpot資訊
func (s *APIJackpotNoLinkService) GetJackpotNoLink(cer string, gameID, flagID, tableNo, levelID int32) ([]*netproto.Slot_JackpotData, error) {
	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	params := map[string]interface{}{
		"isLink":  false,
		"flagID":  flagID,
		"tableNo": tableNo,
		"levelID": levelID,
	}

	ok, res := apiutil.GameSlotService.PostJackpot(rpc, gameID, params)
	if !ok {
		return nil, fmt.Errorf("API POST Jackpot Error, %d %s", res.Code, res.Message)
	}

	var data []*netproto.Slot_JackpotData

	for _, row := range res.Data.Jackpots {
		jp := new(netproto.Slot_JackpotData)
		jp.GameID = proto.Int32(row.GameID)
		jp.FlagID = proto.Int32(row.FlagID)
		jp.TableNo = proto.Int32(row.TableNo)
		jp.Level = proto.Int32(row.LevelID)
		jp.Contribution = proto.Int64(util.MoneyToCent(row.Contribution))
		jp.TotalContribution = proto.Int64(util.MoneyToCent(row.TotalContribution))
		jp.AwardTotal = proto.Int64(util.MoneyToCent(row.AwardTotal))

		data = append(data, jp)
	}

	return data, nil
}

// UpdateJackpotNoLink 取得非連線型Jackpot資訊
func (s *APIJackpotNoLinkService) UpdateJackpotNoLink(cer string, gameID, flagID, tableNo, levelID int32, contribution, totalContribution, awardTotal int64) ([]*netproto.Slot_JackpotData, error) {

	rpc := new(netproto.RPCInfo)
	rpc.Cer = proto.String(cer)

	params := map[string]interface{}{
		"isLink":            false,
		"flagID":            flagID,
		"tableNo":           tableNo,
		"levelID":           levelID,
		"contribution":      util.CentToMoney(contribution),
		"totalContribution": util.CentToMoney(totalContribution),
		"totalAward":        util.CentToMoney(awardTotal),
	}

	ok, res := apiutil.GameSlotService.PostJackpot(rpc, gameID, params)
	if !ok {
		return nil, fmt.Errorf("API POST Jackpot Error, %d %s", res.Code, res.Message)
	}

	var data []*netproto.Slot_JackpotData

	for _, row := range res.Data.Jackpots {
		jp := new(netproto.Slot_JackpotData)
		jp.GameID = proto.Int32(row.GameID)
		jp.FlagID = proto.Int32(row.FlagID)
		jp.TableNo = proto.Int32(row.TableNo)
		jp.Level = proto.Int32(row.LevelID)
		jp.Contribution = proto.Int64(util.MoneyToCent(row.Contribution))
		jp.TotalContribution = proto.Int64(util.MoneyToCent(row.TotalContribution))
		jp.AwardTotal = proto.Int64(util.MoneyToCent(row.AwardTotal))

		data = append(data, jp)
	}

	return data, nil
}
