package jlbpk

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

type JlbpkUser struct {
	User                     *gamelogic.GameUser
	i32UserID                int32 // 玩家ID
	i32SeatNo                int32 // 座位
	iUserType                int   // 玩家類別
	i64Money                 int64 // 玩家金錢
	iUserState               int   // 玩家狀態
	i64BetBase               int64 // 玩家底注
	i64BetSide               int64 // 玩家邊注
	i64ZhuangBetSide         int64 // 莊家邊注
	i64BetDouble             int64 // 玩家加注
	bGetJackpot              bool  // 玩家是否中 Jackpot
	i64JPMoney               int64 // 玩家所中 Jackpot 金額
	bGetZhuangJackpot        bool  // 玩家是否中 莊家Jackpot
	i64ZhuangJPMoney         int64 // 玩家所中 莊家Jackpot 金額
	i64WinLoseMoney          int64 // 玩家輸贏金錢，若值為負數，則為玩家輸掉的金額
	i64BetBaseWinLoseMoney   int64 // 玩家底注輸贏金錢，若值為負數，則為玩家輸掉的金額
	i64BetDoubleWinLoseMoney int64 // 玩家加注輸贏金錢，若值為負數，則為玩家輸掉的金額
	iPlayingState            int   // 遊戲中狀態
	i64Timeout               int64 // 各階段 tumeout
	i64LastMessageTime       int64 // 上一次訊息時間
	i64MoneyChange           int64 // 金錢變化
	i64TotalWinLose          int64 // 總輸贏，client 端顯示用
	iWinState                int   // 玩家輸贏
	i64TotalBet              int64 // 總押注額
	i64MoneyBeforebet        int64 // 押注前金額
	bOffline                 bool  // 玩家是否為離線
	i64MaxBetBase            int64 // 可以下的最大底注值
	bIsVIPFinished           bool  // 是否已結束監控
	i32Weight                int32 // 權重
	i64TargetPoint           int64 // 目標金額
	i64AccumulatedPoint      int64 // 已累積目標金額
	i64Tax                   int64 // 累積稅額
}

func NewJlbpkUser() *JlbpkUser {
	pUser := new(JlbpkUser)
	pUser.Init()

	return pUser
}

// 初始化玩家資料
func (pUser *JlbpkUser) Init() {
	pUser.User = nil
	pUser.iUserType = gamelogic.UserTypeNormal
	pUser.i32UserID = -1
	pUser.i64Money = 0
	pUser.iUserState = JLB_UserState_EmptySeat
	pUser.bOffline = false
	pUser.bIsVIPFinished = true
	pUser.i32Weight = g_i32PlayerWeightDefault
	pUser.i64TargetPoint = 0
	pUser.i64AccumulatedPoint = 0

	pUser.ResetGameData()
}

// 重置遊戲資料
func (pUser *JlbpkUser) ResetGameData() {
	pUser.i64BetBase = 0
	pUser.i64BetSide = 0
	pUser.i64ZhuangBetSide = 0
	pUser.i64BetDouble = 0
	pUser.bGetJackpot = false
	pUser.i64JPMoney = 0
	pUser.bGetZhuangJackpot = false
	pUser.i64ZhuangJPMoney = 0
	pUser.i64WinLoseMoney = 0
	pUser.i64BetBaseWinLoseMoney = 0
	pUser.i64BetDoubleWinLoseMoney = 0
	pUser.iPlayingState = JLB_UserPlayState_WaitingStart
	pUser.i64Timeout = 0
	pUser.i64LastMessageTime = 0
	pUser.i64MoneyChange = 0
	pUser.i64TotalWinLose = 0
	pUser.iWinState = JLB_WinState_None
	pUser.i64TotalBet = 0
	pUser.i64MaxBetBase = pUser.i64Money / 3
	pUser.i64Tax = 0
}

// 設定機器人 timeout
func (pUser *JlbpkUser) SetRobotTimeout(i64Timeout int64) {
	pUser.i64Timeout = 0

	if i64Timeout > 0 {
		if pUser.iUserType == gamelogic.UserTypeRebot && pUser.iUserState == JLB_UserState_Normal {
			iRnd := util.RndInt(1, int(i64Timeout))
			pUser.i64Timeout = int64(iRnd)
		}
	}
}
