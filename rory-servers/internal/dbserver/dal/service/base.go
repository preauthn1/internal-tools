package service

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"example.invalid/ServerTeam/Servers/pkg/module"
)

// Base 服務基類
type Base struct {
	gameAPI        string
	apiCertificate string
}

// NewBase 創建基類
func NewBase(ctx *module.MouduleCtx) *Base {
	base := &Base{}
	base.gameAPI = ctx.GetConfigData("db", "conn", "gameapi")
	base.apiCertificate = ctx.GetConfigData("db", "conn", "api_cer")
	return base
}

// buildGameAPI 產生遊戲接口 URL
func (b *Base) buildGameAPI(apiName string, query map[string]interface{}) string {
	var apiPath = fmt.Sprintf("http://%s", path.Join(b.gameAPI, apiName))
	u, _ := url.Parse(apiPath)
	q := u.Query()
	for k, v := range query {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// buildHeader 產生 http 頭
func (b *Base) buildHeader(cer string) *http.Header {
	header := &http.Header{}
	header.Add("X-API-Certificate", b.apiCertificate)
	if cer != "" {
		header.Add("Access_Token", cer)
	}
	return header
}
