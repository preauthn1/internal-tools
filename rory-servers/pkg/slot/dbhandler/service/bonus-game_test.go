package service

import "testing"

func TestDBBonusGameService_Read(t *testing.T) {

	s := new(DBBonusGameService)

	t.Run("Test User", func(tt *testing.T) {
		data, err := s.Read("", 1, 1, 1)
		tt.Log(err)
		tt.Logf("Length: %d\n", len(data))
		for _, row := range data {
			tt.Log(row)
		}
	})
}

func TestDBBonusGameService_Record(t *testing.T) {

	s := new(DBBonusGameService)

	t.Run("Test User", func(tt *testing.T) {
		tt.Log(s.Record("", 1, 1, 1, 1, 1, 1, 0, 1, 1, ""))
		tt.Log(s.Record("", 1, 1, 1, 1, 1, 2, 0, 1, 1, ""))
		tt.Log(s.Record("", 1, 1, 1, 1, 1, 3, 0, 1, 1, ""))
		tt.Log(s.Record("", 1, 1, 1, 1, 1, 4, 0, 1, 1, ""))
		tt.Log(s.Record("", 1, 1, 1, 1, 1, 5, 0, 1, 1, ""))
	})
}

func TestDBBonusGameService_Reset(t *testing.T) {

	s := new(DBBonusGameService)

	t.Run("Test User", func(tt *testing.T) {
		tt.Log(s.Reset("", 1, 1, 1, 0))
	})
}

func TestAPIBonusGameService_Read(t *testing.T) {

	s := new(APIBonusGameService)

	t.Run("Test User", func(tt *testing.T) {
		data, err := s.Read("", 1, 1, 1)
		tt.Log(err)
		tt.Logf("Length: %d\n", len(data))
		for _, row := range data {
			tt.Log(row)
		}
	})
}

func TestAPIBonusGameService_Record(t *testing.T) {

	s := new(APIBonusGameService)

	t.Run("Test User", func(tt *testing.T) {
		tt.Log(s.Record("", 1, 1, 1, 1, 1, 1, 0, 30, 1, ""))
		tt.Log(s.Record("", 1, 1, 1, 1, 1, 2, 0, 30, 1, ""))
		tt.Log(s.Record("", 1, 1, 1, 1, 1, 3, 0, 30, 1, ""))
		tt.Log(s.Record("", 1, 1, 1, 1, 1, 4, 0, 30, 1, ""))
		tt.Log(s.Record("", 1, 1, 1, 1, 1, 5, 0, 30, 1, ""))
	})
}

func TestAPIBonusGameService_Reset(t *testing.T) {

	s := new(APIBonusGameService)

	t.Run("Test User", func(tt *testing.T) {
		tt.Log(s.Reset("", 1, 1, 1, 1))
	})
}
