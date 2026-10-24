package helper

import (
	"fmt"
	"io"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// OutputHistory 將輸入的Slot_MSG_GameLogAdd輸出
// 每筆; {UserID};{GameID};{ServerID};{CurrentAmount};{BetMoney};{WinMoney};{FinalScore};{GameType};{LogType};{GameData};{CreateTime}
func OutputHistory(w io.Writer, logs []*netproto.Slot_MSG_GameLogAdd) error {

	for _, row := range logs {
		fmt.Fprintf(w, "%d;%d;%d;%d;%d;%d;%d;%d;%d;%s;%s\n",
			row.GetUserID(),
			row.GetGameID(),
			row.GetServerID(),
			row.GetCurrentAmount(),
			row.GetBetMoney(),
			row.GetWinMoney(),
			row.GetFinalScore(),
			row.GetGameType(),
			row.GetLogType(),
			row.GetGameData(),
			row.GetCreateTime(),
		)
	}

	return nil
}

// AnalysisModeConfig 設定分析模式Config
// * 桌數=1
// * 空的 serverinfo
// * 機器人配置
func AnalysisModeConfig(data *config.Group) {
	// 設定遊戲桌數
	{
		section := data.GetConfigSection("gameserver", "gameserver")
		section["tablecount"] = "1"
		section["deallog"] = "0"
	}

	// ServerInfo 配置
	{
		data.Configs["serverinfo"] = new(config.Config)
	}

	// 設定機器人配置
	{
		section := data.GetConfigSection("gameserver", "robot")
		if section == nil {
			section = config.Section{}
			cfg := data.GetConfig("gameserver")
			cfg.Sections["robot"] = section
		}

		section["enable"] = "1"
		section["logoninterval"] = "10"
		section["minmneyamount"] = "100000"
		section["maxmoneyamount"] = "60000000"
		section["queuecheckwaitspan"] = "3,5"
		section["queueagainwaitspan"] = "3,5"
		section["tablemaxusercount"] = "1"
		section["tableminusercount"] = "0"
		section["gameendleaveodds"] = "1"
		section["autositwaitspan"] = "1,3"
		section["gameactionwaitspan"] = "1,2"
		section["gamereadywaitspan"] = "1,5"
		section["maxusergameendleaveodds"] = "0"
		section["maxcardsrndlvl"] = "1"
		section["againstgameing"] = "0"
		section["matchtablemode"] = "0"
		section["tablematchrobotrate"] = "1,0|2,40|3,33|4,27"
		section["stresstest"] = "1"
	}
}
