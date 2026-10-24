package gamelogic

// ActionController 用戶行為控制器
type ActionController interface {
	UserSit(user *GameUser, tableno int32, seatno int32)                                                                             // 用戶坐下
	UserFree(user *GameUser)                                                                                                         // 用戶站立
	UserReady(user *GameUser)                                                                                                        // 用戶舉手
	GetQueueUser() []*GameUser                                                                                                       // 獲取排隊用戶
	TotalUserCount() int32                                                                                                           // 房間總用戶數
	GetFreeTable() *Table                                                                                                            // 取得一個空閒桌子
	GetFreeSeat(user *GameUser, exIP string, exLastGameTableNo int32) (int32, int32)                                                 // 獲取空閒的座位
	GetFreeSeatNew(user *GameUser, exIP string, exLastGameTableNo int32) (int32, int32)                                              // 新的獲取空閒的座位
	GetWinLoseFreeSeat(extableno int32, exIP string, exLastGameTableNo int32, lastWinLose int32, islessthanuser bool) (int32, int32) // 輸贏分隊獲取空閒座位
	GetFreeSeatWithUPM(user *GameUser, exIP string, exLastGameTableNo int32) (int32, int32)                                          // 排除同桌率取得空位
	GetNextTableSeat(user *GameUser) (int32, int32)                                                                                  // 取得指定下局桌空位
}
