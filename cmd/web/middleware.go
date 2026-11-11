package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type contextKey string

const contextUserKey contextKey = "user_ip"

func (a *application) ipFromContext(ctx context.Context) string {
	return ctx.Value(contextUserKey).(string)
}

func (a *application) addIpToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var ctx = context.Background()
		ip, err := getIp(request)
		if err != nil {
			ip, _, _ = net.SplitHostPort(request.RemoteAddr)
			if len(ip) == 0 {
				ip = "unknown"
			}
			ctx = context.WithValue(request.Context(), contextUserKey, ip)
		} else {
			ctx = context.WithValue(request.Context(), contextUserKey, ip)
		}
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func getIp(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		return "unknown", err
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "", fmt.Errorf("userip: %q is not IP:port", r.RemoteAddr)
	}

	forward := r.Header.Get("X-Forwarded-For")
	if len(forward) > 0 {
		ip = forward
	}

	if len(ip) == 0 {
		ip = "forward"
	}

	return ip, nil
}

func (a *application) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.Session.Exists(r.Context(), "user") {
			a.Session.Put(r.Context(), "error", "Log in first")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}
