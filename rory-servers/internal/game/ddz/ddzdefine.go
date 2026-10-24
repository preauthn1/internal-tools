package ddz

//游戏角色类型
type playerRole byte

//游戏角色定义
const (
	roleDiZhu   playerRole = iota //地主
	roleNongMin                   //农民
)

//叫地主类型
type JiaoDiZhuEnum byte

const (
	JDZNone   JiaoDiZhuEnum = iota //none
	JDZJiao                        //叫地主
	JDZBuJiao                      //不叫
)

//输赢类型
type loseEnum byte

//输赢标识
const (
	loseEnumNormal     loseEnum = iota //普通输赢
	loseEnumBankruptcy                 //破产
	loseEnumCap                        //封顶
)

type OddsType int32

const (
	DdzOddsTypeBase OddsType = iota
	DdzOddsTypeJdz
	DdzOddsTypeQdz
	DdzOddsTypeMingPai
	DdzOddsTypeZhaDan
	DdzOddsTypeChunTian
	DdzOddsTypeDiZhuPai
	DdzOddsTypeRocket
	DdzOddsTypeFanChunTian
	DdzOddsTypeMax
)

//游戏阶段类型
type gameProgress int32

//游戏阶段定义
const (
	gameProgressNone      gameProgress = 1
	gameProgressBegin                  = 2 //游戏开始
	gameProgressFapai                  = 3 //发牌阶段
	gameProgressJiaoDiZhu              = 4 //叫牌
	gameProgressRob                    = 5 //抢地主
	gameProgressShow                   = 6 //明牌
	gameProgressStart                  = 7 //出牌
	gameProgressGameEnd                = 8 //结束结算
)

//牌型类型
type cardType byte

//牌型定义
const (
	cardTypeNone            cardType = iota //无牌型
	cardTypeRocket                          //火箭1
	cardTypeBomb                            //炸弹2
	cardTypeSingle                          //单牌3
	cardTypeDouble                          //对子4
	cardTypeThree                           //三张5
	cardTypeThreeAndOne                     //三带一6
	cardTypeStraight                        //单顺（顺子）7
	cardTypeDoubleStraights                 //飞机带翅膀 单8
	cardTypeDoubleStraightd                 //飞机带翅膀 双9
	cardTypeThreeStraight                   //連對(對子順) 10
	//三順 (三條順子)
	cardTypePlaneAndWing //飞机带翅膀 11
	cardTypeFourAndTwod  //四带二双 12
	cardTypeFourAndTwos  //四带二单 13
	cardTypeThreeAndTwo  //三帶二
)

const (
	calMethod       = "calmethod"
	oddsBase        = "base" //炸弹
	oddsMingpai     = "mingpai"
	oddsJiaoDiZhu   = "jiaodizhu"
	oddsQiangDiZhu  = "qiangdizhu"
	oddsZhaDan      = "zhadan"
	oddsChunTian    = "chuntian"
	oddsSanPai      = "sanpai"
	oddsDuiZi       = "duizi"
	oddsShunZi      = "shunzi"
	oddsTongHua     = "tonghua"
	oddsXiaoWang    = "xiaowang"
	oddsDaWang      = "dawang"
	oddsBaoZi       = "baozi"
	oddsRocket      = "rocket"
	oddsFanChunTian = "fanchuntian"
)

const (
	outCardLessCountWeight = "outcardlesscount" //出牌次數最少
	bigKingWeight          = "bigking"
	smallKingWeight        = "smallking"
	cardTwoWeight          = "cardtwo"
	rocketWeight           = "rocket"
	bombWeight             = "bomb"
	doubleStraightWeight   = "doublestraight"
	straightWeight         = "straight"
	lianduiWeight          = "liandui"
)
