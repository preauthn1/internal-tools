package monitorutil

import (
	"fmt"
	"runtime"

	"gopkg.in/yaml.v2"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// MonitorServerConfig 監控配置
type MonitorServerConfig struct {
	Time          int          `yaml:"Time"`          //監控時間頻率(秒)
	ReTryCount    int          `yaml:"ReTryCount"`    //重新連線次數
	ReTryInterval int          `yaml:"ReTryInterval"` //重新連線時間間隔
	NotifyWebHook string       `yaml:"NotifyWebHook"` //推播錯誤訊息路徑
	Serverlist    []ServerList `yaml:"Serverlist"`    //serverlist 大廳
	GetIP         string       `yaml:"GetIp"`
	CheckGame     CheckGame    `yaml:"CheckGame"`
}

// CheckGame check game
type CheckGame struct {
	Enable int `yaml:"Enable"`
	GameID int `yaml:"GameId"`
	Port   int `yaml:"Port"`
}

// ServerList server lsit
type ServerList struct {
	SiteID    int32  `yaml:"SiteId"`
	ServerURL string `yaml:"Url"`
}

// GetMonitorServerConfig get config
func GetMonitorServerConfig(fileName string) *MonitorServerConfig {
	cfData, err := GetConfig(fileName)
	if err != nil {
		mlog.Error("monitor config err : %v", err)
	}

	cfStruct := MonitorServerConfig{}
	err2 := yaml.Unmarshal([]byte(cfData), &cfStruct)
	if err2 != nil {
		mlog.Error("cannot unmarshal data: %v", err2)
	}

	return &cfStruct
}

// GetMemUsage get memory usage
func GetMemUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://example.invalid/pkg/runtime/#MemStats
	s := fmt.Sprintf("Alloc = %v MiB", bToMb(m.Alloc))
	s += fmt.Sprintf(" TotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	s += fmt.Sprintf(" Sys = %v MiB", bToMb(m.Sys))
	s += fmt.Sprintf(" NumGC = %v", m.NumGC)
	return s
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
