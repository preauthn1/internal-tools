package player

import (
	"bytes"
	"fmt"
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// PlayerCmdP 用戶命令
var PlayerCmdP *Cmd

// Cmd 用戶命令
type Cmd struct {
	playerManager *Manager
	server        *msgserver.MsgServer
}

// NewPlayerCmd 創建用戶命令
func NewPlayerCmd(playerManager *Manager, server *msgserver.MsgServer) *Cmd {
	cmd := new(Cmd)
	cmd.playerManager = playerManager
	cmd.server = server

	return cmd
}

// Name 指令名稱
func (pc *Cmd) Name() string {
	return "player"
}

// Help 指令說明
func (pc *Cmd) Help() string {
	return "用户管理，eg:player stat | player kick userid | player find userid | player list 0(1)"
}

// Exec 執行指令
func (pc *Cmd) Exec(args []string) string {
	if len(args) == 0 {
		return "error.seeref usage"
	}

	cmdname := args[0]

	switch cmdname {
	case "stat":
		return fmt.Sprintf("total:%d logon user:%d,unlogon user:%d", pc.PlayerCount(), pc.UserCount(), pc.GuestCount())
	case "kick":
		if len(args) != 2 {
			return "error.eg kick 10001"
		}

		if !util.IsNumber(args[1]) {
			return "error.用户ID必须为数字"
		}

		uid, _ := strconv.Atoi(args[1])

		pr := pc.playerManager.GetPlayerByUserID(int32(uid))
		if pr == nil {
			return "用户不存在"
		}

		pc.server.CloseConn(pr.Conn)
		return "kick success"
	case "find":
		if len(args) != 2 {
			return "error.eg find 10001"
		}

		if !util.IsNumber(args[1]) {
			return "error.用户ID必须为数字"
		}

		uid, _ := strconv.Atoi(args[1])
		pr := pc.playerManager.GetPlayerByUserID(int32(uid))

		if pr == nil {
			return "用户不存在"
		}
		return pr.String()
	case "list":
		if len(args) != 2 {
			return "error.eg list (0/1/2)"
		}

		if args[1] == "0" {
			var buf bytes.Buffer

			pc.playerManager.Range(func(pr *Player) {
				buf.WriteString(fmt.Sprintf("%s\r\n", pr))
			})

			return buf.String()
		} else if args[1] == "1" {
			var buf bytes.Buffer
			prs := pc.playerManager.LogonPlayer()

			for _, pr := range prs {
				buf.WriteString(fmt.Sprintf("%s\r\n", pr))
			}

			return buf.String()

		} else if args[1] == "2" {
			var buf bytes.Buffer
			prs := pc.playerManager.UnLogonPlayer()

			for _, pr := range prs {
				buf.WriteString(fmt.Sprintf("%s\r\n", pr))
			}

			return buf.String()
		} else if args[1] == "3" {
			var buf bytes.Buffer

			pc.playerManager.Range(func(pr *Player) {
				buf.WriteString(fmt.Sprintf("%s\r\n", pr.SimpleString()))
			})

			return buf.String()
		} else {
			return "error.eg list (0/1/2/3)"
		}

	}

	return ""
}

// PlayerCount 所有玩家人數
func (pc *Cmd) PlayerCount() int {
	return pc.playerManager.PlayerCount()
}

// UserCount 登入玩家人數
func (pc *Cmd) UserCount() int {
	return pc.playerManager.UserCount()
}

// GuestCount 未登入玩家人數
func (pc *Cmd) GuestCount() int {
	return pc.playerManager.GuestCount()
}

// QueryUser 查詢玩家
func (pc *Cmd) QueryUser(uid int32) (*Player, error) {
	p := pc.playerManager.GetPlayerByUserID(int32(uid))
	if p == nil {
		return nil, fmt.Errorf("玩家[%d]不存在", uid)
	}
	return p, nil
}

// Kick 剔除玩家
func (pc *Cmd) Kick(uid int32, m string) error {
	p := pc.playerManager.GetPlayerByUserID(int32(uid))
	if p == nil {
		return fmt.Errorf("玩家[%d]不存在", uid)
	}
	pc.server.CloseConn(p.Conn)
	return nil
}

// ListPlayer 列出所有
func (pc *Cmd) ListPlayer() []*Player {
	var players []*Player

	pc.playerManager.Range(func(p *Player) {
		players = append(players, p)
	})

	return players
}

// ListUser 列出所有登入
func (pc *Cmd) ListUser() []*Player {
	return pc.playerManager.LogonPlayer()
}

// ListGuest 列出所有訪客
func (pc *Cmd) ListGuest() []*Player {
	return pc.playerManager.UnLogonPlayer()
}
