package monitorhandler

import (
	"errors"
	"fmt"
	"net"

	"example.invalid/ServerTeam/Servers/internal/monitorserver/monitorhandler/msghandler"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
)

// FishClient ?
func FishClient(tfret *msgserver.TFDataRet, token string, siteid, guestid int32) error {

	selfIP := tfret.Data.UserIP
	serverIP := fmt.Sprintf("%s:%d", tfret.Data.Addr, tfret.Data.TransferPort)
	mlog.System(serverIP)
	msghandler.SelfIP = selfIP
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverIP)
	checkError(err, "ResolveTCPAddr")
	hc := msghandler.NewHallClient(tcpAddr.String())
	hc.SetToken(token)
	mlog.System("token , %s", token)
	defer hc.Stop()

	hc.RegisterMsgHandler()
	if err := hc.Start(); err != nil {
		return err
	}
	if err := hc.RequestHallVersion(); err != nil {
		return err
	}
	loginret, err := hc.RequestGuestLogon(siteid, guestid)
	_ = loginret
	if err != nil {
		return err
	}
	//limiter := time.Tick(time.Duration(5000) * time.Millisecond)
	//for {
	//	<-limiter
	mlog.System("userid %d", loginret.GetUserID())
	serverInfo, err := hc.RequestServerList()
	for _, v := range serverInfo.ServerList {
		mlog.Debug("gameid[%v] flag[%v] name[%v] status[%v] ", v.GetGameID(), v.GetFlag(), v.GetGameName(), v.GetStatus())
	}
	if err != nil {
		return fmt.Errorf("獲取不到遊戲服務器列表, hc isconnected %v", hc.IPReported)
		return err
	}
	_ = serverInfo
	for _, server := range serverInfo.GetServerList() {
		if server.GetGameID() != 14 {
			continue
		}
		_ = server
		gs, _ := hc.RequestServerAddr(server.GetGameID(), server.GetFlag())

		// 轉發時 gameserver addr 必須使用 節點ip
		//gsAddr := fmt.Sprintf("%s:%d", gs.GetIPAddress(), gs.GetPort())
		gsAddr := fmt.Sprintf("%s:%d", tcpAddr.IP, gs.GetPort())
		serverMsg := fmt.Sprintf("gamename: %s, gameid: %d, flag: %s, serverid:%d, serveraddr:%s", server.GetGameName(), server.GetGameID(), server.GetFlag(), gs.GetServerID(), gsAddr)

		if gs.GetServerID() == 0 {
			return fmt.Errorf("獲取不到遊戲服務器連線資訊, %s, hc isconnected %v", serverMsg, hc.ConnClt.IsConnected())
		}

		err := func() error {
			gc := msghandler.NewGameClient(gsAddr)

			defer gc.Stop()
			gc.EnableHeartBeat = false
			gc.RegisterMsgHandler()
			gc.ByLeaveTipHandler()
			if err := gc.Start(); err != nil {
				return errors.New(serverMsg + err.Error())
			}
			if err := gc.RequestGameVer(); err != nil {
				return errors.New(serverMsg + err.Error())
			}

			if server.GetGameID() == 14 && gs.GetServerID() == 88 {
				if err := gc.LoginGameID(loginret, 0); err != nil {
					return errors.New(serverMsg + err.Error())
				}
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}
	//}
	select {}
	return nil
}
