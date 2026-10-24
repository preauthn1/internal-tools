package bjl

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// GameStep
const (
	Bet            = iota // 下注狀態
	BetEnd                // 下注結束動畫
	KaiJiang              // 開獎狀態
	KaiJiangAction        // 開獎動畫
	EndAction             // 結算動畫
	Wait                  // 等待遊戲開始
	DealAction            // 下注前發牌動畫
)

// Animates
const (
	GameDealAction     = "deal"
	GameKaiJiangAction = "kaijiang"
	GameBetEndAction   = "betend"
	GameEndAction      = "end"
)

// 符號
const (
	ZHUANG     = 1  // 莊
	XIAN       = 2  // 閒
	HE         = 3  // 和
	ZHUANGTW   = 4  // 莊天王
	XIANTW     = 5  // 閒天王
	ZHUANGDUI  = 6  // 莊對
	XIANDUI    = 7  // 閒對
	ZHUANGSHUN = 8  // 莊順
	XIANSHUN   = 9  // 閒順
	TONGDIANHE = 10 // 同點和
)

// 百家樂玩法
const (
	BJLModeStandard int = iota // 百家樂
	BJLModeLucky6              // 幸運6百家樂
)

// ?
const (
	OddZhuang6 = 15 // 幸運6百家樂，莊家勝6點倍率 1+0.5 倍，
)

// 第一個用來佔位
var resultstring = []string{"", "庄", "闲", "和", "庄天王", "闲天王", "庄对", "闲对", "庄顺", "闲顺", "同点和"}

// ?
const (
	CtrlBetLimitDefault int64 = 100 // 監控名單預設注額
)

type card struct {
	color  int32
	number int32
}

func (c *card) String() string {
	return fmt.Sprintf("%d", c.number)
}

type userScore struct {
	userdata []*gamelogic.GameUser
}

// BaiJialeLogic ?
type BaiJialeLogic struct {
	*gamelogic.ComputeGameLogicBase
	poker              poker.Poker
	serverfile         *os.File
	opttimer           *event.TimerFunc
	cards              *poker.CardGroup             // 所有的牌
	bettotal           int64                        // 下注總額
	bettotalrobot      int64                        // 機器人下注總額
	bettypetotal       []*bettype                   // 下注分類匯總
	bettyperobot       []*bettype                   // 機器人下注分類匯總
	betvalue           *betconfig                   // 下注配置文件
	userdata           []*baijialeuserdata          // 用戶數據
	cardzhuang         []*card                      // 開獎的莊家牌
	cardxian           []*card                      // 開獎的閒家牌
	gamestate          int32                        // 遊戲狀態
	bettime            int32                        // 下注時間
	jiangtype          []int32                      // 開獎類型
	zuhetype           uint32                       // 開獎的組合
	actionStr          string                       // 當前動畫
	kaijiangactiontime int32                        // 開獎動畫超時
	endactiontime      int32                        // 用戶結算動畫超時
	readytime          int32                        // 用戶準備時間
	zhuangID           int32                        // 當前莊家ID
	zhuanglimit        int64                        // 上莊限制
	lastturn           int32                        // 莊剩餘的局數
	maxturn            int32                        // 上莊的最大局數
	zhuanglist         []int32                      // 上莊列表
	betcount           []int32                      // 下注數量
	zhuangWin          int64                        // 莊家贏錢
	serverid           string                       // 服務器ID
	winuserinfo        []int32                      // 贏錢最多的前三家userid
	resultcardlist     *netproto.BJL_ResultCardList // 開獎的歷史牌
	zhushui            int32                        // 注水率
	maxclearcount      int32                        // 清理數據上限
	betend             int32                        // 結束下注動畫
	maxrobotzhuang     int32                        // 最多有好多個機器人上莊
	robotxiazhuang     int32                        // 機器人下莊概率
	robotzhuang        int32                        // 沒有人當莊的時候機器人當莊的概率
	zhushuimoney       int64                        // 注水錢
	poolline           *poollineconfig              // 水線配置文件
	betamode           int64                        // 是否開啟測試中程式
	startsyspool       int64                        // 開始時的金幣數量
	startmoney         int64                        // 上莊時的金幣
	nextjiangtype      int32                        // 下次開獎類型
	systotalwin        int64                        // 系統收入
	systaxincome       int64                        // 系統稅收
	syslastwin         int64                        // 上一局收益
	currtax            int64                        // 當前稅收
	sametimebetmode    int64                        // 是否開啟庄閒同時下注
	userwinlose        int64                        // 玩家輸贏
	robotwinlose       int64                        // 機器人輸贏
	betlimitmoney      int64                        // 下注限制
	publiczhuang       int64                        // 是否開放用戶上庄
	userscore          *userScore                   // 用戶數據
	luckystar          int32                        // 幸運星玩家id
	luckyfaceid        int32                        // 幸運星玩家頭像id
	luckystarround     int64                        // 幸運星持續局數
	luckystartimes     int64                        // 幸運星第幾局
	notbettimes        int64                        // 沒下注持續局數
	richestuserinfo    []int32                      // 富豪玩家的userid
	richestList        []*netproto.BJL_Richest      // 富豪榜
	stageuserinfo      []int32                      // 標誌玩家的userid
	roundid            int32                        // 當局局數
	zhuangjawin        int32                        // 莊贏局數
	xianwin            int32                        // 閑贏局數
	hewin              int32                        // 和局局數
	newcard            bool                         // 重新發牌
	zhuangxianrange    []int32                      // 莊閑下注限制
	herange            []int32                      // 和下注限制
	duizirange         []int32                      // 對子下注限制
	lastGameUserBet    []*baijialeuserdata          // 上局下注玩家
	//bowProvider        *gamelogic.BlackOrWhiteUserProvider   // 需要放分或殺分的列表
	assignRewardDetail *netproto.SuperToolAssignRewardDetail // 指定開獎
	execAssignReward   bool                                  // true 執行指定開獎
	fangFenTarget      int64                                 // 放分金額
	fangFenComplete    int64                                 // 己放分金額
	execFangFen        bool                                  // true 執行系統放分
	nextzuhetype       int32                                 // 下個組合
	nextzuhetype2      int32                                 // 下個組合2
	SuperToolService   *gamelogic.SuperToolBrGameServiceBase // 超端服務
	ctrlbetlimit       int64                                 // 監控名單會控制的最低注額
	BJLMode            int                                   // 百家樂玩法
	DealActionTime     int32                                 // 下注前發牌動畫時間
}

// NewBaiJiaLeLogic 新建一個百家樂實例
func NewBaiJiaLeLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *BaiJialeLogic {
	glc := new(BaiJialeLogic)
	glc.ComputeGameLogicBase = gamelogic.NewComputeGameLogicBase(t, tfManager, cfgs)
	glc.getbjlconfig(cfgs)
	glc.SuperToolService = glc.GetSuperToolService().(*gamelogic.SuperToolBrGameServiceBase)
	glc.gamestate = Wait
	glc.zhuangID = 0
	glc.nextjiangtype = 0
	glc.syslastwin = 0
	glc.systotalwin = 0
	glc.systaxincome = 0
	glc.roundid = 0
	glc.zhuangjawin = 0
	glc.xianwin = 0
	glc.hewin = 0
	glc.newcard = false
	glc.resultcardlist = new(netproto.BJL_ResultCardList)
	glc.nextzuhetype = 0
	glc.nextzuhetype2 = 0
	//glc.getLastPool()
	glc.startsyspool = glc.SysPool.GetPool()
	glc.addCmd()

	glc.RegisterAfterGameStartFunc(func() {
		glc.clearGameData() //清除数据
		glc.getAssignReward()
		//glc.VipUserService.UpdateVIPUser()
		//glc.getTableBlackUser()
		//glc.getTableWhiteUser()
		glc.fapai()
		glc.getLuckyStar()
		glc.getRichestList()
		glc.getGameRound()
		glc.huangZhuang()
		if glc.publiczhuang == 1 {
			glc.GameLog.AddGameStep("庄家%d，剩余%d", glc.zhuangID, glc.lastturn)
		}
		glc.zhuShui()
		// glc.GameLog.AddGameStep("水池狀態 : " + glc.SysPool.ToString())
		glc.broadCastZhuangInfo()
		// 改為開始下注前增加發牌動畫
		//glc.startBet()
		glc.dealBeforeBet()
	})
	glc.RegisterBeforeGameEndFunc(func() {
		glc.saveLastGameUserBet()

	})
	return glc
}

// 發牌
func (bjl *BaiJialeLogic) fapai() {
	mlog.Debug("roundID:%d", bjl.roundid)
	if bjl.roundid == 0 || bjl.newcard {
		bjl.poker = poker.NewRandomPoker(false, false, 6) // 初始化一副牌
		bjl.poker.Shuffle()
		// 直接發所有的牌
		bjl.cards = bjl.poker.Deal(0, 312)
		bjl.newcard = false // 發完牌了，重置狀態
		mlog.Debug("發牌:%t", bjl.newcard)
	} else {
		bjl.cards.ShuffleOnce()
		//bjl.poker.Shuffle()
	}
	mlog.Debug("len(card):%d", len(bjl.cards.Cards))
}

// 添加命令行
func (bjl *BaiJialeLogic) addCmd() {
	command.AddCommand(NewCmd(bjl))
}

// 廣播莊家
func (bjl *BaiJialeLogic) broadCastZhuangInfo() {
	bjl.BoardcastMessage(bjl.getZhuangInfoMsg())
	bjl.BoardcastMessage(bjl.getRichestListMsg())
	bjl.BoardcastMessage(bjl.getLuckyStarMsg())
}

// 注水到下注裡面去
func (bjl *BaiJialeLogic) zhuShui() {
	bjl.zhushuimoney = 0
	//bjl.GameLog.AddGameStep("系統水池%d, 贏錢線%d", bjl.getSysPool(), bjl.getWinLine())
}

// UserLeave 用戶離開
func (bjl *BaiJialeLogic) UserLeave(gameuser *gamelogic.GameUser) bool {
	if bjl.gamestate != Wait {
		if bjl.getUserDataByID(gameuser.UserID) != nil {
			return false
		}

		if bjl.zhuangID == gameuser.UserID {
			return false
		}
	}

	bjl.delFromZhuangList(gameuser.UserID)

	bjl.ChangeUserStateEnd(gameuser)
	return bjl.ComputeGameLogicBase.UserLeave(gameuser)
}

// UserFree 釋放用戶實例
func (bjl *BaiJialeLogic) UserFree(gameuser *gamelogic.GameUser) {
	if bjl.gamestate != Wait {
		if bjl.getUserDataByID(gameuser.UserID) != nil {
			return
		}

		if bjl.zhuangID == gameuser.UserID {
			return
		}
	}

	bjl.delFromZhuangList(gameuser.UserID)

	bjl.ChangeUserStateEnd(gameuser)
	bjl.ComputeGameLogicBase.UserFree(gameuser)
}

// UserRecover 用戶斷線重連
func (bjl *BaiJialeLogic) UserRecover(gameuser *gamelogic.GameUser) {
	mlog.NewDefaultLogger().Debug("UserRecover")
	bjl.ComputeGameLogicBase.UserRecover(gameuser)
	ud := bjl.getUserDataByID(gameuser.UserID)
	if ud != nil {
		ud.money += ud.user.Score - ud.moneylast
		ud.moneylast = ud.user.Score
	}

	bjl.sendScene(gameuser)
}

// UserEnter 用戶進入遊戲
func (bjl *BaiJialeLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {
	bjl.ComputeGameLogicBase.UserEnter(gameuser, tableno, seatno)
	mlog.NewDefaultLogger().Debug("UserEnter:%v", gameuser)
	bjl.sendScene(gameuser)
}

// 發送場景
func (bjl *BaiJialeLogic) sendScene(gameuser *gamelogic.GameUser) {
	// 發送設置
	bjl.sendSetting(gameuser)
	bjl.SendMessageToUser(gameuser, bjl.getRoundMsg())
	bjl.SendMessageToUser(gameuser, bjl.getResultCardListMsg())
	bjl.SendMessageToUser(gameuser, bjl.getRichestListMsg())
	bjl.SendMessageToUser(gameuser, bjl.getLuckyStarMsg())
	if bjl.gamestate == Wait {
		return
	}
	// 發送階段
	bjl.SendMessageToUser(gameuser, bjl.getZhuangInfoMsg())
	bjl.SendMessageToUser(gameuser, bjl.getBetTotalMsg(0))
	bjl.sendUserBetInfo(gameuser)

	if bjl.gamestate == KaiJiangAction {
		// 發送開獎信息
		bjl.SendMessageToUser(gameuser, bjl.getGameResultPokerMsg())
	} else if bjl.gamestate == EndAction {
		// 發送用戶下注輸贏信息如果有
		ud := bjl.getUserDataByID(gameuser.UserID)
		if ud != nil {
			bjl.SendMessageToUser(gameuser, bjl.getUserWinMsg(ud.getShowWinLose(), bjl.zhuangWin))
		}
	}

	bjl.SendMessageToUser(gameuser, bjl.getStateMsg())
}

// 發送設置
func (bjl *BaiJialeLogic) sendSetting(gameuser *gamelogic.GameUser) {
	mlog.NewDefaultLogger().Debug("sendSetting")
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BJL_GameMessageClassID_BJLSettingID)

	pdata := new(netproto.BJL_Setting)
	pdata.BetCount = append(pdata.BetCount, bjl.betcount...)
	pdata.BetMoney = proto.Int32(int32(bjl.betvalue.getMoneyByType(1)))
	pdata.BetMoneyLimit = proto.Int64(bjl.betvalue.getMaxMoneyByType(1))
	pdata.ZhuangLimit = proto.Int64(bjl.zhuanglimit)
	pdata.MaxZhuang = proto.Int32(bjl.maxturn)
	pdata.BetRange = bjl.getBetRangeData()
	pdata.Betlimitmoney = proto.Int64(bjl.betlimitmoney)
	pmsg.MsgData = pdata
	bjl.SendMessageToUser(gameuser, pmsg)
}

// 發送用戶下注信息
func (bjl *BaiJialeLogic) sendUserBetInfo(gameuser *gamelogic.GameUser) {
	ud := bjl.getUserDataByID(gameuser.UserID)
	if ud == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetInfoID)

	pdata := new(netproto.BJL_UserBetInfo)
	pdata.UserMoney = proto.Int64(ud.money)

	lucky := bjl.getUserDataByID(bjl.luckystar) // 獲取明燈
	for _, v := range ud.userbet {
		betinfo := new(netproto.BJL_BroadCastBet)
		betinfo.Type = proto.Int32(v.baijialetype)
		betinfo.Money = proto.Int64(v.money)
		betinfo.LuckyStar = proto.Bool(false)
		if lucky != nil && lucky.getUserBetByType(v.baijialetype) > 0 {
			betinfo.LuckyStar = proto.Bool(true)
		}
		pdata.BetInfo = append(pdata.BetInfo, betinfo)
	}
	pmsg.MsgData = pdata
	bjl.SendMessageToUser(gameuser, pmsg)
}

// 清理數據
func (bjl *BaiJialeLogic) clearGameData() {
	bjl.userdata = make([]*baijialeuserdata, 0)
	bjl.bettypetotal = make([]*bettype, 0)
	bjl.bettyperobot = make([]*bettype, 0)
	bjl.winuserinfo = make([]int32, 0)
	bjl.cardzhuang = make([]*card, 0)
	bjl.cardxian = make([]*card, 0)
	bjl.jiangtype = make([]int32, 0)
	bjl.stageuserinfo = make([]int32, 0)
	bjl.richestuserinfo = make([]int32, 0)
	bjl.richestList = make([]*netproto.BJL_Richest, 0)
	bjl.bettotalrobot = 0
	bjl.bettotal = 0
	bjl.zhuangWin = 0
	bjl.zuhetype = 0
	bjl.currtax = 0
	bjl.robotwinlose = 0
	bjl.userwinlose = 0
	bjl.jiangtype = nil
	bjl.assignRewardDetail = nil
	bjl.execAssignReward = false
	bjl.execFangFen = false
	rand.Seed(time.Now().UnixNano())

	zlen := len(bjl.zhuanglist)

	if zlen == 0 {
		bjl.robotZhuang()
	}

	bjl.robotXiaZhuang()
}

// 機器人下莊
func (bjl *BaiJialeLogic) robotXiaZhuang() {
	if bjl.zhuangID == 0 {
		return
	}

	ur := bjl.GetUserByUserID(bjl.zhuangID)
	if ur != nil && ur.IsRobot() {
		if rand.Int31n(100) < bjl.robotxiazhuang {
			bjl.zhuangID = 0
		}
	}
}

// 機器人上莊
func (bjl *BaiJialeLogic) robotZhuang() {
	if bjl.publiczhuang == 1 {
		// 概率當莊
		if rand.Int31n(100) > bjl.robotzhuang {
			return
		}

		// 統計機器人個數
		robotarr := make([]int32, 0)
		for _, ur := range bjl.AllTableUser() {
			if ur.IsRobot() {
				robotarr = append(robotarr, ur.UserID)
			}
		}

		if len(robotarr) <= 0 {
			return
		}

		// 隨機有幾個人當莊
		n := int32(0)
		if bjl.maxrobotzhuang > 0 {
			n = rand.Int31n(bjl.maxrobotzhuang) + 1
		}

		for i := int32(0); i < n; i++ {
			rcn := len(robotarr)
			if rcn <= 0 {
				break
			}
			index := rand.Int31n(int32(rcn))

			bjl.addToZhuangList(robotarr[index])
			robotarr = append(robotarr[:index], robotarr[index+1:]...)
		}
	} else {
		mlog.Debug("未開啟上莊機制，機器人不用上莊")
		return
	}
}

// 開始下注
func (bjl *BaiJialeLogic) startBet() {
	bjl.opttimer = bjl.TimerFuncManager.NewTimer(bjl.bettime, bjl.betEnd)
	bjl.setState(Bet)
	bjl.BoardcastMessage(bjl.getStateMsg())
}

// 結束下注
func (bjl *BaiJialeLogic) betEnd() {

	// 統計所有玩家下注
	for _, ud := range bjl.userdata {
		ud.betAmount = 0

		betMap := make(map[int32]int64)
		for _, bet := range ud.userbet {
			betMap[bet.baijialetype] += bet.money
			ud.betAmount += bet.money
		}

		isRobot := ""
		if ud.user.IsRobot() {
			isRobot = "[機器人]"
			continue
		}

		bjl.GameLog.AddGameStep("%s玩家[%d] 押注總計: %d, 莊: %d, 閒: %d, 和: %d, 莊對: %d, 閒對: %d", isRobot, ud.user.UserID, ud.betAmount, betMap[1], betMap[2], betMap[3], betMap[6], betMap[7])
	}

	bjl.actionStr = bjl.RegisterAnimation(GameEndAction, bjl.betend, bjl.kaiJiang)
}

// 開獎
func (bjl *BaiJialeLogic) kaiJiang() {
	bjl.setState(KaiJiang)

	if bjl.zhuangID != 0 {
		if bjl.GetUserByUserID(bjl.zhuangID).IsRobot() {
			bjl.removeRobotBet()
		}
	} else {
		bjl.removeRobotBet()
	}

	// 稅收
	money, tax := bjl.GetUserWinAndTax(bjl.bettotal)
	bjl.currtax = tax
	//bjl.bettotal = Money
	if bjl.zhuangID != 0 {
		ur := bjl.GetUserByUserID(bjl.zhuangID)
		if ur != nil && !ur.IsRobot() {
			money += ur.Score
		}
	}

	bjl.getResultType(money)
	bjl.countGameRound()
	bjl.getPokerByType(bjl.jiangtype)
	bjl.BoardcastMessage(bjl.getGameResultPokerMsg())
	bjl.BoardcastMessage(bjl.getRoundMsg())
	bjl.setState(KaiJiangAction)
	bjl.addResultCardList()
	bjl.actionStr = bjl.RegisterAnimation(GameKaiJiangAction, bjl.kaijiangactiontime, bjl.sendResult)
}

// 添加開獎結果到歷史開獎
func (bjl *BaiJialeLogic) addResultCardList() {
	retcards := new(netproto.BJL_GameResultPoker)
	for i := 0; i < len(bjl.cardzhuang); i++ {
		card := new(netproto.BJL_Poker)
		card.Color = proto.Int32(bjl.cardzhuang[i].color)
		card.Number = proto.Int32(bjl.cardzhuang[i].number)
		retcards.ZhuangPoker = append(retcards.ZhuangPoker, card)
	}

	for i := 0; i < len(bjl.cardxian); i++ {
		card := new(netproto.BJL_Poker)
		card.Color = proto.Int32(bjl.cardxian[i].color)
		card.Number = proto.Int32(bjl.cardxian[i].number)
		retcards.XianPoker = append(retcards.XianPoker, card)
	}

	retcards.RoundID = proto.Int32(bjl.roundid)
	retcards.ZuHe = proto.Uint32(bjl.zuhetype)
	var res string

	res += fmt.Sprintf("%d", bjl.jiangtype[0])

	if bjl.zuhetype&(1<<5) != 0 {
		res += "," + "6"
	}
	if bjl.zuhetype&(1<<6) != 0 {
		res += "," + "7"
	}
	if bjl.assignRewardDetail != nil {
		example.invalid("===========%v", bjl.jiangtype)
		bjl.assignRewardDetailResult(res)
	}
	bjl.resultcardlist.CardList = append(bjl.resultcardlist.CardList, retcards)
	if len(bjl.resultcardlist.CardList) > 100 {
		bjl.resultcardlist.CardList = append(bjl.resultcardlist.CardList[:0], bjl.resultcardlist.CardList[1:]...)
	}
}

// 獲取開獎結果的鍊錶
func (bjl *BaiJialeLogic) getResultCardListMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BJL_GameMessageClassID_BJLResultCardListID)

	pdata := new(netproto.BJL_ResultCardList)
	if len(bjl.resultcardlist.CardList) > 30 {
		pdata.CardList = append(pdata.CardList, bjl.resultcardlist.CardList[len(bjl.resultcardlist.CardList)-30:]...)
	} else {
		pdata.CardList = append(pdata.CardList, bjl.resultcardlist.CardList...)
	}

	//mlog.Debug("getResultCardListMsg:%v", pdata.CardList)
	pmsg.MsgData = pdata
	return pmsg
}

// 獲取開獎的撲克牌
func (bjl *BaiJialeLogic) getGameResultPokerMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BJL_GameMessageClassID_BJLGameResultPokerID)

	pdata := new(netproto.BJL_GameResultPoker)
	for _, v := range bjl.cardzhuang {
		cards := new(netproto.BJL_Poker)
		cards.Color = proto.Int32(v.color)
		cards.Number = proto.Int32(v.number)
		pdata.ZhuangPoker = append(pdata.ZhuangPoker, cards)
	}

	for _, v := range bjl.cardxian {
		cards := new(netproto.BJL_Poker)
		cards.Color = proto.Int32(v.color)
		cards.Number = proto.Int32(v.number)
		pdata.XianPoker = append(pdata.XianPoker, cards)
	}

	pdata.RoundID = proto.Int32(bjl.roundid)
	if bjl.zuhetype == 0 {
		jiangtypeToZuhe := uint32(bjl.jiangtype[0])
		pdata.ZuHe = proto.Uint32(jiangtypeToZuhe)
	} else {
		pdata.ZuHe = proto.Uint32(bjl.zuhetype)
	}

	pmsg.MsgData = pdata

	return pmsg
}

// 計算局數
func (bjl *BaiJialeLogic) countGameRound() {
	if bjl.roundid >= 30 {
		bjl.roundid = 1
	} else {
		bjl.roundid++
		mlog.Debug("bjl.roundid %d", bjl.roundid)
		if bjl.roundid == 30 {
			bjl.newcard = true //30場了，重新發牌
			example.invalid("重新發牌:%t", bjl.newcard)
		}
	}

	if bjl.jiangtype[0] == ZHUANG {
		bjl.zhuangjawin++
	} else if bjl.jiangtype[0] == XIAN {
		bjl.xianwin++
	} else if bjl.jiangtype[0] == HE {
		bjl.hewin++
	}
	mlog.Debug("开奖类型：%v", bjl.jiangtype)
}

// 获取开奖的类型，庄闲和,Betmoney為稅後的下注額+玩家身上的錢
func (bjl *BaiJialeLogic) getResultType(betmoney int64) {
	superToolHandler := bjl.SuperToolService.GetBrGameSuperToolHandler()
	// 指定開獎
	mlog.Debug("bjl.nextjiangtype:%d", bjl.nextjiangtype)
	if bjl.nextjiangtype != 0 {
		bjl.jiangtype = append(bjl.jiangtype, bjl.nextjiangtype)
		if bjl.nextzuhetype != 0 {
			bjl.jiangtype = append(bjl.jiangtype, bjl.nextzuhetype)
		}
		if bjl.nextzuhetype2 != 0 {
			bjl.jiangtype = append(bjl.jiangtype, bjl.nextzuhetype2)
		}
		// 組合需加上贏家
		if bjl.nextjiangtype == ZHUANG {
			bjl.zuhetype |= 1 << uint32(ZHUANG-1) // 莊
		} else if bjl.nextjiangtype == XIAN {
			bjl.zuhetype |= 1 << uint32(XIAN-1) // 閒
		} else if bjl.nextjiangtype == HE {
			bjl.zuhetype |= 1 << uint32(HE-1) // 和
		}

		bjl.zuhetype |= 1 << uint32(bjl.nextzuhetype-1) // 組合
		if bjl.nextzuhetype2 != 0 {
			bjl.zuhetype |= 1 << uint32(bjl.nextzuhetype2-1)
		}

		bjl.nextjiangtype = 0
		bjl.nextzuhetype = 0
		bjl.nextzuhetype2 = 0
		return
	} else if ok, rettype := bjl.isAssignReward(); ok {
		example.invalid("==============================指定開獎:%v", rettype)
		bjl.execAssignReward = true
		if len(rettype) > 0 {
			if rettype[0] == 1 || rettype[0] == 2 || rettype[0] == 3 {
				bjl.jiangtype = []int32{rettype[0]}
				bjl.zuhetype |= 1 << uint32(bjl.jiangtype[0]-1)
			}
			if len(rettype) > 1 {
				bjl.jiangtype = append(bjl.jiangtype, rettype[1])
				bjl.zuhetype |= 1 << uint32(rettype[1]-1)

			}
			if len(rettype) > 2 {
				bjl.jiangtype = append(bjl.jiangtype, rettype[2])
				bjl.zuhetype |= 1 << uint32(rettype[2]-1)
			}
		}
		return
	} else if bjl.fangFenTarget > 0 && bjl.fangFenComplete < bjl.fangFenTarget { // 系統放分
		bjl.execFangFen = true
		alltype := bjl.getmaxjiangtype()
		bjl.jiangtype = []int32{bjl.betvalue.getJiangType(alltype)}
		bjl.zuhetype |= 1 << uint32(bjl.jiangtype[0]-1)
		zuhetypenew := bjl.betvalue.getJiangType(bjl.getFangFenZuheType())
		if zuhetypenew != 0 {
			example.invalid("zuhetype:%d", zuhetypenew)
			bjl.zuhetype |= 1 << uint32(zuhetypenew-1)
		}
		example.invalid("系統放分 : target %d, already %d", bjl.fangFenTarget, bjl.fangFenComplete)
		bjl.GameLog.AddGameStep("系統放分 : target %d, already %d", bjl.fangFenTarget, bjl.fangFenComplete)
		return
	} else if len(bjl.getBowBet()) > 0 { // 監控放殺分
		if bjl.bowJiangType() {
			example.invalid("執行監控放殺分策略 : %s", superToolHandler.GetProiorityListString())
			bjl.GameLog.AddGameStep("執行監控放殺分策略 : %s", superToolHandler.GetProiorityListString())
			return
		}
	}

	// 先開出莊閒和中的一種。
	alltype := bjl.getAllKaiJiangType(betmoney)

	alltype = bjl.getminjiangtype(alltype)

	jiangtype := bjl.betvalue.getJiangType(alltype)

	if bjl.zhuangID != 0 && !bjl.GetUserByUserID(bjl.zhuangID).IsRobot() && len(bjl.bettyperobot) > int(0) {
		// 玩家當莊且機器人有下注的時候
		if bjl.SysPool.GetPool() < bjl.SysPool.GetWinLine() {
			eat := util.RndInt(1, 5)
			mlog.NewDefaultLogger().System("水池小於贏錢線就1/2吃玩家, rand 1-4[%d] 3以上吃", eat)
			if eat >= 3 {
				jiangtype = bjl.getMaxRobotWinType()
			}
		} else {
			eat := util.RndInt(1, 4)
			mlog.NewDefaultLogger().System("水池大於贏錢線就1/3吃玩家, rand 1-3[%d] 3以上吃", eat)
			if eat >= 3 {
				jiangtype = bjl.getMaxRobotWinType()
			}
		}
	}

	// money獲利 =Betmoney:税收后钱总数为,開完庄閒和後還能賠付的錢
	money := betmoney - bjl.getBetMoneyByType(jiangtype)*int64(bjl.betvalue.getOddByType(jiangtype))
	isHe := (jiangtype == 3) //和
	bjl.zuhetype |= 1 << uint32(jiangtype-1)
	bjl.jiangtype = append(bjl.jiangtype, jiangtype)
	// 開出莊天王
	tvalue := int32(4)
	isZhuangTianWang := false // 莊天王
	isXianTianWang := false   // 閒天王
	isZhuangDui := false      // 莊對子
	isXianDui := false        // 閒對子
	isZhuangShun := false     // 莊順子
	isXianShun := false       // 閒順子
	if bjl.getCanKaiJiangType(money, tvalue) && bjl.betvalue.getIsKai(tvalue) && bjl.zuhetype&1 != 0 {
		if isHe {
			needmoney := int64(bjl.betvalue.getOddByType(tvalue)) * bjl.getBetMoneyByType(tvalue)
			needmoney += int64(bjl.betvalue.getOddByType(5)) * bjl.getBetMoneyByType(5)
			// 開出莊天王和條件：稅收後的錢 > 莊天王賠率獎金 ＋ 閒天王賠率獎金
			if money > needmoney {
				money -= int64(bjl.betvalue.getOddByType(tvalue)) * bjl.getBetMoneyByType(tvalue)
				bjl.jiangtype = append(bjl.jiangtype, tvalue)
				bjl.zuhetype |= 1 << uint32(tvalue-1)
				bjl.zuhetype |= 1 << uint32(4)
				isZhuangTianWang = true
			}
		} else { // 莊天王
			money -= int64(bjl.betvalue.getOddByType(tvalue)) * bjl.getBetMoneyByType(tvalue)
			bjl.jiangtype = append(bjl.jiangtype, tvalue)
			bjl.zuhetype |= 1 << uint32(tvalue-1)
			isZhuangTianWang = true
		}
	}

	tvalue = 5 // 閒天王
	if bjl.getCanKaiJiangType(money, tvalue) &&
		bjl.betvalue.getIsKai(tvalue) && bjl.zuhetype&1 == 0 && !isHe {

		money -= int64(bjl.betvalue.getOddByType(tvalue)) * bjl.getBetMoneyByType(tvalue)
		bjl.jiangtype = append(bjl.jiangtype, tvalue)
		isXianTianWang = true
		bjl.zuhetype |= 1 << uint32(tvalue-1)
	}

	tvalue = 6 // 莊對子
	if bjl.getCanKaiJiangType(money, tvalue) && bjl.betvalue.getIsKai(tvalue) {
		money -= int64(bjl.betvalue.getOddByType(tvalue)) * bjl.getBetMoneyByType(tvalue)
		bjl.jiangtype = append(bjl.jiangtype, tvalue)
		isZhuangDui = true
		bjl.zuhetype |= 1 << uint32(tvalue-1)
	}

	tvalue = 7 // 閒對子
	if bjl.getCanKaiJiangType(money, tvalue) && bjl.betvalue.getIsKai(tvalue) {
		money -= int64(bjl.betvalue.getOddByType(tvalue)) * bjl.getBetMoneyByType(tvalue)
		bjl.jiangtype = append(bjl.jiangtype, tvalue)
		isXianDui = true
		bjl.zuhetype |= 1 << uint32(tvalue-1)
	}

	tvalue = 8 // 莊順子
	if bjl.getCanKaiJiangType(money, tvalue) && bjl.betvalue.getIsKai(tvalue) &&
		(!isXianTianWang || !isZhuangTianWang) && !isZhuangDui {
		money -= int64(bjl.betvalue.getOddByType(tvalue)) * bjl.getBetMoneyByType(tvalue)
		bjl.jiangtype = append(bjl.jiangtype, tvalue)
		isZhuangShun = true
		bjl.zuhetype |= 1 << uint32(tvalue-1)
	}

	tvalue = 9 // 閒順子
	if bjl.getCanKaiJiangType(money, tvalue) && bjl.betvalue.getIsKai(tvalue) &&
		(!isXianTianWang || !isZhuangTianWang) && !isXianDui {
		money -= int64(bjl.betvalue.getOddByType(tvalue)) * bjl.getBetMoneyByType(tvalue)
		bjl.jiangtype = append(bjl.jiangtype, tvalue)
		isXianShun = true
		bjl.zuhetype |= 1 << uint32(tvalue-1)
	}

	tvalue = 10 // 同點和

	if isHe && bjl.getCanKaiJiangType(money, tvalue) && bjl.betvalue.getIsKai(tvalue) &&
		(isZhuangTianWang == isXianTianWang || isXianDui == isZhuangDui) &&
		isZhuangShun == isXianShun {
		money -= int64(bjl.betvalue.getOddByType(tvalue)) * bjl.getBetMoneyByType(tvalue)
		bjl.jiangtype = append(bjl.jiangtype, tvalue)
		bjl.zuhetype |= 1 << uint32(tvalue-1)
	}

}

// 發送開獎結果
func (bjl *BaiJialeLogic) sendResult() {
	calWinLose(bjl)
	//bjl.writePoolToFile()
	bjl.updateVipUserWin()
	bjl.actionStr = bjl.RegisterAnimation(GameEndAction, bjl.endactiontime, bjl.gameEnd)
}

// 計算下注輸贏結果
func calWinLose(bjl *BaiJialeLogic) {
	sysTaxIncomeBefore := bjl.systaxincome
	bjl.syslastwin = 0
	bjl.setState(EndAction)
	for _, ud := range bjl.userdata {
		bjl.sendUserResult(ud)
	}

	// 準備計算水池異動值
	changeSysPoolValue := bjl.zhushuimoney
	if bjl.zhuangID == 0 {
		// 系統小妹當莊的時候
		if bjl.userwinlose > 0 {
			// 玩家贏錢，水線下降
			// bjl.syspool -= bjl.userwinlose - bjl.zhushuimoney
			changeSysPoolValue = bjl.zhushuimoney - bjl.userwinlose
		} else {
			// 玩家輸錢，水線上升(稅)
			// bjl.syspool -= int64(math.Floor(float64(bjl.userwinlose)*bjl.GetTax())) - bjl.zhushuimoney
			changeSysPoolValue = bjl.zhushuimoney - int64(math.Floor(float64(bjl.userwinlose)*bjl.GetTax()))
		}
	} else {
		income, tax := int64(0), int64(0)
		if bjl.bettotal > int64(0) { // 賺錢要抽稅
			income, tax = bjl.GetUserWinAndTax(bjl.bettotal)
			bjl.systaxincome += tax
			bjl.GameLog.AddGameStep("莊家贏錢 %d, 抽稅 %d", income, tax)
		} else {
			income = bjl.bettotal
		}

		if len(bjl.userdata) > int(0) {
			bjl.systotalwin += income
			bjl.syslastwin += income
			bjl.zhuangWin += income
		}
		ur := bjl.GetUserByUserID(bjl.zhuangID)
		var finalBetTotal int64
		if ur != nil {
			if ur.IsRobot() {
				finalBetTotal = bjl.bettotal + bjl.bettotalrobot
			} else {
				finalBetTotal = bjl.bettotal
			}
			income, tax = bjl.GetUserWinAndTax(finalBetTotal)
			// 機器人當莊的時候
			if ur.IsRobot() {
				bjl.zhuangWin = income
				bjl.systaxincome = tax
				bjl.systotalwin = income
				// 玩家贏錢水池減少，並將注出的水倒回水池
				if bjl.userwinlose > 0 {
					// 玩家贏錢，水線下降
					// bjl.syspool -= bjl.userwinlose - bjl.zhushuimoney
					changeSysPoolValue = bjl.zhushuimoney - bjl.userwinlose
				} else {
					// 玩家輸錢，水線上升(稅)
					// bjl.syspool -= int64(math.Floor(float64(bjl.userwinlose)*bjl.GetTax())) - bjl.zhushuimoney
					changeSysPoolValue = bjl.zhushuimoney - int64(math.Floor(float64(bjl.userwinlose)*bjl.GetTax()))

				}
			} else {
				// 玩家當莊的時候，機器人贏錢水池就添加，並將注出的水倒回水池
				if bjl.robotwinlose > 0 {
					// 機器人贏錢 水位上升(稅)
					// bjl.syspool += int64(math.Floor(float64(bjl.robotwinlose)*bjl.GetTax())) - bjl.zhushuimoney
					changeSysPoolValue = int64(math.Floor(float64(bjl.robotwinlose)*bjl.GetTax())) - bjl.zhushuimoney
				} else {
					// bjl.syspool += bjl.robotwinlose - bjl.zhushuimoney
					changeSysPoolValue = bjl.robotwinlose - bjl.zhushuimoney
				}
			}
			fmt.Println("====================")
			fmt.Println("====================")
			fmt.Println("====================")
			fmt.Println("====================")
			fmt.Println(ur.UserID, income, tax)
			bjl.UpdateUserScoreFunc(ur.UserID, income, tax) // 更新金幣異動
			bjl.GameLog.Context().UserUpdateScore(ur.UserID, income)
			bjl.GameLog.Context().UserTaxAmount(ur.UserID, tax)

			bjl.SendMessageToUser(ur, bjl.getUserWinMsg(bjl.zhuangWin, bjl.zhuangWin))
		} else {
			mlog.NewDefaultLogger().Debug("庄家数据为空")
		}
	}

	if !bjl.execFangFen {
		example.invalid("系统水池異動, from %d, change %d", bjl.SysPool.GetPool(), changeSysPoolValue)
		bjl.GameLog.AddGameStep("系统水池異動, from %d, change %d", bjl.SysPool.GetPool(), changeSysPoolValue)
		bjl.SysPool.UpdateSysPool(changeSysPoolValue)
	} else {
		before := bjl.fangFenComplete
		changeSysPoolValue = -changeSysPoolValue // 水池下降代表玩家贏錢
		if changeSysPoolValue > 0 {              // 玩家贏錢要抽稅
			changeSysPoolValue = int64(math.Floor(float64(changeSysPoolValue) * bjl.GetTax()))
		} else { // 玩家輸錢要把稅算回來
			changeSysPoolValue = int64(math.Floor(float64(changeSysPoolValue) / bjl.GetTax()))
		}
		bjl.fangFenComplete += changeSysPoolValue
		bjl.GameLog.AddGameStep("系統放分完成, target %d, already %d-%d", bjl.fangFenTarget, before, bjl.fangFenComplete)
		example.invalid("系統放分完成, target %d, already %d-%d", bjl.fangFenTarget, before, bjl.fangFenComplete)
	}

	mlog.Debug("庄家赢钱：%v", bjl.zhuangWin)
	bjl.GameLog.AddGameStep("本局稅收 %d", bjl.systaxincome-sysTaxIncomeBefore)
	// 發送沒下注的user 莊家結果
	for _, u2 := range bjl.AllTableUser() {
		sendBool := true
		for _, ud := range bjl.userdata {
			if u2.UserID == ud.user.UserID {
				sendBool = false
				break
			}
		}
		if sendBool {
			bjl.SendMessageToUser(u2, bjl.getUserWinMsg(0, bjl.zhuangWin))
			if u2.UserID == bjl.luckystar {
				if bjl.getUserDataByID(bjl.luckystar) == nil {
					bjl.notbettimes++
				}
			}
		}
	}
	// 發送下注user結算結果
	for _, ud := range bjl.userdata {
		bjl.SendMessageToUser(ud.user, bjl.getUserWinMsg(ud.getShowWinLose(), bjl.zhuangWin))

		bjl.GameLog.Context().UserUpdateScore(ud.user.UserID, ud.updatemoney)
		bjl.GameLog.Context().UserBetScore(ud.user.UserID, ud.betAmount)
		bjl.GameLog.Context().UserPayoutScore(ud.user.UserID, ud.winAmount)
		bjl.GameLog.Context().UserValidScore(ud.user.UserID, ud.flowAmount)
		bjl.GameLog.Context().UserTaxAmount(ud.user.UserID, ud.taxAmount)
	}
	bjl.GameLog.GameEnd()
	// 更新 pergamelog.detaillog
	bjl.bettotal = 0
}

// 給每個玩家發送他們的開獎輸贏
func (bjl *BaiJialeLogic) sendUserResult(ud *baijialeuserdata) {
	// var betlosemoney int64                                                         // 下注輸給莊的錢
	var updateAmount, flowAmount, betAmount, rtnAmount, winAmount, taxAmount int64 // 流水, 有效投注金額, 返還金額計算

	betMap := make(map[int32]int64)
	for _, bet := range ud.userbet {
		betMap[bet.baijialetype] += bet.money
	}

	isHe, _ := util.InArray(int32(HE), bjl.jiangtype)

	// 依押注類型計算
	for bettype, money := range betMap {
		betAmount += money
		var isWin, isReturn bool
		if ok, _ := util.InArray(bettype, bjl.jiangtype); ok { // 此注類型有開
			isWin = true
			isReturn = true
			win := int64(bjl.betvalue.getOddByType(bettype))*money - money // 設定中的賠率包含押注金額

			if bjl.BJLMode == BJLModeLucky6 { // 幸運6百家樂
				if bettype == ZHUANG {
					i32Point := bjl.getCardsDian(bjl.cardzhuang)
					if i32Point == 6 {
						win = int64(OddZhuang6)*money/10 - money
					}
				}
			}

			income, tax := bjl.GetUserWinAndTax(win)

			winAmount += win
			rtnAmount += money
			rtnAmount += income
			taxAmount += tax
		} else if isHe && (bettype == ZHUANG || bettype == XIAN) { // 開和退 [莊] [閒]
			isReturn = true
			rtnAmount += money
		}

		// 沒贏沒退 = 輸 算流水
		if !isWin && !isReturn {
			flowAmount += money
		}
	}

	updateAmount = rtnAmount - betAmount

	mlog.Debug("User[%d] UpdateAmount: %d, BetAmount: %d, WinAmount: %d, ReturnAmount: %d, FlowAmount: %d, TaxAmount: %d", ud.user.UserID, updateAmount, betAmount, winAmount, rtnAmount, flowAmount, taxAmount)

	// 計算玩家贏錢, 並扣掉派彩池, 計算剩下的派彩池
	// for i := 0; i < len(ud.userbet); i++ {
	// 	betAmount += util.ABSInt64(ud.userbet[i].money)
	// 	betlosemoney += ud.userbet[i].money
	// 	isHe, _ := util.InArray(int32(HE), bjl.jiangtype)
	// 	isWin := false
	// 	if exists, _ := util.InArray(ud.userbet[i].baijialetype, bjl.jiangtype); exists {
	// 		isWin = true
	// 		win := int64(bjl.betvalue.getOddByType(ud.userbet[i].baijialetype)) * ud.userbet[i].money

	// 		if bjl.BJLMode == BJLModeLucky6 { // 幸運6百家樂
	// 			// 莊家勝，檢查是否為 6 點
	// 			if ud.userbet[i].baijialetype == ZHUANG {
	// 				i32Point := bjl.getCardsDian(bjl.cardzhuang)
	// 				if i32Point == 6 {
	// 					win = int64(OddZhuang6) * ud.userbet[i].money / 10
	// 				}
	// 			}
	// 		}

	// 		//flowAmountTotal := win - ud.userbet[i].money
	// 		//flowamountAfterTax, _ := bjl.GetUserWinAndTax(flowAmountTotal)
	// 		//flowAmount += flowamountAfterTax
	// 		if bjl.zhuangID == 0 {
	// 			if !ud.user.IsRobot() {
	// 				bjl.bettotal -= win
	// 			}
	// 		} else {
	// 			if bjl.GetUserByUserID(bjl.zhuangID).IsRobot() {
	// 				if !ud.user.IsRobot() {
	// 					bjl.bettotal -= win
	// 				} else {
	// 					bjl.bettotalrobot -= win
	// 				}
	// 			} else {
	// 				bjl.bettotal -= win
	// 			}
	// 		}
	// 		ud.money += win
	// 	} else if isHe && !isWin && ud.userbet[i].baijialetype != ZHUANGDUI && ud.userbet[i].baijialetype != XIANDUI {
	// 		//開和全退
	// 		win := ud.userbet[i].money
	// 		ud.win = win
	// 		//机器人不计算奖池里面的钱
	// 		if bjl.zhuangID == 0 {
	// 			if !ud.user.IsRobot() {
	// 				bjl.bettotal -= win
	// 			}
	// 		} else {
	// 			ur := bjl.GetUserByUserID(bjl.zhuangID)
	// 			if ur != nil && ur.IsRobot() && ud.user.IsRobot() {
	// 				bjl.bettotalrobot -= win
	// 			} else {
	// 				bjl.bettotal -= win
	// 			}
	// 		}
	// 		ud.money += win
	// 	} else {
	// 		//flowAmount += ud.userbet[i].money
	// 	}
	// }
	// Money 賠付獎金 = 結算後的錢 - 當前餘額
	// money := ud.money - ud.user.Score

	// 將機器人及玩家輸贏減掉稅
	if ud.user.IsRobot() {
		bjl.robotwinlose += updateAmount
	} else {
		bjl.userwinlose += updateAmount
	}

	bjl.systaxincome += taxAmount

	// var income, tax int64
	// if money > int64(0) {
	// 	income, tax = bjl.GetUserWinAndTax(money)
	// 	ud.win = income + betlosemoney
	// 	bjl.systaxincome += tax
	// } else {
	// 	income = money
	// 	ud.win = income + betlosemoney
	// }
	// if income < 0 {
	// 	flowAmount = betAmount
	// } else {
	// 	flowAmount = 0
	// }
	// payAmount = betAmount + income

	ud.betAmount = betAmount
	ud.winAmount = rtnAmount + taxAmount
	ud.flowAmount = flowAmount
	ud.taxAmount = taxAmount
	ud.win = rtnAmount

	ud.updatemoney = updateAmount
	stepmsg := "玩家[%d] 下注贏錢: %d, 抽稅: %d, 流水: %d, 有效投注金額: %d, 返還金額: %d"
	if ud.user.IsRobot() {
		stepmsg = "[機器人]" + stepmsg
	}

	if !ud.user.IsRobot() {
		bjl.GameLog.AddGameStep(stepmsg, ud.user.UserID, updateAmount, taxAmount, flowAmount, betAmount, rtnAmount)
	}

	ud.user.FlowAmount += flowAmount
	ud.user.BetAmount += betAmount
	ud.user.PayAmount += rtnAmount
	bjl.UpdateUserScoreFunc(ud.user.UserID, updateAmount, taxAmount) // 更新金幣異動

	// bjl.GameLog.Context().UserUpdateScore(ud.user.UserID, updateAmount)
	// bjl.GameLog.Context().UserBetScore(ud.user.UserID, betAmount)
	// bjl.GameLog.Context().UserPayoutScore(ud.user.UserID, rtnAmount+taxAmount)
	// bjl.GameLog.Context().UserValidScore(ud.user.UserID, flowAmount)
	// bjl.GameLog.Context().UserTaxAmount(ud.user.UserID, taxAmount)

	winSort(ud, bjl)
}

// 這裡寫入贏家排名列表
func winSort(ud *baijialeuserdata, bjl *BaiJialeLogic) {
	badd := false
	if ud.win > 0 {
		if len(bjl.winuserinfo) > 0 {
			for i := 0; i < len(bjl.winuserinfo); i++ {
				u := bjl.getUserDataByID(bjl.winuserinfo[i])
				if u != nil {
					if ud.win > u.win {
						if len(bjl.winuserinfo) >= 3 {
							if i == 0 {
								bjl.winuserinfo[2] = bjl.winuserinfo[1]
								bjl.winuserinfo[1] = bjl.winuserinfo[0]
								bjl.winuserinfo[0] = ud.user.UserID
							} else if i == 1 {
								bjl.winuserinfo[2] = bjl.winuserinfo[1]
								bjl.winuserinfo[1] = ud.user.UserID
							} else {
								bjl.winuserinfo[2] = ud.user.UserID
							}
						} else if len(bjl.winuserinfo) == 2 {
							bjl.winuserinfo = append(bjl.winuserinfo, ud.user.UserID)
							if i == 0 {
								bjl.winuserinfo[2] = bjl.winuserinfo[1]
								bjl.winuserinfo[1] = bjl.winuserinfo[0]
								bjl.winuserinfo[0] = ud.user.UserID
							} else if i == 1 {
								bjl.winuserinfo[2] = bjl.winuserinfo[1]
								bjl.winuserinfo[1] = ud.user.UserID
							}
						} else {
							bjl.winuserinfo = append(bjl.winuserinfo, ud.user.UserID)
							bjl.winuserinfo[1] = bjl.winuserinfo[0]
							bjl.winuserinfo[0] = ud.user.UserID
						}

						badd = true
						break
					}
				}
			}
		}

		if len(bjl.winuserinfo) < 3 && !badd {
			bjl.winuserinfo = append(bjl.winuserinfo, ud.user.UserID)
		}
	}
}

// 獲取贏家的消息
func (bjl *BaiJialeLogic) getUserWinMsg(selfwin int64, zhuangwin int64) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BJL_GameMessageClassID_BJLGameResultWinID)

	pdata := new(netproto.BJL_GameResultUserMoney)
	pdata.ZhuangWin = proto.Int64(zhuangwin)
	pdata.SelfWin = proto.Int64(selfwin)
	for i := 0; i < len(bjl.winuserinfo); i++ {
		ud := bjl.getUserDataByID(bjl.winuserinfo[i])
		if ud == nil {
			continue
		}
		userwin := new(netproto.BJL_WinUerInfo)
		userid := strconv.Itoa(int(ud.user.UserID))
		userwin.UserName = proto.String(userid)
		userwin.Win = proto.Int64(ud.getShowWinLose())
		pdata.WinUserInfo = append(pdata.WinUserInfo, userwin)
	}
	for j := 0; j < len(bjl.stageuserinfo); j++ {
		ud := bjl.getUserDataByID(bjl.stageuserinfo[j])
		if ud == nil {
			continue
		}
		stageuserwin := new(netproto.BJL_StageUerInfo)
		stageuserwin.UserID = proto.Int32(ud.user.UserID)
		stageuserwin.Win = proto.Int64(ud.getShowWinLose())
		pdata.StageUerInfo = append(pdata.StageUerInfo, stageuserwin)
	}

	pmsg.MsgData = pdata
	return pmsg
}

// 遊戲結束
func (bjl *BaiJialeLogic) gameEnd() {
	bjl.setState(Wait)
	if bjl.zhuangID != 0 {
		user := bjl.GetUserByUserID(bjl.zhuangID)
		if user.IsOffLine() {
			bjl.zhuangID = 0
		}
	}
	bjl.GameEndFunc()
}

// 獲取開獎的玩家撲克牌數據
func (bjl *BaiJialeLogic) getPokerByType(t []int32) {
	if bjl.zuhetype&(1<<3) != 0 || bjl.zuhetype&(1<<4) != 0 {
		// 莊天王和閒天王。只發兩張牌
		if bjl.zuhetype&1 != 0 {
			// 莊贏，莊天王
			bjl.cardzhuang = append(bjl.cardzhuang, bjl.getTianWang(bjl.zuhetype&(1<<5) != 0, true, true)...)
			bjl.cardxian = append(bjl.cardxian, bjl.getTwoCards(bjl.cardzhuang, 1, bjl.zuhetype&(1<<6) != 0, bjl.zuhetype&(1<<4) != 0)...)
		} else {
			// 閒贏，閒天王
			v := int32(1)
			if bjl.zuhetype&(1<<2) != 0 {
				v = 0
			}
			bjl.cardxian = append(bjl.cardxian, bjl.getTianWang(bjl.zuhetype&(1<<6) != 0, true, false)...)
			bjl.cardzhuang = append(bjl.cardzhuang, bjl.getTwoCards(bjl.cardxian, v, bjl.zuhetype&(1<<5) != 0, bjl.zuhetype&(1<<3) != 0)...)
		}
	} else if bjl.zuhetype&(1<<5) != 0 || bjl.zuhetype&(1<<6) != 0 { // 莊對或閒對特殊取牌
		var xiancards []*card
		var zhuangcards []*card
		limit := 1000
		var times = 0
		for {
			bjl.fapai()
			times++
			if times > limit || bjl.cards.Len() < 6 {
				bjl.newcard = true
				times = 0
				continue
			}
			var ok bool
			xiancards, ok = bjl.getXDuiCards(bjl.zuhetype&(1<<1) != 0, bjl.zuhetype&(1<<6) != 0)
			if !ok {
				continue
			}
			zhuangcards, ok = bjl.getZDuiCards(xiancards, bjl.zuhetype&1 != 0, bjl.zuhetype&(1<<5) != 0)
			if ok {
				break
			}
			time.Sleep(time.Duration(15) * time.Millisecond)
		}
		bjl.cardxian = xiancards
		bjl.cardzhuang = zhuangcards
		// mlog.Debug("%b", bjl.zuhetype)
		// mlog.Debug("bjl.cardxian", bjl.cardxian)
		// mlog.Debug("bjl.cardzhuang", bjl.cardzhuang)
		// 1. 取莊對子
		// 2. 取閒對子
		// 3. 根據輸贏補牌
	} else {
		// 其他情況發三張牌
		if bjl.zuhetype&1 != 0 {
			// 莊贏
			for {
				bjl.fapai()
				cards1, ok1 := bjl.getThreefirstCards(bjl.zuhetype&(1<<5) != 0, bjl.zuhetype&(1<<7) != 0, true, true)
				if !ok1 {
					continue
				}
				cards2, ok2 := bjl.getThreeSencondCards(cards1, 1, bjl.zuhetype&(1<<6) != 0, bjl.zuhetype&(1<<8) != 0, false)
				if ok2 {
					mlog.NewDefaultLogger().SetData(fmt.Sprintf("cards1长度为:%d", len(cards1))).Debug("获取开奖的玩家扑克牌数据->庄赢")
					zhuang := poker.NewEmptyCardGroup(0)
					for i := 0; i < len(cards1); i++ {
						card := poker.NewPokerCard(cards1[i].color, cards1[i].number)
						if card.CardNumber == 1 {
							card.CardNumber = 14
						}
						zhuang.AddCard(card)
					}

					bjl.cardzhuang = append(bjl.cardzhuang, cards1...)
					bjl.cardxian = append(bjl.cardxian, cards2...)
					break
				}
			}
		} else {
			v := int32(1)
			if bjl.zuhetype&(1<<2) != 0 {
				v = 0
			}
			for {
				bjl.fapai()
				cards1, ok1 := bjl.getThreefirstCards(bjl.zuhetype&(1<<6) != 0, bjl.zuhetype&(1<<8) != 0, bjl.zuhetype&(1<<2) != 0, false)
				if !ok1 {
					continue
				}
				cards2, ok2 := bjl.getThreeSencondCards(cards1, v, bjl.zuhetype&(1<<5) != 0, bjl.zuhetype&(1<<7) != 0, true)
				if ok2 {

					bjl.cardzhuang = append(bjl.cardzhuang, cards2...)
					bjl.cardxian = append(bjl.cardxian, cards1...)
					break
				}
			}
		}
	}
	cArr := []int32{0, 1, 2, 3}
	for _, c := range append(bjl.cardxian, bjl.cardzhuang...) {
		if ok, _ := util.InArray(c.color, cArr); !ok {
			mlog.Error("沒有花色的牌 %v", bjl.jiangtype)

		}
	}
	zhuang := poker.NewEmptyCardGroup(0)

	for i := 0; i < len(bjl.cardzhuang); i++ {
		card := poker.NewPokerCard(bjl.cardzhuang[i].color, bjl.cardzhuang[i].number)
		if card.CardNumber == 1 {
			card.CardNumber = 14
		}
		zhuang.AddCard(card)
	}
	xian := poker.NewEmptyCardGroup(0)

	for i := 0; i < len(bjl.cardxian); i++ {
		card := poker.NewPokerCard(bjl.cardxian[i].color, bjl.cardxian[i].number)
		if card.CardNumber == 1 {
			card.CardNumber = 14
		}
		xian.AddCard(card)
	}

	bjl.GameLog.AddGameStep("開獎結果是%d,莊家的撲克牌是%v ,閒家的撲克牌是%v", t, zhuang, xian)

}

// 刪除一張牌
func (bjl *BaiJialeLogic) delCard(c *card) {
	delcard := poker.NewEmptyPokerCard()
	delcard.CardNumber = c.number
	delcard.Color = c.color
	if delcard.CardNumber == 1 {
		delcard.CardNumber = 14
	}
	mlog.Debug("roundID:%d", bjl.roundid)
	bjl.cards.DeleteCardForPokers(delcard)
}

// 還回去一張牌
func (bjl *BaiJialeLogic) restoreCard(c *card) {
	addcard := poker.NewEmptyPokerCard()
	addcard.CardNumber = c.number
	addcard.Color = c.color
	if addcard.CardNumber == 1 {
		addcard.CardNumber = 14
	}
	bjl.cards.AddCard(addcard)
}

//通过索引获取一张牌
func (bjl *BaiJialeLogic) getCardbyIndex(index int32) *card {
	card := new(card)
	if index < int32(len(bjl.cards.Cards)) {
		card.number = bjl.cards.Cards[index].CardNumber
		card.color = bjl.cards.Cards[index].Color
	} else {
		return nil
		//mlog.Fatal("已经没牌了")
	}

	if card.number == 14 {
		card.number = 1
	}

	return card
}

//通过扑克点数获取一张扑克点数相同的牌
func (bjl *BaiJialeLogic) getCardByNumber(no int32) *card {
	card := new(card)
	tmpno := no
	if tmpno == 1 {
		tmpno = 14
	}
	for i := 0; i < len(bjl.cards.Cards); i++ {
		if bjl.cards.Cards[i].CardNumber == tmpno {
			card.number = bjl.cards.Cards[i].CardNumber
			card.color = bjl.cards.Cards[i].Color
			break
		}
	}

	if card.number == 14 {
		card.number = 1
	}

	return card
}

//获取一个天王牌型
func (bjl *BaiJialeLogic) getTianWang(isDui bool, isWin bool, isZhuang bool) []*card {
	cards := make([]*card, 0)
	index := int32(0)
	card1 := bjl.getCardbyIndex(index)
	card2 := new(card)
	if isDui {
		//是对子的时候只能是4和9
		v := rand.Int31n(2)
		if v == 0 {
			card1 = bjl.getCardByNumber(4)
		} else {
			card1 = bjl.getCardByNumber(9)
		}

		bjl.delCard(card1)
		card2 = bjl.getCardByNumber(card1.number)
	} else {
		bjl.delCard(card1)
		card2 = bjl.getCardbyIndex(index)
		n1 := card1.number
		if n1 >= 10 {
			n1 = 0
		}
		number := int32(8)
		v := 3
		if isZhuang {
			v = 4
		}
		if isWin && bjl.zuhetype&(1<<uint(v)) != 0 {
			number = 9
		}

		for {
			n2 := card2.number
			if n2 >= 10 {
				n2 = 0
			}

			if card1.number == card2.number {
				index++
				card2 = bjl.getCardbyIndex(index)
				continue
			}

			if ((n1 + n2) % 10) >= number {
				break
			}

			index++
			card2 = bjl.getCardbyIndex(index)
		}
	}

	bjl.delCard(card2)
	cards = append(cards, card1)
	cards = append(cards, card2)
	return cards
}

//获取赢的三张牌型
func (bjl *BaiJialeLogic) getThreefirstCards(isDui bool, isShun bool, isWin bool, isZhuang bool) ([]*card, bool) {
	cards := make([]*card, 0)
	index := int32(0)
	card1 := bjl.getCardbyIndex(index)
	card2 := bjl.getCardbyIndex(index + 1)
	card3 := new(card)
	hasThird := true
	if isDui {

		v := rand.Int31n(3)
		//这里决定牌是和第一张相同还是第二张
		if (card1.number == 4 || card1.number == 9) && v == 0 {
			v++
		}

		if bjl.zuhetype&(1<<9) != 0 {
			//同点和

			tmpcards := make([]*card, 0)
			tmpcards = append(tmpcards, card1)
			tmpcards = append(tmpcards, card2)
			dian := bjl.getCardsDian(tmpcards)
			if dian > 3 {
				//这个地方确定庄家有没得第三张牌
				if dian == 3 {
					if v == 1 && card1.number == 8 {
						n := rand.Int31n(12) + 2
						if n == 8 {
							n = 7
						}
						card1 = bjl.getCardByNumber(n)
					}
				} else if dian == 4 {
					if v == 1 &&
						(card1.number >= 10 || card1.number == 1 ||
							card1.number == 8 || card1.number == 9) {
						tmpn := [6]int32{2, 3, 4, 5, 6, 7}
						n := rand.Int31n(6)
						card1 = bjl.getCardByNumber(tmpn[n])
					}
				} else if dian == 5 {
					if v == 1 && card1.number >= 10 || card1.number == 1 ||
						card1.number == 2 || card1.number == 3 ||
						card1.number == 8 || card1.number == 9 {
						tmpn := [4]int32{4, 5, 6, 7}
						n := rand.Int31n(4)
						card1 = bjl.getCardByNumber(tmpn[n])
					}
				}
			}
		}

		bjl.delCard(card1)

		if v == 0 {
			card2 = bjl.getCardByNumber(card1.number)
			cardarr1 := make([]*card, 0)
			cardarr1 = append(cardarr1, card1)
			cardarr1 = append(cardarr1, card2)
			if bjl.getCardsDian(cardarr1) >= 6 {
				hasThird = false
			}
		} else {
			hasThird = true
			for {
				if card2 == nil {
					return nil, false
				}
				cardarr1 := make([]*card, 0)
				cardarr1 = append(cardarr1, card1)
				cardarr1 = append(cardarr1, card2)
				if card1.number == card2.number {
					index++
					card2 = bjl.getCardbyIndex(index)
					continue
				}

				dian := bjl.getCardsDian(cardarr1)
				if dian >= 8 {
					index++
					card2 = bjl.getCardbyIndex(index)
					continue
				}

				if bjl.zuhetype&(1<<9) != 0 {
					//同点和

					if dian > 3 {
						//这个地方确定庄家有没得第三张牌
						if dian == 3 {
							if v == 2 && card2.number == 8 {
								n := rand.Int31n(12) + 2
								if n == 8 {
									n = 7
								}
								card2 = bjl.getCardByNumber(n)
								continue
							}
						} else if dian == 4 {
							if v == 2 &&
								(card2.number >= 10 || card2.number == 1 ||
									card2.number == 8 || card2.number == 9) {
								tmpn := [6]int32{2, 3, 4, 5, 6, 7}
								n := rand.Int31n(6)
								card2 = bjl.getCardByNumber(tmpn[n])
								continue
							}
						} else if dian == 5 {
							if v == 2 && card2.number >= 10 || card2.number == 1 ||
								card2.number == 2 || card2.number == 3 ||
								card2.number == 8 || card2.number == 9 {
								tmpn := [4]int32{4, 5, 6, 7}
								n := rand.Int31n(4)
								card2 = bjl.getCardByNumber(tmpn[n])
								continue
							}
						}
					}
				}

				if isWin {
					if v == 1 {
						cardarr1 = append(cardarr1, card2)
					} else {
						cardarr1 = append(cardarr1, card1)
					}

					if bjl.getCardsDian(cardarr1) >= 1 {
						break
					}
				}

				index++
				card2 = bjl.getCardbyIndex(index)
			}
		}

		bjl.delCard(card2)
		if hasThird {

			//第二张和第三张是对子
			if v == 1 {
				card3 = bjl.getCardByNumber(card2.number)
			} else {
				//第一张和第三张是对子
				card3 = bjl.getCardByNumber(card1.number)
			}

			if card1.number == card2.number && card2.number == card3.number {
				if bjl.zuhetype&(1<<9) != 0 {
					bjl.zuhetype = bjl.zuhetype & ^(uint32(1 << 9))
				}
			}
		}

	} else if isShun {

		for {

			if card1 == nil {
				return nil, false
			}
			if isWin && card1.number >= 10 {
				index++
				card1 = bjl.getCardbyIndex(index)
				continue
			}

			//加起来大于6了不用发后面的牌了。所以肯定不对
			card2.number = card1.number + 1
			if (card2.number+card1.number)%10 >= 6 {
				index++
				card1 = bjl.getCardbyIndex(index)
				continue
			}

			if bjl.zuhetype&(1<<9) != 0 {
				if card1.number >= 6 {
					index++
					card1 = bjl.getCardbyIndex(index)
					continue
				}
			}

			break
		}
		card2 = bjl.getCardByNumber(card1.number + 1)
		bjl.delCard(card2)
		card3 = bjl.getCardByNumber(card2.number + 1)
	} else {
		bjl.delCard(card1)
		card2 = bjl.getCardbyIndex(index)
		for {

			if card2 == nil {
				return nil, false
			}
			//不能是对子
			if card1.number == card2.number {
				index++
				card2 = bjl.getCardbyIndex(index)
				continue
			}
			cardarr1 := make([]*card, 0)
			cardarr1 = append(cardarr1, card1)
			cardarr1 = append(cardarr1, card2)
			dian := bjl.getCardsDian(cardarr1)

			if !isZhuang && isWin && bjl.zuhetype&(1<<8) != 0 {
				if dian >= 6 {
					index++
					card2 = bjl.getCardbyIndex(index)
					continue
				}
			}

			if dian >= 6 && dian < 8 {
				hasThird = false
				break
			}

			if dian < 8 {
				break
			}

			index++
			card2 = bjl.getCardbyIndex(index)
		}

		bjl.delCard(card2)
		if hasThird {
			index = 0
			card3 = bjl.getCardbyIndex(index)

			for {
				if card3 == nil {
					return nil, false
				}
				//不能是对子
				if card1.number == card3.number || card2.number == card3.number {
					index++
					card3 = bjl.getCardbyIndex(index)
					continue
				}

				cardarr1 := make([]*card, 0)
				cardarr1 = append(cardarr1, card1)
				cardarr1 = append(cardarr1, card2)
				cardarr1 = append(cardarr1, card3)
				bjl.sortPoker(cardarr1)

				//不能是顺子
				if cardarr1[0].number+1 == cardarr1[1].number && cardarr1[1].number+1 == cardarr1[2].number {
					index++
					card3 = bjl.getCardbyIndex(index)
					continue
				}

				if bjl.zuhetype&(1<<9) != 0 {
					//同点和
					tmpcards := make([]*card, 0)
					tmpcards = append(tmpcards, card1)
					tmpcards = append(tmpcards, card2)
					dian := bjl.getCardsDian(tmpcards)
					if dian > 3 {
						if dian == 3 {
							if card3.number == 8 {
								index++
								card3 = bjl.getCardbyIndex(index)
								continue
							}
						} else if dian == 4 {
							if card1.number >= 10 || card1.number == 1 ||
								card1.number == 8 || card1.number == 9 {
								index++
								card3 = bjl.getCardbyIndex(index)
								continue
							}
						} else if dian == 5 {
							if card1.number >= 10 || card1.number == 1 ||
								card1.number == 2 || card1.number == 3 ||
								card1.number == 8 || card1.number == 9 {
								index++
								card3 = bjl.getCardbyIndex(index)
								continue
							}
						}
					}
				}

				if bjl.getCardsDian(cardarr1) >= 1 {
					break
				}

				index++
				card3 = bjl.getCardbyIndex(index)
			}
		}
	}

	cards = append(cards, card1)
	cards = append(cards, card2)
	if hasThird {
		bjl.delCard(card3)
		cards = append(cards, card3)
	}

	return cards, true
}

//获取第二个三张的牌型
func (bjl *BaiJialeLogic) getThreeSencondCards(cards1 []*card, v int32, isDui bool, isShun bool, isZhuang bool) ([]*card, bool) {
	cards2 := make([]*card, 0)
	index := int32(0)
	card1 := bjl.getCardbyIndex(index)
	card2 := new(card)
	card3 := new(card)
	hasThird := true

	if bjl.zuhetype&(1<<9) != 0 {
		//同点和
		card1 = bjl.getCardByNumber(cards1[0].number)
		bjl.delCard(card1)
		card2 = bjl.getCardByNumber(cards1[1].number)
		bjl.delCard(card2)
		if len(cards1) == 3 {
			card3 = bjl.getCardByNumber(cards1[2].number)
			hasThird = true
		} else {
			hasThird = false
		}
	} else if isDui {
		//对子
		n := rand.Int31n(3)
		//这里决定牌是和第一张相同还是第二张
		if card1.number != 4 && card1.number != 9 && n == 0 {
			n++
		}

		if n == 0 {
			for {
				if card1 == nil {
					return cards2, false
				}
				cardarr1 := make([]*card, 0)
				cardarr1 = append(cardarr1, card1)
				cardarr1 = append(cardarr1, card1)

				if bjl.compareCards(cards1, cardarr1) == v {
					dian := bjl.getCardsDian(cardarr1)
					if isZhuang {
						if dian >= 7 {
							hasThird = false
						}
					} else {
						if dian >= 6 {
							hasThird = false
						}
					}

					if isZhuang {
						if len(cards1) == 3 {
							if dian == 3 {
								if cards1[2].number == 8 {
									hasThird = false
								}
							} else if dian == 4 {
								if cards1[2].number >= 10 || cards1[2].number == 1 ||
									cards1[2].number == 8 || cards1[2].number == 9 {
									hasThird = false
								}
							} else if dian == 5 {
								if cards1[2].number >= 10 || cards1[2].number == 1 ||
									cards1[2].number == 2 || cards1[2].number == 3 ||
									cards1[2].number == 8 || cards1[2].number == 9 {
									hasThird = false
								}
							} else if dian == 6 {
								if cards1[2].number != 6 && cards1[2].number != 7 {
									hasThird = false
								}
							}
						} else {
							if dian == 6 {
								hasThird = false
							}
						}
					}
					break
				}

				index++
				card1 = bjl.getCardbyIndex(index)
			}

			bjl.delCard(card1)
			card2 = bjl.getCardByNumber(card1.number)
			bjl.delCard(card2)
		} else {
			cardarr2 := make([]*card, 0)
			cardarr2 = append(cardarr2, cards1[0])
			cardarr2 = append(cardarr2, cards1[1])

			dian2 := bjl.getCardsDian(cardarr2)

			if len(cards1) == 3 {
				for {
					if card1 == nil {
						return cards2, false
					}
					if !isZhuang {
						if dian2 == 3 {
							if card1.number == 8 && v == 1 {
								index++
								card1 = bjl.getCardbyIndex(index)
								mlog.NewDefaultLogger().SetData("1241card1").Debug("获取第二个三张的牌型")
								continue
							}

							break
						} else if dian2 == 4 {
							if (card1.number == 8 || card1.number == 9 ||
								card1.number == 1 || card1.number >= 10) && v == 1 {
								index++
								card1 = bjl.getCardbyIndex(index)
								mlog.NewDefaultLogger().SetData("1251card1").Debug("获取第二个三张的牌型")
								continue
							}

							break
						} else if dian2 == 5 {
							//if (card1.number != int32(4) || card1.number != int32(5) ||
							//	card1.number != int32(6) || card1.number != int32(7))
							if true && v == int32(1) {
								index++
								card1 = bjl.getCardbyIndex(index)
								mlog.NewDefaultLogger().SetData("1261card1").Debug("获取第二个三张的牌型")
								continue
							}
							break
						} else {
							break
						}
					} else {
						break
					}
				}
			}

			bjl.delCard(card1)
			card2 = bjl.getCardbyIndex(index)
			for {

				if card2 == nil {
					return cards2, false
				}
				cardarr1 := make([]*card, 0)
				cardarr1 = append(cardarr1, card1)
				cardarr1 = append(cardarr1, card2)
				dian := bjl.getCardsDian(cardarr1)

				if !isZhuang && len(cards1) == 3 {
					if dian2 == 3 {
						if card2.number == 8 && v != 1 {
							index++
							card2 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("1287card2").Debug("获取第二个三张的牌型")
							continue
						}
					} else if dian2 == 4 {
						if (card2.number == 8 || card2.number == 9 ||
							card2.number == 1 || card2.number >= 10) && v != 1 {
							index++
							card2 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("1295card2").Debug("获取第二个三张的牌型")
							continue
						}
					} else if dian2 == 5 {
						//if (card2.number != int32(4) || card2.number != int32(5) ||
						//	card2.number != int32(6) || card2.number != int32(7))
						if true && v != int32(1) {
							index++
							card2 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("1303card2").Debug("获取第二个三张的牌型")
							continue
						}
					}
				}

				if isZhuang {
					if len(cards1) == 3 {
						if dian == 3 {
							if cards1[2].number == 8 {
								index++
								card2 = bjl.getCardbyIndex(index)
								mlog.NewDefaultLogger().SetData("1389card2").Debug("获取第二个三张的牌型")
								continue
							}
						} else if dian == 4 {
							if cards1[2].number >= 10 || cards1[2].number == 1 ||
								cards1[2].number == 8 || cards1[2].number == 9 {
								index++
								card2 = bjl.getCardbyIndex(index)
								mlog.NewDefaultLogger().SetData("1397card2").Debug("获取第二个三张的牌型")
								continue
							}
						} else if dian == 5 {
							if cards1[2].number >= 10 || cards1[2].number == 1 ||
								cards1[2].number == 2 || cards1[2].number == 3 ||
								cards1[2].number == 8 || cards1[2].number == 9 {
								index++
								card2 = bjl.getCardbyIndex(index)
								mlog.NewDefaultLogger().SetData("1406card2").Debug("获取第二个三张的牌型")
								continue
							}
						} else if dian == 6 {
							if cards1[2].number != 6 && cards1[2].number != 7 {
								index++
								card2 = bjl.getCardbyIndex(index)
								mlog.NewDefaultLogger().SetData("1413card2").Debug("获取第二个三张的牌型")
								continue
							}
						}
					} else {
						if dian2 == 6 {
							index++
							card2 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("1303card2").Debug("获取第二个三张的牌型")
							continue
						}
					}
				}

				if dian >= 6 {
					index++
					card2 = bjl.getCardbyIndex(index)
					mlog.NewDefaultLogger().SetData("1312card2").Debug("获取第二个三张的牌型")
					continue
				}

				if n == 1 {
					cardarr1 = append(cardarr1, card1)
				} else {
					cardarr1 = append(cardarr1, card2)
				}

				if bjl.compareCards(cards1, cardarr1) == v {
					break
				}

				index++
				card2 = bjl.getCardbyIndex(index)
			}
		}

		bjl.delCard(card2)

		if n == 1 {
			card3 = bjl.getCardByNumber(card1.number)
		} else {
			card3 = bjl.getCardByNumber(card2.number)
		}

		if card3 == nil {
			return cards2, false
		}

	} else if isShun {
		mlog.NewDefaultLogger().SetData("是顺子").Debug("获取第二个三张的牌型")

		cardarr2 := make([]*card, 0)
		cardarr2 = append(cardarr2, cards1[0])
		cardarr2 = append(cardarr2, cards1[1])

		dian2 := bjl.getCardsDian(cardarr2)

		for {
			if card1 == nil {
				return cards2, false
			}
			if !isZhuang && len(cards1) == 3 {
				if card1.number >= 6 {
					index++
					card1 = bjl.getCardbyIndex(index)
					continue
				}
			}

			if card1.number > 11 {
				index++
				card1 = bjl.getCardbyIndex(index)
				continue
			}

			card2.number = card1.number + 1

			cardarr1 := make([]*card, 0)
			cardarr1 = append(cardarr1, card1)
			cardarr1 = append(cardarr1, card2)
			dian := bjl.getCardsDian(cardarr1)

			if isZhuang {
				if len(cards1) == 3 {
					if dian == 3 {
						if cards1[2].number == 8 {
							index++
							card1 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("1303card2").Debug("获取第二个三张的牌型")
							continue
						}
					} else if dian == 4 {
						if cards1[2].number >= 10 || cards1[2].number == 1 ||
							cards1[2].number == 8 || cards1[2].number == 9 {
							index++
							card1 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("1303card2").Debug("获取第二个三张的牌型")
							continue
						}
					} else if dian == 5 {
						if cards1[2].number >= 10 || cards1[2].number == 1 ||
							cards1[2].number == 2 || cards1[2].number == 3 ||
							cards1[2].number == 8 || cards1[2].number == 9 {
							index++
							card1 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("1303card2").Debug("获取第二个三张的牌型")
							continue
						}
					} else if dian == 6 {
						if cards1[2].number != 6 && cards1[2].number != 7 {
							index++
							card1 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("1303card2").Debug("获取第二个三张的牌型")
							continue
						}
					}
				} else {
					if dian2 == 6 {
						index++
						card1 = bjl.getCardbyIndex(index)
						mlog.NewDefaultLogger().SetData("1303card2").Debug("获取第二个三张的牌型")
						continue
					}
				}
			}

			if dian >= 6 {
				index++
				card1 = bjl.getCardbyIndex(index)
				continue
			}

			card3.number = card2.number + 1
			cardarr1 = append(cardarr1, card3)

			if bjl.compareCards(cards1, cardarr1) == v {
				break
			}

			index++
			card1 = bjl.getCardbyIndex(index)
		}
		bjl.delCard(card1)
		card2 = bjl.getCardByNumber(card1.number + 1)
		bjl.delCard(card2)
		card3 = bjl.getCardByNumber(card2.number + 1)
	} else {
		mlog.NewDefaultLogger().SetData("是杂牌").Debug("获取第二个三张的牌型")
		bjl.delCard(card1)
		card2 = bjl.getCardbyIndex(index)

		cardarr2 := make([]*card, 0)
		cardarr2 = append(cardarr2, cards1[0])
		cardarr2 = append(cardarr2, cards1[1])

		dian2 := bjl.getCardsDian(cardarr2)

		for {
			if card2 == nil {
				return cards2, false
			}
			hasThird = true
			if card2.number == card1.number {
				index++
				card2 = bjl.getCardbyIndex(index)
				mlog.NewDefaultLogger().SetData("因为和牌一相同所以下一张").Debug("获取第二个三张的牌型")
				continue
			}

			cardarr1 := make([]*card, 0)
			cardarr1 = append(cardarr1, card1)
			cardarr1 = append(cardarr1, card2)
			dian := bjl.getCardsDian(cardarr1)

			if isZhuang {
				if len(cards1) == 3 {
					if dian == 3 {
						if cards1[2].number == 8 {
							index++
							card2 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("1303card2").Debug("获取第二个三张的牌型")
							continue
						}
					} else if dian == 4 {
						if cards1[2].number >= 10 || cards1[2].number == 1 ||
							cards1[2].number == 8 || cards1[2].number == 9 {
							index++
							card2 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("1303card2").Debug("获取第二个三张的牌型")
							continue
						}
					} else if dian == 5 {
						if cards1[2].number >= 10 || cards1[2].number == 1 ||
							cards1[2].number == 2 || cards1[2].number == 3 ||
							cards1[2].number == 8 || cards1[2].number == 9 {
							index++
							card2 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("1303card2").Debug("获取第二个三张的牌型")
							continue
						}
					} else if dian == 6 {
						if cards1[2].number != 6 && cards1[2].number != 7 {
							index++
							card2 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("1303card2").Debug("获取第二个三张的牌型")
							continue
						}
					}
				} else {
					if dian2 == 6 {
						index++
						card2 = bjl.getCardbyIndex(index)
						mlog.NewDefaultLogger().SetData("1303card2").Debug("获取第二个三张的牌型")
						continue
					}
				}
			}

			if dian >= 6 {
				index++
				card2 = bjl.getCardbyIndex(index)
				mlog.NewDefaultLogger().SetData("因为大于6所以重新取牌card2").Debug("获取第二个三张的牌型")
				continue
			}

			break
		}
		bjl.delCard(card2)
		if hasThird {
			cardarr2 := make([]*card, 0)
			cardarr2 = append(cardarr2, cards1[0])
			cardarr2 = append(cardarr2, cards1[1])
			dian2 := bjl.getCardsDian(cardarr2)

			card3 = bjl.getCardbyIndex(index)
			i := 0
			for {
				if card3 == nil || card2 == nil {
					return cards2, false
				}
				if index >= 46 && i == 0 {
					index = 0
					i++
					tmpcard := poker.NewPokerCard(card1.color, card1.number)
					bjl.cards.AddCard(tmpcard)
					tmpcard.Color = card2.color
					tmpcard.CardNumber = card2.number
					bjl.cards.AddCard(tmpcard)
					card1 = bjl.getCardbyIndex(index)
					bjl.delCard(card1)
					card2 = bjl.getCardbyIndex(index)
					for {
						hasThird = true
						if card2.number == card1.number {
							index++
							card2 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("因为和牌一相同所以下一张").Debug("获取第二个三张的牌型")
							continue
						}

						cardarr1 := make([]*card, 0)
						cardarr1 = append(cardarr1, card1)
						cardarr1 = append(cardarr1, card2)
						dian := bjl.getCardsDian(cardarr1)

						if dian >= 6 {
							index++
							card2 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("因为大于6所以重新取牌card2").Debug("获取第二个三张的牌型")
							continue
						}

						break
					}
					bjl.delCard(card2)
				}
				//不能是对子
				if card1.number == card3.number || card2.number == card3.number {
					index++
					card3 = bjl.getCardbyIndex(index)
					mlog.NewDefaultLogger().SetData("不能是对子card3").Debug("获取第二个三张的牌型")
					continue
				}
				cardarr1 := make([]*card, 0)
				cardarr1 = append(cardarr1, card1)
				cardarr1 = append(cardarr1, card2)
				cardarr1 = append(cardarr1, card3)
				bjl.sortPoker(cardarr1)

				//不能是顺子
				if cardarr1[0].number+1 == cardarr1[1].number && cardarr1[1].number+1 == cardarr1[2].number {
					index++
					card3 = bjl.getCardbyIndex(index)
					mlog.NewDefaultLogger().SetData("不能是顺子card3").Debug("获取第二个三张的牌型")
					continue
				}

				if !isZhuang && len(cards1) == 3 {
					if dian2 == 3 {
						if card3.number == 8 {
							index++
							card3 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("1420card3").Debug("获取第二个三张的牌型")
							continue
						}
					} else if dian2 == 4 {
						if card3.number == 8 || card3.number == 9 ||
							card3.number == 1 || card3.number >= 10 {
							index++
							card3 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("1427card3").Debug("获取第二个三张的牌型")
							continue
						}
					} else if dian2 == 5 {
						//if card3.number != int32(4) || card3.number != int32(5) ||
						//	card3.number != int32(6) || card3.number != int32(7) ||
						//	card3.number >= int32(10) {
						if true {
							index++
							card3 = bjl.getCardbyIndex(index)
							mlog.NewDefaultLogger().SetData("1436card3").Debug("获取第二个三张的牌型")
							continue
						}
					} else if dian2 == 6 {
						//if card3.number != int32(6) || card3.number != int32(7) {
						if true {
							index++
							card3 = bjl.getCardbyIndex(index)
							mlog.Debug("1436card3")
							continue
						}
					}
				}

				c := bjl.compareCards(cards1, cardarr1)
				mlog.NewDefaultLogger().SetData(fmt.Sprintf("比牌结果是c=%d", c)).Debug("获取第二个三张的牌型")
				if bjl.compareCards(cards1, cardarr1) == v {
					break
				}

				index++
				card3 = bjl.getCardbyIndex(index)
			}
		}
	}

	cards2 = append(cards2, card1)
	cards2 = append(cards2, card2)
	if hasThird {
		bjl.delCard(card3)
		cards2 = append(cards2, card3)
	}

	return cards2, true
}

//cards1已经确定，获取两张牌，v为-1cards2大，0等于cards1,1cards1大
//回傳閒的牌對子方式：當參數card1是莊，且v是1（此為莊家贏）就能取得閒輸的牌，反之
func (bjl *BaiJialeLogic) getTwoCards(cards1 []*card, v int32, isDui bool, isTianWang bool) []*card {
	cards2 := make([]*card, 0)
	if len(cards1) != 2 {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("严重错误，牌不是两个：%v", cards1)).Debug("cards1已经确定，获取两张牌，v为-1cards2大，0等于cards1,1cards1大")
		return cards2
	}

	index := int32(0)
	card1 := bjl.getCardbyIndex(index)
	card2 := new(card)
	if bjl.zuhetype&(1<<9) != 0 {
		//同点和
		card1 = bjl.getCardByNumber(cards1[0].number)
		bjl.delCard(card1)
		card2 = bjl.getCardByNumber(cards1[1].number)
	} else if isDui {
		//对子
		for {
			cardarr1 := make([]*card, 0)
			cardarr1 = append(cardarr1, card1)
			cardarr1 = append(cardarr1, card1)
			dian := bjl.getCardsDian(cardarr1)
			if isTianWang { //天王, 牌的組合點數<8再要一張牌讓點數變大
				if dian <= 8 {
					index++
					card1 = bjl.getCardbyIndex(index)
					continue
				}
			} else if dian >= 8 { //非天王，牌的組合點數>8則在要牌讓點數變小
				index++
				card1 = bjl.getCardbyIndex(index)
				continue
			}

			if bjl.compareCards(cards1, cardarr1) == v {
				break
			}

			index++
			card1 = bjl.getCardbyIndex(index)
		}

		bjl.delCard(card1)
		card2 = bjl.getCardByNumber(card1.number)
	} else {
		bjl.delCard(card1)
		card2 = bjl.getCardbyIndex(index)
		for {
			cardarr1 := make([]*card, 0)
			cardarr1 = append(cardarr1, card1)
			cardarr1 = append(cardarr1, card2)
			dian := bjl.getCardsDian(cardarr1)
			if isTianWang {
				if dian <= 8 {
					index++
					card1 = bjl.getCardbyIndex(index)
					continue
				}
			} else if dian >= 8 {
				index++
				card1 = bjl.getCardbyIndex(index)
				continue
			}

			if card1.number == card2.number {
				index++
				card1 = bjl.getCardbyIndex(index)
				continue
			}

			if bjl.compareCards(cards1, cardarr1) == v {
				break
			}

			index++
			card2 = bjl.getCardbyIndex(index)
		}
	}

	bjl.delCard(card2)
	cards2 = append(cards2, card1)
	cards2 = append(cards2, card2)
	return cards2
}

//获取扑克的点数
func (bjl *BaiJialeLogic) getCardsDian(cards []*card) int32 {
	total := int32(0)
	if cards == nil {
		return 0
	}
	for i := 0; i < len(cards); i++ {
		if cards[i].number >= 10 {
			total += 0
		} else {
			total += cards[i].number
		}
	}

	total %= 10
	return total
}

//比较两组扑克的大小
func (bjl *BaiJialeLogic) compareCards(cards1 []*card, cards2 []*card) int32 {
	if bjl.getCardsDian(cards1) > bjl.getCardsDian(cards2) {
		return 1
	} else if bjl.getCardsDian(cards1) == bjl.getCardsDian(cards2) {
		return 0
	} else {
		return -1
	}
}

//检查是否能开出这个奖, 當稅收後的錢大於開獎需要支付的獎金時才能開該類型
func (bjl *BaiJialeLogic) getCanKaiJiangType(money int64, bjltype int32) bool {
	//開獎需要支付的獎金 = 賠率(類型) * 下注金額（類型）
	needmoney := int64(bjl.betvalue.getOddByType(bjltype)) * bjl.getBetMoneyByType(bjltype)
	if money > needmoney {
		return true
	}

	return false
}

//上莊是為系統或機器人
func (bjl *BaiJialeLogic) isZhuangRobotAndSystem() bool {
	var iszhuangrobot bool = false
	ur := bjl.GetUserByUserID(bjl.zhuangID)

	if ur != nil && ur.IsRobot() {
		iszhuangrobot = true
	}

	if bjl.zhuangID == 0 {
		iszhuangrobot = true
	}

	return iszhuangrobot
}

//获取所有的开奖结果
func (bjl *BaiJialeLogic) getAllKaiJiangType(betmoney int64) []int32 {
	alltype := make([]int32, 0)
	for i := 0; i < 3; i++ {
		alltype = append(alltype, int32(i+1))
	}

	//檢查上莊是否為系統或機器人
	var iszhuangrobot bool = bjl.isZhuangRobotAndSystem()

	if !iszhuangrobot { //是真人上莊所有獎都可以開
		return alltype
	}

	return alltype
}

//获取赔钱最少的开奖类型
func (bjl *BaiJialeLogic) getminjiangtype(alltype []int32) []int32 {
	bjl.GameLog.AddGameStep(bjl.SysPool.ToString())

	if bjl.SysPool.GetMode() == gamelogic.SysPoolModeNone { //無水池模式
		return alltype
	}

	if (bjl.getBetMoneyByType(ZHUANG) + bjl.getBetMoneyByType(XIAN)) == 0 {
		return alltype
	}
	sysPool := bjl.SysPool.GetPool()

	hasNewJiang := false
	//@todo poolline 找時間移到數據庫裡
	for _, v := range bjl.poolline.pool { //水池線
		if sysPool < v.limitmoney {
			bjl.GameLog.AddGameStep("*** BJL 使用自定義水線表 *** limit money %d, rate %d", v.limitmoney, v.rate)
			if rand.Int31n(100) < v.rate { //最有利開獎, 最少人下注的開獎清單
				//alltype = alltype[:0]
				hasNewJiang = true
				break
			}
		}
	}

	if hasNewJiang {
		if bjl.getBetMoneyByType(ZHUANG) > bjl.getBetMoneyByType(XIAN) {
			alltype = []int32{XIAN}
		} else if bjl.getBetMoneyByType(ZHUANG) < bjl.getBetMoneyByType(XIAN) {
			alltype = []int32{ZHUANG}
		}
		bjl.GameLog.AddGameStep("水線之下, 開獎:%v", alltype)
	}

	return alltype
}

//获取赔钱最多的开奖类型
func (bjl *BaiJialeLogic) getmaxjiangtype() []int32 {
	if bjl.getRealManBetMoneyByType(ZHUANG) > bjl.getRealManBetMoneyByType(XIAN) {
		return []int32{ZHUANG}
	} else if bjl.getRealManBetMoneyByType(ZHUANG) < bjl.getRealManBetMoneyByType(XIAN) {
		return []int32{XIAN}
	}
	return []int32{ZHUANG, XIAN}
}

//设置开奖状态
func (bjl *BaiJialeLogic) setState(state int32) {
	bjl.gamestate = state
	//bjl.BoardcastMessage(bjl.getStateMsg())
}

//游戏消息处理逻辑
func (bjl *BaiJialeLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	if !bjl.CheckGameUserMessage(user) {
		bjl.KickUserFunc(user.UserID, "非法操作")
		return
	}
	//不是给游戏逻辑服务器的消息
	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	bjl.CheckUserAnimationCommited(user, gmsg, bjl.actionStr, nil)

	if gmsg.SClassID == int32(netproto.BJL_GameMessageClassID_BJLGetShangZhuangListID) {
		bjl.SendMessageToUser(user, bjl.getZhuangListMsg())
		return
	}

	if bjl.gamestate != Bet {
		return
	}

	if gmsg.SClassID == int32(netproto.BJL_GameMessageClassID_BJLBetID) {
		bjl.userBet(user, gmsg.MsgData.([]byte))
		return
	}

	if gmsg.SClassID == int32(netproto.BJL_GameMessageClassID_BJLShangZhuangID) {
		bjl.shangZhuang(user, gmsg.MsgData.([]byte))
		return
	}

	if gmsg.SClassID == int32(netproto.BJL_GameMessageClassID_BJLBetClearID) && bjl.gamestate == Bet {
		bjl.clearUserBet(user)
		return
	}
}

//获取配置文件
func (bjl *BaiJialeLogic) getbjlconfig(cfgs *config.Group) {
	if cfgs == nil {
		mlog.Error("桌號[%d] 讀取 config 失敗", bjl.TableNo())
		return
	}

	bjl.betvalue = newBetconfig(cfgs)
	bjl.poolline = newPoollineConfig(cfgs)
	scfgs := cfgs.GetConfig("gameserver")
	bjl.bettime = int32(scfgs.GetConfigInt("gameserver", "bettime"))
	bjl.kaijiangactiontime = int32(scfgs.GetConfigInt("gameserver", "kaijiangactiontime"))
	bjl.endactiontime = int32(scfgs.GetConfigInt("gameserver", "endactiontime"))
	bjl.readytime = int32(scfgs.GetConfigInt("gameserver", "readywaittime"))
	bjl.zhuanglimit = int64(scfgs.GetConfigInt("gameserver", "zhuanglimit"))
	bjl.maxturn = int32(scfgs.GetConfigInt("gameserver", "zhuangturn"))
	bjl.serverid = strconv.Itoa(int(bjl.GetServerConfig().ServerID))
	bjl.maxclearcount = int32(scfgs.GetConfigInt("gameserver", "maxclearcount"))
	bjl.zhushui = int32(scfgs.GetConfigInt("gameserver", "zhushui"))
	bjl.betend = int32(scfgs.GetConfigInt("gameserver", "betend"))
	bjl.maxrobotzhuang = int32(scfgs.GetConfigInt("gameserver", "maxrobotzhuang"))
	bjl.robotxiazhuang = int32(scfgs.GetConfigInt("gameserver", "robotxiazhuang"))
	bjl.robotzhuang = int32(scfgs.GetConfigInt("gameserver", "robotzhuang"))
	bjl.betamode = int64(scfgs.GetConfigInt("gameserver", "betamode"))
	bjl.sametimebetmode = int64(scfgs.GetConfigInt("gameserver", "sametimebetmode"))
	bjl.betlimitmoney = int64(scfgs.GetConfigInt("gameserver", "betlimitmoney"))
	bjl.publiczhuang = int64(scfgs.GetConfigInt("gameserver", "publiczhuang"))
	bjl.luckystarround = int64(scfgs.GetConfigInt("gameserver", "luckystarround"))
	bjl.DealActionTime = int32(scfgs.GetConfigInt("gameserver", "DealBeforeBetTime"))

	betcountstr := scfgs.GetConfigValue("gameserver", "betcount")

	// 百家樂玩法
	iBJLMode := scfgs.GetConfigInt("gameserver", "BJLMode")
	switch iBJLMode {
	case BJLModeLucky6: // 幸運6百家樂
		bjl.BJLMode = BJLModeLucky6
	default:
		bjl.BJLMode = BJLModeStandard
	}

	strarr := strings.Split(betcountstr, ",")
	for i := 0; i < len(strarr); i++ {
		count, _ := strconv.ParseInt(strarr[i], 10, 32)
		bjl.betcount = append(bjl.betcount, int32(count))
	}

	betrangestr := scfgs.GetConfigValue("baijiale", "1")
	strarr = strings.Split(betrangestr, ",")
	for i := 0; i < len(strarr); i++ {
		zhuangxianrange, _ := strconv.ParseInt(strarr[i], 10, 32)
		bjl.zhuangxianrange = append(bjl.zhuangxianrange, int32(zhuangxianrange))
	}
	betrangestr = scfgs.GetConfigValue("baijiale", "3")
	strarr = strings.Split(betrangestr, ",")
	for i := 0; i < len(strarr); i++ {
		herange, _ := strconv.ParseInt(strarr[i], 10, 32)
		bjl.herange = append(bjl.herange, int32(herange))
	}
	betrangestr = scfgs.GetConfigValue("baijiale", "6")
	strarr = strings.Split(betrangestr, ",")
	for i := 0; i < len(strarr); i++ {
		duizirange, _ := strconv.ParseInt(strarr[i], 10, 32)
		bjl.duizirange = append(bjl.duizirange, int32(duizirange))
	}

	bjl.ctrlbetlimit = int64(scfgs.GetConfigInt("supertool", "ctrlbetlimit"))
	if bjl.ctrlbetlimit <= 0 {
		bjl.ctrlbetlimit = CtrlBetLimitDefault
	}
}

//清除玩家下注
func (bjl *BaiJialeLogic) clearUserBet(user *gamelogic.GameUser) {
	ud := bjl.getUserDataByID(user.UserID)
	if ud == nil {
		return
	}

	if ud.clearcount >= bjl.maxclearcount {
		bjl.SendMessageToUser(user, bjl.getClearRetMsg(ud.money, -1))
		return
	}

	ud.clearcount++
	for _, ub := range ud.getUserBet() {
		bjl.updateBet(ub.baijialetype, -ub.money)
	}

	bjl.GameLog.AddGameStep("用户%d点击清除", user.UserID)
	ud.clearData()
	bjl.SendMessageToUser(user, bjl.getClearRetMsg(user.Score, 0))
	bjl.BoardcastMessage(bjl.getBetTotalMsg(user.UserID))
}

//获取清除结果的消息
func (bjl *BaiJialeLogic) getClearRetMsg(money int64, ret int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetClearRetID)

	pdata := new(netproto.BJL_BetClearRet)
	pdata.Money = proto.Int64(money)
	pdata.Ret = proto.Int32(ret)
	pmsg.MsgData = pdata

	return pmsg
}

func (bjl *BaiJialeLogic) getUserBetByType(ud *baijialeuserdata, t int32) int64 {
	money := int64(0)
	for _, bet := range ud.userbet {
		if bet.baijialetype == t {
			money += bet.money
		}
	}
	return money
}

//处理用户下注
func (bjl *BaiJialeLogic) userBet(user *gamelogic.GameUser, data []byte) {
	ud := bjl.getUserDataByID(user.UserID)
	if ud == nil {
		ud = bjl.addUserData(user)
	}

	ust := &netproto.BJL_BetArr{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		return
	}

	//庄家不能下注
	if bjl.zhuangID == user.UserID {
		return
	}

	money := int32(0)
	for _, b := range ust.BetArr {
		*b.BetType = int32(math.Abs(float64(*b.BetType)))
		*b.Count = int32(math.Abs(float64(*b.Count)))

		//玩家身上金額小於下注限制
		if user.Score < bjl.betlimitmoney*100 {
			for i, u := range bjl.userdata {
				if u.user.UserID == ud.user.UserID {
					bjl.userdata = append(bjl.userdata[:i], bjl.userdata[i+1:]...)
				}
			}
			bjl.SendMessageToUser(user, bjl.getBetRetMsg(ud, -6, *b.BetType, ust.BetArr))
			return
		}

		if bjl.sametimebetmode != 1 { //關閉庄閒同時下注的配置
			//庄閒不能同時下注
			if (*b.BetType == int32(ZHUANG) && bjl.getUserBetByType(ud, XIAN) > 0) ||
				(*b.BetType == int32(XIAN) && bjl.getUserBetByType(ud, ZHUANG) > 0) {
				bjl.SendMessageToUser(user, bjl.getBetRetMsg(ud, -5, *b.BetType, ust.BetArr))
				return
			}
		}

		tmpmon := bjl.getMoneyByType(*b.BetType) * (*b.Count)

		// 系統下注上限
		if (int64(tmpmon) + bjl.getBetMoneyByType(*b.BetType)) > bjl.betvalue.getMaxMoneyByType(*b.BetType) {
			bjl.SendMessageToUser(user, bjl.getBetRetMsg(ud, -4, *b.BetType, ust.BetArr))
			return
		}

		// 不是系統當庄, 則壓注上限為庄家身上的錢所能夠賠付的錢
		if bjl.zhuangID != 0 {
			ur := bjl.GetUserByUserID(bjl.zhuangID)
			odd := int64(0)
			if odd = int64(bjl.betvalue.getOddByType(b.GetBetType())) - 1; odd < 0 {
				odd = 0
			}
			maxmoney := (int64(tmpmon) + bjl.getBetMoneyByType(*b.BetType)) * odd
			if maxmoney > ur.Score {
				bjl.SendMessageToUser(user, bjl.getBetRetMsg(ud, -4, *b.BetType, ust.BetArr))
				return
			}
		}

		money += tmpmon
	}

	ret := int32(0)
	if int64(money) > ud.money {
		ret = -1
	} else {
		for _, b := range ust.BetArr {
			betmon := bjl.getMoneyByType(*b.BetType) * (*b.Count)
			ud.setUserBet(*b.BetType, betmon)
			bjl.updateBet(*b.BetType, int64(betmon))
			if user.IsRobot() {
				bjl.updateRobotBet(*b.BetType, int64(betmon))
			}
		}
		if user.UserID == bjl.luckystar {
			bjl.notbettimes = 0
		}
		// bjl.GameLog.AddGameStep("用户%d下注%v", user.UserID, ust.BetArr)
		bjl.BoardcastMessage(bjl.getBetTotalMsg(user.UserID))
	}

	bjl.SendMessageToUser(user, bjl.getBetRetMsg(ud, ret, 0, ust.BetArr))
}

//获取下注的总钱数，通过游戏类型
func (bjl *BaiJialeLogic) getBetMoneyByType(t int32) int64 {
	for i := 0; i < len(bjl.bettypetotal); i++ {
		if bjl.bettypetotal[i].baijialetype == t {
			return bjl.bettypetotal[i].money
		}
	}

	return 0
}

//获取此类型的机器人下注总额
func (bjl *BaiJialeLogic) getRobotBetMoneyByType(t int32) int64 {
	for i := 0; i < len(bjl.bettyperobot); i++ {
		if bjl.bettyperobot[i].baijialetype == t {
			return bjl.bettyperobot[i].money
		}
	}

	return 0
}

//获取下注需要的钱，通过下注类型
func (bjl *BaiJialeLogic) getMoneyByType(t int32) int32 {
	return bjl.betvalue.getMoneyByType(t)
}

//更新下注的钱
func (bjl *BaiJialeLogic) updateBet(bjlType int32, money int64) {
	bjl.bettotal += int64(money)
	//mlog.Debug("BetType=%v,Money=%v,bjl.bettotal=%v", bjlType, Money, bjl.bettotal)
	for _, v := range bjl.bettypetotal {
		if v.baijialetype == bjlType {
			v.money += int64(money)
			if v.money < 0 {
				v.money = 0
			}
			return
		}
	}

	b := new(bettype)
	b.baijialetype = bjlType
	b.money = int64(money)
	if b.money < 0 {
		b.money = 0
	}
	bjl.bettypetotal = append(bjl.bettypetotal, b)
}

//更新机器人下注的钱
func (bjl *BaiJialeLogic) updateRobotBet(t int32, money int64) {
	bjl.bettotalrobot += money
	for _, v := range bjl.bettyperobot {
		if v.baijialetype == t {
			v.money += int64(money)
			if v.money < 0 {
				v.money = 0
			}
			return
		}
	}

	b := new(bettype)
	b.baijialetype = t
	b.money = int64(money)
	if b.money < 0 {
		b.money = 0
	}
	bjl.bettyperobot = append(bjl.bettyperobot, b)
}

//移除一机器人下注的钱
func (bjl *BaiJialeLogic) removeRobotBet() {
	for _, v := range bjl.bettyperobot {
		bjl.updateBet(v.baijialetype, -v.money)
	}
}

//获取机器人赢了多少钱
func (bjl *BaiJialeLogic) getRobotWin() int64 {
	win := int64(0)
	for _, rtd := range bjl.bettyperobot {
		for _, t := range bjl.jiangtype {
			if rtd.baijialetype == t {
				win += rtd.money * int64(bjl.betvalue.getOddByType(t))
				break
			}
		}
	}

	return win
}

//获取机器人下注赢最多钱的类型
func (bjl *BaiJialeLogic) getMaxRobotWinType() int32 {
	example.invalid("获取机器人赢最多钱的")
	maxtype := int32(1)
	maxmoney := int64(0)
	for i := int32(1); i <= 3; i++ {
		for _, rtd := range bjl.bettyperobot {
			if i == rtd.baijialetype {
				money := rtd.money * int64(bjl.betvalue.getOddByType(i))
				if money > maxmoney {
					maxtype = i
					maxmoney = money
				}
				break
			}
		}
	}

	return maxtype
}

//获取状态消息
func (bjl *BaiJialeLogic) getStateMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BJL_GameMessageClassID_BJLStateID)

	pdata := new(netproto.BJL_Sate)
	pdata.State = proto.Int32(bjl.gamestate)
	pdata.Time = proto.Int32(bjl.getBetTime())
	pmsg.MsgData = pdata
	return pmsg
}

func (bjl *BaiJialeLogic) getBetTime() int32 {
	time := int32(0)
	if bjl.opttimer != nil {
		time = bjl.opttimer.WaitSeconds - bjl.opttimer.ElapsedSeconds()
	}
	return time
}

//获取玩家数据通过用户ID
func (bjl *BaiJialeLogic) getUserDataByID(userid int32) *baijialeuserdata {
	for _, ud := range bjl.userdata {
		if ud.user.UserID == userid {
			return ud
		}
	}

	return nil
}

//添加一个玩家数据
func (bjl *BaiJialeLogic) addUserData(user *gamelogic.GameUser) *baijialeuserdata {
	ud := newBaiJiaLeUserData(user)
	bjl.userdata = append(bjl.userdata, ud)
	return ud
}

//删除一个玩家数据
func (bjl *BaiJialeLogic) delUserData(user *gamelogic.GameUser) {
	for i := 0; i < len(bjl.userdata); i++ {
		if bjl.userdata[i].user.UserID == user.UserID {
			bjl.userdata = append(bjl.userdata[:i], bjl.userdata[i+1:]...)
			break
		}
	}
}

//获取庄消息
func (bjl *BaiJialeLogic) getZhuangInfoMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BJL_GameMessageClassID_BJLZhuangInfoID)

	pdata := new(netproto.BJL_ZhuangInfo)
	pdata.UserID = proto.Int32(0)
	pdata.LastTurn = proto.Int32(0)
	ur := bjl.GetUserByUserID(bjl.zhuangID)
	if ur != nil {
		pdata.UserInfo = append(pdata.UserInfo, ur.GetProtoBufGameUser())
		pdata.UserID = proto.Int32(ur.UserID)
		pdata.LastTurn = proto.Int32(bjl.lastturn)
	}
	pmsg.MsgData = pdata
	return pmsg
}

//获取上庄列表消息
func (bjl *BaiJialeLogic) getZhuangListMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BJL_GameMessageClassID_BJLShangZhuangListInfoID)

	pdata := new(netproto.BJL_ZhuangList)
	for i := 0; i < len(bjl.zhuanglist); i++ {
		ur := bjl.GetUserByUserID(bjl.zhuanglist[i])
		pdata.UserInfo = append(pdata.UserInfo, ur.GetProtoBufGameUser())
	}
	pmsg.MsgData = pdata
	return pmsg
}

//检测用户是否能上庄
func (bjl *BaiJialeLogic) checkUserShangZhuang(userid int32) bool {
	ur := bjl.GetUserByUserID(userid)
	if ur != nil && ur.Score >= bjl.zhuanglimit {

		return true
	}

	return false
}

//换庄
func (bjl *BaiJialeLogic) huangZhuang() {
	if bjl.publiczhuang == 1 {
		if bjl.zhuangID == 0 {
			//庄是系统直接换庄
			bjl.setNextZhuang()
			return
		}
		if bjl.lastturn <= 0 {
			//已经到最后一局了
			bjl.setNextZhuang()
			return
		}
		ur := bjl.GetUserByUserID(bjl.zhuangID)
		if ur != nil && ur.Score < bjl.zhuanglimit {
			bjl.setNextZhuang()
			return
		}

		if ur == nil {
			bjl.setNextZhuang()
			return
		}

		bjl.lastturn--
	} else {
		mlog.Debug("未開啟上庄機制，不用換庄")
		return
	}
}

//检测用户是否能上庄
func (bjl *BaiJialeLogic) checkNetUserCanShangZhuang() {
	if len(bjl.zhuanglist) > 0 {
		for {
			user := bjl.GetUserByUserID(bjl.zhuanglist[0])
			if user != nil && user.Score >= int64(bjl.zhuanglimit) {
				break
			}

			bjl.delFromZhuangList(bjl.zhuanglist[0])
			if len(bjl.zhuanglist) <= 0 {
				break
			}
		}
	}
}

//换下一个人上庄
func (bjl *BaiJialeLogic) setNextZhuang() {
	if bjl.publiczhuang == 1 {
		bjl.checkNetUserCanShangZhuang()
		if len(bjl.zhuanglist) > 0 {
			bjl.zhuangID = bjl.zhuanglist[0]
			bjl.lastturn = bjl.maxturn
			bjl.lastturn--
			bjl.delFromZhuangList(bjl.zhuangID)
			user := bjl.GetUserByUserID(bjl.zhuangID)
			bjl.startmoney = user.Score
		} else {
			bjl.zhuangID = 0
			bjl.startmoney = bjl.startsyspool
			//如果需要系統上莊, 則呼叫機器人登入上莊
			//bjl.useRobotShangZhuang()
		}
	} else {
		mlog.Error("未開啟上庄機制，不得上庄")
		return
	}
}

func (bjl *BaiJialeLogic) useRobotShangZhuang() {
	users := bjl.AllTableUser()
	var robotUsers []*gamelogic.GameUser
	robotCount := 0
	for _, user := range users {
		if user.IsRobot() {
			robotUsers = append(robotUsers, user)
			robotCount++
		}
	}
	if robotCount > 0 {
		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(robotCount)
		bjl.zhuangID = robotUsers[x].UserID
		bjl.lastturn = bjl.maxturn
		bjl.lastturn--
		bjl.startmoney = robotUsers[x].Score
	} else {
		// @todo 呼叫機器人進房
	}
}

//获取下注的结果消息
func (bjl *BaiJialeLogic) getBetRetMsg(ud *baijialeuserdata, ret int32, t int32, betinfo []*netproto.BJL_Bet) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetRetID)

	pdata := new(netproto.BJL_BetRet)
	pdata.Money = proto.Int64(ud.money)
	pdata.Ret = proto.Int32(ret)
	pdata.Type = proto.Int32(t)
	pdata.BetArr = betinfo
	pmsg.MsgData = pdata
	return pmsg
}

//获取总下注的消息
func (bjl *BaiJialeLogic) getBetTotalMsg(userid int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBroadCastBetID)

	pdata := new(netproto.BJL_BroadCastBetArr)
	pdata.UserID = proto.Int32(userid)
	lucky := bjl.getUserDataByID(bjl.luckystar) //獲取明燈
	for _, v := range bjl.betvalue.AllType {
		betinfo := new(netproto.BJL_BroadCastBet)
		betinfo.Type = proto.Int32(v)
		betinfo.Money = proto.Int64(bjl.getBetMoneyByType(v))
		betinfo.LuckyStar = proto.Bool(false)
		if lucky != nil && lucky.getUserBetByType(v) > 0 {
			betinfo.LuckyStar = proto.Bool(true)
		}
		pdata.BroadCastBetArr = append(pdata.BroadCastBetArr, betinfo)
	}

	pmsg.MsgData = pdata

	return pmsg
}

//获取上庄列表的索引通过用户ID
func (bjl *BaiJialeLogic) getZhuangListIndexByUserId(userid int32) int32 {
	for i := 0; i < len(bjl.zhuanglist); i++ {
		if bjl.zhuanglist[i] == userid {
			return int32(i)
		}
	}

	return -1
}

//添加到上庄列表
func (bjl *BaiJialeLogic) addToZhuangList(userid int32) {
	bjl.zhuanglist = append(bjl.zhuanglist, userid)
}

//将玩家从上庄列表中删除
func (bjl *BaiJialeLogic) delFromZhuangList(userid int32) {
	for i := 0; i < len(bjl.zhuanglist); i++ {
		if bjl.zhuanglist[i] == userid {
			bjl.zhuanglist = append(bjl.zhuanglist[:i], bjl.zhuanglist[i+1:]...)
			break
		}
	}
}

//获取上庄结果消息
func (bjl *BaiJialeLogic) getShangzhuangRetMsg(ret int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BJL_GameMessageClassID_BJLShangZhuangRetID)

	pdata := new(netproto.BJL_ShangzhuangRet)
	pdata.Ret = proto.Int32(ret)
	pmsg.MsgData = pdata
	return pmsg
}

//处理玩家上庄
func (bjl *BaiJialeLogic) shangZhuang(user *gamelogic.GameUser, data []byte) {
	ust := &netproto.BJL_ShangZhuang{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		return
	}

	if bjl.getZhuangListIndexByUserId(user.UserID) != -1 {
		bjl.SendMessageToUser(user, bjl.getShangzhuangRetMsg(-2))
		return
	}

	if user.Score < bjl.zhuanglimit {
		bjl.SendMessageToUser(user, bjl.getShangzhuangRetMsg(-1))
		return
	}

	if *ust.IsShangZhuang {
		if bjl.publiczhuang == 1 {
			bjl.GameLog.AddGameStep("用户%d上莊", user.UserID)
			bjl.addToZhuangList(user.UserID)
			bjl.SendMessageToUser(user, bjl.getShangzhuangRetMsg(0))
		} else {
			mlog.Error("未開啟上庄機制，不得上庄")
			return
		}
	} else {
		bjl.GameLog.AddGameStep("用户%d下莊", user.UserID)
		bjl.delFromZhuangList(user.UserID)
		bjl.SendMessageToUser(user, bjl.getShangzhuangRetMsg(1))
		bjl.lastturn = -1
	}

	bjl.BoardcastMessage(bjl.getZhuangListMsg())
}

//获取水池的钱
/*
func (bjl *BaiJialeLogic) getLastPool() {
	var folderName = "lastpool"
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err = os.MkdirAll(folderName, os.ModePerm)
	}

	var tablenostr = strconv.Itoa(int(bjl.TableNo()))
	var filename = bjl.serverid + "_" + tablenostr + ".conf"
	filename = path.Join(folderName, filename)

	_, err := os.Stat(filename)
	var f *os.File = nil
	if os.IsNotExist(err) {
		f, err = os.Create(filename)
	} else {
		f, err = os.OpenFile(filename, os.O_RDWR, 0)
	}

	if err != nil {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("%v文件不存在", filename)).Error("获取水池的钱")
		panic(err)
	} else {
		bjl.serverfile = f
		defer bjl.serverfile.Close()
		bjl.readPoolFromFile()
	}
}

//读取上一次存储的水池钱数
func (bjl *BaiJialeLogic) readPoolFromFile() {
	buff := make([]byte, 1024)
	n, _ := bjl.serverfile.Read(buff)
	if n > 0 {
		poolstr := string(buff)[0:n]
		//pool, err := strconv.ParseInt(poolstr, 10, 32)
		pool, _ := strconv.Atoi(poolstr)
		bjl.syspool = int64(pool)
	}
}
*/
//写水池的钱
//func (bjl *BaiJialeLogic) writePoolToFile() {
//	buff := strconv.Itoa(int(bjl.syspool))
//	//bufio.ReadWriter.Peek()
//	tablenostr := strconv.Itoa(int(bjl.TableNo()))
//	filename := "./lastpool/" + bjl.serverid + "_" + tablenostr + ".conf"
//	_, err := os.Stat(filename)
//	var f *os.File = nil
//	f, err = os.Create(filename)
//	if err != nil {
//		mlog.NewDefaultLogger().SetData(fmt.Sprintf("%v文件不存在", filename)).Error("写水池的钱")
//		panic(err)
//	} else {
//		f.WriteString(buff)
//		defer f.Close()
//	}
//}

//改变水池的钱
func (bjl *BaiJialeLogic) changeSysPool(value int64) {
	change := value - bjl.SysPool.GetPool()
	bjl.SysPool.UpdateSysPool(change)
}

//改变注水的百分比
func (bjl *BaiJialeLogic) changeZhuShui(value int32) {
	if value >= 0 && value <= 100 {
		bjl.zhushui = value
	}
}

func (bjl *BaiJialeLogic) checkSysPool() {
	mlog.NewDefaultLogger().Debug("桌子号：%v，系统水池：%v，注水概率：%v", bjl.TableNo(), bjl.SysPool.GetPool(), bjl.zhushui)
}

//对扑克进行排序
func (bjl *BaiJialeLogic) sortPoker(cards []*card) {
	arrlen := len(cards)
	for i := 0; i < arrlen; i++ {
		for j := 0; j+1 < arrlen-i; j++ {
			if cards[j].number > cards[j+1].number {
				cards[j], cards[j+1] = cards[j+1], cards[j]
			}
		}
	}
}

//占位防止崩溃
func (bjl *BaiJialeLogic) ExchangeMaxCard(userid int32) {

}

func (bjl *BaiJialeLogic) OnRecvHTTPMessage(req *http.Request) proto.Message {
	action := req.PostFormValue("action")
	if action == "getserverinfo" {
		return bjl.getTableListMsg()
	} else if action == "gettableinfo" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		return bjl.getTableDetailMsg(int32(tid))
	} else if action == "modresult" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		ret := req.PostFormValue("result")
		adminuser := req.PostFormValue("adminuser")
		return bjl.setNextJiangType(int32(tid), ret, adminuser)
	} else if action == "modpoolamount" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 64)
		typeid, _ := strconv.ParseInt(req.PostFormValue("type"), 10, 32)
		pool, _ := strconv.ParseInt(req.PostFormValue("poolamount"), 10, 64)
		operator := req.PostFormValue("operator")
		if typeid == 1 {
			return bjl.setSysPool(int32(tid), pool, operator)
		} else {
			return bjl.setFangFen(int32(tid), pool, operator)
		}
	} else if action == "modpopwaterrate" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		zhushui, _ := strconv.ParseInt(req.PostFormValue("popwaterrate"), 10, 32)
		return bjl.setZhuShui(int32(tid), int32(zhushui))
	} else if action == "history" {
		pdata := new(netproto.BJL_ResultCardList)

		if len(bjl.resultcardlist.CardList) > 20 {
			pdata.CardList = bjl.resultcardlist.CardList[len(bjl.resultcardlist.CardList)-20 : len(bjl.resultcardlist.CardList)]
		} else {
			pdata = bjl.resultcardlist
		}

		return pdata
	}
	return nil
}

//获取桌子列表数据
func (bjl *BaiJialeLogic) getTableListMsg() proto.Message {
	pdata := new(netproto.Http_BJL_GameInfo)
	for _, table := range bjl.GetTable().Room.Tables {
		pTable := new(netproto.Http_BJL_TableList)
		logic := table.GameLogic.(*BaiJialeLogic)
		pTable.TableNo = proto.Int32(logic.TableNo())
		pTable.Status = proto.Int32(logic.gamestate)
		robotcount := int32(0)
		realusercount := int32(0)
		for _, user := range logic.AllGameUser() {
			if user.IsRobot() {
				robotcount++
			} else {
				realusercount++
			}
		}
		pTable.GameUserCount = proto.Int32(realusercount)
		pTable.RobotCount = proto.Int32(robotcount)
		pTable.PoolCurrentAmount = proto.Int64(logic.SysPool.GetPool())
		lastturn := logic.lastturn
		if lastturn < 0 {
			lastturn = 0
		}

		pTable.SumEarnAmount = proto.Int64(logic.systotalwin)
		pTable.LastEarnAmount = proto.Int32(int32(logic.syslastwin))
		pTable.SysTaxIncome = proto.Int64(logic.systaxincome)

		pdata.TableList = append(pdata.TableList, pTable)
	}

	return pdata
}

//获取桌子详细数据
func (bjl *BaiJialeLogic) getTableDetailMsg(tableid int32) proto.Message {
	t := bjl.GetTable().Room.GetTableByNo(tableid)
	pTable := new(netproto.Http_BJL_TableDetail)
	if t != nil {
		logic := t.GameLogic.(*BaiJialeLogic)
		pTable.TableNo = proto.Int32(logic.TableNo())
		pTable.Status = proto.Int32(logic.gamestate)
		//		pTable.PoolInitAmount = proto.Int64(logic.startsyspool)
		pTable.PoolCurrentAmount = proto.Int64(logic.SysPool.GetPool())
		robotcount := int32(0)
		realusercount := int32(0)
		for _, user := range logic.AllGameUser() {
			if user.IsRobot() {
				robotcount++
			} else {
				realusercount++
			}
		}
		pTable.GameUserCount = proto.Int32(realusercount)
		pTable.RobotCount = proto.Int32(robotcount)
		//		pTable.PoolTaxAmount = proto.Int64(logic.totaltax)
		pTable.PopWaterRate = proto.Int32(logic.zhushui)
		zhuang := logic.GetUserByUserID(logic.zhuangID)
		zhuangtype := int32(0)
		nickname := ""
		currmoney := int64(0)
		if zhuang != nil {
			if zhuang.IsRobot() {
				zhuangtype = 1
			} else {
				zhuangtype = 2
				currmoney = zhuang.Score
			}

			nickname = zhuang.NickName
		}
		pTable.ZhuangType = proto.Int32(zhuangtype)
		pTable.ZhuangNickName = proto.String(nickname)
		pTable.ZhuangUserID = proto.Int32(logic.zhuangID)
		lastturn := logic.lastturn
		if lastturn < 0 {
			lastturn = 0
		}
		pTable.ZhuangSurplusTotal = proto.Int32(lastturn)
		pTable.ZhuangInitMoneyAmount = proto.Int64(logic.startmoney)
		pTable.ZhuangCurrentMoneyAmount = proto.Int64(currmoney)
		pTable.SumEarnAmount = proto.Int64(logic.systotalwin)
		pTable.LastEarnAmount = proto.Int32(int32(logic.syslastwin))
		pTable.SysTaxIncome = proto.Int64(logic.systaxincome)
		pTable.OpenResultZhuang = proto.Int64(logic.getRealUserWinByType(ZHUANG))
		pTable.OpenResultXian = proto.Int64(logic.getRealUserWinByType(XIAN))
		pTable.OpenResultHe = proto.Int64(logic.getRealUserWinByType(HE))
		pTable.BetTime = proto.Int32(bjl.getBetTime())
		if logic.gamestate == KaiJiangAction {
			pTable.CurentOpenResult = proto.String(logic.getResultString())
			win := int64(0)
			for i := 0; i < len(logic.jiangtype); i++ {
				win += logic.getRealUserWinByType(logic.jiangtype[i])
			}
			pTable.CurentOpenResultEarn = proto.Int64(win)
		}

		pTable.UserInfo = logic.responseCurrentUserBetList()
		pTable.LastGameUserBet = logic.responseLastGameUserBetList()
		btwrate := logic.getBetTypeWinRate()
		pTable.ZhuangWinRate = proto.Int32(btwrate[0])
		pTable.HeWinRate = proto.Int32(btwrate[2])
		pTable.XianWinRate = proto.Int32(btwrate[1])
		pTable.ZhuangDZWinRate = proto.Int32(btwrate[3])
		pTable.XianDZWinRate = proto.Int32(btwrate[4])
		pTable.WinLine = proto.Int64(bjl.poolline.pool[0].limitmoney)
		pTable.FanFengTarget = proto.Int64(logic.fangFenTarget)
		pTable.FanFengComplete = proto.Int64(logic.fangFenComplete)

	}

	return pTable
}

func (bjl *BaiJialeLogic) getResultString() string {
	var res string
	if len(bjl.jiangtype) <= 0 {
		return res
	}
	res += resultstring[bjl.jiangtype[0]]

	if bjl.zuhetype&(1<<5) != 0 {
		res += "," + "庄对"
	}
	if bjl.zuhetype&(1<<6) != 0 {
		res += "," + "闲对"
	}

	return res
}

func (bjl *BaiJialeLogic) getRetMsg(ret bool, msg string) proto.Message {
	pdata := new(netproto.Http_Result)
	pdata.Result = proto.Bool(ret)
	pdata.Message = proto.String(msg)
	return pdata
}

//设置下次开奖结果
func (bjl *BaiJialeLogic) setNextJiangType(tableid int32, jiang string, adminuser string) proto.Message {
	t := bjl.GetTable().Room.GetTableByNo(tableid)
	if t != nil {
		var ret []int32
		logic := t.GameLogic.(*BaiJialeLogic)
		if adminuser != "" {
			logic.GameLog.SetNextJiangRecord(adminuser, fmt.Sprintf("%s", jiang))
		}
		ss := strings.Split(jiang, ",")
		for _, s := range ss {
			i := util.StrToInt32(s)
			if i != 0 {
				ret = append(ret, i)
			}
		}
		if len(ret) > 0 {
			//validate ret[0]
			if ret[0] == 1 || ret[0] == 2 || ret[0] == 3 {
				logic.nextjiangtype = ret[0]
			}
			//validate ret[1]
			if len(ret) > 1 {
				logic.nextzuhetype = ret[1]
			}
			if len(ret) > 2 {
				logic.nextzuhetype2 = ret[2]
			}
		}

		example.invalid("修改開獎結果成功 %s", jiang)
		return bjl.getRetMsg(true, "修改开奖结果成功")
	}

	return bjl.getRetMsg(false, "修改开奖结果失败")
}

func (bjl *BaiJialeLogic) setZhuShui(tableid int32, zhushui int32) proto.Message {
	if tableid == -1 {
		for _, table := range bjl.GetTable().Room.Tables {
			logic := table.GameLogic.(*BaiJialeLogic)
			logic.changeZhuShui(zhushui)
		}
	} else {
		t := bjl.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*BaiJialeLogic)
			logic.changeZhuShui(zhushui)
		} else {
			return bjl.getRetMsg(false, "修改注水率失败")
		}
	}

	return bjl.getRetMsg(true, "修改注水率成功")
}

func (bjl *BaiJialeLogic) getRealUserWinByType(t int32) int64 {
	win := int64(0)
	odd := int64(bjl.betvalue.getOddByType(t))
	user := bjl.GetUserByUserID(bjl.zhuangID)
	if user != nil && !user.IsRobot() {
		_, tax := bjl.GetUserWinAndTax(bjl.bettotal)
		if bjl.currtax != 0 {
			tax = bjl.currtax
		}
		robotwin := odd * bjl.getRobotBetMoneyByType(t)
		win = bjl.bettotalrobot - robotwin
		win -= tax
		return -win
	}

	for _, ud := range bjl.userdata {
		if !ud.user.IsRobot() {
			win += odd * ud.getUserBetByType(t)
			win += ud.money - ud.user.Score
		}
	}

	mlog.Debug("玩家赢钱%d", win)
	return -win
}

//排序身上攜帶的錢大到小
func (bjl *BaiJialeLogic) getMoneyRanking() {
	sort.Sort(bjl.userscore)
}
func (sbg *userScore) Len() int {
	return len(sbg.userdata)
}
func (sbg *userScore) Swap(i, j int) {
	sbg.userdata[i], sbg.userdata[j] = sbg.userdata[j], sbg.userdata[i]
}
func (sbg *userScore) Less(i, j int) bool {
	return sbg.userdata[i].Score > sbg.userdata[j].Score
}

//獲取幸運星玩家
func (bjl *BaiJialeLogic) getLuckyStar() {
	t := 0
	gameuser := bjl.AllGameUser()
	if bjl.luckystar == 0 || bjl.notbettimes >= 3 {
		//如果沒有幸運星、連續三場沒下注或是幸運星的持續局數到了，隨機一個玩家
		t = 1
	} else {
		//檢查玩家還在不在場
		for _, u := range gameuser {
			if u.UserID == bjl.luckystar {
				bjl.luckystartimes += 1
				break
			} else {
				t = 1
			}
		}
	}
	if t == 1 || bjl.luckystartimes == bjl.luckystarround {
		if len(gameuser) != 0 {
			userIndex := rand.Intn(len(gameuser))
			bjl.luckystar = gameuser[userIndex].UserID
			bjl.luckyfaceid = gameuser[userIndex].FaceID
		} else {
			//如果沒有明燈，須和client確認要傳什麼
			bjl.luckystar = 0
			bjl.luckyfaceid = 0
		}
		bjl.luckystartimes = 0
		bjl.notbettimes = 0
	}
}

//獲取富豪玩家
func (bjl *BaiJialeLogic) getRichestList() {
	userscore := new(userScore)
	userscore.userdata = bjl.AllTableUser()
	sort.Sort(userscore)
	//計算分數最高的7個
	for index, user := range userscore.userdata {
		if index < 7 {
			rich := new(netproto.BJL_Richest)
			rich.Ranking = proto.Int32(int32(index + 1))
			rich.UserID = proto.Int32(user.UserID)
			rich.UserFaceID = proto.Int32(user.FaceID)
			rich.Money = proto.Int64(user.Score)
			bjl.richestList = append(bjl.richestList, rich)
			bjl.richestuserinfo = append(bjl.richestuserinfo, user.UserID)
		} else {
			break
		}
	}
	bjl.stageuserinfo = append(bjl.richestuserinfo, bjl.luckystar)
}

//获取大富豪列表的消息
func (bjl *BaiJialeLogic) getRichestListMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BJL_GameMessageClassID_BJLRichestListID)
	pdata := new(netproto.BJL_RichestList)
	pdata.RichestList = bjl.richestList
	//如果當下沒有玩家 回傳空值
	pmsg.MsgData = pdata
	return pmsg
}

//获取幸運星的消息
func (bjl *BaiJialeLogic) getLuckyStarMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BJL_GameMessageClassID_BJLLuckyStarListID)
	luckyStarList := new(netproto.BJL_LuckyStar)
	luckyStarList.UserID = proto.Int32(bjl.luckystar)
	luckyStarList.UserFaceID = proto.Int32(bjl.luckyfaceid)
	pmsg.MsgData = luckyStarList
	return pmsg
}
func (bjl *BaiJialeLogic) getGameRound() {
	if bjl.roundid == 30 && bjl.zhuangjawin+bjl.xianwin+bjl.hewin == 30 {
		bjl.zhuangjawin = 0
		bjl.xianwin = 0
		bjl.hewin = 0
	}
}
func (bjl *BaiJialeLogic) getRoundMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BJL_GameMessageClassID_BJLRoundID)
	round := new(netproto.BJL_RoundID)
	round.RoundID = proto.Int32(bjl.roundid)
	round.ZhuangWin = proto.Int32(bjl.zhuangjawin)
	round.XianWin = proto.Int32(bjl.xianwin)
	round.HeWin = proto.Int32(bjl.hewin)
	pmsg.MsgData = round
	return pmsg
}

//獲取下注上限
func (bjl *BaiJialeLogic) getBetRangeData() *netproto.BJL_BetRangeList {
	zhuangxianrange := new(netproto.BJL_Range)
	zhuangxianrange.Max = proto.Int32(bjl.zhuangxianrange[3])
	zhuangxianrange.Min = proto.Int32(bjl.zhuangxianrange[2] * bjl.betcount[0])
	herange := new(netproto.BJL_Range)
	herange.Max = proto.Int32(bjl.herange[3])
	herange.Min = proto.Int32(bjl.herange[2] * bjl.betcount[0])
	duiZi := new(netproto.BJL_Range)
	duiZi.Max = proto.Int32(bjl.duizirange[3])
	duiZi.Min = proto.Int32(bjl.duizirange[2] * bjl.betcount[0])
	rdata := new(netproto.BJL_BetRangeList)
	rdata.ZhuangXian = zhuangxianrange
	rdata.HeRange = herange
	rdata.DuiZi = duiZi

	return rdata
}

//收到RPC消息解析處理
func (bjl *BaiJialeLogic) OnRecvRPCMessage(m *msg.Message) {
	rm, ok := m.MsgData.(*rpcmsg.RPCMessage)
	if !ok {
		return
	}
	mmbyte := rm.MsgData.([]byte)
	if len(mmbyte) == 0 {
		return
	}

	if m.SClassID == int32(netproto.GameMessageClassID_SuperToolGetAssignRewardDetailRetID) {
		bjl.setAssignReward(m)
	}

}

func (bjl *BaiJialeLogic) getFangFenKaiJiangType() []int32 {
	//放分開和條件及機率
	r := (bjl.fangFenTarget * 15 / 10) - bjl.fangFenComplete
	hepay := bjl.getBetMoneyByType(HE) * int64(bjl.betvalue.getOddByType(bjl.jiangtype[0]))
	if r > hepay {
		rnd := util.RndInt(0, 1000)
		if rnd < 75 {
			return []int32{HE}
		}
	}

	if bjl.getBetMoneyByType(ZHUANG) > bjl.getBetMoneyByType(XIAN) {
		return []int32{ZHUANG}
	} else if bjl.getBetMoneyByType(ZHUANG) < bjl.getBetMoneyByType(XIAN) {
		return []int32{XIAN}
	}

	return []int32{ZHUANG, XIAN}
}

//获取此类型的下注总额 真人
func (bjl *BaiJialeLogic) getRealManBetMoneyByType(t int32) int64 {
	allBetTotal := int64(0)
	for i := 0; i < len(bjl.bettypetotal); i++ {
		if bjl.bettypetotal[i].baijialetype == t {
			allBetTotal += bjl.bettypetotal[i].money
		}
	}
	return allBetTotal
}

//獲取放分組合類型
func (bjl *BaiJialeLogic) getFangFenZuheType() (zuhetype []int32) {
	alltype := []int32{ZHUANGDUI, XIANDUI}
	rnd := util.RndInt(0, 1000)
	if rnd < 165 {
		r := (bjl.fangFenTarget * 15 / 10) - bjl.fangFenComplete
		for _, zuhe := range alltype {
			zuhepay := bjl.getRealManBetMoneyByType(3) * int64(bjl.betvalue.getOddByType(zuhe))
			if r > zuhepay {
				zuhetype = append(zuhetype, zuhe)
			}
		}
	}
	//if len(zuhetype) < 1 {
	//	zuhetype = append(alltype, []int32{NIUNIU, NIUNINE, NIUEIGHT, NIUSEVEN, NIUSIX, NIUFIVE, NIUFOUR, NIUTHREE, NIUTWO, NIUONE, WUNIU}...)
	//}
	return
}

// 開始下注前發牌動畫
func (bjl *BaiJialeLogic) dealBeforeBet() {
	bjl.setState(DealAction)
	bjl.BoardcastMessage(bjl.getStateMsg())
	bjl.actionStr = bjl.RegisterAnimation(GameDealAction, bjl.DealActionTime, bjl.startBet)
}
