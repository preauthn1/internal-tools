package msghandler

import (
	"bytes"
	"fmt"
	"strings"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"

	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// AddPerGameLogHandler 寫遊戲記錄詳情
type AddPerGameLogHandler struct {
}

// IsHook 是否匹配
func (h *AddPerGameLogHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_DealLogID)
}

// OnMessage 處理訊息
func (h *AddPerGameLogHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.DealLog{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.DealLog)

	// var details []map[string]interface{}

	// for _, row := range mm.GetUserGameDetails() {

	// 	if row.GetTradeID() == 0 {
	// 		continue
	// 	}

	// 	d := make(map[string]interface{})

	// 	d["playerID"] = row.GetUserID()
	// 	d["tradeID"] = row.GetTradeID()
	// 	d["currentAmount"] = util.CentToMoney(row.GetCurrentScore())
	// 	d["updateAmount"] = util.CentToMoney(row.GetUpdateScore())
	// 	d["betAmount"] = util.CentToMoney(row.GetBetScore())
	// 	d["payoutAmount"] = util.CentToMoney(row.GetWinScore())
	// 	d["validAmount"] = util.CentToMoney(row.GetValidScore())
	// 	d["taxAmount"] = util.CentToMoney(row.GetTaxAmount())

	// 	details = append(details, d)
	// }

	// pergamelog := map[string]interface{}{
	// 	"serverID":           mm.GetServerID(),
	// 	"gameData":           mm.GetGameData(),
	// 	"userGameDetail":     h.convertGameRet(mm.GetUserGameRet()),
	// 	"userGameDetailLogs": details,
	// 	"gameBeginTime":      mm.GetGameBeginTime(),
	// 	"gameEndTime":        mm.GetGameEndTime(),
	// 	"taxAmount":          util.CentToMoney(mm.GetTaxAmount()),
	// 	"gameID":             mm.GetGameID(),
	// 	"winPrizeID":         "", //超端用
	// }

	// if len(details) > 0 {
	// 	delete(pergamelog, "userGameDetail")
	// }

	// var params = map[string]interface{}{
	// 	"perGameLog": pergamelog,
	// }

	if service.GameResultService.PostPerGameLogV2(rpcInfo, mm) {
		return []*msg.Message{dbutil.GetDBRetMsg(1)}
	}
	mlog.Error("寫入遊戲紀錄失敗")
	return []*msg.Message{dbutil.GetDBRetMsg(0)}
}

// 轉換細單的金額單位
func (h *AddPerGameLogHandler) convertGameRet(ret string) string {
	var buffer bytes.Buffer
	for _, row := range strings.Split(ret, ",") {
		if len(row) > 0 {
			var userID, currentAmount, updateAmount, tradeID int64
			_, err := fmt.Sscanf(row, "%d:%d:%d:%d", &userID, &currentAmount, &updateAmount, &tradeID)
			if err != nil {
				mlog.Error("convertGameRet %s, error", row)
			}
			buffer.WriteString(fmt.Sprintf("%d:%s:%s:%d,", userID, util.CentToMoney(currentAmount), util.CentToMoney(updateAmount), tradeID))
		}
	}
	return buffer.String()
}
