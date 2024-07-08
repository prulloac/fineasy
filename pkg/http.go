package pkg

import (
	"net/http"
)

type RequestMeta struct {
	Ip    string
	Agent string
}

func GetClientIp(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func GetRequestMeta(r *http.Request) RequestMeta {
	return RequestMeta{
		Ip:    GetClientIp(r),
		Agent: r.UserAgent(),
	}
}
