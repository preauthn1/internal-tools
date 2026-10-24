package ddz

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
)

//AI強化學習智能
func (lgc *DDZRobotLogic) robotAIAction(robot *robotcomm.Robot) {
	//lgc.updateRoleInfo(robot)
	roleInfo := lgc.userRole[robot.UserInfo.UserID]
	mlog.Debug("AI 操作 robot[%v] seatmine[%v] ", robot.UserInfo.UserID, roleInfo.SeatMine)
	mlog.Debug("AI Action userid[%v] seatno[%v] 本身地主[%v] 上家地主[%v] 下家地主[%v] 本身座位[%v] 上家座位[%v] 下家座位[%v] ",
		robot.UserInfo.UserID, robot.UserInfo.SeatNo,
		roleInfo.bMine, roleInfo.bLast, roleInfo.bNext, roleInfo.SeatMine, roleInfo.SeatLast, roleInfo.SeatNext)

	if roleInfo == nil {
		//lgc.robotPass(robot)
		lgc.robotAction(robot)
		return
	}

	//跟牌, 特殊處理
	more := int32(0)
	if lgc.outInfo.RealSeatNo != -1 && lgc.outInfo.RealSeatNo != robot.UserInfo.SeatNo {
		more = lgc.outInfo.nMax
		if lgc.outInfo.nMax == poker.PokerNumberTwo {
			more = poker.PokerNumberA
		}

		if lgc.outInfo.CardType == int32(cardTypeRocket) {
			lgc.getGameLogic().GameLog.AddGameStep("robotPass1 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "robotAIAction 跟牌，場面牌是王炸")
			lgc.robotPass(robot)
			return
		}

		//場上有人剩下兩手牌
		dc := NewDefineCard(lgc.GameLogic.(*DDZGameLogic))
		up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
		ctg := dc.Combination(up.Cards)
		passTwoHand := false
		passTwoHand = lgc.passTwoHands(ctg, dc, robot, passTwoHand, roleInfo)
		if checkAIWithFlag(robot, RobotAIFlag3) && passTwoHand {
			//當地主手上沒單張, 只剩炸彈, 不能跟單張,這是陷阱.
			lastUr := lgc.getGameLogic().GetUserBySeatNo(lgc.outInfo.RealSeatNo)
			lastUrPoker := lgc.getGameLogic().gamePoker.GetUserPoker(lastUr.UserID)
			lastCtg := dc.Combination(lastUrPoker.Cards)
			if lgc.outInfo.CardType == int32(cardTypeSingle) && len(lastCtg.SingleCards) != 0 {
				lgc.FollowCards(robot, cardTypeSingle)
				return
			}
			lgc.getGameLogic().GameLog.AddGameStep("robotPass2 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "robotAIAction checkAIWithFlag(robot, RobotAIFlag3) 跟牌，場上有人剩下兩手牌")
			lgc.robotPass(robot)
			return
		}
	}

	//主動, 特殊處理
	if lgc.outInfo.RealSeatNo == robot.UserInfo.SeatNo {
		ok := lgc.outMinCard(roleInfo, robot)
		if ok {
			return
		}
	}

	ok := false

	//TODO 上家是地主並且是真人.
	if checkAIWithFlag(robot, RobotAIFlag5) {
		ok = lgc.aiLastDizhuAndRealUser(roleInfo, robot, more)
		if ok {
			lgc.getGameLogic().GameLog.AddGameStep("robotPass3 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "robotAIAction checkAIWithFlag(robot, RobotAIFlag5) 跟牌，上家是地主並且是真人")
			lgc.robotPass(robot)
			return
		}
	}

	//TODO 下家是隊友農民, 地主是真人, 當隊友快贏要做球給隊友
	ok = lgc.aiNextPartner(roleInfo, robot, more)
	if ok {
		return
	}

	//AI出牌邏輯不管主動或跟牌都使用以下情境.

	_, dizhu := lgc.getGameLogic().ext.getDiZhu()
	dizhuUr := lgc.getGameLogic().GetUserByUserID(dizhu.userId)
	//地主是真人, 且沒有組合沒不要先出單張
	if !dizhuUr.IsRobot() && len(lgc.sandBoxUser[dizhu.userId].ctg.SingleCards) > 0 && noGroupCard(lgc.sandBoxUser[dizhu.userId].ctg.Ct) {
		mlog.Warn("robot[%v] 地主是真人, 且沒有組合沒不要先出單張", robot.UserInfo.UserID)
		ok = false
	} else {
		ok = lgc.aiOutDanpai(robot, ok, more) //主動或跟牌都由單牌開始
	}

	if ok {
		return
	}

	//從多張組合開始出牌
	ok = lgc.aiOutCombine(robot, more)
	if ok {
		return
	}

	ok = lgc.aiOutDanpai(robot, ok, more)
	if ok {
		return
	} else {
		mlog.Warn("AI 單牌輸, 切一般robot模式 userid[%v] ", robot.UserInfo.UserID)
		lgc.robotAction(robot)
	}
	return
}

//上家是地主並且是真人
func (lgc *DDZRobotLogic) aiLastDizhuAndRealUser(roleInfo *robotRole, robot *robotcomm.Robot, more int32) bool {
	ok := false

	//主動出牌
	if lgc.outInfo.RealSeatNo == robot.UserInfo.SeatNo {
		return ok
	}

	_, dizhu := lgc.getGameLogic().ext.getDiZhu()
	if dizhu == nil {
		return false
	}
	user := lgc.getGameLogic().GetUserByUserID(dizhu.userId)
	if user == nil {
		return false
	}

	//mlog.Debug("robot[%v] 下家是農民 下家是否地主[%v] 地主是真人[%v] ", robot.UserInfo.UserID, !roleInfo.bNext, !user.IsRobot())

	if !roleInfo.bNext && !user.IsRobot() {

		if lgc.outInfo.RealSeatNo == roleInfo.SeatNext {
			mlog.Debug("最後出牌是下家農民就pass, 讓他先出牌, SeatNext[%v]", roleInfo.SeatNext)
			return true
		}

		//mlog.Debug("AI robot[%v] 下家是農民,上家是地主並且是真人, SeatMine[%v]  ", robot.UserInfo.UserID, roleInfo.SeatMine)

		lgc.addSandBoxAllUser()
		nextUser := lgc.getGameLogic().GetUserBySeatNo(roleInfo.SeatNext)
		aiBoxNUser := lgc.sandBoxUser[nextUser.UserID]

		lastCardType := cardType(lgc.outInfo.CardType) //最後出牌類型
		lastOutMax := lgc.outInfo.nMax
		if lastOutMax == poker.PokerNumberTwo {
			lastOutMax = poker.PokerNumberA
		}
		//檢核下家農民是否有大於地主的單牌, 有的話就pass讓下家農民出牌.
		if lastCardType == cardTypeSingle {
			for _, v := range aiBoxNUser.ctg.SingleCards {
				if v > lastOutMax {
					mlog.Debug("單牌 robot[%v] 下家農民大於地主牌, 下家農民牌[%v], 地主牌[%v]", robot.UserInfo.UserID, v, lastOutMax)
					ok = true
					break
				}
			}
		}

		//檢核組合牌
		for _, v := range aiBoxNUser.ctg.Ct {
			if lastCardType == v.byteCardType && v.nMax > lastOutMax && lastOutMax != poker.PokerNumberTwo ||
				(lastCardType > cardTypeBomb && v.byteCardType <= cardTypeBomb) {
				mlog.Debug("組牌 robot[%v] 下家農民大於地主牌, cardtype[%v], 下家農民牌型[%v] 下家農民牌[%v], 地主牌[%v]",
					robot.UserInfo.UserID, lastCardType, v.byteCardType, v.nMax, lastOutMax)
				ok = true
				break
			}
		}
	}

	return ok
}

//自己是農民, 下家也是農民, 當下家只剩一組牌, 做球給下家贏
func (lgc *DDZRobotLogic) outMinCard(roleInfo *robotRole, robot *robotcomm.Robot) bool {
	lgc.outInfo.CardType = int32(cardTypeNone)
	//自己是農民, 下家也是農民, 當下家只剩一組牌, 做球給下家贏
	if checkAIWithFlag(robot, RobotAIFlag1) && roleInfo.bLast == true { //上家是地主
		lgc.addSandBoxAllUser()
		nextUser := lgc.getGameLogic().GetUserBySeatNo(roleInfo.SeatNext)
		aiBoxUser := lgc.sandBoxUser[nextUser.UserID]
		if aiBoxUser.ctg.hands == 1 {
			mlog.Debug("robot[%v] 做球給下家robot[%v]贏 ", robot.UserInfo.UserID, nextUser.UserID)
			if len(aiBoxUser.ctg.SingleCards) > 0 {
				outCard := lgc.aiOutCardByType(cardTypeSingle, robot.UserInfo.UserID, 0) //產生出牌

				if len(outCard) > 0 {
					mlog.Debug("userid[%v] 做球給下家贏, 出牌為[%v]", robot.UserInfo.UserID, outCard)
					up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
					cg := lgc.list2cardGroup(up.Cards, outCard)
					lgc.S2COutCard(cg, robot.UserInfo)
					return true
				}
			}

			nextCt := cardTypeNone
			for _, v := range lgc.sandBoxUser[nextUser.UserID].ctg.Ct {
				if v.byteCardType != cardTypeNone {
					nextCt = v.byteCardType
					break
				}
			}
			outCard := lgc.aiOutCardByType(nextCt, robot.UserInfo.UserID, 0) //產生出牌

			if len(outCard) > 0 {
				mlog.Debug("userid[%v] 做球給下家贏, 出牌為[%v]", robot.UserInfo.UserID, outCard)
				up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
				cg := lgc.list2cardGroup(up.Cards, outCard)
				lgc.S2COutCard(cg, robot.UserInfo)
				return true
			}
		}
	}
	return false
}

//下家是隊友農民, 地主是真人, 當隊友快贏要做球給隊友
func (lgc *DDZRobotLogic) aiNextPartner(roleInfo *robotRole, robot *robotcomm.Robot, more int32) bool {
	ok := false

	_, dizhu := lgc.getGameLogic().ext.getDiZhu()
	if dizhu == nil {
		return false
	}
	dizhuUr := lgc.getGameLogic().GetUserByUserID(dizhu.userId)
	if dizhuUr == nil || dizhuUr.IsRobot() {
		return false
	}

	nextUser := lgc.getGameLogic().GetUserBySeatNo(roleInfo.SeatNext)
	aiBoxNUser := lgc.sandBoxUser[nextUser.UserID]

	if !roleInfo.bNext && aiBoxNUser.ctg.hands > 5 {
		mlog.Debug("robot[%v] 下家農民超過5手牌先不做球.", robot.UserInfo.UserID)
		return false
	}

	//主動出牌
	if lgc.outInfo.RealSeatNo == robot.UserInfo.SeatNo {
		mlog.Debug("robot[%v] 主動出牌,下家是隊友農民, 地主是真人, 當隊友快贏要做球給隊友", robot.UserInfo.UserID)
		if !roleInfo.bNext && !dizhuUr.IsRobot() {
			mlog.Debug("robot[%v] B下家是隊友農民[%v], 地主是真人, 當隊友快贏要做球給隊友", robot.UserInfo.UserID, nextUser.UserID)
			lgc.addSandBoxAllUser()
			if lgc.outInfo.RealSeatNo == roleInfo.SeatNext {
				mlog.Debug("最後出牌是下家農民就pass, 讓他先出牌, SeatNext[%v]", roleInfo.SeatNext)
				return true
			}

			//下家農民有炸彈先不做球
			for _, v := range aiBoxNUser.ctg.Ct {
				if v.byteCardType == cardTypeBomb || v.byteCardType == cardTypeRocket {
					mlog.Debug("robot[%v] 下家農民有炸彈先不做球", robot.UserInfo.UserID)
					return false
				}
			}

			//下家農民剩下四手牌以內才做球, 從組合牌小牌開始出
			if aiBoxNUser.ctg.hands <= 4 {
				mlog.Debug("robot[%v] 下家農民[%v]剩下四手牌以內做球 ", robot.UserInfo.UserID, nextUser.UserID)

				for _, v := range aiBoxNUser.ctg.Ct {
					mlog.Debug("robot[%v] 產生做球牌型[%v] 點數[%v] 給下家農民", robot.UserInfo.UserID, v.byteCardType, v.nMax)
					outCard := lgc.aiOutCardByType(v.byteCardType, robot.UserInfo.UserID, 0) //產生出牌

					if len(outCard) > 0 {
						mlog.Warn("userid[%v] 做球給下家贏, 出牌為[%v]", robot.UserInfo.UserID, outCard)
						up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
						cg := lgc.list2cardGroup(up.Cards, outCard)
						lgc.S2COutCard(cg, robot.UserInfo)
						return true
					}
				}

				outCard := lgc.aiOutCardByType(cardTypeSingle, robot.UserInfo.UserID, 2) //產生出牌
				mlog.Debug("robot[%v] 產生做單牌[%v]", robot.UserInfo.UserID, outCard)
				for _, v := range aiBoxNUser.ctg.SingleCards {
					if len(outCard) > 0 && v > outCard[0] {
						mlog.Warn("userid[%v] 做球單牌 給下家贏, 出牌為[%v]", robot.UserInfo.UserID, outCard)
						up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
						cg := lgc.list2cardGroup(up.Cards, outCard)
						lgc.S2COutCard(cg, robot.UserInfo)
						return true
					}
				}
			}
		}

		//當真人地主只剩單張要主動出組合牌
		dizhuHand := lgc.sandBoxUser[dizhuUr.UserID].ctg.hands
		if dizhuHand == 1 && !dizhuUr.IsRobot() {
			dizhuCardType := cardTypeSingle
			for _, v := range lgc.sandBoxUser[dizhuUr.UserID].ctg.Ct {
				if v.byteCardType != cardTypeNone {
					dizhuCardType = v.byteCardType
					break
				}
			}

			mlog.Debug("robot[%v] 當真人地主只剩單張要主動出組合牌", robot.UserInfo.UserID)
			for iCardType := 13; iCardType > 0; iCardType-- {
				if cardType(iCardType) == dizhuCardType {
					mlog.Debug("robot[%v] 不出相同牌型給地主, 牌型[%v]", robot.UserInfo.UserID, iCardType)
					continue
				}
				outCard := lgc.aiOutCardByType(cardType(iCardType), robot.UserInfo.UserID, 0) //產生出牌
				if len(outCard) > 0 {
					mlog.Warn("userid[%v]壓牌 不出相同牌型給地主, 出牌為[%v]", robot.UserInfo.UserID, outCard)
					up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
					cg := lgc.list2cardGroup(up.Cards, outCard)
					lgc.S2COutCard(cg, robot.UserInfo)
					return true
				}
			}
		}
	}

	//跟牌
	if lgc.outInfo.RealSeatNo != -1 && lgc.outInfo.RealSeatNo != robot.UserInfo.SeatNo {
		//當下家農民只剩一手牌,  此農民要壓牌拿到出牌權再做球給下家農民.
		if !roleInfo.bNext && aiBoxNUser.ctg.hands == 1 {
			mlog.Debug("robot[%v] 跟牌 當下家農民只剩一手牌, 此農民要壓牌拿到出牌權再做球給下家農民", robot.UserInfo.UserID)
			outCard := lgc.aiOutCardByType(cardTypeBomb, robot.UserInfo.UserID, lgc.outInfo.nMax) //產生出牌
			if lgc.outInfo.nMax == poker.PokerNumberTwo {
				outCard = make([]int32, 0)
			}

			if len(outCard) == 0 {
				cdlist := lgc.outRocket(robot.UserInfo.UserID)
				if len(cdlist) > 0 {
					mlog.Debug("robot[%v] 用火箭壓取得出牌權", robot.UserInfo.UserID)
					up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
					cg := lgc.list2cardGroup(up.Cards, cdlist)
					lgc.S2COutCard(cg, robot.UserInfo)
					return true
				}
			}

			if len(outCard) > 0 {
				mlog.Debug("userid[%v]壓牌 準備做球給下家贏, 出牌為[%v]", robot.UserInfo.UserID, outCard)
				up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
				cg := lgc.list2cardGroup(up.Cards, outCard)
				lgc.S2COutCard(cg, robot.UserInfo)
				return true
			}
		}
	}

	return ok
}

//AI情境出組合牌, 從多張組合開始出牌
func (lgc *DDZRobotLogic) aiOutCombine(robot *robotcomm.Robot, more int32) bool {
	ok := false
	boxTurn := lgc.getGameLogic().ext.turn
	for t := 0; t <= 3; t++ {
		if ok {
			break
		}
		firstOut := false //要主動出牌
		if lgc.outInfo.RealSeatNo == robot.UserInfo.SeatNo || lgc.outInfo.RealSeatNo == -1 {
			firstOut = true
		}

		action1outCard := make([]int32, 0) //產生出牌
		//example.invalid("AI重新組合牌 回合[%v]", t)
		lgc.addSandBoxAllUser()
		lastActionTurn := boxTurn //最後出牌的位置
		retry := 3
		lastCardType := cardType(lgc.outInfo.CardType) //最後出牌類型

		for i := int(13) - t; i >= int(0); i-- {
			if ok {
				break
			}
			if i <= 0 && retry > 0 {
				i = 14
				retry--
				boxTurn = lastActionTurn
				continue
			}

			if i != int(lastCardType) && i > int(cardTypeBomb) && cardType(i) == cardType(cardTypeNone) {
				mlog.Debug("AI robot[%v] 不同牌型不能壓, 牌型i[%v] ", robot.UserInfo.UserID, i)
				continue
			}

			if lgc.outInfo.CardType == int32(cardTypeSingle) && i > int(cardTypeSingle) {
				//mlog.Debug("AI robot[%v] 不同牌型不能壓, 牌型i[%v] ", robot.UserInfo.UserID, i)
				continue
			}

			runLimit := 0
			cdlist := make([]int32, 0) //出牌
			danpai := more             //int32(2)         //輪詢從不同單牌點數開始出牌

			passTimes := 0 //pass次數, 當兩個pass, user1就從最小牌開始

			for {
				if ok {
					break
				}
				if lastCardType >= cardTypeBomb && lastCardType != cardType(i) {
					//mlog.Debug("上一手出牌型[%v] 目前牌型[%v] , break", lastCardType, i)
					break
				}
				//mlog.Debug("上一手出牌型[%v] 目前牌型[%v] 要大於[%v]點, danpai[%v] , break", lastCardType, i, more, danpai)

				//組合牌AI
				actionUser1 := lgc.getGameLogic().ext.getExtInfoByTurn(boxTurn)
				outCard := lgc.aiOutCardByType(cardType(i), actionUser1.userId, danpai) //產生出牌

				if len(outCard) == 0 {
					if firstOut {
						//mlog.Debug("seatno[%v] userid[%v] 主動出牌是空的", lgc.outInfo.RealSeatNo, actionUser1.userId)
						break
					}
					mlog.Debug("AI 組牌 pass userid[%v]", actionUser1.userId)
					passTimes++
				}
				if passTimes == int(poker.PokerNumberTwo) {
					boxTurn = lastActionTurn
					firstOut = true
					example.invalid("AI 2次pass, 最後出牌turn[%v] [%v]", lastActionTurn, boxTurn)
					break
				}
				if len(outCard) > 0 {
					lastCardType = cardType(i)
					firstOut = false
					lastActionTurn = boxTurn

					if len(action1outCard) == 0 {
						action1outCard = outCard
						example.invalid("AI 起始出牌 userid[%v] [%v]", actionUser1.userId, action1outCard)
					}

					mlog.Debug(" AI黑箱 出牌組合 userid[%v] turn[%v] cardType[%v] 出牌[%v] moreDanpai[%v] 剩餘手牌[%v] - [%v]",
						actionUser1.userId, boxTurn, i, outCard, danpai, lgc.sandBoxUser[actionUser1.userId].ctg.Ct, lgc.sandBoxUser[actionUser1.userId].ctg.SingleCards)
					if cardType(i) == cardTypeThreeStraight || cardType(i) == cardTypeDoubleStraightd || cardType(i) == cardTypeDoubleStraights ||
						cardType(i) == cardTypeStraight || cardType(i) == cardTypeBomb ||
						cardType(i) == cardTypeThree || cardType(i) == cardTypeThreeAndOne {
						danpai = outCard[0]
					} else if cardType(i) == cardTypeDouble {
						if outCard[0] == poker.PokerNumberTwo {
							boxTurn++
							mlog.Debug("AI 2對子不能再出牌, +boxTurn[%v]", boxTurn)
							break
						} else {
							danpai = outCard[0]
						}
					} else {
						danpai = outCard[len(outCard)-1]
					}
				}

				//機器人只剩單牌,  先出完獲勝, 用cardTypesCount判斷
				if len(action1outCard) > 0 && noGroupCard(lgc.sandBoxUser[robot.UserInfo.UserID].ctg.Ct) &&
					cardType(i) == cardTypeSingle && action1outCard[0] != lgc.outInfo.nMax {
					mlog.Debug("userid[%v] AI組牌獲勝, 出牌[%v] more[%v] outMax[%v]", robot.UserInfo.UserID, action1outCard, danpai, lgc.outInfo.nMax)
					ok = true
					cdlist = action1outCard
					break
				}

				boxTurn++
				actionUser2 := lgc.getGameLogic().ext.getExtInfoByTurn(boxTurn)
				actionUser3 := lgc.getGameLogic().ext.getExtInfoByTurn(boxTurn + 1)

				if noGroupCard(lgc.sandBoxUser[actionUser2.userId].ctg.Ct) ||
					noGroupCard(lgc.sandBoxUser[actionUser3.userId].ctg.Ct) {
					mlog.Debug("AI 組合牌 地主[%v]機器人輸, userid[%v]單牌張數[%v], userid[%v]單牌張數[%v]",
						actionUser1.userId,
						actionUser2.userId, len(lgc.sandBoxUser[actionUser2.userId].ctg.Ct),
						actionUser3.userId, len(lgc.sandBoxUser[actionUser3.userId].ctg.Ct))
					break
				}

				if runLimit > 5 {
					break
				}
				runLimit++
			}

			if ok {
				if len(cdlist) > 0 {
					mlog.Debug("userid[%v] 組牌機器人獲勝, 出牌為[%v]", robot.UserInfo.UserID, cdlist)
					up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
					cg := lgc.list2cardGroup(up.Cards, cdlist)
					lgc.S2COutCard(cg, robot.UserInfo)
					ok = true
					break
				}
			}
		}
	}
	return ok
}

//AI情境出單牌
func (lgc *DDZRobotLogic) aiOutDanpai(robot *robotcomm.Robot, isWin bool, more int32) bool {
	if len(lgc.sandBoxUser[robot.UserInfo.UserID].ctg.SingleCards) == 0 ||
		lgc.outInfo.CardType != int32(cardTypeSingle) {
		return false
	}

	for iCard := more; iCard <= poker.PokerNumberBigKing; iCard++ {
		lgc.addSandBoxAllUser()
		runLimit := 0
		boxTurn := lgc.getGameLogic().ext.turn
		danpai := int32(iCard)     //輪詢從不同單牌點數開始出牌
		passTimes := 0             //pass次數, 當兩個pass, user1就從最小牌開始
		cdlist := make([]int32, 0) //出牌

		//單牌sandbox
		for {
			if more == poker.PokerNumberA && danpai > poker.PokerNumberTwo && danpai < poker.PokerNumberA {
				break
			}
			//單牌AI
			actionUser1 := lgc.getGameLogic().ext.getExtInfoByTurn(boxTurn)
			outCard := lgc.aiOutCardByType(cardType(cardTypeSingle), actionUser1.userId, danpai) //產生出牌

			if len(outCard) == 0 {
				passTimes++
			}
			if len(outCard) == 1 {
				if outCard[0] == poker.PokerNumberTwo {
					danpai = poker.PokerNumberA //出牌為2, 下一個單牌要超過14
				} else {
					danpai = outCard[0]
				}
				passTimes = 0
			}
			if passTimes == int(poker.PokerNumberTwo) {
				break
			}

			//TODO 檢核是否有人出完單牌,  先出完獲勝, 用cardTypesCount判斷
			if len(lgc.sandBoxUser[robot.UserInfo.UserID].ctg.SingleCards) == 0 {
				isWin = true
				cdlist = outCard
				break
			}

			boxTurn++
			actionUser2 := lgc.getGameLogic().ext.getExtInfoByTurn(boxTurn)
			actionUser3 := lgc.getGameLogic().ext.getExtInfoByTurn(boxTurn + 1)

			if len(lgc.sandBoxUser[actionUser2.userId].ctg.SingleCards) == 0 ||
				len(lgc.sandBoxUser[actionUser3.userId].ctg.SingleCards) == 0 {
				mlog.Debug("AI黑箱 單牌 地主[%v]機器人輸, iCard[%v], userid[%v]單牌張數[%v], userid[%v]單牌張數[%v]",
					actionUser1.userId, iCard,
					actionUser2.userId, len(lgc.sandBoxUser[actionUser2.userId].ctg.SingleCards),
					actionUser3.userId, len(lgc.sandBoxUser[actionUser3.userId].ctg.SingleCards))
				break
			}

			if runLimit > 10 {
				break
			}
			runLimit++
		}

		//組合牌sandbox
		if isWin {
			if len(cdlist) > 0 {
				up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
				cg := lgc.list2cardGroup(up.Cards, cdlist)
				lgc.S2COutCard(cg, robot.UserInfo)
				mlog.Debug("userid[%v] 獲勝單牌, 出牌[%v]", robot.UserInfo.UserID, cdlist)
				break
			}
		}
	}
	return isWin
}

//沒有組合牌
func noGroupCard(c []*cardTypeGroup) bool {
	result := true
	for _, v := range c {
		if v.byteCardType != cardTypeNone && v.byteCardType != cardTypeSingle {
			result = false
			break
		}
	}
	return result
}

//將所有玩家手牌放入黑箱
func (lgc *DDZRobotLogic) addSandBoxAllUser() {
	for _, user := range lgc.getGameLogic().AllTableUser() {
		upk := lgc.getGameLogic().gamePoker.GetUserPoker(user.UserID)
		lgc.aiSandBoxRecombination(user.UserID, upk.Cards)
	}
}

//AI 黑箱重新拆牌
func (lgc *DDZRobotLogic) aiSandBoxRecombination(userid int32, cg *poker.CardGroup) {
	//重新拆牌
	dc := NewDefineCard(lgc.getGameLogic())
	result := dc.Combination(cg)
	ctg := new(robotUserInfoLocal)
	ctg.cardTypesCount = make(map[cardType]int32, 0)
	ctg.ctg = result

	for _, data := range result.Ct {
		ctg.cardTypesCount[data.byteCardType]++
	}
	ctg.cardTypesCount[cardTypeSingle] += int32(len(result.SingleCards))
	lgc.sandBoxUser[userid] = ctg
}

//AI根据牌型出牌
func (lgc *DDZRobotLogic) aiOutCardByType(outCardType cardType, userid int32, more int32) (a []int32) {
	defer func() {
		if len(a) > 0 {
			mlog.Debug("userid[%v] 根据牌型出牌，选定的牌型为[%v],牌%v點 , 要大於[%v]點", userid, int(outCardType), a, more)
		}
	}()

	var result []int32

	if len(result) != int(0) {
		return result
	}

	//三带
	if outCardType == cardTypeThree { //三带
		return lgc.aiboxOutThreeCards(userid, more)
	}
	//飞机
	if outCardType == cardTypePlaneAndWing {
		return lgc.aiBoxOutFeiji(userid, 5, more)
	}
	//连对
	if outCardType == cardTypeThreeStraight {
		return lgc.aiBoxOutLiandui(userid, more)
	}
	//顺子
	if outCardType == cardTypeStraight {
		return lgc.aiBoxOutShunzi(userid, more)
	}
	//对子
	if outCardType == cardTypeDouble {
		return lgc.aiBoxOutDuizi(userid, more)
	}
	//单牌
	if outCardType == cardTypeSingle {
		return lgc.aiBoxOutDanpai(userid, more)
	}
	//炸弹
	if outCardType == cardTypeBomb {
		return lgc.aiBoxOutZhadan(userid, more)
	}
	//火箭
	if outCardType == cardTypeRocket {
		return lgc.outRocket(userid)
	}

	return nil
}

//AI黑箱 移除指定單牌
func (lgc *DDZRobotLogic) aibRemoveDanpai(userid, danpdi int32) {
	newDanpdi := make([]int32, 0)
	for _, v := range lgc.sandBoxUser[userid].ctg.SingleCards {
		if v == danpdi {
			continue
		}
		newDanpdi = append(newDanpdi, v)
	}
	lgc.sandBoxUser[userid].ctg.SingleCards = newDanpdi
}

//单牌, 获取不小于more的最小牌型
func (lgc *DDZRobotLogic) aiBoxOutDanpai(userid int32, more int32) []int32 {
	if more == poker.PokerNumberTwo {
		more = poker.PokerNumberA
	}
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌中-單牌 大於點數[%v]", userid, more)).Debug("AI")
	outCardType := cardTypeSingle
	outCards := make([]int32, 0)
	defer func() {
		if len(outCards) > 0 {
			mlog.Debug("userid[%v]出牌单牌牌型%v", userid, outCards)
		}
	}()

	if more == poker.PokerNumberBigKing {
		return outCards
	}

	var outCardTypeGroup *cardTypeGroup

	aiUser := lgc.sandBoxUser[userid]

	outCardTypeGroup = lgc.getGameLogic().getMinCardSingleType(outCardType, more, aiUser.ctg.SingleCards)
	//mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌中-單牌 單排最小張[%v]", userid, outCardTypeGroup)).Debug("AI")
	if outCardTypeGroup == nil { //手牌點數沒有比牌桌大時, 出手牌的2(如果有)
		outCardTypeGroup = lgc.getMinCardType(outCardType, 0, aiUser.ctg)

		if outCardTypeGroup != nil && outCardTypeGroup.nMax == int32(2) && more != poker.PokerNumberSmallKing {
			outCards = append(outCards, outCardTypeGroup.nMax)
			lgc.eraseResult(aiUser.ctg, outCardTypeGroup)
			lgc.reduceCardType(aiUser, cardTypeSingle, 1)
			lgc.aibRemoveDanpai(userid, outCardTypeGroup.nMax)
			mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌中-單牌出牌B [%v]", userid, outCards)).Debug("AI")
			return outCards
		}

	} else if outCardTypeGroup != nil {
		outCards = append(outCards, outCardTypeGroup.nMax)
		lgc.eraseResult(aiUser.ctg, outCardTypeGroup)
		lgc.reduceCardType(aiUser, cardTypeSingle, 1)
		lgc.aibRemoveDanpai(userid, outCardTypeGroup.nMax)
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌中-單牌出牌 [%v]", userid, outCards)).Debug("託管")
		return outCards
	}
	//mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌中-單牌出牌 [%v]", userid, outCards)).Debug("託管")
	return outCards
}

//对子
func (lgc *DDZRobotLogic) aiBoxOutDuizi(userid int32, more int32) []int32 {
	outCardType := cardTypeDouble
	outCards := make([]int32, 0)
	defer func() {
		mlog.Debug("userid[%v] 获取对子牌型%v", userid, outCards)
	}()

	tuoGuanUser := lgc.sandBoxUser[userid]
	if tuoGuanUser == nil {
		mlog.Debug("获取机器人数据失败，对子")
		return outCards
	}

	outCardTypeGroup := new(cardTypeGroup)

	var handDoubleCount int32 //手牌對子的數量
	for _, v := range tuoGuanUser.ctg.Ct {
		if v.byteCardType == cardTypeDouble {
			handDoubleCount++
		}
	}

	outCardTypeGroup = lgc.getMinCardType(outCardType, more, tuoGuanUser.ctg)
	//mlog.NewDefaultLogger().SetData(fmt.Sprintf("組合自動出牌中-對子 userid[%v] outCardTypeGroup[%v] allCardCtg[%v]", userid, &outCardTypeGroup, tuoGuanUser.ctg.Ct)).Debug("AI組牌")
	if outCardTypeGroup == nil {
		outCardTypeGroup = lgc.getMinCardType(outCardType, 0, tuoGuanUser.ctg)
		if outCardTypeGroup == nil || (outCardTypeGroup.nMax > int32(2) && outCardTypeGroup.nMax <= more) {
			return outCards
		}
	}

	if outCardTypeGroup.nMax > more {
		outCards = append(outCards, outCardTypeGroup.nMax)
		outCards = append(outCards, outCardTypeGroup.nMax)
		//mlog.NewDefaultLogger().SetData(fmt.Sprintf("組合自動出牌中-對子 userid[%v] outCards[%v]", userid, outCards)).Debug("AI組牌")
		lgc.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
		lgc.reduceCardType(tuoGuanUser, cardTypeDouble, 1)
	}

	return outCards
}

//顺子
func (lgc *DDZRobotLogic) aiBoxOutShunzi(userid int32, more int32) []int32 {
	outCardType := cardTypeStraight
	outCards := make([]int32, 0)
	defer func() {
		if len(outCards) > 0 {
			mlog.Debug("获取顺子牌型%v", outCards)
		}
	}()
	tuoGuanUser := lgc.sandBoxUser[userid]
	if tuoGuanUser == nil {
		mlog.Warn("机器人数据获取失败，顺子")
		return outCards
	}
	outCardTypeGroup := lgc.getMinCardType(outCardType, more, tuoGuanUser.ctg)
	if outCardTypeGroup == nil {
		//mlog.Warn("机器人数据获取失败，顺子")
		return outCards
	}

	mlog.NewDefaultLogger().SetData(fmt.Sprintf("出順子 outCardTypeGroup.nMin[%v] nMax[%v] ", outCardTypeGroup.nMin, outCardTypeGroup.nMax)).Debug("AI組牌")
	for i := outCardTypeGroup.nMin; i <= outCardTypeGroup.nMax; i++ {
		outCards = append(outCards, i)
	}

	lgc.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
	lgc.reduceCardType(tuoGuanUser, cardTypeStraight, 1)

	mlog.NewDefaultLogger().SetData(fmt.Sprintf("出順子 %v", outCards)).Debug("AI組牌")
	return outCards
}

//连对
func (lgc *DDZRobotLogic) aiBoxOutLiandui(userid int32, more int32) []int32 {
	outCardType := cardTypeThreeStraight
	outCards := make([]int32, 0)
	defer func() {
		if len(outCards) > 0 {
			mlog.Debug("获取连对牌型%v", outCards)
		}
	}()
	tuoGuanUser := lgc.sandBoxUser[userid]
	if tuoGuanUser == nil {
		mlog.Warn("机器人数据获取失败，连对")
		return outCards
	}
	outCardTypeGroup := lgc.getMinCardType(outCardType, more, tuoGuanUser.ctg)
	if outCardTypeGroup == nil {
		//mlog.Warn("机器人数据获取失败，连对")
		return outCards
	}

	for i := outCardTypeGroup.nMin; i <= outCardTypeGroup.nMax; i++ {
		outCards = append(outCards, i)
		outCards = append(outCards, i)
	}

	lgc.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
	lgc.reduceCardType(tuoGuanUser, cardTypeThreeStraight, 1)
	return outCards
}

//出牌 飞机
func (lgc *DDZRobotLogic) aiBoxOutFeiji(userid, sum, more int32) []int32 {
	outCardType := cardTypePlaneAndWing
	outCards := make([]int32, 0)
	defer func() {
		if len(outCards) > 0 {
			mlog.Debug("获取飞机牌型%v", outCards)
		}
	}()
	robotUser := lgc.sandBoxUser[userid]
	if robotUser == nil {
		mlog.Warn("机器人数据获取失败，飞机")
		return outCards
	}
	outCardTypeGroup := lgc.getMinCardType(outCardType, more, robotUser.ctg)
	if outCardTypeGroup == nil {
		//mlog.Warn("机器人数据获取失败，飞机")
		return outCards
	}
	if sum < 2 {
		return outCards //
	}
	size := outCardTypeGroup.nMax - outCardTypeGroup.nMin + 1
	if size > sum {
		size = sum
	}
	for i := outCardTypeGroup.nMax; i > outCardTypeGroup.nMax-size; i-- {
		outCards = append(outCards, i)
		outCards = append(outCards, i)
		outCards = append(outCards, i)
	}
	lgc.eraseResult(robotUser.ctg, outCardTypeGroup)
	lgc.reduceCardType(robotUser, cardTypePlaneAndWing, 1)
	//有size单牌
	bNBeedReset := false
	if robotUser.cardTypesCount[cardTypeSingle] >= size {
		singleList := make([]int32, 0)
		for i := int32(0); i < size; i++ {
			min := int32(0)
			robotUser.ctg.SingleCards, min = lgc.minOneSingle(robotUser.ctg.SingleCards)
			if min <= 0 {
				if lgc.getHand(userid) <= size+3 {
					singOne := lgc.getMinCardType(cardTypeSingle, 0, robotUser.ctg)
					if singOne == nil {
						bNBeedReset = true
						break
					}
					singleList = append(singleList, singOne.nMax)
					lgc.reduceCardType(robotUser, cardTypeSingle, 1)
					lgc.eraseResult(robotUser.ctg, singOne)
					continue
				} else {
					bNBeedReset = true
					break
				}

			}
			singleList = append(singleList, min)
			lgc.reduceCardType(robotUser, cardTypeSingle, 1)

		}
		if len(singleList) == int(size) {
			outCards = append(outCards, singleList...)
			return outCards
		}
	}
	//有两对的情况
	if robotUser.cardTypesCount[cardTypeDouble] >= size {
		doubleList := make([]int32, 0)
		for i := int32(0); i < size; i++ {
			double := new(cardTypeGroup)
			if lgc.getHand(userid) <= size+3 {
				double = lgc.getMinCardType(cardTypeDouble, 0, robotUser.ctg)
			} else {
				double = lgc.getMinCardType(cardTypeDouble, 2, robotUser.ctg)
			}
			if double == nil {
				mlog.Warn("机器人数据获取失败，飞机")
				bNBeedReset = true
				break
			}
			doubleList = append(doubleList, double.nMax)
			doubleList = append(doubleList, double.nMax)
			lgc.reduceCardType(robotUser, cardTypeDouble, 1)
			lgc.eraseResult(robotUser.ctg, double)
		}
		if len(doubleList) == int(size)*2 {
			outCards = append(outCards, doubleList...)
			return outCards
		}

	}
	//检测手牌时候除了三条之外还有1~2张2 或一个王
	up := lgc.getGameLogic().gamePoker.GetUserPoker(userid)
	if up == nil {
		mlog.Warn("系统错误robotOutCards，获取玩家牌型失败")
		return outCards
	}
	//有1对的情况
	if bNBeedReset {
		bNBeedReset = false
		lgc.recombination(userid, up.Cards)
		lgc.eraseResult(robotUser.ctg, outCardTypeGroup)
		lgc.reduceCardType(robotUser, cardTypePlaneAndWing, 1)
	}
	if robotUser.cardTypesCount[cardTypeDouble]*2 >= size {
		doubleList := make([]int32, 0)
		x := size / 2 * 2
		if x < size {
			for i := int32(0); i < x/2; i++ {
				double := new(cardTypeGroup)
				if lgc.getHand(userid) <= size+3 {
					double = lgc.getMinCardType(cardTypeDouble, 0, robotUser.ctg)
				} else {
					double = lgc.getMinCardType(cardTypeDouble, 2, robotUser.ctg)
				}
				if double == nil {
					mlog.Warn("机器人数据获取失败，飞机")
					bNBeedReset = true
					break
				}
				doubleList = append(doubleList, double.nMax)
				doubleList = append(doubleList, double.nMax)
				lgc.reduceCardType(robotUser, cardTypeDouble, 1)
				lgc.eraseResult(robotUser.ctg, double)
			}
			double := new(cardTypeGroup)
			if lgc.getHand(userid) <= size+3 {
				double = lgc.getMinCardType(cardTypeDouble, 0, robotUser.ctg)
			} else {
				double = lgc.getMinCardType(cardTypeDouble, 2, robotUser.ctg)
			}

			if double == nil {
				mlog.Warn("机器人数据获取失败，飞机")
				bNBeedReset = true
			}
			if double != nil {
				doubleList = append(doubleList, double.nMax)
				lgc.reduceCardType(robotUser, cardTypeDouble, 1)
				lgc.eraseResult(robotUser.ctg, double)
			}
		} else {
			for i := int32(0); i < size/2; i++ {
				double := new(cardTypeGroup)
				if lgc.getHand(userid) <= size+3 {
					double = lgc.getMinCardType(cardTypeDouble, 0, robotUser.ctg)
				} else {
					double = lgc.getMinCardType(cardTypeDouble, 2, robotUser.ctg)
				}
				if double == nil {
					// mlog.Error("机器人数据获取失败，飞机")
					bNBeedReset = true
					break
				}
				doubleList = append(doubleList, double.nMax)
				doubleList = append(doubleList, double.nMax)
				lgc.reduceCardType(robotUser, cardTypeDouble, 1)
				lgc.eraseResult(robotUser.ctg, double)
			}
		}
		if len(doubleList) == int(size) {
			outCards = append(outCards, doubleList...)
			return outCards
		}
	}

	return outCards
}

//出牌 三带
func (lgc *DDZRobotLogic) aiboxOutThreeCards(userid int32, more int32) []int32 {
	//三带
	outCardType := cardTypeThree
	outCards := make([]int32, 0)
	defer func() {
		if len(outCards) > 0 {
			mlog.Debug("获取三带牌型%v", outCards)
		}
	}()
	tuoGuanUser := lgc.sandBoxUser[userid]
	if tuoGuanUser == nil {
		mlog.Debug("机器人数据获取失败，三带")
		return outCards
	}

	outCardTypeGroup := lgc.getMinCardType(outCardType, more, tuoGuanUser.ctg)
	if outCardTypeGroup == nil {
		if lgc.getHand(userid) <= 3 {
			outCardTypeGroup = lgc.getMinCardType(outCardType, 0, tuoGuanUser.ctg)
			if outCardTypeGroup == nil {
				//mlog.Warn("机器人数据获取失败，三带")
				return outCards
			}
		} else {
			//mlog.Warn("机器人数据获取失败，三带")
			return outCards
		}
	}

	if outCardTypeGroup.nMax != poker.PokerNumberTwo && outCardTypeGroup.nMax < more {
		mlog.Debug("機器人取三帶失敗. nMax[%v] more[%v]", outCardTypeGroup.nMax, more)
		return outCards
	}

	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	lgc.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
	lgc.reduceCardType(tuoGuanUser, cardTypeThree, 1)

	//mlog.Debug("最後三帶的出牌長度[%v] 牌[%v] more[%v]", len(lgc.outInfo.cg.Cards), lgc.outInfo.cg.Cards, more)
	if (len(lgc.outInfo.cg.Cards) == len(outCards) && outCardTypeGroup.nMax > more) ||
		(len(lgc.outInfo.cg.Cards) == len(outCards) && outCardTypeGroup.nMax == poker.PokerNumberTwo) {
		mlog.Debug("跟三帶 [%v]", outCards)
		return outCards
	}

	singleMin := int32(17)
	doubleMin := int32(17)
	for _, ct := range tuoGuanUser.ctg.Ct {
		if ct.byteCardType == cardTypeSingle && cardMap[ct.nMax] < cardMap[singleMin] {
			singleMin = ct.nMax
		}
		if ct.byteCardType == cardTypeDouble && cardMap[ct.nMax] < cardMap[doubleMin] {
			doubleMin = ct.nMax
		}
	}
	for _, i := range tuoGuanUser.ctg.SingleCards {
		if cardMap[i] < cardMap[singleMin] {
			singleMin = i
		}
	}

	if cardMap[singleMin] <= cardMap[doubleMin] && singleMin != 17 {
		//出三帶一, 不帶鬼跟2出或最後一手直接出
		if (singleMin != 15 && singleMin != 16 && singleMin != 2) || tuoGuanUser.ctg.hands <= 2 && len(lgc.outInfo.cg.Cards) == len(outCards) {
			mlog.Debug("出三帶一, 不帶鬼跟2出或最後一手直接出 ")
			outCards = append(outCards, singleMin)
			lgc.reduceCardType(tuoGuanUser, cardTypeSingle, 1)
			return outCards
		}
	} else if cardMap[doubleMin] <= cardMap[singleMin] && doubleMin != 17 {
		//出三帶二, 不帶2出或最後一手直接出
		if doubleMin != 2 || tuoGuanUser.ctg.hands <= 2 && len(lgc.outInfo.cg.Cards) == len(outCards) {
			//mlog.Debug("出三帶二, 不帶鬼跟2出或最後一手直接出 ")
			double := lgc.getMinCardType(cardTypeDouble, 0, tuoGuanUser.ctg)
			outCards = append(outCards, doubleMin)
			outCards = append(outCards, doubleMin)
			lgc.eraseResult(tuoGuanUser.ctg, double)
			lgc.reduceCardType(tuoGuanUser, cardTypeDouble, 1)
			return outCards
		}
	}

	return outCards
}

//炸弹
func (lgc *DDZRobotLogic) aiBoxOutZhadan(userid int32, more int32) []int32 {
	outCardType := cardTypeBomb
	outCards := make([]int32, 0)
	defer func() {
		if len(outCards) > 0 {
			mlog.Debug("获取炸弹牌型%v", outCards)
		}
	}()
	tuoGuanUser := lgc.sandBoxUser[userid]
	if tuoGuanUser == nil {
		mlog.Debug("获取机器人数据失败，炸弹")
		return outCards
	}

	outCardTypeGroup := lgc.getMinCardType(outCardType, more, tuoGuanUser.ctg)
	if outCardTypeGroup == nil {
		// mlog.Error("获取机器人数据失败，炸弹")
		return outCards
	}
	if outCardTypeGroup.nMax <= more {
		return outCards
	}

	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	lgc.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
	lgc.reduceCardType(tuoGuanUser, cardTypeBomb, 1)
	return outCards
}
