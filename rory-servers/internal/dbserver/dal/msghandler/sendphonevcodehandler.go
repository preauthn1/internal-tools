package msghandler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// SendPhoneVCodeHandler 發送驗證碼
type SendPhoneVCodeHandler struct{}

// IsHook ?
func (rh *SendPhoneVCodeHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_SendPhoneVCodeID)
}

// OnMessage ?
func (rh *SendPhoneVCodeHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcinfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SendPhoneVCode{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SendPhoneVCode)
	tel := mm.GetTel()
	codetype := mm.GetCodeType()

	userid := rpcinfo.GetUserID()

	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", userid)
	if tel != "" {
		ps.AddVarcharInput("chvTel", tel, 32)
	}
	ps.AddIntInput("tnyCodeType", codetype)
	ps.AddIntOutput("tnyCountdown", 0)
	ps.AddVarcharOutput("chvErrMsg", "")
	ps.AddIntOutput("intCodeID", 0)

	proc := db.NewProcedure("PrPs_User_SendSMSCode", ps)
	//ret, err := DbUtil.ExecProc(proc)
	var ret db.RecordSet
	var err error
	ret, err = dbutil.DBUtil.ExecProc(proc)

	code := 0
	retmessage := ""
	var codeID int32 = 0
	var countdown int32 = 0
	if err != nil {
		ctx.Error("执行存储过程PrPs_User_SendSMSCode错误, %v, %v", err, proc.GetSQLString())
		code = 0
		retmessage = "服务器错误"
	} else {
		if ret.GetReturnValue() == 1 {
			code = 1
		}
		codeID = int32(ret.GetOutputParamValue("intCodeID").(int64))
		countdown = int32(ret.GetOutputParamValue("tnyCountdown").(int64))
		retmessage = ret.GetOutputParamValue("chvErrMsg").(string)
	}

	coderet := &netproto.SendPhoneVCodeRet{}
	coderet.Code = proto.Int32(int32(code))
	coderet.Message = proto.String(retmessage)
	coderet.CountDown = proto.Int32(countdown)
	coderet.CodeID = proto.Int32(codeID)

	retmsg := new(msg.Message)
	retmsg.BClassID = int32(netproto.MessageBClassID_Hall)
	retmsg.SClassID = int32(netproto.HallMsgClassID_SendPhoneVCodeRetID)
	retmsg.MsgData = coderet

	SendSMSCode(codeID)

	return []*msg.Message{retmsg}
}

// SendSMSCode 用Web API發送第三方簡訊服務
func SendSMSCode(pCode int32) {
	codeID := fmt.Sprintf("codeID=%v", pCode)
	URLAPI := fmt.Sprintf("http://%s", dbutil.SMSWebAPI)
	//mlog.Debug(err.Error())
	resp, err := http.Post(URLAPI, "application/x-www-form-urlencoded", strings.NewReader(codeID))

	if err != nil {
		mlog.Error("調用 %s 錯誤, %v", URLAPI, err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// handle error
		mlog.Error("調用 %s 錯誤, %v", URLAPI, err)
	}
	mlog.Debug("調用 API 成功, 返回 %s", body)
}
