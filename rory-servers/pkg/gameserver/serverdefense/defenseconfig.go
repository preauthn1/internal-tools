package serverdefense

import "example.invalid/ServerTeam/Servers/pkg/config"

// 配置節名稱
const (
	ConfigSectionNameDefense = "defense" // 防護方案配置節名稱
)

// DefenseConfig 防護配置
type DefenseConfig struct {
	BlockIPTime       int    // 封IP時間
	MaxConn           int    // 最大連接數
	ConnTimesLimit    string // 連接次數
	SendDataLimint    string // 發送數據包大小
	SendMsgLimit      string // 發送消息數量
	ProtoLimitDefault string // 協議發送頻率默認限制
	ProtoLimitSpecial string // 特殊協議發送頻率
	WhiteIPList       string // ip白名單
	ServerDesc        string // 服務器信息
	Enable            bool   // 是否開啟
}

// ParseDefenseConfig 解析防護協議配置
func ParseDefenseConfig(cfg *config.Config) *DefenseConfig {
	defCfg := new(DefenseConfig)
	blockiptime := cfg.GetConfigInt(ConfigSectionNameDefense, "blocktime")
	whiteiplist := cfg.GetConfigValue(ConfigSectionNameDefense, "whileiplist")

	ccMaxconn := cfg.GetConfigInt(ConfigSectionNameDefense, "cc_ipmaxconn")
	ccConntimes := cfg.GetConfigValue(ConfigSectionNameDefense, "cc_conntimes")
	ccSenddata := cfg.GetConfigValue(ConfigSectionNameDefense, "cc_senddata")
	ccSendmsg := cfg.GetConfigValue(ConfigSectionNameDefense, "cc_sendmsg")

	protoDefault := cfg.GetConfigValue(ConfigSectionNameDefense, "proto_default")
	protoSpecial := cfg.GetConfigValue(ConfigSectionNameDefense, "proto_special")

	serverDesc := cfg.GetConfigValue(ConfigSectionNameDefense, "serverdesc")
	enable := cfg.GetConfigValue(ConfigSectionNameDefense, "enable")

	defCfg.BlockIPTime = blockiptime
	defCfg.MaxConn = ccMaxconn
	defCfg.ConnTimesLimit = ccConntimes
	defCfg.SendDataLimint = ccSenddata
	defCfg.SendMsgLimit = ccSendmsg
	defCfg.ProtoLimitDefault = protoDefault
	defCfg.ProtoLimitSpecial = protoSpecial
	defCfg.WhiteIPList = whiteiplist
	defCfg.ServerDesc = serverDesc

	if enable == "" {
		defCfg.Enable = true
	} else {
		defCfg.Enable = (enable == "1")
	}

	return defCfg
}
