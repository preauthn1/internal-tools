package monitorhandler

import (
	"errors"
	"fmt"
	"net"
	"time"

	"example.invalid/ServerTeam/Servers/internal/monitorserver/monitorhandler/msghandler"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ?
var (
	ReTryCount    int
	ReTryInterval int
)

// StartClient 客戶端啟動函數
//      遠程ip地址和端口 tcpaddr
func StartClient(tfret *msgserver.TFDataRet, token string, siteid, guestid int32, enable, gameid, port int) error {
	selfIP := tfret.Data.UserIP
	serverIP := fmt.Sprintf("%s:%d", tfret.Data.Addr, tfret.Data.TransferPort)
	mlog.System(serverIP)
	msghandler.SelfIP = selfIP
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverIP)
	checkError(err, "ResolveTCPAddr")
	hc := msghandler.NewHallClient(tcpAddr.String())
	hc.SetToken(token)
	mlog.Debug("token , %s", token)
	defer hc.Stop()

	hc.RegisterMsgHandler()
	if err := hc.Start(); err != nil {
		return err
	}
	if err := hc.RequestHallVersion(); err != nil {
		return err
	}
	loginret, err := hc.RequestGuestLogon(siteid, guestid)
	if err != nil {
		return err
	}
	serverInfo, err := hc.RequestServerList()
	if err != nil {
		return err
	}
	for _, server := range serverInfo.GetServerList() {
		gs, _ := hc.RequestServerAddr(server.GetGameID(), server.GetFlag())
		// 轉發時 gameserver addr 必須使用 節點ip
		//gsAddr := fmt.Sprintf("%s:%d", gs.GetIPAddress(), gs.GetPort())
		gsAddr := fmt.Sprintf("%s:%d", tcpAddr.IP, gs.GetPort())
		serverMsg := fmt.Sprintf("gamename: %s, gameid: %d, flag: %s, serverid:%d, serveraddr:%s", server.GetGameName(), server.GetGameID(), server.GetFlag(), gs.GetServerID(), gsAddr)
		if gs.GetServerID() == 0 {
			return fmt.Errorf("獲取不到遊戲服務器連線資訊, %s", serverMsg)
		}
		if enable != 1 {
			mlog.System(serverMsg)
		}

		if enable == 1 && int32(gameid) != server.GetGameID() {
			continue
		} else if enable == 1 && int32(gameid) == server.GetGameID() {
			gsAddr = fmt.Sprintf("%s:%d", tcpAddr.IP, port)
		}

		err := func() error {
			gc := msghandler.NewGameClient(gsAddr)
			gc.SetToken(token)
			defer gc.Stop()

			if err := gc.Start(); err != nil {
				return errors.New(serverMsg + err.Error())
			}

			if err := gc.RequestGameVer(); err != nil {
				return errors.New(serverMsg + err.Error())
			}

			if enable == 1 {
				mlog.System("%v 登入中 ", server.GetGameName())
				if err := gc.LoginGameID(loginret, enable); err != nil {
					return errors.New(serverMsg + err.Error())
				}
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}

// RetryHandler ?
var RetryHandler = func(conn net.Conn, handler func(conn net.Conn) (err error)) error {
	err2 := util.Retry(ReTryCount, time.Duration(ReTryInterval)*time.Second, func() error {
		err := handler(conn)
		if err != nil {
			mlog.Error("%v", err)
			return err
		}
		return nil
	})

	if err2 != nil {
		return err2
	}

	return nil
}

//错误检查
func checkError(err error, info string) (res bool) {
	if err != nil {
		fmt.Println(info + "  " + err.Error())
		return false
	}
	return true
}
