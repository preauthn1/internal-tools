package rpchandler

import (
	"bytes"
	"fmt"
	"strconv"

	"example.invalid/ServerTeam/Servers/internal/hallserver/halllogic/hallutil"
)

// HallVerCmd ?
type HallVerCmd struct {
	hallVerProvider hallutil.HallVerProvider
}

// NewHallVerCmd 創建機器人命令
func NewHallVerCmd(verHandler hallutil.HallVerProvider) *HallVerCmd {
	cvcmd := new(HallVerCmd)
	cvcmd.hallVerProvider = verHandler

	return cvcmd
}

// Name ?
func (hv *HallVerCmd) Name() string {
	return "hallver"
}

// Help /
func (hv *HallVerCmd) Help() string {
	return "客戶端版本號，eg:hallver <areaid>"
}

// Exec ?
func (hv *HallVerCmd) Exec(args []string) string {
	if len(args) < 1 {
		return "error,seeref usage"
	}
	areaID, err := strconv.Atoi(args[0])

	if err == nil {
		var buf bytes.Buffer
		ver := hv.hallVerProvider.GetVer(nil, nil, int64(areaID))
		buf.WriteString(fmt.Sprintf("===大廳版本資訊===\nIOS Version：%s\nAndroid Version：%s\n", ver.GetIOSVersion(), ver.GetAndroidVersion()))
		buf.WriteString(fmt.Sprintf("===大廳版本資訊===\nIOS Version2：%s\nAndroid Version2：%s\n", ver.GetIOSVersion2(), ver.GetAndroidVersion2()))
		buf.WriteString(fmt.Sprintf("===配置資訊===\nAreaConfig：%s\n", ver.GetAreaConfig()))
		buf.WriteString(fmt.Sprintf("GiftMoney：%d\n", ver.GetGiftMoney()))
		buf.WriteString(fmt.Sprintf("DefaultMoney：%d\n", ver.GetDefaultMoney()))
		buf.WriteString(fmt.Sprintf("BillDomain：%s\n", ver.GetBillDomain()))
		buf.WriteString(fmt.Sprintf("ClubDomain：%s\n", ver.GetClubDomain()))
		buf.WriteString(fmt.Sprintf("CustomerDomain：%s\n", ver.GetCustomerDomain()))
		buf.WriteString(fmt.Sprintf("UpdateDomain：%s\n", ver.GetUpdateDomain()))
		buf.WriteString(fmt.Sprintf("UpdateDomain2：%s\n", ver.GetUpdateDomain2()))
		return buf.String()
	}

	return "error,seeref usage"
}
