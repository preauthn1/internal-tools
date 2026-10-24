package sss

import (
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

var animotionTimeOut int32 = 4

/*
//游戏阶段定义
enum SSSGameStep{
	SSSGameStepNull				= 1;
	SSSGameStepBegin 			= 2;//游戏开始阶段
	SSSGameStepCardTypeChoose 	= 3;//玩家牌型选择阶段(普通牌型、特熟牌型)
	SSSGameStepCombine			= 4;//组牌阶段
	SSSGameStepAnimotion        = 5;//比牌开始动画
	SSSGameStepCompareH		   	= 6;//前墩比牌阶段
	SSSGameStepCompareB		   	= 7;//中墩比牌阶段
	SSSGameStepCompareF		   	= 8;//后墩比牌阶段
	SSSGameStepCompareSpecial  	= 9;//特殊牌型比较阶段
	SSSGameStepUserGun			= 10;//玩家打枪阶段
	SSSGameStepSettlement	   	= 11;//结算阶段
	SSSGameStepEnd			   	= 12;//游戏结束阶段
};

*/

//阶段超时时间配置
var stepTimeOut = []int32{0, 0, 1, 10, 30, 0, 0, 0, 0, 0, 0, 0, 0}

var cardNumberWeight = []int32{1, 0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

type gunType int32

const (
	NormalGun gunType = iota
	ALLGun
)

type NSOddsIdx int32 //普通牌型中的特殊牌型配置

const (
	ZhongDunTongHuaShun NSOddsIdx = 1 //中墩同花顺
	DiDunTongHuaShun    NSOddsIdx = 2 //底墩同花顺
	ZhongDunZhaDan      NSOddsIdx = 3 //中墩炸弹
	DiDunZhaDan         NSOddsIdx = 4 //底墩炸弹
	ChongSan            NSOddsIdx = 5 //冲三
	ZhongDunHuLu        NSOddsIdx = 6 //中墩葫芦
)

var NSOddsCfg = []int32{0, 10, 5, 8, 4, 3, 2}

func (s *SSSTable) GetNsOdds(ni NSOddsIdx) int32 {
	return NSOddsCfg[int32(ni)]
}

var oddsCfg = []int32{2, 3}

func (s *SSSTable) GetGunCfg(gt gunType) int32 {
	return oddsCfg[int32(gt)]
}

func getCardNumber(cardNumber int32) int32 {
	return cardNumberWeight[cardNumber%int32(len(cardNumberWeight))]
}

var weight = map[int32]int32{
	int32(netproto.SSSDunCardsNomalType_TongHuaShun):         1,
	int32(netproto.SSSDunCardsNomalType_ZhaDan):              2,
	int32(netproto.SSSDunCardsNomalType_HuLu):                3,
	int32(netproto.SSSDunCardsNomalType_TongHua):             4,
	int32(netproto.SSSDunCardsNomalType_ShunZi):              5,
	int32(netproto.SSSDunCardsNomalType_SanTiao):             6,
	int32(netproto.SSSDunCardsNomalType_ChongSan):            6,
	int32(netproto.SSSDunCardsNomalType_LiangDui):            8,
	int32(netproto.SSSDunCardsNomalType_YiDui):               9,
	int32(netproto.SSSDunCardsNomalType_SanPai):              10,
	int32(netproto.SSSDunCardsSpecialType_TongHuaShiSanShui): 101,
	int32(netproto.SSSDunCardsSpecialType_ShiSanShui):        102,
	int32(netproto.SSSDunCardsSpecialType_ShiErHuangZu):      103,
	int32(netproto.SSSDunCardsSpecialType_SanTongHuaShun):    104,
	int32(netproto.SSSDunCardsSpecialType_SanFenTianXia):     105,
	int32(netproto.SSSDunCardsSpecialType_QuanDa):            106,
	int32(netproto.SSSDunCardsSpecialType_QuanXiao):          107,
	int32(netproto.SSSDunCardsSpecialType_CouYiSe):           108,

	int32(netproto.SSSDunCardsSpecialType_ShuangGuaiChongSan): 110,
	int32(netproto.SSSDunCardsSpecialType_SanChongJiaDan):     111,
	int32(netproto.SSSDunCardsSpecialType_SiTaoSantiao):       112,
	int32(netproto.SSSDunCardsSpecialType_WuDuiSanTiao):       113,
	int32(netproto.SSSDunCardsSpecialType_LiuDuiBan):          114,
	int32(netproto.SSSDunCardsSpecialType_SanTongHua):         115,
	int32(netproto.SSSDunCardsSpecialType_SanShunZi):          116}

func (s *SSSTable) getStepTimeOut() int32 {
	if s.step == netproto.SSSGameStep_SSSGameStepCombine {
		return s.zupaitime
	}
	return stepTimeOut[s.step%(netproto.SSSGameStep_SSSGameStepEnd+1)]
}

func (s *SSSTable) getWeight(t int32) int32 {
	if v, ok := weight[t]; ok {
		return v
	}
	return 0
}

//特殊牌型水数配置
var specialShui = map[int32]int32{
	int32(netproto.SSSDunCardsSpecialType_TongHuaShiSanShui):  104,
	int32(netproto.SSSDunCardsSpecialType_ShiSanShui):         52,
	int32(netproto.SSSDunCardsSpecialType_ShiErHuangZu):       48,
	int32(netproto.SSSDunCardsSpecialType_SanTongHuaShun):     44,
	int32(netproto.SSSDunCardsSpecialType_SanFenTianXia):      40,
	int32(netproto.SSSDunCardsSpecialType_QuanDa):             36,
	int32(netproto.SSSDunCardsSpecialType_QuanXiao):           32,
	int32(netproto.SSSDunCardsSpecialType_CouYiSe):            28,
	int32(netproto.SSSDunCardsSpecialType_ShuangGuaiChongSan): 24,
	int32(netproto.SSSDunCardsSpecialType_SanChongJiaDan):     22,
	int32(netproto.SSSDunCardsSpecialType_SiTaoSantiao):       20,
	int32(netproto.SSSDunCardsSpecialType_WuDuiSanTiao):       18,
	int32(netproto.SSSDunCardsSpecialType_LiuDuiBan):          16,
	int32(netproto.SSSDunCardsSpecialType_SanTongHua):         14,
	int32(netproto.SSSDunCardsSpecialType_SanShunZi):          12}

//获取特殊牌型水数
func (s *SSSTable) getSpecialShui(t int32) int32 {
	if v, ok := specialShui[t]; ok {
		return v
	}
	return 0
}
