package roomlogic

import (
	"fmt"
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/mylog"
)

// RoomActionController 房間動作控制器
type RoomActionController struct {
	roomuserManager         *gamelogic.RoomUserManager
	roomactionManager       *RoomActionManager
	room                    *gamelogic.Room
	userPlayedManager       *gamelogic.UserPlayedManager
	disableEnterGamingTable bool
}

// NewRoomActionController 創建房間動作控制器對象
func NewRoomActionController(roomuserManager *gamelogic.RoomUserManager, roomactionManager *RoomActionManager, room *gamelogic.Room, cfgs *config.Group) *RoomActionController {
	rac := new(RoomActionController)
	rac.roomuserManager = roomuserManager
	rac.roomactionManager = roomactionManager
	rac.room = room
	cfg := cfgs.GetConfig(gamelogic.ConfigFileNameGameServer)
	playLimit := cfg.GetConfigInt32(gamelogic.ConfigSectionNameGameQueue, "playlimit")

	if val, err := strconv.ParseBool(cfg.GetConfigValue(gamelogic.ConfigSectionNameGameQueue, "disableentergamingtable")); err == nil {
		rac.disableEnterGamingTable = val
	}

	rac.userPlayedManager = gamelogic.NewUserPlayedManager(playLimit)
	return rac
}

// GetFreeTable 取得一個空閒桌子
func (rac *RoomActionController) GetFreeTable() *gamelogic.Table {
	for _, t := range rac.room.Tables {
		if t.GameLogic.IsGaming() {
			continue
		}
		if len(t.GetTableUser()) > 0 {
			continue
		}
		return t
	}
	return nil
}

// UserSit 用戶坐下
func (rac *RoomActionController) UserSit(user *gamelogic.GameUser, tableno int32, seatno int32) {
	//mlog.System("UserSit user %d seatNo %d", user.UserID, seatno)
	//	pr := rac.roomuserManager.GetPlayerByUserID(userid)

	//	if pr == nil {
	//		return
	//	}

	//	user := rac.roomuserManager.GetGameUserByUserID(userid)

	//	if user == nil || !user.IsQueue() {
	//		return
	//	}

	//example.invalid("服务器使用户%v坐下.座位号:%d,%d", user, tableno, seatno)
	rac.roomactionManager.onUserSit(user.ConnInfo, user, tableno, seatno)
}

// UserFree 用戶站立
func (rac *RoomActionController) UserFree(user *gamelogic.GameUser) {
	//	pr := rac.roomuserManager.GetPlayerByUserID(userid)

	//	if pr == nil {
	//		return
	//	}

	//	user := rac.roomuserManager.GetGameUserByUserID(userid)

	//	if user == nil || (!user.IsSit() && !user.IsReady()) {
	//		return
	//	}

	//example.invalid("服务器使用户%v站立", user)
	rac.roomactionManager.onUserFree(user.ConnInfo, user)
}

// UserReady 用戶準備
func (rac *RoomActionController) UserReady(user *gamelogic.GameUser) {
	//	pr := rac.roomuserManager.GetPlayerByUserID(userid)

	//	if pr == nil {
	//		return
	//	}

	//	user := rac.roomuserManager.GetGameUserByUserID(userid)

	//	if user == nil || !user.IsSit() {
	//		return
	//	}

	//example.invalid("服务器使用户%v举手", user)
	rac.roomactionManager.onUserReady(user.ConnInfo, user)
}

// GetQueueUser 取得排隊中用戶
func (rac *RoomActionController) GetQueueUser() []*gamelogic.GameUser {
	return rac.roomuserManager.GetQueueUser()
}

// TotalUserCount 取得所有用戶數
func (rac *RoomActionController) TotalUserCount() int32 {
	return rac.roomuserManager.UserCount()
}

// GetFreeSeatNew 取得空閒座位
func (rac *RoomActionController) GetFreeSeatNew(user *gamelogic.GameUser, exIP string, exLastGameTableNo int32) (int32, int32) {
	if rac.userPlayedManager.PlayLimit > 0 {
		// 包含同桌率的過濾方式
		return rac.GetFreeSeatWithUPM(user, exIP, exLastGameTableNo)
	}
	// 舊的獲取空閒座位
	return rac.GetFreeSeat(user, exIP, exLastGameTableNo)
}

// GetFreeSeat 獲取空閒座位
func (rac *RoomActionController) GetFreeSeat(user *gamelogic.GameUser, exIP string, exLastGameTableNo int32) (int32, int32) {
	extableno := user.LastTableNo
	// 對桌子進行排序
	tbs := rac.room.GetSortedTable()

	for _, t := range tbs {
		tno := t.TableNo

		// 排除桌号
		if int32(tno) == extableno {
			continue
		}

		glc := rac.roomactionManager.logics.GetGameLogic(int32(tno))

		if glc == nil {
			continue
		}

		// 排除遊戲中牌桌
		if rac.disableEnterGamingTable && glc.IsGaming() {
			continue
		}

		// 拒絕同IP
		if exIP != "" {
			issameIP := false
			for _, ur := range t.GetTableUser() {
				if !ur.IsRobot() && ur.ConnInfo.IPAddress == exIP {
					issameIP = true
					break
				}
			}

			if issameIP {
				continue
			}
		}
		// 拒絕最後一次遊戲用戶
		if exLastGameTableNo != gamelogic.FreeTableNo {
			isSameTableGame := false
			for _, ur := range t.GetTableUser() {
				if !ur.IsRobot() && ur.LastTableNo == exLastGameTableNo {
					isSameTableGame = true
					break
				}
			}

			if isSameTableGame {
				continue
			}
		}

		tno, sno := glc.GetFreeSeat()

		if tno != gamelogic.FreeTableNo && sno != gamelogic.FreeSeatNo {
			// logmsg := fmt.Sprintf("B 获取空闲座位 userid[%v] tableno[%v] seatno[%v]", user.UserID, user.TableNo, sno)
			// mylog.MongoLog(user.UserID, 0, user.TableNo, logmsg)
			return tno, sno
		}
	}

	logmsg := fmt.Sprintf("所有桌子 無座位 userid[%v] tableno[%v] seatno[%v]", user.UserID, user.TableNo, gamelogic.FreeSeatNo)
	mylog.MongoLog(user.UserID, 0, user.TableNo, logmsg)

	return gamelogic.FreeTableNo, gamelogic.FreeSeatNo
}

// GetWinLoseFreeSeat 輸贏對戰型, 取得桌子空位
func (rac *RoomActionController) GetWinLoseFreeSeat(extableno int32, exIP string, exLastGameTableNo int32, lastWinLose int32, islessthanuser bool) (int32, int32) {
	// 對桌子進行排序
	tbs := rac.room.GetSortedTable()
	// 最後的99張固定不拿來分配給用戶
	limit := int32(len(tbs) - 20)
	for _, t := range tbs {

		tno := t.TableNo
		if tno >= limit {
			continue
		}
		if t.GameLogic.IsGaming() {
			continue
		}

		glc := rac.roomactionManager.logics.GetGameLogic(int32(tno))

		if glc == nil {
			continue
		}

		if !islessthanuser {

			// 排除桌號
			if int32(tno) == extableno {
				continue
			}

			// 拒絕同IP
			if exIP != "" {
				issameIP := false
				for _, ur := range t.GetTableUser() {
					if !ur.IsRobot() && ur.ConnInfo.IPAddress == exIP {
						issameIP = true
						break
					}
				}

				if issameIP {
					continue
				}
			}

			// 贏的跟贏的打
			isAllWinUser := true
			if lastWinLose == int32(1) {
				for _, ur := range t.GetTableUser() {
					if ur.LastWinLose != int32(1) {
						isAllWinUser = false
						break
					}
				}
				if !isAllWinUser {
					continue
				}
			}

			if isAllWinUser && lastWinLose == int32(1) {
				tno, sno := glc.GetFreeSeat()

				if tno != gamelogic.FreeTableNo && sno != gamelogic.FreeSeatNo {
					return tno, sno
				}
			}

			// 輸贏配桌
			if extableno != gamelogic.FreeTableNo {
				isSameTableGame := false
				for _, ur := range t.GetTableUser() {
					//mlog.Debug(fmt.Sprintf("[winlose] 01 userid : %v", ur.UserID))
					if !ur.IsRobot() && ur.LastWinLose != lastWinLose {
						mlog.Debug("[winlose] 04 userid[%v] isSameTableGame = true", ur.UserID)
						isSameTableGame = true
						break
					}
				}

				if isSameTableGame {
					continue
				}
			}
			// 第一次進遊戲
			if extableno == gamelogic.FreeTableNo {
				isSameTableGame := false
				for _, ur := range t.GetTableUser() {
					if ur.LastWinLose != lastWinLose {
						isSameTableGame = true
						break
					}
				}

				if isSameTableGame {
					continue
				}
			}
		}

		tno, sno := glc.GetFreeSeat()

		if tno != gamelogic.FreeTableNo && sno != gamelogic.FreeSeatNo {
			return tno, sno
		}
	}

	return gamelogic.FreeTableNo, gamelogic.FreeSeatNo
}

// GetFreeSeatWithUPM 獲取空閒座位
func (rac *RoomActionController) GetFreeSeatWithUPM(user *gamelogic.GameUser, exIP string, exLastGameTableNo int32) (int32, int32) {
	// 對桌子進行排序
	tbs := rac.room.GetSortedTable()

	for _, t := range tbs {
		tno := t.TableNo

		// 排除桌號
		if int32(tno) == user.LastTableNo {
			continue
		}

		glc := rac.roomactionManager.logics.GetGameLogic(int32(tno))

		if glc == nil {
			continue
		}

		// 排除遊戲中牌桌
		if rac.disableEnterGamingTable && glc.IsGaming() {
			continue
		}

		// 拒絕同IP
		if exIP != "" {
			issameIP := false
			for _, ur := range t.GetTableUser() {
				if !ur.IsRobot() && ur.ConnInfo.IPAddress == exIP {
					issameIP = true
					break
				}
			}

			if issameIP {
				continue
			}
		}
		// 拒絕最後一次遊戲用戶
		if exLastGameTableNo != gamelogic.FreeTableNo {
			isSameTableGame := false
			for _, ur := range t.GetTableUser() {
				if !ur.IsRobot() && ur.LastTableNo == exLastGameTableNo {
					isSameTableGame = true
					break
				}
			}

			if isSameTableGame {
				continue
			}
		}
		// 拒絕同桌率過高的遊戲用戶
		if rac.userPlayedManager.PlayLimit > 0 {
			isExists := false
			for _, ur := range t.GetTableUser() {
				if !ur.IsRobot() && rac.userPlayedManager.IsExists(user.UserID, ur.UserID) {
					isExists = true
				}
			}
			if isExists {
				continue
			}
		}

		tno, sno := glc.GetFreeSeat()

		if tno != gamelogic.FreeTableNo && sno != gamelogic.FreeSeatNo {
			// 同桌用戶寫入表
			if rac.userPlayedManager.PlayLimit > 0 {
				for _, ur := range t.GetTableUser() {
					if !ur.IsRobot() {
						rac.userPlayedManager.AddUserID(user.UserID, ur.UserID)
						rac.userPlayedManager.AddUserID(ur.UserID, user.UserID)
					}
				}
			}
			return tno, sno
		}
	}

	return gamelogic.FreeTableNo, gamelogic.FreeSeatNo
}

// GetNextTableSeat 取得下張桌子座位
func (rac *RoomActionController) GetNextTableSeat(user *gamelogic.GameUser) (int32, int32) {
	handler := rac.room.SuperToolService.GetSuperToolHandler()

	if tno := handler.GetNextTable(user.UserID); tno != -1 {

		table := rac.room.GetTableByNo(tno)
		if table.GameLogic.IsGaming() {
			return gamelogic.FreeTableNo, gamelogic.FreeSeatNo
		}
		mlog.Debug("==========%d", tno)
		tno, sno := table.GameLogic.GetFreeSeat()
		if tno != gamelogic.FreeTableNo && sno != gamelogic.FreeSeatNo {
			mlog.Debug("==================== tno %d, sno %d", tno, sno)
			if tno == tno {
				handler.ResetNextTable(user.UserID)
				handler.SaveVipUser(user.UserID)
			}
			return tno, sno
		}

	}
	return gamelogic.FreeTableNo, gamelogic.FreeSeatNo
}
