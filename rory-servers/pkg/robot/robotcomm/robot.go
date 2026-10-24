package robotcomm

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
)

// Robot 機器人
type Robot struct {
	UserInfo      *gamelogic.GameUser // 用戶信息
	RebotLevel    int32               // 機器人等級
	LeaveState    int32               // 1 = 離開
	BuchuBombCard bool                // 是否不出狀態
}

func (r *Robot) String() string {
	return fmt.Sprintf("%v,level:%d,LeaveState:%d", r.UserInfo.String(), r.RebotLevel, r.LeaveState)
}

// IsLeave 是否被凍結住
func (r *Robot) IsLeave() bool {
	if r.LeaveState == 1 {
		return true
	}
	return false
}

// SetLeave ?
func (r *Robot) SetLeave() {
	//mlog.System("SetLeave %d", r.UserInfo.UserID)
	r.LeaveState = 1
}
