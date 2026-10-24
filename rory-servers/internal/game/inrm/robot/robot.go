package robot

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"example.invalid/dustin/go-humanize"
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/game/inrm"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
)

func pStateName(s netproto.INRM_PState) string {
	switch s {
	default:
		return "?"
	case netproto.INRM_PState_INRM_PState_Wait:
		return "等待中"
	case netproto.INRM_PState_INRM_PState_GameV1:
		return "遊戲中V1"
	case netproto.INRM_PState_INRM_PState_GameV2:
		return "遊戲中V2"
	case netproto.INRM_PState_INRM_PState_FirstGiveup:
		return "棄牌-"
	case netproto.INRM_PState_INRM_PState_Giveup:
		return "棄牌+"
	case netproto.INRM_PState_INRM_PState_Bluff:
		return "詐胡"
	case netproto.INRM_PState_INRM_PState_Victory:
		return "勝利"
	case netproto.INRM_PState_INRM_PState_Loss:
		return "失敗"
	}
}

func printPlayerAll(p *netproto.INRM_Player, indent string) {
	var hand []*poker.Deck
	for _, d := range p.Hands.Decks {
		hand = append(hand, poker.ReversePokerDeck(d.Deck))
	}
	log.Printf("%s玩家[%d][%s] 點數[%d] 贏[%s] 稅[%s] 更新[%s] 手牌: %v", indent, p.ID,
		pStateName(p.State),
		p.LossPoint,
		humanize.Commaf(float64(p.Win)/100),
		humanize.Commaf(float64(p.Tax)/100),
		humanize.Commaf(float64(p.UpdateScore)/100),
		hand)
}

// MessagePureRobot 純封包溝通的機器人
type MessagePureRobot struct {
	client    *msgserver.TCPMsgClient
	evt       *msgserver.MsgEvtChan
	refresher inrm.Refresher

	addr        string
	uid         int32
	cer         string
	actWaitTime time.Duration

	Setting   *netproto.INRM_GameSetting
	GameState *netproto.INRM_GameState

	deck      *poker.Deck
	lastThrow poker.Card
}

// NewMessagePureRobot 純封包溝通的機器人
func NewMessagePureRobot(addr string, uid int32, cer string) *MessagePureRobot {

	packer := network.NewHeadMsgPacker(conf.MSGPACKERBUF, true)
	evt := msgserver.NewMsgEvtChan(conf.EVTCHANSIZE)
	parser := msg.NewProtoMsgParser()

	client := msgserver.NewTCPMsgClient(packer, evt, parser, false)

	r := new(MessagePureRobot)
	r.client = client
	r.evt = evt
	r.refresher = inrm.DefaultRefresher()

	r.addr = addr
	r.uid = uid
	r.cer = cer

	r.actWaitTime = 1 * time.Second

	return r
}

func (r *MessagePureRobot) _refreshHands(hands []*poker.Deck) {

	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
	data := new(netproto.Poker_Hand)
	for _, d := range hands {
		data.Decks = append(data.Decks, poker.TransePokerDeck(d))
	}

	m.MsgData, _ = proto.Marshal(data)

	r.client.SendMessage(m)
}

func (r *MessagePureRobot) refreshHands() {
	_, _, hands := r.refresher.Refresh(int(r.Setting.HandCardLength), poker.Rank(r.GameState.RandJoker), r.deck)
	r._refreshHands(hands)
}

func (r *MessagePureRobot) popCard() {
	time.Sleep(r.actWaitTime)

	// 檢查開牌區的牌加入可不可以胡
	d := poker.NewEmptyDeck(r.deck.Cards...)
	d.AppendCard(poker.ReversePokerCard(r.GameState.FirstOpenDeck))

	// 可以胡就抽開牌區，不能就抽蓋牌區
	if ok, _, _ := r.refresher.Refresh(int(r.Setting.HandCardLength), poker.Rank(r.GameState.RandJoker), d); ok {
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_Game)
		m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopOpenDeckID)
		r.client.SendMessage(m)
		return
	}

	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopCloseDeckID)
	r.client.SendMessage(m)
}

func (r *MessagePureRobot) throwOrVictory() {
	time.Sleep(r.actWaitTime)

	ok, throw, hands := r.refresher.Refresh(int(r.Setting.HandCardLength), poker.Rank(r.GameState.RandJoker), r.deck)

	if ok {
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_Game)
		m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_DeclareID)
		data := new(netproto.INRM_Declare)
		data.Throw = poker.TransePokerCard(throw)
		data.Hands = new(netproto.Poker_Hand)
		for _, d := range hands {
			data.Hands.Decks = append(data.Hands.Decks, poker.TransePokerDeck(d))
		}
		m.MsgData, _ = proto.Marshal(data)

		r.client.SendMessage(m)
		return
	}

	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_ThrowCardID)
	data := new(netproto.INRM_ThrowCard)
	data.Throw = poker.TransePokerCard(throw)
	m.MsgData, _ = proto.Marshal(data)

	r.client.SendMessage(m)

	r.lastThrow = throw

	r.deck.PickCard(throw)

	r._refreshHands(hands)
}

func (r *MessagePureRobot) ready() {
	time.Sleep(r.actWaitTime)

	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_GameRoom)
	m.SClassID = int32(netproto.GameRoomClassID_GameReadyID)
	r.client.SendMessage(m)
}

// Run 運行
//  登入
//  持續打心跳包
//  遊戲操作
func (r *MessagePureRobot) Run() error {

	if err := r.client.Connect(r.addr); err != nil {
		return err
	}
	defer r.client.Close()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Common)
			m.SClassID = int32(netproto.PlatformCommonClassID_HeartBeatID)
			if err := r.client.SendMessage(m); err != nil {
				log.Fatal(err)
			}
		}
	}()

	logined := make(chan bool, 1)
	go func() {
		for m := range r.evt.MsgSig {
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
				log.Printf("玩家[%d] 進入 桌子[%d] 座位[%d] 狀態[%s] 金幣[%s]", data.User.GetUserID(), data.User.GetTableNo(), data.User.GetSeatNo(), gamelogic.UserStateNameByValue(int(data.User.GetStateValue())), humanize.Commaf(float64(data.User.GetScore())/100))
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GameMessageClassID_GameStartID) {
				log.Printf("=== 遊戲開始 ===")
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GameMessageClassID_GameEndID) {
				log.Printf("=== 遊戲結束 ===")
				r.ready()
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GameMessageClassID_UserStateChangeID) {
				data := new(netproto.UserStateChange)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("玩家[%d] 狀態[%s] 金幣[%s]", data.User.GetUserID(), gamelogic.UserStateNameByValue(int(data.User.GetStateValue())), humanize.Commaf(float64(data.User.GetScore())/100))
				if data.User.GetStateValue() == 0 { // 排隊
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_GameRoom)
					m.SClassID = int32(netproto.GameRoomClassID_UserQueueID)
					if err := r.client.SendMessage(m); err != nil {
						log.Fatal(err)
					}
				}
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_GameStartCountDownID) {
				log.Printf("=== 遊戲準備開始(倒數) ===")
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_GameStartCountDownCancelID) {
				log.Printf("=== 遊戲取消準備 ===")
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_PlayerActionCountDownID) {
				data := new(netproto.INRM_Player)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("玩家[%d] 動作倒數 %d 秒", data.ID, data.ActionCD/1000)
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_WaitFinalRefreshCountDownID) {
				data := new(netproto.INRM_TimeCD)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("最後整理牌堆機會 倒數 %d 秒", data.CD/1000)
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_GameSettingID) {
				data := new(netproto.INRM_GameSetting)
				proto.Unmarshal(m.MsgData.([]byte), data)
				r.Setting = data
				log.Printf("=== 遊戲設定 ===")
				log.Printf("開始倒數時間: %v", time.Duration(data.StartWaitTime)*time.Millisecond)
				log.Printf("動作時間: %v", time.Duration(data.ActionWaitTime)*time.Millisecond)
				log.Printf("每點使用的金幣: %s", humanize.Commaf(float64(data.ScorePerPoint)/100))
				log.Printf("每個玩家手牌數: %d", data.HandCardLength)
				log.Printf("遊戲最大點數: %d", data.MaxPoint)
				log.Printf("棄牌點數: %d", data.GiveupPoint)
				log.Printf("最多Miss次數: %d", data.MissLimit)
				log.Printf("=================")
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_GameStateID) {
				data := new(netproto.INRM_GameState)
				proto.Unmarshal(m.MsgData.([]byte), data)
				r.GameState = data
				switch data.State {
				case netproto.INRM_GState_INRM_GState_WaitPlayer:
					log.Printf("=== 遊戲等待玩家中 ===")
				case netproto.INRM_GState_INRM_GState_CountDown:
					start := time.Unix(data.StartAt, 0)
					log.Printf("遊戲將於 %s 開始，倒數 %d 秒", start.Format(time.RFC3339), int(math.Round(float64(data.StartCD)/1000)))
				case netproto.INRM_GState_INRM_GState_Gaming:
					log.Printf("=== 遊戲進行中 ===")
					log.Printf("鬼牌點數[%d] 牌局[%02d] 玩家[%d]", poker.Rank(data.RandJoker), data.Round, data.PlayerID)
					log.Printf("蓋牌區[%d] 開牌區[%d][%v]", data.CloseDeckCount, data.OpenDeckCount, poker.ReversePokerCard(data.FirstOpenDeck))
					var strs []string
					for _, p := range data.Players {
						strs = append(strs, fmt.Sprintf("玩家[%d][%s][%d]", p.ID, pStateName(p.State), p.Miss))
					}
					log.Printf("%s", strings.Join(strs, " ; "))
					log.Printf("=================")
					if data.PlayerID == r.uid {
						r.popCard()
					}
				}
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_GameResultID) {
				data := new(netproto.INRM_GameResult)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("=== 遊戲結算 ===")
				printPlayerAll(data.Winner, "")
				for _, loss := range data.Losser {
					printPlayerAll(loss, "")
				}
				log.Printf("=================")
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_PlayerCardID) {
				data := new(netproto.INRM_PlayerCards)
				proto.Unmarshal(m.MsgData.([]byte), data)
				r.deck = poker.NewEmptyDeck()
				for _, c := range data.Cards {
					r.deck.AppendCard(poker.ReversePokerCard(c))
				}
				log.Printf("獲得手牌 %v", r.deck)
				r.refreshHands()
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_DeckResetID) {
				log.Printf("蓋牌區卡牌已用完，重新洗牌")
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID) {
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_PopOpenDeckHappenedID) {
				data := new(netproto.INRM_Player)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("玩家[%d] 抽取一張開牌區", data.ID)
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_PopCloseDeckHappenedID) {
				data := new(netproto.INRM_Player)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("玩家[%d] 抽取一張蓋牌區", data.ID)
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_GiveupHappenedID) {
				data := new(netproto.INRM_Player)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("玩家[%d] 棄牌", data.ID)
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_ThrowCardHappenedID) {
				data := new(netproto.INRM_ThrowCard)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("玩家[%d] 丟出 [%v]", data.PlayerID, poker.ReversePokerCard(data.Throw))
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_DeclareHappenedID) {
				data := new(netproto.INRM_Declare)
				proto.Unmarshal(m.MsgData.([]byte), data)
				var hands []*poker.Deck
				for _, d := range data.Hands.Decks {
					hands = append(hands, poker.ReversePokerDeck(d))
				}
				log.Printf("玩家[%d] 宣告勝利 丟[%v] 手牌 %v", data.PlayerID, poker.ReversePokerCard(data.Throw), hands)
				r.refreshHands()
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_MissHappenedID) {
				data := new(netproto.INRM_Player)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("玩家[%d] 動作超時", data.ID)
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_BluffHappenedID) {
				data := new(netproto.INRM_Player)
				proto.Unmarshal(m.MsgData.([]byte), data)
				log.Printf("玩家[%d] 詐胡", data.ID)
				continue
			}

			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandRetID) {
				data := new(netproto.INRM_ActionRet)
				proto.Unmarshal(m.MsgData.([]byte), data)
				if data.Success {
					log.Printf("你更新手牌成功")
				} else {
					log.Printf("你更新手牌失敗, %s", data.Message)
				}
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_PopOpenDeckRetID) {
				data := new(netproto.INRM_PopCardRet)
				proto.Unmarshal(m.MsgData.([]byte), data)
				if data.Success {
					log.Printf("你抽取一張開牌區 %v", poker.ReversePokerCard(data.Card))
					r.deck.AppendCard(poker.ReversePokerCard(data.Card))
					r.throwOrVictory()
				} else {
					log.Printf("你抽取一張開牌區失敗, %s", data.Message)
				}
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_PopCloseDeckRetID) {
				data := new(netproto.INRM_PopCardRet)
				proto.Unmarshal(m.MsgData.([]byte), data)
				if data.Success {
					log.Printf("你抽取一張蓋牌區 %v", poker.ReversePokerCard(data.Card))
					r.deck.AppendCard(poker.ReversePokerCard(data.Card))
					r.throwOrVictory()
				} else {
					log.Printf("你抽取一張蓋牌區失敗, %s", data.Message)
				}
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_GiveupRetID) {
				data := new(netproto.INRM_ActionRet)
				proto.Unmarshal(m.MsgData.([]byte), data)
				if data.Success {
					log.Printf("你棄牌成功")
				} else {
					log.Printf("你棄牌失敗, %s", data.Message)
				}
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_ThrowCardRetID) {
				data := new(netproto.INRM_ActionRet)
				proto.Unmarshal(m.MsgData.([]byte), data)
				if data.Success {
					log.Printf("你丟出卡片成功")
				} else {
					log.Printf("你丟出卡片失敗, %s", data.Message)
				}
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_DeclareRetID) {
				data := new(netproto.INRM_ActionRet)
				proto.Unmarshal(m.MsgData.([]byte), data)
				if data.Success {
					log.Printf("你宣告勝利成功")
				} else {
					log.Printf("你宣告勝利失敗, %s", data.Message)
				}
				continue
			}
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_MissRetID) {
				data := new(netproto.INRM_ActionRet)
				proto.Unmarshal(m.MsgData.([]byte), data)
				if data.Success {
					log.Printf("你動作超時成功")
				} else {
					log.Printf("你動作超時失敗, %s", data.Message)
				}
				continue
			}
		}
	}()

	{ // 登入
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_LoginRoomID)
		data := new(netproto.LoginGameRoomInfo)
		data.UserID = proto.Int32(r.uid)
		data.Cer = proto.String(r.cer)
		data.HDCode = proto.String("")
		data.HDType = proto.Int32(0)
		m.MsgData, _ = proto.Marshal(data)
		if err := r.client.SendMessage(m); err != nil {
			log.Fatal(err)
		}
	}

	<-logined

	select {}
}
