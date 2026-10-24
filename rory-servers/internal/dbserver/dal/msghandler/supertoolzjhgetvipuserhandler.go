package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// SuperToolZjhGetVipUserHandler 寫遊戲記錄詳情
type SuperToolZjhGetVipUserHandler struct{}

// IsHook ?
func (rh *SuperToolZjhGetVipUserHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.GameMessageClassID_SuperToolZjhGetVipUserID)
}

// OnMessage ?
func (rh *SuperToolZjhGetVipUserHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	ps := db.NewSQLParameters()
	proc := db.NewProcedure("PrGs_SuperTool_ZJH_GetVIPUser", ps)

	ret, err := dbutil.DBUtil.ExecProc(proc)

	var retcode int32 = 1
	if err != nil {
		ctx.Error("执行存储过程PrGs_SuperTool_ZJH_GetVIPUser错误, %v, %s", err, proc.GetSQLString())
		retcode = 0
	} else {
		retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Game), int32(netproto.GameMessageClassID_SuperToolZjhGetVipUserRetID), nil)
		retdata := new(netproto.SuperToolZjhVipUserDataList)
		for i := 0; i < len(ret[0].Rows); i++ {
			vipuser := new(netproto.SuperToolZjhVipUserData)

			vipuser.UserID = proto.Int32(int32(ret[0].GetValueByColName(i, "UserID").(int64)))
			vipuser.NextType = proto.Int64(ret[0].GetValueByColName(i, "NextType").(int64))
			vipuser.NextWinLose = proto.Int64(ret[0].GetValueByColName(i, "NextWinLose").(int64))
			vipuser.NextWinLoseNum = proto.Int64(ret[0].GetValueByColName(i, "NextWinLoseNum").(int64))
			vipuser.Finish = proto.Int64(ret[0].GetValueByColName(i, "Finish").(int64))
			vipuser.FinishNextWinLoseNum = proto.Int64(ret[0].GetValueByColName(i, "FinishNextWinLoseNum").(int64))
			vipuser.NextTable = proto.Int32(int32(ret[0].GetValueByColName(i, "NextTable").(int64)))
			vipuser.Status = proto.Int32(int32(ret[0].GetValueByColName(i, "Status").(int64)))
			vipuser.TodaySumcharge = proto.Int64(ret[0].GetValueByColName(i, "TodaySumcharge").(int64))
			vipuser.TodaySumConvert = proto.Int64(ret[0].GetValueByColName(i, "TodaySumConvert").(int64))
			vipuser.TodayUserWinLose = proto.Int64(ret[0].GetValueByColName(i, "TodayUserWinLose").(int64))
			retdata.DataList = append(retdata.DataList, vipuser)
		}
		retmsg.MsgData = retdata
		return []*msg.Message{retmsg}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(retcode)}

}
