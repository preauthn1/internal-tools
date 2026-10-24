package sss

import (
	"os"
	"path"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

//校验顺子
func TestCkDunSantiao(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("dunSantiao")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		ok := ckDunSantiao(cards)

		if !ok {
			t.Errorf("校验顺子 失敗")
		}
	}
}

//获取一对
func TestGet1d(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("get1d")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		get1dCardGroup := s.Get1d(cards)

		if get1dCardGroup == nil {
			t.Errorf("获取一对 nil 錯誤")
		}

		if get1dCardGroup != nil {
			for k, v := range get1dCardGroup.Cards {
				t.Logf("k[%v] cardnumber[%v] color[%v]", k, v.CardNumber, v.Color)
			}

			if len(get1dCardGroup.Cards) != 2 {
				t.Errorf("获取一对 錯誤")
			}
		}
	}
}

//获取两对牌型
func TestGet2d(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("get2d")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		get2dCardGroup := s.Get2d(cards)

		if get2dCardGroup == nil {
			t.Errorf("获取两对牌型 nil 錯誤")
		}

		if get2dCardGroup != nil {
			for k, v := range get2dCardGroup.Cards {
				t.Logf("k[%v] cardnumber[%v] color[%v]", k, v.CardNumber, v.Color)
			}

			if len(get2dCardGroup.Cards) != 4 {
				t.Errorf("获取两对牌型 錯誤")
			}
		}
	}
}

//获取三条
func TestGet3(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("get3d")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		get3dCardGroup := s.Get3(cards)

		if get3dCardGroup == nil {
			t.Errorf("获取三条 nil 錯誤")
		}

		if get3dCardGroup != nil {
			for k, v := range get3dCardGroup.Cards {
				t.Logf("k[%v] cardnumber[%v] color[%v]", k, v.CardNumber, v.Color)
			}

			if len(get3dCardGroup.Cards) != 3 {
				t.Errorf("获取三条 錯誤")
			}
		}
	}
}

//获取顺子牌型
func TestGetSZ(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("getSZ")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		get3dCardGroup := s.GetSZ(5, cards)

		if get3dCardGroup == nil {
			t.Errorf("获取顺子牌型 nil 錯誤")
		}

		if get3dCardGroup != nil {
			for k, v := range get3dCardGroup.Cards {
				t.Logf("k[%v] cardnumber[%v] color[%v]", k, v.CardNumber, v.Color)
			}

			if len(get3dCardGroup.Cards) != 5 {
				t.Errorf("获取顺子牌型 錯誤")
			}
		}
	}
}

//获取同花牌型
func TestGetTH(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("getTH")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		get3dCardGroup := s.GetTH(5, cards, false)

		if get3dCardGroup == nil {
			t.Errorf("获取同花牌型 nil 錯誤")
		}

		if get3dCardGroup != nil {
			for k, v := range get3dCardGroup.Cards {
				t.Logf("k[%v] cardnumber[%v] color[%v]", k, v.CardNumber, v.Color)
			}

			if len(get3dCardGroup.Cards) != 5 {
				t.Errorf("获取同花牌型 錯誤")
			}
		}
	}
}

//获取葫芦牌型
func TestGetHL(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("getHL")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		get3dCardGroup := s.GetHL(cards)

		if get3dCardGroup == nil {
			t.Errorf("获取葫芦牌型 nil 錯誤")
		}

		if get3dCardGroup != nil {
			for k, v := range get3dCardGroup.Cards {
				t.Logf("k[%v] cardnumber[%v] color[%v]", k, v.CardNumber, v.Color)
			}

			if len(get3dCardGroup.Cards) != 5 {
				t.Errorf("获取葫芦牌型 錯誤")
			}
		}
	}
}

//获取炸弹牌型
func TestGet4(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("get4")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		get3dCardGroup := s.Get4(cards)

		if get3dCardGroup == nil {
			t.Errorf("获取炸弹牌型 nil 錯誤")
		}

		if get3dCardGroup != nil {
			for k, v := range get3dCardGroup.Cards {
				t.Logf("k[%v] cardnumber[%v] color[%v]", k, v.CardNumber, v.Color)
			}

			if len(get3dCardGroup.Cards) != 4 {
				t.Errorf("获取炸弹牌型 錯誤")
			}
		}
	}
}

//获取同花顺牌型
func TestGetTongHuaShun(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("getTongHuaShun")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		get3dCardGroup := s.GetTongHuaShun(5, cards)

		if get3dCardGroup == nil {
			t.Errorf("获取同花顺牌型 nil 錯誤")
		}

		if get3dCardGroup != nil {
			for k, v := range get3dCardGroup.Cards {
				t.Logf("k[%v] cardnumber[%v] color[%v]", k, v.CardNumber, v.Color)
			}

			if len(get3dCardGroup.Cards) != 5 {
				t.Errorf("获取同花顺牌型 錯誤")
			}
		}
	}
}

//校验三顺子
func TestCheckSanShunZi(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("sanShunZi")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		ok := s.CheckSanShunZi(cards)

		if !ok {
			t.Errorf("校验三顺子 nil 錯誤")
		}
	}
}

//校验三同花
func TestCheckSanTongHua(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("sanTongHua")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)

		ok := s.CheckSanTongHua(cards)

		if !ok {
			t.Errorf("校验三同花 nil 錯誤")
		}
	}
}

//校验六对半
func TestCheckLiuDuiBan(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("liuDuiBan")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		ok := s.CheckLiuDuiBan(cards)
		// fmt.Println(cards)
		if !ok {
			t.Errorf("校验六对半 nil 錯誤")
		}
	}
}

//校验五对三条
func TestCheckWuDuiSanTiao(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("wuDuiSanTiao")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		ok := s.CheckWuDuiSanTiao(cards)
		// fmt.Println(cards)
		if !ok {
			t.Errorf("校验五对三条 nil 錯誤")
		}
	}
}

//校验四套三条
func TestCheckSiTaoSanTiao(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("siTaoSanTiao")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		ok := s.CheckSiTaoSanTiao(cards)
		// fmt.Println(cards)
		if !ok {
			t.Errorf("校验四套三条 nil 錯誤")
		}
	}
}

//校验三冲加弹
func TestCheckSanChongJiaDan(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("sanChongJiaDan")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		ok := s.CheckSanChongJiaDan(cards)

		if !ok {
			t.Errorf("校验三冲加弹 nil 錯誤")
		}
	}
}

//校验双怪冲三
func TestCheckShuangGuaiChongSan(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("shuangGuaiChongSan")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		ok := s.CheckShuangGuaiChongSan(cards)

		if !ok {
			t.Errorf("校验双怪冲三 nil 錯誤")
		}
	}
}

//校验凑一色
func TestCheckCouYiSe(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("couYiSe")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		ok := s.CheckCouYiSe(cards)

		if !ok {
			t.Errorf("校验凑一色 nil 錯誤")
		}
	}
}

//校验全小
func TestCheckQuanXiao(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("quanXiao")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		ok := s.CheckQuanXiao(cards)

		if !ok {
			t.Errorf("校验全小 nil 錯誤")
		}
	}
}

//校验全大
func TestCheckQuanDa(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("quanDa")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		ok := s.CheckQuanDa(cards)

		if !ok {
			t.Errorf("校验全大 nil 錯誤")
		}
	}
}

//校验三分天下
func TestCheckSanFenTianXia(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("sanFenTianXia")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		ok := s.CheckSanFenTianXia(cards)

		if !ok {
			t.Errorf("校验三分天下 nil 錯誤")
		}
	}
}

//校验十二皇族
func TestCheckShiErHuangZu(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("shiErHuangZu")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		ok := s.CheckShiErHuangZu(cards)

		if !ok {
			t.Errorf("校验十二皇族 nil 錯誤")
		}
	}
}

//校验十三水
func TestCheckShiSanShui(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("shiSanShui")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		ok := s.CheckShiSanShui(cards)

		if !ok {
			t.Errorf("校验十三水 nil 錯誤")
		}
	}
}

//校验三同花顺
func TestCheckSanTongHuaShun(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("sanTongHuaShun")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		s := new(SSSTable)
		// fmt.Println(cards)
		ok := s.CheckSanTongHuaShun(cards)

		if !ok {
			t.Errorf("校验三同花顺 nil 錯誤")
		}
	}
}

//OutCardByType 出牌牌型
func TestOutCardByType(t *testing.T) {
	var seat int32 = 0
	sss := NewSSSGameLogicByTest("get5")
	sss.RadomCardType(seat)
}
func NewSSSGameLogicByTest(filePoker string) *SSSTable {
	var seat int32 = 0
	var UserID int32 = 0
	fixPoker := fixedPoker(filePoker) //讀取發牌資料 /testfixpoker/xx.dat
	table := enterTable(fixPoker.UserCardLen())
	tfManager := new(event.TimerFuncManager)

	currentPath, _ := os.Getwd()
	configDir := path.Join(currentPath, "config")
	cfgs := new(config.Group)
	cfgs.LoadConfig(configDir, 0)

	sss := new(SSSTable)
	t := new(gamelogic.Table)
	sss.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(table, tfManager, cfgs)
	sss.turn = 0
	sss.gamePoker = poker.NewGamePoker()
	sss.RegisterAfterGameStartFunc(func() {
		sss.init(t)
		sss.regist()
		sss.logicEngine()

	})
	sss.RegisterBeforeGameStartFunc(func() {
		if sss.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
			sss.poker = poker.NewRandomPoker(false, false, 1)
			sss.poker = poker.NewFixedPoker(sss.AllTableUser(), sss.poker.(*poker.RandomPoker), 13)
		} else { //发牌器发牌
			sss.poker = poker.NewFilePoker(sss.Configs, sss.AllTableUser())
			sss.poker = poker.NewFixedFilePoker(sss.poker.(*poker.FilePoker))
		}
	})

	sss.RegisterAfterGameEndFunc(func() {
		sss.clearData()
	})
	sss.poker = fixPoker
	roles := new(SSSRole)
	roles.roleid = UserID
	roles.seat = seat
	roles.isDo = false
	roles.isZuPai = false
	sss.roles = make([]*SSSRole, 0)
	sss.roles = append(sss.roles, roles)

	for _, ur := range sss.AllGameUser() {
		upk := poker.NewUserPoker(ur, fixPoker.GetUserCards(ur.UserID))
		sss.gamePoker.AddUserPoker(upk)
	}

	for _, regf := range sss.BeforeGameStartFunc {
		regf()
	}

	return sss
}

func enterTable(userLen int) *gamelogic.Table {
	table := gamelogic.NewTable(3, 0)
	for i := 0; i < userLen; i++ {
		user := new(gamelogic.GameUser)
		user.SeatNo = int32(i)
		user.UserID = int32(i)

		table.Seats[i].SeatUser = user
		gameState := new(gamelogic.UserGameState)
		user.UserState = gameState
	}
	return table
}

func fixedPoker(filePoker string) *poker.FixedPoker {
	var p poker.Poker
	p = poker.NewTestFilePoker(filePoker)
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	return fp
}

/*
	先簡單的Run 一下, 還有很多需要優化
*/
func TestGameRound(t *testing.T) {
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewSSSGameLogic(t, tfManager, cfgs)
	}

	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewSSSRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	glm, room := fixtures.StartGameLogicByTest()
	glc := glm.GetGameLogic(0).(*SSSTable)

	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.DynamicGameLogicBase.LogicBase)
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	glc.OnGameStart(func() {
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}

	users := fixtures.GetGameUser()
	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	users[0].UserState.OnUserQueue(users[0], room)
	users[1].UserState.OnUserQueue(users[1], room)
	users[2].UserState.OnUserQueue(users[2], room)
	users[3].UserState.OnUserQueue(users[3], room)
	//用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[3], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])
	glc.UserReady(users[3])

	glc.GameStartFunc()

	//發牌
	glc.gameBeginTimeOut()

	glc.ResumeLoop()
	//普通牌型超時隨機組合
	glc.CardTypeChooseTimeOut()
	//比牌階段
	glc.CompareNormal()
	//結算過程
	glc.SettlementTimeOut()

	glc.GameEnd()
}

func NewGameLogicByTest() *SSSTable {
	currentPath, _ := os.Getwd()
	configDir := path.Join(currentPath, "config")
	cfgs := new(config.Group)
	cfgs.LoadConfig(configDir, 0)
	cfg := cfgs.GetConfig(gamelogic.ConfigFileNameGameServer)
	serverInfo := cfgs.GetConfig(gamelogic.ConfigFileNameServerInfo)

	//获取配置信息
	servercfg := gamelogic.GetGameServerConfig(serverInfo, cfg)
	rumanager := gamelogic.NewRoomUserManager(new(player.Manager))
	room := gamelogic.NewRoom(servercfg, rumanager)
	table := room.GetTableByNo(0)

	evtdispatcher := event.NewAsyncEventDispatcher(1)
	tfManager := event.NewTimerFuncManager(evtdispatcher)

	glc := new(SSSTable)
	glc.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(table, tfManager, cfgs)
	glc.gamePoker = poker.NewGamePoker() //创建游戏牌型
	glc.RegisterAfterGameStartFunc(func() {
		glc.init(table)
		glc.regist()
		glc.logicEngine()
	})
	glc.RegisterBeforeGameStartFunc(func() {
		if glc.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
			glc.poker = poker.NewRandomPoker(false, false, 1)
			glc.poker = poker.NewFixedPoker(glc.AllTableUser(), glc.poker.(*poker.RandomPoker), 13)
		} else { //发牌器发牌
			glc.poker = poker.NewFilePoker(glc.Configs, glc.AllTableUser())
			glc.poker = poker.NewFixedFilePoker(glc.poker.(*poker.FilePoker))
		}
	})

	glc.RegisterAfterGameEndFunc(func() {
		glc.clearData()
	})
	roles := new(SSSRole)
	users := fixtures.GetGameUser()
	roles.roleid = users[0].UserID
	roles.seat = glc.GetTable().GetFreeSeat().SeatNo
	roles.isDo = false
	roles.isZuPai = false
	glc.roles = make([]*SSSRole, 0)
	glc.roles = append(glc.roles, roles)

	var p poker.Poker
	p = poker.NewTestFilePoker("get5")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	for _, ur := range glc.AllGameUser() {
		upk := poker.NewUserPoker(ur, fp.GetUserCards(ur.UserID))
		glc.gamePoker.AddUserPoker(upk)
	}

	for _, regf := range glc.BeforeGameStartFunc {
		regf()
	}
	return glc
}

func TestStartGameLogicServer(t *testing.T) {
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewSSSGameLogic(t, tfManager, cfgs)
	}

	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewSSSRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	fixtures.StartGameLogicByTest()
}
