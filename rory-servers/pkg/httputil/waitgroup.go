package httputil

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/util"
)

/**
Example:
    header := &http.Header{}
	header.Add("X-API-Certificate", "851b-222e5f24024b")
	header.Add("Access_Token", "")
	apiPath := fmt.Sprintf("http://%s", path.Join("example.invalid", "/system/setting/deny-ip"))
*/

// WaitGroupRequest ?
type WaitGroupRequest struct {
	requests  []*TaskRequest
	responses []*TaskResponse
	costtime  time.Duration
}

// AddTask ?
func (r *WaitGroupRequest) AddTask(method string, url string, params map[string]interface{}, header *http.Header, bindClass interface{}) *TaskRequest {
	if ok, _ := util.InArray(method, []string{GET, POST, DELETE, PUT}); !ok {
		panic("method is not GET, POST, DELETE, PUT")
	}
	request, err := NewTaskRequest(method, url, params, header, bindClass)
	if err != nil {
		panic(err)
	}
	r.requests = append(r.requests, request)
	return request
}

// Send ?
func (r *WaitGroupRequest) Send() []*TaskResponse {
	var timebegin = time.Now()

	var wg sync.WaitGroup
	// var retCh = make(chan *TaskResponse, len(r.requests))

	// for _, request := range r.requests {
	// 	wg.Add(1)
	// 	go func(req *TaskRequest) {
	// 		defer wg.Done()
	// 		response, _ := req.Send()
	// 		retCh <- response
	// 	}(request)
	// 	time.Sleep(1)
	// }

	// wg.Wait()
	// close(retCh)

	// for val := range retCh {
	// 	r.responses = append(r.responses, val)
	// }

	var retChans []chan *TaskResponse
	retChans = make([]chan *TaskResponse, len(r.requests))
	for i := range retChans {
		retChans[i] = make(chan *TaskResponse, 1)
	}

	for i, request := range r.requests {
		wg.Add(1)
		go func(index int, req *TaskRequest) {
			defer util.PanicErrStack()
			defer wg.Done()
			response, _ := req.Send()
			retChans[index] <- response
		}(i, request)
		// time.Sleep(1)
	}

	wg.Wait()
	for _, retCh := range retChans {
		select {
		case val := <-retCh:
			close(retCh)
			r.responses = append(r.responses, val)
		case <-time.After(RequestTimeOut):
			panic(fmt.Errorf("request timeout over %v", RequestTimeOut))
		}
	}

	r.costtime = time.Since(timebegin)
	return r.responses
}

// Output ?
func (r *WaitGroupRequest) Output() {
	for _, v := range r.responses {
		fmt.Println(v.String())
	}
	fmt.Println("cost", r.costtime)
}

// IsAllDone ?
func (r *WaitGroupRequest) IsAllDone() bool {
	for _, v := range r.responses {
		if v.Response == nil || v.StatusCode < 200 || v.StatusCode >= 300 {
			return false
		}
	}
	return true
}
