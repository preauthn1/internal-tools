package slot

import (
	"reflect"
	"testing"
)

func TestPayResult_AddBlinkPosition(t *testing.T) {

	t.Run("11111", func(tt *testing.T) {
		r := new(PayResult)
		r.AddBlinkPosition(Position{Reel: 1, Row: 1})
		r.AddBlinkPosition(Position{Reel: 2, Row: 1})
		r.AddBlinkPosition(Position{Reel: 3, Row: 1})
		r.AddBlinkPosition(Position{Reel: 4, Row: 1})
		r.AddBlinkPosition(Position{Reel: 5, Row: 1})

		expected := []Position{
			{Reel: 1, Row: 1},
			{Reel: 2, Row: 1},
			{Reel: 3, Row: 1},
			{Reel: 4, Row: 1},
			{Reel: 5, Row: 1},
		}
		if !reflect.DeepEqual(r.BlinkPositions, expected) {
			tt.Errorf("BlinkPositions want %v but %v", expected, r.BlinkPositions)
		}

		dExpected := [][]Position{
			{
				{Reel: 1, Row: 1},
				{Reel: 2, Row: 1},
				{Reel: 3, Row: 1},
				{Reel: 4, Row: 1},
				{Reel: 5, Row: 1},
			},
		}
		if !reflect.DeepEqual(r.BlinkDetails, dExpected) {
			tt.Errorf("BlinkDetails want %v but %v", dExpected, r.BlinkDetails)
		}
	})

	t.Run("11121", func(tt *testing.T) {
		r := new(PayResult)
		r.AddBlinkPosition(Position{Reel: 1, Row: 1})
		r.AddBlinkPosition(Position{Reel: 2, Row: 1})
		r.AddBlinkPosition(Position{Reel: 3, Row: 1})
		r.AddBlinkPosition(Position{Reel: 4, Row: 1})
		r.AddBlinkPosition(Position{Reel: 4, Row: 2})
		r.AddBlinkPosition(Position{Reel: 5, Row: 1})

		expected := []Position{
			{Reel: 1, Row: 1},
			{Reel: 2, Row: 1},
			{Reel: 3, Row: 1},
			{Reel: 4, Row: 1},
			{Reel: 4, Row: 2},
			{Reel: 5, Row: 1},
		}
		if !reflect.DeepEqual(r.BlinkPositions, expected) {
			tt.Errorf("BlinkPositions want %v but %v", expected, r.BlinkPositions)
		}

		dExpected := [][]Position{
			{
				{Reel: 1, Row: 1},
				{Reel: 2, Row: 1},
				{Reel: 3, Row: 1},
				{Reel: 4, Row: 1},
				{Reel: 5, Row: 1},
			},
			{
				{Reel: 1, Row: 1},
				{Reel: 2, Row: 1},
				{Reel: 3, Row: 1},
				{Reel: 4, Row: 2},
				{Reel: 5, Row: 1},
			},
		}
		if !reflect.DeepEqual(r.BlinkDetails, dExpected) {
			tt.Errorf("BlinkDetails want %v but %v", dExpected, r.BlinkDetails)
		}
	})

	t.Run("11131", func(tt *testing.T) {
		r := new(PayResult)
		r.AddBlinkPosition(Position{Reel: 1, Row: 1})
		r.AddBlinkPosition(Position{Reel: 2, Row: 1})
		r.AddBlinkPosition(Position{Reel: 3, Row: 1})
		r.AddBlinkPosition(Position{Reel: 4, Row: 1})
		r.AddBlinkPosition(Position{Reel: 4, Row: 2})
		r.AddBlinkPosition(Position{Reel: 4, Row: 3})
		r.AddBlinkPosition(Position{Reel: 5, Row: 1})

		expected := []Position{
			{Reel: 1, Row: 1},
			{Reel: 2, Row: 1},
			{Reel: 3, Row: 1},
			{Reel: 4, Row: 1},
			{Reel: 4, Row: 2},
			{Reel: 4, Row: 3},
			{Reel: 5, Row: 1},
		}
		if !reflect.DeepEqual(r.BlinkPositions, expected) {
			tt.Errorf("BlinkPositions want %v but %v", expected, r.BlinkPositions)
		}

		dExpected := [][]Position{
			{
				{Reel: 1, Row: 1},
				{Reel: 2, Row: 1},
				{Reel: 3, Row: 1},
				{Reel: 4, Row: 1},
				{Reel: 5, Row: 1},
			},
			{
				{Reel: 1, Row: 1},
				{Reel: 2, Row: 1},
				{Reel: 3, Row: 1},
				{Reel: 4, Row: 2},
				{Reel: 5, Row: 1},
			},
			{
				{Reel: 1, Row: 1},
				{Reel: 2, Row: 1},
				{Reel: 3, Row: 1},
				{Reel: 4, Row: 3},
				{Reel: 5, Row: 1},
			},
		}
		if !reflect.DeepEqual(r.BlinkDetails, dExpected) {
			tt.Errorf("BlinkDetails want %v but %v", dExpected, r.BlinkDetails)
		}
	})

	t.Run("21121", func(tt *testing.T) {
		r := new(PayResult)
		r.AddBlinkPosition(Position{Reel: 1, Row: 1})
		r.AddBlinkPosition(Position{Reel: 1, Row: 2})
		r.AddBlinkPosition(Position{Reel: 2, Row: 1})
		r.AddBlinkPosition(Position{Reel: 3, Row: 1})
		r.AddBlinkPosition(Position{Reel: 4, Row: 1})
		r.AddBlinkPosition(Position{Reel: 4, Row: 2})
		r.AddBlinkPosition(Position{Reel: 5, Row: 1})

		expected := []Position{
			{Reel: 1, Row: 1},
			{Reel: 1, Row: 2},
			{Reel: 2, Row: 1},
			{Reel: 3, Row: 1},
			{Reel: 4, Row: 1},
			{Reel: 4, Row: 2},
			{Reel: 5, Row: 1},
		}
		if !reflect.DeepEqual(r.BlinkPositions, expected) {
			tt.Errorf("BlinkPositions want %v but %v", expected, r.BlinkPositions)
		}

		dExpected := [][]Position{
			{
				{Reel: 1, Row: 1},
				{Reel: 2, Row: 1},
				{Reel: 3, Row: 1},
				{Reel: 4, Row: 1},
				{Reel: 5, Row: 1},
			},
			{
				{Reel: 1, Row: 2},
				{Reel: 2, Row: 1},
				{Reel: 3, Row: 1},
				{Reel: 4, Row: 1},
				{Reel: 5, Row: 1},
			},
			{
				{Reel: 1, Row: 1},
				{Reel: 2, Row: 1},
				{Reel: 3, Row: 1},
				{Reel: 4, Row: 2},
				{Reel: 5, Row: 1},
			},
			{
				{Reel: 1, Row: 2},
				{Reel: 2, Row: 1},
				{Reel: 3, Row: 1},
				{Reel: 4, Row: 2},
				{Reel: 5, Row: 1},
			},
		}
		if !reflect.DeepEqual(r.BlinkDetails, dExpected) {
			tt.Errorf("BlinkDetails want %v but %v", dExpected, r.BlinkDetails)
		}
	})
}
