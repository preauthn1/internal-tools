package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

// MustDirExists 檢查目錄
func MustDirExists(dirname string) {
	_, err := os.Stat(dirname)

	if err == nil {
		return
	}

	if os.IsNotExist(err) {
		os.MkdirAll(dirname, 0777)
	}
}

// GetCurrentPath 獲取當前程序執行目錄
func GetCurrentPath() (string, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if currentPath != "/" && currentPath != "\\" {
		if runtime.GOOS == "windows" {
			currentPath = currentPath + "\\"
		} else {
			currentPath = currentPath + "/"
		}
	}
	return currentPath, nil
}

// ReadFile 讀取文件
func ReadFile(filename string) string {
	f, err := os.Open(filename)
	defer f.Close()

	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	return string(data)

}

// FilePutContents ?
func FilePutContents(filename string, content string) error {
	bytes := []byte(content)
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)

	if err != nil {
		return err
	}

	defer fp.Close()

	_, err = fp.Write(bytes)

	return err
}

// FilePutContentsNA ?
func FilePutContentsNA(filename string, content string) error {
	bytes := []byte(content)
	fp, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)

	if err != nil {
		return err
	}

	defer fp.Close()
	fp.Truncate(0)
	fp.Seek(0, 0)
	_, err = fp.Write(bytes)

	return err
}

// WriteErrorLog ?
func WriteErrorLog(content string) {
	filepath := "./log/error.log"
	FilePutContents(filepath, content+"\r\n")
}
