package model

import (
	"encoding/json"
)

// GameRoomEnterViewReConnectInfo ?
type GameRoomEnterViewReConnectInfo struct {
	ServerID     int    `json:"serverID"`
	GameID       int    `json:"gameID"`
	ClientAddr   string `json:"clientAddr"`
	ReConnectMsg string `json:"reConnectMsg"`
}

// GameRoomEnterViewPlayerInfo ?
type GameRoomEnterViewPlayerInfo struct {
	PlayerName  string `json:"playerName"`
	DisplayName string `json:"displayName"`
	PlayerID    int32  `json:"playerID"`
	PlayerType  int32  `json:"playerType"`
	FaceID      int32  `json:"faceID"`
	Sex         int32  `json:"sex"`
	WinCount    int32  `json:"winCount"`
	LoseCount   int32  `json:"loseCount"`
	DrawCount   int32  `json:"drawCount"`
	LastWinLose int32  `json:"lastWinLose"`
	TimeSeq     int64  `json:"timeSeq"`
}

// GameRoomEnterViewModel ?
type GameRoomEnterViewModel struct {
	IsReConnect   bool                           `json:"isReConnect"`
	TradeID       int64                          `json:"tradeID"`
	ReconnectInfo GameRoomEnterViewReConnectInfo `json:"reConnectInfo"`
	PlayerInfo    GameRoomEnterViewPlayerInfo    `json:"playerInfo"`
}

// GameRoomEnterBody ?
type GameRoomEnterBody struct {
	BodyBase
	Data GameRoomEnterViewModel `json:"data"`
}

// PlayerWalletViewModel ?
type PlayerWalletViewModel struct {
	Amount             json.Number `json:"amount"`
	MerchantWalletMode int         `json:"merchantWalletMode"`
}

// PlayerWalletBody ?
type PlayerWalletBody struct {
	BodyBase
	Data PlayerWalletViewModel `json:"data"`
}
