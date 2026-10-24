package robot

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/game/inrm"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
)

// Logic 機器人邏輯
type Logic struct {
	*robotlogic.LogicBase
	refresher inrm.Refresher
}

// NewLogic 新建機器人邏輯
func NewLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *Logic {
	rblgc := new(Logic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rblgc.refresher = inrm.DefaultRefresher()
	return rblgc
}

// OnGameMessage 遊戲訊息
func (logic *Logic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {

	gs := logic.GameLogic.(*inrm.GameServer)

	if gs.Game == nil {
		return
	}

	// 遊戲中 換你
	if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_GameStateID) {
		data := m.MsgData.(*netproto.INRM_GameState)
		mlog.Debug("機器人[%d] 收到訊息 %v", robot.UserInfo.UserID, data)
		if data.PlayerID != robot.UserInfo.UserID {
			return
		}
		p, err := gs.Game.Player(robot.UserInfo.UserID)
		if err != nil {
			mlog.Error("機器人[%d] 取得玩家身份失敗, %v", robot.UserInfo.UserID, err)
			return
		}
		logic.popCard(robot, p)
		return
	}

	// 抽完卡
	if (m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_PopOpenDeckRetID)) ||
		(m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_PopCloseDeckRetID)) {
		data := m.MsgData.(*netproto.INRM_PopCardRet)
		if !data.Success {
			mlog.Warn("抽牌失敗, %v", data.Message)
			return
		}
		p, err := gs.Game.Player(robot.UserInfo.UserID)
		if err != nil {
			mlog.Error("機器人[%d] 取得玩家身份失敗, %v", robot.UserInfo.UserID, err)
			return
		}
		logic.throwOrVictory(robot, p)
		return
	}

	// 有人贏了
	if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.INRM_GameMessageClassID_INRM_DeclareHappenedID) {
		p, err := gs.Game.Player(robot.UserInfo.UserID)
		if err != nil {
			mlog.Error("機器人[%d] 取得玩家身份失敗, %v", robot.UserInfo.UserID, err)
			return
		}
		logic.refreshHands(robot, p)
		return
	}

}

func (logic *Logic) refreshHands(robot *robotcomm.Robot, p *inrm.Player) {

	gs := logic.GameLogic.(*inrm.GameServer)
	if gs.Game == nil {
		return
	}

	_, _, hands := logic.refresher.Refresh(gs.Game.HandCardLength, gs.Game.RandJoker, p.Deck)

	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
	data := new(netproto.Poker_Hand)
	for _, d := range hands {
		data.Decks = append(data.Decks, poker.TransePokerDeck(d))
	}

	m.MsgData, _ = proto.Marshal(data)

	logic.SendGameMessage(robot.UserInfo, m)
}

func (logic *Logic) popCard(robot *robotcomm.Robot, p *inrm.Player) {

	gs := logic.GameLogic.(*inrm.GameServer)
	if gs.Game == nil {
		return
	}

	// 檢查開牌區的牌加入可不可以胡
	d := poker.NewEmptyDeck(p.Deck.Cards...)
	d.AppendCard(gs.Game.OpenDeck.Cards[0])

	// 可以胡就抽開牌區，不能就抽蓋牌區
	if ok, _, _ := logic.refresher.Refresh(gs.Game.HandCardLength, gs.Game.RandJoker, d); ok {
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_Game)
		m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopOpenDeckID)
		logic.SendGameMessage(robot.UserInfo, m)
		return
	}

	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopCloseDeckID)

	logic.SendGameMessage(robot.UserInfo, m)
}

func (logic *Logic) throwOrVictory(robot *robotcomm.Robot, p *inrm.Player) {

	gs := logic.GameLogic.(*inrm.GameServer)
	if gs.Game == nil {
		return
	}

	ok, throw, hands := logic.refresher.Refresh(gs.Game.HandCardLength, gs.Game.RandJoker, p.Deck)

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

		logic.SendGameMessage(robot.UserInfo, m)
		return
	}

	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_ThrowCardID)
	data := new(netproto.INRM_ThrowCard)
	data.Throw = poker.TransePokerCard(throw)
	m.MsgData, _ = proto.Marshal(data)

	logic.SendGameMessage(robot.UserInfo, m)

	// logic.refreshHands(robot, p)
}
