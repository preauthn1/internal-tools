package ddz

import (
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

type ddzTimeOutCfg struct {
	tc map[gameProgress]int32 //阶段超时配置
}

//获取游戏阶段对应的等待超时时间
func (ddt *ddzTimeOutCfg) GetTime(gp gameProgress) int32 {
	return ddt.tc[gp]
}

//初始化超时配置
func ddzTimeOutInit() *ddzTimeOutCfg {
	ddzt := new(ddzTimeOutCfg)
	ddzt.tc = map[gameProgress]int32{
		gameProgressBegin:     0,
		gameProgressFapai:     0,
		gameProgressJiaoDiZhu: 10,
		gameProgressRob:       10,
		gameProgressShow:      10,
		gameProgressStart:     30,
		gameProgressGameEnd:   10,
	}
	return ddzt
}

//倍率配置
type ddzOddsCfg struct {
	odds      map[string]int //倍率配置信息
	calMethod int            //倍數計算方式
}

//換牌的牌型積分權重
type ddzExchangeCardWeight struct {
	exchange map[string]int //積分
}

//不洗牌規則
type ddzNoShufflePoker struct {
	groupPoker map[int]int //group by, 相同點數幾張一組
	alloctake  map[int]int //分牌一次取幾張
}

func (exeCfg *ddzNoShufflePoker) init(cfgs *config.Group) {
	cgf := cfgs.GetConfig("ddz")
	strGroupPoker := cgf.GetConfigValue("NoShufflePoker", "grouppoker")
	group := strings.Split(strGroupPoker, ",")
	for k, v := range group {
		num, _ := strconv.Atoi(v)
		exeCfg.groupPoker[k] = num
	}

	strAllocTake := cgf.GetConfigValue("NoShufflePoker", "alloctake")
	take := strings.Split(strAllocTake, ",")
	for k, v := range take {
		num, _ := strconv.Atoi(v)
		exeCfg.alloctake[k] = num
	}
}

func NewDdzNoShufflePoker(cgf *config.Group) *ddzNoShufflePoker {
	d := new(ddzNoShufflePoker)
	d.groupPoker = make(map[int]int)
	d.alloctake = make(map[int]int)
	d.init(cgf)
	return d
}

func (excCfg *ddzExchangeCardWeight) init(cfgs *config.Group) {
	cgf := cfgs.GetConfig("ddz")
	excCfg.exchange[outCardLessCountWeight] = cgf.GetConfigInt("ExchangeCardWeight", outCardLessCountWeight)
	excCfg.exchange[bigKingWeight] = cgf.GetConfigInt("ExchangeCardWeight", bigKingWeight)
	excCfg.exchange[smallKingWeight] = cgf.GetConfigInt("ExchangeCardWeight", smallKingWeight)
	excCfg.exchange[cardTwoWeight] = cgf.GetConfigInt("ExchangeCardWeight", cardTwoWeight)
	excCfg.exchange[rocketWeight] = cgf.GetConfigInt("ExchangeCardWeight", rocketWeight)
	excCfg.exchange[bombWeight] = cgf.GetConfigInt("ExchangeCardWeight", bombWeight)
	excCfg.exchange[doubleStraightWeight] = cgf.GetConfigInt("ExchangeCardWeight", doubleStraightWeight)
	excCfg.exchange[straightWeight] = cgf.GetConfigInt("ExchangeCardWeight", straightWeight)
	excCfg.exchange[lianduiWeight] = cgf.GetConfigInt("ExchangeCardWeight", lianduiWeight)
}

func NewDdzExchangeCardWeight(cfg *config.Group) *ddzExchangeCardWeight {
	d := new(ddzExchangeCardWeight)
	d.exchange = make(map[string]int)
	d.init(cfg)
	return d
}

//檢查單牌類型好牌的權重
func (decw *ddzExchangeCardWeight) CalSingleCardsWeight(cg *poker.CardGroup) int {
	result := int(0)

	for _, v := range cg.Cards {
		switch v.CardNumber {
		case poker.PokerNumberBigKing:
			result += decw.exchange[bigKingWeight]
		case poker.PokerNumberSmallKing:
			result += decw.exchange[smallKingWeight]
		case poker.PokerNumberTwo:
			result += decw.exchange[cardTwoWeight]
		}
	}

	return result
}

//檢查組合類型好牌的權重
func (cg *ddzExchangeCardWeight) CalCardGroupWeight(cfg *CtGroup) int {
	result := int(0)

	for _, v := range cfg.Ct {
		switch v.byteCardType {
		case cardTypeRocket:
			result += cg.exchange[rocketWeight]
		case cardTypeBomb:
			result += cg.exchange[bombWeight]
		case cardTypeDoubleStraights, cardTypeDoubleStraightd:
			result += cg.exchange[doubleStraightWeight]
		case cardTypeStraight:
			result += cg.exchange[straightWeight]
		case cardTypeThreeStraight:
			result += cg.exchange[lianduiWeight]
		}
	}

	return result
}

func (odscfg *ddzOddsCfg) init(pokerType int32, cfgs *config.Group) {
	cfg := cfgs.GetConfig("ddz")
	if pokerType == gamelogic.PokerTypeShuffle { //不洗牌玩法
		odscfg.calMethod = cfg.GetConfigInt("Odds", calMethod)
		odscfg.odds[oddsBase] = cfg.GetConfigInt("Odds", oddsBase)
		odscfg.odds[oddsMingpai] = cfg.GetConfigInt("Odds", oddsMingpai)
		odscfg.odds[oddsJiaoDiZhu] = cfg.GetConfigInt("Odds", oddsJiaoDiZhu)
		odscfg.odds[oddsQiangDiZhu] = cfg.GetConfigInt("Odds", oddsQiangDiZhu)
		odscfg.odds[oddsZhaDan] = cfg.GetConfigInt("Odds", oddsZhaDan)
		odscfg.odds[oddsChunTian] = cfg.GetConfigInt("Odds", oddsChunTian)
		odscfg.odds[oddsFanChunTian] = cfg.GetConfigInt("Odds", oddsFanChunTian)
		odscfg.odds[oddsRocket] = cfg.GetConfigInt("Odds", oddsRocket)

		odscfg.odds[oddsSanPai] = cfg.GetConfigInt("Odds", oddsSanPai)
		odscfg.odds[oddsDuiZi] = cfg.GetConfigInt("Odds", oddsDuiZi)
		odscfg.odds[oddsShunZi] = cfg.GetConfigInt("Odds", oddsShunZi)
		odscfg.odds[oddsTongHua] = cfg.GetConfigInt("Odds", oddsTongHua)
		odscfg.odds[oddsXiaoWang] = cfg.GetConfigInt("Odds", oddsXiaoWang)
		odscfg.odds[oddsDaWang] = cfg.GetConfigInt("Odds", oddsDaWang)
		odscfg.odds[oddsBaoZi] = cfg.GetConfigInt("Odds", oddsBaoZi)
	} else {
		odscfg.calMethod = cfg.GetConfigInt("OddsClassic", calMethod)
		odscfg.odds[oddsBase] = cfg.GetConfigInt("OddsClassic", oddsBase)
		odscfg.odds[oddsMingpai] = cfg.GetConfigInt("OddsClassic", oddsMingpai)
		odscfg.odds[oddsJiaoDiZhu] = cfg.GetConfigInt("OddsClassic", oddsJiaoDiZhu)
		odscfg.odds[oddsQiangDiZhu] = cfg.GetConfigInt("OddsClassic", oddsQiangDiZhu)
		odscfg.odds[oddsZhaDan] = cfg.GetConfigInt("OddsClassic", oddsZhaDan)
		odscfg.odds[oddsChunTian] = cfg.GetConfigInt("OddsClassic", oddsChunTian)
		odscfg.odds[oddsFanChunTian] = cfg.GetConfigInt("OddsClassic", oddsFanChunTian)
		odscfg.odds[oddsRocket] = cfg.GetConfigInt("OddsClassic", oddsRocket)

		odscfg.odds[oddsSanPai] = cfg.GetConfigInt("OddsClassic", oddsSanPai)
		odscfg.odds[oddsDuiZi] = cfg.GetConfigInt("OddsClassic", oddsDuiZi)
		odscfg.odds[oddsShunZi] = cfg.GetConfigInt("OddsClassic", oddsShunZi)
		odscfg.odds[oddsTongHua] = cfg.GetConfigInt("OddsClassic", oddsTongHua)
		odscfg.odds[oddsXiaoWang] = cfg.GetConfigInt("OddsClassic", oddsXiaoWang)
		odscfg.odds[oddsDaWang] = cfg.GetConfigInt("OddsClassic", oddsDaWang)
		odscfg.odds[oddsBaoZi] = cfg.GetConfigInt("OddsClassic", oddsBaoZi)
	}

}

func (odscfg *ddzOddsCfg) GetValue(key string) int {
	return odscfg.odds[key]
}

func NewDdzOddsCfg(pokerType int32, cfg *config.Group) *ddzOddsCfg {
	d := new(ddzOddsCfg)
	d.odds = make(map[string]int)
	d.init(pokerType, cfg)
	return d
}
