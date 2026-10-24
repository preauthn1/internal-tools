package hhdz

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

const (
	Bet            = iota // 下注状态
	KaiJiang              // 開獎状态
	KaiJiangAction        // 開獎动画
	EndAction             // 结算动画
	Wait                  // 等待遊戲開始
)

// 開獎類型
const (
	HUNG       = 1  // 紅
	HEI        = 2  // 黑
	NATURE     = 3  // 自然機率
	BAOZI      = 4  // 豹子
	SHUNJIN    = 5  // 順金
	JINHUA     = 6  // 金花
	SHUNZI     = 7  // 順子
	BIGDUIZI   = 8  // 對子(9-A)
	SMALLDUIZI = 9  // 對子(2-8)
	DANPAI     = 10 // 單張
)

//第一个用来占位
var resultstring = []string{"", "红", "黑", "致命一击"}

const (
	GameStartAnimation  = "gamestart"
	GameKaiJiang_Action = "kaijiang"
	GameEnd_Action      = "end"
)

const (
	CtrlBetLimitDefault int64 = 100 // 監控名單預設注額
)

type userScore struct {
	userdata []*gamelogic.GameUser
}

type HHDZLogic struct {
	*gamelogic.ComputeGameLogicBase
	serverfile         *os.File
	opttimer           *event.TimerFunc
	bettotal           int64                         //下注总额
	bettotalrobot      int64                         //機器人下注总额
	bettypetotal       []*bettype                    //下注分類汇总
	bettyperobot       []*bettype                    //機器人下注分類汇总
	betvalue           *betconfig                    //下注配置文件
	userdata           []*hhdzuuserdata              //用户数据
	gamestate          int32                         //遊戲状态
	bettime            int32                         //下注时间
	jiangtype          int32                         //開獎類型
	action_str         string                        //當前动画
	kaijiangactiontime int32                         //開獎动画超时
	endactiontime      int32                         //用户结算动画超时
	readytime          int32                         //用户准备时间
	zhuangID           int32                         //當前莊家ID
	zhuanglimit        int64                         //上莊限制
	lastturn           int32                         //莊剩余的局数
	maxturn            int32                         //上莊的最大局数
	zhuanglist         []int32                       //上莊列表
	betcount           []int32                       //下注数量
	zhuangWin          int64                         //莊家赢錢
	serverid           string                        //服务器ID
	winuserinfo        []int32                       //赢錢最多的前三家userid
	resultcardlist     *netproto.HHDZ_ResultCardList //開獎的历史牌
	zhushui            int32                         //注水率
	maxclearcount      int32                         //清理数据上限
	maxrobotzhuang     int32                         //最多有好多个機器人上莊
	robotxiazhuang     int32                         //機器人下莊概率
	robotzhuang        int32                         //没有人當莊的时候機器人當莊的概率
	zhushuimoney       int64                         //注水的錢
	he                 int32                         //下和的概率百分比
	startsyspool       int64                         //開始时的金币数量
	startmoney         int64                         //上莊时的金币
	currtax            int64                         //當前抽水
	nextjiangtype      int32                         //下次開獎類型
	nextzuhetype       int32                         //下次開獎牌型
	systotalwin        int64                         //系统收入
	syslastwin         int64                         //上一局收益
	sametimebetmode    int64                         //是否開啟莊閒同時下注
	systaxincome       int64                         //系統稅收, 非該局稅收
	userwinlose        int64                         //玩家輸贏
	robotwinlose       int64                         //機器人輸贏
	betlimitmoney      int64                         //下注限制
	publiczhuang       int64                         //是否開放用戶上莊
	userscore          *userScore                    //用户数据
	luckystar          int32                         //幸運星玩家id
	luckyfaceid        int32                         //幸運星玩家頭像id
	luckystarround     int64                         //幸運星持續局數
	luckystartimes     int64                         //幸運星第幾局
	notbettimes        int64                         //沒下注持續局數
	richestuserinfo    []int32                       //富豪玩家的userid
	richestList        []*netproto.HHDZ_Richest      //富豪榜
	stageuserinfo      []int32                       //標誌玩家的userid
	roundid            int32                         //當局局數
	hungwin            int32                         //紅贏局數
	heiwin             int32                         //黑贏局數
	poker              poker.Poker                   //整副牌
	cards              *poker.CardGroup              //所有的牌
	pokerHelper        *hhdzPokerHelper
	zuhetype           int32            //開獎贏的组合
	zuhehung           int32            //開獎紅的组合
	zuhehei            int32            //開獎黑的组合
	wintype            int32            //開獎贏的组合
	losetype           int32            //開獎輸的组合
	weight             []*betconfig     //权重配置文件
	scoreSection       []int64          //分数段用来加载不同的分数段配置
	lastGameUserBet    []*hhdzuuserdata //上局下注玩家
	//bowProvider        *gamelogic.BlackOrWhiteUserProvider   //需要放分或殺分的列表
	assignRewardDetail *netproto.SuperToolAssignRewardDetail //指定開獎
	execAssignReward   bool                                  //true 執行指定開獎
	fangFenTarget      int64                                 //放分金額
	fangFenComplete    int64                                 //己放分金額
	execFangFen        bool                                  //true 執行系統放分
	prettyrate         int32                                 //前兩張秀出好牌的概率
	SuperToolService   *gamelogic.SuperToolBrGameServiceBase //超端服務
	ctrlbetlimit       int64                                 //監控名單會控制的最低注額
}

// NewHHDZLogic 新建一个紅黑大戰实例
func NewHHDZLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *HHDZLogic {
	glc := new(HHDZLogic)
	glc.ComputeGameLogicBase = gamelogic.NewComputeGameLogicBase(t, tfManager, cfgs)
	glc.gethhdzconfig(cfgs)
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
	glc.resultcardlist = new(netproto.HHDZ_ResultCardList)
	//glc.getLastPool()
	glc.startsyspool = glc.SysPool.GetPool()
	glc.addCmd()
	glc.RegisterAfterGameStartFunc(func() {
		glc.clearGameData()   //清除数据
		glc.getAssignReward() //獲取多局指定開獎
		//glc.VipUserService.UpdateVIPUser() //獲取最新的監控用戶名單
		//glc.getTableBlackUser()
		//glc.getTableWhiteUser()
		glc.huangZhuang() //6/12 移除上莊機制
		glc.getLuckyStar()
		glc.getRichestList()
		glc.getGameRound()
		if glc.publiczhuang == 1 {
			glc.GameLog.AddGameStep("莊家%d，剩余%d", glc.zhuangID, glc.lastturn)
		}
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

func (hhdz *HHDZLogic) fapai() {
	//随机發牌
	if hhdz.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
		hhdz.poker = poker.NewRandomPoker(false, false, 1) //初始化一副牌
		hhdz.poker.Shuffle()
		//開局按概率表發牌
		//pokertype := []int32{SHUNJIN, JINHUA, SHUNZI, BIGDUIZI, SMALLDUIZI, DANPAI}
		//if hhdz.syspool < hhdz.winline {
		//	pokertype = []int32{SMALLDUIZI, DANPAI}
		//}
		//area := hhdz.getArea()
		//hungtype := hhdz.weight[area].getJiangType(pokertype)
		//heitype := hhdz.weight[area].getJiangType(pokertype)
		hhdz.pokerHelper = newPokerHelper(hhdz.poker)
		//hhdz.pokerHelper.initStorePokerType(hungtype, heitype)

		//直接發所有的牌
		//hhdz.cards = hhdz.poker.Deal(0, 52)

	} else { //發牌器發牌
		hhdz.poker = poker.NewFilePoker(hhdz.Configs, hhdz.AllTableUser())
		hhdz.poker = poker.NewFixedFilePoker(hhdz.poker.(*poker.FilePoker))
		hhdz.pokerHelper = newPokerHelper(hhdz.poker)
		//hhdz.pokerHelper.initStorePokerType(NATURE, NATURE)
	}
	//rnd := util.RndInt(0, 100)        // rand 0-99
	//if int32(rnd) < hhdz.prettyrate { //取出好看的兩張
	//	hhdz.pokerHelper.getPrettyFirstCards()
	//} else {
	//	hhdz.pokerHelper.getRandomFirstCards()
	//}
	//hhdz.BoardcastMessage(hhdz.getFirstPokerMsg())
}

func (hhdz *HHDZLogic) getArea() int32 {
	if hhdz.SysPool.GetPool() < hhdz.scoreSection[1] {
		return 0
	} else {
		return 1
	}
}

//获取配置表
func (hhdz *HHDZLogic) gethhdzconfig(cfgs *config.Group) {
	hhdz.betvalue = newBetconfig(cfgs, "hhdz")
	scfgs := cfgs.GetConfig("gameserver")
	hhdz.bettime = int32(scfgs.GetConfigInt("gameserver", "bettime"))
	hhdz.kaijiangactiontime = int32(scfgs.GetConfigInt("gameserver", "kaijiangactiontime"))
	hhdz.endactiontime = int32(scfgs.GetConfigInt("gameserver", "endactiontime"))
	hhdz.readytime = int32(scfgs.GetConfigInt("gameserver", "readywaittime"))
	hhdz.zhuanglimit = int64(scfgs.GetConfigInt("gameserver", "zhuanglimit"))
	hhdz.maxturn = int32(scfgs.GetConfigInt("gameserver", "zhuangturn"))
	hhdz.serverid = strconv.Itoa(int(hhdz.GetServerConfig().ServerID))
	hhdz.maxclearcount = int32(scfgs.GetConfigInt("gameserver", "maxclearcount"))
	hhdz.zhushui = int32(scfgs.GetConfigInt("gameserver", "zhushui"))
	hhdz.maxrobotzhuang = int32(scfgs.GetConfigInt("gameserver", "maxrobotzhuang"))
	hhdz.robotxiazhuang = int32(scfgs.GetConfigInt("gameserver", "robotxiazhuang"))
	hhdz.robotzhuang = int32(scfgs.GetConfigInt("gameserver", "robotzhuang"))
	hhdz.sametimebetmode = int64(scfgs.GetConfigInt("gameserver", "sametimebetmode"))
	hhdz.betlimitmoney = int64(scfgs.GetConfigInt("gameserver", "betlimitmoney"))
	hhdz.publiczhuang = int64(scfgs.GetConfigInt("gameserver", "publiczhuang"))
	hhdz.luckystarround = int64(scfgs.GetConfigInt("gameserver", "luckystarround"))
	hhdz.prettyrate = int32(scfgs.GetConfigInt("gameserver", "prettyrate"))
	betcountstr := scfgs.GetConfigValue("gameserver", "betcount")
	lastindex := 0
	for index, char := range betcountstr {
		if char == ',' {
			substr := betcountstr[lastindex:index]
			lastindex = index + 1
			count, _ := strconv.ParseInt(substr, 10, 32)
			hhdz.betcount = append(hhdz.betcount, int32(count))
		}
	}
	scorestr := scfgs.GetConfigValue("gameserver", "scoresection")
	lastindex = 0
	for index, char := range scorestr {
		if char == ',' {
			substr := scorestr[lastindex:index]
			lastindex = index + 1
			score, _ := strconv.ParseInt(substr, 10, 64)
			hhdz.scoreSection = append(hhdz.scoreSection, score)
			w := newBetconfig(cfgs, substr)
			hhdz.weight = append(hhdz.weight, w)
		}
	}

	hhdz.ctrlbetlimit = int64(scfgs.GetConfigInt("supertool", "ctrlbetlimit"))
	if hhdz.ctrlbetlimit <= 0 {
		hhdz.ctrlbetlimit = CtrlBetLimitDefault
	}
}

//添加命令行
func (hhdz *HHDZLogic) addCmd() {
	command.AddCommand(NewLhdCmd(hhdz))
}

//广播莊家、大富豪、幸運星
func (hhdz *HHDZLogic) broadCastZhuangInfo() {
	hhdz.BoardcastMessage(hhdz.getZhuangInfoMsg())
	hhdz.BoardcastMessage(hhdz.getRichestListMsg())
	hhdz.BoardcastMessage(hhdz.getLuckyStarMsg())
}

func (hhdz *HHDZLogic) startAnimation() {
	//開始下注前+2秒
	hhdz.RegisterAnimation(GameStartAnimation, int32(2), hhdz.startBet)
}

//获取上一次水池的錢
/*
func (hhdz *HHDZLogic) getLastPool() {
	if conf.TestMode {
		return
	}

	var folderName = "lastpool"
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err = os.MkdirAll(folderName, os.ModePerm)
	}

	var tablenostr = strconv.Itoa(int(hhdz.TableNo()))
	var filename = hhdz.serverid + "_" + tablenostr + ".conf"
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
		hhdz.serverfile = f
		defer hhdz.serverfile.Close()
		hhdz.readPoolFromFile()
	}
}

//读取水池文件
func (hhdz *HHDZLogic) readPoolFromFile() {
	if conf.TestMode {
		return
	}

	buff := make([]byte, 1024)
	n, _ := hhdz.serverfile.Read(buff)
	if n > 0 {
		poolstr := string(buff)[0:n]
		//pool, err := strconv.ParseInt(poolstr, 10, 32)
		pool, _ := strconv.Atoi(poolstr)
		hhdz.syspool = int64(pool)
	}
}
*/
//注水到下注额
func (hhdz *HHDZLogic) zhuShui() {
	hhdz.zhushuimoney = 0
	hhdz.GameLog.AddGameStep("系統水池: %d, 贏錢線: %d", hhdz.SysPool.GetPool(), hhdz.SysPool.GetWinLine())
}

//改变水池注水率
func (hhdz *HHDZLogic) changeZhuShui(value int32) {
	//if value >= 0 && value <= 100 {
	//	hhdz.zhushui = value
	//}
}
func (hhdz *HHDZLogic) checkSysPool() {
	mlog.Printf("桌子号：%v，系统水池：%v，注水概率：%v", hhdz.TableNo(), hhdz.SysPool.GetPool(), hhdz.zhushui)
}

//改变水池的錢
func (hhdz *HHDZLogic) changeSysPool(value int64) {
	change := value - hhdz.SysPool.GetPool()
	hhdz.SysPool.UpdateSysPool(change)
}

//改变機器人的和的下注概率
func (hhdz *HHDZLogic) changeHe(value int32) {
	if value >= 0 && value <= 100 {
		hhdz.he = value
	}
}

//清除数据
func (hhdz *HHDZLogic) clearGameData() {
	hhdz.userdata = make([]*hhdzuuserdata, 0)
	hhdz.bettypetotal = make([]*bettype, 0)
	hhdz.bettyperobot = make([]*bettype, 0)
	hhdz.winuserinfo = make([]int32, 0)
	hhdz.stageuserinfo = make([]int32, 0)
	hhdz.richestuserinfo = make([]int32, 0)
	hhdz.richestList = make([]*netproto.HHDZ_Richest, 0)
	hhdz.bettotal = 0
	hhdz.bettotalrobot = 0
	hhdz.zhuangWin = 0
	hhdz.zuhetype = 0
	hhdz.zuhehung = 0
	hhdz.zuhehei = 0
	hhdz.currtax = 0
	hhdz.robotwinlose = 0
	hhdz.userwinlose = 0
	hhdz.jiangtype = 0
	hhdz.assignRewardDetail = nil
	hhdz.execAssignReward = false
	hhdz.execFangFen = false
	rand.Seed(time.Now().UnixNano())
	zlen := len(hhdz.zhuanglist)
	if zlen == 0 {
		hhdz.robotZhuang()
	}

	hhdz.robotXiaZhuang()
}

//获取莊家信息消息
func (hhdz *HHDZLogic) getZhuangInfoMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZZhuangInfoID)

	pdata := new(netproto.HHDZ_ZhuangInfo)
	pdata.UserID = proto.Int32(0)
	pdata.LastTurn = proto.Int32(0)
	ur := hhdz.GetUserByUserID(hhdz.zhuangID)
	if ur != nil {
		pdata.UserInfo = append(pdata.UserInfo, ur.GetProtoBufGameUser())
		pdata.UserID = proto.Int32(ur.UserID)
		pdata.LastTurn = proto.Int32(hhdz.lastturn)
	}
	pmsg.MsgData = pdata
	return pmsg
}

//换莊
func (hhdz *HHDZLogic) huangZhuang() {
	if hhdz.publiczhuang == 1 {
		if hhdz.zhuangID == 0 {
			//莊是系统直接换莊
			hhdz.setNextZhuang()
			return
		}
		if hhdz.lastturn <= 0 {
			//已经到最后一局了
			hhdz.setNextZhuang()
			return
		}
		ur := hhdz.GetUserByUserID(hhdz.zhuangID)
		if ur != nil && ur.Score < hhdz.zhuanglimit {
			hhdz.setNextZhuang()
			return
		}

		if ur == nil {
			hhdz.setNextZhuang()
			return
		}

		hhdz.lastturn--
		mlog.Debug("1015發送莊家剩余局数%d", hhdz.lastturn)
	} else {
		mlog.Debug("未開啟上莊機制，不用換莊")
		return
	}
}

//设置下一个上莊的玩家
func (hhdz *HHDZLogic) setNextZhuang() {
	if hhdz.publiczhuang == 1 {
		hhdz.checkNetUserCanShangZhuang()
		if len(hhdz.zhuanglist) > 0 {
			hhdz.zhuangID = hhdz.zhuanglist[0]
			hhdz.lastturn = hhdz.maxturn
			hhdz.lastturn--
			hhdz.delFromZhuangList(hhdz.zhuangID)
			user := hhdz.GetUserByUserID(hhdz.zhuangID)
			hhdz.startmoney = user.Score
		} else {
			hhdz.zhuangID = 0
			hhdz.startmoney = hhdz.startsyspool
		}
	} else {
		mlog.Error("未開啟上莊機制，不得上莊")
		return
	}
}

//检测下一个用户是否能上莊
func (hhdz *HHDZLogic) checkNetUserCanShangZhuang() {
	if len(hhdz.zhuanglist) > 0 {
		for {
			user := hhdz.GetUserByUserID(hhdz.zhuanglist[0])
			if user != nil && user.Score >= int64(hhdz.zhuanglimit) {
				break
			}

			hhdz.delFromZhuangList(hhdz.zhuanglist[0])
			if len(hhdz.zhuanglist) <= 0 {
				break
			}
		}
	}
}

//删除一个用户从上莊列表
func (hhdz *HHDZLogic) delFromZhuangList(userid int32) {
	for i := 0; i < len(hhdz.zhuanglist); i++ {
		if hhdz.zhuanglist[i] == userid {
			hhdz.zhuanglist = append(hhdz.zhuanglist[:i], hhdz.zhuanglist[i+1:]...)
			break
		}
	}
}

//機器人上莊
func (hhdz *HHDZLogic) robotZhuang() {
	if hhdz.publiczhuang == 1 {
		//概率當莊
		if rand.Int31n(100) > hhdz.robotzhuang {
			return
		}
		//统计機器人个数
		robotarr := make([]int32, 0)
		for _, ur := range hhdz.AllTableUser() {
			if ur.IsRobot() {
				robotarr = append(robotarr, ur.UserID)
			}
		}

		//沒有機器人或沒有設定最多幾個機器人上莊
		if len(robotarr) <= 0 || hhdz.maxrobotzhuang <= 0 {
			return
		}

		//随机有几个人當莊
		n := rand.Int31n(hhdz.maxrobotzhuang) + 1

		for i := int32(0); i < n; i++ {
			rcn := len(robotarr)
			if rcn <= 0 {
				break
			}
			index := rand.Int31n(int32(rcn))
			hhdz.addToZhuangList(robotarr[index])
			robotarr = append(robotarr[:index], robotarr[index+1:]...)
		}
	} else {
		mlog.Debug("未開啟上莊機制，機器人不用上莊")
		return
	}
}

//添加一个用户到上莊列表
func (hhdz *HHDZLogic) addToZhuangList(userid int32) {
	hhdz.zhuanglist = append(hhdz.zhuanglist, userid)
}

//機器人下莊
func (hhdz *HHDZLogic) robotXiaZhuang() {
	if hhdz.zhuangID == 0 {
		return
	}

	ur := hhdz.GetUserByUserID(hhdz.zhuangID)
	if ur != nil && ur.IsRobot() {
		if rand.Int31n(100) < hhdz.robotxiazhuang {
			hhdz.zhuangID = 0
		}
	}
}

//排序身上攜帶的錢大到小
func (hhdz *HHDZLogic) getMoneyRanking() {
	sort.Sort(hhdz.userscore)
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
func (hhdz *HHDZLogic) getLuckyStar() {
	t := 0
	gameuser := hhdz.AllGameUser()
	if hhdz.luckystar == 0 || hhdz.notbettimes >= 3 {
		//如果沒有幸運星、連續三場沒下注或是幸運星的持續局數到了，隨機一個玩家
		t = 1
	} else {
		//檢查玩家還在不在場
		for _, u := range gameuser {
			if u.UserID == hhdz.luckystar {
				hhdz.luckystartimes += 1
				break
			} else {
				t = 1
			}
		}
	}
	if t == 1 || hhdz.luckystartimes == hhdz.luckystarround {
		if len(gameuser) != 0 {
			userIndex := rand.Intn(len(gameuser))
			hhdz.luckystar = gameuser[userIndex].UserID
			hhdz.luckyfaceid = gameuser[userIndex].FaceID
		} else {
			//如果沒有明燈，須和client確認要傳什麼
			hhdz.luckystar = 0
			hhdz.luckyfaceid = 0
		}
		hhdz.luckystartimes = 0
		hhdz.notbettimes = 0
	}
}

//獲取富豪玩家
func (hhdz *HHDZLogic) getRichestList() {
	userscore := new(userScore)
	userscore.userdata = hhdz.AllTableUser()
	sort.Sort(userscore)
	//計算分數最高的7個
	for index, user := range userscore.userdata {
		if index < 5 {
			rich := new(netproto.HHDZ_Richest)
			rich.Ranking = proto.Int32(int32(index + 1))
			rich.UserID = proto.Int32(user.UserID)
			rich.UserFaceID = proto.Int32(user.FaceID)
			rich.Money = proto.Int64(user.Score)
			hhdz.richestList = append(hhdz.richestList, rich)
			hhdz.richestuserinfo = append(hhdz.richestuserinfo, user.UserID)
		} else {
			break
		}
	}
	hhdz.stageuserinfo = append(hhdz.richestuserinfo, hhdz.luckystar)
}
func (hhdz *HHDZLogic) getGameRound() {
	if hhdz.roundid == 30 {
		hhdz.hungwin = 0
		hhdz.heiwin = 0
	}
}

//获取大富豪列表的消息
func (hhdz *HHDZLogic) getRichestListMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZRichestListID)
	pdata := new(netproto.HHDZ_RichestList)
	pdata.RichestList = hhdz.richestList
	//如果當下沒有玩家 回傳空值
	pmsg.MsgData = pdata
	return pmsg
}

//获取幸運星的消息
func (hhdz *HHDZLogic) getLuckyStarMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZLuckyStarListID)
	luckyStarList := new(netproto.HHDZ_LuckyStar)
	luckyStarList.UserID = proto.Int32(hhdz.luckystar)
	luckyStarList.UserFaceID = proto.Int32(hhdz.luckyfaceid)
	pmsg.MsgData = luckyStarList
	return pmsg
}

func (hhdz *HHDZLogic) getRoundMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZRoundID)
	round := new(netproto.HHDZ_RoundID)
	round.RoundID = proto.Int32(hhdz.roundid)
	round.HungWin = proto.Int32(hhdz.hungwin)
	round.HeiWin = proto.Int32(hhdz.heiwin)
	pmsg.MsgData = round
	return pmsg
}

//用户離開
func (hhdz *HHDZLogic) UserLeave(gameuser *gamelogic.GameUser) bool {
	if hhdz.gamestate != Wait {
		if hhdz.getUserDataByID(gameuser.UserID) != nil {
			mlog.Debug("用户在玩遊戲還不能離開")
			return false
		}

		if hhdz.zhuangID == gameuser.UserID {
			mlog.Debug("用户在當莊，所以用户不能離開")
			return false
		}
	}

	hhdz.delFromZhuangList(gameuser.UserID)

	hhdz.ChangeUserStateEnd(gameuser)
	return hhdz.ComputeGameLogicBase.UserLeave(gameuser)
}

//注销用户实例
func (hhdz *HHDZLogic) UserFree(gameuser *gamelogic.GameUser) {
	if hhdz.gamestate != Wait {
		if hhdz.getUserDataByID(gameuser.UserID) != nil {
			mlog.Debug("用户在玩遊戲還不能離開")
			return
		}

		if hhdz.zhuangID == gameuser.UserID {
			mlog.Debug("用户在當莊，所以用户不能離開")
			return
		}
	} else {
		if hhdz.zhuangID == gameuser.UserID {
			hhdz.zhuangID = 0
		}
	}

	hhdz.delFromZhuangList(gameuser.UserID)

	hhdz.ChangeUserStateEnd(gameuser)
	hhdz.ComputeGameLogicBase.UserFree(gameuser)
}

//断线重连
func (hhdz *HHDZLogic) UserRecover(gameuser *gamelogic.GameUser) {
	mlog.Debug("UserRecover")
	hhdz.ComputeGameLogicBase.UserRecover(gameuser)
	ud := hhdz.getUserDataByID(gameuser.UserID)
	if ud != nil {
		ud.money += ud.user.Score - ud.moneylast
		ud.moneylast = ud.user.Score
	}

	hhdz.sendScene(gameuser)
}

//用户进入
func (hhdz *HHDZLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {
	hhdz.ComputeGameLogicBase.UserEnter(gameuser, tableno, seatno)
	mlog.Debug("UserEnter")
	hhdz.sendScene(gameuser)
}

//發送場景
func (hhdz *HHDZLogic) sendScene(gameuser *gamelogic.GameUser) {
	//發送设置

	if gameuser.IsRobot() {
		mlog.Debug("给機器人%d發送場景消息", gameuser.UserID)
	} else {
		mlog.Debug("给玩家%d發送場景消息", gameuser.UserID)
	}
	hhdz.sendSetting(gameuser)
	hhdz.SendMessageToUser(gameuser, hhdz.getRoundMsg())
	hhdz.SendMessageToUser(gameuser, hhdz.getResultCardListMsg())
	hhdz.SendMessageToUser(gameuser, hhdz.getRichestListMsg())
	hhdz.SendMessageToUser(gameuser, hhdz.getLuckyStarMsg())
	if hhdz.gamestate == Wait {
		return
	}
	//發送阶段
	hhdz.SendMessageToUser(gameuser, hhdz.getStateMsg())
	hhdz.SendMessageToUser(gameuser, hhdz.getZhuangInfoMsg())
	hhdz.SendMessageToUser(gameuser, hhdz.getBetTotalMsg(0))
	hhdz.sendUserBetInfo(gameuser)
	if hhdz.gamestate == Bet {
		//hhdz.SendMessageToUser(gameuser, hhdz.getFirstPokerMsg())
	}
	if hhdz.gamestate == KaiJiangAction {
		//發送開獎信息
		hhdz.SendMessageToUser(gameuser, hhdz.getGameResultPokerMsg())
	} else if hhdz.gamestate == EndAction {
		//發送用户下注输赢信息如果有
		ud := hhdz.getUserDataByID(gameuser.UserID)
		if ud != nil {
			hhdz.SendMessageToUser(gameuser, hhdz.getUserWinMsg(ud.getShowWinLose(), hhdz.zhuangWin))
		}
	}
}

//發送设置
func (hhdz *HHDZLogic) sendSetting(gameuser *gamelogic.GameUser) {
	mlog.Debug("sendSetting")
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZSettingID)

	pdata := new(netproto.HHDZ_Setting)
	pdata.BetCount = append(pdata.BetCount, hhdz.betcount...)
	pdata.BetMoney = proto.Int32(int32(hhdz.betvalue.getMoneyByType(1)))
	pdata.ZhuangLimit = proto.Int64(hhdz.zhuanglimit)
	pdata.Betlimitmoney = proto.Int64(hhdz.betlimitmoney)
	for i := 0; i < 3; i++ {
		pdata.BetMax = append(pdata.BetMax, hhdz.betvalue.betvalue[i].maxmoney)
	}
	pmsg.MsgData = pdata
	hhdz.SendMessageToUser(gameuser, pmsg)
}

//获取用户数据通过用户ID
func (hhdz *HHDZLogic) getUserDataByID(userid int32) *hhdzuuserdata {
	for _, ud := range hhdz.userdata {
		if ud.user.UserID == userid {
			return ud
		}
	}

	return nil
}

//获取開獎历史列表的消息
func (hhdz *HHDZLogic) getResultCardListMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZResultCardListID)

	pdata := new(netproto.HHDZ_ResultCardList)
	if len(hhdz.resultcardlist.CardList) > 30 { //大於 30 取最後三十個
		pdata.CardList = append(pdata.CardList, hhdz.resultcardlist.CardList[len(hhdz.resultcardlist.CardList)-30:]...)
	} else {
		pdata.CardList = append(pdata.CardList, hhdz.resultcardlist.CardList...)
	}
	mlog.Debug("getResultCardListMsg:%v", pdata.CardList)
	pmsg.MsgData = pdata
	return pmsg
}

//获取状态消息
func (hhdz *HHDZLogic) getStateMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZStateID)

	pdata := new(netproto.HHDZ_Sate)
	pdata.State = proto.Int32(hhdz.gamestate)

	pdata.Time = proto.Int32(hhdz.getBetTime())
	pmsg.MsgData = pdata
	return pmsg
}

//取得下注時間
func (hhdz *HHDZLogic) getBetTime() int32 {
	time := int32(0)
	if hhdz.opttimer != nil {
		time = hhdz.opttimer.WaitSeconds - hhdz.opttimer.ElapsedSeconds()
	}
	return time
}

//获取下注总额消息
func (hhdz *HHDZLogic) getBetTotalMsg(userid int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBroadCastBetID)

	pdata := new(netproto.HHDZ_BroadCastBetArr)
	pdata.UserID = proto.Int32(userid)

	lucky := hhdz.getUserDataByID(hhdz.luckystar) //獲取明燈

	for _, v := range hhdz.betvalue.alltype {
		betinfo := new(netproto.HHDZ_BroadCastBet)
		betinfo.Type = proto.Int32(v)
		betinfo.Money = proto.Int64(hhdz.getBetMoneyByType(v))
		betinfo.LuckyStar = proto.Bool(false)
		if lucky != nil && lucky.getUserBetByType(v) > 0 {
			betinfo.LuckyStar = proto.Bool(true)
		}
		pdata.BroadCastBetArr = append(pdata.BroadCastBetArr, betinfo)
	}
	pmsg.MsgData = pdata

	return pmsg
}

//获取此類型的下注总额 真人+機器人
func (hhdz *HHDZLogic) getBetMoneyByType(t int32) int64 {
	allBetTotal := int64(0)
	for i := 0; i < len(hhdz.bettypetotal); i++ {
		if hhdz.bettypetotal[i].hhdztype == t {
			allBetTotal += hhdz.bettypetotal[i].money
		}
	}
	for i := 0; i < len(hhdz.bettyperobot); i++ {
		if hhdz.bettyperobot[i].hhdztype == t {
			allBetTotal += hhdz.bettyperobot[i].money
		}
	}
	return allBetTotal
}

//获取此類型的下注总额 真人
func (hhdz *HHDZLogic) getRealManBetMoneyByType(t int32) int64 {
	allBetTotal := int64(0)
	for i := 0; i < len(hhdz.bettypetotal); i++ {
		if hhdz.bettypetotal[i].hhdztype == t {
			allBetTotal += hhdz.bettypetotal[i].money
		}
	}
	return allBetTotal
}

//發送用户的下注信息
func (hhdz *HHDZLogic) sendUserBetInfo(gameuser *gamelogic.GameUser) {
	ud := hhdz.getUserDataByID(gameuser.UserID)
	if ud == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetInfoID)

	pdata := new(netproto.HHDZ_UserBetInfo)
	pdata.UserMoney = proto.Int64(ud.money)

	lucky := hhdz.getUserDataByID(hhdz.luckystar) //獲取明燈

	for _, v := range ud.userbet {
		betinfo := new(netproto.HHDZ_BroadCastBet)
		betinfo.Type = proto.Int32(v.hhdztype)
		betinfo.Money = proto.Int64(v.money)
		betinfo.LuckyStar = proto.Bool(false)
		if lucky != nil && lucky.getUserBetByType(v.hhdztype) > 0 {
			betinfo.LuckyStar = proto.Bool(true)
		}
		pdata.BetInfo = append(pdata.BetInfo, betinfo)
	}
	pmsg.MsgData = pdata
	hhdz.SendMessageToUser(gameuser, pmsg)
}

//获取玩家赢錢消息
func (hhdz *HHDZLogic) getUserWinMsg(selfwin int64, zhuangwin int64) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZGameResultWinID)

	pdata := new(netproto.HHDZ_GameResultUserMoney)
	pdata.ZhuangWin = proto.Int64(zhuangwin)
	pdata.SelfWin = proto.Int64(selfwin)
	for i := 0; i < len(hhdz.winuserinfo); i++ {
		ud := hhdz.getUserDataByID(hhdz.winuserinfo[i])
		if ud == nil {
			continue
		}
		userwin := new(netproto.HHDZ_WinUerInfo)
		userid := strconv.Itoa(int(ud.user.UserID))
		userwin.UserName = proto.String(userid)
		userwin.Win = proto.Int64(ud.getShowWinLose())
		pdata.WinUserInfo = append(pdata.WinUserInfo, userwin)
	}
	for j := 0; j < len(hhdz.stageuserinfo); j++ {
		ud := hhdz.getUserDataByID(hhdz.stageuserinfo[j])
		if ud == nil {
			continue
		}
		stageuserwin := new(netproto.HHDZ_StageUerInfo)
		stageuserwin.UserID = proto.Int32(ud.user.UserID)
		stageuserwin.Win = proto.Int64(ud.getShowWinLose())
		pdata.StageUerInfo = append(pdata.StageUerInfo, stageuserwin)
	}
	pmsg.MsgData = pdata
	return pmsg
}

//获取開獎结果撲克牌消息
func (hhdz *HHDZLogic) getGameResultPokerMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZGameResultPokerID)

	pdata := new(netproto.HHDZ_GameResultPoker)
	for _, v := range hhdz.pokerHelper.hungStore.cards {
		cards := new(netproto.HHDZ_Poker)
		cards.Color = proto.Int32(v.Color)
		if v.CardNumber == 14 {
			cards.Number = proto.Int32(1)
		} else {
			cards.Number = proto.Int32(v.CardNumber)
		}
		pdata.HungPoker = append(pdata.HungPoker, cards)
	}

	for _, v := range hhdz.pokerHelper.heiStore.cards {
		cards := new(netproto.HHDZ_Poker)
		cards.Color = proto.Int32(v.Color)
		if v.CardNumber == 14 {
			cards.Number = proto.Int32(1)
		} else {
			cards.Number = proto.Int32(v.CardNumber)
		}
		pdata.HeiPoker = append(pdata.HeiPoker, cards)
	}

	pdata.RoundID = proto.Int32(hhdz.roundid)
	pdata.HungZuHe = proto.Int32(hhdz.zuhehung)
	pdata.HeiZuHe = proto.Int32(hhdz.zuhehei)
	pdata.Win = proto.Int32(hhdz.jiangtype)
	pmsg.MsgData = pdata
	return pmsg
}

//获取先發兩張撲克牌消息
func (hhdz *HHDZLogic) getFirstPokerMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZFapaiFirstCardsID)

	pdata := new(netproto.HHDZ_FapaiFirstCards)
	for _, v := range hhdz.pokerHelper.hungStore.cards {
		cards := new(netproto.HHDZ_Poker)
		cards.Color = proto.Int32(v.Color)
		if v.CardNumber == 14 {
			cards.Number = proto.Int32(1)
		} else {
			cards.Number = proto.Int32(v.CardNumber)
		}

		pdata.HungPoker = append(pdata.HungPoker, cards)
	}

	for _, v := range hhdz.pokerHelper.heiStore.cards {
		cards := new(netproto.HHDZ_Poker)
		cards.Color = proto.Int32(v.Color)
		if v.CardNumber == 14 {
			cards.Number = proto.Int32(1)
		} else {
			cards.Number = proto.Int32(v.CardNumber)
		}
		pdata.HeiPoker = append(pdata.HeiPoker, cards)
	}

	pmsg.MsgData = pdata
	return pmsg
}

//開始下注
func (hhdz *HHDZLogic) startBet() {
	hhdz.fapai()
	hhdz.opttimer = hhdz.TimerFuncManager.NewTimer(hhdz.bettime, hhdz.kaiJiang)
	hhdz.setState(Bet)
	hhdz.BoardcastMessage(hhdz.getStateMsg())
}

//设置状态
func (hhdz *HHDZLogic) setState(state int32) {
	hhdz.gamestate = state
}

//開始開獎
func (hhdz *HHDZLogic) kaiJiang() {

	// 統計所有玩家下注
	for _, ud := range hhdz.userdata {
		ud.betAmount = 0

		betMap := make(map[int32]int64)
		for _, bet := range ud.userbet {
			betMap[bet.hhdztype] += bet.money
			ud.betAmount += bet.money
		}

		isRobot := ""
		if ud.user.IsRobot() {
			isRobot = "[機器人]"
			continue
		}

		hhdz.GameLog.AddGameStep("%s玩家[%d] 押注總計: %d, 紅方: %d, 幸運一擊: %d, 黑方: %d", isRobot, ud.user.UserID, ud.betAmount, betMap[HUNG], betMap[NATURE], betMap[HEI])
	}

	hhdz.setState(KaiJiang)

	if hhdz.zhuangID != 0 {
		if hhdz.GetUserByUserID(hhdz.zhuangID).IsRobot() {
			//hhdz.removeRobotBet()
		}
	} else {
		//hhdz.removeRobotBet()
	}

	hhdz.getResultType()
	hhdz.getPokerByType(hhdz.jiangtype)
	hhdz.countGameRound()
	hhdz.BoardcastMessage(hhdz.getGameResultPokerMsg())
	hhdz.BoardcastMessage(hhdz.getRoundMsg())
	hhdz.setState(KaiJiangAction)
	hhdz.addResultCardList()
	hhdz.action_str = hhdz.RegisterAnimation(GameKaiJiang_Action, hhdz.kaijiangactiontime, hhdz.sendResult)
}

//移除機器人下注额
func (hhdz *HHDZLogic) removeRobotBet() {
	for _, v := range hhdz.bettyperobot {
		hhdz.updateBet(v.hhdztype, -v.money)
	}
}

//更新下注錢
func (hhdz *HHDZLogic) updateBet(t int32, money int64) {
	hhdz.bettotal += int64(money)
	mlog.Debug("BetType=%v,money=%v,lhd.bettotal=%v", t, money, hhdz.bettotal)
	for _, v := range hhdz.bettypetotal {
		if v.hhdztype == t {
			v.money += int64(money)
			if v.money < 0 {
				mlog.Debug("錢小於0了！")
				v.money = 0
			}
			return
		}
	}

	b := new(bettype)
	b.hhdztype = t
	b.money = int64(money)
	if b.money < 0 {
		mlog.Debug("錢小於0了！")
		b.money = 0
	}
	hhdz.bettypetotal = append(hhdz.bettypetotal, b)
}

//計算回合
func (hhdz *HHDZLogic) countGameRound() {
	if hhdz.roundid >= 30 {
		hhdz.roundid = 1
	} else {
		hhdz.roundid += 1
	}
	if hhdz.jiangtype == HUNG {
		hhdz.hungwin += 1
	} else if hhdz.jiangtype == HEI {
		hhdz.heiwin += 1
	}
	mlog.Debug("開獎類型：%v", hhdz.jiangtype)
}

//添加開獎到历史列表
func (hhdz *HHDZLogic) addResultCardList() {
	retcards := new(netproto.HHDZ_GameResultPoker)
	for _, v := range hhdz.pokerHelper.hungStore.cards {
		cards := new(netproto.HHDZ_Poker)
		cards.Color = proto.Int32(v.Color)
		if v.CardNumber == 14 {
			cards.Number = proto.Int32(1)
		} else {
			cards.Number = proto.Int32(v.CardNumber)
		}
		retcards.HungPoker = append(retcards.HungPoker, cards)
	}

	for _, v := range hhdz.pokerHelper.heiStore.cards {
		cards := new(netproto.HHDZ_Poker)
		cards.Color = proto.Int32(v.Color)
		if v.CardNumber == 14 {
			cards.Number = proto.Int32(1)
		} else {
			cards.Number = proto.Int32(v.CardNumber)
		}
		retcards.HeiPoker = append(retcards.HeiPoker, cards)
	}

	retcards.RoundID = proto.Int32(hhdz.roundid)
	retcards.HungZuHe = proto.Int32(hhdz.zuhehung)
	retcards.HeiZuHe = proto.Int32(hhdz.zuhehei)
	retcards.Win = proto.Int32(hhdz.jiangtype)

	//回報指定多局開獎結果
	if hhdz.assignRewardDetail != nil {
		example.invalid("===========%v,%v", hhdz.jiangtype, hhdz.zuhetype)
		if hhdz.isPowerJiangType(hhdz.zuhetype) {
			hhdz.assignRewardDetailResult(fmt.Sprintf("%d,%d", hhdz.jiangtype, 3))
		} else {
			hhdz.assignRewardDetailResult(fmt.Sprintf("%d", hhdz.jiangtype))
		}
	}

	hhdz.resultcardlist.CardList = append(hhdz.resultcardlist.CardList, retcards)
	if len(hhdz.resultcardlist.CardList) > 100 {
		hhdz.resultcardlist.CardList = append(hhdz.resultcardlist.CardList[:0], hhdz.resultcardlist.CardList[1:]...)
	}
}

//获取開獎類型。
func (hhdz *HHDZLogic) getResultType() {
	superToolHandler := hhdz.SuperToolService.GetBrGameSuperToolHandler()

	if hhdz.nextjiangtype != 0 { //指定開獎
		hhdz.jiangtype = hhdz.nextjiangtype
		hhdz.zuhetype = hhdz.nextzuhetype
		hhdz.GameLog.AddGameStep("指定開獎: %v, %v", hhdz.jiangtype, hhdz.zuhetype)
		hhdz.nextjiangtype = 0
		hhdz.nextzuhetype = 0
		return
	} else if ok, rettype := hhdz.isAssignReward(); ok { //多局指定開獎
		hhdz.execAssignReward = true
		if len(rettype) > 0 {
			hhdz.jiangtype = rettype[0]
		}
		if len(rettype) > 1 {
			if rettype[1] == 3 { //組合牌型
				pokertype := []int32{SHUNJIN, JINHUA, SHUNZI, BIGDUIZI}
				area := hhdz.getArea()
				hhdz.zuhetype = hhdz.weight[area].getJiangType(pokertype)
			}
		}
		hhdz.GameLog.AddGameStep("多局指定開獎:rettype %v, jiangtype %v, zhuetype %v", rettype, hhdz.jiangtype, hhdz.zuhetype)
		return
	} else if hhdz.fangFenTarget > 0 && hhdz.fangFenComplete < hhdz.fangFenTarget { //系統放分
		hhdz.execFangFen = true
		alltype := hhdz.getFangFenJiangType()
		hhdz.jiangtype = hhdz.weight[1].getJiangType(alltype)
		pokertype := []int32{SHUNJIN, JINHUA, SHUNZI, BIGDUIZI, SMALLDUIZI, DANPAI}
		hhdz.zuhetype = hhdz.weight[1].getJiangType(pokertype)
		hhdz.GameLog.AddGameStep("系統放分 : target %d, already %d", hhdz.fangFenTarget, hhdz.fangFenComplete)
		return
	} else if len(hhdz.getBowBet()) > 0 { //監控放殺分
		if hhdz.bowJiangType() {
			example.invalid("執行監控放殺分策略 : %s", superToolHandler.GetProiorityListString())
			hhdz.GameLog.AddGameStep("執行監控放殺分策略 : %s", superToolHandler.GetProiorityListString())
			return
		}
	} else if hhdz.noRealUserBet() { //沒有真人下注時, 開最漂高的機率表
		hhdz.jiangtype = 0
		hhdz.GameLog.AddGameStep("無真人玩家下注, jiangtype [%d], zhueytpe [%d]", hhdz.jiangtype, hhdz.zuhetype)
		return
	}
	//mlog.Debug("=====================%d", len(hhdz.getBowBet()))
	alltype := []int32{HUNG, HEI}
	pokertype := []int32{SHUNJIN, JINHUA, SHUNZI, BIGDUIZI, SMALLDUIZI, DANPAI}
	area := hhdz.getArea()
	hhdz.GameLog.AddGameStep("水池狀態 : " + hhdz.SysPool.ToString())
	if hhdz.SysPool.IsCtrlSystemWin() {
		//當水線低於贏錢線，開最小獎
		alltype = hhdz.getMinKaiJiangType()
		pokertype = []int32{SMALLDUIZI, DANPAI}
		hhdz.GameLog.AddGameStep("水線之下, 開獎:%v", alltype)
		mlog.Warn("水線之下, 開獎:%v", alltype)
	}

	hhdz.jiangtype = hhdz.weight[area].getJiangType(alltype)
	hhdz.zuhetype = hhdz.weight[area].getJiangType(pokertype)
	//根據權重獲得贏家的牌型
	hhdz.GameLog.AddGameStep("開獎類型:%v, 牌型選擇:%v", hhdz.jiangtype, hhdz.zuhetype)
}

//获取赔偿最小的開獎结果
func (hhdz *HHDZLogic) getMinKaiJiangType() []int32 {
	if hhdz.getRealManBetMoneyByType(HUNG) > hhdz.getRealManBetMoneyByType(HEI) {
		return []int32{HEI}
	} else if hhdz.getRealManBetMoneyByType(HUNG) < hhdz.getRealManBetMoneyByType(HEI) {
		return []int32{HUNG}
	}
	return []int32{HUNG, HEI}
}

//获取赔偿最大的開獎结果
func (hhdz *HHDZLogic) getFangFenJiangType() []int32 {
	if hhdz.getRealManBetMoneyByType(HUNG) > hhdz.getRealManBetMoneyByType(HEI) {
		return []int32{HUNG}
	} else if hhdz.getRealManBetMoneyByType(HUNG) < hhdz.getRealManBetMoneyByType(HEI) {
		return []int32{HEI}
	}
	return []int32{HUNG, HEI}
}

//獲取放分組合類型
func (hhdz *HHDZLogic) getFangFenZuheType() (zuhetype []int32) {
	//放分開致命一擊的條件及機率
	alltype := []int32{BAOZI, SHUNJIN, JINHUA, SHUNZI, BIGDUIZI}
	rnd := util.RndInt(0, 1000)
	if rnd < 165 {
		r := (hhdz.fangFenTarget * 15 / 10) - hhdz.fangFenComplete
		for _, zuhe := range alltype {
			zuhepay := hhdz.getRealManBetMoneyByType(3) * int64(hhdz.betvalue.getOddByType(zuhe))
			if r > zuhepay {
				zuhetype = append(zuhetype, zuhe)
			}
		}
	}
	if len(zuhetype) < 1 {
		zuhetype = append(alltype, []int32{SMALLDUIZI, DANPAI}...)
	}
	return
}

//获取開獎的撲克信息
func (hhdz *HHDZLogic) getPokerByType(t int32) {
	//pokertype := []int32{SHUNJIN, JINHUA, SHUNZI, BIGDUIZI, SMALLDUIZI, DANPAI}
	//if hhdz.syspool < hhdz.winline {
	//	pokertype = []int32{SMALLDUIZI, DANPAI}
	//}
	//area := hhdz.getArea()
	//hungtype := hhdz.weight[area].getJiangType(pokertype)
	//heitype := hhdz.weight[area].getJiangType(pokertype)
	//hhdz.pokerHelper = newPokerHelper(hhdz.poker)
	//hhdz.pokerHelper.initStorePokerType(hungtype, heitype)
	hungtype := int32(DANPAI)
	heitype := int32(DANPAI)
	rnd := util.RndInt(0, 2) // 0-1
	winUserID := []int32{_HungUserID, _HeiUserID}[rnd]
	if hhdz.jiangtype == 0 { //jiangtype 0 表示無人下注，開最漂高的牌型機率表
		hungtype = hhdz.weight[1].getJiangType([]int32{SHUNJIN, JINHUA, SHUNZI, BIGDUIZI, SMALLDUIZI, DANPAI})
		heitype = hhdz.weight[1].getJiangType([]int32{SHUNJIN, JINHUA, SHUNZI, BIGDUIZI, SMALLDUIZI, DANPAI})
		winUserID = 0
		hhdz.GameLog.AddGameStep("無真人玩家下注, winUserID [%d], hungtype [%d], heitype [%d]", winUserID, hungtype, heitype)
	} else if hhdz.jiangtype == HUNG {
		winUserID = _HungUserID
		hungtype = hhdz.zuhetype
	} else if hhdz.jiangtype == HEI {
		winUserID = _HeiUserID
		heitype = hhdz.zuhetype
	} else {
		mlog.Error("jiangtype %d 指定錯誤，隨機一方贏", hhdz.jiangtype)
	}

	hhdz.pokerHelper.initStorePokerType(hungtype, heitype)
	hhdz.pokerHelper.getPokerCardByWinUserIDFromStore(winUserID)

	hhdz.jiangtype = int32(hhdz.pokerHelper.getJiangType())
	hungz, heiz, z := hhdz.pokerHelper.getZuheTypes()

	hhdz.zuhehung = int32(hungz)
	hhdz.zuhehei = int32(heiz)
	hhdz.zuhetype = int32(z)

	hhdz.GameLog.AddGameStep("開獎结果是 %d(%d), 紅區的撲克牌是%v, 黑區的撲克牌是%v",
		hhdz.jiangtype, hhdz.zuhetype, hhdz.pokerHelper.hungStore.cards, hhdz.pokerHelper.heiStore.cards)

	mlog.Debug("開獎结果是 %d(%d), 紅區的撲克牌是%v, 黑區的撲克牌是%v", hhdz.jiangtype, hhdz.zuhetype, hhdz.pokerHelper.hungStore.cards, hhdz.pokerHelper.heiStore.cards)
}

//给玩家發送结果
func (hhdz *HHDZLogic) sendResult() {
	hhdz.calWinLose()

	//hhdz.writePoolToFile()
	hhdz.updateVipUserWin()
	hhdz.action_str = hhdz.RegisterAnimation(GameEnd_Action, hhdz.endactiontime, hhdz.gameEnd)
}

//用户消息主处理函数
func (hhdz *HHDZLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	if !hhdz.CheckGameUserMessage(user) {
		hhdz.KickUserFunc(user.UserID, "非法操作")
		return
	}
	//不是给遊戲逻辑服务器的消息
	mlog.Debug("當前遊戲状态为：%v", hhdz.gamestate)
	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	hhdz.CheckUserAnimationCommited(user, gmsg, hhdz.action_str, nil)

	if gmsg.SClassID == int32(netproto.HHDZ_GameMessageClassID_HHDZGetShangZhuangListID) {
		hhdz.SendMessageToUser(user, hhdz.getZhuangListMsg())
		return
	}

	if hhdz.gamestate != Bet {
		return
	}

	if gmsg.SClassID == int32(netproto.HHDZ_GameMessageClassID_HHDZBetID) {
		hhdz.userBet(user, gmsg.MsgData.([]byte))
		return
	}

	if gmsg.SClassID == int32(netproto.HHDZ_GameMessageClassID_HHDZShangZhuangID) {
		hhdz.shangZhuang(user, gmsg.MsgData.([]byte))
		return
	}

	if gmsg.SClassID == int32(netproto.HHDZ_GameMessageClassID_HHDZBetClearID) && hhdz.gamestate == Bet {
		hhdz.clearUserBet(user)
		return
	}
}

//将水池的錢写入文件
/*
func (hhdz *HHDZLogic) writePoolToFile() {
	if conf.TestMode {
		return
	}
	buff := strconv.Itoa(int(hhdz.syspool))
	//bufio.ReadWriter.Peek()
	tablenostr := strconv.Itoa(int(hhdz.TableNo()))
	filename := "./lastpool/" + hhdz.serverid + "_" + tablenostr + ".conf"
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
*/
//遊戲结束
func (hhdz *HHDZLogic) gameEnd() {
	hhdz.setState(Wait)
	if hhdz.zhuangID != 0 {
		user := hhdz.GetUserByUserID(hhdz.zhuangID)
		if user.IsOffLine() {
			hhdz.zhuangID = 0
		}
	}
	hhdz.GameEndFunc()
}

//获取上莊列表消息
func (hhdz *HHDZLogic) getZhuangListMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZShangZhuangListInfoID)

	pdata := new(netproto.HHDZ_ZhuangList)
	for i := 0; i < len(hhdz.zhuanglist); i++ {
		ur := hhdz.GetUserByUserID(hhdz.zhuanglist[i])
		pdata.UserInfo = append(pdata.UserInfo, ur.GetProtoBufGameUser())
	}
	pmsg.MsgData = pdata
	return pmsg
}

//用户下注处理函数
func (hhdz *HHDZLogic) userBet(user *gamelogic.GameUser, data []byte) {
	ud := hhdz.getUserDataByID(user.UserID)
	if ud == nil {
		ud = hhdz.addUserData(user)
	}

	ust := &netproto.HHDZ_BetArr{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Debug("解析用户%v下注错误%v", user, err)
		return
	}

	//莊家不能下注
	if hhdz.zhuangID == user.UserID {
		return
	}

	money := int32(0)
	for _, b := range ust.BetArr {
		*b.BetType = int32(math.Abs(float64(*b.BetType)))
		*b.Count = int32(math.Abs(float64(*b.Count)))

		//玩家身上金額小於下注限制
		if user.Score < hhdz.betlimitmoney*100 {
			for i, u := range hhdz.userdata {
				if u.user.UserID == ud.user.UserID {
					hhdz.userdata = append(hhdz.userdata[:i], hhdz.userdata[i+1:]...)
				}
			}
			hhdz.SendMessageToUser(user, hhdz.getBetRetMsg(ud, -6, *b.BetType, ust.BetArr))
			return
		}

		if hhdz.sametimebetmode != 1 { //關閉莊閒同時下注的配置
			//黑紅不能同時下注
			if (*b.BetType == int32(HUNG) && hhdz.getUserBetByType(ud, HEI) > 0) ||
				(*b.BetType == int32(HEI) && hhdz.getUserBetByType(ud, HUNG) > 0) {
				hhdz.SendMessageToUser(user, hhdz.getBetRetMsg(ud, -5, *b.BetType, ust.BetArr))
				return
			}
		}

		tmpmon := hhdz.getMoneyByType(*b.BetType) * (*b.Count)

		// 系統下注上限
		if (int64(tmpmon) + hhdz.getBetMoneyByType(*b.BetType)) > hhdz.betvalue.getMaxMoneyByType(*b.BetType) {
			hhdz.SendMessageToUser(user, hhdz.getBetRetMsg(ud, -4, *b.BetType, ust.BetArr))
			return
		}

		// 不是系統當莊, 則壓注上限為莊家身上的錢所能夠賠付的錢
		if hhdz.zhuangID != 0 {
			ur := hhdz.GetUserByUserID(hhdz.zhuangID)
			odd := int64(0)
			if odd = int64(hhdz.betvalue.getOddByType(b.GetBetType())) - 1; odd < 0 {
				odd = 0
			}

			allmoney := (int64(tmpmon) + hhdz.getBetMoneyByType(*b.BetType)) * odd
			if allmoney > ur.Score {
				hhdz.SendMessageToUser(user, hhdz.getBetRetMsg(ud, -4, *b.BetType, ust.BetArr))
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
			betmon := hhdz.getMoneyByType(*b.BetType) * (*b.Count)
			ud.setUserBet(*b.BetType, betmon)
			if user.IsRobot() {
				hhdz.updaterobotBet(*b.BetType, int64(betmon))
			} else {
				hhdz.updateBet(*b.BetType, int64(betmon))
			}
			if user.UserID == hhdz.luckystar {
				hhdz.notbettimes = 0
			}
			// hhdz.GameLog.AddGameStep("用户%d下注:%s;金額:%d, %v", user.UserID, hhdz.getBetTypeString(int(b.GetBetType())), betmon, b)
		}
		hhdz.BoardcastMessage(hhdz.getBetTotalMsg(user.UserID))
	}
	hhdz.SendMessageToUser(user, hhdz.getBetRetMsg(ud, ret, 0, ust.BetArr))
}

//添加一个用户数据
func (hhdz *HHDZLogic) addUserData(user *gamelogic.GameUser) *hhdzuuserdata {
	ud := newLongHuDouUserData(user)
	hhdz.userdata = append(hhdz.userdata, ud)
	return ud
}

//获取下注消息
func (hhdz *HHDZLogic) getBetRetMsg(ud *hhdzuuserdata, ret int32, t int32, betinfo []*netproto.HHDZ_Bet) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetRetID)

	pdata := new(netproto.HHDZ_BetRet)
	pdata.Money = proto.Int64(ud.money)
	pdata.Ret = proto.Int32(ret)
	pdata.Type = proto.Int32(t)
	pdata.BetArr = betinfo
	pmsg.MsgData = pdata
	return pmsg
}

func (hhdz *HHDZLogic) getUserBetByType(ud *hhdzuuserdata, t int32) int64 {
	money := int64(0)
	for _, bet := range ud.userbet {
		if bet.hhdztype == t {
			money += bet.money
		}
	}
	return money
}

//获取此類型一次下注额
func (hhdz *HHDZLogic) getMoneyByType(t int32) int32 {
	return hhdz.betvalue.getMoneyByType(t)
}

//更新機器人下注额
func (hhdz *HHDZLogic) updaterobotBet(t int32, money int64) {
	hhdz.bettotalrobot += int64(money)
	for _, v := range hhdz.bettyperobot {
		if v.hhdztype == t {
			v.money += money
			if v.money < 0 {
				v.money = 0
			}
			mlog.Debug("下注類型%d，下注的总錢数为%d", t, v.money)
			return
		}
	}

	b := new(bettype)
	b.hhdztype = t
	b.money = int64(money)
	if b.money < 0 {
		b.money = 0
	}

	hhdz.bettyperobot = append(hhdz.bettyperobot, b)
}

func (hhdz *HHDZLogic) getBetTypeString(idx int) string {
	if idx > 0 && idx < len(resultstring) {
		return resultstring[idx]
	}
	return ""
}

//处理用户上莊请求
func (hhdz *HHDZLogic) shangZhuang(user *gamelogic.GameUser, data []byte) {
	ust := &netproto.HHDZ_ShangZhuang{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Debug("解析用户%v上莊错误%v", user, err)
		return
	}

	if hhdz.getZhuangListIndexByUserId(user.UserID) != -1 {
		hhdz.SendMessageToUser(user, hhdz.getShangzhuangRetMsg(-2))
		return
	}

	if user.Score < hhdz.zhuanglimit {
		hhdz.SendMessageToUser(user, hhdz.getShangzhuangRetMsg(-1))
		return
	}

	if *ust.IsShangZhuang {
		if hhdz.publiczhuang == 1 {
			hhdz.GameLog.AddGameStep("用户%d上莊", user.UserID)
			hhdz.addToZhuangList(user.UserID)
			hhdz.SendMessageToUser(user, hhdz.getShangzhuangRetMsg(0))
		} else {
			mlog.Error("未開啟上莊機制，不得上莊")
			return
		}
	} else {
		hhdz.GameLog.AddGameStep("用户%d下莊", user.UserID)
		hhdz.delFromZhuangList(user.UserID)
		hhdz.SendMessageToUser(user, hhdz.getShangzhuangRetMsg(1))
		hhdz.lastturn = -1
	}
}

//获取上莊列表索引通过用户ID
func (hhdz *HHDZLogic) getZhuangListIndexByUserId(userid int32) int32 {
	for i := 0; i < len(hhdz.zhuanglist); i++ {
		if hhdz.zhuanglist[i] == userid {
			return int32(i)
		}
	}
	return -1
}

//获取上莊结果消息
func (hhdz *HHDZLogic) getShangzhuangRetMsg(ret int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZShangZhuangRetID)

	pdata := new(netproto.HHDZ_ShangzhuangRet)
	pdata.Ret = proto.Int32(ret)
	pmsg.MsgData = pdata
	return pmsg
}

//清除玩家下注信息
func (hhdz *HHDZLogic) clearUserBet(user *gamelogic.GameUser) {
	ud := hhdz.getUserDataByID(user.UserID)
	if ud == nil {
		return
	}

	mlog.Debug("ud.clearcount=%v", ud.clearcount)
	if ud.clearcount >= hhdz.maxclearcount {
		hhdz.SendMessageToUser(user, hhdz.getClearRetMsg(ud.money, -1))
		return
	}

	ud.clearcount++
	for _, ub := range ud.getUserBet() {
		hhdz.updateBet(ub.hhdztype, -ub.money)
	}

	hhdz.GameLog.AddGameStep("用户%d点击清除", user.UserID)
	ud.clearData()
	hhdz.SendMessageToUser(user, hhdz.getClearRetMsg(user.Score, 0))
	hhdz.BoardcastMessage(hhdz.getBetTotalMsg(user.UserID))
}

//获取清除结果消息
func (hhdz *HHDZLogic) getClearRetMsg(money int64, ret int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetClearRetID)

	pdata := new(netproto.HHDZ_BetClearRet)
	pdata.Money = proto.Int64(money)
	pdata.Ret = proto.Int32(ret)
	pmsg.MsgData = pdata

	return pmsg
}

//上莊是為系統或機器人
func (hhdz *HHDZLogic) isZhuangRobotAndSystem() bool {
	var iszhuangrobot bool = false
	ur := hhdz.GetUserByUserID(hhdz.zhuangID)

	if ur != nil && ur.IsRobot() {
		iszhuangrobot = true
	}

	if hhdz.zhuangID == 0 {
		iszhuangrobot = true
	}

	return iszhuangrobot
}

//获取機器人赢最多的開獎方案
func (hhdz *HHDZLogic) getMaxRobotJiang() int32 {
	maxtype := int32(0)
	maxmoney := int64(0)
	for _, rtd := range hhdz.bettyperobot {
		money := int64(hhdz.betvalue.getOddByType(rtd.hhdztype)) * rtd.money
		if money >= maxmoney {
			maxtype = rtd.hhdztype
			maxmoney = money
		}
	}

	return maxtype
}

//計算下注輸贏結果
func (hhdz *HHDZLogic) calWinLose() {
	sysTaxIncomeBefore := hhdz.systaxincome
	hhdz.syslastwin = 0
	hhdz.setState(EndAction)
	for _, ud := range hhdz.userdata {
		hhdz.sendUserResult(ud)
	}

	finalBetTotal := hhdz.bettotal + hhdz.bettotalrobot //真人+機器人下注輸贏總和
	//這裡計算莊家輸贏及稅收

	//準備計算水池異動值
	changeSysPoolValue := hhdz.zhushuimoney
	if hhdz.zhuangID == 0 {
		//系統小妹當莊的時候，玩家贏錢水池減少
		if hhdz.userwinlose > 0 {
			//玩家贏錢，水線下降
			//hhdz.syspool -= hhdz.userwinlose - hhdz.zhushuimoney
			changeSysPoolValue = hhdz.userwinlose
		} else {
			//玩家輸錢，水線上升(稅)
			//hhdz.syspool -= int64(math.Floor(float64(hhdz.userwinlose)*hhdz.GetTax())) - hhdz.zhushuimoney
			changeSysPoolValue = int64(math.Floor(float64(hhdz.userwinlose) * hhdz.GetTax()))
		}
	} else {
		var income, tax int64
		urZhuang := hhdz.GetUserByUserID(hhdz.zhuangID)
		if hhdz.zhuangID == 0 || urZhuang.IsRobot() {
			//莊家是機器人或系統, 莊賺錢只抽真人下注的金額bettotal
			_, tax = hhdz.GetUserWinAndTax(hhdz.bettotal)
		} else {
			//莊家是真人, 莊賺錢要抽真人和機器人下注的金額
			_, tax = hhdz.GetUserWinAndTax(finalBetTotal)
		}
		income, _ = hhdz.GetUserWinAndTax(finalBetTotal)

		hhdz.systaxincome += tax
		_, finalBetTax := hhdz.GetUserWinAndTax(finalBetTotal)
		hhdz.GameLog.AddGameStep("莊家贏錢 %d, 抽稅 %d", income, finalBetTax)
		mlog.Debug("莊家贏錢 %d, 抽稅 %d", income, finalBetTax)

		//這裡計算水池, 系統小妹當莊
		hhdz.zhuangWin += income

		ur := hhdz.GetUserByUserID(hhdz.zhuangID)
		if ur != nil {
			income, tax = hhdz.GetUserWinAndTax(finalBetTotal)
			//機器人當莊的時候
			if ur.IsRobot() {
				//玩家贏錢水池減少，並將注出的水倒回水池
				if hhdz.userwinlose > 0 {
					hhdz.systaxincome += tax //真人贏錢收稅
					//玩家贏錢，水線下降
					hhdz.syslastwin = hhdz.bettotal //系統上一局收入來源只記錄賠給真人的錢
					//hhdz.syspool -= hhdz.userwinlose - hhdz.zhushuimoney
					changeSysPoolValue = hhdz.userwinlose
				} else {
					//玩家輸錢，水線上升(稅)
					//hhdz.syspool -= int64(math.Floor(float64(hhdz.userwinlose)*hhdz.GetTax())) - hhdz.zhushuimoney
					changeSysPoolValue = int64(math.Floor(float64(hhdz.userwinlose) * hhdz.GetTax()))
					realUserIncome, _ := hhdz.GetUserWinAndTax(hhdz.bettotal)
					hhdz.syslastwin = realUserIncome
				}
			} else {

				//玩家當莊的時候，機器人贏錢水池就添加，並將注出的水倒回水池
				if hhdz.robotwinlose > 0 {
					//機器人贏錢 水位上升(稅)
					income, tax = hhdz.GetUserWinAndTax(hhdz.robotwinlose)
					//hhdz.syspool += income - hhdz.robotwinlose
					changeSysPoolValue = -income
					hhdz.syslastwin = income
					hhdz.systaxincome = tax
				} else {
					//hhdz.syspool += hhdz.robotwinlose - hhdz.zhushuimoney
					changeSysPoolValue = -hhdz.robotwinlose
					hhdz.syslastwin += hhdz.robotwinlose
				}

			}

			hhdz.UpdateUserScoreFunc(ur.UserID, hhdz.zhuangWin, tax) // 更新金幣異動
			hhdz.GameLog.Context().UserTaxAmount(ur.UserID, tax)
			hhdz.SendMessageToUser(ur, hhdz.getUserWinMsg(hhdz.zhuangWin, hhdz.zhuangWin))
		} else {
			mlog.Debug("莊家数据为空")
		}

	}
	if !hhdz.execFangFen {
		changeSysPoolValue = -changeSysPoolValue
		example.invalid("系统水池異動, from %d, change %d", hhdz.SysPool.GetPool(), changeSysPoolValue)
		hhdz.GameLog.AddGameStep("系统水池異動, from %d, change %d", hhdz.SysPool.GetPool(), changeSysPoolValue)
		hhdz.SysPool.UpdateSysPool(changeSysPoolValue)
	} else {
		before := hhdz.fangFenComplete
		if changeSysPoolValue > 0 { //玩家贏錢要扣掉稅才是放分
			changeSysPoolValue = int64(math.Floor(float64(changeSysPoolValue) * hhdz.GetTax()))
		} else { //玩家輸錢要把稅算回來
			changeSysPoolValue = int64(math.Floor(float64(changeSysPoolValue) / hhdz.GetTax()))
		}
		hhdz.fangFenComplete += changeSysPoolValue
		hhdz.GameLog.AddGameStep("系統放分完成, target %d, already %d-%d", hhdz.fangFenTarget, before, hhdz.fangFenComplete)
		example.invalid("系統放分完成, target %d, already %d-%d", hhdz.fangFenTarget, before, hhdz.fangFenComplete)
	}

	hhdz.GameLog.AddGameStep("本局稅收 %d", hhdz.systaxincome-sysTaxIncomeBefore)
	hhdz.GameLog.AddGameStep("系統累積稅收 %d", hhdz.systaxincome)

	uscoresRet := make(map[int32]int64)
	//發送沒下注的user 莊家結果
	for _, u2 := range hhdz.AllTableUser() {
		sendBool := true
		for _, ud := range hhdz.userdata {
			if u2.UserID == ud.user.UserID {
				sendBool = false
				break
			}
		}
		if sendBool {
			hhdz.SendMessageToUser(u2, hhdz.getUserWinMsg(0, hhdz.zhuangWin))
			if u2.UserID == hhdz.luckystar {
				if hhdz.getUserDataByID(hhdz.luckystar) == nil {
					hhdz.notbettimes += 1
				}
			}
		}
	}

	// 發送下注user結算結果
	for _, ud := range hhdz.userdata {
		hhdz.SendMessageToUser(ud.user, hhdz.getUserWinMsg(ud.getShowWinLose(), hhdz.zhuangWin))
		uscoresRet[ud.user.UserID] = ud.updatemoney
	}

	// uscoresRet[hhdz.zhuangID] = hhdz.zhuangWin
	// mlog.Debug("莊家贏錢：%v", uscoresRet[hhdz.zhuangID])

	hhdz.GameLog.GameEnd()
	hhdz.bettotal = 0
}

//给每个玩家發送開獎结果
func (hhdz *HHDZLogic) sendUserResult(ud *hhdzuuserdata) {
	// var thirdjiangtype bool //第三區輸贏

	// HUNG       = 1  // 紅
	// HEI        = 2  // 黑
	// NATURE     = 3  // 自然機率
	// BAOZI      = 4  // 豹子
	// SHUNJIN    = 5  // 順金
	// JINHUA     = 6  // 金花
	// SHUNZI     = 7  // 順子
	// BIGDUIZI   = 8  // 對子(9-A)
	// SMALLDUIZI = 9  // 對子(2-8)
	// DANPAI     = 10 // 單張

	/**
	if hhdz.zuhetype >= 4 && hhdz.zuhetype < 9 {
		thirdjiangtype = true
	}
	for i := 0; i < len(ud.userbet); i++ {
		betlosemoney += ud.userbet[i].money
		betAmount += util.ABSInt64(ud.userbet[i].money)
		if ud.userbet[i].hhdztype == 3 && thirdjiangtype {
			win := int64(hhdz.betvalue.getOddByType(hhdz.zuhetype)) * ud.userbet[i].money
			//flowamountTotal := win - ud.userbet[i].money
			//flowAmountAfterTax, _ := hhdz.GetUserWinAndTax(flowamountTotal)
			//flowAmount += flowAmountAfterTax
			//機器人不计算獎池里面的錢
			if hhdz.zhuangID == 0 {
				if !ud.user.IsRobot() {
					hhdz.bettotal -= win
				}
			} else {
				ur := hhdz.GetUserByUserID(hhdz.zhuangID)
				if ur != nil && ur.IsRobot() && ud.user.IsRobot() {
					hhdz.bettotalrobot -= win
				} else {
					hhdz.bettotal -= win
				}
			}
			ud.money += win
		} else if ud.userbet[i].hhdztype == hhdz.jiangtype {
			win := int64(hhdz.betvalue.getOddByType(hhdz.jiangtype)) * ud.userbet[i].money
			//flowamountTotal := win - ud.userbet[i].money
			//flowAmountAfterTax, _ := hhdz.GetUserWinAndTax(flowamountTotal)
			//flowAmount += flowAmountAfterTax
			//機器人不计算獎池里面的錢
			if hhdz.zhuangID == 0 {
				if !ud.user.IsRobot() {
					hhdz.bettotal -= win
				}
			} else {
				ur := hhdz.GetUserByUserID(hhdz.zhuangID)
				if ur != nil && ur.IsRobot() && ud.user.IsRobot() {
					hhdz.bettotalrobot -= win
				} else {
					hhdz.bettotal -= win
				}
			}
			ud.money += win
		} else {
			//flowAmount += ud.userbet[i].money
		}
	}

	money := ud.money - ud.user.Score
	*/

	var updateAmount, flowAmount, betAmount, rtnAmount, winAmount, taxAmount int64

	betMap := make(map[int32]int64)
	for _, bet := range ud.userbet {
		betMap[bet.hhdztype] += bet.money
	}

	// 依押注類型計算
	for bettype, money := range betMap {
		betAmount += money
		var isWin, isReturn bool
		if bettype == hhdz.jiangtype { // 此注類型有開
			isWin = true
			isReturn = true
			win := int64(hhdz.betvalue.getOddByType(bettype))*money - money // 設定中的賠率包含押注金額
			income, tax := hhdz.GetUserWinAndTax(win)

			winAmount += win
			rtnAmount += money
			rtnAmount += income
			taxAmount += tax
		} else if bettype == NATURE && hhdz.zuhetype >= BAOZI && hhdz.zuhetype <= BIGDUIZI { // 押幸運一擊
			isWin = true
			isReturn = true
			win := int64(hhdz.betvalue.getOddByType(hhdz.zuhetype))*money - money // 贏家牌型取倍率
			income, tax := hhdz.GetUserWinAndTax(win)

			winAmount += win
			rtnAmount += money
			rtnAmount += income
			taxAmount += tax
		}

		// 沒贏沒退 = 輸 算流水
		if !isWin && !isReturn {
			flowAmount += money
		}

	}

	updateAmount = rtnAmount - betAmount

	mlog.Debug("User[%d] UpdateAmount: %d, BetAmount: %d, WinAmount: %d, ReturnAmount: %d, FlowAmount: %d, TaxAmount: %d", ud.user.UserID, updateAmount, betAmount, winAmount, rtnAmount, flowAmount, taxAmount)

	hhdz.systaxincome += taxAmount
	ud.win = rtnAmount

	// 將機器人及玩家輸贏減掉稅
	if ud.user.IsRobot() {
		hhdz.robotwinlose += updateAmount
	} else {
		hhdz.userwinlose += updateAmount
	}

	ud.updatemoney = updateAmount

	ud.betAmount = betAmount
	ud.winAmount = rtnAmount + taxAmount
	ud.flowAmount = flowAmount

	stepmsg := "玩家[%d] 下注贏錢: %d, 抽稅: %d, 流水: %d, 有效投注金額: %d, 返還金額: %d"
	if ud.user.IsRobot() {
		stepmsg = "[機器人]" + stepmsg
	}

	if !ud.user.IsRobot() {
		hhdz.GameLog.AddGameStep(stepmsg, ud.user.UserID, updateAmount, taxAmount, flowAmount, betAmount, rtnAmount)
	}

	ud.user.FlowAmount += flowAmount
	ud.user.BetAmount += betAmount
	ud.user.PayAmount += rtnAmount

	hhdz.GameLog.Context().UserUpdateScore(ud.user.UserID, updateAmount)
	hhdz.GameLog.Context().UserBetScore(ud.user.UserID, betAmount)
	hhdz.GameLog.Context().UserPayoutScore(ud.user.UserID, rtnAmount+taxAmount)
	hhdz.GameLog.Context().UserValidScore(ud.user.UserID, flowAmount)
	hhdz.GameLog.Context().UserTaxAmount(ud.user.UserID, taxAmount)

	hhdz.UpdateUserScoreFunc(ud.user.UserID, updateAmount, taxAmount)

	//这里写入赢家排名列表
	winSort(ud, hhdz)
}

//这里写入赢家排名列表
func winSort(ud *hhdzuuserdata, hhdz *HHDZLogic) {
	badd := false
	if ud.win > 0 {
		if len(hhdz.winuserinfo) > 0 {
			for i := 0; i < len(hhdz.winuserinfo); i++ {
				u := hhdz.getUserDataByID(hhdz.winuserinfo[i])
				if u != nil {
					if ud.win > u.win {
						if len(hhdz.winuserinfo) >= 3 {
							if i == 0 {
								hhdz.winuserinfo[2] = hhdz.winuserinfo[1]
								hhdz.winuserinfo[1] = hhdz.winuserinfo[0]
								hhdz.winuserinfo[0] = ud.user.UserID
							} else if i == 1 {
								hhdz.winuserinfo[2] = hhdz.winuserinfo[1]
								hhdz.winuserinfo[1] = ud.user.UserID
							} else {
								hhdz.winuserinfo[2] = ud.user.UserID
							}
						} else if len(hhdz.winuserinfo) == 2 {
							hhdz.winuserinfo = append(hhdz.winuserinfo, ud.user.UserID)
							if i == 0 {
								hhdz.winuserinfo[2] = hhdz.winuserinfo[1]
								hhdz.winuserinfo[1] = hhdz.winuserinfo[0]
								hhdz.winuserinfo[0] = ud.user.UserID
							} else if i == 1 {
								hhdz.winuserinfo[2] = hhdz.winuserinfo[1]
								hhdz.winuserinfo[1] = ud.user.UserID
							}
						} else {
							hhdz.winuserinfo = append(hhdz.winuserinfo, ud.user.UserID)
							hhdz.winuserinfo[1] = hhdz.winuserinfo[0]
							hhdz.winuserinfo[0] = ud.user.UserID
						}

						badd = true
						break
					}
				}
			}
		}

		if len(hhdz.winuserinfo) < 3 && !badd {
			hhdz.winuserinfo = append(hhdz.winuserinfo, ud.user.UserID)
		}
	}
}

func (hhdz *HHDZLogic) OnRecvHTTPMessage(req *http.Request) proto.Message {
	action := req.PostFormValue("action")
	if action == "getserverinfo" {
		return hhdz.getTableListMsg()
	} else if action == "gettableinfo" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		return hhdz.getTableDetailMsg(int32(tid))
	} else if action == "modresult" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		ret := req.PostFormValue("result")
		adminuser := req.PostFormValue("adminuser")
		return hhdz.setNextJiangType(int32(tid), ret, adminuser)
	} else if action == "modpoolamount" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 64)
		typeid, _ := strconv.ParseInt(req.PostFormValue("type"), 10, 32)
		pool, _ := strconv.ParseInt(req.PostFormValue("poolamount"), 10, 64)
		operator := req.PostFormValue("operator")
		if typeid == 1 {
			return hhdz.setSysPool(int32(tid), pool, operator)
		} else {
			return hhdz.setFangFen(int32(tid), pool, operator)
		}
	} else if action == "modpopwaterrate" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		zhushui, _ := strconv.ParseInt(req.PostFormValue("popwaterrate"), 10, 32)
		return hhdz.setZhuShui(int32(tid), int32(zhushui))
	} else if action == "history" {
		pdata := new(netproto.HHDZ_ResultCardList)

		if len(hhdz.resultcardlist.CardList) > 20 {
			pdata.CardList = hhdz.resultcardlist.CardList[len(hhdz.resultcardlist.CardList)-20 : len(hhdz.resultcardlist.CardList)]
		} else {
			pdata = hhdz.resultcardlist
		}

		return pdata
	}

	return nil
}

//获取桌子列表数据
func (hhdz *HHDZLogic) getTableListMsg() proto.Message {
	pdata := new(netproto.Http_HHDZ_GameInfo)
	for _, table := range hhdz.GetTable().Room.Tables {
		pTable := new(netproto.Http_HHDZ_TableList)
		logic := table.GameLogic.(*HHDZLogic)
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
		pTable.SumEarnAmount = proto.Int64(logic.SysPool.GetPool())
		pTable.LastEarnAmount = proto.Int32(int32(logic.syslastwin))
		pTable.SysTaxIncome = proto.Int64(logic.systaxincome)

		pdata.TableList = append(pdata.TableList, pTable)
	}

	return pdata
}

//获取桌子详细数据
func (hhdz *HHDZLogic) getTableDetailMsg(tableid int32) proto.Message {
	t := hhdz.GetTable().Room.GetTableByNo(tableid)
	pTable := new(netproto.Http_HHDZ_TableDetail)
	if t != nil {
		logic := t.GameLogic.(*HHDZLogic)
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
		pTable.OpenResultHung = proto.Int64(logic.getRealUserWinByType(HUNG))
		pTable.OpenResultHei = proto.Int64(logic.getRealUserWinByType(HEI))
		//pTable.OpenResultDraw = proto.Int64(logic.getRealUserWinByType(NATURE))
		pTable.BetTime = proto.Int32(hhdz.getBetTime())

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
		pTable.WinLine = proto.Int64(logic.SysPool.GetWinLine())
		pTable.FanFengTarget = proto.Int64(logic.fangFenTarget)
		pTable.FanFengComplete = proto.Int64(logic.fangFenComplete)
	}

	return pTable
}

//收到RPC消息解析處理
func (hhdz *HHDZLogic) OnRecvRPCMessage(m *msg.Message) {
	rm, ok := m.MsgData.(*rpcmsg.RPCMessage)
	if !ok {
		return
	}
	mmbyte := rm.MsgData.([]byte)
	if len(mmbyte) == 0 {
		return
	}

	if m.SClassID == int32(netproto.GameMessageClassID_SuperToolGetAssignRewardDetailRetID) {
		hhdz.setAssignReward(m)
	}

}

//获取系统赢錢数
func (hhdz *HHDZLogic) getRealUserWinByType(t int32) int64 {
	win := int64(0)
	odd := int64(hhdz.betvalue.getOddByType(t))
	user := hhdz.GetUserByUserID(hhdz.zhuangID)
	if user != nil && !user.IsRobot() {
		_, tax := hhdz.GetUserWinAndTax(hhdz.bettotal)
		if hhdz.currtax != 0 {
			tax = hhdz.currtax
		}

		robotwin := odd * hhdz.getRobotBetMoneyByType(t)
		win = hhdz.bettotalrobot - robotwin - tax
		return -win
	}

	for _, ud := range hhdz.userdata {
		if !ud.user.IsRobot() {
			win += odd * ud.getUserBetByType(t)
			win += ud.money - ud.user.Score
		}
	}

	mlog.Debug("玩家赢錢%d", win)
	return -win
}

//获取此類型的機器人下注总额
func (hhdz *HHDZLogic) getRobotBetMoneyByType(t int32) int64 {
	for i := 0; i < len(hhdz.bettyperobot); i++ {
		if hhdz.bettyperobot[i].hhdztype == t {
			return hhdz.bettyperobot[i].money
		}
	}

	return 0
}

func (hhdz *HHDZLogic) getResultString() string {
	var res string
	res += resultstring[hhdz.jiangtype]
	return res
}

//设置下次開獎结果
func (hhdz *HHDZLogic) setNextJiangType(tableid int32, jiang string, adminuser string) proto.Message {
	t := hhdz.GetTable().Room.GetTableByNo(tableid)
	example.invalid("=======%s", jiang)
	if t != nil {
		var ret []int32
		logic := t.GameLogic.(*HHDZLogic)
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
				logic.nextzuhetype = hhdz.getPowerJiangType()
			}
		}
		example.invalid("修改開獎结果成功 開獎: %d, 贏家組合: %d", logic.nextjiangtype, logic.nextzuhetype)
		return hhdz.getRetMsg(true, "修改開獎结果成功")
	}

	return hhdz.getRetMsg(false, "修改開獎结果失败")
}

func (hhdz *HHDZLogic) getRetMsg(ret bool, msg string) proto.Message {
	pdata := new(netproto.Http_Result)
	pdata.Result = proto.Bool(ret)
	pdata.Message = proto.String(msg)
	return pdata
}

func (hhdz *HHDZLogic) setZhuShui(tableid int32, zhushui int32) proto.Message {
	if tableid == -1 {
		for _, table := range hhdz.GetTable().Room.Tables {
			logic := table.GameLogic.(*HHDZLogic)
			logic.changeZhuShui(zhushui)
		}
	} else {
		t := hhdz.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*HHDZLogic)
			logic.changeZhuShui(zhushui)
		} else {
			return hhdz.getRetMsg(false, "修改注水率失败")
		}
	}

	return hhdz.getRetMsg(true, "修改注水率成功")
}

//無真人下注
func (hhdz *HHDZLogic) noRealUserBet() bool {
	result := true

	for _, ud := range hhdz.userdata {
		if !ud.user.IsRobot() && len(ud.userbet) > 0 {
			return false
		}
	}

	return result
}
