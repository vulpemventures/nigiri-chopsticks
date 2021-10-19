package helpers

import (
	"fmt"
)

func RescanBlockchain(client *RpcClient) error {
	_, _, err := HandleRPCRequest(client, "rescanblockchain", []interface{}{})
	if err != nil {
		return fmt.Errorf("could not rescan: %w", err)
	}

	return nil
}
