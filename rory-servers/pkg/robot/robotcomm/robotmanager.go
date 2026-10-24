package robotcomm

import (
	"sync"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
)

// RobotManager 機器人管理器
type RobotManager struct {
	robots sync.Map //2019.12.24 改用線程安全的結構 by Rory
	//robots map[int32]*Robot
}

// NewRobotManager 創建機器人管理器
func NewRobotManager() *RobotManager {
	rm := new(RobotManager)
	//rm.robots = make(map[int32]*Robot, 0)
	rm.robots = sync.Map{}

	return rm
}

// AddRobot 添加機器人
func (rm *RobotManager) AddRobot(r *Robot) {
	rm.robots.Store(r.UserInfo.UserID, r)
	//rm.robots[r.UserInfo.UserID] = r
}

// RemoveRobot 移除機器人
func (rm *RobotManager) RemoveRobot(userid int32) {
	rm.robots.Delete(userid)
	//delete(rm.robots, userid)
}

// GetRobot 根據ID獲取機器人
func (rm *RobotManager) GetRobot(userid int32) *Robot {
	if r, ok := rm.robots.Load(userid); ok {
		return r.(*Robot)
	}
	return nil
}

// GetAllRobot 獲取所有機器人
func (rm *RobotManager) GetAllRobot() []*Robot {
	robots := make([]*Robot, 0)
	rm.robots.Range(func(k, r interface{}) bool {
		robots = append(robots, r.(*Robot))
		return true
	})
	return robots
}

// RobotCount 機器人數量
func (rm *RobotManager) RobotCount() int {
	length := 0
	rm.robots.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	return length
}

// GameRobotCount 遊戲中的機器人數量
func (rm *RobotManager) GameRobotCount() int {
	length := 0
	rm.robots.Range(func(key, r interface{}) bool {
		if r.(*Robot).UserInfo.IsGameing() {
			length++
		}
		return true
	})
	return length
}

// GetNotGamingRobot 不在遊戲中的機器人數量
func (rm *RobotManager) GetNotGamingRobot() map[int32]*Robot {
	bots := make(map[int32]*Robot, 0)
	rm.robots.Range(func(key, val interface{}) bool {
		r := val.(*Robot)
		if r.UserInfo.IsFree() || r.UserInfo.IsSit() || r.UserInfo.IsReady() {
			bots[r.UserInfo.UserID] = r
		}
		return true
	})
	return bots
}

// FreeRobotCount 空閒的機器人數量
func (rm *RobotManager) FreeRobotCount() int {
	cn := 0
	rm.robots.Range(func(key, val interface{}) bool {
		r := val.(*Robot)
		if r.UserInfo.IsFree() || r.UserInfo.IsSit() || r.UserInfo.IsReady() {
			cn++
		}
		return true
	})
	return cn
}

// GetFreeRobot 獲取一個空閒的機器人
func (rm *RobotManager) GetFreeRobot(exTableNo int32) (robot *Robot) {
	rm.robots.Range(func(key, val interface{}) bool {
		r := val.(*Robot)
		if r.UserInfo.IsFree() {
			if exTableNo == gamelogic.FreeTableNo { // 排除的桌號若為-1
				robot = r
				return false
			}

			if r.UserInfo.LastTableNo == exTableNo { // 排除最後一次遊戲的桌號
				return true
			}
			robot = r
			return false
		}
		return true
	})
	// if robot == nil {
	// 	mlog.Debug("獲取不到空閒機器人")
	// }
	return
}

// ClearFreeRobotLastTable 清除機器人上一次桌號
func (rm *RobotManager) ClearFreeRobotLastTable(tableNo int32) {
	rm.robots.Range(func(key, val interface{}) bool {
		r := val.(*Robot)
		if r.UserInfo.IsFree() && r.UserInfo.LastTableNo == tableNo {
			r.UserInfo.LastTableNo = gamelogic.FreeTableNo
		}
		return true
	})
}
