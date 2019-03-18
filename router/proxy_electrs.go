package router

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// ProxyElectrs forwards every request to the /esplora endpoint to electrs HTTP server
func (r *Router) ProxyElectrs(res http.ResponseWriter, req *http.Request) {
	electrsURL := fmt.Sprintf("http://%s:%s", r.Config.Electrs.Host, r.Config.Electrs.Port)
	endpoint := strings.Split(req.URL.String(), "/esplora/")[1]

	parsedURL, _ := url.Parse(electrsURL)
	parsedEndpoint, _ := url.Parse(endpoint)

	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.URL = parsedEndpoint
	req.Host = parsedURL.Host
	req.RequestURI = endpoint
	req.URL.Host = parsedURL.Host
	req.URL.Scheme = parsedURL.Scheme

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	proxy.ServeHTTP(res, req)
}
