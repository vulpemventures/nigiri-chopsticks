package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Getbalance(url string) (int, float64, error) {
	body := `{
		"jsonrpc": "1.0",
		"id": "rpc",
		"method": "getbalance",
		"params": []
	}`

	status, resp, err := post(url, body, nil)
	if err != nil {
		return status, 0, err
	}
	if status != http.StatusOK {
		err := parseResponseError(resp)
		return status, 0, fmt.Errorf("an error occured while getting balance with status %d: %s", status, err)
	}

	out := map[string]interface{}{}
	json.Unmarshal([]byte(resp), &out)

	return status, out["result"].(float64), nil
}
