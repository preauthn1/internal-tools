package ddz

import (
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

func TestR(t *testing.T) {
	for i := 0; i < 1; i++ {
		mlog.Debug("局數[%v]", i)
		TestAIRobot2(t)
	}
}

var userWin int

//2機器人對1真人, 機器人贏
func TestAIRobot2(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, _, users := comUserReady(endChannel, "ai2Poker", t)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}
	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId && score > 0 {
			userWin++
			mlog.System("真人贏次數[%v] score[%v]", userWin, score)
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//AI對子
func TestAIOutDouble(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, ddz, users := comUserReady(endChannel, "aioutCardByDouble", t)

	initStart := true

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if userid == users[1].UserID {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameStep) {
				glc.ext.lastQiangDiZhu = users[1].SeatNo
				JiaoDiZhuMsg(m, users, glc)
				gameProgressStartMsg(m, initStart, glc, ddz)
			}

		}
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)

				u := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u.userId {
					mlog.Warn("Test userId[%v] 非法操作，不该您出牌 u.isTuoGuan[%v] seatno %v ", userid, u.isTuoGuan, u.userSeat)
					return
				}

				rbglc.OnGameMessage(r, m)
			}

			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId &&
			score != int64(-10)*int64(glc.ext.odds)*2 && taxamount != int64(1)*int64(glc.ext.odds) {
			t.Errorf("userid[%v] 更新金幣錯誤, 期望-40, 實際[%v], tax[%v] , odd[%v]",
				userid, score, taxamount, glc.ext.odds)
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//AI 機器人最後一手,不能K葫蘆拆掉分兩手出
func TestThreeAndTwo(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, ddz, users := comUserReady(endChannel, "aithreeandtwo", t)

	initStart := true

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if userid == users[1].UserID {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameStep) {
				glc.ext.lastQiangDiZhu = users[1].SeatNo
				JiaoDiZhuMsg(m, users, glc)
				gameProgressStartMsg(m, initStart, glc, ddz)
			}

		}
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)

				u2 := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u2.userId {
					return
				}

				rbglc.OnGameMessage(r, m)
			}

			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId &&
			score != int64(-10)*int64(glc.ext.odds)*2 && taxamount != int64(1)*int64(glc.ext.odds) {
			t.Errorf("userid[%v] 更新金幣錯誤, 期望-40, 實際[%v], tax[%v] , odd[%v]",
				userid, score, taxamount, glc.ext.odds)
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//機器人是真人地主, 機器人農民下家策略
func TestAI01(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, ddz, users := comUserReady(endChannel, "ai01", t)

	initStart := true

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if userid == users[1].UserID {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameStep) {
				glc.ext.lastQiangDiZhu = users[1].SeatNo
				JiaoDiZhuMsg(m, users, glc)
				gameProgressStartMsg(m, initStart, glc, ddz)
			}

		}
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)

				u2 := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u2.userId {
					//mlog.Debug("Test userId[%v] 非法操作，不该您出牌 u.isTuoGuan[%v] seatno %v ", userid, u.isTuoGuan, u.userSeat)
					return
				}

				rbglc.OnGameMessage(r, m)
			}

			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId && score > 0 {
			t.Errorf("真人當地主, 機器人農民要勝利, 但卻是真人勝利.")
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//真人當地主, 機器人農民要勝利,
func TestAI02(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, ddz, users := comUserReady(endChannel, "ai02", t)

	initStart := true

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if userid == users[1].UserID {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameStep) {
				glc.ext.lastQiangDiZhu = users[1].SeatNo
				JiaoDiZhuMsg(m, users, glc)
				gameProgressStartMsg(m, initStart, glc, ddz)
			}

		}
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)

				u2 := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u2.userId {
					return
				}

				rbglc.OnGameMessage(r, m)
			}

			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId && score > 0 {
			t.Errorf("真人當地主, 機器人農民要勝利, 但卻是真人勝利.")
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//真人當地主, 機器人農民要勝利,
func TestAI03(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, ddz, users := comUserReady(endChannel, "ai03", t)

	initStart := true

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if userid == users[1].UserID {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameStep) {
				glc.ext.lastQiangDiZhu = users[1].SeatNo
				JiaoDiZhuMsg(m, users, glc)
				gameProgressStartMsg(m, initStart, glc, ddz)
			}

		}
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)

				u2 := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u2.userId {
					return
				}

				rbglc.OnGameMessage(r, m)
			}

			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId && score > 0 {
			t.Logf("真人當地主, 機器人農民要勝利, 但卻是真人勝利.")
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//TODO 未完成 真人當地主, 機器人農民要勝利,
func BTestAI04(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, ddz, users := comUserReady(endChannel, "ai04", t)

	initStart := true

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if userid == users[1].UserID {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameStep) {
				glc.ext.lastQiangDiZhu = users[1].SeatNo
				JiaoDiZhuMsg(m, users, glc)
				gameProgressStartMsg(m, initStart, glc, ddz)
			}

		}
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)

				u2 := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u2.userId {
					return
				}

				rbglc.OnGameMessage(r, m)
			}

			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId && score > 0 {
			t.Errorf("真人當地主, 機器人農民要勝利, 但卻是真人勝利.")
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//TODO 未完成 農民換牌
func BTestAI05(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, ddz, users := comUserReady(endChannel, "ai05", t)

	initStart := true

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if userid == users[1].UserID {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameStep) {
				glc.ext.lastQiangDiZhu = users[1].SeatNo
				JiaoDiZhuMsg(m, users, glc)
				gameProgressStartMsg(m, initStart, glc, ddz)
			}

		}
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)

				u2 := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u2.userId {
					return
				}

				rbglc.OnGameMessage(r, m)
			}

			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId && score > 0 {
			t.Errorf("真人當地主, 機器人農民要勝利, 但卻是真人勝利.")
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//農民機器人做球給下家農民
func TestAI06(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, ddz, users := comUserReady(endChannel, "ai06", t)

	initStart := true

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if userid == users[1].UserID {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameStep) {
				glc.ext.lastQiangDiZhu = users[1].SeatNo
				JiaoDiZhuMsg(m, users, glc)
				gameProgressStartMsg(m, initStart, glc, ddz)
			}

		}
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)

				u2 := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u2.userId {
					return
				}

				rbglc.OnGameMessage(r, m)
			}

			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId && score > 0 {
			t.Errorf("真人當地主, 機器人農民要勝利, 但卻是真人勝利.")
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//AI 農民做球
func TestAI07(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, ddz, users := comUserReady(endChannel, "ai07", t)

	initStart := true

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if userid == users[1].UserID {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameStep) {
				glc.ext.lastQiangDiZhu = users[1].SeatNo
				JiaoDiZhuMsg(m, users, glc)
				gameProgressStartMsg(m, initStart, glc, ddz)
			}

		}
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)

				u2 := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u2.userId {
					return
				}

				rbglc.OnGameMessage(r, m)
			}

			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId && score > 0 {
			t.Errorf("真人當地主, 機器人農民要勝利, 但卻是真人勝利.")
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//手上沒有最大炸彈, 農民有兩個炸彈
func TestAI08(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, ddz, users := comUserReady(endChannel, "ai08", t)

	initStart := true

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if userid == users[1].UserID {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameStep) {
				glc.ext.lastQiangDiZhu = users[1].SeatNo
				JiaoDiZhuMsg(m, users, glc)
				gameProgressStartMsg(m, initStart, glc, ddz)
			}

		}
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)

				u2 := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u2.userId {
					return
				}

				rbglc.OnGameMessage(r, m)
			}

			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId && score > 0 {
			t.Errorf("真人當地主, 機器人農民要勝利, 但卻是真人勝利.")
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//最後一手葫蘆卻拆牌出,最後一手牌是3帶2
func TestAI09(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, ddz, users := comUserReady(endChannel, "ai09", t)

	initStart := true

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if userid == users[1].UserID {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameStep) {
				glc.ext.lastQiangDiZhu = users[1].SeatNo
				JiaoDiZhuMsg(m, users, glc)
				gameProgressStartMsg(m, initStart, glc, ddz)
			}

		}
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)

				u2 := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u2.userId {
					return
				}

				rbglc.OnGameMessage(r, m)
			}

			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId && score > 0 {
			t.Errorf("真人當地主, 機器人農民要勝利, 但卻是真人勝利.")
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//農民有最大炸彈可以做球贏, 農民要做球再跟牌的時候要先壓
func TestAI10(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, ddz, users := comUserReady(endChannel, "ai10", t)

	initStart := true

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if userid == users[1].UserID {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameStep) {
				glc.ext.lastQiangDiZhu = users[1].SeatNo
				JiaoDiZhuMsg(m, users, glc)
				gameProgressStartMsg(m, initStart, glc, ddz)
			}

		}
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)

				u2 := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u2.userId {
					return
				}

				rbglc.OnGameMessage(r, m)
			}

			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId && score > 0 {
			t.Errorf("真人當地主, 機器人農民要勝利, 但卻是真人勝利.")
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//真人地主剩一手牌，機器人有炸彈不能出. 真人要勝利
func TestAI11(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, ddz, users := comUserReady(endChannel, "ai11", t)

	initStart := true

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if userid == users[1].UserID {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameStep) {
				glc.ext.lastQiangDiZhu = users[1].SeatNo
				JiaoDiZhuMsg(m, users, glc)
				gameProgressStartMsg(m, initStart, glc, ddz)
			}

		}
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)

				u2 := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u2.userId {
					return
				}

				rbglc.OnGameMessage(r, m)
			}

			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId && score < 0 {
			t.Logf("真人當地主, 真人要勝利, 但卻是真人輸.")
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//真人剩單張, 機器人要從頂天的牌出
func TestAI12(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, ddz, users := comUserReady(endChannel, "ai12", t)

	initStart := true

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if userid == users[1].UserID {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameStep) {
				glc.ext.lastQiangDiZhu = users[1].SeatNo
				JiaoDiZhuMsg(m, users, glc)
				gameProgressStartMsg(m, initStart, glc, ddz)
			}

		}
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)

				u2 := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u2.userId {
					return
				}

				rbglc.OnGameMessage(r, m)
			}

			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId && score < 0 {
			t.Logf("真人當地主, 真人要勝利, 但卻是真人輸.")
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//机器人顺子
func TestAI13(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		robotCfg.Enable = false //手動產生測試機器人
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, rbglc, ddz, users := comUserReady(endChannel, "ai13", t)

	initStart := true

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if userid == users[1].UserID {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameStep) {
				glc.ext.lastQiangDiZhu = users[1].SeatNo
				JiaoDiZhuMsg(m, users, glc)
				gameProgressStartMsg(m, initStart, glc, ddz)
			}

		}
		for _, u := range glc.AllGameUser() {
			if u.UserID == userid && u.IsRobot() {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)

				u2 := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u2.userId {
					return
				}

				rbglc.OnGameMessage(r, m)
			}

			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) {
				r := tgrs.RobotService.GetRobotManager().GetRobot(u.UserID)
				rbglc.OnGameMessage(r, m)
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
		_, dz := glc.ext.getDiZhu()

		if users[1].UserID == dz.userId && userid == dz.userId && score < 0 {
			t.Logf("真人當地主, 真人要勝利, 但卻是真人輸.")
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("ddz game end")
	}
}

//TODO 真人地主只剩一手牌型, 機器人主動出牌不出相同牌型.

//TODO AI 發三張牌模擬三組牌情境, 最好的牌給機器人

//TODO AI 機器人是真人地主上家策略

//

//共用玩家準備
func comUserReady(endChannel chan bool, filaName string, t *testing.T) (*fixtures.TestGameRoomServer, *DDZGameLogic, *DDZRobotLogic, *DDZGameLogic, []*gamelogic.GameUser) {
	tgrs, glc, rbglc, endChannel, robots := IntegrationTestInit(endChannel, t)
	glc.maxJiaoDiZhuRound = 1
	glc.jiaoDizhuUserRate = 100
	ddz := NewDDZGameLogicByTest(filaName)
	users := fixtures.GetGameUser()
	users[1].UserType = gamelogic.UserTypeNormal
	robots[2].UserType = gamelogic.UserTypeRebot
	robots[3].UserType = gamelogic.UserTypeRebot
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	robots[2].UserState.OnUserQueue(robots[2], tgrs.Room)
	robots[3].UserState.OnUserQueue(robots[3], tgrs.Room)
	//用戶入座
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(robots[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(robots[3], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserReady(users[1])
	glc.UserReady(robots[2])
	glc.UserReady(robots[3])
	return tgrs, glc, rbglc, ddz, users
}

//真人叫地主
func JiaoDiZhuMsg(m *msg.Message, users []*gamelogic.GameUser, glc *DDZGameLogic) {
	gdata := m.MsgData.(*netproto.DDZ_NotifyGameStep)
	u := glc.ext.getExtInfoByTurn(glc.ext.getNowSeat())
	if gdata.GetStep() == gameProgressJiaoDiZhu && u.userId == users[1].UserID {
		//gmsg *msg.Message
		gmsg := new(msg.Message)
		gmsg.SClassID = int32(netproto.MessageBClassID_Game)
		gmsg.BClassID = int32(netproto.DDZ_GameMessageClassID_JiaoDiZhu)
		mdata := new(netproto.DDZ_JiaoDiZhu)
		mdata.IsJiaoDiZhu = proto.Bool(true)
		gmsg.MsgData, _ = proto.Marshal(mdata)
		glc.GameMessage(users[1], gmsg)
	}
}

//遊戲開始替換牌
func gameProgressStartMsg(m *msg.Message, initStart bool, glc *DDZGameLogic, ddz *DDZGameLogic) {
	gdata := m.MsgData.(*netproto.DDZ_NotifyGameStep)
	if initStart && gdata.GetStep() == gameProgressStart {
		initStart = false
		glc.poker = ddz.poker

		for k, v := range glc.gamePoker.UserPokers {
			v.Cards = ddz.gamePoker.UserPokers[k].Cards
		}

	}
}

//outCardByType 出牌牌型 對子
func NoTestAIOutCardByTypeDouble2(t *testing.T) {
	ddz := NewDDZGameLogicByTest("outCardByDouble")

	//先出單牌 user=0
	autoCg, autoCardType, autoCardNumber, autoMainCard := ddz.onAutoChuPaiByTuoGuan()
	ckType := hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	//出對子 user=0
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeDouble || ckType != cardTypeDouble || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 對子, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	ddz.ext.turn++ //user=1

	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeDouble || ckType != cardTypeDouble || autoMainCard != int32(1) || autoCardNumber != int32(6) {
		t.Errorf("出牌牌型 對子, 跟牌 壓牌(2) 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	//user=1
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeDouble || ckType != cardTypeDouble || autoMainCard != int32(1) || autoCardNumber != int32(2) {
		t.Errorf("出牌牌型 對子, 主動 手上牌<=2種對子就從特殊牌最小開始出牌(2) 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	//user=1
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeDouble || ckType != cardTypeDouble || autoMainCard != int32(1) || autoCardNumber != int32(3) {
		t.Errorf("出牌牌型 對子, 跟牌 壓牌(2) 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	ddz = NewDDZGameLogicByTest("outCardByDouble")
	//先出單牌 user=0
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	//出對子 user=0
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	ddz.ext.turn++ //user=1
	ddz.ext.turn++ //user=2
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeDouble || ckType != cardTypeDouble || autoMainCard != int32(1) || autoCardNumber != int32(2) {
		t.Errorf("出牌牌型 對子, 跟牌 壓牌(2) 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

}
