package tfnet

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ReportChannel 转发状态上报通道
type ReportChannel struct {
	*TFChannelBase
	submitURL string
}

// NewReportChannel 创建转发状态上报通道
func NewReportChannel(submitURL string, reportInterval int) *ReportChannel {
	rc := new(ReportChannel)
	rc.submitURL = submitURL

	//定时提交数据到URL地址
	go func() {
		defer util.PanicErrStack()

		for {
			time.Sleep(time.Duration(reportInterval) * time.Second)

			urn := fmt.Sprintf(submitURL, currentConn+GetCurrentWSConnCount())

			resp, err := http.Get(urn)

			if err != nil {
				mlog.Error("向服务器%s上报转发状态数据错误.%v", urn, err)
			} else {
				resdata, _ := ioutil.ReadAll(resp.Body)
				example.invalid("向服务器%s上报转发状态数据成功.获得返回数据%s", urn, resdata)
			}

			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
		}
	}()

	return rc
}

// OnConnected OnConnected
func (channel *ReportChannel) OnConnected(conninfo *ConnInfo) {}

// OnClientClose OnClientClose
func (channel *ReportChannel) OnClientClose(conninfo *ConnInfo) {}
