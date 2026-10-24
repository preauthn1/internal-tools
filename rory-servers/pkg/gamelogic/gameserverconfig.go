package gamelogic

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/module/logicmodule"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// 寫分方式
const (
	ScoreUpdateModeNow         = 0 // 寫分方式, 即時寫分
	ScoreUpdateModeLogout      = 1 // 寫分方式, 退房間寫分
	ScoreUpdateModePiSai       = 2 // 寫分方式, 比賽積分
	ScoreUpdateModeConditional = 3 // 寫分方式, 滿足設定條件才寫分
)

// 發牌方式
const (
	PokerTypeRandom  = 1 // 隨機發牌
	PokerTypeConfig  = 2 // 發牌器發牌
	PokerTypeShuffle = 3 // 不洗牌
)

// 牌局記錄方式
const (
	DealLogTypeNone          int32 = iota // 不記錄
	DealLogTypeFile                       // 文件記錄
	DealLogTypeDB                         // 數據庫記錄
	DealLogTypeFish                       // 捕魚紀錄
	DealLogTypeConditionalDB              // 滿足設定條件才更新數據庫紀錄
	DealLogTypeConsole                    // console 紀錄
)

// GameServerConfig 服務器配置信息
type GameServerConfig struct {
	ProxyIP          string // 代理服務器 ip
	ServerAddr       string // gameserver 偵聽地址
	BServerAddr      string // B服 gameserver ip
	GRPCAddr         string // GRPC 監聽位置
	HallServerAddr   string // 大廳地址,以逗號分隔
	DBServerAddr     string // dbserver地址
	NotifyServerAddr string // 通知程序地址

	ServerID              int32 // 服務器ID
	ScoreMode             int32 // 寫分方式
	MaxUser               int32 // 人數上限
	RegHallServerInterval int32 // 註冊廣場時間間隔

	GameID            int32  // 遊戲ID
	GameName          string // 遊戲名稱
	ServerName        string // 服務器名稱
	Flag              string // 服務器標記
	MinPlayer         int32  // 最小人數
	MaxPlayer         int32  // 最大人數
	TableCount        int32  // 桌子數
	SeatCount         int32  // 椅子數
	LoginMoney        int32  // 入場金幣數
	TaxScale          int32  // 抽水百分比
	TaxAddPool        int32  // 抽取稅收的百分比放入獎金池
	ReadyWaitTime     int32  // 準備超時時間
	BaseScore         int32  // 底分
	PokerType         int32  // 發牌方式
	MatchMakingTime   int32  // 未滿遊戲最小人數的等待時間
	MatchMakingNumber int32  // 開始倒數時間的最低人數
	Kickoutmoney      int32  // 達到此积分，強制退出房間
	MaxOdd            int32
	HTTPServerAddr    string // http服務地址
	HTTPVCode         string // http校驗碼

	DealLogType           int32 // 牌局记录方式
	MaxGameCountSameTable int32 // 同一桌最大游戏次数
	GameTimeOut           int32 // 當遊戲開始後超過N分鐘未結束該局就踢出此張桌子所有人

	Single              int32 // 直接進入取第一個服務器進入, 0=>顯示房間列表, 1=>第一個服務器進入
	HostAndBackup       int32 // AB服. 0主機 1備援轉正主機接手中 2備援
	GameType            int32 // 0真金 1比賽 2休閒
	SuperToolCgf        SuperToolConfig
	SettleScore         int64 // 大於SettleScore分時結算, deallog=4 或 scoremode=3 吃這個設定
	SettleRound         int   // 大於SettleRound局時結算, deallog=4 或 scoremode=3 吃這個設定
	SaveHistoryInterval int32 // 更新歷史紀錄間隔時間(毫秒)
}

// SuperToolConfig 超端設定檔
type SuperToolConfig struct {
	SuperToolModeType int32 // 超端: 1配桌型
	SyncInterval      int32 // 放分更新時間,單位秒
}

// GetServerIPPort 獲取服務器ip和端口
func (cfg *GameServerConfig) GetServerIPPort() (string, int32) {
	ip, port := util.ParseIPAndPort(cfg.ServerAddr)
	if cfg.ProxyIP != "" {
		ip = cfg.ProxyIP
	}
	return ip, int32(port)
}

// GetGameInfo 獲取遊戲信息
func (cfg *GameServerConfig) GetGameInfo() *Game {
	game := NewGame()
	game.GameID = cfg.GameID
	game.GameName = cfg.GameName
	game.MinPlayer = cfg.MinPlayer
	game.MaxPlayer = cfg.MaxPlayer

	return game
}

// GetHTTPServerCfg 獲取Http服務器配置
func (cfg *GameServerConfig) GetHTTPServerCfg() *logicmodule.HTTPServerConfig {
	httpCfg := new(logicmodule.HTTPServerConfig)
	httpCfg.ServerAddr = cfg.HTTPServerAddr
	httpCfg.VCode = cfg.HTTPVCode

	return httpCfg
}

// Config Names
const (
	ConfigSectionNameServer     = "server"
	ConfigSectionNameGameServer = "gameserver"
	ConfigSectionNameGameQueue  = "gamequeue"
	ConfigSectionNameHTTP       = "http"
	ConfigSectionNameGameInfo   = "gameinfo"

	ConfigFileNameGameServer  = "gameserver"
	ConfigFileNameServerInfo  = "serverinfo"
	ConfigFileNameLogicConfig = "logicconfig"
	ConfigFileNameSuperTool   = "supertool"
)

// GetGameServerConfig 獲取服務器配置
func GetGameServerConfig(serverInfo, cfg *config.Config) *GameServerConfig {
	gscfg := new(GameServerConfig)

	proxyip := serverInfo.GetConfigValue(ConfigSectionNameServer, "proxyip")
	addr := serverInfo.GetConfigValue(ConfigSectionNameServer, "serveraddr")
	grpcaddr := serverInfo.GetConfigValue(ConfigSectionNameServer, "grpcaddr")
	dbaddr := serverInfo.GetConfigValue(ConfigSectionNameServer, "dbserveraddr")
	nsaddr := serverInfo.GetConfigValue(ConfigSectionNameServer, "nsserveraddr")
	hsaddr := serverInfo.GetConfigValue(ConfigSectionNameServer, "hallsvraddr")

	servername := serverInfo.GetConfigValue(ConfigSectionNameGameServer, "servername")
	gameid := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "gameid"))
	gamename := cfg.GetConfigValue(ConfigSectionNameGameServer, "gamename")
	minplayer := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "minplayer"))
	maxplayer := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "maxplayer"))
	tablecount := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "tablecount"))
	seatcount := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "seatcount"))
	flag := cfg.GetConfigValue(ConfigSectionNameGameServer, "flag")
	loginmoney := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "loginmoney"))
	maxuser := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "maxuser"))
	serverid := int32(serverInfo.GetConfigInt(ConfigSectionNameGameServer, "serverid"))
	scoremode := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "scoremode"))
	taxscale := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "taxscale"))
	taxaddpool := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "taxaddpool"))
	readywaittime := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "readywaittime"))
	reghallserverinterval := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "reghallserverinterval"))
	basescore := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "basescore"))
	pokertype := int32(cfg.GetConfigInt(ConfigFileNameGameServer, "pokertype"))
	matchmakingtime := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "matchmakingtime"))
	matchmakingnumber := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "matchmakingnumber"))
	kickoutmoney := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "kickoutmoney"))
	maxodd := int32(cfg.GetConfigInt(ConfigSectionNameGameServer, "maxodd"))

	deallog := int32(cfg.GetConfigInt(ConfigFileNameGameServer, "deallog"))
	settlescore := cfg.GetConfigInt64(ConfigSectionNameGameServer, "settlescore")
	settleround := cfg.GetConfigInt(ConfigSectionNameGameServer, "settleround")
	savehistoryinterval := cfg.GetConfigInt32(ConfigSectionNameGameServer, "savehistoryinterval")
	maxgamecountsametable := int32(cfg.GetConfigInt(ConfigFileNameGameServer, "maxgamecountsametable"))
	gametimeout := int32(cfg.GetConfigInt(ConfigFileNameGameServer, "gametimeout"))
	if gametimeout == int32(0) {
		gametimeout = 30
	}

	gscfg.ProxyIP = proxyip
	httpServerAddr := serverInfo.GetConfigValue(ConfigSectionNameHTTP, "addr")
	httpVCode := serverInfo.GetConfigValue(ConfigSectionNameHTTP, "vcode")

	single := int32(cfg.GetConfigInt(ConfigFileNameGameServer, "single"))
	hostandbackup := serverInfo.GetConfigInt32(ConfigSectionNameGameServer, "hostandbackup")
	gametype := serverInfo.GetConfigInt32(ConfigFileNameGameServer, "gametype")

	superToolModeType := serverInfo.GetConfigInt32(ConfigFileNameSuperTool, "modetype")
	superToolSyncInterval := serverInfo.GetConfigInt32(ConfigFileNameSuperTool, "syncinterval")

	gscfg.ServerAddr = addr
	gscfg.GRPCAddr = grpcaddr
	gscfg.HallServerAddr = hsaddr
	gscfg.DBServerAddr = dbaddr
	gscfg.NotifyServerAddr = nsaddr

	gscfg.GameID = gameid
	gscfg.GameName = gamename
	gscfg.ServerID = serverid
	gscfg.ServerName = servername
	gscfg.ScoreMode = scoremode
	gscfg.Flag = flag
	gscfg.MinPlayer = minplayer
	gscfg.MaxPlayer = maxplayer
	gscfg.TableCount = tablecount
	gscfg.SeatCount = seatcount
	gscfg.LoginMoney = loginmoney
	gscfg.MaxUser = maxuser
	gscfg.TaxScale = taxscale
	gscfg.TaxAddPool = taxaddpool
	gscfg.ReadyWaitTime = readywaittime
	gscfg.RegHallServerInterval = reghallserverinterval
	gscfg.BaseScore = basescore
	gscfg.PokerType = pokertype
	gscfg.DealLogType = deallog
	gscfg.SettleScore = settlescore
	gscfg.SettleRound = settleround
	gscfg.SaveHistoryInterval = savehistoryinterval
	gscfg.MaxGameCountSameTable = maxgamecountsametable
	gscfg.GameTimeOut = gametimeout
	gscfg.MatchMakingTime = matchmakingtime
	gscfg.MatchMakingNumber = matchmakingnumber
	gscfg.Kickoutmoney = kickoutmoney
	gscfg.MaxOdd = maxodd
	gscfg.HTTPServerAddr = httpServerAddr
	gscfg.HTTPVCode = httpVCode

	gscfg.Single = single
	gscfg.HostAndBackup = hostandbackup
	gscfg.GameType = gametype
	gscfg.SuperToolCgf.SuperToolModeType = superToolModeType
	gscfg.SuperToolCgf.SyncInterval = superToolSyncInterval

	return gscfg
}

func (cfg *GameServerConfig) String() string {
	return fmt.Sprintf("gameid:%d gamename:%s serverid:%d,servername:%s,listen:%s,flag:%s,scoremode:%d,loginmoney:%d,base score:%d", cfg.GameID, cfg.GameName, cfg.ServerID, cfg.ServerName, cfg.ServerAddr, cfg.Flag, cfg.ScoreMode, cfg.LoginMoney, cfg.BaseScore)
}
