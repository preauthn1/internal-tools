package hallutil

import (
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// HallVerProvider ?
type HallVerProvider interface {
	GetVer(ctx *module.MouduleCtx, clt *network.ConnInfo, areaID int64) *netproto.HallVerInfo
	GetShowGameFlagByArea(areaID int64) []*netproto.AreaShowGameFlag
	GetGameServerListByArea(areaid int64) []*netproto.PlatformGameServer
}

// HallVerLoader 平台版本號加載器
type HallVerLoader interface {
	// 加載平台版本號
	LoadHallVer()
}
