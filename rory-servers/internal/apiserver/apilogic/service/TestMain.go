package service

import (
	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/model"
	"example.invalid/ServerTeam/Servers/pkg/httputil"
)

// GameLaunchBody ?
type GameLaunchBody struct {
	model.BodyBase
	Data struct {
		model.GameLaunchViewModel `json:"info"`
	} `json:"data"`
}

func newBase() *Base {
	base := &Base{}
	base.schema = "http"
	base.backendAPI = "example.invalid"
	base.apiCertificate = "851b-222e5f24024b"
	return base
}

// InitTestService 初始化測試用的Service
func InitTestService() {
	httputil.DebugMode(true)
	GameResultService = &GameResult{}
	GameResultService.Base = newBase()
	SystemSettingService = &SystemSetting{}
	SystemSettingService.Base = newBase()
	ServerSettingService = &ServerSetting{}
	ServerSettingService.Base = newBase()
	GamePlayService = &GamePlay{}
	GamePlayService.Base = newBase()
	GameSlotService = &GameSlot{}
	GameSlotService.Base = newBase()
	SuperToolPoolService = &SuperToolPool{}
	SuperToolPoolService.Base = newBase()
}

// GetTestAccessToken ?
func GetTestAccessToken() (int32, string) {
	base := newBase()
	model := new(GameLaunchBody)

	// 取得 token
	api := base.buildBackendAPI("/fake-test/player/token/generate", map[string]interface{}{
		"PlayerName": "alex001",
	})
	header := base.buildHeader(base.apiCertificate)
	request := new(httputil.WaitGroupRequest)
	request.AddTask(httputil.GET, api, nil, header, model)
	request.Send()
	if request.IsAllDone() {
		return model.Data.PlayerID, model.Data.AccessToken
	}
	return 0, ""
}
