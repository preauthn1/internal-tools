//==================================================================
//文件名： ermjlogic.go
//说明：   麻将战斗逻辑
//修改记录：
//==================================================================
package ermj

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mahjong"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/util"

	"example.invalid/golang/protobuf/proto"
)

//胡牌
type STHuCardInfo struct {
	card       *mahjong.CardGroup //胡牌
	hucardType int32              //胡牌类型
	cardType   int32              //牌型
	rate       int64              //倍率
	score      int64              //得分
	huSeatNo   int32              //胡牌玩家
	fromSeatNo int32              //支付玩家
}

//玩家操作事件，用于响应不同优先级事件
type STUserOpEvent struct {
	opSeatNo  int32   //op玩家座位号
	srcSeatNo int32   //来源
	eventId   int32   //事件id。1pass,2peng,3杠，4胡
	card      []int32 //数据
}

// 检测杠胡
type STCheckGanghu struct {
	bGang      bool  //标志
	seatNoGang int32 //座位号
	seatNosrc  int32 //来源座位号
	times      int32 //出牌次数
}

type STRate struct {
	PH      int32 //屁胡
	DDZ     int32 //大对子
	JGD     int32 //金勾吊
	QDZ     int32 //七对子
	LQD     int32 //龙七对
	QYS     int32 //清一色
	DYJ     int32 //带幺九
	TH      int32 //天胡
	DH      int32 //地胡
	GANG    int32 //杠 跟
	ZIMO    int32 //自摸
	HD      int32 //海底捞月
	GSH     int32 //杠上花
	GSP     int32 //杠上炮
	RATEMAX int32 //番数上限
}

func NewSTRate() *STRate {
	p := new(STRate)
	p = &STRate{0, 1, 2, 2, 3, 2, 2, 6, 6, 1, 1, 2, 3, 2, 32}
	return p
}

type STOutTime struct {
	Choose    int32 //选择倍数 定缺超时时间
	Action    int32 //出牌 吃碰杠出牌时间
	Deal      int32 //发牌动画时间
	BeginGame int32 //游戏开始动画
	Dice      int32 //丢色子动画超时时间
	GameEnd   int32 //游戏结束等待时间
	Huazhu    int32 //查花猪
	Dajiao    int32 //查大叫
	FightEnd  int32 //战斗结束
	Light     int32 //亮牌
}
type STQiangGang struct {
	Card   int32 //操作牌
	SeatNo int32 //座位号
	Score  int64 //得分
}

//麻将
type MahjongSC struct {
	*gamelogic.FixedGameLogicBase
	timerTimeOutAction     *event.TimerFunc
	timerTimeOutEvent      *event.TimerFunc
	timerTimeTuoguanAction *event.TimerFunc
	configs                *config.Group
	randGroup              mahjong.Mahjong  //公共牌池
	gameUserData           []*STUserData    //用户数据
	huCardInfo             []*STHuCardInfo  //胡牌信息
	curOpEvent             []*STUserOpEvent //玩家操作事件
	curCanOpUser           map[int32]bool   //当前可操作玩家
	dice                   []int32          //骰子
	gameIndex              int32            //牌局计数
	animationName          string           //当前动画名字
	baseScore              int64            //底分
	readywaittime          int32            //准备时间
	gameUserCount          int32            //用户人数
	surplusCardNum         int32            //剩余牌张数
	buttonSeatNo           int32            //庄家座位号
	nextButtonSeatNo       int32            //下一局庄家座位号
	nRandomType            int32            //发牌类型,配置发牌
	lastOutCard            int32            //上一个玩家出牌
	lastOpType             int32            //上一个玩家操作类型
	lastOpSeatNo           int32            //上一个玩家座位号
	currOpSeatNo           int32            //当前操作玩家
	gameState              int32            //游戏状态
	endCardSize            int32            //游戏结束时剩余牌张数
	curEventCount          int32            //当前事件所需操作人数
	huNum                  int32            //胡牌玩家数
	nextSeatNoTmp          int32            //临时下一个玩家操作位置
	outcardTimes           int32            //出牌次数
	canDihu                bool             //是否可以地胡
	checkGangHu            *STCheckGanghu   //检测杠胡
	rateSettleMent         *STRate          //结算番数
	outTimeConfig          *STOutTime       //超时时间
	isEnd                  bool             //战斗是否结束
	isEndCheck             bool             //是否查底
	qianggang              *STQiangGang     //抢杠牌
	bOpenTimer             bool             //是否打开定时器
	isDYJ                  int32            //是否開啟帶么九
}

//初始化游戏数据
func (m *MahjongSC) initGameData() {
	m.gameIndex++
	m.gameUserData = nil
	m.huCardInfo = make([]*STHuCardInfo, 0)
	m.curOpEvent = make([]*STUserOpEvent, 0)
	m.dice = nil
	m.animationName = ""
	m.gameUserCount = 0
	m.surplusCardNum = 0
	m.lastOutCard = 0
	m.lastOpType = 0
	m.lastOpSeatNo = -1
	m.currOpSeatNo = -1
	m.gameState = 0
	m.curEventCount = 0
	// m.outTimeBegin = 0
	m.curCanOpUser = make(map[int32]bool, 0)
	m.huNum = 0
	m.setGameUserData() //设置玩家数据
	//更新庄家座位号
	m.updateButton()
	m.timerTimeOutAction = nil
	m.timerTimeOutEvent = nil
	m.outcardTimes = 0
	m.canDihu = true
	m.checkGangHu = new(STCheckGanghu)
	m.checkGangHu.bGang = false
	m.isEnd = false
	m.isEndCheck = false
	//	m.isGameOver = false
	m.qianggang = new(STQiangGang)
	m.qianggang.Card = 0
	m.bOpenTimer = true
}

//创建游戏逻辑
func NewMahjongGame(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *MahjongSC {
	m := new(MahjongSC)
	m.FixedGameLogicBase = gamelogic.NewFixedGameLogicBase(t, tfManager, cfgs)

	m.getGameConfig(cfgs)
	m.gameIndex = 0
	m.nextButtonSeatNo = 0
	m.RegisterAfterGameStartFunc(func() {
		m.initGameData()
		m.begin()

	})
	m.RegisterBeforeGameStartFunc(func() {
		m.makeCardGroup()
	})
	return m
}

//游戏消息
func (m *MahjongSC) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	mlog.Debug("GameMessage, receive message，SeatNo：%d", user.SeatNo)
	//动画
	m.CheckUserAnimationCommited(user, gmsg, m.animationName, nil)
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.GameMessageClassID_UserActionCommitID) {
		mlog.Debug("GameMessage,animation commite SeatNo=%d,Name=%s", user.SeatNo, m.animationName)
		return
	}
	//操作合法性校验
	if !m.oprateCheckLegal(user, gmsg) {
		m.rspError(user, ERROR_STATE, "game end")
		return
	}
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.SCMJ_GameMessageClassID_SCNJActionId) {
		m.actionInfo(user, gmsg.MsgData.([]byte))
	} else if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.SCMJ_GameMessageClassID_C2S_SCMJRemindId) {
		m.remindOp(user, gmsg.MsgData.([]byte))
	} else if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.SCMJ_GameMessageClassID_C2S_TuoGuanId) {
		m.tuoGuan(user, gmsg.MsgData.([]byte))
	}
}

// 构造牌型
func (m *MahjongSC) makeCardGroup() {
	if 1 == m.nRandomType {
		m.randGroup = mahjong.NewFileCards(m.Configs) //初始化一副牌,配置文件
	} else {
		m.addCardGroup()
	}
	m.randGroup.Shuffle()
}
func (m *MahjongSC) addCardGroup() {
	m.randGroup = mahjong.NewRandomEmpty()
	cfg := m.configs.GetConfig("gameserver")
	// 筒
	if cfg.GetConfigInt32("gameprivate", "havetong") == 1 {
		cg := mahjong.NewTong(1)
		m.randGroup.AddCardGroup(cg)
	}
	// 条
	if cfg.GetConfigInt32("gameprivate", "havetiao") == 1 {
		cg := mahjong.NewTiao(1)
		m.randGroup.AddCardGroup(cg)
	}
	// 万
	if cfg.GetConfigInt32("gameprivate", "havewan") == 1 {
		cg := mahjong.NewWan(1)
		m.randGroup.AddCardGroup(cg)
	}
	// 风
	if cfg.GetConfigInt32("gameprivate", "havefeng") == 1 {
		cg := mahjong.NewFeng(1)
		m.randGroup.AddCardGroup(cg)
	}
	// 箭
	if cfg.GetConfigInt32("gameprivate", "havejian") == 1 {
		cg := mahjong.NewJian(1)
		m.randGroup.AddCardGroup(cg)
	}
	// 花
	if cfg.GetConfigInt32("gameprivate", "havehua") == 1 {
		cg := mahjong.NewHua(1)
		m.randGroup.AddCardGroup(cg)
	}
	// 百搭
	if cfg.GetConfigInt32("gameprivate", "havebaida") == 1 {
		cg := mahjong.NewBaida(1)
		m.randGroup.AddCardGroup(cg)
	}
	if m.randGroup.Len() <= 0 {
		m.randGroup = mahjong.NewRandomSimple(1)
	}
}

// 托管
func (m *MahjongSC) tuoGuan(user *gamelogic.GameUser, data []byte) {

	ust := &netproto.C2S_SCMJ_TuoGuan{}
	err := proto.Unmarshal(data, ust)
	//错误校验
	if err != nil {
		mlog.Error("桌子号[%d],解析用户%v用户错误%v", m.TableNo(), user, err)
		return
	}
	ud := m.getUserInfoBySeatNo(user.SeatNo)
	ud.isTuoGuan = *ust.State
	if ud.isTuoGuan {
		mlog.Debug("托管,桌子号[%d],用户id[%d],SeatNo[%d]", m.TableNo(), ud.User.UserID, ud.seatNo)
		m.GameLog.AddGameStep("托管,SeatNo[%d]", ud.seatNo)
	} else {
		mlog.Debug("取消托管,桌子号[%d],用户id[%d],SeatNo[%d]", m.TableNo(), ud.User.UserID, ud.seatNo)
		m.GameLog.AddGameStep("取消托管,SeatNo[%d]", ud.seatNo)
		// 取消托管不再自动出牌
		if m.currOpSeatNo == user.SeatNo && ud.moPai > 0 {
			if m.timerTimeTuoguanAction != nil {
				m.timerTimeTuoguanAction.Stop()
				m.timerTimeTuoguanAction = nil
			}
		}
		if ud.timerTimeTuoguanEvent != nil {
			ud.timerTimeTuoguanEvent.Stop()
			ud.timerTimeTuoguanEvent = nil
		}
		if ud.seatNo == m.currOpSeatNo {
			m.remindInfo(ud.seatNo)
		}
	}

	m.BroadcastTuoGuan(user.SeatNo, ud.isTuoGuan)
	if ud.isTuoGuan {
		if m.gameState == PROGRESS_FIGHT && len(ud.canOpEventInfo.OpInfo) > 0 {
			m.tuoguanOpEvent(ud)
		} else if (m.gameState == PROGRESS_FIGHT && m.currOpSeatNo == user.SeatNo && ud.moPai > 0) || (m.gameState == PROGRESS_FIGHT && m.currOpSeatNo == user.SeatNo && ud.handCards.Len()%3 != 1) {
			OpValue := make([]int32, 0)
			if m.Configs.GetConfig("gameserver").GetConfigInt("gameprivate", "tuoguantype") == 1 {
				bestcd := ud.calBestOutCard()
				OpValue = append(OpValue, bestcd)
			} else {
				OpValue = append(OpValue, ud.moPai)
			}
			m.OutCard(user, OpValue)
			m.mahjongWaitOpt()
		} else if m.gameState == PROGRESS_CHOOSE && ud.color <= 0 {
			dingqueColor := ud.autoDingque()
			OpValue := make([]int32, 0)
			OpValue = append(OpValue, dingqueColor)
			m.ChooseColor(user, OpValue)
		}
	}
}

// 提示
func (m *MahjongSC) remindOp(user *gamelogic.GameUser, data []byte) {
}

// 提示
func (m *MahjongSC) remindInfo(seatNo int32) {
	mlog.Debug("提示计算,桌子号[%d],SeatNo[%d]", m.TableNo(), seatNo)
	ud := m.getUserInfoBySeatNo(seatNo)
	if ud == nil {
		return
	}
	if ud.isTuoGuan {
		return
	}
	if ud.User.IsRobot() {
		return
	}
	cg := new(mahjong.CardGroup)
	if ud.moPai/100 == ud.color {
		cg.AddCard(ud.moPai)
	} else {
		for _, cd := range ud.handCards.Cards {
			if cd/100 == ud.color {
				cg = new(mahjong.CardGroup)
				cg.AddCard(cd)
				break
			} else {
				cg.AddCard(cd)
			}
		}
		cg.AddCard(ud.moPai)
	}

	lastcd := int32(0)
	ud.remindData = make(map[int32][]*STRemindInfo, 0)
	sort.Sort(cg)
	for _, cd := range cg.Cards {
		if lastcd == cd {
			continue
		}

		lastcd = cd
		info := m.calTingInfo(seatNo, lastcd)
		if len(info) > 0 {
			ud.remindData[lastcd] = info
		}

	}
	// if len(ud.remindData) <= 0 {
	// 	// mlog.Debug("桌子号[%d],不能听牌不需要提示", m.TableNo())
	// 	return
	// }
	m.BroadcastRemind(seatNo, ud.remindData)
}

//获取游戏配置
func (m *MahjongSC) getGameConfig(cfgs *config.Group) {
	m.rateSettleMent = NewSTRate()
	m.outTimeConfig = new(STOutTime)
	//mlog.Debug("获取游戏配置")
	m.configs = cfgs
	cfg := cfgs.GetConfig("gameserver")

	m.baseScore = int64(cfg.GetConfigInt("gameserver", "basescore"))
	m.readywaittime = int32(cfg.GetConfigInt("gameserver", "readywaittime"))

	m.nRandomType = int32(cfg.GetConfigInt("gameprivate", "cardtype"))

	m.outTimeConfig.GameEnd = cfg.GetConfigInt32("gameprivate", "outtimegameend")
	m.outTimeConfig.Action = cfg.GetConfigInt32("gameprivate", "outtimeaction")
	m.outTimeConfig.Choose = cfg.GetConfigInt32("gameprivate", "outtimechoose")
	m.outTimeConfig.BeginGame = cfg.GetConfigInt32("gameprivate", "outtimeanibegin")
	m.outTimeConfig.Deal = cfg.GetConfigInt32("gameprivate", "outtimeanideal")
	m.outTimeConfig.Dice = cfg.GetConfigInt32("gameprivate", "outtimeanidice")
	m.outTimeConfig.Huazhu = cfg.GetConfigInt32("gameprivate", "outtimeanihuazhu")
	m.outTimeConfig.Dajiao = cfg.GetConfigInt32("gameprivate", "outtimeanidajiao")
	m.outTimeConfig.FightEnd = cfg.GetConfigInt32("gameprivate", "outtimeanifightend")
	m.outTimeConfig.Light = cfg.GetConfigInt32("gameprivate", "outtimeanilight")

	m.endCardSize = cfg.GetConfigInt32("gameprivate", "endsizecard")

	//番数
	m.rateSettleMent.PH = cfg.GetConfigInt32("gameprivate", "ratepihu")
	m.rateSettleMent.DDZ = cfg.GetConfigInt32("gameprivate", "ratedaduizi")
	m.rateSettleMent.JGD = cfg.GetConfigInt32("gameprivate", "ratejgd")
	m.rateSettleMent.QDZ = cfg.GetConfigInt32("gameprivate", "rateqiduiu")
	m.rateSettleMent.LQD = cfg.GetConfigInt32("gameprivate", "ratelqd")
	m.rateSettleMent.QYS = cfg.GetConfigInt32("gameprivate", "rateqys")
	m.rateSettleMent.DYJ = cfg.GetConfigInt32("gameprivate", "ratedyj")
	m.rateSettleMent.TH = cfg.GetConfigInt32("gameprivate", "ratetianhu")
	m.rateSettleMent.DH = cfg.GetConfigInt32("gameprivate", "ratedihu")
	m.rateSettleMent.GANG = cfg.GetConfigInt32("gameprivate", "rategang")
	m.rateSettleMent.ZIMO = cfg.GetConfigInt32("gameprivate", "ratezimo")
	m.rateSettleMent.HD = cfg.GetConfigInt32("gameprivate", "ratehaidi")
	m.rateSettleMent.GSH = cfg.GetConfigInt32("gameprivate", "rategangsh")
	m.rateSettleMent.GSP = cfg.GetConfigInt32("gameprivate", "rategangsp")
	m.rateSettleMent.RATEMAX = cfg.GetConfigInt32("gameprivate", "ratemax")
}

// 更新庄家座位号
func (m *MahjongSC) updateButton() {
	if m.nextButtonSeatNo >= 0 {
		m.buttonSeatNo = m.nextButtonSeatNo
	} else {
		m.buttonSeatNo = 0
	}
	m.nextButtonSeatNo = -1
	mlog.Debug("桌子号[%d],庄家座位号[%d]", m.TableNo(), m.buttonSeatNo)
	for i := 0; i < 10; i++ {
		ud := m.getUserInfoBySeatNo(m.buttonSeatNo)
		if ud != nil {
			return
		}
	}
	mlog.Error("庄家座位号错误 桌子号[%d]", m.TableNo())
}

//加载用户数据
func (m *MahjongSC) setGameUserData() {

	m.gameUserData = make([]*STUserData, 0)
	m.gameUserCount = 0
	for _, ur := range m.AllTableUser() {
		userdata := NewUserData(ur)
		m.gameUserData = append(m.gameUserData, userdata)
		m.gameUserCount++
	}
	mlog.Debug("设置玩家数据,桌子号[%d],玩家人数：%d人。", m.TableNo(), m.gameUserCount)
}

//开始游戏
func (m *MahjongSC) begin() {
	mlog.Debug("游戏开始,桌子号[%d]", m.TableNo())
	//注册游戏开始动画
	m.animationName = m.RegisterAnimation(ANIMATION_BEGIN, m.outTimeConfig.BeginGame, func() {
		//丢骰子
		m.ExchangeGameState(PROGRESS_DICE)
		m.PlayDice(2)
		//注册色子动画
		m.animationName = m.RegisterAnimation(ANIMATION_DICE, m.outTimeConfig.Dice, func() {
			//摸牌
			m.ExchangeGameState(PROGRESS_DEAL)
			m.Deal()
			//注册发牌动画
			m.animationName = m.RegisterAnimation(ANIMATION_DEAL, m.outTimeConfig.Deal, func() {
				m.TakeCard(m.buttonSeatNo, 1)
				if m.Configs.GetConfig("gameserver").GetConfigInt("gameprivate", "havedingque") == 1 {
					m.ExchangeGameState(PROGRESS_CHOOSE)
				} else if m.Configs.GetConfig("gameserver").GetConfigInt("gameprivate", "havechooserate") == 1 {
					// 需要加个阶段用于选倍数
					m.ExchangeGameState(PROGRESS_CHOOSE)
				} else {
					m.currOpSeatNo = m.buttonSeatNo
					m.nextSeatNoTmp = m.buttonSeatNo
					mlog.Debug("轮到玩家,桌子号[%d],SeatNo[%d]操作", m.TableNo(), m.currOpSeatNo)
					m.ExchangeGameState(PROGRESS_FIGHT)
					//广播玩家可操作信息updataCanOpInfo
					ud := m.getUserInfoBySeatNo(m.buttonSeatNo)
					if ud != nil {
						m.updateCanOpInfo(m.buttonSeatNo, ud.moPai)
					}
				}
				m.remindInfo(m.currOpSeatNo)
				m.mahjongWaitOpt() //等待玩家操作
			})
			mlog.Debug("桌子号[%d],注册动画[%s],Time[%d]", m.TableNo(), ANIMATION_DEAL, m.outTimeConfig.Deal)
		})
		mlog.Debug("桌子号[%d],注册动画[%s],Time[%d]", m.TableNo(), ANIMATION_DICE, m.outTimeConfig.Dice)
	})
	mlog.Debug("桌子号[%d],注册动画[%s],Time[%d]", m.TableNo(), ANIMATION_BEGIN, m.outTimeConfig.BeginGame)

}

//吃牌
func (m *MahjongSC) EatCard(user *gamelogic.GameUser, OpValue []int32) {
	mlog.Debug("吃牌,桌子号[%d],用户id[%d]", m.TableNo(), user.UserID)
}

//明杠牌
func (m *MahjongSC) MingGang(user *gamelogic.GameUser, cd int32) {
	mlog.Debug("玩家操作明杠,桌子号[%d],用户id[%d],SeatNo[%d]", m.TableNo(), user.UserID, user.SeatNo)
	ud := m.getUserInfoBySeatNo(user.SeatNo)
	//mlog.System("MingGang 桌子号[%d],用户id[%d],手牌[%s],摸牌[%d]",  m.TableNo(), user.UserID, ud.handCards.String(), ud.moPai)
	if ud == nil {
		m.rspError(user, 1, "User data is not exist")
		return
	}
	if m.randGroup.Len() <= int(m.endCardSize) {
		m.rspError(user, 1, "card pool error")
		return
	}
	//操作成功，注册记录
	opEvent := new(STUserOpEvent)
	opEvent.opSeatNo = user.SeatNo
	opEvent.srcSeatNo = user.SeatNo
	opEvent.eventId = OP_MINGGANG
	opEvent.card = append(opEvent.card, cd)
	m.curOpEvent = append(m.curOpEvent, opEvent)
	m.curCanOpUser[user.SeatNo] = false
	//记录
	cggp := new(STGangCard)
	cggp.gangCard = cd
	cggp.gangType = OP_MINGGANG
	cggp.seatNo = ud.seatNo
	cggp.valid = false
	ud.gangCards = append(ud.gangCards, cggp)

	mlog.Debug("玩家操作明杠 成功,桌子号[%d],用户id[%d],SeatNo[%d]", m.TableNo(), user.UserID, user.SeatNo)
}

//引杠
func (m *MahjongSC) YinGang(user *gamelogic.GameUser, cd int32) {
	mlog.Debug("玩家操作引杠,桌子号[%d],用户id[%d],SeatNo[%d]", m.TableNo(), user.UserID, user.SeatNo)
	ud := m.getUserInfoBySeatNo(user.SeatNo)
	//mlog.System("YinGang 桌子号[%d],用户id[%d],手牌[%s],摸牌[%d]",  m.TableNo(), user.UserID, ud.handCards.String(), ud.moPai)
	if ud == nil {
		m.rspError(user, 1, "User data is not exist")
		return
	}
	if m.randGroup.Len() <= int(m.endCardSize) {
		m.rspError(user, 1, "card pool error")
		return
	}
	opEvent := new(STUserOpEvent)
	opEvent.opSeatNo = user.SeatNo
	opEvent.srcSeatNo = m.lastOpSeatNo
	opEvent.eventId = OP_YINGANG
	opEvent.card = append(opEvent.card, cd)
	opEvent.card = append(opEvent.card, cd)
	opEvent.card = append(opEvent.card, cd)
	opEvent.card = append(opEvent.card, cd)
	m.curOpEvent = append(m.curOpEvent, opEvent)
	m.curCanOpUser[user.SeatNo] = false

	//记录
	cggp := new(STGangCard)
	cggp.gangCard = cd
	cggp.gangType = OP_YINGANG
	cggp.seatNo = ud.seatNo
	cggp.valid = false
	ud.gangCards = append(ud.gangCards, cggp)
	//需要即时广播吗
	mlog.Debug("玩家操作引杠 成功,桌子号[%d],用户id[%d],SeatNo[%d]cd[%d]", m.TableNo(), user.UserID, user.SeatNo, cd)
}

//暗杠
func (m *MahjongSC) AnGang(user *gamelogic.GameUser, cd int32) {
	mlog.Debug("玩家操作暗杠,桌子号[%d],用户id[%d],SeatNo[%d]", m.TableNo(), user.UserID, user.SeatNo)
	ud := m.getUserInfoBySeatNo(user.SeatNo)
	//mlog.System("AnGang 桌子号[%d],用户id[%d],手牌[%s],摸牌[%d]",  m.TableNo(), user.UserID, ud.handCards.String(), ud.moPai)
	if ud == nil {
		m.rspError(user, 1, "User data is not exist")
		return
	}
	if ud.moPai <= 100 {
		m.rspError(user, 1, "logic error")
		return
	}
	if m.randGroup.Len() <= int(m.endCardSize) {
		m.rspError(user, 1, "card pool error")
		return
	}
	//减少手牌，是否广播手牌减少
	ud.checkErrorHandSize()
	if ud.moPai <= 0 {
		mlog.Error("杠牌错误AnGang,非法摸牌,桌子号[%d],用户id[%d]", m.TableNo(), user.UserID)
	}
	ud.handCards.AddCard(ud.moPai)
	ud.handCards.DeleteCard(cd, 4)
	ud.moPai = 0
	ud.checkErrorHandSize()
	//添加杠牌池
	gc := new(STGangCard)
	gc.gangType = OP_ANGANG
	gc.seatNo = ud.seatNo
	gc.gangCard = cd
	gc.valid = true
	ud.gangCards = append(ud.gangCards, gc)
	m.canDihu = false
	//广播
	cg := make([]int32, 0)
	cg = append(cg, cd)
	cg = append(cg, cd)
	cg = append(cg, cd)
	cg = append(cg, cd)
	m.BoardScenceExchange(ud.seatNo, ud.seatNo, OP_ANGANG, -1, cg)
	m.SettleMentIntimeGang(ud.seatNo, ud.seatNo, OP_ANGANG, cd)
	mlog.Debug("玩家操作暗杠 成功,桌子号[%d],用户id[%d],SeatNo[%d]cd[%d]", m.TableNo(), user.UserID, user.SeatNo, cd)
	m.GameLog.AddGameStep("暗杠,SeatNo[%d],牌[%s]", ud.seatNo, mahjong.CardString(cd))

	m.setGangHuFlag(user.SeatNo, user.SeatNo)
	m.curCanOpUser[user.SeatNo] = false
	//摸一张牌
	m.currOpSeatNo = ud.seatNo
	if m.randGroup.Len() > 0 {
		m.TakeCard(m.currOpSeatNo, 1)
		m.BroadcastWait()
		m.remindInfo(m.currOpSeatNo)
	}
	m.resetTimmerOp()
	m.resetTimmerEvent()
}

//胡牌
func (m *MahjongSC) HuCard(user *gamelogic.GameUser) {
	mlog.Debug("玩家操作胡牌,桌子号[%d],用户id[%d],SeatNo[%d]", m.TableNo(), user.UserID, user.SeatNo)
	ud := m.getUserInfoBySeatNo(user.SeatNo)
	if ud == nil {
		m.rspError(user, 1, "User data is not exist")
		return
	}
	if ud.eventFlag {
		m.rspError(user, 1, "op error")
		return
	}
	// mlog.Debug("HuCard curCanOpUser*********%t*******", m.curCanOpUser[user.SeatNo])
	if !m.curCanOpUser[user.SeatNo] {
		m.rspError(user, 1, "have no this type op")
		return
	}
	// 自摸
	if ud.moPai > 0 && m.currOpSeatNo == ud.seatNo {
		//检测玩家碰牌能否杠
		if !ud.handCards.CanHuBase(ud.moPai) {
			cgcopy := new(mahjong.CardGroup)
			for _, cdcopy := range ud.handCards.Cards {
				cgcopy.AddCard(cdcopy)
			}
			cgcopy.AddCard(ud.moPai)
			if !cgcopy.IsSevenPair() {
				m.rspError(user, 1, "illegaliful Op1")
				return
			}
		}
		if m.checkGangHu.bGang {
			m.EventHu(ud.seatNo, m.checkGangHu.seatNosrc, ud.moPai, OP_HU_GANG)
			mlog.Debug("胡牌,杠上花,桌子号[%d],用户id[%d],SeatNo[%d],SrcSeatNo[%d]", m.TableNo(), user.UserID, ud.seatNo, m.checkGangHu.seatNosrc)
		} else {
			m.EventHu(ud.seatNo, ud.seatNo, ud.moPai, OP_HU_ZIMO)
			mlog.Debug("胡牌,自摸,桌子号[%d],用户id[%d],SeatNo[%d],SrcSeatNo[%d]", m.TableNo(), user.UserID, ud.seatNo, ud.seatNo)
		}
		ud.moPai = 0
	} else {
		//检测玩家碰牌能否杠
		if m.qianggang.Card > 0 {
			// 检测是否能胡
			if !ud.handCards.CanHuBase(m.qianggang.Card) {
				cgcopy := new(mahjong.CardGroup)
				for _, cdcopy := range ud.handCards.Cards {
					cgcopy.AddCard(cdcopy)
				}
				cgcopy.AddCard(m.qianggang.Card)
				if !cgcopy.IsSevenPair() {
					m.rspError(user, 1, "illegaliful Op3")
					return
				}
			}
			// 抢杠胡
			m.EventHu(ud.seatNo, m.qianggang.SeatNo, m.qianggang.Card, OP_HU_QIANGGANG)
			mlog.Debug("胡牌,抢杠胡,桌子号[%d],用户id[%d],SeatNo[%d],SrcSeatNo[%d]", m.TableNo(), user.UserID, ud.seatNo, m.qianggang.SeatNo)
		} else {
			if !ud.handCards.CanHuBase(m.lastOutCard) {
				cgcopy := new(mahjong.CardGroup)
				for _, cdcopy := range ud.handCards.Cards {
					cgcopy.AddCard(cdcopy)
				}
				mlog.Debug("lastOutCard[%d]", m.lastOutCard)
				cgcopy.AddCard(m.lastOutCard)
				if !cgcopy.IsSevenPair() {
					m.rspError(user, 1, "illegaliful Op2")
					return
				}
			}
			if m.checkGangHu.bGang {
				m.EventHu(ud.seatNo, m.lastOpSeatNo, m.lastOutCard, OP_HU_GANGPAO)
				mlog.Debug("胡牌,杠上炮,桌子号[%d],用户id[%d],SeatNo[%d],SrcSeatNo[%d]", m.TableNo(), user.UserID, ud.seatNo, m.lastOpSeatNo)
			} else {
				m.EventHu(ud.seatNo, m.lastOpSeatNo, m.lastOutCard, OP_HU_DIANPAO)
				mlog.Debug("胡牌,点炮胡,桌子号[%d],用户id[%d],SeatNo[%d],SrcSeatNo[%d]", m.TableNo(), user.UserID, ud.seatNo, m.lastOpSeatNo)
			}
		}

	}
	if m.nextButtonSeatNo < 0 {
		m.nextButtonSeatNo = ud.seatNo
		mlog.Debug("第一个胡牌玩家 桌子号[%d],用户id[%d],SeatNo[%d]", m.TableNo(), user.UserID, ud.seatNo)
	}
	ud.isHu = true
	ud.setEnevtFlag(true)
	m.qianggang.Card = 0
}

// 胡牌
func (m *MahjongSC) EventHu(seatNo, srcSeatNo, cd, huType int32) {
	opEvent := new(STUserOpEvent)
	opEvent.opSeatNo = seatNo
	opEvent.srcSeatNo = srcSeatNo
	opEvent.eventId = huType
	opEvent.card = append(opEvent.card, cd)
	m.curOpEvent = append(m.curOpEvent, opEvent)
	m.curCanOpUser[seatNo] = false
	//广播
	mlog.Debug("玩家操作胡牌 成功,桌子号[%d],SeatNo[%d]", m.TableNo(), seatNo)
}

// 重置抢杠标志位
func (m *MahjongSC) resetGangHuFlag() {
	if m.checkGangHu.times != m.outcardTimes {
		m.checkGangHu.bGang = false
	}
	mlog.Debug("重置杠标志位,桌子号[%d],结果[%t]", m.TableNo(), m.checkGangHu.bGang)
}

// 设置抢杠标志位
func (m *MahjongSC) setGangHuFlag(des, src int32) {
	m.checkGangHu.bGang = true
	m.checkGangHu.seatNoGang = des
	m.checkGangHu.seatNosrc = src
	m.checkGangHu.times = m.outcardTimes
	mlog.Debug("设置杠标志位,桌子号[%d],结果[%t]", m.TableNo(), m.checkGangHu.bGang)
}

//出牌
func (m *MahjongSC) OutCard(user *gamelogic.GameUser, OpValue []int32) {
	mlog.Debug("玩家操作出牌,桌子号[%d],用户id[%d],SeatNo[%d],Value%v", m.TableNo(), user.UserID, user.SeatNo, OpValue)
	ud := m.getUserInfoBySeatNo(user.SeatNo)
	if ud == nil {
		m.rspError(user, 1, "User data is not exist")
		return
	}
	if m.isEnd {
		m.rspError(user, 1, "fight is over")
		return
	}
	if len(OpValue) < 1 {
		m.rspError(user, 1, "Op data is empty")
		return
	}
	if m.currOpSeatNo != user.SeatNo {
		m.rspError(user, 1, "It's not your turn")
		mlog.Error("桌子号[%d],用户id[%d],当前操作玩家[%d]", m.TableNo(), user.UserID, m.currOpSeatNo)
		return
	}
	cardsize := ud.handCards.Len()
	if ud.moPai > 0 {
		cardsize += 1
	}
	if cardsize%3 != 2 {
		m.rspError(user, 1, "repeat op outcard")
		return
	}
	cd := OpValue[0]
	//检测出的牌是否为自己的牌
	if !ud.handCards.HaveCardOf(cd) && cd != ud.moPai {
		m.rspError(user, 1, "Op data error")
		mlog.Error("桌子号[%d],用户id[%d],手牌[%s][%s],出牌[%s]", m.TableNo(), user.UserID, ud.handCards.String(), mahjong.CardString(ud.moPai), mahjong.CardString(cd))
		return
	}
	// 抢杠标志位
	m.resetGangHuFlag()

	m.outcardTimes += 1
	m.resetTimmerOp()
	//添加牌河

	ud.handCards.AddCard(ud.moPai)
	ud.handCards.DeleteCard(cd, 1)
	ud.checkErrorHandSize()
	sort.Sort(ud.handCards)
	ud.moPai = 0
	ud.riverCards.Cards = append(ud.riverCards.Cards, cd)
	m.lastOutCard = cd
	m.lastOpSeatNo = ud.seatNo
	m.lastOpType = OP_OUTCARD
	m.BoardScenceExchange(ud.seatNo, -1, OP_OUTCARD, -1, OpValue)
	m.updateCanOpInfos(ud.seatNo, cd)
	mlog.Debug("出牌 成功,桌子号[%d],用户id[%d],SeatNo[%d],Card[%d]", m.TableNo(), user.UserID, user.SeatNo, cd)
	m.GameLog.AddGameStep("出牌,SeatNo[%d],牌[%s]", ud.seatNo, mahjong.CardString(cd))
	m.updateCurOpMessage(ud.seatNo)
}

//碰牌，//考虑操作的优先级，暂时不广播
func (m *MahjongSC) PengCard(user *gamelogic.GameUser, OpValue []int32) {
	mlog.Debug("碰牌,桌子号[%d],用户id[%d],SeatNo[%d]", m.TableNo(), user.UserID, user.SeatNo)
	ud := m.getUserInfoBySeatNo(user.SeatNo)
	if ud == nil {
		m.rspError(user, 1, "User data is not exist")
		return
	}
	if !ud.handCards.HavePairOf(m.lastOutCard) {
		m.rspError(user, 1, "illegaliful Op")
		return
	}
	if ud.eventFlag {
		m.rspError(user, 1, "op error")
		return
	}
	if !m.curCanOpUser[user.SeatNo] {
		m.rspError(user, 1, "have no this type op")
		return
	}
	//添加碰牌池
	pg := new(STPengCard)
	pg.pengCard = m.lastOutCard
	pg.seatNo = m.lastOpSeatNo
	pg.valid = false
	ud.pengCards = append(ud.pengCards, pg)

	opEvent := new(STUserOpEvent)
	opEvent.opSeatNo = user.SeatNo
	opEvent.srcSeatNo = m.lastOpSeatNo
	opEvent.eventId = OP_PENG
	opEvent.card = append(opEvent.card, m.lastOutCard)
	opEvent.card = append(opEvent.card, m.lastOutCard)
	opEvent.card = append(opEvent.card, m.lastOutCard)
	m.curOpEvent = append(m.curOpEvent, opEvent)
	m.curCanOpUser[user.SeatNo] = false
	ud.setEnevtFlag(true)
	//需要即时广播吗
	mlog.Debug("碰牌 成功,桌子号[%d],用户id[%d],SeatNo[%d],cd[%d]", m.TableNo(), user.UserID, user.SeatNo, m.lastOutCard)
}

//过牌
func (m *MahjongSC) PassHand(user *gamelogic.GameUser) {
	mlog.Debug("玩家操作过,桌子号[%d],用户id[%d],SeatNo[%d]", m.TableNo(), user.UserID, user.SeatNo)
	ud := m.getUserInfoBySeatNo(user.SeatNo)
	if ud == nil {
		m.rspError(user, 1, "User data is not exist")
		return
	}
	if ud.eventFlag {
		m.rspError(user, 1, "op error")
		return
	}
	if !m.curCanOpUser[user.SeatNo] {
		m.rspError(user, 1, "have no this type op")
		return
	}
	m.qianggang.Card = 0
	m.resetTimmerEvent()
	ud.resetCanOpInfo()
	m.curCanOpUser[ud.seatNo] = false
	mlog.Debug("桌子号[%d],用户id[%d],玩家SeatNo[%d]操作过成功,curSeatNo[%d],nextSeatNo[%d]", m.TableNo(), user.UserID, user.SeatNo, m.currOpSeatNo, m.nextSeatNoTmp)
	m.GameLog.AddGameStep("过牌,SeatNo[%d]", ud.seatNo)
	ud.setEnevtFlag(true)
	if m.haveEvent() {
		// mlog.Debug("事件不为空，不用摸牌，等待其他玩家操作.")
		return
	}
	if len(m.curOpEvent) != 0 {
		//需要响应优先级事件,不能摸牌
		// mlog.Debug("优先级事件队列不为空，不能摸牌，等待响应")
		return
	}
	if m.currOpSeatNo != m.nextSeatNoTmp {
		m.currOpSeatNo = m.nextSeatNoTmp
		mlog.Debug("pass,桌子号[%d],用户id[%d],seatNo[%d]摸牌", m.TableNo(), user.UserID, m.currOpSeatNo)
		m.TakeCard(m.currOpSeatNo, 1)
		m.BroadcastWait()
		m.remindInfo(m.currOpSeatNo)
	}

}

//选定缺
func (m *MahjongSC) ChooseColor(user *gamelogic.GameUser, OpValue []int32) {
	mlog.Debug("玩家操作选定缺,桌子号[%d],用户id[%d],SeatNo[%d]", m.TableNo(), user.UserID, user.SeatNo)
	ud := m.getUserInfoBySeatNo(user.SeatNo)
	if ud == nil {
		m.rspError(user, 1, "User data is not exist")
		return
	}
	if ud.color > 0 {
		m.rspError(user, 1, "Repeat Op")
		return
	}
	if len(OpValue) < 1 {
		m.rspError(user, 1, "Op data is empty")
		return
	}
	if OpValue[0] < mahjong.ColorTong || OpValue[0] > mahjong.ColorWan {
		m.rspError(user, 1, "Op data error")
		return
	}

	ud.color = OpValue[0]
	//广播给所有用户玩家已经操作过了
	m.BoardcastChooseColor(user.SeatNo)
	mlog.Debug("玩家操作选定缺 成功,桌子号[%d],用户id[%d],SeatNo[%d]", m.TableNo(), user.UserID, user.SeatNo)
	m.GameLog.AddGameStep("定缺,SeatNo[%d],花色[%s]", ud.seatNo, mahjong.ColorString(ud.color))
}

//选倍数
func (m *MahjongSC) ChooseRate(user *gamelogic.GameUser, OpValue []int32) {
	mlog.Debug("玩家操作选倍数,桌子号[%d],用户id[%d],SeatNo[%d]", m.TableNo(), user.UserID, user.SeatNo)
	ud := m.getUserInfoBySeatNo(user.SeatNo)
	if ud == nil {
		m.rspError(user, 1, "User data is not exist")
		return
	}
	if ud.color > 0 {
		m.rspError(user, 1, "Repeat Op")
		return
	}
	if len(OpValue) < 1 {
		m.rspError(user, 1, "Op data is empty")
		return
	}
	if OpValue[0] <= 0 {
		m.rspError(user, 1, "Op data is error")
		return
	}
	ud.rate = OpValue[0]

	//广播给所有用户玩家已经操作过了
	m.BoardcastChooseRate(user.SeatNo)
	mlog.Debug("玩家操作选倍率 成功,桌子号[%d],用户id[%d],SeatNo[%d],Rate[%d]", m.TableNo(), user.UserID, user.SeatNo, ud.rate)
	m.GameLog.AddGameStep("选倍率,SeatNo[%d],倍率[%d]", ud.seatNo, ud.rate)
}

//摸牌，number发牌数量
func (m *MahjongSC) TakeCard(seatNo, number int32) {

	ud := m.getUserInfoBySeatNo(seatNo)
	if ud == nil {
		mlog.Error("用户数据不存在,SeatNo[%d]", seatNo)
		return
	}
	mlog.Debug("玩家摸牌,桌子号[%d],用户id[%d],SeatNo[%d],Number[%d]", m.TableNo(), ud.User.UserID, seatNo, number)
	if m.randGroup.Len() <= int(m.endCardSize) {
		mlog.Debug("桌子号[%d],用户id[%d],牌池用完，可以结束游戏", m.TableNo(), ud.User.UserID)
		m.isEnd = true
		m.resetMopai()
		return
	}
	cd := m.randGroup.Deal(1)
	if cd.Len() == 0 {
		mlog.Error("摸牌失败,桌子号[%d],用户id[%d],SeatNo[%d]", m.TableNo(), ud.User.UserID, seatNo)
		return
	}
	ud.moPai = cd.Cards[0]
	mlog.Debug("玩家摸牌 成功,桌子号[%d],用户id[%d],SeatNo[%d],Number[%d]", m.TableNo(), ud.User.UserID, seatNo, number)
	m.GameLog.AddGameStep("摸牌,SeatNo[%d],牌[%s]", seatNo, mahjong.CardString(ud.moPai))
	m.BoardcastTakeCard(cd, seatNo)
	m.updateCanOpInfo(seatNo, ud.moPai)
	m.nextSeatNoTmp = seatNo
	// m.remindInfo(seatNo)
	m.resetGangHuFlag()
}

// 重置摸牌
func (m *MahjongSC) resetMopai() {
	for _, ud := range m.gameUserData {
		ud.moPai = 0
	}
}

//骰子
func (m *MahjongSC) PlayDice(diceNum int32) {
	if m.gameState != PROGRESS_DICE {
		mlog.Error("游戏不处于丢色子阶段，操作失败 桌子号[%d]", m.TableNo())
		return
	}
	m.dice = make([]int32, 0)
	for i := int32(0); i < diceNum; i++ {
		nSource := rand.NewSource(time.Now().UnixNano() - int64(i))
		rRand := rand.New(nSource)
		nDice := rRand.Intn(6)
		nDice += 1
		m.dice = append(m.dice, int32(nDice))
	}
	//广播丢色子
	mlog.Debug("丢色子 成功 桌子号[%d],Dice%v", m.TableNo(), m.dice)
	m.GameLog.AddGameStep("丢色子%v", m.dice)
	m.BroadcastDice()

}

// 发牌
func (m *MahjongSC) Deal() {
	mlog.Debug("开始发牌 桌子号[%d]", m.TableNo())
	if m.gameState != PROGRESS_DEAL {
		mlog.Error("游戏不处于发牌阶段，操作失败,桌子号[%d]", m.TableNo())
		return
	}
	config := m.Configs.GetConfig("gameserver")
	cardnumber := config.GetConfigInt("gameprivate", "cardnumber")
	if cardnumber <= 0 {
		cardnumber = 13
	}
	mahJongString := fmt.Sprintf("------------ tableno %d --------------\n", m.TableNo())
	for _, ud := range m.gameUserData {
		if ud.seatNo == m.buttonSeatNo {
			cards := m.randGroup.Deal(cardnumber)
			m.BoardcastDeal(cards, ud.seatNo)
			sort.Sort(cards)
			ud.handCards = cards
			m.GameLog.AddGameStep("庄家发牌：%s", cards.String())

		} else {
			cards := m.randGroup.Deal(cardnumber)
			sort.Sort(cards)
			ud.handCards = cards
			m.BoardcastDeal(cards, ud.seatNo)
			m.GameLog.AddGameStep("闲家发牌：%s", cards.String())
		}
		mahJongString += fmt.Sprintf("%v\n", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ud.handCards.Cards)), ","), "[]"))
	}
	mahJongString += fmt.Sprintf("%v\n", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(m.randGroup.GetCardGroup().Cards)), ","), "[]"))
}

//根据座位号获取用户信息
func (m *MahjongSC) getUserInfoBySeatNo(seatno int32) *STUserData {
	for _, urd := range m.gameUserData {
		if urd.User.SeatNo == seatno {
			return urd
		}
	}
	return nil
}

//是否响应优先级事件
func (m *MahjongSC) isAnswerEvent() bool {
	//没有优先级事件，不用响应
	if len(m.curOpEvent) == 0 {
		return false
	}
	//优先级事件没有处理完，不用响应
	for _, data := range m.curCanOpUser {
		if data {
			return false
		}
	}
	return true
}

//是否有优先级事件;curCanOpUser有事件玩家
func (m *MahjongSC) haveEvent() (b bool) {
	defer func() {
		//优先级事件已经处理完
		if !b {
			m.curCanOpUser = make(map[int32]bool, 0)
			m.resetTimmerEvent()
		}
	}()
	//优先级事件没有处理完 返回ture继续响应
	mlog.Debug("%v", m.curCanOpUser)
	for _, data := range m.curCanOpUser {
		if data {
			return true
		}
	}
	return false
}

// 减掉河牌
func (m *MahjongSC) deleteRiverCard(seatNo int32) {

	ud := m.getUserInfoBySeatNo(seatNo)
	if ud != nil {
		if !ud.isDeleteRiver {
			srcLength := ud.riverCards.Len()
			if srcLength <= 1 {
				ud.riverCards = new(mahjong.CardGroup)
			} else {
				ud.riverCards.Cards = ud.riverCards.Cards[:srcLength-1]
			}
			ud.isDeleteRiver = true
		}
	}
}

// 设置河牌flag 判断是否需要删除河牌
func (m *MahjongSC) resetRiverFlag() {
	for _, ud := range m.gameUserData {
		ud.isDeleteRiver = false
	}
}

//响应优先级事件
func (m *MahjongSC) AnswerEvent() {
	if !m.isAnswerEvent() {
		return
	}
	//处理完优先级事件后清空队列
	defer func() {
		m.curOpEvent = make([]*STUserOpEvent, 0)
		m.resetRiverFlag()
		m.resetTimmerOp()
		m.resetTimmerEvent()
		//mlog.System("AnswerEvent defer 桌子号[%d]", m.TableNo())
	}()
	mlog.Debug("响应优先级事件,桌子号[%d]", m.TableNo())
	//获取最大优先级事件，同时存在事件，胡碰、胡杠、胡吃、胡吃碰、胡吃杠、吃碰、吃杠、
	bHaveHu := false
	curHuSeatNo := m.currOpSeatNo
	huCount := 0
	nextButton := int32(0)
	for _, event := range m.curOpEvent {
		if event.eventId/10 == 7 || event.eventId == 7 {
			mlog.Debug("响应优先级事件 胡牌,桌子号[%d],SeatNo[%d],Card[%v]", m.TableNo(), event.opSeatNo, event.card)
			m.GameLog.AddGameStep("胡牌,SeatNo[%d],牌[%s]", event.opSeatNo, mahjong.CardString(event.card[0]))
			bHaveHu = true
			curHuSeatNo = event.opSeatNo
			// m.nextButtonSeatNo = event.opSeatNo
			m.BoardScenceExchange(event.opSeatNo, event.srcSeatNo, event.eventId, 1, event.card) //结算广播
			m.SettleMentIntimeHu(event.opSeatNo, event.srcSeatNo, event.eventId, event.card[0])  //结算
			if event.eventId == OP_HU_QIANGGANG {
				m.SettleMentIntimeHJZY(event.opSeatNo, event.srcSeatNo, event.card[0])
				m.deleteQiangGang(event.srcSeatNo, event.card[0])
			} else {
				//减掉河牌
				m.deleteRiverCard(event.srcSeatNo)
			}
			huCount += 1
			nextButton = event.srcSeatNo
		}
	}
	//如果有人胡，后面的不用再处理,抢杠呢
	if bHaveHu {
		// 一炮多响，庄家为点炮玩家，否则第一个胡牌玩家
		if huCount > 1 {
			mlog.Debug("一炮多响，庄家座位号为点炮玩家,桌子号[%d],SeatNo[%d]", m.TableNo(), nextButton)
			m.nextButtonSeatNo = nextButton
		}

		m.canDihu = false
		// m.resetGangHuFlag()
		m.checkGangHu.bGang = false
		// 重置可操作信息
		for _, reud := range m.gameUserData {
			reud.resetCanOpInfo()
		}
		//游戏是否结束
		if m.endFight() {
			return
		}

		m.currOpSeatNo = m.getNextSeatNo(curHuSeatNo)
		// 胡牌结束动画
		m.animationName = m.RegisterAnimation(ANIMATION_HUPAI, 1, func() {
			m.bOpenTimer = true
			m.TakeCard(m.currOpSeatNo, 1)
			m.BroadcastWait()
			m.remindInfo(m.currOpSeatNo)
			m.updateGameInfo()
		})
		m.bOpenTimer = false

		// m.TakeCard(m.currOpSeatNo, 1)
		// m.BroadcastWait()
		// m.remindInfo(m.currOpSeatNo)
		//抢杠，抢杠单独处理
		return
	}
	bHaveGang := false
	for _, event := range m.curOpEvent {

		if event.eventId/10 == 4 {
			mlog.Debug("响应优先级事件 杠牌,桌子号[%d],SeatNo[%d],nType[%d],Card[%v]", m.TableNo(), event.opSeatNo, event.eventId, event.card)
			ud := m.getUserInfoBySeatNo(event.opSeatNo)
			if ud == nil {
				mlog.Error("用户数据查询失败,桌子号[%d],", m.TableNo())
				return
			}
			bHaveGang = true
			m.canDihu = false
			if m.randGroup.Len() <= int(m.endCardSize) {
				break
			}
			m.setGangHuFlag(event.opSeatNo, event.srcSeatNo)
			ud.checkErrorHandSize()
			if event.eventId == OP_YINGANG {
				ud.handCards.DeleteCard(event.card[0], 3)
			} else if event.eventId == OP_MINGGANG {
				// 只有这一种杠可以抢杠
				if ud.moPai <= 0 {
					mlog.Error("杠牌错误,非法摸牌,桌子号[%d],", m.TableNo())
				}
				ud.handCards.AddCard(ud.moPai)
				ud.moPai = 0
				ud.handCards.DeleteCard(event.card[0], 1)

				//去掉碰牌
				for _, delPeng := range ud.pengCards {
					if delPeng.pengCard == event.card[0] {
						delPeng.valid = false
						break
					}
				}
				//
				ud.checkErrorHandSize()
				gangsize := len(ud.gangCards)
				if gangsize > 0 {
					ud.gangCards[gangsize-1].valid = true
				}
				m.GameLog.AddGameStep("杠牌,SeatNo[%d],牌[%s]", ud.seatNo, mahjong.CardString(event.card[0]))

				m.BoardScenceExchange(event.opSeatNo, event.srcSeatNo, event.eventId, -1, event.card) //结算广播
				m.SettleMentIntimeGang(event.opSeatNo, event.srcSeatNo, event.eventId, event.card[0]) //结算
				//减掉河牌
				m.deleteRiverCard(event.srcSeatNo)
				// 重置可操作信息
				for _, reud := range m.gameUserData {
					reud.resetCanOpInfo()
				}
				m.currOpSeatNo = event.opSeatNo

				// 检测玩家能否抢杠,有强杠胡 不能摸牌，其他玩家可以随时胡
				if m.checkQianggang(ud.seatNo, event.card[0]) {
					mlog.Debug("有玩家可以抢杠胡,桌子号[%d],", m.TableNo())
					// 设置定时器
					m.BroadcastWait()
					m.nextSeatNoTmp = m.getNextSeatNo(m.currOpSeatNo) //抢杠玩家过牌后继续摸牌
					m.qianggang.Card = event.card[0]
					m.qianggang.SeatNo = ud.seatNo
					m.qianggang.Score = 0
					for _, qhud := range m.gameUserData {
						if qhud.isHu {
							continue
						}
						if qhud.seatNo == event.opSeatNo {
							continue
						}
						m.qianggang.Score += m.baseScore
					}
					break
				} else {
					if m.endFight() {
						return
					}
					m.TakeCard(m.currOpSeatNo, 1)
					m.BroadcastWait()
					m.remindInfo(m.currOpSeatNo)
					break
				}
				break
			}
			ud.checkErrorHandSize()
			gangsize := len(ud.gangCards)
			if gangsize > 0 {
				ud.gangCards[gangsize-1].valid = true
			}
			m.GameLog.AddGameStep("杠牌,SeatNo[%d],牌[%s]", ud.seatNo, mahjong.CardString(event.card[0]))

			m.BoardScenceExchange(event.opSeatNo, event.srcSeatNo, event.eventId, -1, event.card) //结算广播
			m.SettleMentIntimeGang(event.opSeatNo, event.srcSeatNo, event.eventId, event.card[0]) //结算
			//减掉河牌
			m.deleteRiverCard(event.srcSeatNo)
			// 重置可操作信息
			for _, reud := range m.gameUserData {
				reud.resetCanOpInfo()
			}
			if m.endFight() {
				return
			}
			m.currOpSeatNo = event.opSeatNo
			m.TakeCard(m.currOpSeatNo, 1)
			m.BroadcastWait()
			m.remindInfo(m.currOpSeatNo)
			break
		}
	}
	//如果有人杠，后面的不用再处理
	if bHaveGang {
		return
	}
	bHavePeng := false
	for _, event := range m.curOpEvent {
		if event.eventId == OP_PENG {
			mlog.Debug("响应优先级事件 碰牌,桌子号[%d],SeatNo[%d],Card[%v]", m.TableNo(), event.opSeatNo, event.card)
			bHavePeng = true
			m.canDihu = false
			m.BoardScenceExchange(event.opSeatNo, event.srcSeatNo, event.eventId, -1, event.card)
			m.currOpSeatNo = event.opSeatNo
			m.GameLog.AddGameStep("碰牌,SeatNo[%d],牌[%s]", event.opSeatNo, mahjong.CardString(event.card[0]))

			ud := m.getUserInfoBySeatNo(event.opSeatNo)
			if ud != nil {
				ud.checkErrorHandSize()
				sizePeng := len(ud.pengCards)
				if sizePeng > 0 {
					ud.pengCards[sizePeng-1].valid = true
				}
				ud.handCards.DeleteCard(event.card[0], 2)
			}
			//减掉河牌
			m.deleteRiverCard(event.srcSeatNo)
			m.resetGangHuFlag()
			// m.remindInfo(m.currOpSeatNo)
			m.BroadcastWait()
			m.remindInfo(m.currOpSeatNo)
			break
		}
	}
	//如果有人碰，后面的不用再处理
	if bHavePeng {
		// 最后一张牌碰后不会打牌，故 注释之
		for _, reud := range m.gameUserData {
			reud.resetCanOpInfo()
		}
		return
	}
	//处理吃的逻辑
	return
}

//即时结算
func (m *MahjongSC) SettleMentIntimeHu(nSeatNo, srcSeatNo, nType, enentCard int32) {
	mlog.Debug("即时胡牌结算,桌子号[%d],SeatNo[%d],srcSeatNo[%d],Type[%d]", m.TableNo(), nSeatNo, srcSeatNo, nType)
	ud := m.getUserInfoBySeatNo(nSeatNo)
	if ud == nil {
		mlog.Error("桌子号[%d],用户[座位号：%d]数据不存在。", m.TableNo(), nSeatNo)
		return
	}
	ud.isEnd = true
	pcg := make([]int32, 0)
	gcg := make([]int32, 0)
	for _, data := range ud.pengCards {
		pcg = append(pcg, data.pengCard)
	}
	for _, data := range ud.gangCards {
		gcg = append(gcg, data.gangCard)
	}
	checkHand := new(mahjong.CardGroup)
	for _, cd := range ud.handCards.Cards {
		checkHand.AddCard(cd)
	}
	ud.huCards = append(ud.huCards, enentCard)
	checkHand.AddCard(enentCard)
	// mlog.Debug("胡牌[%s]", checkHand.String())
	mapCardType := checkHand.GetCardType(pcg, gcg, m.isDYJ)
	if mapCardType[mahjong.CardTypeJGD] > 0 {
		mapCardType[mahjong.CardTypeDDZ] = 0
	}
	gangRate := ud.getRateGang(m.rateSettleMent.GANG, enentCard)
	if mapCardType[mahjong.CardTypeLQD] > 0 {
		mapCardType[mahjong.CardTypeQDZ] = 0
		gangRate -= 1
	}
	siceHuTypes := make([]int32, 0)
	for htype, ct := range mapCardType {
		if ct <= 0 {
			continue
		}
		siceHuTypes = append(siceHuTypes, htype)
	}
	cardRate, strCardType := m.getCardRate(mapCardType)
	// mlog.Debug("结算胡 %d,%s", cardRate, strCardType)
	// gangRate := ud.getRateGang(m.rateSettleMent.GANG, enentCard)

	config := m.Configs.GetConfig("gameserver")
	if nType == OP_HU_QIANGGANG {
		// 抢杠

	} else if nType == OP_HU_GANG {
		cardRate += m.rateSettleMent.GSH
		strCardType += "杠上花,"
	} else if nType == OP_HU_GANGPAO {
		cardRate += m.rateSettleMent.GSP
		strCardType += "杠上炮,"
	} else if nType == OP_HU_ZIMO {
		// 自摸加番
		if config.GetConfigInt("gameprivate", "havezimojiafan") == 1 {
			cardRate += m.rateSettleMent.ZIMO
			strCardType += "自摸,"
		}
	}

	// 天胡
	if config.GetConfigInt("gameprivate", "havetianhu") == 1 {
		if m.outcardTimes == 0 {
			cardRate += m.rateSettleMent.TH
			siceHuTypes = append(siceHuTypes, mahjong.CardTypeTH)
			strCardType += "天胡,"
		}
	}
	// 地胡
	if config.GetConfigInt("gameprivate", "havedihu") == 1 {
		if m.outcardTimes > 0 && m.outcardTimes < m.gameUserCount && m.canDihu {
			cardRate += m.rateSettleMent.DH
			siceHuTypes = append(siceHuTypes, mahjong.CardTypeDH)
			strCardType += "地胡,"
		}
	}
	// 海底捞月
	if config.GetConfigInt("gameprivate", "havehaidi") == 1 {
		if m.randGroup.Len() <= int(m.endCardSize) {
			cardRate += m.rateSettleMent.HD
			siceHuTypes = append(siceHuTypes, mahjong.CardTypeHD)
			strCardType += "海底捞月,"
		}
	}
	if len(siceHuTypes) == 0 {
		siceHuTypes = append(siceHuTypes, mahjong.CardTypePH)
	}

	mlog.System("桌子号[%d],结算胡 %d,%s", m.TableNo(), cardRate, strCardType)
	mlog.System("桌子号[%d],结算胡,牌型[%s=%v],牌型番数[%d],跟番[%d]", m.TableNo(), strCardType, siceHuTypes, cardRate, gangRate)
	rateTotal := cardRate + gangRate
	if rateTotal > m.rateSettleMent.RATEMAX {
		rateTotal = m.rateSettleMent.RATEMAX
	}
	x := math.Pow(float64(2), float64(rateTotal))
	cardRate = int32(x)
	//cardRate *= m.getChooseRate()	//TODO 先不要選擇倍
	if cardRate > int32(x) {
		newx := math.Pow(float64(2), float64(m.rateSettleMent.RATEMAX))
		cardRate = int32(newx)
	}
	score := int64(cardRate) * m.baseScore
	mlog.System("桌子号[%d],底分[%d],牌型倍数[%d],得分[%d]", m.TableNo(), m.baseScore, cardRate, score)
	//加减分
	//srcSeatNo==nSeatNo自摸
	// srcUd := new(STUserData)
	if nSeatNo == srcSeatNo || (nType == OP_HU_GANG && config.GetConfigInt("gameprivate", "ganghuall") == 1) {
		for _, ud1 := range m.gameUserData {
			if ud1.isEnd || ud1.seatNo == nSeatNo || ud1.isHu {
				continue
			}
			ud.intimeScore += score
			ud1.intimeScore -= score
		}
	} else {
		srcUd := m.getUserInfoBySeatNo(srcSeatNo)
		if srcUd == nil {
			mlog.Error("桌子号[%d],用户[座位号：%d]数据不存在。", m.TableNo(), srcSeatNo)
			return
		}
		ud.intimeScore += score
		srcUd.intimeScore -= score
	}

	seatNoWin := make([]int32, 0)
	seatNoWin = append(seatNoWin, nSeatNo)
	seatNoLos := make([]int32, 0)

	if nSeatNo == srcSeatNo || (nType == OP_HU_GANG && config.GetConfigInt("gameprivate", "ganghuall") == 1) {
		for _, udLos := range m.gameUserData {
			if udLos.isHu {
				continue
			}
			if udLos.seatNo == nSeatNo {
				continue
			}
			seatNoLos = append(seatNoLos, udLos.seatNo)
			//Los
			udLos.saveInTimeRecord(TYPE_SETTLEMENT_HU_ZIMO, enentCard, cardRate, udLos.intimeScore, seatNoWin, siceHuTypes)
		}
		ud.saveInTimeRecord(TYPE_SETTLEMENT_HU_ZIMO, enentCard, cardRate, ud.intimeScore, seatNoLos, siceHuTypes)
	} else { //点炮
		udLos := m.getUserInfoBySeatNo(srcSeatNo)
		if udLos == nil {
			mlog.Error("桌子号[%d],用户[座位号：%d]数据不存在。", m.TableNo(), srcSeatNo)
			return
		}
		seatNoLos = append(seatNoLos, srcSeatNo)
		udLos.saveInTimeRecord(TYPE_SETTLEMENT_HU_DIANPAO, enentCard, cardRate, udLos.intimeScore, seatNoWin, siceHuTypes)
		ud.saveInTimeRecord(TYPE_SETTLEMENT_HU_DIANPAO, enentCard, cardRate, ud.intimeScore, seatNoLos, siceHuTypes)
	}
	m.gamelogSettle("胡牌:" + strCardType)
	//广播
	m.BoardcastIntimeSettle(TYPE_SETTLEMENT_HU, cardRate)

	m.updateIntimeScore()
}

//即时结算,呼叫转移
func (m *MahjongSC) SettleMentIntimeHJZY(nSeatNo, srcSeatNo, enentCard int32) {
	mlog.Debug("即时结算,呼叫转移,桌子号[%d],SeatNo[%d],srcSeatNo[%d]", m.TableNo(), nSeatNo, srcSeatNo)
	seatNos := make([]int32, 0)
	siceHuTypes := make([]int32, 0)
	for _, ud := range m.gameUserData {
		if ud.seatNo == nSeatNo {
			ud.intimeScore += m.qianggang.Score
			seatNos = append(seatNos, srcSeatNo)
			cardRate := m.qianggang.Score / m.baseScore
			ud.saveInTimeRecord(TYPE_SETTLEMENT_HJZY, enentCard, int32(cardRate), ud.intimeScore, seatNos, siceHuTypes)
			continue
		}
		if ud.seatNo == srcSeatNo {
			ud.intimeScore -= m.qianggang.Score
			seatNos = append(seatNos, nSeatNo)
			cardRate := m.qianggang.Score / m.baseScore
			ud.saveInTimeRecord(TYPE_SETTLEMENT_HJZY, enentCard, int32(cardRate), ud.intimeScore, seatNos, siceHuTypes)
			continue
		}
	}
	m.qianggang.Score = 0
	m.gamelogSettle("呼叫转移")
	m.BoardcastIntimeSettle(TYPE_SETTLEMENT_HJZY, 1)
	m.updateIntimeScore()
}

// 游戏记录 即时结算
func (m *MahjongSC) gamelogSettle(strName string) {
	for _, ud := range m.gameUserData {
		if ud.intimeScore == 0 {
			continue
		}
		m.GameLog.AddGameStep("即时结算,项目[%s],SeatNo[%d],Score[%d]", strName, ud.seatNo, ud.intimeScore)
	}
}

// 游戏记录 手牌
func (m *MahjongSC) gamelogHandCards() {
	for _, ud := range m.gameUserData {
		m.GameLog.AddGameStep("手牌,SeatNo[%d],牌[%s]", ud.seatNo, ud.handCards.String())
	}
}

// gangLast
func (m *MahjongSC) updateLastGangScore(nSeatNo int32) {
	//分数汇总
	ud := m.getUserInfoBySeatNo(nSeatNo)
	if ud != nil {
		ud.gangLast = ud.intimeScore
	}
}

//重置即时计算分数
func (m *MahjongSC) updateIntimeScore() {
	//分数汇总
	for _, ud2 := range m.gameUserData {

		newx := math.Pow(float64(2), float64(m.rateSettleMent.RATEMAX))

		newscore := int64(newx) * m.baseScore
		ud2.intimeScore = ud2.intimeScore * int64(m.getChooseRate())
		if ud2.winScore+ud2.intimeScore >= newscore {
			ud2.intimeScore = newscore
			ud2.winScore = 0
		}
		if ud2.winScore+ud2.intimeScore <= -newscore {
			ud2.intimeScore = -newscore
			ud2.winScore = 0
		}

		ud2.winScore += ud2.intimeScore
		ud2.userMoney += ud2.intimeScore
		ud2.intimeScore = 0

	}
}

//杠，收取刮风下雨分数
func (m *MahjongSC) SettleMentIntimeGang(nSeatNo, srcSeatNo, nType, enentCard int32) {
	mlog.Debug("即时杠结算,桌子号[%d],SeatNo[%d],srcSeatNo[%d],Type[%d]", m.TableNo(), nSeatNo, srcSeatNo, nType)
	//杠，收取刮风下雨分数
	seatNoWin := make([]int32, 0)
	seatNoWin = append(seatNoWin, nSeatNo)
	seatNoLos := make([]int32, 0)
	cardRate := int32(2) //m.getChooseRate()
	switch nType {
	case OP_MINGGANG:
		ud1 := m.getUserInfoBySeatNo(nSeatNo)
		if ud1 == nil {
			mlog.Error("桌子号[%d],玩家[SeatNo:%d]数据不存在", m.TableNo(), nSeatNo)
		}
		for _, ud2 := range m.gameUserData {
			if ud2.isEnd || ud2.seatNo == nSeatNo {
				continue
			}
			ud1.intimeScore += m.baseScore * int64(cardRate)
			ud2.intimeScore -= m.baseScore * int64(cardRate)
			seatNoLos = append(seatNoLos, ud2.seatNo)
			ud2.saveInTimeRecord(TYPE_SETTLEMENT_GF, enentCard, cardRate, ud2.intimeScore, seatNoWin, make([]int32, 0))
		}
		ud1.saveInTimeRecord(TYPE_SETTLEMENT_GF, enentCard, cardRate, ud1.intimeScore, seatNoLos, make([]int32, 0))
		m.gamelogSettle("明杠")
		m.BoardcastIntimeSettle(TYPE_SETTLEMENT_GF, cardRate)
		m.updateLastGangScore(nSeatNo)
		m.updateIntimeScore()
	case OP_YINGANG:
		ud1 := m.getUserInfoBySeatNo(nSeatNo)
		ud2 := m.getUserInfoBySeatNo(srcSeatNo)
		if ud1 == nil || ud2 == nil {
			mlog.Error("桌子号[%d],玩家[SeatNo:%d][SeatNo2:%d]数据不存在", m.TableNo(), nSeatNo, srcSeatNo)
		}
		ud1.intimeScore += m.baseScore * int64(cardRate)
		ud2.intimeScore -= m.baseScore * int64(cardRate)
		seatNoLos = append(seatNoLos, srcSeatNo)
		ud1.saveInTimeRecord(TYPE_SETTLEMENT_XY, enentCard, cardRate, ud1.intimeScore, seatNoLos, make([]int32, 0))
		ud2.saveInTimeRecord(TYPE_SETTLEMENT_XY, enentCard, cardRate, ud2.intimeScore, seatNoWin, make([]int32, 0))
		m.gamelogSettle("引杠")
		m.BoardcastIntimeSettle(TYPE_SETTLEMENT_XY, cardRate)
		m.updateLastGangScore(nSeatNo)
		m.updateIntimeScore()
	case OP_ANGANG:
		ud1 := m.getUserInfoBySeatNo(nSeatNo)
		if ud1 == nil {
			mlog.Error("桌子号[%d],玩家[SeatNo:%d]数据不存在", m.TableNo(), nSeatNo)
		}

		for _, ud2 := range m.gameUserData {
			if ud2.isEnd || ud2.seatNo == nSeatNo {
				continue
			}
			seatNoLos = append(seatNoLos, ud2.seatNo)
			ud1.intimeScore += m.baseScore * int64(cardRate)
			ud2.intimeScore -= m.baseScore * int64(cardRate)

			ud2.saveInTimeRecord(TYPE_SETTLEMENT_XY, enentCard, cardRate, ud2.intimeScore, seatNoWin, make([]int32, 0))
		}
		ud1.saveInTimeRecord(TYPE_SETTLEMENT_XY, enentCard, cardRate, ud1.intimeScore, seatNoLos, make([]int32, 0))
		m.gamelogSettle("暗杠")
		m.BoardcastIntimeSettle(TYPE_SETTLEMENT_XY, cardRate)
		m.updateLastGangScore(nSeatNo)
		m.updateIntimeScore()
	default:
	}
}

//等待操作
func (m *MahjongSC) mahjongWaitOpt() {
	mlog.Debug("循环等待,桌子号[%d]", m.TableNo())
	//mlog.System("=-= =-= =-= mahjongWaitOpt 1 =-= =-= =-=")
	//处理注册事件
	m.AnswerEvent()
	//mlog.System("=-= =-= =-= mahjongWaitOpt 2 =-= =-= =-=")
	//更新信息
	m.updateGameInfo()

	//mlog.System("=-= =-= =-= mahjongWaitOpt 3 =-= =-= =-=")
	//广播玩家后继操作路线
}

//摸牌
func (m *MahjongSC) MoPai(nSeatNo int32) {
	mlog.Debug("摸牌,桌子号[%d],SeatNo[%d]", m.TableNo(), nSeatNo)
	ud := m.getUserInfoBySeatNo(nSeatNo)
	if ud == nil {
		mlog.Error("ud is nil")
		return
	}
	cards := m.randGroup.Deal(1)
	if cards == nil {
		mlog.Error("cards is nil")
		return
	}
	ud.handCards.AddCardGroup(cards)
	m.BoardcastDeal(cards, ud.seatNo)
}

//检测庄家起牌操作信息
func (m *MahjongSC) updateCanOpInfo(seatNo, cd int32) {
	if m.gameState != PROGRESS_FIGHT {
		return
	}
	//mlog.Debug("检测玩家是否需要提示信息,SeatNo[%d]", seatNo)
	ud := m.getUserInfoBySeatNo(seatNo)
	if ud == nil {
		mlog.Error("用户数据不存在,桌子号[%d],SeatNo[%d]", m.TableNo(), seatNo)
		return
	}
	bCanGang := true
	if m.randGroup.Len() <= int(m.endCardSize) {
		bCanGang = false
	}
	ud.resetCanOpInfo()
	if ud.updataCanOpInfo(cd, false, bCanGang) {
		m.curCanOpUser[ud.seatNo] = true
		// mlog.Debug("curCanOpUser*********true*******")
		if ud.isTuoGuan {
			m.BroadcastUserCanOpInfo(ud)
			m.tuoguanOpEvent(ud)
			return
		} else {
			m.BroadcastUserCanOpInfo(ud)
		}
	}
}

// 托管操作
func (m *MahjongSC) tuoguanOpEvent(ud *STUserData) {
	// mlog.Debug("tuoguanOpEvent[%d]", len(ud.canOpEventInfo.OpInfo))
	if m.gameState != PROGRESS_FIGHT {
		return
	}
	config := m.Configs.GetConfig("gameserver")
	for _, event := range ud.canOpEventInfo.OpInfo {
		// 胡
		if config.GetConfigInt("gameprivate", "tuoguanhu") == 1 {
			if event.nOpType/10 == 7 || event.nOpType == 7 {
				str := fmt.Sprintf("HuAutoTuoguanEvent%d_%d_%d", time.Now().UnixNano(), ud.seatNo, m.TableNo())
				ud.timerTimeTuoguanEvent = m.TimerFuncManager.NewNamedTimer(1, func() {
					m.HuCard(ud.User)
					m.mahjongWaitOpt()
				}, str)
				return
			}
		}
	}
	for _, event := range ud.canOpEventInfo.OpInfo {
		// 杠
		if config.GetConfigInt("gameprivate", "tuoguangang") == 1 {
			if event.nOpType/10 == 4 || event.nOpType == 4 {
				str := fmt.Sprintf("GangAutoTuoguanEvent%d_%d_%d", time.Now().UnixNano(), ud.seatNo, m.TableNo())
				ud.timerTimeTuoguanEvent = m.TimerFuncManager.NewNamedTimer(2, func() {
					m.AllGang(ud.User, event.nValue)
					m.mahjongWaitOpt()
				}, str)
				return
			}
		}
	}
	for _, event := range ud.canOpEventInfo.OpInfo {
		//		mlog.Debug("tuoguanpeng[%d]", config.GetConfigInt("gameprivate", "tuoguanpeng"))
		// 碰
		if config.GetConfigInt("gameprivate", "tuoguanpeng") == 1 {
			if event.nOpType == OP_PENG {
				str := fmt.Sprintf("PengAutoTuoguanEvent%d_%d_%d", time.Now().UnixNano(), ud.seatNo, m.TableNo())
				ud.timerTimeTuoguanEvent = m.TimerFuncManager.NewNamedTimer(2, func() {
					m.PengCard(ud.User, event.nValue)
					m.mahjongWaitOpt()
				}, str)
				return
			}
		}
	}
	for _, _ = range ud.canOpEventInfo.OpInfo {
		// 过
		str := fmt.Sprintf("PassAutoTuoguanEvent%d_%d_%d", time.Now().UnixNano(), ud.seatNo, m.TableNo())
		ud.timerTimeTuoguanEvent = m.TimerFuncManager.NewNamedTimer(2, func() {
			m.PassHand(ud.User)
			m.mahjongWaitOpt()
		}, str)
		return
	}
}

//检测所有玩家起牌操作信息
func (m *MahjongSC) updateCanOpInfos(seatNo, cd int32) {

	//mlog.Debug("检测所有玩家是否需要提示信息,cd[%d]", cd)
	m.curCanOpUser = make(map[int32]bool, 0)
	for _, ud := range m.gameUserData {
		if ud.seatNo == seatNo {
			continue
		}
		if ud.isHu {
			continue
		}
		bCanGang := true
		if m.randGroup.Len() <= int(m.endCardSize) {
			bCanGang = false
		}
		ud.resetCanOpInfo()
		if ud.updataCanOpInfo(cd, true, bCanGang) {
			m.curCanOpUser[ud.seatNo] = true
			// if ud.isTuoGuan {
			// 	m.BroadcastUserCanOpInfo(ud)
			// 	m.tuoguanOpEvent(ud)
			// 	return
			// } else {
			// 	m.BroadcastUserCanOpInfo(ud)
			// }
		}
	}
	for nSeatNo, bFlag := range m.curCanOpUser {
		if !bFlag {
			continue
		}
		ud := m.getUserInfoBySeatNo(nSeatNo)
		if ud == nil {
			continue
		}
		if ud.isTuoGuan {
			m.BroadcastUserCanOpInfo(ud)
			m.tuoguanOpEvent(ud)
		} else {
			m.BroadcastUserCanOpInfo(ud)
		}
	}
}

//取消定时器 操作
func (m *MahjongSC) resetTimmerOp() {
	if m.timerTimeOutAction != nil {
		mlog.Debug("取消动作定时器,桌子号[%d],", m.TableNo())
		m.timerTimeOutAction.Stop()
		m.timerTimeOutAction = nil
	}
}

//取消定时器 事件
func (m *MahjongSC) resetTimmerEvent() {
	if m.timerTimeOutEvent != nil {
		mlog.Debug("取消事件定时器,桌子号[%d],", m.TableNo())
		m.timerTimeOutEvent.Stop()
		m.timerTimeOutEvent = nil
	}
}

//设置自动操作
func (m *MahjongSC) autoAction() {
	if m.gameState > PROGRESS_FIGHT {
		return
	}
	if m.isEnd {
		return
	}
	// 关闭定时器
	if !m.bOpenTimer {
		return
	}
	if m.haveEvent() {
		// mlog.Debug("设置事件自动操作,桌子号[%d],%v", m.TableNo(), m.timerTimeOutEvent)
		m.resetTimmerOp()
		if m.timerTimeOutEvent == nil {
			// m.outTimeBegin = int32(time.Now().Unix())
			str := fmt.Sprintf("autoEvent%d_%d", time.Now().UnixNano(), m.TableNo())
			m.timerTimeOutEvent = m.TimerFuncManager.NewNamedTimer(m.getOutTime(), func() {
				for _, ud := range m.gameUserData {
					if m.curCanOpUser[ud.seatNo] {
						ud.resetCanOpInfo()
						m.PassHand(ud.User) //有问题
					}
				}
				m.mahjongWaitOpt()
			}, str)
		}
	} else {
		//mlog.System("设置动作自动操作,%v", m.timerTimeOutAction)
		if m.timerTimeOutAction == nil {
			m.tuogunOutCards()
			str := fmt.Sprintf("autoAction%d_%d", time.Now().UnixNano(), m.TableNo())
			m.timerTimeOutAction = m.TimerFuncManager.NewNamedTimer(m.getOutTime(), m.autoOutCards, str)
		} else {
			mlog.Debug("%v", m.timerTimeOutAction)
		}
	}
}

// 托管自动出牌
func (m *MahjongSC) tuogunOutCards() {
	ud := m.getUserInfoBySeatNo(m.currOpSeatNo)
	if ud != nil {
		if ud.isTuoGuan {
			if m.timerTimeTuoguanAction != nil {
				m.timerTimeTuoguanAction.Stop()
				m.timerTimeTuoguanAction = nil
			}
			strTuoguan := fmt.Sprintf("TuoguanAction%d_%d", time.Now().UnixNano(), m.TableNo())
			m.timerTimeTuoguanAction = m.TimerFuncManager.NewNamedTimer(2, m.tuoguanAutoOutCards, strTuoguan)
		}
	}
}

// 托管自动出牌
func (m *MahjongSC) tuoguanAutoOutCards() {
	if m.gameState == PROGRESS_CHOOSE {
		if m.Configs.GetConfig("gameserver").GetConfigInt("gameprivate", "havedingque") == 1 {
			m.autoDingQue() //自动选定缺
			m.dingQueEndAction()
		} else if m.Configs.GetConfig("gameserver").GetConfigInt("gameprivate", "havechooserate") == 1 {
			m.autoChooseRate()
			m.rateEndAction()
		}
		m.mahjongWaitOpt()
		return
	} else if m.gameState == PROGRESS_FIGHT {
		//该出牌就自动出牌，其他的自动选择pass
		ud := m.getUserInfoBySeatNo(m.currOpSeatNo)
		if ud != nil {
			value := make([]int32, 0)
			if m.Configs.GetConfig("gameserver").GetConfigInt("gameprivate", "tuoguantype") == 1 {
				bestcd := ud.calBestOutCard()
				value = append(value, bestcd)
			} else {
				value = append(value, ud.moPai)
			}
			m.OutCard(ud.User, value)
			m.mahjongWaitOpt()
		}
		return
	}
}

// 自动出牌
func (m *MahjongSC) autoOutCards() {
	if m.gameState == PROGRESS_CHOOSE {
		if m.Configs.GetConfig("gameserver").GetConfigInt("gameprivate", "havedingque") == 1 {
			m.autoDingQue() //自动选定缺
			m.dingQueEndAction()
		} else if m.Configs.GetConfig("gameserver").GetConfigInt("gameprivate", "havechooserate") == 1 {
			m.autoChooseRate() //自动选倍率
			m.rateEndAction()
		}
		m.mahjongWaitOpt()
		return
	} else if m.gameState == PROGRESS_FIGHT {
		//该出牌就自动出牌，其他的自动选择pass
		ud := m.getUserInfoBySeatNo(m.currOpSeatNo)
		if ud != nil {
			if !ud.isTuoGuan {
				ud.isTuoGuan = true
				m.BroadcastTuoGuan(ud.seatNo, ud.isTuoGuan)
			}
			value := make([]int32, 0)
			if m.Configs.GetConfig("gameserver").GetConfigInt("gameprivate", "tuoguantype") == 1 {
				bestcd := ud.calBestOutCard()
				value = append(value, bestcd)
			} else {
				value = append(value, ud.moPai)
			}
			m.OutCard(ud.User, value)
			m.mahjongWaitOpt()
		}
		return
	}
}

//获取超时时间
func (m *MahjongSC) getOutTime() int32 {
	if m.gameState == PROGRESS_CHOOSE {
		return m.outTimeConfig.Choose
	} else if m.gameState == PROGRESS_FIGHT {
		// ud := m.getUserInfoBySeatNo(m.currOpSeatNo)
		// if ud != nil && ud.isTuoGuan && !m.haveEvent() {
		// 	return 2
		// }
		return m.outTimeConfig.Action
	}
	return 15
}

//获取超时时间
func (m *MahjongSC) getOutTimeWait() int32 {
	if m.gameState == PROGRESS_CHOOSE {
		return m.outTimeConfig.Choose
	} else if m.gameState == PROGRESS_FIGHT {
		return m.outTimeConfig.Action
	}
	return 0
}

//自动选定缺
func (m *MahjongSC) autoDingQue() {

	for _, ud := range m.gameUserData {
		if ud.color <= 0 {
			min, color := int32(15), int32(0)
			for i := mahjong.ColorTong; i <= mahjong.ColorWan; i++ {
				count := ud.handCards.CalColorCount(i)
				if ud.moPai/100 == i {
					count += 1
				}
				if min > count {
					min = count
					color = i
				}
			}
			mlog.Debug("自动选择定缺,桌子号[%d],SeatNo[%d],Color[%d]", m.TableNo(), ud.seatNo, color)
			ud.color = color
			ud.isTuoGuan = true
			m.BroadcastTuoGuan(ud.seatNo, ud.isTuoGuan)
			m.GameLog.AddGameStep("托管,SeatNo[%d]", ud.seatNo)
			m.GameLog.AddGameStep("系统自动选择定缺牌,SeatNo[%d],花色[%s]", ud.seatNo, mahjong.ColorString(ud.color))
		}
	}
}

//自动选倍率
func (m *MahjongSC) autoChooseRate() {

	for _, ud := range m.gameUserData {
		if ud.rate <= 0 {
			ud.rate = 1
			mlog.Debug("自动选择倍率,桌子号[%d],SeatNo[%d],Rate[%d]", m.TableNo(), ud.seatNo, ud.rate)
			ud.isTuoGuan = true
			m.BroadcastTuoGuan(ud.seatNo, ud.isTuoGuan)
			m.GameLog.AddGameStep("托管,SeatNo[%d]", ud.seatNo)
			m.GameLog.AddGameStep("系统自动选择倍率,SeatNo[%d],倍率[%d]", ud.seatNo, ud.rate)
		}
	}
}

//自动操作
func (m *MahjongSC) autoOp() {
	mlog.Debug("自动操作,桌子号[%d]", m.TableNo())
}

//操作合法性校验
func (m *MahjongSC) oprateCheckLegal(user *gamelogic.GameUser, gmsg *msg.Message) bool {
	// mlog.Debug("操作合法性校验")
	if m.gameState >= PROGRESS_END {
		return false
	}
	return true
}

//玩家操作
func (m *MahjongSC) actionInfo(user *gamelogic.GameUser, data []byte) {
	//解析数据

	ust := &netproto.SCMJ_Action{}
	err := proto.Unmarshal(data, ust)
	//错误校验
	if err != nil {
		mlog.Error("桌子号[%d],解析用户%v用户错误%v", m.TableNo(), user, err)
		return
	}
	// mlog.Debug("收到用户操作数据，SeatNo[%d],OpType[%d]", user.SeatNo, *(ust.OpType))
	switch *(ust.OpType) {

	case OP_CHOOSEDINGQUE:
		//定缺
		//mlog.System("桌子号[%d],userid[%d],定缺", m.TableNo(), user.UserID)
		m.ChooseColor(user, ust.OpValue)
	case OP_OUTCARD:
		//出牌
		//mlog.System("桌子号[%d],userid[%d],出牌%d", m.TableNo(), user.UserID, ust.OpValue)
		m.OutCard(user, ust.OpValue)
	case OP_GANG:
		//杠
		//mlog.System("桌子号[%d],userid[%d],杠%d", m.TableNo(), user.UserID, ust.OpValue)
		m.AllGang(user, ust.OpValue)
	case OP_PENG:
		//碰牌
		//mlog.System("桌子号[%d],userid[%d],碰牌", m.TableNo(), user.UserID)
		m.PengCard(user, ust.OpValue)
	case OP_HU:
		//胡牌
		//mlog.System("桌子号[%d],userid[%d],胡牌", m.TableNo(), user.UserID)
		m.HuCard(user)
	case OP_PASS:
		//过牌
		m.PassHand(user)
		//mlog.System("桌子号[%d],userid[%d],过牌", m.TableNo(), user.UserID)
	case OP_CHOOSERATE:
		// 选倍率
		m.ChooseRate(user, ust.OpValue)
		//mlog.System("桌子号[%d],userid[%d],选倍率", m.TableNo(), user.UserID)
	}
	m.mahjongWaitOpt()

}

//gang
func (m *MahjongSC) AllGang(user *gamelogic.GameUser, OpValue []int32) {
	mlog.Debug("玩家操作杠,桌子号[%d],用户id[%d],SeatNo[%d]", m.TableNo(), user.UserID, user.SeatNo)
	ud := m.getUserInfoBySeatNo(user.SeatNo)
	//mlog.System("AllGang 桌子号[%d],用户id[%d],手牌[%s],摸牌[%d],%s",  m.TableNo(), user.UserID, ud.handCards.String(), ud.moPai, ud.canOpEventInfo.String())
	if ud == nil {
		m.rspError(user, 1, "User data is not exist")
		return
	}
	//检测玩家碰牌能否杠
	if len(OpValue) < 1 {
		m.rspError(user, 1, "Op data is empty")
		return
	}
	if ud.eventFlag {
		m.rspError(user, 1, "op error")
		return
	}
	if !m.curCanOpUser[user.SeatNo] {
		m.rspError(user, 1, "have no this type op")
		return
	}
	cd := OpValue[0]
	//检测玩家碰牌能否杠

	gangType := ud.findOpInfo(cd, OP_GANG)

	if gangType == OP_ANGANG {
		if m.currOpSeatNo != user.SeatNo {
			m.rspError(user, 1, "OP_ANGANG it's not your turn")
			return
		}
		if !ud.handCards.HaveFourOf(cd) && !(ud.handCards.HaveThreeOf(cd) && ud.moPai == cd) {
			m.rspError(user, 1, fmt.Sprintf("OP_ANGANG error you have not four of %d", cd))
			return
		}
		if !ud.canAnGang(false) {
			m.rspError(user, 1, "OP_ANGANG error")
			return
		}
		m.AnGang(user, cd)
		ud.setEnevtFlag(true)
	} else if gangType == OP_YINGANG {
		if !ud.canYinGang(cd, false) {
			m.rspError(user, 1, "OP_YINGANG error")
			return
		}
		m.YinGang(user, cd)
		ud.setEnevtFlag(true)
	} else if gangType == OP_MINGGANG {
		if m.currOpSeatNo != user.SeatNo {
			m.rspError(user, 1, "OP_MINGGANG it's not your turn")
			return
		}
		if !ud.canMingGang(false) {
			m.rspError(user, 1, "OP_MINGGANG error")
			return
		}
		m.MingGang(user, cd)
		ud.setEnevtFlag(true)
	} else {
		m.rspError(user, 1, "AllGang illegaliful Op")
		return
	}
	if ud.canHu(ud.moPai) || (len(ud.handCards.HaveFour()) > 0 || ud.handCards.HaveThreeOf(ud.moPai)) {
		ud.resetCanOpInfo()
		m.updateCanOpInfo(user.SeatNo, ud.moPai)
	}

}

//
func (m *MahjongSC) dingQueEndAction() {
	//取消定时器
	m.resetTimmerOp()
	//广播定缺信息
	m.BroadCastChooseFlower()
	m.currOpSeatNo = m.buttonSeatNo
	m.nextSeatNoTmp = m.buttonSeatNo
	mlog.Debug("轮到玩家,桌子号[%d],SeatNo[%d]操作", m.TableNo(), m.currOpSeatNo)
	m.ExchangeGameState(PROGRESS_FIGHT)
	//广播玩家可操作信息updataCanOpInfo
	ud := m.getUserInfoBySeatNo(m.buttonSeatNo)
	if ud != nil {
		m.updateCanOpInfo(m.buttonSeatNo, ud.moPai)
	}
	m.remindInfo(m.currOpSeatNo)
}

// 选倍率结束
func (m *MahjongSC) rateEndAction() {
	//取消定时器
	m.resetTimmerOp()
	//广播定缺信息
	m.BroadCastChooseRate()
	m.currOpSeatNo = m.buttonSeatNo
	m.nextSeatNoTmp = m.buttonSeatNo
	mlog.Debug("轮到玩家,桌子号[%d],SeatNo[%d]操作", m.TableNo(), m.currOpSeatNo)
	// //重新发牌
	// for _, bud := range m.gameUserData {
	// 	cg := new(mahjong.CardGroup)
	// 	cg.Cards = append(cg.Cards, bud.handCards.Cards...)
	// 	cg.AddCard(bud.moPai)
	// 	m.BoardcastDeal(cg, bud.seatNo)
	// }
	m.ExchangeGameState(PROGRESS_FIGHT)
	//广播玩家可操作信息updataCanOpInfo
	ud := m.getUserInfoBySeatNo(m.buttonSeatNo)
	if ud != nil {
		m.updateCanOpInfo(m.buttonSeatNo, ud.moPai)
	}
	m.remindInfo(m.currOpSeatNo)
}

//更新游戏信息
func (m *MahjongSC) updateGameInfo() {
	mlog.Debug("更新游戏信息,桌子号[%d],", m.TableNo())
	switch m.gameState {
	case PROGRESS_CHOOSE:
		if m.Configs.GetConfig("gameserver").GetConfigInt("gameprivate", "havechooserate") == 1 {
			if m.endChooseRate() {
				m.rateEndAction()
			}
		} else if m.Configs.GetConfig("gameserver").GetConfigInt("gameprivate", "havedingque") == 1 {
			if m.endChooseFlower() {
				m.dingQueEndAction()
			}
		}
	case PROGRESS_FIGHT:
		if m.endFight() {
			m.ExchangeGameState(PROGRESS_LIGHT)
			//取消定时器
			// if m.isGameOver {
			// 	return
			// }
			// m.isGameOver = true
			m.resetTimmerOp()
			m.resetTimmerEvent()
			// 打完最后一张牌 进入亮牌延时
			mlog.Debug("战斗结束动画注册,桌子号[%d],time[%d]", m.TableNo(), m.outTimeConfig.FightEnd)
			m.animationName = m.RegisterAnimation(ANIMATION_FIGHTEND, m.outTimeConfig.FightEnd, func() {
				m.BroadcastLight()
				// 亮牌结束等待时间
				mlog.Debug("亮牌结束动画注册,桌子号[%d],time[%d]", m.TableNo(), m.outTimeConfig.Light)
				m.animationName = m.RegisterAnimation(ANIMATION_LIGHT, m.outTimeConfig.Light, func() {
					notHuSize := 0
					for _, ud := range m.gameUserData {
						if !ud.isHu {
							notHuSize++
						}
					}

					if notHuSize >= 2 {
						//查花猪
						m.ExchangeGameState(PROGRESS_PIG)
						m.CheckHuaZhu()
						mlog.Debug("查花猪结束动画注册,桌子号[%d],time[%d]", m.TableNo(), m.outTimeConfig.Huazhu)
						m.animationName = m.RegisterAnimation(ANIMATION_HUAZHU, m.outTimeConfig.Huazhu, func() {
							//查大叫
							m.ExchangeGameState(PROGRESS_DAJIAO)
							m.CheckDajiao()
							mlog.Debug("查大叫结束动画注册,桌子号[%d],time[%d]", m.TableNo(), m.outTimeConfig.Dajiao)
							m.animationName = m.RegisterAnimation(ANIMATION_DAJIAO, m.outTimeConfig.Dajiao, func() {
								//结算
								m.ExchangeGameState(PROGRESS_SETTLE)
								m.SettleMent()
								m.gamelogHandCards()
								// 总分汇总

								m.makeGameScore()
								m.GameLog.GameEnd()
								//m.GameEndFunc
								mlog.Debug("游戏结束动画注册,桌子号[%d],time[%d]", m.TableNo(), m.outTimeConfig.GameEnd)
								m.animationName = m.RegisterAnimation(ANIMATION_END, m.outTimeConfig.GameEnd, m.gameOver)
								m.gameState = PROGRESS_END
							}) //end 大叫
						}) //end 花猪
					} else {
						//结算
						m.ExchangeGameState(PROGRESS_SETTLE)
						m.SettleMent()
						m.gamelogHandCards()
						// 总分汇总
						m.makeGameScore()
						m.GameLog.GameEnd()
						//m.GameEndFunc
						mlog.Debug("游戏结束动画注册,桌子号[%d],time[%d]", m.TableNo(), m.outTimeConfig.GameEnd)
						m.animationName = m.RegisterAnimation(ANIMATION_END, m.outTimeConfig.GameEnd, m.gameOver)
						m.gameState = PROGRESS_END
					}
				}) //end 亮牌
			}) //end 战斗
		}
	default:
		//return
	}
	//注册定时器
	m.autoAction()
}

// 游戏总结算汇总
func (m *MahjongSC) makeGameScore() map[int32]int64 {
	mapScore := make(map[int32]int64, 0)
	// 结算汇总  税收
	for _, ud := range m.gameUserData {
		var flowAmount, betAmount, payAmount int64
		amount, tax := m.GetUserWinAndTax(ud.winScore)
		m.UpdateUserScoreFunc(ud.User.UserID, amount, tax)
		flowAmount = util.ABSInt64(amount)      //流水=派彩
		betAmount = util.ABSInt64(amount) + tax //有效投注=輸贏+稅
		payAmount = betAmount + amount          //返邏金額=有效投注+派彩
		ud.User.FlowAmount += flowAmount
		ud.User.BetAmount += betAmount
		ud.User.PayAmount += payAmount
		m.GameLog.AddGameStep("用戶[%d] 結算金額, win %d, 稅收 %d, 流水 %d, 有效投注 %d, 返還金額 %d", ud.User.UserID, amount, tax, flowAmount, betAmount, payAmount)
		mapScore[ud.User.UserID] = amount

		m.GameLog.Context().UserUpdateScore(ud.User.UserID, amount)
		m.GameLog.Context().UserBetScore(ud.User.UserID, betAmount)
		m.GameLog.Context().UserPayoutScore(ud.User.UserID, payAmount)
		m.GameLog.Context().UserValidScore(ud.User.UserID, flowAmount)
		m.GameLog.Context().UserTaxAmount(ud.User.UserID, tax)
	}
	return mapScore
}

//更新玩家操作
func (m *MahjongSC) updateCurOpMessage(seatNo int32) {

	if len(m.curCanOpUser) == 0 && m.gameState == PROGRESS_FIGHT {
		m.currOpSeatNo = m.getNextSeatNo(seatNo)
		m.TakeCard(m.currOpSeatNo, 1)
		m.BroadcastWait()
		m.remindInfo(m.currOpSeatNo)
	} else {
		m.nextSeatNoTmp = m.getNextSeatNo(seatNo)
	}
}

// 获取下个玩家操作信息
func (m *MahjongSC) getNextSeatNo(seatNo int32) int32 {
	curSeatNo := seatNo
	defer func() {
		mlog.Debug("桌子号[%d],上一个玩家座位号[%d],当前玩家座位号[%d]", m.TableNo(), seatNo, curSeatNo)
	}()
	for i := 0; i < 20; i++ {
		curSeatNo++
		if curSeatNo > 5 {
			curSeatNo = 0
		}
		ud := m.getUserInfoBySeatNo(curSeatNo)
		if ud == nil {
			continue
		}
		if ud.isHu {
			continue
		}
		return curSeatNo
	}
	return curSeatNo
}

//结束选定缺
func (m *MahjongSC) endChooseFlower() bool {
	for _, ud := range m.gameUserData {
		if ud.color <= 0 {
			return false
		}
	}
	mlog.Debug("玩家选择定缺结束,桌子号[%d],", m.TableNo())
	return true
}

// 结束选倍率
func (m *MahjongSC) endChooseRate() bool {
	for _, ud := range m.gameUserData {
		if ud.rate <= 0 {
			return false
		}
	}
	mlog.Debug("玩家选倍率结束,桌子号[%d],", m.TableNo())
	return true
}

//结束战斗
func (m *MahjongSC) endFight() (b bool) {
	defer func() {
		if b {
			m.currOpSeatNo = -1
			// m.isEnd = true
		}
	}()
	huCount := 0
	// 如果打完最后一张牌，还有碰 胡，不能结束游戏
	for _, ud := range m.gameUserData {
		if ud.isHu {
			huCount += 1
		}
		// 玩家需要打完最后一张牌
		if ud.moPai > 0 {
			mlog.Debug("桌子号[%d],seatNo[%d][%d]摸牌不为空，需要打出", m.TableNo(), ud.seatNo, ud.moPai)
			return false
		}
	}
	if int32(huCount) >= m.gameUserCount-1 {
		mlog.Debug("胡牌触发战斗结束,桌子号[%d]", m.TableNo())
		return true
	}

	if int32(m.randGroup.Len()) <= m.endCardSize {
		//保证海底操作
		for _, ud := range m.gameUserData {
			for _, info := range ud.canOpEventInfo.OpInfo {
				if info.nOpType == OP_HU || info.nOpType == OP_PENG || info.nOpType/10 == OP_HU {
					mlog.Debug("桌子号[%d],玩家[%d]有碰或者杠,不能结束游戏", m.TableNo(), ud.seatNo)
					return false
				}
			}
		}
		ud := m.getUserInfoBySeatNo(m.currOpSeatNo)
		if ud != nil {
			if ud.handCards.Len()%3 != 1 {
				mlog.Debug("桌子号[%d],玩家[%d]有碰或者杠2,不能结束游戏", m.TableNo(), ud.seatNo)
				return false
			}
		}

		mlog.Debug("桌子号[%d],剩余牌不足战斗结束", m.TableNo())
		return true
	}
	return false //状态不改变
}

//改变游戏阶段
func (m *MahjongSC) ExchangeGameState(state int32) {
	mlog.Debug("设置游戏状态,桌子号[%d],State[%d]", m.TableNo(), state)
	m.gameState = state
	// m.outTimeBegin = int32(time.Now().Unix())
	m.BoardcastGameState()
	if m.gameState >= PROGRESS_CHOOSE {
		m.BroadcastWait()
	}

}

//查花猪
func (m *MahjongSC) CheckHuaZhu() {
	mlog.Debug("查花猪,桌子号[%d]", m.TableNo())
	for _, ud := range m.gameUserData {
		ud.isHuaZhu = ud.handCards.HaveColor(ud.color)
	}
	if m.isEndCheck {
		return
	}
	//对花猪进行结算
	m.SettleMentHuaZhu()
}

//对花猪进行结算
func (m *MahjongSC) SettleMentHuaZhu() {
	mlog.Debug("花猪结算,桌子号[%d]", m.TableNo())
	//搜索没有胡牌 并且不是花猪的玩家
	siceHuazhu := make([]*STUserData, 0)
	siceNotHuazhu := make([]*STUserData, 0)
	siceSeatNoWin := make([]int32, 0)
	siceSeatNoLos := make([]int32, 0)
	for _, ud := range m.gameUserData {
		if ud.isHu {
			continue
		}
		if ud.isHuaZhu {
			siceHuazhu = append(siceHuazhu, ud)
			siceSeatNoLos = append(siceSeatNoLos, ud.seatNo)
			continue
		}
		siceNotHuazhu = append(siceNotHuazhu, ud)
		siceSeatNoWin = append(siceSeatNoWin, ud.seatNo)
	}
	//退出杠钱backGangScore
	for _, udLos := range siceHuazhu {
		m.backGangScore(udLos)
	}
	if len(siceHuazhu) <= 0 || len(siceNotHuazhu) <= 0 {
		return
	}
	for _, udNot := range siceNotHuazhu {
		for _, udH := range siceHuazhu {
			udNot.intimeScore += 16 * m.baseScore
			udH.intimeScore -= 16 * m.baseScore
		}
	}

	for _, udWin := range siceNotHuazhu {
		m.GameLog.AddGameStep("花猪赢家,SeatNo[%d],得分[%d],手牌[%s]", udWin.seatNo, udWin.intimeScore, udWin.handCards.String())
		udWin.saveInTimeRecord(TYPE_SETTLEMENT_HZ, 0, 16, udWin.intimeScore, siceSeatNoLos, make([]int32, 0))
	}
	for _, udLos := range siceHuazhu {
		m.GameLog.AddGameStep("花猪输家,SeatNo[%d],得分[%d],手牌[%s]", udLos.seatNo, udLos.intimeScore, udLos.handCards.String())
		udLos.saveInTimeRecord(TYPE_SETTLEMENT_HZ, 0, 16, udLos.intimeScore, siceSeatNoWin, make([]int32, 0))
	}
	//广播
	m.BoardcastIntimeSettle(TYPE_SETTLEMENT_HZ, 0)
	m.updateIntimeScore()

	//退出杠钱backGangScore
	// for _, udLos := range siceHuazhu {
	// 	m.backGangScore(udLos)
	// }
	m.isEndCheck = true
}

//查大叫进行结算
func (m *MahjongSC) SettleMentDaJiao() {
	mlog.Debug("查大叫结算,桌子号[%d]", m.TableNo())
	siceDajiao := make([]*STUserData, 0)
	siceNotDajiao := make([]*STUserData, 0)
	for _, ud := range m.gameUserData {
		if ud.isHu {
			continue
		}
		if !ud.isCanTingPai {
			siceDajiao = append(siceDajiao, ud)
			continue
		}
		siceNotDajiao = append(siceNotDajiao, ud)
	}
	//退出杠钱
	for _, udLos := range siceDajiao {
		m.backGangScore(udLos)
	}
	if len(siceDajiao) <= 0 || len(siceNotDajiao) <= 0 {
		return
	}
	for _, udWin := range siceNotDajiao {
		//计算最大听牌倍率
		bigest := m.calBigestHuRate(udWin.seatNo)
		x := math.Pow(2, float64(bigest))
		bigest = int32(x)
		for _, udLos := range siceDajiao {
			score := m.baseScore * int64(bigest)
			udWin.intimeScore += score //听牌最大倍数
			udLos.intimeScore -= score

			seatNoWin := make([]int32, 0)
			seatNoWin = append(seatNoWin, udWin.seatNo)
			seatNoLos := make([]int32, 0)
			seatNoLos = append(seatNoLos, udLos.seatNo)

			udWin.saveInTimeRecord(TYPE_SETTLEMENT_DJ, 0, bigest, score, seatNoLos, make([]int32, 0))
			udLos.saveInTimeRecord(TYPE_SETTLEMENT_DJ, 0, bigest, -score, seatNoWin, make([]int32, 0))
		}
	}
	for _, ud := range siceNotDajiao {
		m.GameLog.AddGameStep("大叫赢家,SeatNo[%d],得分[%d],手牌[%s]", ud.seatNo, ud.intimeScore, ud.handCards.String())
	}
	for _, ud := range siceDajiao {
		m.GameLog.AddGameStep("大叫输家,SeatNo[%d],得分[%d],手牌[%s]", ud.seatNo, ud.intimeScore, ud.handCards.String())
	}
	//广播
	m.BoardcastIntimeSettle(TYPE_SETTLEMENT_DJ, 1)
	m.updateIntimeScore()

	//退出杠钱
	// for _, udLos := range siceDajiao {
	// 	m.backGangScore(udLos)
	// }
	m.isEndCheck = true
}

//查大叫
func (m *MahjongSC) CheckDajiao() {
	mlog.Debug("查大叫,桌子号[%d]", m.TableNo())
	for _, ud := range m.gameUserData {
		ud.isCanTingPai = ud.handCards.CanTingPai()
		if !ud.isCanTingPai {
			ud.isCanTingPai = ud.handCards.CanTingSevenPair()
		}
		if !ud.isCanTingPai {
			mlog.Debug("玩家不能听牌,桌子号[%d],SeatNo[%d],手牌[%s]", m.TableNo(), ud.seatNo, ud.handCards.String())
		}
	}
	if m.isEndCheck {
		return
	}
	m.SettleMentDaJiao()
}

//结算
func (m *MahjongSC) SettleMent() {
	mlog.Debug("结算,桌子号[%d]", m.TableNo())
	m.BroadcastTotalSettleMent()

}

func (m *MahjongSC) gameOver() {
	m.resetTimmerEvent()
	m.resetTimmerOp()
	m.GameEndFunc()
}

//退还所有收到的杠分
func (m *MahjongSC) backGangScore(ud *STUserData) {
	for _, data := range ud.intimeInfo {
		// mlog.Debug("退杠开始[%v]", data)
		if data.valid && data.Score > 0 && (data.Type == TYPE_SETTLEMENT_GF || data.Type == TYPE_SETTLEMENT_XY) {
			data.valid = false
			for _, seatNo := range data.SeatNo {
				// mlog.Debug("返还积分,SeatNo[%d]", seatNo)
				if ud.seatNo == seatNo {
					continue
				}
				udLos := m.getUserInfoBySeatNo(seatNo)
				if udLos == nil {
					continue
				}
				for _, losData := range udLos.intimeInfo {
					if losData.valid && (losData.Type == TYPE_SETTLEMENT_GF || losData.Type == TYPE_SETTLEMENT_XY) && data.card == losData.card {
						losData.valid = false
					}
				}
				//
				score := data.Score / int64(len(data.SeatNo))
				udLos.intimeScore += score
				ud.intimeScore -= score
				// m.GameLog.AddGameStep("退杠,赢家[%d]退%d分,输家[%d]得%d分", ud.seatNo, score, udLos.seatNo, score)
			}
		}
	}
	m.gamelogSettle("退杠")
	m.BoardcastIntimeSettle(TYPE_SETTLEMENT_BACK, 1)
	m.updateIntimeScore()
}
func (m *MahjongSC) resetbackGangScore() {
	for _, ud := range m.gameUserData {
		ud.winScore += ud.backScore
		ud.backScore = 0
	}
}

// 计算玩家听牌提示信息
func (m *MahjongSC) calTingInfo(seatno, outcard int32) []*STRemindInfo {
	// mlog.Debug("计算玩家听牌提示信息,桌子号[%d],SeatNo[%d]", m.TableNo(), seatno)
	mapRemindInfo := make([]*STRemindInfo, 0)
	ud := m.getUserInfoBySeatNo(seatno)
	if ud == nil {
		return mapRemindInfo
	}

	// 计算能听的牌
	for i := int32(10); i < 40; i++ {

		if i/10 == ud.color {
			continue
		}
		if i%10 == 0 {
			continue
		}

		// 复制手牌信息
		handcg := new(mahjong.CardGroup)
		for _, cd := range ud.handCards.Cards {
			handcg.AddCard(cd)
		}
		handcg.AddCard(ud.moPai)

		// 判断牌数量
		if handcg.Len()%3 != 2 {
			mlog.Debug("计算玩家听牌提示信息,数量不对")
			return mapRemindInfo
		}

		// 1.打一张计算能不能胡。 2.不打，计算缺几张
		handcg.DeleteCard(outcard, 1)

		cd := i*10 + (outcard % 100 % 10)
		for _, handcd := range handcg.Cards {
			if mahjong.GetCardColor(handcd) == ud.color {
				return mapRemindInfo
			}
		}
		// mlog.Debug("手牌[%s],检测牌[%d],出牌[%d]", handcg.String(), cd, outcard)
		// 计算能否听牌
		if !handcg.CanHuBase(cd) {
			handcg.AddCard(cd)
			if !handcg.IsSevenPair() {
				continue
			}
		} else {
			handcg.AddCard(cd)
		}
		// mlog.Debug("计算玩家听牌提示信息,听牌,手牌[%s],检测牌[%d],出牌[%d]", handcg.String(), cd, outcard)
		// 计算牌型倍率
		sicePeng := make([]int32, 0)
		siceGang := make([]int32, 0)
		genRate := int32(0) //跟番数
		for _, pengInfo := range ud.pengCards {
			if !pengInfo.valid {
				continue
			}
			sicePeng = append(sicePeng, pengInfo.pengCard)
		}
		for _, gangInfo := range ud.gangCards {
			if !gangInfo.valid {
				continue
			}
			siceGang = append(siceGang, gangInfo.gangCard)
			genRate += 1
		}
		// handcg.AddCard(cd)
		mapCardType := handcg.GetCardType(sicePeng, siceGang, m.isDYJ)
		if mapCardType[mahjong.CardTypeJGD] > 0 {
			mapCardType[mahjong.CardTypeDDZ] = 0
		}
		cardRate, _ := m.getCardRate(mapCardType)
		// mlog.Debug("提示 牌[%s],牌型%v,牌型番数[%d]", handcg.String(), mapCardType, cardRate)
		// 计算跟倍率
		for _, pengInfo := range ud.pengCards {
			if !pengInfo.valid {
				continue
			}
			for _, ratecd := range handcg.Cards {
				if ratecd == pengInfo.pengCard {
					genRate += 1
					break
				}
			}
		}
		// 手牌中单独有四个的
		mapGen := make(map[int32]int32, 0)
		for _, handcd := range handcg.Cards {
			mapGen[mahjong.GetCardValue(handcd)]++
		}
		for _, ct := range mapGen {
			if ct == 4 {
				genRate += 1
			}
		}
		// mlog.Debug("提示 跟番数[%d]", genRate)
		x := math.Pow(float64(2), float64(cardRate+genRate))
		rateTotal := int32(x)
		remindInfo := new(STRemindInfo)
		remindInfo.TingCard = cd
		remindInfo.Rate = rateTotal
		remindInfo.Number = m.calSurPlus(seatno, cd, outcard)
		remindInfo.OutCard = outcard
		mapRemindInfo = append(mapRemindInfo, remindInfo)

	}
	return mapRemindInfo
}

//计算场上可见剩余牌张数
func (m *MahjongSC) calSurPlus(seatno, calValue, outcard int32) int32 {
	count := 0
	for _, ud := range m.gameUserData {
		for _, peng := range ud.pengCards {
			if !peng.valid {
				continue
			}
			if peng.pengCard/10 == calValue/10 {
				count += 3
			}
		}
		for _, river := range ud.riverCards.Cards {
			if river/10 == calValue/10 {
				count += 1
			}
		}
		if ud.seatNo == seatno {
			for _, cd := range ud.handCards.Cards {
				if cd/10 == calValue/10 {
					count += 1
				}
			}
		}
		for _, gang := range ud.gangCards {
			if !gang.valid {
				continue
			}
			if gang.gangCard/10 == calValue/10 {
				count += 4
			}
		}
		for _, hucard := range ud.huCards {
			if hucard/10 == calValue/10 {
				count += 1
			}
		}
		if ud.moPai/10 == calValue/10 {
			count += 1
		}
	}
	if calValue/10 == outcard/10 {
		count += 1
	}
	surplus := 4 - int32(count)
	if surplus < 0 {
		surplus = 0
	}
	return surplus
}

//进入房间
func (m *MahjongSC) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {
	mlog.Debug("进入房间,SeatNo=%d,userId=%d", gameuser.SeatNo, gameuser.UserID)
	m.FixedGameLogicBase.UserEnter(gameuser, tableno, seatno)

	//发送游戏设置
	m.SendGameSettings(gameuser)
	//如果非游戏玩家，发送当前场景
	if m.IsGaming() && !m.IsGameUser(gameuser) {
		mlog.Debug("UserEnter,userId=%d is not game user,send cur sence.", gameuser.UserID)
		m.SendGameScene(gameuser)
	}
}

//用户断线重回
func (m *MahjongSC) UserRecover(gameuser *gamelogic.GameUser) {
	mlog.Debug("UserRecover，用户断线重回")
	m.FixedGameLogicBase.UserRecover(gameuser)

	//发送游戏设置
	m.SendGameSettings(gameuser)

	//发送游戏场景
	if m.IsGaming() {
		m.SendGameScene(gameuser)
	}
}

// 计算听牌玩家能胡的最大牌型，使用前提是玩家可以听牌
func (m *MahjongSC) calBigestHuRate(seatNo int32) int32 {
	ud := m.getUserInfoBySeatNo(seatNo)
	if ud == nil {
		return 0
	}
	if !ud.isCanTingPai {
		return 0
	}
	maxRate := int32(0)
	maxStr := ""
	for i := int32(10); i < 40; i++ {
		if i/10 == ud.color {
			continue
		}
		copyCg := new(mahjong.CardGroup)
		for _, cd := range ud.handCards.Cards {
			copyCg.Cards = append(copyCg.Cards, cd)
		}
		cd := i*10 + 1
		if !copyCg.CanHuBase(cd) {
			// mlog.Debug("查大叫 不能屁胡[%d]", cd)
			copyCg.AddCard(cd)
			if !copyCg.IsSevenPair() {
				// mlog.Debug("查大叫 不能胡七对[%d]", cd)
				continue
			}
		} else {
			copyCg.AddCard(cd)
		}

		sicePeng := make([]int32, 0)
		siceGang := make([]int32, 0)
		for _, pengInfo := range ud.pengCards {
			sicePeng = append(sicePeng, pengInfo.pengCard)
		}
		for _, gangInfo := range ud.gangCards {
			siceGang = append(siceGang, gangInfo.gangCard)
		}
		mapCardType := copyCg.GetCardType(sicePeng, siceGang, m.isDYJ)
		if mapCardType[mahjong.CardTypeJGD] > 0 {
			mapCardType[mahjong.CardTypeDDZ] = 0
		}
		cardRate, strCardType := m.getCardRate(mapCardType)
		// 计算手牌和碰牌能组成的跟
		for _, data := range ud.pengCards {
			if !data.valid {
				continue
			}
			for _, hcd := range copyCg.Cards {
				if hcd == data.pengCard {
					// mlog.Debug("碰牌可以与胡牌或者手牌组成跟[%d]", hcd)
					cardRate += 1
					break
				}
			}
		}
		// 计算手牌中的跟
		mapGen := make(map[int32]int32, 0)
		for _, hcd := range copyCg.Cards {
			mapGen[hcd/10]++
		}
		for _, count := range mapGen {
			if count == 4 {
				// mlog.Debug("手牌组成跟[%d]", mapcd)
				cardRate += 1
			}
		}
		// 计算杠牌的跟
		for _, data := range ud.gangCards {
			if data.valid {
				// mlog.Debug("杠牌组成跟[%d]", data.gangCard)
				cardRate += 1
			}
		}
		//选择最大的番数
		if maxRate < cardRate {
			maxRate = cardRate
			maxStr = strCardType
		}
	}
	mlog.Debug("桌子号[%d],最大胡牌[%s],番数[%d]", m.TableNo(), maxStr, maxRate)
	return maxRate
}

// 检测玩家能否抢杠
func (m *MahjongSC) checkQianggang(seatNo, cd int32) bool {
	bCanhu := false
	for _, ud := range m.gameUserData {
		// 过滤自己
		if seatNo == ud.seatNo {
			continue
		}
		if ud.isHu {
			continue
		}
		if ud.updataCanQiang(cd) {
			m.curCanOpUser[ud.seatNo] = true
			bCanhu = true
			m.BroadcastUserCanOpInfo(ud)
			if ud.isTuoGuan {
				m.tuoguanOpEvent(ud)
			}
		}
	}
	return bCanhu
}

// 删除杠
func (m *MahjongSC) deleteQiangGang(seatNo, card int32) {
	ud := m.getUserInfoBySeatNo(seatNo)
	if ud == nil {
		return
	}
	// 去杠
	for _, gp := range ud.gangCards {
		if gp.gangCard == card {
			gp.valid = false
		}
		break
	}
	// 加碰
	for _, pp := range ud.pengCards {
		if pp.pengCard == card {
			pp.valid = true
		}
		break
	}
}

func (m *MahjongSC) getChooseRate() int32 {
	if len(m.gameUserData) > 2 {
		return 1
	}
	rate := int32(1)
	for _, ud := range m.gameUserData {
		rate *= ud.rate
	}
	if rate <= 0 {
		rate = 1
	}
	return rate
}
