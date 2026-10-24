package model

import (
	"encoding/json"
)

// RobotInfo ?
type RobotInfo struct {
	PlayerID   int32       `json:"playerID"`
	PlayerType int32       `json:"playerType"`
	NickName   string      `json:"nickName"`
	FaceID     int32       `json:"faceID"`
	Sex        int32       `json:"sex"`
	CashAmount json.Number `json:"cashAmount"`
	WinCount   int32       `json:"winCount"`
	LoseCount  int32       `json:"loseCount"`
	DrawCount  int32       `json:"drawCount"`
	LevelID    int32       `json:"levelID"`
}

// RobotListResultModel ?
type RobotListResultModel struct {
	LogoutRobotCount int32       `json:"logoutRobotCount"`
	Robots           []RobotInfo `json:"robots"`
}

// RobotListBody ?
type RobotListBody struct {
	BodyBase
	Data RobotListResultModel `json:"data"`
}
