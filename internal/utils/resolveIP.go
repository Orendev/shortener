package utils

import (
	"fmt"
	"net"
	"net/http"
)

// ResolveIP get the IP address from the header.
func ResolveIP(r *http.Request) (net.IP, error) {
	// смотрим заголовок запроса X-Real-IP
	ipStr := r.Header.Get("X-Real-IP")
	// парсим ip
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("failed parse ip from http header")
	}
	return ip, nil
}
