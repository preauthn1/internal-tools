package service

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/module"
)

// Base ?
type Base struct {
	schema         string
	backendAPI     string
	apiCertificate string
}

// NewBase ?
func NewBase(ctx *module.MouduleCtx) *Base {
	base := &Base{}
	schema := ctx.GetConfigData("api", "conn", "schema")
	if !strings.Contains("https", schema) {
		schema = "http"
	}
	base.schema = schema
	base.backendAPI = ctx.GetConfigData("api", "conn", "backendapi")
	base.apiCertificate = ctx.GetConfigData("api", "conn", "api_cer")

	return base
}

func (b *Base) buildBackendAPI(apiName string, query map[string]interface{}) string {
	var apiPath = fmt.Sprintf("%s://%s", b.schema, path.Join(b.backendAPI, apiName))
	u, _ := url.Parse(apiPath)
	q := u.Query()
	for k, v := range query {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (b *Base) buildHeader(cer string) *http.Header {
	header := &http.Header{}
	header.Add("X-API-Certificate", b.apiCertificate)
	if cer != "" {
		header.Add("Access_Token", cer)
	}

	token := ParseAccessToken(cer)
	if token != nil {
		header.Add("locale", strings.ToLower(token.Lang))
	} else {
		header.Add("locale", "en-us")
	}

	return header
}
