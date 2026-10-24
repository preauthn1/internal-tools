package model

// PlayerSlotProbRecordItemViewModel ?
type PlayerSlotProbRecordItemViewModel struct {
	PlayerID  int32  `json:"playerID"`
	GameID    int32  `json:"gameID"`
	FlagID    int32  `json:"flagID"`
	PatternID int32  `json:"patternID"`
	OriginID  int32  `json:"originID"`
	Times     int32  `json:"times"`
	Data      string `json:"data"`
}

// PlayerSlotProbRecordReadBody ?
type PlayerSlotProbRecordReadBody struct {
	BodyBase
	Data struct {
		SlotProbRecord []PlayerSlotProbRecordItemViewModel `json:"slotProbRecord"`
	}
}
