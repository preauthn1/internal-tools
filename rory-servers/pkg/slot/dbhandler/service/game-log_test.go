package service

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/model"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

func TestDBGameLog_Add(t *testing.T) {
	s := new(DBGameLogService)
	t.Run("No User", func(tt *testing.T) {
		log := &netproto.Slot_MSG_GameLogAdd{
			UserID:        proto.Int32(123),
			GameID:        proto.Int32(1001),
			ServerID:      proto.Int32(398),
			CurrentAmount: proto.Int64(100),
			BetMoney:      proto.Int64(10),
			WinMoney:      proto.Int64(30),
			FinalScore:    proto.Int64(120),
			GameType:      proto.Int32(1),
			LogType:       proto.Int32(1),
			GameData:      proto.String("101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
			CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
		}
		err := s.Add("", log)
		tt.Log(err)
	})
	t.Run("NewLog", func(tt *testing.T) {
		log := &netproto.Slot_MSG_GameLogAdd{
			UserID:        proto.Int32(15267553),
			GameID:        proto.Int32(1004),
			ServerID:      proto.Int32(414),
			CurrentAmount: proto.Int64(100),
			BetMoney:      proto.Int64(10),
			WinMoney:      proto.Int64(30),
			FinalScore:    proto.Int64(120),
			GameType:      proto.Int32(1),
			LogType:       proto.Int32(1),
			GameData:      proto.String("[test add]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
			CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
		}
		err := s.Add("", log)
		tt.Log(err)
	})
}

func TestDBGameLog_BatchAdd(t *testing.T) {
	s := new(DBGameLogService)
	t.Run("NewLogOne", func(tt *testing.T) {
		log := &netproto.Slot_MSG_GameLogAdd{
			UserID:        proto.Int32(15267553),
			GameID:        proto.Int32(1004),
			ServerID:      proto.Int32(414),
			CurrentAmount: proto.Int64(100),
			BetMoney:      proto.Int64(10),
			WinMoney:      proto.Int64(30),
			FinalScore:    proto.Int64(120),
			GameType:      proto.Int32(1),
			LogType:       proto.Int32(1),
			GameData:      proto.String("[test batch one]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
			CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
		}
		err := s.BatchAdd("", []*netproto.Slot_MSG_GameLogAdd{log})
		tt.Log(err)
	})
	t.Run("NewLog Free", func(tt *testing.T) {
		log1 := &netproto.Slot_MSG_GameLogAdd{
			UserID:        proto.Int32(15267553),
			GameID:        proto.Int32(1004),
			ServerID:      proto.Int32(414),
			CurrentAmount: proto.Int64(100),
			BetMoney:      proto.Int64(10),
			WinMoney:      proto.Int64(30),
			FinalScore:    proto.Int64(120),
			GameType:      proto.Int32(1),
			LogType:       proto.Int32(1),
			GameData:      proto.String("[test batch 1]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
			CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
		}
		log2 := &netproto.Slot_MSG_GameLogAdd{
			UserID:        proto.Int32(15267553),
			GameID:        proto.Int32(1004),
			ServerID:      proto.Int32(414),
			CurrentAmount: proto.Int64(100),
			BetMoney:      proto.Int64(10),
			WinMoney:      proto.Int64(30),
			FinalScore:    proto.Int64(120),
			GameType:      proto.Int32(1),
			LogType:       proto.Int32(1),
			GameData:      proto.String("[test batch 2]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
			CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
		}
		log3 := &netproto.Slot_MSG_GameLogAdd{
			UserID:        proto.Int32(15267553),
			GameID:        proto.Int32(1004),
			ServerID:      proto.Int32(414),
			CurrentAmount: proto.Int64(100),
			BetMoney:      proto.Int64(10),
			WinMoney:      proto.Int64(30),
			FinalScore:    proto.Int64(120),
			GameType:      proto.Int32(1),
			LogType:       proto.Int32(1),
			GameData:      proto.String("[test batch 3]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
			CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
		}
		err := s.BatchAdd("", []*netproto.Slot_MSG_GameLogAdd{log1, log2, log3})
		tt.Log(err)
	})
}

func TestDBGameLog_Search(t *testing.T) {
	s := new(DBGameLogService)
	t.Run("No User", func(tt *testing.T) {
		pageCnt, logs, err := s.Search("", 1, 1001, 1, 10, "2019-10-04", "2019-10-04")
		tt.Log(err)
		tt.Logf("All Row Count: %d\n", pageCnt)
		tt.Logf("Logs Length: %d\n", len(logs))
		for _, row := range logs {
			tt.Logf(" - %v\n", row)
		}
	})
	t.Run("No Game", func(tt *testing.T) {
		pageCnt, logs, err := s.Search("", 9936738, 999, 1, 10, "2019-10-04", "2019-10-04")
		tt.Log(err)
		tt.Logf("All Row Count: %d\n", pageCnt)
		tt.Logf("Logs Length: %d\n", len(logs))
		for _, row := range logs {
			tt.Logf(" - %v\n", row)
		}
	})
	t.Run("Pass", func(tt *testing.T) {
		pageCnt, logs, err := s.Search("", 15267553, 1004, 1, 10, "2020-03-20", "2020-03-31")
		tt.Log(err)
		tt.Logf("All Row Count: %d\n", pageCnt)
		tt.Logf("Logs Length: %d\n", len(logs))
		for _, row := range logs {
			tt.Logf(" - %v\n", row)
		}
	})
}

func TestAPIGameLog_Add(t *testing.T) {
	s := new(APIGameLogService)
	t.Run("No User", func(tt *testing.T) {
		log := &netproto.Slot_MSG_GameLogAdd{
			UserID:        proto.Int32(123),
			GameID:        proto.Int32(1001),
			ServerID:      proto.Int32(398),
			CurrentAmount: proto.Int64(100),
			BetMoney:      proto.Int64(10),
			WinMoney:      proto.Int64(30),
			FinalScore:    proto.Int64(120),
			GameType:      proto.Int32(1),
			LogType:       proto.Int32(1),
			GameData:      proto.String("101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
			CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
		}
		err := s.Add("", log)
		tt.Log(err)
	})
	t.Run("NewLog", func(tt *testing.T) {
		log := &netproto.Slot_MSG_GameLogAdd{
			UserID:        proto.Int32(1000174),
			GameID:        proto.Int32(1011),
			ServerID:      proto.Int32(475),
			CurrentAmount: proto.Int64(100),
			BetMoney:      proto.Int64(10),
			WinMoney:      proto.Int64(30),
			FinalScore:    proto.Int64(120),
			GameType:      proto.Int32(1),
			LogType:       proto.Int32(1),
			GameData:      proto.String("[test add]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
			CreateTime:    proto.String(time.Now().Add(1 * time.Hour).Format("2006-01-02 15:04:05.000")),
		}
		err := s.Add("", log)
		tt.Log(err)
	})
}

func TestAPIGameLog_BatchAdd(t *testing.T) {
	s := new(APIGameLogService)
	t.Run("NewLogOne", func(tt *testing.T) {
		log := &netproto.Slot_MSG_GameLogAdd{
			UserID:        proto.Int32(1000174),
			GameID:        proto.Int32(1011),
			ServerID:      proto.Int32(475),
			CurrentAmount: proto.Int64(100),
			BetMoney:      proto.Int64(10),
			WinMoney:      proto.Int64(30),
			FinalScore:    proto.Int64(120),
			GameType:      proto.Int32(1),
			LogType:       proto.Int32(1),
			GameData:      proto.String("[test batch one]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
			CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
		}
		err := s.BatchAdd("", []*netproto.Slot_MSG_GameLogAdd{log})
		tt.Log(err)
	})
	t.Run("NewLog Free", func(tt *testing.T) {
		log1 := &netproto.Slot_MSG_GameLogAdd{
			UserID:        proto.Int32(1000174),
			GameID:        proto.Int32(1011),
			ServerID:      proto.Int32(475),
			CurrentAmount: proto.Int64(100),
			BetMoney:      proto.Int64(10),
			WinMoney:      proto.Int64(30),
			FinalScore:    proto.Int64(120),
			GameType:      proto.Int32(1),
			LogType:       proto.Int32(1),
			GameData:      proto.String("[test batch 1]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
			CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
		}
		log2 := &netproto.Slot_MSG_GameLogAdd{
			UserID:        proto.Int32(1000174),
			GameID:        proto.Int32(1011),
			ServerID:      proto.Int32(475),
			CurrentAmount: proto.Int64(100),
			BetMoney:      proto.Int64(10),
			WinMoney:      proto.Int64(30),
			FinalScore:    proto.Int64(120),
			GameType:      proto.Int32(1),
			LogType:       proto.Int32(1),
			GameData:      proto.String("[test batch 2]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
			CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
		}
		log3 := &netproto.Slot_MSG_GameLogAdd{
			UserID:        proto.Int32(1000174),
			GameID:        proto.Int32(1011),
			ServerID:      proto.Int32(475),
			CurrentAmount: proto.Int64(100),
			BetMoney:      proto.Int64(10),
			WinMoney:      proto.Int64(30),
			FinalScore:    proto.Int64(120),
			GameType:      proto.Int32(1),
			LogType:       proto.Int32(1),
			GameData:      proto.String("[test batch 3]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
			CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
		}
		err := s.BatchAdd("", []*netproto.Slot_MSG_GameLogAdd{log1, log2, log3})
		tt.Log(err)
	})
}

func TestAPIGameLog_Search(t *testing.T) {
	s := new(APIGameLogService)
	t.Run("No User", func(tt *testing.T) {
		pageCnt, logs, err := s.Search("", 1, 1001, 1, 10, "2019-10-04", "2019-10-04")
		tt.Log(err)
		tt.Logf("All Row Count: %d\n", pageCnt)
		tt.Logf("Logs Length: %d\n", len(logs))
		for _, row := range logs {
			tt.Logf(" - %v\n", row)
		}
	})
	t.Run("No Game", func(tt *testing.T) {
		pageCnt, logs, err := s.Search("", 9936738, 999, 1, 10, "2019-10-04", "2019-10-04")
		tt.Log(err)
		tt.Logf("All Row Count: %d\n", pageCnt)
		tt.Logf("Logs Length: %d\n", len(logs))
		for _, row := range logs {
			tt.Logf(" - %v\n", row)
		}
	})
	t.Run("Pass", func(tt *testing.T) {
		pageCnt, logs, err := s.Search("", 1000174, 1011, 1, 10, "2020-03-20", "2020-03-31")
		tt.Log(err)
		tt.Logf("All Row Count: %d\n", pageCnt)
		tt.Logf("Logs Length: %d\n", len(logs))
		for _, row := range logs {
			tt.Logf(" - %v\n", row)
		}
	})
}

func TestQQ(t *testing.T) {
	now := time.Now()

	fmt.Println(now)
	fmt.Println(now.UTC())

	fmt.Println(now.UnixNano())
	fmt.Println(now.UTC().UnixNano())

	fmt.Println(now.Format(time.RFC3339Nano))
	fmt.Println(now.UTC().Format(time.RFC3339Nano))

	q, _ := time.Parse(time.RFC3339Nano, now.UTC().Format(time.RFC3339Nano))

	fmt.Println(q)
	fmt.Println(q.UnixNano())
	fmt.Println(q.UTC().Format(time.RFC3339Nano))
	fmt.Println(q.Local().Format(time.RFC3339Nano))

	jt := model.JSONTime{Time: now}
	{
		b, _ := jt.MarshalJSON()
		fmt.Println(string(b))
	}

	data := map[string]interface{}{
		"b":    1,
		"time": jt,
	}

	b, err := json.Marshal(data)

	fmt.Println(b, err, string(b))

	now.In(time.UTC)
	fmt.Println(now.Zone())
	d := time.Duration(28800 * time.Second)
	fmt.Println(d)
}
