package model

// PlayerSlotUserAssignedModeItemViewModel ?
type PlayerSlotUserAssignedModeItemViewModel struct {
	PlayerID   int32    `json:"playerID"`
	GameID     int32    `json:"gameID"`
	FlagID     int32    `json:"flagID"`
	PatternID  int32    `json:"patternID"`
	CreateTime JSONTime `json:"createTime"`
	Mode       int32    `json:"mode"`
	ExpireTime JSONTime `json:"expireTime"`
}

// PlayerSlotUserAssignedModeReadBody ?
type PlayerSlotUserAssignedModeReadBody struct {
	BodyBase
	Data struct {
		SlotUserAssignedModeItem []PlayerSlotUserAssignedModeItemViewModel `json:"slotUserAssignedMode"`
	}
}
