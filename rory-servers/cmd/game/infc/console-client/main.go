package main

import (
	"bufio"
	"context"
	"flag"
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

var u1 = &user{uid: 1015669, cer: "eyJhbGciOiJBRVMiLCJ0eXAiOiJKV1QifQ==.eyJMb2dpblVzZXJJRCI6MTAxNTY2OSwiTG9naW5OYW1lIjoiZGVtby0xNjA3OTMxODMwMjA5IiwiTWVyY2hhbnRJRCI6MTc1MDEzLCJMYW5nIjoiRU4tVVMiLCJUb2tlbkV4cGlyZWRUaW1lIjoxNjA4Nzk1ODg4LCJSbmRJRCI6IjJiOTRkM2FmLTU4ZTAtNDliNS1iNzJlLWIwZGYwMzZhODc0OSIsIlNpZ25LZXkiOiJhOTQ1MWQ4YWUyYjQwNjg1MGVkZGZlNjAxNzUyNzM0ZSIsIlRva2VuVmVyc2lvbiI6MTAwMH0=.Ob3fFSoNyH1NPXWFYDOOUe+bTDFIZ1dC3fPEXelRt9A+T/s1+7siu8S0oOir8NZj3tupr1G8+hutlV9/dQQee1747JjH+4BYIeFjFMVNJxuYbME9qhEDBxc5H5Q+blgh0QfwQmwlkR+yQQJImQ+0s44ljFliNkOusCJOjR0dFfpW7xlW5B7PyaX0XqT8jOzFfTWXLbg/sD/YD2X6cG/Op2jCFajdZ3/pI0VFhYvnD+zdbNzYwLmZO8GS66qlUPP2z54c2aQZJSNDaBE1dfEapYok4tWtzl+fcpXsDJB3EFt7knKhTW3T+FGo8LiyBB5Dv2EbLxEN+LEIhPgyqQHHBUHpOrcICNEYgA0nMINhjHbEOTJC4+lKM0GK7mrcgbxuf65p+ZiLLUG5Gda1DmxVf/HA0LiDOIRadHjDLyuy0b/WorM6mN5UFXk3SjvL4JlXfMij76EeCLCYYiW13ScpSg=="}
var u2 = &user{uid: 1015671, cer: "eyJhbGciOiJBRVMiLCJ0eXAiOiJKV1QifQ==.eyJMb2dpblVzZXJJRCI6MTAxNTY3MSwiTG9naW5OYW1lIjoiZGVtby0xNjA3OTMxOTc1Mzc1IiwiTWVyY2hhbnRJRCI6MTc1MDEzLCJMYW5nIjoiRU4tVVMiLCJUb2tlbkV4cGlyZWRUaW1lIjoxNjA4Nzk1OTc4LCJSbmRJRCI6IjQwNGYzYTZkLTQyZjQtNDRiYS05ZGRkLTAwOTIyZGJkMGEyZCIsIlNpZ25LZXkiOiI0ZGIxNWFlNTQ5ZWRmN2U5ZWNlNGJjNTc1YmNjMmQ0MCIsIlRva2VuVmVyc2lvbiI6MTAwMH0=.Ob3fFSoNyH1NPXWFYDOOUe+bTDFIZ1dC3fPEXelRt9A+T/s1+7siu8S0oOir8NZj3tupr1G8+hutlV9/dQQee3HA+3T9V45Y3kZdqM2vuQGYbME9qhEDBxc5H5Q+blghSvxRbcv8LalwV4adVu+1hI4ljFliNkOusCJOjR0dFfpW7xlW5B7PyaX0XqT8jOzFfTWXLbg/sD/YD2X6cG/Op2jCFajdZ3/pI0VFhYvnD+zdbNzYwLmZO8GS66qlUPP2E0+uB6gdjkz8ZoLaWi5MmqjL5ocquvpXNuhZMDGEnbXsRQVIX0bkkFsIPuUAqocHRrVe4XC7mXt3kswaMtPZ3++NNVN8z/zcThwZApO63n0KC9ja5zSVN/6s8pZDsy1p/73ifgtGfn1gO+TOxj9OAP/IdVft/P1FUXDm6KjSyPvWorM6mN5UFXk3SjvL4JlXfMij76EeCLCYYiW13ScpSg=="}
var pid int

func init() {
	flag.IntVar(&pid, "p", 0, "user1 ort user2")
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

func main() {

	var u *user

	flag.Parse()

	switch pid {
	case 1:
		u = u1
	case 2:
		u = u2
	}

	if u == nil {
		log.Fatal("Not right user")
	}

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
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INFC_GameMessageClassID_INFC_GameSettingID) {
				data := new(netproto.INFC_GameSetting)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("=== 遊戲設定(%v) ===", m)
				log.Printf("桌號: %d\n", data.TableNo)
				log.Printf("每次扣除: %d\n", data.Bet)
				log.Printf("開始時間: %s\n", time.Unix(data.StartTime, 0).Format(time.RFC3339))
				log.Printf("玩家: (%d)\n", len(data.Players))
				for _, p := range data.Players {
					log.Printf(" - %s(%d) %d\n", p.UserName, p.UserID, p.UserScore)
				}
				log.Printf("====================")
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INFC_GameMessageClassID_INFC_ProtoSettingID) {
				data := new(netproto.INFC_ProtoSetting)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("=== Proto設定(%v) ===", m)
				// log.Printf("Patterns: %d\n", len(data.Patterns))
				// for _, row := range data.Patterns {
				// 	log.Printf(" - PatternID: %d, StageGroupID: %d, IsLoop: %d\n", row.PatternID, row.StageGroupID, row.IsLoop)
				// }
				// log.Printf("FruitAttributes: %d\n", len(data.FruitAttributes))
				// for _, row := range data.FruitAttributes {
				// 	log.Printf(" - ID: %d, GroupID: %d, FruitName: %s, FruitWeight: %.2f\n", row.ID, row.GroupID, row.FruitName, row.FruitWeight)
				// }
				// log.Printf("FruitEmitters: %d\n", len(data.FruitEmitters))
				// for _, row := range data.FruitEmitters {
				// 	log.Printf(" - ID: %d, EmitterGroupID: %d, EmitterID: %d, AngleMin: %.2f, AngleMax: %.2f, ForceMin: %.2f, ForceMax: %.2f, IntervalMin: %.2f, IntervalMax: %.2f, FruitID: %.2f\n",
				// 		row.ID, row.EmitterGroupID, row.EmitterID, row.AngleMin, row.AngleMax, row.ForceMin, row.ForceMax, row.IntervalMin, row.IntervalMax, row.FruitID)
				// }
				// log.Printf("StageGroups: %d\n", len(data.StageGroups))
				// for _, row := range data.StageGroups {
				// 	log.Printf(" - StageGroupID: %d, EmitterGroupID: %d, OrderID: %d, GroupDelayTime: %d\n", row.StageGroupID, row.EmitterGroupID, row.OrderID, row.GroupDelayTime)
				// }
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INFC_GameMessageClassID_INFC_BetStepID) {
				log.Printf("=== 現在扣錢(%v) ===", m)
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INFC_GameMessageClassID_INFC_GameStartID) {
				log.Printf("=== 請開始刷水果(%v) ===", m)
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INFC_GameMessageClassID_INFC_GameRankID) {
				data := new(netproto.INFC_GameRank)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("=== 更新排行榜(%v) ===", m)
				for _, p := range data.Players {
					log.Printf(" - 玩家[%d](%s) %d, 排名: %d, 分數: %d, 狀態: %d, 失誤: %d", p.UserID, p.UserName, p.UserScore, p.Rank, p.Score, p.State, p.Miss)
				}
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INFC_GameMessageClassID_INFC_GameResultID) {
				data := new(netproto.INFC_GameResult)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("=== 遊戲結果(%v) ===", m)
				log.Printf("獎金: %d\n", data.WinMoney)
				log.Printf("贏家[%d](%s) %d, 排名: %d, 分數: %d, 狀態: %d, 失誤: %d", data.Winner.UserID, data.Winner.UserName, data.Winner.UserScore, data.Winner.Rank, data.Winner.Score, data.Winner.State, data.Winner.Miss)
				log.Printf("排行榜: %d\n", len(data.Rank.Players))
				for _, p := range data.Rank.Players {
					log.Printf(" - 玩家[%d](%s) %d, 排名: %d, 分數: %d, 狀態: %d, 失誤: %d", p.UserID, p.UserName, p.UserScore, p.Rank, p.Score, p.State, p.Miss)
				}
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
			case "queue": // 排隊
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_GameRoom)
				m.SClassID = int32(netproto.GameRoomClassID_UserQueueID)
				client.SendMessage(m)
			case "score":
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
				m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_IncScoreID)
				data := new(netproto.INFC_IncScore)
				data.Score = score
				m.MsgData, _ = proto.Marshal(data)
				client.SendMessage(m)
			case "miss": // 傳送Miss封包
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_MissID)
				client.SendMessage(m)
			case "boom": // 傳送Boom封包
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_BoomID)
				client.SendMessage(m)
			case "giveup": // 傳送Giveup封包
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_GiveupID)
				client.SendMessage(m)
			}
		}
	}()

	{ // 登入
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_LoginRoomID)
		data := new(netproto.LoginGameRoomInfo)
		data.UserID = proto.Int32(u.uid)
		data.Cer = proto.String(u.cer)
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
