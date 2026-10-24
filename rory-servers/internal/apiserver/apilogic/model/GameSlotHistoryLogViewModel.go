package model

import (
	"encoding/json"
)

// GameSlotHistoryLogInfo ?
type GameSlotHistoryLogInfo struct {
	RowNum      int32       `json:"rowNum"`
	PlayerID    int32       `json:"playerID"`
	LogID       int64       `json:"logID"`
	BetAmount   json.Number `json:"betAmount"`
	FinalAmount json.Number `json:"finalAmount"`
	GameTime    JSONTime    `json:"gameTime"`
}

// GameSlotHistoryLogViewModel ?
type GameSlotHistoryLogViewModel struct {
	Items       []GameSlotHistoryLogInfo `json:"items"`
	PageIndex   int32                    `json:"pageIndex"`
	PageSize    int32                    `json:"pageSize"`
	TotalCount  int32                    `json:"totalCount"`
	TotalPage   int32                    `json:"totalPage"`
	ExtraObject string                   `json:"extraObject"`
}

// GameSlotHistoryLogBody ?
type GameSlotHistoryLogBody struct {
	BodyBase
	Data GameSlotHistoryLogViewModel `json:"data"`
}
