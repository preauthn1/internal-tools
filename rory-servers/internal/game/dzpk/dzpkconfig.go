//==================================================================
//文件名： dzpkconfig.go
//说明：   德州扑克配置文件
//修改记录：
//==================================================================
package dzpk

import (
	"example.invalid/ServerTeam/Servers/pkg/config"
)

type dzpkTimeOutCfg struct {
	tc map[gameProgress]int32 //阶段超时配置
}

//获取游戏阶段对应的等待超时时间
func (dzpk *dzpkTimeOutCfg) GetTime(gp gameProgress) int32 {
	return dzpk.tc[gp]
}

//初始化超时配置
func dzpkTimeOutInit() *dzpkTimeOutCfg {
	dzpkt := new(dzpkTimeOutCfg)
	dzpkt.tc = map[gameProgress]int32{
		GAMEPROGRESS_BEGIN:   0,
		GAMEPROGRESS_DEAL:    0,
		GAMEPROGRESS_PREFLOP: 10,
		GAMEPROGRESS_FLOP:    10,
		GAMEPROGRESS_TURN:    10,
		GAMEPROGRESS_RIVER:   60,
		GAMEPROGRESS_END:     10,
	}
	return dzpkt
}

//基础配置
type dzpkCfg struct {
	dataInt map[string]int //基础配置信息
	dataStr map[string]string
}

//配置数据初始化，读取游戏配置文件
func (dzcfg *dzpkCfg) init(cfgs *config.Group) {

	cfgSys := cfgs.GetConfig("gameserver")
	dzcfg.dataInt["ReadyTime"] = cfgSys.GetConfigInt("gameserver", "readywaittime")
	dzcfg.dataInt["BaseScore"] = cfgSys.GetConfigInt("gameserver", "basescore")
	dzcfg.dataInt["ActionTime"] = cfgSys.GetConfigInt("gameserver", "opttime")
	dzcfg.dataInt["GameEndTime"] = cfgSys.GetConfigInt("gameserver", "actiongameend")
	dzcfg.dataInt["PokerType"] = cfgSys.GetConfigInt("gameserver", "autotpokertype")
	dzcfg.dataInt["GM"] = cfgSys.GetConfigInt("gameserver", "gm")
}

//获取配置文件整形值
func (dzcfg *dzpkCfg) GetValueInt(key string) int {
	return dzcfg.dataInt[key]
}

//获取配置文件字符值
func (dzcfg *dzpkCfg) GetValueStr(key string) string {
	return dzcfg.dataStr[key]
}

//创建配置文件缓存
func NewDzpkCfg(cfg *config.Group) *dzpkCfg {
	d := new(dzpkCfg)
	d.dataInt = make(map[string]int)
	d.dataStr = make(map[string]string)
	d.init(cfg)
	return d
}
