package by

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/internal/game/by/DataCfg"
	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/mylog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/roomlogic"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

const (
	ExperienceCashFlag = 6
)

const (
	BoomFish       = 3 //炸彈蟹
	LeiLong        = 4
	FishAll        = 5
	Zuantoulongxia = 6 //磚頭龍蝦
	Dujiaojing     = 7 //北冥神鯤
)

const (
	Basic   = 1 //基本子彈
	VIP     = 2 //VIP子彈
	ZlxBoom = 3 //磚頭蝦子彈
)

const (
	Delete = 1
)

type Bylogic struct {
	*gamelogic.ComputeGameLogicBase
	*PiSaiLogic                             //比賽相關資訊
	userdata               []*byuserdata    //用户数据
	fishwavetimer          *event.TimerFunc //鱼潮时间管理器
	fishtimer              *event.TimerFunc //生成鱼时间管理器
	bossfishtimer          *event.TimerFunc //生成BOSS鱼时间管理器
	teshufishtimer         *event.TimerFunc //生成特殊鱼时间管理器
	shoalfishtimer         *event.TimerFunc //生成魚群時間管理棄
	goldfishtimer          *event.TimerFunc
	asynctimer             *event.TimerFunc        //帧同步时间管理器
	updatetimer            *event.TimerFunc        //更新时间管理器
	rpctimer               *event.TimerFunc        //rpc时间管理器
	icetimer               *event.TimerFunc        //冰冻时间管理器
	rewardtimer            *event.TimerFunc        //放分時間管理棄
	fishkilltimer          *event.TimerFunc        //定時刪除魚
	reloadconfigtimer      *event.TimerFunc        //定時更新logic config
	nextsendmsg            *msg.Message            //接下来需要发送的消息
	framedata              *netproto.BY_FrameS2C   //下一帧数据内容
	framemsg               []*netproto.BY_FrameS2C //所有的帧消息
	frameid                int64                   //当前帧ID
	frametime              int32                   //帧同步时间
	maxframetime           int64                   //最长存储帧的时间
	fishmanager            *byfishmanager          //当前鱼管理器
	boommanager            []*byboommanager        //炮彈管理器
	icemanager             *byicemanager           //冰凍管理器
	bid                    []int32                 //炸弹的ID
	flag                   int32                   //房间类型
	boomIdOption           map[int32]int32         //用户可选炮弹倍率
	updatetime             int32                   //更新玩家分数时间
	basescore              int32                   //底分
	rpctime                int32                   //发送RPC消息时间
	dealplayer             int64                   //多久不发子弹将玩家踢出游戏，单位秒
	dealstring             string                  //踢人提示
	robotmsg               *msg.Message            //机器人动作消息
	framecount             int32                   //断线重连的帧数据个数
	framefishodd           int32                   //断线重连的帧数据的魚倍率在多少之上,倍率千分比
	shootinterval          float32                 //发射子弹间隔，单位秒
	boomspeed              float32                 //子弹移动速度
	icetime                int32                   //冰冻时间
	background             string                  //背景
	maxtotalfishlimit      int32                   //總共魚數量最大限制
	pochantime             int32                   //破產超過時間踢出
	robotfightrate         []int32                 //機器人打死魚的機率
	robotsuodingfishodd    int32                   //機器人可瞄準多少倍以上的魚. 倍率千分比
	dealplayertiptime      int64                   //踢出房前幾秒提示訊息
	pochantip              string                  //破產提示訊息
	pochantiptime          int32                   //破產踢出前幾秒提示訊息
	notifywinmoney         int32                   //廣播打中多少倍以上的魚
	fishwavetip            string                  //魚潮進行中，請稍等
	killfishcount          int32                   //達到魚總數量刪掉N隻魚
	killfishtime           int32                   //魚總數量刪掉N隻魚時間
	cfgs                   *config.Group
	maxboom                int32                       //最大的炸弹数
	zuantoulongxiamaxboom  int32                       //钻头龙虾最多1個炸彈
	nosupertoolreward      []int32                     //指定魚不提供個人或系統放分
	urs                    map[int32]*gamelog.UserData //玩家遊戲分數記錄
	updateMoney            int64
	zlxreceivespeed        int32                        //磚頭龍蝦砲彈接收速度(ms)
	robotcantkillfish      []int32                      //機器人不能打死特殊魚. 炸彈, 磚頭龍蝦等
	boompower              *float32                     //子彈威力
	checkSysPoolNum        *int32                       //每N發子彈統計並檢核系統水池
	syspoolTable           []*DataCfg.SyspoolTable      //系統水池區間
	playerpoolTable        []*DataCfg.PlayerpoolTable   //個人水線區間
	newPlayerTable         *NewPlayerPool               //新手優惠區間
	systemWinCommission    *SystemWinCommission         //系統水池贏錢抽水
	syspoolAdjustRate      float32                      //系統增減攻擊率
	fishtypetable          *DataCfg.FishTypeTable_ARRAY //詳細fishtype資料
	luckyTable             *DataCfg.LuckyTable_ARRAY    //幸運倍數Table
	fishSwell              map[int32]*SwellInfo         //魚的膨脹, map[fishid]rate
	freePlayerBoomPower    *FreePlayerPower             //未充值玩家的威力
	newplayeradjustsyspool int32                        //針對新手玩家優惠調整系統水池, 新手入場系統先降, 離場會補差校正.
	newplayerbonuslimit    int32                        //未充值新手優惠已經到贏超過14元,不再提供優惠
}

type FreePlayerPower struct {
	moneyLimit    int32
	monitorFishID []int32
	boomPower     int32
}

type SystemWinCommission struct {
	// * 房间编号
	RoomNumber *int32
	// * 子彈判斷發數
	BulletNumber *int32
	// * 金幣獲利區間
	CoinDivision *int32
	// * 初始抽水百分比
	IniCommissionRate *float32
	// * 越過區間增加百分比
	IncreaseRate *float32
	// * 抽水百分比上限
	MaxRate *float32
}

type NewPlayerPool struct {
	offerLimitCount int32 //優惠發數範圍
	poolTable       []*DataCfg.NewPlayerTable
}

type SwellInfo struct {
	rate      int32 //膨脹
	sumbullet int32 //打中子彈數量
}

//新建捕鱼实例
func NewByLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *Bylogic {
	glc := new(Bylogic)
	glc.ComputeGameLogicBase = gamelogic.NewComputeGameLogicBase(t, tfManager, cfgs)
	glc.init()
	glc.getconfig(cfgs)
	glc.cfgs = cfgs
	glc.AddCmd()
	glc.ReloadLogicConfig()
	glc.RegisterAfterGameStartFunc(func() {
		rand.Seed(time.Now().UnixNano())
		glc.startGame()
	})

	glc.RegisterBeforePiSaiStartFunc(func() {

	})
	glc.RegisterAfterPiSaiStartFunc((func() {
		rand.Seed(time.Now().UnixNano())
		glc.startGame()
	}))

	return glc
}

func (by *Bylogic) AddCmd() {
	command.AddCommand(NewByCmd(by))
	command.AddCommand(NewBySyspoolCmd(by))
}

//定時更新logic config設定檔
//func (by *Bylogic) ReloadConfigTimer() {
//	by.ReloadLogicConfig()
//	//by.reloadconfigtimer = by.TimerFuncManager.NewTimer(by.updatetime, by.ReloadConfigTimer)
//}

func (by *Bylogic) ReloadLogicConfig() {
	cgf := by.cfgs.GetConfig("logicconfig")
	by.updatetime = int32(cgf.GetConfigInt("logicsetting", "updatetime"))
	by.maxboom = int32(cgf.GetConfigInt("logicsetting", "maxboom"))
	by.zuantoulongxiamaxboom = int32(cgf.GetConfigInt32("logicsetting", "zuantoulongxiamaxboom"))
	by.rpctime = int32(cgf.GetConfigInt("logicsetting", "rpctime"))
	by.framecount = int32(cgf.GetConfigInt("logicsetting", "framecount"))
	by.framefishodd = int32(cgf.GetConfigFloat64("logicsetting", "framefishodd") / 1000 * float64(by.basescore))
	by.dealplayer = int64(cgf.GetConfigInt("logicsetting", "dealplayer")) * int64(time.Second)
	by.dealplayertiptime = int64(cgf.GetConfigInt("logicsetting", "dealplayertiptime")) * int64(time.Second)
	by.dealstring = cgf.GetConfigValue("logicsetting", "dealstring")
	by.maxtotalfishlimit = cgf.GetConfigInt32("logicsetting", "maxtotalfishlimit")
	by.pochantime = cgf.GetConfigInt32("logicsetting", "pochanleave")
	by.pochantip = cgf.GetConfigValue("logicsetting", "pochantip")
	by.pochantiptime = int32(cgf.GetConfigInt("logicsetting", "pochantiptime"))
	by.notifywinmoney = int32(cgf.GetConfigInt("logicsetting", "notifywinmoney"))
	by.fishwavetip = cgf.GetConfigValue("logicsetting", "fishwavetip")
	by.killfishcount = cgf.GetConfigInt32("logicsetting", "killfish")
	by.killfishtime = cgf.GetConfigInt32("logicsetting", "killfishtime")
	by.zlxreceivespeed = int32(cgf.GetConfigInt32("logicsetting", "zlxreceivespeed"))
	by.newplayeradjustsyspool = int32(cgf.GetConfigInt32("logicsetting", "newplayeradjustsyspool"))
	by.newplayerbonuslimit = int32(cgf.GetConfigInt32("logicsetting", "newplayerbonuslimit"))
	by.robotsuodingfishodd = int32(cgf.GetConfigFloat64("robotfight", "suodingfishodd") / 1000 * float64(by.basescore))
	robotfightstr := cgf.GetConfigValue("robotfight", "rate")
	robotfightarr := strings.Split(robotfightstr, ",")
	for _, v := range robotfightarr {
		r, _ := strconv.Atoi(v)
		by.robotfightrate = append(by.robotfightrate, int32(r))
	}

	nosupertoolrewardstr := cgf.GetConfigValue("logicsetting", "nosupertoolreward")
	nosupertoolrewardarr := strings.Split(nosupertoolrewardstr, ",")
	for _, v := range nosupertoolrewardarr {
		r, _ := strconv.Atoi(v)
		by.nosupertoolreward = append(by.nosupertoolreward, int32(r))
	}

	cantkillfishstr := cgf.GetConfigValue("logicsetting", "robotcantkillfish")
	cantkillfisharr := strings.Split(cantkillfishstr, ",")
	for _, v := range cantkillfisharr {
		r, _ := strconv.Atoi(v)
		by.robotcantkillfish = append(by.robotcantkillfish, int32(r))
	}

	if by.TableManager.Room.GameType == gamelogic.PiSai {
		by.gameEndLeave = int64(cgf.GetConfigInt("pisai", "gameendleave"))
	}
}

//开始游戏
func (by *Bylogic) startGame() {
	//开始没有鱼潮
	by.fishmanager.fishwave = false
	//开启帧同步
	by.asynctimer = by.TimerFuncManager.NewMilliSecTicker(by.frametime, by.sync)
	//开始生成鱼
	by.newfish()
	//開始生成黃金魚（高倍魚）
	by.startGoldFishTimer()
	//開始生成群群計時
	by.startShoalTimer()
	//开始生成BOSS鱼计时
	by.startBossTimer()
	//开始生成特殊鱼计时
	by.startTeshuTimer()
	//定時刪除魚
	by.startKillTimer()
	//开始鱼潮生成计时
	by.fishwavetimer = by.TimerFuncManager.NewTimer(by.fishmanager.fishwavetime, by.newfishwave)
	//启动更新玩家分数
	updatescop := util.NewIntScope(by.updatetime-1, by.updatetime+1)
	by.updatetimer = by.TimerFuncManager.NewNamedRandomTicker(updatescop, by.updateUserSorce, "update")
	//发送RPC消息
	rpcscop := util.NewIntScope(by.rpctime-60, by.rpctime+60)
	by.rpctimer = by.TimerFuncManager.NewNamedRandomTicker(rpcscop, by.sendWinMoney, "sendWinMoney")
}

func (by *Bylogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	defer util.PanicErrStack()

	if !by.CheckGameUserMessage(user) {
		by.KickUserFunc(user.UserID, "非法操作")
		return
	}

	//不是给游戏逻辑服务器的消息
	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	if gmsg.SClassID == int32(netproto.BY_GameMessageClassID_FrameC2SID) {
		by.userFrame(user, gmsg.MsgData.([]byte))
		return
	}
}

//初始化数据
func (by *Bylogic) init() {
	defer util.PanicErrStack()

	by.urs = make(map[int32]*gamelog.UserData, 0)
	by.framemsg = make([]*netproto.BY_FrameS2C, 0)
	by.nextsendmsg = new(msg.Message)
	by.nextsendmsg.BClassID = int32(netproto.MessageBClassID_Game)
	by.nextsendmsg.SClassID = int32(netproto.BY_GameMessageClassID_FrameS2CID)
	by.framedata = new(netproto.BY_FrameS2C)
	by.frameid = 0
	by.bid = make([]int32, 4)
	for i := int32(0); i < 4; i++ {
		by.bid[i] = 10000000 * i
	}

	by.fishwavetimer = nil
	by.fishmanager = NewFishManager()
	by.boomIdOption = make(map[int32]int32, 0)
	by.robotmsg = new(msg.Message)
	by.robotmsg.BClassID = int32(netproto.MessageBClassID_Game)
	by.robotmsg.SClassID = int32(netproto.BY_GameMessageClassID_RobotID)
	by.fishSwell = make(map[int32]*SwellInfo, 0)

	initPiSai(by)
}

//获取配置文件
func (by *Bylogic) getconfig(cfgs *config.Group) {
	//cgf := cfgs.GetConfig("logicconfig")
	scfgs := cfgs.GetConfig("gameserver")
	by.frametime = int32(scfgs.GetConfigInt("gameserver", "frametime"))
	by.flag = int32(scfgs.GetConfigInt("gameserver", "flag"))
	by.basescore = int32(scfgs.GetConfigInt("gameserver", "basescore"))
	by.shootinterval = scfgs.GetConfigFloat32("gameserver", "shootinterval")
	by.boomspeed = scfgs.GetConfigFloat32("gameserver", "boomspeed")
	by.icetime = scfgs.GetConfigInt32("gameserver", "icetime")

	by.loadTableData()
}

//讀取捕魚table參數資料
func (by *Bylogic) loadTableData() {
	defer util.PanicErrStack()

	by.getRoomTable()
	by.getBackgroundTable()
	by.getLuckyTable()
	by.getSyspoolTable()
	by.getPlayerpoolTable()
	by.getNewPlayerTable()
	by.getSystemWinCommissionTable()
	by.getFishTable()
	by.getFishShoalTable()
	by.getFishWaveTable()
	by.getFreePlayerTable()

	if by.TableManager.Room.GameType != gamelogic.PiSai {
		setmsg := by.getSettingMsg()
		by.BoardcastMessage(setmsg)
	}
}

//获取鱼类配置文件
func (by *Bylogic) getFishTable() {

	fishdata, pcfgerr := by.ReadProtoConfig("FishTable.bytes", func(data []byte) (proto.Message, error) {
		fishtableinfo := &DataCfg.FishTable_ARRAY{}
		err := proto.Unmarshal(data, fishtableinfo)
		return fishtableinfo, err
	})

	if pcfgerr != nil {
		mlog.Error("解析文件错误%v", pcfgerr)
		return
	}

	fishtablearr := fishdata.(*DataCfg.FishTable_ARRAY)
	fishtypearr := by.getFishTypeTable()
	by.fishmanager.setfishtabledata(fishtablearr, by.flag, fishtypearr)
	by.fishtypetable = fishtypearr
	//for _, v := range by.fishtypetable.Items {
	//	mlog.System("fishtypetable[%+v] ", v)
	//}

	by.setFishRewardTable(by.fishmanager.config.fishinfoarr)
}

func (by *Bylogic) getFishTypeTable() *DataCfg.FishTypeTable_ARRAY {
	fishdata, pcfgerr := by.ReadProtoConfig("FishTypeTable.bytes", func(data []byte) (proto.Message, error) {
		fishtableinfo := &DataCfg.FishTypeTable_ARRAY{}
		err := proto.Unmarshal(data, fishtableinfo)
		return fishtableinfo, err
	})

	if pcfgerr != nil {
		mlog.Error("解析FishTypeTable_ARRAY文件错误%v", pcfgerr)
		return nil
	}

	return fishdata.(*DataCfg.FishTypeTable_ARRAY)
}

func (by *Bylogic) setFishRewardTable(fishinfo []*fishinfo) {
	fishdata, pcfgerr := by.ReadProtoConfig("RewardTable.bytes", func(data []byte) (proto.Message, error) {
		fishtableinfo := &DataCfg.RewardTable_ARRAY{}
		err := proto.Unmarshal(data, fishtableinfo)
		return fishtableinfo, err
	})

	if pcfgerr != nil {
		mlog.Error("解析RewardTable文件错误%v", pcfgerr)
		return
	}

	fishtablearr := fishdata.(*DataCfg.RewardTable_ARRAY)
	by.fishmanager.setfishtypetabledata(fishtablearr, by.basescore)
}

//获取鱼群配置文件
func (by *Bylogic) getFishShoalTable() {
	fishdata, pcfgerr := by.ReadProtoConfig("FishShoalTable.bytes", func(data []byte) (proto.Message, error) {
		fishshoaltableinfo := &DataCfg.FishShoalTable_ARRAY{}
		err := proto.Unmarshal(data, fishshoaltableinfo)
		return fishshoaltableinfo, err
	})

	if pcfgerr != nil {
		mlog.Error("解析文件错误%v", pcfgerr)
		return
	}

	fishshoaltablearr := fishdata.(*DataCfg.FishShoalTable_ARRAY)

	by.fishmanager.setfishshoaltabledata(fishshoaltablearr)
}

//获取房间配置文件
func (by *Bylogic) getRoomTable() {
	roomtable, pcfgerr := by.ReadProtoConfig("RoomTable.bytes", func(data []byte) (proto.Message, error) {
		info := &DataCfg.RoomTable_ARRAY{}
		err := proto.Unmarshal(data, info)
		return info, err
	})

	if pcfgerr != nil {
		mlog.Error("解析文件RoomTable.bytes错误%v", pcfgerr)
		return
	}
	rt := roomtable.(*DataCfg.RoomTable_ARRAY)
	by.fishmanager.fishwavecf.fishwaveinfoarr = make([]*fishwaveinfo, 0)
	for i := 0; i < len(rt.GetItems()); i++ {
		//是这个房间的数据
		if int32(rt.GetItems()[i].GetRoomNumber()) == by.flag {
			by.fishmanager.fishwavetime = int32(rt.GetItems()[i].GetFishWaveInterval())
			by.fishmanager.setFishConfig(rt.GetItems()[i])
			by.fishmanager.parsefishwavestring(rt.GetItems()[i].GetFishWaveWeight())
			by.getboomtable(rt.GetItems()[i].GetBulletRateOption())
			by.geticetable(rt.GetItems()[i].GetIceOption())
			by.boompower = proto.Float32(float32(rt.GetItems()[i].GetBoomPower()) / 1000)
			by.checkSysPoolNum = proto.Int32(int32(rt.GetItems()[i].GetCountRate()))
			//mlog.Debug("BoomPower[%v] 每%v發子彈統計並檢核系統水池 ", *by.boompower, *by.checkSysPoolNum)
			break
		}
	}
}

//取得背景
func (by *Bylogic) getBackgroundTable() {
	boomtable, pcfgerr := by.ReadProtoConfig("BackgroundTable.bytes", func(data []byte) (proto.Message, error) {
		info := &DataCfg.BackgroundTable_ARRAY{}
		err := proto.Unmarshal(data, info)
		return info, err
	})

	if pcfgerr != nil {
		mlog.Error("解析文件BackgroundTable.bytes错误%v", pcfgerr)
		return
	}

	bt := boomtable.(*DataCfg.BackgroundTable_ARRAY)
	today := time.Now().Day()
	for i := 0; i < len(bt.GetItems()); i++ {
		if uint32(today) == bt.GetItems()[i].GetDate() {
			by.background = bt.GetItems()[i].GetBglist()
		}
	}
}

//取得未充值玩家的威力
func (by *Bylogic) getFreePlayerTable() {
	freeplayer, pcfgerr := by.ReadProtoConfig("FreePlayerTable.bytes", func(data []byte) (proto.Message, error) {
		info := &DataCfg.FreePlayerTable_ARRAY{}
		err := proto.Unmarshal(data, info)
		return info, err
	})

	if pcfgerr != nil {
		mlog.Error("解析文件FreePlayerTable.bytes错误%v", pcfgerr)
		return
	}

	bt := freeplayer.(*DataCfg.FreePlayerTable_ARRAY)
	//by.freePlayerBoomPower = bt
	by.freePlayerBoomPower = new(FreePlayerPower)
	by.freePlayerBoomPower.moneyLimit = bt.Items[0].GetMoneyLimit()
	by.freePlayerBoomPower.boomPower = bt.Items[0].GetBoomPower()
	arrFishId := strings.Split(bt.Items[0].GetMonitorFish(), ",")

	for _, v := range arrFishId {
		fishid, _ := strconv.Atoi(v)
		by.freePlayerBoomPower.monitorFishID = append(by.freePlayerBoomPower.monitorFishID, int32(fishid))
	}
	//mlog.Debug("freePlayerBoomPower[%+v] ", by.freePlayerBoomPower)
}

//幸運倍數
func (by *Bylogic) getLuckyTable() {
	luckytable, pcfgerr := by.ReadProtoConfig("LuckyTable.bytes", func(data []byte) (proto.Message, error) {
		info := &DataCfg.LuckyTable_ARRAY{}
		err := proto.Unmarshal(data, info)
		return info, err
	})

	if pcfgerr != nil {
		mlog.Error("解析文件LuckyTable.bytes错误%v", pcfgerr)
		return
	}

	bt := luckytable.(*DataCfg.LuckyTable_ARRAY)
	by.luckyTable = bt
}

//取得系統水池區間
func (by *Bylogic) getSyspoolTable() {
	syspooltable, pcfgerr := by.ReadProtoConfig("SyspoolTable.bytes", func(data []byte) (proto.Message, error) {
		info := &DataCfg.SyspoolTable_ARRAY{}
		err := proto.Unmarshal(data, info)
		return info, err
	})

	if pcfgerr != nil {
		mlog.Error("SyspoolTable.bytes错误%v", pcfgerr)
		return
	}

	bt := syspooltable.(*DataCfg.SyspoolTable_ARRAY)
	by.syspoolTable = append(by.syspoolTable, bt.Items...)

	//for _, v := range by.syspoolTable {
	//	mlog.Debug("取得系統水池區間 GetConditionOdds[%v] GetWaterAdjust[%v] GetArea[%v]", v.GetConditionOdds(), v.GetWaterAdjust(), v.GetArea())
	//}

	if len(by.syspoolTable) == int(0) {
		mlog.Error("系統水池區間 is nil")
	}
}

//取得個人水線區間
func (by *Bylogic) getPlayerpoolTable() {
	syspooltable, pcfgerr := by.ReadProtoConfig("PlayerpoolTable.bytes", func(data []byte) (proto.Message, error) {
		info := &DataCfg.PlayerpoolTable_ARRAY{}
		err := proto.Unmarshal(data, info)
		return info, err
	})

	if pcfgerr != nil {
		mlog.Error("PlayerpoolTable.bytes错误%v", pcfgerr)
		return
	}

	bt := syspooltable.(*DataCfg.PlayerpoolTable_ARRAY)
	by.playerpoolTable = append(by.playerpoolTable, bt.Items...)

	//for _, v := range by.playerpoolTable {
	//	mlog.Debug("取得個人水線區間 GetTotalMoneyWin[%v] GetAttackRate[%v]", v.GetTotalMoneyWin(), v.GetAttackRate())
	//}

	if len(by.playerpoolTable) == int(0) {
		example.invalid("個人水線區間 is nil")
	}
}

//取得新手優惠區間
func (by *Bylogic) getNewPlayerTable() {
	newplayertable, pcfgerr := by.ReadProtoConfig("NewPlayerTable.bytes", func(data []byte) (proto.Message, error) {
		info := &DataCfg.NewPlayerTable_ARRAY{}
		err := proto.Unmarshal(data, info)
		return info, err
	})

	if pcfgerr != nil {
		mlog.Error("NewPlayer.bytes错误%v", pcfgerr)
		return
	}

	bt := newplayertable.(*DataCfg.NewPlayerTable_ARRAY)

	by.newPlayerTable = new(NewPlayerPool)
	by.newPlayerTable.offerLimitCount = bt.Items[0].GetOfferLimitCount()
	by.newPlayerTable.poolTable = append(by.newPlayerTable.poolTable, bt.Items...)

	if len(by.newPlayerTable.poolTable) == int(0) {
		example.invalid("新手優惠區間 is nil")
	}
}

//設定系統贏錢抽水訊息
func (by *Bylogic) getSystemWinCommissionTable() {
	sysWinTable, pcfgerr := by.ReadProtoConfig("SystemWinCommissionTable.bytes", func(data []byte) (proto.Message, error) {
		info := &DataCfg.SystemWinCommissionTable_ARRAY{}
		err := proto.Unmarshal(data, info)
		return info, err
	})

	if pcfgerr != nil {
		mlog.Error("SystemWinCommissionTable.bytes错误%v", pcfgerr)
		return
	}

	bt := sysWinTable.(*DataCfg.SystemWinCommissionTable_ARRAY)
	by.systemWinCommission = new(SystemWinCommission)

	for _, v := range bt.GetItems() {
		if v.GetRoomNumber() == by.flag {
			by.systemWinCommission.RoomNumber = proto.Int32(v.GetRoomNumber())
			by.systemWinCommission.BulletNumber = proto.Int32(int32(v.GetBulletNumber()))
			by.systemWinCommission.CoinDivision = proto.Int32(int32(v.GetCoinDivision() * 100))
			by.systemWinCommission.IniCommissionRate = proto.Float32(float32(v.GetIniCommissionRate()) / 100)
			by.systemWinCommission.IncreaseRate = proto.Float32(float32(v.GetIncreaseRate()) / 100)
			by.systemWinCommission.MaxRate = proto.Float32(float32(v.GetMaxRate()) / 100)
		}
	}

	//mlog.System("systemWinCommission[%+v]", by.systemWinCommission)
}

//獲取炮彈設置文件
func (by *Bylogic) getboomtable(bulletRateOption string) {
	boomtable, pcfgerr := by.ReadProtoConfig("BoomTable.bytes", func(data []byte) (proto.Message, error) {
		info := &DataCfg.BoomTable_ARRAY{}
		err := proto.Unmarshal(data, info)
		return info, err
	})

	if pcfgerr != nil {
		mlog.Error("解析文件BoomTable.bytes错误%v", pcfgerr)
		return
	}

	by.boommanager = make([]*byboommanager, 0)
	bt := boomtable.(*DataCfg.BoomTable_ARRAY)
	oddstr := strings.Split(bulletRateOption, ",")
	for j := 0; j < len(oddstr); j++ {
		boomId, _ := strconv.ParseInt(oddstr[j], 10, 32)

		for i := 0; i < len(bt.GetItems()); i++ {
			if int32(boomId) == int32(bt.GetItems()[i].GetBoomNumber()) {
				bm := new(byboommanager)
				bm.boomNumber = int32(bt.GetItems()[i].GetBoomNumber())
				bm.boomName = bt.GetItems()[i].GetBoomName()
				bm.boomUIType = int32(bt.GetItems()[i].GetBoomUIType())
				bm.boomOdds = int32(bt.GetItems()[i].GetBoomRate())
				bm.boomMoney = int32(bt.GetItems()[i].GetBoomMoney()) * by.basescore
				by.boomIdOption[bm.boomNumber-1] = bm.boomNumber
				bm.boomSpeed = bt.GetItems()[i].GetBoomSpeed()
				bm.shootInterval = bt.GetItems()[i].GetShootInterval()
				bm.boomType = int32(bt.GetItems()[i].GetBoomType())
				bm.rebound = int32(bt.GetItems()[i].GetRebound())
				by.boommanager = append(by.boommanager, bm)
				//mlog.Debug("boomNumber[%v] boomtype[%v]", bm.boomNumber, bm.boomType)
			}
		}
	}
}

//取得炮彈資訊
func (by *Bylogic) getBoomInfo(boomNumber int32) *byboommanager {
	boomInfo := new(byboommanager)
	for _, v := range by.boommanager {
		if boomNumber == v.boomNumber {
			boomInfo = v
		}
	}

	return boomInfo
}

func (by *Bylogic) geticetable(iceOption uint32) {
	icetable, pcfgerr := by.ReadProtoConfig("IceTable.bytes", func(data []byte) (proto.Message, error) {
		info := &DataCfg.IceTable_ARRAY{}
		err := proto.Unmarshal(data, info)
		return info, err
	})

	if pcfgerr != nil {
		mlog.Error("解析文件IceTable.bytes错误%v", pcfgerr)
		return
	}

	by.icemanager = new(byicemanager)
	it := icetable.(*DataCfg.IceTable_ARRAY)
	for i := 0; i < len(it.GetItems()); i++ {
		if it.GetItems()[i].GetIceNumber() == iceOption {
			by.icemanager.iceNumber = int32(it.GetItems()[i].GetIceNumber())
			by.icemanager.iceName = it.GetItems()[i].GetIceName()
			by.icemanager.iceUIType = int32(it.GetItems()[i].GetIceUIType())
			by.icemanager.iceTime = int32(it.GetItems()[i].GetIceTime())
			by.icemanager.iceResetTime = int32(it.GetItems()[i].GetIceResetTime())
			by.icemanager.iceMoney = int32(it.GetItems()[i].GetIceMoney())
		}
	}
}

//获取鱼潮配置文件
func (by *Bylogic) getFishWaveTable() {
	fishwavetable, pcfgerr := by.ReadProtoConfig("FishWaveTable.bytes", func(data []byte) (proto.Message, error) {
		info := &DataCfg.FishWaveTable_ARRAY{}
		err := proto.Unmarshal(data, info)
		return info, err
	})

	if pcfgerr != nil {
		mlog.Error("FishWaveTable解析文件错误%v", pcfgerr)
		return
	}

	pathtable, pcfgerr2 := by.ReadProtoConfig("PathTable.bytes", func(data []byte) (proto.Message, error) {
		info := &DataCfg.PathTable_ARRAY{}
		err := proto.Unmarshal(data, info)
		return info, err
	})

	if pcfgerr2 != nil {
		mlog.Error("PathTable解析文件错误%v", pcfgerr)
		return
	}

	dt := fishwavetable.(*DataCfg.FishWaveTable_ARRAY)
	pt := pathtable.(*DataCfg.PathTable_ARRAY)
	by.fishmanager.setFishWaveTalbeData(dt, pt)
}

func (by *Bylogic) GetFishRewardList() []*netproto.FishSystemRewardDataList {
	empty := []*netproto.FishSystemRewardDataList{}

	if by.TableManager.Room.SysRewardList == nil {
		return empty
	}

	return by.TableManager.Room.SysRewardList.GetFishRewardList()
}

func (by *Bylogic) GetFishVIPUserList() []*netproto.FishVIPUserRewardDataList {
	empty := []*netproto.FishVIPUserRewardDataList{}

	if by.TableManager.Room.VIPUserRewardList == nil {
		return empty
	}

	return by.TableManager.Room.VIPUserRewardList.GetFishVIPUserList()
}

func (by *Bylogic) bossEnter() {
	defer util.PanicErrStack()

	//鱼潮不产生新的BOSS鱼
	if !by.canDone() {
		by.bossfishtimer = nil
		return
	}

	//by.framedata.BossNear = proto.Int32(1)
	by.bossfishtimer = by.TimerFuncManager.NewTimer(by.fishmanager.bossentertime, by.newBossFish)
}

//产生一条新的BOSS鱼
func (by *Bylogic) newBossFish() {
	defer util.PanicErrStack()

	if !by.canDone() {
		by.bossfishtimer = nil
		return
	}

	fish := by.fishmanager.newBossFish(by)

	if fish != nil {
		by.framedata.BossNear = proto.Int32(fish.fishtype)
		by.framedata.FishOdds = proto.Int32(fish.fishodd)

		fi, ok := by.fishmanager.config.getfishinfobytype(fish.fishtype)
		if ok && fi.bosstypeinfo.fishtype == Dujiaojing {
			by.framedata.FishOdds = proto.Int32(0)
		}
	}

	//mlog.Debug("newboss near[%v] fishodd[%v]", fish.fishtype, fish.fishodd)
	//fi, ok := by.fishmanager.config.getfishinfobytype(fish.fishtype)
	//if ok {
	//	mlog.System("创建一条BOSS鱼 fishid[%v] luckybonus[%v]", fi.fishid, fi.luckybonus)
	//}

	by.addnewfish(fish)
	by.startBossTimer()
}

func (by *Bylogic) newGoldFish() {
	defer util.PanicErrStack()

	if !by.canDone() {
		by.goldfishtimer = nil
		return
	}

	if !by.isTotalFishMaxLimit() {
		fish := by.fishmanager.newGoldFish(by)
		by.addnewfish(fish)
	}

	by.startGoldFishTimer()
}

//产生一条新的特殊鱼
func (by *Bylogic) newTeshuFish() {
	defer util.PanicErrStack()

	//鱼潮不产生新的特殊鱼
	if !by.canDone() {
		by.teshufishtimer = nil
		return
	}

	fish := by.fishmanager.newTeshuFish(by)
	if fish != nil {
		mlog.Debug("产生一条新的特殊鱼 fishtype[%v] fishodd[%v]", fish.fishtype, fish.fishodd)
		by.framedata.TeshuNear = proto.Int32(fish.fishtype)
		by.framedata.FishOdds = proto.Int32(fish.fishodd)
		by.addnewfish(fish)
	}

	by.startTeshuTimer()
}

//產生一條魚群
func (by *Bylogic) newFishShoal() {
	mlog.Debug("產生一條魚群")
	defer util.PanicErrStack()

	if !by.canDone() {
		by.fishtimer = nil
		return
	}

	if !by.isTotalFishMaxLimit() {
		t := by.fishmanager.getCreateFishShoalType()
		fish := by.fishmanager.createFishShoal(t, by)

		by.addnewfish(fish)
	}

	by.startShoalTimer()
}

//产生一条新的鱼
func (by *Bylogic) newfish() {
	defer util.PanicErrStack()

	//鱼潮不产生新的鱼
	if !by.canDone() {
		by.fishtimer = nil
		return
	}

	fish, t := by.fishmanager.newfish(by)
	by.addnewfish(fish)
	by.fishtimer = by.TimerFuncManager.NewMilliSecTimer(t, by.newfish)
}

//进入鱼潮
func (by *Bylogic) addEnterFishWave() {
	defer util.PanicErrStack()

	by.framedata.AdvanceNoticeType = proto.Int32(1)
}

//添加鱼潮
func (by *Bylogic) addFishWave() {
	defer util.PanicErrStack()

	mlog.Debug("添加鱼潮")
	//fishArrCount := make(map[int32]int32, 0)
	//for _, v := range by.fishmanager.fishdata {
	//	if v.fishSeq > 100000 {
	//		fishArrCount[v.fishtype]++
	//	}
	//}

	fishWave := new(netproto.BY_FishWaveManage)
	fishWave.FishWaveID = proto.Int32(by.fishmanager.fishwavetype)
	fishWave.FishInfo = make([]*netproto.BY_FishInfoWave, 0)
	for _, v := range by.fishmanager.fishwavecf.fishwaveinfoarr {
		if v.fishwavetype == (by.fishmanager.fishwavetype) {

			for _, o := range v.fisharr {
				fishInfoWave := new(netproto.BY_FishInfoWave)
				fishInfoWave.FishId = proto.Int32(o.fishType)
				fishInfoWave.Count = proto.Int32(o.count)
				fishInfoWave.PathName = proto.String(o.pathName)
				fishInfoWave.Interval = proto.Int32(o.interval)
				fishWave.FishInfo = append(fishWave.FishInfo, fishInfoWave)
			}
		}
	}

	by.framedata.FishWave = fishWave
	//for _, v := range by.framedata.FishWave.FishInfo {
	//	mlog.Debug("frameid[%v] 增加魚潮 fishid[%v] count[%v]", by.framedata.GetFrameID(), v.GetFishId(), v.GetCount())
	//}
}

//产生鱼潮
func (by *Bylogic) newfishwave() {
	defer util.PanicErrStack()

	by.framemsg = make([]*netproto.BY_FrameS2C, 0)
	by.killFish()
	by.fishmanager.newFishWave()
	by.fishwavetimer = by.TimerFuncManager.NewTimer(by.fishmanager.enterfishwavetime, by.enterFishWave)
	by.addEnterFishWave()
}

//进入鱼潮
func (by *Bylogic) enterFishWave() {
	defer util.PanicErrStack()

	by.fishwavetimer = by.TimerFuncManager.NewTimer(by.fishmanager.fishwavelasttime, by.endFishWave)
	by.addFishWave()
}

//鱼潮结束
func (by *Bylogic) endFishWave() {
	defer util.PanicErrStack()

	if !by.fishmanager.fishwave || by.fishmanager.ice {
		mlog.Debug("不是鱼潮了")
		return
	}

	if by.fishwavetimer != nil {
		by.fishwavetimer.Stop()
		by.fishwavetimer = nil
	}

	by.framemsg = make([]*netproto.BY_FrameS2C, 0)
	by.fishmanager.endFishWave()
	by.killFish()
	by.newfish()
	by.startGoldFishTimer()
	//開啟魚群計時
	by.startShoalTimer()
	//开启BOSS鱼计时
	by.startBossTimer()
	//开启特殊鱼计时
	by.startTeshuTimer()
	//开启鱼潮计时
	by.fishwavetimer = by.TimerFuncManager.NewTimer(by.fishmanager.fishwavetime, by.newfishwave)
}

//捕鱼帧同步
func (by *Bylogic) sync() {
	defer util.PanicErrStack()
	//defer func() {
	//	err := util.PanicErrStack()
	//	//if err != nil {
	//	//	mlog.Fatal("系統掛掉")
	//	//}
	//}()

	starttime := time.Now().UnixNano()
	//是否有真人
	hasuser := false
	by.getNextSendMsg()
	by.setPiSaiNextFramePlayTime()
	for _, ur := range by.AllTableUser() {
		ud := by.getuserdatabyuserid(ur.UserID)

		//踢出破產
		for _, v := range by.framedata.PoChan {
			if ur.UserType == gamelogic.UserTypeRebot && v == ur.UserID {
				by.KickUserFunc(ur.UserID, "robot 破產")
				by.UserLeave(ur)
			}

			tipPochanTime := by.pochantime - by.pochantiptime
			by.TimerFuncManager.NewTimer(tipPochanTime, func() {
				user := by.GetUserByUserID(v)
				if user == nil || by.pochantip == "" || by.pochantiptime == 0 || user.UserType == gamelogic.UserTypeRebot {
					return
				}
				by.SendMessageToUser(user, by.getLeaveTipMessage(by.pochantip, by.pochantiptime))
			})

			by.TimerFuncManager.NewTimer(by.pochantime, func() {
				user := by.GetUserByUserID(v)
				if user == nil || by.pochantip == "" || by.pochantiptime == 0 || user.UserType == gamelogic.UserTypeRebot {
					return
				}
				by.KickUserFunc(v, "-")
				by.UserLeave(user)
			})
		}

		if !ur.IsRobot() {
			hasuser = true
			//不管机器人，真人才踢出
			//长时间未操作踢出用户

			if starttime-ud.lastshoottime > by.dealplayer-int64(by.dealplayertiptime) {
				t := by.dealplayertiptime / int64(time.Second)
				by.SendMessageToUser(ur, by.getLeaveTipMessage(by.dealstring, int32(t)))
			}
			if starttime-ud.lastshoottime > by.dealplayer {
				mlog.Debug("长时间未操作踢出用户%d", ur.UserID)
				// 正式斷線後，計算完分數會清除
				//by.deluserdatabyuserid(ur.UserID)
				by.KickUserFunc(ur.UserID, "-")
			}
		}

		ud.setLastFrameMoney()
	}

	nextMsg := by.nextsendmsg.MsgData.(*netproto.BY_FrameS2C)
	if nextMsg.Shoot != nil {
		user := by.GetUserByUserID(nextMsg.Shoot[0].GetUserID())
		if !user.IsRobot() {
			mlog.Debug("BoardcastMessage userid[%v] shootX[%v] shootFishId[%v]", nextMsg.Shoot[0].GetUserID(), nextMsg.Shoot[0].GetPos().GetX(), nextMsg.Shoot[0].GetFishID())
			mlog.Debug("BoardcastMessage userid[%v] SuoDing[%v]", nextMsg.Shoot[0].GetUserID(), nextMsg.GetSuoDing())
		}
	}

	//logFrame := by.nextsendmsg.MsgData.(*netproto.BY_FrameS2C)
	//if logFrame.GetBoom() != nil || logFrame.GetShoot() != nil  {
	//	mlog.System("userid[%v] boom[%+v] shoot[%+v]", by.userdata[0].user.UserID, logFrame.GetBoom(), logFrame.GetShoot())
	//}

	by.BoardcastMessage(by.nextsendmsg)

	//本次发送完成清理掉消息
	by.framedata = new(netproto.BY_FrameS2C)
	by.frameid++

	if hasuser {
		nextMsg := by.nextsendmsg.MsgData.(*netproto.BY_FrameS2C)
		if nextMsg.Shoot != nil {
			user := by.GetUserByUserID(nextMsg.Shoot[0].GetUserID())
			by.robotmsg.MsgData = nil
			if !user.IsRobot() {
				f := by.fishmanager.getFishInfoByFishID(nextMsg.Shoot[0].GetFishID())
				if f != nil {
					fishcf, ok := by.fishmanager.config.getfishinfobytype(f.fishtype)
					if ok && fishcf.rewardinfo[0].reward >= by.robotsuodingfishodd {
						//mlog.System("sendToRobot userid[%v] shootX[%v] shootFishId[%v] fishodd[%v]",
						//	nextMsg.Shoot[0].GetUserID(), nextMsg.Shoot[0].GetPos().GetX(), nextMsg.Shoot[0].GetFishID(), fishcf.rewardinfo[0].reward)
						by.robotmsg.MsgData = by.nextsendmsg.MsgData
					}
				}
			}
		}

		by.sendToRobot()
	} else {
		//清理一次数据，保证机器人的子弹数据无问题
		for _, ud := range by.userdata {
			ud.delAllBoom()
		}
	}

	by.GameIsSlave()

	if len(by.AllTableUser()) == 0 {
		if p, ok := by.getPiSaiUserData(); ok {
			totalPlayTime := (p.EndTime - time.Now().UnixNano()) / int64(time.Second)
			if (!p.IsPlaying && !p.IsGameStart) || totalPlayTime <= 0 {
				by.gameEnd()
			}
		} else {
			by.gameEnd()
		}
	}
}

//遊戲切到備援後踢出所有玩家
func (by *Bylogic) GameIsSlave() {
	if by.TableManager.Room.GameType == gamelogic.PiSai {
		return
	}
	if by.TableManager.Room.HostAndBackup == roomlogic.HostTypeSlave {
		by.stopCreatFish()
		by.TimerFuncManager.NewTimer(2, func() {
			for _, ur := range by.AllTableUser() {
				user := by.GetUserByUserID(ur.UserID)
				by.KickUserFunc(user.UserID, "-")
				by.UserLeave(user)
			}
		})
	}
}

func (by *Bylogic) getNextSendMsg() {
	defer util.PanicErrStack()

	for _, ud := range by.userdata {
		userinfo := new(netproto.BY_UserInfo)
		userinfo.UserID = proto.Int32(ud.user.UserID)
		userinfo.UserMoney = proto.Int64(ud.user.Score + ud.money)

		piSaiUser, ok := by.getPiSaiUserData()
		if ok {
			userinfo.UserMoney = proto.Int64(piSaiUser.Score)
			userinfo.GameCompetitionBonus = proto.Int32(int32(piSaiUser.Bonus))
		}

		by.framedata.UserInfo = append(by.framedata.UserInfo, userinfo)
	}

	if by.fishmanager.ice {
		by.framedata.IsIce = proto.Bool(by.fishmanager.ice)
	}

	by.framedata.IsFishWave = proto.Bool(by.fishmanager.fishwave)

	by.framedata.FrameID = proto.Int64(by.frameid)
	by.nextsendmsg.MsgData = by.framedata
	by.framemsg = append(by.framemsg, by.framedata)

	//if len(by.framedata.GetFish()) > 1 {
	//
	//	sum := int64(0)
	//	for _, v := range by.framedata.GetFish() {
	//		sum += v.GetMoney()
	//	}
	//	mlog.Debug("通知client抓到魚 sum[%v]", sum)
	//}

	//if len(by.framedata.Boom) > 0 {
	//	mlog.Debug("Boom[%v]", by.framedata.Boom)
	//}

	//if len(by.framedata.Shoot) > 0 {
	//	mlog.System("shoot [%v]", by.framedata.GetShoot())
	//}

	//if len(by.framedata.GetNewFish()) > 0 || by.framedata.FishWave.GetFishWaveID() > 0 {
	//	mlog.Debug("给用户发送帧消息 GetNewFish[%v]", by.framedata.GetNewFish())
	//	mlog.Debug("给用户发送帧消息 GetNewFish[%v]", by.framedata.GetFish())
	//
	//	for _, f := range by.framedata.NewFish {
	//		mlog.Debug("给用户发送帧消息 NewFish[%v]", f.GetFishID())
	//	}
	//	mlog.Debug("给用户发送帧消息 GetFrameID[%v]", by.framedata.GetFrameID())
	//	mlog.Debug("给用户发送帧消息 FishWaveID[%v]", by.framedata.FishWave.GetFishWaveID())
	//	mlog.Debug("给用户发送帧消息 GetAdvanceNoticeType[%v]", by.framedata.GetAdvanceNoticeType())
	//}

	//if by.framedata.Shoot != nil {
	//	for _, v := range by.framedata.Shoot {
	//		mlog.Debug("tableid=%d,frameid=%d,用户%d发射子弹ID:%d", by.TableNo(), by.frameid, v.GetUserID(), v.GetBoomID())
	//	}
	//}
}

//发送用户消息
func (by *Bylogic) sendUserMsg(hasuser bool) {
	by.framedata.FrameID = proto.Int64(by.frameid)
	by.nextsendmsg.MsgData = by.framedata
	//不给机器人发送帧消息
	for _, ur := range by.AllTableUser() {
		if !ur.IsRobot() {
			by.SendMessageToUser(ur, by.nextsendmsg)
		} else if hasuser {
			by.SendMessageToUser(ur, by.robotmsg)
		}
	}

	by.framemsg = append(by.framemsg, by.framedata)

	//本次发送完成清理掉消息
	by.framedata = new(netproto.BY_FrameS2C)
}

func (by *Bylogic) canUseBoomNumber(ud *byuserdata) bool {
	isCheck := false

	boomInfo := by.getBoomInfo(by.boomIdOption[ud.boomNumber])
	if boomInfo.boomNumber >= 0 {
		isCheck = true
	}

	return isCheck
}

//添加一次用户射击
func (by *Bylogic) addshoot(data *netproto.BY_ShootC2S, ud *byuserdata) {
	defer util.PanicErrStack()

	shootinfo := new(netproto.BY_ShootS2C)
	shootinfo.Pos = data.Pos

	shootinfo.UserID = proto.Int32(ud.user.UserID)
	if data.FishID != nil {
		shootinfo.FishID = data.FishID
	}

	if ud.user.IsRobot() {
		shootinfo.IsRobot = proto.Bool(true)
	}

	shootinfo.BoomNumber = proto.Int32(ud.boomNumber)
	shootinfo.BoomID = proto.Int32(by.getBoomId(ud.user.SeatNo))
	shootinfo.BoomType = proto.Int32(by.getBoomInfo(by.boomIdOption[ud.boomNumber-1]).boomType)
	boomOdd := by.getBoomInfo(by.boomIdOption[ud.boomNumber-1]).boomOdds

	ud.addBoom(shootinfo.GetBoomID(), boomOdd)
	if !ud.user.IsRobot() {
		mlog.Debug("frameid=%d,用户%d射击,BoomID[%d]  BoomNumber[%v] rate[%v] boomtype[%v] ", by.frameid, ud.user.UserID, shootinfo.GetBoomID(), ud.boomNumber, ud.rate, shootinfo.GetBoomType())
	}

	by.framedata.Shoot = append(by.framedata.Shoot, shootinfo)

	boomMoney := by.getBoomInfo(by.boomIdOption[ud.boomNumber-1]).boomMoney
	//mlog.Debug("userid[%v] userShoot boomNumber : %v boomMoney[%v]", ud.user.UserID, int64(ud.boomNumber), boomMoney)

	ud.updateUserMoney(int32(-boomMoney))
	by.updatePiSaiUserScore(-boomMoney)
	by.updateSysPoolWinLoss(ud, int32(-boomMoney))

	if !ud.user.IsRobot() {
		by.addSyspoolMoneyAndBullet(int64(boomMoney), int64(1), ud)
		ud.sumbetmoney += int64(boomMoney)
	}

	example.invalid("frameid[%v] 用户%d的钱%d，用户射擊 BoomId[%v] 砲彈編號[%v] 砲倍率[%v] 打出的钱%v，本次扣钱:%v",
		by.frameid, ud.user.UserID, ud.user.Score, shootinfo.GetBoomID(), shootinfo.GetBoomNumber(), ud.rate, ud.money, boomMoney)
}

//Server先即時計算系統總輸贏
func (by *Bylogic) updateSysPoolWinLoss(ud *byuserdata, money int32) {
	if !ud.user.IsRobot() && by.TableManager.Room.FishSysPoolData != nil {

		by.TableManager.Room.FishSysPoolData.AfterAdjustSyspoolWinLoss = proto.Int64(by.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() + int64(money))

		logMsg := fmt.Sprintf("即時計算系統總輸贏 調整前系統水池[%v] 變動金額[%v] 調整後系統水池[%v] ,userid[%v]",
			-by.TableManager.Room.FishSysPoolData.GetSysSumWinLoss(),
			-money,
			-by.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss(),
			ud.user.UserID)
		mlog.Debug(logMsg)
		mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), logMsg)

		by.TableManager.Room.FishSysPoolData.SysSumWinLoss = proto.Int64(by.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss())
	}
}

//新增系統水池, 該魚場玩家總押注金額, 總子彈發數
func (by *Bylogic) addSyspoolMoneyAndBullet(betMoney int64, bullet int64, ud *byuserdata) {
	//mlog.Debug("新增, 該魚場玩家總押注金額[%v], 總子彈發數[%v]", betMoney, bullet)
	if by.TableManager.Room.FishSysPoolData == nil {
		by.TableManager.Room.FishSysPoolData = new(netproto.FishSysPoolData)
		by.TableManager.Room.FishSysPoolData.SumBetMoney = proto.Int64(betMoney)
		by.TableManager.Room.FishSysPoolData.SumBullet = proto.Int64(bullet)
	} else {
		by.TableManager.Room.FishSysPoolData.SumBetMoney = proto.Int64(by.TableManager.Room.FishSysPoolData.GetSumBetMoney() + betMoney)
		by.TableManager.Room.FishSysPoolData.SumBullet = proto.Int64(by.TableManager.Room.FishSysPoolData.GetSumBullet() + bullet)
	}

	logMsg := fmt.Sprintf("userid[%v] score[%v] 新增, 該魚場玩家總押注金額[%v]betMoney[%v], 總子彈發數[%v]bullet[%v] ",
		ud.user.UserID, ud.user.Score, by.TableManager.Room.FishSysPoolData.GetSumBetMoney(), betMoney, by.TableManager.Room.FishSysPoolData.GetSumBullet(), bullet)
	mlog.Debug(logMsg)
	//mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), logMsg)
}

//公告獲得魚
func (by *Bylogic) notifyaddfish(fish *fishinfo, ud *byuserdata, money int32) {
	boomMoney := float64(by.getBoomInfo(by.boomIdOption[ud.boomNumber-1]).boomMoney) / float64(100)
	fishOdd := float64(money) / float64(boomMoney)

	if by.notifywinmoney <= int32(fishOdd) {

		notifyMoney := float64(money) / float64(100)
		strMoney := strconv.FormatFloat(notifyMoney, 'f', 2, 64)

		boomMoneyStr := ""
		if boomMoney >= 0.1 {
			boomMoneyStr = strconv.FormatFloat(boomMoney, 'f', 1, 64)
		} else {
			boomMoneyStr = strconv.FormatFloat(boomMoney, 'f', 2, 64)
		}

		tip := fmt.Sprintf("恭喜%v用%v倍炮打下%v，获得%v金币！", ud.user.UserID, boomMoneyStr, fish.name, strMoney)
		mlog.Debug("公告 %v", tip)

		if fish.fishid >= 37 {
			mlog.Debug("公告 %v", tip)
		}

		pmsg := new(msg.Message)
		pmsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
		pmsg.SClassID = int32(netproto.NotifyServerClassID_AnnFishRegID)
		pdata := new(netproto.BY_NotifyFishMsg)
		pdata.TipMessage = proto.String(tip)
		pmsg.MsgData = pdata

		by.NotifyFishMessage(pmsg)
	}
}

//個人水池記錄
func (by *Bylogic) recardPlayerPool(money int32, ud *byuserdata) {
	if ud.user.IsRobot() {
		return
	}

	if by.TableManager.Room.FishPlayerRewardDetail == nil {
		//mlog.Debug("個人水池記錄 inti empty ")
		by.TableManager.Room.FishPlayerRewardDetail = make(map[int32]*netproto.FishPlayerPoolDetailList, 0)
	}

	boomMoney := by.getBoomInfo(by.boomIdOption[ud.boomNumber-1]).boomMoney

	if v, ok := by.TableManager.Room.FishPlayerRewardDetail[ud.user.UserID]; ok {
		if v.GetUserID() == ud.user.UserID {
			v.SpendMoney = proto.Int64(v.GetSpendMoney() + int64(boomMoney))
			v.HitMoney = proto.Int64(v.GetHitMoney() + int64(money))

			if money != 0 {
				mlog.Debug("個人水池記錄 打死魚的總金額[%v]boomMoney[%v] 花出的總金額[%v] money[%v] ",
					v.GetHitMoney(), money, v.GetSpendMoney(), boomMoney)
			}

			return
		}
	}

	newPlayer := new(netproto.FishPlayerPoolDetailList)
	newPlayer.UserID = proto.Int32(ud.user.UserID)
	newPlayer.SpendMoney = proto.Int64(int64(boomMoney))
	newPlayer.HitMoney = proto.Int64(int64(money))
	by.TableManager.Room.FishPlayerRewardDetail[ud.user.UserID] = newPlayer
	//mlog.System("新手優惠記錄 %+v", newPlayer)
}

//添加一条捕获到的鱼
func (by *Bylogic) addfish(data *netproto.BY_Fish) {
	defer util.PanicErrStack()
	by.framedata.Fish = append(by.framedata.Fish, data)
}

//打中魚的多階段贏錢表演
func (by *Bylogic) addMoneyShow(fishcf *fishinfo, odd int32) int32 {
	defer util.PanicErrStack()
	money := int32(0)
	//計算規則為，600 + N + M，其中N為浮動倍率且1 ≦ N ≦ 6000；M為補正倍率且1 ≦ M ≦ 99。
	//　ｂ．接著，依上述計算規則骰5次，得到5個隨機數值後，然後取最小的值，該值即為玩家獲得的總倍率。
	smallmoney := int32(999999999)
	for i := 0; i < 4; i++ {
		n := util.RndInt(int(fishcf.rewardinfo[0].reward), int(fishcf.rewardinfo[1].reward))
		mratestr := strings.Split(fishcf.bosstypeinfo.mrate, ",")
		m1, _ := strconv.Atoi(mratestr[0])
		m2, _ := strconv.Atoi(mratestr[1])
		mrate := util.RndInt(m1, m2)
		sum := int32(n+mrate) * odd
		if sum < smallmoney {
			smallmoney = sum
		}
	}

	//得分分配比例
	steparr := strings.Split(fishcf.bosstypeinfo.bodysep, ",")
	stepstr := util.RndInt(0, len(steparr))
	stepLen, _ := strconv.Atoi(steparr[stepstr])
	winseplist := make([]float32, stepLen) //分配比例
	moneylist := make([]int32, stepLen)    //分配Money

	winseparr := strings.Split(fishcf.bosstypeinfo.winsep, ",")
	index := 0
	for _, w := range winseparr {
		winsep, _ := strconv.ParseFloat(w, 32)
		winseplist[index] += float32(winsep)
		mlog.Debug("winseplist[%v] = %v stepLen[%v]", index, winseplist[index], stepLen)
		if index == stepLen-1 {
			index = 0
		} else {
			index++
		}
	}

	check := float32(0)
	for _, v := range winseplist {
		check += v
	}
	if check != 1 {
		mlog.Error("打中鯤 獎金check[%v]", check)
	}

	for k, v := range winseplist {
		moneylist[k] = int32(float32(smallmoney) * v)
	}

	winlist := new(netproto.BY_WinList)
	winlist.Time = proto.Int32(2)
	winlist.Money = moneylist
	by.framedata.FishWinList = winlist
	mlog.Debug("smallmoney[%v] sumMoney[%v] 打中魚的多階段贏錢表演 %+v ", smallmoney, smallmoney, winlist)
	money = smallmoney
	return money
}

//增加魚的體積
func (by *Bylogic) addFishSwell(fishid int32) {
	if f, ok := by.fishmanager.fishdata[fishid]; ok {
		v, _ := by.fishmanager.config.getfishinfobytype(f.fishtype)
		//mlog.System("before膨脹 fishid[%v] rate[%v] typeinfo[%+v]", fishid, v.bosstypeinfo.swellrate, v)
		if v.bosstype == Dujiaojing {
			by.addFishSwellRate(fishid, v.bosstypeinfo)
			by.addFishSwellBulletCount(fishid)

			fishSwell := new(netproto.BY_FishSwellRate)
			fishSwell.FishID = proto.Int32(fishid)
			fishSwell.Rate = proto.Int32(by.fishSwell[fishid].rate)
			by.framedata.FishSwell = fishSwell
			//mlog.System("膨脹 fishid[%v] rate[%v]", fishid, by.fishSwell[fishid])
		}
	}
}

//增加膨脹率
func (by *Bylogic) addFishSwellRate(fishid int32, typeinfo fishtypeinfo) {
	if s, sok := by.fishSwell[fishid]; sok {
		bulletArr := strings.Split(typeinfo.bulletnumber, ",")
		randBullet := util.RndInt(0, len(bulletArr))
		bulletCount, _ := strconv.Atoi(bulletArr[randBullet])
		canSwell := (s.sumbullet % int32(bulletCount)) == 0
		if canSwell && s.rate < typeinfo.swellmaxrate {
			by.fishSwell[fishid].rate = s.rate + typeinfo.swellrate
			//mlog.System("膨脹 fishid[%v] fishSwell[%+v]", fishid, by.fishSwell[fishid])
		}
	} else {
		by.fishSwell[fishid] = &SwellInfo{rate: 0, sumbullet: 1}
	}
}
func (by *Bylogic) addFishSwellBulletCount(fishid int32) {
	if _, sok := by.fishSwell[fishid]; sok {
		by.fishSwell[fishid].sumbullet++
		//mlog.System("BulletCount fishid[%v] fishSwell[%+v]", fishid, by.fishSwell[fishid])
	}
}

//添加一个炸弹
func (by *Bylogic) addboom(data *netproto.BY_BoomS2C) {
	defer util.PanicErrStack()

	by.framedata.Boom = append(by.framedata.Boom, data)
}

//添加一条鱼
func (by *Bylogic) addnewfish(fish *byfishdata) {
	defer util.PanicErrStack()

	if fish == nil {
		return
	}

	data := new(netproto.BY_NewFish)
	data.FishID = proto.Int32(fish.fishSeq)
	data.FishType = proto.Int32(fish.fishtype)
	data.FishNumber = proto.Int32(fish.fishnumber)
	data.FishCount = proto.Int32(fish.fishcount)
	by.framedata.NewFish = append(by.framedata.NewFish, data)

	mlog.NewDefaultLogger().SetCallDepLen(5).Debug("添加一条鱼 fishSeq[%v] FishType[%v] FishNumber[%v] FishCount[%v] ", data.GetFishID(), data.GetFishType(), fish.fishnumber, fish.fishcount)

}

//炮的編號改变
func (by *Bylogic) addRateChange(index int32, ud *byuserdata) {
	defer util.PanicErrStack()

	if ud.diedbosstype == Zuantoulongxia && index != 35 {
		//mlog.Error("打中磚頭不可以換砲彈")
		return
	}

	if index < 0 {
		return
	}

	boomInfo := by.getBoomInfo(by.boomIdOption[index])
	odd := boomInfo.boomOdds
	ud.setUserBoomNumber(boomInfo.boomNumber, boomInfo.boomNumber)
	if odd > 0 {
		ud.currodd = odd
	}

	msglog := fmt.Sprintf("frameid[%v] 炮的編號改變為[%v]號 砲倍率[%v], 改之前的砲倍率為[%v]", by.frameid, boomInfo.boomNumber, odd, ud.currodd)
	mlog.Debug(msglog)
	mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, ud.user.TableNo, msglog)

	rateinfo := new(netproto.BY_UserChangeRateS2C)
	rateinfo.SeatNO = proto.Int32(ud.user.SeatNo)
	rateinfo.RateAfterChange = proto.Int32(index + 1)
	by.framedata.UserChangeRate = append(by.framedata.UserChangeRate, rateinfo)
}

//添加锁定
func (by *Bylogic) addSuoDing(data *netproto.BY_SuoDingS2C) {
	by.framedata.SuoDing = append(by.framedata.SuoDing, data)
}

//添加用户破产消息
func (by *Bylogic) addUserPoChan(userid int32) {
	if p, ok := by.getPiSaiUserData(); ok {
		if p.Score > 0 {
			return
		}
	}

	by.framedata.PoChan = append(by.framedata.PoChan, userid)
}

//处理用户数据帧
func (by *Bylogic) userFrame(user *gamelogic.GameUser, data []byte) {
	//mylog.MongoLog(user.UserID, *by.TableManager.Room.Seq, by.userdata[0].user.TableNo, "C2S userFrame")
	defer util.PanicErrStack()

	ust := &netproto.BY_FrameC2S{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Warn("解析用户%v帧数据错误%v", user, err)
		return
	}

	ud := by.getuserdatabyuserid(user.UserID)
	//先处理用户改变倍率
	if ust.UserChangeRate != nil {
		by.addRateChange(ust.UserChangeRate.GetRateAfterChange()-1, ud)
	}

	//檢查是否有此炮彈編號
	if !by.canUseBoomNumber(ud) {
		mlog.Error("userId[%v] 不可使用炮彈 boomId[%v]", ud.user.UserID, ud.boomNumber)
		return
	}

	//用户施放冰冻技能
	if ust.UseIce != nil && ust.GetUseIce() {
		by.userCastIce(ud)
	}
	//处理用户射击并限制其最大炸弹数
	boomType := by.getBoomInfo(by.boomIdOption[ud.boomNumber-1]).boomType
	boomInfo := by.getBoomInfo(by.boomIdOption[ud.boomNumber-1])

	if ust.Shoot != nil && ud.boomcount < by.maxboom && boomType == int32(Basic) {
		by.basicShoot(boomInfo, ud, ust)
	} else if ust.Shoot != nil && boomType == int32(ZlxBoom) {
		example.invalid("frameid[%v] 磚頭龍蝦砲射擊 boomNumber[%v] ", by.frameid, (ud.boomNumber - 1))
		ud.userShoot(by.frameid, 0)
		by.addshoot(ust.Shoot, ud)
	} else {
		mlog.Debug(" frameid[%v] 砲彈超過上限 or nil userid[%v] ust.Shoot[%v] ud.boomcount[%v]", by.frameid, ud.user.UserID, ust.GetShoot(), ud.boomcount)
	}

	if ust.Boom != nil {
		//磚頭砲
		if ust.Boom.GetBoomType() == ZlxBoom && ust.Boom.GetDelete() != Delete {
			//leiLongfish, _ := by.fishmanager.config.getfishinfobytype(ud.bossfishid)
			//mlog.Debug("磚頭砲 boomtype[%v] delete[%v] fish[%v]  teshufishtype[%v] fishtype[%v] spmaxfishcount[%v] 倍數門檻[%v] ",
			//	ust.Boom.GetBoomType(), ust.Boom.GetDelete(), ust.Boom.GetFishID(), ud.teshufishtype, leiLongfish.bosstypeinfo.fishtype, leiLongfish.bosstypeinfo.spmaxfishcount, leiLongfish.bosstypeinfo.spmaxfishmulti)

			ust.Boom.FishID = by.leiLongOrZtxAllFishDied(ust.Boom.GetFishID(), ud)
			by.isAllFishDied(ust.Boom.GetFishID(), ud)
			return
		} else if ust.Boom.GetBoomType() == ZlxBoom && ust.Boom.GetDelete() == Delete {
			ud.fightspmaxfishcount = 0
			by.clearTeshufish(ud)
			//mlog.Debug("clear boomtype[%v] delete[%v] fish[%v] allfish[%v]", ust.Boom.GetBoomType(), ust.Boom.GetDelete(), ust.Boom.GetFishID(), ust.GetAllFish())
		}

		//验证炸弹ID是否是自己的
		//mlog.Debug("frameid[%v] 檢核打中魚 ust.Boom[%v]", by.frameid, ust.Boom)
		bfindboom := ud.removeBoom(ust.Boom.GetBoomID())

		if ust.Boom.UserID != nil {
			//机器人炸弹
			if by.GetUserByUserID(ust.Boom.GetUserID()) == nil {
				bfindboom = false
			} else {
				//其他玩家打到魚
				ud = by.getuserdatabyuserid(ust.Boom.GetUserID())
				bfindboom = ud.removeBoom(ust.Boom.GetBoomID())
				mlog.Debug("frameid[%v] 其他玩家打到魚 userid[%v] removeBoom[%v] boomcount[%v] ",
					by.frameid, ud.user.UserID, ust.Boom.GetBoomID(), ud.boomcount)
			}
		}

		if bfindboom {
			//炸弹存在检查是否捕中鱼
			boominfo := new(netproto.BY_BoomS2C)
			boominfo.BoomID = ust.Boom.BoomID
			by.addboom(boominfo)
			//每次只打到一条鱼
			fishid := make([]int32, 0)
			if ust.Boom == nil || len(ust.Boom.FishID) == 0 {
				mlog.Debug("沒有魚的資料 boom[%v] len[%v]", ust.Boom, len(ust.Boom.FishID))
				return
			}
			fishid = append(fishid, ust.Boom.FishID[0])

			logMsg := fmt.Sprintf("userid[%v] score[%v] 子彈boomid[%v]打到魚[%v] fishid[%v] 砲倍率[%v]",
				ud.user.UserID, ud.user.Score+ud.money, ust.Boom.GetBoomID(), ust.Boom.GetFishID(), fishid, ud.getUserOdd())
			mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), logMsg)
			mlog.Debug(logMsg)
			allfishid := []int32{}
			if ust.GetAllFish() != nil {
				allfishid = ust.GetAllFish().FishID
			}
			by.fish(fishid, ud, false, allfishid)
		} else {
			//mlog.Debug("其他玩家打到魚 ")
		}

	} else {
		//mlog.Debug("frameid[%v] boom is nil [%v] boomid[%v] boomType[%v]", by.frameid, ust.GetBoom(), ust.Boom.GetBoomID(), boomType)
	}

	//用户锁状态
	by.userSuoDing(ust, ud, user)

	//这里处理是否捕中所有的鱼，炸弹鱼，一网打尽，雷龙鱼，磚頭龍蝦
	by.catchAllFishAndTeshufish(ust, ud)

}

//一般子彈發射
func (by *Bylogic) basicShoot(boomInfo *byboommanager, ud *byuserdata, ust *netproto.BY_FrameC2S) {
	example.invalid("frameid[%v] 用戶userid[%v]射擊 砲彈編號[%v] 砲彈型態[%v] 砲倍率[%v] 場上留存砲彈數[%v]",
		by.frameid, ud.user.UserID, boomInfo.boomNumber, boomInfo.boomType, boomInfo.boomOdds, ud.boomcount)
	boomMoney := by.getBoomInfo(by.boomIdOption[ud.boomNumber-1]).boomMoney
	IsPoChan := !ud.userShoot(by.frameid, boomMoney)
	//当前的倍数是否破产
	IsPoChan = by.isPoChanRateChange(IsPoChan, ud)
	IsPoChan = by.piSaiCheckShoot(IsPoChan, boomMoney)
	if !IsPoChan {
		by.addshoot(ust.Shoot, ud)
	} else {
		by.addUserPoChan(ud.user.UserID)
		//mlog.System("用户%d已经破產", ud.user.UserID)
	}
}

//破產了,改變子彈倍數
func (by *Bylogic) isPoChanRateChange(IsPoChan bool, ud *byuserdata) bool {
	if IsPoChan {
		lastmoney := ud.user.Score + ud.money
		example.invalid("frameid[%v] 用户%d钱不够，切换倍率 剩餘金額[%v] score[%v] money[%v]", by.frameid, ud.user.UserID, lastmoney, ud.user.Score, ud.money)
		for i := ud.boomindex - 1; i > 0; i-- {
			//如果用户钱不够就减小一档倍率，直到能发出炮。如果发不出炮就破产
			boomInfo := by.getBoomInfo(by.boomIdOption[i])
			if (lastmoney - int64(boomInfo.boomOdds)) >= 0 {
				by.addRateChange(int32(i), ud)
				boomMoney := by.getBoomInfo(by.boomIdOption[ud.boomNumber-1]).boomMoney
				IsPoChan = !ud.userShoot(by.frameid, boomMoney)
				break
			}
		}
	}
	return IsPoChan
}

//这里处理是否捕中所有的鱼，炸弹鱼，一网打尽，雷龙鱼，磚頭龍蝦
func (by *Bylogic) catchAllFishAndTeshufish(ust *netproto.BY_FrameC2S, ud *byuserdata) {
	defer util.PanicErrStack()
	if ust.AllFish != nil && ud.teshufishtype != 0 {

		//mlog.System("frameid[%v] 是否捕中所有的鱼 allfishid[%v] teshufishtype[%v] 檢核打中魚[%v] ust[%+v]",
		//	by.frameid, ust.AllFish.GetFishID(), ud.teshufishtype, ust.Boom, ust)

		by.boomFishSetAllOdd(ud, ust)
		if ud.teshufishtype == LeiLong {
			ust.AllFish.FishID = by.leiLongOrZtxAllFishDied(ust.AllFish.FishID, ud)
		}

		by.isAllFishDied(ust.AllFish.FishID, ud)
		ud.fightspmaxfishcount = 0

		teshufishfish, ok := by.fishmanager.config.getfishinfobytype(ud.bossfishid)
		if !ok {
			return
		}
		//mlog.Debug("teshufishfish[%v] ud.diedbosstype[%v] fishid[%v] bossfishid[%v]", teshufishfish.bosstypeinfo.fishtype, ud.diedbosstype, teshufishfish.bosstypeinfo.id, ud.bossfishid)
		if teshufishfish.bosstypeinfo.fishtype == ud.diedbosstype {
			by.clearTeshufish(ud)
		}
	}

	if ud.diedbosstype == Dujiaojing {
		by.clearTeshufish(ud)
	}
}

func (by *Bylogic) isAllFishDied(fishID []int32, ud *byuserdata) {
	if len(fishID) == int(0) {
		//mlog.Debug("isAllFishDied is empty")
		return
	}

	by.fish(fishID, ud, true, nil)
}

//清除特殊魚
func (by *Bylogic) clearTeshufish(ud *byuserdata) {
	//mlog.Debug("清除特殊魚 teshufishtype[%v] teshufishfishodd[%v] allFishOdd[%v] diedbosstype[%v] ", ud.teshufishtype, ud.teshufishfishodd, ud.allFishOdd, ud.diedbosstype)
	ud.teshufishfishodd = 0
	ud.allFishOdd = 0
	ud.diedbosstype = 0
	ud.spmaxfishmulti = -1
	ud.setFishType(0)
	ud.setTeShuFishType(0)
}

//用户锁状态
func (by *Bylogic) userSuoDing(ust *netproto.BY_FrameC2S, ud *byuserdata, user *gamelogic.GameUser) {
	if ust.SuoDing != nil {

		ud.setSuoDing()
		suoding := new(netproto.BY_SuoDingS2C)
		suoding.SeatNO = proto.Int32(user.SeatNo)
		suoding.SuoDing = proto.Bool(ust.GetSuoDing().GetSuoDing())
		ud.suoding = ust.GetSuoDing().GetSuoDing()
		by.addSuoDing(suoding)

		logmsg := fmt.Sprintf("userid[%v] 切換瞄準狀態[%v]", user.UserID, ud.suoding)
		//mlog.Debug(logmsg)
		mylog.MongoLog(user.UserID, *by.TableManager.Room.Seq, ud.user.TableNo, logmsg)
	}
}

//閃電鰻魚死亡清單
func (by *Bylogic) leiLongOrZtxAllFishDied(fid []int32, ud *byuserdata) []int32 {
	defer util.PanicErrStack()
	fishlist := make([]int32, 0)

	teshufishfish, ok := by.fishmanager.config.getfishinfobytype(ud.bossfishid)

	if !ok {
		return fishlist
	}

	if teshufishfish.bosstype != ud.teshufishtype {
		mlog.Debug("非特殊魚 ")
		return fishlist
	}

	if len(fid) == int(0) {
		mlog.Error("閃電鰻魚死亡清單 is empty")
		return fishlist
	}

	for i := 0; i < len(fid); i++ {

		f := by.fishmanager.getFishInfoByFishID(fid[i])
		if f == nil {
			example.invalid("此类鱼未找到 fishSeq[%v]", fid[i])
			continue
		}
		fishcf, _ := by.fishmanager.config.getfishinfobytype(f.fishtype)
		//mlog.Debug("frameid[%v] fishid[%v] fishtype[%v] bosstype[%v]", by.frameid, fishcf.fishid, f.fishtype, fishcf.bosstype)
		if !ok || fishcf == nil {
			example.invalid("此类鱼未找到:%d", f.fishtype)
			continue
		}

		//不会炸死BOSS鱼和特殊鱼
		if fishcf.bosstype != 1 && ud.teshufishtype != FishAll {
			mlog.Debug("不会炸死BOSS鱼和特殊鱼 fishtype[%v]", f.fishtype)
			continue
		}

		//磚頭龍蝦不會炸死其他磚頭龍蝦
		if fishcf.bosstype == 6 && ud.teshufishtype == Zuantoulongxia {
			mlog.Debug("磚頭龍蝦不會炸死其他磚頭龍蝦 fishtype[%v]", f.fishtype)
			continue
		}

		//处理一网打尽
		if ud.teshufishtype == FishAll && f.fishtype != ud.fishtype {
			mlog.Debug("处理一网打尽 continue")
			continue
		}

		//处理北冥神鯤
		if fishcf.bosstype != 1 && ud.teshufishtype == Dujiaojing {
			mlog.Debug("处理北冥神鯤 continue")
			continue
		}

		_, _, fishOdd := by.fishmanager.getFishLuckyOdds(f)
		if fishOdd > ud.spmaxfishmulti {
			mlog.Debug("fishid[%v] 特殊於不能打中%v倍以上的魚 [%v]", fishcf.fishid, ud.spmaxfishmulti, fishOdd)
			continue
		}

		if int32(ud.fightspmaxfishcount) > teshufishfish.bosstypeinfo.spmaxfishcount || ud.teshufishfishodd <= int32(0) {
			logMsg := fmt.Sprintf("userid[%v] 特殊魚%v攻擊無效, 已攻擊%v隻 ", ud.user.UserID, teshufishfish.bosstype, ud.fightspmaxfishcount)
			mlog.Debug(logMsg)
			mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), logMsg)
			break
		}

		n, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", rand.Float64()), 64)

		//mlog.Debug("bossfishid[%v] canDiedRate[%v] rand[%v] spdeathrate[%v] fishid[%v]",
		//	ud.bossfishid, ud.candiedrate, n, teshufishfish.bosstypeinfo.spdeathrate, fishcf.fishid)

		if ud.candiedrate > float32(n) {
			fishOdd := fishcf.rewardinfo[0].reward

			mlog.Debug("count[%v] 閃電鰻魚 or 磚頭蝦 死亡清單 rand[%v] canDiedRate[%v] fishid[%v] fishodd[%v] leiLongfishOdd[%v] ",
				ud.fightspmaxfishcount, n, ud.candiedrate, fishcf.fishid, fishOdd, ud.teshufishfishodd)
			mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), fmt.Sprintf("userid[%v] 抓到第[%v]條魚 閃電鰻魚or磚頭蝦攻擊 死亡清單 亂數隨機[%v] 死亡機率[%v] fishid[%v] 魚倍率[%v] 特殊魚倍率[%v] ",
				ud.user.UserID, ud.fightspmaxfishcount, n, ud.candiedrate, fishcf.fishid, fishOdd, ud.teshufishfishodd))

			fishlist = append(fishlist, int32(fid[i]))
			ud.teshufishfishodd -= fishOdd
			ud.candiedrate = float32(teshufishfish.bosstypeinfo.spdeathrate)
		} else {
			ud.candiedrate += float32(teshufishfish.bosstypeinfo.rateincrease)
			if teshufishfish.bosstypeinfo.increasemax != 0 && ud.candiedrate > float32(teshufishfish.bosstypeinfo.increasemax) {
				ud.candiedrate = float32(teshufishfish.bosstypeinfo.increasemax)
			}
			mlog.Debug("增加攻擊[%v] sum[%v]", float32(teshufishfish.bosstypeinfo.rateincrease), ud.candiedrate)
		}

		ud.fightspmaxfishcount++
		//mlog.Debug("fightspmaxfishcount[%v] spmaxfishcount[%v]", ud.fightspmaxfishcount, teshufishfish.bosstypeinfo.spmaxfishcount)
	}

	return fishlist
}

//炸彈蟹
func (by *Bylogic) boomFishSetAllOdd(ud *byuserdata, ust *netproto.BY_FrameC2S) {
	if ud.teshufishtype == BoomFish {
		if len(ust.AllFish.FishID) == 0 {
			mlog.Error("炸彈蟹 allfish is empty")
		}

		for _, v := range ust.AllFish.FishID {
			if by.fishmanager.findFish(v) {
				fishcf, _, fishOdd := by.fishmanager.getFishLuckyOdds(by.fishmanager.fishdata[v])
				ud.allFishOdd += fishOdd

				mlog.Debug("炸彈蟹 ud.teshufishtype[%v] allFishOdd[%v] fishid[%v] fishodd[%v] fishtype[%v] Lucky[%v]",
					ud.teshufishtype, ud.allFishOdd, v, by.fishmanager.fishdata[v].fishodd, by.fishmanager.fishdata[v].fishtype, fishcf.luckybonus)
				if ud.allFishOdd == 0 {
					mlog.Error("炸彈蟹打中所有魚但倍率是0. QueryMessage[%v] fishid[%v] odd[%v]", by.fishmanager.fishdata[v].fishSeq, by.fishmanager.fishdata[v].fishtype, by.fishmanager.fishdata[v].fishodd)
				}
			}
		}
	}
}

//获取一个炸弹ID
func (by *Bylogic) getBoomId(seatno int32) int32 {
	if by.bid[seatno] >= (seatno+1)*10000000 {
		by.bid[seatno] = seatno * 10000000
	}

	by.bid[seatno]++

	return by.bid[seatno]
}

//通过用户ID获取用户数据
func (by *Bylogic) getuserdatabyuserid(uid int32) *byuserdata {
	defer util.PanicErrStack()

	for _, ud := range by.userdata {
		if ud.user.UserID == uid {
			return ud
		}
	}

	user := by.GetUserByUserID(uid)
	if user == nil {
		mlog.Warn("用户%d数据为空", uid)
	}
	return by.adduserdata(user)
}

//添加用户数据
func (by *Bylogic) adduserdata(user *gamelogic.GameUser) *byuserdata {
	defer util.PanicErrStack()
	ud := newuserdata(user)
	by.userdata = append(by.userdata, ud)
	boomInfo := by.getBoomInfo(by.boomIdOption[0])
	ud.setUserBoomNumber(boomInfo.boomNumber, 0)
	win, _ := strconv.ParseInt(user.GameBuff, 10, 64)
	ud.allwin = int64(win)
	return ud
}

//删除用户数据
func (by *Bylogic) deluserdatabyuserid(uid int32) {
	defer util.PanicErrStack()

	for i := 0; i < len(by.userdata); i++ {
		if by.userdata[i].user.UserID == uid {
			//删除以前更新一次分数
			if by.flag == ExperienceCashFlag && !by.userdata[i].user.IsRobot() {
				by.UpdateExperienceScoreFunc(by.userdata[i].user.UserID, by.userdata[i].money, 0)
			} else if by.TableManager.Room.GameType == gamelogic.PiSai {

				if p, ok := by.getPiSaiUserData(); ok {
					totalPlayTime := (p.EndTime - time.Now().UnixNano()) / int64(time.Second)
					if !p.IsPlaying || totalPlayTime > 0 {
						by.updateContestUser()
						//mlog.System("未完成比賽不能清除 userid[%v]", p.UserId)
						p.Pause = true
						by.userdata = make([]*byuserdata, 0)
						return
					} else if !p.IsPlaying {
						mlog.System("删除用户数据 pisai userid[%v]", p.UserID)
						delete(by.TableManager.Room.OnPiSaiUserInfo, p.UserID)
					}
				}
			} else {
				by.UpdateUserScoreFunc(by.userdata[i].user.UserID, by.userdata[i].money, 0)
			}

			by.userdata = append(by.userdata[:i], by.userdata[i+1:]...)

			if _, ok := by.TableManager.Room.FisnNewPlayerDetail[uid]; ok {
				//	mlog.System("刪除FisnNewPlayerDetail userid[%v]", uid)
				delete(by.TableManager.Room.FisnNewPlayerDetail, uid)
			}
			break
		}
	}
}

//用户进入
func (by *Bylogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {

	logMsg := fmt.Sprintf("用户进入UserEnter userid[%v] IsRobot[%v]", gameuser.UserID, gameuser.IsRobot())
	//mlog.System(logMsg)
	mylog.MongoLog(gameuser.UserID, *by.TableManager.Room.Seq, by.TableNo(), logMsg)

	by.ComputeGameLogicBase.UserEnter(gameuser, tableno, seatno)
	by.findNewPlayerPool(gameuser.UserID)

	by.urs[gameuser.UserID] = new(gamelog.UserData)
	by.urs[gameuser.UserID].Score = gameuser.Score
	by.urs[gameuser.UserID].TradeID = gameuser.TradeID

	by.GameLog.UpdateGameUsers(by.urs)

	//檢查玩家是否有買票
	//by.checkUserTicket()

	by.getContentInfo()
	setmsg := by.getSettingMsg()
	by.SendMessageToUser(gameuser, setmsg)

	//所有用户都需要在进入的时候发送场景数据
	by.sendSence(gameuser)

}

//增減系統水池
func (by *Bylogic) adjustSyspool(ud *byuserdata, money int64) {

	if by.TableManager.Room.FishSysPoolData == nil {
		mlog.Warn("增減系統水池FishSysPoolData 是空的nil. 請確認異常.")
		return
	}

	by.TableManager.Room.FishSysPoolData.AfterAdjustSyspoolWinLoss = proto.Int64(by.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() + int64(money))
	by.TableManager.Room.FishSysPoolData.PoolAmount = proto.Int64(by.TableManager.Room.FishSysPoolData.GetPoolAmount() + int64(money))

	logMsg := fmt.Sprintf("frameid[%v] 增減系統水池 調整前系統總輸贏[%v] 變動金額[%v] 調整後系統總輸贏[%v] ",
		*by.TableManager.Room.Seq,
		-1*by.TableManager.Room.GetSysSumWinLoss(), -money, -by.TableManager.Room.GetAfterAdjustSyspoolWinLoss())
	mlog.Debug(logMsg)
	mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, ud.user.TableNo, logMsg)

	by.TableManager.Room.FishSysPoolData.SysSumWinLoss = proto.Int64(by.TableManager.Room.GetAfterAdjustSyspoolWinLoss())
}

func (by *Bylogic) findNewPlayerPool(userid int32) {
	ud := by.getuserdatabyuserid(userid)
	if ud.user.SumCharge > 0 || ud.user.IsRobot() {
		return
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_GetBeginnerOfferID)
	pdata := new(netproto.FindFishBeginnerOfferByUserID)
	pdata.UserID = proto.Int32(userid)
	pdata.NewPlayerAdjustSyspool = proto.Int32(by.newplayeradjustsyspool)
	pmsg.MsgData = pdata

	by.SendRPCMessageFunc(false, pmsg)
}

//发送场景消息
func (by *Bylogic) sendSence(gameuser *gamelogic.GameUser) {
	defer util.PanicErrStack()
	cloneClearFrame := by.clearFrameData()

	if by.flag == ExperienceCashFlag && !gameuser.IsRobot() {
		gameuser.Score = gameuser.ExperienceCash
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BY_GameMessageClassID_UserRecoverRetID)
	pdata := new(netproto.BY_UserRecoverRet)

	pdata.FrameInfo = cloneClearFrame

	pdata.Isfishwave = proto.Bool(by.fishmanager.fishwave)
	pdata.Fishwavetip = proto.String(by.fishwavetip)

	for _, v := range pdata.FrameInfo {
		v.Shoot = make([]*netproto.BY_ShootS2C, 0)
	}

	//添加用户状态
	for _, ud := range by.userdata {
		us := new(netproto.BY_UserStatus)
		us.UserID = proto.Int32(ud.user.UserID)
		//客户端索引是从1开始的
		//boomMoney := by.boommanager[ud.boomNumber-1].boomMoney
		us.Rate = proto.Int32(ud.boomNumber)
		us.Money = proto.Int64(ud.lastframemoney)
		IsPoChan := false

		//当钱和炮弹数量都为0的时候才破产
		mlog.Debug("当钱和炮弹数量都为0的时候才破产 score[%v] money[%v]", ud.user.Score, ud.money)
		if (ud.user.Score + ud.money) == 0 {
			IsPoChan = true
		}
		if p, ok := by.getPiSaiUserData(); ok {
			IsPoChan = false
			if p.IsGameStart {
				if p.Score <= 0 && p.IsPlaying {
					msglog := fmt.Sprintf("userid[%v] 破產", ud.user.UserID)
					mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, 0, msglog)
					IsPoChan = true
				}
			}
		}

		mlog.Debug("UserRecoverRetID 发给用户%d，用户%d当前的钱为：%d，炮弹数目为：%d", gameuser.UserID, ud.user.UserID, ud.user.Score+ud.money, ud.boomcount)
		us.PoChan = proto.Bool(IsPoChan)
		pdata.UserStatus = append(pdata.UserStatus, us)
	}

	pmsg.MsgData = pdata
	by.SendMessageToUser(gameuser, pmsg)
}

//清理帧数据
func (by *Bylogic) clearFrameData() []*netproto.BY_FrameS2C {
	defer util.PanicErrStack()

	cloneFrame := make([]*netproto.BY_FrameS2C, len(by.framemsg))
	copy(cloneFrame, by.framemsg)
	framecount := int32(len(cloneFrame))

	lenlimit := framecount
	if framecount > by.framecount {
		by.framemsg = append(by.framemsg[:0], by.framemsg[framecount-by.framecount:]...)
		lenlimit = by.framecount
	}

	newFrameMsg := make([]*netproto.BY_FrameS2C, 0)
	checkIsIceStart := true
	checkIsIceEnd := false
	for k := len(cloneFrame) - int(lenlimit); k < len(cloneFrame); k++ {
		if checkIsIceStart && cloneFrame[k].GetIsIce() {
			checkIsIceStart = false
			checkIsIceEnd = true
			cloneFrame[k].IsIce = proto.Bool(true)
			cloneFrame[k].UserInfo = make([]*netproto.BY_UserInfo, 0)
			newFrameMsg = append(newFrameMsg, cloneFrame[k:k+1]...)
			mlog.Debug("UserEnter 冰凍開始 frameid[%v] isice[%v] newFrameMsg[%v]",
				cloneFrame[k].GetFrameID(), cloneFrame[k].GetIsIce(), cloneFrame[k:k+1])
			continue
		}
		if checkIsIceEnd && !cloneFrame[k].GetIsIce() {
			checkIsIceEnd = false
			cloneFrame[k].IsIce = proto.Bool(false)
			cloneFrame[k].UserInfo = make([]*netproto.BY_UserInfo, 0)
			newFrameMsg = append(newFrameMsg, cloneFrame[k:k+1]...)
			mlog.Debug("UserEnter 冰凍結束 frameid[%v] isice[%v] newframe[%v]",
				cloneFrame[k].GetFrameID(), cloneFrame[k].GetIsIce(), cloneFrame[k:k+1])
			checkIsIceStart = true
			continue
		}

		for _, f := range cloneFrame[k].GetNewFish() {
			fishcf, ok := by.fishmanager.config.getfishinfobytype(f.GetFishType())
			if ok && fishcf != nil && fishcf.rewardinfo[0].reward > by.framefishodd {
				ok2 := by.fishmanager.findFish(f.GetFishID())
				if ok2 {
					cloneFrame[k].UserInfo = make([]*netproto.BY_UserInfo, 0)
					newFrameMsg = append(newFrameMsg, cloneFrame[k:k+1]...)
					mlog.Debug("UserEnter frameid[%v] fishtype[%v] fishodd[%v] newfish[%v]", cloneFrame[k].GetFrameID(), f.GetFishType(),
						fishcf.rewardinfo[0].reward, cloneFrame[k].GetNewFish())
				}
			}
		}

		if k == int(lenlimit)-1 {
			cloneFrame[k].UserInfo = make([]*netproto.BY_UserInfo, 0)
			newFrameMsg = append(newFrameMsg, cloneFrame[k:k+1]...)
			mlog.Debug("最後一筆要給client len[%v]", len(newFrameMsg))
		}

	}

	return newFrameMsg
}

func (by *Bylogic) delShootByBoomId(s []*netproto.BY_ShootS2C, id int32) bool {
	for i := 0; i < len(s); i++ {
		if id == s[i].GetBoomID() {
			s = append(s[:i], s[i+1:]...)
			return true
		}
	}

	return false
}

//炮弹ID是否存在
func (by *Bylogic) boomIsFind(boomid int32, UserID int32) bool {
	user := by.GetUserByUserID(UserID)
	if user == nil {
		return false
	}

	ud := by.getuserdatabyuserid(UserID)
	return ud.findBoomById(boomid)
}

//断线重连
func (by *Bylogic) UserRecover(gameuser *gamelogic.GameUser) {
	mlog.System("用户断线重连")
	by.ComputeGameLogicBase.UserRecover(gameuser)
	//断线重连的时候清理一下用户的子弹
	ud := by.getuserdatabyuserid(gameuser.UserID)
	ud.delAllBoom()

	if p, ok := by.getPiSaiUserData(); ok {
		//mlog.System("recover time[%v]", p.EndTime)
		totalPlayTime := (p.EndTime - time.Now().UnixNano()) / int64(time.Second)
		if totalPlayTime <= -15 {
			by.gameEndKickUser()
		}
	}
}

//用户离开
func (by *Bylogic) UserLeave(gameuser *gamelogic.GameUser) bool {
	gametime := time.Now().Sub(gameuser.LoginTime)
	gameuser.GameTime = int32(gametime.Seconds())

	uscoreRet := make(map[int32]int64)
	var buffer bytes.Buffer
	for _, urd := range by.userdata {
		if urd.user.UserID == gameuser.UserID {
			uscoreRet[urd.user.UserID] = urd.winmoney
			by.GameLog.Context().UserUpdateScore(urd.user.UserID, urd.winmoney)

			str := fmt.Sprintf("UserId[%v] 打中魚合計：", urd.user.UserID)
			buffer.WriteString(str)
			for k, f := range urd.fishwinloss {
				fishSum := fmt.Sprintf("FishId[%v]Money[%v]; ", k, f)
				buffer.WriteString(fishSum)
			}
		}
	}

	//by.GameLog.AddGameStep(buffer.String())
	by.GameLog.RemoveGameUsers(gameuser.UserID)
	by.GameLog.GameEnd()
	//by.GameLog.AddGameStep("----------------------------------------------------------------")

	by.updateUserSorce()

	by.newPlayerLeave(gameuser)

	if by.TableManager.Room.GameType == gamelogic.PiSai {
		if by.piSaiLeaveTimer != nil {

			msglog := fmt.Sprintf("by.piSaiLeaveTimer.Stop() userid[%v]", gameuser.UserID)
			//mlog.Error(msglog)
			mylog.MongoLog(gameuser.UserID, *by.TableManager.Room.Seq, gameuser.TableNo, msglog)

			by.piSaiLeaveTimer.Stop()
			by.piSaiLeaveTimer = nil
		}
	}

	if gameuser == nil {
		return true
	}
	by.deluserdatabyuserid(gameuser.UserID)
	by.ChangeUserStateEnd(gameuser)
	by.bid[gameuser.SeatNo] = gameuser.SeatNo * 10000000

	if by.TableManager.Room.GameType == gamelogic.PiSai {
		if p, okPS := by.getPiSaiUserData(); okPS && !p.IsGameStart {
			by.pisaiCancelTimer()
		}
	}

	msglog := fmt.Sprintf("玩家離開 userid[%v] IsRobot[%v] Score[%v]", gameuser.UserID, gameuser.IsRobot(), gameuser.Score)
	mlog.Debug(msglog)
	mylog.MongoLog(gameuser.UserID, *by.TableManager.Room.Seq, gameuser.TableNo, msglog)

	return by.ComputeGameLogicBase.UserLeave(gameuser)
}

//新手優惠玩家離開水池要回調
func (by *Bylogic) newPlayerLeave(gameuser *gamelogic.GameUser) {
	_, ok := by.TableManager.Room.NewPlayerPoolRecord[gameuser.UserID]
	ud := by.getuserdatabyuserid(gameuser.UserID)

	if ud.user.IsRobot() {
		return
	}

	if ok && gameuser.LoginScore+gameuser.UpdateScore <= int64(by.newplayeradjustsyspool) {
		//新手玩家離開加回剩餘身上的錢到水池
		by.adjustSyspool(ud, int64(-gameuser.Score))
		delete(by.TableManager.Room.NewPlayerPoolRecord, gameuser.UserID)
	}
	if ok && gameuser.LoginScore+gameuser.UpdateScore > int64(by.newplayeradjustsyspool) {
		//新手玩家離開加回剩餘身上的錢到水池
		by.adjustSyspool(ud, int64(-by.newplayeradjustsyspool))
		delete(by.TableManager.Room.NewPlayerPoolRecord, gameuser.UserID)
	}
}

//用户捕鱼
func (by *Bylogic) fish(fid []int32, ud *byuserdata, fishall bool, allfishid []int32) {
	defer util.PanicErrStack()

	for i := 0; i < len(fid); i++ {
		f := by.fishmanager.getFishInfoByFishID(fid[i])

		if f == nil {
			mlog.Debug("此类鱼未找到 fishSeq[%v]", fid[i])
			continue
		}

		fishcf, _ := by.fishmanager.config.getfishinfobytype(f.fishtype)

		if fishcf == nil {
			//by.calEmptyBoomFishAmount(-boomMoney)
			mlog.Debug("此类鱼未找到:%d", f.fishtype)
			return
		}

		//不会炸死BOSS鱼和特殊鱼
		if fishall && fishcf.bosstype != 1 && ud.teshufishtype != FishAll {
			_, _, fishOdd := by.fishmanager.getFishLuckyOdds(f)
			if ud.spmaxfishmulti != -1 && fishOdd <= ud.spmaxfishmulti {

			} else if ud.spmaxfishmulti != -1 && fishOdd > ud.spmaxfishmulti {
				logmsg := fmt.Sprintf("不会炸死BOSS鱼和特殊鱼 ,fishid[%v][%v] fishtype[%v] bosstype[%v] 倍率[%v] 超過倍數門檻[%v]不能打死",
					fishcf.fishid, fishcf.name, f.fishtype, fishcf.bosstype, fishOdd, ud.spmaxfishmulti)
				mlog.Debug(logmsg)
				mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, ud.user.TableNo, logmsg)
				continue
			} else {
				logmsg := fmt.Sprintf("不会炸死BOSS鱼和特殊鱼 ,fishid[%v][%v] fishtype[%v] bosstype[%v]  ",
					fishcf.fishid, fishcf.name, f.fishtype, fishcf.bosstype)
				mlog.Debug(logmsg)
				mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, ud.user.TableNo, logmsg)
				continue
			}
		}

		//磚頭龍蝦不會炸死其他磚頭龍蝦
		if fishall && fishcf.bosstype == 6 && ud.teshufishtype == Zuantoulongxia {
			mlog.Debug("磚頭龍蝦不會炸死其他磚頭龍蝦 fishtype[%v]", f.fishtype)
			continue
		}

		//处理一网打尽
		if fishall && ud.teshufishtype == FishAll && f.fishtype != ud.fishtype {
			mlog.Debug("处理一网打尽 continue")
			continue
		}

		//处理北冥神鯤
		if fishall && ud.teshufishtype == Dujiaojing {
			mlog.Debug("处理北冥神鯤 continue")
			continue
		}

		//打中魚的攻擊威力
		fishdata, rate := by.fightFishPowerRate(fid, i, ud, fishall)

		if f.fishtype == 41 {
			mlog.Debug("打到核彈頭 frameid[%v] fishid[%v] fishtype[%v] bosstype[%v] allfish[%v]", by.frameid, fishcf.fishid, f.fishtype, fishcf.bosstype, allfishid)
		}
		if f.fishtype == 42 {
			mlog.Debug("打到昆 frameid[%v] fishid[%v] rate[%v]", by.frameid, fishcf.fishid, rate)
		}

		ok, money := by.fishmanager.canfish(*by.TableManager.Room.Seq, fid[i], rate, fishall, ud, by.basescore, allfishid)

		if fishall && ok {
			mlog.Debug("frameid[%v] 打到allfish 魚編號[%v] money[%v]", by.frameid, f.fishtype, money)
		}

		by.AddGameLog(ok, fishcf, ud, money)

		if ok {
			//如果鱼被打完了，提前结束鱼潮
			if by.fishmanager.fishwave && by.fishmanager.checkFishWaveEnd() {
				by.endFishWave()
			}
			mlog.Debug("捕魚前 money:%v currodd[%v] boomindex[%v] boomdata[%v]", money, ud.currodd, ud.boomindex, ud.boomdata)
			if money <= 0 && fishdata.fishtype != 31 && fishdata.fishtype != 32 && fishdata.fishtype != 36 && fishdata.fishtype != 41 && fishdata.fishtype != 42 {
				mlog.Error("非特殊魚死亡卻沒有金額, fishid[%v]捕魚前 money:%v currodd[%v] boomindex[%v] boomdata[%v]", fishdata.fishtype, money, ud.currodd, ud.boomindex, ud.boomdata)
				mylog.MongoLog(ud.user.UserID, by.frameid, by.TableNo(), fmt.Sprintf("非特殊魚死亡卻沒有金額, fishid[%v]捕魚前 money:%v currodd[%v] boomindex[%v] boomdata[%v]", fishdata.fishtype, money, ud.currodd, ud.boomindex, ud.boomdata))
			}

			money = money * ud.getUserOdd()

			fish := new(netproto.BY_Fish)
			if fishcf.bosstype == Dujiaojing {
				fish.Type = proto.Int32(Dujiaojing)
				fish.TeshuFishID = proto.Int32(ud.bossfishid)

				money = by.addMoneyShow(fishcf, ud.getUserOdd())
				mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, ud.user.TableNo, fmt.Sprintf("userid[%v] 打中北冥神鯤 獲得總獎金[%v], FishWinList[%v]",
					ud.user.UserID, money, by.framedata.FishWinList))
			}

			logMsg := fmt.Sprintf("frameid[%v] userid[%v] 抓到魚後 fishid[%v] lucky[%v] 更新後金幣為[%v] 更新金幣前的金幣[%v] 更新金幣[%v] 砲倍率[%v] 砲攻擊率[%v]",
				by.frameid, ud.user.UserID, fishcf.name, fishcf.luckybonus, (int32(ud.user.Score) + money), ud.user.Score, money, ud.getUserOdd(), ud.rate)
			mlog.Debug(logMsg)
			mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), logMsg)
			mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), fmt.Sprintf("frameid[%v] userid[%v] 抓到魚後 fishid[%v] money[%v] 更新後金幣為[%v] 更新金幣前的金幣[%v] 更新金幣[%v] 砲倍率[%v] 砲攻擊率[%v]",
				by.frameid, ud.user.UserID, fid[i], money, (int32(by.updateMoney)+money), ud.money, by.updateMoney, ud.getUserOdd(), ud.rate))

			ud.updateUserMoney(money)
			boomMoney := by.getBoomInfo(by.boomIdOption[ud.boomNumber-1]).boomMoney
			//fmt.Printf("boomMoney %d, money %d\n", boomMoney, money)
			flowamount := money - boomMoney
			//ud.updateFlowAmount(flowamount)
			ud.updateBetAmount(flowamount, boomMoney, money)
			by.updatePiSaiUserBonus(money)
			by.updateSysPoolWinLoss(ud, money)
			by.recardPlayerPool(money, ud)

			logMsg = fmt.Sprintf("frameid=%d:用户%d的錢%d，本次加錢:%d, 打中%v lucky[%v], 子彈倍率[%v] 子彈威力[%v]",
				by.frameid, ud.user.UserID, ud.user.Score+ud.money, money, fishcf.name, fishcf.luckybonus, ud.getUserOdd(), rate)
			mlog.Debug(logMsg)
			mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), logMsg)

			mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), fmt.Sprintf("frameid=%d:用户%d的钱%d，本次加钱:%d", by.frameid, ud.user.UserID, ud.user.Score+ud.money, money))
			fish.FishID = proto.Int32(fid[i])
			fish.UserID = proto.Int32(ud.user.UserID)
			fish.Money = proto.Int64(int64(money))
			fish.LuckyBonus = proto.Int32(fishcf.luckybonus)
			//告诉客户端鱼因为什么类型的鱼死亡，客户端特效
			if fishall {
				fish.Type = proto.Int32(ud.teshufishtype)
				fish.TeshuFishID = proto.Int32(ud.bossfishid)
			}

			by.addfish(fish)
			by.notifyaddfish(fishcf, ud, money)

			//机器人不会捕中特殊鱼，服务器上无法对机器人捕中特殊鱼进行处理
			//boss类型为3的是特殊鱼
			by.realUserSetTeShuFish(fishcf, ud, f)
		} else {
			by.addFishSwell(fid[i])
			by.recardPlayerPool(money, ud)
			boomMoney := by.getBoomInfo(by.boomIdOption[ud.boomNumber-1]).boomMoney
			//fmt.Printf("boomMoney %d, money %d\n", boomMoney, money)
			ud.updateBetAmount(boomMoney, boomMoney, money)
			//ud.updateFlowAmount(boomMoney)
		}

		by.updateSysReword(ok, fishdata, ud, money)
		//更新新手優惠資料
		by.updateNewPlayer(ud)

		//s := (int32(by.updateMoney) + money)
		//mlog.Debug("捕魚後2 fishid[%v] money:%v after[%v] beforeMoney[%v] updateMoney[%v] boomOdd[%v] ud.rate[%v]",
		//	fid[i], ud.money, s, ud.money, by.updateMoney, ud.getUserOdd(), ud.rate)
	}

}

//更新新手優惠資料
func (by *Bylogic) updateNewPlayer(ud *byuserdata) {
	if ud.user.IsRobot() || by.TableManager.Room.GameType == gamelogic.PiSai {
		return
	}

	type LogData struct {
		Seq       int64
		GroupName string
		FrameId   int64
		UserId    int32 //此玩家賬號ID
		Bullet    int32 //子彈發數
		Score     int32 //輸贏
	}

	playerReward, _ := by.TableManager.Room.FishPlayerRewardDetail[ud.user.UserID]
	//玩家淨贏分=玩家擊殺魚總得分-玩家總壓分
	playerWinLoss := int32((playerReward.GetHitMoney() - playerReward.GetSpendMoney()))
	//mlog.Debug("playerWinLoss[%playerReward] GetHitMoney[%v] GetSpendMoney[%v]", playerWinLoss, playerReward.GetHitMoney(), playerReward.GetSpendMoney())
	if p, ok := by.TableManager.Room.FisnNewPlayerDetail[ud.user.UserID]; ok {
		if p.GetStatus() == 1 || p.GetBullet() > by.newPlayerTable.offerLimitCount { //已經送完新手優惠
			return
		}

		p.Bullet = proto.Int32(p.GetBullet() + 1)
		p.Score = proto.Int32(playerWinLoss)

		msg := LogData{
			Seq:       *by.TableManager.Room.Seq,
			GroupName: "newPlayer",
			FrameId:   by.frameid,
			UserId:    ud.user.UserID,
			Bullet:    p.GetBullet(),
			Score:     p.GetScore(),
		}
		logMsg := fmt.Sprintf("新手優惠 [%+v]", msg)

		mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), logMsg)

	} else {
		by.TableManager.Room.FisnNewPlayerDetail[ud.user.UserID] = &netproto.FishBeginnerOffer{
			UserID: proto.Int32(ud.user.UserID),
			Bullet: proto.Int32(1),
			Score:  proto.Int32(playerWinLoss),
			Status: proto.Int32(0),
		}
		//mlog.Debug("NewPlayer[%+v]", by.TableManager.Room.FisnNewPlayerDetail[ud.user.UserID])
	}

}

//真人玩家打中特殊魚
func (by *Bylogic) realUserSetTeShuFish(fishcf *fishinfo, ud *byuserdata, f *byfishdata) {
	if fishcf.bosstype >= BoomFish && !ud.user.IsRobot() {
		if f != nil {
			if fishcf.bosstype == BoomFish {
				mlog.Debug("炸弹鱼")
				ud.setTeShuFishType(BoomFish)
			} else if fishcf.bosstype == FishAll {
				mlog.Debug("一网打尽鱼")
				ud.setTeShuFishType(FishAll)
				ud.setFishType(f.fishtype)
			} else if fishcf.bosstype == LeiLong {
				mlog.Debug("雷龙鱼")
				ud.setTeShuFishType(LeiLong)
			} else if fishcf.bosstype == Zuantoulongxia {
				mlog.Debug("用户捕鱼 磚頭龍蝦砲")
				ud.setTeShuFishType(Zuantoulongxia)
			}
		}
	}
}

//打中魚的攻擊威力
func (by *Bylogic) fightFishPowerRate(fid []int32, i int, ud *byuserdata, fishall bool) (*byfishdata, float32) {
	rate := float32(1)
	var fishdata *byfishdata
	var ok bool
	if fishdata, ok = by.fishmanager.fishdata[fid[i]]; ok {
		rate = by.getRate(fishdata.fishtype, ud)

		if !ud.user.IsRobot() {

			by.adjustSyspoolWin(ud)

			rate = by.adjectSysPoolRate(rate, ud)
			//mlog.Debug("系統adjectSysPoolRate[%v] ", rate)

			rate = by.adjectPlayerPoolRate(rate, ud)
			//mlog.Debug("個人adjectPlayerPoolRate[%v] ", rate)
			if rate < 0.1 {
				mlog.Error("rate err rate[%v]", rate)
			}

			rate = by.adjectNewPlayerRate(rate, ud)
			mlog.Debug("新手優惠 done Rate[%v] ", rate)

			rate = by.adjustFreePlayerRate(rate, ud, fishdata)
			//mlog.Debug("未充值超過限制金額 rate[%v]", rate)

			if by.TableManager.Room.Seq == nil {
				by.TableManager.Room.Seq = proto.Int64(0)
			}
			by.TableManager.Room.Seq = proto.Int64(*by.TableManager.Room.Seq + 1)
		} else {
			//mlog.Debug("IsRobot rate[%v]", rate)
		}
		//rate = util.Decimal(float64(rate), 3)
		//mlog.System("frameid[%v] 打中魚的攻擊率[%v]", by.frameid, rate)
	}

	if fishall && ud.diedbosstype == BoomFish && !ud.user.IsRobot() {
		rate = by.getBoomFishRate(ud)
		rate = rate / float32(by.basescore)
		//mlog.System("炸彈蟹 fishdiedRate[%v]  teshufishtype[%v] allFishOdd[%v] teshufishfishodd[%v] ",
		//	rate, ud.teshufishtype, ud.allFishOdd, ud.teshufishfishodd )
	}

	rate = rate * float32(by.basescore)
	return fishdata, rate
}

//系統水池贏錢抽水
func (by *Bylogic) adjustSyspoolWin(ud *byuserdata) {
	if ud.user.IsRobot() {
		return
	}
	modNum := by.TableManager.Room.FishSysPoolData.GetSumBullet() % int64(*by.systemWinCommission.BulletNumber) //檢核子彈發數
	sysPool := by.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() * -1
	if modNum == 0 && sysPool > 0 {
		//mlog.System("系統水池贏錢抽水 GetSysSumWinLoss[%v]", by.TableManager.Room.FishSysPoolData.GetSysSumWinLoss()*-1/100)
		adjust := float32(0)
		for i := 0; i < 10; i++ {
			if sysPool <= int64(*by.systemWinCommission.CoinDivision) || (adjust+*by.systemWinCommission.IncreaseRate) >= (*by.systemWinCommission.MaxRate) {
				by.TableManager.Room.HsiTungHuoLiIRate = adjust + *by.systemWinCommission.IniCommissionRate

				if by.TableManager.Room.HsiTungHuoLiIRate > *by.systemWinCommission.MaxRate {
					by.TableManager.Room.HsiTungHuoLiIRate = *by.systemWinCommission.MaxRate
				}

				adjustSysRate := float32(by.TableManager.Room.FishSysPoolData.GetSysSumWinLoss()) * (by.TableManager.Room.HsiTungHuoLiIRate)

				adjustSyspool := (int32(util.Decimal(float64(adjustSysRate), 1)*-1/100 + 0.6)) * 100
				if adjustSyspool == 0 {
					adjustSyspool = 1 * by.basescore
				}

				by.adjustSyspool(ud, int64(adjustSyspool))

				by.TableManager.Room.HsiTungHuoLiChouShuiChingFen = int32(adjustSyspool / 100)

				if by.TableManager.Room.GetSysSumWinLoss() > 0 && adjustSyspool > 0 {
					mlog.Warn("系統輸錢不收水, SysSum[%v] adjustSyspool[%v] syspool[%v]",
						by.TableManager.Room.GetSysSumWinLoss(), adjustSyspool, sysPool)
				}

				logMsg := fmt.Sprintf("抽水前水池[%v] 抽水百分比[%v] 調整水池[%v] 系統水池抽水後[%v] 金幣獲利區間[%v] ",
					by.TableManager.Room.FishSysPoolData.GetSysSumWinLoss()*-1, adjust+*by.systemWinCommission.IniCommissionRate, adjustSyspool, by.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss()*-1,
					*by.systemWinCommission.CoinDivision)
				mlog.Debug(logMsg)
				mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), logMsg)

				break
			} else {
				adjust = adjust + (*by.systemWinCommission.IncreaseRate)
			}
			sysPool = sysPool - int64(*by.systemWinCommission.CoinDivision)
		}
	} else if modNum == 0 {
		by.TableManager.Room.HsiTungHuoLiIRate = float32(0)
		by.TableManager.Room.HsiTungHuoLiChouShuiChingFen = 0
	}
}

func (by *Bylogic) killFish() {
	defer util.PanicErrStack()
	if !by.canDone() {
		by.fishtimer = nil
		return
	}

	//140隻以內不刪除魚
	if len(by.fishmanager.fishdata) < int(140) {
		return
	}

	// To store the keys in slice in sorted order
	var keys []int
	for k := range by.fishmanager.fishdata {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	i := 0
	for _, v := range keys {
		f := by.fishmanager.fishdata[int32(v)]
		if f == nil {
			continue
		}

		//24號以上魚不要刪除
		//if f.fishtype >= 24 {
		//	continue
		//}

		delete(by.fishmanager.fishdata, int32(v))
		i++
		limit := float32(len(by.fishmanager.fishdata)) * (float32(by.killfishcount) / float32(100))

		if int(i) >= int(limit) {
			mlog.Debug("超過數量 刪除前總數[%v]，刪除後 魚總數[%v] 刪除數量[%v]", len(by.fishmanager.fishdata)+int(limit), len(by.fishmanager.fishdata), limit)
			break
		}
	}

}

//檢查是否達到魚總數量上限
func (by *Bylogic) isTotalFishMaxLimit() bool {
	result := false

	if by.maxtotalfishlimit <= int32(len(by.fishmanager.fishdata)) {
		mlog.Debug("魚總數量[%v] 最大上限[%v]", len(by.fishmanager.fishdata), by.maxtotalfishlimit)
		result = true
	}

	return result
}

//get 炸彈蟹 rate
func (by *Bylogic) getBoomFishRate(ud *byuserdata) float32 {
	fishdiedRate := float32(ud.teshufishfishodd) / float32(ud.allFishOdd)
	rate := util.Decimal(float64(fishdiedRate), 8)
	logmsg := fmt.Sprintf("取得炸彈螃蟹(or核彈頭)死亡幾率 [%v] = 炸彈螃蟹倍率[%v] / 所有魚倍數總和[%v]", rate, ud.teshufishfishodd, ud.allFishOdd)
	mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, ud.user.TableNo, logmsg)
	if rate > 10000 {
		mlog.Error("get 炸彈蟹 rate[%v]", rate)
	}
	if rate == 0 && !ud.user.IsRobot() {
		rate = 99999
	}
	return rate
}

//打中魚記錄
func (by *Bylogic) AddGameLog(canFish bool, fish *fishinfo, ud *byuserdata, money int32) {

	boomMoney := by.getBoomInfo(by.boomIdOption[ud.boomNumber-1]).boomMoney

	if canFish {
		money = money*ud.getUserOdd() - boomMoney
		if fishSum, ok := ud.fishwinloss[fish.fishid]; ok {
			mlog.Debug("打中魚增加金額 fishid[%v] fishSum[%v] money[%v] oddType[%v]", fish.fishid, fishSum, money, ud.getUserOdd())
			ud.fishwinloss[fish.fishid] = fishSum + money
		} else {
			mlog.Debug("打中魚增加金額 fishid[%v] money[%v] oddType[%v]", fish.fishid, money, ud.getUserOdd())
			ud.fishwinloss[fish.fishid] = money
		}
	} else {
		if fishSum, ok := ud.fishwinloss[fish.fishid]; ok {
			mlog.Debug("沒打死魚減少金額 fishid[%v] fishSum[%v] money[%v] oddType[%v]", fish.fishid, fishSum, boomMoney, ud.getUserOdd())
			ud.fishwinloss[fish.fishid] = fishSum - boomMoney
		} else {
			mlog.Debug("沒打死魚減少金額 fishid[%v] money[%v] oddType[%v]", fish.fishid, boomMoney, ud.getUserOdd())
			ud.fishwinloss[fish.fishid] = -boomMoney
		}
	}
}

//檢查要打的魚是否在不能放分裡面
func (by *Bylogic) canSupertoolReward(fishid int32) bool {
	result := true
	for _, v := range by.nosupertoolreward {
		if v == fishid {
			result = false
		}
	}

	return result
}

//刷新系統放分
func (by *Bylogic) updateSysReword(ok bool, fishdata *byfishdata, ud *byuserdata, money int32) {
	// 超端控制個人殺送分，不計限制魚
	if ud.user.UserType == gamelogic.UserTypeRebot || (!by.canSupertoolReward(fishdata.fishtype) && (by.GetFishVIPUserList() == nil || len(by.GetFishVIPUserList()) == 0)) {
		return
	}

	boomMoney := by.getBoomInfo(by.boomIdOption[ud.boomNumber-1]).boomMoney
	if ok { //打死魚
		money = money - boomMoney

		//個人放分
		for _, v := range by.GetFishVIPUserList() {
			mlog.Debug("打死魚個人放分名單 userid[%v] status[%v] money[%v] = money[%v] - boomMoney[%v] ",
				ud.user.UserID, v.GetStatus(), (money - boomMoney), money, boomMoney)

			if v.GetUserID() == ud.user.UserID && v.GetStatus() == 1 {
				mlog.Debug("打死魚個人放分 userid[%v] money[%v]", ud.user.UserID, money)
				adjustAmount := v.GetFinishAmount() + int64(money)
				v.FinishAmount = proto.Int64(adjustAmount)
				//by.GameLog.AddGameStep("[超端]個人放分 UserId[%v]抓到FishId[%v] 獲得放分金額[%v]", ud.user.UserID, fishdata.fishtype, adjustAmount)
				example.invalid("frameid[%v] 個人放分 UserId[%v]抓到FishId[%v] 獲得放分金額[%v]", by.frameid, ud.user.UserID, fishdata.fishtype, adjustAmount)
				return
			}
		}

		//系統放分
		for _, v := range by.GetFishRewardList() {
			if v.GetFishID() == 0 && v.GetStatus() == 1 {
				adjustAmount := v.GetFinishAmount() + int64(money)
				v.FinishAmount = proto.Int64(adjustAmount)
				//by.GameLog.AddGameStep("[超端]系統所有魚放分 UserId[%v]抓到FishId[%v] 獲得放分金額[%v]", ud.user.UserID, fishdata.fishtype, adjustAmount)
				example.invalid("frameid[%v] 系統所有魚放分 UserId[%v]抓到FishId[%v] 獲得放分金額[%v]", by.frameid, ud.user.UserID, fishdata.fishtype, adjustAmount)
				break
			}
			if v.GetFishID() == fishdata.fishtype && v.GetStatus() == 1 {
				adjustAmount := v.GetFinishAmount() + int64(money)
				v.FinishAmount = proto.Int64(adjustAmount)
				mlog.Debug("抓到魚 減少可放分金額 fishtype[%v] money[%v] FinishAmount[%v] ", fishdata.fishtype, money, adjustAmount)
				//by.GameLog.AddGameStep("[超端]系統放分 抓到指定魚 UserId[%v]抓到FishId[%v] 獲得放分金額[%v]", ud.user.UserID, fishdata.fishtype, adjustAmount)
				example.invalid("frameid[%v] 系統放分 抓到指定魚 UserId[%v]抓到FishId[%v] 獲得放分金額[%v]", by.frameid, ud.user.UserID, fishdata.fishtype, adjustAmount)
				break
			}
		}

	} else {

		//個人
		for _, v := range by.GetFishVIPUserList() {
			//送分
			if v.GetAssignAmount() > 0 {
				if v.GetUserID() == ud.user.UserID && v.GetStatus() == 1 {
					adjustAmount := v.GetFinishAmount() - int64(boomMoney)
					v.FinishAmount = proto.Int64(adjustAmount)
					//by.GameLog.AddGameStep("[超端]個人放分 UserId[%v] 沒抓到FishId[%v] 異動完成金額[%v]", ud.user.UserID, fishdata.fishtype, adjustAmount)
					example.invalid("frameid[%v] 個人放分 UserId[%v] 沒抓到FishId[%v] 異動完成金額[%v]", by.frameid, ud.user.UserID, fishdata.fishtype, adjustAmount)
					return
				}
			}

			//殺分
			if v.GetAssignAmount() < 0 {
				if v.GetUserID() == ud.user.UserID && v.GetStatus() == 1 {
					adjustAmount := v.GetFinishAmount() - int64(boomMoney)
					v.FinishAmount = proto.Int64(adjustAmount)
					//by.GameLog.AddGameStep("[超端]個人殺分 UserId[%v] 沒抓到FishId[%v] 異動完成金額[%v]", ud.user.UserID, fishdata.fishtype, adjustAmount)
					example.invalid("frameid[%v] 個人殺分 UserId[%v] 沒抓到FishId[%v] 異動完成金額[%v]", by.frameid, ud.user.UserID, fishdata.fishtype, adjustAmount)
					return
				}
			}
		}

		//系統
		for _, v := range by.GetFishRewardList() {
			//送分
			sungFenOk := sungFenRefreshfinishAmount(v, boomMoney, fishdata)
			if sungFenOk {
				break
			}

			//殺分
			shaFenOk := shaFenRefreshFinishAmount(v, boomMoney, fishdata)
			if shaFenOk {
				break
			}
		}
	}
}

//殺分時, 更新完成金額
func shaFenRefreshFinishAmount(v *netproto.FishSystemRewardDataList, boomMoney int32, fishdata *byfishdata) bool {
	if v.GetAmount() < 0 {
		if v.GetFishID() == 0 && v.GetStatus() == 1 {
			adjustAmount := v.GetFinishAmount() - int64(boomMoney)
			v.FinishAmount = proto.Int64(adjustAmount)
			return true
		}
		if v.GetFishID() == fishdata.fishtype && v.GetStatus() == 1 {
			mlog.Debug("殺分時 沒抓到魚 減少可放分金額 fishtype[%v] boomMoney[%v] ", fishdata.fishtype, boomMoney)
			adjustAmount := v.GetFinishAmount() - int64(boomMoney)
			v.FinishAmount = proto.Int64(adjustAmount)
			return true
		}
	}
	return false
}

//送分時, 更新完成金額
func sungFenRefreshfinishAmount(v *netproto.FishSystemRewardDataList, boomMoney int32, fishdata *byfishdata) bool {
	if v.GetAmount() > 0 {
		if v.GetFishID() == 0 && v.GetStatus() == 1 {
			adjustAmount := v.GetFinishAmount() - int64(boomMoney)
			v.FinishAmount = proto.Int64(adjustAmount)
			return true
		}
		if v.GetFishID() == fishdata.fishtype && v.GetStatus() == 1 {
			mlog.Debug("沒抓到魚 增加可放分金額 fishtype[%v] boomMoney[%v] ", fishdata.fishtype, boomMoney)
			adjustAmount := v.GetFinishAmount() - int64(boomMoney)
			v.FinishAmount = proto.Int64(adjustAmount)
			return true
		}
	}
	return false
}

//取得系統放分或個人放分的攻擊倍率
func (by *Bylogic) getRate(fishid int32, ud *byuserdata) float32 {
	result := *by.boompower
	//mlog.Debug("取得系統放分或個人放分的攻擊倍率[%v] ", *by.boompower)
	// 超端控制個人殺送分，不計限制魚
	if !by.canSupertoolReward(fishid) && (by.GetFishVIPUserList() == nil || len(by.GetFishVIPUserList()) == 0) {
		if result < 0.1 {
			mlog.Error("rate [%v]", result)
		}
		return result
	}

	if ud.user.IsRobot() {
		rand := util.RndInt(0, len(by.robotfightrate))

		for _, v := range by.robotcantkillfish {
			if fishid == v {
				return float32(0)
			}
		}

		if (result*float32(by.robotfightrate[rand]))/100 < 0.1 {
			mlog.Warn("rate [%v]", (result*float32(by.robotfightrate[rand]))/100)
		}

		return (result * float32(by.robotfightrate[rand])) / 100
	}

	if by.TableManager.Room.GameType == gamelogic.PiSai {
		return result
	}

	//個人放分
	if by.GetFishVIPUserList() != nil {
		result2, useRate := by.vipUserRate(ud)
		if useRate {
			//mlog.Warn("個人放分 result[%v] = result[%v] * result2[%v] ", (result * result2), result, result2)
			mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), fmt.Sprintf("userid[%v] 個人放分 result[%v] = result[%v] * result2[%v]", ud.user.UserID, (result*result2), result, result2))

			if result*result2 < 0.1 {
				mlog.Error("rate [%v]", result*result2)
			}
			return result * result2
		}
	}

	//系統放分
	if by.GetFishRewardList() != nil {
		result2 := by.sysRate(fishid)
		//mlog.Debug("系統放分 result[%v] = result[%v] * result2[%v] ", (result * result2), result, result2)
		mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), fmt.Sprintf("userid[%v] 系統放分 result[%v] = result[%v] * result2[%v]", ud.user.UserID, (result*result2), result, result2))
		result = result * result2
		if result < 0.1 {
			mlog.Error("rate [%v]", result)
		}
	}

	if result < 0.1 {
		mlog.Error("rate [%v]", result)
	}

	return result
}

//系統水池區間調整攻擊率
//P（子彈平均押注）=  該魚場玩家總押注金額/總子彈發數
//Px淨分區間定值x區間數贏=系统水线
//机器人不算
//"系统水线判断的是整个游戏房间的
//e.g.新手场的整个输赢统计"
func (by *Bylogic) adjectSysPoolRate(initRate float32, ud *byuserdata) float32 {
	defer util.PanicErrStack()

	newRate := initRate

	if initRate < 0.1 {
		mlog.Error("InitRate [%v]", initRate)
	}

	if by.TableManager.Room.FishSysPoolData == nil {
		example.invalid("無系統水池資料")
		return newRate
	}

	modNum := by.TableManager.Room.FishSysPoolData.GetSumBullet() % int64(*by.checkSysPoolNum)
	if modNum == 0 {
		//子彈平均押注
		boomAvgBet := util.Decimal(float64(by.TableManager.Room.FishSysPoolData.GetSumBetMoney())/float64(by.TableManager.Room.FishSysPoolData.GetSumBullet()), 2)
		//mlog.Debug("魚場玩家總押注金額[%v] 總子彈發數[%v] = 子彈平均押注[%v] ",
		//	float64(by.TableManager.Room.FishSysPoolData.GetSumBetMoney()), by.TableManager.Room.FishSysPoolData.GetSumBullet(), boomAvgBet)

		for i := 0; i < len(by.syspoolTable); i++ {
			calAreaUp := float32(by.syspoolTable[i].GetArea()) * float32(by.syspoolTable[i].GetConditionOdds()) * boomAvgBet
			calAreaDown := float32(0)
			if i == len(by.syspoolTable)-1 {
				by.syspoolAdjustRate = float32(by.syspoolTable[i].GetWaterAdjust()) / 1000

				//mlog.Debug("A 淨分區間[%v] 魚場內玩家總平均押注[%v] 淨分上限[%v] 淨分下限[%v] 系統總輸贏[%v]", by.syspoolTable[i].GetConditionOdds(),
				//	boomAvgBet, calAreaUp, calAreaDown, -by.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss())

				by.syspoolLog(boomAvgBet, i, calAreaUp, calAreaDown, initRate, ud)
				break
			}

			calAreaDown = float32(by.syspoolTable[i+1].GetArea()) * float32(by.syspoolTable[i+1].GetConditionOdds()) * boomAvgBet

			if float32(by.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss()*-1) > calAreaUp {
				by.syspoolAdjustRate = float32(by.syspoolTable[i].GetWaterAdjust()) / 1000

				//mlog.Debug("B 淨分區間[%v] 魚場內玩家總平均押注[%v] 淨分上限[%v] 淨分下限[%v] 系統總輸贏[%v]",
				//	by.syspoolTable[i].GetConditionOdds(),
				//	boomAvgBet, calAreaUp, calAreaDown, -by.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss())

				by.syspoolLog(boomAvgBet, i, calAreaUp, calAreaDown, initRate, ud)
				break
			}

			if float32(by.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss()*-1) < calAreaUp &&
				float32(by.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss()*-1) >= calAreaDown {
				by.syspoolAdjustRate = float32(by.syspoolTable[i+1].GetWaterAdjust()) / 1000

				//mlog.Debug("C 淨分區間[%v] 魚場內玩家總平均押注[%v] 淨分上限[%v] 淨分下限[%v] 系統總輸贏[%v] ",
				//	by.syspoolTable[i+1].GetConditionOdds(),
				//	boomAvgBet, calAreaUp, calAreaDown, -by.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss())

				by.syspoolLog(boomAvgBet, i+1, calAreaUp, calAreaDown, initRate, ud)
				break
			}
		}

	}
	if initRate+by.syspoolAdjustRate < 0.1 {
		mlog.Error("InitRate [%v] by.SyspoolAdjustRate[%v]", initRate, by.syspoolAdjustRate)
	}
	return newRate + by.syspoolAdjustRate
}

//pool_all系統水池Log
func (by *Bylogic) syspoolLog(boomAvgBet float32, i int, calAreaUp float32, calAreaDown float32, initRate float32, ud *byuserdata) {
	if by.TableManager.Room.FishSysPoolData == nil || by.syspoolTable == nil {
		return
	}

	type LogData struct {
		Seq                          int64
		GroupName                    string
		SumBullet                    int64   //魚場內子彈ID
		BoomAvgBet                   float32 //魚場內玩家總平均押注
		BoomAvgBetArea               float32 //魚場水位線基準值
		CalAreaUp                    float32 //淨分上限
		CalAreaDown                  float32 //淨分下限
		ConditionOdds                int32   //淨分區間
		InitRate                     float32 //原子彈威力
		SyspoolAdjustRate            float32 //應加扣加威力
		AfterSyspoolAdjustRate       float32 //更改後子彈威力
		SysSumWinLoss                int64   //系統總輸贏
		HsiTungHuoLiChouShui         int64   //系統獲利抽水後總輸贏
		HsiTungHuoLiBulletNumber     int32   //系統獲利每N發抽水
		HsiTungHuoLiIRate            float32 //系統獲利每次抽N%
		HsiTungHuoLiChouShuiChingFen int32   //系統獲利當次抽水淨分
		Allrealusersumfish           float32 //玩家總得分
		Allrealuserbetmoney          float32 //玩家總押分
	}

	msg := LogData{
		Seq:                          *by.TableManager.Room.Seq,
		GroupName:                    "syspool_all",
		SumBullet:                    by.TableManager.Room.FishSysPoolData.GetSumBullet(),
		BoomAvgBet:                   boomAvgBet / 100,
		BoomAvgBetArea:               boomAvgBet * float32(by.syspoolTable[i].GetArea()) / 100,
		CalAreaUp:                    calAreaUp / 100,
		CalAreaDown:                  calAreaDown / 100,
		ConditionOdds:                by.syspoolTable[i].GetConditionOdds(),
		InitRate:                     initRate,
		SyspoolAdjustRate:            by.syspoolAdjustRate,
		AfterSyspoolAdjustRate:       initRate + by.syspoolAdjustRate,
		SysSumWinLoss:                by.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() * -1 / 100,
		HsiTungHuoLiChouShui:         by.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() * -1 / 100,
		HsiTungHuoLiBulletNumber:     *by.systemWinCommission.BulletNumber,
		HsiTungHuoLiIRate:            util.Decimal(float64(by.TableManager.Room.HsiTungHuoLiIRate), 2),
		HsiTungHuoLiChouShuiChingFen: by.TableManager.Room.HsiTungHuoLiChouShuiChingFen,
		Allrealusersumfish:           float32(by.getAllRealUserSumFish()) / 100,
		Allrealuserbetmoney:          float32(by.TableManager.Room.FishSysPoolData.GetSumBetMoney()) / 100,
	}

	if msg.HsiTungHuoLiChouShui < 0 && msg.HsiTungHuoLiChouShuiChingFen > 0 {
		mlog.Error("系統輸錢不能抽水")
	}

	mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), msg)
	//mlog.Debug("系統水池log [%+v]", msg)
}

//all玩家總得分
func (by *Bylogic) getAllRealUserSumFish() int64 {
	sum := int64(0)
	for _, v := range by.TableManager.Room.FishPlayerRewardDetail {
		sum += v.GetHitMoney()
		//mlog.Debug("userid[%v] sum[%v]", v.GetUserID(), sum)
	}
	return sum
}

//個人水池區間調整攻擊率
//玩家總輸贏=玩家在這個魚場中打死魚的總金額/玩家在這個魚場花出的總金額
//判斷玩家捕魚總輸贏是否大於等於TotalMoneyWin
//(子彈威力+增加子彈威力+子彈點殺變化)/(賠率*爆炸命中隻數)=死亡概率
//"e.g.大於等於1200且小於1250
//子彈威力就得-30"
//"各魚場各自計算，同個魚場的輸贏會累計
//e.g.初級場玩家總輸贏為1200，出去再進來不會重新計算"
//打死魚再更新並判斷點殺變化
func (by *Bylogic) adjectPlayerPoolRate(initRate float32, ud *byuserdata) float32 {
	defer util.PanicErrStack()

	newRate := initRate

	if !by.TableManager.Room.FishPlayerPoolStatus {
		return newRate
	}

	if initRate < 0.1 {
		mlog.Error("InitRate [%v]", initRate)
	}
	if by.TableManager.Room.FishPlayerRewardDetail == nil ||
		len(by.TableManager.Room.FishPlayerRewardDetail) == int(0) {
		//mlog.Debug("個人水池區間調整攻擊率[%v] ", newRate)
		return newRate
	}

	//for _, v := range by.TableManager.Room.FishPlayerRewardDetail {
	v, ok := by.TableManager.Room.FishPlayerRewardDetail[ud.user.UserID]
	if !ok {
		mlog.Debug("FishPlayerRewardDetail not find userid[%v]", ud.user.UserID)
		return newRate
	}

	if v.GetUserID() == ud.user.UserID {
		playerWinLoss := float32(0)

		if by.playerpoolTable == nil {
			return newRate
		}

		if v.GetHitMoney() == 0 || v.GetSpendMoney() == 0 {
			//example.invalid("userid[%v] HitMoney[%v] SpendMoney[%v] rate[%v] ", v.GetUserID(), v.GetHitMoney(), v.GetSpendMoney(), newRate)
			return newRate
		} else {
			//玩家總輸贏=玩家在這個魚場中打死魚的總金額/玩家在這個魚場花出的總金額
			playerWinLoss = float32(v.GetHitMoney()) / float32(v.GetSpendMoney()) * 1000
		}

		for i := 0; i < len(by.playerpoolTable); i++ {
			if i == len(by.playerpoolTable)-1 {
				by.playerPoolLog(v, ud, i, newRate, playerWinLoss)
				return newRate + float32(by.playerpoolTable[i].GetAttackRate())/1000
			}

			if playerWinLoss > float32(by.playerpoolTable[i].GetTotalMoneyWin()) {
				by.playerPoolLog(v, ud, i, newRate, playerWinLoss)
				return newRate + float32(by.playerpoolTable[i].GetAttackRate())/1000
			}

			if playerWinLoss < float32(by.playerpoolTable[i].GetTotalMoneyWin()) && playerWinLoss >= float32(by.playerpoolTable[i+1].GetTotalMoneyWin()) {
				by.playerPoolLog(v, ud, i, newRate, playerWinLoss)
				return newRate + float32(by.playerpoolTable[i+1].GetAttackRate())/1000
			}
		}

		//mlog.Debug("個人水池 userid[%v] hit[%v] money[%v] = 玩家總輸贏[%v] ",
		//	v.GetUserID(), v.GetHitMoney(), v.GetSpendMoney(), playerWinLoss)

	}
	//}

	return newRate
}

//個人水池 Log
func (by *Bylogic) playerPoolLog(v *netproto.FishPlayerPoolDetailList, ud *byuserdata, i int, newRate float32, playerWinLoss float32) {
	modNum := by.TableManager.Room.FishSysPoolData.GetSumBullet() % int64(*by.checkSysPoolNum)
	if modNum == 0 && by.playerpoolTable != nil {

		type LogData struct {
			Seq              int64
			GroupName        string
			FrameId          int64
			UserId           int32   //此玩家賬號ID
			Bid              int32   //此玩家子彈ID
			PlayerPoolStatus int     //單一玩家超端狀態
			AttackRate       float32 `json:"單一玩家應加扣威力"`   //單一玩家應加扣威力
			AfterAttackRate  float32 `json:"單一玩家更改後子彈威力"` //單一玩家更改後子彈威力
			PlayerWinLoss    float32 `json:"單一玩家獲利區間"`    //單一玩家獲利區間
			HitMoney         float32 //單一玩家得分
			SpendMoney       float32 //單一玩家押分
			HitSpendScale    float32 `json:"單一玩家獲利比例"` //單一玩家獲利比例
		}

		msg := LogData{
			Seq:              *by.TableManager.Room.Seq,
			GroupName:        "personal_kill",
			FrameId:          by.frameid,
			UserId:           v.GetUserID(),
			Bid:              by.bid[ud.user.SeatNo],
			PlayerPoolStatus: 1,
			AttackRate:       float32(by.playerpoolTable[i].GetAttackRate()) / 1000,
			AfterAttackRate:  newRate + float32(by.playerpoolTable[i].GetAttackRate())/1000,
			PlayerWinLoss:    playerWinLoss,
			HitMoney:         float32(v.GetHitMoney()) / 100,
			SpendMoney:       float32(v.GetSpendMoney()) / 100,
			HitSpendScale:    float32(v.GetHitMoney()) / float32(v.GetSpendMoney()),
		}

		mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), msg)

		jsonMsg, _ := json.Marshal(msg)
		logMsg := fmt.Sprintf("個人水池[%v]", string(jsonMsg))
		mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), logMsg)
		//mlog.Debug("個人水池log[%+v]", logMsg)
	}
}

//針對未充值玩家限制可以打死魚與威力
func (by *Bylogic) adjustFreePlayerRate(initRate float32, ud *byuserdata, fishdata *byfishdata) float32 {
	if by.TableManager.Room.GameType == gamelogic.PiSai ||
		by.freePlayerBoomPower == nil ||
		ud.user.SumCharge != int64(0) {
		return initRate
	}

	cantKill := false
	for _, v := range by.freePlayerBoomPower.monitorFishID {
		if v == fishdata.fishtype {
			cantKill = true
			break
		}
	}
	if !cantKill {
		return initRate
	}

	v, _, fishOdd := by.fishmanager.getFishLuckyOdds(fishdata)

	feturescore := ud.user.Score + ud.winmoney + (int64(fishOdd * ud.getUserOdd()))

	if feturescore >= int64(by.freePlayerBoomPower.moneyLimit) {

		logmsg := fmt.Sprintf("針對未充值玩家userid[%v] 目前金額[%v],限制可以打死魚[%v] 魚倍率[%v]Lucky[%v]與子彈倍率[%v]的攻擊威力限制為[%v], ",
			ud.user.UserID, ud.user.Score+ud.money, v.name, fishOdd, v.luckybonus, ud.getUserOdd(), by.freePlayerBoomPower.boomPower)
		mlog.Debug(logmsg)
		mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, ud.user.TableNo, logmsg)

		return float32(by.freePlayerBoomPower.boomPower) / 1000
	}

	return initRate
}

//新手優惠
func (by *Bylogic) adjectNewPlayerRate(initRate float32, ud *byuserdata) float32 {

	if ud.user.IsRobot() || by.TableManager.Room.GameType == gamelogic.PiSai {
		return initRate
	}

	if by.TableManager.Room.FisnNewPlayerDetail == nil {
		by.TableManager.Room.FisnNewPlayerDetail = make(map[int32]*netproto.FishBeginnerOffer, 0)
	}

	if p, ok := by.TableManager.Room.FisnNewPlayerDetail[ud.user.UserID]; ok {
		if p.GetStatus() == 1 || p.GetBullet() > by.newPlayerTable.offerLimitCount { //已經送完新手優惠
			return initRate
		}

		//未充值新手優惠已經到贏超過14元,不再提供優惠
		if ud.user.SumCharge == 0 && (ud.allwin >= int64(by.newplayerbonuslimit) || ud.winmoney >= int64(by.newplayerbonuslimit)) {
			logMsg := fmt.Sprintf("userid[%v] win[%v] 新手優惠已經到贏超過14元,不再提供優惠", ud.user.UserID, ud.allwin)
			mlog.Debug(logMsg)
			mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, by.TableNo(), logMsg)
			p.Status = proto.Int(1)
			return initRate
		}

		if v, ok2 := by.TableManager.Room.FishPlayerRewardDetail[ud.user.UserID]; ok2 {
			//玩家淨贏分=玩家擊殺魚總得分-玩家總壓分
			playerWinLoss := int32((v.GetHitMoney() - v.GetSpendMoney()))

			for i := 0; i < len(by.newPlayerTable.poolTable); i++ {
				//mlog.System("GetTotalMoneyWin[%v] basescore[%v]", by.NewPlayerTable.poolTable[i].GetTotalMoneyWin(), by.basescore )

				if i == len(by.newPlayerTable.poolTable)-1 {
					//mlog.System("index[%v] NewPlayer[%+v]", by.NewPlayerTable.poolTable[i].GetConditionOdds(), p)
					return float32(by.newPlayerTable.poolTable[i].GetBoomPower()) / 1000
				}

				if playerWinLoss > by.newPlayerTable.poolTable[i].GetTotalMoneyWin()*by.basescore {
					//by.newPlayLog()	TODO
					//mlog.System("index[%v] NewPlayer[%+v]", by.NewPlayerTable.poolTable[i].GetConditionOdds(), p)
					return float32(by.newPlayerTable.poolTable[i].GetBoomPower()) / 1000
				}

				if playerWinLoss < by.newPlayerTable.poolTable[i].GetTotalMoneyWin()*by.basescore && playerWinLoss >= by.newPlayerTable.poolTable[i].GetTotalMoneyWin()*by.basescore {
					//mlog.Debug("index[%v] NewPlayer[%+v]", by.newPlayerTable.poolTable[i].GetConditionOdds(), p)
					return float32(by.newPlayerTable.poolTable[i+1].GetBoomPower()) / 1000
				}
			}
		}

		//mlog.Debug("NewPlayer[%+v]", p)
	} else {
		for _, v := range by.newPlayerTable.poolTable {
			if v.GetConditionOdds() == 0 {
				return float32(v.GetBoomPower()) / 1000
			}
		}
	}

	return initRate
}

//個人超端攻擊率
func (by *Bylogic) vipUserRate(ud *byuserdata) (float32, bool) {
	useUserRate := false
	result := float32(1)
	for _, v := range by.GetFishVIPUserList() {

		//放分
		if v.GetAssignAmount() > 0 {
			if v.GetUserID() == ud.user.UserID && v.GetFinishAmount() < v.GetAssignAmount() {
				mlog.Debug("攻擊魚 個人放分 fishid[%v] rate[%v] status[%v]", v.GetUserID(), v.GetRate(), v.GetStatus())
				result = float32(1) + float32(v.GetRate())/100
				if result < 0 {
					result = 0.1
				}
				useUserRate = true
				//by.GameLog.AddGameStep("[超端]個人放分 UserID[%v] 放分機率[%v]", ud.user.UserID, result)
				break
			} else {
				//v.Status = proto.Int32(0)
			}
		}

		//殺分
		if v.GetAssignAmount() < 0 {
			if v.GetUserID() == ud.user.UserID && v.GetFinishAmount() > v.GetAssignAmount() {
				mlog.Debug("攻擊魚 個人殺分 userid[%v] rate[%v] status[%v]", v.GetUserID(), v.GetRate(), v.GetStatus())
				result = float32(1) - float32(v.GetRate())/100
				if result < 0 {
					result = 0.1
				}
				useUserRate = true
				//by.GameLog.AddGameStep("[超端]個人殺分 UserID[%v] 殺分機率[%v]", ud.user.UserID, result)
				break
			} else {
				//v.Status = proto.Int32(0)
			}
		}

	}
	return result, useUserRate
}

//系統超端攻擊率
func (by *Bylogic) sysRate(fishid int32) float32 {
	result := float32(1)
	for _, v := range by.GetFishRewardList() {

		//放分
		if v.GetAmount() > 0 {
			//全部魚
			if v.GetFishID() == 0 && v.GetFinishAmount() < v.GetAmount() {
				result = float32(v.GetRate()) / 100
				mlog.Debug("攻擊魚 系統放分 fishid[%v] rate[%v] status[%v]", v.GetFishID(), v.GetRate(), v.GetStatus())
				//by.GameLog.AddGameStep("[超端]系統所有魚放分  打中FishID[%v] 放分機率[%v]", fishid, result)
				break
			}

			if v.GetFishID() == fishid && v.GetFinishAmount() < v.GetAmount() {
				mlog.Debug("攻擊魚 系統放分 fishid[%v] rate[%v] status[%v]", v.GetFishID(), v.GetRate(), v.GetStatus())
				result = float32(v.GetRate()) / 100
				//by.GameLog.AddGameStep("[超端]系統放分 打中指定FishID[%v] 放分機率[%v]", fishid, result)
				break
			} else {
				//v.Status = proto.Int32(0)
			}
		}

		//殺分
		if v.GetAmount() < 0 {
			//全部魚
			if v.GetFishID() == 0 && v.GetFinishAmount() > v.GetAmount() {
				result = float32(v.GetRate()) / 100
				//by.GameLog.AddGameStep("[超端]系統所有魚殺分  打中FishID[%v] 殺分機率[%v]", fishid, result)
				break
			}

			if v.GetFishID() == fishid && v.GetFinishAmount() > v.GetAmount() {
				mlog.Debug("攻擊魚 系統殺分 fishid[%v] rate[%v] status[%v]", v.GetFishID(), v.GetRate(), v.GetStatus())
				result = float32(v.GetRate()) / 100
				//by.GameLog.AddGameStep("[超端]系統殺分 打中指定FishID[%v] 殺分機率[%v]", fishid, result)
				break
			} else {
				//v.Status = proto.Int32(0)
			}
		}
	}
	return result
}

//获取设置消息
func (by *Bylogic) getSettingMsg() *msg.Message {
	defer util.PanicErrStack()

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BY_GameMessageClassID_BYSettingID)

	pdata := new(netproto.BY_Setting)
	pdata.BaseScore = proto.Int32(by.basescore)
	pdata.FrameTime = proto.Int32(int32(by.frametime))
	pdata.MaxBoom = proto.Int32(by.maxboom)
	pdata.ZuantoulongxiaMaxBoom = proto.Int32(by.zuantoulongxiamaxboom)
	pdata.ZlxReceiveSpeed = proto.Int32(by.zlxreceivespeed)
	pdata.BoomSetting = make([]*netproto.BY_BoomSetting, 0)
	pdata.GameType = proto.Int32(by.TableManager.Room.GameType)

	for _, v := range by.boommanager {
		bs := new(netproto.BY_BoomSetting)
		bs.BoomNumber = proto.Int32(v.boomNumber)
		bs.BoomName = proto.String(v.boomName)
		bs.BoomUIType = proto.Int32(v.boomUIType)
		bs.BoomRate = proto.Int32(v.boomMoney)
		bs.BoomSpeed = proto.Float32(v.boomSpeed)
		bs.ShootInterval = proto.Float32(v.shootInterval)
		bs.BoomType = proto.Int32(v.boomType)
		bs.Rebound = proto.Int32(v.rebound)
		pdata.BoomSetting = append(pdata.BoomSetting, bs)
	}

	pdata.IceSetting = new(netproto.BY_IceSetting)
	pdata.IceSetting.IceNumber = proto.Int32(by.icemanager.iceNumber)
	pdata.IceSetting.IceName = proto.String(by.icemanager.iceName)
	pdata.IceSetting.IceUIType = proto.Int32(by.icemanager.iceUIType)
	pdata.IceSetting.IceTime = proto.Int32(by.icemanager.iceTime)
	pdata.IceSetting.IceResetTime = proto.Int32(by.icemanager.iceResetTime)
	pdata.IceSetting.IceMoney = proto.Int32(by.icemanager.iceMoney)
	pdata.BackgroundList = proto.String(by.background)

	pdata.FishSetting = make([]*netproto.BY_FishSetting, 0)
	for _, fishInfo := range by.fishmanager.config.fishinfoarr {
		by.addFishInfoSetting(by.fishmanager.fishweightarr, fishInfo, pdata)
		by.addFishInfoSetting(by.fishmanager.bossfishweight, fishInfo, pdata)
		by.addFishInfoSetting(by.fishmanager.teshufishweight, fishInfo, pdata)
		by.addFishInfoSetting(by.fishmanager.goldfishweight, fishInfo, pdata)
	}

	if by.TableManager.Room.GameType == gamelogic.PiSai {
		by.pisaiSettingMsg(pdata)
	}

	pmsg.MsgData = pdata

	return pmsg
}

//設定client可以顯示的魚,有在roomtable才能顯示
func (by *Bylogic) addFishInfoSetting(fishweightarr []*fishweight, fishInfo *fishinfo, pdata *netproto.BY_Setting) {
	for _, v := range fishweightarr {
		if v.fishid == fishInfo.fishid {
			fishInfoSetting := new(netproto.BY_FishSetting)
			fishInfoSetting.FishId = proto.Int32(fishInfo.fishid)
			fishInfoSetting.Reward = append(fishInfoSetting.Reward, fishInfo.rewardinfo[0].reward)
			if len(fishInfo.rewardinfo) > 1 {
				fishInfoSetting.Reward = append(fishInfoSetting.Reward, fishInfo.rewardinfo[len(fishInfo.rewardinfo)-1].reward)
			}
			if fishInfo.bosstypeinfo.fishtype == Dujiaojing {
				fishInfoSetting.Reward = []int32{0}
			}

			fishInfoSetting.FishType = proto.Int32(fishInfo.bosstype)
			fishInfoSetting.SPBoomArea = proto.Float32(fishInfo.bosstypeinfo.spboomarea)
			pdata.FishSetting = append(pdata.FishSetting, fishInfoSetting)
			//mlog.System("fishInfoSetting[%v] bosstype[%v] add type", fishInfoSetting, fishInfo.bosstype)
		}
	}
}

//踢出房提示訊息
func (by *Bylogic) getLeaveTipMessage(tip string, kicktime int32) *msg.Message {
	defer util.PanicErrStack()

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BY_GameMessageClassID_BYLeaveTip)

	pdata := new(netproto.BY_LeaveMsg)
	pdata.TipMessage = proto.String(tip)
	pdata.Time = proto.Int32(kicktime)

	pmsg.MsgData = pdata

	return pmsg
}

//更新用户分数
func (by *Bylogic) updateUserSorce() {
	for _, ud := range by.userdata {
		if ud.money != 0 {

			if by.flag == ExperienceCashFlag && !ud.user.IsRobot() {
				gUser := by.GetUserByUserID(ud.user.UserID)
				gUser.ExperienceCash += ud.money
				by.UpdateExperienceScoreFunc(ud.user.UserID, ud.money, 0)
				gUser.Score = gUser.ExperienceCash
				mlog.Debug("體驗更新金幣 ud.money[%v] Score[%v] ExperienceCash[%v]",
					ud.money, gUser.Score, gUser.ExperienceCash)

			} else if by.TableManager.Room.GameType == gamelogic.PiSai {
				p, k := by.getPiSaiUserData()
				if k {
					by.UpdateUserScoreFunc(ud.user.UserID, p.Bonus, 0)
				}

			} else {
				if ud.user.UserType != gamelogic.UserTypeRebot {
					by.updateMoney += ud.money
				}

				by.UpdateUserScoreFunc(ud.user.UserID, ud.money, 0)
				if !ud.user.IsRobot() {
					mlog.Debug("frameid[%v] 玩家%d更新后的钱为%d money[%v] by.updateMoney[%v]", by.frameid, ud.user.UserID, ud.user.Score, ud.money, by.updateMoney)
				}

			}

			ud.money = 0
		}
	}

	//更新新手優惠完成狀態
	by.updateNewPlayerStatus()
}

//更新新手優惠完成狀態
func (by *Bylogic) updateNewPlayerStatus() {

	if by.TableManager.Room.GameType == gamelogic.PiSai {
		return
	}

	by.updateFishNewPlayerPoolRetID()

	for _, v := range by.TableManager.Room.FisnNewPlayerDetail {
		if v.GetBullet() > by.newPlayerTable.offerLimitCount {
			v.Status = proto.Int32(1)
		}
	}
}

//更新資料庫新手優惠
func (by *Bylogic) updateFishNewPlayerPoolRetID() {
	defer util.PanicErrStack()

	if by.TableManager.Room.GameType == gamelogic.PiSai {
		return
	}

	if by.TableManager.Room.HostAndBackup == 2 {
		return
	}

	if len(by.TableManager.Room.FisnNewPlayerDetail) == 0 {
		return
	}

	if by.PiSaiUserPause() {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_UpdateBeginnerOfferID)
	pdata := new(netproto.FishBeginnerOfferDetail)

	for _, v := range by.userdata {
		if newPlayer, ok := by.TableManager.Room.FisnNewPlayerDetail[v.user.UserID]; ok {
			sumCharge := by.GetUserByUserID(newPlayer.GetUserID()).SumCharge
			if sumCharge == 0 && newPlayer.GetScore() >= by.newplayerbonuslimit {
				newPlayer.Status = proto.Int(1)
			}

			data := new(netproto.FishBeginnerOffer)
			data.UserID = proto.Int32(newPlayer.GetUserID())
			data.Bullet = proto.Int32(newPlayer.GetBullet())
			data.Score = proto.Int32(newPlayer.GetScore())
			data.Status = proto.Int32(newPlayer.GetStatus())
			//mlog.Debug("更新資料庫新手優惠 [%+v]", data)
			pdata.NewPlayerList = append(pdata.NewPlayerList, data)
		}
	}

	pmsg.MsgData = pdata
	by.SendRPCMessageFunc(false, pmsg)
}

func (by *Bylogic) sendToRobot() {
	for _, ur := range by.AllTableUser() {
		if ur.IsRobot() {
			by.SendMessageToUser(ur, by.robotmsg)
		}
	}
}

//個人水池, 寫入個人水池和取得該水池的個人捕魚機率
func (by *Bylogic) sendWinMoney() {
	for _, ud := range by.userdata {
		//机器人不发
		if ud.user.IsRobot() {
			continue
		}

		//by.sendRPCMsg(ud.user.UserID, ud.winmoney)
	}
}

//收到RPC消息解析处理消息
func (by *Bylogic) OnRecvRPCMessage(m *msg.Message) {
	rm, ok := m.MsgData.(*rpcmsg.RPCMessage)
	if !ok {
		return
	}
	mmbyte := rm.MsgData.([]byte)
	if len(mmbyte) == 0 {
		return
	}

	if m.SClassID == int32(netproto.HallMsgClassID_GameCompetioniInfoRuleRetID) {
		by.setContentInfo(m)
		//檢查玩家是否有買票
		by.checkUserTicket()
	}

	if m.SClassID == int32(netproto.DBServerClassID_GetContestUserDataRetID) {
		by.continuContestUser(m)
	}

}

//游戏结束清掉所有的timer
func (by *Bylogic) gameEnd() {

	by.stopCreatFish()

	if by.asynctimer != nil {
		by.asynctimer.Stop()
		by.asynctimer = nil
	}

	if by.updatetimer != nil {
		by.updatetimer.Stop()
		by.updatetimer = nil
	}

	if by.rewardtimer != nil {
		by.rewardtimer.Stop()
		by.rewardtimer = nil
	}

	if by.rpctimer != nil {
		by.rpctimer.Stop()
		by.rpctimer = nil
	}

	if by.fishkilltimer != nil {
		by.fishkilltimer.Stop()
		by.fishkilltimer = nil
	}

	if by.reloadconfigtimer != nil {
		by.reloadconfigtimer.Stop()
		by.reloadconfigtimer = nil
	}

	if by.TableManager.Room.GameType == gamelogic.PiSai {

		if by.piSaiLeaveTimer != nil {
			//mlog.Error("stop piSaiLeaveTimer")
			by.piSaiLeaveTimer.Stop()
			by.piSaiLeaveTimer = nil
		}
	}

	by.clearData()
	by.GameEndFunc()
}

func (by *Bylogic) stopCreatFish() {
	if by.fishtimer != nil {
		by.fishtimer.Stop()
		by.fishtimer = nil
	}
	if by.bossfishtimer != nil {
		by.bossfishtimer.Stop()
		by.bossfishtimer = nil
	}
	if by.teshufishtimer != nil {
		by.teshufishtimer.Stop()
		by.teshufishtimer = nil
	}
	if by.shoalfishtimer != nil {
		by.shoalfishtimer.Stop()
		by.shoalfishtimer = nil
	}
	if by.goldfishtimer != nil {
		by.goldfishtimer.Stop()
		by.goldfishtimer = nil
	}
	if by.fishwavetimer != nil {
		by.fishwavetimer.Stop()
		by.fishwavetimer = nil
	}
}

//清理数据
func (by *Bylogic) clearData() {
	by.framemsg = make([]*netproto.BY_FrameS2C, 0)
	by.nextsendmsg = new(msg.Message)
	by.nextsendmsg.BClassID = int32(netproto.MessageBClassID_Game)
	by.nextsendmsg.SClassID = int32(netproto.BY_GameMessageClassID_FrameS2CID)
	by.framedata = new(netproto.BY_FrameS2C)
	by.frameid = 0
	by.bid = make([]int32, 4)
	for i := int32(0); i < 4; i++ {
		by.bid[i] = 10000000 * i
	}
	by.userdata = make([]*byuserdata, 0)
	by.fishmanager.delAllFish()
}

//啟動定時刪除魚
func (by *Bylogic) startKillTimer() {
	updatescop := util.NewIntScope(by.killfishtime-5, by.killfishtime+5)
	by.fishkilltimer = by.TimerFuncManager.NewNamedRandomTicker(updatescop, by.killFish, "killFish")
}

//启动生成BOSS鱼
func (by *Bylogic) startBossTimer() {
	if by.bossfishtimer != nil {
		by.bossfishtimer.Stop()
		by.bossfishtimer = nil
	}

	by.bossfishtimer = by.TimerFuncManager.NewTimer(by.fishmanager.bossfishtime, by.bossEnter)
}

//启动生成特殊鱼
func (by *Bylogic) startTeshuTimer() {
	if by.teshufishtimer != nil {
		by.teshufishtimer.Stop()
		by.teshufishtimer = nil
	}

	by.teshufishtimer = by.TimerFuncManager.NewTimer(by.fishmanager.teshufishtime, by.newTeshuFish)
}

//生成魚群
func (by *Bylogic) startShoalTimer() {

	if by.shoalfishtimer != nil {
		by.shoalfishtimer.Stop()
		by.shoalfishtimer = nil
	}

	by.shoalfishtimer = by.TimerFuncManager.NewMilliSecTimer(by.fishmanager.fishshoaltime, by.newFishShoal)
}

func (by *Bylogic) startGoldFishTimer() {
	if by.goldfishtimer != nil {
		by.goldfishtimer.Stop()
		by.goldfishtimer = nil
	}

	by.goldfishtimer = by.TimerFuncManager.NewMilliSecTimer(by.fishmanager.goldfishtime, by.newGoldFish)
}

//处理用户施放冰冻技能
func (by *Bylogic) userCastIce(ud *byuserdata) {
	ud.updateUserMoney(-by.icemanager.iceMoney)
	ud.updateBetAmount(by.icemanager.iceMoney, by.icemanager.iceMoney, 0)
	//ud.updateFlowAmount(by.icemanager.iceMoney)
	by.updatePiSaiUserScore(-by.icemanager.iceMoney)
	by.updateSysPoolWinLoss(ud, -by.icemanager.iceMoney)
	logmsg := fmt.Sprintf("用户%v施放冰冻技能 扣金幣:%v", ud.user.UserID, -by.icemanager.iceMoney)
	//mlog.Debug(logmsg)
	mylog.MongoLog(ud.user.UserID, *by.TableManager.Room.Seq, ud.user.TableNo, logmsg)

	endtime := by.icemanager.iceTime
	//停止上一次冰冻
	if by.icetimer != nil {
		endtime -= by.icetimer.WaitSeconds - by.icetimer.ElapsedSeconds()
		by.icetimer.Stop()
	}

	//不产生鱼，将原来产生的鱼删除掉
	by.framedata.NewFish = make([]*netproto.BY_NewFish, 0)
	by.icetimer = by.TimerFuncManager.NewTimer(by.icemanager.iceTime, by.iceEnd)
	by.fishmanager.ice = true
	by.framedata.UseIce = proto.Bool(true)

	//对所有timer进行延迟
	if by.fishwavetimer != nil {
		delaytime := by.fishwavetimer.WaitSeconds - by.fishwavetimer.ElapsedSeconds()
		if delaytime <= 0 {
			delaytime = 1
		}
		by.fishwavetimer.Stop()
		by.fishwavetimer = by.TimerFuncManager.NewTimer(endtime+delaytime, by.fishwavetimer.Func)
	}

	if by.bossfishtimer != nil {
		delaytime := by.bossfishtimer.WaitSeconds - by.bossfishtimer.ElapsedSeconds()
		if delaytime <= 0 {
			delaytime = 1
		}
		by.bossfishtimer.Stop()
		by.bossfishtimer = by.TimerFuncManager.NewTimer(endtime+delaytime, by.bossfishtimer.Func)
	}

	if by.teshufishtimer != nil {
		delaytime := by.teshufishtimer.WaitSeconds - by.teshufishtimer.ElapsedSeconds()
		if delaytime <= 0 {
			delaytime = 1
		}
		by.teshufishtimer.Stop()
		by.teshufishtimer = by.TimerFuncManager.NewTimer(endtime+delaytime, by.teshufishtimer.Func)
	}

	if by.shoalfishtimer != nil {
		delaytime := by.shoalfishtimer.WaitSeconds - by.shoalfishtimer.ElapsedSeconds()
		if delaytime <= 0 {
			delaytime = 1
		}
		by.shoalfishtimer.Stop()
		by.shoalfishtimer = by.TimerFuncManager.NewTimer(endtime+delaytime, by.shoalfishtimer.Func)
	}

	if by.goldfishtimer != nil {
		delaytime := by.goldfishtimer.WaitSeconds - by.goldfishtimer.ElapsedSeconds()
		if delaytime <= 0 {
			delaytime = 1
		}
		by.goldfishtimer.Stop()
		by.goldfishtimer = by.TimerFuncManager.NewTimer(endtime+delaytime, by.goldfishtimer.Func)
	}
}

func (by *Bylogic) iceEnd() {
	by.icetimer = nil
	by.fishmanager.ice = false
	by.newfish()
}

//是否能执行下一步
func (by *Bylogic) canDone() bool {
	//有冰冻和鱼潮都不执行
	if by.fishmanager.fishwave || by.fishmanager.ice {
		return false
	}
	return true
}
