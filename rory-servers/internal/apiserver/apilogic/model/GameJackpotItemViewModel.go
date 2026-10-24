package model

import (
	"encoding/json"
)

// GameJackpotItemViewModel ?
type GameJackpotItemViewModel struct {
	GameID            int32       `json:"gameID"`
	TableNo           int32       `json:"tableNo"`
	FlagID            int32       `json:"flagID"`
	LevelID           int32       `json:"levelID"`
	Contribution      json.Number `json:"contribution"`
	TotalContribution json.Number `json:"totalContribution"`
	AwardTotal        json.Number `json:"totalAward"`
	Type              bool        `json:"type"`
}

// SlotJackpotBody ?
type SlotJackpotBody struct {
	BodyBase
	Data struct {
		Jackpots []GameJackpotItemViewModel `json:"jackpots"`
	}
}
