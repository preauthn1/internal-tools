package gamelogic

import (
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// UserPlayedManager 玩家管理器
type UserPlayedManager struct {
	PlayedUserID map[int32][]int32
	ExpiredTime  int64
	PlayLimit    int32
	lock         sync.Mutex
}

// NewUserPlayedManager 創建玩家管理器
func NewUserPlayedManager(playLimit int32) *UserPlayedManager {
	upm := new(UserPlayedManager)
	upm.PlayedUserID = make(map[int32][]int32)
	upm.ExpiredTime = time.Now().Unix() + (60 * 60 * 12) // now + 12 hours
	upm.PlayLimit = playLimit
	return upm
}

// AddUserID 新增玩家
func (upm *UserPlayedManager) AddUserID(cuid int32, tuid int32) {

	//now := time.Now().Unix()
	//if now > upm.ExpiredTime {
	//	upm.PlayedUserID = make(map[int32][] int32)
	//}
	if upm.PlayedUserID[cuid] == nil {
		upm.PlayedUserID[cuid] = make([]int32, upm.PlayLimit)
	}
	if upm.IsExists(cuid, tuid) {
		return
	}
	upm.PlayedUserID[cuid] = upm.PlayedUserID[cuid][1:]
	upm.PlayedUserID[cuid] = append(upm.PlayedUserID[cuid], tuid)
	mlog.NewDefaultLogger().Debug("cuid %d 加入, %d, 目前同桌表 %v", cuid, tuid, upm.PlayedUserID[cuid])
}

// IsExists 是否存在
func (upm *UserPlayedManager) IsExists(cuid int32, tuid int32) bool {

	arr := upm.PlayedUserID[cuid]
	mlog.NewDefaultLogger().Debug("檢查 tuid %d 是否存在 cuid %d 的同桌表內, %v", tuid, cuid, upm.PlayedUserID[cuid])
	for _, uid := range arr {
		if uid == tuid {
			return true
		}
	}
	return false
}
