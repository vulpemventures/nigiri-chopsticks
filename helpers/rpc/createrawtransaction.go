package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Createrawtransaction(url string, utxo string, address string, changeAddress string) (int, string, error) {
	body := fmt.Sprintf(`{
		"jsonrpc": "1.0",
		"id": "rpc",
		"method": "createrawtransaction",
		"params": [
			[%s],
			[
				{"%s": 1},
				{"%s": 48.99}
			]
		]
	}`, utxo, address, changeAddress)

	status, resp, err := post(url, body, nil)
	if err != nil {
		return status, "", err
	}
	if status != http.StatusOK {
		err := parseResponseError(resp)
		return status, "", fmt.Errorf("an error occured while creating raw transction with status %d: %s", status, err)
	}

	out := map[string]interface{}{}
	json.Unmarshal([]byte(resp), &out)

	return status, out["result"].(string), nil
}
