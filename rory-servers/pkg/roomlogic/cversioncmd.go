package roomlogic

import (
	"bytes"
	"fmt"

	"example.invalid/golang/protobuf/proto"
)

// CVersionCmdP Client 版本指令
var CVersionCmdP *CVersionCmd

// CVersionCmd Client 版本 指令
type CVersionCmd struct {
	verprovider GameVerProvider
}

// NewCVersionCmd Create Client 版本 指令
func NewCVersionCmd(verprovider GameVerProvider) *CVersionCmd {
	cvcmd := new(CVersionCmd)
	cvcmd.verprovider = verprovider
	return cvcmd
}

// Name 指令名稱
func (cv *CVersionCmd) Name() string {
	return "cversion"
}

// Help 指令說明
func (cv *CVersionCmd) Help() string {
	return "客戶端版本號，eg:cversion list | cversion ios 1.0.0 | cversion android 1.0.0"
}

// Exec 執行指令
func (cv *CVersionCmd) Exec(args []string) string {
	if len(args) < 1 {
		return "error,seeref usage"
	}

	cmdname := args[0]

	if cmdname == "list" {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("===客戶端版本資訊===\nIOS Version：%s\nAndroid Version：%s", cv.GetIOS(), cv.GetAndroid()))
		return buf.String()
	}

	if cmdname == "ios" {
		if len(args) > 1 && args[1] != "" {
			cv.verprovider.(*RPCGameVerProvider).gsreg.gameServerRegRet.IOSVersion = proto.String(args[1])
			cv.SetIOS(args[1])
			return "變更版號成功"
		}
	}

	if cmdname == "android" {
		if len(args) > 1 && args[1] != "" {
			cv.SetAndroid(args[1])
			return "變更版號成功"
		}
	}

	return "error,seeref usage"
}

// GetIOS 取得 IOS 版本
func (cv *CVersionCmd) GetIOS() string {
	return cv.verprovider.GetVer().GetIOSVersion()
}

// SetIOS 設定 IOS 版本
func (cv *CVersionCmd) SetIOS(v string) {
	cv.verprovider.(*RPCGameVerProvider).gsreg.gameServerRegRet.IOSVersion = proto.String(v)
}

// GetAndroid 取得 Android 版本
func (cv *CVersionCmd) GetAndroid() string {
	return cv.verprovider.GetVer().GetAndroidVersion()
}

// SetAndroid 設定 Android 版本
func (cv *CVersionCmd) SetAndroid(v string) {
	cv.verprovider.(*RPCGameVerProvider).gsreg.gameServerRegRet.AndroidVersion = proto.String(v)
}
