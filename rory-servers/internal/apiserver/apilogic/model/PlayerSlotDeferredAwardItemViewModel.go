package model

// PlayerSlotDeferredAwardItemViewModel ?
type PlayerSlotDeferredAwardItemViewModel struct {
	PlayerID   int32    `json:"playerID"`
	GameID     int32    `json:"gameID"`
	FlagID     int32    `json:"flagID"`
	AwardType  string   `json:"awardType"`
	PatternID  int32    `json:"patternID"`
	OriginID   int32    `json:"originID"`
	Data       string   `json:"data"`
	CreateTime JSONTime `json:"createTime"`
}

// PlayerSlotDeferredAwardReadBody ?
type PlayerSlotDeferredAwardReadBody struct {
	BodyBase
	Data struct {
		SlotDeferredAward []PlayerSlotDeferredAwardItemViewModel `json:"playerSlotDeferredAward"`
	}
}

// PlayerSlotDeferredAwardDeleteBody ?
type PlayerSlotDeferredAwardDeleteBody struct {
	BodyBase
	Data PlayerSlotDeferredAwardItemViewModel
}
