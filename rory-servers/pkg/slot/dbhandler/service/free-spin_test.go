package service

import "testing"

func TestDBFreeSpinService_Read(t *testing.T) {
	s := new(DBFreeSpinService)
	t.Run("Test User", func(tt *testing.T) {
		data, err := s.Read("", 1, 1, 1)
		tt.Log(data, err)
	})
}

func TestDBFreeSpinService_Record(t *testing.T) {
	s := new(DBFreeSpinService)
	t.Run("Test User", func(tt *testing.T) {
		err := s.Record("", 1, 1, 1, 1, 1, 1, 1, 1, 10, 0, 0, "")
		tt.Log(err)
	})
}

func TestAPIFreeSpinService_Read(t *testing.T) {
	s := new(APIFreeSpinService)
	t.Run("Test User", func(tt *testing.T) {
		data, err := s.Read("", 1, 1, 1)
		tt.Log(data, err)
	})
}

func TestAPIFreeSpinService_Record(t *testing.T) {
	s := new(APIFreeSpinService)
	t.Run("Test User", func(tt *testing.T) {
		err := s.Record("", 1, 1, 1, 1, 1, 1, 1, 1, 10, 100, 0, "")
		tt.Log(err)
	})
}
