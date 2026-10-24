package util

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ?
const (
	IPISV4    = "ipv4"
	IPISV6    = "ipv6"
	IPUNKNOWN = "unknown"
)

// GetIPAndPort 解析地址中的IP和端口
func GetIPAndPort(addr net.Addr) (string, int) {
	return ParseIPAndPort(addr.String())
}

// ParseIPAndPort 從字符串中解析IP和端口
func ParseIPAndPort(addr string) (string, int) {
	if addr == "" {
		return "unknown", -1
	}

	data := strings.Split(addr, ":")
	ip := data[0]
	sport := data[1]

	port, _ := strconv.Atoi(sport)

	return ip, port
}

// ParseIPv4orIPv6 從字符串解析是否為 IPv4 或 IPv6, 正確 err 回傳 nil
func ParseIPv4orIPv6(addr string) (ip net.IP, ipv4orv6 string, err error) {
	var validIP = net.ParseIP(addr)
	if validIP != nil {
		// check IPv4, To4只能對IPv4地址進行轉換
		if validIP.To4() != nil {
			return validIP, IPISV4, nil
		}

		//check IPv6, To16轉換能對IPv4和IPv6
		if validIP.To16() != nil {
			return validIP, IPISV6, nil
		}

	}
	return nil, IPUNKNOWN, fmt.Errorf("not valid IP address. [%s]", addr)
}
