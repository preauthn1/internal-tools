package cd21

import (
	"example.invalid/golang/protobuf/proto"
)

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"	
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

type bettype struct {
	cd21type int32 //下注类型
	money    int64 //下注的钱
}

type cd21userdata struct {
	userCards    *poker.CardGroup //用户手牌
	splitCards   *poker.CardGroup //用户分牌
	user         *gamelogic.GameUser
	userbet      []*bettype
	userwinblock []int32 //用户贏的區塊
	userDiZu     int64   //用户主遊戲底注
	diZuBase     int64   //用户一單位底注
	money        int64   //用户开始的钱
	moneylast    int64   //用户最开始的钱，用来更新用户中途加钱
	win          int64   //用户下注區赢的钱
	win1         int64   //用户21點赢的钱
	win2         int64   //用户保險赢的钱
	win3         int64   //用户莊家爆牌區赢的钱
	win4         int64   //用户主遊戲輸贏
	win5         int64   //用户分牌輸贏
	win6         int64   //用户投降輸錢
	winDUIZI     int64   //用户對子贏的錢
	winLUCKY20   int64   //用户LUCKY20贏的錢
	winBIGGER13  int64   //用户BIGGER13贏的錢
	winSMALLER13 int64   //用户SMALLER13贏的錢
	winMAINGAME	 int64	 //用户主遊戲贏的錢
	winSPILTGAME int64	 //用户分牌贏的錢
	winZHUANGBAO int64	 //用户莊爆贏的錢
	winBLACKJACK int64	 //用户BJ贏的錢
	betmoney     int64   //用户下注的錢
	tax          int64   //用户被抽的稅
	updatemoney  int64   //用户下注更新的钱
	clearcount   int32   //清理次数
	isStand      bool    //是否已停牌
	isStand1     bool    //是否已停牌
	//isSurrender    bool    //是否已投降 (1.10移除投降)
	isInsurance    bool  //是否已保險
	isBust         bool  //是否爆牌
	isBust1        bool  //是否爆牌
	isDouble       bool  //是否加倍
	isDouble1      bool  //是否加倍
	isSplit        bool  //玩家是否分牌
	isBlackJack    bool  //是否為BJ
	userNoWin      bool  //玩家主遊戲 Push
	userSplitNoWin bool  //玩家主遊戲分牌 Push
}

//新建一个玩家数据实例
func newCD21userdata(user *gamelogic.GameUser) *cd21userdata {
	ud := new(cd21userdata)
	ud.userCards = new(poker.CardGroup)
	ud.splitCards = new(poker.CardGroup)
	ud.user = user
	ud.userbet = make([]*bettype, 0)
	ud.userDiZu = 0
	ud.diZuBase = 0
	ud.money = user.Score
	ud.moneylast = user.Score
	ud.win = 0
	ud.win1 = 0
	ud.win2 = 0
	ud.win3 = 0
	ud.win4 = 0
	ud.win5 = 0
	ud.win6 = 0
	ud.winDUIZI = 0
	ud.winLUCKY20 = 0
	ud.winBIGGER13 = 0
	ud.winSMALLER13 = 0
	ud.winMAINGAME = 0
	ud.winSPILTGAME = 0
	ud.winZHUANGBAO = 0
	ud.winBLACKJACK = 0
	ud.tax = 0
	ud.clearcount = 0
	ud.updatemoney = 0
	ud.isSplit = false
	ud.isDouble = false
	ud.isDouble1 = false
	ud.isBlackJack = false
	ud.isStand = false
	ud.isStand1 = false
	//ud.isSurrender = false (1.10移除投降)
	ud.isBust = false
	ud.isBust1 = false
	ud.isBlackJack = false
	ud.userNoWin = false
	ud.userSplitNoWin = false
	ud.betmoney = 0
	return ud
}

//用户是否可操作
func (ud *cd21userdata) canOpt() bool {
	//1.10移除投降
	if ud.isStand || ud.isBust {
		if ud.isSplit {
			if !ud.isStand1 && !ud.isBust1 {
				return true
			}
		}
		return false
	}

	return true
}

//用户分牌是否可操作(前副牌還未操作完成 也不可操作)
func (ud *cd21userdata) cansplitOpt() bool {
	//1.10移除投降
	if ud.isSplit && (ud.isBust || ud.isStand) {
		if !ud.isStand1 && !ud.isBust1 {
			return true
		}
	}
	return false
}

//设置玩家下注
func (ud *cd21userdata) setUserBet(typevalue int32, money int32) {
	for _, v := range ud.userbet {
		if v.cd21type == typevalue {
			v.money += int64(money)
			ud.money -= int64(money)
			return
		}
	}

	v := new(bettype)
	v.cd21type = typevalue
	v.money = int64(money)
	ud.userbet = append(ud.userbet, v)
	ud.money -= int64(money)
}

//获取用户的操作消息
func (ud *cd21userdata) GetOptMsg(time int32, hit bool) *msg.Message {
	stand := true
	var split bool
	var double bool
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21UserOptID)

	//required int32 SitNo=1;     //座位號
	//required bool hit=2;        //要牌
	//required bool stand=3;      //停牌
	//required bool split=4;      //分牌
	//required bool double=5;     //加倍
	//required bool bust=6;       //爆牌
	//required int32 OptTime=7;   //操作時間

	//若兩張牌一樣時，可分牌
	if ud.userCards.Cards[0].CardNumber == ud.userCards.Cards[1].CardNumber && len(ud.userCards.Cards) == 2 && (ud.user.Score-ud.betmoney-ud.userDiZu) > ud.diZuBase && !ud.isSplit {
		split = true
	}

	//還未進行第一次操作時，可加倍
	if len(ud.userCards.Cards) == 2 && (ud.user.Score-ud.betmoney-ud.userDiZu) > ud.diZuBase {
		double = true
	}

	//如果爆牌，任何行為皆不可操作
	if ud.isBust {
		double = false
		split = false
		stand = false
	}
	pdata := new(netproto.CD21_UserInfo)
	pdata.SitNo = proto.Int32(ud.user.SeatNo)
	pdata.Bust = proto.Bool(ud.isBust)
	pdata.Double = proto.Bool(double)
	pdata.Hit = proto.Bool(hit)
	pdata.Split = proto.Bool(split)
	pdata.Stand = proto.Bool(stand)
	pdata.OptTime = proto.Int32(time)

	pmsg.MsgData = pdata
	return pmsg
}

//设置用户停牌
func (ud *cd21userdata) SetStand() {
	ud.isStand = true
}

//设置用户停牌
func (ud *cd21userdata) SetSplitStand() {
	ud.isStand1 = true
}

//获取用户的操作消息
func (ud *cd21userdata) GetSplitOptMsg(time int32, hit bool) *msg.Message {
	stand := true
	split := false
	var double bool
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21UserOptID)

	//required int32 SitNo=1;     //座位號
	//required bool hit=2;        //要牌
	//required bool stand=3;      //停牌
	//required bool split=4;      //分牌
	//required bool double=5;     //加倍
	//required bool bust=6;       //爆牌
	//required int32 OptTime=7;   //操作時間

	//還未進行第一次操作時，可加倍
	if len(ud.splitCards.Cards) == 2 && (ud.user.Score-ud.betmoney-ud.userDiZu) > ud.diZuBase {
		double = true
	}

	//如果爆牌，任何行為皆不可操作
	if ud.isBust1 {
		double = false
		stand = false
	}
	pdata := new(netproto.CD21_UserInfo)
	splitSitNo := int32(0)
	if ud.user.SeatNo == 1 {
		splitSitNo = 4
	} else if ud.user.SeatNo == 2 {
		splitSitNo = 5
	} else {
		splitSitNo = 3
	}
	pdata.SitNo = proto.Int32(splitSitNo)
	pdata.Bust = proto.Bool(ud.isBust1)
	pdata.Double = proto.Bool(double)
	pdata.Hit = proto.Bool(hit)
	pdata.Split = proto.Bool(split)
	pdata.Stand = proto.Bool(stand)
	pdata.OptTime = proto.Int32(time)

	pmsg.MsgData = pdata
	return pmsg
}

//设置玩家下注
func (ud *cd21userdata) getUserBetMoneyByType(bettype int32) int64 {
	money := int64(0)
	for _, v := range ud.userbet {
		if v.cd21type == bettype {
			money += v.money
		}
	}
	return money
}
