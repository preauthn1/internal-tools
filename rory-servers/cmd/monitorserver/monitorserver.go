package main

import (
	"fmt"
	"runtime"
	"time"

	"example.invalid/ServerTeam/Servers/internal/monitorserver/monitorhandler"
	"example.invalid/ServerTeam/Servers/internal/monitorserver/monitorutil"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

func init() {
	conf.CLogLevel = mlog.ErrorLevel
	conf.ClearLog = true
}

func main() {
	conf := monitorutil.GetMonitorServerConfig("monitor")
	monitorhandler.ReTryCount = conf.ReTryCount
	monitorhandler.ReTryInterval = conf.ReTryInterval

	go func() {
		for {
			mlog.System("MemUsage : %s , NumGoroutine: %d", monitorutil.GetMemUsage(), runtime.NumGoroutine())
			Run(conf)
			time.Sleep(time.Second * time.Duration(conf.Time))
		}
	}()

	select {}
}

// Run run server
func Run(conf *monitorutil.MonitorServerConfig) {
	defer util.PanicErrStack()
	var errMsg string
	errCount := 0
	start := time.Now()
	for _, v := range conf.Serverlist {
		tfData := monitorutil.BuildTfData(v.ServerURL, v.SiteID)
		tfQuery := fmt.Sprintf(v.ServerURL, tfData)
		tfRet := ""
		// fetch hallserver
		err := util.Retry(monitorhandler.ReTryCount, time.Duration(monitorhandler.ReTryInterval)*time.Second, func() error {
			var e error
			tfRet, e = monitorutil.Get(tfQuery) //fix for loop
			if e != nil {
				mlog.Error("%v", e)
				return e
			}
			return nil
		})
		if err != nil {
			errMsg = "fetch url [" + tfQuery + "] error msg [" + err.Error() + "]"
			monitorutil.Push(conf.NotifyWebHook, errMsg, time.Since(start))
			continue
		}

		mlog.System(fmt.Sprintf("fetch url [%s]", tfQuery))
		ret, _ := msgserver.NodeDefDecryptor(tfRet)
		token, err := msgserver.ConnTokenDecryptor(tfRet)
		mlog.Debug("token %v, %s", token, err)

		mlog.Debug("%v", ret)

		serverIP := fmt.Sprintf("%s:%d", ret.Data.Addr, ret.Data.TransferPort)
		selfIP := ret.Data.UserIP
		if serverIP == "" {
			errMsg := "fetch url [" + tfQuery + "] error, with data [" + tfData + "]"
			monitorutil.Push(conf.NotifyWebHook, errMsg, time.Since(start))
			continue
		}

		//serverip = "0.0.0.0:10000"
		//serverip = "0.0.0.0:10000"
		//serverip = "0.0.0.0:10000"
		//serverIP = "0.0.0.0:10000"
		//serverIP = "0.0.0.0:10000"
		//v.SiteId = 1
		mlog.System("Connect HallServer [%s], selfip [%s]", serverIP, selfIP)
		err = monitorhandler.StartClient(ret, tfRet, v.SiteID, int32(1), conf.CheckGame.Enable, conf.CheckGame.GameID, conf.CheckGame.Port)
		if err != nil {
			errMsg = "serverIP : " + serverIP + "<br />" + err.Error()
			//monitorutil.Push(conf.NotifyWebHook, errMsg, time.Since(start))
			errCount++
		} else {
			errCount = 0
		}

		if errCount == 0 {
			break
		}

		if conf.CheckGame.Enable == 1 {
			time.Sleep(time.Duration(7) * time.Second)
		}
	}

	if errCount > 0 {
		monitorutil.Push(conf.NotifyWebHook, errMsg, time.Since(start))
	}
}
