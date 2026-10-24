package util

import (
	"fmt"
	"time"
)

// Retry ?
func Retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}
		if attempts--; attempts > 0 {
			time.Sleep(sleep * 3)
			fmt.Println("重新連線 剩餘次數:", attempts)
			return Retry(attempts, sleep, fn)
		}
		return err
	}
	return nil
}

type stop struct {
	error
}
