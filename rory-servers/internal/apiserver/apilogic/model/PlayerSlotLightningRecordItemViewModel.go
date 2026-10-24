package model

// PlayerSlotLightningRecordItemViewModel ?
type PlayerSlotLightningRecordItemViewModel struct {
	PlayerID  int32  `json:"playerID"`
	GameID    int32  `json:"gameID"`
	FlagID    int32  `json:"flagID"`
	PatternID int32  `json:"patternID"`
	Data      string `json:"data"`
}

// PlayerSlotLightningRecordReadBody ?
type PlayerSlotLightningRecordReadBody struct {
	BodyBase
	Data struct {
		SlotLightningRecord []PlayerSlotLightningRecordItemViewModel `json:"slotLightingRecord"`
	}
}
