package ddz

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mylog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

type gameRecord struct {
	PaijuID    int32            //牌局編號
	User       []*rUserInfo     //玩家資料
	DiZhuSeat  int32            //地主座位
	InitCards  []*rUserPoker    //地主初始牌
	DizhuCards *poker.CardGroup //地主牌
	GameStep   []rGameStep      //遊戲歷程
	RobotCount int32            //機器人數
	Odds       int32            //当前所有的倍率综合
	WinnerKind int32            //贏家位置 0.地主 1.農民
	WinnerSeat int32            //贏家位置
}

type rUserInfo struct {
	ID      int32 //ID
	Seat    int32 //座位号
	IsDiZhu bool  //是否地主
	IsRobot bool  //是否機器人
	Level   int32 //機器人智商
}

type rUserPoker struct {
	Seat  int32            //座位号
	Cards *poker.CardGroup //牌型
}

type rGameStep struct {
	Seat     int32            //座位号
	Pass     bool             //是否Pass
	OutCards *poker.CardGroup //出的牌
	OwnCards *poker.CardGroup //剩餘手牌
}

func (record *gameRecord) init() {
	record.PaijuID = 0
	record.User = nil
	record.DiZhuSeat = -1
	record.InitCards = nil
	record.DizhuCards = nil
	record.GameStep = make([]rGameStep, 0)
	record.RobotCount = 0
	record.Odds = 0
	record.WinnerKind = -1
	record.WinnerSeat = -1
}

func (record *gameRecord) setUser(t *gamelogic.Table) {
	record.User = make([]*rUserInfo, len(t.GetTableUser()))
	for _, v := range t.GetTableUser() {
		record.User[v.SeatNo] = new(rUserInfo)
		record.User[v.SeatNo].ID = v.UserID
		record.User[v.SeatNo].Seat = v.SeatNo
		record.User[v.SeatNo].IsDiZhu = false
		record.User[v.SeatNo].IsRobot = v.IsRobot()
		record.User[v.SeatNo].Level = 0
		if record.User[v.SeatNo].IsRobot {
			record.RobotCount++
		}
	}
}

func (record *gameRecord) setDiZhu(UserId int32) {
	for i := 0; i < len(record.User); i++ {
		if record.User[i].ID == UserId {
			record.User[i].IsDiZhu = true
			record.DiZhuSeat = record.User[i].Seat
			break
		}
	}
}

func (record *gameRecord) setInitCards(GamePoker *poker.GamePoker) {
	record.InitCards = make([]*rUserPoker, len(GamePoker.UserPokers))
	for i := 0; i < len(GamePoker.UserPokers); i++ {
		record.InitCards[i] = new(rUserPoker)
		record.InitCards[i].Seat = GamePoker.UserPokers[i].User.SeatNo
		record.InitCards[i].Cards = poker.NewCardGroup(GamePoker.UserPokers[i].Cards.Cards)
		//record.InitCards[i].Cards = new(poker.CardGroup)
		//record.InitCards[i].Cards.Cards = append(GamePoker.UserPokers[i].Cards.Cards[:0:0], GamePoker.UserPokers[i].Cards.Cards...)
	}
}

func (record *gameRecord) setLastCards(Cards []*poker.Card) {
	record.DizhuCards = new(poker.CardGroup)
	record.DizhuCards.Cards = Cards
}

func (record *gameRecord) setRobotLevel(UserId, Level int32) {
	for i := 0; i < len(record.User); i++ {
		if record.User[i].ID == UserId && record.User[i].IsRobot {
			record.User[i].Level = Level
			break
		}
	}
}

func (record *gameRecord) addGameStep(aSeat int32, aPass bool, aOutCards, aOwnCards *poker.CardGroup) {
	step := new(rGameStep)
	step.Seat = aSeat
	step.Pass = aPass
	if aOutCards != nil {
		step.OutCards = poker.NewCardGroup(aOutCards.Cards)
	}
	step.OwnCards = poker.NewCardGroup(aOwnCards.Cards)

	record.GameStep = append(record.GameStep, *step)
}

func (record *gameRecord) checkResult(IsDiZhu bool, WinnerSeat, Odds int32) {
	record.Odds = Odds
	record.WinnerSeat = WinnerSeat

	if IsDiZhu {
		record.WinnerKind = 0
	} else {
		record.WinnerKind = 1
	}

	mylog.AddLog(0, 0, record)

	//err := mylog.GetCollection().Insert(record)
	//if err != nil {
	//	panic(err)
	//}
}
