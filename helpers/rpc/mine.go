package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Mine(url string, blocks int) (int, []string, error) {
	body := fmt.Sprintf(`{
		"jsonrpc": "1.0",
		"id": "rpc",
		"method": "generate",
		"params": [%d]
	}`, blocks)

	status, resp, err := post(url, body, nil)
	if err != nil {
		return status, nil, err
	}
	if status != http.StatusOK {
		err := parseResponseError(resp)
		return status, nil, fmt.Errorf("an error occured while mining blocks with status %d: %s", status, err)
	}

	out := map[string]interface{}{}
	json.Unmarshal([]byte(resp), &out)

	blockHashes := []string{}
	for _, bh := range out["result"].([]interface{}) {
		blockHashes = append(blockHashes, bh.(string))
	}

	return status, blockHashes, nil
}
