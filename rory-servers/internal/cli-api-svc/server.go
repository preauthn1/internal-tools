package cliapisvr

import (
	"context"
	"fmt"
	"time"

	"example.invalid/dig"
	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotaction"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/roomlogic"
	"example.invalid/grpc/codes"
	"example.invalid/grpc/status"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
)

// CLIService 處理服務
type CLIService struct{}

// ServerVer 取得 Server 版號
func (s *CLIService) ServerVer(ctx context.Context, arg *pb.Empty) (*pb.Version, error) {
	v := new(pb.Version)
	v.Version = conf.GitTag
	v.BuildName = conf.BuildNum
	v.BuildTime = conf.BuildTime
	v.Hash = conf.GitHash
	return v, nil
}

// ClientVer 取得 Client 版號
func (s *CLIService) ClientVer(ctx context.Context, arg *pb.Empty) (*pb.CVersion, error) {

	if roomlogic.CVersionCmdP == nil {
		return nil, fmt.Errorf("CVersionCmd not init")
	}

	v := new(pb.CVersion)
	v.IOS = roomlogic.CVersionCmdP.GetIOS()
	v.Android = roomlogic.CVersionCmdP.GetAndroid()

	return v, nil
}

// SetClientVer 設定 Client 版號
func (s *CLIService) SetClientVer(ctx context.Context, v *pb.CVersion) (*pb.Empty, error) {

	if roomlogic.CVersionCmdP == nil {
		return nil, fmt.Errorf("CVersionCmd not init")
	}

	if len(v.IOS) > 0 {
		roomlogic.CVersionCmdP.SetIOS(v.IOS)
	}

	if len(v.Android) > 0 {
		roomlogic.CVersionCmdP.SetAndroid(v.Android)
	}

	return new(pb.Empty), nil
}

// Shotdown 關閉 Server
func (s *CLIService) Shotdown(ctx context.Context, req *pb.ShotdownRequest) (*pb.Empty, error) {

	if roomlogic.RoomStateCmdP == nil {
		return nil, fmt.Errorf("RoomStateCmd not init")
	}

	dur, err := time.ParseDuration(req.Delay)
	if err != nil {
		return nil, fmt.Errorf("Parse Delay Error. %v", err)
	}

	defer roomlogic.RoomStateCmdP.Shoutdown(req.Message, dur)

	return new(pb.Empty), nil
}

// LogStatus 查詢 Log 狀態
func (s *CLIService) LogStatus(ctx context.Context, arg *pb.Empty) (*pb.LogState, error) {
	stat := new(pb.LogState)
	stat.Enable = mlog.GetCLogger().Enabled()
	stat.Lite = conf.ClearLog
	stat.Level = int32(conf.CLogLevel)
	return stat, nil
}

// SetLog 設定 Log
func (s *CLIService) SetLog(ctx context.Context, stat *pb.LogState) (*pb.Empty, error) {
	mlog.GetCLogger().Enable(stat.Enable)
	conf.ClearLog = stat.Lite
	conf.CLogLevel = int(stat.Level)
	return new(pb.Empty), nil
}

// GetPlayerSummary 取得玩家總攬
func (s *CLIService) GetPlayerSummary(ctx context.Context, arg *pb.Empty) (*pb.PlayerSummary, error) {

	if player.PlayerCmdP == nil {
		return nil, fmt.Errorf("PlayerCmd not init")
	}

	summary := new(pb.PlayerSummary)
	summary.PlayerCount = int64(player.PlayerCmdP.PlayerCount())
	summary.UserCount = int64(player.PlayerCmdP.UserCount())
	summary.GuestCount = int64(player.PlayerCmdP.GuestCount())

	return summary, nil
}

// QueryPlayer 查詢玩家
func (s *CLIService) QueryPlayer(ctx context.Context, cond *pb.QueryCondition) (*pb.Player, error) {
	if player.PlayerCmdP == nil {
		return nil, fmt.Errorf("PlayerCmd not init")
	}

	p, err := player.PlayerCmdP.QueryUser(cond.PlayerID)
	if err != nil {
		return nil, err
	}

	return TransPlayer(p), nil
}

// KickPlayer 剔除玩家
func (s *CLIService) KickPlayer(ctx context.Context, cond *pb.QueryCondition) (*pb.Empty, error) {

	if player.PlayerCmdP == nil {
		return nil, fmt.Errorf("PlayerCmd not init")
	}

	if err := player.PlayerCmdP.Kick(cond.PlayerID, ""); err != nil {
		return nil, err
	}

	return new(pb.Empty), nil
}

// ListPlayer 列出所有玩家
func (s *CLIService) ListPlayer(ctx context.Context, cond *pb.QueryCondition) (*pb.PlayerList, error) {

	if player.PlayerCmdP == nil {
		return nil, fmt.Errorf("PlayerCmd not init")
	}

	list := new(pb.PlayerList)

	var players []*player.Player

	if (cond.IsLogined && cond.IsGuest) || (!cond.IsLogined && !cond.IsGuest) {
		players = player.PlayerCmdP.ListPlayer()
	} else if cond.IsLogined && !cond.IsGuest {
		players = player.PlayerCmdP.ListUser()
	} else if !cond.IsLogined && cond.IsGuest {
		players = player.PlayerCmdP.ListGuest()
	}

	for _, p := range players {
		list.Players = append(list.Players, TransPlayer(p))
	}

	return list, nil
}

// GetRobotSummary 查詢機器人概況
func (s *CLIService) GetRobotSummary(ctx context.Context, arg *pb.Empty) (*pb.RobotSummary, error) {
	if robotaction.RobotCmdP == nil {
		return nil, fmt.Errorf("此遊戲沒有設定機器人")
	}

	summary := new(pb.RobotSummary)
	summary.RobotCount = int64(robotaction.RobotCmdP.RobotCount())
	summary.GamingCount = int64(robotaction.RobotCmdP.GameRobotCount())
	summary.FreeCount = int64(robotaction.RobotCmdP.FreeRobotCount())

	return summary, nil
}

// QueryRobot 查詢指定機器人
func (s *CLIService) QueryRobot(ctx context.Context, cond *pb.QueryCondition) (*pb.Robot, error) {
	if robotaction.RobotCmdP == nil {
		return nil, fmt.Errorf("此遊戲沒有設定機器人")
	}

	rbt := robotaction.RobotCmdP.GetRobot(cond.PlayerID)

	return TransRobot(rbt), nil
}

// ListRobot 列出所有機器人
func (s *CLIService) ListRobot(ctx context.Context, cond *pb.QueryCondition) (*pb.RobotList, error) {
	if robotaction.RobotCmdP == nil {
		return nil, fmt.Errorf("此遊戲沒有設定機器人")
	}

	list := new(pb.RobotList)

	var players []*robotcomm.Robot

	rbts := robotaction.RobotCmdP.GetAllRobot()

	for _, rbt := range rbts {
		if cond.IsGaming && !rbt.UserInfo.IsGameing() {
			continue
		}
		if cond.IsFree && !rbt.UserInfo.IsFree() {
			continue
		}
		players = append(players, rbt)
	}

	for _, p := range players {
		list.Robots = append(list.Robots, TransRobot(p))
	}

	return list, nil
}

// RoomInfo 查詢房間資訊
func (s *CLIService) RoomInfo(ctx context.Context, arg *pb.Empty) (room *pb.Room, err error) {
	if app.Context == nil {
		return nil, status.Error(codes.Aborted, "app.Context not set")
	}
	if err := app.Context.Invoke(func(r *gamelogic.Room) {
		room = TransRoom(r)
	}); err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}
	return
}

// TableInfo 查詢桌子資訊
func (s *CLIService) TableInfo(ctx context.Context, cond *pb.QueryCondition) (table *pb.Table, err error) {
	if app.Context == nil {
		return nil, status.Error(codes.Aborted, "app.Context not set")
	}
	if err := app.Context.Invoke(func(r *gamelogic.Room) error {
		t := r.GetTableByNo(cond.Table)
		if t == nil {
			return status.Error(codes.NotFound, fmt.Sprintf("Table %d not found", cond.Table))
		}
		table = TransTable(t)
		return nil
	}); err != nil {
		if dig.CanVisualizeError(err) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return nil, err
	}
	return
}
