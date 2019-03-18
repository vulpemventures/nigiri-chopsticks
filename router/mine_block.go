package router

import (
	"fmt"
	"net/http"
)

// MineBlock calls the JSONRPC method to mine 1 block
func (r *Router) MineBlock() error {
	url := fmt.Sprintf("http://%s:%s@%s:%s", r.Config.Bitcoin.RPCUser, r.Config.Bitcoin.RPCPassword, r.Config.Bitcoin.Host, r.Config.Bitcoin.Port)
	body := `{"jsonrpc": "1.0", "id": "2", "method": "generate", "params": [1]}`
	headers := map[string]string{"Content-Type": "application/json"}

	status, resp, err := post(url, body, headers)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("an error occured while mining block: %s", resp)
	}

	return nil
}
