package hhnn

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/conf"
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

const (
	Bet            = iota //下注状态
	KaiJiang              //开奖状态
	KaiJiangAction        //开奖动画
	EndAction             //结算动画
	Wait                  //等待游戏开始
)

const (
	HUNG      = 1  //紅
	HEI       = 2  //黑
	NATURE    = 3  //自然機率
	WUXIAONIU = 4  //五小牛
	ZHADAN    = 5  //炸彈
	WUHUANIU  = 6  //五花牛
	KNIUNIU   = 7  //含K牛牛
	NIUNIU    = 8  //不含K牛牛
	NIUNINE   = 9  //牛9
	NIUEIGHT  = 10 //牛8
	NIUSEVEN  = 11 //牛7
	NIUSIX    = 12 //牛6
	NIUFIVE   = 13 //牛5
	NIUFOUR   = 14 //牛4
	NIUTHREE  = 15 //牛3
	NIUTWO    = 16 //牛2
	NIUONE    = 17 //牛1
	WUNIU     = 18 //無牛
)

//第一个用来占位
var resultstring = []string{"", "红", "黑", "好运连连"}

const (
	GameStartAnimation  = "gamestart"
	GameKaiJiang_Action = "kaijiang"
	GameEnd_Action      = "end"
)

const (
	CtrlBetLimitDefault int64 = 100 // 監控名單預設注額
)

type JiangTypeWinLossOdd struct {
	betType int32 //下注種類選項
	winLoss int32 //此項和莊家比牌後輸贏結果　-1輸　1贏 , 負數是賠付倍數方便計算
	odds    int32 //賠率, 贏的參考閒家的牛點數當賠率　輸參考莊的牌賠率
}

type userScore struct {
	userdata []*gamelogic.GameUser
}
type HHNNLogic struct {
	*gamelogic.ComputeGameLogicBase
	serverfile         *os.File
	opttimer           *event.TimerFunc
	bettotal           int64                         //下注总额
	bettotalrobot      int64                         //机器人下注总额
	bettypetotal       []*bettype                    //下注分类汇总
	bettyperobot       []*bettype                    //机器人下注分类汇总
	betvalue           *betconfig                    //下注配置文件
	userdata           []*hhnnuserdata               //用户数据
	gamestate          int32                         //游戏状态
	bettime            int32                         //下注时间
	jiangtype          int32                         //开奖类型
	action_str         string                        //当前动画
	kaijiangactiontime int32                         //开奖动画超时
	endactiontime      int32                         //用户结算动画超时
	readytime          int32                         //用户准备时间
	zhuangID           int32                         //当前庄家ID
	zhuanglimit        int64                         //上庄限制
	lastturn           int32                         //庄剩余的局数
	maxturn            int32                         //上庄的最大局数
	zhuanglist         []int32                       //上庄列表
	betcount           []int32                       //下注数量
	syspool            int64                         //当前水池的钱
	zhuangWin          int64                         //庄家赢钱
	serverid           string                        //服务器ID
	winuserinfo        []int32                       //赢钱最多的前三家userid
	resultcardlist     *netproto.HHNN_ResultCardList //开奖的历史牌
	zhushui            int32                         //注水率
	maxclearcount      int32                         //清理数据上限
	maxrobotzhuang     int32                         //最多有好多个机器人上庄
	robotxiazhuang     int32                         //机器人下庄概率
	robotzhuang        int32                         //没有人当庄的时候机器人当庄的概率
	zhushuimoney       int64                         //注水的钱
	he                 int32                         //下和的概率百分比
	startsyspool       int64                         //开始时的金币数量
	startmoney         int64                         //上庄时的金币
	currtax            int64                         //当前抽水
	nextjiangtype      int32                         //下次开奖类型
	nextzuhetype       int32                         //下次開獎牌型
	systotalwin        int64                         //系统收入
	syslastwin         int64                         //上一局收益
	sametimebetmode    int64                         //是否開啟庄閒同時下注
	systaxincome       int64                         //系統稅收, 非該局稅收
	userwinlose        int64                         //玩家輸贏
	robotwinlose       int64                         //機器人輸贏
	betlimitmoney      int64                         //下注限制
	publiczhuang       int64                         //是否開放用戶上庄
	userscore          *userScore                    //用户数据
	luckystar          int32                         //幸運星玩家id
	luckyfaceid        int32                         //幸運星玩家頭像id
	luckystarround     int64                         //幸運星持續局數
	luckystartimes     int64                         //幸運星第幾局
	notbettimes        int64                         //沒下注持續局數
	richestuserinfo    []int32                       //富豪玩家的userid
	richestList        []*netproto.HHNN_Richest      //富豪榜
	stageuserinfo      []int32                       //標誌玩家的userid
	roundid            int32                         //當局局數
	hungwin            int32                         //紅贏局數
	heiwin             int32                         //黑贏局數
	poker              poker.Poker                   //整副牌
	cardhung           *poker.CardGroup              //开奖的紅區牌
	cardhei            *poker.CardGroup              //开奖的黑區牌
	wintype            int32                         //开奖贏的组合
	losetype           int32                         //开奖輸的组合
	weight             []*betconfig                  //权重配置文件
	scoreSection       []int64                       //分数段用来加载不同的分数段配置
	winline            int64                         //赢钱线大于这个值正常開。小于这个值開最小獎
	//firstPoker          map[int]*poker.CardGroup      //第一手牌:紅*1和黑*1
	//secondPokerGroup    map[int]*poker.CardGroup      //第二手牌
	jiangTypeWinLossOdd []*JiangTypeWinLossOdd //賠率, 贏的參考閒家的牛點數當賠率　輸參考莊的牌賠率
	zuhetype            int32                  //开奖贏的组合
	criticalwinline     int64                  //好運連連專屬Winline線
	paytable            []*oddconfig           //权重配置文件
	lastGameUserBet     []*hhnnuserdata        //上局下注玩家
	//bowProvider         *gamelogic.BlackOrWhiteUserProvider   //需要放分或殺分的列表
	assignRewardDetail *netproto.SuperToolAssignRewardDetail //指定開獎
	execAssignReward   bool                                  //true 執行指定開獎
	fangFenTarget      int64                                 //放分金額
	fangFenComplete    int64                                 //己放分金額
	execFangFen        bool                                  //true 執行系統放分
	SuperToolService   *gamelogic.SuperToolBrGameServiceBase //超端服務
	ctrlbetlimit       int64                                 //監控名單會控制的最低注額
}

//新建一个紅黑牛牛实例
func NewHHNNLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *HHNNLogic {
	glc := new(HHNNLogic)
	glc.ComputeGameLogicBase = gamelogic.NewComputeGameLogicBase(t, tfManager, cfgs)
	glc.getHHNNconfig(cfgs)
	glc.SuperToolService = glc.GetSuperToolService().(*gamelogic.SuperToolBrGameServiceBase)
	glc.gamestate = Wait
	glc.zhuangID = 0
	glc.he = 0
	glc.nextjiangtype = 0
	glc.syslastwin = 0
	glc.systotalwin = 0
	glc.roundid = 0
	glc.hungwin = 0
	glc.heiwin = 0
	glc.resultcardlist = new(netproto.HHNN_ResultCardList)
	glc.getLastPool()
	glc.startsyspool = glc.syspool
	glc.addCmd()
	glc.RegisterAfterGameStartFunc(func() {
		glc.clearGameData()   //清除数据
		glc.getAssignReward() //獲取多局指定開獎
		//glc.VipUserService.UpdateVIPUser() //獲取最新的監控用戶名單
		//glc.getTableBlackUser()
		//glc.getTableWhiteUser()
		glc.fapai()
		glc.startFirstCard()
		glc.huangZhuang() //6/12 移除上莊機制
		glc.getLuckyStar()
		glc.getRichestList()
		glc.getGameRound()
		glc.GameLog.AddGameStep("庄家%d，剩余%d", glc.zhuangID, glc.lastturn)
		glc.zhuShui()
		glc.broadCastZhuangInfo()
		glc.startAnimation()
	})

	glc.RegisterBeforeGameEndFunc(func() {
		glc.saveLastGameUserBet() //上局投注用戶記錄
		//glc.updateBowUserWin()    //更新監控殺放分額度
	})

	return glc
}

//廣播第一手三張牌，並开始下注，下注時間結束才發第二手兩張牌並接著開獎
func (hhnn *HHNNLogic) startFirstCard() {
	//根據預設開獎類型進行組牌及留牌
	//hhnn.fapeiFirstCardGroup()
	//廣播第一手三張牌
	hhnn.boardcastFirstCardGroup()
}

//廣播第一手牌
func (hhnn *HHNNLogic) boardcastFirstCardGroup() {
	hhnn.BoardcastMessage(hhnn.getFirstCardGroup())
}

//
////第一手三張牌, 發三張空牌給客戶端
func (hhnn *HHNNLogic) getFirstCardGroup() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNBroadCastFapaiThreeCard)

	pdata := new(netproto.HHNN_BrodcastPokerCard)

	card := new(poker.Card)
	card.CardNumber = 0
	card.Color = 0

	cg1 := new(poker.CardGroup)
	for i := 0; i < 3; i++ {
		cg1.Cards = append(cg1.Cards, card)
	}

	for i := int32(0); i <= int32(1); i++ {
		urPokerCard := new(netproto.HHNN_BrodcastUserPokerCard)
		urPokerCard.SitNo = proto.Int32(i + 1)
		for _, cd := range cg1.Cards {
			card := new(netproto.PokerCard)
			card.Color = proto.Int32(cd.Color)
			card.CardNumber = proto.Int32(cd.CardNumber)
			urPokerCard.CardGroup = append(urPokerCard.CardGroup, card)
		}
		pdata.UserPokerCard = append(pdata.UserPokerCard, urPokerCard)
	}

	pmsg.MsgData = pdata
	return pmsg
}

//取得第一手三張牌的指定位置牌組
//func (hhnn *HHNNLogic) getCalFirstPokerGroup(position int32) *poker.CardGroup {
//	return hhnn.firstPoker[int(position)]
//}

//產生第一手三張牌
//func (hhnn *HHNNLogic) fapeiFirstCardGroup() {
//
//}

func (hhnn *HHNNLogic) fapai() {
	//随机发牌
	if hhnn.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
		hhnn.poker = poker.NewRandomPoker(false, false, 1) //初始化一副牌
		hhnn.poker.Shuffle()
	} else { //发牌器发牌
		hhnn.poker = poker.NewFilePoker(hhnn.Configs, hhnn.AllTableUser())
		hhnn.poker = poker.NewFixedFilePoker(hhnn.poker.(*poker.FilePoker))
		prepare := poker.NewRandomPoker(false, false, 1) //初始化一副牌
		prepare.Shuffle()
	}
}

//获取配置表
func (hhnn *HHNNLogic) getHHNNconfig(cfgs *config.Group) {
	hhnn.betvalue = newBetconfig(cfgs, "hhnn")
	scfgs := cfgs.GetConfig("gameserver")
	hhnn.bettime = int32(scfgs.GetConfigInt("gameserver", "bettime"))
	hhnn.kaijiangactiontime = int32(scfgs.GetConfigInt("gameserver", "kaijiangactiontime"))
	hhnn.endactiontime = int32(scfgs.GetConfigInt("gameserver", "endactiontime"))
	hhnn.readytime = int32(scfgs.GetConfigInt("gameserver", "readywaittime"))
	hhnn.zhuanglimit = int64(scfgs.GetConfigInt("gameserver", "zhuanglimit"))
	hhnn.maxturn = int32(scfgs.GetConfigInt("gameserver", "zhuangturn"))
	hhnn.serverid = strconv.Itoa(int(hhnn.GetServerConfig().ServerID))
	hhnn.maxclearcount = int32(scfgs.GetConfigInt("gameserver", "maxclearcount"))
	hhnn.zhushui = int32(scfgs.GetConfigInt("gameserver", "zhushui"))
	hhnn.maxrobotzhuang = int32(scfgs.GetConfigInt("gameserver", "maxrobotzhuang"))
	hhnn.robotxiazhuang = int32(scfgs.GetConfigInt("gameserver", "robotxiazhuang"))
	hhnn.robotzhuang = int32(scfgs.GetConfigInt("gameserver", "robotzhuang"))
	hhnn.winline = int64(scfgs.GetConfigInt("gameserver", "winline"))
	hhnn.sametimebetmode = int64(scfgs.GetConfigInt("gameserver", "sametimebetmode"))
	hhnn.betlimitmoney = int64(scfgs.GetConfigInt("gameserver", "betlimitmoney"))
	hhnn.publiczhuang = int64(scfgs.GetConfigInt("gameserver", "publiczhuang"))
	hhnn.luckystarround = int64(scfgs.GetConfigInt("gameserver", "luckystarround"))
	hhnn.criticalwinline = int64(scfgs.GetConfigInt("gameserver", "criticalwinline"))
	betcountstr := scfgs.GetConfigValue("gameserver", "betcount")
	lastindex := 0
	for index, char := range betcountstr {
		if char == ',' {
			substr := betcountstr[lastindex:index]
			lastindex = index + 1
			count, _ := strconv.ParseInt(substr, 10, 32)
			hhnn.betcount = append(hhnn.betcount, int32(count))
		}
	}
	scorestr := scfgs.GetConfigValue("gameserver", "scoresection")
	lastindex = 0
	for index, char := range scorestr {
		if char == ',' {
			substr := scorestr[lastindex:index]
			lastindex = index + 1
			score, _ := strconv.ParseInt(substr, 10, 64)
			hhnn.scoreSection = append(hhnn.scoreSection, score)
			w := newBetconfig(cfgs, substr)
			hhnn.weight = append(hhnn.weight, w)
		}
	}
	p := newOddconfig(cfgs, "paytable")
	hhnn.paytable = append(hhnn.paytable, p)

	hhnn.ctrlbetlimit = int64(scfgs.GetConfigInt("supertool", "ctrlbetlimit"))
	if hhnn.ctrlbetlimit <= 0 {
		hhnn.ctrlbetlimit = CtrlBetLimitDefault
	}
}

//添加命令行
func (hhnn *HHNNLogic) addCmd() {
	command.AddCommand(NewHHNNCmd(hhnn))
}

//广播庄家、大富豪、幸運星
func (hhnn *HHNNLogic) broadCastZhuangInfo() {
	hhnn.BoardcastMessage(hhnn.getZhuangInfoMsg())
	hhnn.BoardcastMessage(hhnn.getRichestListMsg())
	hhnn.BoardcastMessage(hhnn.getLuckyStarMsg())
}

func (hhnn *HHNNLogic) startAnimation() {
	hhnn.setState(Bet)
	//hhnn.startBet()
	//開始下注前+2秒
	hhnn.RegisterAnimation(GameStartAnimation, int32(2), hhnn.startBet)
}

//获取上一次水池的钱
func (hhnn *HHNNLogic) getLastPool() {
	if conf.TestMode {
		return
	}

	var folderName = "lastpool"
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err = os.MkdirAll(folderName, os.ModePerm)
	}

	var tablenostr = strconv.Itoa(int(hhnn.TableNo()))
	var filename = hhnn.serverid + "_" + tablenostr + ".conf"
	filename = path.Join(folderName, filename)

	_, err := os.Stat(filename)
	var f *os.File = nil
	if os.IsNotExist(err) {
		f, err = os.Create(filename)
	} else {
		f, err = os.OpenFile(filename, os.O_RDWR, 0)
	}

	if err != nil {
		mlog.Error("%v文件不存在", filename)
		panic(err)
	} else {
		hhnn.serverfile = f
		defer hhnn.serverfile.Close()
		hhnn.readPoolFromFile()
	}
}

//读取水池文件
func (hhnn *HHNNLogic) readPoolFromFile() {
	if conf.TestMode {
		return
	}

	buff := make([]byte, 1024)
	n, _ := hhnn.serverfile.Read(buff)
	if n > 0 {
		poolstr := string(buff)[0:n]
		//pool, err := strconv.ParseInt(poolstr, 10, 32)
		pool, _ := strconv.Atoi(poolstr)
		hhnn.syspool = int64(pool)
	}
}

//注水到下注额
func (hhnn *HHNNLogic) zhuShui() {
	if hhnn.zhuangID == 0 {
		hhnn.bettotal = hhnn.syspool * int64(hhnn.zhushui) / 100
		hhnn.syspool -= hhnn.bettotal
	} else {
		ur := hhnn.GetUserByUserID(hhnn.zhuangID)
		if ur != nil && ur.IsRobot() {
			hhnn.bettotal = hhnn.syspool * int64(hhnn.zhushui) / 100
			hhnn.syspool -= hhnn.bettotal
		}
	}
	hhnn.zhushuimoney = hhnn.bettotal
	hhnn.GameLog.AddGameStep("系統水池%d, 系統注水%d", hhnn.syspool+hhnn.zhushuimoney, hhnn.zhushuimoney)
}

//改变水池注水率
func (hhnn *HHNNLogic) changeZhuShui(value int32) {
	if value >= 0 && value <= 100 {
		hhnn.zhushui = value
	}
}
func (hhnn *HHNNLogic) checkSysPool() {
	mlog.Printf("桌子号：%v，系统水池：%v，注水概率：%v", hhnn.TableNo(), hhnn.syspool, hhnn.zhushui)
}

//改变水池的钱
func (hhnn *HHNNLogic) changeSysPool(value int64) {
	hhnn.startsyspool += value - hhnn.syspool
	hhnn.syspool = value
}

//改变机器人的和的下注概率
func (hhnn *HHNNLogic) changeHe(value int32) {
	if value >= 0 && value <= 100 {
		hhnn.he = value
	}
}

//清除数据
func (hhnn *HHNNLogic) clearGameData() {
	hhnn.userdata = make([]*hhnnuserdata, 0)
	hhnn.bettypetotal = make([]*bettype, 0)
	hhnn.bettyperobot = make([]*bettype, 0)
	hhnn.winuserinfo = make([]int32, 0)
	hhnn.stageuserinfo = make([]int32, 0)
	hhnn.richestuserinfo = make([]int32, 0)
	hhnn.cardhung = new(poker.CardGroup)
	hhnn.cardhei = new(poker.CardGroup)
	hhnn.richestList = make([]*netproto.HHNN_Richest, 0)
	hhnn.bettotal = 0
	hhnn.bettotalrobot = 0
	hhnn.zhuangWin = 0
	hhnn.currtax = 0
	hhnn.robotwinlose = 0
	hhnn.userwinlose = 0
	hhnn.zuhetype = 0
	hhnn.jiangtype = 0
	hhnn.assignRewardDetail = nil
	hhnn.execAssignReward = false
	hhnn.execFangFen = false
	//hhnn.firstPoker = make(map[int]*poker.CardGroup)
	//hhnn.secondPokerGroup = make(map[int]*poker.CardGroup)
	hhnn.jiangTypeWinLossOdd = make([]*JiangTypeWinLossOdd, 0)
	rand.Seed(time.Now().UnixNano())
	zlen := len(hhnn.zhuanglist)
	if zlen == 0 {
		hhnn.robotZhuang()
	}

	hhnn.robotXiaZhuang()
}

//获取庄家信息消息
func (hhnn *HHNNLogic) getZhuangInfoMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNZhuangInfoID)

	pdata := new(netproto.HHNN_ZhuangInfo)
	pdata.UserID = proto.Int32(0)
	pdata.LastTurn = proto.Int32(0)
	ur := hhnn.GetUserByUserID(hhnn.zhuangID)
	if ur != nil {
		pdata.UserInfo = append(pdata.UserInfo, ur.GetProtoBufGameUser())
		pdata.UserID = proto.Int32(ur.UserID)
		pdata.LastTurn = proto.Int32(hhnn.lastturn)
	}
	pmsg.MsgData = pdata
	return pmsg
}

//换庄
func (hhnn *HHNNLogic) huangZhuang() {
	if hhnn.publiczhuang == 1 {
		if hhnn.zhuangID == 0 {
			//庄是系统直接换庄
			hhnn.setNextZhuang()
			return
		}
		if hhnn.lastturn <= 0 {
			//已经到最后一局了
			hhnn.setNextZhuang()
			return
		}
		ur := hhnn.GetUserByUserID(hhnn.zhuangID)
		if ur != nil && ur.Score < hhnn.zhuanglimit {
			hhnn.setNextZhuang()
			return
		}

		if ur == nil {
			hhnn.setNextZhuang()
			return
		}

		hhnn.lastturn--
		mlog.Debug("1015发送庄家剩余局数%d", hhnn.lastturn)
	} else {
		mlog.Debug("未開啟上庄機制，不用換庄")
		return
	}
}

//设置下一个上庄的玩家
func (hhnn *HHNNLogic) setNextZhuang() {
	if hhnn.publiczhuang == 1 {
		hhnn.checkNetUserCanShangZhuang()
		if len(hhnn.zhuanglist) > 0 {
			hhnn.zhuangID = hhnn.zhuanglist[0]
			hhnn.lastturn = hhnn.maxturn
			hhnn.lastturn--
			hhnn.delFromZhuangList(hhnn.zhuangID)
			user := hhnn.GetUserByUserID(hhnn.zhuangID)
			hhnn.startmoney = user.Score
		} else {
			hhnn.zhuangID = 0
			hhnn.startmoney = hhnn.startsyspool
		}
	} else {
		mlog.Error("未開啟上庄機制，不得上庄")
		return
	}
}

//检测下一个用户是否能上庄
func (hhnn *HHNNLogic) checkNetUserCanShangZhuang() {
	if len(hhnn.zhuanglist) > 0 {
		for {
			user := hhnn.GetUserByUserID(hhnn.zhuanglist[0])
			if user != nil && user.Score >= int64(hhnn.zhuanglimit) {
				break
			}

			hhnn.delFromZhuangList(hhnn.zhuanglist[0])
			if len(hhnn.zhuanglist) <= 0 {
				break
			}
		}
	}
}

//删除一个用户从上庄列表
func (hhnn *HHNNLogic) delFromZhuangList(userid int32) {
	for i := 0; i < len(hhnn.zhuanglist); i++ {
		if hhnn.zhuanglist[i] == userid {
			hhnn.zhuanglist = append(hhnn.zhuanglist[:i], hhnn.zhuanglist[i+1:]...)
			break
		}
	}
}

//机器人上庄
func (hhnn *HHNNLogic) robotZhuang() {
	if hhnn.publiczhuang == 1 {
		//概率当庄
		if rand.Int31n(100) > hhnn.robotzhuang {
			return
		}
		//统计机器人个数
		robotarr := make([]int32, 0)
		for _, ur := range hhnn.AllTableUser() {
			if ur.IsRobot() {
				robotarr = append(robotarr, ur.UserID)
			}
		}

		//沒有機器人或沒有設定最多幾個機器人上庄
		if len(robotarr) <= 0 || hhnn.maxrobotzhuang <= 0 {
			return
		}

		//随机有几个人当庄
		n := rand.Int31n(hhnn.maxrobotzhuang) + 1

		for i := int32(0); i < n; i++ {
			rcn := len(robotarr)
			if rcn <= 0 {
				break
			}
			index := rand.Int31n(int32(rcn))
			hhnn.addToZhuangList(robotarr[index])
			robotarr = append(robotarr[:index], robotarr[index+1:]...)
		}
	} else {
		mlog.Debug("未開啟上庄機制，機器人不用上庄")
		return
	}
}

//添加一个用户到上庄列表
func (hhnn *HHNNLogic) addToZhuangList(userid int32) {
	hhnn.zhuanglist = append(hhnn.zhuanglist, userid)
}

//机器人下庄
func (hhnn *HHNNLogic) robotXiaZhuang() {
	if hhnn.zhuangID == 0 {
		return
	}

	ur := hhnn.GetUserByUserID(hhnn.zhuangID)
	if ur != nil && ur.IsRobot() {
		if rand.Int31n(100) < hhnn.robotxiazhuang {
			hhnn.zhuangID = 0
		}
	}
}

//排序身上攜帶的錢大到小
func (hhnn *HHNNLogic) getMoneyRanking() {
	sort.Sort(hhnn.userscore)
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
func (hhnn *HHNNLogic) getLuckyStar() {
	t := 0
	gameuser := hhnn.AllGameUser()
	if hhnn.luckystar == 0 || hhnn.notbettimes >= 3 {
		//如果沒有幸運星、連續三場沒下注或是幸運星的持續局數到了，隨機一個玩家
		t = 1
	} else {
		//檢查玩家還在不在場
		for _, u := range gameuser {
			if u.UserID == hhnn.luckystar {
				hhnn.luckystartimes += 1
				break
			} else {
				t = 1
			}
		}
	}
	if t == 1 || hhnn.luckystartimes == hhnn.luckystarround {
		if len(gameuser) != 0 {
			userIndex := rand.Intn(len(gameuser))
			hhnn.luckystar = gameuser[userIndex].UserID
			hhnn.luckyfaceid = gameuser[userIndex].FaceID
		} else {
			//如果沒有明燈，須和client確認要傳什麼
			hhnn.luckystar = 0
			hhnn.luckyfaceid = 0
		}
		hhnn.luckystartimes = 0
		hhnn.notbettimes = 0
	}
}

//獲取富豪玩家
func (hhnn *HHNNLogic) getRichestList() {
	userscore := new(userScore)
	userscore.userdata = hhnn.AllTableUser()
	sort.Sort(userscore)
	//計算分數最高的7個
	for index, user := range userscore.userdata {
		if index < 5 {
			rich := new(netproto.HHNN_Richest)
			rich.Ranking = proto.Int32(int32(index + 1))
			rich.UserID = proto.Int32(user.UserID)
			rich.UserFaceID = proto.Int32(user.FaceID)
			rich.Money = proto.Int64(user.Score)
			hhnn.richestList = append(hhnn.richestList, rich)
			hhnn.richestuserinfo = append(hhnn.richestuserinfo, user.UserID)
		} else {
			break
		}
	}
	hhnn.stageuserinfo = append(hhnn.richestuserinfo, hhnn.luckystar)
}
func (hhnn *HHNNLogic) getGameRound() {
	if hhnn.roundid == 30 {
		hhnn.hungwin = 0
		hhnn.heiwin = 0
	}
}

//获取大富豪列表的消息
func (hhnn *HHNNLogic) getRichestListMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNRichestListID)
	pdata := new(netproto.HHNN_RichestList)
	pdata.RichestList = hhnn.richestList
	//如果當下沒有玩家 回傳空值
	pmsg.MsgData = pdata
	return pmsg
}

//获取幸運星的消息
func (hhnn *HHNNLogic) getLuckyStarMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNLuckyStarListID)
	luckyStarList := new(netproto.HHNN_LuckyStar)
	luckyStarList.UserID = proto.Int32(hhnn.luckystar)
	luckyStarList.UserFaceID = proto.Int32(hhnn.luckyfaceid)
	pmsg.MsgData = luckyStarList
	return pmsg
}

func (hhnn *HHNNLogic) getRoundMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNRoundID)
	round := new(netproto.HHNN_RoundID)
	round.RoundID = proto.Int32(hhnn.roundid)
	round.HungWin = proto.Int32(hhnn.hungwin)
	round.HeiWin = proto.Int32(hhnn.heiwin)
	pmsg.MsgData = round
	return pmsg
}

//用户离开
func (hhnn *HHNNLogic) UserLeave(gameuser *gamelogic.GameUser) bool {
	if hhnn.gamestate != Wait {
		if hhnn.getUserDataByID(gameuser.UserID) != nil {
			mlog.Debug("用户在玩游戏还不能离开")
			return false
		}

		if hhnn.zhuangID == gameuser.UserID {
			mlog.Debug("用户在当庄，所以用户不能离开")
			return false
		}
	}

	hhnn.delFromZhuangList(gameuser.UserID)

	hhnn.ChangeUserStateEnd(gameuser)
	return hhnn.ComputeGameLogicBase.UserLeave(gameuser)
}

//注销用户实例
func (hhnn *HHNNLogic) UserFree(gameuser *gamelogic.GameUser) {
	if hhnn.gamestate != Wait {
		if hhnn.getUserDataByID(gameuser.UserID) != nil {
			mlog.Debug("用户在玩游戏还不能离开")
			return
		}

		if hhnn.zhuangID == gameuser.UserID {
			mlog.Debug("用户在当庄，所以用户不能离开")
			return
		}
	} else {
		if hhnn.zhuangID == gameuser.UserID {
			hhnn.zhuangID = 0
		}
	}

	hhnn.delFromZhuangList(gameuser.UserID)

	hhnn.ChangeUserStateEnd(gameuser)
	hhnn.ComputeGameLogicBase.UserFree(gameuser)
}

//断线重连
func (hhnn *HHNNLogic) UserRecover(gameuser *gamelogic.GameUser) {
	mlog.Debug("UserRecover")
	hhnn.ComputeGameLogicBase.UserRecover(gameuser)
	ud := hhnn.getUserDataByID(gameuser.UserID)
	if ud != nil {
		ud.money += ud.user.Score - ud.moneylast
		ud.moneylast = ud.user.Score
	}

	hhnn.sendScene(gameuser)
}

//用户进入
func (hhnn *HHNNLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {
	hhnn.ComputeGameLogicBase.UserEnter(gameuser, tableno, seatno)
	mlog.Debug("UserEnter")
	hhnn.sendScene(gameuser)
}

//发送场景
func (hhnn *HHNNLogic) sendScene(gameuser *gamelogic.GameUser) {
	//发送设置

	if gameuser.IsRobot() {
		mlog.Debug("给机器人%d发送场景消息", gameuser.UserID)
	} else {
		mlog.Debug("给玩家%d发送场景消息", gameuser.UserID)
	}
	hhnn.sendSetting(gameuser)
	hhnn.SendMessageToUser(gameuser, hhnn.getRoundMsg())
	hhnn.SendMessageToUser(gameuser, hhnn.getResultCardListMsg())
	hhnn.SendMessageToUser(gameuser, hhnn.getRichestListMsg())
	hhnn.SendMessageToUser(gameuser, hhnn.getLuckyStarMsg())
	if hhnn.gamestate == Wait {
		hhnn.BoardcastMessage(hhnn.getStateMsg())
		return
	}
	//发送阶段
	hhnn.SendMessageToUser(gameuser, hhnn.getStateMsg())
	hhnn.SendMessageToUser(gameuser, hhnn.getZhuangInfoMsg())
	hhnn.SendMessageToUser(gameuser, hhnn.getBetTotalMsg(0))
	hhnn.sendUserBetInfo(gameuser)

	if hhnn.gamestate == Bet {
		hhnn.sendToUserFirstCardGroup(gameuser)
	}

	if hhnn.gamestate == KaiJiangAction {
		hhnn.BoardcastMessage(hhnn.getStateMsg())
		//发送开奖信息
		hhnn.SendMessageToUser(gameuser, hhnn.getGameResultPokerMsg())
	} else if hhnn.gamestate == EndAction {
		hhnn.BoardcastMessage(hhnn.getStateMsg())
		//发送用户下注输赢信息如果有
		ud := hhnn.getUserDataByID(gameuser.UserID)
		if ud != nil {
			hhnn.SendMessageToUser(gameuser, hhnn.getUserWinMsg(ud.win, hhnn.zhuangWin))
		}
	}
}

//发送设置
func (hhnn *HHNNLogic) sendSetting(gameuser *gamelogic.GameUser) {
	mlog.Debug("sendSetting")
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNSettingID)

	pdata := new(netproto.HHNN_Setting)
	pdata.BetCount = append(pdata.BetCount, hhnn.betcount...)
	pdata.BetMoney = proto.Int32(int32(hhnn.betvalue.getMoneyByType(1)))
	pdata.ZhuangLimit = proto.Int64(hhnn.zhuanglimit)
	pdata.Betlimitmoney = proto.Int64(hhnn.betlimitmoney)
	pmsg.MsgData = pdata
	hhnn.SendMessageToUser(gameuser, pmsg)
}

//获取用户数据通过用户ID
func (hhnn *HHNNLogic) getUserDataByID(userid int32) *hhnnuserdata {
	for _, ud := range hhnn.userdata {
		if ud.user.UserID == userid {
			return ud
		}
	}

	return nil
}

//获取开奖历史列表的消息
func (hhnn *HHNNLogic) getResultCardListMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNResultCardListID)

	pdata := new(netproto.HHNN_ResultCardList)
	if len(hhnn.resultcardlist.CardList) > 30 { //大於 30 取最後三十個
		pdata.CardList = append(pdata.CardList, hhnn.resultcardlist.CardList[len(hhnn.resultcardlist.CardList)-30:]...)
	} else {
		pdata.CardList = append(pdata.CardList, hhnn.resultcardlist.CardList...)
	}

	mlog.Debug("getResultCardListMsg:%v", pdata.CardList)
	pmsg.MsgData = pdata
	return pmsg
}

//获取状态消息
func (hhnn *HHNNLogic) getStateMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNStateID)

	pdata := new(netproto.HHNN_Sate)
	pdata.State = proto.Int32(hhnn.gamestate)
	pdata.Time = proto.Int32(hhnn.getBetTime())
	pmsg.MsgData = pdata
	return pmsg
}

//下注時間
func (hhnn *HHNNLogic) getBetTime() int32 {
	time := int32(0)
	if hhnn.opttimer != nil {
		time = hhnn.opttimer.WaitSeconds - hhnn.opttimer.ElapsedSeconds()
	}
	return time
}

//获取下注总额消息
func (hhnn *HHNNLogic) getBetTotalMsg(userid int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNBroadCastBetID)

	pdata := new(netproto.HHNN_BroadCastBetArr)
	pdata.UserID = proto.Int32(userid)

	lucky := hhnn.getUserDataByID(hhnn.luckystar) //獲取明燈

	for _, v := range hhnn.betvalue.alltype {
		betinfo := new(netproto.HHNN_BroadCastBet)
		betinfo.Type = proto.Int32(v)
		betinfo.Money = proto.Int64(hhnn.getBetMoneyByType(v))
		betinfo.LuckyStar = proto.Bool(false)
		if lucky != nil && lucky.getUserBetByType(v) > 0 {
			betinfo.LuckyStar = proto.Bool(true)
		}
		pdata.BroadCastBetArr = append(pdata.BroadCastBetArr, betinfo)
	}
	pmsg.MsgData = pdata

	return pmsg
}

//获取此类型的下注总额 真人+機器人
func (hhnn *HHNNLogic) getBetMoneyByType(t int32) int64 {
	allBetTotal := int64(0)
	for i := 0; i < len(hhnn.bettypetotal); i++ {
		if hhnn.bettypetotal[i].hhnntype == t {
			allBetTotal += hhnn.bettypetotal[i].money
		}
	}
	for i := 0; i < len(hhnn.bettyperobot); i++ {
		if hhnn.bettyperobot[i].hhnntype == t {
			allBetTotal += hhnn.bettyperobot[i].money
		}
	}
	return allBetTotal
}

//获取此类型的下注总额 真人
func (hhnn *HHNNLogic) getRealManBetMoneyByType(t int32) int64 {
	allBetTotal := int64(0)
	for i := 0; i < len(hhnn.bettypetotal); i++ {
		if hhnn.bettypetotal[i].hhnntype == t {
			allBetTotal += hhnn.bettypetotal[i].money
		}
	}
	return allBetTotal
}

//发送用户的下注信息
func (hhnn *HHNNLogic) sendUserBetInfo(gameuser *gamelogic.GameUser) {
	ud := hhnn.getUserDataByID(gameuser.UserID)
	if ud == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNBetInfoID)

	pdata := new(netproto.HHNN_UserBetInfo)
	pdata.UserMoney = proto.Int64(ud.money)

	lucky := hhnn.getUserDataByID(hhnn.luckystar) //獲取明燈

	for _, v := range ud.userbet {
		betinfo := new(netproto.HHNN_BroadCastBet)
		betinfo.Type = proto.Int32(v.hhnntype)
		betinfo.Money = proto.Int64(v.money)
		betinfo.LuckyStar = proto.Bool(false)
		if lucky != nil && lucky.getUserBetByType(v.hhnntype) > 0 {
			betinfo.LuckyStar = proto.Bool(true)
		}
		pdata.BetInfo = append(pdata.BetInfo, betinfo)
	}
	pmsg.MsgData = pdata
	hhnn.SendMessageToUser(gameuser, pmsg)
}

//获取玩家赢钱消息
func (hhnn *HHNNLogic) getUserWinMsg(selfwin int64, zhuangwin int64) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNGameResultWinID)

	pdata := new(netproto.HHNN_GameResultUserMoney)
	pdata.ZhuangWin = proto.Int64(zhuangwin)
	pdata.SelfWin = proto.Int64(selfwin)
	for i := 0; i < len(hhnn.winuserinfo); i++ {
		ud := hhnn.getUserDataByID(hhnn.winuserinfo[i])
		if ud == nil {
			continue
		}
		userwin := new(netproto.HHNN_WinUerInfo)
		userid := strconv.Itoa(int(ud.user.UserID))
		userwin.UserName = proto.String(userid)
		userwin.Win = proto.Int64(ud.win)
		pdata.WinUserInfo = append(pdata.WinUserInfo, userwin)
	}
	for j := 0; j < len(hhnn.stageuserinfo); j++ {
		ud := hhnn.getUserDataByID(hhnn.stageuserinfo[j])
		if ud == nil {
			continue
		}
		stageuserwin := new(netproto.HHNN_StageUerInfo)
		stageuserwin.UserID = proto.Int32(ud.user.UserID)
		stageuserwin.Win = proto.Int64(ud.win)
		pdata.StageUerInfo = append(pdata.StageUerInfo, stageuserwin)
	}
	pmsg.MsgData = pdata
	return pmsg
}

//获取开奖结果扑克牌消息
func (hhnn *HHNNLogic) getGameResultPokerMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNGameResultPokerID)

	pdata := new(netproto.HHNN_GameResultPoker)
	pokerhung := new(netproto.HHNN_Poker)
	pokerhei := new(netproto.HHNN_Poker)

	for _, v := range hhnn.cardhung.Cards {
		cards := new(netproto.PokerCard)
		cards.Color = proto.Int32(v.Color)
		cards.CardNumber = proto.Int32(v.CardNumber)
		pokerhung.Cards = append(pokerhung.Cards, cards)
	}
	pdata.HungPoker = pokerhung
	hungret := getNiuNiuRet(hhnn.cardhung)
	pdata.HungPoker.NiuNiuRet = proto.Int32(int32(hungret.niuRet))
	if hungret.niuRet == nnCard_NiuNiu {
		for _, v := range hhnn.cardhung.Cards {
			if v.CardNumber == 13 {
				hungret.number = 11
			}
		}
	}
	pdata.HungPoker.NiuNiuNumber = proto.Int32(int32(hungret.number))

	for _, v := range hhnn.cardhei.Cards {
		cards := new(netproto.PokerCard)
		cards.Color = proto.Int32(v.Color)
		cards.CardNumber = proto.Int32(v.CardNumber)
		pokerhei.Cards = append(pokerhei.Cards, cards)
	}
	pdata.HeiPoker = pokerhei
	heiret := getNiuNiuRet(hhnn.cardhei)

	pdata.HeiPoker.NiuNiuRet = proto.Int32(int32(heiret.niuRet))
	if heiret.niuRet == nnCard_NiuNiu {
		for _, v := range hhnn.cardhei.Cards {
			if v.CardNumber == 13 {
				heiret.number = 11
			}
		}
	}
	pdata.HeiPoker.NiuNiuNumber = proto.Int32(int32(heiret.number))

	betret := compareCards(hungret, heiret)
	if betret > int32(0) { //紅贏黑輸
		pdata.Win = proto.Int32(1)
	} else {
		pdata.Win = proto.Int32(2)
	}

	pdata.RoundID = proto.Int32(hhnn.roundid)
	pmsg.MsgData = pdata
	return pmsg
}

//廣播第二手兩張牌
//func (hhnn *HHNNLogic) boardcastSecondCardGroup() {
//	hhnn.BoardcastMessage(hhnn.getSecondCardGroup())
//}

//第二手兩張牌
//func (hhnn *HHNNLogic) getSecondCardGroup() *msg.Message {
//	pmsg := new(msg.Message)
//	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
//	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNBroadCastFapaiTwoCard)
//
//	pdata := new(netproto.HHNN_BrodcastPokerCard)
//
//	for i := int32(0); i <= int32(1); i++ {
//		cg := hhnn.getCalSecondPokerGroup(i)
//
//		urPokerCard := new(netproto.HHNN_BrodcastUserPokerCard)
//		urPokerCard.SitNo = proto.Int32(i + 1)
//
//		for _, cd := range cg.Cards {
//			card := new(netproto.PokerCard)
//			card.Color = proto.Int32(cd.Color)
//			card.CardNumber = proto.Int32(cd.CardNumber)
//			urPokerCard.CardGroup = append(urPokerCard.CardGroup, card)
//		}
//
//		pdata.UserPokerCard = append(pdata.UserPokerCard, urPokerCard)
//	}
//
//	pmsg.MsgData = pdata
//	return pmsg
//}

//取得每家的第二張牌(含計算賺賠機率的牌型)
//func (hhnn *HHNNLogic) getCalSecondPokerGroup(position int32) *poker.CardGroup {
//	return hhnn.secondPokerGroup[int(position)]
//}

//开始下注
func (hhnn *HHNNLogic) startBet() {
	hhnn.opttimer = hhnn.TimerFuncManager.NewTimer(hhnn.bettime, hhnn.kaiJiang)
	hhnn.BoardcastMessage(hhnn.getStateMsg())
}

//设置状态
func (hhnn *HHNNLogic) setState(state int32) {
	hhnn.gamestate = state
}

//开始开奖
func (hhnn *HHNNLogic) kaiJiang() {
	hhnn.setState(KaiJiang)
	hhnn.BoardcastMessage(hhnn.getStateMsg())

	if hhnn.zhuangID != 0 {
		if hhnn.GetUserByUserID(hhnn.zhuangID).IsRobot() {
			//hhnn.removeRobotBet()
		}
	} else {
		//hhnn.removeRobotBet()
	}

	if hhnn.GetServerConfig().PokerType != gamelogic.PokerTypeConfig {
		hhnn.getResultType()
	}
	hhnn.getPokerByType(hhnn.jiangtype) //獲得牌型
	hhnn.countGameRound()
	//hhnn.boardcastSecondCardGroup()
	hhnn.BoardcastMessage(hhnn.getGameResultPokerMsg())
	hhnn.BoardcastMessage(hhnn.getRoundMsg())
	hhnn.setState(KaiJiangAction)
	hhnn.BoardcastMessage(hhnn.getStateMsg())
	hhnn.addResultCardList()
	hhnn.calWinLose()
	hhnn.updateVipUserWin()
	hhnn.setState(EndAction)
	hhnn.action_str = hhnn.RegisterAnimation(GameKaiJiang_Action, hhnn.kaijiangactiontime, hhnn.sendResult)
}

//移除机器人下注额
func (hhnn *HHNNLogic) removeRobotBet() {
	for _, v := range hhnn.bettyperobot {
		hhnn.updateBet(v.hhnntype, -v.money)
	}
}

//更新下注钱
func (hhnn *HHNNLogic) updateBet(t int32, money int64) {
	hhnn.bettotal += int64(money)
	mlog.Debug("BetType=%v,money=%v,lhd.bettotal=%v", t, money, hhnn.bettotal)
	for _, v := range hhnn.bettypetotal {
		if v.hhnntype == t {
			v.money += int64(money)
			if v.money < 0 {
				mlog.Debug("钱小于0了！")
				v.money = 0
			}
			return
		}
	}

	b := new(bettype)
	b.hhnntype = t
	b.money = int64(money)
	if b.money < 0 {
		mlog.Debug("钱小于0了！")
		b.money = 0
	}
	hhnn.bettypetotal = append(hhnn.bettypetotal, b)
}

//計算回合
func (hhnn *HHNNLogic) countGameRound() {
	if hhnn.roundid >= 30 {
		hhnn.roundid = 1
	} else {
		hhnn.roundid += 1
	}
	if hhnn.jiangtype == HUNG {
		hhnn.hungwin += 1
	} else if hhnn.jiangtype == HEI {
		hhnn.heiwin += 1
	}
	mlog.Debug("开奖类型：%v", hhnn.jiangtype)
}

//添加开奖到历史列表
func (hhnn *HHNNLogic) addResultCardList() {
	retcards := new(netproto.HHNN_GameResultPoker)
	pokerhung := new(netproto.HHNN_Poker)
	pokerhei := new(netproto.HHNN_Poker)

	for _, v := range hhnn.cardhung.Cards {
		cards := new(netproto.PokerCard)
		cards.Color = proto.Int32(v.Color)
		cards.CardNumber = proto.Int32(v.CardNumber)
		pokerhung.Cards = append(pokerhung.Cards, cards)
	}
	retcards.HungPoker = pokerhung
	hungret := getNiuNiuRet(hhnn.cardhung)
	retcards.HungPoker.NiuNiuRet = proto.Int32(int32(hungret.niuRet))
	if hungret.niuRet == nnCard_NiuNiu {
		for _, v := range hhnn.cardhung.Cards {
			if v.CardNumber == 13 {
				hungret.number = 11
			}
		}
	}
	retcards.HungPoker.NiuNiuNumber = proto.Int32(int32(hungret.number))

	for _, v := range hhnn.cardhei.Cards {
		cards := new(netproto.PokerCard)
		cards.Color = proto.Int32(v.Color)
		cards.CardNumber = proto.Int32(v.CardNumber)
		pokerhei.Cards = append(pokerhei.Cards, cards)
	}
	retcards.HeiPoker = pokerhei
	heiret := getNiuNiuRet(hhnn.cardhei)
	retcards.HeiPoker.NiuNiuRet = proto.Int32(int32(heiret.niuRet))
	if heiret.niuRet == nnCard_NiuNiu {
		for _, v := range hhnn.cardhei.Cards {
			if v.CardNumber == 13 {
				heiret.number = 11
			}
		}
	}
	retcards.HeiPoker.NiuNiuNumber = proto.Int32(int32(heiret.number))

	betret := compareCards(hungret, heiret)
	if betret > int32(0) { //紅贏黑輸
		retcards.Win = proto.Int32(1)
		hhnn.recordTypeWinLossOdds(1, 0, hungret) //下注紅的賠率
		hhnn.recordTypeWinLossOdds(2, 1, hungret) //下注黑的賠率
		hhnn.recordSpecialTypeOdds(hungret)       //下注致命一擊的賠率

		mlog.Debug("紅贏 添加开奖到历史列表 , 紅[%v] 黑[%v]", hungret, heiret)
		hhnn.GameLog.AddGameStep("紅贏 紅牌 添加开奖到历史列表 , 紅[%v] 黑[%v]", hungret, heiret)
		//retcards.HungPoker.NiuNiuNumber = proto.Int32(hhnn.zuhetype)
	} else {
		retcards.Win = proto.Int32(2)
		hhnn.recordTypeWinLossOdds(1, 1, heiret) //下注紅的
		hhnn.recordTypeWinLossOdds(2, 0, heiret)
		hhnn.recordSpecialTypeOdds(heiret) //下注致命一擊的賠率

		mlog.Debug("黑贏 添加开奖到历史列表 , 紅[%v] 黑[%v]", hungret, heiret)
		hhnn.GameLog.AddGameStep("黑贏 紅牌 添加开奖到历史列表 , 紅[%v] 黑[%v]", hungret, heiret)
		//retcards.HeiPoker.NiuNiuNumber = proto.Int32(hhnn.zuhetype)
	}

	//回報指定多局開獎結果
	if hhnn.assignRewardDetail != nil {
		mlog.Debug("===========指定開獎%v,組合%v", hhnn.jiangtype, hhnn.zuhetype)
		if hhnn.isPowerJiangType(hhnn.zuhetype) {
			hhnn.assignRewardDetailResult(fmt.Sprintf("%d,%d", hhnn.jiangtype, 3))
		} else {
			hhnn.assignRewardDetailResult(fmt.Sprintf("%d", hhnn.jiangtype))
		}
	}

	retcards.RoundID = proto.Int32(hhnn.roundid)
	hhnn.resultcardlist.CardList = append(hhnn.resultcardlist.CardList, retcards)
	if len(hhnn.resultcardlist.CardList) > 100 {
		hhnn.resultcardlist.CardList = append(hhnn.resultcardlist.CardList[:0], hhnn.resultcardlist.CardList[1:]...)
	}

}

//好運連連的賠率紀錄
func (hhnn *HHNNLogic) recordSpecialTypeOdds(wincardret *niuniuCardRet) {

	jiangWinLossOdd := new(JiangTypeWinLossOdd)
	jiangWinLossOdd.betType = int32(3)
	jiangWinLossOdd.odds = getSpecialNiuNiuOdds(wincardret, hhnn)

	if jiangWinLossOdd.odds > 0 {
		jiangWinLossOdd.winLoss = int32(1)

	} else {
		jiangWinLossOdd.winLoss = int32(-1)
		jiangWinLossOdd.odds = 0
	}

	hhnn.jiangTypeWinLossOdd = append(hhnn.jiangTypeWinLossOdd, jiangWinLossOdd)
}

//紀錄開獎後下注種類的和紅黑比牌後的輸贏，並且記錄賠率
func (hhnn *HHNNLogic) recordTypeWinLossOdds(indexType int, winLoss int32, wincardret *niuniuCardRet) {
	jiangWinLossOdd := new(JiangTypeWinLossOdd)
	jiangWinLossOdd.betType = int32(indexType)

	if winLoss == int32(0) { //贏
		jiangWinLossOdd.winLoss = int32(1)
		jiangWinLossOdd.odds = getNiuNiuOdds(wincardret, hhnn)
	} else if winLoss == int32(1) { //輸
		jiangWinLossOdd.winLoss = int32(-1)
		jiangWinLossOdd.odds = getNiuNiuOdds(wincardret, hhnn) - 2 // -1是把投注成本拿回
	}

	hhnn.jiangTypeWinLossOdd = append(hhnn.jiangTypeWinLossOdd, jiangWinLossOdd)
}

//確認开奖类型。
func (hhnn *HHNNLogic) getResultType() {
	//獲取監控殺放分清單
	superToolHandler := hhnn.SuperToolService.GetBrGameSuperToolHandler()
	//指定开奖
	if hhnn.nextjiangtype != 0 {
		hhnn.jiangtype = hhnn.nextjiangtype
		hhnn.zuhetype = hhnn.nextzuhetype
		hhnn.GameLog.AddGameStep("指定开奖:%v, %v", hhnn.jiangtype, hhnn.zuhetype)
		hhnn.nextjiangtype = 0
		hhnn.nextzuhetype = 0
		return
	} else if ok, rettype := hhnn.isAssignReward(); ok { //多局指定开奖

		hhnn.execAssignReward = true
		if len(rettype) > 0 {
			hhnn.jiangtype = rettype[0]
			if len(rettype) > 1 && rettype[1] == 3 {
				hhnn.zuhetype = hhnn.getPowerJiangType()
			} else {
				pokertype := []int32{WUXIAONIU, ZHADAN, WUHUANIU, KNIUNIU, NIUNIU, NIUNINE, NIUEIGHT, NIUSEVEN, NIUSIX, NIUFIVE, NIUFOUR, NIUTHREE, NIUTWO, NIUONE, WUNIU}
				hhnn.zuhetype = hhnn.weight[1].getJiangType(pokertype)
			}
		}
		hhnn.GameLog.AddGameStep("多局指定开奖: rettype %v, jiangtype %v, zhuetype %v", rettype, hhnn.jiangtype, hhnn.zuhetype)
		return
	} else if hhnn.fangFenTarget > 0 && hhnn.fangFenComplete < hhnn.fangFenTarget { //系統放分
		hhnn.execFangFen = true
		alltype := hhnn.getMaxKaiJiangType()
		hhnn.jiangtype = hhnn.weight[1].getJiangType(alltype)
		zuhetype := hhnn.getFangFenZuheType()
		hhnn.zuhetype = hhnn.weight[1].getJiangType(zuhetype)
		hhnn.GameLog.AddGameStep("系統放分 : target %d, already %d", hhnn.fangFenTarget, hhnn.fangFenComplete)
		return
	} else if len(hhnn.getBowBet()) > 0 { //監控放殺分
		if hhnn.bowJiangType() {
			pokertype := []int32{NIUNIU, NIUNINE, NIUEIGHT, NIUSEVEN, NIUSIX, NIUFIVE, NIUFOUR, NIUTHREE, NIUTWO, NIUONE, WUNIU}
			hhnn.zuhetype = hhnn.weight[1].getJiangType(pokertype)
			example.invalid("執行監控放殺分策略 : %s", superToolHandler.GetProiorityListString())
			hhnn.GameLog.AddGameStep("執行監控放殺分策略 : %s", superToolHandler.GetProiorityListString())
			return
		}
	} else if hhnn.noRealUserBet() { //沒有真人下注時, 開最漂高的機率表
		alltype := []int32{HUNG, HEI, NATURE}
		pokertype := []int32{WUXIAONIU, ZHADAN, WUHUANIU, KNIUNIU, NIUNIU, NIUNINE, NIUEIGHT, NIUSEVEN, NIUSIX, NIUFIVE, NIUFOUR, NIUTHREE, NIUTWO, NIUONE, WUNIU}
		area := int32(len(hhnn.scoreSection) - 1)
		hhnn.jiangtype = hhnn.weight[area].getJiangType(alltype)
		hhnn.zuhetype = hhnn.weight[area].getJiangType(pokertype)
		example.invalid("無真人玩家下注, jiangtype [%d], zhueytpe [%d]", hhnn.jiangtype, hhnn.zuhetype)
		hhnn.GameLog.AddGameStep("無真人玩家下注, jiangtype [%d], zhueytpe [%d]", hhnn.jiangtype, hhnn.zuhetype)
		return
	}

	alltype := []int32{HUNG, HEI, NATURE}
	pokertype := make([]int32, 0)
	if hhnn.syspool < hhnn.winline { //贏錢線下不開特殊獎
		pokertype = []int32{NIUNIU, NIUNINE, NIUEIGHT, NIUSEVEN, NIUSIX, NIUFIVE, NIUFOUR, NIUTHREE, NIUTWO, NIUONE, WUNIU}
	} else {
		pokertype = []int32{WUXIAONIU, ZHADAN, WUHUANIU, KNIUNIU, NIUNIU, NIUNINE, NIUEIGHT, NIUSEVEN, NIUSIX, NIUFIVE, NIUFOUR, NIUTHREE, NIUTWO, NIUONE, WUNIU}
	}
	area := int32(0)
	if hhnn.syspool > hhnn.scoreSection[len(hhnn.scoreSection)-1] {
		area = int32(len(hhnn.scoreSection) - 1)
	} else { //依水線獲取機率表
		for i := int32(0); i < int32(len(hhnn.scoreSection))-1; i++ {
			if hhnn.syspool > hhnn.scoreSection[i] && hhnn.syspool < hhnn.scoreSection[i+1] {
				area = i
			}
		}
	}

	hhnn.jiangtype = hhnn.weight[area].getJiangType(alltype)
	//根據權重獲得贏家的牌型
	hhnn.zuhetype = hhnn.weight[area].getJiangType(pokertype)

	mlog.Debug("預設開獎類型:%d", hhnn.jiangtype)
	mlog.Debug("預設開獎牌型:%d", hhnn.zuhetype)
	hhnn.GameLog.AddGameStep("預設開獎類型:%v", hhnn.jiangtype)
	hhnn.GameLog.AddGameStep("預設開獎牌型:%v", hhnn.zuhetype)

	//(2)低於贏錢線
	if hhnn.syspool < hhnn.winline {
		//當水線低於贏錢線，開最小獎
		alltype := hhnn.getMinKaiJiangType()
		example.invalid("開最小獎 alltype [%v], pokertype [%v]", alltype, pokertype)
		hhnn.jiangtype = hhnn.weight[0].getJiangType(alltype)
		if hhnn.zuhetype >= WUXIAONIU && hhnn.zuhetype <= KNIUNIU {
			pokertype := []int32{NIUEIGHT, NIUSEVEN, NIUSIX, NIUFIVE, NIUFOUR, NIUTHREE, NIUTWO, NIUONE, WUNIU}
			hhnn.zuhetype = hhnn.weight[0].getJiangType(pokertype)
		}
	}

	//(3)開大獎的驗算
	if hhnn.zuhetype >= WUXIAONIU && hhnn.zuhetype <= KNIUNIU {
		odd := int64(hhnn.paytable[0].betvalue[hhnn.zuhetype-3].criticalodd)
		if hhnn.syspool-hhnn.getBetMoneyByType(3)*odd < hhnn.criticalwinline {
			pokertype := []int32{NIUEIGHT, NIUSEVEN, NIUSIX, NIUFIVE, NIUFOUR, NIUTHREE, NIUTWO, NIUONE, WUNIU}
			hhnn.zuhetype = hhnn.weight[0].getJiangType(pokertype)
		}
	}

	mlog.Debug("最終開獎類型:%d", hhnn.jiangtype)
	mlog.Debug("最終開獎牌型:%d", hhnn.zuhetype)
	hhnn.GameLog.AddGameStep("最終開獎類型:%v", hhnn.jiangtype)
	hhnn.GameLog.AddGameStep("最終開獎牌型:%v", hhnn.zuhetype)
}

//获取預設開獎类型。
func (hhnn *HHNNLogic) getFirstResultType() {

}

//获取赔偿最小的开奖结果
func (hhnn *HHNNLogic) getMinKaiJiangType() []int32 {
	mlog.Debug("開不出最小獎, Long:%v HU: %v", hhnn.getBetMoneyByType(HUNG), hhnn.getBetMoneyByType(HEI))
	if hhnn.getRealManBetMoneyByType(HUNG) > hhnn.getRealManBetMoneyByType(HEI) {
		return []int32{HEI}
	} else if hhnn.getRealManBetMoneyByType(HUNG) < hhnn.getRealManBetMoneyByType(HEI) {
		return []int32{HUNG}
	}

	return []int32{hhnn.jiangtype}
}

//获取赔偿最大的开奖结果
func (hhnn *HHNNLogic) getMaxKaiJiangType() []int32 {
	if hhnn.getRealManBetMoneyByType(HUNG) > hhnn.getRealManBetMoneyByType(HEI) {
		return []int32{HUNG}
	} else if hhnn.getRealManBetMoneyByType(HUNG) < hhnn.getRealManBetMoneyByType(HEI) {
		return []int32{HEI}
	}
	return []int32{HUNG, HEI}
}

//获取开奖的扑克信息
func (hhnn *HHNNLogic) getPokerByType(t int32) {
	winCard := new(niuniuCardRet)
	loseCard := new(niuniuCardRet)
	if hhnn.GetServerConfig().PokerType == gamelogic.PokerTypeConfig {
		hhnn.jiangtype = NATURE
	} else {
		hhnn.poker.Shuffle()
		winCardType := getCardsbyType(int(hhnn.zuhetype), hhnn)
		winCard, loseCard = getLoseCard(hhnn, winCardType)
	}
	//根據預設開獎類型進行組牌及留牌
	if hhnn.jiangtype == HUNG {
		hhnn.cardhung = winCard.cards
		hhnn.cardhei = loseCard.cards
	} else if hhnn.jiangtype == HEI {
		hhnn.cardhung = loseCard.cards
		hhnn.cardhei = winCard.cards
	} else {
		for i := int32(0); i <= int32(1); i++ {
			cg := hhnn.poker.Deal(0, 5) //固定userid 0 只為了取得隨機三張牌

			urPokerCard := new(netproto.BRNN_BrodcastUserPokerCard)
			urPokerCard.SitNo = proto.Int32(i)

			for _, cd := range cg.Cards {
				card := new(netproto.PokerCard)
				card.Color = proto.Int32(cd.Color)
				card.CardNumber = proto.Int32(cd.CardNumber)
				urPokerCard.CardGroup = append(urPokerCard.CardGroup, card)
			}
			if i == 0 {
				hhnn.cardhung = cg
			} else {
				hhnn.cardhei = cg
			}
		}
	}
	hungNiuRet := getNiuNiuRet(hhnn.cardhung)
	heiNiuRet := getNiuNiuRet(hhnn.cardhei)

	betret := compareCards(hungNiuRet, heiNiuRet)
	if betret > int32(0) { //紅贏黑輸
		hhnn.jiangtype = HUNG
	} else {
		hhnn.jiangtype = HEI
	}

	//hungFirstPoker := new(poker.CardGroup)
	//heiFirstPoker := new(poker.CardGroup)
	//hungSecondPoker := new(poker.CardGroup)
	//heiSecondPoker := new(poker.CardGroup)
	//for index := 0; index < 5; index++ {
	//	//將cardhung及hei分別放入前三張及後兩張牌堆
	//	if index < 3 {
	//		hungFirstPoker.Cards = append(hungFirstPoker.Cards, hhnn.cardhung.Cards[index])
	//		heiFirstPoker.Cards = append(heiFirstPoker.Cards, hhnn.cardhei.Cards[index])
	//	} else {
	//		hungSecondPoker.Cards = append(hungSecondPoker.Cards, hhnn.cardhung.Cards[index])
	//		heiSecondPoker.Cards = append(heiSecondPoker.Cards, hhnn.cardhei.Cards[index])
	//	}
	//}
	//
	//hhnn.firstPoker[int(0)] = hungFirstPoker
	//hhnn.firstPoker[int(1)] = heiFirstPoker
	//hhnn.secondPokerGroup[int(0)] = hungSecondPoker
	//hhnn.secondPokerGroup[int(1)] = heiSecondPoker

	//hhnn.GameLog.AddGameStep("紅方手牌:%v, %v ", hhnn.firstPoker[0], hhnn.secondPokerGroup[int(0)])
	//hhnn.GameLog.AddGameStep("黑方手牌:%v, %v ", hhnn.firstPoker[1], hhnn.secondPokerGroup[int(1)])
	//mlog.Debug("紅方手牌:%v ", hhnn.firstPoker)
	//mlog.Debug("黑方手牌:%v ", hhnn.firstPoker)

	//if !hhnn.secondPokerGroup[0].Cards[0].IsEqual(hhnn.cardhung.Cards[3]) {
	//	mlog.Fatal("===================== fatal")
	//}

	hhnn.GameLog.AddGameStep("开奖结果是%d,紅區的扑克牌是%v ,黑區的扑克牌是%v", hhnn.jiangtype, hhnn.cardhung, hhnn.cardhei)
	example.invalid("开奖结果是%d,紅區的扑克牌是%v ,黑區的扑克牌是%v", t, hhnn.cardhung, hhnn.cardhei)
}

//func (hhnn *HHNNLogic) appendSecondPokerGroup(index int, niuCardRet *niuniuCardRet) bool {
//	var pokerGroup []*poker.CardGroup
//	niuCardRet.cards.Cards = niuCardRet.cards.Cards[3:5]
//	pokerGroup = append(pokerGroup, niuCardRet.cards)
//	hhnn.secondPokerGroup[index] = niuCardRet.cards
//	if len(hhnn.secondPokerGroup) >= int(5) {
//		return false
//	}
//	return true
//}

//隨機自然發牌
//func (hhnn *HHNNLogic) natureResultType() {
//	for i := int32(0); i <= int32(1); i++ {
//		hhnn.poker.Shuffle()
//		cg := hhnn.poker.Deal(0, 2)
//
//		urCard := new(poker.CardGroup)
//
//		for _, cd := range cg.Cards {
//			pCard := new(poker.PokerCard)
//			pCard.Color = cd.Color
//			pCard.CardNumber = cd.CardNumber
//			urCard.Cards = append(urCard.Cards, pCard)
//		}
//
//		hhnn.secondPokerGroup[int(i)] = urCard
//	}
//
//	hhnn.cardhung.Cards = append(hhnn.cardhung.Cards, hhnn.firstPoker[0].Cards...)
//	hhnn.cardhung.Cards = append(hhnn.cardhung.Cards, hhnn.secondPokerGroup[0].Cards...)
//	hhnn.cardhei.Cards = append(hhnn.cardhei.Cards, hhnn.firstPoker[1].Cards...)
//	hhnn.cardhei.Cards = append(hhnn.cardhei.Cards, hhnn.secondPokerGroup[1].Cards...)
//}

//给玩家发送结果
func (hhnn *HHNNLogic) sendResult() {
	//hhnn.calWinLose()

	hhnn.writePoolToFile()
	hhnn.action_str = hhnn.RegisterAnimation(GameEnd_Action, hhnn.endactiontime, hhnn.gameEnd)
}

//用户消息主处理函数
func (hhnn *HHNNLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	if !hhnn.CheckGameUserMessage(user) {
		hhnn.KickUserFunc(user.UserID, "非法操作")
		return
	}
	//不是给游戏逻辑服务器的消息
	mlog.Debug("当前游戏状态为：%v", hhnn.gamestate)
	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	hhnn.CheckUserAnimationCommited(user, gmsg, hhnn.action_str, nil)

	if gmsg.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNGetShangZhuangListID) {
		hhnn.SendMessageToUser(user, hhnn.getZhuangListMsg())
		return
	}

	if hhnn.gamestate != Bet {
		return
	}

	if gmsg.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNBetID) {
		hhnn.userBet(user, gmsg.MsgData.([]byte))
		return
	}

	if gmsg.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNShangZhuangID) {
		hhnn.shangZhuang(user, gmsg.MsgData.([]byte))
		return
	}

	if gmsg.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNBetClearID) && hhnn.gamestate == Bet {
		hhnn.clearUserBet(user)
		return
	}

	//if gmsg.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNGetFapeiThreeCard) && hhnn.gamestate == Bet {
	//	hhnn.sendToUserFirstCardGroup(user)
	//}
}

func (hhnn *HHNNLogic) sendToUserFirstCardGroup(gameuser *gamelogic.GameUser) {
	hhnn.SendMessageToUser(gameuser, hhnn.getFirstCardGroup())
}

//将水池的钱写入文件
func (hhnn *HHNNLogic) writePoolToFile() {
	if conf.TestMode {
		return
	}

	buff := strconv.Itoa(int(hhnn.syspool))
	//bufio.ReadWriter.Peek()
	tablenostr := strconv.Itoa(int(hhnn.TableNo()))
	filename := "./lastpool/" + hhnn.serverid + "_" + tablenostr + ".conf"
	_, err := os.Stat(filename)
	var f *os.File = nil
	f, err = os.Create(filename)
	if err != nil {
		mlog.Error("%v文件不存在", filename)
		panic(err)
	} else {
		f.WriteString(buff)
		defer f.Close()
	}
}

//游戏结束
func (hhnn *HHNNLogic) gameEnd() {
	hhnn.setState(Wait)
	hhnn.BoardcastMessage(hhnn.getStateMsg())
	if hhnn.zhuangID != 0 {
		user := hhnn.GetUserByUserID(hhnn.zhuangID)
		if user.IsOffLine() {
			hhnn.zhuangID = 0
		}
	}
	hhnn.GameEndFunc()
}

//获取上庄列表消息
func (hhnn *HHNNLogic) getZhuangListMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNShangZhuangListInfoID)

	pdata := new(netproto.HHNN_ZhuangList)
	for i := 0; i < len(hhnn.zhuanglist); i++ {
		ur := hhnn.GetUserByUserID(hhnn.zhuanglist[i])
		pdata.UserInfo = append(pdata.UserInfo, ur.GetProtoBufGameUser())
	}
	pmsg.MsgData = pdata
	return pmsg
}

//用户下注处理函数
func (hhnn *HHNNLogic) userBet(user *gamelogic.GameUser, data []byte) {
	ud := hhnn.getUserDataByID(user.UserID)
	if ud == nil {
		ud = hhnn.addUserData(user)
	}

	ust := &netproto.HHNN_BetArr{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Debug("解析用户%v下注错误%v", user, err)
		return
	}

	//庄家不能下注
	if hhnn.zhuangID == user.UserID {
		return
	}

	money := int32(0)
	for _, b := range ust.BetArr {
		*b.BetType = int32(math.Abs(float64(*b.BetType)))
		*b.Count = int32(math.Abs(float64(*b.Count)))

		//玩家身上金額小於下注限制
		if user.Score < hhnn.betlimitmoney*100 {
			for i, u := range hhnn.userdata {
				if u.user.UserID == ud.user.UserID {
					hhnn.userdata = append(hhnn.userdata[:i], hhnn.userdata[i+1:]...)
				}
			}
			hhnn.SendMessageToUser(user, hhnn.getBetRetMsg(ud, -6, *b.BetType, ust.BetArr))
			return
		}

		//下注不可超過還可下注金額和玩家的十分之一金額
		userBetTotal := int64(0)
		if ud != nil && ud.userbet != nil {
			for i := int32(0); i < int32(len(ud.userbet)); i++ {
				if ud.userbet[i] != nil {
					userBetTotal += ud.userbet[i].money
				}
			}
		}

		ableBet := user.Score / int64(10)
		thisBet := int64(0)
		for _, b := range ust.BetArr {
			ableBet = ableBet - int64(*b.Count)
			thisBet += int64(*b.Count)
		}

		if (ableBet - userBetTotal) <= int64(0) { //超過可下注的十分之一
			hhnn.SendMessageToUser(user, hhnn.getBetRetMsg(ud, -6, *b.BetType, ust.BetArr))
			return
		}

		if hhnn.sametimebetmode != 1 { //關閉庄閒同時下注的配置
			//黑紅不能同時下注
			if (*b.BetType == int32(HUNG) && hhnn.getUserBetByType(ud, HEI) > 0) ||
				(*b.BetType == int32(HEI) && hhnn.getUserBetByType(ud, HUNG) > 0) {
				hhnn.SendMessageToUser(user, hhnn.getBetRetMsg(ud, -5, *b.BetType, ust.BetArr))
				return
			}
		}

		tmpmon := hhnn.getMoneyByType(*b.BetType) * (*b.Count)

		// 系統下注上限
		if (int64(tmpmon) + hhnn.getBetMoneyByType(*b.BetType)) > hhnn.betvalue.getMaxMoneyByType(*b.BetType) {
			hhnn.SendMessageToUser(user, hhnn.getBetRetMsg(ud, -4, *b.BetType, ust.BetArr))
			return
		}

		// 不是系統當庄, 則壓注上限為庄家身上的錢所能夠賠付的錢
		if hhnn.zhuangID != 0 {
			ur := hhnn.GetUserByUserID(hhnn.zhuangID)
			odd := int64(0)
			if odd = int64(hhnn.betvalue.getOddByType(b.GetBetType())) - 1; odd < 0 {
				odd = 0
			}

			allmoney := (int64(tmpmon) + hhnn.getBetMoneyByType(*b.BetType)) * odd
			if allmoney > ur.Score {
				hhnn.SendMessageToUser(user, hhnn.getBetRetMsg(ud, -4, *b.BetType, ust.BetArr))
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
			betmon := hhnn.getMoneyByType(*b.BetType) * (*b.Count)
			ud.setUserBet(*b.BetType, betmon)
			if user.IsRobot() {
				hhnn.updaterobotBet(*b.BetType, int64(betmon))
			} else {
				hhnn.updateBet(*b.BetType, int64(betmon))
			}
			if user.UserID == hhnn.luckystar {
				hhnn.notbettimes = 0
			}
			hhnn.GameLog.AddGameStep("用户%d下注:%s;金額:%d, %v", user.UserID, hhnn.getBetTypeString(int(b.GetBetType())), betmon, b)
		}
		hhnn.BoardcastMessage(hhnn.getBetTotalMsg(user.UserID))
	}
	hhnn.SendMessageToUser(user, hhnn.getBetRetMsg(ud, ret, 0, ust.BetArr))
}

//添加一个用户数据
func (hhnn *HHNNLogic) addUserData(user *gamelogic.GameUser) *hhnnuserdata {
	ud := newLongHuDouUserData(user)
	hhnn.userdata = append(hhnn.userdata, ud)
	return ud
}

//获取下注消息
func (hhnn *HHNNLogic) getBetRetMsg(ud *hhnnuserdata, ret int32, t int32, betinfo []*netproto.HHNN_Bet) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNBetRetID)

	pdata := new(netproto.HHNN_BetRet)
	pdata.Money = proto.Int64(ud.money)
	pdata.Ret = proto.Int32(ret)
	pdata.Type = proto.Int32(t)
	pdata.BetArr = betinfo
	pmsg.MsgData = pdata
	return pmsg
}

func (hhnn *HHNNLogic) getUserBetByType(ud *hhnnuserdata, t int32) int64 {
	money := int64(0)
	for _, bet := range ud.userbet {
		if bet.hhnntype == t {
			money += bet.money
		}
	}
	return money
}

//获取此类型一次下注额
func (hhnn *HHNNLogic) getMoneyByType(t int32) int32 {
	return hhnn.betvalue.getMoneyByType(t)
}

//更新机器人下注额
func (hhnn *HHNNLogic) updaterobotBet(t int32, money int64) {
	hhnn.bettotalrobot += int64(money)
	for _, v := range hhnn.bettyperobot {
		if v.hhnntype == t {
			v.money += money
			if v.money < 0 {
				v.money = 0
			}
			mlog.Debug("下注类型%d，下注的总钱数为%d", t, v.money)
			return
		}
	}

	b := new(bettype)
	b.hhnntype = t
	b.money = int64(money)
	if b.money < 0 {
		b.money = 0
	}

	hhnn.bettyperobot = append(hhnn.bettyperobot, b)
}

func (hhnn *HHNNLogic) getBetTypeString(idx int) string {
	if idx > 0 && idx < len(resultstring) {
		return resultstring[idx]
	}
	return ""
}

//处理用户上庄请求
func (hhnn *HHNNLogic) shangZhuang(user *gamelogic.GameUser, data []byte) {
	ust := &netproto.HHNN_ShangZhuang{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Debug("解析用户%v上庄错误%v", user, err)
		return
	}

	if hhnn.getZhuangListIndexByUserId(user.UserID) != -1 {
		hhnn.SendMessageToUser(user, hhnn.getShangzhuangRetMsg(-2))
		return
	}

	if user.Score < hhnn.zhuanglimit {
		hhnn.SendMessageToUser(user, hhnn.getShangzhuangRetMsg(-1))
		return
	}

	if *ust.IsShangZhuang {
		if hhnn.publiczhuang == 1 {
			hhnn.GameLog.AddGameStep("用户%d上庄", user.UserID)
			hhnn.addToZhuangList(user.UserID)
			hhnn.SendMessageToUser(user, hhnn.getShangzhuangRetMsg(0))
		} else {
			mlog.Error("未開啟上庄機制，不得上庄")
			return
		}
	} else {
		hhnn.GameLog.AddGameStep("用户%d下庄", user.UserID)
		hhnn.delFromZhuangList(user.UserID)
		hhnn.SendMessageToUser(user, hhnn.getShangzhuangRetMsg(1))
		hhnn.lastturn = -1
	}
}

//获取上庄列表索引通过用户ID
func (hhnn *HHNNLogic) getZhuangListIndexByUserId(userid int32) int32 {
	for i := 0; i < len(hhnn.zhuanglist); i++ {
		if hhnn.zhuanglist[i] == userid {
			return int32(i)
		}
	}
	return -1
}

//获取上庄结果消息
func (hhnn *HHNNLogic) getShangzhuangRetMsg(ret int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNShangZhuangRetID)

	pdata := new(netproto.HHNN_ShangzhuangRet)
	pdata.Ret = proto.Int32(ret)
	pmsg.MsgData = pdata
	return pmsg
}

//清除玩家下注信息
func (hhnn *HHNNLogic) clearUserBet(user *gamelogic.GameUser) {
	ud := hhnn.getUserDataByID(user.UserID)
	if ud == nil {
		return
	}

	mlog.Debug("ud.clearcount=%v", ud.clearcount)
	if ud.clearcount >= hhnn.maxclearcount {
		hhnn.SendMessageToUser(user, hhnn.getClearRetMsg(ud.money, -1))
		return
	}

	ud.clearcount++
	for _, ub := range ud.getUserBet() {
		hhnn.updateBet(ub.hhnntype, -ub.money)
	}

	hhnn.GameLog.AddGameStep("用户%d点击清除", user.UserID)
	ud.clearData()
	hhnn.SendMessageToUser(user, hhnn.getClearRetMsg(user.Score, 0))
	hhnn.BoardcastMessage(hhnn.getBetTotalMsg(user.UserID))
}

//获取清除结果消息
func (hhnn *HHNNLogic) getClearRetMsg(money int64, ret int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNBetClearRetID)

	pdata := new(netproto.HHNN_BetClearRet)
	pdata.Money = proto.Int64(money)
	pdata.Ret = proto.Int32(ret)
	pmsg.MsgData = pdata

	return pmsg
}

//上莊是為系統或機器人
func (hhnn *HHNNLogic) isZhuangRobotAndSystem() bool {
	var iszhuangrobot bool = false
	ur := hhnn.GetUserByUserID(hhnn.zhuangID)

	if ur != nil && ur.IsRobot() {
		iszhuangrobot = true
	}

	if hhnn.zhuangID == 0 {
		iszhuangrobot = true
	}

	return iszhuangrobot
}

//获取机器人赢最多的开奖方案
func (hhnn *HHNNLogic) getMaxRobotJiang() int32 {
	maxtype := int32(0)
	maxmoney := int64(0)
	for _, rtd := range hhnn.bettyperobot {
		money := int64(hhnn.betvalue.getOddByType(rtd.hhnntype)) * rtd.money
		if money >= maxmoney {
			maxtype = rtd.hhnntype
			maxmoney = money
		}
	}

	return maxtype
}

//計算下注輸贏結果
func (hhnn *HHNNLogic) calWinLose() {
	sysTaxIncomeBefore := hhnn.systaxincome
	hhnn.syslastwin = 0

	for _, ud := range hhnn.userdata {
		hhnn.sendUserResult(ud)
	}

	finalBetTotal := hhnn.bettotal + hhnn.bettotalrobot //真人+機器人下注輸贏總和
	//這裡計算庄家輸贏及稅收

	var income, tax int64
	urZhuang := hhnn.GetUserByUserID(hhnn.zhuangID)
	if hhnn.zhuangID == 0 || urZhuang.IsRobot() {
		//莊家是機器人或系統, 莊賺錢只抽真人下注的金額bettotal
		_, tax = hhnn.GetUserWinAndTax(hhnn.bettotal)
	} else {
		//莊家是真人, 莊賺錢要抽真人和機器人下注的金額
		_, tax = hhnn.GetUserWinAndTax(finalBetTotal)
	}
	income, _ = hhnn.GetUserWinAndTax(finalBetTotal)

	hhnn.systaxincome += tax
	_, finalBetTax := hhnn.GetUserWinAndTax(finalBetTotal)
	hhnn.GameLog.AddGameStep("莊家贏錢 %d, 抽稅 %d", income, finalBetTax)
	mlog.Debug("莊家贏錢 %d, 抽稅 %d", income, finalBetTax)

	//這裡計算水池, 系統小妹當莊
	hhnn.zhuangWin += income
	//準備計算水池異動值
	changeSysPoolValue := hhnn.zhushuimoney
	if hhnn.zhuangID == 0 {
		hhnn.syslastwin = income

		//系統小妹當莊的時候，玩家贏錢水池減少
		if hhnn.userwinlose > 0 {
			//玩家贏錢，水線下降
			//hhnn.syspool -= hhnn.userwinlose - hhnn.zhushuimoney
			changeSysPoolValue = hhnn.zhushuimoney - hhnn.userwinlose
		} else {
			//玩家輸錢，水線上升(稅)
			//hhnn.syspool -= int64(math.Floor(float64(hhnn.userwinlose)*hhnn.GetTax())) - hhnn.zhushuimoney
			changeSysPoolValue = hhnn.zhushuimoney - int64(math.Floor(float64(hhnn.userwinlose)*hhnn.GetTax()))
		}
		mlog.Debug("系统水池異動%d：%d", -hhnn.userwinlose, hhnn.syspool)
	} else {
		ur := hhnn.GetUserByUserID(hhnn.zhuangID)
		if ur != nil {
			income, tax = hhnn.GetUserWinAndTax(finalBetTotal)
			//機器人當莊的時候
			if ur.IsRobot() {
				//玩家贏錢水池減少，並將注出的水倒回水池
				if hhnn.userwinlose > 0 {
					hhnn.systaxincome += tax //真人贏錢收稅
					//玩家贏錢，水線下降
					hhnn.syslastwin = hhnn.bettotal //系統上一局收入來源只記錄賠給真人的錢
					//hhnn.syspool -= hhnn.userwinlose - hhnn.zhushuimoney
					changeSysPoolValue = hhnn.zhushuimoney - hhnn.userwinlose
				} else {
					//玩家輸錢，水線上升(稅)
					//hhnn.syspool -= int64(math.Floor(float64(hhnn.userwinlose)*hhnn.GetTax())) - hhnn.zhushuimoney
					changeSysPoolValue = hhnn.zhushuimoney - int64(math.Floor(float64(hhnn.userwinlose)*hhnn.GetTax()))
					realUserIncome, _ := hhnn.GetUserWinAndTax(hhnn.bettotal)
					hhnn.syslastwin = realUserIncome
				}
				mlog.Debug("系统水池異動%d：%d", -hhnn.userwinlose, hhnn.syspool)
			} else {

				//玩家當莊的時候，機器人贏錢水池就添加，並將注出的水倒回水池
				if hhnn.robotwinlose > 0 {
					//機器人贏錢 水位上升(稅)
					income, tax = hhnn.GetUserWinAndTax(hhnn.robotwinlose)
					//hhnn.syspool += income - hhnn.robotwinlose
					changeSysPoolValue = hhnn.zhushuimoney + income
					hhnn.syslastwin = income
					hhnn.systaxincome = tax
				} else {
					//hhnn.syspool += hhnn.robotwinlose - hhnn.zhushuimoney
					changeSysPoolValue = hhnn.zhushuimoney + hhnn.robotwinlose
					hhnn.syslastwin += hhnn.robotwinlose
				}
				mlog.Debug("系统水池異動%d：%d", hhnn.robotwinlose, hhnn.syspool)

			}

			hhnn.UpdateUserScoreFunc(ur.UserID, hhnn.zhuangWin, tax) // 更新金幣異動
			hhnn.SendMessageToUser(ur, hhnn.getUserWinMsg(hhnn.zhuangWin, hhnn.zhuangWin))
		} else {
			mlog.Debug("庄家数据为空")
		}
	}
	if !hhnn.execFangFen {
		hhnn.syspool += changeSysPoolValue
		mlog.Debug("系统水池異動%d：%d", changeSysPoolValue, hhnn.syspool)
	} else {
		before := hhnn.fangFenComplete
		changeSysPoolValue = -changeSysPoolValue //水池下降代表玩家贏錢
		if changeSysPoolValue > 0 {              //玩家贏錢要抽稅
			changeSysPoolValue = int64(math.Floor(float64(changeSysPoolValue) * hhnn.GetTax()))
		}
		hhnn.fangFenComplete += changeSysPoolValue
		hhnn.GameLog.AddGameStep("系統放分完成, target %d, already %d-%d", hhnn.fangFenTarget, before, hhnn.fangFenComplete)
		example.invalid("系統放分完成, target %d, already %d-%d", hhnn.fangFenTarget, before, hhnn.fangFenComplete)
	}

	hhnn.GameLog.AddGameStep("本局稅收 %d", hhnn.systaxincome-sysTaxIncomeBefore)
	hhnn.GameLog.AddGameStep("系統累積稅收 %d", hhnn.systaxincome)

	uscoresRet := make(map[int32]int64)
	//發送沒下注的user 莊家結果
	for _, u2 := range hhnn.AllTableUser() {
		sendBool := true
		for _, ud := range hhnn.userdata {
			if u2.UserID == ud.user.UserID {
				sendBool = false
				break
			}
		}
		if sendBool {
			hhnn.SendMessageToUser(u2, hhnn.getUserWinMsg(0, hhnn.zhuangWin))
			if u2.UserID == hhnn.luckystar {
				if hhnn.getUserDataByID(hhnn.luckystar) == nil {
					hhnn.notbettimes += 1
				}
			}
		}
	}

	//發送下注user結算結果
	for _, ud := range hhnn.userdata {
		hhnn.SendMessageToUser(ud.user, hhnn.getUserWinMsg(ud.win, hhnn.zhuangWin))
		uscoresRet[ud.user.UserID] = ud.updatemoney
		hhnn.GameLog.Context().UserTaxAmount(ud.user.UserID, ud.updatemoney)
	}

	uscoresRet[hhnn.zhuangID] = hhnn.zhuangWin
	hhnn.GameLog.Context().UserTaxAmount(hhnn.zhuangID, hhnn.zhuangWin)
	mlog.Debug("莊家贏錢：%v", uscoresRet[hhnn.zhuangID])
	hhnn.GameLog.GameEnd()
	hhnn.bettotal = 0
}

//给每个玩家发送开奖结果
func (hhnn *HHNNLogic) sendUserResult(ud *hhnnuserdata) {
	var betlosemoney int64                     //下注輸給莊的錢
	var flowAmount, betAmount, payAmount int64 //流水, 有效投注金額, 返還金額計算
	for i := 0; i < len(ud.userbet); i++ {
		betlosemoney += ud.userbet[i].money
		betAmount = betlosemoney
		for _, v := range hhnn.jiangTypeWinLossOdd {
			if ud.userbet[i].hhnntype == v.betType { //有下注到該項目
				winlossmoney := calWinLossMoney(ud, i, v)
				//if winlossmoney > 0 {
				//	flowamountTotal := winlossmoney - ud.userbet[i].money
				//	flowAmountAfterTax, _ := hhnn.GetUserWinAndTax(flowamountTotal)
				//	flowamount += flowAmountAfterTax
				//} else {
				//	flowamount += -(winlossmoney - ud.userbet[i].money)
				//}
				//機器人不計算獎池裡面的錢
				if hhnn.zhuangID == 0 {
					if !ud.user.IsRobot() {
						hhnn.bettotal -= winlossmoney
					}
				} else {
					ur := hhnn.GetUserByUserID(hhnn.zhuangID) //當莊是否為機器人
					if ur != nil && ur.IsRobot() && ud.user.IsRobot() {
						hhnn.bettotalrobot -= winlossmoney
					} else {
						if ud.user.IsRobot() {
							hhnn.bettotalrobot -= winlossmoney
						} else {
							hhnn.bettotal -= winlossmoney
						}
					}
				}

				ud.money += winlossmoney
			}
		}
	}

	money := ud.money - ud.user.Score //ud.money已經提早先扣除下注成本了

	//將機器人及玩家輸贏減掉稅
	if ud.user.IsRobot() {
		hhnn.robotwinlose += money
	} else {
		hhnn.userwinlose += money
	}

	var income, tax int64
	if money > int64(0) {
		income, tax = hhnn.GetUserWinAndTax(money)
		ud.win = income + betlosemoney //client動畫顯示輸贏金額
		if !ud.user.IsRobot() {        //真人贏錢才紀錄系統稅收, 機器人下注贏錢不紀錄系統稅收
			hhnn.systaxincome += tax
		}
	} else {
		income = money
		ud.win = money + betlosemoney //client動畫顯示輸贏金額
	}
	flowAmount = util.ABSInt64(income)
	payAmount = betAmount + income
	ud.updatemoney = income
	stepmsg := "用户%d下注赢钱%d, 抽稅%d, 流水%d, 有效投注金額%d, 返還金額%d"
	if ud.user.IsRobot() {
		stepmsg = "[機器人]" + stepmsg
	}

	hhnn.GameLog.AddGameStep(stepmsg, ud.user.UserID, income, tax, flowAmount, betAmount, payAmount)
	ud.user.FlowAmount += flowAmount
	ud.user.BetAmount += betAmount
	ud.user.PayAmount += payAmount
	hhnn.UpdateUserScoreFunc(ud.user.UserID, income, tax)

	mlog.Debug("用户%d，实际赢钱：%v, 顯示下注赢钱：%v, 更新下注贏錢:%v", ud.user.UserID, money, ud.win, income)

	//这里写入赢家排名列表
	winSort(ud, hhnn)

}

//計算投注選項的輸贏金額
func calWinLossMoney(ud *hhnnuserdata, indexType int, v *JiangTypeWinLossOdd) int64 {
	return ud.userbet[indexType].money * int64(v.odds) * int64(v.winLoss)
}

//这里写入赢家排名列表
func winSort(ud *hhnnuserdata, HHNN *HHNNLogic) {
	badd := false
	if ud.win > 0 {
		if len(HHNN.winuserinfo) > 0 {
			for i := 0; i < len(HHNN.winuserinfo); i++ {
				u := HHNN.getUserDataByID(HHNN.winuserinfo[i])
				if u != nil {
					if ud.win > u.win {
						if len(HHNN.winuserinfo) >= 3 {
							if i == 0 {
								HHNN.winuserinfo[2] = HHNN.winuserinfo[1]
								HHNN.winuserinfo[1] = HHNN.winuserinfo[0]
								HHNN.winuserinfo[0] = ud.user.UserID
							} else if i == 1 {
								HHNN.winuserinfo[2] = HHNN.winuserinfo[1]
								HHNN.winuserinfo[1] = ud.user.UserID
							} else {
								HHNN.winuserinfo[2] = ud.user.UserID
							}
						} else if len(HHNN.winuserinfo) == 2 {
							HHNN.winuserinfo = append(HHNN.winuserinfo, ud.user.UserID)
							if i == 0 {
								HHNN.winuserinfo[2] = HHNN.winuserinfo[1]
								HHNN.winuserinfo[1] = HHNN.winuserinfo[0]
								HHNN.winuserinfo[0] = ud.user.UserID
							} else if i == 1 {
								HHNN.winuserinfo[2] = HHNN.winuserinfo[1]
								HHNN.winuserinfo[1] = ud.user.UserID
							}
						} else {
							HHNN.winuserinfo = append(HHNN.winuserinfo, ud.user.UserID)
							HHNN.winuserinfo[1] = HHNN.winuserinfo[0]
							HHNN.winuserinfo[0] = ud.user.UserID
						}

						badd = true
						break
					}
				}
			}
		}

		if len(HHNN.winuserinfo) < 3 && !badd {
			HHNN.winuserinfo = append(HHNN.winuserinfo, ud.user.UserID)
		}
	}
}

func (hhnn *HHNNLogic) OnRecvHTTPMessage(req *http.Request) proto.Message {
	action := req.PostFormValue("action")
	if action == "getserverinfo" {
		return hhnn.getTableListMsg()
	} else if action == "gettableinfo" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		return hhnn.getTableDetailMsg(int32(tid))
	} else if action == "modresult" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		ret := req.PostFormValue("result")
		adminuser := req.PostFormValue("adminuser")
		return hhnn.setNextJiangType(int32(tid), ret, adminuser)
	} else if action == "modpoolamount" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 64)
		typeid, _ := strconv.ParseInt(req.PostFormValue("type"), 10, 32)
		pool, _ := strconv.ParseInt(req.PostFormValue("poolamount"), 10, 64)
		operator := req.PostFormValue("operator")
		if typeid == 1 {
			return hhnn.setSysPool(int32(tid), pool, operator)
		} else {
			return hhnn.setFangFen(int32(tid), pool, operator)
		}
	} else if action == "modpopwaterrate" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		zhushui, _ := strconv.ParseInt(req.PostFormValue("popwaterrate"), 10, 32)
		return hhnn.setZhuShui(int32(tid), int32(zhushui))
	} else if action == "history" {
		pdata := new(netproto.Http_HHNN_ResultCardList)
		historyList := make([]*netproto.HHNN_GameResultPoker, 0)

		if len(hhnn.resultcardlist.CardList) > 20 {
			historyList = hhnn.resultcardlist.CardList[len(hhnn.resultcardlist.CardList)-20 : len(hhnn.resultcardlist.CardList)]
		} else {
			historyList = hhnn.resultcardlist.CardList
		}

		for _, v := range historyList {
			resultPoker := new(netproto.Http_HHNN_GameResultPoker)
			hungPoker := new(netproto.Http_HHNN_Poker)
			heiPoker := new(netproto.Http_HHNN_Poker)
			for _, hp := range v.HungPoker.Cards {
				hungPoker.Cards = append(hungPoker.Cards, &netproto.Http_PokerCard{CardNumber: hp.CardNumber, Color: hp.Color})
			}
			hungPoker.NiuNiuNumber = v.HungPoker.NiuNiuNumber
			hungPoker.NiuNiuRet = v.HungPoker.NiuNiuRet

			for _, hp := range v.HeiPoker.Cards {
				heiPoker.Cards = append(heiPoker.Cards, &netproto.Http_PokerCard{CardNumber: hp.CardNumber, Color: hp.Color})
			}
			heiPoker.NiuNiuNumber = v.HeiPoker.NiuNiuNumber
			heiPoker.NiuNiuRet = v.HeiPoker.NiuNiuRet

			resultPoker.HungPoker = hungPoker
			resultPoker.HeiPoker = heiPoker
			resultPoker.Win = v.Win
			resultPoker.RoundID = v.RoundID

			pdata.CardList = append(pdata.CardList, resultPoker)
		}

		return pdata
	}

	return nil
}

//获取桌子列表数据
func (hhnn *HHNNLogic) getTableListMsg() proto.Message {
	pdata := new(netproto.Http_HHNN_GameInfo)
	for _, table := range hhnn.GetTable().Room.Tables {
		pTable := new(netproto.Http_HHNN_TableList)
		logic := table.GameLogic.(*HHNNLogic)
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
		pTable.PoolCurrentAmount = proto.Int64(logic.syspool)
		lastturn := logic.lastturn
		if lastturn < 0 {
			lastturn = 0
		}
		pTable.SumEarnAmount = proto.Int64(logic.syspool)
		pTable.LastEarnAmount = proto.Int32(int32(logic.syslastwin))
		//pTable.SysTaxIncome = proto.Int64(logic.systaxincome)

		pdata.TableList = append(pdata.TableList, pTable)
	}

	return pdata
}

//获取桌子详细数据
func (hhnn *HHNNLogic) getTableDetailMsg(tableid int32) proto.Message {
	t := hhnn.GetTable().Room.GetTableByNo(tableid)
	pTable := new(netproto.Http_HHNN_TableDetail)

	if t != nil {

		logic := t.GameLogic.(*HHNNLogic)
		pTable.TableNo = proto.Int32(logic.TableNo())
		pTable.Status = proto.Int32(logic.gamestate)
		//		pTable.PoolInitAmount = proto.Int64(logic.startsyspool)
		pTable.PoolCurrentAmount = proto.Int64(logic.syspool)
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
		pTable.OpenResultHung = proto.Int64(logic.getRealUserWinByType(HUNG))
		pTable.OpenResultHei = proto.Int64(logic.getRealUserWinByType(HEI))
		//pTable.OpenResultDraw = proto.Int64(logic.getRealUserWinByType(NATURE))
		pTable.BetTime = proto.Int32(hhnn.getBetTime())

		if logic.gamestate == KaiJiangAction {
			pTable.CurentOpenResult = proto.String(logic.getResultString())
			pTable.CurentOpenResultEarn = proto.Int64(logic.getRealUserWinByType(logic.jiangtype))
		}
		pTable.UserInfo = logic.responseCurrentUserBetList()
		pTable.LastGameUserBet = logic.responseLastGameUserBetList()
		btwrate := logic.getBetTypeWinRate()
		pTable.HungWinRate = proto.Int32(btwrate[0])
		pTable.PowerWinRate = proto.Int32(btwrate[1])
		pTable.HeiWinRate = proto.Int32(btwrate[2])
		pTable.WinLine = proto.Int64(logic.winline)
		pTable.FanFengTarget = proto.Int64(logic.fangFenTarget)
		pTable.FanFengComplete = proto.Int64(logic.fangFenComplete)
		//mlog.Debug("===============================%d", pTable.GetHungWinRate())
	}

	return pTable
}

//收到RPC消息解析處理
func (hhnn *HHNNLogic) OnRecvRPCMessage(m *msg.Message) {
	rm, ok := m.MsgData.(*rpcmsg.RPCMessage)
	if !ok {
		return
	}
	mmbyte := rm.MsgData.([]byte)
	if len(mmbyte) == 0 {
		return
	}

	if m.SClassID == int32(netproto.GameMessageClassID_SuperToolGetAssignRewardDetailRetID) {
		hhnn.setAssignReward(m)
	}

}

//获取系统赢钱数
func (hhnn *HHNNLogic) getRealUserWinByType(t int32) int64 {
	win := int64(0)
	odd := int64(hhnn.betvalue.getOddByType(t))
	user := hhnn.GetUserByUserID(hhnn.zhuangID)
	if user != nil && !user.IsRobot() {
		_, tax := hhnn.GetUserWinAndTax(hhnn.bettotal)
		if hhnn.currtax != 0 {
			tax = hhnn.currtax
		}

		robotwin := odd * hhnn.getRobotBetMoneyByType(t)
		win = hhnn.bettotalrobot - robotwin - tax
		return -win
	}

	for _, ud := range hhnn.userdata {
		if !ud.user.IsRobot() {
			win += odd * ud.getUserBetByType(t)
			win += ud.money - ud.user.Score
		}
	}

	mlog.Debug("玩家赢钱%d", win)
	return -win
}

//获取此类型的机器人下注总额
func (hhnn *HHNNLogic) getRobotBetMoneyByType(t int32) int64 {
	for i := 0; i < len(hhnn.bettyperobot); i++ {
		if hhnn.bettyperobot[i].hhnntype == t {
			return hhnn.bettyperobot[i].money
		}
	}

	return 0
}

func (hhnn *HHNNLogic) getResultString() string {
	var res string
	res += resultstring[hhnn.jiangtype]
	return res
}

//设置下次开奖结果
func (hhnn *HHNNLogic) setNextJiangType(tableid int32, jiang string, adminuser string) proto.Message {
	t := hhnn.GetTable().Room.GetTableByNo(tableid)
	example.invalid("=======%s", jiang)
	if t != nil {
		var ret []int32
		logic := t.GameLogic.(*HHNNLogic)
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
			if ret[0] == 1 || ret[0] == 2 {
				logic.nextjiangtype = ret[0]
			}
			//validate ret[1]
			if len(ret) > 1 && ret[1] == 3 {
				logic.nextzuhetype = hhnn.getPowerJiangType()
			}
		}
		return hhnn.getRetMsg(true, "修改开奖结果成功")
	}

	return hhnn.getRetMsg(false, "修改开奖结果失败")
}

func (hhnn *HHNNLogic) getRetMsg(ret bool, msg string) proto.Message {
	pdata := new(netproto.Http_Result)
	pdata.Result = proto.Bool(ret)
	pdata.Message = proto.String(msg)
	return pdata
}

func (hhnn *HHNNLogic) setZhuShui(tableid int32, zhushui int32) proto.Message {
	if tableid == -1 {
		for _, table := range hhnn.GetTable().Room.Tables {
			logic := table.GameLogic.(*HHNNLogic)
			logic.changeZhuShui(zhushui)
		}
	} else {
		t := hhnn.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*HHNNLogic)
			logic.changeZhuShui(zhushui)
		} else {
			return hhnn.getRetMsg(false, "修改注水率失败")
		}
	}

	return hhnn.getRetMsg(true, "修改注水率成功")
}

//獲取放分組合類型
func (hhnn *HHNNLogic) getFangFenZuheType() (zuhetype []int32) {
	//放分開致命一擊的條件及機率
	alltype := []int32{WUXIAONIU, ZHADAN, WUHUANIU, KNIUNIU}
	rnd := util.RndInt(0, 1000)
	if rnd < 165 {
		r := (hhnn.fangFenTarget * 15 / 10) - hhnn.fangFenComplete
		for _, zuhe := range alltype {
			zuhepay := hhnn.getRealManBetMoneyByType(3) * int64(hhnn.betvalue.getOddByType(zuhe))
			if r > zuhepay {
				zuhetype = append(zuhetype, zuhe)
			}
		}
	}
	if len(zuhetype) < 1 {
		zuhetype = append(alltype, []int32{NIUNIU, NIUNINE, NIUEIGHT, NIUSEVEN, NIUSIX, NIUFIVE, NIUFOUR, NIUTHREE, NIUTWO, NIUONE, WUNIU}...)
	}
	return
}

//無真人下注
func (hhnn *HHNNLogic) noRealUserBet() bool {
	result := true

	for _, ud := range hhnn.userdata {
		if !ud.user.IsRobot() && len(ud.userbet) > 0 {
			return false
		}
	}

	return result
}
