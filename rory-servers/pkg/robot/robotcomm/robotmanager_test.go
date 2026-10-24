package robotcomm

import (
	"fmt"
	"strconv"
	"testing"

	"example.invalid/stretchr/testify/assert"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
)

func TestRobotManager_AddRobot(t *testing.T) {
	rm := NewRobotManager()
	robots := GetRobots(10)
	assert.Equal(t, len(robots), 10, "robots len should be 10")
	// test add
	for _, r := range robots {
		rm.AddRobot(r)
	}
	// test robot count
	assert.Equal(t, rm.RobotCount(), len(robots), "they should be equal")
	// test get free robot
	freeRobot := rm.GetFreeRobot(gamelogic.FreeTableNo)
	assert.Equal(t, freeRobot.UserInfo.IsGameing(), false, "freeRobot IsGameing should be false")
	// test get not gaming robot count
	freeRobot.UserInfo.UserState = new(gamelogic.UserGameState)
	assert.Equal(t, len(rm.GetNotGamingRobot()), 9, "GetNotGamingRobot should be 9")
	assert.Equal(t, rm.FreeRobotCount(), 9, "FreeRobotCount should be 9")
	assert.Equal(t, rm.GameRobotCount(), 1, "GameRobotCount should be 1")
	// test get all robot
	assert.Equal(t, len(rm.GetAllRobot()), 10, "GetAllRobot should be 10")
	// get robot by userid
	assert.Equal(t, rm.GetRobot(freeRobot.UserInfo.UserID), freeRobot, "they should be equal")
	// test remove free robot
	for i := 0; i < 9; i++ {
		r := rm.GetFreeRobot(gamelogic.FreeTableNo)
		rm.RemoveRobot(r.UserInfo.UserID)
	}
	assert.Equal(t, rm.FreeRobotCount(), 0, "FreeRobotCount should be 0")
	assert.Equal(t, rm.GameRobotCount(), 1, "GameRobotCount should be 1")
	robotList := rm.GetAllRobot()
	rm.RemoveRobot(robotList[0].UserInfo.UserID)
	assert.Equal(t, rm.RobotCount(), 0, "RobotCount should be 0")
}

func GetRobots(cn int) (robots []*Robot) {
	indexUserIdx := int32(0)
	genUser := func(nickname string, score int64) (nur *gamelogic.GameUser) {
		indexUserIdx++
		baseuid, _ := strconv.Atoi(fmt.Sprintf("9%d00000", gamelogic.UserTypeRebot))
		userid := int32(baseuid) + indexUserIdx
		nur = gamelogic.NewGameUser()
		nur.ConnInfo = nil
		nur.WinCount = 0
		nur.DrawCount = 0
		nur.LoseCount = 0

		nur.HDCode = "000000"
		nur.HDType = 1

		nur.UserID = userid
		nur.NickName = nickname
		nur.Score = score
		nur.LoginScore = score
		nur.UserSex = 0
		//设定用户类型为机器人
		nur.UserType = gamelogic.UserTypeRebot // urdata.GetUserType()
		nur.FaceID = 3
		return
	}
	for i := 0; i < cn; i++ {
		//foo := new robotcomm.Robot{genUser("台中市", 1000000), ""}
		foo := &Robot{genUser("台中市", 1000000), 5, 0, false}
		robots = append(robots, foo)
	}
	return
}
