package gamelogic

import (
	"bytes"
	"fmt"
	"sort"

	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// Room 房間
type Room struct {
	ServerID      int32    // 服務器ID
	ServerName    string   // 服務器名稱
	Flag          string   // 標記
	GameType      int32    // 0真金 1比賽 2休閒
	GameInfo      *Game    // 遊戲信息
	Tables        []*Table // 桌子
	SeatCount     int32    // 座位數
	TableCount    int32    // 桌子數
	LoginMoney    int32    // 進入房間最低金幣數
	MaxUser       int32    // 最大人數
	BaseScore     int32    // 房間底分
	Single        int32    // 直接進入取第一個服務器進入, 0=>顯示房間列表, 1=>第一個服務器進入
	HostAndBackup int32    // AB服. 0主機 1備援轉正主機接手中 2備援
	KickOutMoney  int32    // 達到此積分，強制退出房間
	MaxOdd        int32    // 最大倍數

	RoomUserManager   *RoomUserManager  // 用户管理器
	IPAddress         string            // 服務器IP
	Port              int32             // 服務器端口
	*Fish                               // 捕魚放分資料
	SuperToolService  ISuperToolService // 超端服務
	NewGameServerAddr string            // 切AB服後, 新的B服game ip
	*netproto.FishSysPoolData
	HsiTungHuoLiIRate            float32                                      // 系統贏錢抽百分比
	HsiTungHuoLiChouShuiChingFen int32                                        // 系統獲利當次抽水淨分
	FishPlayerRewardDetail       map[int32]*netproto.FishPlayerPoolDetailList // 玩家個人水池詳細記錄
	NewPlayerPoolRecord          map[int32]int32                              // 第一次玩的新手記錄, 第一次進遊戲系統水池會先扣0.7, 離開出去會將低於0.7的剩餘金額加回系統水池
	FishPlayerPoolStatus         bool
	FisnNewPlayerDetail          map[int32]*netproto.FishBeginnerOffer // 新手優惠
	Seq                          *int64                                // 唯一流水號
	OnPiSaiUserInfo              map[int32]*PiSaiUserInfo              // 是否在進行比賽中
}

// PiSaiUserInfo ?
type PiSaiUserInfo struct {
	UserID      int32
	IsPlaying   bool  // 報名參加比賽是否成功
	IsGameStart bool  // 可進行遊戲
	TableNo     int32 // 桌子
	EndTime     int64 // 比賽結束時間
	Pause       bool  // 斷線或離開將暫時
	Score       int64 // 總積分
	Bonus       int64 // 獲得
	ContestNo   int32 // 第幾局
}

// NewRoom 創建房間對象
func NewRoom(serverCfg *GameServerConfig, roomUserManager *RoomUserManager) *Room {
	ip, port := serverCfg.GetServerIPPort()

	r := new(Room)
	r.ServerID = serverCfg.ServerID
	r.ServerName = serverCfg.ServerName
	r.Flag = serverCfg.Flag
	r.GameType = serverCfg.GameType
	r.SeatCount = serverCfg.SeatCount
	r.TableCount = serverCfg.TableCount
	r.IPAddress = ip
	r.Port = int32(port)
	r.RoomUserManager = roomUserManager
	r.LoginMoney = serverCfg.LoginMoney
	r.MaxUser = serverCfg.MaxUser
	r.BaseScore = serverCfg.BaseScore
	r.Single = serverCfg.Single
	r.HostAndBackup = serverCfg.HostAndBackup
	r.KickOutMoney = serverCfg.Kickoutmoney
	r.MaxOdd = serverCfg.MaxOdd
	r.GameInfo = serverCfg.GetGameInfo()
	r.Fish = new(Fish)

	var i int32
	for i = 0; i < r.TableCount; i++ {
		t := NewTable(r.SeatCount, i)
		t.Room = r

		r.Tables = append(r.Tables, t)
	}

	// 添加房間命令
	RoomCmdP = NewRoomCmd(r)
	command.AddCommand(RoomCmdP)

	return r
}

//func (r *Room) SetSuperToolService(service ISuperToolService) {
//	r.SuperToolService = service
//}

// GetTableByNo 根據桌子號獲取桌子
func (r *Room) GetTableByNo(tableno int32) *Table {
	for _, t := range r.Tables {
		if t.TableNo == tableno {
			return t
		}
	}

	return nil
}

// GetSeat 獲取椅子
func (r *Room) GetSeat(tableno int32, seatno int32) *Seat {
	t := r.GetTableByNo(tableno)

	if t != nil {
		return t.GetSeatByNo(seatno)
	}

	return nil
}

// GetTableUser 獲取桌子上的用戶
func (r *Room) GetTableUser(tableno int32) []*GameUser {
	t := r.GetTableByNo(tableno)
	if t == nil {
		return nil
	}

	return t.GetTableUser()
}

func (r *Room) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("serverid:%d,servername:%s,tablecount:%d,seatcount:%d,game:%v\r\n", r.ServerID, r.ServerName, r.TableCount, r.SeatCount, r.GameInfo))
	for _, t := range r.Tables {
		buf.WriteString(t.String())
	}

	return buf.String()
}

// GetSortedTable ?
func (r *Room) GetSortedTable() []*Table {
	st := newTableSorter(r.Tables)
	sort.Sort(st)

	return st.tables
}

// GetRobotSortedTable ?
func (r *Room) GetRobotSortedTable() []*Table {
	st := newRobotTableSorter(r.Tables)
	sort.Sort(st)

	return st.tables
}

// GetRoomRealUserCount ?
func (r *Room) GetRoomRealUserCount() int32 {
	var realUserTotal int32
	for _, table := range r.Tables {
		realUserTotal += int32(len(table.GetTableRealUser()))
	}
	return realUserTotal
}

// 默認桌子排序
type tableSorter struct {
	tables []*Table
}

func newTableSorter(tables []*Table) *tableSorter {
	ts := new(tableSorter)
	ts.tables = make([]*Table, 0)
	ts.tables = append(ts.tables, tables...)

	return ts
}

func (ts *tableSorter) Len() int {
	return len(ts.tables)
}

func (ts *tableSorter) Swap(i, j int) {
	ts.tables[i], ts.tables[j] = ts.tables[j], ts.tables[i]
}

func (ts *tableSorter) Less(i, j int) bool {
	leni := len(ts.tables[i].GetTableUser())
	lenj := len(ts.tables[j].GetTableUser())

	if leni != lenj {
		return leni > lenj
	}

	return ts.tables[i].TableNo < ts.tables[j].TableNo
}

// 機器人桌子排序
type robotTableSorter struct {
	tables []*Table
}

func newRobotTableSorter(tables []*Table) *robotTableSorter {
	ts := new(robotTableSorter)
	ts.tables = make([]*Table, 0)
	ts.tables = append(ts.tables, tables...)

	return ts
}

func (ts *robotTableSorter) Len() int {
	return len(ts.tables)
}

func (ts *robotTableSorter) Swap(i, j int) {
	ts.tables[i], ts.tables[j] = ts.tables[j], ts.tables[i]
}

func (ts *robotTableSorter) Less(i, j int) bool {
	tablei := ts.tables[i]
	tablej := ts.tables[j]

	if tablei.TableUserCount() == 0 {
		return false
	}

	leni := tablei.TableRealUserCount()
	lenj := tablej.TableRealUserCount()

	if !tablei.GameLogic.IsGaming() {
		if tablej.GameLogic.IsGaming() {
			return true
		}
		return leni > lenj
	}
	if tablej.GameLogic.IsGaming() {
		return tablei.TableUserCount() > tablej.TableUserCount()
	}
	if tablej.TableUserCount() == 0 {
		return true
	}
	return false

	// return ts.tables[i].TableNo < ts.tables[j].TableNo
}
