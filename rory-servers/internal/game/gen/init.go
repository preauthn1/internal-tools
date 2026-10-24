package gen

import (
	"sort"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/pachislot"
)

// RoomTableInfoCreator 房間資訊產生器
var RoomTableInfoCreator = func(glm *gamelogic.LogicManager) interface{} {
	data := new(netproto.GEN_RoomInfo)
	oddsMap := make(map[int64]bool)

	glm.Range(func(g gamelogic.GameLogic) {
		game := g.(*Game)
		oddsMap[game.Bet] = true

		info := new(netproto.GEN_TableInfo)
		info.ID = game.TableManager.TableNo
		info.Odds = game.Bet
		info.RoadMap = game.mDataRoadMap(game.PachiSlotGame.RoadMap)
		info.Occupied = len(game.AllGameUser()) != 0

		data.Tables = append(data.Tables, info)
	})

	for odds := range oddsMap {
		data.Odds = append(data.Odds, odds)
	}

	sort.Slice(data.Odds, func(i, j int) bool {
		return data.Odds[i] < data.Odds[j]
	})

	sort.Slice(data.Tables, func(i, j int) bool {
		return data.Tables[i].ID < data.Tables[j].ID
	})

	return data
}

// GameLogicCreator 遊戲邏輯產生器
var GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
	var logic gamelogic.GameLogic
	err := app.Context.Invoke(func(mdm *pachislot.MachineDistributionManager, proto *pachislot.Proto) {
		odds, pattern := mdm.Pick()
		logic = NewGame(t, tfManager, cfgs, proto, pattern, odds)
	})
	if err != nil {
		mlog.Fatal("%v", err)
	}
	return logic
}

func init() {
	app.Context.Provide(func() pachislot.PachiGameFunc {
		return func(g gamelogic.GameLogic) *pachislot.Game {
			game := g.(*Game)
			return game.PachiSlotGame
		}
	})
	msg.InstanceProvider.Register(int32(netproto.MessageBClassID_Game), int32(netproto.GEN_GameMessageClassID_GEN_BetID), func() proto.Message { return new(netproto.GEN_Bet) })
}
