package model

// GameLaunchViewModel ?
type GameLaunchViewModel struct {
	PlayerID         int32  `json:"playerID"`
	PlayerName       string `json:"playerName"`
	Lang             string `json:"lang"`
	ServerID         int32  `json:"serverID"`
	ServerConnection string `json:"serverConnection"`
	AccessToken      string `json:"access_Token"`
}
