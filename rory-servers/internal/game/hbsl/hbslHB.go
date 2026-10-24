package hbsl

type HB struct {
	i64Chip       int64 // 紅包金額
	i32BombNumber int32 // 雷號
	i32UserID     int32 // 發紅包用户ID
	i32FaceID     int32 // 發紅包用户头像ID
	pUser         *User // 發紅包玩家資訊
}

type HBRecord struct {
	pHB              *HB           // 紅包資料
	ary_i64SHB       []int64       // 子紅包
	ary_pUserRecord  []*UserRecord // 搶紅包玩家紀錄
	i64ZhuangWinLose int64         // 莊家輸贏
}

// 清除HB數據
func (pHB *HB) Init() {
	pHB.pUser = nil
}

// 取得搶紅包紀錄
func (pHBRecord *HBRecord) GetUserRecordByID(i32UserID int32) *UserRecord {
	for i := 0; i < len(pHBRecord.ary_pUserRecord); i++ {
		if pHBRecord.ary_pUserRecord[i].iUserID == i32UserID {
			return pHBRecord.ary_pUserRecord[i]
		}
	}
	return nil
}
