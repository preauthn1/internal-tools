package monitorutil

import (
	"encoding/json"
	"strings"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// BuildTfData build tf data
func BuildTfData(url string, siteID int32) string {
	var tfData *msgserver.TFDataReq
	if strings.Contains(url, "AppleReview") {
		tfData = &msgserver.TFDataReq{
			UserID: 1,
			AreaID: int(siteID),
			VIP:    0,
			GUI:    "xxxxxxxxxxxxxxxx",
			Time:   int(time.Now().Unix()),
			Ver:    3,
		}
	} else {
		tfData = &msgserver.TFDataReq{
			UserID:     1,
			Channel:    "1",
			GUI:        "xxxxxxxxxxxxxxxx",
			Simulation: 0,
			VIP:        0,
			Action:     1,
			Time:       int(time.Now().Unix())}
	}

	tfBytes, _ := json.Marshal(tfData)
	dataString, err := util.AesEncryptString(string(tfBytes), "V1M5WE7X8WE2TWEK", "DNW123EXGMTR83K7")
	if err != nil {

	}
	mlog.Debug("%v", dataString)
	return dataString
}
