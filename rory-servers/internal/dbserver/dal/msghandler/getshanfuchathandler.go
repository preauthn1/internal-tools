package msghandler

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// GetShanFuChatHandler 閃付聊天
type GetShanFuChatHandler struct{}

// IsHook ?
func (rh *GetShanFuChatHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_GetShanFuChatID)
}

// OnMessage ?
func (rh *GetShanFuChatHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GetShanFuChatMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.GetShanFuChatMsg)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	proc := db.NewProcedure("PrPs_GetShanFuKey", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetShanFuChatRetID), nil)
	retdata := new(netproto.GetShanFuChatRet)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程PrPs_GetShanFuKey错误,%v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 {
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				accountid := tb.GetValueByColNameString(0, "AccountID")
				publickey := tb.GetValueByColNameString(0, "Sign")
				sumamount := tb.GetValueByColNameInt64(0, "SumAmount")
				retdata.ClientUrl = proto.String(tb.GetValueByColNameString(0, "PayAgentClientUrl"))
				url := tb.GetValueByColName(0, "PayAgentUrl").(string)
				url = fmt.Sprintf("%v/agent/agents", url)

				stramount := "000"
				if sumamount > 0 {
					stramount = strconv.Itoa(int(sumamount))
				}
				var md5data string
				var sec int64
				if mm.GetType() == "1" { //ios md5
					//md5data, sec = iosMd5Data(mm, accountid, stramount, publickey)
					md5data, sec = androidMd5Data(mm, accountid, stramount, publickey)
				} else if mm.GetType() == "2" { //android md5
					md5data, sec = androidMd5Data(mm, accountid, stramount, publickey)
				}
				//md5data, sec := SendShanFuMd5Data(accountid, stramount, publickey)
				retdata.UserID = mm.UserID
				retdata.AgentID = mm.AgentID
				retdata.AppKey = proto.String(accountid)
				retdata.AppID = mm.AppID
				retdata.DeviceID = mm.DeviceID
				retdata.UserName = mm.UserName
				retdata.TimeStamp = proto.String(strconv.Itoa(int(sec)))
				retdata.SecretID = proto.String(md5data)
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}

func androidMd5Data(mm *netproto.GetShanFuChatMsg, account string, sumamount string, publickey string) (string, int64) {
	sec := time.Now().Unix()

	var buf bytes.Buffer
	buf.WriteString("app_id=")
	buf.WriteString(mm.GetAppID())
	buf.WriteString("&app_key=")
	buf.WriteString(account)
	buf.WriteString("&device_id=")
	buf.WriteString(mm.GetDeviceID())
	buf.WriteString("&id=")
	buf.WriteString(mm.GetAgentID())
	buf.WriteString("&timestamp=")
	buf.WriteString(strconv.Itoa(int(sec)))
	buf.WriteString("&user_id=")
	buf.WriteString(strconv.Itoa(int(mm.GetUserID())))
	buf.WriteString("&user_name=")
	buf.WriteString(mm.GetUserName())
	buf.WriteString(publickey)

	signVal := buf.String()

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signVal))
	cipherStr := md5Ctx.Sum(nil)
	sign := hex.EncodeToString(cipherStr)
	sign = strings.ToLower(sign)
	return sign, sec
}

func iosMd5Data(mm *netproto.GetShanFuChatMsg, account string, sumamount string, publickey string) (string, int64) {
	sec := time.Now().Unix()

	var buf bytes.Buffer
	buf.WriteString("app_key=")
	buf.WriteString(account)
	buf.WriteString("&app_id=")
	buf.WriteString(mm.GetAppID())
	buf.WriteString("&device_id=")
	buf.WriteString(mm.GetDeviceID())
	buf.WriteString("&id=")
	buf.WriteString(mm.GetAppID())
	buf.WriteString("&timestamp=")
	buf.WriteString(strconv.Itoa(int(sec)))
	buf.WriteString("&user_id=")
	buf.WriteString(strconv.Itoa(int(mm.GetUserID())))
	buf.WriteString("&user_name=")
	buf.WriteString(mm.GetUserName())
	buf.WriteString(publickey)

	signVal := buf.String()
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signVal))
	cipherStr := md5Ctx.Sum(nil)
	sign := hex.EncodeToString(cipherStr)
	sign = strings.ToLower(sign)

	return sign, sec
}
