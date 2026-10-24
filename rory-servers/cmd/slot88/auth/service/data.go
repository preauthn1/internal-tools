package service

import (
	"example.invalid/dgrijalva/jwt-go"
)

// Info save user & server
type Info struct {
	RandomID   string     `json:"randomID"`
	Expire     string     `json:"expiredTime"`
	PlayerInfo PlayerInfo `json:"playerInfo"`
	ServerInfo ServerInfo `json:"serverInfo"`
	jwt.StandardClaims
}

// PlayerInfo user info
type PlayerInfo struct {
	PlayerID    int32  `json:"playerID"`
	PlayerName  string `json:"playerName"`
	Language    string `json:"lang"`
	AccessToken string `json:"access_Token"`
}

// ServerInfo info
type ServerInfo struct {
	GameID           int32  `json:"gameID,omitempty"`
	ServerID         int32  `json:"serverID,omitempty"`
	ServerConnection string `json:"serverConnection,omitempty"`
	ConnType         string `json:"connType,omitempty"`
	GameVendorType   int32  `json:"gameVendorType,omitempty"`
}
