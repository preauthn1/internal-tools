package monitorutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// Get get message
func Get(urlstr string) (string, error) {
	u, err := url.Parse(urlstr)
	if err != nil {
		panic(err)
	}

	URLAPI := fmt.Sprintf("http://" + u.Host + u.Path + "?" + u.RawQuery)

	resp, err := http.Get(URLAPI)
	if err != nil {
		mlog.Error("http Get Error : %v", err)
		return "", errors.New("http Get Error :" + err.Error())
	}
	if resp.StatusCode != 200 {
		mlog.Error("%v連線失敗 StatusCode:%v", urlstr, resp.StatusCode)
		return "", errors.New(urlstr + "連線失敗 StatusCode:" + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		mlog.Error("http Read Error : %v", err)
		return "", errors.New("http Read Error : " + err.Error())
	}

	return string(body), nil
}
