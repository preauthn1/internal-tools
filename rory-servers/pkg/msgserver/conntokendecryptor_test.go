package msgserver

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

func TestConnTokenDecryptor(t *testing.T) {

	data := "CC233FAB00555F0DB6446745323457E30EFFA60FB6BEC1E68D016198941099AC727B2428E0C7EF43A018537C6FB9AC78B083E9199D6CF28E3B12BCB412BCB6C221C28E60CFF976B487CB28EE0F9F84412E39BD07B817E12472C3BFE9313137324476EA15D925A880E0D76F0296267A24E6F932BA098C14FB10251134348B2B98A53B0C858013C3F93EA8AE7E1DCC4DD85A96D35C599A59F7C7C4D61802CEC3C9335679EBFBB306F3796BF23DDD9655C741F8161913F628E82C353376CE02075D5878BAA4F50016781CE1AFBA4C754AA20013777B78A12445FCDBD2CA169003A69D60A4986E623CBC0DF2B02A8C39B75A"
	token, _ := ConnTokenDecryptor(data)
	mlog.System("value %s, time %d", token.Value, token.Time)
	tm := time.Unix(int64(token.Time), 0)
	fmt.Println(tm)
	tm = time.Unix(int64(token.Time), +int64(1*time.Minute))

	fmt.Println(time.Now().Before(tm))
}

func TestConnTokenEncryptor(t *testing.T) {
	tfData := &TFDataReq{
		UserID:     1,
		Channel:    "1",
		GUI:        "xxxxxxxxxxxxxxxx",
		Simulation: 0,
		VIP:        0,
		Action:     1,
		Time:       int(time.Now().Unix())}
	tfBytes, err := json.Marshal(tfData)
	if err != nil {
		t.Errorf("%v", err)
	}

	dataString, err := util.AesEncryptString(string(tfBytes), "V1M5WE7X8WE2TWEK", "DNW123EXGMTR83K7")
	mlog.System("%v", string(dataString))
	if err != nil {
		t.Errorf("%v", err)
	}
	data, err := util.AesDecryptString(dataString, "V1M5WE7X8WE2TWEK", "DNW123EXGMTR83K7")
	if err != nil {
		t.Errorf("%v", err)
	}
	res := TFDataReq{}
	json.Unmarshal([]byte(data), &res)
	if res.UserID != tfData.UserID {
		t.Errorf("userid error")
	}
	if res.Time != tfData.Time {
		t.Errorf("time error")
	}
}
