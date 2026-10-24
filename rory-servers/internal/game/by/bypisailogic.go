package by

import (
	"fmt"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/mylog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

//比賽相關資訊
type PiSaiLogic struct {
	preGameStatus bool                             //準備遊戲前
	gameEndLeave  int64                            //比賽活動結束後幾秒踢出去
	ruleInfo      *netproto.GameCompetitionRuleMsg //比賽規則資訊
	//piSaiGameStartTimer *event.TimerFunc                 //遊戲開始整個時間管理棄
	//piSaiGameEndTimer *event.TimerFunc //遊戲結束时间管理器
	piSaiLeaveTimer *event.TimerFunc //遊戲離開时间管理器
}

//初始化比賽資訊
func initPiSai(by *Bylogic) {
	if by.TableManager.Room.GameType == gamelogic.PiSai {
		by.PiSaiLogic = new(PiSaiLogic)
		by.preGameStatus = false
	}
}

//取得比賽資訊
func (by *Bylogic) getContentInfo() {

	if by.TableManager.Room.GameType != gamelogic.PiSai {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HallMsgClassID_GameCompetioniInfoRuleID)
	pdata := new(netproto.GetGameCompetitionRuleMsg)
	pdata.GameID = proto.Int32(by.TableManager.Room.GameInfo.GameID)
	pmsg.MsgData = pdata
	by.SendRPCMessageFunc(true, pmsg)
}

func (by *Bylogic) setContentInfo(m *msg.Message) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GameCompetitionRuleMsg{}
		return mm
	})
	detail := rmm.(*netproto.GameCompetitionRuleMsg)

	by.ruleInfo = detail
	logmsg := fmt.Sprintf("設定比賽資訊:%+v ", detail)
	mlog.Debug(logmsg)
	mylog.MongoLog(by.userdata[0].user.UserID, *by.TableManager.Room.Seq, by.userdata[0].user.TableNo, logmsg)

	if _, ok := by.getPiSaiUserData(); ok {
		by.userdata[0].user.ExperienceCash = int64(by.ruleInfo.GetCoin())
		by.userdata[0].user.ContestID = by.ruleInfo.GetContestID()
		//mlog.System("取得比賽Rule[%+v] ", by.ruleInfo)
	}
}

//比賽活動消息
func (by *Bylogic) PiSaiMessage(gameuser *gamelogic.GameUser, gmsg *msg.Message) {
	defer util.PanicErrStack()

	if !by.CheckGameUserMessage(gameuser) {
		by.KickUserFunc(gameuser.UserID, "非法操作")
		return
	}

	//不是给游戏逻辑服务器的消息
	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	if gmsg.SClassID == int32(netproto.BY_GameMessageClassID_BYGameStartC2SID) {
		by.PiSaiStart(gameuser)
		return
	}

	p, ok := by.getPiSaiUserData()
	if !ok || !p.IsGameStart {
		mlog.Debug("return,  遊戲尚未開始 ")
		return
	}

	if gmsg.SClassID == int32(netproto.BY_GameMessageClassID_FrameC2SID) {
		by.userFrame(gameuser, gmsg.MsgData.([]byte))
		return
	}
}

//遊戲開始時間管理,完全不能停或刪掉, 當時間到則遊戲就結束
func (by *Bylogic) pisaiCancelTimer() {

	for _, v := range by.TableManager.Room.OnPiSaiUserInfo {
		if v.TableNo == by.TableNo() {
			v.IsGameStart = false
			mlog.System("before isplaying[%v] gamestart[%v] preGameStatus[%v] 當時間到則遊戲就結束 userid[%v]",
				v.IsPlaying, v.IsGameStart, by.preGameStatus, v.UserID)
			if v.IsPlaying && !v.IsGameStart && !by.preGameStatus {

				if _, ok := by.getPiSaiUserData(); !ok {
					delete(by.TableManager.Room.OnPiSaiUserInfo, v.UserID)
					by.KickUserFunc(v.UserID, "-")
				}
			}
		}
	}

}

//比賽結束
func (by *Bylogic) PiSaiGameEnd() {
	//if by.piSaiGameEndTimer == nil {
	//	return
	//}

	by.updateUserSorce()

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BY_GameMessageClassID_BYGameCompetitionEndID)
	pdata := new(netproto.BY_GameCompetitionEndMsg)

	p, ok := by.getPiSaiUserData()
	if ok {
		pdata.UserID = proto.Int32(by.userdata[0].user.UserID)
		pdata.Score = proto.Int32(int32(p.Bonus))
		pmsg.MsgData = pdata
		gameuser := new(gamelogic.GameUser)
		gameuser.UserID = by.userdata[0].user.UserID
		by.SendMessageToUser(gameuser, pmsg)

		p.EndTime = -15
		by.updateContestUser()
		logMsg := fmt.Sprintf("比賽結束 userid[%v] p.EndTime[%v] Score[%v] Bonus[%v] IsPlaying[%v]",
			by.userdata[0].user.UserID, p.EndTime, by.userdata[0].user.Score, p.Bonus, p.IsPlaying)
		mlog.Debug(logMsg)
		mylog.MongoLog(by.userdata[0].user.UserID, *by.TableManager.Room.Seq, by.userdata[0].user.TableNo, logMsg)

	}

	p.IsGameStart = false
	by.preGameStatus = false

	msglog := fmt.Sprintf("比賽結束timmer  踢人 gameEndLeave time[%v] isGameStart[%v]", by.gameEndLeave, p.IsGameStart)
	mlog.Debug(msglog)
	mylog.MongoLog(by.userdata[0].user.UserID, *by.TableManager.Room.Seq, by.userdata[0].user.TableNo, msglog)

	if by.userdata != nil && len(by.userdata) > 0 {
		mylog.MongoLog(by.userdata[0].user.UserID, *by.TableManager.Room.Seq, by.userdata[0].user.TableNo, msglog)
	}

	updatescop := util.NewIntScope(int32(by.gameEndLeave)-1, int32(by.gameEndLeave)+1)
	by.piSaiLeaveTimer = by.TimerFuncManager.NewNamedRandomTicker(updatescop, by.gameEndKickUser, "piSaiLeaveTimer")
}

func (by *Bylogic) gameEndKickUser() {
	if by.piSaiLeaveTimer == nil {
		return
	}

	if len(by.userdata) <= 0 {
		return
	}

	ur := by.userdata[0].user
	if p, ok := by.getPiSaiUserData(); ok {
		if p.IsGameStart {
			return
		}

		p.IsPlaying = false
		p.Pause = true
		logmsg := fmt.Sprintf("比賽結束刪除玩家 userid[%v] [%+v]", p.UserID, p)
		mylog.MongoLog(by.userdata[0].user.UserID, *by.TableManager.Room.Seq, by.userdata[0].user.TableNo, logmsg)
		delete(by.TableManager.Room.OnPiSaiUserInfo, by.userdata[0].user.UserID)
	}

	logMsg := fmt.Sprintf("gameEnd 踢出 userid[%v] ",
		by.userdata[0].user.UserID)
	mylog.MongoLog(by.userdata[0].user.UserID, *by.TableManager.Room.Seq, by.userdata[0].user.TableNo, logMsg)
	by.KickUserFunc(ur.UserID, "-")
	by.UserLeave(ur)
}

//比賽開始
func (by *Bylogic) PiSaiStart(gameuser *gamelogic.GameUser) {
	defer util.PanicErrStack()

	p, ok := by.getPiSaiUserData()
	if !ok {
		mlog.Error("沒有userdata無法開始比賽. user[%+v]", by.userdata[0])
	}
	if p.IsGameStart {
		mlog.Debug("return,  比賽進行中 ")
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.BY_GameMessageClassID_BYGameStartS2CID)
	by.SendMessageToUser(gameuser, pmsg)

	//by.totalTime = time.Now().UnixNano() + int64(by.ruleInfo.GetPlayTime()) * int64(time.Second)

	if p, ok := by.getPiSaiUserData(); ok {
		if p.IsPlaying {
			logmsg := fmt.Sprintf("比賽開始 isGameStart[%v] ", p.IsGameStart)
			mlog.Debug(logmsg)
			mylog.MongoLog(gameuser.UserID, *by.TableManager.Room.Seq, gameuser.TableNo, logmsg)

			p.Pause = false
			p.TableNo = by.userdata[0].user.TableNo

			by.userdata[0].user.Score = int64(by.ruleInfo.GetCoin())
			p.Score = by.userdata[0].user.Score
			p.Bonus = 0

			//by.ruleInfo.PlayTime = proto.Int32(10)
			p.EndTime = time.Now().UnixNano() + int64(by.ruleInfo.GetPlayTime())*int64(time.Second)
			totalPlayTime := (p.EndTime - time.Now().UnixNano()) / int64(time.Second)

			p.IsGameStart = true
			//by.piSaiGameEndTimer = by.TimerFuncManager.NewTimer(int32(totalPlayTime), by.PiSaiGameEnd)
			//by.piSaiGameStartTimer = by.TimerFuncManager.NewTimer(int32(totalPlayTime), by.pisaiCancelTimer)

			logMsg := fmt.Sprintf("userid[%v] isplaying[%v] 比賽開始 isGameStart[%v] 是否有票[%v] score[%v] bonus[%v] totalPlayTime[%v] userScore[%v]",
				p.UserID, p.IsPlaying, p.IsGameStart, p.IsPlaying, p.Score, p.Bonus, totalPlayTime, by.userdata[0].user.Score)
			mlog.Debug(logMsg)
			mylog.MongoLog(by.userdata[0].user.UserID, *by.TableManager.Room.Seq, by.userdata[0].user.TableNo, logMsg)

			by.preGameStatus = false
		} else if !p.IsPlaying || (p.EndTime-time.Now().UnixNano()) <= 0 {
			by.KickUserFunc(p.UserID, "-")
			by.UserLeave(by.userdata[0].user)
			return
		} else {
			by.userdata[0].user.Score = p.Score
			mlog.System("斷線重連比賽 score[%v] bonus[%v] userScore[%v]", p.Score, p.Bonus, by.userdata[0].user.Score)
		}

		if by.nextsendmsg != nil && by.nextsendmsg.MsgData != nil {
			nextMsg := by.nextsendmsg.MsgData.(*netproto.BY_FrameS2C)
			nextMsg.GameCompetitionTime = proto.Int32(by.ruleInfo.GetPlayTime())
			by.nextsendmsg.MsgData = nextMsg
		}
	}

}

//比賽玩家離開或斷線
func (by *Bylogic) PiSaiUserPause() bool {
	if p, ok := by.getPiSaiUserData(); ok {
		return p.Pause
	}
	return false
}

//設定比賽時間frmae
func (by *Bylogic) setPiSaiNextFramePlayTime() {
	if by.TableManager.Room.GameType != gamelogic.PiSai {
		return
	}

	if p, ok := by.getPiSaiUserData(); ok {
		if !p.IsGameStart {
			//mlog.Debug(" 設定比賽時間frmae stop")
			return
		}

		if !p.IsPlaying {
			logmsg := fmt.Sprintf("B 遊戲結束斷線重連進來的 userid[%v] isGameStart[%v] IsPlaying[%v]", p.UserID, p.IsGameStart, p.IsPlaying)
			mlog.System(logmsg)
			mylog.MongoLog(p.UserID, *by.TableManager.Room.Seq, p.TableNo, logmsg)
			by.gameEndKickUser()
			return
		}
	} else {
		return
	}

	lastTime := int64(time.Now().UnixNano())
	playTime := int64(0)

	if p, ok := by.getPiSaiUserData(); ok {
		if p.Pause {
			return
		}
		playTime = (p.EndTime - lastTime) / int64(time.Second)
		boomCount := by.userdata[0].boomcount
		ud := by.getuserdatabyuserid(p.UserID)
		if p.Score <= 0 && boomCount <= 0 && p.IsPlaying && ud.teshufishtype == 0 {
			logMsg := fmt.Sprintf(" gameEnd 設定比賽時間frmae userid[%v] playTime[%v] Score[%v] Bonus[%v] p.EndTime[%v] isplaying[%v]",
				by.userdata[0].user.UserID, playTime, by.userdata[0].user.Score, p.Bonus, p.EndTime, p.IsPlaying)
			mlog.Debug(logMsg)
			mylog.MongoLog(by.userdata[0].user.UserID, *by.TableManager.Room.Seq, by.userdata[0].user.TableNo, logMsg)

			by.PiSaiGameEnd()
		}

		if playTime <= 0 {
			mlog.Debug("playTime[%v]", playTime)
			by.PiSaiGameEnd()
		}

	}

	if playTime <= 0 {
		return
	}

	nextMsg := by.nextsendmsg.MsgData.(*netproto.BY_FrameS2C)
	nextMsg.GameCompetitionTime = proto.Int32(int32(playTime))
	by.nextsendmsg.MsgData = nextMsg

}

//更新比賽玩家積分
func (by *Bylogic) updatePiSaiUserScore(money int32) {
	if p, ok := by.getPiSaiUserData(); ok && p.Score > 0 {
		p.Score = p.Score + int64(money)
	}
}

//更新比賽玩家總獲得魚
func (by *Bylogic) updatePiSaiUserBonus(money int32) {

	if p, ok := by.getPiSaiUserData(); ok {
		p.Bonus = p.Bonus + int64(money)
		mlog.Debug("userid[%v] 玩家總獲得魚 [%v]", by.userdata[0].user.UserID, p.Bonus)
	}

}

//取得比賽玩家資訊
func (by *Bylogic) getPiSaiUserData() (*gamelogic.PiSaiUserInfo, bool) {
	if len(by.userdata) > 0 && by.TableManager.Room.OnPiSaiUserInfo != nil {
		if by.TableManager.Room.GameType == gamelogic.PiSai {
			if p, ok := by.TableManager.Room.OnPiSaiUserInfo[by.userdata[0].user.UserID]; ok {
				return p, ok
			}
		}
	}
	return nil, false
}

//取得玩家未完成遊戲資訊
func (by *Bylogic) getContestUser() {
	if by.TableManager.Room.GameType != gamelogic.PiSai {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DBServerClassID_GetContestUserDataID)
	pdata := new(netproto.GetConestUser)
	pdata.GameID = proto.Int32(by.TableManager.Room.GameInfo.GameID)
	pdata.UserID = proto.Int32(by.userdata[0].user.UserID)
	pmsg.MsgData = pdata
	by.SendRPCMessageFunc(true, pmsg)
}

//玩家離開時儲存遊戲資訊
func (by *Bylogic) updateContestUser() {
	if by.TableManager.Room.GameType != gamelogic.PiSai {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DBServerClassID_UpdateContestUserDataID)
	pdata := new(netproto.UpdateContestUser)
	pdata.GameID = proto.Int32(by.TableManager.Room.GameInfo.GameID)
	pdata.UserID = proto.Int32(by.userdata[0].user.UserID)
	pdata.ContestID = proto.Int32(by.ruleInfo.GetContestID())

	pdata.FlagID = proto.Int32(by.flag)
	pdata.TableNo = proto.Int32(by.TableNo())
	p, _ := by.getPiSaiUserData()
	totalPlayTime := (p.EndTime - time.Now().UnixNano()) / int64(time.Second)
	pdata.PlayTime = proto.Int32(int32(totalPlayTime))
	pdata.RemainScore = proto.Int32(int32(p.Bonus))
	pdata.Score = proto.Int32(int32(p.Score))
	pdata.ContestNo = proto.Int32(p.ContestNo)

	logmsg := fmt.Sprintf("更新用戶比賽資料:%+v", pdata)
	mlog.Debug(logmsg)
	mylog.MongoLog(by.userdata[0].user.UserID, *by.TableManager.Room.Seq, by.userdata[0].user.TableNo, logmsg)

	pmsg.MsgData = pdata
	by.SendRPCMessageFunc(true, pmsg)
}

//比賽设置消息
func (by *Bylogic) pisaiSettingMsg(pdata *netproto.BY_Setting) {
	pdata.Competition = new(netproto.BY_GameCompetitionSetting)
	pdata.Competition.Time = proto.Int32(by.ruleInfo.GetPlayTime())
	if by.ruleInfo.GetPlayTime() == 0 { //斷線重連setting還沒取得時間先預設99999
		pdata.Competition.Time = proto.Int32(99999)
	}
	if p, ok := by.getPiSaiUserData(); ok {
		logmsg := fmt.Sprintf("userid[%v] 比賽设置消息 IsPlaying[%v] isGameStart[%v] GetPlayTime[%v]", p.UserID, p.IsPlaying, p.IsGameStart, by.ruleInfo.GetPlayTime())
		mlog.Debug(logmsg)
		mylog.MongoLog(by.userdata[0].user.UserID, *by.TableManager.Room.Seq, by.userdata[0].user.TableNo, logmsg)

		if p.IsPlaying && p.IsGameStart {
			p.IsPlaying = true
			pdata.Competition.OnGamingRecover = proto.Int32(1)
			by.getContestUser()
		} else if p.IsPlaying {
			by.preGameStatus = true
			pdata.Competition.OnGamingRecover = proto.Int32(0)
			p.Bonus = 0

			totalPlayTime := (p.EndTime - time.Now().UnixNano()) / int64(time.Second)
			logMsg := fmt.Sprintf("新比賽 getsetting userid[%v] 是否有票[%v] preGameStatus[%v] isGameStart[%v] totalPlayTime[%v] Score[%v] Bonus[%v] ",
				by.userdata[0].user.UserID, p.IsPlaying, by.preGameStatus, p.IsGameStart, totalPlayTime, by.userdata[0].user.Score, p.Bonus)
			mlog.Debug(logMsg)
			mylog.MongoLog(by.userdata[0].user.UserID, *by.TableManager.Room.Seq, by.userdata[0].user.TableNo, logMsg)
		} else if !p.IsPlaying && !p.IsGameStart {
			mlog.Error("getSettingMsg 遊戲結束強制踢出")
			by.KickUserFunc(p.UserID, "-")
		}
	}
}

//繼續未完成的比賽
func (by *Bylogic) continuContestUser(m *msg.Message) {
	defer util.PanicErrStack()

	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UpdateContestUser{}
		return mm
	})
	ust := rmm.(*netproto.UpdateContestUser)

	if p, ok := by.getPiSaiUserData(); ok || (ust.GetPlayTime() > 0 && ust.GetPlayTime() < 99999) {
		endTime := time.Now().UnixNano() + int64(ust.GetPlayTime())*int64(time.Second)
		totalPlayTime := (endTime - time.Now().UnixNano()) / int64(time.Second)
		//mlog.Error("check totalPlayTime[%v] EndTime[%v] GetPlayTime[%v]", totalPlayTime, p.EndTime, ust.GetPlayTime())

		if by.TableManager.Room.OnPiSaiUserInfo == nil {

			by.TableManager.Room.OnPiSaiUserInfo = make(map[int32]*gamelogic.PiSaiUserInfo, 0)

			if ust.GetPlayTime() < 99999 {
				by.TableManager.Room.OnPiSaiUserInfo[ust.GetUserID()] = &gamelogic.PiSaiUserInfo{
					UserID:      ust.GetUserID(),
					IsPlaying:   true,
					IsGameStart: true,
					TableNo:     by.userdata[0].user.TableNo,
					EndTime:     endTime,
					Pause:       false,
					Score:       int64(ust.GetScore()),
					Bonus:       int64(ust.GetRemainScore()),
					ContestNo:   ust.GetContestNo(),
				}

				gameuser := new(gamelogic.GameUser)
				gameuser.UserID = ust.GetUserID()
				setmsg := by.getSettingMsg()
				by.SendMessageToUser(gameuser, setmsg)
			}
		} else {
			if p == nil {
				by.TableManager.Room.OnPiSaiUserInfo[ust.GetUserID()] = &gamelogic.PiSaiUserInfo{
					UserID:      ust.GetUserID(),
					IsPlaying:   true,
					IsGameStart: true,
					TableNo:     0,
					EndTime:     endTime,
					Pause:       false,
					Score:       int64(ust.GetScore()),
					Bonus:       int64(ust.GetRemainScore()),
					ContestNo:   ust.GetContestNo(),
				}
			} else {
				p.Score = int64(ust.GetScore())
				p.Bonus = int64(ust.GetRemainScore())
				p.IsPlaying = true
				p.Pause = false
				p.ContestNo = ust.GetContestNo()
				p.EndTime = endTime
				by.userdata[0].user.Score = p.Score
				logMsg := fmt.Sprintf("比賽進行中 userid[%v] 剩餘秒數[%v] 積分[%v] 得分[%v] 報名參成功[%v] ",
					ust.GetUserID(), totalPlayTime, by.userdata[0].user.Score, p.Bonus, p.IsPlaying)
				mlog.Debug(logMsg)
				mylog.MongoLog(by.userdata[0].user.UserID, *by.TableManager.Room.Seq, by.userdata[0].user.TableNo, logMsg)
			}
		}

		//by.piSaiGameEndTimer = by.TimerFuncManager.NewTimer(int32(totalPlayTime), by.PiSaiGameEnd)

		if totalPlayTime <= -12 {
			by.gameEndKickUser()
		}

	} else {
		playtime := by.ruleInfo.GetPlayTime()
		if ust.GetPlayTime() != 99999 {
			playtime = ust.GetPlayTime()
		}

		endTime := time.Now().UnixNano() + int64(playtime)*int64(time.Second)
		totalPlayTime := (endTime - time.Now().UnixNano()) / int64(time.Second)
		isplaying := false
		if totalPlayTime > 0 {
			isplaying = true
		} else {
			logmsg := fmt.Sprintf("時間到, 剩餘秒數[%v], 踢出userid[%v]", totalPlayTime, by.userdata[0].user.UserID)
			mylog.MongoLog(by.userdata[0].user.UserID, *by.TableManager.Room.Seq, by.userdata[0].user.TableNo, logmsg)
			by.KickUserFunc(by.userdata[0].user.UserID, "-")
			by.UserLeave(by.userdata[0].user)
			return
		}

		if by.TableManager.Room.OnPiSaiUserInfo == nil {
			by.TableManager.Room.OnPiSaiUserInfo = make(map[int32]*gamelogic.PiSaiUserInfo, 0)
		}

		by.TableManager.Room.OnPiSaiUserInfo[ust.GetUserID()] = &gamelogic.PiSaiUserInfo{
			UserID:    ust.GetUserID(),
			IsPlaying: isplaying,
			TableNo:   0,
			EndTime:   endTime,
			Pause:     false,
			Score:     0,
			Bonus:     0,
			ContestNo: ust.GetContestNo(),
		}

		logmsg := fmt.Sprintf("建立比賽玩家資料 剩餘秒數[%v] 買票[%v] [%+v]",
			totalPlayTime, isplaying, by.TableManager.Room.OnPiSaiUserInfo[ust.GetUserID()])
		mylog.MongoLog(ust.GetUserID(), *by.TableManager.Room.Seq, 0, logmsg)
	}

}

//檢查玩家是否有買票
func (by *Bylogic) checkUserTicket() {
	if by.TableManager.Room.GameType != gamelogic.PiSai {
		return
	}

	by.getContestUser()

}

//是否破產
func (by *Bylogic) piSaiCheckShoot(ispochan bool, boomMoney int32) bool {
	if p, ok := by.getPiSaiUserData(); ok {
		if p.Score-int64(boomMoney) < 0 {
			return true
		} else {
			return false
		}
	}

	return ispochan
}
