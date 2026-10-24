package tfconfig

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// Configs
const (
	ConfigDirName                   = "config"
	ConfigTFCenterServerSectionName = "centerserver"
)

// TFConfig 转发配置
type TFConfig struct {
	ReportURL              string // 上报转发状态地址
	PortMapURL             string // 端口映射地址
	ReportInterval         int    // 上报状态时间间隔
	RefreshPortMapInterval int    // 刷新端口映射时间间隔
	CheckTFPortInterval    int    // 转发端口检查时间间隔
	ReportIP               int    // 建立通道時是否由TFSV上報連線來源IP
	NeedCheckKey           int    // 需要驗證key的埠號
	NeedCheckKeySourceIP   string // 需要驗證key的來源埠號
	ConnWatchInterval      int    // 每 n 秒偵測一次連線數量
	ConnWatchExceeded      int    // 每一分鐘超過50個連線數量的波動
	AlertBotAPI            string // 報警群
	WebSocketServerPort    int    // WebSocketServer port number
	Slot88WSPort           int    // slot88
}

// ConfigCtx config context
var ConfigCtx *TFConfig

func init() {
	ConfigCtx = loadTFConfig()
}

func loadTFConfig() *TFConfig {
	cfgs := config.GetConfigGroup(ConfigDirName)
	cfg := cfgs.GetConfig("tfserver")

	tfcfg := new(TFConfig)
	tfcfg.PortMapURL = cfg.GetConfigValue(ConfigTFCenterServerSectionName, "portmapurl")
	tfcfg.ReportURL = cfg.GetConfigValue(ConfigTFCenterServerSectionName, "reporturl")
	tfcfg.ReportInterval = cfg.GetConfigInt(ConfigTFCenterServerSectionName, "reportinterval")
	tfcfg.RefreshPortMapInterval = cfg.GetConfigInt(ConfigTFCenterServerSectionName, "refreshportmapinterval")
	tfcfg.CheckTFPortInterval = cfg.GetConfigInt(ConfigTFCenterServerSectionName, "checktfportinterval")
	tfcfg.ReportIP = cfg.GetConfigInt(ConfigTFCenterServerSectionName, "reportip")
	tfcfg.NeedCheckKey = cfg.GetConfigInt(ConfigTFCenterServerSectionName, "needcheckkey")
	tfcfg.NeedCheckKeySourceIP = cfg.GetConfigValue(ConfigTFCenterServerSectionName, "needcheckkeysourceip")

	tfcfg.ConnWatchInterval = cfg.GetConfigInt(ConfigTFCenterServerSectionName, "connwatchinterval")
	tfcfg.ConnWatchExceeded = cfg.GetConfigInt(ConfigTFCenterServerSectionName, "connwatchexceeded")
	tfcfg.AlertBotAPI = cfg.GetConfigValue(ConfigTFCenterServerSectionName, "alertbotapi")

	tfcfg.PortMapURL = fmt.Sprintf("http://%s", tfcfg.PortMapURL)
	tfcfg.ReportURL = fmt.Sprintf("http://%s", tfcfg.ReportURL)
	tfcfg.AlertBotAPI = fmt.Sprintf("http://%s", tfcfg.AlertBotAPI)
	tfcfg.WebSocketServerPort, _ = strconv.Atoi(cfg.GetConfigValue(ConfigTFCenterServerSectionName, "websocketserverport"))

	tfcfg.Slot88WSPort, _ = strconv.Atoi(cfg.GetConfigValue(ConfigTFCenterServerSectionName, "slot88wsport"))
	return tfcfg

}

// AlertMsg 訊息
func AlertMsg(m string) {
	s := url.QueryEscape(m)
	url := fmt.Sprintf(ConfigCtx.AlertBotAPI, s)
	mlog.Error(url)
	http.Get(url)
}
