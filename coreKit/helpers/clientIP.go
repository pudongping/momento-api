package helpers

import (
	"net"
	"net/http"
	"strings"
)

// GetClientIP 从请求中提取客户端 IP，优先使用 X-Forwarded-For，再使用 X-Real-Ip，最后使用 RemoteAddr。
func GetClientIP(r *http.Request) string {
	// X-Forwarded-For 可能包含多个 IP，取第一个（最原始的客户端 IP）
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			if ip := strings.TrimSpace(parts[0]); ip != "" {
				return ip
			}
		}
	}

	// X-Real-Ip
	if xr := strings.TrimSpace(r.Header.Get("X-Real-Ip")); xr != "" {
		return xr
	}

	// 回退到 RemoteAddr（格式可能是 "ip:port"）
	if host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}
