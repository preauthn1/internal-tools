package model

// BodyBase ?
type BodyBase struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// MapDataBody ?
type MapDataBody struct {
	BodyBase
	Data map[string]interface{} `json:"data"`
}
