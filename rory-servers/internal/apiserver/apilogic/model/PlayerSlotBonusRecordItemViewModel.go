package model

import (
	"encoding/json"
)

// PlayerSlotBonusRecordItemViewModel ?
type PlayerSlotBonusRecordItemViewModel struct {
	PlayerID  int32       `json:"playerID"`
	GameID    int32       `json:"gameID"`
	FlagID    int32       `json:"flagID"`
	PatternID int32       `json:"patternID"`
	OriginID  int32       `json:"originID"`
	TurnID    int32       `json:"turnID"`
	Location  int32       `json:"location"`
	BetAmount json.Number `json:"betAmount"`
	BetRate   int32       `json:"betRate"`
	Data      string      `json:"data"`
}

// PlayerSlotBonusRecordReadBody ?
type PlayerSlotBonusRecordReadBody struct {
	BodyBase
	Data struct {
		SlotBonusRecord []PlayerSlotBonusRecordItemViewModel `json:"slotBonusRecord"`
	}
}
