package httputil

import (
	"fmt"
	"net/http"
	"path"
	"testing"

	"example.invalid/stretchr/testify/assert"
)

func TestNewTaskRequest(t *testing.T) {
	// systemSettingDenyIP(t)
	// gameJackpotLink(t)
}

func systemSettingDenyIP(t *testing.T) {
	header := &http.Header{}
	header.Add("X-API-Certificate", "851b-222e5f24024b")
	header.Add("Access_Token", "")
	apiPath := fmt.Sprintf("http://%s", path.Join("example.invalid", "/system/setting/deny-ip"))
	req, err := NewTaskRequest("GET", apiPath, nil, header, nil)
	assert.Nil(t, err)
	response, err := req.Send()
	t.Log(response.String())
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, 200, response.String())
	req, err = NewTaskRequest("POST", apiPath, nil, header, nil)
	assert.Nil(t, err)
	response, err = req.Send()
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, 405, response.String())
	t.Log(response.String())
}

func gameJackpotLink(t *testing.T) {
	header := &http.Header{}
	header.Add("X-API-Certificate", "851b-222e5f24024b")
	header.Add("Access_Token", "")
	apiPath := fmt.Sprintf("http://%s", path.Join("example.invalid", fmt.Sprintf("/game/%d/jackpot", 25)))
	var params = map[string]interface{}{
		"flagID":  1,
		"tableNo": 0,
		"levelID": 1,
		"isLink":  true,
	}
	req, err := NewTaskRequest("GET", apiPath, params, header, nil)
	assert.Nil(t, err)
	response, err := req.Send()
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, 405)
	t.Log(response.String())
	req, err = NewTaskRequest("POST", apiPath, params, header, nil)
	assert.Nil(t, err)
	response, err = req.Send()
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, 200)
	t.Log(response.String())
}
