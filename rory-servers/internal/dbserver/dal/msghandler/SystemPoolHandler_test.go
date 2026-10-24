package msghandler

import (
	"encoding/json"
	"fmt"
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"testing"
)

// 遊戲輸贏結果 DB 儲存資料
type GameDataDB struct {
	WeightTableNoNow int   // 目前所使用的莊家權重表
	DealerNotOKCount int32 // 莊家不成局
	DealerWinCount   int32 // 莊家勝利
	DealerLoseCount  int32 // 玩家比牌勝利（莊家成局）
	GameRound        int32 // 進行局數
}

func TestSystemPoolHandler_IsHook(t *testing.T) {
	type args struct {
		bclsid int32
		sclsid int32
	}
	tests := []struct {
		name string
		h    *SystemPoolHandler
		args args
		want bool
	}{
		{
			name: "Hooked GetSysPoolWinLose",
			h:    &SystemPoolHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.GameMessageClassID_SuperToolGetSysPoolWinLose),
			},
			want: true,
		},
		{
			name: "Hooked UpdateSysPoolWinLose",
			h:    &SystemPoolHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.GameMessageClassID_SuperToolUpdateSysPoolWinLose),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SystemPoolHandler{}
			if got := h.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("SystemPoolHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 要求 遊戲輸贏結果
func TestSystemPoolHandler_OnMessage_GetSysPoolWinLose(t *testing.T) {
	initDBConnection()

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(9936738),
	}

	msgData, _ := proto.Marshal(&netproto.SuperTool_Msg_GetSysPoolWinLose{
		GameID: proto.Int32(2001),
		FlagID: proto.Int32(1),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.GameMessageClassID_SuperToolGetSysPoolWinLose),
		MsgData:  rpcMessage,
	}

	handler := new(SystemPoolHandler)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.GameMessageClassID_SuperToolGetSysPoolWinLoseRet) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.GameMessageClassID_SuperToolGetSysPoolWinLoseRet), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.SuperTool_Msg_GetSysPoolWinLoseRet)

	if !ok {
		t.Errorf("Return Message is not *netproto.SuperTool_Msg_GetSysPoolWinLoseRet")
	}

	fmt.Printf("GameID: %d\n", msg.GetGameID())           // 遊戲編號
	fmt.Printf("FlagID: %d\n", msg.GetFlagID())           // 房間編號
	fmt.Printf("PoolAmount: %d\n", msg.GetPoolAmount())   // 水線
	fmt.Printf("WinLine: %d\n", msg.GetWinLine())         // 贏錢線
	// 遊戲輸贏結果
	pGameData := new(GameDataDB)
	err := json.Unmarshal([]byte(msg.GetWinLoseData()), pGameData)
	if err == nil {
		fmt.Printf("WeightTableNoNow:%d, DealerNotOKCount:%d, DealerWinCount:%d, DealerLoseCount:%d, GameRound:%d\n",
			pGameData.WeightTableNoNow, pGameData.DealerNotOKCount, pGameData.DealerWinCount, pGameData.DealerLoseCount, pGameData.GameRound)
	}
}

// 更新 遊戲輸贏結果
func TestSystemPoolHandler_OnMessage_UpdateSysPoolWinLose(t *testing.T) {
	initDBConnection()

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(9936738),
	}

	gameData := &GameDataDB{
		WeightTableNoNow: 1,
		DealerNotOKCount: 2,
		DealerWinCount:   1,
		DealerLoseCount:  2,
		GameRound:        5,
	}

	byGameData, _ := json.Marshal(gameData)

	fmt.Printf("strGameData: %s\n", byGameData)

	msgData, _ := proto.Marshal(&netproto.SuperTool_Msg_UpdateSysPoolWinLose{
		GameID:      proto.Int32(2001),
		FlagID:      proto.Int32(1),
		PoolAmount:  proto.Int64(1500100),
		WinLine:     proto.Int64(1500000),
		WinLoseData: proto.String(string(byGameData)),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.GameMessageClassID_SuperToolUpdateSysPoolWinLose),
		MsgData:  rpcMessage,
	}

	handler := new(SystemPoolHandler)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}
}
