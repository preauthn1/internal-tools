package msghandler

// import (
// 	"example.invalid/golang/protobuf/proto"
// )

// import (
// 	"example.invalid/ServerTeam/Servers/pkg/db"
// 	"example.invalid/ServerTeam/Servers/pkg/module"
// 	"example.invalid/ServerTeam/Servers/pkg/msg"
// 	"example.invalid/ServerTeam/Servers/pkg/netproto"
// 	"example.invalid/ServerTeam/Servers/pkg/network"
// 	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
// )

// import (
// 	"sync"
// 	"time"
// )

// // ExchangeCouponHandler 兌換禮碼
// type ExchangeCouponHandler struct {
// 	UserBlock map[int32]bool
// 	mu        sync.Mutex
// }

// // IsHook 是否有對應到這隻handler
// func (h *ExchangeCouponHandler) IsHook(bclsid int32, sclsid int32) bool {
// 	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_ExchangeCouponID)
// }

// // OnMessage 處理
// func (h *ExchangeCouponHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

// 	const PSName = "PrPs_Coupon_UpdateCouponByUserID"

// 	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
// 		return new(netproto.ExchangeCouponMsg)
// 	})

// 	if rmm == nil || rpcInfo == nil {
// 		return []*msg.Message{GetDBRetMsg(0)}
// 	}

// 	mm := rmm.(*netproto.ExchangeCouponMsg)

// 	// 檢查UserID不可在10秒內重複兌換
// 	if h.IsBlocked(rpcInfo.GetUserID()) {
// 		retMsg := new(msg.Message)
// 		retMsg.BClassID = int32(netproto.MessageBClassID_Hall)
// 		retMsg.SClassID = int32(netproto.HallMsgClassID_ExchangeCouponRetID)
// 		retMsg.MsgData = &netproto.ExchangeCouponRetMsg{
// 			RetCode: proto.Int32(0),
// 			ChvMsg:  proto.String("请勿短时间连续输入，稍后再试"),
// 		}
// 		return []*msg.Message{retMsg}
// 	}

// 	h.BlockUser(rpcInfo.GetUserID(), 10*time.Second)

// 	ps := db.NewSqlParameters()

// 	// 在這裡設定輸入參數，須和資料庫輸出欄位相同
// 	ps.AddIntInput("intUserID", rpcInfo.GetUserID())                  // 使用者ID
// 	ps.AddVarcharInput("chvCouponCode", mm.GetExchangeSequence(), 12) // 兌換碼
// 	ps.AddVarcharOutput("chvMsg", "")                                 // chvMsg的輸出

// 	// 開始準備SP字串
// 	proc := db.NewProcedure(PSName, ps)

// 	ret, err := DbUtil.ExecProc(proc)

// 	var code int32
// 	var message string

// 	if err != nil {
// 		ctx.Error("執行存儲過程[%s]錯誤, %v, %s", PSName, err, proc.GetSQLString())
// 		code = 0
// 		message = "服务器错误"
// 	} else {
// 		code = int32(ret.GetReturnValue())
// 		message = ret.GetOutputParamValue("chvMsg").(string)
// 	}

// 	retMsg := new(msg.Message)
// 	retMsg.BClassID = int32(netproto.MessageBClassID_Hall)
// 	retMsg.SClassID = int32(netproto.HallMsgClassID_ExchangeCouponRetID)
// 	retMsg.MsgData = &netproto.ExchangeCouponRetMsg{
// 		RetCode: proto.Int32(code),
// 		ChvMsg:  proto.String(message),
// 	}

// 	return []*msg.Message{retMsg}
// }

// // IsBlocked 檢查User是否被Block
// func (h *ExchangeCouponHandler) IsBlocked(userID int32) bool {
// 	h.mu.Lock()
// 	defer h.mu.Unlock()

// 	if h.UserBlock == nil {
// 		return false
// 	}

// 	if _, ok := h.UserBlock[userID]; !ok {
// 		return false
// 	}

// 	return true
// }

// // BlockUser 鎖住玩家
// func (h *ExchangeCouponHandler) BlockUser(userID int32, expire time.Duration) {
// 	h.mu.Lock()
// 	defer h.mu.Unlock()

// 	if h.UserBlock == nil {
// 		h.UserBlock = make(map[int32]bool)
// 	}

// 	h.UserBlock[userID] = true

// 	go func() {
// 		time.Sleep(expire)

// 		h.mu.Lock()
// 		defer h.mu.Unlock()

// 		delete(h.UserBlock, userID)
// 	}()
// }
