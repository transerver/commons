package utils

import (
	"net/http"
	"net/netip"
	"strings"
)

func FetchIp(r *http.Request) string {
	ip := r.Header.Get("Cdn-Src-Ip")
	if len(ip) != 0 {
		return ip
	}

	ip = r.Header.Get("X-Connecting-IP")
	if len(ip) != 0 {
		return ip
	}

	ip = r.Header.Get("X-FORWARDED-FOR")
	if len(ip) == 0 || strings.EqualFold(ip, "unknown") {
		ip = r.RemoteAddr
	} else {
		ip = strings.SplitAfterN(ip, ",", 2)[0]
	}

	addrPort, err := netip.ParseAddrPort(ip)
	if err != nil || !addrPort.IsValid() {
		return ""
	}
	return addrPort.Addr().String()
}
