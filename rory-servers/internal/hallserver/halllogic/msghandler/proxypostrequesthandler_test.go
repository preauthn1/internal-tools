package msghandler

import (
	"testing"
)

func TestProxyPost(t *testing.T) {
	htmldata := "pay_md5sign=B7D217277EADD47D38FEE0E9070131C1&" +
		"pay_orderid=1379362&" +
			"pay_memberid=10011&" +
				"pay_notifyurl=http://0.0.0.0/Pay/ThPay/JuXin.aspx&" +
					"pay_amount=100.00&" +
						"pay_applydate=2019-05-13 14:30:58&" +
							"pay_productname=order1379362&" +
								"pay_bankcode=904&" +
									"pay_callbackurl=http://0.0.0.0/callback/M190511001&"
	pxy := new(ProxyPOSTRequestHandler)
	urlstring := "http://example.invalid/Pay_Index.html"

	pxy.proxyPostRequest(urlstring, htmldata)
}
