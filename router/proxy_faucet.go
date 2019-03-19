package router

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ProxyFaucet forwards every request to the /faucet endpoint to faucet HTTP server
func (r *Router) ProxyFaucet(res http.ResponseWriter, req *http.Request) {
	faucetURL := fmt.Sprintf("http://%s:%s", r.Config.Faucet.Host, r.Config.Faucet.Port)
	parsedURL, _ := url.Parse(faucetURL)

	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = parsedURL.Host
	req.URL.Host = parsedURL.Host
	req.URL.Scheme = parsedURL.Scheme

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	proxy.ServeHTTP(res, req)
}
