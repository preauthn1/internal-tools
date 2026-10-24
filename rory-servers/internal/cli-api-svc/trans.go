package cliapisvr

import (
	"strconv"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
)

// TransGameUser 切換玩家資料
func TransGameUser(u *gamelogic.GameUser, _p *pb.Player) *pb.Player {
	if _p == nil {
		_p = new(pb.Player)
	}
	_p.PlayerID = u.UserID

	_p.Table = u.TableNo
	_p.Seat = u.SeatNo
	_p.State = u.UserState.StateName()

	_p.PlayerName = u.NickName
	_p.PlayerScore = u.Score
	_p.PlayerBetAmount = u.BetAmount
	_p.PlayerPayoutAmount = u.PayAmount
	_p.PlayerFlowAmount = u.FlowAmount
	_p.PlayerTaxAmount = u.TaxAmount
	_p.Type = u.UserType

	return _p
}

// TransPlayer 切換玩家資料
func TransPlayer(p *player.Player) *pb.Player {
	_p := new(pb.Player)
	_p.PlayerID = p.UserID
	_p.Conn = new(pb.ConnectionInfo)
	_p.Conn.IP = p.Conn.IPAddress
	_p.Conn.Port = int32(p.Conn.Port)
	_p.Conn.CreateTime = p.Conn.ConnTime.Format(time.RFC3339Nano)

	if p.UserData != nil {
		TransGameUser(p.UserData.(*gamelogic.GameUser), _p)
	}

	return _p
}

// TransRobot 切換機器人資料
func TransRobot(rbt *robotcomm.Robot) *pb.Robot {
	r := new(pb.Robot)
	r.Level = rbt.RebotLevel
	r.Player = TransGameUser(rbt.UserInfo, nil)
	return r
}

// TransRoom 轉換房間資料
func TransRoom(room *gamelogic.Room) *pb.Room {
	r := new(pb.Room)

	r.ID = room.ServerID
	r.Name = room.ServerName
	r.GameID = room.GameInfo.GameID
	r.GameName = room.GameInfo.GameName
	f, _ := strconv.ParseInt(room.Flag, 0, 32)
	r.Flag = int32(f)

	for _, t := range room.Tables {
		r.Tables = append(r.Tables, TransTable(t))
	}

	return r
}

// TransTable 轉換牌桌資料
func TransTable(tb *gamelogic.Table) *pb.Table {
	t := new(pb.Table)

	t.ID = tb.TableNo
	t.Gaming = tb.GameLogic.IsGaming()

	for _, u := range tb.GetTableUser() {
		t.Players = append(t.Players, TransGameUser(u, nil))
	}

	t.GameSummary = tb.GameLogic.Summary()

	return t
}
