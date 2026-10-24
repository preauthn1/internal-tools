package hallutil

import (
	"fmt"
	"strconv"
	"strings"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/config"
)

// HallServerConfig 大廳配置
type HallServerConfig struct {
	HallAddr            string                     // 大廳地址
	GSMgmAddr           string                     // 遊戲服務器註冊地址
	NotifyServerAddr    string                     // 協調服務器地址
	ServerID            int32                      // 服務器ID
	DBServerAddr        string                     // ?
	StrategyList        []*ServerAllocStrategyList // ?
	MaxUser             int32                      // ?
	RefreshCategoryTime int32                      // 遊戲分類刷新時間(秒)
}

// ServerAllocStrategyList 服务器分配策略, 百人桌, 配桌型
type ServerAllocStrategyList struct {
	GameIDList map[int32]*int32
	MaxUser    *int32
}

// GetHallServerConfig 獲取大廳配置
func GetHallServerConfig(cfg *config.Config) *HallServerConfig {
	halladdr := cfg.GetConfigValue("server", "halladdr")
	gsaddr := cfg.GetConfigValue("server", "gsaddr")
	dbaddr := cfg.GetConfigValue("dbserver", "addr")
	serverid := cfg.GetConfigInt("server", "serverid")
	nsaddr := cfg.GetConfigValue("server", "nsaddr")

	size := cfg.GetConfigInt("serverallocstrategy", "size")
	maxuser := cfg.GetConfigInt("serverallocstrategy", "maxuser")

	serverCfg := new(HallServerConfig)
	serverCfg.HallAddr = halladdr
	serverCfg.GSMgmAddr = gsaddr
	serverCfg.ServerID = int32(serverid)
	serverCfg.DBServerAddr = dbaddr
	serverCfg.NotifyServerAddr = nsaddr
	serverCfg.MaxUser = int32(maxuser)

	for i := 1; i <= size; i++ {
		cname := fmt.Sprintf("c%d", i)
		gameidval := cfg.GetConfigValue(cname, "gameidlist")
		gameidlist := strings.Split(gameidval, ",")
		strategy := new(ServerAllocStrategyList)
		strategy.GameIDList = make(map[int32]*int32, 0)
		for _, v := range gameidlist {
			num, _ := strconv.Atoi(v)
			strategy.GameIDList[int32(num)] = proto.Int32(int32(num))
		}
		strategy.MaxUser = proto.Int32(cfg.GetConfigInt32(cname, "maxuser"))
		serverCfg.StrategyList = append(serverCfg.StrategyList, strategy)
	}

	serverCfg.RefreshCategoryTime = int32(cfg.GetConfigInt("category", "reflash_time"))
	return serverCfg
}
