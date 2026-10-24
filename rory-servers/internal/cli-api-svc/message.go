package cliapisvr

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/dig"
	"example.invalid/grpc/codes"
	"example.invalid/grpc/status"

	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/roomlogic"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
)

// MessageServiceServer Message Service
type MessageServiceServer struct{}

// Trace 追蹤封包
func (s *MessageServiceServer) Trace(cond *pb.QueryCondition, stream pb.MessageService_TraceServer) error {

	if app.Context == nil {
		return status.Error(codes.Aborted, "app.Context not set")
	}

	err := app.Context.Invoke(func(room *gamelogic.Room, mgr *roomlogic.RoomActionManager) error {
		fmt.Println(room, mgr)
		t := room.GetTableByNo(cond.Table)
		if t == nil {
			return status.Error(codes.NotFound, fmt.Sprintf("Could not find Table[%d]", cond.Table))
		}

		if cond.PlayerID > 0 {

			var hasUser bool
			for _, u := range t.GetTableUser() {
				if u.UserID == cond.PlayerID {
					hasUser = true
					break
				}
			}

			if !hasUser {
				return status.Error(codes.NotFound, fmt.Sprintf("Could not find Player[%d]", cond.PlayerID))
			}

			example.invalid("開始追蹤遊戲封包, 桌號: %d, 玩家: %d", cond.Table, cond.PlayerID)
		} else {
			example.invalid("開始追蹤遊戲封包, 桌號: %d", cond.Table)
		}

		c2sMessageCh := make(chan *pb.Message)
		s2cMessageCh := make(chan *pb.Message)

		fn := func(user *gamelogic.GameUser, m *msg.Message) {
			if cond.PlayerID > 0 && user.UserID != cond.PlayerID {
				return
			}
			// 過濾不同桌
			if user.TableNo != t.TableNo {
				return
			}

			_m := new(pb.Message)
			_m.PlayerID = user.UserID
			_m.Src = pb.Endpoint_Client
			_m.Dest = pb.Endpoint_Server
			_m.BClassID = m.BClassID
			_m.SClassID = m.SClassID
			_m.Data = fmt.Sprintf("%v", m.MsgData)
			if msg.InstanceProvider != nil {
				if transFn := msg.InstanceProvider.Load(m.BClassID, m.SClassID); transFn != nil {
					obj := transFn()
					if err := proto.Unmarshal(m.MsgData.([]byte), obj); err == nil {
						_m.Data = fmt.Sprintf("%v", obj)
					}
				}
			}

			c2sMessageCh <- _m
		}

		mgr.OnGameMessage = fn

		game := t.GameLogic
		msgfunc := game.GetSendMessageFunc()

		transFunc := func(uid int32, m *msg.Message) {
			if cond.PlayerID > 0 && uid != cond.PlayerID {
				return
			}
			// 過濾機器人
			for _, u := range t.GetTableUser() {
				if u.UserID == uid && u.IsRobot() {
					return
				}
			}

			_m := new(pb.Message)
			_m.PlayerID = uid
			_m.Src = pb.Endpoint_Server
			_m.Dest = pb.Endpoint_Client
			_m.BClassID = m.BClassID
			_m.SClassID = m.SClassID
			_m.Data = fmt.Sprintf("%v", m.MsgData)

			s2cMessageCh <- _m
			msgfunc(uid, m)
		}

		game.OnSendMessage(transFunc)

		go func(stream pb.MessageService_TraceServer, ch <-chan *pb.Message) {
			for m := range ch {
				stream.Send(m)
			}
		}(stream, s2cMessageCh)

		go func(stream pb.MessageService_TraceServer, ch <-chan *pb.Message) {
			for m := range ch {
				stream.Send(m)
			}
		}(stream, c2sMessageCh)

		select {
		case <-stream.Context().Done():
			close(s2cMessageCh)
			close(c2sMessageCh)
			game.OnSendMessage(msgfunc)
			mgr.OnGameMessage = nil

			if cond.PlayerID > 0 {
				example.invalid("停止追蹤遊戲封包, 桌號: %d, 玩家: %d", cond.Table, cond.PlayerID)
			} else {
				example.invalid("停止追蹤遊戲封包, 桌號: %d", cond.Table)
			}
		}

		return nil
	})
	if err != nil {
		if dig.CanVisualizeError(err) {
			return status.Error(codes.Internal, fmt.Sprintf("Invoke error, %v", err))
		}
		return err
	}

	return nil
}
