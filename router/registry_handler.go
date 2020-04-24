package router

import (
	"encoding/json"
	"net/http"
)

// HandleRegistryRequest accepts a list of asset ids and returns info retrieved from
// the asset registry about them
func (r *Router) HandleRegistryRequest(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Methods", "POST")

	body := parseRequestBody(req.Body)
	assets := body["assets"]
	if assets == nil {
		http.Error(res, "Malformed Request", http.StatusBadRequest)
		return
	}

	entries, err := r.Registry.GetEntries(assets.([]interface{}))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(res).Encode(entries)
	return
}
