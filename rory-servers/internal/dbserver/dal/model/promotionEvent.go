package model

// Event 活動資訊
type Event struct {
	EventID   int32    `json:"eventID"`
	EventName string   `json:"eventName"`
	EventType int32    `json:"eventType"`
	StartTime JSONTime `json:"startTime"`
	EndTime   JSONTime `json:"endTime"`
}

// EventList 獲取活動列表
type EventList struct {
	Events []Event `json:"events"`
}

// GetEventListBody response body of query /event
type GetEventListBody struct {
	BodyBase
	Data EventList `json:"data"`
}

// LuckyRouletteInfo 輪盤資訊
type LuckyRouletteInfo struct {
	TicketCount int32       `json:"ticketCount"`
	EventMsg    string      `json:"eventMsg"`
	TipsMsg     string      `json:"tipsMsg"`
	BonusItems  []BonusItem `json:"bonusItems"`
}

// BonusItem 獎項
type BonusItem struct {
	BonusID   int32  `json:"bonusID"`
	BonusName string `json:"bonusName"`
	Amount    int32  `json:"amount"`
	BonusType int32  `json:"bonusType"`
	ImageURL  string `json:"imageUrl"`
}

// GetEventInfoBody response body of query /events/{eventID}
type GetEventInfoBody struct {
	Data struct {
		BasicInfo         Event             `json:"basicInfo"`
		LuckyRouletteInfo LuckyRouletteInfo `json:"luckyRouletteInfo"`
	} `json:"data"`
}

// LuckyRouletteData 抽獎結果
type LuckyRouletteData struct {
	IsSuccess        bool   `json:"isSuccess"`
	ResultMsg        string `json:"resultMsg"`
	BonusID          int32  `json:"bonusID"`
	LotteryRemaining int32  `json:"lotteryRemaining"`
}

// PostLuckyRouletteResponseBody response body of post /events/lucky-roulette/request
type PostLuckyRouletteResponseBody struct {
	BodyBase
	Data LuckyRouletteData `json:"data"`
}
