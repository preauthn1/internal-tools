package service

import "testing"

func TestDBDeferredAwardService_Read(t *testing.T) {
	s := new(DBDeferredAwardService)
	t.Run("Test User", func(tt *testing.T) {
		data, err := s.Read("", 1, 1, 1, "Test")
		tt.Log(err)
		tt.Logf("Length: %d\n", len(data))
		for _, row := range data {
			tt.Log(row)
		}
	})
}

func TestDBDeferredAwardService_Record(t *testing.T) {
	s := new(DBDeferredAwardService)
	t.Run("Test User", func(tt *testing.T) {
		tt.Log(s.Record("", 1, 1, 1, "Test", 1, 1, "qq"))
	})
}

func TestDBDeferredAwardService_Pop(t *testing.T) {
	s := new(DBDeferredAwardService)
	t.Run("Test User", func(tt *testing.T) {
		tt.Log(s.Pop("", 1, 1, 1, "Test"))
	})
}

func TestAPIDeferredAwardService_Read(t *testing.T) {
	s := new(APIDeferredAwardService)
	t.Run("Test User", func(tt *testing.T) {
		data, err := s.Read("", 1, 1, 1, "Test")
		tt.Log(err)
		tt.Logf("Length: %d\n", len(data))
		for _, row := range data {
			tt.Log(row)
		}
	})
}

func TestAPIDeferredAwardService_Record(t *testing.T) {
	s := new(APIDeferredAwardService)
	t.Run("Test User", func(tt *testing.T) {
		tt.Log(s.Record("", 1, 1, 1, "Test", 1, 1, "qq"))
	})
}

func TestAPIDeferredAwardService_Pop(t *testing.T) {
	s := new(APIDeferredAwardService)
	t.Run("Test User", func(tt *testing.T) {
		tt.Log(s.Pop("", 1, 1, 1, "Test"))
	})
}
