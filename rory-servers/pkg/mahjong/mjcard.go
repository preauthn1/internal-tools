package mahjong

import (
	"fmt"
)

// Card int32=100*Color+10*Number+index(第幾副牌)
const (
	ColorTong  int32 = 1 // 筒子
	ColorTiao  int32 = 2 // 條子
	ColorWan   int32 = 3 // 萬子
	ColorFeng  int32 = 4 // 風牌
	ColorJain  int32 = 5 // 箭牌
	ColorHua   int32 = 6 // 花牌
	ColorBaida int32 = 7 // 百搭牌
	ColorQita  int32 = 8 // 其他
)

// ?
const (
	// 風牌
	NumberEast  int32 = 1 // 東
	NumberWest  int32 = 2 // 西
	NumberSouth int32 = 3 // 南
	NumberNorth int32 = 4 // 北
	// 箭牌
	NumberZhong int32 = 1 //中
	NumberFa    int32 = 2 //發
	NumberBai   int32 = 3 //白
	// 花牌
	NumberChun int32 = 1 // 春
	NumberXia  int32 = 2 // 夏
	NumberQiu  int32 = 3 // 秋
	NumberDong int32 = 4 // 冬
	NumberMei  int32 = 5 // 梅
	NumberLan  int32 = 6 // 蘭
	NumberZhu  int32 = 7 // 竹
	NumberJu   int32 = 8 // 菊
	// 百搭
	NumberCaishen   int32 = 1 // 財神
	NumberMao       int32 = 2 // 貓
	NumberLaoshu    int32 = 3 // 老鼠
	NumberJubaopeng int32 = 4 // 聚寶盆
	// 序牌1~9
)

// CardString ?
func CardString(cd int32) string {
	color := ""

	if cd/100 == ColorTong {
		color = "筒"
		return fmt.Sprintf("%d%s", cd%100/10, color)
	} else if cd/100 == ColorTiao {
		color = "条"
		return fmt.Sprintf("%d%s", cd%100/10, color)
	} else if cd/100 == ColorWan {
		color = "万"
		return fmt.Sprintf("%d%s", cd%100/10, color)
	} else if cd/100 == ColorFeng {
		if cd%100/10 == NumberEast {
			return "东"
		} else if cd%100/10 == NumberWest {
			return "西"
		} else if cd%100/10 == NumberSouth {
			return "南"
		} else if cd%100/10 == NumberNorth {
			return "北"
		}
	} else if cd/100 == ColorJain {
		if cd%100/10 == NumberZhong {
			return "红中"
		} else if cd%100/10 == NumberFa {
			return "发财"
		} else if cd%100/10 == NumberBai {
			return "白板"
		}
	} else if cd/100 == ColorHua {
		if cd%100/10 == NumberChun {
			return "春"
		} else if cd%100/10 == NumberXia {
			return "夏"
		} else if cd%100/10 == NumberQiu {
			return "秋"
		} else if cd%100/10 == NumberDong {
			return "冬"
		} else if cd%100/10 == NumberMei {
			return "梅"
		} else if cd%100/10 == NumberLan {
			return "兰"
		} else if cd%100/10 == NumberZhu {
			return "竹"
		} else if cd%100/10 == NumberJu {
			return "菊"
		}
	} else if cd/100 == ColorBaida {
		if cd%100/10 == NumberCaishen {
			return "财神"
		} else if cd%100/10 == NumberMao {
			return "猫"
		} else if cd%100/10 == NumberLaoshu {
			return "老鼠"
		} else if cd%100/10 == NumberJubaopeng {
			return "聚宝盆"
		}
	}
	return "未知牌"

}

// ColorString ?
func ColorString(color int32) string {
	if color == 1 {
		return "筒子"
	} else if color == 2 {
		return "条子"
	} else if color == 3 {
		return "万子"
	} else if color == 4 {
		return "风牌"
	} else if color == 5 {
		return "箭牌"
	} else if color == 6 {
		return "花牌"
	} else if color == 7 {
		return "百搭牌"
	} else {
		return "未知颜色"
	}
}

// GetCardIndex 根據牌值 獲取牌索引
func GetCardIndex(cd int32) int32 {
	return cd % 100 % 10
}

// GetCardNumber 根據牌值獲取牌點數
func GetCardNumber(cd int32) int32 {
	return cd / 10 % 10
}

// GetCardColor 根據牌值獲取牌花色
func GetCardColor(cd int32) int32 {
	return cd / 100
}

// GetCardValue 根據牌值獲取牌有效值，點數和花色
func GetCardValue(cd int32) int32 {
	return cd / 10
}

// GetNumberByValue 根據牌有效值 獲取點數
func GetNumberByValue(value int32) int32 {
	return value % 10
}

// GetColorByValue 根據牌有效值 獲取花色
func GetColorByValue(value int32) int32 {
	return value / 10
}
