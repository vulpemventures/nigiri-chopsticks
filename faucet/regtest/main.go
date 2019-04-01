package regtestfaucet

import "github.com/vulpemventures/nigiri-chopsticks/faucet"

type regtestfaucet struct {
	URL string
}

// NewFaucet initialize a regtest faucet and returns it as interface
func NewFaucet(url string) faucet.Faucet {
	return &regtestfaucet{url}
}

func (f *regtestfaucet) NewTransaction(address string) (int, string, error) {
	status, utxo, err := listunspent(f.URL)
	if err != nil {
		return status, "", err
	}

	status, changeAddress, err := getrawchangeaddress(f.URL)
	if err != nil {
		return status, "", err
	}

	status, tx, err := createrawtransaction(f.URL, utxo, address, changeAddress)
	if err != nil {
		return status, "", err
	}

	status, signedTx, err := signrawtransaction(f.URL, tx)
	if err != nil {
		return status, "", err
	}

	return status, signedTx, nil
}
