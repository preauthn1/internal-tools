package httputil

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

// TaskRequest ?
type TaskRequest struct {
	*http.Request
	bindClass      interface{}
	OnSuccess      func(resp *TaskResponse)
	CallerFileLine string
}

// NewTaskRequest ?
// Example:
//   NewTaskRequest(method, url, params, headers, callback)
func NewTaskRequest(method string, url string, params map[string]interface{}, header *http.Header, bindClass interface{}) (*TaskRequest, error) {
	var callerFileLine string
	for i := 0; i < 10; i++ {
		_, file, line, _ := runtime.Caller(i)
		if !strings.Contains(file, "httputil") {
			callerFileLine = fmt.Sprintf("%v:%v", filepath.Base(file), line)
			break
		}
	}
	var body *bytes.Buffer
	if params == nil || len(params) == 0 {
		body = bytes.NewBuffer([]byte(""))
	} else {
		body = bytes.NewBuffer([]byte(util.ToJSON(params)))
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header = *header
	return &TaskRequest{req, bindClass, nil, callerFileLine}, nil
}

// Send ?
func (request *TaskRequest) Send() (taskResponse *TaskResponse, err error) {
	start := time.Now()

	defer func() {
		elapsed := time.Since(start)
		if err != nil {
			str := request.CallerFileLine + " " + request.String() + " " + err.Error() + "\n"
			mlog.Error("%s", err.Error())
			history.AddHistory(str, elapsed)
			return
		}
		str := request.CallerFileLine + " " + taskResponse.String()
		example.invalid(str)
		history.AddHistory(str, elapsed)
		request.onSuccess(taskResponse)
	}()
	var resp *http.Response
	var client = &http.Client{}
	var body []byte
	resp, err = client.Do(request.Request)

	if err != nil {
		taskResponse = &TaskResponse{
			Response: resp,
			Status:   http.StatusText(http.StatusBadRequest),
		}
		return
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		taskResponse = &TaskResponse{
			Response: resp,
			Status:   http.StatusText(http.StatusBadRequest),
		}
		return
	}
	taskResponse = &TaskResponse{
		Response: resp,
		Status:   resp.Status,
		Body:     string(body),
	}
	return
}

func (request *TaskRequest) onSuccess(resp *TaskResponse) {
	json.Unmarshal([]byte(resp.Body), request.bindClass)
	if request.OnSuccess != nil {
		request.OnSuccess(resp)
	}
}

func (request *TaskRequest) String() string {
	body, _ := request.Request.GetBody()
	var p []byte
	if body != nil {
		p, _ = ioutil.ReadAll(body)
	}
	s := fmt.Sprintf("%s %s header:\"%v\" body:\"%s\"\n", request.Request.Method, request.Request.URL, request.Request.Header, string(p))
	return s
}
