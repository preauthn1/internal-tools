package service

import "example.invalid/ServerTeam/Servers/pkg/module"


// PromotionEventService 活動相關接口服務
var PromotionEventService *PromotionEvent
// GameInterfaceService 遊戲相關接口服務
var GameInterfaceService *GameInterface

// InitService 初始化服務
func InitService(ctx *module.MouduleCtx) {
	PromotionEventService = NewPromotionEvent(ctx)
	GameInterfaceService = NewGameInterface(ctx)
}