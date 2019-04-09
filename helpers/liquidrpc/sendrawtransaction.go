package liquidrpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Sendrawtransaction(url string, address string) (int, string, error) {
	body := fmt.Sprintf(`{
		"jsonrpc": "1.0",
		"id": "liquidrpc",
		"method": "sendrawtransaction",
		"params": ["%s"]
	}`, address)

	status, resp, err := post(url, body, nil)
	if err != nil {
		return status, "", err
	}
	if status != http.StatusOK {
		err := parseResponseError(resp)
		return status, "", fmt.Errorf("an error occured while broadcasting transaction with status %d: %s", status, err)
	}

	out := map[string]interface{}{}
	json.Unmarshal([]byte(resp), &out)

	return status, out["result"].(string), nil
}
