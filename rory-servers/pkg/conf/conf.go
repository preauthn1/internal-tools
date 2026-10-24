package conf

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"example.invalid/spf13/cobra"

	"example.invalid/ServerTeam/Servers/pkg/util"
)

// Platform 平台
type Platform string

const (
	// PlatformNone 無平台(預設)
	PlatformNone Platform = "NONE"
	// PlatformPhnum 金邊平台
	PlatformPhnum Platform = "PHNUM"
)

var (
	// CLogLevel 控制台日志等级
	CLogLevel int = 3
	// FLogLevel 檔案日志等級
	FLogLevel int = 1
	// ConfigDirName 配置文件目录名
	ConfigDirName = "config"
	// EVTCHANSIZE 事件通道大小
	EVTCHANSIZE = 2000
	// MSGPACKERBUF 消息缓冲区大小
	MSGPACKERBUF int32 = 1024 * 1000
	// READBUFLEN 接收消息缓冲区大小
	READBUFLEN int = 1024 * 1000
	// DefaultEventThrottle 默认的事件并发数
	DefaultEventThrottle int = 1000
	// ClearLog 是否使用乾淨的Log
	ClearLog bool = false
	// Uptime 不知道幹嘛的
	Uptime bool = false
	// ServicePlatform 營運廠商
	ServicePlatform Platform = PlatformNone
)

// 金邊專用參數
var (
	// LianyunID 聯運ID (金邊用)
	LianyunID int = 0
	RealAddr  string
)

// TCPSERVERNAME
var (
	TCPSERVERNAMEHALLCLIENT     = "hallTCP_Client"   // 大廳面向客戶端TCP服務器名稱
	TCPSERVERNAMEHALLGAMESERVER = "hallTCP_GSMANAGE" // 大廳遊戲服務器管理TCP服務器名稱
	TCPSERVERNAMEROOMCLIENT     = "roomTcp_Client"   // 房間面向客戶端TCP服務器名稱

	TCPSERVERNAMEDBSERVER = "dbserverTCP" // dbserver tcp服務器名稱

	TCPSERVERCONNHALLCLIENT     = 10000 // 大廳客戶端連接最大連接數
	TCPSERVERCONNHALLGAMESERVER = 2000  // 大廳遊戲服務器連接最大連接數
	TCPSERVERCONNROOMCLIENT     = 2000  // 遊戲客戶端連接最大連接數

	TCPSERVERCONNDBSERVER = 2000 // dbserverTCP連接數

	TCPSERVERCONNNOTIFYSERVER = 2000              // notifyserver tcp連接數
	TCPSERVERNAMENOTIFYSERVER = "notifyserverTCP" // notifyserver tcp服務器名稱
)

// 運行測試案例時使用的一些配置
var (
	TestMode      bool = false // 搭配視別是否為執行測試案例
	AnalysisMode  bool = false // 是否為分析模式
	PrintMessage  bool = false //
	TimerFastMode int  = 1     // 計時器的快轉倍數
)

// Password 在Cmd指令內所要用到的操控密碼
var Password string = "ㄎㄢˋㄙㄢㄒㄧㄠˇ"

func init() {
	type version struct {
		BuildTime string
		GitTag    string
		GitHash   string
		BuildNum  string
	}
	if BuildNum != "{BuildNum}" {
		v := version{BuildTime, GitTag, GitHash, BuildNum}
		b, _ := json.Marshal(v)
		util.FilePutContentsNA("version.txt", string(b))
		fmt.Fprintf(os.Stdout, "version : (%s) %s\n", GitTag, GitHash)
		fmt.Fprintf(os.Stdout, "buildtime : %s\n", BuildTime)
		fmt.Fprintf(os.Stdout, "buildNum : %s\n", BuildNum)
		FLogLevel = 1 // 正式版本的 filelog level 為 info
	} else {
		FLogLevel = 1 // 開發版本的 filelog level 為 debug
	}

	if flag.Lookup("test.v") != nil {
		TestMode = true
	}

	// 使 cobra command 在 windows terminal 不會crash
	cobra.MousetrapHelpText = ""
}

// 服務器類型
const (
	ServerTypeGameServer = iota // gameserver
	ServerTypeHallServer        // hallserver
)

// 服務器狀態
const (
	ServerStatusRuning  int32 = 1 // 運行中
	ServerStatusStoping int32 = 0 // 停服中
	ServerStatusClose   int32 = 2 // 已關閉
	ServerStatusRecover int32 = 3 // 維護中
)

// Version Information
var (
	BuildTime = "{BuildTime}"
	GitTag    = "{GitTag}"
	GitHash   = "{GitHash}"
	BuildNum  = "{BuildNum}"
)

// Version 取得版本資訊
func Version() string {
	return fmt.Sprintf("=== 版本資訊 ===\n當前版本：%s GitHash：%s\n建置時間：%s\n建置編號：%s", GitTag, GitHash, BuildTime, BuildNum)
}
