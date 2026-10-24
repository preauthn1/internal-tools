

enum DFDC_GameMessageClassID {
    DFDCSettingID = 10001,  //设置相关 S2C
    DFDCGameStartID = 10002,    //遊戲開始C2S
    DFDCGameResult = 10003,     //遊戲結果S2C
    DFDCStateID = 10004,    //遊戲狀態S2C
    DFDCFreeOrPayID = 10005,    //遊戲要付費或免費
    DFDCGetJackpotPoolID = 10006,       //取得jackpot pool 水池
    DFDCGetJackpotPoolRetID = 10007,        //回傳jackpot pool 水池
    DFDCUpdateJackpotPoolID = 10008,        //更新jackpot pool水池
    DFDCJackpotPoolS2CID = 10009,   //S2C 更新Client Jackpot 值
    DFDCBetRetID = 10010,   //S2C下注結果
    DFDCGoldenSymbolID = 10011, //S2C 黃金符號列表 (DFDC_GoldenSymbolList)
    DFDCBonusGameID = 10012,    //S2C 進入Bonus Game (DFDC_BonusGame)
    DFDCBonusGameIconID = 10013,	//C2S 玩家Bonus Game選擇 (DFDC_BonusGameDetail)
};
