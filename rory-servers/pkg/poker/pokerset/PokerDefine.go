package pokerset

// ?
const (
	POKEROneColorCardsCount    int32 = 13                                         // 一種花色有幾張牌
	POKERColorsCount           int32 = 4                                          // 幾種花色
	POKERSetCardsCount         int32 = POKERColorsCount * POKEROneColorCardsCount // 一副牌的數量(不含Joker)
	POKERJokerCardsCount       int32 = 2                                          // Joker數量
	POKERMaxSet                int32 = 10                                         // 最多用幾副牌洗牌
	POKERMinStraightCount      int   = 2                                          // 檢查順面所需最小張數
	POKERPatternBasicCardCount int   = 5                                          // 牌型基本張數
)

// ?
const (
	LENSingleCardCode int = 2 // 一張牌的代碼文字長度
)

// 比牌結果
const (
	POKERCompareEven   int = iota // 平手
	POKERCompareDealer            // 莊家贏
	POKERComparePlayer            // 玩家贏
)

// 牌型
const (
	POKERPatternNone          int32 = iota // 無
	POKERPatternHighCard                   // 高牌
	POKERPatternPair                       // 對子
	POKERPatternTwoPairs                   // 兩對
	POKERPatternThreeOfAKind               // 三條
	POKERPatternStraight                   // 順子
	POKERPatternFlush                      // 同花
	POKERPatternFullHouse                  // 葫蘆
	POKERPatternFourOfAKind                // 四條
	POKERPatternStraightFlush              // 同花順
	POKERPatternRoyalFlush                 // 皇家同花順
	POKERPatternDragon                     // 一條龍
	POKERPatternDragonFlush                // 同花一條龍
	POKERPatternMax
)

// 花色
const (
	POKERCardColorClub    int32 = iota // 梅花
	POKERCardColorDiamond              // 方塊
	POKERCardColorHeart                // 紅桃
	POKERCardColorSpade                // 黑桃
	POKERCardColorOther                // 其他
	POKERCardColorMax
)

// 數字
const (
	POKERCardNoNone int32 = iota
	POKERCardNo1
	POKERCardNo2
	POKERCardNo3
	POKERCardNo4
	POKERCardNo5
	POKERCardNo6
	POKERCardNo7
	POKERCardNo8
	POKERCardNo9
	POKERCardNo10
	POKERCardNoJ
	POKERCardNoQ
	POKERCardNoK
	POKERCardNoA
	POKERCardNoJoker1 // 小王
	POKERCardNoJoker2 // 大王
	POKERCardNoHide
	POKERCardNoMax
)

// 牌的資料以編號形式表達
const (
	POKERCardIDNone int32 = iota
	// Club
	POKERCardIDC1
	POKERCardIDC2
	POKERCardIDC3
	POKERCardIDC4
	POKERCardIDC5
	POKERCardIDC6
	POKERCardIDC7
	POKERCardIDC8
	POKERCardIDC9
	POKERCardIDC0
	POKERCardIDCJ
	POKERCardIDCQ
	POKERCardIDCK
	// Diamond
	POKERCardIDD1
	POKERCardIDD2
	POKERCardIDD3
	POKERCardIDD4
	POKERCardIDD5
	POKERCardIDD6
	POKERCardIDD7
	POKERCardIDD8
	POKERCardIDD9
	POKERCardIDD0
	POKERCardIDDJ
	POKERCardIDDQ
	POKERCardIDDK
	// Heart
	POKERCardIDH1
	POKERCardIDH2
	POKERCardIDH3
	POKERCardIDH4
	POKERCardIDH5
	POKERCardIDH6
	POKERCardIDH7
	POKERCardIDH8
	POKERCardIDH9
	POKERCardIDH0
	POKERCardIDHJ
	POKERCardIDHQ
	POKERCardIDHK
	// Spade
	POKERCardIDS1
	POKERCardIDS2
	POKERCardIDS3
	POKERCardIDS4
	POKERCardIDS5
	POKERCardIDS6
	POKERCardIDS7
	POKERCardIDS8
	POKERCardIDS9
	POKERCardIDS0
	POKERCardIDSJ
	POKERCardIDSQ
	POKERCardIDSK
	// Joker
	POKERCardIDJ1
	POKERCardIDJ2
	// 蓋牌
	POKERCardIDHide
	POKERCardIDMax
)

//-----------------------------------------------------------------------------
// 花色字串
var gMapStrCardColor = map[int32]string{
	POKERCardColorClub:    "梅花",
	POKERCardColorDiamond: "方块",
	POKERCardColorHeart:   "红桃",
	POKERCardColorSpade:   "黑桃",
	POKERCardColorOther:   "",
}

// 數字字串
var gMapStrCardNumber = map[int32]string{
	POKERCardNo1:      "1",
	POKERCardNo2:      "2",
	POKERCardNo3:      "3",
	POKERCardNo4:      "4",
	POKERCardNo5:      "5",
	POKERCardNo6:      "6",
	POKERCardNo7:      "7",
	POKERCardNo8:      "8",
	POKERCardNo9:      "9",
	POKERCardNo10:     "10",
	POKERCardNoJ:      "J",
	POKERCardNoQ:      "Q",
	POKERCardNoK:      "K",
	POKERCardNoA:      "A",
	POKERCardNoJoker1: "小王",
	POKERCardNoJoker2: "大王",
	POKERCardNoHide:   "盖牌",
}

// 花色代碼
var gMapStrColorCode = map[int32]string{
	POKERCardColorClub:    "c",
	POKERCardColorDiamond: "d",
	POKERCardColorHeart:   "h",
	POKERCardColorSpade:   "s",
	POKERCardColorOther:   "j",
}

var gMapI32ColorByCode = map[string]int32{
	"c": POKERCardColorClub,
	"d": POKERCardColorDiamond,
	"h": POKERCardColorHeart,
	"s": POKERCardColorSpade,
	"j": POKERCardColorOther,
}

// 數字代碼
var gMapStrNumberCode = map[int32]string{
	POKERCardNoNone:   "d",
	POKERCardNo1:      "1",
	POKERCardNo2:      "2",
	POKERCardNo3:      "3",
	POKERCardNo4:      "4",
	POKERCardNo5:      "5",
	POKERCardNo6:      "6",
	POKERCardNo7:      "7",
	POKERCardNo8:      "8",
	POKERCardNo9:      "9",
	POKERCardNo10:     "0",
	POKERCardNoJ:      "J",
	POKERCardNoQ:      "Q",
	POKERCardNoK:      "K",
	POKERCardNoA:      "A",
	POKERCardNoJoker1: "a",
	POKERCardNoJoker2: "b",
	POKERCardNoHide:   "c",
}

var gMapI32NumberByCode = map[string]int32{
	"d": POKERCardNoNone,
	"1": POKERCardNo1,
	"2": POKERCardNo2,
	"3": POKERCardNo3,
	"4": POKERCardNo4,
	"5": POKERCardNo5,
	"6": POKERCardNo6,
	"7": POKERCardNo7,
	"8": POKERCardNo8,
	"9": POKERCardNo9,
	"0": POKERCardNo10,
	"J": POKERCardNoJ,
	"Q": POKERCardNoQ,
	"K": POKERCardNoK,
	"A": POKERCardNoA,
	"a": POKERCardNoJoker1,
	"b": POKERCardNoJoker2,
	"c": POKERCardNoHide,
}

//牌型字串
var gMapStrPattern = map[int32]string{
	POKERPatternNone:          "无",
	POKERPatternHighCard:      "高牌",
	POKERPatternPair:          "对子",
	POKERPatternTwoPairs:      "两对",
	POKERPatternThreeOfAKind:  "三条",
	POKERPatternStraight:      "顺子",
	POKERPatternFlush:         "同花",
	POKERPatternFullHouse:     "葫芦",
	POKERPatternFourOfAKind:   "四条",
	POKERPatternStraightFlush: "同花顺",
	POKERPatternRoyalFlush:    "皇家同花顺",
	POKERPatternDragon:        "一条龙",
	POKERPatternDragonFlush:   "同花一条龙",
}

//-----------------------------------------------------------------------------
// 各牌型張數
var gMapIPatternCardCount = map[int32]int{
	POKERPatternNone:          0,  // 無
	POKERPatternHighCard:      1,  // 高牌
	POKERPatternPair:          2,  // 對子
	POKERPatternTwoPairs:      4,  // 兩對
	POKERPatternThreeOfAKind:  3,  // 三條
	POKERPatternStraight:      3,  // 順子
	POKERPatternFlush:         3,  // 同花
	POKERPatternFullHouse:     5,  // 葫蘆
	POKERPatternFourOfAKind:   4,  // 四條
	POKERPatternStraightFlush: 3,  // 同花順
	POKERPatternRoyalFlush:    5,  // 皇家同花順
	POKERPatternDragon:        13, // 一條龍
	POKERPatternDragonFlush:   13, // 同花一條龍
}

//-----------------------------------------------------------------------------
