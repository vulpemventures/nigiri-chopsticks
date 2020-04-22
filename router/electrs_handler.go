package router

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/vulpemventures/nigiri-chopsticks/helpers"
)

type transport struct {
	r *helpers.Registry
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// check when request path is /asset/{asset_id}
	if strings.HasPrefix(req.URL.Path, "/asset/") {
		if s := strings.Split(req.URL.Path, "/"); len(s) == 3 {
			// parse response body
			payload, _ := ioutil.ReadAll(resp.Body)
			body := map[string]interface{}{}
			json.Unmarshal(payload, &body)

			// get registry entry for asset
			asset := body["asset_id"].(string)
			entry, _ := t.r.GetEntry(asset)

			// if entry exist add extra info to response
			if len(entry) > 0 {
				body["name"] = entry["name"]
				body["ticker"] = entry["ticker"]
				payload, _ = json.Marshal(body)
			}

			newBody := ioutil.NopCloser(bytes.NewReader(payload))
			resp.Body = newBody
			resp.ContentLength = int64(len(payload))
			resp.Header.Set("Content-Length", strconv.Itoa(len(payload)))
		}
	}

	return resp, nil
}

// HandleElectrsRequest forwards every request to the electrs HTTP server
func (r *Router) HandleElectrsRequest(res http.ResponseWriter, req *http.Request) {
	electrsURL := r.Config.ElectrsURL()
	parsedURL, _ := url.Parse(electrsURL)

	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = parsedURL.Host
	req.URL.Host = parsedURL.Host
	req.URL.Scheme = parsedURL.Scheme

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	proxy.Transport = &transport{r.Registry}
	proxy.ServeHTTP(res, req)
}
