package pirate

// import (
// 	"testing"

// 	"example.invalid/golang/protobuf/proto"

// 	"example.invalid/ServerTeam/Servers/pkg/conf"
// 	"example.invalid/ServerTeam/Servers/pkg/config"
// 	"example.invalid/ServerTeam/Servers/pkg/event"
// 	"example.invalid/ServerTeam/Servers/pkg/fixtures"
// 	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
// 	"example.invalid/ServerTeam/Servers/pkg/roomserver"
// 	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
// )

// func TestLogic_Game_Spin_WinFreeSpin(t *testing.T) {

// 	conf.CLogLevel = mlog.DebugLevel
// 	conf.TimerFastMode = 1000 //40
// 	conf.TestMode = true

// 	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.ConfigGroup) gamelogic.GameLogic {
// 		return NewLogic(t, tfManager, cfgs)
// 	}

// 	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

// 	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)
// 	{
// 		logic.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(402), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(402), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(402), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(402), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(402), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(402), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(402), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(402), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(402), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(402), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(402), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(402), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(402), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(402), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(402), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
// 		}
// 	}

// 	logic.Spin(logic.DefaultSymbolGroupID)

// 	if logic.FreeSpin == nil {
// 		t.Errorf("Game should win free spin")
// 		return
// 	}

// }

// func TestLogic_Game_Spin_WinJackpot(t *testing.T) {

// 	conf.CLogLevel = mlog.DebugLevel
// 	conf.TimerFastMode = 1000 //40
// 	conf.TestMode = true

// 	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.ConfigGroup) gamelogic.GameLogic {
// 		return NewLogic(t, tfManager, cfgs)
// 	}

// 	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

// 	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)
// 	{
// 		logic.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
// 			&protodata.SlotSymbolDistribution{GroupId: proto.Uint32(401), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
// 		}
// 	}

// 	logic.GameBetMoney = int64(logic.MaxBetCount)
// 	logic.Spin(logic.DefaultSymbolGroupID)

// 	if logic.GameWinJackpot == 0 {
// 		t.Errorf("Game should win jackpot")
// 	}

// }
