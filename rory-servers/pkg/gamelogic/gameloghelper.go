package gamelogic

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// GameLogHelper 牌局記錄輔助器
type GameLogHelper struct {
	innerGameLog gamelog.GameLog
	glc          *LogicBase
}

// NewGameLogHelper 創建牌局記錄輔助器
func NewGameLogHelper(gl gamelog.GameLog, glc *LogicBase) *GameLogHelper {
	glh := new(GameLogHelper)
	glh.innerGameLog = gl
	glh.glc = glc

	return glh
}

// Context 當下牌局context
func (glh *GameLogHelper) Context() *gamelog.RoundContext {
	return glh.innerGameLog.Context()
}

// GameBegin 遊戲開始
func (glh *GameLogHelper) GameBegin() *gamelog.RoundContext {
	urs := glh.glc.AllGameUser()

	ctx := glh.innerGameLog.BeginGameLog()

	stepdesc := "遊戲開始."
	stepdesc += fmt.Sprintf("桌子號: %d.", glh.glc.TableNo())
	for _, ur := range urs {
		ctx.AddUser(ur.UserID, ur.TradeID)
		ctx.UserScore(ur.UserID, ur.Score)
	}

	step := gamelog.NewGameStep(stepdesc, gamelog.GameStepTimeTypeFullTime)

	glh.innerGameLog.AddGameStep(step)

	return ctx
}

// GameEnd 遊戲結束
func (glh *GameLogHelper) GameEnd() {
	stepdesc := "遊戲結束."
	step := gamelog.NewGameStep(stepdesc, gamelog.GameStepTimeTypeFullTime)
	glh.innerGameLog.AddGameStep(step)
	glh.innerGameLog.EndGameLog()
}

// AddGameStep 添加遊戲動作
func (glh *GameLogHelper) AddGameStep(format string, a ...interface{}) {
	stepdesc := fmt.Sprintf(format, a...)

	step := gamelog.NewGameStep(stepdesc, gamelog.GameStepTimeTypeSimpleTime)
	glh.innerGameLog.AddGameStep(step)
}

// UpdateGameUsers 更新用戶
func (glh *GameLogHelper) UpdateGameUsers(urs map[int32]*gamelog.UserData) {
	// glh.innerGameLog.UpdateGameUsers(urs)
}

// RemoveGameUsers 移除牌局記錄gameuser
func (glh *GameLogHelper) RemoveGameUsers(userID int32) {
	// glh.innerGameLog.RemoveGameUsers(userID)
}

// Fapai 發牌
func (glh *GameLogHelper) Fapai(uid int32, carddesc string) {
	stepdesc := fmt.Sprintf("發牌,userid:%d,牌型:%s", uid, carddesc)
	step := gamelog.NewGameStep(stepdesc, gamelog.GameStepTimeTypeSimpleTime)
	glh.innerGameLog.AddGameStep(step)
}

// SetNextJiangRecord 配置指定開獎信息
func (glh *GameLogHelper) SetNextJiangRecord(adminuser string, jiangtype string) {
	record := new(netproto.NextJiangRecord)
	record.AdminUser = proto.String(adminuser)
	record.JiangType = proto.String(jiangtype)

	glh.innerGameLog.SetNextJiangRecord(record)
}
