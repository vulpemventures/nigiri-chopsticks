package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// HandleElectrsRequest forwards every request to the electrs HTTP server
func (r *Router) HandleElectrsRequest(res http.ResponseWriter, req *http.Request) {
	electrsURL := r.Config.ElectrsURL()
	parsedURL, _ := url.Parse(electrsURL)

	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = parsedURL.Host
	req.URL.Host = parsedURL.Host
	req.URL.Scheme = parsedURL.Scheme

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	proxy.ServeHTTP(res, req)
}
