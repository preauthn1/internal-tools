package fkeb

import (
	"math"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	proto "example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
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
	ZHUANG = 0 //庄
	SHUN   = 1 //顺
	TIAN   = 2 //天
	DI     = 3 //地
)

//開幾門
const (
	TONGSHA = 0
	XIAWIN1 = 1
	XIAWIN2 = 2
	TONGPEI = 3
	NATURL  = 99
)

var resultstring = []string{"庄家通杀", "顺门", "天门", "地门"}

const (
	GameKaiJiang_Action = "kaijiang"
	GameEnd_Action      = "end"
)

const (
	CtrlBetLimitDefault int64 = 100 // 監控名單預設注額
)

type majiang struct {
	majiang1 int32
	majiang2 int32
}

type fkebweight struct {
	count  []int32 //开奖的个数
	weight []int32 //开奖的权重
}

type FKEBLogic struct {
	*gamelogic.ComputeGameLogicBase
	serverfile         *os.File
	opttimer           *event.TimerFunc
	bettotal           int64                         //下注总额
	bettotalrobot      int64                         //机器人下注总额
	bettypetotal       []*bettype                    //下注分类汇总
	bettyperobot       []*bettype                    //机器人下注分类汇总
	betvalue           *betconfig                    //下注配置文件
	userdata           []*fkebuserdata               //用户数据
	majiangresult      []*majiang                    //开奖的牌
	gamestate          int32                         //游戏状态
	bettime            int32                         //下注时间
	jiangtype          []int32                       //开奖类型
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
	resultcardlist     *netproto.FKEB_ResultCardList //开奖的历史牌
	zhushui            int32                         //注水率
	maxclearcount      int32                         //清理数据上限
	maxrobotzhuang     int32                         //最多有好多个机器人上庄
	robotxiazhuang     int32                         //机器人下庄概率
	robotzhuang        int32                         //没有人当庄的时候机器人当庄的概率
	winline            int64                         //机器人赢钱线大于这个值随机开。小于这个值机器人赢最多的开
	zhushuimoney       int64                         //注水的钱
	startsyspool       int64                         //开始时的金币数量
	startmoney         int64                         //上庄时的金币
	currtax            int64                         //当前抽水
	nextjiangtype      []int32                       //下次开奖类型
	systotalwin        int64                         //系统收入
	syslastwin         int64                         //上一局收益
	scoresection       []int64                       //分数段用来加载不同的分数段配置
	weight             []*fkebweight                 //权重配置文件
	majiang            []int32                       //麻将点数1-9
	betlimitmoney      int64                         //下注限制
	zrusertongpei      int32                         //真人上莊通賠
	realuserzhuangflow int64                         //真人上莊累計流水總額
	ruztotalwinlose    int64                         //真人上莊累計輸贏
	jiangtyperate
	assignRewardDetail *netproto.SuperToolAssignRewardDetail //指定開獎
	SuperToolService   *gamelogic.SuperToolBrGameServiceBase //超端服務
	ctrlbetlimit       int64                                 //監控名單會控制的最低注額
}

type jiangtyperate struct {
	naturl    int32 //自然機率
	tongsha   int32 //通殺機率
	tongpei   int32 //通陪機率
	win1loss2 int32 //庄贏一輸二機率
	win2loss1 int32 //庄贏二輸一機率
}

//新建一个疯狂二八实例
func NewFKEBLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *FKEBLogic {
	glc := new(FKEBLogic)
	glc.ComputeGameLogicBase = gamelogic.NewComputeGameLogicBase(t, tfManager, cfgs)
	glc.SuperToolService = glc.GetSuperToolService().(*gamelogic.SuperToolBrGameServiceBase)
	glc.getfkebconfig(cfgs)
	glc.gamestate = Wait
	glc.zhuangID = 0
	glc.syslastwin = 0
	glc.systotalwin = 0
	glc.nextjiangtype = make([]int32, 0)
	glc.resultcardlist = new(netproto.FKEB_ResultCardList)
	glc.getLastPool()
	glc.startsyspool = glc.syspool
	glc.addCmd()
	glc.RegisterAfterGameStartFunc(func() {
		glc.clearGameData() //清除数据
		glc.huangZhuang()
		glc.GameLog.AddGameStep("庄家%d，剩余%d", glc.zhuangID, glc.lastturn)
		glc.zhuShui()
		glc.broadCastZhuangInfo()
		glc.startBet()
	})

	return glc
}

//添加命令行
func (fkeb *FKEBLogic) addCmd() {
	command.AddCommand(NewFKEBCmd(fkeb))
}

//广播庄家
func (fkeb *FKEBLogic) broadCastZhuangInfo() {
	fkeb.BoardcastMessage(fkeb.getZhuangInfoMsg())
}

//注水到下注额
func (fkeb *FKEBLogic) zhuShui() {
	if fkeb.zhuangID == 0 {
		fkeb.bettotal = fkeb.syspool * int64(fkeb.zhushui) / 100
		fkeb.syspool -= fkeb.bettotal
	} else {
		ur := fkeb.GetUserByUserID(fkeb.zhuangID)
		if ur != nil && ur.IsRobot() {
			fkeb.bettotal = fkeb.syspool * int64(fkeb.zhushui) / 100
			fkeb.syspool -= fkeb.bettotal
		}
	}

	fkeb.zhushuimoney = fkeb.bettotal
}

//用户离开
func (fkeb *FKEBLogic) UserLeave(gameuser *gamelogic.GameUser) bool {
	if fkeb.gamestate != Wait {
		if fkeb.getUserDataByID(gameuser.UserID) != nil {
			mlog.Debug("用户在玩游戏还不能离开")
			return false
		}

		if fkeb.zhuangID == gameuser.UserID {
			mlog.Debug("用户在当庄，所以用户不能离开")
			return false
		}
	}

	fkeb.delFromZhuangList(gameuser.UserID)

	fkeb.ChangeUserStateEnd(gameuser)
	return fkeb.ComputeGameLogicBase.UserLeave(gameuser)
}

//注销用户实例
func (fkeb *FKEBLogic) UserFree(gameuser *gamelogic.GameUser) {
	if fkeb.gamestate != Wait {
		if fkeb.getUserDataByID(gameuser.UserID) != nil {
			mlog.Debug("用户在玩游戏还不能离开")
			return
		}

		if fkeb.zhuangID == gameuser.UserID {
			mlog.Debug("用户在当庄，所以用户不能离开")
			return
		}
	} else {
		if fkeb.zhuangID == gameuser.UserID {
			fkeb.zhuangID = 0
		}
	}

	fkeb.delFromZhuangList(gameuser.UserID)

	fkeb.ChangeUserStateEnd(gameuser)
	fkeb.ComputeGameLogicBase.UserFree(gameuser)
}

//断线重连
func (fkeb *FKEBLogic) UserRecover(gameuser *gamelogic.GameUser) {
	mlog.Debug("UserRecover")
	fkeb.ComputeGameLogicBase.UserRecover(gameuser)
	ud := fkeb.getUserDataByID(gameuser.UserID)
	if ud != nil {
		ud.money += ud.user.Score - ud.moneylast
		ud.moneylast = ud.user.Score
	}

	fkeb.sendScene(gameuser)
}

//用户进入
func (fkeb *FKEBLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {
	gameuser.NickName = strconv.Itoa(int(gameuser.UserID))
	fkeb.ComputeGameLogicBase.UserEnter(gameuser, tableno, seatno)
	fkeb.sendScene(gameuser)
}

//发送场景
func (fkeb *FKEBLogic) sendScene(gameuser *gamelogic.GameUser) {
	//发送设置

	if gameuser.IsRobot() {
		mlog.Debug("给机器人%d发送场景消息", gameuser.UserID)
	} else {
		mlog.Debug("给玩家%d发送场景消息", gameuser.UserID)
	}
	fkeb.sendSetting(gameuser)
	fkeb.SendMessageToUser(gameuser, fkeb.getResultCardListMsg())
	if fkeb.gamestate == Wait {
		return
	}
	//发送阶段
	fkeb.SendMessageToUser(gameuser, fkeb.getStateMsg())
	fkeb.SendMessageToUser(gameuser, fkeb.getZhuangInfoMsg())
	fkeb.SendMessageToUser(gameuser, fkeb.getBetTotalMsg(0))
	fkeb.sendUserBetInfo(gameuser)

	if fkeb.gamestate == KaiJiangAction {
		//发送开奖信息
		fkeb.SendMessageToUser(gameuser, fkeb.getGameResultMajiangMsg())
	} else if fkeb.gamestate == EndAction {
		//发送用户下注输赢信息如果有
		ud := fkeb.getUserDataByID(gameuser.UserID)
		if ud != nil {
			fkeb.SendMessageToUser(gameuser, fkeb.getUserWinMsg(ud.win, fkeb.zhuangWin))
		}
	}
}

//发送设置
func (fkeb *FKEBLogic) sendSetting(gameuser *gamelogic.GameUser) {
	mlog.Debug("sendSetting")
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.FKEB_GameMessageClassID_FKEBSettingID)

	pdata := new(netproto.FKEB_Setting)
	pdata.BetCount = append(pdata.BetCount, fkeb.betcount...)
	pdata.BetMoney = proto.Int32(int32(fkeb.betvalue.getMoneyByType(1)))
	pdata.ZhuangLimit = proto.Int64(fkeb.zhuanglimit)
	pdata.Betlimitmoney = proto.Int64(fkeb.betlimitmoney)
	pmsg.MsgData = pdata
	fkeb.SendMessageToUser(gameuser, pmsg)
}

//发送用户的下注信息
func (fkeb *FKEBLogic) sendUserBetInfo(gameuser *gamelogic.GameUser) {
	ud := fkeb.getUserDataByID(gameuser.UserID)
	if ud == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.FKEB_GameMessageClassID_FKEBBetInfoID)

	pdata := new(netproto.FKEB_UserBetInfo)
	pdata.UserMoney = proto.Int64(ud.money)
	for _, v := range ud.userbet {
		betinfo := new(netproto.FKEB_BroadCastBet)
		betinfo.Type = proto.Int32(v.fkebtype)
		betinfo.Money = proto.Int64(v.money)
		pdata.BetInfo = append(pdata.BetInfo, betinfo)
	}
	pmsg.MsgData = pdata
	fkeb.SendMessageToUser(gameuser, pmsg)
}

//清除数据
func (fkeb *FKEBLogic) clearGameData() {
	fkeb.userdata = make([]*fkebuserdata, 0)
	fkeb.bettypetotal = make([]*bettype, 0)
	fkeb.bettyperobot = make([]*bettype, 0)
	fkeb.majiangresult = make([]*majiang, 0)
	fkeb.winuserinfo = make([]int32, 0)
	fkeb.jiangtype = make([]int32, 0)
	fkeb.bettotal = 0
	fkeb.bettotalrobot = 0
	fkeb.zhuangWin = 0
	fkeb.currtax = 0
	rand.Seed(time.Now().UnixNano())
	zlen := len(fkeb.zhuanglist)
	if zlen == 0 {
		fkeb.robotZhuang()
	}

	fkeb.robotXiaZhuang()
}

//初始化麻将1-9
func (fkeb *FKEBLogic) initmajiang() {
	fkeb.majiang = make([]int32, 0)
	//麻将编号1-10
	for i := int32(1); i <= 10; i++ {
		for n := 0; n < 4; n++ {
			fkeb.majiang = append(fkeb.majiang, i)
		}
	}
	//洗牌
	for i := 0; i < 40; i++ {
		n := rand.Int31n(40)
		fkeb.majiang[i], fkeb.majiang[n] = fkeb.majiang[n], fkeb.majiang[i]
	}
}

//机器人下庄
func (fkeb *FKEBLogic) robotXiaZhuang() {
	if fkeb.zhuangID == 0 {
		return
	}

	ur := fkeb.GetUserByUserID(fkeb.zhuangID)
	if ur != nil && ur.IsRobot() {
		if rand.Int31n(100) < fkeb.robotxiazhuang {
			fkeb.zhuangID = 0
		}
	}
}

//机器人上庄
func (fkeb *FKEBLogic) robotZhuang() {
	//概率当庄
	if rand.Int31n(100) > fkeb.robotzhuang {
		return
	}
	//统计机器人个数
	robotarr := make([]int32, 0)
	for _, ur := range fkeb.AllTableUser() {
		if ur.IsRobot() {
			robotarr = append(robotarr, ur.UserID)
		}
	}

	if len(robotarr) <= 0 {
		return
	}

	//随机有几个人当庄
	n := rand.Int31n(fkeb.maxrobotzhuang) + 1

	for i := int32(0); i < n; i++ {
		rcn := len(robotarr)
		if rcn <= 0 {
			break
		}
		index := rand.Int31n(int32(rcn))
		fkeb.addToZhuangList(robotarr[index])
		robotarr = append(robotarr[:index], robotarr[index+1:]...)
	}
}

//开始下注
func (fkeb *FKEBLogic) startBet() {
	fkeb.opttimer = fkeb.TimerFuncManager.NewTimer(fkeb.bettime, fkeb.kaiJiang)
	fkeb.setState(Bet)
	fkeb.BoardcastMessage(fkeb.getStateMsg())
}

//开奖动画后的消息
func (fkeb *FKEBLogic) actionend() {

}

//开始开奖
func (fkeb *FKEBLogic) kaiJiang() {
	// fmt.Println(fkeb.nextjiangtype[0])
	// fkeb.nextjiangtype = append(fkeb.nextjiangtype, SHUN, TIAN) // 測試用的指定開獎
	fkeb.setState(KaiJiang) // 中斷點設在這
	key := len(fkeb.nextjiangtype)

	if !fkeb.GetUserByUserID(fkeb.zhuangID).IsRobot() {
		fkeb.realuserzhuangflow += fkeb.bettotal
	}

	if fkeb.zhuangID != 0 {
		if fkeb.GetUserByUserID(fkeb.zhuangID).IsRobot() {
			fkeb.removeRobotBet()
		}
	} else {
		fkeb.removeRobotBet()
	}

	//税收
	mlog.Debug("钱总数为：%v", fkeb.bettotal)
	money, tax := fkeb.GetUserWinAndTax(fkeb.bettotal)
	fkeb.currtax = tax
	mlog.Debug("税收后钱总数为：%v", money)

	//fkeb.bettotal = money
	if fkeb.zhuangID != 0 {
		ur := fkeb.GetUserByUserID(fkeb.zhuangID)
		if ur != nil && !ur.IsRobot() {
			money += ur.Score
		}
	}

	fkeb.getResultType(money)
	fkeb.getMaJiang(key)
	fkeb.opttimer = fkeb.TimerFuncManager.NewTimer(fkeb.kaijiangactiontime+fkeb.endactiontime, fkeb.actionend)
	fkeb.setState(KaiJiangAction)
	fkeb.BoardcastMessage(fkeb.getStateMsg())
	fkeb.BoardcastMessage(fkeb.getGameResultMajiangMsg())
	fkeb.addResultCardList()
	fkeb.action_str = fkeb.RegisterAnimation(GameKaiJiang_Action, fkeb.kaijiangactiontime, fkeb.sendResult)
}

//添加开奖到历史列表
func (fkeb *FKEBLogic) addResultCardList() {
	cards := new(netproto.FKEB_Majiangs)
	for i := 0; i < len(fkeb.majiangresult); i++ {
		card := new(netproto.FKEB_Majiang)
		card.Majiang1 = proto.Int32(fkeb.majiangresult[i].majiang1)
		card.Majiang2 = proto.Int32(fkeb.majiangresult[i].majiang2)
		cards.Majiangs = append(cards.Majiangs, card)
	}

	fkeb.resultcardlist.MaJiangList = append(fkeb.resultcardlist.MaJiangList, cards)
	if len(fkeb.resultcardlist.MaJiangList) > 20 {
		fkeb.resultcardlist.MaJiangList = append(fkeb.resultcardlist.MaJiangList[:0], fkeb.resultcardlist.MaJiangList[1:]...)
	}
}

//获取开奖历史列表的消息
func (fkeb *FKEBLogic) getResultCardListMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.FKEB_GameMessageClassID_FKEBResultCardListID)

	pdata := new(netproto.FKEB_ResultCardList)
	pdata.MaJiangList = append(pdata.MaJiangList, fkeb.resultcardlist.MaJiangList...)
	mlog.Debug("getResultCardListMsg:%v", pdata.MaJiangList)
	pmsg.MsgData = pdata
	return pmsg
}

//获取开奖结果扑克牌消息
func (fkeb *FKEBLogic) getGameResultMajiangMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.FKEB_GameMessageClassID_FKEBGameResultMajiangID)

	pdata := new(netproto.FKEB_GameResultMajiang)
	pdata.RetType = append(pdata.RetType, fkeb.jiangtype...)
	pdata.SeZi = append(pdata.SeZi, rand.Int31n(6)+1)
	pdata.SeZi = append(pdata.SeZi, rand.Int31n(6)+1)
	for _, v := range fkeb.majiangresult {
		cards := new(netproto.FKEB_Majiang)
		cards.Majiang1 = proto.Int32(v.majiang1)
		cards.Majiang2 = proto.Int32(v.majiang2)
		pdata.Majiang = append(pdata.Majiang, cards)
	}

	pmsg.MsgData = pdata
	return pmsg
}

//获取开奖类型。
func (fkeb *FKEBLogic) getResultType(betmoney int64) {
	//指定开奖

	if len(fkeb.nextjiangtype) > 0 {
		fkeb.jiangtype = append(fkeb.jiangtype, fkeb.nextjiangtype...)
		fkeb.nextjiangtype = make([]int32, 0)
		return
	}

	//TODO 系統庄 OR 機器人上莊 純自然幾率
	if fkeb.zhuangID == 0 || fkeb.GetUserByUserID(fkeb.zhuangID).IsRobot() {
		naturlKaijang(fkeb)
	}

	//真人上莊, (贏-輸)>0時,((贏-輸))/流水*0.25=通賠機率,其他狀況皆是自然機率
	if !fkeb.GetUserByUserID(fkeb.zhuangID).IsRobot() && fkeb.ruztotalwinlose != 0 && fkeb.realuserzhuangflow != 0 && fkeb.zrusertongpei != 0 {
		tongpeirate := float32(fkeb.ruztotalwinlose) / float32(fkeb.realuserzhuangflow) * (float32(fkeb.zrusertongpei) / 100)
		rnd := util.RndFloat()
		if rnd < tongpeirate {
			fkeb.jiangtype = append(fkeb.jiangtype, SHUN, TIAN, DI)
		}
	}

	if len(fkeb.jiangtype) == 0 {
		naturlKaijang(fkeb)
	}

	//fkeb.jiangtype = append(fkeb.jiangtype, 1)
	//fkeb.jiangtype = append(fkeb.jiangtype, 2)
	//fkeb.jiangtype = append(fkeb.jiangtype, 3)

	//t := 0
	//if fkeb.zhuangID != 0 && !fkeb.GetUserByUserID(fkeb.zhuangID).IsRobot() &&
	//	fkeb.syspool < fkeb.winline {
	//	t = 1
	//}
	//
	////机器人和系统当庄的时候方案一
	//if t == 0 {
	//	//玩家上庄无通杀
	//	if fkeb.zhuangID <= 0 || fkeb.GetUserByUserID(fkeb.zhuangID).IsRobot() {
	//		if rand.Int31n(100) < fkeb.tongsha {
	//			fkeb.jiangtype = append(fkeb.jiangtype, 0)
	//			mlog.Debug("通杀")
	//			return
	//		}
	//	}
	//
	//	alltype := fkeb.getAllKaiJiangType(betmoney)
	//
	//	if len(alltype) == 0 {
	//		mlog.Debug("钱不够陪，通杀")
	//		fkeb.jiangtype = append(fkeb.jiangtype, 0)
	//		return
	//	}
	//	fkeb.getJiangType(alltype, betmoney)
	//} else {
	//	fkeb.getMaxRobotJiang()
	//	//fkeb.jiangtype = append(fkeb.jiangtype, fkeb.getMaxRobotJiang())
	//	//fkeb.jiangtype = fkeb.getMaxRobotJiang()
	//}
	//
	//mlog.Debug("开奖类型：%v", fkeb.jiangtype)
}

//自然機率開獎
func naturlKaijang(fkeb *FKEBLogic) {
	time.Sleep(1 * time.Nanosecond)
	rand.Seed(time.Now().UnixNano())
	//jiangCount := util.RndInt(0, 4)

	jiangCount := 0
	sumrate := int32(0)
	rndRate := util.RndInt(1, 1001)

	totalWeight := make(map[int32]int32, 0)
	totalWeight[TONGSHA] = fkeb.tongsha
	totalWeight[XIAWIN1] = fkeb.win2loss1
	totalWeight[XIAWIN2] = fkeb.win1loss2
	totalWeight[TONGPEI] = fkeb.tongpei
	totalWeight[NATURL] = fkeb.naturl

	for k, v := range totalWeight {
		sumrate += v
		if sumrate >= int32(rndRate) {
			jiangCount = int(k)
			break
		}
	}

	if jiangCount > 0 && jiangCount != NATURL {
		for i := 1; i <= jiangCount; i++ {
			time.Sleep(1 * time.Nanosecond)
			rand.Seed(time.Now().UnixNano())
			jiangtype := util.RndInt(1, 4)
			hasJiangType := false
			for c := 0; c < len(fkeb.jiangtype); c++ {
				if fkeb.jiangtype[c] == int32(jiangtype) {
					hasJiangType = true
				}
			}
			if !hasJiangType {
				fkeb.jiangtype = append(fkeb.jiangtype, int32(jiangtype))
			}

		}
	} else {
		fkeb.jiangtype = append(fkeb.jiangtype, ZHUANG)
	}
}

//获取机器人赢最多的开奖方案，开两个组合
func (fkeb *FKEBLogic) getMaxRobotJiang() int32 {
	maxtype := int32(0)
	maxmoney := int64(0)
	secondtype := int32(0)
	//secondmoney := int64(0)
	for _, rtd := range fkeb.bettyperobot {
		money := int64(fkeb.betvalue.getOddByType(rtd.fkebtype)) * rtd.money
		if money >= maxmoney {
			secondtype = maxtype
			maxtype = rtd.fkebtype
			maxmoney = money
		}
	}

	fkeb.jiangtype = append(fkeb.jiangtype, maxtype)
	if secondtype != 0 {
		fkeb.jiangtype = append(fkeb.jiangtype, secondtype)
	}
	if len(fkeb.jiangtype) <= 0 {
		mlog.Fatal("1开奖出问题了，未开出奖")
	}
	return maxtype
}

//获取机器人赢钱
func (fkeb *FKEBLogic) getRobotWin() int64 {
	win := int64(0)
	for _, rtd := range fkeb.bettyperobot {
		for _, t := range fkeb.jiangtype {
			if rtd.fkebtype == t {
				win += rtd.money * int64(fkeb.betvalue.getOddByType(t))
				break
			}
		}
	}

	return win
}

//给玩家发送结果
func (fkeb *FKEBLogic) sendResult() {
	fkeb.syslastwin = 0
	fkeb.setState(EndAction)
	correctNumberWin := int64(0) // 機器人在非開獎區域的下注金額，一開始設定為零，之後會更新

	for _, ud := range fkeb.userdata {
		fkeb.sendUserResult(ud, &correctNumberWin)
	}

	fkeb.zhuangWin = int64(0)
	var tax int64
	if fkeb.zhuangID == 0 {
		income, _ := fkeb.GetUserWinAndTax(fkeb.bettotal)
		fkeb.zhuangWin = fkeb.bettotalrobot + fkeb.bettotal
		fkeb.syspool += income
		mlog.System("无人当庄系统水池：%d, bettotal[%v]", fkeb.syspool, fkeb.bettotal)
	} else {
		ur := fkeb.GetUserByUserID(fkeb.zhuangID)
		if ur != nil {
			//finalBetTotal := lh.bettotal + lh.bettotalrobot
			//income, tax = lh.GetUserWinAndTax(finalBetTotal)
			if ur.IsRobot() {
				zhuangTotalWin := fkeb.bettotal + fkeb.bettotalrobot
				zhuangIncome, zhuangTax := fkeb.GetUserWinAndTax(zhuangTotalWin)
				tax = zhuangTax
				fkeb.zhuangWin = zhuangIncome

				syspoolIncome := zhuangTotalWin
				if syspoolIncome < 0 {
					syspoolIncome = fkeb.bettotal - correctNumberWin
				}
				fkeb.syspool += syspoolIncome
				fkeb.UpdateUserScoreFunc(ur.UserID, zhuangIncome, zhuangTax)
			} else {
				//robotwin := fkeb.getRobotWin()
				fkeb.syspool += -fkeb.bettotalrobot
				mlog.System("真人庄userid[%v] 当前水池的钱：%d bettotalrobot[%v]", ur.UserID, fkeb.syspool, fkeb.bettotalrobot)

				if fkeb.bettotal > ur.Score {
					fkeb.syspool += fkeb.bettotal - ur.Score
					mlog.System("庄家赢钱超过系统水池：%d", fkeb.syspool)
					fkeb.zhuangWin = ur.Score
				} else {
					//fkeb.zhuangWin = fkeb.bettotalrobot
				}
			}

			if !ur.IsRobot() {
				zhuangTotalWin := fkeb.bettotal
				zhuangIncome, zhuangTax := fkeb.GetUserWinAndTax(zhuangTotalWin)
				tax = zhuangTax
				fkeb.zhuangWin = zhuangIncome
				fkeb.UpdateUserScoreFunc(ur.UserID, zhuangIncome, zhuangTax)
				fkeb.ruztotalwinlose += zhuangTotalWin
			}
			//計算流水
			var flowAmount, betAmount, payAmount int64
			flowAmount = util.ABSInt64(fkeb.zhuangWin)
			betAmount = flowAmount + tax
			payAmount = betAmount + flowAmount
			fkeb.GameLog.AddGameStep("莊家贏錢 %d, 抽稅 %d, 流水 %d, 有效投注 %d, 返還金額 %d", fkeb.zhuangWin, tax, flowAmount, betAmount, payAmount)
			ur.FlowAmount += flowAmount
			ur.BetAmount += betAmount
			ur.PayAmount += payAmount

			fkeb.SendMessageToUser(ur, fkeb.getUserWinMsg(fkeb.zhuangWin, fkeb.zhuangWin))
		} else {
			mlog.System("庄家数据为空")
		}
	}

	uscoresRet := make(map[int32]int64)
	cscoresRet := make(map[int32]int64)
	for _, ud := range fkeb.userdata {
		fkeb.SendMessageToUser(ud.user, fkeb.getUserWinMsg(ud.win, fkeb.zhuangWin))
		uscoresRet[ud.user.UserID] = ud.updatemoney
		fkeb.GameLog.Context().UserUpdateScore(ud.user.UserID, ud.updatemoney)
		cscoresRet[ud.user.UserID] = 0
		// 在這裡更新超端的完成金額
		if !ud.user.IsRobot() && ud != nil {
			fkeb.updateVipUserWin(ud.user, ud.updatemoney)
			// 超端的LOG記錄在updateVipUserWin裡面
		}
	}
	fkeb.boardcastZhuangWinMsg(fkeb.zhuangWin)
	uscoresRet[fkeb.zhuangID] = fkeb.zhuangWin
	cscoresRet[fkeb.zhuangID] = 0
	fkeb.GameLog.GameEnd()
	//fkeb.bettotal = 0
	fkeb.writePoolToFile()
	fkeb.action_str = fkeb.RegisterAnimation(GameEnd_Action, fkeb.endactiontime, fkeb.gameEnd)
}

//给每个玩家发送开奖结果
func (fkeb *FKEBLogic) sendUserResult(ud *fkebuserdata, correctNumberWin *int64) {
	var flowamount, betAmount, payAmount int64 //流水, 有效投注金額, 返還金額計算
	betmoney := ud.user.Score - ud.money
	win := int64(0)
	afterTaxWin := int64(0)
	totalTax := int64(0)
	for i := 0; i < len(ud.userbet); i++ {
		isZhung := false //中獎
		betAmount += ud.userbet[i].money
		for _, t := range fkeb.jiangtype {
			if ud.userbet[i].fkebtype == t {
				isZhung = true
				if !fkeb.GetUserByUserID(fkeb.zhuangID).IsRobot() {
					*correctNumberWin += ud.userbet[i].money
					//win = ud.userbet[i].money/int64(100)*int64(100-fkeb.GetServerConfig().TaxScale) + ud.userbet[i].money
					win = ud.userbet[i].money * int64(fkeb.betvalue.getOddByType(t))
				} else {
					//win = ud.userbet[i].money/int64(100)*int64(100-fkeb.GetServerConfig().TaxScale) + ud.userbet[i].money
					win = ud.userbet[i].money * int64(fkeb.betvalue.getOddByType(t))
				}

				flowamountTotal := win - ud.userbet[i].money
				flowAmountAfterTax, _ := fkeb.GetUserWinAndTax(flowamountTotal)
				flowamount += flowAmountAfterTax
				payAmount += ud.userbet[i].money + flowAmountAfterTax
				mlog.System("用户%d,fkebtype=%d，win=%d", ud.user.UserID, t, win)
				//机器人不计算奖池里面的钱
				if fkeb.zhuangID == 0 {
					if !ud.user.IsRobot() {
						fkeb.bettotal -= ud.userbet[i].money * int64(fkeb.betvalue.getOddByType(t))
					}
				} else {
					ur := fkeb.GetUserByUserID(fkeb.zhuangID)
					if ur != nil && ur.IsRobot() && ud.user.IsRobot() {
						fkeb.bettotalrobot -= win
						//fkeb.bettotal -= win
					}
					if ur != nil && !ur.IsRobot() && ud.user.IsRobot() {
						fkeb.bettotalrobot -= win
						if !ur.IsRobot() {
							fkeb.bettotal -= win
						}
					} else {
						if !ud.user.IsRobot() {
							fkeb.bettotal -= win
						}

					}
				}
				afterTaxWin, totalTax = fkeb.GetUserWinAndTax(win - ud.userbet[i].money)
				ud.money += afterTaxWin + ud.userbet[i].money
				break
			}
		}
		if win == 0 && !fkeb.GetUserByUserID(fkeb.zhuangID).IsRobot() {
			afterTaxWin, totalTax = fkeb.GetUserWinAndTax(ud.userbet[i].money)
			*correctNumberWin -= afterTaxWin
		}

		if !isZhung {
			flowamount += ud.userbet[i].money
		}
	}

	money := ud.money - ud.user.Score
	//ud.win = money + betmoney
	//ud.updatemoney = money
	//fkeb.GameLog.AddGameStep("用户%d下注赢钱%d", ud.user.UserID, ud.win)

	ud.win = money + betmoney
	ud.updatemoney = money

	fkeb.UpdateUserScoreFunc(ud.user.UserID, ud.updatemoney, totalTax)
	mlog.System("用户%d，实际赢钱：%v,下注赢钱：%v", ud.user.UserID, ud.updatemoney, ud.win)

	ud.user.FlowAmount += flowamount
	ud.user.BetAmount += betAmount
	ud.user.PayAmount += payAmount
	stepmsg := "用户%d下注赢钱%d, 抽稅%d, 流水%d, 有效投注金額%d, 返還金額%d"
	if ud.user.IsRobot() {
		stepmsg = "[機器人]" + stepmsg
	}
	fkeb.GameLog.AddGameStep(stepmsg, ud.user.UserID, ud.updatemoney, totalTax, flowamount, betAmount, payAmount)

	//这里写入赢家排名列表
	badd := false
	if ud.win > 0 {
		if len(fkeb.winuserinfo) > 0 {
			for i := 0; i < len(fkeb.winuserinfo); i++ {
				u := fkeb.getUserDataByID(fkeb.winuserinfo[i])
				if u != nil {
					if ud.win > u.win {
						if len(fkeb.winuserinfo) >= 3 {
							if i == 0 {
								fkeb.winuserinfo[2] = fkeb.winuserinfo[1]
								fkeb.winuserinfo[1] = fkeb.winuserinfo[0]
								fkeb.winuserinfo[0] = ud.user.UserID
							} else if i == 1 {
								fkeb.winuserinfo[2] = fkeb.winuserinfo[1]
								fkeb.winuserinfo[1] = ud.user.UserID
							} else {
								fkeb.winuserinfo[2] = ud.user.UserID
							}
						} else if len(fkeb.winuserinfo) == 2 {
							fkeb.winuserinfo = append(fkeb.winuserinfo, ud.user.UserID)
							if i == 0 {
								fkeb.winuserinfo[2] = fkeb.winuserinfo[1]
								fkeb.winuserinfo[1] = fkeb.winuserinfo[0]
								fkeb.winuserinfo[0] = ud.user.UserID
							} else if i == 1 {
								fkeb.winuserinfo[2] = fkeb.winuserinfo[1]
								fkeb.winuserinfo[1] = ud.user.UserID
							}
						} else {
							fkeb.winuserinfo = append(fkeb.winuserinfo, ud.user.UserID)
							fkeb.winuserinfo[1] = fkeb.winuserinfo[0]
							fkeb.winuserinfo[0] = ud.user.UserID
						}

						badd = true
						break
					}
				}
			}
		}

		if len(fkeb.winuserinfo) < 3 && !badd {
			fkeb.winuserinfo = append(fkeb.winuserinfo, ud.user.UserID)
		}
	}
}

//廣播莊家贏錢消息
func (fkeb *FKEBLogic) boardcastZhuangWinMsg(zhuangwin int64) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.FKEB_GameMessageClassID_FKEBGetZhuangWinRetID)

	pdata := new(netproto.FKEB_GameZhuangWinRet)
	pdata.ZhuangWin = proto.Int64(zhuangwin)
	pmsg.MsgData = pdata

	fkeb.BoardcastMessage(pmsg)
}

//获取玩家赢钱消息
func (fkeb *FKEBLogic) getUserWinMsg(selfwin int64, zhuangwin int64) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.FKEB_GameMessageClassID_FKEBGameResultWinID)

	pdata := new(netproto.FKEB_GameResultUserMoney)
	pdata.ZhuangWin = proto.Int64(zhuangwin)
	pdata.SelfWin = proto.Int64(selfwin)
	for i := 0; i < len(fkeb.winuserinfo); i++ {
		ud := fkeb.getUserDataByID(fkeb.winuserinfo[i])
		if ud == nil {
			continue
		}
		userwin := new(netproto.FKEB_WinUerInfo)
		userwin.UserName = proto.String(ud.user.NickName)
		userwin.Win = proto.Int64(ud.win)
		pdata.WinUserInfo = append(pdata.WinUserInfo, userwin)
	}

	pmsg.MsgData = pdata
	return pmsg
}

//游戏结束
func (fkeb *FKEBLogic) gameEnd() {
	fkeb.setState(Wait)
	fkeb.BoardcastMessage(fkeb.getStateMsg())
	if fkeb.zhuangID != 0 {
		user := fkeb.GetUserByUserID(fkeb.zhuangID)
		if user.IsOffLine() {
			fkeb.zhuangID = 0
		}
	}
	fkeb.GameEndFunc()
}

//获取开奖的麻将信息
func (fkeb *FKEBLogic) getMaJiang(key int) {
	jianglen := len(fkeb.jiangtype)
	var m []*majiang
	//fkeb.majiangresult = make([]*majiang, 0)
	//这里循环找出符合开奖的牌
	for {
		fkeb.initmajiang()
		m = make([]*majiang, 0)
		for i := 0; i < 4; i++ {
			m = append(m, fkeb.getMajiangCard())
		}

		//按从大到小进行排序
		fkeb.sortMajiang(m)
		//庄家牌不能和赢庄家的相同。
		if jianglen == 1 && fkeb.jiangtype[0] != 0 {
			if fkeb.getMajiangsType(m[0]) == fkeb.getMajiangsType(m[1]) {
				continue
			}
		} else if jianglen == 2 {
			if fkeb.getMajiangsType(m[1]) == fkeb.getMajiangsType(m[2]) {
				continue
			}
		} else if jianglen == 3 {
			if fkeb.getMajiangsType(m[2]) == fkeb.getMajiangsType(m[3]) {
				continue
			}
		}

		break
	}

	var shunret *majiang //顺的牌
	var tianret *majiang //天的牌
	var diret *majiang   //地的牌
	//将牌以庄，顺，天，地的顺序放入结果列表中
	if fkeb.jiangtype[0] == 0 {
		//通杀
		fkeb.majiangresult = append(fkeb.majiangresult, m[0])
		m = append(m[:0], m[1:]...)
	} else if jianglen >= 3 {
		//通赔
		fkeb.majiangresult = append(fkeb.majiangresult, m[3])
		m = append(m[:3], m[4:]...)
	} else if jianglen == 2 {
		//开两个的情况
		fkeb.majiangresult = append(fkeb.majiangresult, m[2])
		m = append(m[:2], m[3:]...)
		//这里随机一次，以免出现牌看起来是从大到小排序的
		i := rand.Int31n(2)
		if fkeb.jiangtype[0] == SHUN {
			shunret = m[i]
			m = append(m[:i], m[i+1:]...)
		} else if fkeb.jiangtype[0] == TIAN {
			tianret = m[i]
			m = append(m[:i], m[i+1:]...)
		} else if fkeb.jiangtype[0] == DI {
			diret = m[i]
			m = append(m[:i], m[i+1:]...)
		}

		if fkeb.jiangtype[1] == SHUN {
			shunret = m[0]
			m = append(m[:0], m[1:]...)
		} else if fkeb.jiangtype[1] == TIAN {
			tianret = m[0]
			m = append(m[:0], m[1:]...)
		} else if fkeb.jiangtype[1] == DI {
			diret = m[0]
			m = append(m[:0], m[1:]...)
		}
	} else {
		//其他开一个奖的情况
		fkeb.majiangresult = append(fkeb.majiangresult, m[1])
		m = append(m[:1], m[2:]...)
		if fkeb.jiangtype[0] == SHUN {
			shunret = m[0]
			m = append(m[:0], m[1:]...)
		} else if fkeb.jiangtype[0] == TIAN {
			tianret = m[0]
			m = append(m[:0], m[1:]...)
		} else if fkeb.jiangtype[0] == DI {
			diret = m[0]
			m = append(m[:0], m[1:]...)
		}
	}

	if shunret != nil {
		fkeb.majiangresult = append(fkeb.majiangresult, shunret)
	} else {
		//这里随机一次，以免出现牌看起来是从大到小排序的
		i := rand.Int31n(int32(len(m)))
		fkeb.majiangresult = append(fkeb.majiangresult, m[i])
		m = append(m[:i], m[i+1:]...)
	}

	if tianret != nil {
		fkeb.majiangresult = append(fkeb.majiangresult, tianret)
	} else {
		//这里随机一次，以免出现牌看起来是从大到小排序的
		i := rand.Int31n(int32(len(m)))
		fkeb.majiangresult = append(fkeb.majiangresult, m[i])
		m = append(m[:i], m[i+1:]...)
	}

	if diret != nil {
		fkeb.majiangresult = append(fkeb.majiangresult, diret)
	} else {
		//这里随机一次，以免出现牌看起来是从大到小排序的
		if len(m) > 0 { // 如果len(m)程式會當掉，防呆
			i := rand.Int31n(int32(len(m)))
			fkeb.majiangresult = append(fkeb.majiangresult, m[i])
			m = append(m[:i], m[i+1:]...)
		}
	}

	// 如果有指定開獎，以指定開獎為主
	if key <= 0 { // 如果沒有指定開獎，就以超端為主
		// 在這裡進行超端控制
		findSmallAssignAmount := int64(900000000)
		currentUser := int32(0)
		for i := 0; i < len(fkeb.userdata); i++ {
			if !fkeb.userdata[i].user.IsRobot() {
				// 二人同時被監控時，以勝率最低的被優先監控
				vipuser := fkeb.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(fkeb.userdata[i].user.UserID)
				if vipuser.GetAssignAmount() < findSmallAssignAmount {
					findSmallAssignAmount = vipuser.GetAssignAmount()
					currentUser = fkeb.userdata[i].user.UserID
				}
			}
		}
		if currentUser != int32(0) {
			fkeb.bowJiangType(fkeb.GetUserByUserID(currentUser))
		}
	}

	// fkeb.testbowJiangType() // 模擬超端控制

	for i := 0; i < len(fkeb.majiangresult); i++ {
		mlog.System("开出的麻将结果为%v", fkeb.majiangresult[i])
		fkeb.GameLog.AddGameStep("開獎區域[%v] 开出的麻将结果为%v", i, fkeb.majiangresult[i])
	}
}

//获取牌
func (fkeb *FKEBLogic) getMajiangCard() *majiang {
	m := new(majiang)
	m.majiang1 = fkeb.majiang[0]
	m.majiang2 = fkeb.majiang[1]
	fkeb.majiang = append(fkeb.majiang[:0], fkeb.majiang[2:]...)
	return m
}

func (fkeb *FKEBLogic) sortMajiang(m []*majiang) {
	for i := 0; i < len(m); i++ {
		for j := 0; j+1 < len(m)-i; j++ {
			if fkeb.getMajiangsType(m[j+1]) > fkeb.getMajiangsType(m[j]) {
				m[j], m[j+1] = m[j+1], m[j]
			}
		}
	}
}

func (fkeb *FKEBLogic) getMajiangsType(m *majiang) int32 {
	// 參考 https://example.invalid/9vo3dSag_xI4khGkpoWK1HF6hhy/baike/c0%3Dbaike80%2C5%2C5%2C80%2C26/sign=84f3f797b6b7d0a26fc40ccfaa861d6c/50da81cb39dbb6fd7ba17e030d24ab18972b3718.jpg
	// 二八槓接近有快六十種組合，對子10種，28為1種，九點5種，八點4種，七點5種，六點4種，五點5種，四點4種，三點3種，二點4種，一點5種，0點3種
	// 為了要預防錯誤，所以牌型分數設定成以二十分為一組，並間距加大，28的分數為300分，普通牌最多到0~190幾分，對子為400分以上
	// 因為 一張牌 單牌 分數 可能有0分，也有可能有10幾分，有必要把組距調成20分以上，去防止分數重複

	//对子
	if m.majiang1 == m.majiang2 {
		return m.majiang1 + 400
	}

	//天杠
	if (m.majiang1 == 2 && m.majiang2 == 8) || (m.majiang1 == 8 && m.majiang2 == 2) { // 這裡進行 2 8 和 8 2 的牌型修正
		return 300
	}

	//普通点
	maxMajiang := m.majiang1
	if maxMajiang < m.majiang2 {
		maxMajiang = m.majiang2
	}

	return ((m.majiang1+m.majiang2)%10)*20 + maxMajiang
}

//获取所有的开奖结果
func (fkeb *FKEBLogic) getAllKaiJiangType(betmoney int64) []int32 {
	alltype := fkeb.betvalue.getAllConfigType()

	for _, v := range fkeb.bettypetotal {
		money := int64(fkeb.betvalue.getOddByType(v.fkebtype)) * v.money

		if money > betmoney {
			mlog.Debug("类型：%v，需要钱为：%v", v.fkebtype, money)
			for i := 0; i < len(alltype); i++ {
				if alltype[i] == v.fkebtype {
					alltype = append(alltype[:i], alltype[i+1:]...)
					break
				}
			}
		}
	}

	return alltype
}

func (fkeb *FKEBLogic) getJiangType(alltype []int32, betmoney int64) {
	//开三个需要的钱
	totalmoney := int64(0)
	//开顺门和天门需要的钱
	money12 := int64(0)
	//开地门和天门需要的钱
	money23 := int64(0)
	//开顺门和地门需要的钱
	money13 := int64(0)
	for _, v := range fkeb.bettypetotal {
		money := int64(fkeb.betvalue.getOddByType(v.fkebtype)) * v.money
		totalmoney += money
		if v.fkebtype == SHUN {
			money12 += money
			money13 += money
		} else if v.fkebtype == TIAN {
			money23 += money
			money13 += money
		} else if v.fkebtype == DI {
			money12 += money
			money23 += money
		}
	}

	count := int32(1)
	if totalmoney < betmoney {
		count = 3
	} else if money12 <= betmoney || money23 <= betmoney || money13 <= betmoney {
		count = 2
	}

	index := 0
	for i := 0; i < len(fkeb.scoresection); i++ {
		if fkeb.syspool < fkeb.scoresection[i] {
			index = i
		}
	}

	//总权重
	totalweight := fkeb.weight[index].weight[0] + fkeb.weight[index].weight[1]
	c := 1
	if count == 2 {
		w := rand.Int31n(totalweight)
		if w > fkeb.weight[index].weight[0] {
			//开两个奖的情况
			c = 2
		}
	} else if count == 3 {
		totalweight += fkeb.weight[index].weight[2]
		w := rand.Int31n(totalweight)

		if w > fkeb.weight[index].weight[0] {
			//开两个奖的情况
			c = 2
		} else if w > fkeb.weight[index].weight[0]+fkeb.weight[index].weight[1] {
			//开三个奖的情况
			c = 3
		}
	}

	if c == 1 {
		mlog.Debug("开一个的奖")
		fkeb.jiangtype = append(fkeb.jiangtype, fkeb.betvalue.getJiangType(alltype))
	} else if c == 2 {
		mlog.Debug("开两个的奖")
		n := int32(0)
		if money12 <= betmoney {
			n++
		}
		if money23 <= betmoney {
			n++
		}
		if money13 <= betmoney {
			n++
		}

		if n == 0 {
			mlog.Fatal("不能开出两个奖")
		}
		tmp := rand.Int31n(n) + 1
		j := int32(0)
		if money12 <= betmoney {
			j++
			if tmp == j {
				fkeb.jiangtype = append(fkeb.jiangtype, 1)
				fkeb.jiangtype = append(fkeb.jiangtype, 2)
			}
		}

		if money23 <= betmoney {
			j++
			if tmp == j {
				fkeb.jiangtype = append(fkeb.jiangtype, 2)
				fkeb.jiangtype = append(fkeb.jiangtype, 3)
			}
		}

		if money13 <= betmoney {
			j++
			if tmp == j {
				fkeb.jiangtype = append(fkeb.jiangtype, 1)
				fkeb.jiangtype = append(fkeb.jiangtype, 3)
			}
		}
	} else {
		fkeb.jiangtype = append(fkeb.jiangtype, 1)
		fkeb.jiangtype = append(fkeb.jiangtype, 2)
		fkeb.jiangtype = append(fkeb.jiangtype, 3)
	}

	if len(fkeb.jiangtype) <= 0 {
		mlog.Fatal("2开奖出问题了，未开出奖")
	}
}

//获取赔偿最小的开奖结果
func (fkeb *FKEBLogic) getMinKaiJiangType() int32 {
	mintype := int32(0)
	minmoney := int64(0)
	for i := int32(1); i <= 3; i++ {
		money := int64(fkeb.betvalue.getOddByType(i)) * fkeb.getBetMoneyByType(i)
		if i == 1 {
			minmoney = money
			mintype = i
			continue
		}

		if money < minmoney {
			minmoney = money
			mintype = i
		}
	}

	return mintype
}

//设置状态
func (fkeb *FKEBLogic) setState(state int32) {
	fkeb.gamestate = state
}

//用户消息主处理函数
func (fkeb *FKEBLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	if !fkeb.CheckGameUserMessage(user) {
		fkeb.KickUserFunc(user.UserID, "非法操作")
		return
	}
	//不是给游戏逻辑服务器的消息
	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	fkeb.CheckUserAnimationCommited(user, gmsg, fkeb.action_str, nil)

	if gmsg.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBGetShangZhuangListID) {
		fkeb.SendMessageToUser(user, fkeb.getZhuangListMsg())
		return
	}

	if fkeb.gamestate != Bet {
		return
	}

	if gmsg.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBBetID) {
		fkeb.userBet(user, gmsg.MsgData.([]byte))
		return
	}

	if gmsg.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBShangZhuangID) {
		fkeb.shangZhuang(user, gmsg.MsgData.([]byte))
		return
	}

	if gmsg.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBBetClearID) && fkeb.gamestate == Bet {
		fkeb.clearUserBet(user)
		return
	}
}

//获取配置表
func (fkeb *FKEBLogic) getfkebconfig(cfgs *config.Group) {
	fkeb.scoresection = make([]int64, 0)
	fkeb.betvalue = newBetconfig(cfgs)
	scfgs := cfgs.GetConfig("gameserver")
	fkeb.bettime = int32(scfgs.GetConfigInt("gameserver", "bettime"))
	fkeb.kaijiangactiontime = int32(scfgs.GetConfigInt("gameserver", "kaijiangactiontime"))
	fkeb.endactiontime = int32(scfgs.GetConfigInt("gameserver", "endactiontime"))
	fkeb.readytime = int32(scfgs.GetConfigInt("gameserver", "readywaittime"))
	fkeb.zhuanglimit = int64(scfgs.GetConfigInt("gameserver", "zhuanglimit"))
	fkeb.maxturn = int32(scfgs.GetConfigInt("gameserver", "zhuangturn"))
	fkeb.serverid = strconv.Itoa(int(fkeb.GetServerConfig().ServerID))
	fkeb.maxclearcount = int32(scfgs.GetConfigInt("gameserver", "maxclearcount"))
	fkeb.zhushui = int32(scfgs.GetConfigInt("gameserver", "zhushui"))
	fkeb.maxrobotzhuang = int32(scfgs.GetConfigInt("gameserver", "maxrobotzhuang"))
	fkeb.robotxiazhuang = int32(scfgs.GetConfigInt("gameserver", "robotxiazhuang"))
	fkeb.robotzhuang = int32(scfgs.GetConfigInt("gameserver", "robotzhuang"))
	fkeb.winline = int64(scfgs.GetConfigInt("gameserver", "winline"))
	fkeb.naturl = int32(scfgs.GetConfigInt("jiangtyperate", "naturl"))
	fkeb.tongsha = int32(scfgs.GetConfigInt("jiangtyperate", "tongsha"))
	fkeb.tongpei = int32(scfgs.GetConfigInt("jiangtyperate", "tongpei"))
	fkeb.win1loss2 = int32(scfgs.GetConfigInt("jiangtyperate", "win1loss2"))
	fkeb.win2loss1 = int32(scfgs.GetConfigInt("jiangtyperate", "win2loss1"))
	fkeb.zrusertongpei = int32(scfgs.GetConfigInt("gameserver", "zrusertongpei"))
	betcountstr := scfgs.GetConfigValue("gameserver", "betcount")
	fkeb.betlimitmoney = int64(scfgs.GetConfigInt("gameserver", "betlimitmoney"))

	betarr := strings.Split(betcountstr, ",")
	for _, str := range betarr {
		count, _ := strconv.ParseInt(str, 10, 32)
		fkeb.betcount = append(fkeb.betcount, int32(count))
	}

	sectionstr := scfgs.GetConfigValue("gameserver", "scoresection")
	sectionarr := strings.Split(sectionstr, ",")

	for _, str := range sectionarr {
		if str == "0" {
			continue
		}
		s, _ := strconv.ParseInt(str, 10, 64)
		fkeb.scoresection = append(fkeb.scoresection, s)
		weightconfig := new(fkebweight)
		for i := 1; i <= 3; i++ {
			weightstr := scfgs.GetConfigValue(str, strconv.Itoa(i))
			w, _ := strconv.ParseInt(weightstr, 10, 32)

			weightconfig.count = append(weightconfig.count, int32(i))
			weightconfig.weight = append(weightconfig.weight, int32(w))
		}

		fkeb.weight = append(fkeb.weight, weightconfig)
	}

	fkeb.ctrlbetlimit = int64(scfgs.GetConfigInt("supertool", "ctrlbetlimit"))
	if fkeb.ctrlbetlimit <= 0 {
		fkeb.ctrlbetlimit = CtrlBetLimitDefault
	}
}

//清除玩家下注信息
func (fkeb *FKEBLogic) clearUserBet(user *gamelogic.GameUser) {
	ud := fkeb.getUserDataByID(user.UserID)
	if ud == nil {
		return
	}

	mlog.Debug("ud.clearcount=%v", ud.clearcount)
	if ud.clearcount >= fkeb.maxclearcount {
		fkeb.SendMessageToUser(user, fkeb.getClearRetMsg(ud.money, -1))
		return
	}

	ud.clearcount++
	for _, ub := range ud.getUserBet() {
		fkeb.updateBet(ub.fkebtype, -ub.money)
	}

	fkeb.GameLog.AddGameStep("用户%d点击清除", user.UserID)
	ud.clearData()
	//fkeb.bettotal -= ud.user.Score - ud.money
	//fkeb.delUserData(user)
	fkeb.SendMessageToUser(user, fkeb.getClearRetMsg(user.Score, 0))
	fkeb.BoardcastMessage(fkeb.getBetTotalMsg(user.UserID))
}

//获取清除结果消息
func (fkeb *FKEBLogic) getClearRetMsg(money int64, ret int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.FKEB_GameMessageClassID_FKEBBetClearRetID)

	pdata := new(netproto.FKEB_BetClearRet)
	pdata.Money = proto.Int64(money)
	pdata.Ret = proto.Int32(ret)
	pmsg.MsgData = pdata

	return pmsg
}

//用户下注处理函数
func (fkeb *FKEBLogic) userBet(user *gamelogic.GameUser, data []byte) {
	ud := fkeb.getUserDataByID(user.UserID)
	if ud == nil {
		ud = fkeb.addUserData(user)
	}

	ust := &netproto.FKEB_BetArr{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Debug("解析用户%v下注错误%v", user, err)
		return
	}

	//庄家不能下注
	if fkeb.zhuangID == user.UserID {
		return
	}

	money := int32(0)
	for _, b := range ust.BetArr {
		*b.BetType = int32(math.Abs(float64(*b.BetType)))
		*b.Count = int32(math.Abs(float64(*b.Count)))

		//玩家身上金額小於下注限制
		if user.Score < fkeb.betlimitmoney*100 {
			for i, u := range fkeb.userdata {
				if u.user.UserID == ud.user.UserID {
					fkeb.userdata = append(fkeb.userdata[:i], fkeb.userdata[i+1:]...)
				}
			}
			fkeb.SendMessageToUser(user, fkeb.getBetRetMsg(ud, -6, *b.BetType, ust.BetArr))
			return
		}

		tmpmon := fkeb.getMoneyByType(*b.BetType) * (*b.Count)
		if (int64(tmpmon) + fkeb.getBetMoneyByType(*b.BetType)) > fkeb.betvalue.getMaxMoneyByType(*b.BetType) {
			fkeb.SendMessageToUser(user, fkeb.getBetRetMsg(ud, -4, *b.BetType, ust.BetArr))
			return
		}

		if fkeb.zhuangID != 0 {
			ur := fkeb.GetUserByUserID(fkeb.zhuangID)
			allmoney := (int64(tmpmon) + fkeb.getBetMoneyByAllType())
			if allmoney > (ur.Score) {
				fkeb.SendMessageToUser(user, fkeb.getBetRetMsg(ud, -4, *b.BetType, ust.BetArr))
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
			betmon := fkeb.getMoneyByType(*b.BetType) * (*b.Count)
			mlog.Debug("收到下注消息%d", betmon)
			ud.setUserBet(*b.BetType, betmon)
			fkeb.updateBet(*b.BetType, int64(betmon))
			if user.IsRobot() {
				fkeb.updaterobotBet(*b.BetType, int64(betmon))
			}
		}

		fkeb.GameLog.AddGameStep("用户%d下注%v", user.UserID, ust.BetArr)
		fkeb.BoardcastMessage(fkeb.getBetTotalMsg(user.UserID))
	}

	fkeb.SendMessageToUser(user, fkeb.getBetRetMsg(ud, ret, 0, ust.BetArr))

}

//获取此类型的下注总额
func (fkeb *FKEBLogic) getBetMoneyByType(t int32) int64 {
	for i := 0; i < len(fkeb.bettypetotal); i++ {
		if fkeb.bettypetotal[i].fkebtype == t {
			return fkeb.bettypetotal[i].money
		}
	}

	return 0
}

//全部類型的下注總額
func (fkeb *FKEBLogic) getBetMoneyByAllType() int64 {
	sum := int64(0)
	for i := 0; i < len(fkeb.bettypetotal); i++ {
		sum += fkeb.bettypetotal[i].money
	}

	return sum
}

//获取此类型的机器人下注总额
func (fkeb *FKEBLogic) getRobotBetMoneyByType(t int32) int64 {
	for i := 0; i < len(fkeb.bettyperobot); i++ {
		if fkeb.bettyperobot[i].fkebtype == t {
			return fkeb.bettyperobot[i].money
		}
	}

	return 0
}

//获取此类型一次下注额
func (fkeb *FKEBLogic) getMoneyByType(t int32) int32 {
	return fkeb.betvalue.getMoneyByType(t)
}

//更新下注钱
func (fkeb *FKEBLogic) updateBet(t int32, money int64) {
	fkeb.bettotal += int64(money)
	mlog.System("更新下注钱 BetType=%v,money=%v,fkeb.bettotal=%v", t, money, fkeb.bettotal)
	for _, v := range fkeb.bettypetotal {
		if v.fkebtype == t {
			v.money += int64(money)
			if v.money < 0 {
				mlog.Debug("钱小于0了！")
				v.money = 0
			}
			return
		}
	}

	b := new(bettype)
	b.fkebtype = t
	b.money = int64(money)
	if b.money < 0 {
		mlog.Debug("钱小于0了！")
		b.money = 0
	}
	fkeb.bettypetotal = append(fkeb.bettypetotal, b)
}

//更新机器人下注额
func (fkeb *FKEBLogic) updaterobotBet(t int32, money int64) {
	fkeb.bettotalrobot += int64(money)
	for _, v := range fkeb.bettyperobot {
		if v.fkebtype == t {
			v.money += money
			if v.money < 0 {
				v.money = 0
			}
			mlog.Debug("下注类型%d，下注的总钱数为%d", t, v.money)
			return
		}
	}

	b := new(bettype)
	b.fkebtype = t
	b.money = int64(money)
	if b.money < 0 {
		b.money = 0
	}

	fkeb.bettyperobot = append(fkeb.bettyperobot, b)
}

//移除机器人下注额
func (fkeb *FKEBLogic) removeRobotBet() {
	for _, v := range fkeb.bettyperobot {
		fkeb.updateBet(v.fkebtype, -v.money)
	}
}

//获取状态消息
func (fkeb *FKEBLogic) getStateMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.FKEB_GameMessageClassID_FKEBStateID)

	pdata := new(netproto.FKEB_Sate)
	pdata.State = proto.Int32(fkeb.gamestate)
	pdata.Time = proto.Int32(fkeb.getBetTime())
	pmsg.MsgData = pdata
	return pmsg
}

func (fkeb *FKEBLogic) getBetTime() int32 {
	time := int32(0)
	if fkeb.opttimer != nil {
		time = fkeb.opttimer.WaitSeconds - fkeb.opttimer.ElapsedSeconds()
	}
	return time
}

//获取用户数据通过用户ID
func (fkeb *FKEBLogic) getUserDataByID(userid int32) *fkebuserdata {
	for _, ud := range fkeb.userdata {
		if ud.user.UserID == userid {
			return ud
		}
	}

	return nil
}

//添加一个用户数据
func (fkeb *FKEBLogic) addUserData(user *gamelogic.GameUser) *fkebuserdata {
	ud := newLongHuDouUserData(user)
	fkeb.userdata = append(fkeb.userdata, ud)
	return ud
}

//删除一个用户数据
func (fkeb *FKEBLogic) delUserData(user *gamelogic.GameUser) {
	for i := 0; i < len(fkeb.userdata); i++ {
		if fkeb.userdata[i].user.UserID == user.UserID {
			fkeb.userdata = append(fkeb.userdata[:i], fkeb.userdata[i+1:]...)
			break
		}
	}
}

//获取庄家信息消息
func (fkeb *FKEBLogic) getZhuangInfoMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.FKEB_GameMessageClassID_FKEBZhuangInfoID)

	pdata := new(netproto.FKEB_ZhuangInfo)
	pdata.UserID = proto.Int32(0)
	pdata.LastTurn = proto.Int32(0)
	ur := fkeb.GetUserByUserID(fkeb.zhuangID)
	if ur != nil {
		pdata.UserInfo = append(pdata.UserInfo, ur.GetProtoBufGameUser())
		pdata.UserID = proto.Int32(ur.UserID)
		pdata.LastTurn = proto.Int32(fkeb.lastturn)
	}
	pmsg.MsgData = pdata
	return pmsg
}

//获取上庄列表消息
func (fkeb *FKEBLogic) getZhuangListMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.FKEB_GameMessageClassID_FKEBShangZhuangListInfoID)

	pdata := new(netproto.FKEB_ZhuangList)
	for i := 0; i < len(fkeb.zhuanglist); i++ {
		ur := fkeb.GetUserByUserID(fkeb.zhuanglist[i])
		pdata.UserInfo = append(pdata.UserInfo, ur.GetProtoBufGameUser())
	}
	pmsg.MsgData = pdata
	return pmsg
}

//检测用户是否能上庄
func (fkeb *FKEBLogic) checkUserShangZhuang(userid int32) bool {
	ur := fkeb.GetUserByUserID(userid)
	if ur != nil && ur.Score >= fkeb.zhuanglimit {
		return true
	}

	return false
}

//换庄
func (fkeb *FKEBLogic) huangZhuang() {
	if fkeb.zhuangID == 0 {
		//庄是系统直接换庄
		fkeb.setNextZhuang()
		return
	}
	if fkeb.lastturn <= 0 {
		//已经到最后一局了
		fkeb.setNextZhuang()
		return
	}
	ur := fkeb.GetUserByUserID(fkeb.zhuangID)
	if ur != nil && ur.Score < fkeb.zhuanglimit {
		fkeb.setNextZhuang()
		return
	}

	if ur == nil {
		fkeb.setNextZhuang()
		return
	}

	fkeb.lastturn--
	mlog.Debug("1015发送庄家剩余局数%d", fkeb.lastturn)
}

//检测下一个用户是否能上庄
func (fkeb *FKEBLogic) checkNetUserCanShangZhuang() {
	if len(fkeb.zhuanglist) > 0 {
		for {
			user := fkeb.GetUserByUserID(fkeb.zhuanglist[0])
			if user != nil && user.Score >= int64(fkeb.zhuanglimit) {
				break
			}

			fkeb.delFromZhuangList(fkeb.zhuanglist[0])
			if len(fkeb.zhuanglist) <= 0 {
				break
			}
		}
	}
}

//设置下一个上庄的玩家
func (fkeb *FKEBLogic) setNextZhuang() {
	fkeb.checkNetUserCanShangZhuang()
	if len(fkeb.zhuanglist) > 0 {
		fkeb.zhuangID = fkeb.zhuanglist[0]
		fkeb.lastturn = fkeb.maxturn
		fkeb.lastturn--
		fkeb.delFromZhuangList(fkeb.zhuangID)
		user := fkeb.GetUserByUserID(fkeb.zhuangID)
		fkeb.startmoney = user.Score

		if !user.IsRobot() {
			fkeb.realuserzhuangflow = 0
			fkeb.ruztotalwinlose = 0
		}

	} else {
		fkeb.zhuangID = 0
		fkeb.startmoney = fkeb.startsyspool
	}
}

//获取下注消息
func (fkeb *FKEBLogic) getBetRetMsg(ud *fkebuserdata, ret int32, t int32, betinfo []*netproto.FKEB_Bet) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.FKEB_GameMessageClassID_FKEBBetRetID)

	pdata := new(netproto.FKEB_BetRet)
	pdata.Money = proto.Int64(ud.money)
	pdata.Ret = proto.Int32(ret)
	pdata.Type = proto.Int32(t)
	pdata.BetArr = betinfo
	pmsg.MsgData = pdata
	return pmsg
}

//获取下注总额消息
func (fkeb *FKEBLogic) getBetTotalMsg(userid int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.FKEB_GameMessageClassID_FKEBBroadCastBetID)

	pdata := new(netproto.FKEB_BroadCastBetArr)
	pdata.UserID = proto.Int32(userid)
	for _, v := range fkeb.betvalue.alltype {
		betinfo := new(netproto.FKEB_BroadCastBet)
		betinfo.Type = proto.Int32(v)
		betinfo.Money = proto.Int64(fkeb.getBetMoneyByType(v))
		pdata.BroadCastBetArr = append(pdata.BroadCastBetArr, betinfo)
	}

	pmsg.MsgData = pdata

	return pmsg
}

//获取上庄列表索引通过用户ID
func (fkeb *FKEBLogic) getZhuangListIndexByUserId(userid int32) int32 {
	for i := 0; i < len(fkeb.zhuanglist); i++ {
		if fkeb.zhuanglist[i] == userid {
			return int32(i)
		}
	}

	return -1
}

//添加一个用户到上庄列表
func (fkeb *FKEBLogic) addToZhuangList(userid int32) {
	checkOK := fkeb.checkUserShangZhuang(userid)
	if checkOK {
		fkeb.zhuanglist = append(fkeb.zhuanglist, userid)
	}
}

//删除一个用户从上庄列表
func (fkeb *FKEBLogic) delFromZhuangList(userid int32) {
	for i := 0; i < len(fkeb.zhuanglist); i++ {
		if fkeb.zhuanglist[i] == userid {
			fkeb.zhuanglist = append(fkeb.zhuanglist[:i], fkeb.zhuanglist[i+1:]...)
			break
		}
	}
}

//获取上庄结果消息
func (fkeb *FKEBLogic) getShangzhuangRetMsg(ret int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.FKEB_GameMessageClassID_FKEBShangZhuangRetID)

	pdata := new(netproto.FKEB_ShangzhuangRet)
	pdata.Ret = proto.Int32(ret)
	pmsg.MsgData = pdata
	return pmsg
}

//处理用户上庄请求
func (fkeb *FKEBLogic) shangZhuang(user *gamelogic.GameUser, data []byte) {
	ust := &netproto.FKEB_ShangZhuang{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Debug("解析用户%v上庄错误%v", user, err)
		return
	}

	if fkeb.getZhuangListIndexByUserId(user.UserID) != -1 {
		fkeb.GameLog.AddGameStep("用户%d下庄", user.UserID)
		fkeb.delFromZhuangList(user.UserID)
		fkeb.SendMessageToUser(user, fkeb.getShangzhuangRetMsg(1))
		return
	}

	if user.Score < fkeb.zhuanglimit {
		fkeb.SendMessageToUser(user, fkeb.getShangzhuangRetMsg(-1))
		return
	}

	if *ust.IsShangZhuang {
		fkeb.GameLog.AddGameStep("用户%d上庄", user.UserID)
		fkeb.addToZhuangList(user.UserID)
		fkeb.SendMessageToUser(user, fkeb.getShangzhuangRetMsg(0))
	} else {
		fkeb.GameLog.AddGameStep("用户%d下庄", user.UserID)
		fkeb.delFromZhuangList(user.UserID)
		fkeb.SendMessageToUser(user, fkeb.getShangzhuangRetMsg(1))
		fkeb.lastturn = -1
	}
}

//获取上一次水池的钱
func (fkeb *FKEBLogic) getLastPool() {
	if conf.TestMode { // 這裡是要方便進行單元測試
		return
	}
	var folderName = "lastpool"
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err = os.MkdirAll(folderName, os.ModePerm)
	}

	var tablenostr = strconv.Itoa(int(fkeb.TableNo()))
	var filename = fkeb.serverid + "_" + tablenostr + ".conf"
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
		fkeb.serverfile = f
		defer fkeb.serverfile.Close()
		fkeb.readPoolFromFile()
	}
}

//读取水池文件
func (fkeb *FKEBLogic) readPoolFromFile() {
	buff := make([]byte, 1024)
	n, _ := fkeb.serverfile.Read(buff)
	if n > 0 {
		poolstr := string(buff)[0:n]
		//pool, err := strconv.ParseInt(poolstr, 10, 32)
		pool, _ := strconv.Atoi(poolstr)
		fkeb.syspool = int64(pool)
	}
}

//将水池的钱写入文件
func (fkeb *FKEBLogic) writePoolToFile() {
	if conf.TestMode { // 這裡是要方便進行單元測試
		return
	}
	buff := strconv.Itoa(int(fkeb.syspool))
	//bufio.ReadWriter.Peek()
	tablenostr := strconv.Itoa(int(fkeb.TableNo()))
	filename := "./lastpool/" + fkeb.serverid + "_" + tablenostr + ".conf"
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

//添加系统水池的钱
func (fkeb *FKEBLogic) addSysPool(value int64) {
	fkeb.syspool += value
}

//减去系统水池的钱
func (fkeb *FKEBLogic) subtractSysPool(value int64) {
	fkeb.syspool -= value
}

//改变水池的钱
func (fkeb *FKEBLogic) changeSysPool(value int64) {
	fkeb.startsyspool += value - fkeb.syspool
	fkeb.syspool = value
}

//改变水池注水率
func (fkeb *FKEBLogic) changeZhuShui(value int32) {
	if value >= 0 && value <= 100 {
		fkeb.zhushui = value
	}
}

func (fkeb *FKEBLogic) checkSysPool() {
	mlog.Printf("桌子号：%v，系统水池：%v，注水概率：%v", fkeb.TableNo(), fkeb.syspool, fkeb.zhushui)
}

//占位用防止崩溃
func (fkeb *FKEBLogic) ExchangeMaxCard(userid int32) {

}

func (fkeb *FKEBLogic) OnRecvHTTPMessage(req *http.Request) proto.Message {
	action := req.PostFormValue("action")
	if action == "getserverinfo" {
		return fkeb.getTableListMsg()
	} else if action == "gettableinfo" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		return fkeb.getTableDetailMsg(int32(tid))
	} else if action == "modresult" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		ret, _ := strconv.ParseInt(req.PostFormValue("result"), 10, 32)
		return fkeb.setNextJiangType(int32(tid), int32(ret))
	} else if action == "modpoolamount" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		pool, _ := strconv.ParseInt(req.PostFormValue("poolamount"), 10, 64)
		return fkeb.setSysPool(int32(tid), pool)
	} else if action == "modpopwaterrate" {
		tid, _ := strconv.ParseInt(req.PostFormValue("tableno"), 10, 32)
		zhushui, _ := strconv.ParseInt(req.PostFormValue("popwaterrate"), 10, 32)
		return fkeb.setZhuShui(int32(tid), int32(zhushui))
	}

	return nil
}

//获取桌子列表数据
func (fkeb *FKEBLogic) getTableListMsg() proto.Message {
	pdata := new(netproto.Http_FKEB_GameInfo)
	for _, table := range fkeb.GetTable().Room.Tables {
		pTable := new(netproto.Http_FKEB_TableList)
		logic := table.GameLogic.(*FKEBLogic)
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
		pTable.SumEarnAmount = proto.Int64(logic.systotalwin)
		pTable.LastEarnAmount = proto.Int32(int32(logic.syslastwin))

		pdata.TableList = append(pdata.TableList, pTable)
	}

	return pdata
}

//获取桌子详细数据
func (fkeb *FKEBLogic) getTableDetailMsg(tableid int32) proto.Message {
	t := fkeb.GetTable().Room.GetTableByNo(tableid)
	pTable := new(netproto.Http_FKEB_TableDetail)
	if t != nil {
		logic := t.GameLogic.(*FKEBLogic)
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
		pTable.OpenResultZhuang = proto.Int64(logic.getRealUserWinByType(ZHUANG))
		pTable.OpenResultShun = proto.Int64(logic.getRealUserWinByType(SHUN))
		pTable.OpenResultTian = proto.Int64(logic.getRealUserWinByType(TIAN))
		pTable.OpenResultDi = proto.Int64(logic.getRealUserWinByType(DI))
		pTable.OpenResultShunTaiDi = proto.Int64(logic.getRealUserWinByType(DI) +
			logic.getRealUserWinByType(SHUN) + logic.getRealUserWinByType(TIAN))
		pTable.OpenResultShunTai = proto.Int64(logic.getRealUserWinByType(SHUN) + logic.getRealUserWinByType(TIAN))
		pTable.OpenResultShunDi = proto.Int64(logic.getRealUserWinByType(DI) + logic.getRealUserWinByType(SHUN))
		pTable.OpenResultTianDi = proto.Int64(logic.getRealUserWinByType(DI) + logic.getRealUserWinByType(TIAN))
		pTable.BetTime = proto.Int32(fkeb.getBetTime())
		if logic.gamestate == KaiJiangAction {
			pTable.CurentOpenResult = proto.String(logic.getResultString())
			win := int64(0)
			for i := 0; i < len(logic.jiangtype); i++ {
				win += logic.getRealUserWinByType(logic.jiangtype[i])
			}
			pTable.CurentOpenResultEarn = proto.Int64(win)
		}

		for _, ud := range logic.userdata {
			userinfo := new(netproto.Http_FKEB_UserInfo)
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
			userinfo.BetInfoShun = proto.Int64(ud.getUserBetByType(SHUN))
			userinfo.BetInfoTian = proto.Int64(ud.getUserBetByType(TIAN))
			userinfo.BetInfoDi = proto.Int64(ud.getUserBetByType(DI))
			pTable.UserInfo = append(pTable.UserInfo, userinfo)
		}
	}

	return pTable
}

func (fkeb *FKEBLogic) getResultString() string {
	var res string
	for i := 0; i < len(fkeb.jiangtype); i++ {
		res += resultstring[fkeb.jiangtype[i]]
	}
	return res
}

func (fkeb *FKEBLogic) getRetMsg(ret bool, msg string) proto.Message {
	pdata := new(netproto.Http_Result)
	pdata.Result = proto.Bool(ret)
	pdata.Message = proto.String(msg)
	return pdata
}

//设置下次开奖结果
func (fkeb *FKEBLogic) setNextJiangType(tableid int32, jiang int32) proto.Message {
	t := fkeb.GetTable().Room.GetTableByNo(tableid)

	if t != nil {
		logic := t.GameLogic.(*FKEBLogic)
		if jiang <= 3 {
			logic.nextjiangtype = append(logic.nextjiangtype, jiang)
		} else if jiang == 4 {
			//通赔
			logic.nextjiangtype = append(logic.nextjiangtype, 1)
			logic.nextjiangtype = append(logic.nextjiangtype, 2)
			logic.nextjiangtype = append(logic.nextjiangtype, 3)
		} else if jiang == 5 {
			logic.nextjiangtype = append(logic.nextjiangtype, 1)
			logic.nextjiangtype = append(logic.nextjiangtype, 2)
		} else if jiang == 6 {
			logic.nextjiangtype = append(logic.nextjiangtype, 1)
			logic.nextjiangtype = append(logic.nextjiangtype, 3)
		} else if jiang == 7 {
			logic.nextjiangtype = append(logic.nextjiangtype, 2)
			logic.nextjiangtype = append(logic.nextjiangtype, 3)
		}

		return fkeb.getRetMsg(true, "修改开奖结果成功")
	}

	return fkeb.getRetMsg(false, "修改开奖结果失败")
}

func (fkeb *FKEBLogic) setSysPool(tableid int32, pool int64) proto.Message {
	if tableid == -1 {
		for _, table := range fkeb.GetTable().Room.Tables {
			logic := table.GameLogic.(*FKEBLogic)
			logic.changeSysPool(pool)
		}
	} else {
		t := fkeb.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*FKEBLogic)
			logic.changeSysPool(pool)
		} else {
			return fkeb.getRetMsg(false, "修改系统水池失败")
		}
	}

	return fkeb.getRetMsg(true, "修改系统水池成功")
}

func (fkeb *FKEBLogic) setZhuShui(tableid int32, zhushui int32) proto.Message {
	if tableid == -1 {
		for _, table := range fkeb.GetTable().Room.Tables {
			logic := table.GameLogic.(*FKEBLogic)
			logic.changeZhuShui(zhushui)
		}
	} else {
		t := fkeb.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*FKEBLogic)
			logic.changeZhuShui(zhushui)
		} else {
			return fkeb.getRetMsg(false, "修改注水率失败")
		}
	}

	return fkeb.getRetMsg(true, "修改注水率成功")
}

//获取系统赢钱数
func (fkeb *FKEBLogic) getRealUserWinByType(t int32) int64 {
	win := int64(0)
	odd := int64(fkeb.betvalue.getOddByType(t))
	user := fkeb.GetUserByUserID(fkeb.zhuangID)
	if user != nil && !user.IsRobot() {
		_, tax := fkeb.GetUserWinAndTax(fkeb.bettotal)
		if fkeb.currtax != 0 {
			tax = fkeb.currtax
		}

		robotwin := odd * fkeb.getRobotBetMoneyByType(t)
		win = fkeb.bettotalrobot - robotwin - tax
		return -win
	}

	for _, ud := range fkeb.userdata {
		if !ud.user.IsRobot() {
			win += odd * ud.getUserBetByType(t)
			win += ud.money - ud.user.Score
		}
	}

	mlog.Debug("玩家赢钱%d", win)
	return -win
}
