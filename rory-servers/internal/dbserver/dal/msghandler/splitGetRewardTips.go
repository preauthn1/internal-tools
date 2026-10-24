package msghandler

import (
	"fmt"
	"strconv"

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

// SplitGetRewardTipsHandler 裂變3說明
type SplitGetRewardTipsHandler struct{}

// IsHook ?
func (rh *SplitGetRewardTipsHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:SplitGetRewardTipsHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_SplitGetRewardTipsID)
}

// OnMessage ?
func (rh *SplitGetRewardTipsHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	// mlog.System("HANDLER名稱:SplitGetRewardTipsHandler,進入OnMessage") // DEBUG專用

	// 從CLIENT接收傳入的封包，並讀取封包裡面的訊息
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SplitGetRewardTipsMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SplitGetRewardTipsMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數，須和資料庫輸出欄位相同
	ps.AddIntInput("intUserID", mm.GetUserID()) //使用者ID

	//開始準備SP字串
	proc := db.NewProcedure("PrPs_Split_GetRewardTips", ps)
	mlog.System("HANDLER名稱:SplitGetRewardTipsHandler,SP字串:%v", proc.GetSQLString()) // DEBUG專用
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	//資料庫回傳值
	retcode := int32(dbret.GetReturnValue()) //retcode必須為1
	// mlog.System("HANDLER名稱:SplitGetRewardTipsHandler,retcode:%v", retcode) // DEBUG專用

	// 準備回傳給CLIENT的封包
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitGetRewardTipsRetID), nil)
	retdata := new(netproto.SplitGetRewardTipsRet) // 回傳封包

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("執行存儲過程 PrPs_Split_GetRewardTips 錯誤 %v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 { // 準備回傳資料
			// 第一張表
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				for i := 0; i < len(tb.Rows); i++ {
					subRetData := new(netproto.SplitGetRewardTipsList) // 產生新物件，並APPEND
					subRetData.Level = proto.Int32(int32(tb.GetValueByColName(i, "Level").(int64)))

					if i < 1 {
						subRetData.RecruitB = proto.String(strconv.Itoa(int(tb.GetValueByColName(i, "RecruitB").(int64))))
						subRetData.RecruitC = proto.String(strconv.Itoa(int(tb.GetValueByColName(i, "RecruitC").(int64))))
						subRetData.RecruitD = proto.String(strconv.Itoa(int(tb.GetValueByColName(i, "RecruitD").(int64))))

					} else {
						subRetData.RecruitB = proto.String(fmt.Sprintf("%v＋", strconv.Itoa(int(tb.GetValueByColName(i, "RecruitB").(int64)))))
						subRetData.RecruitC = proto.String(fmt.Sprintf("%v＋", strconv.Itoa(int(tb.GetValueByColName(i, "RecruitC").(int64)))))
						subRetData.RecruitD = proto.String(fmt.Sprintf("%v＋", strconv.Itoa(int(tb.GetValueByColName(i, "RecruitD").(int64)))))
					}

					subRetData.RewardB = proto.String(strconv.Itoa(int(tb.GetValueByColName(i, "RewardB").(int64))))
					subRetData.RewardC = proto.String(strconv.Itoa(int(tb.GetValueByColName(i, "RewardC").(int64))))
					subRetData.RewardD = proto.String(strconv.Itoa(int(tb.GetValueByColName(i, "RewardD").(int64))))

					retdata.SplitGetRewardTips = append(retdata.SplitGetRewardTips, subRetData)
				}
			} else {
				for i := 0; i < 5; i++ { // 要回傳五組數據
					subRetData := new(netproto.SplitGetRewardTipsList) // 產生新物件，並APPEND
					subRetData.Level = proto.Int32(0)
					subRetData.RecruitB = proto.String("0")
					subRetData.RecruitC = proto.String("0")
					subRetData.RecruitD = proto.String("0")
					subRetData.RewardB = proto.String("0")
					subRetData.RewardC = proto.String("0")
					subRetData.RewardD = proto.String("0")
					retdata.SplitGetRewardTips = append(retdata.SplitGetRewardTips, subRetData)
				}
			}
			tb = dbret[1]
			if len(tb.Rows) > 0 {
				retdata.ReturenAmount = proto.Int32(int32(tb.GetValueByColName(0, "TaxAmount").(int64)))
			} else {
				retdata.ReturenAmount = proto.Int32(0)
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
