package brnn

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
	"example.invalid/ServerTeam/Servers/pkg/mylog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

const (
	block1 = iota + 1 //下注區塊1
	block2            //下注區塊2
	block3
	block4
)

var resultstring = []string{"隨機", "通殺", "通陪", "贏一輸三", "贏二輸二", "贏三輸一"}

const (
	Bet            = iota //下注状态
	KaiJiang              //开奖状态
	KaiJiangAction        //开奖动画
	EndAction             //结算动画
	Wait                  //等待游戏开始
)

const (
	GameKaiJiang_Poker  = "kaijiangpoker"
	GameKaiJiang_Action = "kaijiang"
	GameEnd_Action      = "end"
)

//服务器状态
const (
	KaijianType_Normal             = iota //初始
	KaijianType_RobotZhuangWin            //機器人當莊贏
	KaijianType_RealUserZhuangLoss        //真人當莊, 莊輸
	KaijianType_RealUserZhuangWin         //真人当装, 庄赢
	KaijianType_RobotZhuangLoss           //莊輸
)

type jiangTypeEnum int

//開獎結果類型
const (
	naturl    jiangTypeEnum = iota //隨機
	tongsha                        //通殺
	tongpei                        //通陪
	win1loss3                      //贏一輸三
	win2loss2                      //贏二輸二
	win3loss1                      //贏三輸一
	assign                         //指定開獎
)

//动画名称
const (
	animationName_GameStart  = "gamestart"  //游戏开始
	animationName_FirstCard  = "firstcard"  //第一手牌
	animationName_SecondCard = "secondcard" //第二手牌
	animationName_GameRet    = "gameret"    //结算

	animationTimeOut_GameStart  = 3
	animationTimeOut_FirstCard  = 4
	animationTimeout_SecondCard = 3
	animationTimeout_GameRet    = 4
)

const (
	CtrlBetLimitDefault int64 = 100 // 監控名單預設注額
)

type userScore struct {
	userdata []*gamelogic.GameUser
}

type BRNNLogic struct {
	*gamelogic.ComputeGameLogicBase
	poker                   poker.Poker //游戏扑克牌
	serverfile              *os.File
	opttimer                *event.TimerFunc
	bettotal                int64                         //下注总额
	bettotalrobot           int64                         //机器人下注总额
	bettypetotal            *betTypeBlockGroup            //下注分类汇总
	bettyperobot            *betTypeBlockGroup            //机器人下注分类汇总
	betvalue                *betconfig                    //下注配置文件
	userdata                []*brnnuserdata               //用户数据
	brnnpokerresult         []*int32                      //开奖的牌
	gamestate               int32                         //游戏状态
	bettime                 int32                         //下注时间
	jiangtype               jiangTypeEnum                 //开奖类型
	action_str              string                        //当前动画
	kaijiangactiontime      int32                         //开奖动画超时
	secondfapeitimeinterval int32                         //第二手牌的時間間距. 0秒全發or其他秒數
	endactiontime           int32                         //用户结算动画超时
	readytime               int32                         //用户准备时间
	zhuangID                int32                         //当前庄家ID
	zhuanglimit             int64                         //上庄限制
	lastturn                int32                         //庄剩余的局数
	maxturn                 int32                         //上庄的最大局数
	zhuanglist              []int32                       //上庄列表
	betcount                []int32                       //下注数量
	syspool                 int64                         //当前水池的钱
	zhuangWin               int64                         //庄家赢钱
	serverid                string                        //服务器ID
	winuserinfo             []int32                       //赢钱最多的前三家userid
	resultcardlist          *netproto.BRNN_ResultCardList //开奖的历史牌
	zhushui                 int32                         //注水率
	maxclearcount           int32                         //清理数据上限
	maxrobotzhuang          int32                         //最多有好多个机器人上庄
	robotxiazhuang          int32                         //机器人下庄概率
	robotzhuang             int32                         //没有人当庄的时候机器人当庄的概率
	winline                 int64                         //机器人赢钱线大于这个值随机开。小于这个值机器人赢最多的开
	zhushuimoney            int64                         //注水的钱
	startsyspool            int64                         //开始时的金币数量
	startmoney              int64                         //上庄时的金币
	currtax                 int64                         //当前抽水
	nextjiangtype           jiangTypeEnum                 //下次开奖类型
	nextjiangitem           string                        //下次開獎
	syslastwin              int64                         //上一局收益
	allKindPoker            map[int]*poker.CardGroup      //牌:莊*1和閒*4
	secondPokerGroup        map[int]*poker.CardGroup      //第二手牌
	systaxincome            int64                         //系統稅收, 非該局稅收
	userWinLose             int64                         //玩家輸贏
	robotWinLose            int64                         //機器人輸贏
	jiangTypeWinLossOdd     []*JiangTypeWinLossOdd        //賠率, 贏的參考閒家的牛點數當賠率　輸參考莊的牌賠率
	betlimitmoney           int64                         //下注限制
	publiczhuang            int64                         //是否開放用戶上庄
	userscore               *userScore                    //用户数据
	luckystar               int32                         //幸運星玩家id
	luckyfaceid             int32                         //幸運星玩家頭像id
	luckystarround          int64                         //幸運星持續局數
	luckystartimes          int64                         //幸運星第幾局
	notbettimes             int64                         //沒下注持續局數
	richestuserinfo         []int32                       //富豪玩家的userid
	richestList             []*netproto.BRNN_Richest      //富豪榜
	stageuserinfo           []int32                       //標誌玩家的userid
	lastGameUserBet         []*brnnuserdata               //上局下注玩家
	//bowProvider             *gamelogic.BlackOrWhiteUserProvider   //需要放分或殺分的列表
	assignRewardDetail *netproto.SuperToolAssignRewardDetail //指定開獎
	execAssignReward   bool                                  //true 執行指定開獎
	fangFenTarget      int64                                 //放分金額
	fangFenComplete    int64                                 //己放分金額
	execFangFen        bool                                  //true 執行系統放分
	SuperToolService   *gamelogic.SuperToolBrGameServiceBase
	kaijiangbiao       jiangTypeEnum //哪張開獎表
	zhunagPoint        int           //莊家牌點
	xianPoint          []int         //閒家牌點
	ctrlbetlimit       int64         //監控名單會控制的最低注額
}

type JiangTypeWinLossOdd struct {
	betType int32 //下注種類選項
	winLoss int32 //此項和莊家比牌後輸贏結果　-1輸　1贏 , 負數是賠付倍數方便計算
	odds    int32 //賠率, 贏的參考閒家的牛點數當賠率　輸參考莊的牌賠率
}

//新建一个百人牛牛实例
func NewBRNNLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *BRNNLogic {
	glc := new(BRNNLogic)
	glc.ComputeGameLogicBase = gamelogic.NewComputeGameLogicBase(t, tfManager, cfgs)
	glc.getbrnnconfig(cfgs)
	glc.SuperToolService = glc.GetSuperToolService().(*gamelogic.SuperToolBrGameServiceBase)
	glc.gamestate = Wait
	glc.zhuangID = 0
	glc.syslastwin = 0
	glc.nextjiangtype = naturl
	glc.resultcardlist = new(netproto.BRNN_ResultCardList)
	glc.getLastPool()
	glc.startsyspool = glc.syspool
	glc.addCmd()
	glc.RegisterBeforeGameStartFunc(func() {
		glc.clearGameData() //清除数据
	})

	glc.RegisterAfterGameStartFunc(func() {
		//glc.fangFenTarget = 100000
		glc.betvalue = newBetconfig(glc.Configs)

		//随机发牌
		if glc.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
			glc.poker = poker.NewRandomPoker(false, false, 1) //初始化一副牌
			glc.poker.Shuffle()
		} else { //发牌器发牌
			glc.poker = poker.NewFilePoker(glc.Configs, glc.AllTableUser())
			glc.poker = poker.NewFixedFilePoker(glc.poker.(*poker.FilePoker))
		}

		glc.getAssignReward()
		//glc.getTableBlackUser()
		//glc.getTableWhiteUser()
		//glc.VipUserService.UpdateVIPUser()
		glc.huangZhuang() //7/9 上莊機制配置
		glc.getLuckyStar()
		glc.getRichestList()
		glc.GameLog.AddGameStep("庄家%d，剩余%d", glc.zhuangID, glc.lastturn)
		glc.zhuShui()
		glc.broadCastZhuangInfo()
		glc.BroadCastAbleBetTotal()
		glc.startFirstCard()
		glc.startBet()
		mlog.Debug("九局四門勝率 : %v", glc.getBetTypeWinRate(9))
	})

	glc.RegisterBeforeGameEndFunc(func() {
		glc.saveLastGameUserBet()
		//glc.updateBowUserWin()
	})

	return glc
}

//开始下注
func (brnn *BRNNLogic) startBet() {
	brnn.opttimer = brnn.TimerFuncManager.NewTimer(3, func() {
		brnn.setState(Bet)
		brnn.BoardcastMessage(brnn.getStateMsg(true))
		brnn.BroadCastAbleBetTotal()
		//brnn.intervalBetTotolInfo()
	})

	brnn.opttimer = brnn.TimerFuncManager.NewTimer(brnn.bettime, brnn.kaiJiang)
}

//获取配置表
func (brnn *BRNNLogic) getbrnnconfig(cfgs *config.Group) {
	brnn.betvalue = newBetconfig(cfgs)
	scfgs := cfgs.GetConfig("gameserver")
	brnn.bettime = int32(scfgs.GetConfigInt("gameserver", "bettime"))
	brnn.kaijiangactiontime = int32(scfgs.GetConfigInt("gameserver", "kaijiangactiontime"))
	brnn.endactiontime = int32(scfgs.GetConfigInt("gameserver", "endactiontime"))
	brnn.readytime = int32(scfgs.GetConfigInt("gameserver", "readywaittime"))
	brnn.zhuanglimit = int64(scfgs.GetConfigInt("gameserver", "zhuanglimit"))
	brnn.maxturn = int32(scfgs.GetConfigInt("gameserver", "zhuangturn"))
	brnn.serverid = strconv.Itoa(int(brnn.GetServerConfig().ServerID))
	brnn.maxclearcount = int32(scfgs.GetConfigInt("gameserver", "maxclearcount"))
	brnn.zhushui = int32(scfgs.GetConfigInt("gameserver", "zhushui"))
	brnn.maxrobotzhuang = int32(scfgs.GetConfigInt("gameserver", "maxrobotzhuang"))
	brnn.robotxiazhuang = int32(scfgs.GetConfigInt("gameserver", "robotxiazhuang"))
	brnn.robotzhuang = int32(scfgs.GetConfigInt("gameserver", "robotzhuang"))
	brnn.winline = int64(scfgs.GetConfigInt("gameserver", "winline"))
	betcountstr := scfgs.GetConfigValue("gameserver", "betcount")
	brnn.betlimitmoney = int64(scfgs.GetConfigInt("gameserver", "betlimitmoney"))
	brnn.publiczhuang = int64(scfgs.GetConfigInt("gameserver", "publiczhuang"))
	brnn.luckystarround = int64(scfgs.GetConfigInt("gameserver", "luckystarround"))

	betarr := strings.Split(betcountstr, ",")
	for _, str := range betarr {
		count, _ := strconv.ParseInt(str, 10, 32)
		brnn.betcount = append(brnn.betcount, int32(count))
	}

	//client setting
	cgf := cfgs.GetConfig("logicconfig")
	brnn.secondfapeitimeinterval = int32(cgf.GetConfigInt("client", "secondfapeitimeinterval"))

	brnn.ctrlbetlimit = int64(scfgs.GetConfigInt("supertool", "ctrlbetlimit"))
	if brnn.ctrlbetlimit <= 0 {
		brnn.ctrlbetlimit = CtrlBetLimitDefault
	}
}

//改变水池注水率
func (brnn *BRNNLogic) changeZhuShui(value int32) {
	if value >= 0 && value <= 100 {
		brnn.zhushui = value
	}
}

//获取上一次水池的钱
func (brnn *BRNNLogic) getLastPool() {
	var folderName = "lastpool"
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err = os.MkdirAll(folderName, os.ModePerm)
	}

	var tablenostr = strconv.Itoa(int(brnn.TableNo()))
	var filename = brnn.serverid + "_" + tablenostr + ".conf"
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
		brnn.serverfile = f
		defer brnn.serverfile.Close()
		brnn.readPoolFromFile()
	}
}

//读取水池文件
func (brnn *BRNNLogic) readPoolFromFile() {
	if conf.TestMode {
		return
	}

	buff := make([]byte, 1024)
	n, _ := brnn.serverfile.Read(buff)
	if n > 0 {
		poolstr := string(buff)[0:n]
		//pool, err := strconv.ParseInt(poolstr, 10, 32)
		pool, _ := strconv.Atoi(poolstr)
		brnn.syspool = int64(pool)
	}
}

//添加命令行
func (brnn *BRNNLogic) addCmd() {
	command.AddCommand(NewBRNNCmd(brnn))
}

//清除数据
func (brnn *BRNNLogic) clearGameData() {
	brnn.userdata = make([]*brnnuserdata, 0)
	brnn.bettypetotal = new(betTypeBlockGroup)
	brnn.bettyperobot = new(betTypeBlockGroup)
	brnn.brnnpokerresult = make([]*int32, 0)
	brnn.winuserinfo = make([]int32, 0)
	brnn.stageuserinfo = make([]int32, 0)
	brnn.richestuserinfo = make([]int32, 0)
	brnn.richestList = make([]*netproto.BRNN_Richest, 0)
	brnn.jiangtype = naturl
	brnn.kaijiangbiao = naturl
	brnn.zhunagPoint = 0
	brnn.xianPoint = make([]int, 0)
	brnn.bettotal = 0
	brnn.bettotalrobot = 0
	brnn.zhuangWin = 0
	brnn.currtax = 0
	brnn.allKindPoker = make(map[int]*poker.CardGroup)
	brnn.secondPokerGroup = make(map[int]*poker.CardGroup)
	brnn.userWinLose = 0
	brnn.robotWinLose = 0
	brnn.jiangTypeWinLossOdd = make([]*JiangTypeWinLossOdd, 0)
	brnn.nextjiangitem = ""
	brnn.assignRewardDetail = nil
	brnn.execAssignReward = false
	brnn.execFangFen = false
	rand.Seed(time.Now().UnixNano())
	zlen := len(brnn.zhuanglist)
	if zlen == 0 {
		brnn.robotZhuang()
	}

	brnn.robotXiaZhuang()
}

//机器人上庄
func (brnn *BRNNLogic) robotZhuang() {
	if brnn.publiczhuang == 1 {
		//概率当庄
		if rand.Int31n(100) > brnn.robotzhuang {
			return
		}
		//统计机器人个数
		robotarr := make([]int32, 0)
		for _, ur := range brnn.AllTableUser() {
			if ur.IsRobot() {
				robotarr = append(robotarr, ur.UserID)
			}
		}

		if len(robotarr) <= 0 {
			return
		}

		//随机有几个人当庄
		n := rand.Int31n(brnn.maxrobotzhuang) + 1

		for i := int32(0); i < n; i++ {
			rcn := len(robotarr)
			if rcn <= 0 {
				break
			}
			index := rand.Int31n(int32(rcn))
			brnn.addToZhuangList(robotarr[index])
			robotarr = append(robotarr[:index], robotarr[index+1:]...)
		}
	} else {
		mlog.Debug("未開啟上庄機制，機器人不用上庄")
		return
	}
}

//机器人下庄
func (brnn *BRNNLogic) robotXiaZhuang() {
	if brnn.zhuangID == 0 {
		return
	}

	ur := brnn.GetUserByUserID(brnn.zhuangID)
	if ur != nil && ur.IsRobot() {
		if rand.Int31n(100) < brnn.robotxiazhuang {
			brnn.zhuangID = 0
		}
	}
}

//添加一个用户到上庄列表
func (brnn *BRNNLogic) addToZhuangList(userid int32) {
	brnn.zhuanglist = append(brnn.zhuanglist, userid)
}

//换庄
func (brnn *BRNNLogic) huangZhuang() {
	if brnn.publiczhuang == 1 {
		if brnn.zhuangID == 0 {
			//庄是系统直接换庄
			brnn.setNextZhuang()
			return
		}
		if brnn.lastturn <= 0 {
			//已经到最后一局了
			brnn.setNextZhuang()
			return
		}
		ur := brnn.GetUserByUserID(brnn.zhuangID)
		if ur != nil && ur.Score < brnn.zhuanglimit {
			brnn.setNextZhuang()
			return
		}

		if ur == nil {
			brnn.setNextZhuang()
			return
		}

		brnn.lastturn--
		mlog.Debug("1015发送庄家剩余局数%d", brnn.lastturn)
	} else {
		mlog.Debug("未開啟上庄機制，不用換庄")
		return
	}
}

//设置下一个上庄的玩家
func (brnn *BRNNLogic) setNextZhuang() {
	if brnn.publiczhuang == 1 {
		brnn.checkNetUserCanShangZhuang()
		if len(brnn.zhuanglist) > 0 {
			brnn.zhuangID = brnn.zhuanglist[0]
			brnn.lastturn = brnn.maxturn
			brnn.lastturn--
			brnn.delFromZhuangList(brnn.zhuangID)
			user := brnn.GetUserByUserID(brnn.zhuangID)
			brnn.startmoney = user.Score
		} else {
			brnn.zhuangID = 0
			brnn.startmoney = brnn.startsyspool
		}
	} else {
		mlog.Error("未開啟上庄機制，不得上庄")
		return
	}
}

//检测下一个用户是否能上庄
func (brnn *BRNNLogic) checkNetUserCanShangZhuang() {
	if len(brnn.zhuanglist) > 0 {
		for {
			user := brnn.GetUserByUserID(brnn.zhuanglist[0])
			if user != nil && user.Score >= int64(brnn.zhuanglimit) {
				break
			}

			brnn.delFromZhuangList(brnn.zhuanglist[0])
			if len(brnn.zhuanglist) <= 0 {
				break
			}
		}
	}
}

//删除一个用户从上庄列表
func (brnn *BRNNLogic) delFromZhuangList(userid int32) {
	for i := 0; i < len(brnn.zhuanglist); i++ {
		if brnn.zhuanglist[i] == userid {
			brnn.zhuanglist = append(brnn.zhuanglist[:i], brnn.zhuanglist[i+1:]...)
			break
		}
	}
}

//注水到下注额
func (brnn *BRNNLogic) zhuShui() {
	if brnn.zhuangID == 0 {
		brnn.bettotal = brnn.syspool * int64(brnn.zhushui) / 100
		brnn.syspool -= brnn.bettotal
	} else {
		ur := brnn.GetUserByUserID(brnn.zhuangID)
		if ur != nil && ur.IsRobot() {
			brnn.bettotal = brnn.syspool * int64(brnn.zhushui) / 100
			brnn.syspool -= brnn.bettotal
		}
	}

	brnn.zhushuimoney = brnn.bettotal
	brnn.GameLog.AddGameStep("系統水池%d, 系統注水%d", brnn.syspool+brnn.zhushuimoney, brnn.zhushuimoney)
}

//广播庄家
func (brnn *BRNNLogic) broadCastZhuangInfo() {
	brnn.BoardcastMessage(brnn.getZhuangInfoMsg())
	brnn.BoardcastMessage(brnn.getRichestListMsg())
	brnn.BoardcastMessage(brnn.getLuckyStarMsg())
}

//廣播已下注和還可下注
func (brnn *BRNNLogic) BroadCastAbleBetTotal() {
	brnn.BoardcastMessage(brnn.getAbleBetTotalInfo())
}

//取得還可下注和已下注資訊
func (brnn *BRNNLogic) getAbleBetTotalInfo() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNBroadCastTotalBet)

	score := int64(0)

	if brnn.zhuangID == int32(0) { //系統當莊, 可下注總額
		score = int64(brnn.betvalue.betvalue[0].maxmoney +
			brnn.betvalue.betvalue[1].maxmoney +
			brnn.betvalue.betvalue[2].maxmoney +
			brnn.betvalue.betvalue[3].maxmoney)
	} else {
		ur := brnn.GetUserByUserID(brnn.zhuangID)
		if ur != nil {
			score = ur.Score
		}
	}

	pdata := new(netproto.BRNN_BroadCastTotal)
	pdata.AlreadyTotalBet = proto.Int64(brnn.bettotal + brnn.bettotalrobot)

	ableTotal := int64((int64(score/10)-*pdata.AlreadyTotalBet)/100) * 100
	pdata.AbleTotalBet = proto.Int64(ableTotal)

	pmsg.MsgData = pdata
	return pmsg
}

//获取庄家信息消息
func (brnn *BRNNLogic) getZhuangInfoMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNZhuangInfoID)

	pdata := new(netproto.BRNN_ZhuangInfo)
	pdata.UserID = proto.Int32(0)
	pdata.LastTurn = proto.Int32(0)
	ur := brnn.GetUserByUserID(brnn.zhuangID)
	if ur != nil {
		pdata.UserInfo = append(pdata.UserInfo, ur.GetProtoBufGameUser())
		pdata.UserID = proto.Int32(ur.UserID)
		pdata.LastTurn = proto.Int32(brnn.lastturn)
	}
	pmsg.MsgData = pdata
	mlog.Debug("莊家資訊 : %v , %v", *pdata.UserID, pdata.UserInfo)

	return pmsg
}

func (brnn *BRNNLogic) sendToUserFirstCardGroup(gameuser *gamelogic.GameUser) {
	brnn.SendMessageToUser(gameuser, brnn.getFirstCardGroup())
}

//第一手三張牌
func (brnn *BRNNLogic) getFirstCardGroup() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNBroadCastFapaiThreeCard)

	pdata := new(netproto.BRNN_BrodcastPokerCard)

	for i := int32(0); i <= int32(len(brnn.betvalue.alltype)); i++ {
		cg := brnn.getCalFirstPokerGroup(i)

		urPokerCard := new(netproto.BRNN_BrodcastUserPokerCard)
		urPokerCard.SitNo = proto.Int32(i)

		for _, cd := range cg.Cards {
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
func (brnn *BRNNLogic) getCalFirstPokerGroup(position int32) *poker.CardGroup {
	return brnn.allKindPoker[int(position)]
}

//廣播第二手兩張牌
func (brnn *BRNNLogic) boardcastSecondCardGroup() {
	brnn.BoardcastMessage(brnn.getSecondCardGroup())
}

//第二手兩張牌
func (brnn *BRNNLogic) getSecondCardGroup() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNBroadCastFapaiTwoCard)

	pdata := new(netproto.BRNN_BrodcastPokerCard)

	for i := int32(0); i <= int32(len(brnn.betvalue.alltype)); i++ {
		cg := brnn.getCalSecondPokerGroup(i)
		urPokerCard := new(netproto.BRNN_BrodcastUserPokerCard)
		urPokerCard.SitNo = proto.Int32(i)

		urCard := new(poker.CardGroup)

		if cg == nil || cg.Cards == nil {
			mlog.Debug("index : [%v] 沒取到第二張牌", i)
			break
		}

		for _, cd := range cg.Cards {
			card := new(netproto.PokerCard)
			card.Color = proto.Int32(cd.Color)
			card.CardNumber = proto.Int32(cd.CardNumber)
			urPokerCard.CardGroup = append(urPokerCard.CardGroup, card)

			pCard := new(poker.Card)
			pCard.Color = cd.Color
			pCard.CardNumber = cd.CardNumber
			urCard.Cards = append(urCard.Cards, pCard)
		}

		pdata.UserPokerCard = append(pdata.UserPokerCard, urPokerCard)
		brnn.allKindPoker[int(i)].Cards = append(brnn.allKindPoker[int(i)].Cards, urCard.Cards...)

	}
	pmsg.MsgData = pdata

	return pmsg
}

//依照開獎的牌
func (brnn *BRNNLogic) getKaiJiangPokerByType() {
	switch brnn.jiangtype {
	case naturl: //隨機開獎
		brnn.natureResultType()
	case tongsha: //通殺
		mlog.Debug("通殺")
		brnn.tongshaResultType(10)
	case tongpei:
		mlog.Debug("通賠")
		brnn.tongpeiResultType(10)
	case win1loss3:
		mlog.Debug("贏一輸三")
		brnn.win1loss3ResultType(10)
	case win2loss2:
		mlog.Debug("贏二輸二")
		brnn.win2loss2ResultType(10)
	case win3loss1:
		mlog.Debug("贏三輸一")
		brnn.win3loss1ResultType(10)
	case assign:
		mlog.Debug("指定開獎")
		brnn.assignKajianResultType(10)
	}

}

//計算每家的第二張牌(含計算賺賠機率的牌型)
func (brnn *BRNNLogic) getResultType() {
	superToolHandler := brnn.SuperToolService.GetBrGameSuperToolHandler()
	//指定开奖
	if brnn.nextjiangtype > naturl {
		brnn.jiangtype = brnn.nextjiangtype
		brnn.nextjiangtype = naturl
		mlog.Debug("指定開獎 : %v", brnn.jiangtype)
		brnn.GameLog.AddGameStep("指定開獎 : %v", brnn.jiangtype)
		return
	} else if ok, rettype := brnn.isAssignReward(); ok {
		example.invalid("九局指定開獎 - %d : %v", brnn.assignRewardDetail.GetIdx()+1, rettype)
		brnn.GameLog.AddGameStep("九局指定開獎 - %d : %v", brnn.assignRewardDetail.GetIdx()+1, rettype)
		brnn.execAssignReward = true
		brnn.jiangtype = assign
		for i, v := range rettype {
			if v == 1 {
				brnn.nextjiangitem += fmt.Sprintf("%d,", i+1)
			}
		}
		if len(brnn.nextjiangitem) > 0 {
			brnn.nextjiangitem = brnn.nextjiangitem[:len(brnn.nextjiangitem)-1]
		}

		return
	}
	//系統放分
	if brnn.fangFenTarget > 0 && brnn.fangFenComplete < brnn.fangFenTarget {
		brnn.execFangFen = true
		brnn.jiangtype = assign
		items := []int{0, 1, 2, 3}
		items = util.IntArrayShuffle(items)[:len(items)-1]
		for _, i := range items {
			brnn.nextjiangitem += fmt.Sprintf("%d,", i+1)
		}
		brnn.nextjiangitem = brnn.nextjiangitem[:len(brnn.nextjiangitem)-1]
		brnn.GameLog.AddGameStep("系統放分 : target %d, already %d", brnn.fangFenTarget, brnn.fangFenComplete)
		//mlog.Debug("系統放分 : target %d, already %d", brnn.fangFenTarget, brnn.fangFenComplete)
		return
	}
	t := KaijianType_Normal

	rand.Seed(time.Now().Unix())
	//監控放殺分
	if len(brnn.getBowBet()) > 0 { //監控放殺分
		if brnn.bowJiangType() {
			example.invalid("執行監控放殺分策略 : %s", superToolHandler.GetProiorityListString())
			brnn.GameLog.AddGameStep("執行監控放殺分策略 : %s", superToolHandler.GetProiorityListString())
			return
		}
	}

	if t == KaijianType_Normal {
		if brnn.noRealUsreBet() {
			t = KaijianType_RobotZhuangLoss
			if !brnn.GetUserByUserID(brnn.zhuangID).IsRobot() {

				if brnn.syspool < brnn.winline {
					//真人當莊, 吃莊
					t = KaijianType_RealUserZhuangLoss
				} else {
					//真人當莊, 庄赢
					t = KaijianType_RealUserZhuangWin
				}
			}
		} else if brnn.zhuangID != 0 && brnn.GetUserByUserID(brnn.zhuangID).IsRobot() &&
			brnn.syspool < brnn.winline {
			//機器人當莊
			t = KaijianType_RobotZhuangWin //是否要作弊
			mlog.Debug("機器人當莊")
		} else if !brnn.GetUserByUserID(brnn.zhuangID).IsRobot() {

			if brnn.syspool < brnn.winline {
				//真人當莊, 吃莊
				t = KaijianType_RealUserZhuangLoss
			} else {
				//真人當莊, 庄赢
				t = KaijianType_RealUserZhuangWin
			}
		} else {
			if brnn.getRealUserAllBet(t) {
				t = KaijianType_RobotZhuangWin
			} else {
				t = KaijianType_RobotZhuangLoss
			}
		}
		if brnn.zhuangID == 0 && brnn.syspool < brnn.winline {
			t = KaijianType_RobotZhuangWin
		}
	}

	if t == KaijianType_Normal { //所有都可開獎
		brnn.kaijiangbiao = KaijianType_Normal
		brnn.jiangtype = naturl
	} else if t == KaijianType_RobotZhuangWin { //作弊 机器人和系统当庄的时候方案一
		brnn.kaijiangbiao = KaijianType_RobotZhuangWin
		rand.Seed(time.Now().Unix())
		randNum := rand.Int31n(100)
		sum := int32(0)

		for k := 0; k <= 5; k++ {
			kaijanKind := brnn.betvalue.jiangWeight
			sum += kaijanKind[k].kaijianWeight
			if randNum <= sum {
				brnn.jiangtype = jiangTypeEnum(k)
				mlog.Debug("机器人和系统当庄 低於水池 開獎:%v", brnn.jiangtype)
				break
			}
		}

	} else if t == KaijianType_RealUserZhuangLoss {
		brnn.kaijiangbiao = KaijianType_RealUserZhuangLoss
		//真人當莊
		rand.Seed(time.Now().Unix())
		randNum := rand.Int31n(100)
		sum := int32(0)
		for k := 0; k <= 5; k++ {
			kaijanKind := brnn.betvalue.realUserZhuangSyspoolJiangWeight
			sum += kaijanKind[k].kaijianWeight
			if randNum <= sum {
				brnn.jiangtype = jiangTypeEnum(k)
				brnn.GameLog.AddGameStep("真人當莊 低於水池 開獎:%v", brnn.jiangtype)
				example.invalid("真人當莊 莊輸 低於水池 開獎:%v v.kaijianWeight[%v]", brnn.jiangtype, kaijanKind[k].kaijianWeight)
				break
			}
		}
		if brnn.jiangtype == naturl {
			mlog.Debug("未作弊")
		}
	} else if t == KaijianType_RealUserZhuangWin {
		brnn.kaijiangbiao = KaijianType_RealUserZhuangWin
		//真人當莊
		rand.Seed(time.Now().Unix())
		randNum := rand.Int31n(100)
		sum := int32(0)

		for k := 0; k <= 5; k++ {
			kaijanKind := brnn.betvalue.realUserZhuangWinlineJiangWeight
			sum += kaijanKind[k].kaijianWeight

			if randNum <= sum {
				brnn.jiangtype = jiangTypeEnum(k)
				brnn.GameLog.AddGameStep("真人當莊 低於水池 開獎:%v", brnn.jiangtype)
				example.invalid("真人當莊 莊輸 低於水池 開獎:%v v.kaijianWeight[%v]", brnn.jiangtype, kaijanKind[k].kaijianWeight)
				break
			}
		}
		if brnn.jiangtype == naturl {
			mlog.Debug("未作弊")
		}
	} else if t == KaijianType_RobotZhuangLoss {
		brnn.kaijiangbiao = KaijianType_RobotZhuangLoss
		rand.Seed(time.Now().Unix())
		randNum := rand.Int31n(100)
		sum := int32(0)
		for k := 0; k <= 5; k++ {
			kaijanKind := brnn.betvalue.robotZhuangLoss
			sum += kaijanKind[k].kaijianWeight
			if randNum <= sum {
				brnn.jiangtype = jiangTypeEnum(k)
				brnn.GameLog.AddGameStep("高於水池 開獎:%v", brnn.jiangtype)
				example.invalid("高於水池 開獎:%v v.kaijianWeight[%v], randNum[%v]", brnn.jiangtype, kaijanKind[k].kaijianWeight, randNum)
				break
			}
		}
		if brnn.jiangtype == naturl {
			mlog.Debug("未作弊")
		}

	} else {
		brnn.jiangtype = naturl
	}

	mlog.Debug("开奖类型：%v", brnn.jiangtype)
}

//無真人下注
func (brnn *BRNNLogic) noRealUsreBet() bool {
	result := true

	for _, ud := range brnn.userdata {
		if !brnn.GetUserByUserID(ud.user.UserID).IsRobot() && len(ud.userbet) > 0 {
			return false
		}
	}

	return result
}

//四門全下吃水線下的機率表
func (brnn *BRNNLogic) getRealUserAllBet(t int) bool {
	result := false

	rand.Seed(time.Now().Unix())
	randExec := rand.Int31n(100)
	if randExec > brnn.betvalue.allbetsyspool {
		return false
	}

	if t == 1 {
		return true
	}

	for _, ud := range brnn.userdata {
		if !brnn.GetUserByUserID(ud.user.UserID).IsRobot() && len(ud.userbet) == int(4) {
			return true
		}
	}

	return result
}

//隨機自然發牌
func (brnn *BRNNLogic) natureResultType() {
	for i := int32(0); i <= int32(len(brnn.betvalue.alltype)); i++ {
		cg := brnn.poker.Deal(0, 2)

		urCard := new(poker.CardGroup)

		for _, cd := range cg.Cards {
			pCard := new(poker.Card)
			pCard.Color = cd.Color
			pCard.CardNumber = cd.CardNumber
			urCard.Cards = append(urCard.Cards, pCard)
		}

		brnn.secondPokerGroup[int(i)] = urCard
	}
}

//通殺牌,莊家 機率 牛七以上開獎 或 牛七以下的牛牛, 閒家小於莊
func (brnn *BRNNLogic) tongshaResultType(runtime int32) {
	brnn.zhuangTongshaCard(runtime)
	brnn.xianTongshaCard(runtime)
}

//莊通賠
func (brnn *BRNNLogic) tongpeiResultType(runtime int32) {
	brnn.secondPokerGroup = make(map[int]*poker.CardGroup, 0)

	blockItem := brnn.assignBlockLossItem(0)

	lossNumber := len(brnn.betvalue.alltype) - len(blockItem)
	ok := brnn.checkAssignWinLoss(lossNumber)
	if !ok { //重新組牌
		if runtime <= int32(0) {
			return
		}
		for k, _ := range brnn.secondPokerGroup {
			brnn.poker.AddCardGroup(0, brnn.secondPokerGroup[k]) //把牌放回去
		}
		if runtime > int32(0) {
			runtime--
			brnn.win1loss3ResultType(runtime)
		}
	}
}

//贏一輸三
func (brnn *BRNNLogic) win1loss3ResultType(runtime int32) {
	brnn.secondPokerGroup = make(map[int]*poker.CardGroup, 0)

	blockItem := brnn.assignBlockLossItem(1)

	lossNumber := len(brnn.betvalue.alltype) - len(blockItem) - 1
	ok := brnn.checkAssignWinLoss(lossNumber)
	if !ok { //重新組牌
		if runtime <= int32(0) {
			return
		}
		for k, _ := range brnn.secondPokerGroup {
			brnn.poker.AddCardGroup(0, brnn.secondPokerGroup[k]) //把牌放回去
		}
		if runtime > int32(0) {
			runtime--
			brnn.win1loss3ResultType(runtime)
		}
	}
}

//指定開獎項目
func (brnn *BRNNLogic) assignKajianResultType(runtime int32) {
	brnn.secondPokerGroup = make(map[int]*poker.CardGroup, 0)
	//莊
	//隨機取莊的點數之上
	rand.Seed(time.Now().Unix())
	randNum := util.RndInt(3, 5)
	niuCardRet := brnn.assignNewNiuNiuCardNumber(brnn.allKindPoker[0], randNum)
	//niuCardRet := brnn.assignNewNiuNiuCardMoreThen(brnn.allKindPoker[0], niuRet, randNum)
	brnn.appendSecondPokerGroup(0, niuCardRet)

	//閒
	var blockItem []int32 //要開出下注金額大到小的第幾個項目
	nextItemList := strings.Split(brnn.nextjiangitem, ",")
	for i := 1; i < len(brnn.allKindPoker); i++ {
		isLose := true
		for _, v := range nextItemList {
			if v == "" {
				continue
			}
			item, _ := strconv.Atoi(v)
			if i == item {
				mlog.Debug("修改開獎結果 win :%v", item)
				isLose = false
			}
		}

		if isLose {
			mlog.Debug("修改開獎結果 lose :%v", i)
			blockItem = append(blockItem, int32(i))
		}
	}

	mlog.Debug("指定開獎block : %v 輸", blockItem)
	brnn.assignBlockLoss(blockItem, niuCardRet.number)
	brnn.assignBlockWin(true)
}

//贏三輸一
func (brnn *BRNNLogic) win3loss1ResultType(runtime int32) {
	brnn.secondPokerGroup = make(map[int]*poker.CardGroup, 0)

	blockItem := brnn.assignBlockLossItem(3)

	lossNumber := len(brnn.betvalue.alltype) - len(blockItem) - 3
	ok := brnn.checkAssignWinLoss(lossNumber)
	if !ok { //重新組牌
		if runtime <= int32(0) {
			return
		}
		for k, _ := range brnn.secondPokerGroup {
			brnn.poker.AddCardGroup(0, brnn.secondPokerGroup[k]) //把牌放回去
		}
		if runtime > int32(0) {
			runtime--
			brnn.win3loss1ResultType(runtime)
		}
	}
}

//贏二輸二
func (brnn *BRNNLogic) win2loss2ResultType(runtime int32) {
	brnn.secondPokerGroup = make(map[int]*poker.CardGroup, 0)

	blockItem := brnn.assignBlockLossItem(2)

	lossNumber := len(brnn.betvalue.alltype) - len(blockItem) - 2
	ok := brnn.checkAssignWinLoss(lossNumber)
	if !ok { //重新組牌
		if runtime <= int32(0) {
			return
		}
		for k, _ := range brnn.secondPokerGroup {
			brnn.poker.AddCardGroup(0, brnn.secondPokerGroup[k]) //把牌放回去
		}
		if runtime > int32(0) {
			runtime--
			brnn.win2loss2ResultType(runtime)
		}

	}
}

//指定開幾個輸
func (brnn *BRNNLogic) assignBlockLossItem(lossCount int) []int32 {
	randNumWinLoss := rand.Int31n(100)
	//莊
	niuRet := nnCard_Niu
	//隨機取莊的點數之上
	randNum, isTakeMore := brnn.getZhungWinCardNumber(randNumWinLoss)
	niuCardRet := new(niuniuCardRet)
	if isTakeMore {
		niuCardRet = brnn.assignNewNiuNiuCardMoreThen(brnn.allKindPoker[0], niuRet, randNum) //指定莊牌
	} else {
		niuCardRet = brnn.assignNewNiuNiuCardLessThen(brnn.allKindPoker[0], niuRet, randNum, 10)
	}

	brnn.appendSecondPokerGroup(0, niuCardRet)
	//閒
	brnn.getBlockOrderBy()

	blockItem := brnn.getKajianTopBlock(lossCount, randNumWinLoss) //選取要開獎的下注金額大到小的第幾個
	mlog.Debug("指定開獎block : %v 輸", blockItem)
	brnn.assignBlockLoss(blockItem, niuCardRet.number)
	brnn.assignBlockWin(isTakeMore)
	return blockItem
}

func (brnn *BRNNLogic) getZhungWinCardNumber(randNum int32) (int, bool) {
	sum := int32(0)

	sum += brnn.betvalue.jiangWeight[int(brnn.jiangtype)].bigwin
	if randNum <= sum {
		result := brnn.getWinLossAssignCardNumber("bigwin")
		return int(result), true
	}

	sum += brnn.betvalue.jiangWeight[int(brnn.jiangtype)].smallwin
	if randNum <= sum {
		result := brnn.getWinLossAssignCardNumber("smallwin")
		return int(result), true
	}

	sum += brnn.betvalue.jiangWeight[int(brnn.jiangtype)].smallloss
	if randNum <= sum {
		result := brnn.getWinLossAssignCardNumber("smallloss")
		return int(result), false
	}

	sum += brnn.betvalue.jiangWeight[int(brnn.jiangtype)].bigloss
	if randNum <= sum {
		result := brnn.getWinLossAssignCardNumber("bigloss")
		return int(result), false
	}

	return 3, true
}

func (brnn *BRNNLogic) getWinLossAssignCardNumber(mapName string) int32 {
	lenCardNumber := len(brnn.betvalue.winLossCardNumber[mapName].cardNumber)
	randCardIndex := util.RndInt(0, lenCardNumber)
	result := brnn.betvalue.winLossCardNumber[mapName].cardNumber[randCardIndex]
	return result
}

func (brnn *BRNNLogic) checkAssignWinLoss(lossNumber int) bool {

	if len(brnn.secondPokerGroup) != int(5) {
		mlog.Warn("沒發滿poker")
		return false
	}

	copyZhuangCard := copyPokerCards(brnn.allKindPoker[0])
	copyZhuangCard.Cards = append(copyZhuangCard.Cards, brnn.secondPokerGroup[0].Cards...)
	zhuangcardret := getNiuNiuRet(copyZhuangCard)

	winLoss := int32(0)
	for i := int32(1); i <= int32(len(brnn.betvalue.alltype)); i++ { //閒家牌

		copyXianCard := copyPokerCards(brnn.allKindPoker[int(i)])
		copyXianCard.Cards = append(copyXianCard.Cards, brnn.secondPokerGroup[int(i)].Cards...)
		xiancardret := getNiuNiuRet(copyXianCard)
		betret := compareCards(zhuangcardret, xiancardret)

		if betret > int32(0) { //莊贏閒輸
			winLoss -= 1
		} else {
			winLoss += 1
		}
	}

	if winLoss != int32(lossNumber) {
		mlog.Debug("檢查莊輸的條件失敗, result winloss : %d", winLoss)
		return false
	}

	return true
}

//排序下注區塊大到小
func (brnn *BRNNLogic) getBlockOrderBy() {
	sort.Sort(brnn.bettypetotal)
	sort.Sort(brnn.bettyperobot)
}

func (sbg *betTypeBlockGroup) Len() int {
	return len(sbg.block)
}

func (sbg *betTypeBlockGroup) Swap(i, j int) {
	sbg.block[i], sbg.block[j] = sbg.block[j], sbg.block[i]
}

func (sbg *betTypeBlockGroup) Less(i, j int) bool {
	return sbg.block[i].money > sbg.block[j].money
}

//隨機選取要開獎的下注金額大到小的第幾個
func (brnn *BRNNLogic) getKajianTopBlock(times int, randNum int32) []int32 {
	rand.Seed(time.Now().Unix())
	var orderbyItem []int32 //要開出下注金額大到小的第幾個項目

	hasBet := len(brnn.bettypetotal.block) > 0

	for i := 1; i <= times; i++ {

		if randNum <= brnn.betvalue.jiangWeight[int(brnn.jiangtype)].bigwin && hasBet && i <= len(brnn.bettypetotal.block) {
			orderbyItem = append(orderbyItem, brnn.bettypetotal.block[i-1].brnntype)
		} else if randNum <= brnn.betvalue.jiangWeight[int(brnn.jiangtype)].smallwin && hasBet && i <= len(brnn.bettypetotal.block)-1 {
			orderbyItem = append(orderbyItem, brnn.bettypetotal.block[i].brnntype)
		} else if randNum <= brnn.betvalue.jiangWeight[int(brnn.jiangtype)].smallloss && hasBet && false {
			//暫時不要小輸
			if i+1 <= len(brnn.bettypetotal.block) {
				orderbyItem = append(orderbyItem, brnn.bettypetotal.block[i+1].brnntype)
			}
		} else if randNum <= brnn.betvalue.jiangWeight[int(brnn.jiangtype)].bigloss && false {
			//暫時不要大輸
			index := len(brnn.betvalue.alltype) - i
			orderbyItem = append(orderbyItem, brnn.bettypetotal.block[index].brnntype)
		} else {

			randBlock := int32(util.RndInt(1, 5))

			existItem := false
			for _, v := range orderbyItem {
				if v == randBlock {
					existItem = true
				}
			}

			if existItem {
				i--
				continue
			} else {
				orderbyItem = append(orderbyItem, randBlock)
			}
		}
	}

	return orderbyItem
}

//產生閒家被通殺的牌
func (brnn *BRNNLogic) xianTongshaCard(runtime int32) {
	//zhuangCard := getNiuNiuRet(brnn.secondPokerGroup[0])

	for i := int32(1); i <= int32(len(brnn.betvalue.alltype)); i++ {
		niuRet := nnCard_Niu
		niuCardRet := brnn.assignNewNiuNiuCardLessThen(brnn.allKindPoker[int(i)], niuRet, 5, 3)
		ok := brnn.appendSecondPokerGroup(int(i), niuCardRet)

		if !ok {
			break
		}
	}

	ok := brnn.checkTongshaCard()
	if !ok { //通殺失敗 重新組牌
		if runtime <= int32(0) {
			return
		}
		for k, _ := range brnn.secondPokerGroup {
			brnn.poker.AddCardGroup(0, brnn.secondPokerGroup[k]) //把牌放回去
		}
		if runtime > int32(0) {
			runtime--
			brnn.tongshaResultType(runtime)
		}
	}

}

func (brnn *BRNNLogic) appendSecondPokerGroup(index int, niuCardRet *niuniuCardRet) bool {
	var pokerGroup []*poker.CardGroup
	if niuCardRet == nil || len(niuCardRet.cards.Cards) < 5 {
		mlog.Warn("niuCardRet is nil")
		return false
	}
	niuCardRet.cards.Cards = niuCardRet.cards.Cards[3:5]
	pokerGroup = append(pokerGroup, niuCardRet.cards)
	brnn.secondPokerGroup[index] = niuCardRet.cards
	if len(brnn.secondPokerGroup) >= int(5) {
		return false
	}
	return true
}

//檢查通殺莊是否小於閒
func (brnn *BRNNLogic) checkTongshaCard() bool {
	if len(brnn.secondPokerGroup) != int(5) {
		mlog.Warn("沒發滿poker")
		return false
	}

	copyZhuangCard := copyPokerCards(brnn.allKindPoker[0])
	copyZhuangCard.Cards = append(copyZhuangCard.Cards, brnn.secondPokerGroup[0].Cards...)
	zhuangcardret := getNiuNiuRet(copyZhuangCard)

	for i := int32(1); i <= int32(len(brnn.betvalue.alltype)); i++ { //閒家牌
		copyXianCard := copyPokerCards(brnn.allKindPoker[int(i)])
		copyXianCard.Cards = append(copyXianCard.Cards, brnn.secondPokerGroup[int(i)].Cards...)
		xiancardret := getNiuNiuRet(copyXianCard)

		betret := compareCards(zhuangcardret, xiancardret)
		if betret < int32(0) { //莊輸
			return false
		}
	}
	return true
}

//產生莊家通殺的牌牛牛
func (brnn *BRNNLogic) zhuangTongshaCard(runtime int32) {
	brnn.secondPokerGroup = make(map[int]*poker.CardGroup, 0)
	niuCardRet := new(niuniuCardRet)

	//莊, 牛牛
	niniuCardGroup := brnn.takeNiuNiuByNumber(brnn.allKindPoker[0])

	copyZhuangCard := copyPokerCards(brnn.allKindPoker[0])
	copyZhuangCard.Cards = append(copyZhuangCard.Cards, niniuCardGroup.Cards...)
	niuCardRet = getNiuNiuRet(copyZhuangCard)

	brnn.appendSecondPokerGroup(0, niuCardRet)

	if niuCardRet.niuRet == nnCard_Niu || niuCardRet.niuRet == nnCard_WuNiu {
		mlog.Debug("莊家拿小牌")
	}

	if niuCardRet.niuRet == nnCard_WuNiu || niuCardRet.niuRet == nnCard_Niu {
		if runtime <= int32(0) {
			return
		}

		mlog.Debug("莊拿到無牛牛從新取牌")
		for k, _ := range brnn.secondPokerGroup { //牌放回去
			brnn.poker.AddCardGroup(0, brnn.secondPokerGroup[k])
		}

		if runtime > int32(0) {
			runtime--
			brnn.tongshaResultType(runtime)
		}
	}

	if niuCardRet.cards.Cards[0].CardNumber == int32(0) || niuCardRet.cards.Cards[1].CardNumber == int32(0) {
		mlog.Debug("莊家沒拿到第二手牌")
	}

	if len(niuCardRet.cards.Cards) != int(2) {
		mlog.Debug("莊家 第二手牌 超過兩張")
	}
}

//指定區塊項目輸
func (brnn *BRNNLogic) assignBlockLoss(blockItem []int32, zhungNumber int) {
	for _, v := range blockItem {
		niuCardRet := brnn.assignNewNiuNiuCardLessThen(brnn.allKindPoker[int(v)], nnCard_Niu, zhungNumber, 3)
		brnn.appendSecondPokerGroup(int(v), niuCardRet)
		mlog.Debug("指定區塊項目%v 輸 %v", v, niuCardRet)
	}
}

//指定區塊項目贏
func (brnn *BRNNLogic) assignBlockWin(isTakeMore bool) {
	copyZhuangCard := copyPokerCards(brnn.allKindPoker[0])
	copyZhuangCard.Cards = append(copyZhuangCard.Cards, brnn.secondPokerGroup[0].Cards...)
	zhuangCardRet := getNiuNiuRet(copyZhuangCard)

	for i := 1; i < len(brnn.allKindPoker); i++ {
		_, ok := brnn.secondPokerGroup[i]
		if ok {
			continue
		}
		niuRet := nnCard_Niu
		number := 2
		if zhuangCardRet.niuRet == nnCard_Niu {
			number = zhuangCardRet.number
		} else {
			niuRet = zhuangCardRet.niuRet
		}

		niniuCardGroup := new(niuniuCardRet)
		if isTakeMore {
			niniuCardGroup = brnn.assignNewNiuNiuCardMoreThen(brnn.allKindPoker[i], niuRet, number)
		} else {
			if niuRet == nnCard_WuNiu {
				niuRet = nnCard_Niu
			}
			niniuCardGroup = brnn.assignNewNiuNiuCardLessThen(brnn.allKindPoker[i], niuRet, number+4, 10)
		}

		brnn.appendSecondPokerGroup(int(i), niniuCardGroup)
	}
}

//取得每家的第二張牌(含計算賺賠機率的牌型)
func (brnn *BRNNLogic) getCalSecondPokerGroup(position int32) *poker.CardGroup {
	return brnn.secondPokerGroup[int(position)]
}

func (brnn *BRNNLogic) checkSysPool() {
	mlog.Printf("桌子号：%v，系统水池：%v，注水概率：%v", brnn.TableNo(), brnn.syspool, brnn.zhushui)
}

//改变水池的钱
func (brnn *BRNNLogic) changeSysPool(value int64) {
	brnn.startsyspool += value - brnn.syspool
	brnn.syspool = value
}

//设置状态
func (brnn *BRNNLogic) setState(state int32) {
	brnn.gamestate = state
}

//开始开奖
func (brnn *BRNNLogic) kaiJiang() {
	brnn.action_str = brnn.RegisterAnimation(GameKaiJiang_Poker, int32(1), brnn.boardcastKaijianPoker)
	//brnn.action_str = brnn.RegisterAnimation(GameKaiJiang_Action, brnn.kaijiangactiontime, brnn.sendResult)
	//brnn.action_str = brnn.RegisterAnimation(GameKaiJiang_Poker, int32(1), brnn.boardcastKaijianPoker)
}

//廣播開獎牌
func (brnn *BRNNLogic) boardcastKaijianPoker() {
	brnn.setState(KaiJiang)
	brnn.BoardcastMessage(brnn.getStateMsg(false))

	if brnn.zhuangID != 0 {
		if brnn.GetUserByUserID(brnn.zhuangID).IsRobot() {
			//brnn.removeRobotBet()
		}
	} else {
		//brnn.removeRobotBet()
	}

	brnn.getResultType()
	brnn.getKaiJiangPokerByType()
	brnn.boardcastSecondCardGroup()

	brnn.setState(KaiJiangAction)
	brnn.BoardcastMessage(brnn.getStateMsg(false))
	brnn.boardcastGameResultPokerMsg()

	brnn.addResultCardList()

	brnn.calWinLose()
	brnn.updateVipUserWin()
	brnn.action_str = brnn.RegisterAnimation(GameKaiJiang_Action, brnn.kaijiangactiontime, brnn.sendResult)
}

//给玩家发送结果
func (brnn *BRNNLogic) sendResult() {
	//brnn.calWinLose()

	brnn.writePoolToFile()
	brnn.action_str = brnn.RegisterAnimation(GameEnd_Action, brnn.endactiontime, brnn.gameEnd)
}

//計算下注輸贏結果
func (brnn *BRNNLogic) calWinLose() {
	sysTaxIncomeBefore := brnn.systaxincome
	brnn.syslastwin = 0
	brnn.setState(EndAction)
	brnn.BoardcastMessage(brnn.getStateMsg(false))

	for _, ud := range brnn.userdata {
		brnn.sendUserResult(ud)
	}

	finalBetTotal := brnn.bettotal + brnn.bettotalrobot //真人+機器人下注輸贏總和
	//這裡計算庄家輸贏及稅收

	var income, tax int64
	urZhuang := brnn.GetUserByUserID(brnn.zhuangID)
	if brnn.zhuangID == 0 || urZhuang.IsRobot() {
		//莊家是機器人或系統, 莊賺錢只抽真人下注的金額bettotal
		_, tax = brnn.GetUserWinAndTax(brnn.bettotal)
	} else {
		//莊家是真人, 莊賺錢要抽真人和機器人下注的金額
		_, tax = brnn.GetUserWinAndTax(finalBetTotal)
	}
	income, _ = brnn.GetUserWinAndTax(finalBetTotal)

	brnn.systaxincome += tax
	_, finalBetTax := brnn.GetUserWinAndTax(finalBetTotal)

	if finalBetTotal > 0 {
		finalBetTotal -= tax
	} else {
		finalBetTotal = -(finalBetTotal)
	}
	// 玩家上莊的有效投注及返還金額用對戰類的算法
	betAmount := util.ABSInt64(income)
	payAmount := betAmount + income
	if urZhuang != nil {
		urZhuang.FlowAmount += finalBetTotal
		urZhuang.BetAmount += betAmount
		urZhuang.PayAmount += payAmount
	}
	brnn.GameLog.AddGameStep("莊家贏錢 %d, 抽稅 %d, 流水 %d, 有效投注 %d, 返還金額 %d", income, finalBetTax, finalBetTotal, betAmount, payAmount)

	mlog.Debug("莊家贏錢 %d, 抽稅 %d", income, finalBetTax)

	//這裡計算水池, 系統小妹當莊
	brnn.zhuangWin += income
	changeSysPoolValue := int64(0)
	if brnn.zhuangID == 0 {
		brnn.syslastwin = income

		//系統小妹當莊的時候，玩家贏錢水池減少
		if brnn.userWinLose > 0 {
			//玩家贏錢，水線下降
			changeSysPoolValue = brnn.userWinLose
			//brnn.syspool -= brnn.userWinLose - brnn.zhushuimoney
		} else {
			//玩家輸錢，水線上升(稅)
			changeSysPoolValue = int64(math.Floor(float64(brnn.userWinLose) * brnn.GetTax()))
			//brnn.syspool -= int64(math.Floor(float64(brnn.userWinLose)*brnn.GetTax())) - brnn.zhushuimoney
		}
	} else {
		ur := brnn.GetUserByUserID(brnn.zhuangID)
		if ur != nil {
			income, tax = brnn.GetUserWinAndTax(finalBetTotal)
			//機器人當莊的時候
			if ur.IsRobot() {
				//玩家贏錢水池減少，並將注出的水倒回水池
				if brnn.userWinLose > 0 {
					brnn.systaxincome += tax //真人贏錢收稅
					//玩家贏錢，水線下降
					brnn.syslastwin = brnn.bettotal //系統上一局收入來源只記錄賠給真人的錢
					changeSysPoolValue = brnn.userWinLose
					//brnn.syspool -= brnn.userWinLose - brnn.zhushuimoney

				} else {
					//玩家輸錢，水線上升(稅)
					changeSysPoolValue = int64(math.Floor(float64(brnn.userWinLose) * brnn.GetTax()))
					//brnn.syspool -= int64(math.Floor(float64(brnn.userWinLose)*brnn.GetTax())) - brnn.zhushuimoney
					realUserIncome, _ := brnn.GetUserWinAndTax(brnn.bettotal)
					brnn.syslastwin = realUserIncome
				}
				mlog.Debug("系统水池異動%d：%d", -brnn.userWinLose, brnn.syspool)
			} else {

				//玩家當莊的時候，機器人贏錢水池就添加，並將注出的水倒回水池
				if brnn.robotWinLose > 0 {
					//機器人贏錢 水位上升(稅)
					income, tax = brnn.GetUserWinAndTax(brnn.robotWinLose)
					changeSysPoolValue = -income
					//brnn.syspool += income - brnn.zhushuimoney
					brnn.syslastwin = income
					brnn.systaxincome = tax
				} else {
					changeSysPoolValue = -brnn.robotWinLose
					//brnn.syspool += brnn.robotWinLose - brnn.zhushuimoney
					brnn.syslastwin += brnn.robotWinLose
				}
				mlog.Debug("系统水池異動%d：%d", brnn.robotWinLose, brnn.syspool)

			}

			brnn.UpdateUserScoreFunc(ur.UserID, brnn.zhuangWin, tax) // 更新金幣異動
			brnn.SendMessageToUser(ur, brnn.getUserWinMsg(brnn.zhuangWin, brnn.zhuangWin))
		} else {
			mlog.Debug("庄家数据为空")
		}
	}
	if !brnn.execFangFen {
		brnn.syspool -= changeSysPoolValue - brnn.zhushuimoney
		mlog.Debug("系统水池異動%d：%d", -changeSysPoolValue, brnn.syspool)
	} else {
		before := brnn.fangFenComplete
		if changeSysPoolValue > 0 {
			changeSysPoolValue = int64(math.Floor(float64(changeSysPoolValue) * brnn.GetTax()))
		} else { //玩家輸錢要把稅算回來
			changeSysPoolValue = int64(math.Floor(float64(changeSysPoolValue) / brnn.GetTax()))
		}
		brnn.fangFenComplete += changeSysPoolValue
		//mlog.Debug("系統放分完成, target %d, already %d-%d", brnn.fangFenTarget, before, brnn.fangFenComplete)
		brnn.GameLog.AddGameStep("系統放分完成, target %d, already %d-%d", brnn.fangFenTarget, before, brnn.fangFenComplete)
	}

	brnn.GameLog.AddGameStep("本局稅收 %d", brnn.systaxincome-sysTaxIncomeBefore)
	brnn.GameLog.AddGameStep("系統累積稅收 %d", brnn.systaxincome)

	uscoresRet := make(map[int32]int64)
	//發送沒下注的user 莊家結果
	for _, u2 := range brnn.AllTableUser() {
		sendBool := true
		for _, ud := range brnn.userdata {
			if u2.UserID == ud.user.UserID {
				sendBool = false
				break
			}
		}
		if sendBool {
			brnn.SendMessageToUser(u2, brnn.getUserWinMsg(0, brnn.zhuangWin))
		}
	}
	cscoresRet := make(map[int32]int64)
	//發送下注user結算結果
	for _, ud := range brnn.userdata {
		brnn.SendMessageToUser(ud.user, brnn.getUserWinMsg(ud.win, brnn.zhuangWin))
		uscoresRet[ud.user.UserID] = ud.updatemoney
		cscoresRet[ud.user.UserID] = 0

		brnn.GameLog.Context().UserUpdateScore(ud.user.UserID, ud.updatemoney)
	}

	brnn.syspoolLog(-changeSysPoolValue, brnn.syspool, brnn.kaijiangbiao, brnn.jiangtype, brnn.zhunagPoint, brnn.xianPoint)
	uscoresRet[brnn.zhuangID] = brnn.zhuangWin
	cscoresRet[brnn.zhuangID] = 0
	mlog.Debug("莊家贏錢：%v", uscoresRet[brnn.zhuangID])
	brnn.GameLog.GameEnd()
	brnn.bettotal = 0
}

//游戏结束
func (brnn *BRNNLogic) gameEnd() {
	brnn.setState(Wait)
	brnn.BoardcastMessage(brnn.getStateMsg(false))
	if brnn.zhuangID != 0 {
		user := brnn.GetUserByUserID(brnn.zhuangID)
		if user.IsOffLine() {
			brnn.zhuangID = 0
		}
	}
	brnn.GameEndFunc()
}

//将水池的钱写入文件
func (brnn *BRNNLogic) writePoolToFile() {
	if conf.TestMode {
		return
	}

	buff := strconv.Itoa(int(brnn.syspool))
	//bufio.ReadWriter.Peek()
	tablenostr := strconv.Itoa(int(brnn.TableNo()))
	filename := "./lastpool/" + brnn.serverid + "_" + tablenostr + ".conf"
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

//計算给每个玩家发送开奖结果
func (brnn *BRNNLogic) sendUserResult(ud *brnnuserdata) {
	var betlosemoney int64                     //下注輸給莊的錢
	var flowamount, betAmount, payAmount int64 //流水, 有效投注金額, 返還金額計算
	for i := 0; i < len(ud.userbet); i++ {
		betlosemoney += util.ABSInt64(ud.userbet[i].money)
		betAmount += ud.userbet[i].money
		for _, v := range brnn.jiangTypeWinLossOdd {
			if ud.userbet[i].brnntype == v.betType { //有下注到該項目
				winlossmoney := calWinLossMoney(ud, i, v)
				//if winlossmoney > 0 {
				//	flowamountTotal := winlossmoney - ud.userbet[i].money
				//	flowAmountAfterTax, _ := brnn.GetUserWinAndTax(flowamountTotal)
				//	flowamount += flowAmountAfterTax
				//} else {
				//	flowamount += -(winlossmoney - ud.userbet[i].money)
				//}

				//機器人不計算獎池裡面的錢
				if brnn.zhuangID == 0 {
					if !ud.user.IsRobot() {
						brnn.bettotal -= winlossmoney
					}
				} else {
					ur := brnn.GetUserByUserID(brnn.zhuangID) //當莊是否為機器人
					if ur != nil && ur.IsRobot() && ud.user.IsRobot() {
						brnn.bettotalrobot -= winlossmoney
					} else {
						if ud.user.IsRobot() {
							brnn.bettotalrobot -= winlossmoney
						} else {
							brnn.bettotal -= winlossmoney
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
		brnn.robotWinLose += money
	} else {
		brnn.userWinLose += money
	}

	var income, tax int64
	if money > int64(0) {
		income, tax = brnn.GetUserWinAndTax(money)
		ud.win = income + betlosemoney
		if !ud.user.IsRobot() { //真人贏錢才紀錄系統稅收, 機器人下注贏錢不紀錄系統稅收
			brnn.systaxincome += tax
		}
	} else {
		income = money
		ud.win = money + betlosemoney
	}
	flowamount = util.ABSInt64(income)
	payAmount = betAmount + income
	ud.updatemoney = income
	stepmsg := "用户%d下注赢钱%d, 抽稅%d, 流水%d, 有效投注金額%d, 返還金額%d"
	if ud.user.IsRobot() {
		stepmsg = "[機器人]" + stepmsg
	}

	brnn.GameLog.AddGameStep(stepmsg, ud.user.UserID, income, tax, flowamount, betAmount, payAmount)
	ud.user.FlowAmount += flowamount
	ud.user.BetAmount += betAmount
	ud.user.PayAmount += payAmount
	brnn.UpdateUserScoreFunc(ud.user.UserID, income, tax)

	mlog.Debug("用户%d，实际赢钱：%v, 顯示下注赢钱：%v, 更新下注贏錢:%v", ud.user.UserID, money, ud.win, income)

	//这里写入赢家排名列表
	winSort(ud, brnn)
}

//計算投注選項的輸贏金額
func calWinLossMoney(ud *brnnuserdata, indexType int, v *JiangTypeWinLossOdd) int64 {
	return ud.userbet[indexType].money * int64(v.odds) * int64(v.winLoss)
}

//这里写入赢家排名列表
func winSort(ud *brnnuserdata, brnn *BRNNLogic) {
	badd := false
	if ud.win > 0 {
		if len(brnn.winuserinfo) > 0 {
			for i := 0; i < len(brnn.winuserinfo); i++ {
				u := brnn.getUserDataByID(brnn.winuserinfo[i])
				if u != nil {
					if ud.win > u.win {
						if len(brnn.winuserinfo) >= 3 {
							if i == 0 {
								brnn.winuserinfo[2] = brnn.winuserinfo[1]
								brnn.winuserinfo[1] = brnn.winuserinfo[0]
								brnn.winuserinfo[0] = ud.user.UserID
							} else if i == 1 {
								brnn.winuserinfo[2] = brnn.winuserinfo[1]
								brnn.winuserinfo[1] = ud.user.UserID
							} else {
								brnn.winuserinfo[2] = ud.user.UserID
							}
						} else if len(brnn.winuserinfo) == 2 {
							brnn.winuserinfo = append(brnn.winuserinfo, ud.user.UserID)
							if i == 0 {
								brnn.winuserinfo[2] = brnn.winuserinfo[1]
								brnn.winuserinfo[1] = brnn.winuserinfo[0]
								brnn.winuserinfo[0] = ud.user.UserID
							} else if i == 1 {
								brnn.winuserinfo[2] = brnn.winuserinfo[1]
								brnn.winuserinfo[1] = ud.user.UserID
							}
						} else {
							brnn.winuserinfo = append(brnn.winuserinfo, ud.user.UserID)
							brnn.winuserinfo[1] = brnn.winuserinfo[0]
							brnn.winuserinfo[0] = ud.user.UserID
						}

						badd = true
						break
					}
				}
			}
		}

		if len(brnn.winuserinfo) < 3 && !badd {
			brnn.winuserinfo = append(brnn.winuserinfo, ud.user.UserID)
		}
	}
}

//添加开奖到历史列表
func (brnn *BRNNLogic) addResultCardList() {

	if brnn.gamestate != KaiJiangAction {
		return
	}

	cards := new(netproto.BRNN_Poker)
	zhuangcardret := getNiuNiuRet(brnn.allKindPoker[0])

	for i := 1; i < len(brnn.brnnpokerresult); i++ {
		xiancardret := getNiuNiuRet(brnn.allKindPoker[i])
		betret := compareCards(zhuangcardret, xiancardret)

		winLoss := int32(-1)

		if betret > int32(0) { //莊贏閒輸
			winLoss = 0

			mlog.Debug("添加开奖到历史列表 [%v] 閒輸 , 莊[%v] 閒[%v]", i, zhuangcardret, xiancardret)

			brnn.GameLog.AddGameStep("添加开奖到历史列表 [%v] 閒輸 , 莊[%v] 閒[%v]", i, zhuangcardret, xiancardret)
		} else {
			winLoss = 1

			mlog.Debug("添加开奖到历史列表 [%v] 閒贏 , 莊[%v] 閒[%v]", i, zhuangcardret, xiancardret)

			brnn.GameLog.AddGameStep("添加开奖到历史列表 [%v] 閒贏 , 莊[%v] 閒[%v]", i, zhuangcardret, xiancardret)
		}

		cards.RetType = append(cards.RetType, winLoss)
		brnn.recordTypeWinLossOdds(i, winLoss, zhuangcardret, xiancardret)
	}

	if brnn.assignRewardDetail != nil {
		example.invalid("===========%v", cards.RetType)
		brnn.assignRewardDetailResult(fmt.Sprintf("%d,%d,%d,%d", cards.RetType[0], cards.RetType[1], cards.RetType[2], cards.RetType[3]))
	}

	brnn.resultcardlist.PokerList = append(brnn.resultcardlist.PokerList, cards)
	if len(brnn.resultcardlist.PokerList) > 20 {
		brnn.resultcardlist.PokerList = append(brnn.resultcardlist.PokerList[:0], brnn.resultcardlist.PokerList[1:]...)
	}
}

//紀錄開獎後下注種類的和莊家比牌後的輸贏，並且記錄賠率
func (brnn *BRNNLogic) recordTypeWinLossOdds(indexType int, winLoss int32, zhuangcardret *niuniuCardRet, xiancardret *niuniuCardRet) {
	jiangWinLossOdd := new(JiangTypeWinLossOdd)
	jiangWinLossOdd.betType = int32(indexType)

	if winLoss == int32(0) { //閒loss, 莊家的牌當賠率
		jiangWinLossOdd.winLoss = int32(-1)
		jiangWinLossOdd.odds = getNiuNiuOdds(zhuangcardret) - 2 // -1是把投注成本拿回\

	} else if winLoss == int32(1) { //閒win 閒家的牌當賠率
		jiangWinLossOdd.winLoss = int32(1)
		jiangWinLossOdd.odds = getNiuNiuOdds(xiancardret)
	}

	brnn.jiangTypeWinLossOdd = append(brnn.jiangTypeWinLossOdd, jiangWinLossOdd)
}

//更新下注钱
func (brnn *BRNNLogic) updateBet(t int32, money int64) {
	brnn.bettotal += int64(money)
	mlog.Debug("BetType=%v,money=%v,brnn.bettotal=%v", t, money, brnn.bettotal)
	for _, v := range brnn.bettypetotal.block {
		if v.brnntype == t {
			v.money += int64(money)
			if v.money < 0 {
				mlog.Debug("钱小于0了！")
				v.money = 0
			}
			return
		}
	}

	b := new(bettype)
	b.brnntype = t
	b.money = int64(money)
	if b.money < 0 {
		mlog.Debug("钱小于0了！")
		b.money = 0
	}
	brnn.bettypetotal.block = append(brnn.bettypetotal.block, b)
}

//获取状态消息
func (brnn *BRNNLogic) getStateMsg(isScene bool) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNStateID)

	pdata := new(netproto.BRNN_Sate)
	pdata.State = proto.Int32(brnn.gamestate)

	pdata.IsScene = proto.Bool(isScene)
	pdata.Time = proto.Int32(brnn.getBetTime())
	pmsg.MsgData = pdata
	return pmsg
}

//用户进入
func (brnn *BRNNLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {
	brnn.ComputeGameLogicBase.UserEnter(gameuser, tableno, seatno)
	brnn.sendScene(gameuser)
}

//发送场景
func (brnn *BRNNLogic) sendScene(gameuser *gamelogic.GameUser) {
	//发送设置

	if gameuser.IsRobot() {
		mlog.Debug("给机器人%d发送场景消息", gameuser.UserID)
	} else {
		mlog.Debug("给玩家%d发送场景消息", gameuser.UserID)
	}
	brnn.sendSetting(gameuser)
	brnn.SendMessageToUser(gameuser, brnn.getResultCardListMsg())
	brnn.SendMessageToUser(gameuser, brnn.getRichestListMsg())
	brnn.SendMessageToUser(gameuser, brnn.getLuckyStarMsg())
	brnn.SendMessageToUser(gameuser, brnn.getZhuangInfoMsg())

	if brnn.gamestate == Wait {
		brnn.BoardcastMessage(brnn.getStateMsg(true))
		return
	}
	//发送阶段
	brnn.SendMessageToUser(gameuser, brnn.getBetTotalMsg(0))
	brnn.sendUserBetInfo(gameuser)
	brnn.SendMessageToUser(gameuser, brnn.getStateMsg(true))

	if brnn.gamestate == Bet {

		brnn.BroadCastAbleBetTotal()
		brnn.sendToUserFirstCardGroup(gameuser)
	}

	if brnn.gamestate == KaiJiangAction {
		brnn.BoardcastMessage(brnn.getStateMsg(true))
		//发送开奖信息
		brnn.sendToUserGameResultPokerMsg(gameuser)
	} else if brnn.gamestate == EndAction {
		brnn.BoardcastMessage(brnn.getStateMsg(true))
		//发送用户下注输赢信息如果有
		ud := brnn.getUserDataByID(gameuser.UserID)
		if ud != nil {
			brnn.SendMessageToUser(gameuser, brnn.getUserWinMsg(ud.win, brnn.zhuangWin))
		}
	}
}

//获取玩家赢钱消息
func (brnn *BRNNLogic) getUserWinMsg(selfwin int64, zhuangwin int64) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNGameResultWinID)

	pdata := new(netproto.BRNN_GameResultUserMoney)

	pdata.ZhuangWin = proto.Int64(zhuangwin)
	pdata.SelfWin = proto.Int64(selfwin)
	for i := 0; i < len(brnn.winuserinfo); i++ {
		ud := brnn.getUserDataByID(brnn.winuserinfo[i])
		if ud == nil {
			continue
		}
		userwin := new(netproto.BRNN_WinUerInfo)
		userwin.UserName = proto.String(ud.user.NickName)
		userwin.Win = proto.Int64(ud.win)
		pdata.WinUserInfo = append(pdata.WinUserInfo, userwin)

	}
	for j := 0; j < len(brnn.stageuserinfo); j++ {
		ud := brnn.getUserDataByID(brnn.stageuserinfo[j])
		if ud == nil {
			continue
		}
		stageuserwin := new(netproto.BRNN_StageUerInfo)
		stageuserwin.UserID = proto.Int32(ud.user.UserID)
		stageuserwin.Win = proto.Int64(ud.win)
		pdata.StageUerInfo = append(pdata.StageUerInfo, stageuserwin)
	}

	pmsg.MsgData = pdata
	return pmsg
}

//發送百人牛牛所有手牌的資料
func (brnn *BRNNLogic) boardcastGameResultPoker() {
	msg := brnn.getGameResultPokerMsg()
	brnn.appendPokerResult(msg)
	brnn.BoardcastMessage(msg)
}

func (brnn *BRNNLogic) appendPokerResult(pdata *msg.Message) {
	msgData := pdata.MsgData.(*netproto.BRNN_GameResultPoker) //須將interface轉[]byte
	for i := int32(0); i <= int32(len(brnn.betvalue.alltype)); i++ {
		brnn.brnnpokerresult = append(brnn.brnnpokerresult, msgData.NiuNiuRet) //紀錄開獎的牌結果點數
	}
}

//获取开奖结果扑克牌消息
func (brnn *BRNNLogic) getGameResultPokerMsg() *msg.Message {
	zhuangcardret := getNiuNiuRet(brnn.allKindPoker[0])
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNGameResultPokerID)
	//獲得玩家牌點(niuniuNumber)
	for i := int32(0); i <= int32(len(brnn.betvalue.alltype)); i++ {
		pdata := new(netproto.BRNN_GameResultPoker)
		mlog.Debug("發送百人牛牛所有手牌的手牌 allKindPoker[%v] = %v", i, brnn.allKindPoker[int(i)])
		zjcardret := getNiuNiuRet(brnn.allKindPoker[int(i)])
		mlog.Debug("發送百人牛牛所有手牌的結果 zjcardret[%v] = %v", i, zjcardret)

		brnn.GameLog.AddGameStep("index: %v 所有手牌的結果:%v ", i, zjcardret)

		pdata.SitNo = proto.Int32(i)
		pdata.NiuNiuNumber = proto.Int32(int32(zjcardret.number))
		pdata.NiuNiuRet = proto.Int32(int32(zjcardret.niuRet))
		pdata.SitPoker = zjcardret.getNetCardGroup()

		//輸贏結果
		xiancardret := getNiuNiuRet(brnn.allKindPoker[int(i)])
		betret := compareCards(zhuangcardret, xiancardret)

		winLoss := int32(-1)
		if betret > int32(0) { //莊贏閒輸
			winLoss = 0
		} else {
			winLoss = 1
		}
		pdata.WinLoss = proto.Int32(winLoss)

		pmsg.MsgData = pdata
	}

	return pmsg
}

//获取开奖结果扑克牌消息, 發送百人牛牛所有手牌的資料
func (brnn *BRNNLogic) boardcastGameResultPokerMsg() {
	zhuangcardret := getNiuNiuRet(brnn.allKindPoker[0])
	for i := int32(0); i <= int32(len(brnn.betvalue.alltype)); i++ {
		pmsg := new(msg.Message)
		pmsg.BClassID = int32(netproto.MessageBClassID_Game)
		pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNGameResultPokerID)

		pdata := new(netproto.BRNN_GameResultPoker)
		zjcardret := getNiuNiuRet(brnn.allKindPoker[int(i)])

		brnn.GameLog.AddGameStep("index: %v 所有手牌的結果:%v ", i, zjcardret)

		pdata.SitNo = proto.Int32(i)
		pdata.NiuNiuNumber = proto.Int32(int32(zjcardret.number))
		pdata.NiuNiuRet = proto.Int32(int32(zjcardret.niuRet))
		pdata.SitPoker = zjcardret.getNetCardGroup()

		//輸贏結果
		xiancardret := getNiuNiuRet(brnn.allKindPoker[int(i)])
		if i != 0 {
			//非莊家，記錄牌點
			brnn.xianPoint = append(brnn.xianPoint, xiancardret.number)
		} else {
			//記錄莊家牌點
			brnn.zhunagPoint = zhuangcardret.number
		}
		betret := compareCards(zhuangcardret, xiancardret)

		winLoss := int32(-1)
		if betret > int32(0) { //莊贏閒輸
			winLoss = 0
		} else {
			winLoss = 1
		}
		pdata.WinLoss = proto.Int32(winLoss)

		pmsg.MsgData = pdata
		brnn.BoardcastMessage(pmsg)

		brnn.brnnpokerresult = append(brnn.brnnpokerresult, pdata.NiuNiuRet) //紀錄開獎的牌結果點數
	}
}

func (brnn *BRNNLogic) sendToUserGameResultPokerMsg(gameuser *gamelogic.GameUser) {
	zhuangcardret := getNiuNiuRet(brnn.allKindPoker[0])
	for i := int32(0); i <= int32(len(brnn.betvalue.alltype)); i++ {
		pmsg := new(msg.Message)
		pmsg.BClassID = int32(netproto.MessageBClassID_Game)
		pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNGameResultPokerID)

		pdata := new(netproto.BRNN_GameResultPoker)
		zjcardret := getNiuNiuRet(brnn.allKindPoker[int(i)])

		brnn.GameLog.AddGameStep("index: %v 所有手牌的結果:%v ", i, zjcardret)

		pdata.SitNo = proto.Int32(i)
		pdata.NiuNiuNumber = proto.Int32(int32(zjcardret.number))
		pdata.NiuNiuRet = proto.Int32(int32(zjcardret.niuRet))
		pdata.SitPoker = zjcardret.getNetCardGroup()

		//輸贏結果
		xiancardret := getNiuNiuRet(brnn.allKindPoker[int(i)])
		betret := compareCards(zhuangcardret, xiancardret)

		winLoss := int32(-1)
		if betret > int32(0) { //莊贏閒輸
			winLoss = 0
		} else {
			winLoss = 1
		}
		pdata.WinLoss = proto.Int32(winLoss)

		pmsg.MsgData = pdata
		brnn.SendMessageToUser(gameuser, pmsg)

	}
}

//发送用户的下注信息
func (brnn *BRNNLogic) sendUserBetInfo(gameuser *gamelogic.GameUser) {
	ud := brnn.getUserDataByID(gameuser.UserID)
	if ud == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNBetInfoID)

	pdata := new(netproto.BRNN_UserBetInfo)
	pdata.UserMoney = proto.Int64(ud.money)
	lucky := brnn.getUserDataByID(brnn.luckystar) //獲取明燈
	for _, v := range ud.userbet {
		betinfo := new(netproto.BRNN_BroadCastBet)
		betinfo.Type = proto.Int32(v.brnntype)
		betinfo.Money = proto.Int64(v.money)
		betinfo.LuckyStar = proto.Bool(false)
		if lucky != nil && lucky.getUserBetByType(v.brnntype) > 0 {
			betinfo.LuckyStar = proto.Bool(true)
		}
		pdata.BetInfo = append(pdata.BetInfo, betinfo)
	}
	pmsg.MsgData = pdata
	brnn.SendMessageToUser(gameuser, pmsg)
}

//获取用户数据通过用户ID
func (brnn *BRNNLogic) getUserDataByID(userid int32) *brnnuserdata {
	for _, ud := range brnn.userdata {
		if ud.user.UserID == userid {
			return ud
		}
	}

	return nil
}

//获取其他玩家下注
func (brnn *BRNNLogic) getBetTotalMsg(userid int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNBroadCastBetID)

	pdata := new(netproto.BRNN_BroadCastBetArr)
	pdata.UserID = proto.Int32(userid)
	lucky := brnn.getUserDataByID(brnn.luckystar) //獲取明燈
	for _, v := range brnn.betvalue.alltype {
		betinfo := new(netproto.BRNN_BroadCastBet)
		betinfo.Type = proto.Int32(v)
		betinfo.Money = proto.Int64(brnn.getBetMoneyByType(v))
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
func (brnn *BRNNLogic) getBetMoneyByType(t int32) int64 {
	allBetTotal := int64(0)
	for i := 0; i < len(brnn.bettypetotal.block); i++ {
		if brnn.bettypetotal.block[i].brnntype == t {
			allBetTotal += brnn.bettypetotal.block[i].money
		}
	}

	for i := 0; i < len(brnn.bettyperobot.block); i++ {
		if brnn.bettyperobot.block[i].brnntype == t {
			allBetTotal += brnn.bettyperobot.block[i].money
		}
	}

	return allBetTotal
}

//获取开奖历史列表的消息
func (brnn *BRNNLogic) getResultCardListMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNResultCardListID)

	pdata := new(netproto.BRNN_ResultCardList)
	pdata.PokerList = append(pdata.PokerList, brnn.resultcardlist.PokerList...)
	mlog.Debug("getResultCardListMsg:%v", pdata.PokerList)
	pmsg.MsgData = pdata
	return pmsg
}

//发送设置
func (brnn *BRNNLogic) sendSetting(gameuser *gamelogic.GameUser) {
	mlog.Debug("sendSetting")
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNSettingID)

	pdata := new(netproto.BRNN_Setting)
	pdata.BetCount = append(pdata.BetCount, brnn.betcount...)
	pdata.BetMoney = proto.Int32(int32(brnn.betvalue.getMoneyByType(1)))
	pdata.ZhuangLimit = proto.Int64(brnn.zhuanglimit)
	pdata.MaxZhuang = proto.Int32(brnn.maxturn)
	pdata.OnceFapei = proto.Int32(brnn.secondfapeitimeinterval)
	pdata.Betlimitmoney = proto.Int64(brnn.betlimitmoney)
	pmsg.MsgData = pdata
	brnn.SendMessageToUser(gameuser, pmsg)
}

//注销用户实例
func (brnn *BRNNLogic) UserFree(gameuser *gamelogic.GameUser) {
	if brnn.gamestate != Wait {
		if brnn.getUserDataByID(gameuser.UserID) != nil {
			mlog.Debug("用户在玩游戏还不能离开")
			return
		}

		if brnn.zhuangID == gameuser.UserID {
			mlog.Debug("用户在当庄，所以用户不能离开")
			return
		}
	} else {
		if brnn.zhuangID == gameuser.UserID {
			brnn.zhuangID = 0
		}
	}

	brnn.delFromZhuangList(gameuser.UserID)

	brnn.ChangeUserStateEnd(gameuser)
	brnn.ComputeGameLogicBase.UserFree(gameuser)
}

//断线重连
func (brnn *BRNNLogic) UserRecover(gameuser *gamelogic.GameUser) {
	mlog.Debug("UserRecover")
	brnn.ComputeGameLogicBase.UserRecover(gameuser)
	ud := brnn.getUserDataByID(gameuser.UserID)
	if ud != nil {
		ud.money += ud.user.Score - ud.moneylast
		ud.moneylast = ud.user.Score
	}

	brnn.sendScene(gameuser)
}

//用户离开
func (brnn *BRNNLogic) UserLeave(gameuser *gamelogic.GameUser) bool {
	if brnn.gamestate != Wait {
		if brnn.getUserDataByID(gameuser.UserID) != nil {
			mlog.Debug("用户在玩游戏还不能离开")
			return false
		}

		if brnn.zhuangID == gameuser.UserID {
			mlog.Debug("用户在当庄，所以用户不能离开")
			return false
		}
	}

	brnn.delFromZhuangList(gameuser.UserID)

	brnn.ChangeUserStateEnd(gameuser)
	return brnn.ComputeGameLogicBase.UserLeave(gameuser)
}

//用户消息主处理函数
func (brnn *BRNNLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	defer util.PanicErrStack()
	if !brnn.CheckGameUserMessage(user) {
		brnn.KickUserFunc(user.UserID, "非法操作")
		return
	}
	//不是给游戏逻辑服务器的消息
	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	brnn.CheckUserAnimationCommited(user, gmsg, brnn.action_str, nil)

	if gmsg.SClassID == int32(netproto.BRNN_GameMessageClassID_BRNNGetShangZhuangListID) {
		brnn.SendMessageToUser(user, brnn.getZhuangListMsg())
		return
	}

	if gmsg.SClassID == int32(netproto.BRNN_GameMessageClassID_BRNNBetID) {
		brnn.userBet(user, gmsg.MsgData.([]byte))
		return
	}

	if gmsg.SClassID == int32(netproto.BRNN_GameMessageClassID_BRNNShangZhuangID) {
		brnn.shangZhuang(user, gmsg.MsgData.([]byte))
		return
	}

	if gmsg.SClassID == int32(netproto.BRNN_GameMessageClassID_BRNNBetClearID) && brnn.gamestate == Bet {
		brnn.clearUserBet(user)
		return
	}

	if gmsg.SClassID == int32(netproto.BRNN_GameMessageClassID_BRNNGetFapeiThreeCard) && brnn.gamestate == Bet {
		brnn.sendToUserFirstCardGroup(user)
	}
}

//获取上庄列表消息
func (brnn *BRNNLogic) getZhuangListMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNShangZhuangListInfoID)

	pdata := new(netproto.BRNN_ZhuangList)
	for i := 0; i < len(brnn.zhuanglist); i++ {
		ur := brnn.GetUserByUserID(brnn.zhuanglist[i])
		pdata.UserInfo = append(pdata.UserInfo, ur.GetProtoBufGameUser())
	}
	pmsg.MsgData = pdata
	return pmsg
}

//清除玩家下注信息
func (brnn *BRNNLogic) clearUserBet(user *gamelogic.GameUser) {
	ud := brnn.getUserDataByID(user.UserID)
	if ud == nil {
		return
	}

	mlog.Debug("ud.clearcount=%v", ud.clearcount)
	if ud.clearcount >= brnn.maxclearcount {
		brnn.SendMessageToUser(user, brnn.getClearRetMsg(ud.money, -1))
		return
	}

	ud.clearcount++
	for _, ub := range ud.getUserBet() {
		brnn.updateBet(ub.brnntype, -ub.money)
	}

	brnn.GameLog.AddGameStep("用户%d点击清除", user.UserID)
	ud.clearData()
	//brnn.bettotal -= ud.user.Score - ud.money
	//brnn.delUserData(user)
	brnn.SendMessageToUser(user, brnn.getClearRetMsg(user.Score, 0))
	brnn.BoardcastMessage(brnn.getBetTotalMsg(user.UserID))
	brnn.BroadCastAbleBetTotal()
}

//获取清除结果消息
func (brnn *BRNNLogic) getClearRetMsg(money int64, ret int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNBetClearRetID)

	pdata := new(netproto.BRNN_BetClearRet)
	pdata.Money = proto.Int64(money)
	pdata.Ret = proto.Int32(ret)
	pmsg.MsgData = pdata

	return pmsg
}

//处理用户上庄请求
func (brnn *BRNNLogic) shangZhuang(user *gamelogic.GameUser, data []byte) {
	ust := &netproto.BRNN_ShangZhuang{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Debug("解析用户%v上庄错误%v", user, err)
		return
	}

	if brnn.getZhuangListIndexByUserId(user.UserID) != -1 {
		brnn.GameLog.AddGameStep("用户%d下庄", user.UserID)
		brnn.delFromZhuangList(user.UserID)
		brnn.SendMessageToUser(user, brnn.getShangzhuangRetMsg(1))
		return
	}

	if user.Score < brnn.zhuanglimit {
		brnn.SendMessageToUser(user, brnn.getShangzhuangRetMsg(-1))
		return
	}

	if *ust.IsShangZhuang {
		if brnn.publiczhuang == 1 {
			brnn.GameLog.AddGameStep("用户%d上庄", user.UserID)
			brnn.addToZhuangList(user.UserID)
			brnn.SendMessageToUser(user, brnn.getShangzhuangRetMsg(0))
		} else {
			mlog.Error("未開啟上庄機制，不得上庄")
			return
		}
	} else {
		brnn.GameLog.AddGameStep("用户%d下庄", user.UserID)
		brnn.delFromZhuangList(user.UserID)
		brnn.SendMessageToUser(user, brnn.getShangzhuangRetMsg(1))
		brnn.lastturn = -1
	}
}

//获取上庄结果消息
func (brnn *BRNNLogic) getShangzhuangRetMsg(ret int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNShangZhuangRetID)

	pdata := new(netproto.BRNN_ShangzhuangRet)
	pdata.Ret = proto.Int32(ret)
	pmsg.MsgData = pdata
	return pmsg
}

//获取上庄列表索引通过用户ID
func (brnn *BRNNLogic) getZhuangListIndexByUserId(userid int32) int32 {
	for i := 0; i < len(brnn.zhuanglist); i++ {
		if brnn.zhuanglist[i] == userid {
			return int32(i)
		}
	}

	return -1
}

//用户下注处理函数
func (brnn *BRNNLogic) userBet(user *gamelogic.GameUser, data []byte) {
	defer util.PanicErrStack()
	if brnn.gamestate != Bet {
		mlog.Debug("解析用户%v 非下注階段", user)
		return
	}

	ust := &netproto.BRNN_BetArr{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Debug("解析用户%v下注错误%v", user, err)
		return
	}

	ud := brnn.getUserDataByID(user.UserID)
	if ud == nil {
		ud = brnn.addUserData(user)
	}

	ret := netproto.RetType(1)

	ur := brnn.GetUserByUserID(user.UserID)

	//玩家身上金額小於下注限制
	if user.Score < brnn.betlimitmoney*100 {
		for i, u := range brnn.userdata {
			if u.user.UserID == ud.user.UserID {
				brnn.userdata = append(brnn.userdata[:i], brnn.userdata[i+1:]...)
			}
		}
		brnn.SendMessageToUser(user, brnn.getBetRetMsg(ud, 6, 0, ust.BetArr))
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

	score := int64(0)
	alreadyTotalBet := int64(brnn.bettotal + brnn.bettotalrobot)
	ableBet := ur.Score / int64(10)
	thisBet := int64(0)
	for _, b := range ust.BetArr {
		ableBet = ableBet - int64(*b.Count)
		thisBet += int64(*b.Count)
	}

	if (ableBet - userBetTotal) <= int64(0) { //超過可下注的十分之一
		ret = netproto.RetType(2)
		brnn.SendMessageToUser(user, brnn.getBetRetMsg(ud, ret, 0, ust.BetArr))
		return
	}
	if brnn.zhuangID == int32(0) { //系統當莊, 可下注總額
		score = int64(brnn.betvalue.betvalue[0].maxmoney +
			brnn.betvalue.betvalue[1].maxmoney +
			brnn.betvalue.betvalue[2].maxmoney +
			brnn.betvalue.betvalue[3].maxmoney)
	} else {
		ur := brnn.GetUserByUserID(brnn.zhuangID)
		score = ur.Score
	}

	//是否超過下注限額
	ableTotal := (score / 10) - (alreadyTotalBet + thisBet)
	if ableTotal < int64(0) {
		ret = netproto.RetType(5)
		brnn.SendMessageToUser(user, brnn.getBetRetMsg(ud, ret, 0, ust.BetArr))
		return
	}

	//庄家不能下注
	if brnn.zhuangID == user.UserID {
		return
	}

	money := int32(0)
	for _, b := range ust.BetArr {
		*b.BetType = int32(math.Abs(float64(*b.BetType)))
		*b.Count = int32(math.Abs(float64(*b.Count)))

		tmpmon := brnn.getMoneyByType(*b.BetType) * (*b.Count)
		if (int64(tmpmon) + brnn.getBetMoneyByType(*b.BetType)) > brnn.betvalue.getMaxMoneyByType(*b.BetType) {
			brnn.SendMessageToUser(user, brnn.getBetRetMsg(ud, netproto.RetType(5), *b.BetType, ust.BetArr))
			return
		}

		if brnn.zhuangID != 0 {
			ur := brnn.GetUserByUserID(brnn.zhuangID)
			if !ur.IsRobot() {
				allmoney := (int64(tmpmon) + brnn.getBetMoneyByType(*b.BetType)) * int64(brnn.betvalue.getOddByType(*b.BetType))
				if allmoney > (ur.Score / 3) {
					brnn.SendMessageToUser(user, brnn.getBetRetMsg(ud, netproto.RetType(5), *b.BetType, ust.BetArr))
					return
				}
			}
		}

		money += tmpmon
	}

	if int64(money) > ud.money {
		ret = netproto.RetType(2)
	} else {
		for _, b := range ust.BetArr {
			betmon := brnn.getMoneyByType(*b.BetType) * (*b.Count)
			mlog.Debug("收到下注消息%d", betmon)
			ud.setUserBet(*b.BetType, betmon)

			if user.IsRobot() {
				brnn.updaterobotBet(*b.BetType, int64(betmon))
			} else {
				brnn.updateBet(*b.BetType, int64(betmon))
			}
			brnn.BoardcastMessage(brnn.getBetTotalMsg(user.UserID))
			brnn.BroadCastAbleBetTotal()
		}

		brnn.GameLog.AddGameStep("用户%d下注%v", user.UserID, ust.BetArr)
	}

	brnn.SendMessageToUser(user, brnn.getBetRetMsg(ud, ret, 0, ust.BetArr))
}

//更新机器人下注额
func (brnn *BRNNLogic) updaterobotBet(t int32, money int64) {
	brnn.bettotalrobot += int64(money)
	for _, v := range brnn.bettyperobot.block {
		if v.brnntype == t {
			v.money += money
			if v.money < 0 {
				v.money = 0
			}
			mlog.Debug("下注类型%d，下注的总钱数为%d", t, v.money)
			return
		}
	}

	b := new(bettype)
	b.brnntype = t
	b.money = int64(money)
	if b.money < 0 {
		b.money = 0
	}

	brnn.bettyperobot.block = append(brnn.bettyperobot.block, b)
}

//添加一个用户数据
func (brnn *BRNNLogic) addUserData(user *gamelogic.GameUser) *brnnuserdata {
	ud := newLongHuDouUserData(user)
	brnn.userdata = append(brnn.userdata, ud)
	return ud
}

//获取下注消息
func (brnn *BRNNLogic) getBetRetMsg(ud *brnnuserdata, ret netproto.RetType, t int32, betinfo []*netproto.BRNN_Bet) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNBetRetID)

	pdata := new(netproto.BRNN_BetRet)
	pdata.Money = proto.Int64(ud.money)
	pdata.Ret = &ret
	pdata.Type = proto.Int32(t)
	pdata.BetArr = betinfo
	pmsg.MsgData = pdata
	return pmsg
}

//获取此类型一次下注额
func (brnn *BRNNLogic) getMoneyByType(t int32) int32 {
	return brnn.betvalue.getMoneyByType(t)
}

func (brnn *BRNNLogic) OnRecvHTTPMessage(req *http.Request) proto.Message {
	action := req.PostFormValue("action")
	if action == "getserverinfo" {
		return brnn.getTableListMsg()
	} else if action == "gettableinfo" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		return brnn.getTableDetailMsg(int32(tid))
	} else if action == "modresult" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		ret := req.PostFormValue("result")
		adminuser := req.PostFormValue("adminuser")
		jiang := int32(assign)
		if ret == "5" {
			jiang = int32(tongsha)
		} else if ret == "6" {
			jiang = int32(tongpei)
		} else if ret == "7" {
			jiang = int32(win3loss1)
		} else if ret == "8" {
			jiang = int32(win1loss3)
		} else if ret == "9" {
			jiang = int32(win2loss2)
		}
		return brnn.setNextJiangType(int32(tid), jiang, ret, adminuser)
	} else if action == "modpoolamount" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 64)
		typeid, _ := strconv.ParseInt(req.PostFormValue("type"), 10, 32)
		pool, _ := strconv.ParseInt(req.PostFormValue("poolamount"), 10, 64)
		operator := req.PostFormValue("operator")
		if typeid == 1 {
			return brnn.setSysPool(int32(tid), pool, operator)
		} else {
			return brnn.setFangFen(int32(tid), pool, operator)
		}

	} else if action == "modpopwaterrate" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		zhushui, _ := strconv.ParseInt(req.PostFormValue("popwaterrate"), 10, 32)
		return brnn.setZhuShui(int32(tid), int32(zhushui))
	} else if action == "history" {
		pdata := new(netproto.BRNN_ResultCardList)

		if len(brnn.resultcardlist.PokerList) > 9 {
			pdata.PokerList = brnn.resultcardlist.PokerList[len(brnn.resultcardlist.PokerList)-9 : len(brnn.resultcardlist.PokerList)]
		} else {
			pdata = brnn.resultcardlist
		}

		return pdata
	}

	return nil
}

func (brnn *BRNNLogic) setZhuShui(tableid int32, zhushui int32) proto.Message {
	if tableid == -1 {
		for _, table := range brnn.GetTable().Room.Tables {
			logic := table.GameLogic.(*BRNNLogic)
			logic.changeZhuShui(zhushui)
		}
	} else {
		t := brnn.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*BRNNLogic)
			logic.changeZhuShui(zhushui)
		} else {
			return brnn.getRetMsg(false, "修改注水率失败")
		}
	}

	return brnn.getRetMsg(true, "修改注水率成功")
}

//获取桌子详细数据
func (brnn *BRNNLogic) getTableDetailMsg(tableid int32) proto.Message {
	t := brnn.GetTable().Room.GetTableByNo(tableid)
	pTable := new(netproto.Http_BRNN_TableDetail)
	if t != nil {
		logic := t.GameLogic.(*BRNNLogic)
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
		currmoney := logic.syspool
		if zhuang != nil {
			if zhuang.IsRobot() {
				zhuangtype = 1
				currmoney = zhuang.Score
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
		pTable.SumEarnAmount = proto.Int64(logic.syspool)
		pTable.LastEarnAmount = proto.Int32(int32(logic.syslastwin))
		pTable.SysTaxIncome = proto.Int64(logic.systaxincome)
		pTable.OpenResultChingLong = proto.Int64(logic.getRealUserWinByType(block1))
		pTable.OpenResultBaiHu = proto.Int64(logic.getRealUserWinByType(block2))
		pTable.OpenResultChuChueh = proto.Int64(logic.getRealUserWinByType(block3))
		pTable.OpenResultHsuanWu = proto.Int64(logic.getRealUserWinByType(block4))
		pTable.BetTime = proto.Int32(brnn.getBetTime())
		if logic.gamestate == KaiJiangAction {
			pTable.CurentOpenResult = proto.String(logic.getResultString())
			pTable.CurentOpenResultEarn = proto.Int64(logic.getRealUserWinByType(int32(logic.jiangtype)))
		}

		pTable.UserInfo = logic.responseCurrentUserBetList()
		pTable.LastGameUserBet = logic.responseLastGameUserBetList()

		btwrate := logic.getBetTypeWinRate(9)
		pTable.ChingLongWinRate = proto.Int32(btwrate[0])
		pTable.BaiHuWinRate = proto.Int32(btwrate[1])
		pTable.ChuChuehWinRate = proto.Int32(btwrate[2])
		pTable.HsuanWuWinRate = proto.Int32(btwrate[3])
		pTable.WinLine = proto.Int64(logic.winline)
		pTable.FanFengTarget = proto.Int64(logic.fangFenTarget)
		pTable.FanFengComplete = proto.Int64(logic.fangFenComplete)
	}

	mlog.Debug("获取桌子%v详细数据   Data: %v", tableid, pTable)

	return pTable
}

func (brnn *BRNNLogic) getResultString() string {
	var res string
	res += resultstring[brnn.jiangtype]
	return res
}

//获取系统赢钱数
func (brnn *BRNNLogic) getRealUserWinByType(t int32) int64 {
	win := int64(0)
	odd := int64(brnn.betvalue.getOddByType(t))
	user := brnn.GetUserByUserID(brnn.zhuangID)
	if user != nil && !user.IsRobot() {
		_, tax := brnn.GetUserWinAndTax(brnn.bettotal)
		if brnn.currtax != 0 {
			tax = brnn.currtax
		}

		robotwin := odd * brnn.getRobotBetMoneyByType(t)
		win = brnn.bettotalrobot - robotwin - tax
		return -win
	}

	for _, ud := range brnn.userdata {
		if !ud.user.IsRobot() {
			win += odd * ud.getUserBetByType(t)
			win += ud.money - ud.user.Score
		}
	}

	mlog.Debug("玩家赢钱%d", win)
	return -win
}

//获取此类型的机器人下注总额
func (brnn *BRNNLogic) getRobotBetMoneyByType(t int32) int64 {
	for i := 0; i < len(brnn.bettyperobot.block); i++ {
		if brnn.bettyperobot.block[i].brnntype == t {
			return brnn.bettyperobot.block[i].money
		}
	}

	return 0
}

//获取桌子列表数据
func (brnn *BRNNLogic) getTableListMsg() proto.Message {
	pdata := new(netproto.Http_LHD_GameInfo)
	for _, table := range brnn.GetTable().Room.Tables {
		pTable := new(netproto.Http_LHD_TableList)
		logic := table.GameLogic.(*BRNNLogic)
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
		pTable.SysTaxIncome = proto.Int64(logic.systaxincome)
		pdata.TableList = append(pdata.TableList, pTable)
	}

	return pdata
}

//取得下注時間
func (brnn *BRNNLogic) getBetTime() int32 {
	time := int32(0)
	if brnn.opttimer != nil {
		time = brnn.opttimer.WaitSeconds - brnn.opttimer.ElapsedSeconds()
	}
	return time
}

//设置下次开奖结果
func (brnn *BRNNLogic) setNextJiangType(tableid int32, jiang int32, nextitem string, adminuser string) proto.Message {
	t := brnn.GetTable().Room.GetTableByNo(tableid)

	if t != nil {
		logic := t.GameLogic.(*BRNNLogic)
		if adminuser != "" {
			logic.GameLog.SetNextJiangRecord(adminuser, nextitem)
		}
		logic.nextjiangtype = jiangTypeEnum(jiang)
		logic.nextjiangitem = nextitem
		mlog.Debug("修改開獎結果:%v", nextitem)
		return brnn.getRetMsg(true, "修改开奖结果成功")
	}

	return brnn.getRetMsg(false, "修改开奖结果失败")
}

func (brnn *BRNNLogic) getRetMsg(ret bool, msg string) proto.Message {
	pdata := new(netproto.Http_Result)
	pdata.Result = proto.Bool(ret)
	pdata.Message = proto.String(msg)
	return pdata
}

//換牌
func (brnn *BRNNLogic) ExchangeMaxCard(userid int32) {

}

//排序身上攜帶的錢大到小
func (brnn *BRNNLogic) getMoneyRanking() {
	sort.Sort(brnn.userscore)
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
func (brnn *BRNNLogic) getLuckyStar() {
	t := 0
	gameuser := brnn.AllGameUser()
	if brnn.luckystar == 0 || brnn.notbettimes >= 3 {
		//如果沒有幸運星、連續三場沒下注或是幸運星的持續局數到了，隨機一個玩家
		t = 1
	} else {
		//檢查玩家還在不在場
		for _, u := range gameuser {
			if u.UserID == brnn.luckystar {
				brnn.luckystartimes += 1
				break
			} else {
				t = 1
			}
		}
		//檢查幸運星是不是莊
		if brnn.zhuangID == brnn.luckystar {
			t = 1
		}
	}
	if t == 1 || brnn.luckystartimes == brnn.luckystarround {
		if len(gameuser) != 0 {
			userIndex := rand.Intn(len(gameuser))
			brnn.luckystar = gameuser[userIndex].UserID
			brnn.luckyfaceid = gameuser[userIndex].FaceID
			if brnn.luckystar == brnn.zhuangID {
				if (userIndex - 1) >= 0 { //往前找
					brnn.luckystar = gameuser[userIndex-1].UserID
					brnn.luckyfaceid = gameuser[userIndex-1].FaceID
				} else if (userIndex + 1) < len(gameuser) { //往後找
					brnn.luckystar = gameuser[userIndex+1].UserID
					brnn.luckyfaceid = gameuser[userIndex+1].FaceID
				} else {
					brnn.luckystar = 0
					brnn.luckyfaceid = 0
				}
			}
		} else {
			//如果沒有明燈，須和client確認要傳什麼
			brnn.luckystar = 0
			brnn.luckyfaceid = 0
		}
		brnn.luckystartimes = 0
		brnn.notbettimes = 0
	}
}

//獲取富豪玩家
func (brnn *BRNNLogic) getRichestList() {
	userscore := new(userScore)
	userscore.userdata = brnn.AllTableUser()
	sort.Sort(userscore)
	//計算分數最高的7個
	for index, user := range userscore.userdata {
		if index < 6 {
			rich := new(netproto.BRNN_Richest)
			rich.Ranking = proto.Int32(int32(index + 1))
			rich.UserID = proto.Int32(user.UserID)
			rich.UserFaceID = proto.Int32(user.FaceID)
			rich.Money = proto.Int64(user.Score)
			brnn.richestList = append(brnn.richestList, rich)
			brnn.richestuserinfo = append(brnn.richestuserinfo, user.UserID)
		} else {
			break
		}
	}
	brnn.stageuserinfo = append(brnn.richestuserinfo, brnn.luckystar)
}

//获取大富豪列表的消息
func (brnn *BRNNLogic) getRichestListMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNRichestListID)
	pdata := new(netproto.BRNN_RichestList)
	pdata.RichestList = brnn.richestList
	//如果當下沒有玩家 回傳空值
	pmsg.MsgData = pdata
	return pmsg
}

//获取幸運星的消息
func (brnn *BRNNLogic) getLuckyStarMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BRNN_GameMessageClassID_BRNNLuckyStarListID)
	luckyStarList := new(netproto.BRNN_LuckyStar)
	luckyStarList.UserID = proto.Int32(brnn.luckystar)
	luckyStarList.UserFaceID = proto.Int32(brnn.luckyfaceid)
	pmsg.MsgData = luckyStarList
	return pmsg
}

//收到RPC消息解析處理
func (brnn *BRNNLogic) OnRecvRPCMessage(m *msg.Message) {
	rm, ok := m.MsgData.(*rpcmsg.RPCMessage)
	if !ok {
		return
	}
	mmbyte := rm.MsgData.([]byte)
	if len(mmbyte) == 0 {
		return
	}

	if m.SClassID == int32(netproto.GameMessageClassID_SuperToolGetAssignRewardDetailRetID) {
		brnn.setAssignReward(m)
	}

}

//pool_all系統水池Log
func (brnn *BRNNLogic) syspoolLog(changeSysValue int64, syspoolNow int64, jiangType jiangTypeEnum, jiangWeight jiangTypeEnum, zhuangPoint int, xianPoint []int) {

	if conf.TestMode {
		return
	}

	type LogData struct {
		GroupName      string        //系統水池
		StartSyspool   int64         //原始水池值
		ChangeSysValue int64         //異動水池值
		SyspoolNow     int64         //更新後水池值
		JiangType      jiangTypeEnum //開獎表名稱
		JiangWeight    jiangTypeEnum //贏分類型(通殺賠)
		ZhuangPoint    int           //莊家牌點
		Xian1Point     int           //閒家1牌點
		Xian2Point     int           //閒家2牌點
		Xian3Point     int           //閒家3牌點
		Xian4Point     int           //閒家4牌點
	}

	msg := LogData{
		GroupName:      "syspool_all",
		StartSyspool:   syspoolNow + changeSysValue,
		ChangeSysValue: changeSysValue,
		SyspoolNow:     syspoolNow,
		JiangType:      jiangType,   //brnn.kaijiangbiao,
		JiangWeight:    jiangWeight, //brnn.jiangtype,
		ZhuangPoint:    zhuangPoint,
		Xian1Point:     xianPoint[0],
		Xian2Point:     xianPoint[1],
		Xian3Point:     xianPoint[2],
		Xian4Point:     xianPoint[3],
	}

	mylog.MongoLog(0, 0, brnn.TableNo(), msg)
	//mlog.Debug("系統水池log [%+v]", msg)
}
