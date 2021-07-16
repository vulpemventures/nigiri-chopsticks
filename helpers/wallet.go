package helpers

import (
	"fmt"
	"net/http"
)

func CreateWalletIfNotExists(client *RpcClient) error {

	status, resp, err := HandleRPCRequest(client, "listwallets", []interface{}{})
	if err != nil {
		return fmt.Errorf("could not list wallets: %w", err)
	}

	numOfWallets, ok := resp.([]interface{})
	if !ok {
		return fmt.Errorf("could not list wallets: %w", err)
	}

	if status == http.StatusOK && len(numOfWallets) == 0 {
		_, _, err = HandleRPCRequest(client, "createwallet", []interface{}{""})
		if err != nil {
			return fmt.Errorf("could not create wallet: %w", err)
		}
	}

	return nil
}
