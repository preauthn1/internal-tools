//==================================================================
//文件名： ttzconfig.go
//说明：   推筒子读配置文件
//修改记录：
//==================================================================
package ttz

import (
	"example.invalid/ServerTeam/Servers/pkg/config"
)

//基础配置
type ttzCfg struct {
	dataInt map[string]int    //基础配置信息
	dataStr map[string]string //基础配置信息
}

//配置数据初始化，读取游戏配置文件
func (gameCfg *ttzCfg) init(cfgs *config.Group) {
	cfg := cfgs.GetConfig("gameserver")

	gameCfg.dataInt["Rob"] = cfg.GetConfigInt("gamelogic", "Rob")
	gameCfg.dataInt["Bet"] = cfg.GetConfigInt("gamelogic", "Bet")
	gameCfg.dataInt["Light"] = cfg.GetConfigInt("gamelogic", "Light")
	gameCfg.dataInt["AniTimeEndGame"] = cfg.GetConfigInt("gamelogic", "AnimationEndGame")
	gameCfg.dataInt["AniTimeDealOne"] = cfg.GetConfigInt("gamelogic", "DealOneTime")
	gameCfg.dataInt["AniTimeDealTwo"] = cfg.GetConfigInt("gamelogic", "DealTwoTime")
	gameCfg.dataInt["AniTimeBeginGame"] = cfg.GetConfigInt("gamelogic", "BeginGameTime")
	gameCfg.dataInt["PokerType"] = cfg.GetConfigInt("gamelogic", "PokerType")

	gameCfg.dataStr["AniNameGameEnd"] = cfg.GetConfigValue("gamelogic", "GameEnd")
	gameCfg.dataStr["AniNameDealOne"] = cfg.GetConfigValue("gamelogic", "DealOne")
	gameCfg.dataStr["AniNameDealTwo"] = cfg.GetConfigValue("gamelogic", "DealTwo")
	gameCfg.dataStr["AniNameBeginGame"] = cfg.GetConfigValue("gamelogic", "BeginGame")

	gameCfg.dataInt["ReadyWaitTime"] = cfg.GetConfigInt("gameserver", "readywaittime")
	gameCfg.dataInt["BaseScore"] = cfg.GetConfigInt("gameserver", "basescore")

}

//获取配置文件整形值
func (gameCfg *ttzCfg) GetValueInt(key string) int {
	return gameCfg.dataInt[key]
}

//获取配置文件字符值
func (gameCfg *ttzCfg) GetValueStr(key string) string {
	return gameCfg.dataStr[key]
}

//创建配置文件缓存
func NewTtzCfg(cfg *config.Group) *ttzCfg {
	d := new(ttzCfg)
	d.dataInt = make(map[string]int)
	d.dataStr = make(map[string]string)
	d.init(cfg)
	return d
}
