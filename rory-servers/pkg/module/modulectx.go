package module

import (
	"fmt"
	"os"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// MouduleCtx 模塊運行上下文
type MouduleCtx struct {
	BackendAPI     string
	APICertificate string
	EvtDispatcher  *event.Dispatcher    // 事件分派器
	Logger         mlog.Logger          // 日誌記錄器
	Config         *config.Group        // 配置信息
	Server         *msgserver.MsgServer // 服務器
}

// Printf 向控制台輸出
func (ctx *MouduleCtx) Printf(format string, args ...interface{}) {
	fmt.Printf(format+"\r\n", args...)
}

// Println 向控制台輸出
func (ctx *MouduleCtx) Println(v ...interface{}) {
	fmt.Println(v...)
}

// WriteLog 寫日誌
func (ctx *MouduleCtx) WriteLog(lvl int, format string, a ...interface{}) {
	if lvl < conf.CLogLevel {
		return
	}

	// mlog.WriteLog(lvl int, format string, a ...interface{})
	// ctx.Logger.WriteLog(lvl, format, a...)
	mlog.WriteLog(lvl, format, a...)

	if lvl >= int(mlog.FatalLevel) {
		os.Exit(1)
	}
}

// Trace 調試信息
func (ctx *MouduleCtx) Trace(format string, a ...interface{}) {
	ctx.WriteLog(int(mlog.TraceLevel), format, a...)
}

// Debug 調試信息
func (ctx *MouduleCtx) Debug(format string, a ...interface{}) {
	ctx.WriteLog(int(mlog.DebugLevel), format, a...)
}

// Info 日志信息
func (ctx *MouduleCtx) Info(format string, a ...interface{}) {
	ctx.WriteLog(int(mlog.InfoLevel), format, a...)
}

// Error 錯誤信息
func (ctx *MouduleCtx) Error(format string, a ...interface{}) {
	ctx.WriteLog(int(mlog.ErrorLevel), format, a...)
}

// Fatal 崩潰信息
func (ctx *MouduleCtx) Fatal(format string, a ...interface{}) {
	ctx.WriteLog(int(mlog.FatalLevel), format, a...)
}

// GetConfig 獲取配置信息
func (ctx *MouduleCtx) GetConfig(mname string) *config.Config {
	cfg := ctx.Config.GetConfig(mname)

	return cfg
}

// GetConfigSection 獲取配置節
func (ctx *MouduleCtx) GetConfigSection(mname, secname string) config.Section {
	cfg := ctx.GetConfig(mname)
	if cfg == nil {
		return nil
	}

	return cfg.Sections[secname]
}

// GetConfigData 獲取配置數據
func (ctx *MouduleCtx) GetConfigData(mname, secname, key string) string {
	cfgsex := ctx.GetConfigSection(mname, secname)

	if cfgsex == nil {
		return ""
	}

	return cfgsex[key]
}

// TriggerEvent 觸發事件
func (ctx *MouduleCtx) TriggerEvent(evt *event.Event) *event.Ret {
	return ctx.EvtDispatcher.OnEvent(evt)
}

// TriggerAsyncEventCallback 觸發異步事件
func (ctx *MouduleCtx) TriggerAsyncEventCallback(evt *event.Event, callback event.CallbackFun) {
	ctx.EvtDispatcher.OnAsyncEvent(evt, callback)
}

// TriggerAsyncEvent 觸發異步事件
func (ctx *MouduleCtx) TriggerAsyncEvent(evt *event.Event) {
	ctx.EvtDispatcher.OnAsyncEvent(evt, nil)
}

// RegisterMsgParser 註冊消息解析器
func (ctx *MouduleCtx) RegisterMsgParser(bclsID int32, sclsid int32, bcp msg.Parser) {
	ctx.Server.GetMsgParserGroup().RegisterParser(bclsID, sclsid, bcp)
}

// RegisterMsgParserFunc 註冊單消息解析器函數
func (ctx *MouduleCtx) RegisterMsgParserFunc(bclsID int32, sclsid int32, f func() proto.Message) {
	umf := func(data []byte) (*msg.Message, error) {
		m1 := new(msg.Message)
		m1.BClassID = bclsID
		m1.SClassID = sclsid

		mm := f()

		err := proto.Unmarshal(data[8:], mm)
		if err != nil {
			return nil, err
		}
		m1.MsgData = mm

		return m1, nil
	}

	mf := func(msg *msg.Message) ([]byte, error) {
		b1 := util.IntToBytes(int32(bclsID), true)
		b2 := util.IntToBytes(int32(sclsid), true)

		d := msg.MsgData.(proto.Message)
		b3, err := proto.Marshal(d)
		if err != nil {
			return nil, err
		}

		return append(append(b1, b2...), b3...), nil
	}

	ctx.Server.GetMsgParserGroup().RegisterParser(bclsID, sclsid, msg.NewCommonMsgParser(mf, umf))
}

// RegisterEventHandler 註冊事件處理器
func (ctx *MouduleCtx) RegisterEventHandler(name string, f func(evt *event.Event, ret *event.Ret) bool, evts ...string) {
	h := event.NewCommonHandler(name, f, evts...)

	ctx.EvtDispatcher.PushHandler(h)
}

// OnConnect 客戶端連接
func (ctx *MouduleCtx) OnConnect(f func(clt *network.ConnInfo) bool) {
	h := event.NewCommonHandler("conn", func(evt *event.Event, evtret *event.Ret) bool {
		clt := msgserver.GetClientInfo(evt.Args)
		return f(clt)
	}, msgserver.EventClientConnect)

	ctx.EvtDispatcher.PushHandler(h)
}

// OnDisConnect 客戶端斷開連接
func (ctx *MouduleCtx) OnDisConnect(f func(clt *network.ConnInfo) bool) {
	h := event.NewCommonHandler("disconn", func(evt *event.Event, evtret *event.Ret) bool {
		clt := msgserver.GetClientInfo(evt.Args)

		return f(clt)
	}, msgserver.EventClientDisconnect)

	ctx.EvtDispatcher.PushHandler(h)
}

// OnMessage 客戶端消息
func (ctx *MouduleCtx) OnMessage(f func(clt *network.ConnInfo, msg *msg.Message) bool) {
	h := event.NewCommonHandler("msg", func(evt *event.Event, evtret *event.Ret) bool {
		clt := msgserver.GetClientInfo(evt.Args)
		m := msgserver.GetClientMsg(evt.Args)

		return f(clt, m)
	}, msgserver.EventClientMessage)

	ctx.EvtDispatcher.PushHandler(h)
}
