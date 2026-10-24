//13水
package main

import (
	"example.invalid/ServerTeam/Servers/internal/game/sss"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

//牌型组合，可以全局只执行一次
/*
func CombInit() {
	for i := int(sssgamelogic.FCardTypeMin) + 1; i < int(sssgamelogic.FCardTypeMax); i++ {
		for j := int(sssgamelogic.NCardTypeMin) + int(i); j < int(sssgamelogic.NCardTypeMax); j++ {
			for k := j; k < int(sssgamelogic.NCardTypeMax); k++ {
				//ct := new(sssgamelogic.CardTemp)
				//ct.ct1 = sssgamelogic.FCardType(i)
				//ct.ct2 = sssgamelogic.NCardType(j)
				//ct.ct3 = sssgamelogic.NCardType(k)
				fmt.Printf("%d,%d,%d\r\n", i, j, k)
				//s.globalCardComb.PushBack(make([]))
			}
		}
	}
}*/

func main() {
	//CombInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return sss.NewSSSGameLogic(t, tfManager, cfgs)
	}

	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return sss.NewSSSRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	roomserver.StartGameLogicServer()
}
