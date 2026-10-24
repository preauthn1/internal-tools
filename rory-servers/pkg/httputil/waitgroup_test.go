package httputil

import (
	"fmt"
	"net/http"
	"path"
	"testing"

	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/model"
)

type LoadBlockIPListBody struct {
	model.BodyBase
	Data []model.DenyIPAddressSettingViewModel `json:"data"`
}

type PlayerSlotBonusRecordReadBody struct {
	model.BodyBase
	Data struct {
		SlotBonusRecord []model.PlayerSlotBonusRecordItemViewModel `json:"slotBonusRecord"`
	}
}

func TestWaitGroupRequest_Output(t *testing.T) {
	header := &http.Header{}
	header.Add("X-API-Certificate", "851b-222e5f24024b")
	header.Add("Access_Token", "")
	apiPath := fmt.Sprintf("http://%s", path.Join("example.invalid", "/system/setting/deny-ip"))
	var groupRequest = new(WaitGroupRequest)
	body := &LoadBlockIPListBody{}
	groupRequest.AddTask(GET, apiPath, nil, header, body)
	body1 := &LoadBlockIPListBody{}
	groupRequest.AddTask(GET, apiPath, nil, header, body1)
	body2 := &LoadBlockIPListBody{}
	groupRequest.AddTask(GET, apiPath, nil, header, body2)
	//groupRequest.AddTask(GET, apiPath, nil, *header)
	//groupRequest.AddTask(GET, apiPath, nil, *header)

	groupRequest.Send()
	groupRequest.Output() // 印資訊
}
