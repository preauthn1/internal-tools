package zjh

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

const opttimerstr = "zjh_opt_timeout"

const (
	actionName_FaPai        = "fapai" //发牌动画
	actionName_Bet          = "bet"   //下注跟注加注动画
	actionName_Allin        = "allin" //全押动画
	actionName_AllinCompare = "allincompare"
	actionName_Compare      = "compare" //比牌动画
	actionName_End          = "end"     //结束动画
)

const (
	BAOZI       = 1  //豹子
	SHUNJIN     = 2  //順金
	JINHUA      = 3  //金花
	SHUNZI      = 4  //順子
	BIGDUIZI    = 5  //對子(9-A)
	SMALLDUIZI  = 6  //對子(2-8)
	DANPAI      = 7  //單張
	BIGJINHUA   = 11 //A/K金花
	SMALLJINHUA = 12 //非A/K金花
	DUIZI       = 13 //對子
)

// Added by Joecool Wang at 2019-05-29
// 監控用戶用參數
// 下局輸贏
const (
	vipNextWinLose_Lose = 0 // 輸
	vipNextWinLose_Win  = 1 // 贏
)

// 開始/暫停
const (
	vipStatus_Suspend = 0 // 暫停
	vipStatus_Start   = 1 // 開始
)

//是否完成
const (
	vipFinish_No  = 0 // 否
	vipFinish_Yes = 1 // 是
)

// 未指定
const (
	vipNotSpecified = 99
)

// Added by Joecool Wang at 2019-05-29 end

type ZhaJinHuaLogic struct {
	*gamelogic.DynamicGameLogicBase
	opttimer             *event.TimerFunc
	poker                poker.Poker           //游戏扑克牌
	zhuangSeat           int32                 //庄家座位号
	initSeat             int32                 //遊戲開始的初始座位號, 預設是莊家
	userData             []*ZjhUserData        //用户数据
	actionName_curr      string                //当前播放的动画
	userCount            int32                 //用户人数统计
	currSeat             int32                 //当前操作的座位号
	optTime              int32                 //可操作时间
	betCurr              int32                 //当前的下注额
	betTotal             int64                 //当前下注总额
	bCompare             bool                  //是否能比牌
	bFirst               bool                  //是否是第一个操作的玩家
	bAllin               bool                  //是否能全押
	winmoney             int64                 //赢家赢钱
	bGameEnd             bool                  //是否游戏结束
	bSendGameResult      bool                  //是否发送了结算消息
	currturn             int32                 //当前进行了多少轮
	canOpt               bool                  //当前是否可操作
	isadd                bool                  //本轮是否有人加注
	disRobotExchangeCard bool                  //關閉機器人換牌
	beforeGameStartStep  []string              //前置處理的記錄
	userinthisgame       []*netproto.ZJH_Poker //此局玩家
	SuperToolService     *ZjhSuperToolService  //超端服務
	gameTimeOutTimer     *event.TimerFunc      //遊戲進行太久的控制器
	map_iVIPInitCardType map[int32]int64       // 每局開始時 VIPUser 最初指定的牌型，Added by Joecool Wang at 2019-06-03
	//以下配置不会随游戏改变
	betBase           int32         //底注
	betAdd            []int32       //加注数组
	betLast           int32         //封顶注
	bFirstTurn        bool          //是否第一局
	maxplayer         int32         //最多玩家
	pokerType         int32         //1为随机牌，2为发牌器
	readyTime         int32         //准备时间
	retAction         int32         //结束动画超时
	fapaiAction       int32         //发牌动画超时
	betAction         int32         //下注跟注的动画
	compareAction     int32         //比牌动画
	maxturn           int32         //经过了多少轮以后用户掉线将不会弃牌而找人比牌
	checkturn         int32         //可以看牌的轮数
	dapairandom       []int32       //大牌隨機組數
	dapaiweight       *weightconfig //权重配置文件
	dapaiNumber       int           //大牌數量
	normalweight      *weightconfig //权重配置文件
	cardType          []int         //大牌種類
	radomnickname     int32         //特殊名稱的開關
	reversepokerorder int32         //發牌順序 順時針 逆時針 設置 ，0為順時針，1為逆時針
}

//创建诈金花逻辑接口
func NewZhaJinHuaLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *ZhaJinHuaLogic {
	glc := new(ZhaJinHuaLogic)
	glc.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(t, tfManager, cfgs)
	glc.getzjhConfig(cfgs)
	glc.bFirstTurn = true
	glc.bGameEnd = true
	glc.SuperToolService = glc.GetTable().Room.SuperToolService.(*ZjhSuperToolService)
	glc.map_iVIPInitCardType = make(map[int32]int64, vipNotSpecified)
	glc.RegisterAfterGameStartFunc(func() {
		//因遊戲開始前無法添加操作記錄, 故將遊戲開始前緩存的操作記錄添加進去
		glc.GameLog.AddGameStep("------ before game start ------")
		for _, log := range glc.beforeGameStartStep {
			glc.GameLog.AddGameStep(log)
		}
		glc.GameLog.AddGameStep("------ before game start ------")
		//因遊戲開始前已經 setuserdata 、但是因為即時到帳在之後才更新, 只能在這裡再做一次金額的異動
		for _, ud := range glc.userData {
			ud.userMoney = ud.User.Score
		}
		//glc.GameTimeOutStopSig = make(chan bool, 1)
		glc.gameTimeOutTimer = glc.TimerFuncManager.NewNamedTimer(glc.GetServerConfig().GameTimeOut*60, func() {
			mlog.Error("遊戲超時異常 取消遊戲 tableno %d", glc.TableNo())
			glc.GameCancelFunc("游戏异常，请稍候再试")
		}, fmt.Sprintf("gametimeoutcd_%d", glc.TableNo()))
		//go glc.GameTimeOutKickUser()
		glc.betDi()
		glc.fapai() //发牌
		glc.broadCastCurrturn()
	})

	glc.RegisterBeforeGameStartFunc(func() {
		//因為機器人要進行換牌所以，牌型必需要在 beforeGameStart 定義好
		glc.clearGameData()
		glc.defineCardType()
	})

	glc.RegisterBeforeGameEndFunc(func() {
		//glc.GameTimeOutStopSig <- true
		glc.gameTimeOutTimer.Stop()
	})

	glc.RegisterAfterGameEndFunc(func() {
		for _, u := range t.GetTableUser() {
			glc.SuperToolService.GetZjhSuperToolHandler().DelUserPokerCard(u.UserID)
			glc.bGameEnd = true
			glc.UserLeave(u)
		}
	})

	return glc
}

//获取当前配置
func (zjh *ZhaJinHuaLogic) getzjhConfig(cfgs *config.Group) {
	scfgs := cfgs.GetConfig("gameserver")
	zjh.betBase = int32(scfgs.GetConfigInt("gameserver", "basescore"))
	zjh.optTime = int32(scfgs.GetConfigInt("gameserver", "opttime"))
	zjh.retAction = int32(scfgs.GetConfigInt("gameserver", "retaction"))
	zjh.readyTime = int32(scfgs.GetConfigInt("gameserver", "readywaittime"))
	zjh.maxplayer = int32(scfgs.GetConfigInt("gameserver", "maxplayer"))
	zjh.fapaiAction = int32(scfgs.GetConfigInt("gameserver", "fapaiaction"))
	zjh.pokerType = int32(scfgs.GetConfigInt("gameserver", "pokertype"))
	zjh.betAction = int32(scfgs.GetConfigInt("gameserver", "betaction"))
	zjh.compareAction = int32(scfgs.GetConfigInt("gameserver", "compareaction"))
	zjh.maxturn = int32(scfgs.GetConfigInt("gameserver", "maxturn"))
	zjh.checkturn = int32(scfgs.GetConfigInt("gameserver", "checkturn"))
	zjh.radomnickname = int32(scfgs.GetConfigInt("gameserver", "radomnickname"))
	zjh.reversepokerorder = int32(scfgs.GetConfigInt("gameserver", "reversepokerorder"))
	betstr := scfgs.GetConfigValue("gameserver", "bet")
	lastindex := 0
	for index, char := range betstr {
		if char == ',' {
			substr := betstr[lastindex:index]
			lastindex = index + 1
			bet, _ := strconv.ParseInt(substr, 10, 32)
			zjh.betLast = int32(bet)
			zjh.betAdd = append(zjh.betAdd, zjh.betLast)
		}
	}

	//dapai setting
	cgf := cfgs.GetConfig("logicconfig")
	dapaistr := cgf.GetConfigValue("logicsetting", "dapai")
	lastindex = 0
	for index, char := range dapaistr {
		if char == ',' {
			substr := dapaistr[lastindex:index]
			lastindex = index + 1
			dapai, _ := strconv.ParseInt(substr, 10, 32)
			zjh.dapairandom = append(zjh.dapairandom, int32(dapai))
		}
	}
	zjh.dapaiweight = newDapaiWeightConfig(cfgs)
	zjh.normalweight = newNormalWeightConfig(cfgs)
}
func (zjh *ZhaJinHuaLogic) defineCardType() {
	if zjh.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
		provider := zjh.SuperToolService.GetZjhSuperToolHandler()
		zjh.poker = poker.NewRandomPoker(false, false, 1)
		zjh.poker.Shuffle()
		var uids []int32
		// Modified by Joecool Wang at 2019-05-30
		// 給予初值
		/*
			var LoserUserID int32
			var WinnerUserID int32
		*/
		var LoserUserID int32 = 0
		var WinnerUserID int32 = 0
		// Modified by Joecool Wang at 2019-05-30 end
		var assignCardType = make(map[int32]int)
		for _, user := range zjh.userData {
			uids = append(uids, user.User.UserID)
		}

		for _, userid := range provider.Proiority { //按優先順序將監控用戶取出
			if ok, _ := util.InArray(userid, uids); ok { //只有存在於桌內的用戶
				example.invalid("設定牌型：VIPUserID: %d", userid)
				if stcVIPUser_d := provider.GetVipUser(userid); stcVIPUser_d != nil {
					example.invalid("GetStatus() = %d, GetFinish() = %d, GetNextWinLose() = %d, GetNextWinLoseNum() = %d, GetFinishNextWinLoseNum() = %d, GetNextType() = %d",
						stcVIPUser_d.GetStatus(), stcVIPUser_d.GetFinish(), stcVIPUser_d.GetNextWinLose(), stcVIPUser_d.GetNextWinLoseNum(), stcVIPUser_d.GetFinishNextWinLoseNum(), stcVIPUser_d.GetNextType())
				}

				// Modified by Joecool Wang at 2019-05-30
				// 判斷是否需對監控用戶進行換牌
				// 修正因判斷錯誤導致監控條件僅執行一次便結束的情況，Finish 應該在遊戲結算時統計
				/*
					if vipuser := provider.GetVipUser(userid); vipuser != nil && vipuser.GetStatus() == 1 && vipuser.GetFinish() == 0 {
						if WinnerUserID == 0 && vipuser.GetNextWinLose() == 1 {
							WinnerUserID = userid
							vipuser.Finish = proto.Int64(1)
						} else if LoserUserID == 0 && vipuser.GetNextWinLose() == 0 {
							LoserUserID = userid
							vipuser.Finish = proto.Int64(1)
						} else if WinnerUserID == 0 && vipuser.GetNextWinLoseNum() != 99 && vipuser.GetNextWinLoseNum() > 0 {
							WinnerUserID = userid
							vipuser.FinishNextWinLoseNum = proto.Int64(vipuser.GetFinishNextWinLoseNum() + 1)
							if vipuser.GetFinishNextWinLoseNum() >= vipuser.GetNextWinLoseNum() {
								vipuser.Finish = proto.Int64(1)
							}
						} else if WinnerUserID == 0 && vipuser.GetNextWinLoseNum() != 99 && vipuser.GetNextWinLoseNum() < 0 {
							LoserUserID = userid
							//mlog.Debug("userid %d, GetFinishNextWinLoseNum %d", LoserUserID, vipuser.GetFinishNextWinLoseNum())
							vipuser.FinishNextWinLoseNum = proto.Int64(vipuser.GetFinishNextWinLoseNum() + 1)
							if vipuser.GetFinishNextWinLoseNum() >= -vipuser.GetNextWinLoseNum() {
								vipuser.Finish = proto.Int64(1)
							}
						}
						if vipuser.GetNextType() != 99 {
							assignCardType[userid] = int(vipuser.GetNextType())
						}
					}
				*/
				if stcVIPUser := provider.GetVipUser(userid); stcVIPUser != nil {
					// 設定要贏還是輸
					if stcVIPUser.GetStatus() == vipStatus_Start && stcVIPUser.GetFinish() == vipFinish_No {
						// 取輸贏局數的絕對值
						iNextWinLoseCount := stcVIPUser.GetNextWinLoseNum()
						if iNextWinLoseCount < 0 {
							iNextWinLoseCount = -iNextWinLoseCount
						}

						if stcVIPUser.GetNextWinLoseNum() == vipNotSpecified {
							// 若輸贏局數被設定為［未指定］則不做換牌處理，並將 Finish 設定為［完成］
							WinnerUserID = 0
							LoserUserID = 0
							stcVIPUser.Finish = proto.Int64(vipFinish_Yes)
							stcVIPUser.Status = proto.Int32(vipStatus_Suspend)
							zjh.SuperToolService.GetZjhSuperToolHandler().SaveVipUser(userid)
						} else {
							if stcVIPUser.GetFinishNextWinLoseNum() >= iNextWinLoseCount {
								// 如果進行局數已達到輸贏局數，且有指定牌型，則派發指定牌型
								WinnerUserID = 0
								LoserUserID = 0
							} else if WinnerUserID == 0 && stcVIPUser.GetNextWinLoseNum() > 0 {
								// 設定為贏
								WinnerUserID = userid
							} else if LoserUserID == 0 && stcVIPUser.GetNextWinLoseNum() < 0 {
								// 設定為輸
								LoserUserID = userid
							}
						}
					}

					// 設定指定牌型
					if stcVIPUser.GetNextType() != vipNotSpecified {
						assignCardType[userid] = int(stcVIPUser.GetNextType())
						zjh.map_iVIPInitCardType[userid] = stcVIPUser.GetNextType()
					}
				}
				// Modified by Joecool Wang at 2019-05-30 end
			}
		}

		if WinnerUserID != 0 || LoserUserID != 0 { //如果需要指定輸贏就不指定牌型了
			zjh.disRobotExchangeCard = true
			assignCardType = make(map[int32]int)
		} else if len(assignCardType) > 0 { //如果指定牌型, 機器人不可以換牌
			zjh.disRobotExchangeCard = true
		}
		zjh.beforeGameStartStep = append(zjh.beforeGameStartStep, fmt.Sprintf("超端指定參數 Winner[%d] Loser[%d] assignCardType[%v]", WinnerUserID, LoserUserID, assignCardType))
		userCardType := make(map[int32]int)
		if zjh.SuperToolService.fapaiRule == 1 { //第一種發牌規則, 分配大牌打亂大牌順序分配給用戶
			userCardType = zjh.fapaiDapaiWeight()
		} else { //第二種牌型規則
			userCardType = zjh.fapaiNormalWeight()
		}
		//檢查用戶是否有被指定牌型
		for uid, ct := range userCardType {
			if cardtype, ok := assignCardType[uid]; ok {
				userCardType[uid] = cardtype
				if vipuser := provider.GetVipUser(uid); vipuser != nil {
					// Remarked by Joecool Wang at 2019-06-01
					// 變更控制模式，改在 calWinMoney() 作是否完成判斷
					//vipuser.Finish = proto.Int64(1)
					//zjh.SuperToolService.GetZjhSuperToolHandler().SaveVipUser(uid)
					// Remarked by Joecool Wang at 2019-06-01 end
					zjh.beforeGameStartStep = append(zjh.beforeGameStartStep, fmt.Sprintf("用戶 %d 原牌型 %d 執行指定牌型 %d", uid, ct, cardtype))
				}
				continue
			}
		}
		example.invalid("userCardType %v", userCardType)
		zjh.poker = poker.NewZJHBigFixedPoker(userCardType, zjh.poker.(*poker.RandomPoker))
		zjh.beforeGameStartStep = append(zjh.beforeGameStartStep, fmt.Sprintf("初始牌組\r\n%s", zjh.poker))
		zjh.beforeGameStartStep = append(zjh.beforeGameStartStep, fmt.Sprintf("disRobotExchangeCard %v", zjh.disRobotExchangeCard))

		if WinnerUserID != 0 {
			zjh.beforeGameStartStep = append(zjh.beforeGameStartStep, fmt.Sprintf("用戶 %d 執行指定贏", WinnerUserID))
			zjh.VipExchangeCard(WinnerUserID, 1)
			// Remarked by Joecool Wang at 2019-05-30，在遊戲結束時才存檔 VIPUser 資料
			//zjh.SuperToolService.GetZjhSuperToolHandler().SaveVipUser(WinnerUserID)
			example.invalid("用戶 %d 執行指定贏", WinnerUserID)
		} else if LoserUserID != 0 {
			zjh.beforeGameStartStep = append(zjh.beforeGameStartStep, fmt.Sprintf("用戶 %d 執行指定輸", LoserUserID))
			zjh.VipExchangeCard(LoserUserID, 2)
			// Remarked by Joecool Wang at 2019-05-30，在遊戲結束時才存檔 VIPUser 資料
			//zjh.SuperToolService.GetZjhSuperToolHandler().SaveVipUser(LoserUserID)
			example.invalid("用戶 %d 執行指定輸", LoserUserID)
		}
		if WinnerUserID != 0 || LoserUserID != 0 {
			zjh.beforeGameStartStep = append(zjh.beforeGameStartStep, fmt.Sprintf("超端執行後的牌組\r\n %s", zjh.poker))
		}
	} else { //发牌器发牌
		zjh.poker = poker.NewFilePoker(zjh.Configs, zjh.AllTableUser())
		zjh.poker = poker.NewFixedFilePoker(zjh.poker.(*poker.FilePoker))
	}
}

//清理用户数据
func (zjh *ZhaJinHuaLogic) clearGameData() {
	zjh.zhuangSeat = -1
	zjh.initSeat = -1
	zjh.betTotal = 0
	zjh.winmoney = 0
	zjh.canOpt = false
	zjh.bCompare = false
	zjh.bFirst = true
	zjh.opttimer = nil
	zjh.bGameEnd = false
	zjh.bSendGameResult = false
	zjh.bAllin = false
	zjh.currturn = 1
	zjh.isadd = false
	zjh.cardType = nil
	zjh.userinthisgame = nil
	zjh.beforeGameStartStep = make([]string, 0)
	zjh.setUserData()
	//優化找莊方法
	i := util.RndInt(0, len(zjh.userData)) //隨機找莊
	example.invalid("初始用戶列表 %v 莊隨機 %d", zjh.userData, i)
	zjh.zhuangSeat = zjh.userData[i].User.SeatNo
	//
	//
	//
	//if zjh.zhuangSeat == -1 || zjh.getUserBySeatNo(zjh.zhuangSeat) == nil {
	//	zjh.zhuangSeat = rand.Int31n(zjh.maxplayer)
	//	for i := 0; i < 10; i++ {
	//		if zjh.getUserBySeatNo(zjh.zhuangSeat) != nil {
	//			break
	//		}
	//
	//		zjh.zhuangSeat = rand.Int31n(zjh.maxplayer)
	//	}
	//
	//	//10次都没找到庄家就用第一家作为庄家
	//	if zjh.getUserBySeatNo(zjh.zhuangSeat) == nil {
	//		zjh.zhuangSeat = zjh.userData[0].User.SeatNo
	//	}
	//}
	zjh.initSeat = zjh.zhuangSeat

	zjh.betCurr = zjh.betBase
	zjh.currSeat = zjh.zhuangSeat
	zjh.disRobotExchangeCard = false
}

//开始下底
func (zjh *ZhaJinHuaLogic) betDi() {
	for _, urd := range zjh.userData {
		zjh.setUserMoney(urd, false)
		zjh.GameLog.AddGameStep("用户%d投入底注%d", urd.User.UserID, zjh.betCurr)
	}

	zjh.sendbetBaseMsg()
}

//发送下底消息
func (zjh *ZhaJinHuaLogic) sendbetBaseMsg() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_BetBaseId)

	pdata := new(netproto.ZJH_BetBase)
	pdata.Betbase = proto.Int32(zjh.betBase)
	pdata.Betcount = proto.Int32(zjh.userCount)
	pmsg.MsgData = pdata
	zjh.BoardcastMessage(pmsg)
}

//开始发牌
func (zjh *ZhaJinHuaLogic) fapai() {
	provider := zjh.SuperToolService.GetZjhSuperToolHandler()

	//確定大牌份數
	mlog.Debug("开始发牌")
	//给所有用户发牌
	for _, ur := range zjh.userData {
		cards := zjh.poker.Deal(ur.User.UserID, 3)
		zjh.setUserCards(ur, cards)
		cardtype := GetType(cards)
		//mlog.Debug("===============userid %d cardtype %d cards %v",ur.User.UserID, cardtype, cards)
		zjh.SuperToolService.GetZjhSuperToolHandler().SaveCardType(int(cardtype))
		mlog.Debug("cards=%v", cards)
		provider.AddUserPokerCard(ur.User.UserID, cards.String(), typeStringMap[int(cardtype)])
		zjh.GameLog.Fapai(ur.User.UserID, cards.String())
		zjh.userinthisgame = append(zjh.userinthisgame, zjh.getPokerTypeMsg(ur))
	}

	//通知用户发牌
	zjh.BoardcastMessage(zjh.getFapaiMsg())
	//注册发牌动画
	zjh.actionName_curr = zjh.RegisterAnimation(actionName_FaPai, zjh.fapaiAction, zjh.waitOpt)
}

//给牌排序
func (zjh *ZhaJinHuaLogic) sortCard(cards *poker.CardGroup) {
	arrlen := len(cards.Cards)
	for i := 0; i < arrlen; i++ {
		for j := 0; j+1 < arrlen-i; j++ {
			card1 := cards.Cards[j].CardNumber*10 + cards.Cards[j].Color
			card2 := cards.Cards[j+1].CardNumber*10 + cards.Cards[j+1].Color
			if card1 > card2 {
				cards.Swap(j, j+1)
			}
		}
	}
}

//将用户的牌放到用户数据里面去。
func (zjh *ZhaJinHuaLogic) setUserCards(ud *ZjhUserData, cards *poker.CardGroup) {
	zjh.sortCard(cards)
	for i := 0; i < len(cards.Cards); i++ {
		mlog.Debug("排序后的第%v牌:%v", i+1, cards.Cards[i])
	}

	ud.userCards = poker.NewCardGroup(cards.Cards)
}

//用户离开 replace dynamicleavegamelogicbase.UserLeave
func (glc *ZhaJinHuaLogic) UserLeave(gameuser *gamelogic.GameUser) bool {
	mlog.Debug("zjh user leave %v", gameuser)
	urd := glc.getUserBySeatNo(gameuser.SeatNo)
	//if !glc.bGameEnd && glc.userCount > 1 { // 這判斷會造成剩一人未結算時用戶離開遊戲的寫分問題
	if glc.IsGaming() {
		if urd != nil && (!urd.canEarlyChange() || !urd.isClear) {
			return false
		}
	}
	glc.SuperToolService.GetZjhSuperToolHandler().DelUserPokerCard(gameuser.UserID)
	return glc.LogicBase.UserLeave(gameuser)
}

//用户消息处理逻辑
func (zjh *ZhaJinHuaLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	if !zjh.CheckGameUserMessage(user) {
		zjh.KickUserFunc(user.UserID, "非法操作")
		return
	}

	//不是给游戏逻辑服务器的消息
	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	if gmsg.SClassID == int32(netproto.GameMessageClassID_UserActionCommitID) {
		zjh.CheckUserAnimationCommited(user, gmsg, zjh.actionName_curr, nil)
		return
	}

	//没有找到相关玩家数据
	urd := zjh.getUserBySeatNo(user.SeatNo)
	if urd == nil {
		mlog.Debug("用户未找到%v", user)
		return
	}

	//用戶提前換桌
	if gmsg.SClassID == int32(netproto.ZJH_GameMessageClassID_EarlyChangeTable) {
		zjh.earlyChangeTable(urd, "change")
		return
	}

	//用戶提前離桌
	if gmsg.SClassID == int32(netproto.ZJH_GameMessageClassID_EarlyLeaveTable) {
		zjh.earlyChangeTable(urd, "leave")
		return
	}

	if zjh.userCount <= 1 {
		mlog.Debug("游戏处于结算状态")
		return
	}

	//用户不可操作
	if !urd.canOpt() {
		mlog.Debug("用户已经放弃或比牌已输%v", urd)
		return
	}

	if gmsg.SClassID == int32(netproto.ZJH_GameMessageClassID_FollowEndID) {
		zjh.GameLog.AddGameStep("在第%d轮用户%d选择跟到底", zjh.currturn, user.UserID)
		zjh.setUserFollowEnd(urd)
		return
	}
	if gmsg.SClassID == int32(netproto.ZJH_GameMessageClassID_AllinID) {
		zjh.GameLog.AddGameStep("在第%d轮用户%d选择全押", zjh.currturn, user.UserID)
		zjh.setUserAllin(urd)
		return
	}

	//用户弃牌
	if gmsg.SClassID == int32(netproto.ZJH_GameMessageClassID_DropID) {
		mlog.Debug("收到弃牌消息")
		zjh.userDrop(user)
		return
	}

	//用戶看牌
	if gmsg.SClassID == int32(netproto.ZJH_GameMessageClassID_CheckCardsID) {
		if urd.isCheck {
			mlog.Warn("已经看过牌了")
			return
		}

		if urd.User.Score-urd.userMoney <= 0 {
			if zjh.currturn < zjh.checkturn {
				mlog.Error("目前階段還不能看牌")
				return
			}
		}

		zjh.GameLog.AddGameStep("在第%d轮用户%d看牌", zjh.currturn, user.UserID)
		urd.isCheck = true
		zjh.SendMessageToUser(user, zjh.getPokerType(urd))
		zjh.SendMessageToUser(user, zjh.getPokerMsg(urd))
		zjh.broadCastCheckCard(user.SeatNo)
		zjh.sendOptMsg(false)
		return
	}

	if user.SeatNo != zjh.currSeat {
		mlog.Error("用户座位号非当前操作用户座位号%v", user)
		return
	}

	if urd.userMoney <= 0 {
		mlog.Debug("用户的钱都为负数了%v", user)
		return
	}

	//用户加注
	if gmsg.SClassID == int32(netproto.ZJH_GameMessageClassID_BetAddID) {
		if ok, _ := zjh.userBetAdd(urd, gmsg.MsgData.([]byte)); !ok {
			mlog.Debug("加注不对所以弃牌")
			zjh.userDrop(user)
		} else {
			zjh.bFirst = false
			zjh.waitOpt()
		}

		return
	}

	//略过用户动画
	if gmsg.SClassID != int32(netproto.GameMessageClassID_UserActionCommitID) {
		zjh.bFirst = false
	}

	//用户比牌
	if gmsg.SClassID == int32(netproto.ZJH_GameMessageClassID_CompareID) {
		if urd.isAllin {
			mlog.Error("用戶 %d , 目前為全押狀態，無法比牌", user.UserID)
			return
		}
		if zjh.currturn < 3 {
			mlog.Error("第三回合才能比牌%v", user)
			return
		}

		//操作扣钱
		zjh.compareCardsByplayer(urd, gmsg.MsgData.([]byte))
		return
	}

	//用户跟注
	if gmsg.SClassID == int32(netproto.ZJH_GameMessageClassID_FollowID) {
		zjh.userFollow(urd)
		return
	}
}

//用户提前換桌
func (zjh *ZhaJinHuaLogic) earlyChangeTable(urd *ZjhUserData, early string) {
	if !urd.canEarlyChange() && !urd.isClear {
		return
	}

	if !urd.isClear {
		var flowAmount, betAmount, payAmount int64
		betMoney := urd.userMoney - urd.User.Score //下注的金額
		win, tax := zjh.GetUserWinAndTax(betMoney)
		flowAmount = util.ABSInt64(win)
		betAmount = util.ABSInt64(betMoney)
		payAmount = betAmount + win
		urd.win = win
		mlog.Debug("玩家%v 提前結算下注的金額 userMoney %d, 進房分數 %d, win %d, tax %d", urd.User.UserID, urd.userMoney, urd.User.Score, win, tax)
		//mlog.Debug("玩家%v 提前結算下注的金額 %d, %d_%d", urd.User.UserID, win, urd.userMoney, urd.User.Score)
		zjh.UpdateUserScoreFunc(urd.User.UserID, win, tax)
		urd.isClear = true
		urd.User.FlowAmount += flowAmount
		urd.User.BetAmount += betAmount
		urd.User.PayAmount += payAmount
		zjh.GameLog.AddGameStep("在第%d轮用户%d提前換桌，提前寫分 %d, 稅 %d, 流水 %d, 有效投注 %d, 返還金額 %d", zjh.currturn, urd.User.UserID, win, tax, flowAmount, betAmount, payAmount)
		zjh.GameLog.Context().UserUpdateScore(urd.User.UserID, win)
		zjh.GameLog.Context().UserBetScore(urd.User.UserID, betAmount)
		zjh.GameLog.Context().UserPayoutScore(urd.User.UserID, payAmount)
		zjh.GameLog.Context().UserValidScore(urd.User.UserID, flowAmount)
		zjh.GameLog.Context().UserTaxAmount(urd.User.UserID, tax)
	} else {
		zjh.GameLog.AddGameStep("在第%d轮用户%d提前換桌", zjh.currturn, urd.User.UserID)
	}
	//提前離桌
	zjh.UserEarlyChangeTable(urd.User, early)
}

//用户跟注
func (zjh *ZhaJinHuaLogic) userFollow(urd *ZjhUserData) {
	if !zjh.canOpt {
		example.invalid("现在用户还不能操作")
		return
	}

	//if urd.isAllin {
	//	mlog.Debug("Allin後不能跟注")
	//	return
	//}

	if zjh.opttimer != nil {
		zjh.opttimer.Stop()
		zjh.opttimer = nil
	}
	odd := zjh.getOdds(urd, false)
	zjh.GameLog.AddGameStep("在第%d轮用户%d跟注，扣钱%d", zjh.currturn, urd.User.UserID, zjh.betCurr*odd)
	zjh.setUserMoney(urd, false)
	zjh.sendFollowMsg(urd)
	zjh.canOpt = false
	zjh.actionName_curr = zjh.RegisterAnimation(actionName_Bet, zjh.betAction, zjh.waitOpt)
	//zjh.waitOpt()
}

//用户全押
func (zjh *ZhaJinHuaLogic) userAllinAction(urd *ZjhUserData) {
	if zjh.opttimer != nil {
		zjh.opttimer.Stop()
		zjh.opttimer = nil
	}

	//選擇身上钱較少玩家的金額進行全押()

	betCurr := urd.userMoney
	for _, u := range zjh.userData {
		if u.User.UserID != urd.User.UserID && u.canOpt() {
			if u.userMoney < betCurr {
				betCurr = u.userMoney
			}
			u.isAllin = true
		}
	}
	odd := zjh.getOdds(urd, false)
	zjh.betCurr = int32(betCurr) / 2
	zjh.GameLog.AddGameStep("在第%d轮用户%d全押，扣钱%d", zjh.currturn, urd.User.UserID, zjh.betCurr*odd)
	zjh.setUserMoney(urd, false)
	zjh.sendAllinMsg(urd)
	zjh.canOpt = false
	zjh.actionName_curr = zjh.RegisterAnimation(actionName_Allin, zjh.betAction, zjh.waitOpt)
}

//设置用户跟到底
func (zjh *ZhaJinHuaLogic) setUserFollowEnd(user *ZjhUserData) {
	user.setFollowEnd()
	if zjh.currSeat == user.User.SeatNo {
		zjh.userFollowEnd()
	}
}

//设置用户全押
func (zjh *ZhaJinHuaLogic) setUserAllin(user *ZjhUserData) {
	//確認用戶是否可以全押
	if !zjh.bAllin {
		return
	}
	if !user.canOpt() {
		return
	}
	if !zjh.canOpt {
		mlog.Debug("现在用户还不能操作")
		return
	}
	user.setAllin()
	if zjh.currSeat == user.User.SeatNo {
		zjh.userAllinAction(user)
	}
}

//获取扑克消息
func (zjh *ZhaJinHuaLogic) getPokerMsg(ud *ZjhUserData) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_CardsInfo)

	pmsg.MsgData = zjh.getCardsNetMsg(ud)

	return pmsg
}

//获取发牌消息
func (zjh *ZhaJinHuaLogic) getFapaiMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_FaPaiId)
	pdata := new(netproto.ZJH_FaiPai)
	pdata.Seatno = proto.Int32(zjh.zhuangSeat)
	pdata.Bfirst = proto.Bool(zjh.bFirstTurn)
	for i := 0; i < len(zjh.betAdd); i++ {
		pdata.BetArr = append(pdata.BetArr, zjh.betAdd[i])
	}
	pmsg.MsgData = pdata
	zjh.bFirstTurn = false
	return pmsg
}

//设置用户数据
func (zjh *ZhaJinHuaLogic) setUserData() {
	zjh.userData = make([]*ZjhUserData, 0)
	zjh.userCount = 0
	for _, ur := range zjh.AllTableUser() {
		userdata := NewZjhUserData(ur)
		zjh.userData = append(zjh.userData, userdata)
		zjh.userCount++
	}
}

// 获取下一个可操作的用户
func (zjh *ZhaJinHuaLogic) waitOpt() {
	if zjh.reversepokerorder == int32(0) {
		zjh.waitoptClockwise()
	}
	if zjh.reversepokerorder == int32(1) {
		zjh.waitoptReverse()
	}
}

//获取下一个可操作的用户 順時針
func (zjh *ZhaJinHuaLogic) waitoptClockwise() {
	//只剩一个人的时候结算
	if zjh.userCount <= 1 {
		zjh.sendGameRetMsg()
		return
	}

	//当前已经是庄家操作了,能比牌了。
	if zjh.currSeat == zjh.initSeat && !zjh.bFirst {
		zjh.currturn++
		if zjh.currturn > 2 { //第三輪的時候可以比牌和Allin
			zjh.bCompare = true
			zjh.bAllin = true
		}
		zjh.isadd = false
		zjh.broadCastCurrturn()
	}

	//遊戲開始的初始座位號, 預設是莊家
	initUser := zjh.getUserBySeatNo(zjh.initSeat)
	if initUser == nil || !initUser.canOpt() {
		//新的初始座位號 -> 遞增取餘數為座位號
		for i := int32(0); i < zjh.maxplayer; i++ {
			seat := (zjh.initSeat + i) % zjh.maxplayer
			urd := zjh.getUserBySeatNo(seat)
			if urd != nil && urd.canOpt() {
				zjh.initSeat = seat
				break
			}
		}
	}

	nextSeat := zjh.userData[0].User.SeatNo
	usercount := len(zjh.userData)
	for i := 0; i < usercount; i++ {
		if zjh.userData[i].canOpt() {
			nextSeat = zjh.userData[i].User.SeatNo
			break
		}
	}

	if nextSeat == -1 {
		return
	}

	for _, urd := range zjh.userData {
		mlog.Debug("遍历座位号:%v,是否可操作%v", urd.User.SeatNo, urd.canOpt())
		if urd.User.SeatNo > zjh.currSeat && urd.canOpt() {
			nextSeat = urd.User.SeatNo
			break
		}
	}

	mlog.Debug("可操作的座位号为：%v", nextSeat)
	zjh.currSeat = nextSeat

	//當所有人都跟全押時, 輪到主動全押玩家就進入結算
	if zjh.getUserBySeatNo(nextSeat).User.UserID == zjh.getUserBySeatNo(nextSeat).AllInUserId {
		zjh.allInSettlement()
	} else {
		zjh.sendOptMsg(false)
		zjh.canOpt = true
		zjh.userFollowEnd()
	}
}

//获取下一个可操作的用户 逆時針
func (zjh *ZhaJinHuaLogic) waitoptReverse() {

	//此程式是把原本的發牌順序為座位01234修改改成04321

	// userdataReverse 為 zjh.userData 左右倒
	// userdataReverse := [10]*ZjhUserData{} // 因為GO語言規定陣列大小要先定好，只能使用以下方式
	userdataReverse := make([]*ZjhUserData, zjh.maxplayer+2) //  先保守一點，避免userdataReverse在操作時會oute of range所以設定為最多人數 zjh.maxplayer +2

	// 開始在 userdataReverse 裡塞資料
	operateTimes := (len(zjh.userData) - len(zjh.userData)%2) / 2 // operateTimes為需要對userdataReverse塞資料的操作次數
	for i := 0; i < operateTimes; i++ {
		userdataReverse[0+i] = zjh.userData[len(zjh.userData)-1-i] // 塞入資料到userdataReverse
		userdataReverse[len(zjh.userData)-1-i] = zjh.userData[0+i] // 塞入資料到userdataReverse
	}

	// 塞入中間玩家的資料
	hasMiddle := len(zjh.userData) % 2 // 但是假設玩家資料為5份，中間玩家的資料之前還沒塞入 userdataReverse，所以進行以下操作
	if hasMiddle == 1 {                // 如果有中間玩家的資料
		userdataReverse[((len(zjh.userData) - 1) / 2)] = zjh.userData[((len(zjh.userData) - 1) / 2)] // 塞入資料到 userdataReverse
	}

	//只剩一个人的时候结算
	if zjh.userCount <= 1 {
		zjh.sendGameRetMsg()
		return
	}

	//当前已经是庄家操作了,能比牌了。
	if zjh.currSeat == zjh.initSeat && !zjh.bFirst {
		zjh.currturn++
		if zjh.currturn > 2 { //第三輪的時候可以比牌和Allin
			zjh.bCompare = true
			zjh.bAllin = true
		}
		zjh.isadd = false
		zjh.broadCastCurrturn()
	}

	//遊戲開始的初始座位號, 預設是莊家
	initUser := zjh.getUserBySeatNo(zjh.initSeat)
	if initUser == nil || !initUser.canOpt() {
		//新的初始座位號 -> 遞 減 取餘數為座位號，順序為 0 4 3 2 1
		seat := zjh.initSeat
		for i := int32(0); i < zjh.maxplayer; i++ {
			seat = (seat - 1)
			if seat == int32(-1) { // 這裡在處理 初始座位號 由 0 直接 跳到最大座位號
				seat = int32(len(zjh.userData) - 1) // 直接 跳到最大座位號
			} // 直接 跳到最大座位號
			urd := zjh.getUserBySeatNo(seat)
			if urd != nil && urd.canOpt() { // 到 zjh.userData (未顛倒的實實資料) 去查詢使用者可否再下注
				zjh.initSeat = seat // 如果這位玩家還可以下注，那這位玩家就設定為起始位置
				break
			}
		}
	}

	//第一次獲得nextSeat
	// 如果現在座位是0，第一次會取得最大可以 行動( 下注 棄牌 ) 的最大座位值，以座位 4 3 2 1 0 為例，第一次的nextSeat值為 4 或者是 3
	nextSeat := userdataReverse[0].User.SeatNo // 43210 起始為4
	usercount := len(zjh.userData)
	for i := 0; i < usercount; i++ {
		if userdataReverse[0+i].canOpt() {
			nextSeat = userdataReverse[0+i].User.SeatNo
			break
		}
	}

	if nextSeat == -1 { // 不能有座位小於 0
		return
	}

	//第二次獲得nextSeat
	// 如果現在座位是0，這個迴圈將沒有作用，因為沒有任何一個座位值小於0，除了 -1 ，但是在第一次就取得最大可以 行動( 下注 棄牌 ) 的最大座位值，剛好形成 0 的下一個座位為 最大座位值
	for _, urd := range userdataReverse[0:len(zjh.userData)] {
		if urd.User.SeatNo < zjh.currSeat && urd.canOpt() {
			nextSeat = urd.User.SeatNo
			break
		}
	}

	mlog.Debug("可操作的座位号为：%v", nextSeat)
	zjh.currSeat = nextSeat // 在這裡正式決定下一個座位號

	//當所有人都跟全押時, 輪到主動全押玩家就進入結算
	if zjh.getUserBySeatNo(nextSeat).User.UserID == zjh.getUserBySeatNo(nextSeat).AllInUserId {
		zjh.allInSettlement()
	} else {
		zjh.sendOptMsg(false)
		zjh.canOpt = true
		zjh.userFollowEnd()
	}
}

//當所有人都跟全押時, 輪到主動全押玩家就進入結算
func (zjh *ZhaJinHuaLogic) allInSettlement() {
	//comret := comPareCards(user.userCards, dstuser.userCards)
	comret := 0
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_AllInSettlement)
	pdata := new(netproto.ZJH_ComparePoker)

	for _, ud1 := range zjh.userData {
		ud1ok := !ud1.isDrop && !ud1.isFaild
		if !ud1ok {
			continue
		}

		pdata.OtherPoker = append(pdata.OtherPoker, zjh.getPokerTypeMsg(ud1))

		for _, ud2 := range zjh.userData {
			ud2ok := !ud2.isDrop && !ud2.isFaild && ud1.User.UserID != ud2.User.UserID
			if ud2ok {
				comret = comPareCards(ud1.userCards, ud2.userCards)
				example.invalid("比牌 ud1[%v] ud2[%v] ret[%v]", ud1.User.UserID, ud2.User.UserID, comret)

				if comret == 1 {
					ud2.isFaild = true
					//pdata.OtherPoker = append(pdata.OtherPoker, zjh.getPokerTypeMsg(ud2))
					zjh.GameLog.AddGameStep("所有全押比牌 userid[%v]card[%v] userid[%v]card[%v] ", ud1.User.UserID, ud1.userCards.Cards, ud2.User.UserID, ud2.userCards.Cards)
					example.invalid("所有全押比牌 userid[%v]card[%v] userid[%v]card[%v] ", ud1.User.UserID, ud1.userCards.Cards, ud2.User.UserID, ud2.userCards.Cards)
				}
				//儲值達500的玩家才可看牌
				if ud2.User.SumCharge >= 50000 {
					pdata.OtherPoker = zjh.userinthisgame
				} else {
					//自己的扑克
					pdata.OtherPoker = append(pdata.OtherPoker, zjh.getPokerTypeMsg(ud2))
				}
			}
		}

	}

	for _, v := range pdata.GetOtherPoker() {
		example.invalid("全押win PokerType[%v] Userpokergroup[%v]  ", v.GetPokerType(), v.GetUserpokergroup())
	}

	pmsg.MsgData = pdata
	zjh.BoardcastMessage(pmsg)

	zjh.userCount = 0
	zjh.canOpt = false
	zjh.actionName_curr = zjh.RegisterAnimation(actionName_AllinCompare, 1, zjh.sendGameRetMsg)
}

//处理用户跟到底
func (zjh *ZhaJinHuaLogic) userFollowEnd() {
	//第一个操作的人不能跟到底
	if zjh.bFirst {
		return
	}

	ud := zjh.getUserBySeatNo(zjh.currSeat)
	if ud == nil {
		return
	}

	if !ud.isFollowEnd {
		return
	}
	mlog.Debug("3 ud.isFollowEnd=%v, id = %v", ud.isFollowEnd, ud.User.UserID)
	//用户掉线了，那么不进行跟到底了
	if ud.User.IsOffLine() {
		ud.setFollowEnd()
		return
	}

	//看钱是否够跟到底
	betadd := (zjh.betAdd[len(zjh.betAdd)-1]) * 2   //注額最大值*看牌倍率
	usercount := (float64(zjh.userCount) - 0.5) * 2 //玩家數量*比牌倍率
	betlimit := int64(betadd) * int64(usercount)    //9/12 betlimit 移除扣底注
	if ud.userCanFollow(betlimit) {
		//够用户直接跟注
		zjh.userFollow(ud)
	} else {
		//钱不够找一个人比牌
		zjh.autoCompare()
	}
}

//发送操作消息
func (zjh *ZhaJinHuaLogic) sendOptMsg(bRecover bool) {
	ud := zjh.getUserBySeatNo(zjh.currSeat)

	if ud == nil {
		mlog.Error("玩家为空")
		return
	}

	addCount := len(zjh.betAdd)

	var arrBet = make([]int32, 0)
	for i := 0; i < addCount; i++ {
		if zjh.bFirst {
			arrBet = append(arrBet, zjh.betAdd[i])
			mlog.Debug("zjh.bFirst")
		} else if zjh.betAdd[i] > zjh.betCurr {
			arrBet = append(arrBet, zjh.betAdd[i])
		}
	}

	opttime := zjh.optTime
	//如果不是场景恢复
	if !bRecover {
		msgall := zjh.getCurrOptMsg()
		zjh.BoardcastMessage(msgall)
		//延迟一秒生效，防止网络卡顿造成弃牌的情况

		//移除操作超时Timer
		if zjh.opttimer != nil {
			zjh.opttimer.Stop()
			zjh.opttimer = nil
		}
		if zjh.currturn >= zjh.maxturn {
			zjh.opttimer = zjh.TimerFuncManager.NewNamedTimer(zjh.optTime+1, zjh.autoCompare, "")
		} else {
			zjh.opttimer = zjh.TimerFuncManager.NewNamedTimer(zjh.optTime+1, zjh.userTimeOutDrop, "")
		}

	} else {
		if zjh.opttimer != nil {
			opttime -= zjh.opttimer.ElapsedSeconds()
		}
	}
	bcheck := false
	if zjh.currturn >= zjh.checkturn {
		bcheck = true
	}
	betadd := (zjh.betAdd[len(zjh.betAdd)-1]) * 2   //注額最大值*看牌倍率
	usercount := (float64(zjh.userCount) - 0.5) * 2 //玩家數量*比牌倍率
	betlimit := int64(betadd) * int64(usercount)    //9/12 betlimit 移除扣底注
	msg := ud.GetOptMsg(arrBet, opttime, zjh.bCompare, zjh.bFirst, bcheck, betlimit, zjh.bAllin)
	zjh.SendMessageToUser(ud.User, msg)
}

//当前操作消息
func (zjh *ZhaJinHuaLogic) getCurrOptMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_CurrOptID)

	pdata := new(netproto.ZJH_WaitOptSeat)
	pdata.Seatno = proto.Int32(zjh.currSeat)
	pdata.OptTime = proto.Int32(zjh.optTime)
	pmsg.MsgData = pdata

	return pmsg
}

//获取玩家数据通过座位号
func (zjh *ZhaJinHuaLogic) getUserBySeatNo(seatno int32) *ZjhUserData {
	for _, urd := range zjh.userData {
		// 這裡很坑, TableNo 要一起檢查才行
		if zjh.TableNo() == urd.User.TableNo && urd.User.SeatNo == seatno {
			return urd
		}
	}

	return nil
}

//处理玩家的比牌消息
func (zjh *ZhaJinHuaLogic) compareCardsByplayer(user *ZjhUserData, data []byte) bool {
	if !zjh.bCompare {
		return false
	}
	ust := &netproto.ZJH_SeatNO{}

	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析用户%v用户座位号错误%v", user, err)
		return false
	}

	dstuser := zjh.getUserBySeatNo(*ust.Seatno)
	if dstuser == nil {
		mlog.Debug("解析用户%v用户座位号不存在", user)
		return false
	}
	if !dstuser.canOpt() {
		return false
	}

	if zjh.opttimer != nil {
		zjh.opttimer.Stop()
		zjh.opttimer = nil
	}

	odd := zjh.getOdds(user, true)

	zjh.GameLog.AddGameStep("在第%d轮用户%d和用户%d比牌，扣钱%d", zjh.currturn, user.User.UserID, dstuser.User.UserID, zjh.betCurr*odd)
	zjh.userCompareOpt(user, dstuser)
	return true
}

//处理玩家比牌消息
func (zjh *ZhaJinHuaLogic) userCompareOpt(user *ZjhUserData, dstuser *ZjhUserData) {
	if !zjh.canOpt {
		mlog.Debug("现在用户还不能操作")
		return
	}
	zjh.setUserMoney(user, true)
	zjh.sendUserCompareResult(user, dstuser)
	zjh.userCount--
	zjh.canOpt = false
	zjh.actionName_curr = zjh.RegisterAnimation(actionName_Compare, zjh.compareAction, zjh.waitOpt)
}

//发送比牌结果
func (zjh *ZhaJinHuaLogic) sendUserCompareResult(user *ZjhUserData, dstuser *ZjhUserData) {
	user.setCompareUserSeatNo(dstuser.User.SeatNo)
	dstuser.setCompareUserSeatNo(user.User.SeatNo)
	comret := comPareCards(user.userCards, dstuser.userCards)

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_CompareRetID)

	pdata := new(netproto.ZJH_CompareCardsRet)

	if comret == 1 {
		mlog.Debug("玩家1赢")
		dstuser.isFaild = true
		pdata.Winseatno = proto.Int32(user.User.SeatNo)
		pdata.Faildseatno = proto.Int32(dstuser.User.SeatNo)
	} else {
		mlog.Debug("玩家2赢")
		user.isFaild = true
		pdata.Winseatno = proto.Int32(dstuser.User.SeatNo)
		pdata.Faildseatno = proto.Int32(user.User.SeatNo)
	}

	odd := zjh.getOdds(user, true)
	pdata.CompareMoney = proto.Int32(zjh.betCurr)
	pdata.Odd = proto.Int32(odd)
	pmsg.MsgData = pdata
	zjh.BoardcastMessage(pmsg)
	//faildUser := zjh.getUserBySeatNo(pdata.GetFaildseatno())
	//zjh.earlyChangeTable(faildUser)
	if dstuser.User.IsOffLine() && dstuser.isFaild {
		zjh.earlyChangeTable(dstuser, "leave")
	}

	if user.User.IsOffLine() && user.isFaild {
		zjh.earlyChangeTable(user, "leave")
	}
}

//发送跟注消息
func (zjh *ZhaJinHuaLogic) sendFollowMsg(user *ZjhUserData) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_FollowBrodCastID)

	odd := zjh.getOdds(user, false)
	pdata := new(netproto.ZJH_Follow)
	pdata.Seatno = proto.Int32(user.User.SeatNo)
	pdata.Follow = proto.Int32(zjh.betCurr)
	pdata.Odd = proto.Int32(odd)
	pmsg.MsgData = pdata
	zjh.BoardcastMessage(pmsg)
}

//发送全押消息
func (zjh *ZhaJinHuaLogic) sendAllinMsg(user *ZjhUserData) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_AllinRetID)

	odd := zjh.getOdds(user, false)
	pdata := new(netproto.ZJH_AllinBrodCast)
	pdata.Seatno = proto.Int32(user.User.SeatNo)
	pdata.Betadd = proto.Int32(zjh.betCurr)
	pdata.Odd = proto.Int32(odd)

	pmsg.MsgData = pdata
	zjh.BoardcastMessage(pmsg)
	mlog.Debug("发送全押消息 userid[%v]", user.User.UserID)
}

//处理用户加注
func (zjh *ZhaJinHuaLogic) userBetAdd(user *ZjhUserData, data []byte) (bool, error) {
	ust := &netproto.ZJH_BetAdd{}

	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析用户%v用户座位号错误%v", user, err)
		return false, nil
	}
	*ust.Betadd = int32(math.Abs(float64(*ust.Betadd)))

	odd := zjh.getOdds(user, false)

	hasvalue := false
	for i := 0; i < len(zjh.betAdd); i++ {
		if zjh.betAdd[i]*odd == *ust.Betadd {
			hasvalue = true
			break
		}
	}

	if !hasvalue {
		mlog.Debug("没有找到这个值：%v", *ust.Betadd)
		return false, util.ErrIsNotBigThenZero
	}

	//加注不能比当前下注还小
	if zjh.betCurr*odd > *ust.Betadd {
		mlog.Debug("加注小于当前值：%v，%v", *ust.Betadd, zjh.betCurr*odd)
		return false, nil
	}

	if !zjh.canOpt {
		mlog.Debug("现在用户还不能操作")
		return false, nil
	}

	zjh.canOpt = false
	if *ust.Betadd > zjh.betCurr*odd {
		zjh.isadd = true
	}
	zjh.betCurr = *ust.Betadd / odd
	zjh.setUserMoney(user, false)

	zjh.GameLog.AddGameStep("在第%d轮用户%d加注%d，扣錢%d", zjh.currturn, user.User.UserID, zjh.betCurr, zjh.betCurr*odd)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_BetAddBrodCastID)

	pdata := new(netproto.ZJH_BetAddBrodCast)
	pdata.Seatno = proto.Int32(user.User.SeatNo)
	pdata.Betadd = proto.Int32(zjh.betCurr)
	pdata.Odd = proto.Int32(odd)
	pdata.Isfirst = proto.Bool(zjh.bFirst)
	pmsg.MsgData = pdata
	zjh.BoardcastMessage(pmsg)

	return true, nil
}

//发送游戏结果
func (zjh *ZhaJinHuaLogic) sendGameRetMsg() {
	//找出赢家和下一次的庄家
	//移除操作超时Timer
	if zjh.opttimer != nil {
		example.invalid("移除超时timer")
		zjh.opttimer.Stop()
		zjh.opttimer = nil
	}

	//分别结算
	for _, urd := range zjh.userData {
		if urd.canOpt() {
			zjh.zhuangSeat = urd.User.SeatNo
			//urd.SetMoney(zjh.betTotal)
			break
		}
	}

	for _, urd := range zjh.userData {
		zjh.calWinMoney(urd)
	}

	for _, urd := range zjh.userData {
		if urd.isAllin {
			break
		}
		zjh.sendComparePoker(urd)
	}

	uscoresRet := make(map[int32]int64)
	for _, urd := range zjh.userData {
		uscoresRet[urd.User.UserID] = urd.win
	}

	zjh.GameLog.GameEnd()

	zjh.broadcastGameRet()

	//zjh.actionName_curr = zjh.RegisterAnimation(actionName_End, zjh.retAction, zjh.sendReadyTime)
	zjh.TimerFuncManager.NewTimer(zjh.retAction, zjh.sendReadyTime)
}

//扑克信息
func (zjh *ZhaJinHuaLogic) getPokerTypeMsg(ud *ZjhUserData) *netproto.ZJH_Poker {
	pokermsg := zjh.getCardsNetMsg(ud)

	zjhpoker := new(netproto.ZJH_Poker)
	zjhpoker.PokerType = proto.Int32(GetType(ud.userCards))
	zjhpoker.Userpokergroup = pokermsg

	return zjhpoker
}

//广播结算消息
func (zjh *ZhaJinHuaLogic) broadcastGameRet() {
	zjh.BoardcastMessage(zjh.getGameRetMsg())
}

//获取游戏结果消息
func (zjh *ZhaJinHuaLogic) getGameRetMsg() *msg.Message {
	zjh.bSendGameResult = true
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_GameResultID)
	pdata := new(netproto.ZJH_GameRet)

	pdata.Money = proto.Int64(zjh.winmoney)
	pdata.WinSeatNO = proto.Int32(zjh.zhuangSeat)
	pmsg.MsgData = pdata
	mlog.Debug("WinSeatNO[%v] winmoney[%v]", zjh.zhuangSeat, zjh.winmoney)
	return pmsg
}

//結算金額
func (zjh *ZhaJinHuaLogic) calWinMoney(ud *ZjhUserData) {
	//己提前結算過
	if ud.isClear {
		return
	}
	mlog.Debug("游戏开始的钱：%v，游戏结束的钱：%v", ud.User.Score, ud.userMoney)
	var userBet, amount, tax, flowAmount, betAmount, payAmount int64
	userBet = ud.userMoney - ud.User.Score //下注的金額 (負數)

	if ud.User.SeatNo == zjh.zhuangSeat { //最後的贏家贏得牌桌上的錢
		amount, tax = zjh.GetUserWinAndTax(zjh.betTotal) //牌桌上的錢需要抽稅
		zjh.winmoney = amount
		amount = amount + userBet // 派彩需扣掉自己下注的金額
	} else { //輸家輸掉下注的錢
		amount = userBet
	}

	flowAmount = util.ABSInt64(amount)      //流水=派彩
	betAmount = util.ABSInt64(amount) + tax //有效投注=輸贏+稅
	payAmount = betAmount + amount          //返邏金額=有效投注+派彩
	ud.User.FlowAmount += flowAmount
	ud.User.BetAmount += betAmount
	ud.User.PayAmount += payAmount
	zjh.GameLog.AddGameStep("用戶%d 結算金額, win %d, 稅收 %d, 流水 %d, 有效投注 %d, 返還金額 %d", ud.User.UserID, amount, tax, flowAmount, betAmount, payAmount)
	//zjh.GameLog.AddGameStep("用戶%d 寫分, win %d, 稅收 %d", ud.User.UserID, win, tax)
	zjh.UpdateUserScoreFunc(ud.User.UserID, amount, tax)

	zjh.GameLog.Context().UserUpdateScore(ud.User.UserID, amount)
	zjh.GameLog.Context().UserBetScore(ud.User.UserID, betAmount)
	zjh.GameLog.Context().UserPayoutScore(ud.User.UserID, payAmount)
	zjh.GameLog.Context().UserValidScore(ud.User.UserID, flowAmount)
	zjh.GameLog.Context().UserTaxAmount(ud.User.UserID, tax)
	ud.isClear = true

	ud.win = amount

	// Added by Joecool Wang at 2019-05-30
	// 若遊戲結果符合［下局輸贏］設定條件，則計完成次數一次並存檔
	stcVipProvider := zjh.SuperToolService.GetZjhSuperToolHandler()

	if ok, _ := util.InArray(ud.User.UserID, stcVipProvider.Proiority); ok { // 檢查 User 是否為監控用戶
		if stcVIPUser := stcVipProvider.GetVipUser(ud.User.UserID); stcVIPUser != nil {
			if stcVIPUser.GetStatus() == vipStatus_Suspend {
				// 控制操作完成，如果有指定牌型，則將指定牌型變更為［未指定］
				if zjh.map_iVIPInitCardType[ud.User.UserID] != vipNotSpecified {
					stcVIPUser.NextType = proto.Int64(vipNotSpecified)
					zjh.SuperToolService.GetZjhSuperToolHandler().SaveVipUser(ud.User.UserID)
					zjh.map_iVIPInitCardType[ud.User.UserID] = vipNotSpecified
				}
			} else {
				// 取輸贏局數的絕對值
				iNextWinLoseCount := stcVIPUser.GetNextWinLoseNum()
				if iNextWinLoseCount < 0 {
					iNextWinLoseCount = -iNextWinLoseCount
				}

				if stcVIPUser.GetFinishNextWinLoseNum() >= iNextWinLoseCount {
					if zjh.map_iVIPInitCardType[ud.User.UserID] != vipNotSpecified {
						// 如果進行局數已達到輸贏局數，且有設定指定牌型，則將指定牌型設為［未指定］，並結束控制
						stcVIPUser.NextType = proto.Int64(vipNotSpecified)
						stcVIPUser.Finish = proto.Int64(vipFinish_Yes)
						stcVIPUser.Status = proto.Int32(vipStatus_Suspend)
						zjh.SuperToolService.GetZjhSuperToolHandler().SaveVipUser(ud.User.UserID)
					}
				} else {
					// 進行局數尚未達到輸贏局數
					if stcVIPUser.GetNextWinLoseNum() < 0 {
						// 指定為輸
						if ud.win < 0 {
							stcVIPUser.FinishNextWinLoseNum = proto.Int64(stcVIPUser.GetFinishNextWinLoseNum() + 1)
							if stcVIPUser.GetFinishNextWinLoseNum() >= iNextWinLoseCount {
								// 進行局數達到輸贏局數，且指定牌型為［未指定］，則結束控制
								if zjh.map_iVIPInitCardType[ud.User.UserID] == vipNotSpecified {
									if stcVIPUser.GetNextType() == vipNotSpecified {
										stcVIPUser.Finish = proto.Int64(vipFinish_Yes)
										stcVIPUser.Status = proto.Int32(vipStatus_Suspend)
									}
								} else {
									if stcVIPUser.GetNextType() == vipNotSpecified {
										zjh.map_iVIPInitCardType[ud.User.UserID] = vipNotSpecified
										stcVIPUser.Finish = proto.Int64(vipFinish_Yes)
										stcVIPUser.Status = proto.Int32(vipStatus_Suspend)
									}
								}
							}
							zjh.SuperToolService.GetZjhSuperToolHandler().SaveVipUser(ud.User.UserID)
						}
					} else if stcVIPUser.GetNextWinLoseNum() > 0 {
						// 指定為贏
						if ud.win > 0 {
							stcVIPUser.FinishNextWinLoseNum = proto.Int64(stcVIPUser.GetFinishNextWinLoseNum() + 1)
							if stcVIPUser.GetFinishNextWinLoseNum() >= iNextWinLoseCount {
								// 進行局數達到輸贏局數，且指定牌型為［未指定］，則結束控制
								if zjh.map_iVIPInitCardType[ud.User.UserID] == vipNotSpecified {
									if stcVIPUser.GetNextType() == vipNotSpecified {
										stcVIPUser.Finish = proto.Int64(vipFinish_Yes)
										stcVIPUser.Status = proto.Int32(vipStatus_Suspend)
									}
								} else {
									if stcVIPUser.GetNextType() == vipNotSpecified {
										zjh.map_iVIPInitCardType[ud.User.UserID] = vipNotSpecified
										stcVIPUser.Finish = proto.Int64(vipFinish_Yes)
										stcVIPUser.Status = proto.Int32(vipStatus_Suspend)
									}
								}
							}
							zjh.SuperToolService.GetZjhSuperToolHandler().SaveVipUser(ud.User.UserID)
						}
					}
				}
			}
		}
	}
	// Added by Joecool Wang at 2019-05-30 end
}

//分别发送和玩家比牌的消息
func (zjh *ZhaJinHuaLogic) sendComparePoker(ud *ZjhUserData) {
	////己提前結算過
	//if ud.isClear {
	//	return
	//}
	//winmoney := ud.userMoney - ud.User.Score //下注的金額
	////mlog.Debug("游戏开始的钱：%v，游戏结束的钱：%v", ud.User.Score, ud.userMoney)
	//if ud.User.SeatNo == zjh.zhuangSeat { //最後的贏家贏得牌桌上的錢
	//	winmoney = zjh.betTotal
	//}
	//
	//money, tax := zjh.GetUserWinAndTax(winmoney) //牌桌上的錢需要抽稅
	//win := money
	//if win > 0 {
	//	win += ud.userMoney - ud.User.Score
	//}
	//zjh.UpdateUserScoreFunc(ud.User.UserID, win, tax)
	//ud.isClear = true
	//if ud.User.SeatNo == zjh.zhuangSeat {
	//	zjh.winmoney = money
	//}
	//
	//ud.win = win

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_ComparePokerID)
	pdata := new(netproto.ZJH_ComparePoker)

	//儲值達500的玩家才可看牌
	if ud.User.SumCharge >= 50000 {
		pdata.OtherPoker = zjh.userinthisgame
	} else {
		//自己的扑克
		//pdata.OtherPoker = append(pdata.OtherPoker, zjh.getPokerTypeMsg(ud))
		//和我比牌人的扑克
		for _, i := range ud.compareUserSeatNo {
			comud := zjh.getUserBySeatNo(i)
			if comud == nil {
				mlog.Debug("和我比牌的人, 己經離開桌子了")
				continue
			}
			pdata.OtherPoker = append(pdata.OtherPoker, zjh.getPokerTypeMsg(comud))
		}
		myself := zjh.getUserBySeatNo(ud.User.SeatNo)
		if myself != nil {
			pdata.OtherPoker = append(pdata.OtherPoker, zjh.getPokerTypeMsg(myself))
		}
	}

	pmsg.MsgData = pdata
	zjh.SendMessageToUser(ud.User, pmsg)
}

//自动找人比牌
func (zjh *ZhaJinHuaLogic) autoCompare() {
	ud := zjh.getUserBySeatNo(zjh.currSeat)
	dstuser := zjh.userData[0]
	for i := 0; i < len(zjh.userData); i++ {
		dstuser = zjh.userData[i]
		if dstuser.User.SeatNo != ud.User.SeatNo && dstuser.canOpt() {
			break
		}
	}

	odd := zjh.getOdds(ud, true)
	zjh.GameLog.AddGameStep("在第%d轮用户%d和用户%d自动比牌，扣钱%d", zjh.currturn, ud.User.UserID, dstuser.User.UserID, zjh.betCurr*odd)
	zjh.userCompareOpt(ud, dstuser)
}

//用户弃牌
func (zjh *ZhaJinHuaLogic) userDrop(user *gamelogic.GameUser) {
	urd := zjh.getUserBySeatNo(user.SeatNo)
	if urd == nil {
		mlog.Error("用户弃牌 Error ")
		return
	}

	if urd.AllInUserId == urd.User.UserID {
		mlog.Error("主動全押的玩家不能棄牌")
		return
	}

	urd.SetDrop()

	zjh.GameLog.AddGameStep("在第%d轮用户%d弃牌", zjh.currturn, urd.User.UserID)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_DropBroadCastID)
	pdata := new(netproto.ZJH_SeatNO)
	pdata.Seatno = proto.Int32(user.SeatNo)
	pmsg.MsgData = pdata

	zjh.BoardcastMessage(pmsg)

	zjh.userCount--
	if zjh.currSeat == user.SeatNo || zjh.userCount <= 1 { //當前操作不是自已，且用戶只剩一個人
		zjh.waitOpt()
	}
	//zjh.earlyChangeTable(urd)

	if urd.User.IsOffLine() {
		zjh.earlyChangeTable(urd, "leave")
	}
}

//处理用户进入
func (zjh *ZhaJinHuaLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {
	//隨機名字
	if zjh.radomnickname == 0 {
		gameuser.NickName = fmt.Sprintf("%d", gameuser.UserID)
	} else {
		gameuser.NickName = zjh.getNickName()
	}
	//隨機頭像

	zjh.DynamicGameLogicBase.UserEnter(gameuser, tableno, seatno)

	//没有庄家游戏还没开始
	if zjh.zhuangSeat == -1 {
		return
	}
	//如果非游戏玩家，发送当前场景
	if !zjh.IsGameUser(gameuser) {
		zjh.sendGameScene(gameuser)
	}
}

//用户断线重连
func (zjh *ZhaJinHuaLogic) UserRecover(gameuser *gamelogic.GameUser) {
	zjh.DynamicGameLogicBase.UserRecover(gameuser)

	//发送游戏场景
	zjh.sendGameScene(gameuser)
}

//发送当前游戏场景
func (zjh *ZhaJinHuaLogic) sendGameScene(user *gamelogic.GameUser) {
	if zjh.bGameEnd {
		mlog.Debug("游戏结束了，不发送场景数据了")
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_ZJHGameSceneID)
	pdata := new(netproto.ZJH_UserInfoArr)
	pdata.CurrOptSeatno = proto.Int32(zjh.currSeat)
	opt := int32(0)
	if zjh.opttimer != nil {
		opt = zjh.opttimer.WaitSeconds - zjh.opttimer.ElapsedSeconds() - 1
	}

	ud := zjh.getUserBySeatNo(user.SeatNo)
	if ud != nil {
		pdata.FollowEnd = proto.Bool(ud.isFollowEnd)
		pdata.Allin = proto.Bool(ud.isAllin)
	} else {
		pdata.FollowEnd = proto.Bool(false)
		pdata.Allin = proto.Bool(false)
	}

	pdata.OptTime = proto.Int32(opt)
	pdata.OptTimeTotal = proto.Int32(zjh.optTime)
	pdata.ZhuangSeatNO = proto.Int32(zjh.zhuangSeat)
	pdata.BetTotal = proto.Int64(zjh.betTotal)
	pdata.Currturn = proto.Int32(zjh.currturn)
	//pdata.AllinStatus = proto.Bool()

	for i := 0; i < len(zjh.betAdd); i++ {
		pdata.BetArr = append(pdata.BetArr, zjh.betAdd[i])
	}

	for _, ud := range zjh.userData {

		userinfo := new(netproto.ZJH_UserInfo)
		if ud.isClear {
			userinfo.Seatno = proto.Int32(-1)
		} else {
			userinfo.Seatno = proto.Int32(ud.User.SeatNo)
		}

		userinfo.Check = proto.Bool(ud.isCheck)
		userinfo.Drop = proto.Bool(ud.isDrop)
		userinfo.Comparefailed = proto.Bool(ud.isFaild)
		betmoney := ud.userMoney - ud.User.Score //下注的金額
		userinfo.BetMoney = proto.Int64(-betmoney)
		pdata.UserInfo = append(pdata.UserInfo, userinfo)
	}

	pmsg.MsgData = pdata
	zjh.SendMessageToUser(user, pmsg)

	if ud != nil {
		if ud.isCheck && ud.canOpt() {
			zjh.SendMessageToUser(user, zjh.getPokerMsg(ud))
		}

		if zjh.userCount <= 1 && zjh.bSendGameResult == true {
			zjh.calWinMoney(ud)
			zjh.sendComparePoker(ud)
			zjh.SendMessageToUser(user, zjh.getGameRetMsg())
		}
	}

	if user.SeatNo == zjh.currSeat && zjh.userCount > 1 {
		zjh.sendOptMsg(true)
	}
}

//广播用户看牌
func (zjh *ZhaJinHuaLogic) broadCastCheckCard(seatno int32) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_CheckCardsBroadCastID)
	pdata := new(netproto.ZJH_SeatNO)
	pdata.Seatno = proto.Int32(seatno)
	pmsg.MsgData = pdata

	zjh.BoardcastMessage(pmsg)
}

//扣用户操作钱
func (zjh *ZhaJinHuaLogic) setUserMoney(urd *ZjhUserData, bCompare bool) {
	odd := zjh.getOdds(urd, bCompare)
	money := zjh.betCurr * odd
	zjh.betTotal += int64(money)
	urd.SetMoney(int64(-money))
}

//获取牌数据消息
func (zjh *ZhaJinHuaLogic) getCardsNetMsg(urd *ZjhUserData) *netproto.UserPoker {
	pokermsg := new(netproto.UserPoker)
	pokermsg.UserID = proto.Int32(urd.User.UserID)
	pokermsg.SeatNo = proto.Int32(urd.User.SeatNo)
	cardslen := len(urd.userCards.Cards)
	for i := 0; i < cardslen; i++ {
		npk := new(netproto.PokerCard)
		npk.Color = proto.Int32(urd.userCards.Cards[i].Color)
		npk.CardNumber = proto.Int32(urd.userCards.Cards[i].CardNumber)
		pokermsg.Cards = append(pokermsg.Cards, npk)
	}

	return pokermsg
}

//发送准备消息
func (zjh *ZhaJinHuaLogic) sendReadyTime() {
	zjh.bGameEnd = true
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_ReadyTime)

	pdata := new(netproto.ZJH_ReadyTime)
	pdata.ReadyTime = proto.Int32(zjh.readyTime)
	pmsg.MsgData = pdata
	zjh.BoardcastMessage(pmsg)
	mlog.Debug("GameEndFunc 发送倒计时 readyTime[%v]", zjh.readyTime)
	zjh.GameEndFunc()
}

//换牌
func (zjh *ZhaJinHuaLogic) ExchangeMaxCard(userid int32) {
	if zjh.disRobotExchangeCard {
		return
	}
	zjh.LogicBase.ExchangeMaxCard(userid)
	var maxRetUserID int32

	var maxcg *poker.CardGroup
	if zjh.poker.(*poker.FixedPoker) == nil {
		example.invalid("*poker.FixedPoker:%v", zjh.poker.(*poker.FixedPoker))
	}

	fpoker := zjh.poker.(*poker.FixedPoker)
	for _, ur := range zjh.AllTableUser() {
		cg := fpoker.GetUserCards(ur.UserID)
		zjh.sortCard(cg)
		if maxcg == nil {
			maxcg = cg
			maxRetUserID = ur.UserID

			continue
		}

		//比较牌型

		if comPareCards(maxcg, cg) != 1 {
			maxcg = cg
			maxRetUserID = ur.UserID
		}
	}

	if maxRetUserID == 0 {
		return
	}

	mlog.Debug("maxcg=%v", maxcg.Cards)
	example.invalid("交换用户%d和用户%d牌型", userid, maxRetUserID)
	zjh.GameLog.AddGameStep("交换用户%d和用户%d牌型", userid, maxRetUserID)
	zjh.disRobotExchangeCard = true //己經有機器人換過就不用再換了
	fpoker.ExChange(userid, maxRetUserID)
}

func (zjh *ZhaJinHuaLogic) VipExchangeCard(userid int32, max int) { //1是最大, 2是最小
	zjh.GameLog.AddGameStep("交換%d交換用戶牌型 %v, (1 最大, 2 最小)", userid, max)
	var retUserID int32

	var retcg *poker.CardGroup
	if zjh.poker.(*poker.FixedPoker) == nil {
		example.invalid("*poker.FixedPoker:%v", zjh.poker.(*poker.FixedPoker))
	}

	fpoker := zjh.poker.(*poker.FixedPoker)
	for _, ur := range zjh.AllTableUser() {
		cg := fpoker.GetUserCards(ur.UserID)
		zjh.sortCard(cg)
		if retcg == nil {
			retcg = cg
			retUserID = ur.UserID
			continue
		}

		//比较牌型

		if comPareCards(retcg, cg) != max {
			retcg = cg
			retUserID = ur.UserID
		}
	}

	if retUserID == 0 {
		return
	}

	example.invalid("retcg=%v", retcg.Cards)
	example.invalid("交换用户%d和用户%d牌型%v", userid, retUserID, retcg)
	zjh.beforeGameStartStep = append(zjh.beforeGameStartStep, fmt.Sprintf("交换用户%d和用户%d牌型%v(1 最大, 2 最小) %v", userid, retUserID, max, retcg))
	fpoker.ExChange(userid, retUserID)
}

func (zjh *ZhaJinHuaLogic) getOdds(ud *ZjhUserData, bCompare bool) int32 {
	odd := int32(1)
	if ud.isCheck {
		odd = odd * 2
	}
	if bCompare {
		odd = odd * 2
	}

	return odd
}

func (zjh *ZhaJinHuaLogic) broadCastCurrturn() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_GameCurrturn)
	pdata := new(netproto.ZJH_Currturn)
	pdata.Currturn = proto.Int32(zjh.currturn)
	//如果當下沒有玩家 回傳空值
	pmsg.MsgData = pdata
	zjh.BoardcastMessage(pmsg)
}

func (zjh *ZhaJinHuaLogic) userTimeOutDrop() {
	urd := zjh.getUserBySeatNo(zjh.currSeat)
	if urd == nil {
		mlog.Error("用户弃牌 Error ")
		return
	}

	urd.SetDrop()

	zjh.GameLog.AddGameStep("在第%d轮用户%d弃牌", zjh.currturn, urd.User.UserID)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_DropBroadCastID)
	pdata := new(netproto.ZJH_SeatNO)
	pdata.Seatno = proto.Int32(zjh.currSeat)
	pmsg.MsgData = pdata

	zjh.BoardcastMessage(pmsg)

	zjh.userCount--
	zjh.waitOpt()
	//zjh.earlyChangeTable(urd)

	if urd.User.IsOffLine() {
		zjh.earlyChangeTable(urd, "leave")
	}
}

//获取扑克消息
func (zjh *ZhaJinHuaLogic) getPokerType(ud *ZjhUserData) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_UserCardType)

	pmsg.MsgData = zjh.getPokerTypeMsg(ud)

	return pmsg
}

func (zjh *ZhaJinHuaLogic) getNickName() string {
	nameString := []string{"大欣becky", "Flomama", "一楠_jyn", "Uni颖儿", "深海贝类的围脖", "雨歌阿", "主持人-王不凡", "周大静_", "卢伟阿", "Mx沫浠", "郑迦文", "杰西酷克", "张大麦和马包包", "王者荣耀神吐槽", "一米小厨娘", "张猫要练嘴皮子", "KatAndSid思凯传媒", "光灿", "达人说", "婉锐", "快美妆BeautyQ", "青藤文化", "papitube", "洋葱视频", "歪果仁研究协会", "橘子娱乐灵药LoveLing ", "黄文煜小盆友", "nG家的猫", "李明霖 ", "-王小强- ", "马克Malik", "铁馆教练", "kanahoooo", "造型师小邱Elvan微小微是你李公子啊", "许墨Zzz", "小威廉的情书", "吕浩男", "林森男神经", "陈吃吃好饿", "智博君李蒲赫", "永远阳光的小宝", "张哲睿bikinir", "la子豪俊平大魔王", "杨霞-Sunny", "王岳鹏Niko", "陈子蜜", "潘雨润PanYR_", "小猪姐姐zz", "扇", "NO_FAN_NO_FUN", "Rika0_0", "Iam方美丽同学", "余潇潇JessicaLinForeverGirl ", "delicious大金", "美美de夏夏啊", "小刘小粒赵大喜", "梦姣--于MOMO", "谢梦Tse", "阿希哥VCRUAN", "滕雨佳Amiu", "Fancystyle", "魏妮妮Yanni大胃王密子君", "美食家大雄", "中华小鸣仔", "锅包肉博士", "香喷喷的小烤鸡", "王大厨", "Amanda", "小厨房", "李子柒", "野食小哥", "办公室小野回忆专用小马甲", "国民老岳父公", "6岛岛", "大锅只是喝醉", "LEO它爹", "大胖大胖", "王白菜-", "酥饼大人", "小H的成长日记", "野老多周玥", "昼夜", "1022女声", "安子与九妹", "沈虫虫", "_胖胖胖_", "李延亮", "阿细-CHEUNG", "奇然liya", "上海彩虹室内合唱团Miss", "Misaya若风lol", "小苍cany", "嗨氏", "dsm大司马解说", "ixm解说小漠", "Esports海涛", "骚男丶注意马叉虫", "靠脸吃饭的徐大王", "起小点同道大叔", "软软其实不太硬", "董新尧", "我的女友嘴很贱", "英国报姐", "王尼玛", "杀马特强子", "小野妹子学吐槽", "追风少年刘全有", "妖妖小精苏小白Alex", "盖达巍", "活力女毒姜允儿", "yero鱼哥", "赵一婷1358", "萌小美Monkey", "冯朗朗", "十三Ace", "一只乔玉鼓山文化", "楼氏传播", "如涵", "宸帆", "蜂群文化", "魔力TV", "大禹网络", "壹加娱乐", "忘忧传媒", "缇苏张沫凡MOMO", "悉尼Cherie", "papi酱", "张大奕eve", "回忆专用小马甲", "同道大叔", "林珊珊_Sunny", "吴大伟DvWooooo", "冯提莫", "陈一发", "周二珂", "章泽天", "莉哥o3o", "摩登兄弟", "一禅小和尚", "会说话的刘二豆", "虎神", "阿哲"}
	//nameString := []string{"资阳市"}
	name := nameString[rand.Intn(len(nameString))]
	return name
}
