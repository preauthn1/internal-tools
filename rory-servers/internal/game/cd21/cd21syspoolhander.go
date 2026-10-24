package cd21

import (
	"math"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

//收到RPC消息解析處理
func (cd21 *CD21gameLogic) OnRecvRPCMessage(m *msg.Message) {
	rm, ok := m.MsgData.(*rpcmsg.RPCMessage)
	if !ok {
		return
	}
	mmbyte := rm.MsgData.([]byte)
	if len(mmbyte) == 0 {
		return
	}

	if m.SClassID == int32(netproto.CD21_GameMessageClassID_Cd21SysPoolRetID) {
		cd21.getSyspoolData(m)
	}

}

//從DB獲得的水池資訊
func (cd21 *CD21gameLogic) getSyspoolData(m *msg.Message) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Cd21SysPoolData{}
		return mm
	})
	detail := rmm.(*netproto.Cd21SysPoolData)
	if cd21.sysroundid%cd21.roundidlimit == 1 {
		if cd21.sysroundid >= cd21.roundidlimit {
			cd21.sysroundid = 0
		}
		cd21.syspoolData = detail
	}
}

//重新計算過後的水池資訊存入DB
func (cd21 *CD21gameLogic) setSyspoolData(m *netproto.Cd21SysPoolData) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_Cd21SysPoolID)
	pdata := new(netproto.Cd21SysPoolData)
	pdata.FlagID = proto.Int32(util.StrToInt32(cd21.GetServerConfig().Flag))
	pdata.SumBetMoney = m.SumBetMoney
	pdata.SumBetCount = m.SumBetCount
	pdata.PoolAmount = m.PoolAmount
	pmsg.MsgData = pdata

	cd21.SendRPCMessageFunc(false, pmsg)
}

//該放水囉
func (cd21 *CD21gameLogic) poolWater() {
	if cd21.syspoolData != nil {
		cd21.GameLog.AddGameStep("第 %d 局，判定局數: %d, 水池金額: %d", cd21.sysroundid, cd21.roundidlimit, *cd21.syspoolData.PoolAmount)
		cd21.GameLog.AddGameStep("此局水池變化量 %d", cd21.poolAmount)

	} else {
		cd21.GameLog.AddGameStep("第一局水池讀取中")
	}
	poolWater := new(netproto.Cd21SysPoolData)
	poolWater.FlagID = proto.Int32(util.StrToInt32(cd21.GetServerConfig().Flag))
	poolWater.SumBetMoney = proto.Int64(cd21.sumBetMoney)
	poolWater.SumBetCount = proto.Int64(cd21.sumBetCount)
	poolWater.PoolAmount = proto.Int64(cd21.poolAmount)
	cd21.setSyspoolData(poolWater)
}

//根據水池均衡校正
func (cd21 *CD21gameLogic) syspoolcorrect() {
	cd21.GameLog.AddGameStep("第 %d 回合", cd21.sysroundid)
	ust := &netproto.Cd21SysPoolData{}
	data, _ := proto.Marshal(cd21.syspoolData)
	err := proto.Unmarshal(data, ust)
	if err != nil {
		//自然機率
		cd21.GameLog.AddGameStep("水池讀取中")
		mlog.System("抓不到水池值")
		//cd21.sendCardToBanker()
		return
	} else {
		cd21.GameLog.AddGameStep("判定局數: %d, 水池金額: %d", cd21.roundidlimit, *cd21.syspoolData.PoolAmount)
		cd21.GameLog.AddGameStep("判定局數內總下注金額: %d", *cd21.syspoolData.SumBetMoney)
		cd21.GameLog.AddGameStep("判定局數內總下注次數: %d", *cd21.syspoolData.SumBetCount)
	}
	sumBetMoney := ust.SumBetMoney
	sumBetCount := ust.SumBetCount
	poolAmount := ust.PoolAmount
	//檢查牌點
	for _, judge := range cd21.judgeRange {

		if *poolAmount >= int64(int32(*sumBetMoney / *sumBetCount)*cd21.defined*judge.value.numrange) {
			cd21.zhuangdianvariable = judge.value.param
			cd21.playervariable = judge.value.judge
			break
		}
	}
	if cd21.zhuangdianvariable == 0 && cd21.playervariable == 0 && cd21.poolAmount < int64(int32(*sumBetMoney / *sumBetCount)*cd21.defined*cd21.judgeRange[10].value.numrange) {
		cd21.zhuangdianvariable = cd21.judgeRange[10].value.param
		cd21.playervariable = cd21.judgeRange[10].value.judge
	}
}

//根據水池補牌
func (cd21 *CD21gameLogic) syspoolbupai() {
	getdian := int32(0)
	for i := 0; i < len(cd21.dianJudge); i++ {
		if cd21.getCardPoint(cd21.zhungCard) == cd21.dianJudge[i] {
			random := util.RndInt(1, 101)
			for _, judge := range cd21.zhuangDianJudge {
				if judge.index == int32(i) {
					ju := int32(0)
					for j := 0; j < 5; j++ {
						if j < 4 {
							ju += judge.weight[j] + cd21.zhuangdianvariable
							if random < int(ju) {
								getdian = judge.dian[j]
								mlog.System("getdian:%d", getdian)
								break
							}
						} else {
							//自然機率
							getdian = 0
						}
					}
				}
			}
		}
	}
	mlog.System("getdian after:%d", getdian)
	//由 getdian 得到莊家欲要得點數
	if getdian == 0 {
		cd21.GameLog.AddGameStep("根據水線，期望自然要牌")
		cd21.sendCardToBanker()
		return
	} else {
		mlog.System("getdian:%d", getdian)
		thirdcardpoint := getdian - cd21.getCardPoint(cd21.zhungCard)
		cd21.GameLog.AddGameStep("根據水線，莊家期望要牌結果 %v", getdian)
		ok, takeCard := cd21.poker.TakeCardNumber(0, int32(thirdcardpoint))
		mlog.System("thirdcardpoint:%d", thirdcardpoint)
		pokerCard := new(poker.Card)
		if !ok {
			//取牌失敗 隨機發牌
			cd21.sendCardToBanker()
			cd21.GameLog.AddGameStep("取牌失敗 自然要牌")
			return
		} else {
			pokerCard.CardNumber = takeCard.CardNumber
			pokerCard.Color = takeCard.Color
		}
		cd21.zhungCard.Cards = append(cd21.zhungCard.Cards, pokerCard)
	}
	cd21.GameLog.AddGameStep("莊家要牌結果 %v", cd21.zhungCard.String())
	if cd21.getCardPoint(cd21.zhungCard) > 21 {
		cd21.zhuangBust = true
		cd21.GameLog.AddGameStep("莊家爆牌, 牌點總和 %d", cd21.getCardPoint(cd21.zhungCard))

	}
	time := math.Ceil(float64(len(cd21.zhungCard.Cards)) * 0.7)
	cd21.zhuangFanPaiTime = int32(time) + 1
	cd21.BoardcastMessage(cd21.getFaZhuangPaiMsg(cd21.zhungCard))
}

//根據水池發玩家牌
func (cd21 *CD21gameLogic) syspoolfapai(user *cd21userdata) bool {
	natureDuiZi := true
	natureLucy20 := true
	//決定玩家是否可開對子或lucky20
	random1 := util.RndInt(1, 101)
	random2 := util.RndInt(1, 101)
	if int32(random1) <= cd21.duiziconfig+cd21.playervariable {
		natureDuiZi = false //不能開對子
	}
	if int32(random2) <= cd21.lucky20config+cd21.playervariable {
		natureLucy20 = false //不能開拉起20
	}
	if natureDuiZi && natureLucy20 {
		return true //自然發牌true, 已控制發牌false
	}

	luckycontroll := false
	duizicontroll := false
	for _, bettype := range user.userbet {
		if bettype.cd21type == LUCKY20 {
			cd21.GameLog.AddGameStep("玩家%d ,有下lucky20，根據水線需控制發牌", user.user.UserID)
			luckycontroll = true //玩家有下lucky20，需控制
		}
		if bettype.cd21type == DUIZI {
			cd21.GameLog.AddGameStep("玩家%d ,有下對子，根據水線需控制發牌", user.user.UserID)
			duizicontroll = true //玩家有下對子，需控制
		}
	}
	if !luckycontroll && !duizicontroll {
		cd21.GameLog.AddGameStep("玩家%d ,沒有下注lucky20及對子，不需控制發牌", user.user.UserID)
		return true //自然發牌true, 已控制發牌false
	}
	controllPoker := new(poker.CardGroup)
	if duizicontroll {
		//獲取非 Duizi 牌型, 最多找24次
		noDuiPoker := cd21.assignNoDuiPoker(24)
		if noDuiPoker.Cards[0].CardNumber == noDuiPoker.Cards[1].CardNumber {
			cd21.GameLog.AddGameStep("玩家 %d ,換牌失敗, 牌型 %v ", user.user.UserID, noDuiPoker.String())
		} else {
			cd21.GameLog.AddGameStep("玩家%d ,控制發牌, 牌型 %v", user.user.UserID, noDuiPoker.String())
		}
		controllPoker = noDuiPoker
	} else if luckycontroll {
		//獲取非 Lucky20 牌型, 最多找24次
		noLuckyPoker := cd21.assignNoLuckyPoker(24)
		if cd21.getCardPoint(noLuckyPoker) == 20 {
			cd21.GameLog.AddGameStep("玩家 %d ,換牌失敗, 牌型 %v ", user.user.UserID, noLuckyPoker.String())
		} else {
			cd21.GameLog.AddGameStep("玩家%d ,控制發牌, 牌型 %v", user.user.UserID, noLuckyPoker.String())
		}
		controllPoker = noLuckyPoker
	}

	cd21.setUserCards(user, controllPoker)
	return false //自然發牌true, 已控制發牌false
}

//確認牌型不可為 Duizi
func (cd21 *CD21gameLogic) assignNoDuiPoker(number int) *poker.CardGroup {
	result := new(poker.CardGroup)
	for runLimit := 0; runLimit < number; runLimit++ {
		result = cd21.poker.Deal(0, 2)
		if result.Cards[0].CardNumber != result.Cards[1].CardNumber {
			break
		}
		if result != nil && runLimit < number-1 {
			cd21.poker.AddCardGroup(0, result) //把牌放回去
		}
	}
	return result
}

//確認牌型不可為 Lucky20
func (cd21 *CD21gameLogic) assignNoLuckyPoker(number int) *poker.CardGroup {
	result := new(poker.CardGroup)
	//玩家沒有下對子
	for runLimit := 0; runLimit < number; runLimit++ {
		result = cd21.poker.Deal(0, 2)
		if cd21.getCardPoint(result) != 20 {
			break
		}
		if result != nil && runLimit < number-1 {
			cd21.poker.AddCardGroup(0, result) //把牌放回去
		}
	}
	return result
}
