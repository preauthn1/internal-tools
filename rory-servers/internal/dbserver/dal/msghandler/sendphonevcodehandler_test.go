package msghandler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.invalid/golang/mock/gomock"
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

func TestSendPhoneVCode(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)
	mockretset := db.RecordSet{}

	dataSetResult1 := db.RecordTable{}
	dataSetResult1.Cols = append(dataSetResult1.Cols, "tnyCountdown")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "chvErrMsg")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "intCodeID")
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{int64(180), "发送验证码成功", int64(5)})
	mockretset = append(mockretset, &dataSetResult1)

	dataSetResult2 := db.RecordTable{}
	dataSetResult2.Rows = append(dataSetResult2.Rows, db.RecordRow{int64(1)})
	mockretset = append(mockretset, &dataSetResult2)

	mockDbUtil.EXPECT().ExecProc(gomock.Any()).Return(mockretset, nil)
	dbutil.DBUtil = mockDbUtil

	//mock http
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()
	dbutil.SMSWebAPI = strings.Replace(ts.URL, "http://", "", -1)

	process := new(UserInitBindProcess)
	//Client傳入資料的協議號
	msgData := new(netproto.SendPhoneVCode)
	msgData.Tel = proto.String(process.Tel)
	msgData.CodeType = proto.Int32(1)

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
