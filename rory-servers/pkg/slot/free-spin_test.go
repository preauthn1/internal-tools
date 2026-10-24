package slot

import "testing"

func TestFreeSpin_GetMultipleValue(t *testing.T) {

	freeSpin := new(FreeSpin)

	t.Run("No Multiple", func(tt *testing.T) {

		actual := freeSpin.GetMultipleValue()

		if actual != 1 {
			tt.Errorf("MultipleValue want %d, but %d", 1, actual)
		}

	})

	freeSpin.WinMultiple = []int32{1, 2, 2, 3, 3, 4}

	t.Run("1'th Win", func(tt *testing.T) {

		freeSpin.WinCount = 0

		actual := freeSpin.GetMultipleValue()

		if actual != 1 {
			tt.Errorf("MultipleValue want %d, but %d", 1, actual)
		}

	})

	t.Run("2'th Win", func(tt *testing.T) {

		freeSpin.WinCount = 1

		actual := freeSpin.GetMultipleValue()

		if actual != 2 {
			tt.Errorf("MultipleValue want %d, but %d", 2, actual)
		}

	})

	t.Run("3'th Win", func(tt *testing.T) {

		freeSpin.WinCount = 2

		actual := freeSpin.GetMultipleValue()

		if actual != 2 {
			tt.Errorf("MultipleValue want %d, but %d", 2, actual)
		}

	})

	t.Run("4'th Win", func(tt *testing.T) {

		freeSpin.WinCount = 3

		actual := freeSpin.GetMultipleValue()

		if actual != 3 {
			tt.Errorf("MultipleValue want %d, but %d", 3, actual)
		}

	})

	t.Run("5'th Win", func(tt *testing.T) {

		freeSpin.WinCount = 4

		actual := freeSpin.GetMultipleValue()

		if actual != 3 {
			tt.Errorf("MultipleValue want %d, but %d", 3, actual)
		}

	})

	t.Run("6'th Win", func(tt *testing.T) {

		freeSpin.WinCount = 5

		actual := freeSpin.GetMultipleValue()

		if actual != 4 {
			tt.Errorf("MultipleValue want %d, but %d", 4, actual)
		}

	})

	t.Run("7'th Win", func(tt *testing.T) {

		freeSpin.WinCount = 6

		actual := freeSpin.GetMultipleValue()

		if actual != 4 {
			tt.Errorf("MultipleValue want %d, but %d", 4, actual)
		}

	})

	t.Run("8'th Win", func(tt *testing.T) {

		freeSpin.WinCount = 7

		actual := freeSpin.GetMultipleValue()

		if actual != 4 {
			tt.Errorf("MultipleValue want %d, but %d", 4, actual)
		}

	})

	t.Run("9'th Win", func(tt *testing.T) {

		freeSpin.WinCount = 8

		actual := freeSpin.GetMultipleValue()

		if actual != 4 {
			tt.Errorf("MultipleValue want %d, but %d", 4, actual)
		}

	})

	t.Run("10'th Win", func(tt *testing.T) {

		freeSpin.WinCount = 9

		actual := freeSpin.GetMultipleValue()

		if actual != 4 {
			tt.Errorf("MultipleValue want %d, but %d", 4, actual)
		}

	})

	t.Run("11'th Win", func(tt *testing.T) {

		freeSpin.WinCount = 10

		actual := freeSpin.GetMultipleValue()

		if actual != 4 {
			tt.Errorf("MultipleValue want %d, but %d", 4, actual)
		}

	})

	t.Run("12'th Win", func(tt *testing.T) {

		freeSpin.WinCount = 11

		actual := freeSpin.GetMultipleValue()

		if actual != 4 {
			tt.Errorf("MultipleValue want %d, but %d", 4, actual)
		}

	})

	t.Run("13'th Win", func(tt *testing.T) {

		freeSpin.WinCount = 12

		actual := freeSpin.GetMultipleValue()

		if actual != 4 {
			tt.Errorf("MultipleValue want %d, but %d", 4, actual)
		}

	})

	t.Run("14'th Win", func(tt *testing.T) {

		freeSpin.WinCount = 13

		actual := freeSpin.GetMultipleValue()

		if actual != 4 {
			tt.Errorf("MultipleValue want %d, but %d", 4, actual)
		}

	})

	t.Run("15'th Win", func(tt *testing.T) {

		freeSpin.WinCount = 14

		actual := freeSpin.GetMultipleValue()

		if actual != 4 {
			tt.Errorf("MultipleValue want %d, but %d", 4, actual)
		}

	})

}
