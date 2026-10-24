

enum GameMessageClassID {
	TestGameChatID = 1,
	TestGameScoreID = 2,
	TestGameEndID = 3,
	UserEnterID=4,
	UserLeaveID=5,
	UserStateChangeID=6,
	GameStartID=7,
	GameEndID=8,
	UserPokerID=9,
	GamePokerID=10,
	UserActionCommitID=11,
	AllUserActionCommitID=12,
	UserRecoverID=13,
	LookOnUserNotifyID=14,
	UserLuckScoreID=15,
	UserLuckScoreRetID=16,
	GetTableBlackUserID=17,
	GetTableBlackUserRetID=18,
	UpdateTableBlackUserID=19,
	UpdateTableBlackUserRetID=20,
	GetTableWhiteUserID=21,
	GetTableWhiteUserRetID=22,
	UpdateTableWhiteUserID=23,
	UpdateTableWhiteUserRetID=24,
	GameBeginCountDownID=25,
	GameBeginCountDownStopID=26,

	SuperToolAddChangeLogID=27,
	SuperToolAssignRewardID=28,
	SuperToolGetAssignRewardDetailID=29,
	SuperToolGetAssignRewardDetailRetID=30,
	SuperToolUpdAssignRewardDetailID=31,
	SuperToolGetVipUserID=32,
	SuperToolGetVipUserRetID=33,
	SuperToolUpdVipUserID=34,
	SuperToolFishGetSystemRewardDetailID=35,
	SuperToolFishGetSystemRewardDetailRet=36,
	SuperToolFishUpdSystemRewardID=37,
	SuperToolFishUpdSystemRewardRetID=38,
	SuperToolFishUpdVIPUserRewardID=39,
	SuperToolFishUpdVIPUserRewardRetID=40,
	SuperToolZjhGetVipUserID=41,
	SuperToolZjhGetVipUserRetID=42,
    
    // Server to Client 20xxx
    Slot_SettingID = 20001,  // 設置相關
    Slot_StatusID = 20002,   // 遊戲狀態 0: 玩家可下注, 1: 開獎, 2: 開獎動畫, 3: 結算動畫, 4: 等待遊戲開始
    Slot_BetStatusID = 20003,  // 下注狀態 0: 下注成功, -1: 用戶金額不足, -2: 下注的類型未找到, -3: 莊家不能下注, -4: 此類下注已達到上限, -5: 非下注階段, -6: 用戶金額小於下注限制
    Slot_GameResultID = 20004,  // 遊戲結果 (轉出內容，中獎項目，中獎金額)

    // Client to Server 21xxx
    Slot_ActionSpinID = 21001,  // 玩家Spin

    // 會一來一回的資訊 22xxx
    Slot_LoadHistoryID = 22001,     // 玩家查詢記錄
    Slot_LoadHistoryRetID = 22002,  // 回傳玩家查詢記錄
    Slot_JackpotPoolID = 22003,     // 查詢獎池資訊
    Slot_JackpotPoolRetID = 22004,  // 回傳獎池資訊
    Slot_BonusGameID = 22005,       // 查詢BonusGame資訊
    Slot_BonusGameRetID = 22006,    // 回傳查詢BonusGame資訊
};
