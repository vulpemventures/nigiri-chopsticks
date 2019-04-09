package liquidrpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Sendtoaddress(url string, address string, amount float64) (int, string, error) {
	body := fmt.Sprintf(`{
		"jsonrpc": "1.0",
		"id": "liquidrpc",
		"method": "sendtoaddress",
		"params": ["%s", %f]
	}`, address, amount)

	status, resp, err := post(url, body, nil)
	if err != nil {
		return status, "", err
	}
	if status != http.StatusOK {
		err := parseResponseError(resp)
		return status, "", fmt.Errorf("an error occured while sending funds to address with status %d: %s", status, err)
	}

	out := map[string]interface{}{}
	json.Unmarshal([]byte(resp), &out)

	return status, out["result"].(string), nil
}
