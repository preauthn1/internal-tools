package gamelogic

import (
	"bytes"
	"fmt"
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/util"
)

// RoomCmdP 房間命令
var RoomCmdP *RoomCmd

// RoomCmd 房間命令
type RoomCmd struct {
	room *Room
}

// NewRoomCmd 創建房間命令對象
func NewRoomCmd(room *Room) *RoomCmd {
	rcmd := new(RoomCmd)
	rcmd.room = room

	return rcmd
}

// Name 指令名稱
func (cmd *RoomCmd) Name() string {
	return "room"
}

// Help 指令說明
func (cmd *RoomCmd) Help() string {
	return "房間管理，eg:room stat | room tablelist | room tablefind index | room finduser userid"
}

// Exec 執行指令
func (cmd *RoomCmd) Exec(args []string) string {
	if len(args) == 0 {
		return "error.seeref usage"
	}

	cmdname := args[0]

	switch cmdname {
	case "stat":
		formatstr := "gameid:%d,gamename:%s,serverid:%d,servername:%s,flag:%s,addr:%s:%d,tablecount:%d,seatcount:%d"
		return fmt.Sprintf(formatstr, cmd.room.GameInfo.GameID, cmd.room.GameInfo.GameName, cmd.room.ServerID, cmd.room.ServerName, cmd.room.Flag, cmd.room.IPAddress, cmd.room.Port, cmd.room.TableCount, cmd.room.SeatCount)
	case "tablelist":
		var buf bytes.Buffer

		for _, t := range cmd.room.Tables {
			buf.WriteString(fmt.Sprintf("tableno:%d,", t.TableNo))

			for _, st := range t.Seats {
				buf.WriteString(fmt.Sprintf("seatno:%d,user:%v", st.SeatNo, st.SeatUser))
			}

			buf.WriteString("\r\n")
		}

		return buf.String()
	case "tablefind":
		if len(args) != 2 {
			return "error,room tablefind index"
		}
		sno := args[1]

		if !util.IsNumber(sno) {
			return "error,TableManager index 必須為數字"
		}

		tno, _ := strconv.Atoi(sno)

		if int32(tno) >= cmd.room.TableCount {
			return fmt.Sprintf("error,桌號不能超出%d", cmd.room.TableCount-1)
		}

		table := cmd.room.Tables[tno]

		var buf bytes.Buffer

		buf.WriteString(fmt.Sprintf("tableno:%d,", table.TableNo))

		for _, st := range table.Seats {
			buf.WriteString(fmt.Sprintf("seatno:%d,user:%v", st.SeatNo, st.SeatUser))
		}

		buf.WriteString("\r\n")

		return buf.String()
	case "finduser":
		if len(args) != 2 {
			return "error,room finduser userid"
		}
		suserid := args[1]
		if !util.IsNumber(suserid) {
			return "error,用戶id必須為數字"
		}

		userid, _ := strconv.Atoi(suserid)

		user := cmd.room.RoomUserManager.GetGameUserByUserID(int32(userid))

		if user == nil {
			return "用戶不在房間中"
		}

		tno := user.TableNo

		if tno == FreeTableNo {
			return fmt.Sprintf("%v", user)
		}

		table := cmd.room.Tables[tno]

		var buf bytes.Buffer

		buf.WriteString(fmt.Sprintf("tableno:%d,", table.TableNo))

		for _, st := range table.Seats {
			buf.WriteString(fmt.Sprintf("seatno:%d,user:%v", st.SeatNo, st.SeatUser))
		}

		buf.WriteString("\r\n")

		return buf.String()

	default:
		return ""
	}

}

// Room 取得整個房間
func (cmd *RoomCmd) Room() *Room {
	return cmd.room
}

// RoomName 房間名稱
func (cmd *RoomCmd) RoomName() string {
	return cmd.room.ServerName
}

// TableCount 查詢桌數
func (cmd *RoomCmd) TableCount() int {
	return len(cmd.room.Tables)
}

// Table 取得桌子
func (cmd *RoomCmd) Table() []*Table {
	return cmd.room.Tables
}
