package model

// ServerStartProcessViewModel ?
type ServerStartProcessViewModel struct {
	GameID     int32 `json:"gameID"`
	GameStatus int32 `json:"gameStatus"`
}

// ServerStartProcessBody ?
type ServerStartProcessBody struct {
	BodyBase
	Data ServerStartProcessViewModel `json:"data"`
}
