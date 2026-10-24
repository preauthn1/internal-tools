package msghandler

import (
	"fmt"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// GetGameCategoryListHandler 遊戲分類清單
type GetGameCategoryListHandler struct{}

// IsHook ?
func (rh *GetGameCategoryListHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:GetGameCategoryListHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_DBServer) &&
		sclsid == int32(netproto.DBServerClassID_GetGameCategoryListID)
}

// OnMessage ?
func (rh *GetGameCategoryListHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	// rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
	// 	mm := &netproto.GetGameCategoryListMsg{}
	// 	return mm
	// })

	// if rmm == nil || rpcInfo == nil {
	// 	return []*msg.Message{GetDBRetMsg(0)}
	// }

	// mm := rmm.(*netproto.GetGameCategoryListMsg)

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_GetGameCategoryListRetID), nil)

	retdata, err := rh.GetGameCategoryList()
	if err != nil {
		mlog.Error("%s", err.Error())
		return []*msg.Message{retmsg}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}

// GetGameCategoryList ?
func (rh *GetGameCategoryListHandler) GetGameCategoryList() (*netproto.GameCategoryListDB, error) {
	ps := db.NewSQLParameters()

	// ps.AddIntInput("intAreaID", int32(areaid))

	proc := db.NewProcedure("PrGs_Game_GetGameCategory", ps)

	retdata := new(netproto.GameCategoryListDB)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		return nil, fmt.Errorf("执行存储过程PrGs_Game_GetGameCategory错误, %s", err.Error())
	}

	if dbret.GetRetTableCount() > 1 {
		tb := dbret[0]
		if len(tb.Rows) > 0 {
			for i := 0; i < len(tb.Rows); i++ {
				var pdata = new(netproto.CategoryInfoDB)
				pdata.CategoryID = proto.Int32(tb.GetValueByColNameInt32(i, "CategoryID"))
				pdata.Name = proto.String(tb.GetValueByColNameString(i, "Name"))
				retdata.GameTabs = append(retdata.GameTabs, pdata)
			}
		}

		tb = dbret[1]
		if len(tb.Rows) > 0 {
			for i := 0; i < len(tb.Rows); i++ {
				var pdata = new(netproto.GameCategoryIDsDB)
				pdata.GameID = proto.Int32(tb.GetValueByColNameInt32(i, "GameID"))
				pdata.CategoryIDs = proto.String(strings.ReplaceAll(tb.GetValueByColNameString(i, "CategoryIDs"), " ", ""))
				retdata.CategoryList = append(retdata.CategoryList, pdata)
			}
		}

		tb = dbret[2]
		if len(tb.Rows) > 0 {
			for i := 0; i < len(tb.Rows); i++ {
				var pdata = new(netproto.GameCategoryOptionDB)
				pdata.AreaID = proto.Int64(int64(tb.GetValueByColNameInt32(i, "AreaID")))
				// fmt.Println(fmt.Sprintf("%v, %v", tb.GetValueByColNameInt32(i, "AreaID"), tb.GetValueByColName(i, "IsVisible")))
				var isvisible = fmt.Sprintf("%v", tb.GetValueByColName(i, "IsVisible"))
				if isvisible == "true" {
					pdata.IsVisible = proto.Int32(1)
				} else {
					pdata.IsVisible = proto.Int32(0)
				}
				// fmt.Println(fmt.Sprintf("%v, %v", pdata.GetAreaID(), pdata.GetIsVisible()))
				retdata.Options = append(retdata.Options, pdata)
			}
		}
		// fmt.Println("")
	}
	// fmt.Println(retdata.GetGameTabs())
	// fmt.Println(retdata.GetCategoryList())
	// fmt.Println(retdata.GetOptions())

	return retdata, nil
}
