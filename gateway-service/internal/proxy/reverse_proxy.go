package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/RadekKusiak71/splitwise/gateway/internal/auth"
)

func New(target string) http.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del("X-User-ID")

		if userID, ok := r.Context().Value(auth.UserIDKey).(string); ok {
			r.Header.Set("X-User-ID", userID)
		}

		r.Host = url.Host
		proxy.ServeHTTP(w, r)
	}
}
