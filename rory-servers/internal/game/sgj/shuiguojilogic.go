package sgj

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

const (
	betState      = iota //下注阶段
	kaiJiangState        //开奖阶段
	actionResult         //动画阶段
	waitState            //等待阶段
)

const (
	action_KaiJiang = "kaijiang" //开奖动画
	action_Ret      = "retaction"
)

const (
	//类型对应，苹果，橘子，芒果，铃铛，西瓜，星星，77，BAR
	PG  = 1
	JZ  = 2
	MG  = 3
	LD  = 4
	XG  = 5
	SX  = 6
	QQ  = 7
	BAR = 8
)

const (
	JZMAX  = 1
	JZMIN  = 2
	SXMAX  = 3
	SXMIN  = 4
	BARMAX = 5
	PGMAX  = 6
	LDMAX  = 7
	LDMIN  = 8
	MGMAX  = 9
	MGMIN  = 10
	LOCK   = 11
	PGMIN  = 14
	QQMAX  = 15
	QQMIN  = 16
	BARMIN = 17
	XGMAX  = 19
	XGMIN  = 20
)

//第一个用来占位
var resultstring = []string{"", "华为×10", "华为×3", "腾讯×30", "腾讯×3",
	"阿里×100", "网易×5", "360×10", "360×3", "京东×10", "京东×3", "崩盘",
	"网易×5", "华为×10", "网易×3", "百度×40", "百度×3", "阿里×50", "网易×5", "小米×20",
	"小米×3", "京东×10", "360×10", "lucky", "网易×5"}

type userScore struct {
	userdata []*gamelogic.GameUser
}

type ShuiGuoJiLogic struct {
	*gamelogic.ComputeGameLogicBase
	serverfile         *os.File
	userdata           []*sgjUserData
	opttimer           *event.TimerFunc
	sysPool            int64
	betTotal           int64
	betRobotTotal      int64
	gameState          int32
	betData            []*sgjBetData      //下注类型和对应的钱
	betrobot           []*sgjBetData      //机器人下注分类汇总
	weight             []*sgjWeightConfig //权重配置文件
	showTypeData       *sgjShowTypeConfig //显示水果的配置文件
	betTypeConfig      *sgjBetTypeConfig  //下注类型的配置文件
	betResultType      int32
	resultTypeList     []int32
	lastResultTypeList []*netproto.SGJ_JiangTypeOdd
	chiShui            int32   //吃水概率百分比，用命令行或者文件重置这个数据
	zhuangList         []int32 //想要上庄人的userid
	zhuangUserID       int32   //当前上庄的玩家ID
	huanZhuang         bool    //下一轮是否换庄
	actionName_curr    string  //当前等待的动画
	scoreSection       []int64 //分数段用来加载不同的分数段配置
	zhhuangLimit       int64   //上庄资金限制
	serverid           string  //配置的serverid
	bettime            int32   //下注时间
	kaiJiangAction     int32   //开奖动画超时
	zhuangturn         int32   //当庄的局数
	restturn           int32   //当庄的剩余局数
	betcount           []int32 //一次下注的数量下发给客户端用
	zhushui            int32   //注水百分比
	maxclearcount      int32   //最大清理次数
	maxrobotzhuang     int32   //最大机器人当庄
	robotxiazhuang     int32   //机器人下庄概率
	robotzhuang        int32   //没有人当庄的时候机器人当庄的概率百分比
	winline            int64   //机器人盈亏线，和系统水池进行比较
	zhushuimoney       int64   //注水的钱
	startsyspool       int64   //开始时的金币数量
	startmoney         int64   //上庄时的金币
	nextjiangtype      int32   //下次开奖类型
	systotalwin        int64   //系统收入
	syslastwin         int64   //上一局收益
	currtax            int64   //当前税收
	systaxincome       int64   //系統稅收
	zhuangWin          int64
	userwinlose        int64
	robotwinlose       int64
	retaction          int32
	betlimitmoney      int64                   //下注限制
	publiczhuang       int64                   //是否開放用戶上庄
	userscore          *userScore              //用户数据
	luckystar          int32                   //幸運星玩家id
	luckyfaceid        int32                   //幸運星玩家頭像id
	luckystarround     int64                   //幸運星持續局數
	luckystartimes     int64                   //幸運星第幾局
	notbettimes        int64                   //沒下注持續局數
	richestuserinfo    []int32                 //富豪玩家的userid
	richestList        []*netproto.SGJ_Richest //富豪榜
	stageuserinfo      []int32                 //標誌玩家的userid
}

func NewShuiGuoJiLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *ShuiGuoJiLogic {
	glc := new(ShuiGuoJiLogic)
	glc.ComputeGameLogicBase = gamelogic.NewComputeGameLogicBase(t, tfManager, cfgs)
	glc.getsgjConfig(cfgs)
	glc.chiShui = 0
	glc.huanZhuang = false //7/2移除上莊機制
	glc.zhuangUserID = 0
	glc.syslastwin = 0
	glc.systotalwin = 0
	glc.systaxincome = 0
	glc.zhuangWin = 0
	glc.nextjiangtype = -1
	glc.gameState = waitState
	glc.zhuangList = make([]int32, 0)
	glc.getLastPool()
	glc.startsyspool = glc.sysPool
	glc.lastResultTypeList = make([]*netproto.SGJ_JiangTypeOdd, 0)
	glc.addCmd()
	glc.betTotal = 0
	glc.RegisterAfterGameStartFunc(func() {
		glc.clearGameData() //清除数据
		glc.getLuckyStar()
		glc.getRichestList()
		glc.setZhuang()
		glc.GameLog.AddGameStep("庄家%d，剩余%d", glc.zhuangUserID, glc.restturn)
		glc.zhuShui()
		glc.broadCastGameState()
	})

	return glc
}

//添加命令行
func (sgj *ShuiGuoJiLogic) addCmd() {
	command.AddCommand(NewSgjCmd(sgj))
}

//释放用户
func (sgj *ShuiGuoJiLogic) UserFree(gameuser *gamelogic.GameUser) {
	if sgj.gameState != waitState {
		if sgj.findPlayerByUserID(gameuser.UserID) != nil {
			return
		}

		if sgj.zhuangUserID == gameuser.UserID {
			mlog.Debug("用户在当庄，所以用户不能离开")
			return
		}
	} else {
		if sgj.zhuangUserID == gameuser.UserID {
			sgj.zhuangUserID = 0
		}
	}

	sgj.delPlayerFromZhuangList(sgj.getIndexByUserId(gameuser.UserID))
	sgj.ChangeUserStateEnd(gameuser)
	sgj.ComputeGameLogicBase.UserFree(gameuser)
}

//用户离开
func (sgj *ShuiGuoJiLogic) UserLeave(gameuser *gamelogic.GameUser) bool {
	if sgj.gameState != waitState {
		if sgj.findPlayerByUserID(gameuser.UserID) != nil {
			mlog.Debug("用户在玩游戏还不能离开")
			return false
		}

		if sgj.zhuangUserID == gameuser.UserID {
			mlog.Debug("用户在当庄，所以用户不能离开")
			return false
		}
	}

	sgj.delPlayerFromZhuangList(sgj.getIndexByUserId(gameuser.UserID))

	sgj.ChangeUserStateEnd(gameuser)
	return sgj.ComputeGameLogicBase.UserLeave(gameuser)
}

//断线重连
func (sgj *ShuiGuoJiLogic) UserRecover(gameuser *gamelogic.GameUser) {
	sgj.ComputeGameLogicBase.UserRecover(gameuser)
	sgj.sendScene(gameuser)
}

//用户进入
func (sgj *ShuiGuoJiLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {
	sgj.ComputeGameLogicBase.UserEnter(gameuser, tableno, seatno)
	sgj.sendScene(gameuser)
}

//发送当前场景
func (sgj *ShuiGuoJiLogic) sendScene(gameuser *gamelogic.GameUser) {
	sgj.SendMessageToUser(gameuser, sgj.getSettingMsg())
	sgj.SendMessageToUser(gameuser, sgj.getGameStateMsg())
	sgj.SendMessageToUser(gameuser, sgj.getRichestListMsg())
	sgj.SendMessageToUser(gameuser, sgj.getLuckyStarMsg())
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SGJ_GameMessageClassID_SGJLastGameResultID)

	pdata := new(netproto.SGJ_LastGameResultType)
	pdata.ResultTypeOdd = append(pdata.ResultTypeOdd, sgj.lastResultTypeList...)
	pdata.BetMoney = proto.Int64(sgj.betTotal)
	pmsg.MsgData = pdata
	sgj.SendMessageToUser(gameuser, pmsg)
	sgj.SendMessageToUser(gameuser, sgj.getZhuangInfoMsg())
	sgj.SendMessageToUser(gameuser, sgj.getBetTotalInfoMsg(nil))

	userbetinfomsg := sgj.getUserBetInfoMsg(gameuser)
	if userbetinfomsg != nil {
		sgj.SendMessageToUser(gameuser, userbetinfomsg)
	}
}

//获取玩家下注消息
func (sgj *ShuiGuoJiLogic) getUserBetInfoMsg(gameuser *gamelogic.GameUser) *msg.Message {
	ud := sgj.findPlayerByUserID(gameuser.UserID)
	if ud == nil {
		return nil
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SGJ_GameMessageClassID_SGJUserBetInfo)

	pdata := new(netproto.SGJ_UserBetInfo)
	pdata.UserMoney = proto.Int64(ud.money)
	for _, v := range ud.betData {
		betinfo := new(netproto.SGJ_BetInfo)
		betinfo.BetType = proto.Int32(v.betType)
		betinfo.BetMoney = proto.Int64(v.bet)
		pdata.BetInfoGroup = append(pdata.BetInfoGroup, betinfo)
	}

	pmsg.MsgData = pdata
	return pmsg
}

//获取庄家信息
func (sgj *ShuiGuoJiLogic) getZhuangInfoMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SGJ_GameMessageClassID_SGJZhuangInfoID)

	user := sgj.GetUserByUserID(sgj.zhuangUserID)
	pdata := new(netproto.SGJ_ZhuangInfo)
	pdata.LastTurn = proto.Int32(sgj.restturn)
	mlog.Debug("剩余庄局数:%v", sgj.restturn)
	if user != nil {
		pdata.UserID = proto.Int32(user.UserID)
		pdata.UserInfo = append(pdata.UserInfo, user.GetProtoBufGameUser())
	} else {
		pdata.UserID = proto.Int32(0)
	}

	pmsg.MsgData = pdata

	return pmsg
}

//获取庄家列表消息
func (sgj *ShuiGuoJiLogic) getZhuangListinfo() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SGJ_GameMessageClassID_SGJShangZhuangListID)

	pdata := new(netproto.SGJ_ZhuangList)
	for i := 0; i < len(sgj.zhuangList); i++ {
		user := sgj.GetUserByUserID(sgj.zhuangList[i])
		if user != nil {
			pdata.UserInfo = append(pdata.UserInfo, user.GetProtoBufGameUser())
		}
	}

	pmsg.MsgData = pdata

	return pmsg
}

//检测换庄
func (sgj *ShuiGuoJiLogic) checkZhuang() {
	if sgj.publiczhuang == 1 {
		if sgj.zhuangUserID == 0 {
			sgj.huanZhuang = true
			return
		}

		if sgj.restturn <= 0 {
			sgj.huanZhuang = true
			return
		}

		ur := sgj.GetUserByUserID(sgj.zhuangUserID)
		if ur != nil && ur.Score < int64(sgj.zhhuangLimit) {
			sgj.huanZhuang = true
		}

		if ur == nil {
			sgj.huanZhuang = true
		}
	} else {
		mlog.Debug("未開啟上庄機制，機器人不用上庄")
		return
	}
}

//通过用户ID获取庄的索引
func (sgj *ShuiGuoJiLogic) getIndexByUserId(id int32) int32 {
	zhuangListCount := len(sgj.zhuangList)
	for i := 0; i < zhuangListCount; i++ {
		if sgj.zhuangList[i] == id {
			return int32(i)
		}
	}

	return -1
}

//获取上一次系统水池的钱
func (sgj *ShuiGuoJiLogic) getLastPool() {
	var folderName = "lastpool"
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err = os.MkdirAll(folderName, os.ModePerm)
	}

	var tablenostr = strconv.Itoa(int(sgj.TableNo()))
	var filename = sgj.serverid + "_" + tablenostr + ".conf"
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
		sgj.serverfile = f
		defer sgj.serverfile.Close()
		sgj.readPoolFromFile()
	}
}

//读取系统水池的钱
func (sgj *ShuiGuoJiLogic) readPoolFromFile() {
	buff := make([]byte, 1024)
	n, _ := sgj.serverfile.Read(buff)
	if n > 0 {
		poolstr := string(buff)[0:n]
		//pool, err := strconv.ParseInt(poolstr, 10, 32)
		pool, _ := strconv.Atoi(poolstr)
		sgj.sysPool = int64(pool)
	}
}

//写系统水池
func (sgj *ShuiGuoJiLogic) writePoolToFile() {
	buff := strconv.Itoa(int(sgj.sysPool))
	//bufio.ReadWriter.Peek()
	tablenostr := strconv.Itoa(int(sgj.TableNo()))
	filename := "./lastpool/" + sgj.serverid + "_" + tablenostr + ".conf"
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

//将玩家从庄家列表删除
func (sgj *ShuiGuoJiLogic) delPlayerFromZhuangList(index int32) {
	zhuangListCount := int32(len(sgj.zhuangList))
	if index >= 0 && index < zhuangListCount {
		sgj.zhuangList = append(sgj.zhuangList[:index], sgj.zhuangList[index+1:]...)
	}
}

//设置庄家
func (sgj *ShuiGuoJiLogic) setZhuang() {
	sgj.checkZhuang()
	if !sgj.huanZhuang {
		sgj.restturn--
		sgj.broadCastZhuang()
		return
	}

	zhuangListCount := len(sgj.zhuangList)
	if len(sgj.zhuangList) > 0 {
		//检查下一个要上庄的用户是否存在
		for {
			user := sgj.GetUserByUserID(sgj.zhuangList[0])
			if user != nil && user.Score >= int64(sgj.zhhuangLimit) {
				break
			}

			sgj.delPlayerFromZhuangList(0)
			if len(sgj.zhuangList) <= 0 {
				break
			}
		}
	}

	zhuangListCount = len(sgj.zhuangList)
	if sgj.zhuangUserID == 0 && zhuangListCount > 0 {
		sgj.zhuangUserID = sgj.zhuangList[0]
		user := sgj.GetUserByUserID(sgj.zhuangUserID)
		sgj.startmoney = user.Score
		sgj.restturn = sgj.zhuangturn - 1
		sgj.huanZhuang = false
	} else {
		if sgj.huanZhuang {
			if zhuangListCount > 0 {
				sgj.zhuangUserID = sgj.zhuangList[0]
				sgj.restturn = sgj.zhuangturn - 1
				user := sgj.GetUserByUserID(sgj.zhuangUserID)
				sgj.startmoney = user.Score
			} else {
				sgj.zhuangUserID = 0
				sgj.startmoney = sgj.startsyspool
			}

			sgj.huanZhuang = false
		}
	}

	sgj.delPlayerFromZhuangList(0)
	sgj.broadCastZhuang()
}

//将系统水池的钱注入下注额里面
func (sgj *ShuiGuoJiLogic) zhuShui() {
	if sgj.zhuangUserID == 0 {
		if sgj.betTotal == 0 {
			sgj.betTotal = sgj.sysPool * int64(sgj.zhushui) / 100
			sgj.sysPool -= sgj.betTotal
		}
	} else {
		ur := sgj.GetUserByUserID(sgj.zhuangUserID)
		if ur != nil {
			if ur.IsRobot() {
				if sgj.betTotal == 0 {
					sgj.betTotal = sgj.sysPool * int64(sgj.zhushui) / 100
					sgj.sysPool -= sgj.betTotal
				}
			}
		}
	}

	sgj.zhushuimoney = sgj.betTotal
	mlog.Debug("游戏开始前庄的userid为:%v，系统水池注水为：%v", sgj.zhuangUserID, sgj.betTotal)
}

//获取状态消息
func (sgj *ShuiGuoJiLogic) getGameStateMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SGJ_GameMessageClassID_SGJStateID)

	pdata := new(netproto.SGJ_Sate)
	pdata.State = proto.Int32(sgj.gameState)

	pdata.Time = proto.Int32(sgj.getBetTime())
	pmsg.MsgData = pdata
	return pmsg
}

func (sgj *ShuiGuoJiLogic) getBetTime() int32 {
	time := int32(0)
	if sgj.opttimer != nil {
		time = sgj.opttimer.WaitSeconds - sgj.opttimer.ElapsedSeconds()
	}
	return time
}

//广播状态
func (sgj *ShuiGuoJiLogic) broadCastGameState() {
	pmsg := sgj.getGameStateMsg()
	sgj.BoardcastMessage(pmsg)
}

//获取设置消息
func (sgj *ShuiGuoJiLogic) getSettingMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SGJ_GameMessageClassID_SGJSettingID)

	pdata := new(netproto.SGJ_Setting)
	pdata.BetCount = append(pdata.BetCount, sgj.betcount...)
	pdata.BetMoney = proto.Int32(sgj.betTypeConfig.getBetMoneyByBetType(1))
	pdata.ZhuangLimit = proto.Int64(sgj.zhhuangLimit)
	pdata.Betlimitmoney = proto.Int64(sgj.betlimitmoney)
	pmsg.MsgData = pdata
	//mlog.Debug("getSettingMsg")
	return pmsg
}

//获取水果机的配置文件
func (sgj *ShuiGuoJiLogic) getsgjConfig(cfgs *config.Group) {
	sgj.scoreSection = make([]int64, 0)
	sgj.weight = make([]*sgjWeightConfig, 0)
	scfgs := cfgs.GetConfig("gameserver")
	sgj.zhhuangLimit = int64(scfgs.GetConfigInt("gameserver", "zhuanglimit"))
	sgj.serverid = strconv.Itoa(int(sgj.GetServerConfig().ServerID))
	sgj.bettime = int32(scfgs.GetConfigInt("gameserver", "bettime"))
	sgj.kaiJiangAction = int32(scfgs.GetConfigInt("gameserver", "kaijiang"))
	sgj.zhuangturn = int32(scfgs.GetConfigInt("gameserver", "zhuangturn"))
	sgj.maxclearcount = int32(scfgs.GetConfigInt("gameserver", "maxclearcount"))
	sgj.zhushui = int32(scfgs.GetConfigInt("gameserver", "zhushui"))
	sgj.maxrobotzhuang = int32(scfgs.GetConfigInt("gameserver", "maxrobotzhuang"))
	sgj.robotxiazhuang = int32(scfgs.GetConfigInt("gameserver", "robotxiazhuang"))
	sgj.robotzhuang = int32(scfgs.GetConfigInt("gameserver", "robotzhuang"))
	sgj.winline = int64(scfgs.GetConfigInt("gameserver", "winline"))
	sgj.retaction = int32(scfgs.GetConfigInt("gameserver", "retaction"))
	sgj.betlimitmoney = int64(scfgs.GetConfigInt("gameserver", "betlimitmoney"))
	sgj.publiczhuang = int64(scfgs.GetConfigInt("gameserver", "publiczhuang"))
	sgj.luckystarround = int64(scfgs.GetConfigInt("gameserver", "luckystarround"))
	betcountstr := scfgs.GetConfigValue("gameserver", "betcount")
	lastindex := 0
	for index, char := range betcountstr {
		if char == ',' {
			substr := betcountstr[lastindex:index]
			lastindex = index + 1
			count, _ := strconv.ParseInt(substr, 10, 32)
			sgj.betcount = append(sgj.betcount, int32(count))
		}
	}
	scorestr := scfgs.GetConfigValue("gameserver", "scoresection")
	lastindex = 0
	for index, char := range scorestr {
		if char == ',' {
			substr := scorestr[lastindex:index]
			lastindex = index + 1
			score, _ := strconv.ParseInt(substr, 10, 64)
			sgj.scoreSection = append(sgj.scoreSection, score)
			w := newSgjWeightConfig(cfgs, substr)
			sgj.weight = append(sgj.weight, w)
		}
	}

	rand := newSgjWeightConfig(cfgs, "rand")
	Pb := newSgjWeightConfig(cfgs, "Pb")
	showtype := newSgjWeightConfig(cfgs, "showtype")
	sgj.weight = append(sgj.weight, rand)
	sgj.weight = append(sgj.weight, Pb)
	sgj.weight = append(sgj.weight, showtype)

	sgj.showTypeData = newSgjShowTyeConfig(cfgs)
	sgj.betTypeConfig = newBetTypeConfig(cfgs)
}

//清理数据
func (sgj *ShuiGuoJiLogic) clearGameData() {
	sgj.userdata = make([]*sgjUserData, 0)
	sgj.gameState = betState
	sgj.betData = make([]*sgjBetData, 0)
	sgj.betrobot = make([]*sgjBetData, 0)
	sgj.stageuserinfo = make([]int32, 0)
	sgj.richestuserinfo = make([]int32, 0)
	sgj.richestList = make([]*netproto.SGJ_Richest, 0)
	sgj.betRobotTotal = 0
	//sgj.betResultType = -1
	sgj.betTotal = 0
	sgj.currtax = 0
	sgj.userwinlose = 0
	sgj.robotwinlose = 0
	sgj.resultTypeList = nil
	sgj.opttimer = sgj.TimerFuncManager.NewNamedTimer(sgj.bettime, sgj.gameResult, "")
	rand.Seed(time.Now().UnixNano())
	zlen := len(sgj.zhuangList)
	if zlen == 0 {
		sgj.robotZhuang()
	}

	sgj.robotXiaZhuang()
}

//机器人下庄
func (sgj *ShuiGuoJiLogic) robotXiaZhuang() {
	if sgj.zhuangUserID == 0 {
		return
	}

	ur := sgj.GetUserByUserID(sgj.zhuangUserID)
	if ur != nil && ur.IsRobot() {
		if rand.Int31n(100) < sgj.robotxiazhuang {
			sgj.zhuangUserID = 0
		}
	}
}

//机器人上庄
func (sgj *ShuiGuoJiLogic) robotZhuang() {
	if sgj.publiczhuang == 1 {
		//概率当庄
		if rand.Int31n(100)+1 > sgj.robotzhuang {
			return
		}
		//统计机器人个数
		robotarr := make([]int32, 0)
		for _, ur := range sgj.AllTableUser() {
			if ur.IsRobot() {
				robotarr = append(robotarr, ur.UserID)
			}
		}

		if len(robotarr) <= 0 {
			return
		}

		//随机有几个人当庄
		n := rand.Int31n(sgj.maxrobotzhuang) + 1

		for i := int32(0); i < n; i++ {
			rcn := len(robotarr)
			if rcn <= 0 {
				break
			}
			index := rand.Int31n(int32(rcn))
			sgj.zhuangList = append(sgj.zhuangList, robotarr[index])
			robotarr = append(robotarr[:index], robotarr[index+1:]...)
		}
	} else {
		mlog.Debug("未開啟上庄機制，機器人不用上庄")
		return
	}
}

//添加玩家数据
func (sgj *ShuiGuoJiLogic) addPlayer(user *gamelogic.GameUser) *sgjUserData {
	ud := newUserData(user)
	sgj.userdata = append(sgj.userdata, ud)
	return ud
}

//删除玩家数据
func (sgj *ShuiGuoJiLogic) delPlayer(user *gamelogic.GameUser) {
	for i := 0; i < len(sgj.userdata); i++ {
		if sgj.userdata[i].user.UserID == user.UserID {
			sgj.userdata = append(sgj.userdata[:i], sgj.userdata[i+1:]...)
			break
		}
	}
}

//通过用户ID找到玩家数据
func (sgj *ShuiGuoJiLogic) findPlayerByUserID(UserID int32) *sgjUserData {
	for _, ud := range sgj.userdata {
		if ud.user.UserID == UserID {
			return ud
		}
	}

	return nil
}

//处理玩家消息
func (sgj *ShuiGuoJiLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	//是不是游戏玩家，这个游戏可能不需要这个
	if !sgj.CheckGameUserMessage(user) {
		sgj.KickUserFunc(user.UserID, "非法操作")
		return
	}

	//不是给游戏逻辑服务器的消息
	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	sgj.CheckUserAnimationCommited(user, gmsg, sgj.actionName_curr, nil)
	//阶段是否正确
	if sgj.gameState != betState {
		//mlog.Debug("用户在错误的阶段下注了，阶段为%v", sgj.gameState)
		return
	}

	if gmsg.SClassID == int32(netproto.SGJ_GameMessageClassID_SGJShangZhuangID) {
		sgj.userShangZhuang(user, gmsg.MsgData.([]byte))
		return
	}

	if gmsg.SClassID == int32(netproto.SGJ_GameMessageClassID_SGJBetID) {
		//庄家不能下注
		if sgj.zhuangUserID == user.UserID {
			mlog.Debug("庄家不能下注")
			return
		}

		ud := sgj.findPlayerByUserID(user.UserID)
		if ud == nil {
			ud = sgj.addPlayer(user)
		}

		sgj.userBet(ud, gmsg.MsgData.([]byte))
		return
	}

	if gmsg.SClassID == int32(netproto.SGJ_GameMessageClassID_SGJGetZhuangListID) {
		sgj.SendMessageToUser(user, sgj.getZhuangListinfo())
		return
	}

	if gmsg.SClassID == int32(netproto.SGJ_GameMessageClassID_SGJClearUserBetID) {
		sgj.clearUserBet(user)
		return
	}
}

//删除用户数据
func (sgj *ShuiGuoJiLogic) delUserData(user *gamelogic.GameUser) {
	for i := 0; i < len(sgj.userdata); i++ {
		if sgj.userdata[i].user.UserID == user.UserID {
			sgj.userdata = append(sgj.userdata[:i], sgj.userdata[i+1:]...)
			break
		}
	}
}

//清理用户下注
func (sgj *ShuiGuoJiLogic) clearUserBet(user *gamelogic.GameUser) {
	ud := sgj.findPlayerByUserID(user.UserID)
	if ud == nil {
		return
	}

	mlog.Debug("ud.clearcount=%v", ud.clearcount)
	if ud.clearcount >= sgj.maxclearcount {
		sgj.SendMessageToUser(user, sgj.getClearRetMsg(user.Score-ud.money, -1))
		return
	}

	ud.clearcount++
	for _, ub := range ud.getUserBet() {
		sgj.addTotalBet(ub.betType, (-ub.bet))
	}

	sgj.GameLog.AddGameStep("用户%d点击清除", user.UserID)
	ud.clearData()

	//sgj.delUserData(user)
	sgj.SendMessageToUser(user, sgj.getClearRetMsg(user.Score, 0))
	sgj.BoardcastMessage(sgj.getBetTotalInfoMsg(nil))
}

//获取清理结果
func (sgj *ShuiGuoJiLogic) getClearRetMsg(money int64, ret int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SGJ_GameMessageClassID_SGJBetClearRetID)

	pdata := new(netproto.SGJ_UserMoney)
	pdata.Money = proto.Int64(money)
	pdata.Ret = proto.Int32(ret)
	pmsg.MsgData = pdata

	return pmsg
}

//检测用户是否够上庄
func (sgj *ShuiGuoJiLogic) canShangZhuang(ur *gamelogic.GameUser) bool {
	if -1 == sgj.getIndexByUserId(ur.UserID) && ur.Score >= int64(sgj.zhhuangLimit) {
		return true
	}

	return false
}

//处理用户上庄
func (sgj *ShuiGuoJiLogic) userShangZhuang(user *gamelogic.GameUser, data []byte) {
	ust := &netproto.SGJ_ShangZhuang{}

	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Debug("解析用户%v用户上庄类型错误%v", user, err)
		return
	}

	ret := int32(0)
	if *ust.IsShangZhuang {
		if sgj.canShangZhuang(user) {
			sgj.GameLog.AddGameStep("玩家%d点击上庄", user.UserID)
			sgj.zhuangList = append(sgj.zhuangList, user.UserID)
		} else {
			if sgj.getIndexByUserId(user.UserID) != -1 {
				ret = -2
			} else {
				ret = -1
			}
		}
	} else if sgj.zhuangUserID == user.UserID {
		sgj.huanZhuang = true
		mlog.Debug("玩家下庄：%v", sgj.zhuangUserID)
		ret = 1
	} else {
		mlog.Debug("将玩家从庄列表中删除：%v", sgj.zhuangUserID)
		sgj.GameLog.AddGameStep("玩家%d点击下庄", user.UserID)
		sgj.delPlayerFromZhuangList(user.UserID)
		ret = 1
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SGJ_GameMessageClassID_SGJShangzhuangRetID)

	pdata := new(netproto.SGJ_ShangzhuangRet)
	pdata.Ret = proto.Int32(ret)
	pmsg.MsgData = pdata

	sgj.SendMessageToUser(user, pmsg)
}

//获取下注的钱，通过下注类型
func (sgj *ShuiGuoJiLogic) getbetMoneyByBetType(bettype int32) int64 {
	for i := 0; i < len(sgj.betData); i++ {
		if sgj.betData[i].betType == bettype {
			return sgj.betData[i].bet
		}
	}

	return 0
}

//获取机器人下注的钱，通过下注类型
func (sgj *ShuiGuoJiLogic) getRobotBetMoneyByBetType(bettype int32) int64 {
	for i := 0; i < len(sgj.betrobot); i++ {
		if sgj.betrobot[i].betType == bettype {
			return sgj.betrobot[i].bet
		}
	}

	return 0
}

//用户下注
func (sgj *ShuiGuoJiLogic) userBet(ud *sgjUserData, data []byte) {
	ust := &netproto.SGJ_BetTypeGroup{}

	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Debug("解析用户%v用户下注类型错误%v", ud.user, err)
		return
	}

	money := int32(0)
	//先算有多少钱
	for _, v := range ust.GetBetGroup() {
		*v.BetType = int32(math.Abs(float64(*v.BetType)))
		*v.BetCount = int32(math.Abs(float64(*v.BetCount)))

		if *v.BetCount == 0 {
			//下注个数是0
			return
		}
		// 下注次數 x 一次下注多少錢
		tmpmon := sgj.betTypeConfig.getBetMoneyByBetType(*v.BetType) * (*v.BetCount)
		if !ud.canBet(tmpmon) {
			sgj.SendMessageToUser(ud.user, sgj.getCannotBetMsg(ud, -1, 0, ust.GetBetGroup()))
			return
		}

		//玩家身上金額小於下注限制
		if ud.user.Score < sgj.betlimitmoney*100 {
			for i, u := range sgj.userdata {
				if u.user.UserID == ud.user.UserID {
					sgj.userdata = append(sgj.userdata[:i], sgj.userdata[i+1:]...)
				}
			}
			sgj.SendMessageToUser(ud.user, sgj.getCannotBetMsg(ud, -6, 0, ust.GetBetGroup()))
			return
		}

		// 這個類型目前共下注了多少錢
		betmoney := sgj.getbetMoneyByBetType(*v.BetType) + int64(tmpmon)

		// 是否超過系統下注上限
		if betmoney > sgj.betTypeConfig.getMaxBetMoneyByBetType(*v.BetType) {
			sgj.SendMessageToUser(ud.user, sgj.getCannotBetMsg(ud, -2, *v.BetType, ust.GetBetGroup()))
			return
		}

		// 不是系統當庄, 則壓注上限為庄家身上的錢所能夠賠付的錢
		if sgj.zhuangUserID != 0 {
			ur := sgj.GetUserByUserID(sgj.zhuangUserID)
			if ur == nil {
				mlog.Fatal("庄家没有了，庄ID=%d", sgj.zhuangUserID)
				return
			}

			allmoney := int64(0)
			for _, t := range sgj.showTypeData.configData {
				if t.sgjSecondType == v.GetBetType() {
					odd := int64(0)
					if odd = int64(sgj.showTypeData.getOddByType(t.sgjType)) - 1; odd < 0 {
						odd = 0
					}

					allmoney += betmoney * odd
				}
			}

			if allmoney > ur.Score {
				sgj.SendMessageToUser(ud.user, sgj.getCannotBetMsg(ud, -2, *v.BetType, ust.GetBetGroup()))
				return
			}

		}

		money += tmpmon
	}

	if !ud.canBet(money) {
		sgj.SendMessageToUser(ud.user, sgj.getCannotBetMsg(ud, -1, 0, ust.GetBetGroup()))
		return
	}

	//下面才是真的下注
	money = 0
	for _, v := range ust.GetBetGroup() {
		tmpmon := sgj.betTypeConfig.getBetMoneyByBetType(*v.BetType) * (*v.BetCount)

		ud.bet(*v.BetType, tmpmon)
		sgj.addTotalBet(*v.BetType, int64(tmpmon))
		if ud.user.IsRobot() {
			sgj.addRobotBet(*v.BetType, int64(tmpmon))
		}
		money += tmpmon
	}
	if ud.user.UserID == sgj.luckystar {
		sgj.notbettimes = 0
	}

	sgj.GameLog.AddGameStep("玩家%d下注%v", ud.user.UserID, ust.GetBetGroup())
	if money > 0 {
		sgj.broadCastUserBet(ud, 0, money)
		sgj.SendMessageToUser(ud.user, sgj.getCannotBetMsg(ud, 0, 0, ust.GetBetGroup()))
	}
}

//广播当前的庄，如果为0就表示是系统庄
func (sgj *ShuiGuoJiLogic) broadCastZhuang() {
	pmsg := sgj.getZhuangInfoMsg()
	sgj.BoardcastMessage(pmsg)
	sgj.BoardcastMessage(sgj.getRichestListMsg())
	sgj.BoardcastMessage(sgj.getLuckyStarMsg())
}

//广播用户下注
func (sgj *ShuiGuoJiLogic) broadCastUserBet(ud *sgjUserData, betType int32, money int32) {
	//mlog.Debug("broadCastUserBet")
	sgj.BoardcastMessage(sgj.getBetTotalInfoMsg(ud))
}

//获取用户不能下注的消息
func (sgj *ShuiGuoJiLogic) getCannotBetMsg(ud *sgjUserData, ret int32, bettype int32, betinfo []*netproto.SGJ_BetType) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SGJ_GameMessageClassID_SGJCannotBetID)

	pdata := new(netproto.SGJ_UserBetRet)
	pdata.Money = proto.Int64(ud.user.Score - ud.money)
	pdata.Ret = proto.Int32(ret)
	pdata.Type = proto.Int32(bettype)
	pdata.BetGroup = betinfo
	pmsg.MsgData = pdata
	return pmsg
}

//获取总下注额消息
func (sgj *ShuiGuoJiLogic) getBetTotalInfoMsg(ud *sgjUserData) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SGJ_GameMessageClassID_SGJBroadCastBetID)

	pdata := new(netproto.SGJ_BroadCastBetType)
	if ud != nil {
		pdata.BetUserID = proto.Int32(ud.user.UserID)
	} else {
		pdata.BetUserID = proto.Int32(0)
	}
	lucky := sgj.getUserDataByID(sgj.luckystar) //獲取明燈
	for _, t := range sgj.betTypeConfig.betType {
		bet := new(netproto.SGJ_BetTotalInfo)
		bet.BetType = proto.Int32(t)
		bet.Money = proto.Int64(sgj.getbetMoneyByBetType(t))
		bet.LuckyStar = proto.Bool(false)
		if lucky != nil && lucky.getUserBetByType(t) > 0 {
			bet.LuckyStar = proto.Bool(true)
		}
		pdata.BetInfo = append(pdata.BetInfo, bet)
	}

	pmsg.MsgData = pdata

	return pmsg
}

//将加注的钱加在统计类型上
func (sgj *ShuiGuoJiLogic) addTotalBet(betType int32, bet int64) {
	sgj.betTotal += bet
	for _, bd := range sgj.betData {
		if bd.betType == betType {
			bd.bet += int64(bet)
			if bd.bet < 0 {
				bd.bet = 0
			}
			return
		}
	}

	bd := new(sgjBetData)
	bd.betType = betType
	bd.bet = int64(bet)
	if bd.bet < 0 {
		bd.bet = 0
	}
	sgj.betData = append(sgj.betData, bd)
}

//加机器人下注的钱
func (sgj *ShuiGuoJiLogic) addRobotBet(betType int32, bet int64) {
	sgj.betRobotTotal += bet
	for _, bd := range sgj.betrobot {
		if bd.betType == betType {
			bd.bet += int64(bet)
			if bd.bet < 0 {
				bd.bet = 0
			}
			return
		}
	}

	bd := new(sgjBetData)
	bd.betType = betType
	bd.bet = int64(bet)
	if bd.bet < 0 {
		bd.bet = 0
	}
	sgj.betrobot = append(sgj.betrobot, bd)
}

//移除机器人下注的钱
func (sgj *ShuiGuoJiLogic) removeRobotBet() {
	for _, bd := range sgj.betrobot {
		sgj.addTotalBet(bd.betType, -bd.bet)
	}
}

//广播开奖类型
func (sgj *ShuiGuoJiLogic) broadCastJiangType() {
	sgj.BoardcastMessage(sgj.getJiangTypeMsg())
}

//获取开奖消息
func (sgj *ShuiGuoJiLogic) getJiangTypeMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SGJ_GameMessageClassID_SGJBroadCastJiangType)

	pdata := new(netproto.SGJ_JiangType)
	pdata.JiangType = proto.Int32(sgj.betResultType)
	pdata.JiangTypeList = append(pdata.JiangTypeList, sgj.resultTypeList...)
	pmsg.MsgData = pdata
	mlog.Debug("开奖类型：%v", pdata)
	return pmsg
}

//获取机器人赢多少钱
func (sgj *ShuiGuoJiLogic) getRobotWin() int64 {
	win := int64(0)
	for _, rtd := range sgj.betrobot {
		for _, t := range sgj.resultTypeList {
			if rtd.betType == t {
				win += rtd.bet * int64(sgj.showTypeData.getOddByType(t))
				break
			}
		}
	}

	return win
}

//给玩家发送结果
func (sgj *ShuiGuoJiLogic) sendUserResult() {
	calWinLose(sgj)

	sgj.betTotal = 0
	sgj.writePoolToFile()
	sgj.gameState = waitState
	mlog.Debug("游戏结束后庄的userid为:%v，系统水池水为：%v", sgj.zhuangUserID, sgj.sysPool)

	if sgj.zhuangUserID != 0 {
		user := sgj.GetUserByUserID(sgj.zhuangUserID)
		if user.IsOffLine() {
			sgj.zhuangUserID = 0
		}
	}
	//sgj.GameEndFunc()
	//sgj.RegisterAnimation(action_Ret, sgj.retaction, sgj.GameEndFunc)
	sgj.TimerFuncManager.NewTimer(sgj.retaction, sgj.GameEndFunc)
}

//計算下注輸贏結果
func calWinLose(sgj *ShuiGuoJiLogic) {
	sgj.syslastwin = 0
	for _, ud := range sgj.userdata {
		sgj.updateUserMoney(ud)
	}

	var income, tax int64
	if sgj.betTotal > int64(0) { //賺錢要抽稅
		income, tax = sgj.GetUserWinAndTax(sgj.betTotal)
		sgj.systaxincome += tax
	} else {
		income = sgj.betTotal
	}
	sgj.systotalwin += income
	sgj.syslastwin = income
	sgj.zhuangWin = income

	zhuang := sgj.zhuangUserID
	ur := sgj.GetUserByUserID(zhuang)
	var finalBetTotal int64
	if ur != nil {
		if ur.IsRobot() {
			finalBetTotal = sgj.betTotal + sgj.betRobotTotal
		} else {
			finalBetTotal = sgj.betTotal
		}
		income, tax = sgj.GetUserWinAndTax(finalBetTotal)
		if ur.IsRobot() {
			//機器人當莊，水線根據玩家輸贏改變，並將注出來的水倒回水池
			if sgj.userwinlose > 0 {
				//玩家贏錢，水線下降
				sgj.sysPool -= sgj.userwinlose - sgj.zhushuimoney
			} else {
				//玩家輸錢，水線上升(稅)
				sgj.sysPool -= int64(math.Floor(float64(sgj.userwinlose)*sgj.GetTax())) - sgj.zhushuimoney
			}
			mlog.Debug("桌子号：%d发完将以后剩余的下注钱%d,机器人下注剩余的钱%d", sgj.GetTable().TableNo, sgj.betTotal, sgj.betRobotTotal)

			sgj.zhuangWin = income
			sgj.systotalwin += tax
			sgj.systaxincome += income
			//sgj.sysPool += sgj.zhuangWin
			mlog.Debug("桌子号：%d,水池机器人当庄变化:%d", sgj.GetTable().TableNo, sgj.sysPool)
		} else {
			//玩家當莊，水線根據機器人輸贏改變，並將注出來的水倒回水池
			if sgj.robotwinlose > 0 {
				//機器人贏錢 水位上升(稅)
				sgj.sysPool += int64(math.Floor(float64(sgj.robotwinlose)*sgj.GetTax())) - sgj.zhushuimoney
			} else {
				sgj.sysPool += sgj.robotwinlose - sgj.zhushuimoney
			}
			mlog.Debug("系统水池異動%d：%d", sgj.robotwinlose, sgj.sysPool)

		}
		sgj.UpdateUserScoreFunc(ur.UserID, income, tax) // 更新金幣異動
		sgj.GameLog.Context().UserUpdateScore(ur.UserID, income)
		sgj.GameLog.Context().UserTaxAmount(ur.UserID, tax)
		sgj.sendUserGameResultMsg(ur, sgj.zhuangWin, sgj.zhuangWin)
	} else {
		if sgj.zhuangUserID == 0 {
			//無人當莊，水線根據玩家輸贏改變，並將注出來的水倒回水池
			if sgj.userwinlose > 0 {
				//玩家贏錢，水線下降
				sgj.sysPool -= sgj.userwinlose - sgj.zhushuimoney
			} else {
				//玩家輸錢，水線上升(稅)
				sgj.sysPool -= int64(math.Floor(float64(sgj.userwinlose)*sgj.GetTax())) - sgj.zhushuimoney
			}
		} else {
			mlog.Error("no zhuang user")
		}
	}

	//發送沒下注的user 莊家結果
	for _, u := range sgj.AllTableUser() {
		sendBool := true
		for _, ud := range sgj.userdata {
			if u.UserID == ud.user.UserID {
				sendBool = false
				break
			}
		}
		if sendBool {
			sgj.sendUserGameResultMsg(u, 0, sgj.zhuangWin)
			if u.UserID == sgj.luckystar {
				if sgj.getUserDataByID(sgj.luckystar) == nil {
					sgj.notbettimes += 1
				}
			}
		}
	}

	uscoresRet := make(map[int32]int64)
	for _, ud := range sgj.userdata {
		sgj.sendUserGameResultMsg(ud.user, ud.win, sgj.zhuangWin)
		uscoresRet[ud.user.UserID] = ud.updatemoney
		sgj.GameLog.Context().UserUpdateScore(ud.user.UserID, ud.updatemoney)
	}
	uscoresRet[sgj.zhuangUserID] = sgj.zhuangWin
	sgj.GameLog.GameEnd()
}

//添加开奖结果到历史开奖列表
func (sgj *ShuiGuoJiLogic) addTolastResultTypeList() {
	result := new(netproto.SGJ_JiangTypeOdd)
	result.ResultType = proto.Int32(sgj.betResultType)
	odd := int32(0)
	for i := 0; i < len(sgj.resultTypeList); i++ {
		odd += sgj.showTypeData.getOddByType(sgj.resultTypeList[i])
	}

	result.Odd = proto.Int32(odd)
	if len(sgj.lastResultTypeList) >= 20 {
		sgj.lastResultTypeList = append(sgj.lastResultTypeList[:0], sgj.lastResultTypeList[1:]...)
	}

	sgj.lastResultTypeList = append(sgj.lastResultTypeList, result)
}

//游戏结算
func (sgj *ShuiGuoJiLogic) gameResult() {
	mlog.Debug("开始开奖")
	//抽水池的水
	if sgj.zhuangUserID != 0 {
		if sgj.GetUserByUserID(sgj.zhuangUserID).IsRobot() {
			sgj.removeRobotBet()
			mlog.Debug("桌子号：%d，移除机器人下注", sgj.TableNo())
		}
	} else {
		sgj.removeRobotBet()
	}

	money, tax := sgj.GetUserWinAndTax(sgj.betTotal)
	sgj.currtax = tax
	mlog.Debug("抽税前%d，抽税后%d", sgj.betTotal, money)
	//sgj.betTotal = money
	sgj.gameState = kaiJiangState
	sgj.getKaiJiangType()
	sgj.addTolastResultTypeList()
	mlog.Debug("桌子号：%d,开奖结果是：%d，%v", sgj.TableNo(), sgj.betResultType, sgj.resultTypeList)
	sgj.GameLog.AddGameStep("开奖结果是：%d，%v", sgj.betResultType, sgj.resultTypeList)
	sgj.broadCastJiangType()
	sgj.actionName_curr = sgj.RegisterAnimation(action_KaiJiang, sgj.kaiJiangAction, sgj.sendUserResult)

	sgj.gameState = actionResult
}

//更新玩家的钱
func (sgj *ShuiGuoJiLogic) updateUserMoney(ud *sgjUserData) {
	win := sgj.getUserWin(ud)
	flowAmount, betAmount, payAmount := sgj.getUserFlowAmount(ud)
	if sgj.zhuangUserID == 0 {
		if !ud.user.IsRobot() {
			sgj.betTotal -= win
		}
	} else {
		if sgj.GetUserByUserID(sgj.zhuangUserID).IsRobot() {
			if !ud.user.IsRobot() {
				sgj.betTotal -= win
				mlog.Debug("给玩家发奖：%d", win)
			} else {
				sgj.betRobotTotal -= win
				mlog.Debug("给机器人发奖：%d", win)
			}
		} else {
			sgj.betTotal -= win
		}
	}

	money := win - ud.money

	if ud.user.IsRobot() {
		sgj.robotwinlose += money
	} else {
		sgj.userwinlose += money
	}

	var income, tax int64
	if money > int64(0) {
		income, tax = sgj.GetUserWinAndTax(money)
		ud.win = income + ud.money
		sgj.systaxincome += tax
	} else {
		income = money
	}
	ud.updatemoney = income

	sgj.GameLog.AddGameStep("用户%d, 下注赢钱%d, 系統稅收%d, 流水%d, 有效投注金額%d, 返還金額%d", ud.user.UserID, income, tax, flowAmount, betAmount, payAmount)
	ud.user.FlowAmount += flowAmount
	ud.user.BetAmount += betAmount
	ud.user.PayAmount += payAmount
	sgj.UpdateUserScoreFunc(ud.user.UserID, income, tax)
	sgj.GameLog.Context().UserUpdateScore(ud.user.UserID, ud.updatemoney)
	sgj.GameLog.Context().UserBetScore(ud.user.UserID, betAmount)
	sgj.GameLog.Context().UserPayoutScore(ud.user.UserID, payAmount)
	sgj.GameLog.Context().UserValidScore(ud.user.UserID, flowAmount)
	sgj.GameLog.Context().UserTaxAmount(ud.user.UserID, tax)
	ud.setUserWinMoney(win - tax) //顯示win
}

//获取玩家赢钱
func (sgj *ShuiGuoJiLogic) getUserWin(ud *sgjUserData) int64 {
	win := int64(0)
	for i := 0; i < len(sgj.resultTypeList); i++ {
		bettype := sgj.showTypeData.getSecondTypeByType(sgj.resultTypeList[i])
		odd := sgj.showTypeData.getOddByType(sgj.resultTypeList[i])
		win += ud.faJiang(bettype, odd)
	}

	return win
}

//获取玩家流水, 有效投注, 返還金額
func (sgj *ShuiGuoJiLogic) getUserFlowAmount(ud *sgjUserData) (flowAmount, betAmount, payAmount int64) {
	secondTypeMap := make(map[int32]int32)
	for i := 0; i < len(sgj.resultTypeList); i++ {
		bettype := sgj.showTypeData.getSecondTypeByType(sgj.resultTypeList[i])
		secondTypeMap[bettype] = sgj.resultTypeList[i]
	}
	var win int64
	for _, betdata := range ud.betData {
		betAmount += betdata.bet
		if resultType, ok := secondTypeMap[betdata.betType]; !ok {
			win -= betdata.bet
			flowAmount += betdata.bet
		} else {
			odd := sgj.showTypeData.getOddByType(resultType)
			win += int64(odd)*betdata.bet - betdata.bet
			flowAmount += int64(odd) * betdata.bet
		}
	}
	amount, tax := sgj.GetUserWinAndTax(win)
	flowAmount -= tax
	payAmount = betAmount + amount
	return
}

//发送用户游戏结果
func (sgj *ShuiGuoJiLogic) sendUserGameResultMsg(ur *gamelogic.GameUser, selfmoney int64, zhuangmoney int64) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SGJ_GameMessageClassID_SGJGameResultID)

	pdata := new(netproto.SGJ_GameResult)
	pdata.ZhuangWin = proto.Int64(zhuangmoney)
	pdata.SelfWin = proto.Int64(selfmoney)
	pdata.JiangTypeList = append(pdata.JiangTypeList, sgj.resultTypeList...)
	pdata.Retaction = proto.Int32(sgj.retaction)
	for i := 0; i < len(sgj.stageuserinfo); i++ {
		ud := sgj.getUserDataByID(sgj.stageuserinfo[i])
		if ud == nil {
			continue
		}
		stageuserwin := new(netproto.SGJ_StageUerInfo)
		stageuserwin.UserID = proto.Int32(ud.user.UserID)
		stageuserwin.Win = proto.Int64(ud.win)
		pdata.StageUerInfo = append(pdata.StageUerInfo, stageuserwin)
	}
	pmsg.MsgData = pdata

	//mlog.Debug("sendUserGameResultMsg")
	sgj.SendMessageToUser(ur, pmsg)
}

//获取开奖的类型
func (sgj *ShuiGuoJiLogic) getKaiJiangType() {
	if sgj.nextjiangtype != -1 {
		sgj.betResultType = sgj.nextjiangtype
		sgj.resultTypeList = make([]int32, 0)
		sgj.resultTypeList = append(sgj.resultTypeList, sgj.betResultType)
		sgj.nextjiangtype = -1
		return
	}

	money := sgj.betTotal
	mlog.Debug("桌子号：%d,sgj.betTotal=%v", sgj.TableNo(), sgj.betTotal)
	if sgj.zhuangUserID != 0 {
		ur := sgj.GetUserByUserID(sgj.zhuangUserID)
		if ur != nil && !ur.IsRobot() {
			money += ur.Score
		}
		mlog.Debug("桌子号：%d,sgj.zhuangUserID=%v,money=%v", sgj.TableNo(), sgj.zhuangUserID, money)
	}
	//這裡開始決定機率表
	index := sgj.getWeightIndex()
	if sgj.zhuangUserID != 0 {
		ur := sgj.GetUserByUserID(sgj.zhuangUserID)
		if ur != nil {
			if !ur.IsRobot() {
				index = len(sgj.scoreSection) - 1
				mlog.Debug("玩家当庄直接选最后一个")
			}
		}
	}

	example.invalid("桌子号：%d,选择的index=%v", sgj.TableNo(), index)
	t := 0
	//玩家當莊
	if sgj.zhuangUserID != 0 && !sgj.GetUserByUserID(sgj.zhuangUserID).IsRobot() {
		//機器人有下注的時候
		//if len(sgj.betrobot) > int(0) {
		//贏錢線控制，5/24修正，暫時先不咬玩家
		//	if sgj.sysPool < sgj.winline {
		//		//若玩家當莊，機率隨機
		//		eat := util.RndInt(1, 5)
		//		mlog.NewDefaultLogger().System("水池小於贏錢線1/2吃玩家, rand 1-4[%d], >=3 eat", eat)
		//		if eat >= 3 {
		//			t=1
		//		}
		//
		//	} else {
		//		eat := util.RndInt(1, 5)
		//		mlog.NewDefaultLogger().System("水池大於贏錢線1/4吃玩家, rand 1-4[%d], >=4 eat", eat)
		//		if eat >= 4 {
		//			t=1
		//		}
		//	}
		//} else {
		//玩家當莊且機器人沒有下注的時候
		t = 2
		//}
	}

	// t=1 機器人獎 t=2 Pb配置開獎，除此之外依照水線開獎
	if t == 1 {
		sgj.getMaxRobotJiang()
	} else if t == 2 {
		index = 5
	} else {
		var userbet int64
		for _, v := range sgj.betData {
			bet := v.bet
			userbet = userbet + bet
		}
		//當機器人或系統當莊的時候，玩家未投注，開ShowType
		example.invalid("userbet = %v", userbet)
		if (sgj.zhuangUserID == 0 || (sgj.zhuangUserID != 0 && sgj.GetUserByUserID(sgj.zhuangUserID).IsRobot())) && userbet == 0 {
			index = len(sgj.weight) - 1
			//sgj.getShowJiangType(index)
		}
	}
	//決定好開獎機率表後進行開獎
	for idx, w := range sgj.weight {
		if idx == index {
			example.invalid("***idx=%d, w=[%v]", idx, w)
		} else {
			example.invalid("   idx=%d, w=[%v]", idx, w)
		}
	}

	//開獎邏輯, a. 高於贏錢線才可以開大獎, b. 獲取是否吃水或大獎, c. 一般開獎
	//dWinline := (sgj.zhushuimoney + sgj.sysPool) < sgj.winline
	flag := sgj.weight[index].getLuckyOrLose() //1 吃水, 2 Lucky

	if flag == 1 { //吃水
		//吃水的情况
		mlog.Debug("吃水的情况")
		sgj.betResultType = 0
	} else if flag == 2 { //大獎
		if !sgj.getDaJiangType(money, index) {
			sgj.getNormalJiangType(money, index)
		}
	} else { //一般開獎
		sgj.getNormalJiangType(money, index)
	}

	//指定大獎
	//sgj.betResultType = 30
	//sgj.resultTypeList = []int32{16, 17, 18, 19, 20}
	sgj.GameLog.AddGameStep("水線：%d，開獎表：%v", sgj.zhushuimoney, index+1)
	example.invalid("注水Money：%d，開獎表：%v", sgj.zhushuimoney, index+1)
}

//获取机器人赢最多钱的结果
func (sgj *ShuiGuoJiLogic) getMaxRobotJiang() {
	maxtype := int32(2)
	maxmoney := int64(0)
	//1~24主類別
	for _, cfdata := range sgj.showTypeData.configData {
		// 略過 Bar 獎
		if cfdata.sgjSecondType == 8 {
			continue
		}
		var b *sgjBetData
		for _, b = range sgj.betrobot {
			if b.betType == cfdata.sgjSecondType {
				money := b.bet * int64(cfdata.oddValue)
				if money > maxmoney {
					maxtype = cfdata.sgjType
					maxmoney = money
				}
			}
		}
	}
	sgj.betResultType = maxtype
	sgj.resultTypeList = make([]int32, 0)
	sgj.resultTypeList = append(sgj.resultTypeList, sgj.betResultType)
}

//获取大奖的类型
func (sgj *ShuiGuoJiLogic) getDaJiangType(money int64, index int) bool {

	fivelist := sgj.getFiveJiangTypeList(money)
	fourlist := sgj.getFourJiangTypeList(money)
	maxthreelist := sgj.getMaxThreeJiangTypeList(money)
	mixthreelist := sgj.getMinThreeJiangTypeList(money)
	maxtwolist := sgj.getMaxTwoJiangTypeList(money)
	mintwolist := sgj.getMixTwoJiangTypeList(money)
	allType := make([]int32, 0)
	if len(*mintwolist) > 0 {
		allType = append(allType, sgj.weight[index].djType[0])
	}
	if len(*maxtwolist) > 0 {
		allType = append(allType, sgj.weight[index].djType[1])
	}
	if len(*mixthreelist) > 0 {
		allType = append(allType, sgj.weight[index].djType[2])
	}
	if len(*maxthreelist) > 0 {
		allType = append(allType, sgj.weight[index].djType[3])
	}
	if len(*fourlist) > 0 {
		allType = append(allType, sgj.weight[index].djType[4])
	}
	if len(*fivelist) > 0 {
		allType = append(allType, sgj.weight[index].djType[5])
	}

	if len(allType) > 0 {
		mlog.Debug("allType=%v", allType)
		jtype := sgj.weight[index].getDaJiangType(allType)
		var dajianglist []int32 = nil
		if jtype == sgj.weight[index].djType[0] {
			dajianglist = sgj.weight[index].getDaJiangTypeList(mintwolist)
		} else if jtype == sgj.weight[index].djType[1] {
			dajianglist = sgj.weight[index].getDaJiangTypeList(maxtwolist)
		} else if jtype == sgj.weight[index].djType[2] {
			dajianglist = sgj.weight[index].getDaJiangTypeList(mixthreelist)
		} else if jtype == sgj.weight[index].djType[3] {
			dajianglist = sgj.weight[index].getDaJiangTypeList(maxthreelist)
		} else if jtype == sgj.weight[index].djType[4] {
			dajianglist = sgj.weight[index].getDaJiangTypeList(fourlist)
		} else if jtype == sgj.weight[index].djType[5] {
			dajianglist = sgj.weight[index].getDaJiangTypeList(fivelist)
		}

		if dajianglist == nil {
			mlog.Error("大奖没开出来!!!!!!!!!!!")
		} else {
			sgj.betResultType = jtype
			sgj.resultTypeList = dajianglist
			return true
		}
	}

	return false
}

//获取五连将的列表
func (sgj *ShuiGuoJiLogic) getFiveJiangTypeList(money int64) *[][]int32 {
	//通算所有奖项
	allType := make([]int32, 0)
	allType = append(allType, sgj.showTypeData.allType...)
	typelen := len(allType)
	typelist := make([][]int32, 0)
	for i := 0; i < typelen; i++ {
		typearr := make([]int32, 0)
		for j := 0; j < 5; j++ {
			n := i + j
			if n >= typelen {
				n %= typelen
			}

			typearr = append(typearr, allType[n])
		}

		jiang := int64(0)
		for m := 0; m < 5; m++ {
			data := sgj.showTypeData.getShowTypeDataByType(typearr[m])
			jiang += sgj.getUsersBetMoneyByBetType(data.sgjSecondType) * int64(data.oddValue)
		}

		if jiang <= money {
			typelist = append(typelist, typearr)
		}
	}

	return &typelist
}

//获取开出4个大奖的列表
func (sgj *ShuiGuoJiLogic) getFourJiangTypeList(money int64) *[][]int32 {
	jiang := int64(0)
	typelist := make([][]int32, 0)
	tmplist := make([]int32, 0)
	for i := 0; i < len(sgj.showTypeData.fiveOdd); i++ {
		jiang += sgj.getUsersBetMoneyByBetType(sgj.showTypeData.fiveOdd[i].sgjSecondType) * int64(sgj.showTypeData.fiveOdd[i].oddValue)
		tmplist = append(tmplist, sgj.showTypeData.fiveOdd[i].sgjType)
	}

	if jiang <= money {
		typelist = append(typelist, tmplist)
	}

	return &typelist
}

//获取开出3个大奖的列表
func (sgj *ShuiGuoJiLogic) getMaxThreeJiangTypeList(money int64) *[][]int32 {
	jiang := int64(0)
	typelist := make([][]int32, 0)
	tmplist := make([]int32, 0)
	for i := 0; i < len(sgj.showTypeData.twentyOdd); i++ {
		jiang += sgj.getUsersBetMoneyByBetType(sgj.showTypeData.twentyOdd[i].sgjSecondType) * int64(sgj.showTypeData.twentyOdd[i].oddValue)
		tmplist = append(tmplist, sgj.showTypeData.twentyOdd[i].sgjType)
	}

	if jiang <= money {
		typelist = append(typelist, tmplist)
	}

	return &typelist
}

//获取开出3个大奖的列表
func (sgj *ShuiGuoJiLogic) getMinThreeJiangTypeList(money int64) *[][]int32 {
	jiang := int64(0)
	secondtype := make([]int32, 0)
	list2 := make([]int32, 0)
	for i := 0; i < len(sgj.showTypeData.tenOdd); i++ {
		bfind := false
		for j := 0; j < len(secondtype); j++ {
			if sgj.showTypeData.tenOdd[i].sgjSecondType == secondtype[j] {
				bfind = true
				break
			}
		}

		if bfind {
			list2 = append(list2, sgj.showTypeData.tenOdd[i].sgjType)
			continue
		}

		secondtype = append(secondtype, sgj.showTypeData.tenOdd[i].sgjSecondType)
		jiang += sgj.getUsersBetMoneyByBetType(sgj.showTypeData.tenOdd[i].sgjSecondType) * int64(sgj.showTypeData.tenOdd[i].oddValue)
	}

	typelist := make([][]int32, 0)
	if jiang <= money {
		//typelist = append(typelist, secondtype)
		typelist = append(typelist, list2)
		mlog.Debug("typelist=%v", typelist)
	}

	return &typelist
}

//获取开出2个大奖的列表
func (sgj *ShuiGuoJiLogic) getMaxTwoJiangTypeList(money int64) *[][]int32 {
	allType := make([]int32, 0)
	allType = append(allType, sgj.showTypeData.allType...)
	typelist := make([][]int32, 0)
	for i := 0; i < len(allType); i++ {
		firstdata := sgj.showTypeData.getShowTypeDataByType(allType[i])
		if firstdata.oddValue == 0 {
			continue
		}
		for j := i + 1; j < len(allType); j++ {
			seconddata := sgj.showTypeData.getShowTypeDataByType(allType[j])
			if seconddata.oddValue == 0 {
				continue
			}
			for m := j + 1; m < len(allType); m++ {
				thriddata := sgj.showTypeData.getShowTypeDataByType(allType[m])
				if thriddata.oddValue == 0 {
					continue
				}
				jiang := int64(0)
				jiang += sgj.getUsersBetMoneyByBetType(firstdata.sgjSecondType) * int64(firstdata.oddValue)
				jiang += sgj.getUsersBetMoneyByBetType(seconddata.sgjSecondType) * int64(seconddata.oddValue)
				jiang += sgj.getUsersBetMoneyByBetType(thriddata.sgjSecondType) * int64(thriddata.oddValue)

				if jiang <= money {
					tmplist := make([]int32, 0)
					tmplist = append(tmplist, allType[i], allType[j], allType[m])
					typelist = append(typelist, tmplist)
				}
			}
		}
	}

	return &typelist
}

//获取开出2个大奖的列表
func (sgj *ShuiGuoJiLogic) getMixTwoJiangTypeList(money int64) *[][]int32 {
	allType := make([]int32, 0)
	allType = append(allType, sgj.showTypeData.allType...)
	typelist := make([][]int32, 0)
	for i := 0; i < len(allType); i++ {
		firstdata := sgj.showTypeData.getShowTypeDataByType(allType[i])
		if firstdata.oddValue == 0 {
			continue
		}
		for j := i + 1; j < len(allType); j++ {
			seconddata := sgj.showTypeData.getShowTypeDataByType(allType[j])
			if seconddata.oddValue == 0 {
				continue
			}
			jiang := int64(0)
			jiang += sgj.getUsersBetMoneyByBetType(firstdata.sgjSecondType) * int64(firstdata.oddValue)
			jiang += sgj.getUsersBetMoneyByBetType(seconddata.sgjSecondType) * int64(seconddata.oddValue)

			if jiang <= money {
				tmplist := make([]int32, 0)
				tmplist = append(tmplist, allType[i], allType[j])
				typelist = append(typelist, tmplist)
			}
		}
	}

	return &typelist
}

//获取权重列表索引值(index:分數區間)
func (sgj *ShuiGuoJiLogic) getWeightIndex() int {
	mlog.Debug("sgj.scoreSection.len=%v", len(sgj.scoreSection))
	index := len(sgj.scoreSection) - 1
	mlog.Debug("当前index=%v", index)
	for i := 0; i < index; i++ {
		example.invalid("总下注额度=%v，sgj.scoreSection[i+1]=%v,sgj.scoreSection[%v]=%v", sgj.betTotal, sgj.scoreSection[i+1], i, sgj.scoreSection[i])
		if sgj.betTotal < sgj.scoreSection[i+1] && sgj.betTotal >= sgj.scoreSection[i] {
			example.invalid("当前选择的分数段为：%v,下注金额为：%v", sgj.scoreSection[i], sgj.betTotal)
			return i
		} else if sgj.betTotal < sgj.scoreSection[0] {
			return 0
		} else if sgj.betTotal == sgj.scoreSection[len(sgj.scoreSection)-1] {
			return len(sgj.scoreSection) - 1
		}
	}
	example.invalid("当前选择的分数段为：%v,下注金额为：%v", sgj.scoreSection[index], sgj.betTotal)
	return len(sgj.scoreSection)
}

//获取普通奖的类型 (money:押注總額 index:權重區間)
func (sgj *ShuiGuoJiLogic) getNormalJiangType(money int64, index int) {
	allType := make([]int32, 0)
	allType = append(allType, sgj.showTypeData.allType...)

	for _, v := range sgj.showTypeData.allType {
		//data:1~24主類別
		data := sgj.showTypeData.getShowTypeDataByType(v)
		if data.oddValue == 0 {
			//去掉倍率為零的選項
			for i := 0; i < len(allType); i++ {
				if allType[i] == v {
					allType = append(allType[:i], allType[i+1:]...)
					break
				}
			}
			continue
		}

		jiang := sgj.getUsersBetMoneyByBetType(data.sgjSecondType) * int64(data.oddValue)
		if jiang > money {
			//去掉不符合的类型
			for i := 0; i < len(allType); i++ {
				if allType[i] == v {
					allType = append(allType[:i], allType[i+1:]...)
					break
				}
			}
		}
	}

	if len(allType) == 0 {
		allType = append(allType, sgj.getMinJiangTypeV2())
		example.invalid("開最小獎, %v", allType)
	}
	sgj.betResultType = sgj.weight[index].getJiangType(allType)
	sgj.resultTypeList = make([]int32, 0)
	sgj.resultTypeList = append(sgj.resultTypeList, sgj.betResultType)
}

//获取表演奖的类型 (index:showtype權重區間)
func (sgj *ShuiGuoJiLogic) getShowJiangType(index int) {
	allType := make([]int32, 0)
	allType = append(allType, sgj.showTypeData.allType...)

	for _, v := range sgj.showTypeData.allType {
		//data:1~24主類別
		data := sgj.showTypeData.getShowTypeDataByType(v)
		if data.oddValue == 0 {
			//去掉倍率為零的選項
			for i := 0; i < len(allType); i++ {
				if allType[i] == v {
					allType = append(allType[:i], allType[i+1:]...)
					break
				}
			}
			continue
		}

	}

	if len(allType) == 0 {
		allType = append(allType, sgj.getMinJiangTypeV2())
	}
	sgj.betResultType = sgj.weight[index].getJiangType(allType)
	sgj.resultTypeList = make([]int32, 0)
	sgj.resultTypeList = append(sgj.resultTypeList, sgj.betResultType)
}

//获取赔偿最小的将的类型
//func (sgj *ShuiGuoJiLogic) getMinJiangType() int32 {
//	//minmoney:最小金額 mintype:最小種類
//	minmoney := int64(-1)
//	mintype := int32(0)
//	for _, v := range sgj.showTypeData.allType {
//		data := sgj.showTypeData.getShowTypeDataByType(v)
//		if data.oddValue == 0 {
//			continue
//		}
//		//計算玩家下注的錢*倍率
//		jiang := sgj.getUsersBetMoneyByBetType(data.sgjSecondType) * int64(data.oddValue)
//		if minmoney == -1 {
//			minmoney = jiang
//			mintype = data.sgjType
//			continue
//		}
//		if jiang < minmoney {
//			minmoney = jiang
//			mintype = data.sgjType
//		} else if jiang == minmoney && rand.Int31n(2) == 1 {
//			minmoney = jiang
//			mintype = data.sgjType
//		}
//	}
//	return mintype
//}

//获取赔偿最小的将的类型
func (sgj *ShuiGuoJiLogic) getMinJiangTypeV2() int32 {
	alltype := make([]int32, 0)
	type2jiang := make(map[int32]int64)
	for _, v := range sgj.showTypeData.allType {
		data := sgj.showTypeData.getShowTypeDataByType(v)
		if data.oddValue == 0 {
			continue
		}
		//計算玩家下注的錢*倍率
		jiang := sgj.getUsersBetMoneyByBetType(data.sgjSecondType) * int64(data.oddValue)
		type2jiang[data.sgjType] = jiang
		alltype = append(alltype, data.sgjType)
	}
	mlog.Debug("==================== alltype before sort %v", alltype)
	//按賠付金額排序 愈小的排愈前面, 金額一樣則隨機
	sort.Slice(alltype, func(i, j int) bool {
		if type2jiang[alltype[i]] == type2jiang[alltype[j]] && rand.Int31n(2) == 1 {
			return true
		}
		return type2jiang[alltype[i]] < type2jiang[alltype[j]]
	})
	mlog.Debug("==================== jiangArray %v", type2jiang)
	mlog.Debug("==================== alltype %v", alltype)
	//取三個賠付最小的然後選一個
	alltype = alltype[0:3]
	mlog.Debug("==================== alltype %v", alltype)
	v := rand.Int31n(2) + 1
	return alltype[v]
}

//获取用户下注的钱，通过类型
func (sgj *ShuiGuoJiLogic) getUsersBetMoneyByBetType(bettype int32) int64 {
	for _, v := range sgj.betData {
		if v.betType == bettype {
			return v.bet
		}
	}

	return 0
}

//增加系统水池的钱
func (sgj *ShuiGuoJiLogic) addSysPool(value int64) {
	sgj.sysPool += value
}

//减少系统水池的钱
func (sgj *ShuiGuoJiLogic) subtractSysPool(value int64) {
	sgj.sysPool -= value
}

//改变系统水池的钱
func (sgj *ShuiGuoJiLogic) changeSysPool(value int64) {
	sgj.startsyspool += value - sgj.sysPool
	sgj.sysPool = value
}

//改变注水百分比
func (sgj *ShuiGuoJiLogic) changeZhuShui(value int32) {
	if value >= 0 && value <= 100 {
		sgj.zhushui = value
	}
}

//查看系统水池有多少钱
func (sgj *ShuiGuoJiLogic) checkSysPool() {
	mlog.Printf("桌子号：%v，系统水池：%v，注水概率：%v，注水Money：%v", sgj.TableNo(), sgj.sysPool, sgj.zhushui, sgj.zhushuimoney)
}

//占位用防止崩溃
func (sgj *ShuiGuoJiLogic) ExchangeMaxCard(userid int32) {

}

func (sgj *ShuiGuoJiLogic) OnRecvHTTPMessage(req *http.Request) proto.Message {
	action := req.PostFormValue("action")
	if action == "getserverinfo" {
		return sgj.getTableListMsg()
	} else if action == "gettableinfo" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		return sgj.getTableDetailMsg(int32(tid))
	} else if action == "modresult" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		ret, _ := strconv.ParseInt(req.PostFormValue("result"), 10, 32)
		adminuser := req.PostFormValue("adminuser")
		return sgj.setNextJiangType(int32(tid), int32(ret), adminuser)
	} else if action == "modpoolamount" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		pool, _ := strconv.ParseInt(req.PostFormValue("poolamount"), 10, 64)
		return sgj.setSysPool(int32(tid), pool)
	} else if action == "modpopwaterrate" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		zhushui, _ := strconv.ParseInt(req.PostFormValue("popwaterrate"), 10, 32)
		return sgj.setZhuShui(int32(tid), int32(zhushui))
	} else if action == "history" {
		pdata := new(netproto.SGJ_LastGameResultType)
		pdata.BetMoney = proto.Int64(0)

		resultList := make([]*netproto.SGJ_JiangTypeOdd, 0)

		if len(sgj.lastResultTypeList) > 20 {
			resultList = sgj.lastResultTypeList[len(sgj.lastResultTypeList)-20 : len(sgj.lastResultTypeList)]
		} else {
			resultList = sgj.lastResultTypeList
		}
		pdata.ResultTypeOdd = append(pdata.ResultTypeOdd, resultList...)

		return pdata
	}

	return nil
}

//获取桌子列表数据
func (sgj *ShuiGuoJiLogic) getTableListMsg() proto.Message {
	pdata := new(netproto.Http_SGJ_GameInfo)
	for _, table := range sgj.GetTable().Room.Tables {
		pTable := new(netproto.Http_SGJ_TableList)
		logic := table.GameLogic.(*ShuiGuoJiLogic)
		pTable.TableNo = proto.Int32(logic.TableNo())
		pTable.Status = proto.Int32(logic.gameState)
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
		pTable.PoolCurrentAmount = proto.Int64(logic.sysPool)
		pTable.ZhuShuiMoney = proto.Int64(logic.zhushuimoney)
		lastturn := logic.restturn
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
func (sgj *ShuiGuoJiLogic) getTableDetailMsg(tableid int32) proto.Message {
	t := sgj.GetTable().Room.GetTableByNo(tableid)
	pTable := new(netproto.Http_SGJ_TableDetail)
	if t != nil {
		logic := t.GameLogic.(*ShuiGuoJiLogic)
		pTable.TableNo = proto.Int32(logic.TableNo())
		pTable.Status = proto.Int32(logic.gameState)
		//		pTable.PoolInitAmount = proto.Int64(logic.startsyspool)
		pTable.PoolCurrentAmount = proto.Int64(logic.sysPool)
		pTable.ZhuShuiMoney = proto.Int64(logic.zhushuimoney)
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
		zhuang := logic.GetUserByUserID(logic.zhuangUserID)
		zhuangtype := int32(0)
		nickname := ""
		currmoney := logic.sysPool
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
		pTable.ZhuangUserID = proto.Int32(logic.zhuangUserID)
		lastturn := logic.restturn
		if lastturn < 0 {
			lastturn = 0
		}
		pTable.ZhuangSurplusTotal = proto.Int32(lastturn)
		pTable.ZhuangInitMoneyAmount = proto.Int64(logic.startmoney)
		pTable.ZhuangCurrentMoneyAmount = proto.Int64(currmoney)
		pTable.SumEarnAmount = proto.Int64(logic.systotalwin)
		pTable.LastEarnAmount = proto.Int32(int32(logic.syslastwin))
		pTable.SysTaxIncome = proto.Int64(logic.systaxincome)
		pTable.OpenResultApple = proto.Int64(logic.getRealUserWinByType(PGMAX))
		pTable.OpenResultAppleMin = proto.Int64(logic.getRealUserWinByType(PGMIN))
		pTable.OpenResultOrange = proto.Int64(logic.getRealUserWinByType(JZMAX))
		pTable.OpenResultOrangeMin = proto.Int64(logic.getRealUserWinByType(JZMIN))
		pTable.OpenResultMango = proto.Int64(logic.getRealUserWinByType(MGMAX))
		pTable.OpenResultMangoMin = proto.Int64(logic.getRealUserWinByType(MGMIN))
		pTable.OpenResultBell = proto.Int64(logic.getRealUserWinByType(LDMAX))
		pTable.OpenResultBellMin = proto.Int64(logic.getRealUserWinByType(LDMIN))
		pTable.OpenResultWatermelon = proto.Int64(logic.getRealUserWinByType(XGMAX))
		pTable.OpenResultWatermelonMin = proto.Int64(logic.getRealUserWinByType(XGMIN))
		pTable.OpenResultDobuleStar = proto.Int64(logic.getRealUserWinByType(SXMAX))
		pTable.OpenResultDobuleStarMin = proto.Int64(logic.getRealUserWinByType(SXMIN))
		pTable.OpenResult77 = proto.Int64(logic.getRealUserWinByType(QQMAX))
		pTable.OpenResult77Min = proto.Int64(logic.getRealUserWinByType(QQMIN))
		pTable.OpenResultBar = proto.Int64(logic.getRealUserWinByType(BARMAX))
		pTable.OpenResultBarMin = proto.Int64(logic.getRealUserWinByType(BARMIN))
		pTable.OpenResultLuckLose = proto.Int64(logic.getRealUserWinByType(0))
		pTable.BetTime = proto.Int32(sgj.getBetTime())

		if logic.gameState == actionResult {
			pTable.CurentOpenResult = proto.String(logic.getResultString())
			win := int64(0)
			for i := 0; i < len(logic.resultTypeList); i++ {
				win += logic.getRealUserWinByType(logic.resultTypeList[i])
			}
			pTable.CurentOpenResultEarn = proto.Int64(win)
		}

		for _, ud := range logic.userdata {
			userinfo := new(netproto.Http_SGJ_UserInfo)
			userinfo.UserID = proto.Int32(ud.user.UserID)
			ut := int32(1)
			if ud.user.IsRobot() {
				ut = 0
			}
			userinfo.UserType = proto.Int32(ut)
			userinfo.NickName = proto.String(ud.user.NickName)
			userinfo.GameScore = proto.Int64(ud.user.UpdateScore)
			userinfo.InitMoney = proto.Int64(ud.user.LoginScore)
			userinfo.CurrentMoney = proto.Int64(ud.user.Score)
			userinfo.BetInfoUserApple = proto.Int64(ud.getUserBetByType(PG))
			userinfo.BetInfoUserOrange = proto.Int64(ud.getUserBetByType(JZ))
			userinfo.BetInfoUserMango = proto.Int64(ud.getUserBetByType(MG))
			userinfo.BetInfoUserBell = proto.Int64(ud.getUserBetByType(LD))
			userinfo.BetInfoUserWatermelon = proto.Int64(ud.getUserBetByType(XG))
			userinfo.BetInfoUserDobuleStar = proto.Int64(ud.getUserBetByType(SX))
			userinfo.BetInfoUser77 = proto.Int64(ud.getUserBetByType(QQ))
			userinfo.BetInfoUserBar = proto.Int64(ud.getUserBetByType(BAR))
			userinfo.SumCharge = proto.Int64(ud.user.SumCharge)
			userinfo.SumConvert = proto.Int64(ud.user.SumConvert)
			pTable.UserInfo = append(pTable.UserInfo, userinfo)
		}
	}

	return pTable
}

func (sgj *ShuiGuoJiLogic) getResultString() string {
	var res string
	if len(sgj.resultTypeList) <= 0 {
		return res
	}

	res += resultstring[sgj.resultTypeList[0]]
	for i := 1; i < len(sgj.resultTypeList); i++ {
		res += "," + resultstring[sgj.resultTypeList[i]]
	}

	return res
}

func (sgj *ShuiGuoJiLogic) getRetMsg(ret bool, msg string) proto.Message {
	pdata := new(netproto.Http_Result)
	pdata.Result = proto.Bool(ret)
	pdata.Message = proto.String(msg)
	return pdata
}

//设置下次开奖结果
func (sgj *ShuiGuoJiLogic) setNextJiangType(tableid int32, jiang int32, adminuser string) proto.Message {
	t := sgj.GetTable().Room.GetTableByNo(tableid)

	if t != nil {
		logic := t.GameLogic.(*ShuiGuoJiLogic)
		if adminuser != "" {
			logic.GameLog.SetNextJiangRecord(adminuser, fmt.Sprintf("%d", jiang))
		}
		logic.nextjiangtype = jiang
		mlog.Debug("设置开奖结果为%s", resultstring[jiang])
		if jiang == LOCK {
			logic.nextjiangtype = 0
		}
		return sgj.getRetMsg(true, "修改开奖结果成功")
	}

	return sgj.getRetMsg(false, "修改开奖结果失败")
}

func (sgj *ShuiGuoJiLogic) setSysPool(tableid int32, pool int64) proto.Message {
	if tableid == -1 {
		for _, table := range sgj.GetTable().Room.Tables {
			logic := table.GameLogic.(*ShuiGuoJiLogic)
			logic.changeSysPool(pool)
		}
	} else {
		t := sgj.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*ShuiGuoJiLogic)
			logic.changeSysPool(pool)
		} else {
			return sgj.getRetMsg(false, "修改系统水池失败")
		}
	}

	return sgj.getRetMsg(true, "修改系统水池成功")
}

func (sgj *ShuiGuoJiLogic) setZhuShui(tableid int32, zhushui int32) proto.Message {
	if tableid == -1 {
		for _, table := range sgj.GetTable().Room.Tables {
			logic := table.GameLogic.(*ShuiGuoJiLogic)
			logic.changeZhuShui(zhushui)
		}
	} else {
		t := sgj.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*ShuiGuoJiLogic)
			logic.changeZhuShui(zhushui)
		} else {
			return sgj.getRetMsg(false, "修改注水率失败")
		}
	}

	return sgj.getRetMsg(true, "修改注水率成功")
}

func (sgj *ShuiGuoJiLogic) getRealUserWinByType(t int32) int64 {
	win := int64(0)
	odd := int64(sgj.showTypeData.getOddByType(t))
	user := sgj.GetUserByUserID(sgj.zhuangUserID)
	st := sgj.showTypeData.getSecondTypeByType(t)

	if user != nil && !user.IsRobot() {
		_, tax := sgj.GetUserWinAndTax(sgj.betTotal)
		if sgj.currtax != 0 {
			tax = sgj.currtax
		}
		robotwin := odd * sgj.getRobotBetMoneyByBetType(st)
		win = sgj.betRobotTotal - robotwin - tax
		return -win
	}

	for _, ud := range sgj.userdata {
		if !ud.user.IsRobot() {
			win += odd * ud.getUserBetByType(st)
			win -= ud.money
		}
	}

	mlog.Debug("玩家赢钱%d", win)
	return -win
}

//排序身上攜帶的錢大到小
func (sgj *ShuiGuoJiLogic) getMoneyRanking() {
	sort.Sort(sgj.userscore)
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
func (sgj *ShuiGuoJiLogic) getLuckyStar() {
	t := 0
	gameuser := sgj.AllGameUser()
	if sgj.luckystar == 0 || sgj.notbettimes >= 3 {
		//如果沒有幸運星、連續三場沒下注或是幸運星的持續局數到了，隨機一個玩家
		t = 1
	} else {
		//檢查玩家還在不在場
		for _, u := range gameuser {
			if u.UserID == sgj.luckystar {
				sgj.luckystartimes += 1
				break
			} else {
				t = 1
			}
		}
	}
	if t == 1 || sgj.luckystartimes == sgj.luckystarround {
		if len(gameuser) != 0 {
			userIndex := rand.Intn(len(gameuser))
			sgj.luckystar = gameuser[userIndex].UserID
			sgj.luckyfaceid = gameuser[userIndex].FaceID
		} else {
			//如果沒有明燈，須和client確認要傳什麼
			sgj.luckystar = 0
			sgj.luckyfaceid = 0
		}
		sgj.luckystartimes = 0
		sgj.notbettimes = 0
	}
}

//獲取富豪玩家
func (sgj *ShuiGuoJiLogic) getRichestList() {
	userscore := new(userScore)
	userscore.userdata = sgj.AllTableUser()
	sort.Sort(userscore)
	//計算分數最高的7個
	for index, user := range userscore.userdata {
		if index < 5 {
			rich := new(netproto.SGJ_Richest)
			rich.Ranking = proto.Int32(int32(index + 1))
			rich.UserID = proto.Int32(user.UserID)
			rich.UserFaceID = proto.Int32(user.FaceID)
			rich.Money = proto.Int64(user.Score)
			sgj.richestList = append(sgj.richestList, rich)
			sgj.richestuserinfo = append(sgj.richestuserinfo, user.UserID)
		} else {
			break
		}
	}
	sgj.stageuserinfo = append(sgj.richestuserinfo, sgj.luckystar)
}

//获取大富豪列表的消息
func (sgj *ShuiGuoJiLogic) getRichestListMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SGJ_GameMessageClassID_SGJRichestListID)
	pdata := new(netproto.SGJ_RichestList)
	pdata.RichestList = sgj.richestList
	//如果當下沒有玩家 回傳空值
	pmsg.MsgData = pdata
	return pmsg
}

//获取幸運星的消息
func (sgj *ShuiGuoJiLogic) getLuckyStarMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SGJ_GameMessageClassID_SGJLuckyStarListID)
	luckyStarList := new(netproto.SGJ_LuckyStar)
	luckyStarList.UserID = proto.Int32(sgj.luckystar)
	luckyStarList.UserFaceID = proto.Int32(sgj.luckyfaceid)
	pmsg.MsgData = luckyStarList
	return pmsg
}

//获取用户数据通过用户ID
func (sgj *ShuiGuoJiLogic) getUserDataByID(userid int32) *sgjUserData {
	for _, ud := range sgj.userdata {
		if ud.user.UserID == userid {
			return ud
		}
	}
	return nil
}

//收到RPC消息解析處理
func (sgj *ShuiGuoJiLogic) OnRecvRPCMessage(m *msg.Message) {

}
