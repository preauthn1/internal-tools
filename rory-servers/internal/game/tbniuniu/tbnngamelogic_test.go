package tbnn

import (
	"os"
	"path"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/player"
)

func TestGameRound(t *testing.T) {
	glc := NewTbniuniuLogicByTest()
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	glc.OnGameStart(func() {
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}
	glc.GameStartFunc()
	users := fixtures.GetGameUser()
	//玩家入座
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//底注設置
	glc.computeUserBetTimes(users[1])
	glc.baseScoreSetting() //Step=basecsore
	glc.baseScoreConfirm()
	//	用戶組牌(gamelogic.GameUser, data []byte)
	//	GameStep為nil導致組牌未完成
	glc.onUserCardGroup(users[1], nil)

	//glc.cardGroupFinish()

	//牌型計算(niuniucard:compare)
	//compare(ret1, ret2) 比較牌型，返回玩家1 贏玩家2的倍數
	//牌型結果(niuniucardret:compare) 比較兩牌型結果(1 大于 0 等于 -1 小于)
	//遊戲結算(niuniugameret:getGameRet)(*netproto.TBNN_GameRet)
	//Logic:gameRet 遊戲結果
	//創建結果牌型>比較>贏家獲得倍率
}
func NewTbniuniuLogicByTest() *TBNNGameLogic {
	currentPath, _ := os.Getwd()
	configDir := path.Join(currentPath, "config")
	cfgs := new(config.Group)
	cfgs.LoadConfig(configDir, 0)
	cfg := cfgs.GetConfig(gamelogic.ConfigFileNameGameServer)
	serverInfo := cfgs.GetConfig(gamelogic.ConfigFileNameServerInfo)

	servercfg := gamelogic.GetGameServerConfig(serverInfo, cfg)
	rumanager := gamelogic.NewRoomUserManager(new(player.Manager))
	room := gamelogic.NewRoom(servercfg, rumanager)
	table := room.GetTableByNo(0)

	evtdispatcher := event.NewAsyncEventDispatcher(1)
	tfManager := event.NewTimerFuncManager(evtdispatcher)
	glc := new(TBNNGameLogic)

	glc.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(table, tfManager, cfgs)
	parseNNConfig(cfg)
	glc.RegisterAfterGameEndFunc(func() {
		//通知用户抢庄倍数
		for _, ur := range glc.AllGameUser() {
			glc.notifyUserBetTimes(ur)
		}
		//发牌动画
		glc.currentWatchAnimationName = glc.RegisterAnimation(animationName_GameStart, animationTimeOut_GameStart, glc.baseScoreSetting)
	})
	glc.RegisterBeforeGameStartFunc(func() {
		//清除数据
		glc.clearGameData()
	})
	return glc
}

func TestCal(t *testing.T) {
	taxscale := 5
	amount := int64(2)
	taxamount := int64(amount * int64(taxscale) / 100)
	println(taxamount)
}
