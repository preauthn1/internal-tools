package service

import "example.invalid/ServerTeam/Servers/pkg/module"

// ?
var (
	GameResultService    *GameResult
	SystemSettingService *SystemSetting
	ServerSettingService *ServerSetting
	GamePlayService      *GamePlay
	GameSlotService      *GameSlot
	SuperToolPoolService *SuperToolPool
)

// InitService ?
func InitService(ctx *module.MouduleCtx) {
	GameResultService = NewGameResult(ctx)
	SystemSettingService = NewSystemSetting(ctx)
	ServerSettingService = NewServerSetting(ctx)
	GamePlayService = NewGamePlay(ctx)
	GameSlotService = NewGameSlot(ctx)
	SuperToolPoolService = NewSuperToolPool(ctx)
}
