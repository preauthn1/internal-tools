package msghandler

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"
	uuid "example.invalid/satori/go.uuid"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

type UserInitBindProcess struct {
	UserID    int32
	Cer       string
	HDCode    string
	HDType    int32
	Tel       string
	UserData  netproto.UserHallLogonData
	Password  string
	LoginName string
}

func TestFoo(t *testing.T) {
	data := []byte{0, 0, 0, 1}
	fmt.Println(data)
	bytesBuffer := bytes.NewBuffer(data)
	var order binary.ByteOrder
	order = binary.BigEndian

	var x int32
	g := binary.Read(bytesBuffer, order, x)
	fmt.Println(g)
}

/**
part1 : 從遊客登入綁定帳號到兌換帳號綁定流程測試
part2 : 綁定帳號後離開大廳，重新利用綁定帳號登入
part3 : 修改保險箱密碼，修改頭像，修改密碼
*/

//func TestUserInitBindProcess(t *testing.T) {
//
//	initDBConnection()
//	process := new(UserInitBindProcess)
//	process.setUser(1021234, "88898498081", "")
//
//	if process.UserID == 0 {
//		process.guestLogon(t)
//		process.guestSendBindSMS(t)
//		process.bindAccount(t)
//	}
//	process.requestUserInfo(t)
//
//	process.sendPhoneVCode(t, 4)
//	process.bindZhifubao(t, 4)
//
//	process.userLeaveHall(t)
//	process.userLogon(t)
//
//	process.modifyBankPassword(t)
//	process.modifyFace(t)
//	process.modifyPassword(t)
//
//	//process.sendPhoneVCode(t, 5)
//	//process.bindZhifubao(t, 5)
//
//	//process.sendPhoneVCode(t, 5)
//	//process.bindBankCard(t, 5)
//
//}

func (process *UserInitBindProcess) setUser(userID int32, userTel string, userCer string) {
	process.UserID = userID
	process.Tel = userTel
	process.Cer = userCer
}

func (process *UserInitBindProcess) requestUserInfo(t *testing.T) {
	msgData := new(netproto.RequestUserHallInfo)
	msgData.SiteID = proto.Int32(1)
	rpcInfo := process.getRPCInfo()
	//rpc詳細資料, MsgData要放入Marshal後的資料
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) //MsgData要放入Marshal後的資料
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_RequestUserHallInfoID), rpcMessage)
	handler := new(RequestUserInfoHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_RequestUserHallInfoID hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_UserHallInfoID) {
		t.Error("HallMsgClassID_GuestLoginID OnMessage 錯誤")
	}
	userLoginRet := ret[0].MsgData.(*netproto.UserLoginRet)
	process.Cer = *userLoginRet.Cer
	process.UserData = *userLoginRet.UserData
	t.Log(ret[0].MsgData)
}

func (process *UserInitBindProcess) bindBankCard(t *testing.T, codeType int32) {
	vcode, err := process.getVCodeByCodeType(codeType)
	if err != nil {
		t.Error(err)
	}
	msgData := new(netproto.BankCardInfo)
	msgData.BankCardNum = proto.String("5544332211")
	msgData.RealName = proto.String("銀行姓名")
	msgData.BankCode = proto.String("001")
	msgData.CodeType = proto.Int32(codeType)
	msgData.VCode = proto.String(vcode)
	rpcInfo := process.getRPCInfo()

	//rpc詳細資料, MsgData要放入Marshal後的資料
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) //MsgData要放入Marshal後的資料

	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_BindBankCardID), rpcMessage)
	handler := new(BindBankCardHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_BindBankCardID hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_BindBankCardRetID) {
		t.Error("HallMsgClassID_BindBankCardID OnMessage 錯誤")
	}
	t.Log(ret[0].MsgData)

}

func (process *UserInitBindProcess) bindZhifubao(t *testing.T, codeType int32) {
	vcode, err := process.getVCodeByCodeType(codeType)
	if err != nil {
		t.Error(err)
	}
	msgData := new(netproto.ZhifubaoInfo)
	msgData.Zhifubao = proto.String("1122334455")
	msgData.RealName = proto.String("真假")
	msgData.CodeType = proto.Int32(codeType)
	msgData.VCode = proto.String(vcode)
	rpcInfo := process.getRPCInfo()

	//rpc詳細資料, MsgData要放入Marshal後的資料
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) //MsgData要放入Marshal後的資料

	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_BindZhifubaoID), rpcMessage)
	handler := new(BindZhifubaoHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_BindZhifubaoID hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_BindZhifubaoRetID) {
		t.Error("HallMsgClassID_BindZhifubaoID OnMessage 錯誤")
	}
	t.Log(ret[0].MsgData)

}

func (process *UserInitBindProcess) bindAccount(t *testing.T) {

	vcode, err := process.getVCodeByCodeType(1)
	if err != nil {
		t.Error(err)
	}

	msgData := new(netproto.BindGuestAccount)
	msgData.Tel = proto.String(process.Tel)

	msgData.VCode = proto.String(vcode)
	msgData.Password = [REDACTED]"123456")
	rpcInfo := process.getRPCInfo()

	//rpc詳細資料, MsgData要放入Marshal後的資料
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) //MsgData要放入Marshal後的資料
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_BindGuestAccountID), rpcMessage)
	handler := new(BindGuestAccountHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_BindGuestAccountID hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_BindGuestAccountRetID) {
		t.Error("HallMsgClassID_GuestLoginID OnMessage 錯誤")
	}
	t.Log(ret[0].MsgData)
}

func (process *UserInitBindProcess) getVCodeByCodeType(codeType int32) (string, error) {

	sql := fmt.Sprintf("SELECT Tel, VCode, CodeType FROM Users.UserCode WHERE Tel = '%s' AND CodeType = %d ORDER BY CreateCodeTime DESC", process.Tel, codeType)
	ret, err := dbutil.DBUtil.ExecSQL(sql)
	if err != nil {
		return "", err
	}
	vcode := ret[0].GetSingleValue("VCode").(string)
	if vcode == "" {
		return "", errors.New("獲取驗證碼失敗")
	}

	return vcode, nil
}
func (process *UserInitBindProcess) guestLogon(t *testing.T) {

	//Client傳入資料的協議號
	msgData := new(netproto.GuestLogin)
	msgData.Version = proto.String("0.0.0")
	msgData.HDType = proto.Int32(1)
	msgData.HDCode = proto.String(genUUID())
	msgData.SiteID = proto.Int32(1)
	msgData.PlatformID = proto.Int32(0)
	msgData.ServerID = proto.Int32(0)
	//rpc資料必填欄位
	rpcInfo := process.getRPCInfo()

	//rpc詳細資料, MsgData要放入Marshal後的資料
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) //MsgData要放入Marshal後的資料
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GuestLoginID), rpcMessage)

	handler := new(GuestLogonHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_GuestLoginID hook 錯誤")
	}

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_LoginRetID) {
		t.Error("HallMsgClassID_GuestLoginID OnMessage 錯誤")
	}
	userLoginRet := ret[0].MsgData.(*netproto.UserLoginRet)
	process.UserID = *userLoginRet.UserID
	process.Cer = *userLoginRet.Cer
	process.HDCode = *userLoginRet.HDCode
	process.HDType = *userLoginRet.HDType
	process.UserData = *userLoginRet.UserData
	t.Log(ret[0].MsgData)
}

func (process *UserInitBindProcess) guestSendBindSMS(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	process.Tel = "888" + fmt.Sprintf("%08d", rand.Intn(100000000))
	process.sendPhoneVCode(t, 1)
}

func (process *UserInitBindProcess) sendPhoneVCode(t *testing.T, codeType int32) {
	//Client傳入資料的協議號
	msgData := new(netproto.SendPhoneVCode)
	msgData.Tel = proto.String(process.Tel)
	msgData.CodeType = proto.Int32(codeType)

	//rpc資料必填欄位
	rpcInfo := process.getRPCInfo()

	//rpc詳細資料, MsgData要放入Marshal後的資料
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) //MsgData要放入Marshal後的資料

	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SendPhoneVCodeID), rpcMessage)

	handler := new(SendPhoneVCodeHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_SendPhoneVCodeID hook 錯誤")
	}

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_SendPhoneVCodeRetID) {
		t.Error("OnMessage 錯誤")
	}
	t.Log(ret[0].MsgData)
}

func (process *UserInitBindProcess) getRPCInfo() *netproto.RPCInfo {
	ipaddr := "0.0.0.0"
	userid := int32(process.UserID)
	rpcInfo := new(netproto.RPCInfo)
	rpcInfo.IPAddress = proto.String(ipaddr)
	rpcInfo.UserID = proto.Int32(userid)
	rpcInfo.ConnID = proto.String("fake")
	rpcInfo.Cer = proto.String(process.Cer)
	rpcInfo.RouteServerID = proto.Int32(1)
	return rpcInfo
}

func genUUID() string {
	foo := uuid.NewV4()
	uid := fmt.Sprintf("%s", foo)
	return uid
}

func (process *UserInitBindProcess) userLeaveHall(t *testing.T) {
	//Client傳入資料的協議號
	msgData := new(netproto.UserLogout)
	msgData.UserID = proto.Int32(process.UserID)
	msgData.Cer = proto.String(process.Cer)

	//rpc資料必填欄位
	rpcInfo := process.getRPCInfo()

	//rpc詳細資料, MsgData要放入Marshal後的資料
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) //MsgData要放入Marshal後的資料
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserLogoutID), rpcMessage)

	handler := new(UserLeaveHallHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_UserLeaveHall hook 錯誤")
	}

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_DBRetID) {
		t.Error("HallMsgClassID_UserLeaveHall OnMessage 錯誤")
	}
	t.Log(ret[0].MsgData)
}

func (process *UserInitBindProcess) userLogon(t *testing.T) {
	//Client傳入資料的協議號
	msgData := new(netproto.UserLogin)
	msgData.HDCode = proto.String(genUUID())
	msgData.HDType = proto.Int32(1)
	msgData.SiteID = proto.Int32(1)
	msgData.Version = proto.String("0.0.0")
	msgData.PlatformID = proto.Int32(0)
	msgData.Cer = proto.String(process.Cer)
	msgData.LoginName = proto.String(process.LoginName)
	msgData.Password = [REDACTED]
	msgData.UserID = proto.Int32(process.UserID)

	//rpc資料必填欄位
	rpcInfo := process.getRPCInfo()

	//rpc詳細資料, MsgData要放入Marshal後的資料
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) //MsgData要放入Marshal後的資料
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserLoginID), rpcMessage)

	handler := new(UserLogonHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_UserLoginID hook 錯誤")
	}

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_LoginRetID) {
		t.Error("HallMsgClassID_UserLoginID OnMessage 錯誤")
	}
	t.Log(ret[0].MsgData)
}
func (process *UserInitBindProcess) modifyBankPassword(t *testing.T) {
	msgData := new(netproto.ModifyBankPassword)
	msgData.NewPassword = [REDACTED]"7654321")
	msgData.OldPwd = [REDACTED]
	rpcInfo := process.getRPCInfo()
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) //MsgData要放入Marshal後的資料
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ModifyBankPasswordID), rpcMessage)
	handler := new(ModifyBankPasswordHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_ModifyBankPasswordID hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_ModifyBankPasswordRetID) {
		t.Error("HallMsgClassID_ModifyBankPasswordID OnMessage 錯誤")
	}
	t.Log(ret[0].MsgData)
}
func (process *UserInitBindProcess) modifyFace(t *testing.T) {
	msgData := new(netproto.ModifyFace)
	msgData.FaceID = proto.Int32(123)
	rpcInfo := process.getRPCInfo()

	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) //MsgData要放入Marshal後的資料

	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ModifyFaceID), rpcMessage)
	handler := new(ModifyFaceHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_ModifyFaceID hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_ModifyFaceRetID) {
		t.Error("HallMsgClassID_ModifyFaceID OnMessage 錯誤")
	}
	t.Log(ret[0].MsgData)
}
func (process *UserInitBindProcess) modifyPassword(t *testing.T) {
	msgData := new(netproto.ModifyPassword)
	msgData.NewPassword = [REDACTED]"76543210")
	msgData.OldPwd = [REDACTED]"12345678")
	rpcInfo := process.getRPCInfo()

	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) //MsgData要放入Marshal後的資料

	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ModifyPasswordID), rpcMessage)
	handler := new(ModifyPasswordHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_ModifyPasswordID hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_ModifyPasswordRetID) {
		t.Error("HallMsgClassID_ModifyPasswordID OnMessage 錯誤")
	}
	t.Log(ret[0].MsgData)
}
