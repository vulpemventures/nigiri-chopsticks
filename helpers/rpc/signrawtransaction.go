package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Signrawtransaction(url string, tx string) (int, string, error) {
	body := fmt.Sprintf(`{
		"jsonrpc": "1.0",
		"id": "rpc",
		"method": "signrawtransactionwithwallet",
		"params": ["%s"]
	}`, tx)

	status, resp, err := post(url, body, nil)
	if err != nil {
		return status, "", err
	}
	if status != http.StatusOK {
		err := parseResponseError(resp)
		return status, "", fmt.Errorf("an error occured while signing transaction with status %d: %s", status, err)
	}

	out := map[string]map[string]interface{}{}
	json.Unmarshal([]byte(resp), &out)

	return status, out["result"]["hex"].(string), nil
}
