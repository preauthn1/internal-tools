package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"example.invalid/dustin/go-humanize"
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

type user struct {
	uid int32
	cer string
}

// UserStateNameByValue ?
func UserStateNameByValue(v int) string {
	// 用戶狀態
	switch v {
	default:
		return ""
	case gamelogic.UserStateFree:
		return gamelogic.StateNameFree
	case gamelogic.UserStateQueue:
		return gamelogic.StateNameQueue
	case gamelogic.UserStateSit:
		return gamelogic.StateNameSit
	case gamelogic.UserStateReady:
		return gamelogic.StateNameReady
	case gamelogic.UserStateGame:
		return gamelogic.StateNameGame
	case gamelogic.UserStateOffline:
		return gamelogic.StateNameOffline
	}
}

var u1 = &user{uid: 1015909, cer: "eyJhbGciOiJBRVMiLCJ0eXAiOiJKV1QifQ==.eyJMb2dpblVzZXJJRCI6MTAxNTkwOSwiTG9naW5OYW1lIjoiZGVtby0xNjEyMDk4Njg0NjEwIiwiTWVyY2hhbnRJRCI6MTc1MDEzLCJMYW5nIjoiRU4tVVMiLCJUb2tlbkV4cGlyZWRUaW1lIjoxNjEyOTYyNjg5LCJSbmRJRCI6IjljMGUxZWUzLWRhZDktNDY3Yy04OGJlLWU5ZjY3NWExNmY4ZiIsIlNlc3Npb25Ub2tlbiI6ImVhMzUxZWJhLTIyMmMtNGI5Zi04MzU2LTIyMzMwOGFiYWI3MyIsIlNpZ25LZXkiOiIyNmQyNDNhNWQzODE0OTkyNTZiMzYzOGM5OTYyOTAyMyIsIlRva2VuVmVyc2lvbiI6MTAwMH0=.Ob3fFSoNyH1NPXWFYDOOUe+bTDFIZ1dC3fPEXelRt9A+T/s1+7siu8S0oOir8NZj+y/j+cZglfyueeAlO6OxOVhT5mlcBNYUjPJHeg60m/udJBYAz9rUFGBfa+t5rSIS5vjhjHWk5sXzOq6rv9cjzo4ljFliNkOusCJOjR0dFfpW7xlW5B7PyaX0XqT8jOzFfTWXLbg/sD/YD2X6cG/Op2jCFajdZ3/pI0VFhYvnD+xwkImZSqPczNB3l7uxmVIS3JDJ/47U7Xl01XIll17nXxiCtuHw4ruATjCazNaKSf1hho6vFbQckz/HoH0l+HgD8TWQRDz9s30tn59nuG2wsyxI4YU9tq+oVz0O3hicsWvSuyoMkAM0l5+7pNqyjuqgkh4ceHqHLvIiRB9yr/eL4wr7M/poYPK+XnTKCqHddlVbU0F/V5XFaVN9FcLc0Em8YV9Eyh6Kaq3+JbASZcdSSnoFYBs89G74KaqjlA4HDb6KvdszgaotF7hEBi8F53nyiAT1DOpVEk+SECSfzEHI794MaLtMAeE+rtw72hrpMsEE+/JGqq2tacREC1f3g5Di"}

func main() {

	ctx, cancel := context.WithCancel(context.Background())

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

	logined := make(chan bool, 1)
	go func() {
		for m := range evt.MsgSig {
			// fmt.Println(m)
			if m.BClassID == int32(netproto.MessageBClassID_Common) && m.SClassID == int32(netproto.PlatformCommonClassID_TipMessageID) {
				data := new(netproto.TipMessage)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("你被踢了... Message: %v, Type: %d, Param: %d", data.GetMessage(), data.GetMsgType(), data.GetParam())
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_GameRoom) && m.SClassID == int32(netproto.GameRoomClassID_LoginRoomRetID) {
				data := new(netproto.LoginGameRoomRet)
				proto.Unmarshal(m.MsgData.([]byte), data)
				if data.GetCode() != 1 {
					log.Fatalf("登入失敗 %s", data.GetMessage())
				}
				logined <- true
				log.Printf("%d(%s) 登入成功", data.GetUserData().GetUserID(), data.GetUserData().GetNickName())
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GameMessageClassID_UserEnterID) {
				data := new(netproto.UserEnter)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("玩家[%d] 進入 桌子[%d] 座位[%d] 狀態[%s] 金幣[%s]", data.User.GetUserID(), data.User.GetTableNo(), data.User.GetSeatNo(), UserStateNameByValue(int(data.User.GetStateValue())), humanize.Commaf(float64(data.User.GetScore())/100))
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GameMessageClassID_GameStartID) {
				log.Printf("=== 遊戲開始 ===")
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GameMessageClassID_GameEndID) {
				log.Printf("=== 遊戲結束 ===")
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GameMessageClassID_UserStateChangeID) {
				data := new(netproto.UserStateChange)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("玩家[%d] 狀態[%s] 金幣[%s]", data.User.GetUserID(), UserStateNameByValue(int(data.User.GetStateValue())), humanize.Commaf(float64(data.User.GetScore())/100))
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_BetStatusID) {
				data := new(netproto.Slot_MSG_BetStatus)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("玩家[%d] 下注結果[%d] 花費[%s]", u1.uid, data.GetRet(), humanize.Commaf(float64(data.GetMoney())/100))
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_GameResultID) {
				data := new(netproto.Slot_MSG_GameResult)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("玩家[%d] 遊戲結果 連線贏[%s] Jackpot贏[%s] 最後[%s]", u1.uid, humanize.Commaf(float64(data.GetWinMoney())/100), humanize.Commaf(float64(data.GetJackpotMoney())/100), humanize.Commaf(float64(data.GetFinalMoney())/100))
				continue
			}
		}
	}()

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				log.Println(rec)
			}
		}()
		for {
			reader := bufio.NewReader(os.Stdin)
			data, _, _ := reader.ReadLine()
			sline := string(data)

			sline = strings.ToLower(sline)
			sline = strings.TrimSpace(sline)

			sarr := strings.Split(sline, " ")
			cname := sarr[:1][0]
			cargs := sarr[1:]

			log.Printf("收到指令: %s, 參數: %s", cname, cargs)

			switch cname {
			default:
			case "exit":
				cancel()
			case "spin":
				if len(cargs) == 0 {
					mlog.Error("分數輸入錯誤, 請輸入一個數字")
					continue
				}
				score, err := strconv.ParseInt(cargs[0], 0, 64)
				if err != nil {
					mlog.Error("分數輸入錯誤 %v", err)
					continue
				}
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)
				data := new(netproto.Slot_MSG_Spin)
				data.Money = proto.Int64(score)
				m.MsgData, _ = proto.Marshal(data)
				client.SendMessage(m)
			}
		}
	}()

	{ // 登入
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_LoginRoomID)
		data := new(netproto.LoginGameRoomInfo)
		data.UserID = proto.Int32(u1.uid)
		data.Cer = proto.String(u1.cer)
		data.HDCode = proto.String("")
		data.HDType = proto.Int32(0)
		m.MsgData, _ = proto.Marshal(data)
		if err := client.SendMessage(m); err != nil {
			log.Fatal(err)
		}
	}

	<-logined

	select {
	case <-ctx.Done():
		log.Printf("結束")
	}
}
