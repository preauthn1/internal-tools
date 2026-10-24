package mahjong

import (
	"testing"
)

//筒子
func TestNewTong(t *testing.T) {
	cg := NewTong(1)
	if cg.String() != "1筒,2筒,3筒,4筒,5筒,6筒,7筒,8筒,9筒,1筒,2筒,3筒,4筒,5筒,6筒,7筒,8筒,9筒,1筒,2筒,3筒,4筒,5筒,6筒,7筒,8筒,9筒,1筒,2筒,3筒,4筒,5筒,6筒,7筒,8筒,9筒," {
		t.Errorf("筒子錯誤內容 : %+v", cg.String())
	}
}

//条
func TestNewTiao(t *testing.T) {
	cg := NewTiao(1)
	if cg.String() != "1条,2条,3条,4条,5条,6条,7条,8条,9条,1条,2条,3条,4条,5条,6条,7条,8条,9条,1条,2条,3条,4条,5条,6条,7条,8条,9条,1条,2条,3条,4条,5条,6条,7条,8条,9条," {
		t.Errorf("条錯誤內容 : %+v", cg.String())
	}
}

//万
func TestNewWan(t *testing.T) {
	cg := NewWan(1)
	if cg.String() != "1万,2万,3万,4万,5万,6万,7万,8万,9万,1万,2万,3万,4万,5万,6万,7万,8万,9万,1万,2万,3万,4万,5万,6万,7万,8万,9万,1万,2万,3万,4万,5万,6万,7万,8万,9万," {
		t.Errorf("万錯誤內容 : %+v", cg.String())
	}
}

//风牌
func TestNewFeng(t *testing.T) {
	cg := NewFeng(1)
	if cg.String() != "东,西,南,北,东,西,南,北,东,西,南,北,东,西,南,北," {
		t.Errorf("风牌錯誤內容 : %+v", cg.String())
	}
}

//箭牌
func TestNewJian(t *testing.T) {
	cg := NewJian(1)
	if cg.String() != "红中,发财,白板,红中,发财,白板,红中,发财,白板,红中,发财,白板," {
		t.Errorf("箭牌錯誤內容 : %+v", cg.String())
	}
}

//花牌
func TestNewHua(t *testing.T) {
	cg := NewHua(1)
	if cg.String() != "春,夏,秋,冬,梅,兰,竹,菊,春,夏,秋,冬,梅,兰,竹,菊,春,夏,秋,冬,梅,兰,竹,菊,春,夏,秋,冬,梅,兰,竹,菊," {
		t.Errorf("花牌錯誤內容 : %+v", cg.String())
	}
}

//百搭
func TestNewBaida(t *testing.T) {
	cg := NewBaida(1)
	if cg.String() != "春,夏,秋,冬,春,夏,秋,冬,春,夏,秋,冬,春,夏,秋,冬," {
		t.Errorf("百搭錯誤內容 : %+v", cg.String())
	}
}

//是否可以碰牌
func TestHavePairOf(t *testing.T) {
	handCards := new(CardGroup)
	handCards = NewTong(1)
	result := handCards.HavePairOf(111)
	if !result {
		t.Errorf("錯誤:是否可以碰牌[%v] , cards[%+v]", result, handCards.String())
	}
}

//是否可以引杠
func TestHaveThreeOf(t *testing.T) {
	handCards := new(CardGroup)
	handCards = NewTong(1)
	result := handCards.HaveThreeOf(111)
	if !result {
		t.Errorf("錯誤:是否可以引杠[%v] , cards[%+v]", result, handCards.String())
	}
}

//查询暗杠牌,暗杠
func TestHaveFour(t *testing.T) {
	handCards := new(CardGroup)
	handCards = NewTong(1)
	result := handCards.HaveFour()
	if len(result) != 9 {
		t.Errorf("錯誤:查询暗杠牌,暗杠[%v] , cards[%+v]", result, handCards.String())
	}
}

//是否可以杠牌,暗杠
func TestHaveFourOf(t *testing.T) {
	handCards := new(CardGroup)
	handCards = NewTong(1)
	result := handCards.HaveFourOf(111)
	if !result {
		t.Errorf("錯誤:是否可以杠牌,暗杠[%v] , cards[%+v]", result, handCards.String())
	}
}

//检测牌中是否有某种花色的牌
func TestHaveColor(t *testing.T) {
	handCards := new(CardGroup)
	handCards = NewTong(1)
	result := handCards.HaveColor(1)
	if !result {
		t.Errorf("錯誤:检测牌中是否有某种花色的牌[%v] , cards[%+v]", result, handCards.String())
	}
}
