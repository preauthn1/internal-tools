package hljbp

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/util"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/global"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// Logic 歡樂聚寶盆主要結構
type Logic struct {
	*gamelogic.DynamicGameLogicBase
	user          []*User       // 玩家資料
	gameState     int           // 遊戲狀態
	nowActionName string        // 遊戲狀態名稱
	luckyNumber   int           // 場上開出數字
	manual        []string      // 遊戲說明
	timeout       map[int]int32 // 狀態運行時間
	chip          []int64       // 籌碼面額
	BetLimit      int64         // 押注上限
	kickUserRound int           // 多少回合沒下注將玩家踢出遊戲，單位回合
}

// NewLogic 創建歡樂聚寶盆邏輯
func NewLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, configs *config.Group) *Logic {
	glc := new(Logic)
	glc.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(table, tfManager, configs)

	glc.GetGameConfig(configs)

	glc.user = make([]*User, glc.TableManager.Room.SeatCount)
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		glc.user[i] = NewUser()
	}

	glc.gameState = StateWaitingStart
	glc.nowActionName = ActionName[StateWaitingStart]
	glc.luckyNumber = 0

	// 遊戲說明

	glc.manual = append(glc.manual, "通过下注获取抽奖机会，下注额越多，中奖率就越高，开奖时，转盘转动，幸运轮随之滚动，转盘停止时，指针停在代表用户的颜色上，拥有该颜色的用户赢取奖池全部筹码")

	glc.RegisterBeforeGameStartFunc(func() {
		//fmt.Printf("BeforeGameStart\n")
		glc.InitData()
	})
	glc.RegisterAfterGameStartFunc(func() {
		//fmt.Printf("AfterGameStart\n")
		glc.Betting()
	})
	glc.RegisterBeforeGameEndFunc(func() {
		//fmt.Printf("BeforeGameEnd\n")
	})
	glc.RegisterAfterGameEndFunc(func() {
		//fmt.Printf("AfterGameEnd\n")

		glc.gameState = StateGameOver

		for _, ur := range glc.AllTableUser() {
			// 用户状态变为准备状态
			glc.DynamicGameLogicBase.UserReady(ur)
			//ur.UserState.OnUserReady(ur, glc.TableManager.Room)
		}
	})

	return glc
}

/*interface*/

// UserEnter d用戶進入
func (glc *Logic) UserEnter(gameUser *gamelogic.GameUser, tableno int32, seatno int32) {
	//fmt.Printf("UserEnter %d\n", gameUser.UserID)

	if gameUser == nil {
		return
	}

	mlog.Debug("TableNo[%d], UserEnter Guest, User[%d]", glc.TableNo(), gameUser.UserID)

	//fmt.Printf("UserEnter OK %d\n", gameUser.UserID)

	glc.DynamicGameLogicBase.UserEnter(gameUser, tableno, seatno)

	glc.AddUser(gameUser)

	// UserReady 擺這 一但人數足夠觸發遊戲開始
	glc.DynamicGameLogicBase.UserReady(gameUser)

	if glc.IsGaming() {
		gameUser.UserState.OnUserGame(gameUser, glc.TableManager.Room)
	}

	//fmt.Printf("all player count: %d\n", glc.GetUserCount())
	//fmt.Printf("player1: %+v\n", glc.user[0].gameUser)
	//fmt.Printf("player2: %+v\n", glc.user[1].gameUser)
}

// UserFree 用戶站立
func (glc *Logic) UserFree(gameUser *gamelogic.GameUser) {
	//fmt.Printf("UserFree %d\n", gameUser.UserID)
	mlog.Debug("UserFree%d", gameUser.UserID)

	if gameUser == nil {
		return
	}

	mlog.Debug("Table[%d], UserFree, User[%d], Seat[%d], IsOffLine()[%t]", glc.TableNo(), gameUser.UserID, gameUser.SeatNo, gameUser.IsOffLine())

	if gameUser.SeatNo != gamelogic.FreeSeatNo {
		user := glc.user[gameUser.SeatNo]

		if user.totalBet > 0 && glc.IsGaming() {
			return
		}

		gameUser.UserState.OnGameEnd(gameUser, glc.TableManager.Room)

		if gameUser.SeatNo >= 0 && gameUser.SeatNo < glc.TableManager.Room.SeatCount {
			glc.user[int(gameUser.SeatNo)].Init()
		}
	}

	mlog.System("UserFree %s", gameUser)

	glc.SendAllUserInfo()

	glc.LogicBase.UserFree(gameUser)
	//fmt.Printf("UserFree OK %d\n", gameUser.UserID)

	//fmt.Printf("all player count: %d\n", glc.GetUserCount())
	//fmt.Printf("player1: %+v\n", glc.user[0].gameUser)
	//fmt.Printf("player2: %+v\n", glc.user[1].gameUser)
}

// UserRecover 用户斷線重連
func (glc *Logic) UserRecover(gameUser *gamelogic.GameUser) {
	//fmt.Printf("UserRecover %d\n", gameUser.UserID)
	mlog.Debug("UserRecover %d", gameUser.UserID)

	if gameUser == nil {
		return
	}

	mlog.Debug("Table[%d], UserRecover, User[%d], Seat[%d], IsOffLine()[%t]", glc.TableNo(), gameUser.UserID, gameUser.SeatNo, gameUser.IsOffLine())

	glc.DynamicGameLogicBase.UserRecover(gameUser)

	glc.SendGameInfo(gameUser)
	glc.SendAllUserInfo()
	//fmt.Printf("UserRecover OK %d\n", gameUser.UserID)
}

// UserLeave 用戶離開
func (glc *Logic) UserLeave(gameUser *gamelogic.GameUser) bool {
	//fmt.Printf("UserLeave %d\n", gameUser.UserID)
	mlog.Debug("UserLeave %d", gameUser.UserID)
	if gameUser == nil {
		return false
	}

	mlog.Debug("Table[%d], UserLeave, User[%d], Seat[%d], IsOffLine()[%t]", glc.TableNo(), gameUser.UserID, gameUser.SeatNo, gameUser.IsOffLine())

	if gameUser.SeatNo != gamelogic.FreeSeatNo {
		pUser := glc.user[gameUser.SeatNo]

		if pUser.totalBet > 0 && glc.IsGaming() {
			return false
		}

		gameUser.UserState.OnGameEnd(gameUser, glc.TableManager.Room)

		if gameUser.SeatNo >= 0 && gameUser.SeatNo < glc.TableManager.Room.SeatCount {
			glc.user[int(gameUser.SeatNo)].Init()
		}
	}

	glc.SendAllUserInfo()

	//fmt.Printf("UserLeave OK %d\n", gameUser.UserID)
	return glc.DynamicGameLogicBase.UserLeave(gameUser)
}

// GameMessage 遊戲封包
func (glc *Logic) GameMessage(gameUser *gamelogic.GameUser, msg *msg.Message) {
	if gameUser == nil || msg == nil {
		return
	}

	if msg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	user := glc.GetUserData(gameUser.UserID)

	if user == nil {
		return
	}

	glc.CheckUserAnimationCommited(gameUser, msg, glc.nowActionName, nil)

	switch msg.SClassID {
	case int32(netproto.HLJBP_GameMessageClassID_HLJBP_Bet): // C2S 玩家下注
		glc.MessageBet(user, msg)
	}
}

/*Interface*/

/*GameLogic*/

/**遊戲邏輯**/

// GetGameConfig 讀取遊戲設定
func (glc *Logic) GetGameConfig(configs *config.Group) {
	mlog.Debug("GetGameConfig")
	if configs == nil {
		mlog.Error("桌號[%d] 讀取 config 失敗", glc.TableNo())
		return
	}

	var str string = ""

	// 讀取 gameserver 設定值
	cfgs := configs.GetConfig(gamelogic.ConfigFileNameGameServer)
	if cfgs == nil {
		mlog.Error("獲取配置 [gameserver] 失敗")
	} else {
		//----- timeout -----
		glc.timeout = make(map[int]int32)
		glc.timeout[StateBetting] = cfgs.GetConfigInt32(ConfigLabel, TimeOut[StateBetting])
		glc.timeout[StateSettlement] = cfgs.GetConfigInt32(ConfigLabel, TimeOut[StateSettlement])

		// 籌碼面額
		str = cfgs.GetConfigValue("GameSetting", "ChipValue")
		strChip := strings.Split(str, ",")
		for i := 0; i < len(strChip); i++ {
			i64chip, _ := strconv.ParseInt(strChip[i], 10, 32)
			glc.chip = append(glc.chip, i64chip)
		}

		// 押注上限
		glc.BetLimit = cfgs.GetConfigInt64(ConfigLabel, "BetLimit")

		// 踢人參數
		glc.kickUserRound = cfgs.GetConfigInt(ConfigLabel, "KickUserRound")
	}
}

// InitData 初始化遊戲參數
func (glc *Logic) InitData() {
	//fmt.Printf("InitData\n")
	mlog.Debug("InitData")

	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		gameUser := glc.GetUserBySeatNo(i)
		if gameUser == nil {
			glc.user[i].Init()
		} else {
			glc.user[i].ResetGameData()
			glc.user[i].state = UserStateNormal
		}
	}

	glc.gameState = StateWaitingStart
	glc.nowActionName = ActionName[StateWaitingStart]
}

// AddUser 加入新玩家
func (glc *Logic) AddUser(gameUser *gamelogic.GameUser) {
	if gameUser == nil {
		return
	}

	seatNo := gameUser.SeatNo
	if seatNo < 0 || seatNo >= glc.TableManager.Room.SeatCount {
		return
	}

	user := glc.user[seatNo]
	if user.state == UserStateEmptySeat {
		user.Init()
		user.gameUser = gameUser
		user.money = gameUser.Score
		user.ResetGameData()

		if glc.gameState > StateBetting {
			user.state = UserStateWatch
		} else {
			user.state = UserStateNormal
		}
	} else {
		mlog.Error("AddUser is not empty, %+v", user)
	}

	glc.SendGameInfo(gameUser)
	glc.SendAllUserInfo()
	example.invalid("玩家入座[%d], 遊戲狀態[%s], 座位[%d], 類別[%d], 玩家狀態[%s], 玩家總人數[%d]",
		user.gameUser.UserID,
		StrGameStateDebug[glc.gameState],
		seatNo,
		user.gameUser.UserType,
		StrUserStateDebug[user.state],
		glc.TableRealUserCount())

	glc.GameLog.AddGameStep("玩家[%d] 入座，狀態：%s", gameUser.UserID, StrUserStateDebug[user.state])
}

// SendGameInfo 傳送牌局資訊
func (glc *Logic) SendGameInfo(gameUser *gamelogic.GameUser) {
	if gameUser == nil {
		return
	}

	mlog.Debug("sendGameInfo, GameState[%s]",
		StrGameStateDebug[glc.gameState])

	message := new(msg.Message)
	message.BClassID = int32(netproto.MessageBClassID_Game)
	message.SClassID = int32(netproto.HLJBP_GameMessageClassID_HLJBP_GameInfo)

	data := new(netproto.HLJBP_Msg_GameInfo)

	// 遊戲說明
	for i := 0; i < len(glc.manual); i++ {
		data.Manual = append(data.Manual, glc.manual[i])
	}

	// 所有玩家資訊
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		userInfo := glc.UserInfoPkg(i)
		data.UserInfo = append(data.UserInfo, userInfo)
	}

	// 籌碼面額
	for i := 0; i < len(glc.chip); i++ {
		data.ChipValue = append(data.ChipValue, glc.chip[i])
	}

	data.Bet_Limit = proto.Int64(glc.BetLimit)
	data.GameState = proto.Int32(int32(glc.gameState))
	data.TimeRemain = proto.Int32(glc.timeout[glc.gameState] - glc.PlayerActionManager.GetActionTimerElapsedSeconds(glc.nowActionName))
	data.LuckyNumber = proto.Int32(int32(glc.luckyNumber))

	message.MsgData = data
	glc.SendMessageToUser(gameUser, message)

	switch glc.gameState {
	case StateBetting: // 押注
	case StateSettlement: // 結算
	}
}

// SendAllUserInfo 傳送玩家資訊
func (glc *Logic) SendAllUserInfo() {
	message := new(msg.Message)
	message.BClassID = int32(netproto.MessageBClassID_Game)
	message.SClassID = int32(netproto.HLJBP_GameMessageClassID_HLJBP_AllUserInfo)

	data := new(netproto.HLJBP_Msg_AllUserInfo)

	for seat := int32(0); seat < glc.TableManager.Room.SeatCount; seat++ {
		pUserInfo := glc.UserInfoPkg(seat)
		data.UserInfo = append(data.UserInfo, pUserInfo)
	}

	message.MsgData = data

	glc.BoardcastMessage(message)
}

// UserInfoPkg 玩家資訊打包資料，seat：此座位的玩家資訊
func (glc *Logic) UserInfoPkg(seat int32) *netproto.HLJBP_Msg_UserInfo {
	if seat < 0 || seat >= glc.TableManager.Room.SeatCount {
		return nil
	}

	user := glc.user[seat]

	userInfo := new(netproto.HLJBP_Msg_UserInfo)

	if user.gameUser != nil {
		userInfo.UserID = proto.Int32(user.gameUser.UserID) // 玩家ID，若是空位時，玩家ID = -1
		userInfo.FaceID = proto.Int32(user.gameUser.FaceID) // 玩家頭像ID
	} else {
		userInfo.UserID = proto.Int32(-1)
		userInfo.FaceID = proto.Int32(0) // 玩家頭像ID
	}
	userInfo.State = proto.Int32(int32(user.state))                  // 玩家狀態，0:空位，1:正常，2:觀戰
	userInfo.SeatNo = proto.Int32(seat + int32(ClientStartUpSeatNo)) // 玩家座位，client 座位號從 1 開始
	userInfo.Money = proto.Int64(user.money)                         // 玩家金錢
	userInfo.Bet = proto.Int64(user.totalBet)                        // 總押注額

	userInfo.Winner = proto.Int32(int32(user.winState))    // 玩家輸贏
	userInfo.WinLoseMoney = proto.Int64(user.winLoseMoney) // 玩家輸贏金錢(不扣稅收)，若值為負數，則為玩家輸掉的金額，client 端顯示用
	userInfo.TotalWinLose = proto.Int64(user.totalWinLose) // 總輸贏(已扣稅收)，client 端顯示用

	return userInfo
}

// SelfInfoPkg 玩家自己資料
func (glc *Logic) SelfInfoPkg(gameUser *gamelogic.GameUser) *netproto.HLJBP_Msg_UserInfo {

	seat := glc.GetUserSeat(gameUser.UserID)
	userInfo := new(netproto.HLJBP_Msg_UserInfo)

	if seat >= 0 || seat < glc.TableManager.Room.SeatCount {
		user := glc.user[seat]

		if user.gameUser != nil {
			userInfo.UserID = proto.Int32(user.gameUser.UserID) // 玩家ID，若是空位時，玩家ID = -1
			userInfo.FaceID = proto.Int32(user.gameUser.FaceID) // 玩家頭像ID
		} else {
			userInfo.UserID = proto.Int32(-1)
			userInfo.FaceID = proto.Int32(0) // 玩家頭像ID
		}
		userInfo.State = proto.Int32(int32(user.state))                  // 玩家狀態，0:空位，1:正常，2:觀戰，3:棄牌
		userInfo.SeatNo = proto.Int32(seat + int32(ClientStartUpSeatNo)) // 玩家座位，client 座位號從 1 開始
		userInfo.Money = proto.Int64(user.money)                         // 玩家金錢
		userInfo.Bet = proto.Int64(user.totalBet)                        // 總押注額

		userInfo.Winner = proto.Int32(int32(user.winState))    // 玩家輸贏
		userInfo.WinLoseMoney = proto.Int64(user.winLoseMoney) // 玩家輸贏金錢(不扣稅收)，若值為負數，則為玩家輸掉的金額，client 端顯示用
		userInfo.TotalWinLose = proto.Int64(user.totalWinLose) // 總輸贏(已扣稅收)，client 端顯示用

		return userInfo
	}

	return nil
}

// IsNoBet 確認是否無人下注
func (glc *Logic) IsNoBet() bool {
	TotalBet := int64(0)
	for seat := int32(0); seat < glc.TableManager.Room.SeatCount; seat++ {
		user := glc.user[seat]

		if user == nil {
			continue
		}

		if user.state == UserStateNormal {
			TotalBet += user.totalBet
		}
	}
	if TotalBet > 0 {
		return false
	}

	return true
}

// KickTipMessage 踢觀戰提示訊息
func (glc *Logic) KickTipMessage(kickRound int32) *msg.Message {
	message := new(msg.Message)
	message.BClassID = int32(netproto.MessageBClassID_Game)
	message.SClassID = int32(netproto.HLJBP_GameMessageClassID_HLJBP_KickTip)

	data := new(netproto.HLJBP_Msg_KickTip)
	data.Round = proto.Int32(kickRound)
	message.MsgData = data

	return message
}

/**遊戲流程**/

// Betting 開始下注
func (glc *Logic) Betting() {
	//fmt.Printf("Betting\n")
	glc.GameLog.AddGameStep("系統水池%d, 贏錢線%d", glc.SysPool.GetPool(), glc.SysPool.GetWinLine())
	mlog.Debug("Betting")
	glc.GameLog.AddGameStep("開始下注")
	glc.gameState = StateBetting
	glc.nowActionName = glc.PlayerActionManager.NewAction(ActionName[glc.gameState], 0, glc.timeout[glc.gameState], func() {
		glc.Settlement()
	})

	// 傳送玩家狀態 觀戰或正在玩
	glc.SendAllUserInfo()

	message := new(msg.Message)
	message.BClassID = int32(netproto.MessageBClassID_Game)
	message.SClassID = int32(netproto.HLJBP_GameMessageClassID_HLJBP_UserBettingStart)
	pData := new(netproto.HLJBP_Msg_UserBettingStart)
	pData.TimeRemain = proto.Int32(glc.timeout[glc.gameState] - glc.PlayerActionManager.GetActionTimerElapsedSeconds(glc.nowActionName))

	message.MsgData = pData

	glc.BoardcastMessage(message)
}

// Settlement 總結算
func (glc *Logic) Settlement() {
	//fmt.Printf("Settlement\n")
	mlog.Debug("Settlement")
	glc.GameLog.AddGameStep("總结算")
	glc.gameState = StateSettlement
	glc.nowActionName = glc.RegisterAnimationByUserCount(ActionName[glc.gameState], int32(glc.TableRealUserCount()), glc.timeout[glc.gameState], func() {
		glc.GameOver()
	})

	robotUsers := make([]*User, 0)
	for seat := int32(0); seat < glc.TableManager.Room.SeatCount; seat++ {
		user := glc.user[seat]
		if user.state == UserStateNormal {
			if user.totalBet <= 0 {
				user.state = UserStateWatch
				user.noBetRound++ // 踢人-未操作回合增加
			} else {
				user.noBetRound = 0 // 踢人-有操作回合重置
				if user.gameUser.IsRobot() {
					robotUsers = append(robotUsers, user)
				}
			}

			if user.noBetRound >= glc.kickUserRound && !user.gameUser.IsRobot() {
				// 踢人-未操作回合到
				glc.KickUserFunc(user.gameUser.UserID, "您即将返回大厅")
			} else if user.noBetRound >= glc.kickUserRound-1 && !user.gameUser.IsRobot() {
				// 踢人-未操作回合到 提示剩一回合
				glc.SendMessageToUser(user.gameUser, glc.KickTipMessage(1))
			}
		}
	}

	IsSettlement := false
	// 有人下注才結算
	if !glc.IsNoBet() {
		// 結算
		IsSettlement = true

		allUserBet := int64(0)
		for i := 0; i < len(glc.user); i++ {
			if glc.user[i].state == UserStateNormal {
				allUserBet += glc.user[i].totalBet
			}
		}

		mlog.Debug("總押注: %s", global.MoneyString(allUserBet))
		glc.GameLog.AddGameStep("總押注: %s", global.MoneyString(allUserBet))

		rnd := rand.Int63n(allUserBet)

		tempValue := int64(0)
		checkWinner := false
		userCount := 0
		// 水池用 贏的玩家ID
		cheatUserID := int32(0)
		if len(robotUsers) > 0 && glc.SysPool.IsCtrlSystemWin() {
			randomIndex := rand.Intn(len(robotUsers))
			cheatUserID = robotUsers[randomIndex].gameUser.UserID
			glc.GameLog.AddGameStep("水線之下, 指定機器人贏:%d", cheatUserID)
		}

		// 水池用 累積真人輸錢
		PoolAmountUser := int64(0)
		// 水池用 累積機器人輸錢
		PoolAmountRobot := int64(0)
		// 水池用 紀錄贏的人是機器人
		winnerIsRobot := false

		for i := 0; i < len(glc.user); i++ {
			user := glc.user[i]
			if user.state == UserStateNormal {
				userCount++

				// 水池用，若啟動作弊機制，則選中的人必贏
				if cheatUserID == 0 {
					tempValue += user.totalBet
				} else {
					if cheatUserID == user.gameUser.UserID {
						tempValue = allUserBet
					} else {
						tempValue = 0
					}
				}

				winMoney := int64(0)
				tax := int64(0)
				if tempValue > rnd && !checkWinner {
					// 贏
					checkWinner = true
					winMoney = allUserBet - user.totalBet
					winMoney, tax = glc.GetUserWinAndTax(winMoney)

					user.money += winMoney + user.totalBet
					user.moneyChange = winMoney
					user.winLoseMoney = winMoney + tax
					user.totalWinLose = winMoney + user.totalBet
					user.taxAmount = tax

					user.winState = WinStateWin

					if user.gameUser.IsRobot() {
						winnerIsRobot = true
					}
				} else {
					// 輸
					winMoney = 0 - user.totalBet

					// money 輸的錢在押注時扣過了
					user.moneyChange = winMoney
					user.winLoseMoney = winMoney
					// totalWinLose 輸的錢在押注時扣過了

					user.winState = WinStateLose

					// 水池用 累積輸錢
					if user.gameUser.IsRobot() {
						PoolAmountRobot += user.totalBet
					} else {
						PoolAmountUser += user.totalBet
					}
				}

				glc.UpdateUserScoreFunc(user.gameUser.UserID, user.moneyChange, tax)
				// UserUpdateScore 会与GameOver 时的 UserUpdateScore 累加一起, 这边这个多做的完全不需要
				//glc.GameLog.Context().UserUpdateScore(user.gameUser.UserID, tax)

				var flowAmount, betAmount, payAmount int64
				flowAmount = util.ABSInt64(user.moneyChange)
				betAmount = user.totalBet
				payAmount = betAmount + user.moneyChange

				user.betAmount = user.totalBet
				user.winAmount = winMoney + user.totalBet + tax
				if winMoney < 0 {
					user.flowAmount = user.totalBet
				} else {
					user.flowAmount = 0
				}

				user.gameUser.FlowAmount += flowAmount
				user.gameUser.BetAmount += betAmount
				user.gameUser.PayAmount += payAmount

				str := fmt.Sprintf("玩家[%d] 押注 %s, 中獎 %s, 稅收 %s, 輸贏 %s, 結算後金錢 %s, 結算流水 %s, 有效投注 %s, 返還金額 %s",
					user.gameUser.UserID,
					global.MoneyString(user.totalBet),
					global.MoneyString(user.winAmount),
					global.MoneyString(tax),
					global.MoneyString(user.moneyChange),
					global.MoneyString(user.money),
					global.MoneyString(user.flowAmount),
					global.MoneyString(user.totalBet),
					global.MoneyString(payAmount))

				//fmt.Printf("%s\n", str)
				mlog.Debug("%s", str)
				glc.GameLog.AddGameStep("%s", str)
			}
		}
		if winnerIsRobot {
			changeSysPool, _ := glc.GetUserWinAndTax(PoolAmountUser)
			example.invalid("系统水池異動, from %d, change %d", glc.SysPool.GetPool(), changeSysPool)
			glc.GameLog.AddGameStep("系统水池異動, from %d, change %d", glc.SysPool.GetPool(), changeSysPool)
			glc.SysPool.UpdateSysPool(changeSysPool)
		} else {
			changeSysPool := -PoolAmountRobot
			example.invalid("系统水池異動, from %d, change %d", glc.SysPool.GetPool(), changeSysPool)
			glc.GameLog.AddGameStep("系统水池異動, from %d, change %d", glc.SysPool.GetPool(), changeSysPool)
			glc.SysPool.UpdateSysPool(changeSysPool)
		}

		// 大於1個玩家 亂數產生數字暫存
		if userCount > 1 {
			glc.luckyNumber = rand.Intn(10000)
		}
	}

	//----- 傳送訊息給 client -----

	message := new(msg.Message)
	message.BClassID = int32(netproto.MessageBClassID_Game)
	message.SClassID = int32(netproto.HLJBP_GameMessageClassID_HLJBP_Settlement)

	data := new(netproto.HLJBP_Msg_Settlement)
	data.IsSettlement = proto.Bool(IsSettlement)
	data.LuckyNumber = proto.Int32(int32(glc.luckyNumber))

	// 所有玩家資訊
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		pUserInfo := glc.UserInfoPkg(i)
		data.UserInfo = append(data.UserInfo, pUserInfo)
	}

	data.ActionName = proto.String(ActionName[glc.gameState])

	message.MsgData = data
	glc.BoardcastMessage(message)
}

// GameOver 遊戲結束
func (glc *Logic) GameOver() {
	//fmt.Printf("GameOver\n")
	mlog.Debug("GameOver")
	glc.GameLog.AddGameStep("游戏结束")

	glc.gameState = StateGameOver

	scoreRets := make(map[int32]int64)
	for seat := int32(0); seat < glc.TableManager.Room.SeatCount; seat++ {
		user := glc.user[seat]
		if user.gameUser != nil {
			scoreRets[user.gameUser.UserID] = user.moneyChange

			glc.GameLog.Context().UserUpdateScore(user.gameUser.UserID, user.moneyChange)
			glc.GameLog.Context().UserTaxAmount(user.gameUser.UserID, user.taxAmount)
			glc.GameLog.Context().UserBetScore(user.gameUser.UserID, user.betAmount)
			glc.GameLog.Context().UserPayoutScore(user.gameUser.UserID, user.winAmount)
			glc.GameLog.Context().UserValidScore(user.gameUser.UserID, user.flowAmount)

		}
	}

	glc.GameLog.GameEnd()

	glc.GameEndFunc()
}

/**封包處理**/

// MessageBet 訊息處理 - 玩家下注
func (glc *Logic) MessageBet(user *User, message *msg.Message) {
	if user == nil || message == nil {
		return
	}

	msgInfo := &netproto.HLJBP_Msg_Bet{}
	proto.Unmarshal(message.MsgData.([]byte), msgInfo)
	bet := msgInfo.GetBet()

	example.invalid("玩家[%d] 下注 %s", user.gameUser.UserID, global.MoneyString(bet))
	glc.GameLog.AddGameStep("玩家[%d] 下注 %s", user.gameUser.UserID, global.MoneyString(bet))

	if glc.gameState != StateBetting {
		glc.SendBetResult(user, bet, BetResultBettingStopped)
		return
	}

	result := glc.UserActionBet(user, bet)
	glc.SendBetResult(user, bet, result)
}

// UserActionBet 玩家下注
func (glc *Logic) UserActionBet(user *User, bet int64) int {
	if user == nil {
		return BetResultBetFailed
	}

	result := BetResultSuccess

	if bet <= 0 {
		return BetResultBetSmallThanZero
	}

	// 剩餘金錢是否足以下注
	if user.money < bet {
		return BetResultNotEnoughMoney
	}

	// 下注是否超過Config上限
	if user.totalBet+bet > glc.BetLimit && glc.BetLimit > 0 {
		return BetResultBetOverUpperLimit
	}

	user.totalBet += bet
	user.totalWinLose -= bet

	user.money -= bet
	return result
}

// SendBetResult 傳送玩家下注結果
func (glc *Logic) SendBetResult(user *User, bet int64, result int) {
	if user == nil {
		mlog.Debug("sendBetResult(), 無玩家資料")
		return
	}

	mlog.Debug("sendBetResult(), User[%d], SeatNo[%d], bet[%d], result[%d]", user.gameUser.UserID, user.gameUser.SeatNo, bet, result)

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.HLJBP_GameMessageClassID_HLJBP_BetResult)

	pData := new(netproto.HLJBP_Msg_BetResult)
	pData.SeatNo = proto.Int32(user.gameUser.SeatNo + int32(ClientStartUpSeatNo)) // 下注玩家座位
	pData.Bet = proto.Int64(bet)                                                  // 下注金額
	pData.MoneyRemain = proto.Int64(user.money)                                   // 下注後玩家剩餘金錢
	pData.Result = proto.Int32(int32(result))                                     // 下注結果

	pMsg.MsgData = pData
	if result == BetResultSuccess {
		glc.BoardcastMessage(pMsg)
	} else {
		glc.SendMessageToUser(user.gameUser, pMsg)
	}

	//fmt.Printf("玩家[%d] 下注 %s, 结果：%s, 注额 %s, 剩余金钱 %s \n",
	//	user.gameUser.UserID,
	//	global.MoneyString(bet),
	//	StrBetResult[result],
	//	global.MoneyString(user.totalBet),
	//	global.MoneyString(user.money))
	mlog.Debug("玩家[%d] 下注 %s, 结果：%s, 注額 %s, 剩餘金錢 %s",
		user.gameUser.UserID,
		global.MoneyString(bet),
		StrBetResult[result],
		global.MoneyString(user.totalBet),
		global.MoneyString(user.money))
	glc.GameLog.AddGameStep("玩家[%d] 下注 %s, 结果：%s, 注額 %s, 剩餘金錢 %s",
		user.gameUser.UserID,
		global.MoneyString(bet),
		StrBetResult[result],
		global.MoneyString(user.totalBet),
		global.MoneyString(user.money))

}

/*GameLogic*/
