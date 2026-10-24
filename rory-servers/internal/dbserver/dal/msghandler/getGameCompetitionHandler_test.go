package msghandler

import (
	"example.invalid/golang/protobuf/proto"
)

import (	
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

import (
	"fmt"
	"testing"
)

func TestGameCompetitionHandler_IsHook(t *testing.T) {
	type args struct {
		bclsid int32
		sclsid int32
	}
	tests := []struct {
		name string
		rh   *GameCompetitionHandler
		args args
		want bool
	}{
		{
			name: "Hooked",
			rh:   &GameCompetitionHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Hall),
				sclsid: int32(netproto.HallMsgClassID_GameCompetitionID),
			},
			want: true,
		},
		{
			name: "Hooked",
			rh:   &GameCompetitionHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_GameRoom),
				sclsid: int32(netproto.HallMsgClassID_GameCompetitionID),
			},
			want: true,
		},
		{
			name: "Not Hooked BClass",
			rh:   &GameCompetitionHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Common),
				sclsid: int32(netproto.HallMsgClassID_GameCompetitionID),
			},
			want: false,
		},
		{
			name: "Not Hooked SClass",
			rh:   &GameCompetitionHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Hall),
				sclsid: int32(netproto.HallMsgClassID_GameCompetioniInfoRuleID),
			},
			want: false,
		},
		{
			name: "Not Hooked SClass",
			rh:   &GameCompetitionHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_GameRoom),
				sclsid: int32(netproto.HallMsgClassID_GameCompetioniInfoRuleID),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rh := &GameCompetitionHandler{}
			if got := rh.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("GameCompetitionHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameCompetitionHandler_OnMessage(t *testing.T) {
	initDBConnection()

	t.Run("Has Game", func(tt *testing.T) {
		rpcInfo := &netproto.RPCInfo{
			RouteServerID: proto.Int32(1),
			ConnID:        proto.String(""),
			IPAddress:     proto.String("local"),
			Cer:           proto.String(""),
			UserID:        proto.Int32(9936738),
		}

		msgData, _ := proto.Marshal(&netproto.GetGameCompetitionMsg{
			Area_ID: proto.Int32(2),
			GameID:  proto.Int32(14),
			FlagID:  proto.Int32(1),
			UserID:  proto.Int32(9936738),
		})

		rpcMessage := &rpcmsg.RPCMessage{
			RPCInfo: rpcInfo,
			MsgData: msgData,
		}

		m := &msg.Message{
			BClassID: int32(netproto.MessageBClassID_Hall),
			SClassID: int32(netproto.HallMsgClassID_GameCompetitionID),
			MsgData:  rpcMessage,
		}

		handler := new(GameCompetitionHandler)

		ctx := new(module.MouduleCtx)
		clt := new(network.ConnInfo)

		ret := handler.OnMessage(ctx, clt, m)

		if len(ret) == 0 {
			tt.Errorf("No Message Return")
		}

		if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) {
			tt.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Hall), ret[0].BClassID)
		}

		if ret[0].SClassID != int32(netproto.HallMsgClassID_GameCompetitionRetID) {
			tt.Errorf("Message SClassID want %d, but %d", int32(netproto.HallMsgClassID_GameCompetitionRetID), ret[0].SClassID)
		}

		msg, ok := ret[0].MsgData.(*netproto.GameCompetitionMsg)

		if !ok {
			tt.Errorf("Return Message is not *netproto.GameCompetitionMsg")
		}

		fmt.Printf("AreaID: %v BestScore: %v, Charge: %v, FreeTimes: %v, PiSaiTime: %v, MsgInfo: %v, UserID: %v, UserRanking: %v\n",
			msg.GetArea_ID(), msg.GetBestScore(), msg.GetCharge(), msg.GetFreeTimes(), msg.GetPiSaiTime(), msg.GetMsgInfo(), msg.GetUserID(), msg.GetUserRanking())

		for _, row := range msg.GetLeaderboard() {
			fmt.Printf("Leaderboard => Area_ID: %v, TitleName: %v\n", row.GetArea_ID(), row.GetTitleName())
			for _, _row := range row.GetLeaderboardList() {
				fmt.Printf("LeaderboardList => Area_ID: %v, UserID: %v, Score: %v, SortID: %v\n", _row.GetArea_ID(), _row.GetUserID(), _row.GetScore(), _row.GetSortID())
			}
		}
	})

	t.Run("Game Not Found", func(tt *testing.T) {
		rpcInfo := &netproto.RPCInfo{
			RouteServerID: proto.Int32(1),
			ConnID:        proto.String(""),
			IPAddress:     proto.String("local"),
			Cer:           proto.String(""),
			UserID:        proto.Int32(9936738),
		}

		msgData, _ := proto.Marshal(&netproto.GetGameCompetitionMsg{
			Area_ID: proto.Int32(2),
			GameID:  proto.Int32(25),
			FlagID:  proto.Int32(1),
			UserID:  proto.Int32(9936738),
		})

		rpcMessage := &rpcmsg.RPCMessage{
			RPCInfo: rpcInfo,
			MsgData: msgData,
		}

		m := &msg.Message{
			BClassID: int32(netproto.MessageBClassID_Hall),
			SClassID: int32(netproto.HallMsgClassID_GameCompetitionID),
			MsgData:  rpcMessage,
		}

		handler := new(GameCompetitionHandler)

		ctx := new(module.MouduleCtx)
		clt := new(network.ConnInfo)

		ret := handler.OnMessage(ctx, clt, m)

		if len(ret) == 0 {
			t.Errorf("No Message Return")
		}

		if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) {
			t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Hall), ret[0].BClassID)
		}

		if ret[0].SClassID != int32(netproto.HallMsgClassID_GameCompetitionRetID) {
			t.Errorf("Message SClassID want %d, but %d", int32(netproto.HallMsgClassID_GameCompetitionRetID), ret[0].SClassID)
		}

		msg, ok := ret[0].MsgData.(*netproto.GameCompetitionMsg)

		if !ok {
			t.Errorf("Return Message is not *netproto.GameCompetitionMsg")
		}

		fmt.Printf("AreaID: %v BestScore: %v, Charge: %v, FreeTimes: %v, PiSaiTime: %v, MsgInfo: %v, UserID: %v, UserRanking: %v\n",
			msg.GetArea_ID(), msg.GetBestScore(), msg.GetCharge(), msg.GetFreeTimes(), msg.GetPiSaiTime(), msg.GetMsgInfo(), msg.GetUserID(), msg.GetUserRanking())

		for _, row := range msg.GetLeaderboard() {
			fmt.Printf("Leaderboard => Area_ID: %v, TitleName: %v\n", row.GetArea_ID(), row.GetTitleName())
			for _, _row := range row.GetLeaderboardList() {
				fmt.Printf("LeaderboardList => Area_ID: %v, UserID: %v, Score: %v, SortID: %v\n", _row.GetArea_ID(), _row.GetUserID(), _row.GetScore(), _row.GetSortID())
			}
		}

	})

}
