package hbsl

import "example.invalid/ServerTeam/Servers/pkg/gamelogic"

type User struct {
	pGameUser       *gamelogic.GameUser
	i64Money        int64               // 玩家金錢
	i64MoneyChange  int64               // 金錢變化
	i64TotalWinLose int64               // 總輸贏，client 端顯示用
	iWinRounds      int                 // 贏的場次
	pFaHB           *HB                 // 發紅包资讯
	pGrabHB         *HBRecord           // 搶紅包資訊
	pHBRecord       []*HBRecord         // 搶紅包紀錄
	pTempState      gamelogic.UserState // 暫存離場紀錄
	bTokenHB        bool                // 是否已拿到紅包
	i64HBMoney      int64               // 玩家搶得的紅包金額
}

type UserRecord struct {
	iUserID      int32 // 用户ID
	iFaceID      int32 // 頭像ID
	i64HBMoney   int64 // 紅包搶到金額
	i64Money     int64 // 玩家金錢
	i64LoseMoney int64 // 踩雷賠付金額
	pUser        *User // 玩家資訊
}

// 新增User數據
func NewUser(pGameUser *gamelogic.GameUser) *User {
	pUser := new(User)
	pUser.pGameUser = pGameUser
	pUser.i64Money = pGameUser.Score
	pUser.i64MoneyChange = 0
	pUser.i64TotalWinLose = 0
	pUser.iWinRounds = 0
	pUser.pFaHB = nil
	pUser.pGrabHB = nil
	pUser.pHBRecord = make([]*HBRecord, 0)
	pUser.pTempState = nil
	pUser.bTokenHB = false

	return pUser
}

// 清除User數據
func (pUser *User) Init(bClearFaHB bool) {
	pUser.i64MoneyChange = 0
	pUser.pGrabHB = nil
	pUser.bTokenHB = false
	pUser.i64HBMoney = 0
	if bClearFaHB {
		pUser.i64Money = pUser.pGameUser.Score
		pUser.pFaHB = nil
	}
}

// 清除UserRecord數據
func (pUserRecord *UserRecord) Init() {
	pUserRecord.pUser = nil
}
