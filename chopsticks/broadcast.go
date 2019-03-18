package router

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Broadcast pushes an incoming transaction to the network and mines a block to get it confirmed
func (r *Router) Broadcast(res http.ResponseWriter, req *http.Request) {
	reqBody := parseRequestBody(req.Body)
	tx := reqBody["tx"]

	url := fmt.Sprintf("http://%s:%s@%s:%s", r.Config.Bitcoin.RPCUser, r.Config.Bitcoin.RPCPassword, r.Config.Bitcoin.Host, r.Config.Bitcoin.Port)
	headers := copyHeaders(req.Header)
	body := fmt.Sprintf(`{"jsonrpc": "1.0", "id": "2", "method": "sendrawtransaction", "params": ["%s"]}`, tx)

	status, resp, err := post(url, body, headers)
	if err != nil {
		http.Error(res, err.Error(), status)
		return
	}
	if status != http.StatusOK {
		http.Error(res, resp, status)
	}

	out := map[string]interface{}{}
	json.Unmarshal([]byte(resp), &out)

	response := map[string]string{
		"txId": out["result"].(string),
	}

	err = r.MineBlock()
	if err != nil {
		http.Error(res, err.Error(), status)
		return
	}

	json.NewEncoder(res).Encode(response)
}
