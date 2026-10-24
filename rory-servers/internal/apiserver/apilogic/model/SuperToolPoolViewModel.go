package model

import "encoding/json"

// SuperToolPoolViewModel ?
type SuperToolPoolViewModel struct {
	GameID     int32       `json:"gameID"`
	FlagID     int32       `json:"flagID"`
	WinLine    json.Number `json:"winLine"`
	PoolAmount json.Number `json:"poolAmount"`
	GameData   string      `json:"gameData"`
}

// SuperToolPoolBody ?
type SuperToolPoolBody struct {
	BodyBase
	Data SuperToolPoolViewModel `json:"data"`
}
