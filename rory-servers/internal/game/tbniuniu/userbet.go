package tbnn

//用户投注
type userBet struct {
	userID   int32 //用户ID
	betTimes int32 //底注
}

//创建用户投注对象
func newUserBet(userid int32, times int32) *userBet {
	ubt := new(userBet)
	ubt.userID = userid
	ubt.betTimes = times

	return ubt
}

//游戏底注
type gameBet struct {
	userbets []*userBet
}

//创建游戏投注
func newGameBet() *gameBet {
	gbt := new(gameBet)

	return gbt
}

//添加用户底注
func (gbt *gameBet) addUserBet(userid int32, times int32) bool {
	//检查用户是否已经提交底注
	if gbt.checkUserBet(userid) {
		return false
	}

	//检查用户底注是否合法
	if !gbt.checkBetTimes(times) {
		return false
	}

	ubt := newUserBet(userid, times)
	gbt.userbets = append(gbt.userbets, ubt)

	return true
}

//检查用户是否投注
func (gbt *gameBet) checkUserBet(userid int32) bool {
	for _, ubt := range gbt.userbets {
		if ubt.userID == userid {
			return true
		}
	}

	return false
}

//确定用户投注,如果用户未投注,则设置最低底注
func (gbt *gameBet) ensureUserBet(userid int32) {
	if gbt.checkUserBet(userid) {
		return
	}

	gbt.addUserBet(userid, getNNConfig().minBaseScoreTimes())
}

//获取用户底注倍数
func (gbt *gameBet) getUserBetTimes(userid int32) int32 {
	for _, ubt := range gbt.userbets {
		if ubt.userID == userid {
			return ubt.betTimes
		}
	}

	return 0
}

//检查投注备注
func (gbt *gameBet) checkBetTimes(times int32) bool {
	for _, bt := range getNNConfig().baseScoreTimes {
		if bt == times {
			return true
		}
	}

	return false
}
