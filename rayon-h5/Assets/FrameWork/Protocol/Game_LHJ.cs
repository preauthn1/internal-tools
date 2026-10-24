

enum LHJ_GameMessageClassID {
    LHJSettingID = 10001,   //设置相关 S2C
    LHJGameStartID = 10002, //遊戲開始C2S
    LHJGameResult = 10003,      //遊戲結果S2C
    LHJStateID = 10004, //遊戲狀態S2C
    LHJFreeOrPayID = 10005, //遊戲要付費或免費
    LHJGetJackpotPoolID = 10006,        //取得jackpot pool 水池
    LHJGetJackpotPoolRetID = 10007,     //回傳jackpot pool 水池
    LHJUpdateJackpotPoolID = 10008,     //更新jackpot pool水池
    LHJJackpotPoolS2CID = 10009,	//S2C 更新Client Jackpot 值
    LHJBetRetID = 10010 //S2C下注結果
};
