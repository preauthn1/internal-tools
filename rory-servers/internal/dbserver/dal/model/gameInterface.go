package model

// GameLaunchResponseBody 三方遊戲登入 response body
type GameLaunchResponseBody struct {
	BodyBase
	Data struct {
		RedirectURL    string `json:"redirectUrl"`
		TransferAmount int    `json:"transferAmount"`
	} `json:"data"`
}
