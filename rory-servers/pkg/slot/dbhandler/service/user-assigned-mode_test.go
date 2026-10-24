package service

import "testing"

func TestDBUserAssignedModeService_Read(t *testing.T) {
	s := new(DBUserAssignedModeService)
	t.Run("Test User", func(tt *testing.T) {
		data, err := s.Read("", 1, 1, 1)
		tt.Log(data, err)
	})
}

func TestDBUserAssignedModeService_Recored(t *testing.T) {
	s := new(DBUserAssignedModeService)
	t.Run("Test User", func(tt *testing.T) {
		err := s.Record("", 1, 1, 1, 1, "2020-03-03 13:23:45.876")
		tt.Log(err)
	})
}

func TestDBUserAssignedModeService_Delete(t *testing.T) {
	s := new(DBUserAssignedModeService)
	t.Run("Test User", func(tt *testing.T) {
		err := s.Delete("", 1, 1, 1)
		tt.Log(err)
	})
}

func TestAPIUserAssignedModeService_Read(t *testing.T) {
	s := new(APIUserAssignedModeService)
	t.Run("Test User", func(tt *testing.T) {
		data, err := s.Read("", 1, 1, 1)
		tt.Log(data, err)
	})
}

func TestAPIUserAssignedModeService_Recored(t *testing.T) {
	s := new(APIUserAssignedModeService)
	t.Run("Test User", func(tt *testing.T) {
		err := s.Record("", 1, 1, 1, 1, "2020-03-03 13:23:45.876")
		tt.Log(err)
	})
}

func TestAPIUserAssignedModeService_Delete(t *testing.T) {
	s := new(APIUserAssignedModeService)
	t.Run("Test User", func(tt *testing.T) {
		err := s.Delete("", 1, 1, 1)
		tt.Log(err)
	})
}
