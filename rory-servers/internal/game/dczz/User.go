package dczz

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	Pokerset "example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

type User struct {
	User               *gamelogic.GameUser
	i32UserID          int32                  // 玩家ID
	i32SeatNo          int32                  // 座位
	iUserType          int                    // 玩家類別
	i64Money           int64                  // 玩家金錢
	iUserState         int                    // 玩家狀態
	i64Timeout         int64                  // 各階段 tumeout
	i64LastMessageTime int64                  // 上一次訊息時間
	pHandCard          *Pokerset.PokerCardSet // 玩家手牌
	i64MoneyChange     int64                  // 金錢變化
	i64TotalWinLose    int64                  // 總輸贏，client 端顯示用
	iPlayingState      int                    // 遊戲中狀態
	bOffline           bool                   // 玩家是否為離線
	pRobot             *RobotSetting          // 機器人設定
	i64Bet_Base        int64                  // 底注
	i64Bet_Side        int64                  // 平手注
	i64Bet_Double      int64                  // 戰爭階段加注
	i64FoldRecover     int64                  // 於戰爭階段選擇棄牌收回的注碼
	bIsTied            bool                   // 是否跟莊家平手
	bGetTieAward       bool                   // 玩家是否於一般階段贏得平手賠付
	bCanEnterWar       bool                   // 可否進入戰爭階段
	bIsFolded          bool                   // 玩家是否已棄牌
	iRobotBetTimes     int                    // 已下注次數，用於機器人
	i32Point           int32                  // 點數
	bOperationDone     bool                   // 操作是否完成
	i64TieAward        int64                  // 平手賠付
	iSettlementWinLose int                    // 玩家輸贏
	i64Tax             int64                  // 稅
	bGoForWar          bool                   // 加入戰爭
	iWeight            int                    // 權重
	totalBet           int64                  // 總下注金額, 記錄下來
}

func NewUser() *User {
	pUser := new(User)

	pUser.pHandCard = Pokerset.NewPokerCardSet()
	pUser.pRobot = NewRobotSetting()

	pUser.Init()

	return pUser
}

// 初始化玩家資料
func (pUser *User) Init() {
	pUser.User = nil
	pUser.iUserType = gamelogic.UserTypeNormal
	pUser.i32UserID = -1
	pUser.i64Money = 0
	pUser.iUserState = UserState_EmptySeat
	pUser.iWeight = DefaultPlayerWeight

	pUser.ResetGameData()
}

// 重置遊戲資料
func (pUser *User) ResetGameData() {
	pUser.pHandCard.SetCardSet("")

	pUser.i64Timeout = 0
	pUser.i64LastMessageTime = 0
	pUser.i64MoneyChange = 0
	pUser.i64TotalWinLose = 0
	pUser.iPlayingState = UserPlayState_WaitingStart
	pUser.bOffline = false
	pUser.i64Bet_Base = 0
	pUser.i64Bet_Side = 0
	pUser.i64Bet_Double = 0
	pUser.i64FoldRecover = 0
	pUser.bIsTied = false
	pUser.bGetTieAward = false
	pUser.bCanEnterWar = false
	pUser.bIsFolded = false
	pUser.iRobotBetTimes = 0
	pUser.i32Point = 0
	pUser.bOperationDone = false
	pUser.i64TieAward = 0
	pUser.iSettlementWinLose = -1
	pUser.i64Tax = 0
	pUser.bGoForWar = false
}

// 設定機器人 timeout
func (pUser *User) SetRobotTimeout(i64Timeout int64) {
	pUser.i64Timeout = 0

	if i64Timeout > 0 {
		if pUser.iUserType == gamelogic.UserTypeRebot && pUser.iUserState == UserState_Normal {
			iRnd := util.RndInt(10, int(i64Timeout/2))
			pUser.i64Timeout = int64(iRnd)
		}
	}

	mlog.Debug("機器人[%d] i64Timeout[%d], timeout %d, PlayingState[%s]", pUser.i32UserID, i64Timeout, pUser.i64Timeout, g_map_strPlayingState_debug[pUser.iPlayingState])
}

// 押注
func (pUser *User) Bet(iBetType int, i64Bet int64) {
	mlog.Debug("User[%d] 押注, iBetType[%s], Bet[%d]", pUser.i32UserID, g_map_strBetType[iBetType], i64Bet)

	if iBetType < BetType_Base || iBetType >= BetType_Max {
		return
	}

	if i64Bet <= 0 || i64Bet > pUser.i64Money {
		return
	}

	switch iBetType {
	case BetType_Base: // 底注
		pUser.i64Bet_Base += i64Bet
	case BetType_Side: // 平手注
		pUser.i64Bet_Side += i64Bet
	}

	pUser.i64Money -= i64Bet
}
