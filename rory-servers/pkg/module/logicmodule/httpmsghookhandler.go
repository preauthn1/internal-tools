package logicmodule

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// HTTPMsgHookHandler http消息處理器接口定義
type HTTPMsgHookHandler interface {
	OnHTTPMsg(semaphore util.Semaphore, w http.ResponseWriter, req *http.Request) bool
}

// HTTPMsgHookHandlerBase http消息處理器基類定義
type HTTPMsgHookHandlerBase struct{}

// NewHTTPMsgHookHandlerBase 創建Http消息處理器基類對象
func NewHTTPMsgHookHandlerBase() *HTTPMsgHookHandlerBase {
	hb := new(HTTPMsgHookHandlerBase)
	return hb
}

// OnHTTPMsg 收到HTTP訊息
func (hb *HTTPMsgHookHandlerBase) OnHTTPMsg(semaphore util.Semaphore, w http.ResponseWriter, req *http.Request) bool {
	ret := hb.HTTPHandler(w, req)
	if ret != nil {
		data, _ := proto.Marshal(ret)

		if data != nil {
			encodeStr := base64.StdEncoding.EncodeToString(data)
			io.WriteString(w, encodeStr)
		}
	}

	semaphore.Release()

	return false
}

// GetParamValue 獲取參數值
func (hb *HTTPMsgHookHandlerBase) GetParamValue(req *http.Request, pName string) string {
	return req.PostFormValue(pName)
}

// HTTPHandler Handler
func (hb *HTTPMsgHookHandlerBase) HTTPHandler(w http.ResponseWriter, req *http.Request) proto.Message {
	return nil
}

// HTTPMsgHookHandlerGroup http消息處理器組
type HTTPMsgHookHandlerGroup struct {
	handlers []HTTPMsgHookHandler
}

// NewHTTPMsgHookHandlerGroup 創建Http消息處理器組對象
func NewHTTPMsgHookHandlerGroup() *HTTPMsgHookHandlerGroup {
	hlg := new(HTTPMsgHookHandlerGroup)
	hlg.handlers = make([]HTTPMsgHookHandler, 0)

	return hlg
}

// AddHandler 添加處理器
func (hlg *HTTPMsgHookHandlerGroup) AddHandler(h HTTPMsgHookHandler) {
	hlg.handlers = append(hlg.handlers, h)
}

// OnHTTPMsg 收到HTTP訊息
func (hlg *HTTPMsgHookHandlerGroup) OnHTTPMsg(semaphore util.Semaphore, w http.ResponseWriter, req *http.Request) bool {
	defer func() {
		semaphore.Release()
	}()

	for _, h := range hlg.handlers {
		ret := h.OnHTTPMsg(semaphore, w, req)

		if !ret {
			return true
		}
	}

	return true
}

// ValidateHTTPMsgHookHandler http消息校驗處理器
type ValidateHTTPMsgHookHandler struct {
	code string
}

// NewValidateHTTPMsgHookHandler 創建http消息校驗處理器對象
func NewValidateHTTPMsgHookHandler(code string) *ValidateHTTPMsgHookHandler {
	vh := new(ValidateHTTPMsgHookHandler)
	vh.code = code

	return vh
}

// OnHTTPMsg 收到HTTP訊息
func (vh *ValidateHTTPMsgHookHandler) OnHTTPMsg(semaphore util.Semaphore, w http.ResponseWriter, req *http.Request) bool {

	// 獲取校驗碼
	pvcode := req.PostFormValue("vcode")
	if pvcode == "" {
		io.WriteString(w, "vcode error.")
		return false
	}

	// 對參數排序
	psarr := make([]string, 0)
	for p := range req.PostForm {
		if strings.ToLower(p) == "vcode" {
			continue
		}

		psarr = append(psarr, p)
	}
	sort.Sort(sort.StringSlice(psarr))

	// 構造參數加密字符串
	strData := ""
	for _, ps := range psarr {
		strData += ps
		strData += req.PostFormValue(ps)
	}
	strData += vh.code

	// 生成校驗碼
	cryptW := md5.New()
	io.WriteString(cryptW, strData)
	vcode := fmt.Sprintf("%x", cryptW.Sum(nil))

	// 檢查校驗碼
	if strings.ToUpper(vcode) == strings.ToUpper(pvcode) {
		return true
	}

	mlog.Debug("http vlidate,orign:%s,md5:%s ", strData, vcode)

	io.WriteString(w, "vcode error.")

	return false
}
