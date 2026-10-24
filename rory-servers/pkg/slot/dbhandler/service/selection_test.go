package service

import "testing"

func TestDBSelectionService_Read(t *testing.T) {
	s := new(DBSelectionService)
	t.Run("Test User", func(tt *testing.T) {
		data, err := s.Read("", 1, 1, 1)
		tt.Log(data, err)
	})
}

func TestDBSelectionService_Record(t *testing.T) {
	s := new(DBSelectionService)
	t.Run("Test User", func(tt *testing.T) {
		err := s.Record("", 1, 1, 1, 1, 1, 5, "Test")
		tt.Log(err)
	})
}

func TestDBSelectionService_Reset(t *testing.T) {
	s := new(DBSelectionService)
	t.Run("Test User", func(tt *testing.T) {
		err := s.Reset("", 1, 1, 1)
		tt.Log(err)
	})
}

func TestAPISelectionService_Read(t *testing.T) {
	s := new(APISelectionService)
	t.Run("Test User", func(tt *testing.T) {
		data, err := s.Read("", 1, 1, 1)
		tt.Log(data, err)
	})
}

func TestAPISelectionService_Record(t *testing.T) {
	s := new(APISelectionService)
	t.Run("Test User", func(tt *testing.T) {
		err := s.Record("", 1, 1, 1, 1, 1, 5, "Test")
		tt.Log(err)
	})
}

func TestAPISelectionService_Reset(t *testing.T) {
	s := new(APISelectionService)
	t.Run("Test User", func(tt *testing.T) {
		err := s.Reset("", 1, 1, 1)
		tt.Log(err)
	})
}
