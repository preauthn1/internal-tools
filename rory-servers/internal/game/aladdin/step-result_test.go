package aladdin

import (
	"fmt"
	"testing"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestNewStepResult_OneReel_OneSymbol(t *testing.T) {

	attrs := []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsExpand: proto.Int32(1)},
	}

	result := slot.GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("b1")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p8")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("b1")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
		},
	}

	manager, _ := slot.NewGameResultManager(5, 3)
	manager.GameResult = result

	step := NewStepResult(1, false, 10, 100, attrs, manager, nil)

	fmt.Println("ExpandSymbols", step.expandSymbols)
	fmt.Println("IsRespin", step.IsRespin)
	fmt.Println("IsExpand", step.IsExpand)
	fmt.Println("ExpandReel", step.ExpandReel)
	fmt.Println("ExpandSymbol", step.ExpandSymbol)
	fmt.Println("ExpandOrder", step.ExpandOrder)

	fmt.Println("AllExpandReels", step.AllExpandReels)

	fmt.Println("UnexpandResult", step.UnexpandResult)
	fmt.Println("ExpandResult", step.ExpandResult)
	fmt.Println("GamePayResult", step.GamePayResult)
	fmt.Println("GameBetMoney", step.GameBetMoney)
	fmt.Println("GamePayMoney", step.GamePayMoney)
	fmt.Println("GameWinMoney", step.GameWinMoney)

}

func TestNewStepResult_ExpandedReel_AllSymbol(t *testing.T) {

	attrs := []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsExpand: proto.Int32(1)},
	}

	result := slot.GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("b1")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p8")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
		},
	}

	manager, _ := slot.NewGameResultManager(5, 3)
	manager.GameResult = result

	step := NewStepResult(2, true, 10, 100, attrs, manager, nil)

	fmt.Println("ExpandSymbols", step.expandSymbols)
	fmt.Println("IsRespin", step.IsRespin)
	fmt.Println("IsExpand", step.IsExpand)
	fmt.Println("ExpandReel", step.ExpandReel)
	fmt.Println("ExpandSymbol", step.ExpandSymbol)
	fmt.Println("ExpandOrder", step.ExpandOrder)

	fmt.Println("AllExpandReels", step.AllExpandReels)

	fmt.Println("UnexpandResult", step.UnexpandResult)
	fmt.Println("ExpandResult", step.ExpandResult)
	fmt.Println("GamePayResult", step.GamePayResult)
	fmt.Println("GameBetMoney", step.GameBetMoney)
	fmt.Println("GamePayMoney", step.GamePayMoney)
	fmt.Println("GameWinMoney", step.GameWinMoney)

}

func TestNewStepResult_TwoExpandedReel_Symbol(t *testing.T) {

	attrs := []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsExpand: proto.Int32(1)},
	}

	result := slot.GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p3")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p3")},
		},
	}

	manager, _ := slot.NewGameResultManager(5, 3)
	manager.GameResult = result

	step := NewStepResult(3, true, 10, 100, attrs, manager, nil)

	fmt.Println("ExpandSymbols", step.expandSymbols)
	fmt.Println("IsRespin", step.IsRespin)
	fmt.Println("IsExpand", step.IsExpand)
	fmt.Println("ExpandReel", step.ExpandReel)
	fmt.Println("ExpandSymbol", step.ExpandSymbol)
	fmt.Println("ExpandOrder", step.ExpandOrder)

	fmt.Println("AllExpandReels", step.AllExpandReels)

	fmt.Println("UnexpandResult", step.UnexpandResult)
	fmt.Println("ExpandResult", step.ExpandResult)
	fmt.Println("GamePayResult", step.GamePayResult)
	fmt.Println("GameBetMoney", step.GameBetMoney)
	fmt.Println("GamePayMoney", step.GamePayMoney)
	fmt.Println("GameWinMoney", step.GameWinMoney)

}

func TestNewStepResult_ExpandedReel_NoSymbol(t *testing.T) {

	attrs := []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsExpand: proto.Int32(1)},
	}

	result := slot.GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("b1")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p8")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
		},
	}

	manager, _ := slot.NewGameResultManager(5, 3)
	manager.GameResult = result

	step := NewStepResult(2, true, 10, 100, attrs, manager, nil)

	fmt.Println("ExpandSymbols", step.expandSymbols)
	fmt.Println("IsRespin", step.IsRespin)
	fmt.Println("IsExpand", step.IsExpand)
	fmt.Println("ExpandReel", step.ExpandReel)
	fmt.Println("ExpandSymbol", step.ExpandSymbol)
	fmt.Println("ExpandOrder", step.ExpandOrder)

	fmt.Println("AllExpandReels", step.AllExpandReels)

	fmt.Println("UnexpandResult", step.UnexpandResult)
	fmt.Println("ExpandResult", step.ExpandResult)
	fmt.Println("GamePayResult", step.GamePayResult)
	fmt.Println("GameBetMoney", step.GameBetMoney)
	fmt.Println("GamePayMoney", step.GamePayMoney)
	fmt.Println("GameWinMoney", step.GameWinMoney)

}

func TestNewStepResult_TwoReel(t *testing.T) {

	attrs := []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsExpand: proto.Int32(1)},
	}

	result := slot.GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("b1")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
		},
	}

	manager, _ := slot.NewGameResultManager(5, 3)
	manager.GameResult = result

	step := NewStepResult(1, true, 10, 100, attrs, manager, nil)

	fmt.Println("ExpandSymbols", step.expandSymbols)
	fmt.Println("IsRespin", step.IsRespin)
	fmt.Println("IsExpand", step.IsExpand)
	fmt.Println("ExpandReel", step.ExpandReel)
	fmt.Println("ExpandSymbol", step.ExpandSymbol)
	fmt.Println("ExpandOrder", step.ExpandOrder)

	fmt.Println("AllExpandReels", step.AllExpandReels)

	fmt.Println("UnexpandResult", step.UnexpandResult)
	fmt.Println("ExpandResult", step.ExpandResult)
	fmt.Println("GamePayResult", step.GamePayResult)
	fmt.Println("GameBetMoney", step.GameBetMoney)
	fmt.Println("GamePayMoney", step.GamePayMoney)
	fmt.Println("GameWinMoney", step.GameWinMoney)

}

func TestNewStepResult_FirsExpandedReel_ThreeSymbol(t *testing.T) {

	attrs := []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsExpand: proto.Int32(1)},
	}

	result := slot.GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("b1")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p8")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
		},
	}

	manager, _ := slot.NewGameResultManager(5, 3)
	manager.GameResult = result

	step := NewStepResult(2, true, 10, 100, attrs, manager, nil)

	fmt.Println("ExpandSymbols", step.expandSymbols)
	fmt.Println("IsRespin", step.IsRespin)
	fmt.Println("IsExpand", step.IsExpand)
	fmt.Println("ExpandReel", step.ExpandReel)
	fmt.Println("ExpandSymbol", step.ExpandSymbol)
	fmt.Println("ExpandOrder", step.ExpandOrder)

	fmt.Println("AllExpandReels", step.AllExpandReels)

	fmt.Println("UnexpandResult", step.UnexpandResult)
	fmt.Println("ExpandResult", step.ExpandResult)
	fmt.Println("GamePayResult", step.GamePayResult)
	fmt.Println("GameBetMoney", step.GameBetMoney)
	fmt.Println("GamePayMoney", step.GamePayMoney)
	fmt.Println("GameWinMoney", step.GameWinMoney)

	result2 := slot.GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("b1")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p8")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
		},
	}

	manager2, _ := slot.NewGameResultManager(5, 3)
	manager2.GameResult = result2

	step2 := NewStepResult(2, true, 10, 100, attrs, manager2, step)

	fmt.Println("ExpandSymbols", step2.expandSymbols)
	fmt.Println("IsRespin", step2.IsRespin)
	fmt.Println("IsExpand", step2.IsExpand)
	fmt.Println("ExpandReel", step2.ExpandReel)
	fmt.Println("ExpandSymbol", step2.ExpandSymbol)
	fmt.Println("ExpandOrder", step2.ExpandOrder)

	fmt.Println("AllExpandReels", step2.AllExpandReels)

	fmt.Println("UnexpandResult", step2.UnexpandResult)
	fmt.Println("ExpandResult", step2.ExpandResult)
	fmt.Println("GamePayResult", step2.GamePayResult)
	fmt.Println("GameBetMoney", step2.GameBetMoney)
	fmt.Println("GamePayMoney", step2.GamePayMoney)
	fmt.Println("GameWinMoney", step2.GameWinMoney)

}
