package ddz

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

func TestFapaiBombRate(t *testing.T) {
	conf.CLogLevel = mlog.ErrorLevel
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	glm, _ := fixtures.StartGameLogicByTest()
	users, _ := fixtures.GetGameUsers()
	users = users[0:3]
	glc := glm.GetGameLogic(0).(*DDZGameLogic)

	glc.faPai(true, 17)

	// totalRound 10000, WeightCountAvg 19.43, HandCountAvg 6.59, BigCountAvg 1.28, jdzCountAvg 0.07
	//noShufflePokerCfg := &ddzNoShufflePoker{
	//	map[int]int{0:3,1:4,2:4},
	//	map[int]int{0:5,1:6,2:7},
	//}
	// sumBigCount 25251, totalRound 10000, BigCountAvg 2.53
	//noShufflePokerCfg := &ddzNoShufflePoker{
	//	map[int]int{0:3,1:4,2:4},
	//	map[int]int{0:3,1:4,2:5},
	//}
	// totalRound 10000, WeightCountAvg 17.18, HandCountAvg 6.79, BigCountAvg 0.62, jdzCountAvg 0.10
	noShufflePokerCfg := &ddzNoShufflePoker{
		map[int]int{0: 2, 1: 3, 2: 4},
		map[int]int{0: 3, 1: 4, 2: 5},
	}

	round := func(users []*gamelogic.GameUser) (int32, int32, int32, int32) {
		var p poker.Poker
		p = poker.NewNoShufflePoker(true, false, 1, noShufflePokerCfg.groupPoker)
		p = poker.NewFixedNoShufflePoker(users, p.(*poker.NoShufflePoker), 17, noShufflePokerCfg.alloctake)
		var weightCount, handCount, bigCount, jdzCount int32
		for _, ur := range users {
			cg := p.Deal(ur.UserID, 17)
			dc := NewDefineCard(glc)
			ctg := dc.Combination(cg)
			weight, hands, bigs := dc.CalWeight(ctg, ctg.SingleCards)
			if weight > 20 && hands < 8 {
				jdzCount++
			}
			weightCount += weight
			handCount += hands
			bigCount += bigs
		}
		return weightCount, handCount, bigCount, jdzCount
	}
	totalRound := int32(10000)
	var sumWeightCount, sumHandCount, sumBigCount, jdzCount int32

	for i := int32(0); i < totalRound; i++ {
		wc, hc, bc, jc := round(users)
		sumWeightCount += wc
		sumHandCount += hc
		sumBigCount += bc
		jdzCount += jc
	}

	t.Logf("totalRound %d, WeightCountAvg %.2f, HandCountAvg %.2f, BigCountAvg %.2f, jdzCountAvg %.2f",
		totalRound, float32(sumWeightCount)/float32(totalRound)/float32(3), float32(sumHandCount)/float32(totalRound)/float32(3),
		float32(sumBigCount)/float32(totalRound)/float32(3), float32(jdzCount)/float32(totalRound)/float32(3))
}

func TestNewCombination(t *testing.T) {
	conf.CLogLevel = mlog.ErrorLevel
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	glm, _ := fixtures.StartGameLogicByTest()
	glc := glm.GetGameLogic(0).(*DDZGameLogic)
	var p poker.Poker
	p = poker.NewTestFilePoker("NewCombination")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()
	for i := 0; i < userLen; i++ {
		userId := int32(i)
		up := fp.GetUserCards(userId)
		if up == nil {
			mlog.Error("系统错误robotOutCards，获取玩家牌型失败")
		}
		dc := NewDefineCard(glc)
		uc := dc.Combination(up)
		mlog.Debug("============= %v", uc)

	}
}
