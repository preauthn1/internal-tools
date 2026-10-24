package main

import (
	"fmt"
	"log"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

func main() {

	conf.ClearLog = true
	conf.CLogLevel = mlog.InfoLevel

	packer := network.NewHeadMsgPacker(conf.MSGPACKERBUF, true)
	evt := msgserver.NewMsgEvtChan(conf.EVTCHANSIZE)
	parser := msg.NewProtoMsgParser()

	client := msgserver.NewTCPMsgClient(packer, evt, parser, false)

	if err := client.Connect(":3000"); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Common)
			m.SClassID = int32(netproto.PlatformCommonClassID_HeartBeatID)
			if err := client.SendMessage(m); err != nil {
				log.Fatal(err)
			}
		}
	}()

	var BetOdds int64

	go func() {
		for m := range evt.MsgSig {
			fmt.Println(m)
			if m.BClassID == int32(netproto.MessageBClassID_GameRoom) && m.SClassID == int32(netproto.GameRoomClassID_ROOM_TableInfoID) {
				data := new(netproto.GEN_RoomInfo)
				proto.Unmarshal(m.MsgData.([]byte), data)
				fmt.Println(data.Odds)
				for _, table := range data.Tables {
					fmt.Println(table)
				}
			}
			if m.BClassID == int32(netproto.MessageBClassID_GameRoom) && m.SClassID == int32(netproto.GameRoomClassID_LoginRoomRetID) {
				data := new(netproto.LoginGameRoomRet)
				proto.Unmarshal(m.MsgData.([]byte), data)
				fmt.Println(data)
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GameMessageClassID_UserStateChangeID) {
				data := new(netproto.UserStateChange)
				proto.Unmarshal(m.MsgData.([]byte), data)
				fmt.Println("=== UserStateChange ===")
				fmt.Printf("StateValue: %d\n", data.User.GetStateValue())
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GEN_GameMessageClassID_GEN_GameSettingID) {
				data := new(netproto.GEN_GameSetting)
				proto.Unmarshal(m.MsgData.([]byte), data)
				fmt.Println("=== GEN_GameSetting ===")
				fmt.Printf("TableNo: %d\n", data.TableNo)
				BetOdds = data.BetOdds
				fmt.Printf("BetOdds: %d\n", data.BetOdds)
				fmt.Println("ReelStrips:")
				for _, reel := range data.ReelStrips {
					fmt.Printf("\tReel: %d, Symbols: %v\n", reel.Reel, reel.Symbols)
				}
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GEN_GameMessageClassID_GEN_RoadMapID) {
				data := new(netproto.GEN_RoadMap)
				proto.Unmarshal(m.MsgData.([]byte), data)
				fmt.Println(data)
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				data := new(netproto.GEN_GameStatus)
				proto.Unmarshal(m.MsgData.([]byte), data)
				fmt.Println(data)
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				data := new(netproto.GEN_GameResult)
				proto.Unmarshal(m.MsgData.([]byte), data)
				fmt.Println(data)
			}
		}
	}()

	{ // 登入
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_LoginRoomID)
		data := new(netproto.LoginGameRoomInfo)
		data.UserID = proto.Int32(1014943)
		data.Cer = proto.String("eyJhbGciOiJBRVMiLCJ0eXAiOiJKV1QifQ==.eyJMb2dpblVzZXJJRCI6MTAxNDk0MywiTG9naW5OYW1lIjoiZGVtby0xNjA0MzEzOTU4NjkyIiwiTWVyY2hhbnRJRCI6MTc1MDEzLCJMYW5nIjoiRU4tVVMiLCJUb2tlbkV4cGlyZWRUaW1lIjoxNjA1MTc3OTYyLCJSbmRJRCI6ImMyZmYyNjdkLWYxODEtNDViZC1iNjBhLWE4YjVlMjQyNjIxYSIsIlNpZ25LZXkiOiJmNDVlM2U2N2NjNmY2Zjk3MjhiNjFmMTIxYWI3MDUxMSIsIlRva2VuVmVyc2lvbiI6MTAwMH0=.Ob3fFSoNyH1NPXWFYDOOUe+bTDFIZ1dC3fPEXelRt9A+T/s1+7siu8S0oOir8NZjfo2BylO6qn96m5L6BPk0PDIR1Wr8Jc58sTjNVW4+JreYbME9qhEDBxc5H5Q+blghXa3YBwpZr2g6KG6K/A/eB44ljFliNkOusCJOjR0dFfpW7xlW5B7PyaX0XqT8jOzFfTWXLbg/sD/YD2X6cG/Op2jCFajdZ3/pI0VFhYvnD+yq5qlxvtAorvzoHOQ8RpZ1j5eAWFMEl0Z5XzEAyf2Kj3it7eFfkathIpG/Xzx/MzAA4NBTkFJHTIA7t4N0d81Z0sNONrfrVo3+kRIL8mL3DW9wl8I8vKR8L4am8vKYkx4Jbw32t3A6FuP9Dp+nOFwYCZCt0GQY/G/5TtbeEbk5v5NO2qfe4nL6sSIUqkIVAxnWorM6mN5UFXk3SjvL4JlXfMij76EeCLCYYiW13ScpSg==")
		// data.Cer = proto.String("eyJhbGciOiJBRVMiLCJ0eXAiOiJKV1QifQ==.eyJMb2dpblVzZXJJRCI6MTAxNDk0MywiTG9naW5OYW1lIjoiZGVtby0xNjA0MzEzOTU4NjkyIiwiTWVyY2hhbnRJRCI6MTc1MDEzLCJMYW5nIjoiRU4tVVMiLCJUb2tlbkV4cGlyZWRUaW1lIjoxNjA1MTc5NTI5LCJSbmRJRCI6IjRjZmJhZmIzLWQ3ZjAtNDJkNS04ZWZjLWQyOTU3MjU0Nzk1NCIsIlNpZ25LZXkiOiJlMzNmYjRkNjljODBkZDlhMzgzYjkxZjczMmQ0NWQ3OCIsIlRva2VuVmVyc2lvbiI6MTAwMH0=.Ob3fFSoNyH1NPXWFYDOOUe+bTDFIZ1dC3fPEXelRt9A+T/s1+7siu8S0oOir8NZjfo2BylO6qn96m5L6BPk0PDIR1Wr8Jc58sTjNVW4+JreYbME9qhEDBxc5H5Q+blghXa3YBwpZr2g6KG6K/A/eB44ljFliNkOusCJOjR0dFfpW7xlW5B7PyaX0XqT8jOzFfTWXLbg/sD/YD2X6cG/Op2jCFajdZ3/pI0VFhYvnD+yq5qlxvtAorvzoHOQ8RpZ1GXAnEIeIefy34wikDyoFl6DZYVaDLmn1UqfHLAvWNB4L54cK+wUa7H3LKsXkENB+B34HlZvIfb3nWD/EIrhSxBpPfECjBXtQ+hGFYRN9pw8+pfN4Uf+3ksPu0lvB23egjUhLkgVWVTmWvj1nd86BLnwTHIV0eW4eCC9xLDymkZfWorM6mN5UFXk3SjvL4JlXfMij76EeCLCYYiW13ScpSg==")
		data.HDCode = proto.String("")
		data.HDType = proto.Int32(0)
		m.MsgData, _ = proto.Marshal(data)
		if err := client.SendMessage(m); err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)
			// Spin
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = BetOdds
			m.MsgData, _ = proto.Marshal(data)
			if err := client.SendMessage(m); err != nil {
				log.Fatal(err)
			}
		}
	}()

	// {
	// 	m := new(msg.Message)
	// 	m.BClassID = int32(netproto.MessageBClassID_GameRoom)
	// 	m.SClassID = int32(netproto.GameRoomClassID_ROOM_Read_TableInfoID)

	// 	if err := client.SendMessage(m); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	select {}
}
