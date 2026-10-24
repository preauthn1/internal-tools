package pachislot

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// PachiGameFunc 取得PachiSlotGame方式
type PachiGameFunc func(gamelogic.GameLogic) *Game

// DefaultModule 預設模組
var DefaultModule = pb.Scenario_NG

// DefaultMode 預設模式
var DefaultMode = pb.ModeType_A

// DefaultRT 預設RT次數
var DefaultRT = 1000

// DefaultBonus 預設大獎
var DefaultBonus = pb.Scenario_BB

// DefaultMinRTLockTime 預設RT次數小於此數字不觸發RT解除
var DefaultMinRTLockTime = 33

// DefaultJACLimit 預設最大JAC在BB累積次數
var DefaultJACLimit = 3

// DefaultBonusSureWinRate 預設在Bonus確定時進入Bonus機率 百分比
var DefaultBonusSureWinRate = 50

// DefaultRoadMapBuf 轉數表預設暫存空間
var DefaultRoadMapBuf = 10

// DefaultGameConfig 預設遊戲設定
var DefaultGameConfig *pb.GameConfig = &pb.GameConfig{ID: 6001, Name: "靠", PayLineGroupID: 6001, PayTableGroupID: 6001, PatternGroupID: 6001, ReelStripsGroupID: 6001}

// DefaultPayLine 預設賠線
var DefaultPayLine []*pb.PayLine = []*pb.PayLine{
	{ID: 600101, GroupID: 6001, LineID: 1, Positions: []*pb.Position{
		{Reel: 1, Row: 1},
		{Reel: 2, Row: 1},
		{Reel: 3, Row: 1},
	}},
	{ID: 600102, GroupID: 6001, LineID: 2, Positions: []*pb.Position{
		{Reel: 1, Row: 2},
		{Reel: 2, Row: 2},
		{Reel: 3, Row: 2},
	}},
	{ID: 600103, GroupID: 6001, LineID: 3, Positions: []*pb.Position{
		{Reel: 1, Row: 3},
		{Reel: 2, Row: 3},
		{Reel: 3, Row: 3},
	}},
	{ID: 600104, GroupID: 6001, LineID: 4, Positions: []*pb.Position{
		{Reel: 1, Row: 1},
		{Reel: 2, Row: 2},
		{Reel: 3, Row: 3},
	}},
	{ID: 600105, GroupID: 6001, LineID: 5, Positions: []*pb.Position{
		{Reel: 1, Row: 3},
		{Reel: 2, Row: 2},
		{Reel: 3, Row: 1},
	}},
}

// DefaultPayTable 預設賠付表
var DefaultPayTable []*pb.PayoutOption = []*pb.PayoutOption{
	{ID: 6001001, GroupID: 6001, ResultType: pb.ResultType_BellLink, PayType: pb.PayType_PTChip, Payout: 10, UseType: pb.Scenario_NG},
	{ID: 6001002, GroupID: 6001, ResultType: pb.ResultType_ReplayLink, PayType: pb.PayType_PTReplay, Payout: 1, UseType: pb.Scenario_NG},
	{ID: 6001003, GroupID: 6001, ResultType: pb.ResultType_WaterMelonLink, PayType: pb.PayType_PTChip, Payout: 5, UseType: pb.Scenario_NG},
	{ID: 6001004, GroupID: 6001, ResultType: pb.ResultType_BerryLeftAny, PayType: pb.PayType_PTChip, Payout: 2, UseType: pb.Scenario_NG},
	{ID: 6001005, GroupID: 6001, ResultType: pb.ResultType_BlueSevenLink, PayType: pb.PayType_PTBonus, Payout: 6001001, UseType: pb.Scenario_NG},
	{ID: 6001006, GroupID: 6001, ResultType: pb.ResultType_DoubleSevenLink, PayType: pb.PayType_PTBonus, Payout: 6001001, UseType: pb.Scenario_NG},
	{ID: 6001007, GroupID: 6001, ResultType: pb.ResultType_RedSevenLink, PayType: pb.PayType_PTBonus, Payout: 6001001, UseType: pb.Scenario_NG},
	{ID: 6001008, GroupID: 6001, ResultType: pb.ResultType_BarLink, PayType: pb.PayType_PTBonus, Payout: 6001002, UseType: pb.Scenario_NG},
	{ID: 6001009, GroupID: 6001, ResultType: pb.ResultType_BlueSevenLink, PayType: pb.PayType_PTChip, Payout: 15, UseType: pb.Scenario_NG},
	{ID: 6001010, GroupID: 6001, ResultType: pb.ResultType_DoubleSevenLink, PayType: pb.PayType_PTChip, Payout: 15, UseType: pb.Scenario_NG},
	{ID: 6001011, GroupID: 6001, ResultType: pb.ResultType_RedSevenLink, PayType: pb.PayType_PTChip, Payout: 15, UseType: pb.Scenario_NG},
	{ID: 6001012, GroupID: 6001, ResultType: pb.ResultType_BarLink, PayType: pb.PayType_PTChip, Payout: 15, UseType: pb.Scenario_NG},
	{ID: 6001013, GroupID: 6001, ResultType: pb.ResultType_BellLink, PayType: pb.PayType_PTChip, Payout: 15, UseType: pb.Scenario_BB},
	{ID: 6001014, GroupID: 6001, ResultType: pb.ResultType_ReplayLink, PayType: pb.PayType_PTChip, Payout: 15, UseType: pb.Scenario_BB},
	{ID: 6001015, GroupID: 6001, ResultType: pb.ResultType_BlueSevenLink, PayType: pb.PayType_PTChip, Payout: 15, UseType: pb.Scenario_BB},
	{ID: 6001016, GroupID: 6001, ResultType: pb.ResultType_DoubleSevenLink, PayType: pb.PayType_PTChip, Payout: 15, UseType: pb.Scenario_BB},
	{ID: 6001017, GroupID: 6001, ResultType: pb.ResultType_RedSevenLink, PayType: pb.PayType_PTChip, Payout: 15, UseType: pb.Scenario_BB},
	{ID: 6001018, GroupID: 6001, ResultType: pb.ResultType_ReplayLink, PayType: pb.PayType_PTBonus, Payout: 6001003, UseType: pb.Scenario_BB},
	{ID: 6001019, GroupID: 6001, ResultType: pb.ResultType_BlueSevenLink, PayType: pb.PayType_PT1G, Payout: 1, UseType: pb.Scenario_BB},
	{ID: 6001020, GroupID: 6001, ResultType: pb.ResultType_DoubleSevenLink, PayType: pb.PayType_PT1G, Payout: 2, UseType: pb.Scenario_BB},
	{ID: 6001021, GroupID: 6001, ResultType: pb.ResultType_RedSevenLink, PayType: pb.PayType_PT1G, Payout: 1, UseType: pb.Scenario_BB},
	{ID: 6001022, GroupID: 6001, ResultType: pb.ResultType_FeatureNone, PayType: pb.PayType_PT1G, Payout: 3, UseType: pb.Scenario_BB},
	{ID: 6001023, GroupID: 6001, ResultType: pb.ResultType_BellLink, PayType: pb.PayType_PTChip, Payout: 15, UseType: pb.Scenario_RB},
	{ID: 6001024, GroupID: 6001, ResultType: pb.ResultType_ReplayLink, PayType: pb.PayType_PTChip, Payout: 15, UseType: pb.Scenario_RB},
	{ID: 6001025, GroupID: 6001, ResultType: pb.ResultType_JACNone, PayType: pb.PayType_PT1G, Payout: 1, UseType: pb.Scenario_RB},
	{ID: 6001026, GroupID: 6001, ResultType: pb.ResultType_BellLink, PayType: pb.PayType_PTChip, Payout: 15, UseType: pb.Scenario_JAC},
	{ID: 6001027, GroupID: 6001, ResultType: pb.ResultType_ReplayLink, PayType: pb.PayType_PTChip, Payout: 15, UseType: pb.Scenario_JAC},
	{ID: 6001028, GroupID: 6001, ResultType: pb.ResultType_JACNone, PayType: pb.PayType_PT1G, Payout: 1, UseType: pb.Scenario_JAC},
}

// DefaultPattern 預設Pattern
var DefaultPattern []*pb.PatternOption = []*pb.PatternOption{
	{ID: 600101, GroupID: 6001, ModeGroupID: 600101, NGResultGroupID: 6001011, ModuleGroupID: 600101, RTReleaseGroupID: 6001},
	{ID: 600102, GroupID: 6001, ModeGroupID: 600102, NGResultGroupID: 6001021, ModuleGroupID: 600101, RTReleaseGroupID: 6001},
	{ID: 600103, GroupID: 6001, ModeGroupID: 600103, NGResultGroupID: 6001031, ModuleGroupID: 600101, RTReleaseGroupID: 6001},
	{ID: 600104, GroupID: 6001, ModeGroupID: 600104, NGResultGroupID: 6001041, ModuleGroupID: 600101, RTReleaseGroupID: 6001},
	{ID: 600105, GroupID: 6001, ModeGroupID: 600105, NGResultGroupID: 6001051, ModuleGroupID: 600101, RTReleaseGroupID: 6001},
	{ID: 600106, GroupID: 6001, ModeGroupID: 600106, NGResultGroupID: 6001061, ModuleGroupID: 600101, RTReleaseGroupID: 6001},
}

// DefaultResultDistribution 預設牌型機率表
var DefaultResultDistribution []*pb.ResultDistributionOption = []*pb.ResultDistributionOption{
	{ID: 1, GroupID: 6001011, Result: pb.ResultType_BellLink, Weight: 9009},
	{ID: 2, GroupID: 6001011, Result: pb.ResultType_ReplayLink, Weight: 13889},
	{ID: 3, GroupID: 6001011, Result: pb.ResultType_WaterMelonLink, Weight: 1099},
	{ID: 4, GroupID: 6001011, Result: pb.ResultType_BerryLeftAny, Weight: 1221},
	{ID: 5, GroupID: 6001011, Result: pb.ResultType_FeatureNone, Weight: 31},
	{ID: 6, GroupID: 6001011, Result: pb.ResultType_ReleaseRT, Weight: 413},
	{ID: 7, GroupID: 6001011, Result: pb.ResultType_None, Weight: 74338},
	{ID: 8, GroupID: 6001021, Result: pb.ResultType_BellLink, Weight: 9009},
	{ID: 9, GroupID: 6001021, Result: pb.ResultType_ReplayLink, Weight: 13889},
	{ID: 10, GroupID: 6001021, Result: pb.ResultType_WaterMelonLink, Weight: 1099},
	{ID: 11, GroupID: 6001021, Result: pb.ResultType_BerryLeftAny, Weight: 1221},
	{ID: 12, GroupID: 6001021, Result: pb.ResultType_FeatureNone, Weight: 27},
	{ID: 13, GroupID: 6001021, Result: pb.ResultType_ReleaseRT, Weight: 415},
	{ID: 14, GroupID: 6001021, Result: pb.ResultType_None, Weight: 74340},
	{ID: 15, GroupID: 6001031, Result: pb.ResultType_BellLink, Weight: 9009},
	{ID: 16, GroupID: 6001031, Result: pb.ResultType_ReplayLink, Weight: 13889},
	{ID: 17, GroupID: 6001031, Result: pb.ResultType_WaterMelonLink, Weight: 1099},
	{ID: 18, GroupID: 6001031, Result: pb.ResultType_BerryLeftAny, Weight: 1221},
	{ID: 19, GroupID: 6001031, Result: pb.ResultType_FeatureNone, Weight: 40},
	{ID: 20, GroupID: 6001031, Result: pb.ResultType_ReleaseRT, Weight: 417},
	{ID: 21, GroupID: 6001031, Result: pb.ResultType_None, Weight: 74326},
	{ID: 22, GroupID: 6001041, Result: pb.ResultType_BellLink, Weight: 9009},
	{ID: 23, GroupID: 6001041, Result: pb.ResultType_ReplayLink, Weight: 13889},
	{ID: 24, GroupID: 6001041, Result: pb.ResultType_WaterMelonLink, Weight: 1099},
	{ID: 25, GroupID: 6001041, Result: pb.ResultType_BerryLeftAny, Weight: 1221},
	{ID: 26, GroupID: 6001041, Result: pb.ResultType_FeatureNone, Weight: 58},
	{ID: 27, GroupID: 6001041, Result: pb.ResultType_ReleaseRT, Weight: 417},
	{ID: 28, GroupID: 6001041, Result: pb.ResultType_None, Weight: 74308},
	{ID: 29, GroupID: 6001051, Result: pb.ResultType_BellLink, Weight: 9009},
	{ID: 30, GroupID: 6001051, Result: pb.ResultType_ReplayLink, Weight: 13889},
	{ID: 31, GroupID: 6001051, Result: pb.ResultType_WaterMelonLink, Weight: 1099},
	{ID: 32, GroupID: 6001051, Result: pb.ResultType_BerryLeftAny, Weight: 1221},
	{ID: 33, GroupID: 6001051, Result: pb.ResultType_FeatureNone, Weight: 69},
	{ID: 34, GroupID: 6001051, Result: pb.ResultType_ReleaseRT, Weight: 418},
	{ID: 35, GroupID: 6001051, Result: pb.ResultType_None, Weight: 74295},
	{ID: 36, GroupID: 6001061, Result: pb.ResultType_BellLink, Weight: 9009},
	{ID: 37, GroupID: 6001061, Result: pb.ResultType_ReplayLink, Weight: 13889},
	{ID: 38, GroupID: 6001061, Result: pb.ResultType_WaterMelonLink, Weight: 1099},
	{ID: 39, GroupID: 6001061, Result: pb.ResultType_BerryLeftAny, Weight: 1221},
	{ID: 40, GroupID: 6001061, Result: pb.ResultType_FeatureNone, Weight: 92},
	{ID: 41, GroupID: 6001061, Result: pb.ResultType_ReleaseRT, Weight: 420},
	{ID: 42, GroupID: 6001061, Result: pb.ResultType_None, Weight: 74270},
	{ID: 43, GroupID: 6001002, Result: pb.ResultType_BellLink, Weight: 77160},
	{ID: 44, GroupID: 6001002, Result: pb.ResultType_ReplayLink, Weight: 22247},
	{ID: 45, GroupID: 6001002, Result: pb.ResultType_FeatureNone, Weight: 5},
	{ID: 46, GroupID: 6001002, Result: pb.ResultType_BlueSevenLink, Weight: 1028},
	{ID: 47, GroupID: 6001002, Result: pb.ResultType_DoubleSevenLink, Weight: 73},
	{ID: 48, GroupID: 6001002, Result: pb.ResultType_BlueSevenLink, Weight: 1028},
	{ID: 49, GroupID: 6001003, Result: pb.ResultType_BellLink, Weight: 50000},
	{ID: 50, GroupID: 6001003, Result: pb.ResultType_ReplayLink, Weight: 50000},
	{ID: 51, GroupID: 6001003, Result: pb.ResultType_JACNone, Weight: 76},
}

// DefaultModuleConfigs 預設遊戲模式設定
var DefaultModuleConfigs []*pb.Module = []*pb.Module{
	{ID: 6001001, GroupID: 600101, ResultProbGroupID: 6001002, Type: pb.Scenario_BB, Count: 30, Cost: 3},
	{ID: 6001002, GroupID: 600101, ResultProbGroupID: 6001003, Type: pb.Scenario_RB, Count: 8, Cost: 1},
	{ID: 6001003, GroupID: 600101, ResultProbGroupID: 6001003, Type: pb.Scenario_JAC, Count: 8, Cost: 1},
}

// DefaultModeConfigs 預設模式設定
var DefaultModeConfigs []*pb.Mode = []*pb.Mode{
	{ID: 6001001, GroupID: 600101, ModeType: pb.ModeType_A, NextModeGroupID: 600111, RTGroupID: 600111, BRGroupID: 600111},
	{ID: 6001002, GroupID: 600101, ModeType: pb.ModeType_B, NextModeGroupID: 600121, RTGroupID: 600121, BRGroupID: 600121},
	{ID: 6001003, GroupID: 600101, ModeType: pb.ModeType_Heaven, NextModeGroupID: 600131, RTGroupID: 600131, BRGroupID: 600131},

	{ID: 6001004, GroupID: 600102, ModeType: pb.ModeType_A, NextModeGroupID: 600112, RTGroupID: 600111, BRGroupID: 600112},
	{ID: 6001005, GroupID: 600102, ModeType: pb.ModeType_B, NextModeGroupID: 600122, RTGroupID: 600121, BRGroupID: 600122},
	{ID: 6001006, GroupID: 600102, ModeType: pb.ModeType_Heaven, NextModeGroupID: 600132, RTGroupID: 600131, BRGroupID: 600132},

	{ID: 6001007, GroupID: 600103, ModeType: pb.ModeType_A, NextModeGroupID: 600113, RTGroupID: 600111, BRGroupID: 600113},
	{ID: 6001008, GroupID: 600103, ModeType: pb.ModeType_B, NextModeGroupID: 600123, RTGroupID: 600121, BRGroupID: 600123},
	{ID: 6001009, GroupID: 600103, ModeType: pb.ModeType_Heaven, NextModeGroupID: 600133, RTGroupID: 600131, BRGroupID: 600133},

	{ID: 6001010, GroupID: 600104, ModeType: pb.ModeType_A, NextModeGroupID: 600114, RTGroupID: 600114, BRGroupID: 600114},
	{ID: 6001011, GroupID: 600104, ModeType: pb.ModeType_B, NextModeGroupID: 600124, RTGroupID: 600124, BRGroupID: 600124},
	{ID: 6001012, GroupID: 600104, ModeType: pb.ModeType_Heaven, NextModeGroupID: 600134, RTGroupID: 600131, BRGroupID: 600134},

	{ID: 6001013, GroupID: 600105, ModeType: pb.ModeType_A, NextModeGroupID: 600115, RTGroupID: 600114, BRGroupID: 600115},
	{ID: 6001014, GroupID: 600105, ModeType: pb.ModeType_B, NextModeGroupID: 600125, RTGroupID: 600124, BRGroupID: 600125},
	{ID: 6001015, GroupID: 600105, ModeType: pb.ModeType_Heaven, NextModeGroupID: 600135, RTGroupID: 600131, BRGroupID: 600135},

	{ID: 6001016, GroupID: 600106, ModeType: pb.ModeType_A, NextModeGroupID: 600116, RTGroupID: 600114, BRGroupID: 600116},
	{ID: 6001017, GroupID: 600106, ModeType: pb.ModeType_B, NextModeGroupID: 600126, RTGroupID: 600124, BRGroupID: 600126},
	{ID: 6001018, GroupID: 600106, ModeType: pb.ModeType_Heaven, NextModeGroupID: 600136, RTGroupID: 600131, BRGroupID: 600136},
}

// DefaultNextModeDistribution 預設模式切換機率表
var DefaultNextModeDistribution []*pb.NextModeDistributionOption = []*pb.NextModeDistributionOption{
	{ID: 60010001, GroupID: 600111, Mode: pb.ModeType_A, Weight: 625},
	{ID: 60010002, GroupID: 600111, Mode: pb.ModeType_B, Weight: 8438},
	{ID: 60010003, GroupID: 600111, Mode: pb.ModeType_Heaven, Weight: 938},

	{ID: 60010004, GroupID: 600112, Mode: pb.ModeType_A, Weight: 625},
	{ID: 60010005, GroupID: 600112, Mode: pb.ModeType_B, Weight: 7891},
	{ID: 60010006, GroupID: 600112, Mode: pb.ModeType_Heaven, Weight: 1484},

	{ID: 60010007, GroupID: 600113, Mode: pb.ModeType_A, Weight: 625},
	{ID: 60010008, GroupID: 600113, Mode: pb.ModeType_B, Weight: 7266},
	{ID: 60010009, GroupID: 600113, Mode: pb.ModeType_Heaven, Weight: 2019},

	{ID: 60010010, GroupID: 600114, Mode: pb.ModeType_A, Weight: 625},
	{ID: 60010011, GroupID: 600114, Mode: pb.ModeType_B, Weight: 7188},
	{ID: 60010012, GroupID: 600114, Mode: pb.ModeType_Heaven, Weight: 2188},

	{ID: 60010013, GroupID: 600115, Mode: pb.ModeType_A, Weight: 625},
	{ID: 60010014, GroupID: 600115, Mode: pb.ModeType_B, Weight: 7344},
	{ID: 60010015, GroupID: 600115, Mode: pb.ModeType_Heaven, Weight: 2031},

	{ID: 60010016, GroupID: 600116, Mode: pb.ModeType_A, Weight: 625},
	{ID: 60010017, GroupID: 600116, Mode: pb.ModeType_B, Weight: 7031},
	{ID: 60010018, GroupID: 600116, Mode: pb.ModeType_Heaven, Weight: 2344},

	{ID: 60010019, GroupID: 600121, Mode: pb.ModeType_A, Weight: 6641},
	{ID: 60010020, GroupID: 600121, Mode: pb.ModeType_B, Weight: 625},
	{ID: 60010021, GroupID: 600121, Mode: pb.ModeType_Heaven, Weight: 2734},

	{ID: 60010022, GroupID: 600122, Mode: pb.ModeType_A, Weight: 6016},
	{ID: 60010023, GroupID: 600122, Mode: pb.ModeType_B, Weight: 625},
	{ID: 60010024, GroupID: 600122, Mode: pb.ModeType_Heaven, Weight: 3359},

	{ID: 60010025, GroupID: 600123, Mode: pb.ModeType_A, Weight: 6406},
	{ID: 60010026, GroupID: 600123, Mode: pb.ModeType_B, Weight: 625},
	{ID: 60010027, GroupID: 600123, Mode: pb.ModeType_Heaven, Weight: 2969},

	{ID: 60010028, GroupID: 600124, Mode: pb.ModeType_A, Weight: 4375},
	{ID: 60010029, GroupID: 600124, Mode: pb.ModeType_B, Weight: 625},
	{ID: 60010030, GroupID: 600124, Mode: pb.ModeType_Heaven, Weight: 5000},

	{ID: 60010031, GroupID: 600125, Mode: pb.ModeType_A, Weight: 6094},
	{ID: 60010032, GroupID: 600125, Mode: pb.ModeType_B, Weight: 625},
	{ID: 60010033, GroupID: 600125, Mode: pb.ModeType_Heaven, Weight: 3281},

	{ID: 60010034, GroupID: 600126, Mode: pb.ModeType_A, Weight: 4375},
	{ID: 60010035, GroupID: 600126, Mode: pb.ModeType_B, Weight: 625},
	{ID: 60010036, GroupID: 600126, Mode: pb.ModeType_Heaven, Weight: 5000},

	{ID: 60010037, GroupID: 600131, Mode: pb.ModeType_A, Weight: 4298},
	{ID: 60010038, GroupID: 600131, Mode: pb.ModeType_B, Weight: 2734},
	{ID: 60010039, GroupID: 600131, Mode: pb.ModeType_Heaven, Weight: 2968},

	{ID: 60010040, GroupID: 600132, Mode: pb.ModeType_A, Weight: 3438},
	{ID: 60010041, GroupID: 600132, Mode: pb.ModeType_B, Weight: 2188},
	{ID: 60010042, GroupID: 600132, Mode: pb.ModeType_Heaven, Weight: 4375},

	{ID: 60010043, GroupID: 600133, Mode: pb.ModeType_A, Weight: 3906},
	{ID: 60010044, GroupID: 600133, Mode: pb.ModeType_B, Weight: 3125},
	{ID: 60010045, GroupID: 600133, Mode: pb.ModeType_Heaven, Weight: 2969},

	{ID: 60010046, GroupID: 600134, Mode: pb.ModeType_A, Weight: 4141},
	{ID: 60010047, GroupID: 600134, Mode: pb.ModeType_B, Weight: 2422},
	{ID: 60010048, GroupID: 600134, Mode: pb.ModeType_Heaven, Weight: 3438},

	{ID: 60010049, GroupID: 600135, Mode: pb.ModeType_A, Weight: 3828},
	{ID: 60010050, GroupID: 600135, Mode: pb.ModeType_B, Weight: 2891},
	{ID: 60010051, GroupID: 600135, Mode: pb.ModeType_Heaven, Weight: 3281},

	{ID: 60010052, GroupID: 600136, Mode: pb.ModeType_A, Weight: 3672},
	{ID: 60010053, GroupID: 600136, Mode: pb.ModeType_B, Weight: 2422},
	{ID: 60010054, GroupID: 600136, Mode: pb.ModeType_Heaven, Weight: 3906},
}

// DefaultBRDistribution 預設大獎機率表
var DefaultBRDistribution []*pb.BRDistributionOption = []*pb.BRDistributionOption{
	{ID: 60010001, GroupID: 600111, Type: pb.Scenario_BB, Weight: 61},
	{ID: 60010002, GroupID: 600111, Type: pb.Scenario_RB, Weight: 39},
	{ID: 60010003, GroupID: 600121, Type: pb.Scenario_BB, Weight: 30},
	{ID: 60010004, GroupID: 600121, Type: pb.Scenario_RB, Weight: 70},
	{ID: 60010005, GroupID: 600131, Type: pb.Scenario_BB, Weight: 75},
	{ID: 60010006, GroupID: 600131, Type: pb.Scenario_RB, Weight: 25},

	{ID: 60010007, GroupID: 600112, Type: pb.Scenario_BB, Weight: 64},
	{ID: 60010008, GroupID: 600112, Type: pb.Scenario_RB, Weight: 36},
	{ID: 60010009, GroupID: 600122, Type: pb.Scenario_BB, Weight: 33},
	{ID: 60010010, GroupID: 600122, Type: pb.Scenario_RB, Weight: 67},
	{ID: 60010011, GroupID: 600132, Type: pb.Scenario_BB, Weight: 61},
	{ID: 60010012, GroupID: 600132, Type: pb.Scenario_RB, Weight: 39},

	{ID: 60010013, GroupID: 600113, Type: pb.Scenario_BB, Weight: 59},
	{ID: 60010014, GroupID: 600113, Type: pb.Scenario_RB, Weight: 41},
	{ID: 60010015, GroupID: 600123, Type: pb.Scenario_BB, Weight: 39},
	{ID: 60010016, GroupID: 600123, Type: pb.Scenario_RB, Weight: 61},
	{ID: 60010017, GroupID: 600133, Type: pb.Scenario_BB, Weight: 75},
	{ID: 60010018, GroupID: 600133, Type: pb.Scenario_RB, Weight: 25},

	{ID: 60010019, GroupID: 600114, Type: pb.Scenario_BB, Weight: 61},
	{ID: 60010020, GroupID: 600114, Type: pb.Scenario_RB, Weight: 39},
	{ID: 60010021, GroupID: 600124, Type: pb.Scenario_BB, Weight: 59},
	{ID: 60010022, GroupID: 600124, Type: pb.Scenario_RB, Weight: 70},
	{ID: 60010023, GroupID: 600134, Type: pb.Scenario_BB, Weight: 61},
	{ID: 60010024, GroupID: 600134, Type: pb.Scenario_RB, Weight: 39},

	{ID: 60010025, GroupID: 600115, Type: pb.Scenario_BB, Weight: 59},
	{ID: 60010026, GroupID: 600115, Type: pb.Scenario_RB, Weight: 41},
	{ID: 60010027, GroupID: 600125, Type: pb.Scenario_BB, Weight: 39},
	{ID: 60010028, GroupID: 600125, Type: pb.Scenario_RB, Weight: 67},
	{ID: 60010029, GroupID: 600135, Type: pb.Scenario_BB, Weight: 75},
	{ID: 60010030, GroupID: 600135, Type: pb.Scenario_RB, Weight: 25},

	{ID: 60010031, GroupID: 600116, Type: pb.Scenario_BB, Weight: 61},
	{ID: 60010032, GroupID: 600116, Type: pb.Scenario_RB, Weight: 39},
	{ID: 60010033, GroupID: 600126, Type: pb.Scenario_BB, Weight: 30},
	{ID: 60010034, GroupID: 600126, Type: pb.Scenario_RB, Weight: 70},
	{ID: 60010035, GroupID: 600136, Type: pb.Scenario_BB, Weight: 61},
	{ID: 60010036, GroupID: 600136, Type: pb.Scenario_RB, Weight: 39},
}

// DefaultRTDistribution 預設RT分配機率表
var DefaultRTDistribution []*pb.RTDistributionOption = []*pb.RTDistributionOption{
	{ID: 60010001, GroupID: 600111, Min: 0, Max: 1, Weight: 50},
	{ID: 60010002, GroupID: 600111, Min: 2, Max: 33, Weight: 200},
	{ID: 60010003, GroupID: 600111, Min: 34, Max: 65, Weight: 200},
	{ID: 60010004, GroupID: 600111, Min: 66, Max: 97, Weight: 200},
	{ID: 60010005, GroupID: 600111, Min: 98, Max: 129, Weight: 20},
	{ID: 60010006, GroupID: 600111, Min: 130, Max: 193, Weight: 20},
	{ID: 60010007, GroupID: 600111, Min: 194, Max: 257, Weight: 20},
	{ID: 60010008, GroupID: 600111, Min: 258, Max: 385, Weight: 100},
	{ID: 60010009, GroupID: 600111, Min: 386, Max: 513, Weight: 20},
	{ID: 60010010, GroupID: 600111, Min: 502, Max: 565, Weight: 450},
	{ID: 60010011, GroupID: 600111, Min: 541, Max: 641, Weight: 20},
	{ID: 60010012, GroupID: 600111, Min: 642, Max: 769, Weight: 20},
	{ID: 60010013, GroupID: 600111, Min: 770, Max: 897, Weight: 20},
	{ID: 60010014, GroupID: 600111, Min: 898, Max: 1025, Weight: 20},
	{ID: 60010015, GroupID: 600111, Min: 902, Max: 965, Weight: 500},
	{ID: 60010016, GroupID: 600111, Min: 1026, Max: 1153, Weight: 20},
	{ID: 60010017, GroupID: 600111, Min: 1154, Max: 1281, Weight: 20},
	{ID: 60010018, GroupID: 600111, Min: 1282, Max: 1409, Weight: 20},
	{ID: 60010019, GroupID: 600111, Min: 1302, Max: 1365, Weight: 550},
	{ID: 60010020, GroupID: 600111, Min: 1410, Max: 1537, Weight: 20},
	{ID: 60010021, GroupID: 600111, Min: 1538, Max: 1665, Weight: 20},
	{ID: 60010022, GroupID: 600111, Min: 1666, Max: 1793, Weight: 700},
	{ID: 60010023, GroupID: 600111, Min: 1794, Max: 1921, Weight: 886},

	{ID: 60010024, GroupID: 600114, Min: 0, Max: 1, Weight: 50},
	{ID: 60010025, GroupID: 600114, Min: 2, Max: 33, Weight: 120},
	{ID: 60010026, GroupID: 600114, Min: 34, Max: 65, Weight: 120},
	{ID: 60010027, GroupID: 600114, Min: 66, Max: 97, Weight: 120},
	{ID: 60010028, GroupID: 600114, Min: 98, Max: 129, Weight: 200},
	{ID: 60010029, GroupID: 600114, Min: 130, Max: 193, Weight: 200},
	{ID: 60010030, GroupID: 600114, Min: 194, Max: 257, Weight: 20},
	{ID: 60010031, GroupID: 600114, Min: 258, Max: 385, Weight: 120},
	{ID: 60010032, GroupID: 600114, Min: 386, Max: 513, Weight: 20},
	{ID: 60010033, GroupID: 600114, Min: 502, Max: 565, Weight: 450},
	{ID: 60010034, GroupID: 600114, Min: 541, Max: 641, Weight: 20},
	{ID: 60010035, GroupID: 600114, Min: 642, Max: 769, Weight: 30},
	{ID: 60010036, GroupID: 600114, Min: 770, Max: 897, Weight: 20},
	{ID: 60010037, GroupID: 600114, Min: 898, Max: 1025, Weight: 20},
	{ID: 60010038, GroupID: 600114, Min: 902, Max: 965, Weight: 500},
	{ID: 60010039, GroupID: 600114, Min: 1026, Max: 1153, Weight: 20},
	{ID: 60010040, GroupID: 600114, Min: 1154, Max: 1281, Weight: 20},
	{ID: 60010041, GroupID: 600114, Min: 1282, Max: 1409, Weight: 20},
	{ID: 60010042, GroupID: 600114, Min: 1302, Max: 1365, Weight: 550},
	{ID: 60010043, GroupID: 600114, Min: 1410, Max: 1537, Weight: 20},
	{ID: 60010044, GroupID: 600114, Min: 1538, Max: 1665, Weight: 20},
	{ID: 60010045, GroupID: 600114, Min: 1666, Max: 1793, Weight: 700},
	{ID: 60010046, GroupID: 600114, Min: 1794, Max: 1921, Weight: 746},

	{ID: 60010047, GroupID: 600121, Min: 0, Max: 1, Weight: 50},
	{ID: 60010048, GroupID: 600121, Min: 2, Max: 33, Weight: 250},
	{ID: 60010049, GroupID: 600121, Min: 34, Max: 65, Weight: 250},
	{ID: 60010050, GroupID: 600121, Min: 66, Max: 97, Weight: 250},
	{ID: 60010051, GroupID: 600121, Min: 98, Max: 129, Weight: 50},
	{ID: 60010052, GroupID: 600121, Min: 130, Max: 193, Weight: 50},
	{ID: 60010053, GroupID: 600121, Min: 194, Max: 257, Weight: 20},
	{ID: 60010054, GroupID: 600121, Min: 258, Max: 385, Weight: 100},
	{ID: 60010055, GroupID: 600121, Min: 386, Max: 513, Weight: 20},
	{ID: 60010056, GroupID: 600121, Min: 402, Max: 465, Weight: 450},
	{ID: 60010057, GroupID: 600121, Min: 541, Max: 641, Weight: 20},
	{ID: 60010058, GroupID: 600121, Min: 642, Max: 769, Weight: 20},
	{ID: 60010059, GroupID: 600121, Min: 770, Max: 897, Weight: 20},
	{ID: 60010060, GroupID: 600121, Min: 802, Max: 865, Weight: 500},
	{ID: 60010061, GroupID: 600121, Min: 898, Max: 1025, Weight: 20},
	{ID: 60010062, GroupID: 600121, Min: 1026, Max: 1153, Weight: 20},
	{ID: 60010063, GroupID: 600121, Min: 1154, Max: 1281, Weight: 20},
	{ID: 60010064, GroupID: 600121, Min: 1202, Max: 1265, Weight: 550},
	{ID: 60010065, GroupID: 600121, Min: 1282, Max: 1409, Weight: 20},
	{ID: 60010066, GroupID: 600121, Min: 1410, Max: 1537, Weight: 20},
	{ID: 60010067, GroupID: 600121, Min: 1538, Max: 1665, Weight: 450},
	{ID: 60010068, GroupID: 600121, Min: 1666, Max: 1793, Weight: 200},
	{ID: 60010069, GroupID: 600121, Min: 1794, Max: 1921, Weight: 746},

	{ID: 60010070, GroupID: 600124, Min: 0, Max: 1, Weight: 50},
	{ID: 60010071, GroupID: 600124, Min: 2, Max: 33, Weight: 120},
	{ID: 60010072, GroupID: 600124, Min: 34, Max: 65, Weight: 120},
	{ID: 60010073, GroupID: 600124, Min: 66, Max: 97, Weight: 120},
	{ID: 60010074, GroupID: 600124, Min: 98, Max: 129, Weight: 250},
	{ID: 60010075, GroupID: 600124, Min: 130, Max: 193, Weight: 250},
	{ID: 60010076, GroupID: 600124, Min: 194, Max: 257, Weight: 20},
	{ID: 60010077, GroupID: 600124, Min: 258, Max: 385, Weight: 120},
	{ID: 60010078, GroupID: 600124, Min: 386, Max: 513, Weight: 20},
	{ID: 60010079, GroupID: 600124, Min: 402, Max: 465, Weight: 450},
	{ID: 60010080, GroupID: 600124, Min: 541, Max: 641, Weight: 20},
	{ID: 60010081, GroupID: 600124, Min: 642, Max: 769, Weight: 20},
	{ID: 60010082, GroupID: 600124, Min: 770, Max: 897, Weight: 20},
	{ID: 60010083, GroupID: 600124, Min: 802, Max: 865, Weight: 500},
	{ID: 60010084, GroupID: 600124, Min: 898, Max: 1025, Weight: 20},
	{ID: 60010085, GroupID: 600124, Min: 1026, Max: 1153, Weight: 20},
	{ID: 60010086, GroupID: 600124, Min: 1154, Max: 1281, Weight: 20},
	{ID: 60010087, GroupID: 600124, Min: 1202, Max: 1265, Weight: 550},
	{ID: 60010088, GroupID: 600124, Min: 1282, Max: 1409, Weight: 20},
	{ID: 60010089, GroupID: 600124, Min: 1410, Max: 1537, Weight: 20},
	{ID: 60010090, GroupID: 600124, Min: 1538, Max: 1665, Weight: 450},
	{ID: 60010091, GroupID: 600124, Min: 1666, Max: 1793, Weight: 200},
	{ID: 60010092, GroupID: 600124, Min: 1794, Max: 1921, Weight: 716},

	{ID: 60010093, GroupID: 600131, Min: 0, Max: 1, Weight: 96},
	{ID: 60010094, GroupID: 600131, Min: 2, Max: 33, Weight: 750},
	{ID: 60010095, GroupID: 600131, Min: 34, Max: 65, Weight: 750},
	{ID: 60010096, GroupID: 600131, Min: 66, Max: 97, Weight: 750},
	{ID: 60010097, GroupID: 600131, Min: 98, Max: 129, Weight: 750},
	{ID: 60010098, GroupID: 600131, Min: 130, Max: 193, Weight: 1000},
}

// DefaultRTReleaseDistribution 預設RT表解除機率表
var DefaultRTReleaseDistribution []*pb.RTReleaseOption = []*pb.RTReleaseOption{
	{ID: 60010001, GroupID: 6001, ResultType: pb.ResultType_WaterMelonLink, Rate: 10, Denominator: 1638, Times: []int32{2}, UseType: pb.Scenario_ALL},
	{ID: 60010002, GroupID: 6001, ResultType: pb.ResultType_BerryLeftAny, Rate: 10, Denominator: 369, Times: []int32{2, 3, 4, 5}, UseType: pb.Scenario_ALL},
	{ID: 60010003, GroupID: 6001, ResultType: pb.ResultType_FeatureNone, Rate: 100, Denominator: 100, Times: []int32{33}, UseType: pb.Scenario_ALL},
	{ID: 60010004, GroupID: 6001, ResultType: pb.ResultType_ReleaseRTonBB, Rate: 10, Denominator: 248, Times: []int32{4}, UseType: pb.Scenario_ALL},
	{ID: 60010005, GroupID: 6001, ResultType: pb.ResultType_ReleaseRTonRB, Rate: 10, Denominator: 247, Times: []int32{2}, UseType: pb.Scenario_ALL},
}

// DefaultReelStrips 預設排列帶
var DefaultReelStrips []*pb.ReelStripSymbol = []*pb.ReelStripSymbol{
	{ID: 60010001, GroupID: 6001, Reel: 1, Row: 1, Symbol: pb.Symbol_Bell},
	{ID: 60010002, GroupID: 6001, Reel: 1, Row: 2, Symbol: pb.Symbol_Replay},
	{ID: 60010003, GroupID: 6001, Reel: 1, Row: 3, Symbol: pb.Symbol_RedSeven},
	{ID: 60010004, GroupID: 6001, Reel: 1, Row: 4, Symbol: pb.Symbol_Berry},
	{ID: 60010005, GroupID: 6001, Reel: 1, Row: 5, Symbol: pb.Symbol_Bell},
	{ID: 60010006, GroupID: 6001, Reel: 1, Row: 6, Symbol: pb.Symbol_Replay},
	{ID: 60010007, GroupID: 6001, Reel: 1, Row: 7, Symbol: pb.Symbol_WaterMelon},
	{ID: 60010008, GroupID: 6001, Reel: 1, Row: 8, Symbol: pb.Symbol_BlueSeven},
	{ID: 60010009, GroupID: 6001, Reel: 1, Row: 9, Symbol: pb.Symbol_BlueSeven},
	{ID: 60010010, GroupID: 6001, Reel: 1, Row: 10, Symbol: pb.Symbol_Bell},
	{ID: 60010011, GroupID: 6001, Reel: 1, Row: 11, Symbol: pb.Symbol_Replay},
	{ID: 60010012, GroupID: 6001, Reel: 1, Row: 12, Symbol: pb.Symbol_Bar},
	{ID: 60010013, GroupID: 6001, Reel: 1, Row: 13, Symbol: pb.Symbol_Berry},
	{ID: 60010014, GroupID: 6001, Reel: 1, Row: 14, Symbol: pb.Symbol_Replay},
	{ID: 60010015, GroupID: 6001, Reel: 1, Row: 15, Symbol: pb.Symbol_Bell},
	{ID: 60010016, GroupID: 6001, Reel: 1, Row: 16, Symbol: pb.Symbol_WaterMelon},
	{ID: 60010017, GroupID: 6001, Reel: 1, Row: 17, Symbol: pb.Symbol_Bar},
	{ID: 60010018, GroupID: 6001, Reel: 1, Row: 18, Symbol: pb.Symbol_Replay},
	{ID: 60010019, GroupID: 6001, Reel: 1, Row: 19, Symbol: pb.Symbol_RedSeven},
	{ID: 60010020, GroupID: 6001, Reel: 1, Row: 20, Symbol: pb.Symbol_Bell},
	{ID: 60010021, GroupID: 6001, Reel: 1, Row: 21, Symbol: pb.Symbol_WaterMelon},

	{ID: 60010022, GroupID: 6001, Reel: 2, Row: 1, Symbol: pb.Symbol_Berry},
	{ID: 60010023, GroupID: 6001, Reel: 2, Row: 2, Symbol: pb.Symbol_Bell},
	{ID: 60010024, GroupID: 6001, Reel: 2, Row: 3, Symbol: pb.Symbol_RedSeven},
	{ID: 60010025, GroupID: 6001, Reel: 2, Row: 4, Symbol: pb.Symbol_Replay},
	{ID: 60010026, GroupID: 6001, Reel: 2, Row: 5, Symbol: pb.Symbol_Berry},
	{ID: 60010027, GroupID: 6001, Reel: 2, Row: 6, Symbol: pb.Symbol_Bell},
	{ID: 60010028, GroupID: 6001, Reel: 2, Row: 7, Symbol: pb.Symbol_Replay},
	{ID: 60010029, GroupID: 6001, Reel: 2, Row: 8, Symbol: pb.Symbol_WaterMelon},
	{ID: 60010030, GroupID: 6001, Reel: 2, Row: 9, Symbol: pb.Symbol_Berry},
	{ID: 60010031, GroupID: 6001, Reel: 2, Row: 10, Symbol: pb.Symbol_WaterMelon},
	{ID: 60010032, GroupID: 6001, Reel: 2, Row: 11, Symbol: pb.Symbol_Bell},
	{ID: 60010033, GroupID: 6001, Reel: 2, Row: 12, Symbol: pb.Symbol_Replay},
	{ID: 60010034, GroupID: 6001, Reel: 2, Row: 13, Symbol: pb.Symbol_Bar},
	{ID: 60010035, GroupID: 6001, Reel: 2, Row: 14, Symbol: pb.Symbol_Berry},
	{ID: 60010036, GroupID: 6001, Reel: 2, Row: 15, Symbol: pb.Symbol_Bell},
	{ID: 60010037, GroupID: 6001, Reel: 2, Row: 16, Symbol: pb.Symbol_BlueSeven},
	{ID: 60010038, GroupID: 6001, Reel: 2, Row: 17, Symbol: pb.Symbol_Replay},
	{ID: 60010039, GroupID: 6001, Reel: 2, Row: 18, Symbol: pb.Symbol_Berry},
	{ID: 60010040, GroupID: 6001, Reel: 2, Row: 19, Symbol: pb.Symbol_Bell},
	{ID: 60010041, GroupID: 6001, Reel: 2, Row: 20, Symbol: pb.Symbol_Replay},
	{ID: 60010042, GroupID: 6001, Reel: 2, Row: 21, Symbol: pb.Symbol_WaterMelon},

	{ID: 60010043, GroupID: 6001, Reel: 3, Row: 1, Symbol: pb.Symbol_Bell},
	{ID: 60010044, GroupID: 6001, Reel: 3, Row: 2, Symbol: pb.Symbol_Bell},
	{ID: 60010045, GroupID: 6001, Reel: 3, Row: 3, Symbol: pb.Symbol_RedSeven},
	{ID: 60010046, GroupID: 6001, Reel: 3, Row: 4, Symbol: pb.Symbol_Replay},
	{ID: 60010047, GroupID: 6001, Reel: 3, Row: 5, Symbol: pb.Symbol_Berry},
	{ID: 60010048, GroupID: 6001, Reel: 3, Row: 6, Symbol: pb.Symbol_Bell},
	{ID: 60010049, GroupID: 6001, Reel: 3, Row: 7, Symbol: pb.Symbol_Bar},
	{ID: 60010050, GroupID: 6001, Reel: 3, Row: 8, Symbol: pb.Symbol_Replay},
	{ID: 60010051, GroupID: 6001, Reel: 3, Row: 9, Symbol: pb.Symbol_Bell},
	{ID: 60010052, GroupID: 6001, Reel: 3, Row: 10, Symbol: pb.Symbol_BlueSeven},
	{ID: 60010053, GroupID: 6001, Reel: 3, Row: 11, Symbol: pb.Symbol_BlueSeven},
	{ID: 60010054, GroupID: 6001, Reel: 3, Row: 12, Symbol: pb.Symbol_WaterMelon},
	{ID: 60010055, GroupID: 6001, Reel: 3, Row: 13, Symbol: pb.Symbol_Replay},
	{ID: 60010056, GroupID: 6001, Reel: 3, Row: 14, Symbol: pb.Symbol_Bell},
	{ID: 60010057, GroupID: 6001, Reel: 3, Row: 15, Symbol: pb.Symbol_Bar},
	{ID: 60010058, GroupID: 6001, Reel: 3, Row: 16, Symbol: pb.Symbol_Replay},
	{ID: 60010059, GroupID: 6001, Reel: 3, Row: 17, Symbol: pb.Symbol_Bell},
	{ID: 60010060, GroupID: 6001, Reel: 3, Row: 18, Symbol: pb.Symbol_WaterMelon},
	{ID: 60010061, GroupID: 6001, Reel: 3, Row: 19, Symbol: pb.Symbol_Berry},
	{ID: 60010062, GroupID: 6001, Reel: 3, Row: 20, Symbol: pb.Symbol_Bar},
	{ID: 60010063, GroupID: 6001, Reel: 3, Row: 21, Symbol: pb.Symbol_Replay},
}

// DefaultProto 預設機率表設定
var DefaultProto *Proto = &Proto{
	GameConfig:            DefaultGameConfig,
	PayLine:               DefaultPayLine,
	PayTable:              DefaultPayTable,
	Pattern:               DefaultPattern,
	ResultDistribution:    DefaultResultDistribution,
	ModuleConfigs:         DefaultModuleConfigs,
	ModeConfigs:           DefaultModeConfigs,
	NextModeDistribution:  DefaultNextModeDistribution,
	BRDistribution:        DefaultBRDistribution,
	RTDistribution:        DefaultRTDistribution,
	RTReleaseDistribution: DefaultRTReleaseDistribution,
	ReelStrips:            DefaultReelStrips,
}
