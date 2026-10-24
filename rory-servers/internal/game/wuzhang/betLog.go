package wzp

import (
	"math"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

const (
	TrunNo_BaseScore = 0 //底注轮次
)

//用户投注记录
type userBetLog struct {
	userID int32 //用户ID
	amount int64 //投注金额
	trun   int32 //轮次
}

//游戏投注记录
type gameBetLog struct {
	betlogs              []*userBetLog   //投注记录
	totalAmount          int64           //总投注金额
	userBets             map[int32]int64 //用户投注金额
	currentTurnNo        int32           //当前轮次
	currentTurnBetAmount int64           //本轮投注金额
	gamelogic            *WZPGameLogic   //游戏逻辑
}

//创建游戏投注记录对象
func newGameBetLog(glc *WZPGameLogic) *gameBetLog {
	logs := new(gameBetLog)
	logs.totalAmount = 0
	logs.userBets = make(map[int32]int64)
	logs.gamelogic = glc

	return logs
}

//添加投注记录
func (gl *gameBetLog) addBetLog(ur *gamelogic.GameUser, amount int64, trun int32) {
	mlog.Debug("添加投注记录")

	// 2020.02.12 跟注金額可能為非籌碼上數字，這邊先注解
	// //检查投注金额
	// if trun > 0 && !gl.checkTurnBetAmount(amount, trun, ur) {
	// 	mlog.Debug("检查本轮投注错误,%d,%d", amount, trun)
	// 	return
	// }

	//检查底注
	if trun == 0 && amount != int64(gameconfig.baseAmount) {
		mlog.Debug("检查本轮投注错误,%d,%d", amount, trun)
		return
	}

	userid := ur.UserID

	l := new(userBetLog)
	l.userID = userid
	l.amount = amount
	l.trun = trun

	gl.totalAmount = gl.totalAmount + amount

	if cn, ok := gl.userBets[userid]; ok {
		gl.userBets[userid] = cn + amount
	} else {
		gl.userBets[userid] = amount
	}

	gl.betlogs = append(gl.betlogs, l)

	//设置当前投注金额
	if gl.currentTurnNo != trun {
		gl.currentTurnBetAmount = amount
	} else if gl.currentTurnNo == trun && amount > gl.currentTurnBetAmount {
		gl.currentTurnBetAmount = amount
	}

	//设置当前轮次
	gl.currentTurnNo = trun

	mlog.Debug("本轮最低投注金额%d", gl.currentTurnBetAmount)

}

//添加底注投注记录
func (gl *gameBetLog) addBaseScoreBetLog(ur *gamelogic.GameUser) {
	gl.addBetLog(ur, int64(gameconfig.baseAmount), TrunNo_BaseScore)
}

//检查投注金额
func (gl *gameBetLog) checkBetAmount(amount int64) bool {
	for _, amt := range gameconfig.betAmount {
		if int64(amt) == amount {
			return true
		}
	}

	return false
}

//检查本轮投注金额
func (gl *gameBetLog) checkTurnBetAmount(amount int64, turn int32, ur *gamelogic.GameUser) bool {
	if turn != 4 && !gl.checkBetAmount(amount) {
		return false
	}
	hasbetamount, ok := gl.userBets[ur.UserID]

	urbetamount := amount
	if ok {
		urbetamount += hasbetamount
	}

	if int64(urbetamount) > ur.Score {
		return false
	}

	//新一轮
	if turn > gl.currentTurnNo {
		return true
	}

	//本轮
	if amount < gl.currentTurnBetAmount {
		return false
	}

	return true

}

//获取梭哈金额
func (gl *gameBetLog) getSuohaAmount() (int64, *gamelogic.GameUser) {
	//获取用户身上最小金币数
	var minamount int64
	minamount = math.MaxInt64
	var minuser *gamelogic.GameUser

	for _, ur := range gl.gamelogic.AllGameUser() {
		//检查用户是否弃牌
		if gl.gamelogic.turnManager.isGiveUp(ur.UserID) {
			continue
		}

		score := ur.Score
		if hasBetAmount, ok := gl.userBets[ur.UserID]; ok {
			score -= int64(hasBetAmount)
		}

		if score <= int64(minamount) {
			minamount = int64(score)
			minuser = ur
		}
	}

	return minamount, minuser

}
