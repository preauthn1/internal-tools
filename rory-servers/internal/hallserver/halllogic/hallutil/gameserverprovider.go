package hallutil

import (
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// GameServerProvider ?
type GameServerProvider interface {
	GetGameServerList() []*netproto.GameServerInfo                                                       // 獲取服務器列表
	GetGameServerAddrByFlag(gameid int32, flag string) []*netproto.GameServerAddr                        // 根據flag獲取服務器地址
	GetGameServerListByAreaID(areaid int64) []*netproto.GameServerInfo                                   // 取得渠道全部可用的遊戲服務器
	GetGameServerAddrByAreaIDAndFlag(areaid int64, gameid int32, flag string) []*netproto.GameServerAddr // 根據areaid及flag獲取服務器地址
}
