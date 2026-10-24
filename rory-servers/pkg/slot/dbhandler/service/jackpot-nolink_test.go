package service

import (
	"testing"
)

func TestDBJackpotNoLinkService_GetJackpotNoLink(t *testing.T) {
	s := new(DBJackpotNoLinkService)
	t.Run("No Game", func(tt *testing.T) {
		result, err := s.GetJackpotNoLink("", 0, 0, 0, 0)
		if len(result) > 0 {
			tt.Errorf("GetJackpotNoLink results length want %d, but %d", 0, len(result))
		}
		if err == nil {
			tt.Errorf("GetJackpotNoLink want a error")
		}
	})
	t.Run("DFDC Game", func(tt *testing.T) {
		result, err := s.GetJackpotNoLink("", 25, 1, 0, 1)
		tt.Log(err)
		for _, row := range result {
			tt.Log(row)
		}
	})
}

func TestDBJackpotNoLinkService_UpdateJackpotNoLink(t *testing.T) {
	s := new(DBJackpotNoLinkService)
	t.Run("No Game", func(tt *testing.T) {
		result, err := s.UpdateJackpotNoLink("", 0, 0, 0, 0, 0, 0, 0)
		if len(result) > 0 {
			tt.Errorf("GetJackpotNoLink results length want %d, but %d", 0, len(result))
		}
		if err == nil {
			tt.Errorf("GetJackpotNoLink want a error")
		}
	})
	t.Run("DFDC Game", func(tt *testing.T) {
		result, err := s.UpdateJackpotNoLink("", 25, 1, 0, 1, 100, 200, 300)
		tt.Log(err)
		for _, row := range result {
			tt.Log(row)
		}
	})
}

func TestAPIJackpotNoLinkService_GetJackpotNoLink(t *testing.T) {
	s := new(APIJackpotNoLinkService)
	t.Run("No Game", func(tt *testing.T) {
		result, err := s.GetJackpotNoLink("", 0, 0, 0, 0)
		if len(result) > 0 {
			tt.Errorf("GetJackpotNoLink results length want %d, but %d", 0, len(result))
		}
		if err == nil {
			tt.Errorf("GetJackpotNoLink want a error")
		}
	})
	t.Run("DFDC Game", func(tt *testing.T) {
		result, err := s.GetJackpotNoLink("", 25, 1, 0, 1)
		tt.Log(err)
		for _, row := range result {
			tt.Log(row)
		}
	})
}

func TestAPIJackpotNoLinkService_UpdateJackpotNoLink(t *testing.T) {
	s := new(APIJackpotNoLinkService)
	t.Run("No Game", func(tt *testing.T) {
		result, err := s.UpdateJackpotNoLink("", 0, 0, 0, 0, 0, 0, 0)
		if len(result) > 0 {
			tt.Errorf("GetJackpotNoLink results length want %d, but %d", 0, len(result))
		}
		if err == nil {
			tt.Errorf("GetJackpotNoLink want a error")
		}
	})
	t.Run("DFDC Game", func(tt *testing.T) {
		result, err := s.UpdateJackpotNoLink("", 25, 1, 0, 1, 100, 200, 300)
		tt.Log(err)
		for _, row := range result {
			tt.Log(row)
		}
	})
}
