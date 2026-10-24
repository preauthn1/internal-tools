package roomlogic

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/gamequeue"
	"example.invalid/ServerTeam/Servers/pkg/gameserver"
	"example.invalid/ServerTeam/Servers/pkg/gameserver/serverhandler"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/module/logicmodule"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/player"
	robotsvr "example.invalid/ServerTeam/Servers/pkg/robot"
	robot "example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/roomutil"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/scorestorage"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// RoomTableInfoCreator 房間牌桌資訊處理器
var RoomTableInfoCreator func(*gamelogic.LogicManager) interface{}

// RoomActionManager 房間動作管理器
type RoomActionManager struct {
	roomUserManager *gamelogic.RoomUserManager       // 房間用戶管理器
	room            *gamelogic.Room                  // 房間
	logics          *gamelogic.LogicManager          // 房間邏輯
	robotService    *robotsvr.Service                // 機器人服務
	serverCtx       *serverhandler.GameServerCtx     // 服務器環境
	notifyCtx       *serverhandler.NotifyHandler     // Notify服務環境
	scoreStorage    scorestorage.UserScoreStorage    // 積分存儲
	gameQueue       gamequeue.GameQueue              // 遊戲撮合器
	serverCfg       *gamelogic.GameServerConfig      // 遊戲服務器配置
	serverStatus    logicmodule.ServerStatusProvider // 服務器狀態

	OnGameMessage func(u *gamelogic.GameUser, m *msg.Message)
}

// NewRoomActionManager 創建房間動作管理器對象
func NewRoomActionManager(roomusermanager *gamelogic.RoomUserManager, room *gamelogic.Room, logics *gamelogic.LogicManager, robotService *robotsvr.Service, gs *gameserver.GameServer, scoreStorage scorestorage.UserScoreStorage, gmq gamequeue.GameQueue, cfg *gamelogic.GameServerConfig, serverStatus logicmodule.ServerStatusProvider) *RoomActionManager {
	ram := new(RoomActionManager)
	ram.roomUserManager = roomusermanager
	ram.room = room
	ram.logics = logics
	ram.robotService = robotService
	ram.serverCtx = gs.GameServerRunCtx
	ram.notifyCtx = gs.NotifyHandler
	ram.scoreStorage = scoreStorage
	ram.serverCfg = cfg
	ram.serverStatus = serverStatus

	if gmq != nil {
		ram.gameQueue = gmq
		ram.gameQueue.Start()
	}

	ram.logics.Range(func(g gamelogic.GameLogic) {
		// 踢人
		g.OnKickUser(func(userid int32, msg string) {
			user := ram.getGameUserByID(userid)
			if user == nil {
				return
			}

			if user.TableNo == g.TableNo() {
				if user.IsRobot() {
					ram.robotService.KickUser(user)
				} else {
					ram.kickUser(user.ConnInfo.ConnID(), msg)
				}
			}
		})

		// 遊戲取消
		g.OnGameCancel(func(str string) {
			if str == "" {
				str = "游戏异常解散"
			}

			tno := g.TableNo()

			urs := ram.getGameUserByTableNo(tno)

			if urs != nil {
				for _, ur := range urs {
					if ur.IsRobot() {
						continue
					}

					ram.sendTipMessage(ur.ConnInfo.ConnID(), str, roomutil.MsgTypeIDRoomKick)

					// 如果是斷線狀態的用戶
					if ur.IsOffLine() {
						ram.OnUserLeave(ur.ConnInfo.ConnID())
					}

					// 通知用戶狀態改變
					if !ur.IsOffLine() {
						ram.gameQueue.UserGameEnd(ur)

						// 通知用戶狀態改變
						ram.notifyUserStateChange(ur)
					}

					// 檢查用戶金幣數
					ram.checkUserMoney(ur)
				}
			}
		})

		// 遊戲結束
		g.OnGameEnd(func(gametime int32) {
			tno := g.TableNo()

			//tbl := ram.room.GetTableByNo(tno)
			//users := tbl.GetTableUser()
			//seatsHasUserCount := 0
			//s := ""
			//for seatno, seat := range tbl.Seats {
			//	if seat.SeatUser != nil {
			//		exists := false
			//		seatsHasUserCount++
			//		for _, robot := range ram.robotService.GetRobotManager().GetAllRobot() {
			//			if robot.UserInfo.UserID == seat.SeatUser.UserID {
			//				exists = true
			//				break
			//			}
			//		}
			//		if !exists {
			//			s += fmt.Sprintf("SeatNo %d, UserID %d,", seatno, seat.SeatUser.UserID)
			//		}
			//	}
			//}
			//mlog.System(fmt.Sprintf("userCount %d, seatsHasUserCount %d, not exsists userList %s", len(users), seatsHasUserCount, s))
			//if len(s) > 0 {
			//	mlog.Error("not exsists userList", s)
			//}

			urs := ram.getGameUserByTableNo(tno)

			if urs != nil {
				for _, ur := range urs {
					ur.GameTime += gametime

					if ur.IsRobot() {
						ram.robotService.GameEnd(ur)

						continue
					}

					// 如果是斷線狀態的用戶
					if ur.IsOffLine() {
						ram.OnUserLeave(ur.ConnInfo.ConnID())
					}

					// 通知用戶狀態改變
					if !ur.IsOffLine() {
						// 通知用戶狀態改變
						ram.notifyUserStateChange(ur)

						ram.gameQueue.UserGameEnd(ur)
					}

					// 檢查用戶金幣數
					ram.checkUserMoney(ur)
				}
			}
		})

		// 遊戲開始
		g.OnGameStart(func() {
			tno := g.TableNo()

			urs := ram.getGameUserByTableNo(tno)
			if urs != nil {
				for _, ur := range urs {
					// 檢查遊戲狀態
					ram.checkServerStatus(ur)

					if ur.IsRobot() {
						ram.robotService.GameStart(ur)

						continue
					}
					ram.scoreStorage.OnUserReady(ur) // 遊戲開始時再確認一次
					ram.notifyUserStateChange(ur)
				}
			}
		})

		// 寫分
		g.OnUpdateUserScore(func(userid int32, score int64, taxamount int64) {
			ur := ram.getGameUserByID(userid)
			if ur == nil {
				mlog.Error("用戶 %d 寫分 %d, %d, 但未找到用戶信息", userid, score, taxamount)
				return
			}

			if score > 0 {
				ur.WinCount++
				ur.UpdateWinCount++
				ur.LastWinLose = 1
			} else if score == 0 {
				ur.DrawCount++
				ur.UpdateDrawCount++
				ur.LastWinLose = 0
			} else {
				ur.LoseCount++
				ur.UpdateLoseCount++
				ur.LastWinLose = 0
			}

			ram.scoreStorage.UpdateUserScore(ur, score, taxamount)
		})

		// 體驗金幣寫分
		g.OnUpdateExperienceScore(func(userid int32, score int64, taxamount int64) {
			ur := ram.getGameUserByID(userid)
			if ur == nil {
				mlog.Error("用户%d體驗金幣寫分%d,%d,但未找到用户信息", userid, score, taxamount)
				return
			}

			ram.scoreStorage.UpdateExperienceScore(ur, score, taxamount)
		})

		// 遊戲獎金池寫分
		g.OnUpdateGameJackpotPool(func(userid int32, amount int64) {
			ur := ram.getGameUserByID(userid)
			if ur == nil {
				mlog.Error("用户%d遊戲獎金池寫分%d,但未找到用户信息", userid, amount)
				return
			}

			ram.scoreStorage.UpdateGameJackpotPool(ur, amount)
		})

		// 發送消息
		g.OnSendMessage(func(userid int32, m *msg.Message) {
			ur := ram.getGameUserByID(userid)
			if ur == nil {
				return
			}

			if ur.TableNo != g.TableNo() {
				return
			}

			if m == nil {
				return
			}

			// 機器人
			if ur.IsRobot() {
				ram.robotService.SendGameMessage(ur, m)
			} else { // 真實用戶
				ram.serverCtx.MsgServer.SendMessage(ur.ConnInfo, m)
			}
		})

		// 捕魚公告
		g.OnNotifyFishMessage(func(m *msg.Message) {
			ram.notifyCtx.TCPMsgClient.SendMessage(m)
		})

		// 發送RPC消息
		g.OnSendRPCMessage(func(safe bool, m *msg.Message) {
			rpcinfo := netproto.CreateEmptyRPCInfo(ram.serverCfg.ServerID)
			// 設置RPC消息桌號
			rpcinfo.ExtendInfo = proto.String(strconv.Itoa(int(g.TableNo())))
			ram.serverCtx.DBMsgQueue.PushMessage(rpcinfo, safe, m)
		})

		{ // 牌局記錄器
			var glog gamelog.GameLog
			if ram.serverCfg.DealLogType == gamelogic.DealLogTypeNone {
				glog = gamelog.NewEmptyGameLog()
			} else if ram.serverCfg.DealLogType == gamelogic.DealLogTypeFile {
				glog = gamelog.NewFileGameLog("deallog", ram.serverCfg.ServerID)
			} else if ram.serverCfg.DealLogType == gamelogic.DealLogTypeDB {
				glog = gamelog.NewDBGameLog(ram.serverCtx.DBMsgQueue, ram.serverCfg.GameID, ram.serverCfg.ServerID, g.TableNo())
			} else if ram.serverCfg.DealLogType == gamelogic.DealLogTypeFish {
				glog = gamelog.NewFishGameLog(ram.serverCtx.DBMsgQueue, ram.serverCfg.GameID, ram.serverCfg.ServerID, g.TableNo())
			} else if ram.serverCfg.DealLogType == gamelogic.DealLogTypeConditionalDB {
				glog = gamelog.NewConditonalDBGameLog(ram.serverCtx.DBMsgQueue, ram.serverCfg.GameID, ram.serverCfg.ServerID, g.TableNo(), map[string]interface{}{
					"SettleRound": ram.serverCfg.SettleRound,
					"SettleScore": ram.serverCfg.SettleScore,
				})
			} else if ram.serverCfg.DealLogType == gamelogic.DealLogTypeConsole {
				glog = gamelog.NewConsoleGameLog(ram.serverCfg.GameID, ram.serverCfg.ServerID, g.TableNo())
			}
			g.SetGameLog(glog)

			app.Context.Provide(func() gamelog.GameLog { return glog })
		}

		// 提前換桌
		g.OnUserEarlyChangeTable(func(ur *gamelogic.GameUser, early string) {
			ram.gameQueue.UserGameEnd(ur)
			g.UserLeave(ur)

			// 通知用戶狀態改變
			if early == "change" {
				ram.notifyUserStateChange(ur)
				if !ur.IsRobot() {
					ram.onUserQueue(ur.ConnInfo, ur)
				}
			} else if early == "leave" {
				ram.OnUserLeave(ur.ConnInfo.ConnID())
				//ram.room.RoomUserManager.RemovePlayer(ur.ConnInfo.ConnID())
			}

		})

		// 異步調用
		g.OnAsyncCall(func(asyncCallFunc func(), callBackFunc func()) {
			go func() {
				defer util.PanicErrStack()

				asyncCallFunc()

				// 發布異步調用事件
				asyncCallEvt := msgserver.NewAsyncCallEvent(callBackFunc)
				ram.serverCtx.MsgServer.EventChannel.OnAsyncEvent(asyncCallEvt, nil)
			}()
		})
	})

	return ram
}

// SetRobotService 設定機器人服務
func (ram *RoomActionManager) SetRobotService(robotService *robotsvr.Service) {
	ram.robotService = robotService
}

// SetGameQueue 設置遊戲撮合邏輯
func (ram *RoomActionManager) SetGameQueue(gameQueue gamequeue.GameQueue) {
	if ram.gameQueue != nil {
		ram.gameQueue.Stop()
	}
	ram.gameQueue = gameQueue
	ram.gameQueue.Start()
}

// IsHook 是否攔截消息
func (ram *RoomActionManager) IsHook(bclsid int32, sclsid int32) bool {
	return (bclsid == int32(netproto.MessageBClassID_GameRoom) || bclsid == int32(netproto.MessageBClassID_Game))
}

// OnMessage 響應遊戲消息
func (ram *RoomActionManager) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	defer util.PanicErrStack()
	mlog.Trace("client %v msg %v", clt, m)

	// 是否買票
	if m.BClassID == int32(netproto.MessageBClassID_GameRoom) && m.SClassID == int32(netproto.HallMsgClassID_GameCompetitionID) {
		// 檢查連接是否已經登錄
		pr := ram.roomUserManager.GetPlayerByConnID(clt.ConnID())
		if pr != nil && !pr.IsGuest() {
			mlog.Error("客户端%v發送登錄數據包,但客戶端已經登錄了賬號%d", clt, pr.UserID)
			//ram.kickUser(clt.ConnID(), "非法操作")
			return nil
		}

		dd := m.MsgData.([]byte)
		mm := &netproto.GetGameCompetitionMsg{}

		err := proto.Unmarshal(dd, mm)
		if err != nil {
			mlog.Error("解析是否買票數據包錯誤: %v,client:%v", err, clt)
			ram.kickUser(clt.ConnID(), "非法操作")
			return nil
		}

		rpcinfo := netproto.CreateEmptyRPCInfo(ram.room.ServerID) //setLoginRPC(clt, mm, ram)
		rpcinfo.ConnID = proto.String(clt.ConnID())
		rpcinfo.ConnToken = [REDACTED]
		ram.serverCtx.DBMsgQueue.PushMessage(rpcinfo, false, m)
	}

	// 用戶登錄
	if m.BClassID == int32(netproto.MessageBClassID_GameRoom) && m.SClassID == int32(netproto.GameRoomClassID_LoginRoomID) {
		// 檢查連接是否已經登錄
		pr := ram.roomUserManager.GetPlayerByConnID(clt.ConnID())
		if pr != nil && !pr.IsGuest() {
			mlog.Error("客户端%v發送登錄數據包,但客戶端已經登錄了賬號%d", clt, pr.UserID)
			//ram.kickUser(clt.ConnID(), "非法操作")
			return nil
		}

		dd := m.MsgData.([]byte)
		mm := &netproto.LoginGameRoomInfo{}

		err := proto.Unmarshal(dd, mm)
		if err != nil {
			mlog.Error("解析登錄數據包錯誤: %v,client:%v", err, clt)
			ram.kickUser(clt.ConnID(), "非法操作")
			return nil
		}

		rpcinfo := setLoginRPC(clt, mm, ram)
		ram.onUserLogin(rpcinfo, mm, m)
	}

	// 用戶排隊
	if m.BClassID == int32(netproto.MessageBClassID_GameRoom) && m.SClassID == int32(netproto.GameRoomClassID_UserQueueID) {
		user := ram.roomUserManager.GetGameUserByConnID(clt.ConnID())

		if user == nil {
			mlog.Error("接受到用户%v排隊消息,但用户信息不存在", clt)
			//ram.kickUser(clt.ConnID(), "非法操作")
			return nil
		}

		// 檢查服務器狀態
		if !ram.checkServerStatus(user) {
			return nil
		}
		// 不是站立或者坐下狀態
		if !user.IsFree() && !user.IsSit() {
			mlog.Error("接受到用户%v排隊消息,但用戶狀態不是坐下或者站立", user)
			//ram.kickUser(clt.ConnID(), "非法操作")
			return nil
		}

		ram.onUserQueue(clt, user)
	}

	// 用戶準備
	if m.BClassID == int32(netproto.MessageBClassID_GameRoom) && m.SClassID == int32(netproto.GameRoomClassID_GameReadyID) {
		user := ram.roomUserManager.GetGameUserByConnID(clt.ConnID())
		if user == nil {
			mlog.Error("接受到用户%v準備消息,但用戶信息不存在", clt)
			//ram.kickUser(clt.ConnID(), "非法操作")
			return nil
		}

		// 檢查服務器狀態
		if !ram.checkServerStatus(user) {
			return nil
		}

		// 如果用戶是free狀態
		if user.IsFree() {
			ram.onUserQueue(clt, user)
		}

		// 如果用戶是坐下狀態
		if !user.IsSit() {
			//mlog.Error("接受到用户%v准备消息,但用户%v不是坐下状态", clt, user)
			//ram.kickUser(clt.ConnID(), "非法操作")
			return nil
		}

		ram.onUserReady(clt, user)
	}

	// 查詢牌桌資訊
	if m.BClassID == int32(netproto.MessageBClassID_GameRoom) && m.SClassID == int32(netproto.GameRoomClassID_ROOM_Read_TableInfoID) {
		if RoomTableInfoCreator == nil {
			return nil
		}

		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_ROOM_TableInfoID)
		m.MsgData = RoomTableInfoCreator(ram.logics)

		return []*msg.Message{m}
	}

	// 玩家自行選桌
	if m.BClassID == int32(netproto.MessageBClassID_GameRoom) && m.SClassID == int32(netproto.GameRoomClassID_ROOM_ChooseTableID) {
		if ram.gameQueue == nil {
			return nil
		}
		if _, ok := ram.gameQueue.(*gamequeue.UserChooseGameQueue); !ok {
			return nil
		}
		user := ram.roomUserManager.GetGameUserByConnID(clt.ConnID())
		if user == nil {
			return nil
		}
		data := new(netproto.UserChooseTable)
		err := proto.Unmarshal(m.MsgData.([]byte), data)
		if err != nil {
			mlog.Error("用戶%v選擇牌桌, 但封包解析錯誤, %v", clt, err)
			return nil
		}
		ram.onUserSit(clt, user, data.GetTableNo(), data.GetSeatNo())
	}

	// 遊戲消息
	if m.BClassID == int32(netproto.MessageBClassID_Game) {
		user := ram.roomUserManager.GetGameUserByConnID(clt.ConnID())

		if user == nil {
			mlog.Error("接受到用戶%v遊戲消息,但用戶不在房間中", clt)
			//ram.kickUser(clt.ConnID(), "非法操作")
			return nil
		}

		if !user.IsGameing() {
			mlog.Error("接受到用戶%v遊戲消息%s，但不在遊戲狀態", user, m)
			//ram.kickUser(clt.ConnID(), "非法操作")
			return nil
		}

		if ram.OnGameMessage != nil {
			ram.OnGameMessage(user, m)
		}
		ram.onGameMessage(user, m)
	}

	return nil
}

func setLoginRPC(clt *network.ConnInfo, mm *netproto.LoginGameRoomInfo, ram *RoomActionManager) *netproto.RPCInfo {
	// 構造rpcinfo
	rpcinfo := new(netproto.RPCInfo)
	rpcinfo.ConnID = proto.String(clt.ConnID())
	rpcinfo.UserID = proto.Int32(mm.GetUserID())
	rpcinfo.Cer = proto.String(mm.GetCer())
	rpcinfo.RouteServerID = proto.Int32(ram.room.ServerID)
	rpcinfo.IPAddress = proto.String(clt.IPAddress)
	return rpcinfo
}

// OnRPCMessage 響應RPC消息
func (ram *RoomActionManager) OnRPCMessage(m *msg.Message, ctx *module.MouduleCtx) {
	defer util.PanicErrStack()
	// 登錄房間返回結果
	if m.BClassID == int32(netproto.MessageBClassID_GameRoom) && m.SClassID == int32(netproto.GameRoomClassID_LoginRoomRetID) {
		rpcinfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
			mm := &netproto.LoginGameRoomRet{}
			return mm
		})

		mm := rmm.(*netproto.LoginGameRoomRet)

		// 通知用戶登錄結果
		notifyLoginRetFunc := func() {
			retmsg := new(msg.Message)
			retmsg.BClassID = int32(netproto.MessageBClassID_GameRoom)
			retmsg.SClassID = int32(netproto.GameRoomClassID_LoginRoomRetID)
			retmsg.MsgData = mm

			// 房間信息
			rinfo := new(netproto.RoomInfo)
			rinfo.GameID = proto.Int32(ram.serverCfg.GameID)
			rinfo.GameName = proto.String(ram.serverCfg.GameName)
			rinfo.Flag = proto.String(ram.serverCfg.Flag)
			rinfo.BaseScore = proto.Int32(ram.serverCfg.BaseScore)
			mm.RoomInfo = rinfo

			mlog.Debug("通知用戶登錄結果,%v", mm)
			if ctx != nil && ctx.Server != nil {
				ctx.Server.SendMessageByConnID(rpcinfo.GetConnID(), retmsg)
			}
		}

		// 通知數據庫解鎖
		notifyDbReleaseFunc := func(userid int32, tradeId int64) {
			// 通知dbserver 登出用户
			quitur := gamelogic.NewGameUser()
			quitur.UserID = userid
			quitur.TradeID = tradeId

			ram.scoreStorage.OnUserLeave(quitur)
		}

		// 設置用戶數據
		setPlayerFunc := func(p *player.Player) {

			nur := gamelogic.NewGameUser()

			nur.ConnInfo = p.Conn
			nur.WinCount = mm.UserData.GetWinCount()
			nur.DrawCount = mm.UserData.GetDrawCount()
			nur.LoseCount = mm.UserData.GetLoseCount()
			nur.LastWinLose = mm.UserData.GetLastWinLose()
			nur.TimeSeq = mm.UserData.GetTimeSeq()

			nur.HDCode = mm.GetLoginRequestData().GetHDCode()
			nur.HDType = mm.GetLoginRequestData().GetHDType()

			nur.UserID = mm.UserData.GetUserID()
			nur.NickName = mm.UserData.GetNickName()
			nur.Score = mm.UserData.GetCashAmount()
			nur.LoginScore = mm.UserData.GetCashAmount()
			nur.UserSex = mm.GetUserData().GetSex()
			nur.UserType = mm.GetUserData().GetUserType()
			nur.FaceID = mm.GetUserData().GetFaceID()
			nur.GameBuff = mm.GetUserData().GetGameBuff()
			nur.TradeID = mm.GetUserData().GetTradeId()
			nur.SumCharge = mm.GetUserData().GetSumCharge()
			nur.SumConvert = mm.GetUserData().GetSumConvert()
			nur.HasBindPay = mm.GetUserData().GetHasBindPay()
			nur.LoginName = mm.GetUserData().GetLoginName()
			nur.AlipayAccount = mm.GetUserData().GetAlipayAccount()
			nur.ExperienceCash = mm.GetUserData().GetExperienceCash()

			// 設置player數據
			p.Cer = mm.GetLoginRequestData().GetCer()
			p.HDCode = mm.GetLoginRequestData().GetHDCode()
			p.HDType = mm.GetLoginRequestData().GetHDType()
			p.UserID = mm.GetLoginRequestData().GetUserID()

			p.UserData = nur
		}

		// 獲取客戶端信息
		pr := ram.roomUserManager.GetPlayerByConnID(rpcinfo.GetConnID())

		// 如果客戶端已經登錄
		if pr != nil && !pr.IsGuest() {
			if mm.GetCode() == 1 {
				newuserid := mm.GetUserData().GetUserID()
				if ram.roomUserManager.GetGameUserByUserID(newuserid) == nil {
					// 數據庫解鎖
					notifyDbReleaseFunc(newuserid, mm.GetUserData().GetTradeId())
				}
			}
			return
		}

		// 登錄成功
		if mm.GetCode() == 1 {

			// 檢查用戶信息
			exur := ram.roomUserManager.GetGameUserByUserID(mm.GetLoginRequestData().GetUserID())

			// 檢查房間人數上限
			if ram.roomUserManager.UserCount() >= ram.room.MaxUser {
				// 賬號不在房間中
				if exur == nil {
					mlog.Error("用戶%d進入房間,房間人數已滿,已踢出用户.", mm.GetUserData().GetUserID())
					ram.kickUser(rpcinfo.GetConnID(), "房间已满.")
					// 數據庫解鎖
					notifyDbReleaseFunc(mm.GetUserData().GetUserID(), mm.GetUserData().GetTradeId())
					return
				}
			}

			// 賬號在房間中
			if exur != nil {

				// 當前連接已經斷開
				if pr == nil {
					return
				} else if pr.Conn.ConnID() == exur.ConnInfo.ConnID() { // 同一個連接重發登錄數據包
					return
				}

				// 賬號正在遊戲中
				if exur.IsGameing() {
					exconnid := exur.ConnInfo.ConnID()

					// 通知之前用戶離線
					ram.OnUserLeave(exconnid)

					// 踢出之前用戶
					ram.kickUserExit(exconnid, "您的账号在别处登录,您已被强制下线.")

					//example.invalid("ex:%s,new:%s", exconnid, rpcinfo.GetConnID())
				}

				/**
				服務器內存在用戶數據有兩種情況, 分別做不同的處理.
				1. 用戶在遊戲中的話, 斷線重連並獲取內存的用戶數據.
				2. 用戶不在遊戲中, 將原數據做退房寫分及現金記錄, 然後新的用戶數據的金額有可能在退房寫分時才被異動, 所以需要用原數據的金額來取代
				*/
				if exur.IsOffLine() { // 用戶斷線狀態
					mlog.Debug("檢測到用戶%v斷線重連.%v", exur.UserID, ram.roomUserManager.GetPlayerByUserID(exur.UserID).Conn)
					extableno := exur.TableNo
					exur.ConnInfo = pr.Conn

					// 設置服務器中用戶金幣數
					ram.scoreStorage.OnUserRecover(exur, mm.GetUserData().GetCashAmount())

					// 設置返回數據包中用戶的金幣數 (為確保用戶的金幣數及TimeSeq 不被斷線重連的數據取代)

					// 有可能在斷線到重連的中間，金幣產生過異動
					mm.GetUserData().CashAmount = proto.Int64(exur.Score)
					// 需返回上次登入的時間序號給客戶端
					mm.GetUserData().TimeSeq = proto.Int64(exur.TimeSeq)

					// 設置用戶斷線重回
					ram.roomUserManager.SetPlayerRecover(rpcinfo.GetConnID(), exur.UserID)

					mlog.Debug("通知用戶%v斷線重連.%v", exur.UserID, ram.roomUserManager.GetPlayerByUserID(exur.UserID).Conn)

					// 通知用戶登錄結果
					notifyLoginRetFunc()

					// 通知遊戲邏輯用戶斷線重回
					glc := ram.logics.GetGameLogic(extableno)

					if glc != nil {
						glc.UserRecover(exur)
					}

					// 通知用戶狀態改變
					ram.notifyUserStateChange(exur)
				} else { //不在遊戲中且在線的狀況
					/** 2019.10.22 Rory 去掉以下返回原數據代碼
					//设置服务器中用户金币数
					//ram.scoreStorage.OnUserRecover(exur, mm.GetUserData().GetCashAmount())

					//设置返回数据包中用户的金币数
					//mm.GetUserData().CashAmount = proto.Int64(exur.Score)
					//需返回上次登入的時間序號給客戶端
					//mm.GetUserData().TimeSeq = proto.Int64(exur.TimeSeq)
					*/

					// 讓老用戶站立
					ram.onUserFree(exur.ConnInfo, exur)

					/**
					不在遊戲中且在線的狀況
					2019.10.22 Rory 去掉以上返回原數據代碼, 將原數據進行退房寫分記錄.
					2020.1.17 Rory 雖然己經將原數據進行退房寫分, 但因為這中間依然有可能產生金額異動,
					所以還是需要在寫分之後取用原數據的金額來替換掉從db取出的金額

					*/
					ram.scoreStorage.OnUserLeave(exur)

					// 清除連接數據
					ram.roomUserManager.ClearPlayerInfo(exur.ConnInfo.ConnID())

					ram.kickUserExit(exur.ConnInfo.ConnID(), "您的账号在别处登录,您已被强制下线.")

					// 2020.1.17 因為後踼前會將原數據進行退房寫分及現金記錄, 有可能在斷線到重連的中間, 金幣產生過異動
					// 所以在這邊取用了原數據的金額來替換掉從db取出的金額
					mm.GetUserData().CashAmount = proto.Int64(exur.Score)

					// 設置當前連接用戶數據
					ram.roomUserManager.SetPlayerData(rpcinfo.GetConnID(), setPlayerFunc)

					// 通知用戶登錄結果
					notifyLoginRetFunc()

					// 通知用戶排隊
					nur := ram.roomUserManager.GetGameUserByConnID(pr.Conn.ConnID())
					ram.onUserQueue(nur.ConnInfo, nur)
				}
			} else {
				// 當前連接已經斷開
				if pr == nil {
					mlog.Warn("用戶連接己經斷開進行資料庫解鎖, %d_%d", mm.GetUserData().GetUserID(), mm.GetUserData().GetTradeId())
					// 數據庫解鎖
					notifyDbReleaseFunc(mm.GetUserData().GetUserID(), mm.GetUserData().GetTradeId())
					return
				}

				// 檢查用戶金幣數
				if !ram.checkLogonMoney(rpcinfo.GetConnID(), mm.GetUserData().GetCashAmount()) {
					mlog.Warn("用戶金幣不足進行資料庫解鎖, %d_%d", mm.GetUserData().GetUserID(), mm.GetUserData().GetTradeId())
					// 數據庫解鎖
					notifyDbReleaseFunc(mm.GetUserData().GetUserID(), mm.GetUserData().GetTradeId())
					return
				}

				// 設置用戶數據
				ram.roomUserManager.SetPlayerData(rpcinfo.GetConnID(), setPlayerFunc)

				// 通知用戶登錄結果
				notifyLoginRetFunc()

				// 通知用戶排隊 (不是手動就自動)
				nur := ram.roomUserManager.GetGameUserByConnID(pr.Conn.ConnID())
				if ram.gameQueue == nil {
					return
				}
				if _, ok := ram.gameQueue.(*gamequeue.UserReadyQueue); ok {
					ram.notifyUserStateChange(nur)
					return
				}
				ram.onUserQueue(nur.ConnInfo, nur)
			}

			return
		}

		// 登錄失敗,關閉客戶端連接
		mlog.Warn("客戶端%s登錄失敗%s,關閉連接", rpcinfo.GetConnID(), mm.GetMessage())

		// 通知用戶登錄結果
		notifyLoginRetFunc()

		ram.kickUserOff(rpcinfo.GetConnID(), mm.GetMessage())
		return
	}

	// 檢查是否買票
	if m.BClassID == int32(netproto.MessageBClassID_GameRoom) && m.SClassID == int32(netproto.HallMsgClassID_GameCompetitionRetID) {
		rpcinfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
			mm := &netproto.GameCompetitionMsg{}
			return mm
		})

		mm := rmm.(*netproto.GameCompetitionMsg)
		// 房間信息
		rinfo := new(netproto.RoomInfo)
		rinfo.GameID = proto.Int32(ram.serverCfg.GameID)
		rinfo.GameName = proto.String(ram.serverCfg.GameName)
		rinfo.Flag = proto.String(ram.serverCfg.Flag)
		rinfo.BaseScore = proto.Int32(ram.serverCfg.BaseScore)

		retmsg := new(msg.Message)
		retmsg.BClassID = int32(netproto.MessageBClassID_GameRoom)
		retmsg.SClassID = int32(netproto.HallMsgClassID_GameCompetitionRetID)
		ret := new(netproto.GameCompetitionMsg)
		ret.PiSaiTime = proto.Int32(mm.GetPiSaiTime())
		ret.UserID = proto.Int32(mm.GetUserID())
		ret.BestScore = proto.Int32(mm.GetBestScore())
		ret.FreeTimes = proto.Int32(mm.GetFreeTimes())
		ret.Charge = proto.Int32(mm.GetCharge())
		ret.UserRanking = proto.Int32(mm.GetUserRanking())
		retmsg.MsgData = ret

		ctx.Server.SendMessageByConnID(rpcinfo.GetConnID(), retmsg)

		if mm.GetPiSaiTime() <= 0 {
			ram.kickUserOff(rpcinfo.GetConnID(), "比賽結束")
			return
		}
		return
	}

	// 解析RPC消息
	rmm, ok := m.MsgData.(*rpcmsg.RPCMessage)
	if ok {
		exinfo := rmm.RPCInfo.GetExtendInfo()
		if exinfo != "" {
			tno, err := strconv.Atoi(exinfo)
			if err == nil {
				glc := ram.logics.GetGameLogic(int32(tno))
				if glc != nil {
					glc.OnRecvRPCMessage(m)
				}
			}
		}
	}
}

// OnUserLeave 用戶離開
func (ram *RoomActionManager) OnUserLeave(connid string) {
	defer util.PanicErrStack()
	user := ram.roomUserManager.GetGameUserByConnID(connid)

	// 用户不存在
	if user == nil {
		mlog.Warn("用户不存在即離開")
		ram.room.RoomUserManager.RemovePlayer(connid)
		return
	}

	// 站立或者排隊狀態
	if user.IsFree() || user.IsQueue() {

		ram.room.RoomUserManager.RemovePlayer(connid)
		// 寫分
		ram.scoreStorage.OnUserLeave(user)

		return
	}

	lgc := ram.logics.GetGameLogic(user.TableNo)

	if lgc == nil {
		mlog.Warn("尚未完成登入行為即離開 %d", user.UserID)
		ram.room.RoomUserManager.RemovePlayer(connid)
		// 寫分
		ram.scoreStorage.OnUserLeave(user)
		return
	}

	// 通知遊戲邏輯用戶離開
	if lgc.UserLeave(user) || ram.room.HostAndBackup == 2 {
		ram.room.RoomUserManager.RemovePlayer(connid)
		// 寫分
		ram.scoreStorage.OnUserLeave(user)
	} else {
		// 通知用戶斷線
		lgc.UserOffLine(user)

		// 設置用戶離線
		ram.roomUserManager.SetPlayerOffLine(connid)
	}
}

// 用戶登錄
func (ram *RoomActionManager) onUserLogin(rpcinfo *netproto.RPCInfo, logmsg *netproto.LoginGameRoomInfo, orignmsg *msg.Message) {
	var m *msg.Message
	if orignmsg != nil {
		m = orignmsg
	} else {
		m = new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_LoginRoomID)
		m.MsgData, _ = proto.Marshal(logmsg)
	}
	ram.serverCtx.DBMsgQueue.PushMessage(rpcinfo, false, m)
}

// 用戶排隊
func (ram *RoomActionManager) onUserQueue(clt *network.ConnInfo, user *gamelogic.GameUser) {
	defer util.PanicErrStack()
	// 檢查用戶金幣數
	if !ram.checkUserMoney(user) {
		return
	}
	// 檢查服務器狀態
	if !ram.checkServerStatus(user) {
		return
	}
	// 通知遊戲邏輯用戶離開
	if user.IsSit() {
		glc := ram.logics.GetGameLogic(user.TableNo)
		if glc != nil {
			glc.UserLeave(user)
		}
	}
	user.UserState.OnUserQueue(user, ram.room)
	// 通知用戶狀態改變
	ram.notifyUserStateChange(user)
	// 通知撮合邏輯用戶排隊
	if ram.gameQueue != nil {
		ram.gameQueue.UserQueue(user)
	}
	mlog.Debug("用戶%v排隊", user)
}

// 用戶坐下
func (ram *RoomActionManager) onUserSit(clt *network.ConnInfo, user *gamelogic.GameUser, tableno int32, seatno int32) {
	defer util.PanicErrStack()
	//msglog := fmt.Sprintf("用户坐下onUserSit userid[%v]", user.UserID)
	//mylog.MongoLog(user.UserID, 0, 0, msglog)
	//mlog.System(msglog)
	if !user.IsQueue() && !user.IsRobot() {
		mlog.Error("用戶[%d] 狀態不是 %v, 不能坐下", user.UserID, user.UserState.StateName())
		return
	}
	// 檢查服務器狀態
	if !ram.checkServerStatus(user) {
		return
	}
	// 機器人只可以由站立到坐下
	if user.IsRobot() && !user.IsFree() {
		mlog.Error("機器人狀態是 %v, 不能坐下", user.UserState.StateName())
		return
	}

	// 檢查桌號
	if tableno >= ram.room.TableCount || tableno < 0 {
		mlog.Error("用戶[%d] 桌號 %d 超過範圍", user.UserID, tableno)
		return
	}

	// 檢查椅子號
	if seatno >= ram.room.SeatCount || seatno < 0 {
		mlog.Error("用戶[%d] 椅子號 %d 超過範圍", user.UserID, seatno)
		return
	}

	// 檢查用戶金幣數
	if !ram.checkUserMoney(user) {
		return
	}

	lgc := ram.logics.GetGameLogic(tableno)
	if lgc == nil {
		mlog.Error("用戶[%d] 坐下 桌號 %d 對應的遊戲邏輯不存在", user.UserID, tableno)
		return
	}

	lgc.UserEnter(user, tableno, seatno)

	// 通知用戶狀態改變
	if !user.IsRobot() {
		ram.notifyUserStateChange(user)
	}
	mlog.Debug("用戶%v坐下，%d,%d", user, tableno, seatno)
}

// 用戶站立
func (ram *RoomActionManager) onUserFree(clt *network.ConnInfo, user *gamelogic.GameUser) {
	tableno := user.TableNo
	if tableno == gamelogic.FreeTableNo {
		return
	}

	lgc := ram.logics.GetGameLogic(tableno)
	if lgc == nil {
		return
	}

	lgc.UserFree(user)

	// 通知用戶狀態改變
	if !user.IsRobot() {
		ram.notifyUserStateChange(user)
	}

	mlog.Debug("用戶%v站立", user)
}

// 用戶舉手
func (ram *RoomActionManager) onUserReady(clt *network.ConnInfo, user *gamelogic.GameUser) {
	// 如果已經停服
	if !ram.serverStatus.IsServerRuning() {
		return
	}

	tableno := user.TableNo

	if tableno == gamelogic.FreeTableNo {
		return
	}

	lgc := ram.logics.GetGameLogic(tableno)
	if lgc == nil {
		return
	}

	if !user.IsSit() {
		return
	}

	// 檢查用戶金幣數
	ram.scoreStorage.OnUserReady(user)
	if !ram.checkUserMoney(user) {
		return
	}

	// 檢查同桌最大遊戲局數
	if ram.serverCfg.MaxGameCountSameTable > 0 {
		if user.GameCountOnTheTable >= ram.serverCfg.MaxGameCountSameTable {
			example.invalid("用戶%v同桌遊戲超過%d次,重新配桌", user, ram.serverCfg.MaxGameCountSameTable)
			// 用戶站立
			ram.onUserFree(clt, user)

			// 如果不是機器人,再次配桌
			if !user.IsRobot() {
				ram.onUserQueue(clt, user)
			}
		}
	}

	lgc.UserReady(user)

	// 通知用戶狀態改變
	if !user.IsRobot() {
		ram.notifyUserStateChange(user)
	}

	mlog.Debug("用戶%v舉手", user)
}

// 遊戲消息
func (ram *RoomActionManager) onGameMessage(user *gamelogic.GameUser, m *msg.Message) {
	defer util.PanicErrStack()

	if user.IsFree() {
		return
	}

	tableno := user.TableNo

	glc := ram.logics.GetGameLogic(tableno)

	if glc == nil {
		mlog.Warn("用戶%v發送遊戲消息,用戶的桌號為%d，對應的遊戲邏輯不存在", user, tableno)
		ram.kickUser(user.ConnInfo.ConnID(), "非法操作")
		return
	}

	if !glc.IsGaming() {
		mlog.Warn("用戶%v發送遊戲消息,但當前遊戲邏輯沒有進入遊戲狀態", user)
		ram.kickUser(user.ConnInfo.ConnID(), "非法操作")
		return
	}

	if ram.room.GameType == gamelogic.ChenChien {
		glc.GameMessage(user, m)
	}

	if ram.room.GameType == gamelogic.PiSai {
		glc.PiSaiMessage(user, m)
	}

}

// 通知用戶狀態改變
func (ram *RoomActionManager) notifyUserStateChange(user *gamelogic.GameUser) {
	if user.IsOffLine() {
		return
	}

	// 如果用户是机器人
	if user.IsRobot() {
		return
	}

	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.GameMessageClassID_UserStateChangeID)

	ust := new(netproto.UserStateChange)
	ust.User = user.GetProtoBufGameUser()

	m.MsgData = ust

	ram.serverCtx.MsgServer.SendMessage(user.ConnInfo, m)
}

//发送提示消息
func (ram *RoomActionManager) sendTipMessage(connid string, str string, typeid int32) {
	if robot.IsRobotConn(connid) {
		return
	}

	m, _ := roomutil.GetTipMessage(str, typeid)

	ram.serverCtx.MsgServer.SendMessageByConnID(connid, m)
}

//踢人,踢出房间
func (ram *RoomActionManager) kickUser(connid string, str string) {
	mlog.Warn("%s 斷開客戶端%s", str, connid)
	if robot.IsRobotConn(connid) {
		return
	}

	strLen := len(str)
	if strLen > 2 {
		roomutil.KickUserParam(str, connid, roomutil.MsgTypeIDRoomKick, 5, ram.serverCtx.MsgServer)
	} else if ram.room.NewGameServerAddr != "" {
		//roomutil.KickUserToNewRoom(ram.room.NewGameServerAddr, connid, roomutil.MsgTypeID_RoomKick, 6, ram.serverCtx.MsgServer)
		roomutil.KickUserParam(str, connid, roomutil.MsgTypeIDRoomKick, 4, ram.serverCtx.MsgServer)
	} else {
		roomutil.KickUserParam(str, connid, roomutil.MsgTypeIDRoomKick, 4, ram.serverCtx.MsgServer)
	}

}

//踢人,踢出房间
func (ram *RoomActionManager) kickUserWithParam(connid string, str string, param int32) {
	mlog.Warn("%s 斷開客戶端%s", str, connid)
	if robot.IsRobotConn(connid) {
		return
	}

	roomutil.KickUserParam(str, connid, roomutil.MsgTypeIDRoomKick, param, ram.serverCtx.MsgServer)

}

//踢人,退出APP
func (ram *RoomActionManager) kickUserExit(connid string, str string) {
	mlog.Warn("%s 斷開客戶端%s", str, connid)
	if robot.IsRobotConn(connid) {
		return
	}

	roomutil.KickUser(str, connid, roomutil.MsgTypeIDHallKick, ram.serverCtx.MsgServer)

}

//踢人,踢出房间但不发送提示消息
func (ram *RoomActionManager) kickUserOff(connid string, str string) {
	mlog.Warn("%s 斷開客戶端%s", str, connid)
	if robot.IsRobotConn(connid) {
		return
	}

	roomutil.KickUser("", connid, roomutil.MsgTypeIDRoomKick, ram.serverCtx.MsgServer)

}

//检查用户金币数
func (ram *RoomActionManager) checkUserMoney(user *gamelogic.GameUser) bool {
	//检查用户金币数
	usermoney := user.Score

	if usermoney < int64(ram.room.LoginMoney) {
		ram.kickUserWithParam(user.ConnInfo.ConnID(), fmt.Sprintf("您的金币不足!请补充金币后再进行游戏."), roomutil.MsgParamMoneyNotEnough)

		return false
	}

	if ram.room.KickOutMoney != 0 && usermoney >= int64(ram.room.KickOutMoney) {
		ram.kickUserWithParam(user.ConnInfo.ConnID(), fmt.Sprintf("您的金币超出房间限制!请前往『菁英场』进行游戏."), roomutil.MsgParamMoneyOverLoad)

		return false
	}

	return true
}

// 檢查進入房間時的金幣數
func (ram *RoomActionManager) checkLogonMoney(connid string, usermoney int64) bool {
	if usermoney < int64(ram.room.LoginMoney) {
		ram.kickUserWithParam(connid, fmt.Sprintf("您的金币不足!请补充金币后再进行游戏."), roomutil.MsgParamMoneyNotEnough)
		return false
	}

	if ram.room.KickOutMoney != 0 && usermoney >= int64(ram.room.KickOutMoney) {
		ram.kickUserWithParam(connid, fmt.Sprintf("您的金币超出房间限制!请前往『菁英场』进行游戏"), roomutil.MsgParamMoneyOverLoad)

		return false
	}

	return true
}

// 檢查服務器狀態
func (ram *RoomActionManager) checkServerStatus(ur *gamelogic.GameUser) bool {
	//glc := ram.logics.GetGameLogic(ur.TableNo)
	//if glc != nil && glc.IsGaming() {
	//	return true
	//}
	if !ram.serverStatus.IsServerRuning() {
		mlog.Debug("服務器狀態停止[%v] serverid[%v] ProxyIP[%v] ServerAddr[%v] ", !ram.serverStatus.IsServerRuning(), ram.serverCfg.ServerID,
			ram.serverCfg.ProxyIP, ram.serverCfg.ServerAddr)
		ram.kickUser(ur.ConnInfo.ConnID(), ram.serverStatus.GetStopServerTipMsg())
		return false
	}

	return true
}

// 根據用戶ID獲取用戶信息
func (ram *RoomActionManager) getGameUserByID(userid int32) *gamelogic.GameUser {
	ur := ram.roomUserManager.GetGameUserByUserID(userid)
	if ur == nil {
		robot := ram.robotService.GetRobotByID(userid)
		if robot != nil {
			ur = robot.UserInfo
		}
	}

	return ur
}

// 根據桌子號獲取用戶
func (ram *RoomActionManager) getGameUserByTableNo(tableno int32) []*gamelogic.GameUser {
	return ram.room.GetTableUser(tableno)
}

// OnHTTPMsg 響應http消息
func (ram *RoomActionManager) OnHTTPMsg(semaphore util.Semaphore, w http.ResponseWriter, req *http.Request) bool {
	defer util.PanicErrStack()

	tableno := -1

	//解析桌号
	stableno := req.PostFormValue("tableno")
	argRoom := req.PostFormValue("room")
	//	if stableno == "" {
	//		stableno = req.URL.Query().Get("tableno")
	//	}

	if stableno != "" {
		var err error
		tableno, err = strconv.Atoi(stableno)

		if err != nil {
			io.WriteString(w, "tableno error")
			return false
		}
	}

	mlog.Debug("收到http消息,桌号:%d", tableno)
	//mlog.System("===================== room %v, tableno %d", argRoom, tableno)
	if tableno == -1 && argRoom != "1" {
		ram.logics.Range(func(glc gamelogic.GameLogic) {
			glc.OnRecvHTTPMessage(req)
		})
		io.WriteString(w, "ok")
	} else if argRoom == "1" { //這是對於整個房間的操作
		//mlog.System("----------------------")
		ram.room.SuperToolService.HTTPService(w, req)
	} else {
		glc := ram.logics.GetGameLogic(int32(tableno))
		if glc == nil {
			io.WriteString(w, "tableno error")
		} else {
			ret := glc.OnRecvHTTPMessage(req)
			if ret != nil {
				data, _ := proto.Marshal(ret)

				if data != nil {
					encodeStr := base64.StdEncoding.EncodeToString(data)
					io.WriteString(w, encodeStr)
				}
			} else {
				io.WriteString(w, "ok")
			}
		}
	}

	return false
}

// NotifyConn ?
func NotifyConn() {

	//注册协调服务器
	//ntfHandler := serverhandler.NewNotifyHandler(nsserveraddr, serverid, servertype, gameid, hostAndBackup)

	//	//创建到NsServer的TCP连接
	//	msgparser := msg.NewProtoMsgParser()
	//	msgpacker := network.NewHeadMsgPacker(conf.MSGPACKERBUF, true)
	//
	//	//nshandler.evtchan = msgserver.NewMsgEvtChan(100)
	//	//nshandler.tcpMsgClient = msgserver.NewTcpMsgClient(msgpacker, nshandler.evtchan, msgparser, true)
	//	//nshandler.cgx = cgx
	//	//
	//	//go nshandler.tcpEventChanFunc(nshandler.evtchan)
	//	//
	//	//err := nshandler.tcpMsgClient.Connect(nshandler.nsAddr)
	//	//
	//	//if err != nil {
	//	//	mlog.Error("向协调服务器%s注册失败%v", nshandler.nsAddr, err)
	//	//	return
	//	//}
}
