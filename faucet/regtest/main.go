package regtestfaucet

import (
	"github.com/vulpemventures/nigiri-chopsticks/faucet"
	"github.com/vulpemventures/nigiri-chopsticks/helpers/rpc"
)

type regtestfaucet struct {
	URL string
}

// NewFaucet initialize a regtest faucet and returns it as interface
func NewFaucet(url string) faucet.Faucet {
	return &regtestfaucet{url}
}

func (f *regtestfaucet) NewTransaction(address string) (int, string, error) {
	status, utxo, err := rpc.Listunspent(f.URL)
	if err != nil {
		return status, "", err
	}

	status, changeAddress, err := rpc.Getrawchangeaddress(f.URL)
	if err != nil {
		return status, "", err
	}

	status, tx, err := rpc.Createrawtransaction(f.URL, utxo, address, changeAddress)
	if err != nil {
		return status, "", err
	}

	status, signedTx, err := rpc.Signrawtransaction(f.URL, tx)
	if err != nil {
		return status, "", err
	}

	return status, signedTx, nil
}

func (f *regtestfaucet) Fund() (int, []string, error) {
	status, balance, err := rpc.Getbalance(f.URL)
	if err != nil {
		return status, nil, err
	}

	if balance <= 0 {
		return rpc.Mine(f.URL, 200)
	}

	return 200, nil, nil
}
