package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Getrawchangeaddress(url string) (int, string, error) {
	body := `{
		"jsonrpc": "1.0",
		"id": "rpc",
		"method": "getrawchangeaddress",
		"params": []
	}`

	status, resp, err := post(url, body, nil)
	if err != nil {
		return status, "", err
	}
	if status != http.StatusOK {
		err := parseResponseError(resp)
		return status, "", fmt.Errorf("an error occured while getting change address with status %d: %s", status, err)
	}

	out := map[string]interface{}{}
	json.Unmarshal([]byte(resp), &out)

	return status, out["result"].(string), nil
}
