package model

// BodyBase ResponseBody 基類
type BodyBase struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
