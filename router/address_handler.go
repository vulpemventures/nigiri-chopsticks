package router

import (
	"encoding/json"
	"net/http"

	"github.com/vulpemventures/nigiri-chopsticks/helpers"
)

// HandleAddressRequest calls the `getnewaddress` and returns the native segwit one
func (r *Router) HandleAddressRequest(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")

	status, resp, err := helpers.HandleRPCRequest(r.RPCClient, "getnewaddress", []interface{}{"", "bech32"})
	if err != nil {
		http.Error(res, err.Error(), status)
		return
	}

	json.NewEncoder(res).Encode(map[string]string{"address": resp.(string)})
	return
}
