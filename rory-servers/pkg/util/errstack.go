package util

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"time"
)

// GetStackInfo ?
func GetStackInfo() string {
	kb := 4

	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) //4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")

	return string(stack)
}

// PanicErrStack ?
func PanicErrStack() error {
	if err := recover(); err != nil {
		errorStr := GetStackInfo()
		s := fmt.Sprintf("%s %v\r\n%s", time.Now().Format("2006/01/02 15:04:05.000000"), err, errorStr)
		WriteErrorLog(s)
		fmt.Println(" ======== " + time.Now().Format("2006/01/02 15:04:05.000000") + " ======== ")
		fmt.Println(err)
		fmt.Println(errorStr)
		return errors.New(errorStr)
	}

	return nil
}

// SafeCall 安全調用
func SafeCall(f func()) (err error) {
	err = nil

	defer func() {
		if err1 := recover(); err1 != nil {
			errorStr := GetStackInfo()
			s := fmt.Sprintf("%s %v\r\n%s", time.Now().Format("2006/01/02 15:04:05.000000"), err1, errorStr)
			WriteErrorLog(s)
			fmt.Println(s)
			err = errors.New(errorStr)
		}
	}()

	f()

	return
}

// AsyncSafeCall 異步安全調用
func AsyncSafeCall(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				errorStr := GetStackInfo()
				s := fmt.Sprintf("%s %v\r\n%s", time.Now().Format("2006/01/02 15:04:05.000000"), err, errorStr)
				WriteErrorLog(s)
				fmt.Println(s)
			}
		}()

		f()
	}()

}
