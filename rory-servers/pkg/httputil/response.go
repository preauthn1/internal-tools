package httputil

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// TaskResponse ?
type TaskResponse struct {
	*http.Response
	Status string
	Body   string
}

// RequestString ?
func (t *TaskResponse) RequestString() string {
	body, _ := t.Request.GetBody()
	var p []byte
	if body != nil {
		p, _ = ioutil.ReadAll(body)
	}
	s := fmt.Sprintf("%s %s body:\"%s\" header:\"%v\"\n", t.Request.Method, t.Request.URL, string(p), t.Request.Header)
	return s
}

// String ?
func (t *TaskResponse) String() string {
	s := t.RequestString()
	s += t.RequestCurlString()
	s += fmt.Sprintf("%s %s\n", t.Status, t.Body)
	return s
}

// RequestCurlString ?
func (t *TaskResponse) RequestCurlString() string {
	var header = func() (s []string) {
		for key, v := range t.Request.Header {
			if len(v) > 0 {
				s = append(s, fmt.Sprintf("-H \"%s: %s\"", key, v[0]))
			}
		}
		return
	}
	var query = func() (s string) {
		s += fmt.Sprintf("-X %s %s", t.Request.Method, t.Request.URL)
		body, _ := t.Request.GetBody()
		var p []byte
		if body != nil {
			p, _ = ioutil.ReadAll(body)
		}
		if len(p) > 0 {
			var ss = strconv.QuoteToASCII(string(p))
			//var ss, _ = json.Marshal(string(p))
			//var ss = strings.Replace(string(foo), "\"", "\\\"", -1)
			s += fmt.Sprintf(" -d %s", ss)
		}
		return
	}

	s := fmt.Sprintf("curl %s %s\n",
		query(),
		strings.Join(header(), " "),
	)
	return s
}
