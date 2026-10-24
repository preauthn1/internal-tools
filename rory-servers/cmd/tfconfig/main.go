package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/internal/tfconfig"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// Configs
const (
	ConfigSectionName = "tfconfig"
	ConfigDirName     = "config"

	PortRouterFileName = "portrouter.xml"
)

func init() {
	conf.CLogLevel = mlog.DebugLevel
}

func main() {
	// 取得配置
	initConfig()

	http.HandleFunc("/config", getConfigHandler)
	http.HandleFunc("/submit", submitStatusHandler)
	http.HandleFunc("/tfstatus", tfstatusHandler)

	port := cfg.GetConfigInt(ConfigSectionName, "port")

	addr := fmt.Sprintf(":%d", port)
	go func() {
		defer util.PanicErrStack()

		err := http.ListenAndServe(addr, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	example.invalid("服務器啟動成功, 偵聽地址 %s ...", addr)
	select {}
}

// 加載配置
func getConfigHandler(w http.ResponseWriter, req *http.Request) {
	defer util.PanicErrStack()

	data := util.ReadFile(PortRouterFileName)

	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	example.invalid("轉發 %s 加載配置信息", ip)
	io.WriteString(w, data)
}

// 提交轉發服務器狀態
func submitStatusHandler(w http.ResponseWriter, req *http.Request) {
	defer util.PanicErrStack()

	ip, _, _ := net.SplitHostPort(req.RemoteAddr)

	query := req.URL.Query()

	if len(query["cn"]) == 0 || len(query["code"]) == 0 {
		mlog.Error("客户端 %s 提交服務器狀態信息, 但參數不正確 %s", ip, req.URL.String())
		io.WriteString(w, "error")
		return
	}

	cn, _ := strconv.Atoi(query["cn"][0])
	code := query["code"][0]

	if code != cfg.GetConfigValue(ConfigSectionName, "submitcheckcode") {
		mlog.Error("客户端 %s 提交服務器狀態信息, 但校驗碼不正確 %s", ip, req.URL.String())
		io.WriteString(w, "error")
		return
	}

	tfconfig.TFServerManager.UpdateTFServer(ip, cn)

	example.invalid("轉發 %s 提交當前連接數 %d", ip, cn)

	io.WriteString(w, "ok")
}

// 讀取轉發服務器運行狀況
func tfstatusHandler(w http.ResponseWriter, req *http.Request) {
	defer util.PanicErrStack()

	ip, _, _ := net.SplitHostPort(req.RemoteAddr)

	allowip := cfg.GetConfigValue(ConfigSectionName, "statusqueryip")
	if strings.Index(allowip, ip) == -1 && strings.Index(allowip, "*") == -1 {
		mlog.Error("客户端 %s 請求查詢轉發列表, 但不在允許的 IP 範圍", ip)
		io.WriteString(w, "forbidden")
		return
	}

	tfs := tfconfig.TFServerManager.GetServerList()

	data := "<table>"
	data += "<tr><td>IP Address</td><td>當前連接數</td><td>最後更新時間</td></tr>"
	for _, t := range tfs {
		data += fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%s</td></tr>", t.IPAddress, t.CurrentConn, t.LastRegTime.Format("2006-01-02 15:04:05"))
	}
	data += "</table>"

	io.WriteString(w, data)

}

var cfg *config.Config

// 初始化配置
func initConfig() {
	cfgs := config.GetConfigGroup(ConfigDirName)
	cfg = cfgs.GetConfig("tfconfig")
}
