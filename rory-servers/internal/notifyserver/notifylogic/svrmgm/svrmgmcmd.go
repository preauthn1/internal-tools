package svrmgm

import (
	"fmt"
)

// Cmd ?
type Cmd struct {
	serverList *ServerList
}

// NewSvrMgmCmd ?
func NewSvrMgmCmd(svrlst *ServerList) *Cmd {
	svrcmd := new(Cmd)
	svrcmd.serverList = svrlst

	return svrcmd
}

// Name ?
func (svrcmd *Cmd) Name() string {
	return "ns"
}

// Exec ?
func (svrcmd *Cmd) Exec(args []string) string {
	if len(args) == 0 {
		return "error,seeref usage"
	}

	if args[0] == "list" {
		if len(args) != 2 {
			return "error,seeref ueage"
		}

		listparam := args[1]

		var servers []*ServerInfo
		if listparam == "0" {
			servers = svrcmd.serverList.GetAllServer()
		} else if listparam == "1" {
			servers = svrcmd.serverList.GetHallServer()
		} else if listparam == "2" {
			servers = svrcmd.serverList.GetGameServer()
		}

		str := ""
		for _, svr := range servers {
			str += fmt.Sprintf("%s\r\n", svr.String())
		}

		return str
	}

	return ""
}

// Help ?
func (svrcmd *Cmd) Help() string {
	return "註冊服務器查詢.usage:ns list 0(1/2)"
}
