package h22p

import (
	"net"
	"net/http"
	"strings"
)

func GetRealIp(r *http.Request) string {
	s := r.Header.Get("X-Real-IP")
	if s != "" {
		return s
	}

	s = r.Header.Get("X-Forwarded-For")
	if s != "" {
		return strings.TrimSpace(strings.Split(s, ",")[0])
	}

	s, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return s
}
