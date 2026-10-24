package msghandler

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"example.invalid/ServerTeam/Servers/pkg/module"
)

// BuildHeader ?
func BuildHeader(apiCer, token string) *http.Header {
	header := &http.Header{}
	header.Add("X-API-Certificate", apiCer)
	header.Add("Access_Token", token)
	return header
}

// BuildAPI ?
func BuildAPI(ctx *module.MouduleCtx, apiName string, query map[string]interface{}) string {
	var apiPath = fmt.Sprintf("http://%s", path.Join(ctx.BackendAPI, apiName))
	u, _ := url.Parse(apiPath)
	q := u.Query()
	for k, v := range query {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	u.RawQuery = q.Encode()
	return u.String()
}
