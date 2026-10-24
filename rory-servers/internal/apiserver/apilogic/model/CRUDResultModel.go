package model

// CRUDResultModel ?
type CRUDResultModel struct {
	IsSuccess    bool   `json:"isSuccess"`
	Message      string `json:"message"`
	ResultID     int    `json:"resultId"`
	ResultObject string `json:"resultObject"`
}

// PlayerSlotProbRecordUpdateBody ?
type PlayerSlotProbRecordUpdateBody struct {
	BodyBase
	Data struct {
		SlotProbRecord []CRUDResultModel `json:"slotProbRecord"`
	}
}
