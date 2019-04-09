package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Listunspent(url string) (int, string, error) {
	body := `{
		"jsonrpc": "1.0",
		"id": "rpc",
		"method": "listunspent",
		"params": []
	}`

	status, resp, err := post(url, body, nil)
	if err != nil {
		return status, "", err
	}
	if status != http.StatusOK {
		err := parseResponseError(resp)
		return status, "", fmt.Errorf("an error occured while listing unspents with status %d: %s", status, err)
	}

	out := map[string]interface{}{}
	err = json.Unmarshal([]byte(resp), &out)
	if err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("an error occured while unmarshaling unspents: %s", err)
	}
	utxo, _ := json.Marshal(out["result"].([]interface{})[0])

	return status, string(utxo), nil
}
