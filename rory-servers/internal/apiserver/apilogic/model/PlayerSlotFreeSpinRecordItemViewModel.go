package model

import (
	"encoding/json"
)

// PlayerSlotFreeSpinRecordItemViewModel ?
type PlayerSlotFreeSpinRecordItemViewModel struct {
	PlayerID      int32       `json:"playerID"`
	GameID        int32       `json:"gameID"`
	FlagID        int32       `json:"flagID"`
	PatternID     int32       `json:"patternID"`
	OriginID      int64       `json:"originID"`
	FreeSpinGroup int32       `json:"freeSpinGroup"`
	ChooseRate    int32       `json:"chooseRate"`
	BetAmount     json.Number `json:"betAmount"`
	Times         int32       `json:"times"`
	WinAmount     json.Number `json:"winAmount"`
	WinTimes      int32       `json:"winTimes"`
	Data          string      `json:"data"`
}

// PlayerSlotFreeSpinRecordReadBody ?
type PlayerSlotFreeSpinRecordReadBody struct {
	BodyBase
	Data struct {
		SlotFreeSpinRecord []PlayerSlotFreeSpinRecordItemViewModel `json:"slotFreeSpinRecord"`
	}
}
