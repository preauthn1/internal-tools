package wzp

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

//游戏轮次管理
type gameTurnManager struct {
	gamePoker       *poker.GamePoker    //玩家牌型
	wzplogic        *WZPGameLogic       //游戏逻辑实例
	currentTurnUser *gamelogic.GameUser //当前轮次用户
	firstUser       *gamelogic.GameUser //本轮下注用户
	turnNo          int32               //轮次
	addBet          bool                // 是否有人加注
	Suoha           bool                //是否梭哈

	giveupUsers []int32 //弃牌的用户
	newTurnFunc func(turnNo int32)
}

//创建游戏轮次管理对象
func newGameTurnManager(gamepoker *poker.GamePoker, wzpglc *WZPGameLogic) *gameTurnManager {
	gtm := new(gameTurnManager)
	gtm.gamePoker = gamepoker
	gtm.wzplogic = wzpglc
	gtm.giveupUsers = make([]int32, 0)

	return gtm
}

//游戏开始
func (gtm *gameTurnManager) gameBegin() {
	//轮次
	gtm.turnNo = 1
	//加注
	gtm.addBet = false
	//梭哈
	gtm.Suoha = false
	up, _ := getGameMaxCard(gtm.gamePoker, gtm.turnNo, gtm.wzplogic)

	gtm.currentTurnUser = up.User
	gtm.firstUser = up.User

}

//检查当前是否指定的用户操作
func (gtm *gameTurnManager) checkUserAction(user *gamelogic.GameUser) bool {
	if gtm.currentTurnUser.UserID == user.UserID {
		return true
	}

	return false
}

//收到用户动作
func (gtm *gameTurnManager) onUserAction(user *gamelogic.GameUser, giveup bool, actionid int32) {
	if !gtm.checkUserAction(user) {
		return
	}

	//如果是弃牌
	if giveup {
		mlog.Debug("用户%v弃牌", user)
		gtm.giveupUsers = append(gtm.giveupUsers, user.UserID)

	}
	//判断上轮用户是否加注
	if actionid == actionID_AddBet {
		gtm.addBet = true
	} else {
		gtm.addBet = false
	}

	//判断是否梭哈
	if actionid == actionID_Suoha {
		gtm.Suoha = true
		gtm.firstUser = user
	}

	// 2020.02.11 用戶加注後設為回合第一個人，使其他用戶再次進行是否要跟注或加注，做完回合才結束 by Alex
	if actionid == actionID_AddBet {
		gtm.firstUser = user
	}

	nextActSeatNo := gtm.currentTurnUser.SeatNo
	for {
		nextActSeatNo++

		if nextActSeatNo > 4 {
			nextActSeatNo = 5 - nextActSeatNo
		}

		if nextActSeatNo == gtm.firstUser.SeatNo {
			gtm.turnNo++
			if gtm.newTurnFunc != nil {
				gtm.newTurnFunc(gtm.turnNo)
			}

			//获取最大牌型
			up, _ := getGameMaxCard(gtm.gamePoker, gtm.turnNo, gtm.wzplogic)

			gtm.currentTurnUser = up.User
			gtm.firstUser = up.User

			mlog.Debug("进入下一轮.%d,下注用户%v", gtm.turnNo, gtm.firstUser)
			gtm.addBet = false

			return
		}

		ur := gtm.wzplogic.GetUserBySeatNo(nextActSeatNo)
		if ur == nil || (!ur.IsGameing() && !ur.IsOffLine()) || gtm.isGiveUp(ur.UserID) {
			continue
		}
		gtm.currentTurnUser = ur
		mlog.Debug("下一操作用户%v", gtm.currentTurnUser)

		//如果下注用户弃牌
		if giveup && user == gtm.firstUser {
			gtm.firstUser = gtm.currentTurnUser
		}
		return
	}
}

//是否弃牌
func (gtm *gameTurnManager) isGiveUp(userid int32) bool {
	for _, uid := range gtm.giveupUsers {
		if uid == userid {
			return true
		}
	}

	return false
}

//是否第一轮
func (gtm *gameTurnManager) isFirstTurn() bool {
	return gtm.currentTurnUser.UserID == gtm.firstUser.UserID
}

//是否最后一轮
func (gtm *gameTurnManager) isLastTurn() bool {
	nextActSeatNo := gtm.currentTurnUser.SeatNo
	for {
		nextActSeatNo++

		if nextActSeatNo >= 5 {
			nextActSeatNo = nextActSeatNo - 5
		}
		if nextActSeatNo == gtm.firstUser.SeatNo {
			return true
		}

		ur := gtm.wzplogic.GetUserBySeatNo(nextActSeatNo)

		if ur == nil || gtm.isGiveUp(ur.UserID) || !gtm.wzplogic.IsGameUser(ur) {
			continue
		} else {
			return false
		}
	}
}

//是否游戏结束
func (gtm *gameTurnManager) isGameEnd() bool {
	return gtm.turnNo >= 5 || len(gtm.giveupUsers) >= len(gtm.wzplogic.AllGameUser())-1
}

//是否其他用户都弃牌
func (gtm *gameTurnManager) isAllGiveUp() bool {
	return len(gtm.giveupUsers) >= len(gtm.wzplogic.AllGameUser())-1
}
